package controller

import (
	"fmt"
	"net/http"
	"sftpbackend/jwt"
	"sftpbackend/models"
	"sftpbackend/tools"

	"github.com/gin-gonic/gin"
)

// !LDAP登录
func Login(c *gin.Context) {
	var userInput models.LoginUser
	// 1. 绑定并校验请求参数
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

	// 赋值解密后的密码给userInput
	userInput.Password = decryptedPassword

	// 2. 初始化userInfo为nil（本地登录时保持nil）
	var userInfo map[string][]string

	// 3. 按登录类型处理，增加类型合法性校验
	switch userInput.LoginType {
	case "ldap":
		// LDAP登录：赋值userInfo
		ldapUserInfo, statusCode, err := models.AuthenticateLDAP(userInput.Name, userInput.Password)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    statusCode,
				"message": err.Error(),
			})
			return
		}
		userInfo = ldapUserInfo // 赋值给外层变量，消除未使用告警

	case "local":
		// 本地登录：不赋值userInfo（保持nil）
		if err := models.AuthenticateLocal(userInput.Name, userInput.Password); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    400,
				"message": err.Error(),
			})
			return
		}

	default:
		// 非法登录类型：直接返回错误
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "不支持的登录类型，仅支持ldap/local",
		})
		return
	}

	// 4. 生成JWT Token
	token, err := jwt.GenerateToken(userInput.Name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "生成Token失败",
		})
		return
	}

	// 5. 返回响应（本地登录时userInfo为nil，前端可识别为空）
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"message": "登录成功",
		"data": gin.H{
			"token":    token,
			"userInfo": userInfo, // LDAP登录有值，本地登录为nil
		},
	})

	// 6. 记录操作日志（独立逻辑，不影响主流程）
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

// ! 获取用户信息
func GetUserInfo(c *gin.Context) {
	// 从gin.Context上下文中获取之前中间件token验证后保存的用户信息
	userName := c.MustGet("username")
	// 返回200 OK响应给前端，包含欢迎消息以及当前用户ID信息，展示已通过验证并成功访问受保护路由
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "欢迎访问受保护的路由，当前用户名为 " + fmt.Sprint(userName),
		"data": gin.H{
			// 可以访问的前端路由
			"routes":  []string{"Product", "Trademark", "Attr", "Spu", "Sku", "Acl", "User", "Role", "RoleAuth", "Permission"},
			"buttons": []string{""},
			"roles":   []string{""},
			"name":    userName,
			"avatar":  "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		},
	})
}

// ! 退出登录
func Logout(c *gin.Context) {
	// 标记Token为无效，实现退出登录功能
	tokenString := c.GetHeader("Token")
	// 去除 "Bearer " 前缀（通常JWT的Bearer认证方式有这个前缀），获取真正的token字符串部分
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}
	jwt.MarkTokenExpired(tokenString)
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"message": "退出登录",
	})
	// 创建操作日志
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
