package models

import (
	"os"
	"path/filepath"
	"sftpbackend/config"
	"sftpbackend/dao"
	"strings"
	"time"
)

// SftpLog SFTP登录与浏览器操作日志
type SftpLog struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time `gorm:"not null" json:"created_at"`
	Username      string    `gorm:"not null" json:"username"` // 操作者：域账号（hotlabel/chinaunicom）或SFTP账号
	Reviewer      string    `gorm:"not null" json:"reviewer"` // 双控复核人（需双控的操作记录，无需双控操作为空）
	IP            string    `gorm:"not null" json:"ip"`
	Action        string    `gorm:"not null" json:"action"` // Login/Logout/Upload/Download/Delete/BatchDelete/Rename/Mkdir
	Message       string    `gorm:"not null" json:"message"`
	Path          string    `gorm:"not null" json:"path"`    // 操作路径
	LogSourceType string    `gorm:"not null;default:normal" json:"log_source_type"` // 日志来源标记：normal（普通登录）/hotlabel（标签上传）/chinaunicom（中国联通）
}

type SftpLogslist struct {
	ID            uint   `json:"id"`
	CreatedAt     string `json:"created_at"`
	Username      string `json:"username"`
	Reviewer      string `json:"reviewer"`
	IP            string `json:"ip"`
	Action        string `json:"action"`
	Message       string `json:"message"`
	Path          string `json:"path"`
	LogSourceType string `json:"log_source_type"`
}

// ! 创建一条SFTP日志
func (l *SftpLog) CreateSftpLog() error {
	return dao.DB.Create(l).Error
}

// ! 获取SFTP日志列表（按时间、用户名模糊过滤，分页，时间倒序）
func (l *SftpLog) GetSftpLogList(page, limit int, date, username string) (logslist []SftpLogslist, totalCount int64, err error) {
	var logs []SftpLog
	// 根据页码和每页数量计算偏移量
	offset := (page - 1) * limit

	// 构建查询条件
	query := dao.DB.Offset(offset).Limit(limit).Order("created_at desc")
	// 先构建一个用于统计总数的查询条件副本，不添加分页相关设置
	countQuery := dao.DB.Model(&SftpLog{})

	// 根据是否传入date和username来动态添加相应的查询条件
	if date != "" && username != "" {
		// 带分页的条件，同时添加时间和用户名条件
		query = query.Where("created_at LIKE? AND username LIKE?", "%"+date+"%", "%"+username+"%")
		// 不带分页的条件，同时添加时间和用户名条件
		countQuery = countQuery.Where("created_at LIKE? AND username LIKE?", "%"+date+"%", "%"+username+"%")
	} else if date != "" {
		// 带分页的条件，仅添加时间条件
		query = query.Where("created_at LIKE?", "%"+date+"%")
		// 不带分页的条件，仅添加时间条件
		countQuery = countQuery.Where("created_at LIKE?", "%"+date+"%")
	} else if username != "" {
		// 带分页的条件，仅添加用户名条件
		query = query.Where("username LIKE?", "%"+username+"%")
		// 不带分页的条件，仅添加用户名条件
		countQuery = countQuery.Where("username LIKE?", "%"+username+"%")
	}
	// 查询满足条件的所有日志
	if err := query.Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	for _, log := range logs {
		formattedTime := log.CreatedAt.Format("2006-01-02 15:04:05")
		logslist = append(logslist, SftpLogslist{
			ID:            log.ID,
			CreatedAt:     formattedTime,
			Username:      log.Username,
			Reviewer:      log.Reviewer,
			IP:            log.IP,
			Action:        log.Action,
			Message:       log.Message,
			Path:          log.Path,
			LogSourceType: log.LogSourceType,
		})
	}
	// 查询满足条件的日志总数
	if err := countQuery.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	return logslist, totalCount, nil
}

// ! 获取中国联通日志列表（按 log_source_type='chinaunicom' 过滤，支持时间、用户名模糊查询，分页，时间倒序）
func (l *SftpLog) GetChinaUnicomLogList(page, limit int, date, username string) (logslist []SftpLogslist, totalCount int64, err error) {
	var logs []SftpLog
	// 根据页码和每页数量计算偏移量
	offset := (page - 1) * limit

	// 构建查询条件（仅中国联通来源日志）
	query := dao.DB.Where("log_source_type = ?", "chinaunicom").Offset(offset).Limit(limit).Order("created_at desc")
	// 先构建一个用于统计总数的查询条件副本，不添加分页相关设置
	countQuery := dao.DB.Model(&SftpLog{}).Where("log_source_type = ?", "chinaunicom")

	// 根据是否传入date和username来动态添加相应的查询条件
	if date != "" && username != "" {
		// 带分页的条件，同时添加时间和用户名条件
		query = query.Where("created_at LIKE? AND username LIKE?", "%"+date+"%", "%"+username+"%")
		// 不带分页的条件，同时添加时间和用户名条件
		countQuery = countQuery.Where("created_at LIKE? AND username LIKE?", "%"+date+"%", "%"+username+"%")
	} else if date != "" {
		// 带分页的条件，仅添加时间条件
		query = query.Where("created_at LIKE?", "%"+date+"%")
		// 不带分页的条件，仅添加时间条件
		countQuery = countQuery.Where("created_at LIKE?", "%"+date+"%")
	} else if username != "" {
		// 带分页的条件，仅添加用户名条件
		query = query.Where("username LIKE?", "%"+username+"%")
		// 不带分页的条件，仅添加用户名条件
		countQuery = countQuery.Where("username LIKE?", "%"+username+"%")
	}
	// 查询满足条件的所有日志
	if err := query.Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	for _, log := range logs {
		formattedTime := log.CreatedAt.Format("2006-01-02 15:04:05")
		logslist = append(logslist, SftpLogslist{
			ID:            log.ID,
			CreatedAt:     formattedTime,
			Username:      log.Username,
			Reviewer:      log.Reviewer,
			IP:            log.IP,
			Action:        log.Action,
			Message:       log.Message,
			Path:          log.Path,
			LogSourceType: log.LogSourceType,
		})
	}
	// 查询满足条件的日志总数
	if err := countQuery.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	return logslist, totalCount, nil
}

// ============ P0: 核心数据统计方法 ============

// GetTotalLoginCount 获取累计登录总次数（访问量总数）
func GetTotalLoginCount() (int64, error) {
	var count int64
	err := dao.DB.Model(&SftpLog{}).
		Where("action = ?", "Login").
		Count(&count).Error
	return count, err
}

// GetTotalTransferCount 获取累计传输总次数（传输数总数，Upload + Download）
func GetTotalTransferCount() (int64, error) {
	var count int64
	err := dao.DB.Model(&SftpLog{}).
		Where("action IN ?", []string{"Upload", "Download"}).
		Count(&count).Error
	return count, err
}

// GetTodayLoginCount 获取今日登录总次数
func GetTodayLoginCount() (int64, error) {
	var count int64
	today := time.Now().Format("2006-01-02")
	err := dao.DB.Model(&SftpLog{}).
		Where("created_at LIKE ?", today+"%").
		Where("action = ?", "Login").
		Count(&count).Error
	return count, err
}

// GetYesterdayLoginCount 获取昨日登录总次数（用于计算增长率）
func GetYesterdayLoginCount() (int64, error) {
	var count int64
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	err := dao.DB.Model(&SftpLog{}).
		Where("created_at LIKE ?", yesterday+"%").
		Where("action = ?", "Login").
		Count(&count).Error
	return count, err
}

// GetTodayTransferCount 获取今日传输总次数（Upload + Download）
func GetTodayTransferCount() (int64, error) {
	var count int64
	today := time.Now().Format("2006-01-02")
	err := dao.DB.Model(&SftpLog{}).
		Where("created_at LIKE ?", today+"%").
		Where("action IN ?", []string{"Upload", "Download"}).
		Count(&count).Error
	return count, err
}

// GetYesterdayTransferCount 获取昨日传输总次数（用于计算增长率）
func GetYesterdayTransferCount() (int64, error) {
	var count int64
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	err := dao.DB.Model(&SftpLog{}).
		Where("created_at LIKE ?", yesterday+"%").
		Where("action IN ?", []string{"Upload", "Download"}).
		Count(&count).Error
	return count, err
}

// GetAuthMethodDistribution 获取认证方式分布（基于 SFTP 日志真实认证记录）
// 返回：map[string]int，key="密码登录"/"密钥登录"
func GetAuthMethodDistribution() (map[string]int, error) {
	result := make(map[string]int)
	passwordCount, keyCount := countAuthMethodsFromLogs()
	result["密码登录"] = passwordCount
	result["密钥登录"] = keyCount
	return result, nil
}

// countAuthMethodsFromLogs 从今日日志与历史每日日志中统计密码/密钥登录次数
// 日志格式：sshd: Accepted password for xxx ... / Accepted publickey for xxx ...
func countAuthMethodsFromLogs() (passwordCount, keyCount int) {
	// 收集日志文件：今日日志 + 历史每日日志
	var files []string
	logFile := config.GlobalConfig.LogFiles.LogFile
	dailyLogFile := config.GlobalConfig.LogFiles.DailyLogFile
	if logFile != "" {
		files = append(files, logFile)
	}
	if dailyLogFile != "" {
		// daily_logfile 格式如 /data/sftplogs/sftp.log-%s，转为 glob 模式匹配历史文件
		dir := filepath.Dir(dailyLogFile)
		pattern := filepath.Join(dir, "sftp.log-*")
		matches, err := filepath.Glob(pattern)
		if err == nil {
			files = append(files, matches...)
		}
	}

	// 统计每份日志中的认证记录
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		content := string(data)
		passwordCount += strings.Count(content, "Accepted password")
		keyCount += strings.Count(content, "Accepted publickey")
	}
	return passwordCount, keyCount
}

// ActiveUserStat 活跃用户统计结构体
type ActiveUserStat struct {
	Username string `json:"username"`
	Count    int64  `json:"count"`
}

// GetActiveUsersTop6 获取活跃用户 Top6（按登录次数排序）
func GetActiveUsersTop6() ([]ActiveUserStat, error) {
	var stats []ActiveUserStat
	err := dao.DB.Model(&SftpLog{}).
		Select("username, COUNT(*) as count").
		Where("action = ?", "Login").
		Group("username").
		Order("COUNT(*) DESC").
		Limit(6).
		Scan(&stats).Error
	return stats, err
}

// GetTopTransferUsers 获取传输量排行 Top6（按 Upload + Download 次数排序）
func GetTopTransferUsers() ([]ActiveUserStat, error) {
	var stats []ActiveUserStat
	err := dao.DB.Model(&SftpLog{}).
		Select("username, COUNT(*) as count").
		Where("action IN ?", []string{"Upload", "Download"}).
		Group("username").
		Order("COUNT(*) DESC").
		Limit(6).
		Scan(&stats).Error
	return stats, err
}
