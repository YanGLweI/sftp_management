package utils

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// dualAuthEntry 双控验证凭证条目
type dualAuthEntry struct {
	SftpToken string    // 绑定的SFTP连接Token
	Reviewer  string    // 双控复核人（另一产业部账号）
	ExpireAt  time.Time // 过期时间
}

// 双控凭证默认有效期：60秒（供上传并发3路复用同一凭证）
const DualAuthTokenTTL = 60 * time.Second

// dualAuthManager 双控验证凭证管理器（内存存储）
type dualAuthManager struct {
	tokenMap map[string]dualAuthEntry // key: 凭证Token, value: 绑定信息
	mu       sync.RWMutex
}

// 全局双控凭证管理器实例
var DualAuthManager = &dualAuthManager{
	tokenMap: make(map[string]dualAuthEntry),
}

// IssueToken 为指定SFTP连接签发双控验证凭证（可复用，有效期 DualAuthTokenTTL）
// reviewer: 通过双控验证的另一产业部账号
func (m *dualAuthManager) IssueToken(sftpToken, reviewer string) string {
	m.mu.Lock()
	defer m.mu.Unlock()

	token := uuid.New().String()
	m.tokenMap[token] = dualAuthEntry{
		SftpToken: sftpToken,
		Reviewer:  reviewer,
		ExpireAt:  time.Now().Add(DualAuthTokenTTL),
	}
	return token
}

// VerifyToken 校验双控凭证：存在、未过期、且绑定当前SFTP连接
func (m *dualAuthManager) VerifyToken(sftpToken, dualToken string) bool {
	if dualToken == "" {
		return false
	}
	m.mu.RLock()
	defer m.mu.RUnlock()

	entry, exists := m.tokenMap[dualToken]
	if !exists || entry.SftpToken != sftpToken {
		return false
	}
	if time.Now().After(entry.ExpireAt) {
		return false
	}
	return true
}

// GetReviewer 获取双控凭证对应的复核人账号
func (m *dualAuthManager) GetReviewer(dualToken string) string {
	if dualToken == "" {
		return ""
	}
	m.mu.RLock()
	defer m.mu.RUnlock()

	entry, exists := m.tokenMap[dualToken]
	if !exists {
		return ""
	}
	return entry.Reviewer
}

// CleanExpiredTokens 定期清理过期凭证（每5分钟检查一次）
func (m *dualAuthManager) CleanExpiredTokens() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		m.mu.Lock()
		now := time.Now()
		for token, entry := range m.tokenMap {
			if now.After(entry.ExpireAt) {
				delete(m.tokenMap, token)
			}
		}
		m.mu.Unlock()
	}
}
