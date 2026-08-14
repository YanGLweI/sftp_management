package models

import (
	"sftpbackend/dao"
	"time"

	"gorm.io/gorm"
)

// LocalUser 平台本地账号
type LocalUser struct {
	gorm.Model
	Username            string     `json:"username" gorm:"column:username;type:varchar(64);uniqueIndex;not null"`
	Password            string     `json:"-" gorm:"column:password;type:varchar(256);not null"`                                   // bcrypt 哈希值，不返回
	MustChangePassword  bool       `json:"mustChangePassword" gorm:"column:must_change_password;default:false"`                       // 登录后需改密（新建/编辑勾选生效，不限首次）
	PasswordNeverExpires bool      `json:"passwordNeverExpires" gorm:"column:password_never_expires;default:false"`                   // 密码永不过期（跳过过期检查）
	PasswordChangedAt   *time.Time `json:"passwordChangedAt" gorm:"column:password_changed_at"`                                    // 用于密码过期判断
	Enabled             bool       `json:"enabled" gorm:"column:enabled;default:true"`                                            // 是否启用
	FailedAttempts      int        `json:"failedAttempts" gorm:"column:failed_attempts;default:0"`                                 // 连续失败登录次数
	RoleID              *uint      `json:"roleId" gorm:"column:role_id"`                                                           // 外键关联 Role
	LastLoginAt         *time.Time `json:"lastLoginAt" gorm:"column:last_login_at"`
}

// GetLocalUserByUsername 通过用户名查询本地账号
func GetLocalUserByUsername(username string) (*LocalUser, error) {
	var user LocalUser
	err := dao.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetLocalUserByID 通过ID查询本地账号
func GetLocalUserByID(id uint) (*LocalUser, error) {
	var user LocalUser
	err := dao.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}