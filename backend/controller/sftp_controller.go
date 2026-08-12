package controller

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"sftpbackend/models"
	"sftpbackend/tools"
	"sftpbackend/utils"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 从请求中获取SFTP连接Token（优先从Header获取，也可从Query/Body获取）
func getSFTPToken(c *gin.Context) string {
	// 1. 优先从Header获取
	token := c.GetHeader("X-SFTP-Token")
	if token != "" {
		return token
	}

	// 2. 备选：从Query参数获取
	token = c.Query("sftp_token")
	if token != "" {
		return token
	}

	// 3. 备选：从Form/Body获取（根据前端传参方式调整）
	token = c.PostForm("sftp_token")
	return token
}

// ! LoginSFTP 处理SFTP登录（返回Token）
func LoginSFTP(c *gin.Context) {
	var conn *utils.SFTPConnection
	var err error

	// （密钥登录）获取上传的密钥文件
	file, fileErr := c.FormFile("file")
	if fileErr == nil && file != nil {
		srcFile, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "无法打开上传的密钥文件: " + err.Error(),
			})
			return
		}
		defer srcFile.Close()

		username := c.PostForm("username")
		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "用户名不能为空",
			})
			return
		}

		keyBytes, err := io.ReadAll(srcFile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "无法读取密钥文件内容: " + err.Error(),
			})
			return
		}

		// 创建连接实例
		conn, err = utils.NewSFTPConnectionByKey(username, keyBytes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "<密钥错误> " + err.Error(),
			})
			return
		}
	} else {
		// （密码登录）
		var sftpLogin models.SFTPLogin
		if err := c.ShouldBindJSON(&sftpLogin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请求数据格式错误",
			})
			return
		}

		// 解密密码
		decryptedPassword, decerr := tools.DecryptPassword(sftpLogin.Password)
		if decerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "<密码错误> " + decerr.Error(),
			})
			return
		}

		// 创建连接实例
		conn, err = utils.NewSFTPConnection(sftpLogin.Username, decryptedPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "<密码错误> " + err.Error(),
			})
			return
		}
	}

	// 将连接添加到管理器，获取唯一Token
	token := utils.SFTPConnManager.AddConn(conn)

	// 登录成功，返回Token给前端
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "SFTP登录成功",
		"data": gin.H{
			"sftp_token":  token, // 前端需要保存这个Token
			"expire_tips": "连接有效期8小时，超时需重新登录",
		},
	})
}

// ! LogoutSFTP 处理SFTP登出（删除Token和连接）
func LogoutSFTP(c *gin.Context) {
	token := getSFTPToken(c)
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "SFTP Token不能为空",
		})
		return
	}

	// 从管理器中删除并关闭连接
	if err := utils.SFTPConnManager.RemoveConn(token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "登出失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "SFTP登出成功",
	})
}

// ! 获取目录下的文件列表（通过Token获取连接）
func GetFiles(c *gin.Context) {
	// 1. 获取Token
	token := getSFTPToken(c)
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "SFTP Token不能为空",
		})
		return
	}

	// 2. 通过Token获取对应的连接实例
	conn, err := utils.SFTPConnManager.GetConn(token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    50014,
			"message": "SFTP连接失效: " + err.Error(),
		})
		return
	}

	// 3. 执行具体操作（使用当前连接实例）
	path := c.Query("path")
	if path == "" {
		path = "/"
	}

	files, err := conn.SftpClient.ReadDir(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "无法读取目录: " + err.Error(),
		})
		return
	}

	var dirs []models.FileInfo
	var fileList []models.FileInfo
	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			continue // 跳过隐藏项，不加入结果集
		}
		fileInfo := models.FileInfo{
			Name:     file.Name(),
			Path:     filepath.Join(path, file.Name()),
			IsDir:    file.IsDir(),
			Size:     file.Size(),
			Modified: file.ModTime().Format(time.RFC3339),
		}
		if file.IsDir() {
			dirs = append(dirs, fileInfo)

		} else {
			fileList = append(fileList, fileInfo)
		}
	}
	// 按照Modified排序
	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i].Modified > dirs[j].Modified
	})
	sort.Slice(fileList, func(i, j int) bool {
		return fileList[i].Modified > fileList[j].Modified
	})

	dirsCount := len(dirs)
	filesCount := len(fileList)
	var description string
	if dirsCount == 0 && filesCount == 0 {
		description = "空目录"
	} else if dirsCount == 0 && filesCount > 0 {
		description = fmt.Sprintf("%d 个文件", filesCount)
	} else if dirsCount > 0 && filesCount == 0 {
		description = fmt.Sprintf("%d 个目录", dirsCount)
	} else {
		description = fmt.Sprintf("%d 个文件 和 %d 个目录", filesCount, dirsCount)
	}

	sortedFiles := append(dirs, fileList...)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "文件列表获取成功",
		"data": gin.H{
			"path":        path,
			"description": description,
			"files":       sortedFiles,
		},
	})
}

// ! UploadFile 上传文件（适配Token，流式直写SFTP，避免整体缓冲导致进度虚高）
func UploadFile(c *gin.Context) {
	// 1. 获取Token（仅从 Header/Query 获取，避免解析请求体破坏流式读取）
	token := c.GetHeader("X-SFTP-Token")
	if token == "" {
		token = c.Query("sftp_token")
	}
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "SFTP Token不能为空",
		})
		return
	}

	// 2. 获取连接实例
	conn, err := utils.SFTPConnManager.GetConn(token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    50014,
			"message": "SFTP连接失效: " + err.Error(),
		})
		return
	}

	// 3. 流式读取 multipart 请求体（不能用 FormFile/PostForm，否则会先把整个文件缓冲到本地）
	mr, err := c.Request.MultipartReader()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "文件上传失败: " + err.Error(),
		})
		return
	}

	var targetPath string
	uploaded := false
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "文件上传失败: " + err.Error(),
			})
			return
		}

		switch part.FormName() {
		case "path":
			data, err := io.ReadAll(part)
			part.Close()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    400,
					"message": "读取目标路径失败: " + err.Error(),
				})
				return
			}
			targetPath = string(data)
		case "file":
			if targetPath == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    400,
					"message": "目标路径不能为空",
				})
				return
			}
			filename := part.FileName()
			if filename == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    400,
					"message": "上传文件名为空",
				})
				return
			}
			dstPath := filepath.Join(targetPath, filename)
			// 边接收边写入SFTP，SFTP写入速度通过TCP背压传导到浏览器，进度条反映真实进度
			if err := conn.CreateUploadFile(dstPath, part); err != nil {
				part.Close()
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "文件上传失败: " + err.Error(),
				})
				return
			}
			part.Close()
			uploaded = true
		default:
			// 忽略其他字段，读完丢弃以保持 multipart 流推进
			io.Copy(io.Discard, part)
			part.Close()
		}
	}

	if !uploaded {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "未接收到上传文件",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "文件上传成功",
	})
}

// ! CreateFolder 创建目录（适配Token）
func CreateFolder(c *gin.Context) {
	// 1. 获取Token
	token := getSFTPToken(c)
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "SFTP Token不能为空",
		})
		return
	}

	// 2. 获取连接实例
	conn, err := utils.SFTPConnManager.GetConn(token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    50014,
			"message": "SFTP连接失效: " + err.Error(),
		})
		return
	}

	// 3. 执行创建目录
	var req struct {
		Path string `json:"path"`
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求数据格式错误",
		})
		return
	}

	if req.Path == "" || req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "路径和名称不能为空",
		})
		return
	}

	if err := conn.CreateFolder(req.Path, req.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建目录失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "目录创建成功",
	})
}

// ! DownloadFile 下载文件（适配Token）
func DownloadFile(c *gin.Context) {
	// 1. 获取Token
	token := getSFTPToken(c)
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "SFTP Token不能为空",
		})
		return
	}

	// 2. 获取连接实例
	conn, err := utils.SFTPConnManager.GetConn(token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    50014,
			"message": "SFTP连接失效: " + err.Error(),
		})
		return
	}

	// 3. 执行下载操作
	filePath := c.Query("path")
	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "文件路径不能为空",
		})
		return
	}

	// 检查文件是否存在
	fileInfo, err := conn.SftpClient.Stat(filePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "文件不存在: " + err.Error(),
		})
		return
	}

	// 4. 检查是否为目录
	if fileInfo.IsDir() {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "路径是目录不是文件",
		})
		return
	}

	// 5. 打开文件并传输
	file, err := conn.SftpClient.Open(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "打开文件失败: " + err.Error(),
		})
		return
	}
	defer file.Close()

	// fileName := filepath.Base(filePath)
	fileName := fileInfo.Name() // 获取SFTP服务器上的文件名
	// 二进制流文件，让浏览器识别为下载文件
	c.Header("Content-Type", "application/octet-stream")
	// c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"; filename*=UTF-8''%s",
		escapeFileName(fileName),  // 兼容低版本浏览器的文件名转义
		url.QueryEscape(fileName), // 标准UTF-8转义，解决中文文件名乱码
	))
	// c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	c.Header("Content-Length", strconv.FormatInt(fileInfo.Size(), 10)) // 文件大小，让浏览器显示下载进度
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")   // 禁止缓存，避免大文件缓存问题
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.Status(http.StatusOK) // 设置HTTP状态码为200

	// if _, err := io.Copy(c.Writer, file); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"code":    500,
	// 		"message": "文件传输失败: " + err.Error(),
	// 	})
	// 	return
	// }
	// 流式传输（SFTP流 → HTTP响应流）
	// io.Copy会自动分块复制数据（默认缓冲区32KB），不会将整个文件加载到后端内存
	// 直接将SFTP读取的字节流式写入gin的ResponseWriter，前端浏览器实时接收
	_, err = io.Copy(c.Writer, file)
	if err != nil {
		// 流式传输过程中出错，记录日志即可（浏览器会自动中断下载）
		log.Printf("SFTP文件流式下载失败，路径：%s，错误：%v", filePath, err)
		return
	}
}

// 辅助函数：转义文件名，兼容IE/Edge等低版本浏览器（解决中文乱码）
func escapeFileName(name string) string {
	// 替换特殊字符，避免响应头解析错误
	replacer := strings.NewReplacer(`"`, `\"`, `\`, `\\`, `/`, `\/`)
	return replacer.Replace(name)
}

// ! DeletePath 删除文件/目录（适配Token）
func DeletePath(c *gin.Context) {
	// 1. 获取Token
	token := getSFTPToken(c)
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "SFTP Token不能为空",
		})
		return
	}

	// 2. 获取连接实例
	conn, err := utils.SFTPConnManager.GetConn(token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    50014,
			"message": "SFTP连接失效: " + err.Error(),
		})
		return
	}

	// 3. 执行删除操作
	var req struct {
		Path string `json:"path"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求数据格式错误",
		})
		return
	}

	if req.Path == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "路径不能为空",
		})
		return
	}

	fileInfo, err := conn.SftpClient.Stat(req.Path)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "路径不存在: " + err.Error(),
		})
		return
	}

	var errDelete error
	if fileInfo.IsDir() {
		errDelete = conn.DeleteDirectory(req.Path)
	} else {
		errDelete = conn.SftpClient.Remove(req.Path)
	}

	if errDelete != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除失败: " + errDelete.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

// ! 批量删除目录和文件
func BatchDelete(c *gin.Context) {
	// 获取Token
	token := getSFTPToken(c)
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "SFTP Token不能为空",
		})
		return
	}

	// 获取连接实例
	conn, err := utils.SFTPConnManager.GetConn(token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    50014,
			"message": "SFTP连接失效: " + err.Error(),
		})
		return
	}

	var req []models.FileInfo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求数据格式错误",
		})
		return
	}

	if len(req) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "路径不能为空",
		})
		return
	}

	// 校验所有路径是否存在
	type pathCheckResult struct {
		Path     string // 路径
		IsDir    bool   // 是否目录
		NotExist bool   // 是否不存在
		Error    string // 错误信息
	}
	var checkResults []pathCheckResult

	for _, file := range req {
		path := file.Path
		fileInfo, err := conn.SftpClient.Stat(path)

		// 路径不存在：记录错误，继续下一个
		if err != nil {
			checkResults = append(checkResults, pathCheckResult{
				Path:     path,
				NotExist: true,
				Error:    "路径不存在: " + err.Error(),
			})
			continue
		}

		// 路径存在：记录类型（文件/目录）
		checkResults = append(checkResults, pathCheckResult{
			Path:  path,
			IsDir: fileInfo.IsDir(),
		})
	}

	// 收集所有校验失败的路径
	var validateErrors []string
	for _, res := range checkResults {
		if res.NotExist {
			validateErrors = append(validateErrors, res.Error)
		}
	}

	// 如果有校验失败，直接返回所有错误，不执行删除
	if len(validateErrors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "部分路径校验失败",
			"errors":  validateErrors,
		})
		return
	}

	// 所有路径都存在，执行批量删除（容错不中断）
	var successCount int      // 成功删除数量
	var deleteErrors []string // 删除失败的错误信息

	for _, res := range checkResults {
		path := res.Path
		var errDel error

		// 目录 / 文件分别删除
		if res.IsDir {
			errDel = conn.DeleteDirectory(path)
		} else {
			errDel = conn.SftpClient.Remove(path)
		}

		// 删除失败：记录错误，继续下一个
		if errDel != nil {
			deleteErrors = append(deleteErrors, "删除失败: "+errDel.Error())
			continue
		}

		successCount++
	}

	if len(deleteErrors) == 0 {
		// 全部删除成功
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "批量删除成功",
		})
	} else if successCount == 0 {
		// 全部删除失败（重点补充）
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "批量删除失败:" + strings.Join(deleteErrors, ", "),
		})
	} else {
		// 部分删除成功，部分失败
		c.JSON(http.StatusPartialContent, gin.H{
			"code":    206,
			"message": "部分删除成功:" + strings.Join(deleteErrors, ", "),
		})
	}
}

// ! RenamePath 重命名（适配Token）
func RenamePath(c *gin.Context) {
	// 1. 获取Token
	token := getSFTPToken(c)
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "SFTP Token不能为空",
		})
		return
	}

	// 2. 获取连接实例
	conn, err := utils.SFTPConnManager.GetConn(token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    50014,
			"message": "SFTP连接失效: " + err.Error(),
		})
		return
	}

	// 3. 执行重命名
	var req struct {
		OldPath string `json:"oldPath"`
		NewName string `json:"newName"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求数据格式错误",
		})
		return
	}

	if req.OldPath == "" || req.NewName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "原路径和新名称不能为空",
		})
		return
	}

	if strings.Contains(req.NewName, "/") || req.NewName == "." || req.NewName == ".." {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "名称包含非法字符",
		})
		return
	}

	oldDir := filepath.Dir(req.OldPath)
	newPath := filepath.Join(oldDir, req.NewName)

	if err := conn.SftpClient.Rename(req.OldPath, newPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "重命名失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "重命名成功",
	})
}

// ! DownloadDirectory 下载目录
func DownloadDirectory(c *gin.Context) {
	// 获取Token
	token := getSFTPToken(c)
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "SFTP Token不能为空",
		})
		return
	}

	// 获取连接实例
	conn, err := utils.SFTPConnManager.GetConn(token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    50014,
			"message": "SFTP连接失效: " + err.Error(),
		})
		return
	}

	// 获取要下载的远程目录路径
	remoteDir := c.Query("path")
	if remoteDir == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "目录路径不能为空",
		})
		return
	}

	// 校验路径是否存在且是目录
	fileInfo, err := conn.SftpClient.Stat(remoteDir)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "目录不存在: " + err.Error(),
		})
		return
	}
	if !fileInfo.IsDir() {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "该路径不是目录，请使用文件下载接口",
		})
		return
	}

	// 设置下载响应头
	zipName := filepath.Base(remoteDir) + ".zip"
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", fmt.Sprintf(
		"attachment; filename=\"%s\"; filename*=UTF-8''%s",
		escapeFileName(zipName),
		url.QueryEscape(zipName),
	))
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.Status(http.StatusOK)

	// 创建 zip 流式写入器
	zipWriter := zip.NewWriter(c.Writer)
	defer zipWriter.Close()

	// 递归遍历并压缩目录
	err = addDirToZip(zipWriter, conn, remoteDir, remoteDir)
	if err != nil {
		log.Printf("目录打包下载失败：%s，错误：%v", remoteDir, err)
		return
	}

	log.Printf("目录打包下载完成：%s", remoteDir)
}

// --- DownloadDirectory辅助函数：递归把SFTP目录写入zip ---
func addDirToZip(zipW *zip.Writer, conn *utils.SFTPConnection, remoteRoot, currentDir string) error {
	// 读取当前目录内容
	files, err := conn.SftpClient.ReadDir(currentDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		// 跳过隐藏文件（和你文件列表逻辑一致）
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}

		remotePath := filepath.Join(currentDir, file.Name())
		// zip 内部路径（保持目录结构）
		zipPath, _ := filepath.Rel(remoteRoot, remotePath)

		if file.IsDir() {
			// 递归子目录
			if err := addDirToZip(zipW, conn, remoteRoot, remotePath); err != nil {
				return err
			}
		} else {
			// 写入文件到 zip
			if err := writeFileToZip(zipW, conn, remotePath, zipPath); err != nil {
				return err
			}
		}
	}
	return nil
}

// 写入单个文件到 zip
func writeFileToZip(zipW *zip.Writer, conn *utils.SFTPConnection, remotePath, zipPath string) error {
	// 打开远程文件
	srcFile, err := conn.SftpClient.Open(remotePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 创建 zip 中的文件
	zipFile, err := zipW.Create(zipPath)
	if err != nil {
		return err
	}

	// 流式复制（不占内存）
	_, err = io.Copy(zipFile, srcFile)
	return err
}

// ! SearchFiles 递归搜索文件和目录（不区分大小写模糊匹配）
func SearchFiles(c *gin.Context) {
	// 1. 获取Token
	token := getSFTPToken(c)
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "SFTP Token不能为空",
		})
		return
	}

	// 2. 通过Token获取对应的连接实例
	conn, err := utils.SFTPConnManager.GetConn(token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    50014,
			"message": "SFTP连接失效: " + err.Error(),
		})
		return
	}

	// 3. 获取搜索参数
	searchPath := c.Query("path")
	if searchPath == "" {
		searchPath = "/"
	}
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "搜索关键字不能为空",
		})
		return
	}

	// 4. 验证搜索路径是否存在
	fileInfo, err := conn.SftpClient.Stat(searchPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "搜索路径不存在: " + err.Error(),
		})
		return
	}
	if !fileInfo.IsDir() {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "搜索路径不是目录",
		})
		return
	}

	// 5. 递归搜索
	var results []models.SearchFileInfo
	recursiveSearch(conn, searchPath, keyword, &results)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "搜索完成",
		"data": gin.H{
			"total":   len(results),
			"results": results,
		},
	})
}

// recursiveSearch 递归搜索目录，匹配文件名（不区分大小写）
func recursiveSearch(conn *utils.SFTPConnection, dirPath, keyword string, results *[]models.SearchFileInfo) error {
	files, err := conn.SftpClient.ReadDir(dirPath)
	if err != nil {
		// 跳过无权限或无法读取的目录
		return nil
	}

	lowerKeyword := strings.ToLower(keyword)

	for _, f := range files {
		// 跳过隐藏文件/目录
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}

		fullPath := filepath.Join(dirPath, f.Name())

		// 模糊匹配文件名（不区分大小写）
		if strings.Contains(strings.ToLower(f.Name()), lowerKeyword) {
			*results = append(*results, models.SearchFileInfo{
				FileInfo: models.FileInfo{
					Name:     f.Name(),
					Path:     fullPath,
					IsDir:    f.IsDir(),
					Size:     f.Size(),
					Modified: f.ModTime().Format(time.RFC3339),
				},
				ParentPath: dirPath,
			})
		}

		// 递归进入子目录
		if f.IsDir() {
			recursiveSearch(conn, fullPath, keyword, results)
		}
	}
	return nil
}
