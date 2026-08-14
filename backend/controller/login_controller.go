package controller

import (
	"fmt"
	"net/http"
	"sftpbackend/dao"
	"sftpbackend/jwt"
	"sftpbackend/models"
	"sftpbackend/tools"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// getRoleRoutes 获取角色ID对应的可访问路由名称列表
func getRoleRoutes(roleID *uint) []string {
	if roleID == nil {
		return []string{}
	}
	role, err := models.GetRoleByID(*roleID)
	if err != nil {
		return []string{}
	}
	routes := make([]string, 0, len(role.Menus))
	for _, menu := range role.Menus {
		routes = append(routes, menu.RouteName)
	}
	return routes
}

// Login 用户登录
func Login(c *gin.Context) {
	var userInput models.LoginUser
	if err := c.BindJSON(&userInput); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "请求数据格式错误",
		})
		return
	}

	// 解密密码
	decryptedPassword, err := tools.DecryptPassword(userInput.Password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	userInput.Password = decryptedPassword

	var userInfo map[string][]string
	var roleID *uint
	var routes []string
	var mustChangePassword bool
	var passwordExpired bool

	switch userInput.LoginType {
	case "ldap":
		ldapUserInfo, statusCode, err := models.AuthenticateLDAP(userInput.Name, userInput.Password)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    statusCode,
				"message": err.Error(),
			})
			return
		}
		userInfo = ldapUserInfo

		// 查询LDAP用户的安全组匹配的角色
		userGroups := ldapUserInfo["memberOf"]
		// 查找角色
		var roleLinks []models.RoleLDAPGroup
		if err := dao.DB.Find(&roleLinks).Error; err == nil {
			for _, link := range roleLinks {
				for _, g := range userGroups {
					if link.GroupDN == g {
						roleID = &link.RoleID
						break
					}
				}
				if roleID != nil {
					break
				}
			}
		}
		routes = getRoleRoutes(roleID)

	case "local":
		localUser, expired, err := models.AuthenticateLocal(userInput.Name, userInput.Password)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    400,
				"message": err.Error(),
			})
			return
		}
		mustChangePassword = localUser.MustChangePassword
		passwordExpired = expired
		roleID = localUser.RoleID
		routes = getRoleRoutes(roleID)

	default:
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "不支持的登录类型，仅支持ldap/local",
		})
		return
	}

	// 生成JWT Token
	claims := &jwt.CustomClaims{
		Username:  userInput.Name,
		LoginType: userInput.LoginType,
		RoleID:    roleID,
		Routes:    routes,
	}
	token, err := jwt.GenerateTokenWithClaims(claims)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "生成Token失败",
		})
		return
	}

	// 返回响应
	response := gin.H{
		"code":    20000,
		"message": "登录成功",
		"data": gin.H{
			"token":    token,
			"userInfo": userInfo,
			"routes":   routes,
		},
	}

	// 本地账号特殊标记
	if userInput.LoginType == "local" {
		if mustChangePassword {
			response["data"] = gin.H{
				"token":                token,
				"must_change_password": true,
				"routes":               routes,
			}
		} else if passwordExpired {
			// 密码过期，生成受限Token
			limitedToken, err := jwt.GenerateLimitedToken(userInput.Name)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code":    500,
					"message": "生成Token失败",
				})
				return
			}
			response["data"] = gin.H{
				"token":            limitedToken,
				"password_expired": true,
				"routes":           routes,
			}
		}
	}

	c.JSON(http.StatusOK, response)

	// 记录操作日志
	log := models.Log{
		Username: userInput.Name,
		IP:       c.ClientIP(),
		Action:   "Login",
		Message:  "用户登录: " + userInput.Name,
	}
	if err := log.CreateLog(); err != nil {
		fmt.Println("日志创建失败:", err)
	}
}

// GetUserInfo 获取用户信息
func GetUserInfo(c *gin.Context) {
	userName := c.MustGet("username").(string)
	claims := c.MustGet("claims").(*jwt.CustomClaims)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "欢迎访问受保护的路由，当前用户名为 " + userName,
		"data": gin.H{
			"routes":  claims.Routes,
			"buttons": []string{""},
			"roles":   []string{""},
			"name":    userName,
			"avatar":  "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		},
	})
}

// ChangePassword 修改密码（本地账号）
func ChangePassword(c *gin.Context) {
	userName := c.MustGet("username").(string)
	claims := c.MustGet("claims").(*jwt.CustomClaims)

	// 仅本地账号可修改密码
	if claims.LoginType != "local" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "仅本地账号支持修改密码",
		})
		return
	}

	var req struct {
		OldPassword string `json:"oldPassword" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "请求数据格式错误",
		})
		return
	}

	// 解密密码
	decryptedOld, err := tools.DecryptPassword(req.OldPassword)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "旧密码解密失败: " + err.Error(),
		})
		return
	}
	decryptedNew, err := tools.DecryptPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "新密码解密失败: " + err.Error(),
		})
		return
	}

	// 查询用户
	user, err := models.GetLocalUserByUsername(userName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "用户不存在",
		})
		return
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(decryptedOld)); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "旧密码错误",
		})
		return
	}

	// 验证新密码符合密码策略
	if err := models.ValidatePasswordPolicy(decryptedNew); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	// 检查历史密码
	reused, err := models.CheckPasswordHistory(user.ID, decryptedNew)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "检查密码历史失败: " + err.Error(),
		})
		return
	}
	if reused {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "新密码不能与最近使用过的密码相同",
		})
		return
	}

	// 哈希新密码
	hashedNew, err := bcrypt.GenerateFromPassword([]byte(decryptedNew), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "密码加密失败",
		})
		return
	}

	// 保存旧密码到历史记录
	if err := models.SavePasswordHistory(user.ID, user.Password); err != nil {
		fmt.Println("保存密码历史失败:", err)
	}

	// 更新密码
	now := time.Now()
	user.Password = string(hashedNew)
	user.MustChangePassword = false
	user.PasswordChangedAt = &now
	if err := dao.DB.Save(user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "密码更新失败",
		})
		return
	}

	// 生成新的完整Token
	claims.Routes = getRoleRoutes(user.RoleID)
	newToken, err := jwt.GenerateTokenWithClaims(claims)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "修改密码成功，但生成新Token失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"message": "密码修改成功",
		"data": gin.H{
			"token": newToken,
		},
	})

	// 记录操作日志
	log := models.Log{
		Username: userName,
		IP:       c.ClientIP(),
		Action:   "ChangePassword",
		Message:  "用户修改密码: " + userName,
	}
	if err := log.CreateLog(); err != nil {
		fmt.Println("日志创建失败:", err)
	}
}

// Logout 退出登录
func Logout(c *gin.Context) {
	tokenString := c.GetHeader("Token")
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}
	jwt.MarkTokenExpired(tokenString)
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"message": "退出登录",
	})
	log := models.Log{
		Username: c.MustGet("username").(string),
		IP:       c.ClientIP(),
		Action:   "Logout",
		Message:  "用户注销: " + c.MustGet("username").(string),
	}
	if err := log.CreateLog(); err != nil {
		fmt.Println("日志创建失败:", err)
	}
}