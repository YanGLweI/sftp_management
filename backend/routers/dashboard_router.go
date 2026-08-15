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
		
		// P0: 核心数据统计 API
		dashboardGroup.GET("/access/total", controller.GetTotalAccessCount)
		dashboardGroup.GET("/access/today", controller.GetTodayAccessCount)
		dashboardGroup.GET("/access/growth", controller.GetAccessGrowth)
		dashboardGroup.GET("/transfer/total", controller.GetTotalTransferCount)
		dashboardGroup.GET("/transfer/today", controller.GetTodayTransferCount)
		dashboardGroup.GET("/transfer/growth", controller.GetTransferGrowth)
		dashboardGroup.GET("/auth/distribution", controller.GetAuthDistribution)
		dashboardGroup.GET("/users/active-top10", controller.GetActiveUsersTop10)

	}
}
