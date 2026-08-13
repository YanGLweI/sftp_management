package models

import (
	"sftpbackend/dao"
	"time"
)

// SftpLog SFTP登录与浏览器操作日志
type SftpLog struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	Username  string    `gorm:"not null" json:"username"` // 操作者：域账号（hotlabel/chinaunicom）或SFTP账号
	Reviewer  string    `gorm:"not null" json:"reviewer"` // 双控复核人（需双控的操作记录，无需双控操作为空）
	IP        string    `gorm:"not null" json:"ip"`
	Action    string    `gorm:"not null" json:"action"` // Login/Logout/Upload/Download/Delete/BatchDelete/Rename/Mkdir
	Message   string    `gorm:"not null" json:"message"`
	Path      string    `gorm:"not null" json:"path"` // 操作路径
}

type SftpLogslist struct {
	ID        uint   `json:"id"`
	CreatedAt string `json:"created_at"`
	Username  string `json:"username"`
	Reviewer  string `json:"reviewer"`
	IP        string `json:"ip"`
	Action    string `json:"action"`
	Message   string `json:"message"`
	Path      string `json:"path"`
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
			ID:        log.ID,
			CreatedAt: formattedTime,
			Username:  log.Username,
			Reviewer:  log.Reviewer,
			IP:        log.IP,
			Action:    log.Action,
			Message:   log.Message,
			Path:      log.Path,
		})
	}
	// 查询满足条件的日志总数
	if err := countQuery.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	return logslist, totalCount, nil
}
