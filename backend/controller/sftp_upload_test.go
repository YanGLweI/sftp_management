package controller

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"sftpbackend/utils"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// buildUploadBody 按给定顺序构造 multipart 请求体（顺序决定 part 到达后端的先后）
func buildUploadBody(t *testing.T, fields map[string]string, order []string, fileName, fileContent string) (*bytes.Buffer, string) {
	t.Helper()
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	for _, name := range order {
		if name == "file" {
			part, err := w.CreateFormFile("file", fileName)
			if err != nil {
				t.Fatalf("CreateFormFile: %v", err)
			}
			if _, err := io.WriteString(part, fileContent); err != nil {
				t.Fatalf("write file part: %v", err)
			}
			continue
		}
		if err := w.WriteField(name, fields[name]); err != nil {
			t.Fatalf("WriteField %s: %v", name, err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatalf("writer close: %v", err)
	}
	return body, w.FormDataContentType()
}

func doUpload(t *testing.T, body *bytes.Buffer, contentType, token string) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/sftp/upload", UploadFile)

	req := httptest.NewRequest(http.MethodPost, "/sftp/upload", body)
	req.Header.Set("Content-Type", contentType)
	if token != "" {
		req.Header.Set("X-SFTP-Token", token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// TestUploadFileNoToken 缺少 Token 时返回 400
func TestUploadFileNoToken(t *testing.T) {
	body, ct := buildUploadBody(t,
		map[string]string{"path": "/tmp"},
		[]string{"path", "file"}, "a.txt", "abc")
	w := doUpload(t, body, ct, "")
	if w.Code != http.StatusBadRequest || !strings.Contains(w.Body.String(), "SFTP Token不能为空") {
		t.Fatalf("expected 400 token error, got %d: %s", w.Code, w.Body.String())
	}
}

// TestUploadFileStreamedAfterPath path 先于 file 到达时，应读取 path 并进入 SFTP 写入阶段。
// 注入的伪连接 SftpClient 为 nil，CreateUploadFile 返回"SFTP连接未初始化"，
// 该错误只在流式解析成功、进入写入阶段后才会出现，以此证明流式路径生效。
func TestUploadFileStreamedAfterPath(t *testing.T) {
	token := utils.SFTPConnManager.AddConn(&utils.SFTPConnection{})
	defer utils.SFTPConnManager.RemoveConn(token)

	body, ct := buildUploadBody(t,
		map[string]string{"path": "/tmp"},
		[]string{"path", "file"}, "a.txt", "hello-stream")
	w := doUpload(t, body, ct, token)
	if w.Code != http.StatusInternalServerError || !strings.Contains(w.Body.String(), "SFTP连接未初始化") {
		t.Fatalf("expected 500 with SFTP未初始化 (streaming reached write stage), got %d: %s", w.Code, w.Body.String())
	}
}

// TestUploadFileBeforePathRejected file 先于 path 到达时应返回 400，
// 防止在未知目标路径的情况下开始接收文件内容
func TestUploadFileBeforePathRejected(t *testing.T) {
	token := utils.SFTPConnManager.AddConn(&utils.SFTPConnection{})
	defer utils.SFTPConnManager.RemoveConn(token)

	body, ct := buildUploadBody(t,
		map[string]string{"path": "/tmp"},
		[]string{"file", "path"}, "a.txt", "abc")
	w := doUpload(t, body, ct, token)
	if w.Code != http.StatusBadRequest || !strings.Contains(w.Body.String(), "目标路径不能为空") {
		t.Fatalf("expected 400 path order error, got %d: %s", w.Code, w.Body.String())
	}
}

// TestUploadFileMissingFilePart 请求中没有 file part 时返回 400
func TestUploadFileMissingFilePart(t *testing.T) {
	token := utils.SFTPConnManager.AddConn(&utils.SFTPConnection{})
	defer utils.SFTPConnManager.RemoveConn(token)

	body, ct := buildUploadBody(t,
		map[string]string{"path": "/tmp"},
		[]string{"path"}, "", "")
	w := doUpload(t, body, ct, token)
	if w.Code != http.StatusBadRequest || !strings.Contains(w.Body.String(), "未接收到上传文件") {
		t.Fatalf("expected 400 missing file error, got %d: %s", w.Code, w.Body.String())
	}
}
