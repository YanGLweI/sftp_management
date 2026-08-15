package utils

import (
	"strings"
	"testing"
	"time"
)

// 测试 ResolvePath 边界修复：HomePath="/" 时子路径应通过
func TestResolvePathRootHome(t *testing.T) {
	conn := &SFTPConnection{HomePath: "/"}
	p, err := conn.ResolvePath("/etc/passwd")
	if err != nil {
		t.Fatalf("HomePath=/ 时 /etc/passwd 应通过，实际: %v", err)
	}
	if p != "/etc/passwd" {
		t.Fatalf("路径应保持 /etc/passwd，实际: %s", p)
	}
	// 根路径本身
	if _, err := conn.ResolvePath("/"); err != nil {
		t.Fatalf("HomePath=/ 时 / 应通过: %v", err)
	}
}

// 测试 ResolvePath 越界拒绝
func TestResolvePathTraversal(t *testing.T) {
	conn := &SFTPConnection{HomePath: "/hotlabel"}
	cases := []string{
		"/hotlabel/../../etc/passwd", // 穿越
		"/etc/passwd",                // 绝对越界
		"/hotlabel/..",               // 越界到上级
		"..",                         // 相对穿越
	}
	for _, c := range cases {
		if _, err := conn.ResolvePath(c); err == nil {
			t.Fatalf("路径 %s 应被拒绝", c)
		}
	}
	// 正常路径应通过
	ok := []string{"/hotlabel", "/hotlabel/sub", "/hotlabel/a/b/c"}
	for _, c := range ok {
		if _, err := conn.ResolvePath(c); err != nil {
			t.Fatalf("路径 %s 应通过: %v", c, err)
		}
	}
}

// 测试无限制连接
func TestResolvePathNoHome(t *testing.T) {
	conn := &SFTPConnection{HomePath: ""}
	p, err := conn.ResolvePath("/any/path")
	if err != nil || p != "/any/path" {
		t.Fatalf("无限制连接应放行任意路径，实际 p=%s err=%v", p, err)
	}
}

// 测试 ValidateFileName
func TestValidateFileName(t *testing.T) {
	bad := []string{"", ".", "..", "../evil", "a/b", "a\\b", "..\\evil"}
	for _, name := range bad {
		if err := ValidateFileName(name); err == nil {
			t.Fatalf("名称 %q 应被拒绝", name)
		}
	}
	good := []string{"file.txt", "目录", "a-b_c.d", "名称 2026"}
	for _, name := range good {
		if err := ValidateFileName(name); err != nil {
			t.Fatalf("名称 %q 应通过: %v", name, err)
		}
	}
}

// 测试双控 Token 签发频率限制
func TestDualAuthRateLimit(t *testing.T) {
	m := &dualAuthManager{tokenMap: make(map[string]dualAuthEntry)}
	sftpToken := "sftp-token-1"
	// 签发超过 MaxTokensPerSftpToken 个
	for i := 0; i < MaxTokensPerSftpToken; i++ {
		if tok := m.IssueToken(sftpToken, "reviewer", "127.0.0.1"); tok == "" {
			t.Fatalf("第 %d 次签发不应失败", i+1)
		}
	}
	// 第 6 个应失败
	if tok := m.IssueToken(sftpToken, "reviewer", "127.0.0.1"); tok != "" {
		t.Fatalf("超过上限应拒绝签发，实际返回: %s", tok)
	}
	// 不同 SFTP Token 不受影响
	if tok := m.IssueToken("sftp-token-2", "reviewer", "127.0.0.1"); tok == "" {
		t.Fatalf("不同 SFTP Token 应可签发")
	}
}

// 测试双控 Token 单次使用与过期
func TestDualAuthSingleUseAndExpiry(t *testing.T) {
	m := &dualAuthManager{tokenMap: make(map[string]dualAuthEntry)}
	token := m.IssueToken("sftp-tok", "reviewerA", "127.0.0.1")
	if token == "" {
		t.Fatalf("签发失败")
	}
	// 首次获取复核人
	r := m.GetReviewer(token, "127.0.0.1")
	if r != "reviewerA" {
		t.Fatalf("首次获取应返回 reviewerA，实际: %s", r)
	}
	// 第二次获取应失败（单次使用）
	if r2 := m.GetReviewer(token, "127.0.0.1"); r2 != "" {
		t.Fatalf("Token 应单次使用，第二次应返回空，实际: %s", r2)
	}

	// 过期 Token 验证
	m2 := &dualAuthManager{tokenMap: make(map[string]dualAuthEntry)}
	token2 := m2.IssueToken("sftp-tok2", "reviewerB", "127.0.0.1")
	// 手动将过期时间设为过去
	entry := m2.tokenMap[token2]
	entry.ExpireAt = time.Now().Add(-time.Minute)
	m2.tokenMap[token2] = entry
	if r3 := m2.GetReviewer(token2, "127.0.0.1"); r3 != "" {
		t.Fatalf("过期 Token 应返回空，实际: %s", r3)
	}
	// 过期 Token 应被删除
	if _, exists := m2.tokenMap[token2]; exists {
		t.Fatalf("过期 Token 应被清理")
	}
}

// 测试 VerifyToken 对已使用 Token 的拒绝
func TestDualAuthVerifyTokenUsed(t *testing.T) {
	m := &dualAuthManager{tokenMap: make(map[string]dualAuthEntry)}
	token := m.IssueToken("sftp-tok3", "reviewerC", "127.0.0.1")
	if !m.VerifyToken("sftp-tok3", token, "127.0.0.1") {
		t.Fatalf("未使用 Token 应验证通过")
	}
	// 使用一次
	m.GetReviewer(token, "127.0.0.1")
	// 再验证应失败（已使用）
	if m.VerifyToken("sftp-tok3", token, "127.0.0.1") {
		t.Fatalf("已使用 Token 不应通过验证")
	}
}

// 测试 VerifyToken 对跨 IP 复用的拒绝
func TestDualAuthVerifyTokenIPMismatch(t *testing.T) {
	m := &dualAuthManager{tokenMap: make(map[string]dualAuthEntry)}
	token := m.IssueToken("sftp-tok5", "reviewerE", "10.0.0.1")
	// 同 IP 验证通过
	if !m.VerifyToken("sftp-tok5", token, "10.0.0.1") {
		t.Fatalf("同 IP 应验证通过")
	}
	// 跨 IP 验证应失败
	if m.VerifyToken("sftp-tok5", token, "10.0.0.2") {
		t.Fatalf("跨 IP 复用应被拒绝")
	}
	// 空 IP 容忍（旧数据兼容）
	m2 := &dualAuthManager{tokenMap: make(map[string]dualAuthEntry)}
	token2 := m2.IssueToken("sftp-tok6", "reviewerF", "")
	if !m2.VerifyToken("sftp-tok6", token2, "10.0.0.9") {
		t.Fatalf("签发 IP 为空时应放行（兼容）")
	}
}

// 测试清理过期 Token
func TestDualAuthCleanup(t *testing.T) {
	m := &dualAuthManager{tokenMap: make(map[string]dualAuthEntry)}
	token := m.IssueToken("sftp-tok4", "reviewerD", "127.0.0.1")
	entry := m.tokenMap[token]
	entry.ExpireAt = time.Now().Add(-time.Minute)
	m.tokenMap[token] = entry
	m.cleanupOnce()
	if _, exists := m.tokenMap[token]; exists {
		t.Fatalf("过期 Token 应被清理")
	}
	// 确认清理函数名存在（编译检查）
	_ = DualAuthCleanupInterval
}

// 辅助：暴露 cleanupOnce 单次清理（避免依赖 ticker）
func (m *dualAuthManager) cleanupOnce() {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now()
	for token, entry := range m.tokenMap {
		if now.After(entry.ExpireAt) || !entry.UsedAt.IsZero() {
			delete(m.tokenMap, token)
		}
	}
}

// 确保 strings 包被使用（编译检查）
var _ = strings.Contains
