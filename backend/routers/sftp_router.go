package routers

import (
	"sftpbackend/controller"

	"github.com/gin-gonic/gin"
)

// registerSFTPRouter 注册SFTP相关路由
func registerSFTPRouter(r *gin.Engine) {
	sftpGroup := r.Group("/sftp")
	// sftpGroup.Use(middleware.AuthMiddleware())
	{
		// SFTP登录:请求体{"username":"string","password":"string"}
		sftpGroup.POST("/login", controller.LoginSFTP)
		// 获取目录下的文件列表:请求参数 /files?path=string
		sftpGroup.GET("/files", controller.GetFiles)
		// SFTP登出
		sftpGroup.GET("/logout", controller.LogoutSFTP)
		// 上传文件
		sftpGroup.POST("/upload", controller.UploadFile)
		// 创建目录:请求体{"path":"string","name":"string"}
		sftpGroup.POST("/mkdir", controller.CreateFolder)
		// 下载文件:请求参数 /download?path=string
		sftpGroup.GET("/download", controller.DownloadFile)
		// sftpGroup.Handle("GET", "/download", controller.DownloadFile)
		// sftpGroup.Handle("HEAD", "/download", controller.DownloadFile)
		// 删除目录或文件：请求体{"path":"string"}
		sftpGroup.POST("/delete", controller.DeletePath)
		// 重命名文件或目录：请求体{"oldPath":"string","newName":"string"}
		sftpGroup.POST("/rename", controller.RenamePath)
		// 下载目录：请求参数 /downloaddir?path=string
		sftpGroup.GET("/downloaddir", controller.DownloadDirectory)
		// 批量删除
		sftpGroup.POST("/batchdelete", controller.BatchDelete)
		// 递归搜索文件或目录：请求参数 /search?path=string&keyword=string
		sftpGroup.GET("/search", controller.SearchFiles)
	}
}
