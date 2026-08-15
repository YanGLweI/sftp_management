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

	// 参数边界校验，防止无效配置导致安全问题
	if req.MinLength < 8 || req.MinLength > 128 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "密码最小长度应在 8-128 位之间"})
		return
	}
	if req.ExpiryDays < 0 || req.ExpiryDays > 3650 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "密码过期天数应在 0-3650 天之间（0 表示永不过期）"})
		return
	}
	if req.PasswordHistory < 0 || req.PasswordHistory > 50 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "密码历史记录数应在 0-50 条之间（0 表示不限制）"})
		return
	}
	if req.MaxLoginAttempts < 1 || req.MaxLoginAttempts > 20 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "最大登录失败次数应在 1-20 次之间"})
		return
	}
	// 复杂度规则全关闭时拒绝（防止出现无约束的弱密码策略）
	if !req.RequireUppercase && !req.RequireLowercase && !req.RequireDigit && !req.RequireSpecialChar {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "至少需要启用一种复杂度要求"})
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