package controller

import (
	"fmt"
	"net/http"
	"sftpbackend/dao"
	"sftpbackend/models"
	"sftpbackend/tools"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// GetLocalUserList 获取本地账号列表
func GetLocalUserList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	username := c.Query("username")

	var users []models.LocalUser
	var total int64

	query := dao.DB.Model(&models.LocalUser{})
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	query.Count(&total)

	if err := query.Offset((page - 1) * limit).Limit(limit).
		Order("created_at DESC").Find(&users).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "查询账号列表失败"})
		return
	}

	// 补充角色名称
	type UserResp struct {
		models.LocalUser
		RoleName string `json:"roleName"`
	}
	var respList []UserResp
	for _, user := range users {
		roleName := ""
		if user.RoleID != nil {
			role, err := models.GetRoleByID(*user.RoleID)
			if err == nil {
				roleName = role.Name
			}
		}
		// 不返回密码哈希
		user.Password = ""
		respList = append(respList, UserResp{
			LocalUser: user,
			RoleName:  roleName,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"list":  respList,
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// CreateLocalUser 创建本地账号
func CreateLocalUser(c *gin.Context) {
	var req struct {
		Username             string `json:"username" binding:"required"`
		Password             string `json:"password" binding:"required"`
		RoleID               *uint  `json:"roleId"`
		Enabled              *bool  `json:"enabled"`
		MustChangePassword   *bool  `json:"mustChangePassword"`   // 登录后需改密（默认 false）
		PasswordNeverExpires *bool  `json:"passwordNeverExpires"` // 密码永不过期（默认 false）
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "请求数据格式错误"})
		return
	}

	// 检查用户名是否已存在
	if _, err := models.GetLocalUserByUsername(req.Username); err == nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "用户名已存在"})
		return
	}

	// 清理可能残留的软删除同名记录（避免用户名唯一索引冲突，历史数据可能为软删）
	dao.DB.Unscoped().Where("username = ?", req.Username).Delete(&models.LocalUser{})

	// 解密密码
	decryptedPassword, err := tools.DecryptPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "密码解密失败: " + err.Error()})
		return
	}

	// 验证密码策略
	if err := models.ValidatePasswordPolicy(decryptedPassword); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": err.Error()})
		return
	}

	// bcrypt哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(decryptedPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "密码加密失败"})
		return
	}

	now := time.Now()
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	mustChangePassword := false
	if req.MustChangePassword != nil {
		mustChangePassword = *req.MustChangePassword
	}
	passwordNeverExpires := false
	if req.PasswordNeverExpires != nil {
		passwordNeverExpires = *req.PasswordNeverExpires
	}

	user := models.LocalUser{
		Username:             req.Username,
		Password:             string(hashedPassword),
		MustChangePassword:   mustChangePassword,
		PasswordNeverExpires: passwordNeverExpires,
		PasswordChangedAt:    &now,
		Enabled:              enabled,
		RoleID:               req.RoleID,
	}

	if err := dao.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "创建账号失败: " + err.Error()})
		return
	}

	user.Password = "" // 不返回密码
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建账号成功",
		"data":    user,
	})
}

// UpdateLocalUser 更新本地账号
func UpdateLocalUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "无效的账号ID"})
		return
	}

	var req struct {
		RoleID               *uint `json:"roleId"`
		Enabled              *bool `json:"enabled"`
		MustChangePassword   *bool `json:"mustChangePassword"`   // 登录后需改密
		PasswordNeverExpires *bool `json:"passwordNeverExpires"` // 密码永不过期
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "请求数据格式错误"})
		return
	}

	user, err := models.GetLocalUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "账号不存在"})
		return
	}

	updates := map[string]interface{}{}
	if req.RoleID != nil {
		updates["role_id"] = *req.RoleID
	}
	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
		// 如果启用账号，重置失败次数
		if *req.Enabled {
			updates["failed_attempts"] = 0
		}
	}
	if req.MustChangePassword != nil {
		updates["must_change_password"] = *req.MustChangePassword
	}
	if req.PasswordNeverExpires != nil {
		updates["password_never_expires"] = *req.PasswordNeverExpires
	}

	if err := dao.DB.Model(user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "更新账号失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新账号成功"})
}

// ResetLocalUserPassword 管理员重置密码
func ResetLocalUserPassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "无效的账号ID"})
		return
	}

	var req struct {
		Password string `json:"password" binding:"required"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "请求数据格式错误"})
		return
	}

	user, err := models.GetLocalUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "账号不存在"})
		return
	}

	// 解密密码
	decryptedPassword, err := tools.DecryptPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "密码解密失败: " + err.Error()})
		return
	}

	// 验证密码策略
	if err := models.ValidatePasswordPolicy(decryptedPassword); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": err.Error()})
		return
	}

	// bcrypt哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(decryptedPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "密码加密失败"})
		return
	}

	// 保存旧密码到历史
	if err := models.SavePasswordHistory(user.ID, user.Password); err != nil {
		fmt.Println("保存密码历史失败:", err)
	}

	now := time.Now()
	user.Password = string(hashedPassword)
	user.MustChangePassword = true
	user.PasswordChangedAt = &now
	user.FailedAttempts = 0
	user.Enabled = true

	if err := dao.DB.Save(user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "重置密码失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "密码重置成功，下次登录需修改密码"})
}

// DeleteLocalUser 删除本地账号
func DeleteLocalUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "无效的账号ID"})
		return
	}

	user, err := models.GetLocalUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "账号不存在"})
		return
	}

	// 禁止删除默认admin
	if user.Username == "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "不能删除默认管理员账号"})
		return
	}

	// 删除密码历史
	dao.DB.Where("local_user_id = ?", user.ID).Delete(&models.PasswordHistory{})

	// 硬删除账号（释放用户名唯一索引，允许之后重新创建同名账号）
	if err := dao.DB.Unscoped().Delete(user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "删除账号失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除账号成功"})
}