package routers

import (
	"sftpbackend/controller"
	"sftpbackend/middleware"

	"github.com/gin-gonic/gin"
)

// registerUserRouter 注册用户相关路由
func registerUserRouter(r *gin.Engine) {
	userGroup := r.Group("/user")
	userGroup.Use(middleware.AuthMiddleware()) // 统一鉴权中间件
	{
		// ! 用户管理（账号管理菜单权限）
		// 添加一个用户：请求体{"name":"string","password":"string","emailOrNot":bool,"NoExpire":bool}
		userGroup.POST("/account", middleware.RequireRoute("SftpUser"), controller.AddUser)
		// 获取全部用户和用户搜索：请求参数 /account/:page/:limit?username=string
		userGroup.GET("/account/:page/:limit", middleware.RequireRoute("SftpUser"), controller.GetUserList)
		// 修改一个用户密码：请求体{"name":"string","password":"string","emailOrNot":bool,"noExpire":bool}
		userGroup.PUT("/account", middleware.RequireRoute("SftpUser"), controller.UpdateUser)
		// 删除一个用户：根据 ID 删除 /account/:id
		userGroup.DELETE("/account/:id", middleware.RequireRoute("SftpUser"), controller.DeleteUser)
		// 批量删除用户：请求体{"ids":[1,2,3]}
		userGroup.DELETE("/account", middleware.RequireRoute("SftpUser"), controller.DeleteUsers)
		// 下载私钥 /download-key/:username
		userGroup.GET("/download-key/:username", middleware.RequireRoute("SftpUser"), controller.DownloadPrivateKey)
		// 获取登录用户信息（公共）
		userGroup.GET("/info", controller.GetUserInfo)
		// ! 邮件发送：请求体：{"to":[]string,"cc":[]string,"subject":"string"}（通讯邮箱菜单权限）
		userGroup.POST("/email", middleware.RequireRoute("Contacts"), controller.SendEmail)
		// ! 主动退出登录（公共）
		userGroup.GET("/logout", controller.Logout)
		// ! 平台日志:/log/:page/:limit?datetime=string&username=string&logtype=string
		userGroup.GET("/log/:page/:limit", middleware.RequireRoute("PlatformLog"), controller.GetLogList)
		// ! sftp 传输日志
		userGroup.GET("/log/sftplog/:date", middleware.RequireRoute("SftpLog"), controller.GetSftpLog)
		// ! SFTP登录与操作日志：/log/sftploglist/:page/:limit?datetime=string&username=string
		userGroup.GET("/log/sftploglist/:page/:limit", middleware.RequireRoute("SftpLog"), controller.GetSftpLogList)
		// ! 中国联通日志：/log/chinaunicomloglist/:page/:limit?datetime=string&username=string
		userGroup.GET("/log/chinaunicomloglist/:page/:limit", middleware.RequireRoute("ChinaUnicomLog"), controller.GetChinaUnicomLogList)
	}
}
