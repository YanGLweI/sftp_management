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
		dashboardGroup.GET("/users/active-top6", controller.GetActiveUsersTop6)   // 登录数排行（活跃用户）
		dashboardGroup.GET("/users/top-transfers", controller.GetTopTransferUsers)  // 传输量排行（SFTP传输量排行）

	}
}
