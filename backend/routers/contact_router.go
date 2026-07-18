package routers

import (
	"sftpbackend/controller"
	"sftpbackend/middleware"

	"github.com/gin-gonic/gin"
)

// registerContactRouter 注册通讯录相关路由
func registerContactRouter(r *gin.Engine) {
	contantGroup := r.Group("/contact")
	contantGroup.Use(middleware.AuthMiddleware())
	{
		// 获取所有用户列表: /address/:page/:limit?name=string
		contantGroup.GET("/address/:page/:limit", controller.GetContactList)
		// 添加一个联系人:请求体{"name":"string","email":"string"}
		contantGroup.POST("/address", controller.AddContact)
		// 更新一个联系人:请求体{"id":uint,"name":"string","email":"string"}
		contantGroup.PUT("/address", controller.UpdateContact)
		// 删除一个联系人:根据ID删除 /address/:id
		contantGroup.DELETE("/address/:id", controller.DeleteContact)
		// 批量删除联系人:请求体{"ids":[1,2,3]}
		contantGroup.DELETE("/address", controller.DeleteContacts)
		// 获取所有联系人姓名和邮箱列表
		contantGroup.GET("/options", controller.GetContactoptions)
	}
}
