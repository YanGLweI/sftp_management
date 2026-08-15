package controller

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"sftpbackend/models"
	"sftpbackend/tools"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-ldap/ldap/v3"
	"github.com/sirupsen/logrus"
)

// LDAPConfigController LDAP 配置管理控制器
type LDAPConfigController struct{}

var LdapConfigController = &LDAPConfigController{}

// GetLDAPConfig 获取当前 LDAP 配置（用户名解密后供表单回显，密码不返回）
func (ctrl *LDAPConfigController) GetLDAPConfig(c *gin.Context) {
	config, err := models.GetLDAPConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取 LDAP 配置失败：" + err.Error(),
		})
		return
	}

	// 解密用户名供表单回显
	username, err := tools.Decrypt(config.Username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "用户名解密失败",
		})
		return
	}

	// 屏蔽 data 中的敏感字段，密码出于安全考虑不返回
	config.Username = ""
	config.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"code":     200,
		"message":  "success",
		"data":     config,
		"username": username,
	})
}

// SaveLDAPConfig 保存 LDAP 配置（接收前加密敏感字段）
func (ctrl *LDAPConfigController) SaveLDAPConfig(c *gin.Context) {
	var req struct {
		Server       string `json:"server" binding:"required"`
		BaseDN       string `json:"base_dn" binding:"required"`
		UseTLS       bool   `json:"use_tls"`
		Insecure     bool   `json:"insecure"`
		UserFilter   string `json:"user_filter" binding:"required"`
		Username     string `json:"username" binding:"required"`
		Password     string `json:"password"`
		CertBase64   string `json:"cert_base64"`
		CertFilename string `json:"cert_filename"` // 证书文件名
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误：" + err.Error(),
		})
		return
	}

	// 加密敏感字段
	encryptedUsername, err := tools.Encrypt(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "用户名加密失败：" + err.Error(),
		})
		return
	}

	encryptedPassword := ""
	if req.Password != "" {
		// 密码由前端 RSA 加密传输（与登录一致），先解密明文再加密存储
		plainPassword, decErr := tools.DecryptPassword(req.Password)
		if decErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "密码解密失败：" + decErr.Error(),
			})
			return
		}
		encryptedPassword, err = tools.Encrypt(plainPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "密码加密失败：" + err.Error(),
			})
			return
		}
	} else {
		// 密码留空表示不修改，保留数据库中已有的加密密码
		existingConfig, e := models.GetLDAPConfig()
		if e == nil && existingConfig.ID > 0 {
			encryptedPassword = existingConfig.Password
			logrus.Printf("LDAP 配置：密码未修改，保留原有值 (ID=%d)", existingConfig.ID)
		} else {
			// 严重错误：无法获取现有配置，拒绝保存
			logrus.Errorf("LDAP 配置：无法获取现有密码 (e=%v)，拒绝保存", e)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "无法验证现有配置，请重试",
			})
			return
		}
	}

	// 构建配置对象
	config := &models.LDAPConfig{
		Server:       req.Server,
		BaseDN:       req.BaseDN,
		UseTLS:       req.UseTLS,
		Insecure:     req.Insecure,
		UserFilter:   req.UserFilter,
		Username:     encryptedUsername,
		Password:     encryptedPassword,
		CertBase64:   req.CertBase64,
		CertFilename: req.CertFilename, // 保存文件名
	}

	// 获取当前用户名（AuthMiddleware 注入），认证与菜单授权已由中间件保证
	username := c.GetString("username")
	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权访问",
		})
		return
	}

	// 审计字段：本地账号取 ID，LDAP 账号等回退 0
	var userID uint
	if u, err := models.GetLocalUserByUsername(username); err == nil {
		userID = u.ID
	}

	// 保存配置
	if err := models.SaveLDAPConfig(config, uint(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "保存 LDAP 配置失败：" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "LDAP 配置保存成功",
	})
}

// TestLDAPConnection 测试 LDAP 连接
func (ctrl *LDAPConfigController) TestLDAPConnection(c *gin.Context) {
	var req struct {
		Server     string `json:"server" binding:"required"`
		BaseDN     string `json:"base_dn" binding:"required"`
		UseTLS     bool   `json:"use_tls"`
		Insecure   bool   `json:"insecure"`
		Username   string `json:"username"`
		Password   string `json:"password"`
		CertBase64 string `json:"cert_base64"`
		UserFilter string `json:"user_filter"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误：" + err.Error(),
		})
		return
	}

	// 密码由前端 RSA 加密传输（与登录一致），先解密明文
	if req.Password != "" {
		plainPassword, err := tools.DecryptPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "密码解密失败：" + err.Error(),
			})
			return
		}
		req.Password = plainPassword
	}

	// 测试连接
	err := testLDAPConnectionInternal(req.Server, req.BaseDN, req.UseTLS, req.Insecure, req.Username, req.Password, req.CertBase64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "连接测试失败：" + err.Error(),
			"data": gin.H{
				"connected": false,
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "LDAP 连接测试成功",
		"data": gin.H{
			"connected": true,
		},
	})
}

// testLDAPConnectionInternal 内部实现 LDAP 连接测试逻辑
func testLDAPConnectionInternal(server, baseDN string, useTLS, insecure bool, username, password, certBase64 string) error {
	dialer := &net.Dialer{Timeout: 10 * time.Second}
	var l *ldap.Conn
	var ldapErr error

	if useTLS {
		if certBase64 == "" {
			return fmt.Errorf("使用 TLS 时需要上传 CA 证书")
		}

		// 解析证书（兼容 PEM/DER 格式），使用缓存避免重复解析
		certPool, err := models.ParseCACertPoolCached(certBase64)
		if err != nil {
			return err
		}

		// 建立 TLS 连接
		l, ldapErr = ldap.DialURL(server,
			ldap.DialWithDialer(dialer),
			ldap.DialWithTLSConfig(&tls.Config{
				InsecureSkipVerify: insecure,
				RootCAs:            certPool,
				MinVersion:         tls.VersionTLS12,
			}))
	} else {
		l, ldapErr = ldap.DialURL(server, ldap.DialWithDialer(dialer))
	}

	if ldapErr != nil {
		return fmt.Errorf("无法连接到 LDAP 服务器：%v", ldapErr)
	}
	defer l.Close()

	// 如果有用户名和密码，则进行绑定测试
	if username != "" && password != "" {
		if err := l.Bind(username, password); err != nil {
			return fmt.Errorf("LDAP 管理员绑定失败：%v", err)
		}
	} else {
		// 如果没有提供凭据，尝试匿名绑定（仅用于测试连接）
		// 如果 LDAP 服务器要求认证，这里会失败并给出提示
	}

	// 搜索测试
	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(objectClass=*)",
		[]string{"cn"},
		nil,
	)

	_, err := l.Search(searchRequest)
	if err != nil {
		return fmt.Errorf("LDAP 搜索失败：%v", err)
	}

	return nil
}
