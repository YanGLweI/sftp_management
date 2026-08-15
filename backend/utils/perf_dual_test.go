package utils

import (
	"testing"
	"time"
)

// 基准测试：双控 Token 签发耗时（要求 < 100ms）
func BenchmarkDualAuthIssueToken(b *testing.B) {
	m := &dualAuthManager{tokenMap: make(map[string]dualAuthEntry)}
	for i := 0; i < b.N; i++ {
		sftpToken := "bench-sftp-token"
		// 每轮先清空 tokenMap，避免并发数限制干扰基准
		if i%MaxTokensPerSftpToken == 0 {
			m.tokenMap = make(map[string]dualAuthEntry)
		}
		tok := m.IssueToken(sftpToken, "reviewer", "127.0.0.1")
		if tok == "" {
			b.Fatalf("签发失败")
		}
	}
}

// 基准测试：GetReviewer 单次使用耗时
func BenchmarkDualAuthGetReviewer(b *testing.B) {
	m := &dualAuthManager{tokenMap: make(map[string]dualAuthEntry)}
	for i := 0; i < b.N; i++ {
		tok := m.IssueToken("bench-tok", "reviewer", "127.0.0.1")
		if tok == "" {
			b.Fatalf("签发失败")
		}
		start := time.Now()
		if r := m.GetReviewer(tok, "127.0.0.1"); r == "" {
			b.Fatalf("获取复核人失败")
		}
		_ = start
	}
}

// 验证：批量签发耗时（1000 次整体耗时应远小于 100ms 每次）
func TestDualAuthIssueTiming(t *testing.T) {
	m := &dualAuthManager{tokenMap: make(map[string]dualAuthEntry)}
	start := time.Now()
	const rounds = 1000
	for i := 0; i < rounds; i++ {
		if i%MaxTokensPerSftpToken == 0 {
			m.tokenMap = make(map[string]dualAuthEntry)
		}
		tok := m.IssueToken("timing-tok", "reviewer", "127.0.0.1")
		if tok == "" {
			t.Fatalf("第 %d 次签发失败", i)
		}
	}
	elapsed := time.Since(start)
	perOp := elapsed / rounds
	t.Logf("1000 次签发总耗时: %v，单次平均: %v", elapsed, perOp)
	if perOp > 100*time.Millisecond {
		t.Fatalf("单次签发耗时 %v 超过 100ms 阈值", perOp)
	}
}
