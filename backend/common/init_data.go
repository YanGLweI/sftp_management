package common

import (
	"errors"
	"log"
	"sftpbackend/config"
	"sftpbackend/dao"
	"sftpbackend/models"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func InitData() {
	// ========== 初始化LDAP配置默认数据（单例，表为空时创建） ==========
	var ldapConfigCount int64
	dao.DB.Model(&models.LDAPConfig{}).Count(&ldapConfigCount)
	if ldapConfigCount == 0 {
		if err := models.CreateLDAPConfig("", ""); err != nil {
			logrus.Printf("InitData Create LDAPConfig failed: %v", err)
		} else {
			logrus.Println("InitData Create LDAPConfig success")
		}
	}

	// ========== 新增Scheduler调度器数据初始化逻辑 ==========
	newScheduler := make([]*models.Scheduler, 0)
	scheduler := []*models.Scheduler{
		{
			// 检查卡巴斯基隔离区
			TaskName: "CheckKasperskyIsolateZone",
			Cron:     config.GlobalConfig.Scheduler.KIsolateZoneTime,
			Enable:   true,
		},
		{
			// 发送卡巴斯基安全报告
			TaskName: "SendKasperskyReport",
			Cron:     config.GlobalConfig.Scheduler.KReportTime,
			Enable:   true,
		},
		{
			// 系统更新
			TaskName: "SystemUpdate",
			Cron:     config.GlobalConfig.Scheduler.SystemUpdateTime,
			Enable:   true,
		},
		{
			// 系统加固检查
			TaskName: "SystemSecurityCheck",
			Cron:     config.GlobalConfig.Scheduler.SystemSecurityCheckTime,
			Enable:   true,
		},
		{
			// 更新报告发送
			TaskName: "SendUpdateReport",
			Cron:     config.GlobalConfig.Scheduler.UpdateReportTime,
			Enable:   true,
		},
		{
			// 加固报告发送
			TaskName: "SendHardeningReport",
			Cron:     config.GlobalConfig.Scheduler.HardeningReportTime,
			Enable:   true,
		},
	}

	for _, s := range scheduler {
		err := dao.DB.Where("task_name = ?", s.TaskName).First(&models.Scheduler{}).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newScheduler = append(newScheduler, s)
		}
	}
	if len(newScheduler) > 0 {
		if err := dao.DB.Create(&newScheduler).Error; err != nil {
			log.Printf("InitData Create Scheduler failed: %v", err)
		}
	}

	// ========== 初始化PasswordPolicy默认密码策略 ==========
	var policyCount int64
	dao.DB.Model(&models.PasswordPolicy{}).Count(&policyCount)
	if policyCount == 0 {
		defaultPolicy := &models.PasswordPolicy{
			MinLength:          14,
			RequireUppercase:   true,
			RequireLowercase:   true,
			RequireDigit:       true,
			RequireSpecialChar: true,
			ExpiryDays:         90,
			PasswordHistory:    5,
			MaxLoginAttempts:   5,
		}
		if err := dao.DB.Create(defaultPolicy).Error; err != nil {
			logrus.Printf("InitData Create PasswordPolicy failed: %v", err)
		} else {
			logrus.Println("InitData Create PasswordPolicy success")
		}
	}

	// ========== 初始化默认角色和本地管理员账号 ==========
	var localUserCount int64
	dao.DB.Model(&models.LocalUser{}).Count(&localUserCount)
	if localUserCount > 0 {
		// 检查是否存在默认 admin 账号
		var admin models.LocalUser
		err := dao.DB.Model(&models.LocalUser{}).Where("username = ?", "admin").First(&admin).Error
		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("InitData: 检测到缺少默认 admin 账号，将重新创建")
		} else if err == nil {
			return // admin 账号已存在，直接返回
		}
	}
	
	// 获取或创建超级管理员角色
	var superRole models.Role
	err := dao.DB.Where("name = ?", "超级管理员").First(&superRole).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果不存在则创建
		superRole = models.Role{
			Name:        "超级管理员",
			Description: "拥有系统所有权限",
		}
		if err := dao.DB.Create(&superRole).Error; err != nil {
			logrus.Printf("InitData Create super role failed: %v", err)
			return
		}
		logrus.Println("InitData Create super role success (new)")
	} else if err == nil {
		// 已存在，使用现有的角色 ID
		logrus.Println("InitData Use existing super role")
	} else {
		logrus.Printf("InitData query super role error: %v", err)
		return
	}
	
	// 定义需要创建的菜单权限
	requiredMenus := []struct {
		RouteName string
		MenuTitle string
	}{
		{RouteName: "Dashboard", MenuTitle: "首页"},
		{RouteName: "Sftp", MenuTitle: "传输管理"},
		{RouteName: "SftpUser", MenuTitle: "账号管理"},
		{RouteName: "Contacts", MenuTitle: "通讯邮箱"},
		{RouteName: "Log", MenuTitle: "日志管理"},
		{RouteName: "PlatformLog", MenuTitle: "平台日志"},
		{RouteName: "SftpLog", MenuTitle: "SFTP 日志"},
		{RouteName: "System", MenuTitle: "系统安全"},
		{RouteName: "SystemUpdate", MenuTitle: "系统更新"},
		{RouteName: "Antivirus", MenuTitle: "病毒管理"},
		{RouteName: "SystemHardening", MenuTitle: "系统加固"},
		{RouteName: "Settings", MenuTitle: "平台设置"},
		{RouteName: "RoleManagement", MenuTitle: "角色管理"},
		{RouteName: "LocalUserManagement", MenuTitle: "本地账号"},
		{RouteName: "PasswordPolicy", MenuTitle: "密码策略"},
		{RouteName: "SftpModuleManagement", MenuTitle: "SFTP 管理"},
		{RouteName: "HotLabelConfig", MenuTitle: "标签上传配置"},
		{RouteName: "ChinaUnicomConfig", MenuTitle: "中国联通配置"},
		{RouteName: "LDAPManagement", MenuTitle: "LDAP 管理"},
	}
	
	// 添加超级管理员角色的菜单权限（如果尚未存在）
	for _, menu := range requiredMenus {
		var count int64
		dao.DB.Model(&models.RoleMenu{}).Where("role_id = ? AND route_name = ?", superRole.ID, menu.RouteName).Count(&count)
		if count == 0 {
			if err := dao.DB.Create(&models.RoleMenu{
				RoleID:    superRole.ID,
				RouteName: menu.RouteName,
				MenuTitle: menu.MenuTitle,
			}).Error; err != nil {
				logrus.Printf("InitData Create role menu %s failed: %v", menu.RouteName, err)
			} else {
				logrus.Printf("InitData Add role menu: %s", menu.MenuTitle)
			}
		}
	}
	
	// 创建默认管理员 admin
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin1234567890."), bcrypt.DefaultCost)
	if err != nil {
		logrus.Printf("InitData hash password failed: %v", err)
		return
	}
	now := time.Now()
	admin := &models.LocalUser{
		Username:           "admin",
		Password:           string(hashedPassword),
		MustChangePassword: true,
		PasswordChangedAt:  &now,
		Enabled:            true,
		RoleID:             &superRole.ID,
	}
	if err := dao.DB.Create(admin).Error; err != nil {
		logrus.Printf("InitData Create admin user failed: %v", err)
	} else {
		logrus.Println("InitData Create admin user success")
	}

	// ========== 新增SystemSecurityStandard标准数据初始化逻辑 ==========
	standardSecurity := &models.SystemSecurityStandard{
		// DNF/Repo配置
		DnfConfGpgcheck:    "1",
		RedhatRepoGpgcheck: "1",
		// 密码策略
		PASSMAXDAYS: "30",
		PASSMINDAYS: "1",
		PASSMINLEN:  "14",
		PASSWARNAGE: "7",
		INACTIVE:    "30",
		GID:         "0",
		TMOUT:       "180",
		// Cron/At配置
		Cron:        "enabled",
		Crontab:     "Access:(0600/-rw-------) Uid：( 0/ root) Gid：( 0/ root)",
		CronHourly:  "Access:(0700/drwx------) Uid：( 0/ root) Gid：( 0/ root)",
		CronDaily:   "Access:(0700/drwx------) Uid：( 0/ root) Gid：( 0/ root)",
		CronWeekly:  "Access:(0700/drwx------) Uid：( 0/ root) Gid：( 0/ root)",
		CronMonthly: "Access:(0700/drwx------) Uid：( 0/ root) Gid：( 0/ root)",
		CronDeny:    "Access:(0600/-rw-------) Uid：( 0/ root) Gid：( 0/ root)",
		AtDeny:      "Access:(0600/-rw-------) Uid：( 0/ root) Gid：( 0/ root)",
		CronAllow:   "No such file or directory",
		AtAllow:     "No such file or directory",
		// SSHD配置
		SshdConfig:              "Access:(0600/-rw-------) Uid：( 0/ root) Gid：( 0/ root)",
		LogLevel:                "INFO",
		X11Forwarding:           "no",
		MaxAuthTries:            "4",
		IgnoreRhosts:            "yes",
		HostbasedAuthentication: "no",
		PermitRootLogin:         "no",
		PermitEmptyPasswords:    "no",
		PermitUserEnvironment:   "no",
		ClientAliveInterval:     "60",
		ClientAliveCountMax:     "3",
		LoginGraceTime:          "60",
		// 密码复杂度
		Minlen:           "14",
		Minclass:         "4",
		Dcredit:          "-1",
		Ucredit:          "-1",
		Lcredit:          "-1",
		Ocredit:          "-1",
		PasswordRemember: "24",
		// 系统文件内容
		Passwd:      "Access:(0644/-rw-r--r--) Uid：( 0/ root) Gid：( 0/ root)",
		PasswdDash:  "Access:(0644/-rw-r--r--) Uid：( 0/ root) Gid：( 0/ root)",
		Group:       "Access:(0644/-rw-r--r--) Uid：( 0/ root) Gid：( 0/ root)",
		GroupDash:   "Access:(0644/-rw-r--r--) Uid：( 0/ root) Gid：( 0/ root)",
		Shadow:      "Access:(0000/----------) Uid：( 0/ root) Gid：( 0/ root)",
		ShadowDash:  "Access:(0000/----------) Uid：( 0/ root) Gid：( 0/ root)",
		Gshadow:     "Access:(0000/----------) Uid：( 0/ root) Gid：( 0/ root)",
		GshadowDash: "Access:(0000/----------) Uid：( 0/ root) Gid：( 0/ root)",
		// 其他配置
		CryptoPolicies: "DEFAULT:NO-SHA1:NO-WEAKMAC:NO-SSHCBC:NO-SSHCHACHA20:NO-SSHETM:NO-SSHWEAKCIPHERS:NO-SSHWEAKMACS",
		NtpServer:      "server 10.60.254.252 iburst",
	}

	// 如果表是空的就插入
	var count int64
	dao.DB.Model(&models.SystemSecurityStandard{}).Count(&count)
	if count == 0 {
		// 不存在则创建
		if err := dao.DB.Create(standardSecurity).Error; err != nil {
			log.Printf("InitData Create SystemSecurityStandard failed: %v", err)
		} else {
			log.Println("InitData Create SystemSecurityStandard success")
		}
	}
}

// EnsureSuperAdminLDAPManagementMenu 确保超级管理员角色拥有新增的 LDAP 管理菜单权限
// 兼容已有部署：旧版本的超级管理员角色菜单列表中不包含新菜单，此处补充
func EnsureSuperAdminLDAPManagementMenu() {
	role, err := models.GetRoleByName("超级管理员")
	if err != nil {
		logrus.Printf("EnsureSuperAdminLDAPManagementMenu: 超级管理员角色不存在：%v", err)
		return
	}

	// 需要确保存在的菜单
	requiredMenus := []struct {
		RouteName string
		MenuTitle string
	}{
		{RouteName: "LDAPManagement", MenuTitle: "LDAP 管理"},
	}

	// 补充缺失的菜单（直接查询 RoleMenu 表判断是否存在，避免 Preload 缺失导致重复插入）
	for _, m := range requiredMenus {
		var count int64
		dao.DB.Model(&models.RoleMenu{}).Where("role_id = ? AND route_name = ?", role.ID, m.RouteName).Count(&count)
		if count == 0 {
			if err := dao.DB.Create(&models.RoleMenu{
				RoleID:    role.ID,
				RouteName: m.RouteName,
				MenuTitle: m.MenuTitle,
			}).Error; err != nil {
				logrus.Printf("EnsureSuperAdminLDAPManagementMenu: 添加菜单 %s 失败：%v", m.RouteName, err)
			} else {
				logrus.Printf("EnsureSuperAdminLDAPManagementMenu: 为超级管理员添加菜单 %s", m.RouteName)
			}
		}
	}
}

// EnsureSuperAdminSftpModuleMenus 确保超级管理员角色拥有新增的 SFTP 模块管理菜单权限
// 兼容已有部署：旧版本的超级管理员角色菜单列表中不包含新菜单，此处补充
func EnsureSuperAdminSftpModuleMenus() {
	role, err := models.GetRoleByName("超级管理员")
	if err != nil {
		logrus.Printf("EnsureSuperAdminSftpModuleMenus: 超级管理员角色不存在: %v", err)
		return
	}

	// 需要确保存在的菜单
	requiredMenus := []struct {
		RouteName string
		MenuTitle string
	}{
		{RouteName: "SftpModuleManagement", MenuTitle: "SFTP 管理"},
		{RouteName: "HotLabelConfig", MenuTitle: "标签上传配置"},
		{RouteName: "ChinaUnicomConfig", MenuTitle: "中国联通配置"},
	}

	// 补充缺失的菜单（直接查询 RoleMenu 表判断是否存在，避免 Preload 缺失导致重复插入）
	for _, m := range requiredMenus {
		var count int64
		dao.DB.Model(&models.RoleMenu{}).Where("role_id = ? AND route_name = ?", role.ID, m.RouteName).Count(&count)
		if count == 0 {
			if err := dao.DB.Create(&models.RoleMenu{
				RoleID:    role.ID,
				RouteName: m.RouteName,
				MenuTitle: m.MenuTitle,
			}).Error; err != nil {
				logrus.Printf("EnsureSuperAdminSftpModuleMenus: 添加菜单 %s 失败: %v", m.RouteName, err)
			} else {
				logrus.Printf("EnsureSuperAdminSftpModuleMenus: 为超级管理员添加菜单 %s", m.RouteName)
			}
		}
	}
}
