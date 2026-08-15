package models

import (
	"sftpbackend/dao"
	"strings"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// SFTPModuleConfig SFTP 模块配置模型
type SFTPModuleConfig struct {
	gorm.Model
	ModuleName      string `json:"moduleName" gorm:"column:module_name;type:varchar(64);uniqueIndex;not null"` // hotlabel, chinaunicom
	LoginType       string `json:"loginType" gorm:"column:login_type;type:varchar(32);not null"`              // local, ldap
	EnabledRoles    string `json:"enabledRoles" gorm:"column:enabled_roles;type:json;not null"`               // JSON 数组，存储角色 ID 列表
	DualAuthEnabled bool   `json:"dualAuthEnabled" gorm:"column:dual_auth_enabled;type:tinyint"`              // 是否启用双控
}

const (
	// ModuleNameHotLabel 标签上传模块
	ModuleNameHotLabel = "hotlabel"
	// ModuleNameChinaUnicom 中国联通模块
	ModuleNameChinaUnicom = "chinaunicom"
)

const (
	// LoginTypeLocal 本地登录
	LoginTypeLocal = "local"
	// LoginTypeLDAP LDAP 登录
	LoginTypeLDAP = "ldap"
)

// GetSFTPModuleConfig 根据模块名称获取配置
func GetSFTPModuleConfig(moduleName string) (*SFTPModuleConfig, error) {
	var config SFTPModuleConfig
	err := dao.DB.Where("module_name = ?", moduleName).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// GetAllSFTPModuleConfigs 获取所有配置
func GetAllSFTPModuleConfigs() ([]SFTPModuleConfig, error) {
	var configs []SFTPModuleConfig
	err := dao.DB.Find(&configs).Error
	if err != nil {
		return nil, err
	}
	return configs, nil
}

// UpdateSFTPModuleConfig 更新配置
func UpdateSFTPModuleConfig(config *SFTPModuleConfig) error {
	result := dao.DB.Save(config)
	return result.Error
}

// CreateSFTPModuleConfig 创建配置（用于初始化）
func CreateSFTPModuleConfig(config *SFTPModuleConfig) error {
	result := dao.DB.Create(config)
	return result.Error
}

// InitDefaultConfigs 初始化默认配置（支持并发）
func InitDefaultConfigs() error {
	configs := []SFTPModuleConfig{
		{ModuleName: ModuleNameHotLabel, LoginType: LoginTypeLDAP, EnabledRoles: "[]", DualAuthEnabled: false},
		{ModuleName: ModuleNameChinaUnicom, LoginType: LoginTypeLDAP, EnabledRoles: "[]", DualAuthEnabled: true},
	}
	
	for i := range configs {
		var existing SFTPModuleConfig
		err := dao.DB.Where("module_name = ?", configs[i].ModuleName).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			// 尝试创建
			if err := dao.DB.Create(&configs[i]).Error; err != nil {
				// 捕获唯一索引冲突（多进程同时创建）
				if strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "Duplicate key") {
					logrus.Warnf("模块 %s 已被其他实例创建 (忽略)", configs[i].ModuleName)
					continue
				}
				return err
			}
			logrus.Printf("初始化模块配置成功：%s", configs[i].ModuleName)
		} else if err != nil {
			return err
		} else {
			logrus.Printf("模块配置已存在：%s (ID=%d)", existing.ModuleName, existing.ID)
		}
	}
	return nil
}
