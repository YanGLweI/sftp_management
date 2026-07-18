package middleware

import (
	"net/http"

	"sftpbackend/jwt"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取Token字段，该字段通常包含JWT token（格式一般为 "Bearer <token>"）
		tokenString := c.GetHeader("Token")
		// 如果请求头中未提供Authorization字段，即tokenString为空，说明未提供授权token
		// 返回401 Unauthorized错误响应给前端，并终止请求继续处理
		if tokenString == "" {
			c.JSON(http.StatusOK, gin.H{
				"code":    400,
				"message": "未提供授权token",
			})
			c.Abort()
			return
		}

		// 去除 "Bearer " 前缀（通常JWT的Bearer认证方式有这个前缀），获取真正的token字符串部分
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}
		// 检查该Token是否已被标记为失效
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
		// 将解析后的claims信息保存到请求上下文中，以便后续处理函数使用
		c.Set("username", claims.Username)
		// 继续处理请求
		c.Next()
	}
}
