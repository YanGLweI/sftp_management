package controller

import (
	"net/http"
	"sftpbackend/tools"
	"github.com/gin-gonic/gin"
)

// GetPublicKey 获取 RSA 公钥（公共接口，无需鉴权）
func GetPublicKey(c *gin.Context) {
	publicKey, err := tools.GetPublicKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取公钥失败：" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    publicKey,
		"message": "成功",
	})
}
