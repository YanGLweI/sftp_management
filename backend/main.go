package main

import (
	"fmt"
	"net/http"
	"sftpbackend/common"
	"sftpbackend/config"
	"sftpbackend/dao"
	"sftpbackend/graceful"
	"sftpbackend/models"
	"sftpbackend/routers"
	"sftpbackend/scheduler"
	"sftpbackend/utils"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	// 连接数据库
	err := dao.InitMariaDB()
	if err != nil {
		logrus.Fatal("连接数据库失败:", err)
	}
	defer dao.Close()

	// 数据库表结构迁移
	err = dao.DB.AutoMigrate(
		&models.Log{},
		&models.SftpUsers{},
		&models.Contact{},
		&models.UpdateMain{},
		&models.UpdateDetails{},
		&models.Scheduler{},
		&models.SystemSecurity{},
		&models.SystemSecurityStandard{},
	)
	if err != nil {
		logrus.Fatal("数据库表迁移失败:", err)
	}

	// 初始化数据
	common.InitData()

	// 启动调度器
	go scheduler.Run()

	// 启动过期连接清理（8小时未使用则清理）
	go utils.SFTPConnManager.CleanExpiredConns(8 * time.Hour)

	// 设置路由
	r := routers.SetupRouter()
	// 从配置文件中获取系统配置
	sysConfig := config.GlobalConfig.System
	addr := fmt.Sprintf(":%d", sysConfig.Port)

	// 替换 r.Run() 为 http.Server（支持优雅关闭）
	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// 异步启动HTTP服务
	go func() {
		logrus.Infof("Gin服务启动成功，监听地址: %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("HTTP服务启动失败: %v", err)
		}
	}()

	// 调用封装的优雅停止函数（停止逻辑）
	// 可选：自定义超时时间，如 graceful.Shutdown(server, graceful.Option{Timeout: 40 * time.Second})
	graceful.Shutdown(server)

}
