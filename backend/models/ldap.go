package models

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"
	"sftpbackend/config"
	"sftpbackend/dao"
	"slices"
	"time"
	"unicode"

	"github.com/go-ldap/ldap/v3"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthenticateLDAP 通过LDAP验证用户的用户名和密码（平台登录）
// 返回：用户属性, 状态码, 错误信息
func AuthenticateLDAP(username, password string) (map[string][]string, int, error) {
	return AuthenticateLDAPWithGroup(username, password, "")
}

/*
AuthenticateLDAPWithGroup 通过LDAP验证用户，并检查其安全组是否匹配某个角色的安全组配置
参数：
  - username: 用户名
  - password: 密码
  - securityGroupDN: 额外安全组DN（如SFTP模块专用），为空则不额外校验

返回：
  - userAttributes: 用户属性 "cn", "WhenChanged", "memberOf"
  - statusCode: 状态码
  - error: 错误信息
*/
func AuthenticateLDAPWithGroup(username, password, securityGroupDN string) (map[string][]string, int, error) {
	config := config.GlobalConfig.LDAP
	var l *ldap.Conn
	var ldapErr error

	dialer := &net.Dialer{
		Timeout: 10 * time.Second,
	}

	if config.UseTLS {
		caCert, err := os.ReadFile(config.CertPath)
		if err != nil {
			return nil, 500, fmt.Errorf("无法读取证书文件: %v", err)
		}
		cert, err := x509.ParseCertificate(caCert)
		if err != nil {
			return nil, 500, fmt.Errorf("无法解析证书: %v", err)
		}
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

	if err := l.Bind(config.Username, config.Password); err != nil {
		return nil, 500, fmt.Errorf("LDAP管理员绑定失败: %v", err)
	}

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

	if err := l.Bind(userDN, password); err != nil {
		return nil, 401, fmt.Errorf("用户名或密码错误")
	}

	// 认证成功，获取用户的memberOf属性
	userGroups := sr.Entries[0].GetAttributeValues("memberOf")

	// 校验安全组归属
	if securityGroupDN != "" {
		// 额外安全组校验（如SFTP模块）
		if !slices.Contains(userGroups, securityGroupDN) {
			return nil, 403, fmt.Errorf("用户不属于授权组，拒绝访问。")
		}
	} else {
		// 平台登录：检查用户的LDAP安全组是否匹配某个角色的配置
		var roleLinks []RoleLDAPGroup
		if err := dao.DB.Find(&roleLinks).Error; err != nil {
			return nil, 500, fmt.Errorf("查询角色安全组配置失败: %v", err)
		}
		matched := false
		for _, link := range roleLinks {
			if slices.Contains(userGroups, link.GroupDN) {
				matched = true
				break
			}
		}
		if !matched {
			return nil, 403, fmt.Errorf("用户不属于任何授权角色，拒绝访问。")
		}
	}

	userAttributes := make(map[string][]string)
	for _, attr := range sr.Entries[0].Attributes {
		userAttributes[attr.Name] = attr.Values
	}

	return userAttributes, 200, nil
}

// AuthenticateLocal 验证平台本地账号密码
// 参数：
//
//	username - 用户名
//	password - 明文密码
//
// 返回：localUser, passwordExpired(bool), error
func AuthenticateLocal(username, password string) (*LocalUser, bool, error) {
	// 查询本地账号
	user, err := GetLocalUserByUsername(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, fmt.Errorf("账号或密码错误")
		}
		return nil, false, fmt.Errorf("查询账号失败: %w", err)
	}

	// 检查账号是否启用
	if !user.Enabled {
		return nil, false, fmt.Errorf("账号已被禁用")
	}

	// 校验密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		// 密码错误，记录失败次数
		user.FailedAttempts++
		dao.DB.Model(user).Update("failed_attempts", user.FailedAttempts)

		// 检查是否超过最大失败次数
		policy, _ := GetPasswordPolicy()
		if policy != nil && policy.MaxLoginAttempts > 0 && user.FailedAttempts >= policy.MaxLoginAttempts {
			dao.DB.Model(user).Update("enabled", false)
			return nil, false, fmt.Errorf("连续登录失败次数过多，账号已被锁定")
		}
		return nil, false, fmt.Errorf("账号或密码错误")
	}

	// 登录成功，重置失败次数
	user.FailedAttempts = 0
	dao.DB.Model(user).Updates(map[string]interface{}{
		"failed_attempts": 0,
		"last_login_at":   time.Now(),
	})

	// 检查是否需要强制改密
	if user.MustChangePassword {
		return user, false, nil
	}

	// 检查密码是否过期（永不过期账号豁免）
	passwordExpired := false
	if !user.PasswordNeverExpires && user.PasswordChangedAt != nil {
		policy, _ := GetPasswordPolicy()
		if policy != nil && policy.ExpiryDays > 0 {
			expiryTime := user.PasswordChangedAt.Add(time.Duration(policy.ExpiryDays) * 24 * time.Hour)
			if time.Now().After(expiryTime) {
				passwordExpired = true
			}
		}
	}

	return user, passwordExpired, nil
}

// ValidatePasswordPolicy 校验密码是否符合当前密码策略
// 返回：符合返回nil，不符合返回具体错误信息
func ValidatePasswordPolicy(password string) error {
	policy, err := GetPasswordPolicy()
	if err != nil {
		return fmt.Errorf("获取密码策略失败: %w", err)
	}

	if len(password) < policy.MinLength {
		return fmt.Errorf("密码长度不能少于%d位", policy.MinLength)
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		default:
			hasSpecial = true
		}
	}

	if policy.RequireUppercase && !hasUpper {
		return fmt.Errorf("密码必须包含大写字母")
	}
	if policy.RequireLowercase && !hasLower {
		return fmt.Errorf("密码必须包含小写字母")
	}
	if policy.RequireDigit && !hasDigit {
		return fmt.Errorf("密码必须包含数字")
	}
	if policy.RequireSpecialChar && !hasSpecial {
		return fmt.Errorf("密码必须包含特殊字符")
	}

	return nil
}

// CheckPasswordHistory 检查新密码是否在历史密码中
// 返回：true表示在历史记录中（禁止使用），false表示可用
func CheckPasswordHistory(localUserID uint, newPassword string) (bool, error) {
	policy, err := GetPasswordPolicy()
	if err != nil {
		return false, fmt.Errorf("获取密码策略失败: %w", err)
	}

	if policy.PasswordHistory <= 0 {
		return false, nil
	}

	var histories []PasswordHistory
	if err := dao.DB.Where("local_user_id = ?", localUserID).
		Order("created_at DESC").
		Limit(policy.PasswordHistory).
		Find(&histories).Error; err != nil {
		return false, fmt.Errorf("查询密码历史失败: %w", err)
	}

	for _, h := range histories {
		if err := bcrypt.CompareHashAndPassword([]byte(h.Password), []byte(newPassword)); err == nil {
			return true, nil // 匹配到历史密码
		}
	}
	return false, nil
}

// SavePasswordHistory 保存密码到历史记录（保留最近N条，超出的删除）
func SavePasswordHistory(localUserID uint, hashedPassword string) error {
	policy, err := GetPasswordPolicy()
	if err != nil {
		return err
	}

	// 插入新记录
	history := PasswordHistory{
		LocalUserID: localUserID,
		Password:    hashedPassword,
	}
	if err := dao.DB.Create(&history).Error; err != nil {
		return err
	}

	// 删除超出限制的旧记录
	if policy.PasswordHistory > 0 {
		var count int64
		dao.DB.Model(&PasswordHistory{}).Where("local_user_id = ?", localUserID).Count(&count)
		if count > int64(policy.PasswordHistory) {
			// 获取需要保留的最小ID
			var ids []uint
			dao.DB.Model(&PasswordHistory{}).
				Where("local_user_id = ?", localUserID).
				Order("created_at DESC").
				Limit(policy.PasswordHistory).
				Pluck("id", &ids)
			if len(ids) > 0 {
				dao.DB.Where("local_user_id = ? AND id NOT IN ?", localUserID, ids).
					Delete(&PasswordHistory{})
			}
		}
	}
	return nil
}