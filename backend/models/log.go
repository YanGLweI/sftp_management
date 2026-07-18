package models

import (
	"fmt"
	"sftpbackend/dao"
	"time"
)

type Log struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	Username  string    `gorm:"not null" json:"username"`
	IP        string    `gorm:"not null" json:"ip"`
	Action    string    `gorm:"not null" json:"action"`
	Message   string    `gorm:"not null" json:"message"`
}
type Logslist struct {
	ID        uint   `json:"id"`
	CreatedAt string `json:"created_at"`
	Username  string `json:"username"`
	IP        string `json:"ip"`
	Action    string `json:"action"`
	Message   string `json:"message"`
}

// ! 创建一条日志
func (l *Log) CreateLog() error {
	return dao.DB.Create(l).Error
}

// ! 获取日志列表
func (l *Log) GetLogList(page, limit int, time string, username string, logtype string) (logslist []Logslist, totalCount int64, err error) {
	var logs []Log
	// 根据页码和每页数量计算偏移量
	offset := (page - 1) * limit
	var logtypePattern []string
	switch logtype {
	case "login":
		logtypePattern = []string{"Login", "Logout"}
	case "operation":
		logtypePattern = []string{"Add", "Update", "Delete"}
	default:
		err = fmt.Errorf("logtype error")
		return nil, 0, err
	}
	// 构建查询条件
	query := dao.DB.Offset(offset).Limit(limit).Order("created_at desc").
		Where("action in ?", logtypePattern)
	// 先构建一个用于统计总数的查询条件副本，不添加分页相关设置
	countQuery := dao.DB.Model(&Log{}).Where("action in ?", logtypePattern)

	// 根据是否传入time和username来动态添加相应的查询条件
	if time != "" && username != "" {
		likePattern := "%" + time + "%"
		// 带分页的条件，同时添加时间和用户名条件
		query = query.Where("created_at LIKE? AND username LIKE?", likePattern, "%"+username+"%")
		// 不带分页的条件，同时添加时间和用户名条件
		countQuery = countQuery.Where("created_at LIKE? AND username LIKE?", likePattern, "%"+username+"%")
	} else if time != "" {
		likePattern := "%" + time + "%"
		// 带分页的条件，仅添加时间条件
		query = query.Where("created_at LIKE?", likePattern)
		// 不带分页的条件，仅添加时间条件
		countQuery = countQuery.Where("created_at LIKE?", likePattern)
	} else if username != "" {
		// 带分页的条件，仅添加用户名条件
		query = query.Where("username LIKE?", "%"+username+"%")
		// 不带分页的条件，仅添加用户名条件
		countQuery = countQuery.Where("username LIKE?", "%"+username+"%")
	}
	// 查询满足条件的所有用户
	if err := query.Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	for _, log := range logs {
		formattedTime := log.CreatedAt.Format("2006-01-02 15:04:05")
		logslist = append(logslist, Logslist{
			ID:        log.ID,
			CreatedAt: formattedTime,
			Username:  log.Username,
			IP:        log.IP,
			Action:    log.Action,
			Message:   log.Message,
		})
	}
	// 查询满足条件的用户总数
	if err := countQuery.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	return logslist, totalCount, nil
}
