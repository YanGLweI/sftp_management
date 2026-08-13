package middleware

import (
	"net/http"
	"sftpbackend/utils"

	"github.com/gin-gonic/gin"
)

// DualAuthMiddleware 双控验证中间件
// 仅对 loginType == "chinaunicom" 的连接强制要求携带有效的双控验证凭证（X-Dual-Token），
// 其他登录方式（密码/密钥/标签上传）直接放行
func DualAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取SFTP连接
		token := c.GetHeader("X-SFTP-Token")
		if token == "" {
			token = c.Query("sftp_token")
		}
		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "SFTP Token不能为空",
			})
			c.Abort()
			return
		}

		conn, err := utils.SFTPConnManager.GetConn(token)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    50014,
				"message": "SFTP连接失效: " + err.Error(),
			})
			c.Abort()
			return
		}

		// 2. 非中国联通连接无需双控验证
		if conn.LoginType != "chinaunicom" {
			c.Next()
			return
		}

		// 3. 校验双控凭证（X-Dual-Token，60秒有效，绑定当前连接）
		dualToken := c.GetHeader("X-Dual-Token")
		if !utils.DualAuthManager.VerifyToken(token, dualToken) {
			c.JSON(http.StatusOK, gin.H{
				"code":    428,
				"message": "需要双控验证",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
