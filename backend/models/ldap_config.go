package models

import (
	"crypto/x509"
	"database/sql/driver"
	"encoding/base64"
	"fmt"
	"gorm.io/gorm"
	"sftpbackend/dao"
	"strings"
	"time"
)

// LDAPConfig LDAP 连接配置表（从配置文件迁移到数据库存储）
type LDAPConfig struct {
	ID          uint   `json:"id" gorm:"primaryKey"`                          // 主键 ID
	Server      string `json:"server" gorm:"size:255;not null;"`              // LDAP 服务器地址（如 ldaps://10.60.254.252:636）
	BaseDN      string `json:"base_dn" gorm:"size:255;not null;"`             // 基础 DN（如 dc=hot,dc=local）
	UseTLS      bool   `json:"use_tls" gorm:"type:tinyint(1);default:0"`      // 是否使用 TLS 连接
	Insecure    bool   `json:"insecure" gorm:"type:tinyint(1);default:1"`     // 是否跳过 TLS 证书验证
	UserFilter  string `json:"user_filter" gorm:"size:255"`
	Username    string `json:"username" gorm:"type:text;encrypt:true"`        // 绑定用户名（加密存储）
	Password    string `json:"password" gorm:"type:text;encrypt:true"`        // 绑定密码（加密存储）
	CertBase64  string `json:"cert_base64" gorm:"type:text;"`                 // CA 证书 Base64 编码
	CertFilename string `json:"cert_filename" gorm:"size:255;"`                // CA 证书文件名
	UpdatedBy   uint   `json:"updated_by" gorm:"index;default:0"`             // 最后更新人 ID
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`            // 创建时间
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`            // 更新时间
}

// TableName 指定表名
func (LDAPConfig) TableName() string {
	return "t_ldap_config"
}

// MigrateLDAPConfig 触发 LDAP 配置表的自动迁移（包含新增字段）
func MigrateLDAPConfig() error {
	return dao.DB.AutoMigrate(&LDAPConfig{})
}

// GetLDAPConfig 获取当前 LDAP 配置（单例模式，只返回第一条记录）
func GetLDAPConfig() (*LDAPConfig, error) {
	var config LDAPConfig
	
	err := dao.DB.Order("id ASC").First(&config).Error
	if err != nil {
		return nil, fmt.Errorf("查询 LDAP 配置失败：%v", err)
	}
	
	return &config, nil
}

// CreateLDAPConfig 创建新的 LDAP 配置记录（单例，首次初始化用）
func CreateLDAPConfig(defaultServer, defaultBaseDN string) error {
	config := &LDAPConfig{
		Server:     defaultServer,
		BaseDN:     defaultBaseDN,
		UseTLS:     false,
		Insecure:   true,
		UserFilter: "(sAMAccountName=%s)",
		UpdatedBy:  0,
	}
	
	return dao.DB.Create(config).Error
}

// SaveLDAPConfig 保存或更新 LDAP 配置（单例更新）
func SaveLDAPConfig(config *LDAPConfig, updatedBy uint) error {
	existingConfig, err := GetLDAPConfig()
	if err == nil && existingConfig.ID > 0 {
		// 更新现有记录
		config.ID = existingConfig.ID
		config.UpdatedBy = updatedBy
		
		return dao.DB.Save(config).Error
	} else if err == gorm.ErrRecordNotFound {
		// 首次创建
		config.UpdatedBy = updatedBy
		
		return dao.DB.Create(config).Error
	}
	
	return fmt.Errorf("保存 LDAP 配置失败：%v", err)
}

// loadCACertFromBase64 解析 CA 证书内容，兼容两种存储格式：
// 1. 明文 PEM 文本（新格式，包含 -----BEGIN CERTIFICATE-----）
// 2. Base64 编码（旧格式，为 PEM/DER 原始字节的 Base64）
func loadCACertFromBase64(base64Cert string) ([]byte, error) {
	if base64Cert == "" {
		return nil, fmt.Errorf("CA 证书不能为空")
	}

	// 新格式：明文 PEM 文本直接使用
	if strings.Contains(base64Cert, "-----BEGIN") {
		return []byte(base64Cert), nil
	}

	// 旧格式：Base64 解码
	certBytes, err := base64.StdEncoding.DecodeString(base64Cert)
	if err != nil {
		return nil, fmt.Errorf("证书 Base64 解码失败：%v", err)
	}

	return certBytes, nil
}

// ParseCACertPool 解析 Base64 编码的 CA 证书，兼容 PEM 与 DER 两种格式
func ParseCACertPool(certBase64 string) (*x509.CertPool, error) {
	certBytes, err := loadCACertFromBase64(certBase64)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	// 优先按 PEM 解析
	if certPool.AppendCertsFromPEM(certBytes) {
		return certPool, nil
	}

	// 回退按 DER 解析
	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, fmt.Errorf("无效的 CA 证书格式（支持 PEM/DER）：%v", err)
	}
	certPool.AddCert(cert)
	return certPool, nil
}

// Value 实现 sql.Scanner 接口（用于加密字段序列化时解密）
func (l *LDAPConfig) Value() (driver.Value, error) {
	// 这里不实现，因为我们通过手动加密/解密处理
	return nil, nil
}

// Scan 实现 sql.Scanner 接口
func (l *LDAPConfig) Scan(value interface{}) error {
	return nil
}
