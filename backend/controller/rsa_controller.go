package controller

import (
	"log"
	"net/http"
	"sftpbackend/tools"
	"github.com/gin-gonic/gin"
)

// GetPublicKey 获取 RSA 公钥（公共接口，无需鉴权）
func GetPublicKey(c *gin.Context) {
	publicKey, err := tools.GetPublicKey()
	if err != nil {
		// 详细错误记录到 server logs
		log.Printf("获取公钥失败：%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "公钥服务暂时不可用，请稍后重试",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    publicKey,
		"message": "成功",
	})
}
