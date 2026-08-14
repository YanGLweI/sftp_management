package controller

import (
	"net/http"
	"sftpbackend/dao"
	"sftpbackend/models"
	"sftpbackend/tools"

	"github.com/gin-gonic/gin"
)

// GetPasswordPolicy 获取当前密码策略
func GetPasswordPolicy(c *gin.Context) {
	policy, err := models.GetPasswordPolicy()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "获取密码策略失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    policy,
	})
}

// UpdatePasswordPolicy 更新密码策略
func UpdatePasswordPolicy(c *gin.Context) {
	var req struct {
		MinLength          int  `json:"minLength"`
		RequireUppercase   bool `json:"requireUppercase"`
		RequireLowercase   bool `json:"requireLowercase"`
		RequireDigit       bool `json:"requireDigit"`
		RequireSpecialChar bool `json:"requireSpecialChar"`
		ExpiryDays         int  `json:"expiryDays"`
		PasswordHistory    int  `json:"passwordHistory"`
		MaxLoginAttempts   int  `json:"maxLoginAttempts"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "请求数据格式错误"})
		return
	}

	policy, err := models.GetPasswordPolicy()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "获取密码策略失败"})
		return
	}

	policy.MinLength = req.MinLength
	policy.RequireUppercase = req.RequireUppercase
	policy.RequireLowercase = req.RequireLowercase
	policy.RequireDigit = req.RequireDigit
	policy.RequireSpecialChar = req.RequireSpecialChar
	policy.ExpiryDays = req.ExpiryDays
	policy.PasswordHistory = req.PasswordHistory
	policy.MaxLoginAttempts = req.MaxLoginAttempts

	if err := dao.DB.Save(policy).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "更新密码策略失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新密码策略成功"})
}

// ValidatePassword 验证密码是否符合当前策略（供前端调用）
func ValidatePassword(c *gin.Context) {
	var req struct {
		Password string `json:"password" binding:"required"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "请求数据格式错误"})
		return
	}

	// 解密密码
	decryptedPassword, err := tools.DecryptPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "密码解密失败: " + err.Error()})
		return
	}

	if err := models.ValidatePasswordPolicy(decryptedPassword); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": err.Error(),
			"data": gin.H{
				"valid": false,
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "密码符合策略",
		"data": gin.H{
			"valid": true,
		},
	})
}