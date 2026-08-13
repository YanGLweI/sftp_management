package models

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"
	"sftpbackend/config"
	"slices"
	"strings"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/msteinert/pam"
)

// AuthenticateLDAP 通过LDAP验证用户的用户名和密码（平台登录，使用 ldap.security_group_dn 安全组）
// 参数：
//
//	username - 用户名
//	password - 密码
//
// 返回：
//
//	userAttributes - 用户属性  "cn", "WhenChanged", "memberOf"
//	statusCode - 状态码
//	error - 错误信息
func AuthenticateLDAP(username, password string) (map[string][]string, int, error) {
	return AuthenticateLDAPWithGroup(username, password, config.GlobalConfig.LDAP.SecurityGroupDN)
}

/*
* AuthenticateLDAPWithGroup 函数用于通过LDAP验证用户的用户名和密码，并校验其属于指定安全组
* 参数：
* @param username 用户名
* @param password 密码
* @param securityGroupDN 允许的安全组DN（如 ldap.sftp_security_group_dn），后续SFTP模块可通用
* 返回：
* @return userAttributes 用户属性  "cn", "WhenChanged", "memberOf"
* @return statusCode 状态码
* @return error 错误信息
 */
func AuthenticateLDAPWithGroup(username, password, securityGroupDN string) (map[string][]string, int, error) {
	// 从配置文件中获取LDAP信息
	config := config.GlobalConfig.LDAP
	// 创建LDAP连接
	var l *ldap.Conn
	var ldapErr error

	// 创建带超时的Dialer
	dialer := &net.Dialer{
		Timeout: 10 * time.Second, // 设置10秒连接超时
	}

	// 判读是否使用TLS连接
	if config.UseTLS {
		// 加载CA证书
		caCert, err := os.ReadFile(config.CertPath)
		if err != nil {
			return nil, 500, fmt.Errorf("无法读取证书文件: %v", err)
		}

		// 解析 DER 编码的证书
		cert, err := x509.ParseCertificate(caCert)
		if err != nil {
			return nil, 500, fmt.Errorf("无法解析证书: %v", err)
		}

		// 创建证书池并添加CA证书
		certPool := x509.NewCertPool()
		certPool.AddCert(cert)

		l, ldapErr = ldap.DialURL(config.Server,
			ldap.DialWithDialer(dialer),
			ldap.DialWithTLSConfig(&tls.Config{
				InsecureSkipVerify: config.Insecure,
				RootCAs:            certPool,
				MinVersion:         tls.VersionTLS12,
			}))
	} else {
		l, ldapErr = ldap.DialURL(config.Server, ldap.DialWithDialer(dialer))
	}

	if ldapErr != nil {
		return nil, 504, fmt.Errorf("无法连接到LDAP服务器: %v", ldapErr)
	}
	defer l.Close()

	// 添加管理员绑定
	if err := l.Bind(config.Username, config.Password); err != nil {
		return nil, 500, fmt.Errorf("LDAP管理员绑定失败: %v", err)
	}

	// 根据用户名搜索用户DN
	searchRequest := ldap.NewSearchRequest(
		config.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(config.UserFilter, username),
		[]string{"cn", "WhenChanged", "memberOf"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return nil, 500, fmt.Errorf("LDAP搜索失败: %v", err)
	}

	if len(sr.Entries) != 1 {
		return nil, 401, fmt.Errorf("用户名或密码错误")
	}

	userDN := sr.Entries[0].DN

	// 使用用户DN和密码进行绑定
	if err := l.Bind(userDN, password); err != nil {
		return nil, 401, fmt.Errorf("用户名或密码错误")
	}

	// 认证成功，检查用户是否属于指定的安全组
	userGroups := sr.Entries[0].GetAttributeValues("memberOf")
	isMember := slices.Contains(userGroups, securityGroupDN)

	if !isMember {
		return nil, 403, fmt.Errorf("用户不属于授权组，拒绝访问。")
	}

	// 认证成功，返回用户属性
	userAttributes := make(map[string][]string)
	for _, attr := range sr.Entries[0].Attributes {
		userAttributes[attr.Name] = attr.Values
	}

	return userAttributes, 200, nil
}

// 验证本地账号
// func AuthenticateLocal(username, password string) error {
// 	conf := config.GlobalConfig.LocalUser
// 	if username == conf.Username && password == conf.Password {
// 		return nil
// 	}
// 	return fmt.Errorf("本地验证失败，账号或密码错误")
// }

// AuthenticateLocal 验证RHEL9系统本地账号密码
// username: 前端传入的系统用户名
// password: 前端传入的系统密码
// 返回值: nil表示验证通过，非nil表示验证失败
// AuthenticateLocal 验证RHEL9系统本地账号密码
// 参数：
//
//	username - 系统用户名
//	password - 待验证的密码
//
// 返回：验证通过返回nil，失败返回具体错误
func AuthenticateLocal(username, password string) error {
	// 初始化PAM会话（service名推荐用"login"，兼容系统默认认证规则）
	s, err := pam.StartFunc("login", username, func(s pam.Style, msg string) (string, error) {
		// PAM回调：当需要密码时返回传入的密码
		switch s {
		case pam.PromptEchoOff:
			return password, nil
		default:
			return "", fmt.Errorf("不支持的PAM交互类型: %v", s)
		}
	})
	if err != nil {
		return fmt.Errorf("PAM初始化失败: %w", err)
	}

	// 执行PAM认证（flags=0 使用默认规则）
	if err := s.Authenticate(0); err != nil {
		return fmt.Errorf("账号或密码错误/账号锁定: %w", err)
	}

	// 可选：验证账号是否有权限登录（如非空shell、未过期等）
	if err := s.AcctMgmt(0); err != nil {
		return fmt.Errorf("账号无登录权限/过期: %w", err)
	}
	// 读取/etc/passwd解析用户Shell（替代os/user的错误写法）
	shell, err := getUserShellFromPasswd(username)
	if err != nil {
		return fmt.Errorf("读取用户Shell失败: %w", err)
	}

	// 4. 校验Shell是否为禁止登录类型
	forbiddenShells := map[string]bool{
		"/sbin/nologin":     true,
		"/usr/sbin/nologin": true, // RHEL9部分环境的nologin路径
		"/bin/false":        true,
		"":                  true, // 空Shell
		"/bin/nologin":      true, // 兼容其他发行版
		"/usr/bin/false":    true,
	}
	if forbiddenShells[shell] {
		return fmt.Errorf("账号无登录权限(%s)，拒绝登录", shell)
	}

	return nil // 验证通过
}

// getUserShellFromPasswd 从/etc/passwd解析指定用户的登录Shell
// /etc/passwd格式：用户名:密码占位符:UID:GID:注释:家目录:Shell
func getUserShellFromPasswd(username string) (string, error) {
	// 读取/etc/passwd（所有用户均可读取，无需强制root）
	data, err := os.ReadFile("/etc/passwd")
	if err != nil {
		return "", fmt.Errorf("读取/etc/passwd失败: %w", err)
	}

	// 按行解析
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") { // 跳过空行/注释行
			continue
		}

		parts := strings.Split(line, ":")
		// 确保行格式正确（至少7个字段）且用户名匹配
		if len(parts) >= 7 && parts[0] == username {
			fmt.Printf("用户%s的Shell为:%s\n", username, parts[6])
			return parts[6], nil // 第7个字段是Shell
		}

	}

	return "", fmt.Errorf("用户%s不存在于/etc/passwd", username)
}
