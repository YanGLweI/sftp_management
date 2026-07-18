package routers

import (
	"sftpbackend/controller"

	"github.com/gin-gonic/gin"
)

func registerDashboardRouter(r *gin.Engine) {
	dashboardGroup := r.Group("/dashboard")
	{
		// ! 看板
		// 获取看板数据
		dashboardGroup.GET("/user/total", controller.GetDashboard)

	}
}
