package routers

import (
	"sftpbackend/controller"
	"sftpbackend/middleware"

	"github.com/gin-gonic/gin"
)

// registerSFTPRouter 注册 SFTP 相关路由
func registerSFTPRouter(r *gin.Engine) {
	sftpGroup := r.Group("/sftp")
	// sftpGroup.Use(middleware.AuthMiddleware())
	{
		// 公共模块配置（/file 登录页无 token 渲染表单用，不返回角色白名单）
		sftpGroup.GET("/module-configs", controller.GetPublicModuleConfigs)
		// SFTP 登录：请求体{"username":"string","password":"string","loginType":"string"}
		sftpGroup.POST("/login", controller.LoginSFTP)
		// 双控验证：请求体{"username":"string","password":"string"} + 请求头 X-SFTP-Token
		sftpGroup.POST("/dualverify", controller.DualVerify)
		// 获取目录下的文件列表：请求参数 /files?path=string
		sftpGroup.GET("/files", controller.GetFiles)
		// SFTP 登出
		sftpGroup.GET("/logout", controller.LogoutSFTP)
		// 写操作（中国联通连接需双控验证，中间件强制校验 X-Dual-Token）
		writeGroup := sftpGroup.Group("")
		writeGroup.Use(middleware.DualAuthMiddleware())
		{
			// 上传文件
			writeGroup.POST("/upload", controller.UploadFile)
			// 创建目录：请求体{"path":"string","name":"string"}
			writeGroup.POST("/mkdir", controller.CreateFolder)
			// 删除目录或文件：请求体{"path":"string"}
			writeGroup.POST("/delete", controller.DeletePath)
			// 重命名文件或目录：请求体{"oldPath":"string","newName":"string"}
			writeGroup.POST("/rename", controller.RenamePath)
			// 批量删除
			writeGroup.POST("/batchdelete", controller.BatchDelete)
		}
		// 下载文件：请求参数 /download?path=string
		sftpGroup.GET("/download", controller.DownloadFile)
		// 下载目录：请求参数 /downloaddir?path=string
		sftpGroup.GET("/downloaddir", controller.DownloadDirectory)
		// 递归搜索文件或目录：请求参数 /search?path=string&keyword=string
		sftpGroup.GET("/search", controller.SearchFiles)
	}
}
