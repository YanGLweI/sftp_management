package routers

import (
	"sftpbackend/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Student struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func SetupRouter() *gin.Engine {
	// 设置gin模式
	if config.GlobalConfig.System.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// ! 跨域配置
	// config := cors.DefaultConfig()
	// config.AllowOrigins = []string{"http://localhost:8080"}
	// r.Use(cors.Default()) // 允许来自任何源请求
	// 只允许config配置的指定域名跨域
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://10.60.1.191:9528", "http://localhost:9528", "https://localhost", "https://hot-sftp.it.local", "*"}, // 允许前端域名
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 注册所有子路由
	registerDashboardRouter(r) // 看板路由
	registerLoginRouter(r)     // 登录路由
	registerUserRouter(r)      // 用户模块路由
	registerSFTPRouter(r)      // SFTP模块路由
	registerContactRouter(r)   // 通讯录模块路由
	registerSystemRouter(r)    // 系统安全模块路由
	registerSettingsRouter(r)  // 平台设置模块路由

	return r
}
