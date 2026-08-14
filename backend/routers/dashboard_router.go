package routers

import (
	"sftpbackend/controller"
	"sftpbackend/middleware"

	"github.com/gin-gonic/gin"
)

func registerDashboardRouter(r *gin.Engine) {
	dashboardGroup := r.Group("/dashboard")
	dashboardGroup.Use(middleware.AuthMiddleware(), middleware.RequireRoute("Dashboard"))
	{
		// ! 看板
		// 获取看板数据
		dashboardGroup.GET("/user/total", controller.GetDashboard)

	}
}
