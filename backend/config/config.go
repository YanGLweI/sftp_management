package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

// SystemConfig结构体用于存储系统相关配置信息
type SystemConfig struct {
	Port              int             `yaml:"port"`
	Mode              string          `yaml:"mode"`            // debug or release
	RSAPrivateKeyPath string          `yaml:"rsa-private-key"` // RSA私钥路径（用于解密前端敏感信息）
	ConfigKeyPath     string          `yaml:"config-key"`      // 配置加密私钥路径（用于解密 config.yml 中的 ENC[] 字段）
	RSAPrivateKey     *rsa.PrivateKey `yaml:"-"`               // 解析RSA私钥
}

// DatabaseConfig结构体用于存储数据库相关配置信息
type DatabaseConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	Prefix   string `yaml:"prefix"`
	Charset  string `yaml:"charset"`
	Query    string `yaml:"query"`
}

// EmailConfig结构体用于存储邮件相关配置信息
type EmailConfig struct {
	From     string   `yaml:"from"`
	Password string   `yaml:"password"`
	Host     string   `yaml:"host"`
	Port     int      `yaml:"port"`
	Subject  string   `yaml:"subject"`
	Content  string   `yaml:"content"` // 邮件内容模板
	Body     string   `yaml:"body"`    // 邮件正文
	Tos      []string `yaml:"tos"`     // 接收方邮箱列表
}

// JWTConfig结构体用于存储JWT相关配置信息
type JWTConfig struct {
	Secret string `yaml:"secret"`
	Expire int    `yaml:"expire"`
	Issuer string `yaml:"issuer"`
}

// LDAPConfig结构体用于存储LDAP相关配置信息
type LDAPConfig struct {
	Server              string `yaml:"server"`
	BaseDN              string `yaml:"base_dn"`
	UseTLS              bool   `yaml:"use_tls"`
	Insecure            bool   `yaml:"insecure"`
	UserFilter          string `yaml:"user_filter"`
	Username            string `yaml:"username"`
	Password            string `yaml:"password"`
	SecurityGroupDN     string `yaml:"security_group_dn"`      // 平台登录安全组
	SftpSecurityGroupDN string `yaml:"sftp_security_group_dn"` // SFTP相关模块安全组（标签上传等），后续SFTP模块可通用
	CertPath            string `yaml:"cert_path"`              // 可选，客户端证书路径
}

// 脚本路径
type ScriptPath struct {
	AdduserScript      string `yaml:"adduser_script"`
	UpdateuserScript   string `yaml:"updateuser_script"`
	DeleteuserScript   string `yaml:"deleteuser_script"`
	BatchDeleteScript  string `yaml:"batch_delete_script"`
	SetupKeyScript     string `yaml:"setupkey_script"`
	SystemUpdateScript string `yaml:"system_update_script"`
	// 系统加固检查脚本
	SystemSecurityCheckScript string `yaml:"system_security_check_script"`
	// 修复ssh日志脚本
	FixSSHLogScript string `yaml:"fix_ssh_log_script"`
	// 传输次数统计脚本
	TransCountScript string `yaml:"trans_count_script"`
	// 访问次数统计脚本
	AccessCountScript string `yaml:"access_count_script"`
}

// sftp日志文件路径
type LogFiles struct {
	LogFile      string `yaml:"logfile"`       // 总日志文件路径
	DailyLogFile string `yaml:"daily_logfile"` // 切割的每日日志文件路径
}

// SFTPConfig结构体用于存储SFTP相关配置信息
type SFTPConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type LocalUserConf struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// SftpAccountConfig 公共SFTP服务账号配置（标签上传、中国联通等模块共用，后续模块可复用）
type SftpAccountConfig struct {
	SFTPUsername string `yaml:"sftp_username"` // 专用SFTP账号
	SFTPPassword string `yaml:"sftp_password"` // SFTP账号密码（支持 ENC[] 加密）
}

// HotLabelConfig 标签上传模块配置（仅配置允许访问的根路径，SFTP账号复用公共 sftp_account）
type HotLabelConfig struct {
	RootPath string `yaml:"root_path"` // 标签上传允许访问的根路径
}

// ChinaUnicomConfig 中国联通模块配置（仅配置允许访问的根路径，SFTP账号复用公共 sftp_account）
type ChinaUnicomConfig struct {
	RootPath string `yaml:"root_path"` // 中国联通允许访问的根路径
}

// 计划任务配置
type SchedulerConfig struct {
	KReportTime      string `yaml:"k_report_time"`
	KIsolateZoneTime string `yaml:"k_isolate_zone_time"`
	SystemUpdateTime string `yaml:"system_update_time"`
	// 系统加固检查时间，格式为："分 时 日 月 周"
	SystemSecurityCheckTime string `yaml:"system_security_check_time"`
	// 更新报告发送时间，格式为："分 时 日 月 周"
	UpdateReportTime string `yaml:"update_report_time"`
	// 加固报告发送时间，格式为："分 时 日 月 周"
	HardeningReportTime string `yaml:"hardening_report_time"`
}

// Config结构体用于整体存储配置信息，包含系统和数据库配置
type Config struct {
	System      SystemConfig      `yaml:"system"`
	Database    DatabaseConfig    `yaml:"database"`
	Email       EmailConfig       `yaml:"email"`
	JWT         JWTConfig         `yaml:"jwt"`
	LDAP        LDAPConfig        `yaml:"ldap"`
	Script      ScriptPath        `yaml:"script"`
	LogFiles    LogFiles          `yaml:"logfiles"`
	SFTP        SFTPConfig        `yaml:"sftp"`
	LocalUser   LocalUserConf     `yaml:"localuser"`
	SftpAccount SftpAccountConfig `yaml:"sftp_account"` // 公共SFTP服务账号
	HotLabel    HotLabelConfig    `yaml:"hotlabel"`
	ChinaUnicom ChinaUnicomConfig `yaml:"chinaunicom"`
	Scheduler   SchedulerConfig   `yaml:"scheduler"`
}

var (
	GlobalConfig Config
)

// init()函数是一直特殊的内置函数，程序运行时会自动执行，无需显式调用
func init() {
	// 读取yml配置文件
	configFile, err := os.Open("config.yml")
	if err != nil {
		fmt.Printf("打开文件失败, err: %v \n", err)
		return
	}
	defer configFile.Close()
	// 解析配置文件到结构体
	decoder := yaml.NewDecoder(configFile)
	err = decoder.Decode(&GlobalConfig)
	if err != nil {
		fmt.Printf("解析配置失败,err: %v \n", err)
	}

	// 解析RSA私钥并存储到全局变量中，供后续解密使用
	GlobalConfig.System.RSAPrivateKey, err = parsePrivateKey()
	if err != nil {
		fmt.Printf("解析RSA私钥失败, err: %v \n", err)
		return
	}

	// 解密配置中的加密字段（ENC[...] 格式）
	if err := DecryptConfig(); err != nil {
		fmt.Printf("解密配置字段失败, err: %v \n", err)
		return
	}

}

// ! 从文件中读取并解析RSA私钥(用于解密前端敏感信息)
func parsePrivateKey() (*rsa.PrivateKey, error) {
	// 从配置文件中获取RSA私钥路径
	privateKeyPath := GlobalConfig.System.RSAPrivateKeyPath
	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// decryptConfigField 解密单个配置字段
// 如果值是 ENC[base64] 格式，则使用提供的私钥解密
// 如果不是加密格式，则原样返回
func decryptConfigField(value string, privKey *rsa.PrivateKey) (string, error) {
	if !isEncrypted(value) {
		return value, nil
	}

	if privKey == nil {
		return "", fmt.Errorf("配置字段已加密，但未加载配置解密密钥")
	}

	// 提取 ENC[...] 中的密文
	ciphertext := value[4 : len(value)-1]

	// Base64 解码
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("base64 解码失败: %w", err)
	}

	// RSA-OAEP 解密
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privKey, data, nil)
	if err != nil {
		return "", fmt.Errorf("RSA-OAEP 解密失败: %w", err)
	}

	return string(plaintext), nil
}

// isEncrypted 判断字符串是否为 ENC[...] 加密格式
func isEncrypted(s string) bool {
	return strings.HasPrefix(s, "ENC[") && strings.HasSuffix(s, "]")
}

// loadConfigPrivateKey 加载配置加密专用私钥
// 路径从 config.yml 的 system.config-key 读取
func loadConfigPrivateKey() (*rsa.PrivateKey, error) {
	keyPath := GlobalConfig.System.ConfigKeyPath
	if keyPath == "" {
		// 未配置则跳过
		return nil, nil
	}
	data, err := os.ReadFile(keyPath)
	if err != nil {
		// 文件不存在时不报错，只是不启用配置解密
		return nil, nil
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("无法解码配置私钥 PEM 块")
	}

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("解析配置私钥失败: %w", err)
	}

	return privKey, nil
}

// DecryptConfig 自动解密配置中所有 ENC[...] 格式的字段
// 通过反射遍历 GlobalConfig 所有嵌套结构体的 string 字段
// 新增加密字段无需修改代码，只要值是 ENC[...] 格式即自动解密
func DecryptConfig() error {
	// 加载配置加密私钥
	privKey, err := loadConfigPrivateKey()
	if err != nil {
		return fmt.Errorf("加载配置加密私钥失败: %w", err)
	}

	// 如果没有配置私钥，跳过解密（兼容未加密的配置）
	if privKey == nil {
		return nil
	}

	// 通过反射遍历并解密所有 ENC[...] 字段
	return decryptStructFields(reflect.ValueOf(&GlobalConfig).Elem(), privKey, "")
}

// ! decryptStructFields 递归遍历结构体，解密所有 ENC[...] 格式的 string 字段
//
// 参数说明：
//   - v:        当前要遍历的 reflect.Value（结构体值）
//   - privKey:  RSA 私钥，用于解密 ENC[...] 密文
//   - path:     当前字段的路径前缀，如 "Database"，用于拼接错误信息
//
// 工作原理：
//  1. 拿到结构体后，逐个字段遍历
//  2. 如果字段是 string 且值为 "ENC[xxx]" 格式 → 解密并写回
//  3. 如果字段是嵌套结构体 → 递归进入下一层
//  4. 如果字段是 []string 切片 → 逐个元素检查并解密
func decryptStructFields(v reflect.Value, privKey *rsa.PrivateKey, path string) error {
	// 第一步：如果传入的是指针，先解引用拿到实际值
	// 例如 *Config → Config
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// 第二步：安全检查，确保当前值确实是结构体
	// 如果不是结构体（如 int、string 等基础类型），无需遍历，直接返回
	if v.Kind() != reflect.Struct {
		return nil
	}

	// 第三步：获取结构体的类型信息（用于读取字段名、判断是否导出等）
	t := v.Type()
	// firstErr 记录第一个遇到的解密错误，后续错误不会覆盖它
	// 这样能保证所有字段都尝试解密，同时至少返回一个错误信息
	var firstErr error

	// 第四步：遍历结构体的每一个字段
	for i := 0; i < v.NumField(); i++ {
		// field: 字段的值（reflect.Value），可以通过 SetString 修改
		field := v.Field(i)
		// fieldType: 字段的元信息（名称、类型、tag 等）
		fieldType := t.Field(i)

		// 第五步：跳过未导出字段（即小写字母开头的字段）
		// 未导出字段无法通过反射修改（CanSet() 为 false），直接跳过
		if !fieldType.IsExported() {
			continue
		}

		// 第六步：构建字段的完整路径，用于错误提示
		// 例如：顶层字段 "Password" → "Password"
		//       嵌套字段 → "Database.Password"
		fieldPath := fieldType.Name
		if path != "" {
			fieldPath = path + "." + fieldType.Name
		}

		// 第七步：根据字段的类型分类处理
		switch field.Kind() {
		case reflect.String:
			//! 【字符串字段】这是最常见的加密字段类型
			// 先取出原始值，检查是否为 ENC[...] 格式
			original := field.String()
			if isEncrypted(original) {
				// 是加密格式，调用解密函数
				decrypted, err := decryptConfigField(original, privKey)
				if err != nil {
					// 解密失败：记录第一个错误（不中断，继续处理其他字段）
					if firstErr == nil {
						firstErr = fmt.Errorf("解密 %s 失败: %w", fieldPath, err)
					}
				} else {
					// 解密成功：将明文写回字段（覆盖原来的 ENC[...] 密文）
					field.SetString(decrypted)
				}
			}
			// 如果不是 ENC[...] 格式，说明是普通明文，无需处理

		case reflect.Struct:
			//! 【嵌套结构体】例如 GlobalConfig.Database 是 DatabaseConfig 类型
			// 递归调用自身，进入下一层结构体继续遍历
			// 例如：Config → Database → Password（string）
			if err := decryptStructFields(field, privKey, fieldPath); err != nil && firstErr == nil {
				firstErr = err
			}

		case reflect.Slice:
			//! 【切片字段】例如 []string 类型的邮箱列表
			// 先判断切片元素类型是否为 string
			if field.Type().Elem().Kind() == reflect.String {
				// 逐个检查切片中的每个元素
				for j := 0; j < field.Len(); j++ {
					elem := field.Index(j) // 获取第 j 个元素的 reflect.Value
					// CanSet() 确保元素可修改，isEncrypted() 检查是否为加密格式
					if elem.CanSet() && isEncrypted(elem.String()) {
						decrypted, err := decryptConfigField(elem.String(), privKey)
						if err != nil {
							if firstErr == nil {
								firstErr = fmt.Errorf("解密 %s[%d] 失败: %w", fieldPath, j, err)
							}
						} else {
							elem.SetString(decrypted)
						}
					}
				}
			}
		}
	}

	// 第八步：返回第一个遇到的错误（如果有）
	// 注意：即使某个字段解密失败，其他字段仍然会被尝试解密
	return firstErr
}
