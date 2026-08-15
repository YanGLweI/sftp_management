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
		&models.SftpLog{},
		&models.SftpUsers{},
		&models.Contact{},
		&models.UpdateMain{},
		&models.UpdateDetails{},
		&models.Scheduler{},
		&models.SystemSecurity{},
		&models.SystemSecurityStandard{},
		&models.LocalUser{},
		&models.Role{},
		&models.RoleMenu{},
		&models.RoleLDAPGroup{},
		&models.PasswordHistory{},
		&models.PasswordPolicy{},
		&models.SFTPModuleConfig{}, // 新增：SFTP 模块配置表
		&models.LDAPConfig{},        // 新增：LDAP 配置表（包含 CertFilename 字段）
	)
	if err != nil {
		logrus.Fatal("数据库表迁移失败:", err)
	}

	// 初始化数据（调度器、密码策略、超级管理员角色、默认admin账号、系统安全标准）
	common.InitData()

	// 初始化 SFTP 模块默认配置
	err = models.InitDefaultConfigs()
	if err != nil {
		logrus.Fatalf("初始化 SFTP 模块配置失败：%v", err)
	}

	// 确保超级管理员角色拥有新增的 SFTP 模块管理菜单权限（兼容已有部署）
	common.EnsureSuperAdminSftpModuleMenus()

	// 确保超级管理员角色拥有新增的 LDAP 管理菜单权限（兼容已有部署）
	common.EnsureSuperAdminLDAPManagementMenu()

	// 启动调度器
	go scheduler.Run()

	// 启动过期连接清理（8小时未使用则清理）
	go utils.SFTPConnManager.CleanExpiredConns(8 * time.Hour)

	// 启动双控凭证过期清理
	go utils.DualAuthManager.CleanupExpiredTokens()

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
