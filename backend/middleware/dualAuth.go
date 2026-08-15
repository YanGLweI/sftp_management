package middleware

import (
	"net/http"
	"sftpbackend/models"
	"sftpbackend/utils"

	"github.com/gin-gonic/gin"
)

// DualAuthMiddleware 双控验证中间件
// 根据模块配置决定是否要求双控验证（t_sftp_module_config.dual_auth_enabled），
// 未配置时默认中国联通模块启用双控，其他模块关闭双控
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

		// 2. 读取模块配置决定是否启用双控
		requireDualAuth := false
		moduleConfig, cfgErr := models.GetSFTPModuleConfig(conn.LoginType)
		if cfgErr == nil && moduleConfig != nil {
			// 配置存在：以配置为准
			requireDualAuth = moduleConfig.DualAuthEnabled
		} else {
			// 配置不存在时兼容默认行为：仅中国联通默认启用双控
			requireDualAuth = conn.LoginType == "chinaunicom"
		}

		// 3. 无需双控验证时直接放行
		if !requireDualAuth {
			c.Next()
			return
		}

		// 4. 校验双控凭证（X-Dual-Token，60 秒有效，绑定当前连接与签发 IP）
		dualToken := c.GetHeader("X-Dual-Token")
		if !utils.DualAuthManager.VerifyToken(token, dualToken, c.ClientIP()) {
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
