package models

import (
	"sftpbackend/dao"
	"sftpbackend/tools"
	"time"

	"github.com/gin-gonic/gin"
)

type SystemSecurity struct {
	// 主键（建议保留自增主键，不纳入"前2个业务字段"范围）
	ID uint64 `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement;comment:自增主键" json:"id"`

	// 前2个核心字段（保留原有类型）
	Date     time.Time `gorm:"column:Date;type:datetime;not null;comment:检查日期" json:"date"`
	Hostname string    `gorm:"column:Hostname;type:varchar(128);not null;comment:主机名" json:"hostname"`

	// 以下所有字段统一为 string 类型，数据库类型 varchar(255)
	Operasystem string `gorm:"column:operasystem;type:varchar(255);comment:操作系统版本" json:"operasystem"`
	Kernel      string `gorm:"column:Kernel;type:varchar(255);comment:内核版本" json:"kernel"`
	IP          string `gorm:"column:IP;type:varchar(255);comment:主机IP地址" json:"ip"`

	// DNF/Repo 配置
	DnfConfGpgcheck    string `gorm:"column:dnf_conf_gpgcheck;type:varchar(255);comment:dnf.conf中gpgcheck配置值" json:"dnf_conf_gpgcheck"`
	RedhatRepoGpgcheck string `gorm:"column:redhat_repo_gpgcheck;type:varchar(255);comment:redhat.repo中gpgcheck配置值" json:"redhat_repo_gpgcheck"`

	// 密码策略（原int类型改为string）
	PASSMAXDAYS string `gorm:"column:PASS_MAX_DAYS;type:varchar(255);comment:密码最大有效期（天）" json:"pass_max_days"`
	PASSMINDAYS string `gorm:"column:PASS_MIN_DAYS;type:varchar(255);comment:密码最小修改间隔（天）" json:"pass_min_days"`
	PASSMINLEN  string `gorm:"column:PASS_MIN_LEN;type:varchar(255);comment:密码最小长度" json:"pass_min_len"`
	PASSWARNAGE string `gorm:"column:PASS_WARN_AGE;type:varchar(255);comment:密码过期警告天数" json:"pass_warn_age"`
	INACTIVE    string `gorm:"column:INACTIVE;type:varchar(255);comment:账户非活动锁定天数" json:"inactive"`
	GID         string `gorm:"column:GID;type:varchar(255);comment:默认GID" json:"gid"`
	TMOUT       string `gorm:"column:TMOUT;type:varchar(255);comment:终端自动超时时间（秒）" json:"tmout"`

	// Cron/At 任务配置
	Cron        string `gorm:"column:Cron;type:varchar(255);comment:cron主配置文件内容" json:"cron"`
	Crontab     string `gorm:"column:crontab;type:varchar(255);comment:crontab文件内容" json:"crontab"`
	CronHourly  string `gorm:"column:cron_hourly;type:varchar(255);comment:cron.hourly目录下任务" json:"cron_hourly"`
	CronDaily   string `gorm:"column:cron_daily;type:varchar(255);comment:cron.daily目录下任务" json:"cron_daily"`
	CronWeekly  string `gorm:"column:cron_weekly;type:varchar(255);comment:cron.weekly目录下任务" json:"cron_weekly"`
	CronMonthly string `gorm:"column:cron_monthly;type:varchar(255);comment:cron.monthly目录下任务" json:"cron_monthly"`
	CronDeny    string `gorm:"column:cron_deny;type:varchar(255);comment:cron.deny文件内容" json:"cron_deny"`
	AtDeny      string `gorm:"column:at_deny;type:varchar(255);comment:at.deny文件内容" json:"at_deny"`
	CronAllow   string `gorm:"column:cron_allow;type:varchar(255);comment:cron.allow文件内容" json:"cron_allow"`
	AtAllow     string `gorm:"column:at_allow;type:varchar(255);comment:at.allow文件内容" json:"at_allow"`

	// SSHD 配置（原int/string统一改为string）
	SshdConfig              string `gorm:"column:sshd_config;type:varchar(255);comment:sshd_config完整配置内容" json:"sshd_config"`
	LogLevel                string `gorm:"column:LogLevel;type:varchar(255);comment:SSH日志级别" json:"log_level"`
	X11Forwarding           string `gorm:"column:X11Forwarding;type:varchar(255);comment:X11转发开关（yes/no）" json:"x11_forwarding"`
	MaxAuthTries            string `gorm:"column:MaxAuthTries;type:varchar(255);comment:SSH最大认证尝试次数" json:"max_auth_tries"`
	IgnoreRhosts            string `gorm:"column:IgnoreRhosts;type:varchar(255);comment:是否忽略rhosts（yes/no）" json:"ignore_rhosts"`
	HostbasedAuthentication string `gorm:"column:HostbasedAuthentication;type:varchar(255);comment:基于主机的认证开关" json:"hostbased_authentication"`
	PermitRootLogin         string `gorm:"column:PermitRootLogin;type:varchar(255);comment:是否允许root登录SSH" json:"permit_root_login"`
	PermitEmptyPasswords    string `gorm:"column:PermitEmptyPasswords;type:varchar(255);comment:是否允许空密码登录" json:"permit_empty_passwords"`
	PermitUserEnvironment   string `gorm:"column:PermitUserEnvironment;type:varchar(255);comment:是否允许用户自定义环境" json:"permit_user_environment"`
	ClientAliveInterval     string `gorm:"column:ClientAliveInterval;type:varchar(255);comment:SSH客户端存活间隔（秒）" json:"client_alive_interval"`
	ClientAliveCountMax     string `gorm:"column:ClientAliveCountMax;type:varchar(255);comment:SSH客户端存活最大次数" json:"client_alive_count_max"`
	LoginGraceTime          string `gorm:"column:LoginGraceTime;type:varchar(255);comment:SSH登录宽限时间（秒）" json:"login_grace_time"`

	// 密码复杂度（原int改为string）
	Minlen           string `gorm:"column:minlen;type:varchar(255);comment:密码最小长度（pam）" json:"minlen"`
	Minclass         string `gorm:"column:minclass;type:varchar(255);comment:密码字符类别最小数" json:"minclass"`
	Dcredit          string `gorm:"column:dcredit;type:varchar(255);comment:数字字符信用值" json:"dcredit"`
	Ucredit          string `gorm:"column:ucredit;type:varchar(255);comment:大写字符信用值" json:"ucredit"`
	Lcredit          string `gorm:"column:lcredit;type:varchar(255);comment:小写字符信用值" json:"lcredit"`
	Ocredit          string `gorm:"column:ocredit;type:varchar(255);comment:特殊字符信用值" json:"ocredit"`
	PasswordRemember string `gorm:"column:password_remember;type:varchar(255);comment:密码记忆次数/天数" json:"password_remember"`

	// 系统文件内容
	Passwd      string `gorm:"column:passwd;type:varchar(255);comment:passwd文件内容" json:"passwd"`
	PasswdDash  string `gorm:"column:passwd-;type:varchar(255);comment:passwd-备份文件内容" json:"passwd_dash"`
	Group       string `gorm:"column:group;type:varchar(255);comment:group文件内容" json:"group"`
	GroupDash   string `gorm:"column:group-;type:varchar(255);comment:group-备份文件内容" json:"group_dash"`
	Shadow      string `gorm:"column:shadow;type:varchar(255);comment:shadow文件内容" json:"shadow"`
	ShadowDash  string `gorm:"column:shadow-;type:varchar(255);comment:shadow-备份文件内容" json:"shadow_dash"`
	Gshadow     string `gorm:"column:gshadow;type:varchar(255);comment:gshadow文件内容" json:"gshadow"`
	GshadowDash string `gorm:"column:gshadow-;type:varchar(255);comment:gshadow-备份文件内容" json:"gshadow_dash"`

	// 其他配置
	CryptoPolicies string `gorm:"column:crypto_policies;type:varchar(255);comment:系统加密策略" json:"crypto_policies"`
	NtpServer      string `gorm:"column:ntp_server;type:varchar(255);comment:NTP服务器地址" json:"ntp_server"`

	// 结果
	Result string `gorm:"column:result;type:varchar(255);comment:检查结果" json:"result"`
}

type SystemSecurityStandard struct {
	// 主键（建议保留自增主键，不纳入"前2个业务字段"范围）
	ID uint64 `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement;comment:自增主键" json:"id"`

	// DNF/Repo 配置
	DnfConfGpgcheck    string `gorm:"column:dnf_conf_gpgcheck;type:varchar(255);comment:dnf.conf中gpgcheck配置值" json:"dnf_conf_gpgcheck"`
	RedhatRepoGpgcheck string `gorm:"column:redhat_repo_gpgcheck;type:varchar(255);comment:redhat.repo中gpgcheck配置值" json:"redhat_repo_gpgcheck"`

	// 密码策略（原int类型改为string）
	PASSMAXDAYS string `gorm:"column:PASS_MAX_DAYS;type:varchar(255);comment:密码最大有效期（天）" json:"pass_max_days"`
	PASSMINDAYS string `gorm:"column:PASS_MIN_DAYS;type:varchar(255);comment:密码最小修改间隔（天）" json:"pass_min_days"`
	PASSMINLEN  string `gorm:"column:PASS_MIN_LEN;type:varchar(255);comment:密码最小长度" json:"pass_min_len"`
	PASSWARNAGE string `gorm:"column:PASS_WARN_AGE;type:varchar(255);comment:密码过期警告天数" json:"pass_warn_age"`
	INACTIVE    string `gorm:"column:INACTIVE;type:varchar(255);comment:账户非活动锁定天数" json:"inactive"`
	GID         string `gorm:"column:GID;type:varchar(255);comment:默认GID" json:"gid"`
	TMOUT       string `gorm:"column:TMOUT;type:varchar(255);comment:终端自动超时时间（秒）" json:"tmout"`

	// Cron/At 任务配置
	Cron        string `gorm:"column:Cron;type:varchar(255);comment:cron主配置文件内容" json:"cron"`
	Crontab     string `gorm:"column:crontab;type:varchar(255);comment:crontab文件内容" json:"crontab"`
	CronHourly  string `gorm:"column:cron_hourly;type:varchar(255);comment:cron.hourly目录下任务" json:"cron_hourly"`
	CronDaily   string `gorm:"column:cron_daily;type:varchar(255);comment:cron.daily目录下任务" json:"cron_daily"`
	CronWeekly  string `gorm:"column:cron_weekly;type:varchar(255);comment:cron.weekly目录下任务" json:"cron_weekly"`
	CronMonthly string `gorm:"column:cron_monthly;type:varchar(255);comment:cron.monthly目录下任务" json:"cron_monthly"`
	CronDeny    string `gorm:"column:cron_deny;type:varchar(255);comment:cron.deny文件内容" json:"cron_deny"`
	AtDeny      string `gorm:"column:at_deny;type:varchar(255);comment:at.deny文件内容" json:"at_deny"`
	CronAllow   string `gorm:"column:cron_allow;type:varchar(255);comment:cron.allow文件内容" json:"cron_allow"`
	AtAllow     string `gorm:"column:at_allow;type:varchar(255);comment:at.allow文件内容" json:"at_allow"`

	// SSHD 配置（原int/string统一改为string）
	SshdConfig              string `gorm:"column:sshd_config;type:varchar(255);comment:sshd_config完整配置内容" json:"sshd_config"`
	LogLevel                string `gorm:"column:LogLevel;type:varchar(255);comment:SSH日志级别" json:"log_level"`
	X11Forwarding           string `gorm:"column:X11Forwarding;type:varchar(255);comment:X11转发开关（yes/no）" json:"x11_forwarding"`
	MaxAuthTries            string `gorm:"column:MaxAuthTries;type:varchar(255);comment:SSH最大认证尝试次数" json:"max_auth_tries"`
	IgnoreRhosts            string `gorm:"column:IgnoreRhosts;type:varchar(255);comment:是否忽略rhosts（yes/no）" json:"ignore_rhosts"`
	HostbasedAuthentication string `gorm:"column:HostbasedAuthentication;type:varchar(255);comment:基于主机的认证开关" json:"hostbased_authentication"`
	PermitRootLogin         string `gorm:"column:PermitRootLogin;type:varchar(255);comment:是否允许root登录SSH" json:"permit_root_login"`
	PermitEmptyPasswords    string `gorm:"column:PermitEmptyPasswords;type:varchar(255);comment:是否允许空密码登录" json:"permit_empty_passwords"`
	PermitUserEnvironment   string `gorm:"column:PermitUserEnvironment;type:varchar(255);comment:是否允许用户自定义环境" json:"permit_user_environment"`
	ClientAliveInterval     string `gorm:"column:ClientAliveInterval;type:varchar(255);comment:SSH客户端存活间隔（秒）" json:"client_alive_interval"`
	ClientAliveCountMax     string `gorm:"column:ClientAliveCountMax;type:varchar(255);comment:SSH客户端存活最大次数" json:"client_alive_count_max"`
	LoginGraceTime          string `gorm:"column:LoginGraceTime;type:varchar(255);comment:SSH登录宽限时间（秒）" json:"login_grace_time"`

	// 密码复杂度（原int改为string）
	Minlen           string `gorm:"column:minlen;type:varchar(255);comment:密码最小长度（pam）" json:"minlen"`
	Minclass         string `gorm:"column:minclass;type:varchar(255);comment:密码字符类别最小数" json:"minclass"`
	Dcredit          string `gorm:"column:dcredit;type:varchar(255);comment:数字字符信用值" json:"dcredit"`
	Ucredit          string `gorm:"column:ucredit;type:varchar(255);comment:大写字符信用值" json:"ucredit"`
	Lcredit          string `gorm:"column:lcredit;type:varchar(255);comment:小写字符信用值" json:"lcredit"`
	Ocredit          string `gorm:"column:ocredit;type:varchar(255);comment:特殊字符信用值" json:"ocredit"`
	PasswordRemember string `gorm:"column:password_remember;type:varchar(255);comment:密码记忆次数/天数" json:"password_remember"`

	// 系统文件内容
	Passwd      string `gorm:"column:passwd;type:varchar(255);comment:passwd文件内容" json:"passwd"`
	PasswdDash  string `gorm:"column:passwd-;type:varchar(255);comment:passwd-备份文件内容" json:"passwd_dash"`
	Group       string `gorm:"column:group;type:varchar(255);comment:group文件内容" json:"group"`
	GroupDash   string `gorm:"column:group-;type:varchar(255);comment:group-备份文件内容" json:"group_dash"`
	Shadow      string `gorm:"column:shadow;type:varchar(255);comment:shadow文件内容" json:"shadow"`
	ShadowDash  string `gorm:"column:shadow-;type:varchar(255);comment:shadow-备份文件内容" json:"shadow_dash"`
	Gshadow     string `gorm:"column:gshadow;type:varchar(255);comment:gshadow文件内容" json:"gshadow"`
	GshadowDash string `gorm:"column:gshadow-;type:varchar(255);comment:gshadow-备份文件内容" json:"gshadow_dash"`

	// 其他配置
	CryptoPolicies string `gorm:"column:crypto_policies;type:varchar(255);comment:系统加密策略" json:"crypto_policies"`
	NtpServer      string `gorm:"column:ntp_server;type:varchar(255);comment:NTP服务器地址" json:"ntp_server"`
}

// 查询最新的系统加固结果
func (s *SystemSecurity) FindLatest() (*SystemSecurity, error) {
	var latest SystemSecurity
	if err := dao.DB.Order("id desc").First(&latest).Error; err != nil {
		return nil, err
	}
	return &latest, nil
}

// 查询标准值
func (s *SystemSecurityStandard) FindStandard() (*SystemSecurityStandard, error) {
	var standard SystemSecurityStandard
	if err := dao.DB.Order("id desc").First(&standard).Error; err != nil {
		return nil, err
	}
	return &standard, nil
}

// 更新实际结果
func (s *SystemSecurity) Update() error {
	return dao.DB.Save(s).Error
}

// 分页获取系统加固检查列表
func (s *SystemSecurity) GetSystemSecurityCheck(c *gin.Context) ([]SystemSecurity, int64, error) {
	var pageOption tools.PageOption
	if err := c.ShouldBindQuery(&pageOption); err != nil {
		return nil, 0, err
	}

	page := tools.NewPageOption(pageOption.PageNum, pageOption.PageSize)

	var checklist []SystemSecurity
	var total int64

	// 分页查询
	query := dao.DB.Model(&SystemSecurity{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("id desc").
		Limit(page.PageSize).
		Offset(page.PageNum).
		Find(&checklist).Error; err != nil {
		return nil, 0, err
	}

	return checklist, total, nil
}
