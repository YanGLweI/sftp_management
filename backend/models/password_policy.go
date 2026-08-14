package models

import (
	"sftpbackend/dao"

	"gorm.io/gorm"
)

// PasswordHistory 记录用户历史密码哈希值，防止密码复用
type PasswordHistory struct {
	gorm.Model
	LocalUserID uint   `gorm:"column:local_user_id;index;not null"` // 关联 local_users.id
	Password    string `gorm:"column:password;type:varchar(256);not null"` // 旧密码的 bcrypt 哈希
}

// PasswordPolicy 密码策略（系统级全局配置，仅一条记录）
type PasswordPolicy struct {
	gorm.Model
	// 密码复杂度
	MinLength          int  `json:"minLength" gorm:"column:min_length;default:14"`
	RequireUppercase   bool `json:"requireUppercase" gorm:"column:require_uppercase;default:true"`
	RequireLowercase   bool `json:"requireLowercase" gorm:"column:require_lowercase;default:true"`
	RequireDigit       bool `json:"requireDigit" gorm:"column:require_digit;default:true"`
	RequireSpecialChar bool `json:"requireSpecialChar" gorm:"column:require_special_char;default:true"`
	// 密码过期
	ExpiryDays int `json:"expiryDays" gorm:"column:expiry_days;default:90"`
	// 历史与锁定
	PasswordHistory  int `json:"passwordHistory" gorm:"column:password_history;default:5"`
	MaxLoginAttempts int `json:"maxLoginAttempts" gorm:"column:max_login_attempts;default:5"`
}

// GetPasswordPolicy 获取当前密码策略
func GetPasswordPolicy() (*PasswordPolicy, error) {
	var policy PasswordPolicy
	err := dao.DB.First(&policy).Error
	if err != nil {
		return nil, err
	}
	return &policy, nil
}