package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// SystemConfig结构体用于存储系统相关配置信息
type SystemConfig struct {
	Port              int             `yaml:"port"`
	Mode              string          `yaml:"mode"`            // debug or release
	RSAPrivateKeyPath string          `yaml:"rsa-private-key"` // RSA私钥路径
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
	Server          string `yaml:"server"`
	BaseDN          string `yaml:"base_dn"`
	UseTLS          bool   `yaml:"use_tls"`
	Insecure        bool   `yaml:"insecure"`
	UserFilter      string `yaml:"user_filter"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	SecurityGroupDN string `yaml:"security_group_dn"` // 可选，LDAP安全组DN
	CertPath        string `yaml:"cert_path"`         // 可选，客户端证书路径
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
	System    SystemConfig    `yaml:"system"`
	Database  DatabaseConfig  `yaml:"database"`
	Email     EmailConfig     `yaml:"email"`
	JWT       JWTConfig       `yaml:"jwt"`
	LDAP      LDAPConfig      `yaml:"ldap"`
	Script    ScriptPath      `yaml:"script"`
	LogFiles  LogFiles        `yaml:"logfiles"`
	SFTP      SFTPConfig      `yaml:"sftp"`
	LocalUser LocalUserConf   `yaml:"localuser"`
	Scheduler SchedulerConfig `yaml:"scheduler"`
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

}

// ! 从文件中读取并解析RSA私钥
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
