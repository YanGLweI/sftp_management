package routers

import (
	"sftpbackend/controller"

	"github.com/gin-gonic/gin"
)

// registerLoginRouter 注册登录相关路由
func registerLoginRouter(r *gin.Engine) {
	// 登录并生成token
	r.POST("/login", controller.Login)
}
