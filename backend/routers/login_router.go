package routers

import (
	"sftpbackend/controller"
	"sftpbackend/middleware"

	"github.com/gin-gonic/gin"
)

// registerLoginRouter 注册登录相关路由
func registerLoginRouter(r *gin.Engine) {
	// 获取 RSA 公钥（公共接口）
	r.GET("/rsa/public-key", controller.GetPublicKey)
	
	// 登录并生成 token
	r.POST("/login", controller.Login)

	// 修改密码（需鉴权）
	r.POST("/user/change-password", middleware.AuthMiddleware(), controller.ChangePassword)
}
