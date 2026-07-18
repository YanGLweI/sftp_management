// sftpbackend/graceful/graceful.go
package graceful

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sftpbackend/dao"
	"sftpbackend/scheduler"

	"github.com/sirupsen/logrus"
)

// Option 可选配置（用于自定义超时时间）
type Option struct {
	Timeout time.Duration // 优雅停止超时时间，默认30秒
}

// 默认配置
var defaultOption = Option{
	Timeout: 30 * time.Second,
}

// Shutdown 封装优雅停止的核心逻辑
// 参数：
//
//	server: HTTP服务实例（Gin的http.Server）
//	opts: 可选配置（如自定义超时时间）
func Shutdown(server *http.Server, opts ...Option) {
	// 合并配置（使用默认值或自定义值）
	cfg := defaultOption
	if len(opts) > 0 {
		if opts[0].Timeout > 0 {
			cfg.Timeout = opts[0].Timeout
		}
	}

	// 1. 捕获停止信号（SIGINT: Ctrl+C；SIGTERM: kill命令）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // 阻塞，直到收到停止信号
	logrus.Info("开始执行优雅停止流程...")

	// 2. 优雅停止超时控制
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	// 3. 停止调度器
	scheduler.Stop()

	// 4. 优雅关闭HTTP服务
	if server != nil {
		if err := server.Shutdown(ctx); err != nil {
			logrus.Fatalf("HTTP服务优雅关闭失败: %v", err)
		}
		logrus.Info("HTTP服务已优雅关闭")
	} else {
		logrus.Warn("HTTP服务实例为nil，无需关闭")
	}

	// 5. 关闭数据库连接
	dao.Close()
	logrus.Info("数据库连接已关闭")

	// 6. 所有资源关闭完成
	logrus.Info("优雅停止完成，程序正常退出")
}
