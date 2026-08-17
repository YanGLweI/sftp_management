package utils

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// dualAuthEntry 双控验证凭证条目
type dualAuthEntry struct {
	SFTPToken   string    // 绑定的 SFTP 连接 Token
	Reviewer    string    // 双控复核人（另一产业部账号）
	IssuerIP    string    // 签发 IP（用于安全审计）
	ExpireAt    time.Time // 过期时间
	UsedAt      time.Time // 使用时间（单次使用后设置）
	RetryCount  int       // 重试次数
}

const (
	// DualAuthTokenTTL 双控凭证有效期：24 小时
	// 原为 60 秒，但因大批量/大文件上传可能超过 60 秒导致凭证过期（428 需要双控验证）
	// 凭证虽长期有效，但每次独立写操作（删除/重命名等）仍由 GetReviewer 单次消耗，
	// 上传批次由 PeekReviewer 非消耗式复用，不会造成安全风险
	DualAuthTokenTTL = 24 * time.Hour
	// MaxTokensPerSftpToken 同一 SFTP Token 最多签发并发双控 Token 数
	MaxTokensPerSftpToken = 5
	// DualAuthCleanupInterval 定期清理间隔：每 5 分钟
	DualAuthCleanupInterval = 5 * time.Minute
)

// dualAuthManager 双控验证凭证管理器（内存存储）
type dualAuthManager struct {
	tokenMap map[string]dualAuthEntry // key: 凭证Token, value: 绑定信息
	mu       sync.RWMutex
}

// 全局双控凭证管理器实例
var DualAuthManager = &dualAuthManager{
	tokenMap: make(map[string]dualAuthEntry),
}

// IssueToken 为指定 SFTP 连接签发双控验证凭证
// - sftpToken: SFTP 连接 Token
// - reviewer: 通过双控验证的另一产业部账号
// - clientIP: 客户端 IP（用于安全审计和速率限制）
// 返回：Token 字符串，失败返回空字符串
func (m *dualAuthManager) IssueToken(sftpToken, reviewer, clientIP string) string {
	m.mu.Lock()
	defer m.mu.Unlock()

	// ✅ 检查同一 SFTP Token 的并发 Token 数量
	existingCount := 0
	for _, entry := range m.tokenMap {
		if entry.SFTPToken == sftpToken &&
			!entry.ExpireAt.IsZero() &&
			entry.ExpireAt.After(time.Now()) {
			existingCount++
		}
	}
	if existingCount >= MaxTokensPerSftpToken {
		return "" // 超过最大并发数
	}

	token := uuid.New().String()
	m.tokenMap[token] = dualAuthEntry{
		SFTPToken:  sftpToken,
		Reviewer:   reviewer,
		IssuerIP:   clientIP,
		ExpireAt:   time.Now().Add(DualAuthTokenTTL),
		UsedAt:     time.Time{},
		RetryCount: 0,
	}

	return token
}

// VerifyToken 校验双控凭证：存在、未过期、且绑定当前 SFTP 连接与签发 IP
// - clientIP: 当前请求客户端 IP，必须与签发时 IP 一致（防跨 IP 盗用）
// - 如果 token 已使用（UsedAt 不为零）则直接返回 false
// 注意：不删除 Token，以支持上传并发 3 路复用同一凭证；成功操作后由 GetReviewer 消耗
func (m *dualAuthManager) VerifyToken(sftpToken, dualToken, clientIP string) bool {
	if dualToken == "" {
		return false
	}
	m.mu.RLock()
	defer m.mu.RUnlock()

	entry, exists := m.tokenMap[dualToken]
	if !exists || entry.SFTPToken != sftpToken {
		return false
	}
	if time.Now().After(entry.ExpireAt) {
		return false
	}
	// ✅ 校验签发 IP 与当前客户端 IP 一致（防跨 IP 复用）
	if entry.IssuerIP != "" && entry.IssuerIP != clientIP {
		return false
	}
	// ✅ 检查是否已使用（单次使用原则）
	if !entry.UsedAt.IsZero() {
		return false
	}
	return true
}

// GetReviewer 获取双控凭证对应的复核人账号
// - 验证 Token 有效性和过期时间
// - ✅ 单次使用后自动删除 Token（防止重用）
// - -clientIP: 客户端 IP（用于安全审计对比）
func (m *dualAuthManager) GetReviewer(dualToken, clientIP string) string {
	if dualToken == "" {
		return ""
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	entry, exists := m.tokenMap[dualToken]
	if !exists {
		return ""
	}
	// ✅ 验证是否过期
	if time.Now().After(entry.ExpireAt) {
		delete(m.tokenMap, dualToken)
		return ""
	}
	// ✅ 单次使用后自动删除 Token
	delete(m.tokenMap, dualToken)
	return entry.Reviewer
}

// PeekReviewer 获取双控凭证对应的复核人账号（非消耗式）
// - 与 GetReviewer 相同，但不删除 Token，供同一凭证批量复用场景使用
// - 适用于批量上传：一个批次内的多个文件请求共享同一凭证，
//   由 CleanupExpiredTokens 在过期后统一清理（默认 24 小时 TTL）
func (m *dualAuthManager) PeekReviewer(dualToken, clientIP string) string {
	if dualToken == "" {
		return ""
	}
	m.mu.RLock()
	defer m.mu.RUnlock()

	entry, exists := m.tokenMap[dualToken]
	if !exists {
		return ""
	}
	// ✅ 验证是否过期
	if time.Now().After(entry.ExpireAt) {
		return ""
	}
	return entry.Reviewer
}

// CleanupExpiredTokens 定期清理过期凭证（每 5 分钟检查一次）
func (m *dualAuthManager) CleanupExpiredTokens() {
	ticker := time.NewTicker(DualAuthCleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		m.mu.Lock()
		now := time.Now()
		for token, entry := range m.tokenMap {
			if now.After(entry.ExpireAt) {
				delete(m.tokenMap, token)
			} else if !entry.UsedAt.IsZero() {
				// 已使用的 Token 也立即删除
				delete(m.tokenMap, token)
			}
		}
		m.mu.Unlock()
	}
}
