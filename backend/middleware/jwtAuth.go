package middleware

import (
	"net/http"
	"sftpbackend/jwt"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Token")
		if tokenString == "" {
			c.JSON(http.StatusOK, gin.H{
				"code":    400,
				"message": "未提供授权token",
			})
			c.Abort()
			return
		}

		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		if _, ok := jwt.InvalidTokens.Load(tokenString); ok {
			c.JSON(http.StatusOK, gin.H{
				"code":    50014,
				"message": "该token已失效，请重新登录",
			})
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    50014,
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		// 检查是否为受限Token（仅允许修改密码）
		if len(claims.Routes) == 1 && claims.Routes[0] == "ChangePasswordOnly" {
			// 仅允许访问 /user/change-password 接口
			if c.Request.URL.Path != "/user/change-password" {
				c.JSON(http.StatusOK, gin.H{
					"code":    403,
					"message": "密码已过期，请先修改密码",
				})
				c.Abort()
				return
			}
		}

		c.Set("username", claims.Username)
		c.Set("claims", claims)
		c.Next()
	}
}

// RequireRoute 路由权限中间件：校验当前用户是否拥有指定路由权限
// 用法：group.Use(middleware.AuthMiddleware(), middleware.RequireRoute("System"))
func RequireRoute(routeName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := c.MustGet("claims").(*jwt.CustomClaims)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"code":    400,
				"message": "认证信息缺失",
			})
			c.Abort()
			return
		}

		// 受限Token（密码过期）已在AuthMiddleware拦截，这里兜底
		if len(claims.Routes) == 1 && claims.Routes[0] == "ChangePasswordOnly" {
			c.JSON(http.StatusOK, gin.H{
				"code":    403,
				"message": "密码已过期，请先修改密码",
			})
			c.Abort()
			return
		}

		// 空权限列表：放行（兼容旧Token）
		if len(claims.Routes) == 0 {
			c.Next()
			return
		}

		// 校验路由名是否在用户权限列表中
		for _, r := range claims.Routes {
			if r == routeName {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    403,
			"message": "无权限访问该功能",
		})
		c.Abort()
	}
}