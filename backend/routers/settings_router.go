package routers

import (
	"sftpbackend/controller"
	"sftpbackend/middleware"

	"github.com/gin-gonic/gin"
)

// registerSettingsRouter 注册平台设置相关路由
func registerSettingsRouter(r *gin.Engine) {
	settingsGroup := r.Group("/settings")
	settingsGroup.Use(middleware.AuthMiddleware())
	{
		// 获取所有菜单树（角色管理/本地账号页面共用，仅鉴权）
		settingsGroup.GET("/menus", controller.GetAllMenus)

		// 角色管理
		settingsGroup.GET("/roles", middleware.RequireRoute("RoleManagement"), controller.GetRoleList)
		settingsGroup.POST("/role", middleware.RequireRoute("RoleManagement"), controller.CreateRole)
		settingsGroup.PUT("/role/:id", middleware.RequireRoute("RoleManagement"), controller.UpdateRole)
		settingsGroup.DELETE("/role/:id", middleware.RequireRoute("RoleManagement"), controller.DeleteRole)
		settingsGroup.GET("/role/:id", middleware.RequireRoute("RoleManagement"), controller.GetRoleDetail)
		// 角色下拉列表（本地账号页面共用，仅鉴权）
		settingsGroup.GET("/role/select", controller.GetRoleSelect)

		// 本地账号管理
		settingsGroup.GET("/localusers", middleware.RequireRoute("LocalUserManagement"), controller.GetLocalUserList)
		settingsGroup.POST("/localuser", middleware.RequireRoute("LocalUserManagement"), controller.CreateLocalUser)
		settingsGroup.PUT("/localuser/:id", middleware.RequireRoute("LocalUserManagement"), controller.UpdateLocalUser)
		settingsGroup.PUT("/localuser/:id/reset-password", middleware.RequireRoute("LocalUserManagement"), controller.ResetLocalUserPassword)
		settingsGroup.DELETE("/localuser/:id", middleware.RequireRoute("LocalUserManagement"), controller.DeleteLocalUser)

		// 密码策略
		settingsGroup.GET("/password-policy", middleware.RequireRoute("PasswordPolicy"), controller.GetPasswordPolicy)
		settingsGroup.PUT("/password-policy", middleware.RequireRoute("PasswordPolicy"), controller.UpdatePasswordPolicy)
		// 密码校验（修改密码弹框共用，仅鉴权）
		settingsGroup.POST("/password-policy/validate", controller.ValidatePassword)

		// LDAP 配置管理（超级管理员专用）
		ldapConfig := settingsGroup.Group("/ldap")
		ldapConfig.Use(middleware.RequireRoute("LDAPManagement"))
		{
			ldapConfig.GET("/config", controller.LdapConfigController.GetLDAPConfig)
			ldapConfig.PUT("/config", controller.LdapConfigController.SaveLDAPConfig)
			ldapConfig.POST("/test", controller.LdapConfigController.TestLDAPConnection)
		}

		// SFTP 模块配置管理（新增）
		settingsGroup.GET("/sftp-modules/all", middleware.RequireRoute("SftpModuleManagement"), controller.GetAllModuleConfigs)
		settingsGroup.GET("/sftp-modules/:name", middleware.RequireRoute("SftpModuleManagement"), controller.GetModuleConfig)
		settingsGroup.PUT("/sftp-modules/:name", middleware.RequireRoute("SftpModuleManagement"), controller.UpdateModuleConfig)
	}
}