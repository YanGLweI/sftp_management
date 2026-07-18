package routers

import (
	"sftpbackend/controller"
	"sftpbackend/middleware"

	"github.com/gin-gonic/gin"
)

// registerSystemRouter 注册系统安全相关路由
func registerSystemRouter(r *gin.Engine) {
	systemGroup := r.Group("/system")

	// 无需鉴权的系统更新接口
	systemGroup.GET("/ws/update", controller.SystemUpdate)

	// 需要鉴权的系统模块接口
	systemGroup.Use(middleware.AuthMiddleware())
	{
		// 获取更新历史 : /update/history?pageNum=1&pageSize=10
		systemGroup.GET("/update/history", controller.GetUpdateHistory)
		// 查看系统更新详情 : /update/detail/:id
		systemGroup.GET("/update/detail/:id", controller.GetUpdateDetail)
		// 获取系统更新计划 : /update/schedule
		systemGroup.GET("/update/schedule", controller.GetUpdateSchedule)
		// 设置系统更新计划 : /update/schedule
		systemGroup.POST("/update/schedule", controller.SetUpdateSchedule)
		// 获取更新报告计划
		systemGroup.GET("/update/report/schedule", controller.GetUpdateReportSchedule)
		// 设置更新报告计划
		systemGroup.POST("/update/report/schedule", controller.SetUpdateReportSchedule)

		// ! 反病毒模块
		// 获取反病毒库信息
		systemGroup.GET("/antivirus", controller.Antivirus)
		// 获取计划任务信息
		systemGroup.GET("/antivirus/schedule", controller.Schedule)
		// 设置计划任务
		systemGroup.POST("/antivirus/schedule", controller.SetSchedule)
		// 启动扫描: /antivirus/scan
		systemGroup.GET("/antivirus/scan", controller.StartScan)
		// 获取隔离区状态
		systemGroup.GET("/antivirus/isolationzone", controller.IsolationZone)
		// 获取卡巴斯基报告计划任务
		systemGroup.GET("/antivirus/report/schedule", controller.GetKReportSchedule)
		// 设置卡巴斯基报告计划任务
		systemGroup.POST("/antivirus/report/schedule", controller.SetKReportSchedule)

		// ! 系统加固模块
		// 获取系统加固检查列表
		systemGroup.GET("/security/checklist", controller.GetSystemSecurityCheck)
		// 立即启动系统加固任务
		systemGroup.GET("/security/start", controller.StartSystemSecurityCheck)
		// 获取系统加固计划 : /security/schedule
		systemGroup.GET("/security/schedule", controller.GetSystemSecuritySchedule)
		// 设置系统加固计划 : /security/schedule
		systemGroup.POST("/security/schedule", controller.SetSystemSecuritySchedule)
		// 获取系统加固报告计划任务
		systemGroup.GET("/security/report/schedule", controller.GetSystemSecurityReportSchedule)
		// 设置系统加固报告计划任务
		systemGroup.POST("/security/report/schedule", controller.SetSystemSecurityReportSchedule)
	}
}
