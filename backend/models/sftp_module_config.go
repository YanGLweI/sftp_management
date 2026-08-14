package models

import (
	"sftpbackend/dao"

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

// InitDefaultConfigs 初始化默认配置
func InitDefaultConfigs() error {
	hotlabelConfig := &SFTPModuleConfig{
		ModuleName:      ModuleNameHotLabel,
		LoginType:       LoginTypeLDAP,
		EnabledRoles:    "[]", // 初始为空，待管理员配置
		DualAuthEnabled: false,
	}
	
	chinaUnicomConfig := &SFTPModuleConfig{
		ModuleName:      ModuleNameChinaUnicom,
		LoginType:       LoginTypeLDAP,
		EnabledRoles:    "[]",
		DualAuthEnabled: true,
	}

	// 检查是否已存在配置
	existingHotLabel, _ := GetSFTPModuleConfig(ModuleNameHotLabel)
	if existingHotLabel == nil {
		// 使用 Select 强制写入 DualAuthEnabled=false（避免 GORM 对 bool 零值使用数据库默认值）
		if err := dao.DB.Select("ModuleName", "LoginType", "EnabledRoles", "DualAuthEnabled").Create(hotlabelConfig).Error; err != nil {
			return err
		}
	}

	existingChinaUnicom, _ := GetSFTPModuleConfig(ModuleNameChinaUnicom)
	if existingChinaUnicom == nil {
		if err := CreateSFTPModuleConfig(chinaUnicomConfig); err != nil {
			return err
		}
	}

	return nil
}
