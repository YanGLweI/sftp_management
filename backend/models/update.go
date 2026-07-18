package models

import (
	"log"
	"sftpbackend/dao"
	"sftpbackend/tools"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UpdateMain 系统更新核心日志表
type UpdateMain struct {
	gorm.Model
	Hostname    string    `gorm:"column:hostname;type:varchar(50);not null;comment:主机名"`
	UpdateTime  time.Time `gorm:"column:update_time;not null;comment:更新执行时间"`
	Status      string    `gorm:"column:status;type:enum('成功','失败','暂无可用更新');not null;comment:更新状态"`
	Duration    float64   `gorm:"column:duration;type:decimal(10,2);not null;comment:更新耗时（秒）"`
	UpdateBrief string    `gorm:"column:update_brief;type:text;not null;comment:更新摘要"`
	// 关联关系
	// Details UpdateDetails `gorm:"foreignKey:MainID;references:ID"`
}

// UpdateDetails 系统更新详情表
type UpdateDetails struct {
	gorm.Model
	MainID        int    `gorm:"column:main_id;not null;comment:关联主表ID;index"`
	UpdateDetails string `gorm:"column:update_details;type:LONGTEXT;not null;comment:更新完整详情"`

	// 关联关系
	UpdateMain UpdateMain `gorm:"foreignKey:MainID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}

type UpdateMainResponse struct {
	ID          uint    `json:"id"`
	Hostname    string  `json:"hostname"`
	UpdateTime  string  `json:"update_time"`
	Status      string  `json:"status"`
	Duration    float64 `json:"duration"`
	UpdateBrief string  `json:"update_brief"`
}

// GetUpdateHistory 获取系统更新历史记录
func (u *UpdateMain) GetUpdateHistory(c *gin.Context) (updateHistory *[]UpdateMain, total int64, err error) {
	// 获取分页参数
	var pageReq tools.PageOption
	if err := c.ShouldBindQuery(&pageReq); err != nil {
		log.Println(err.Error())
		return nil, 0, err
	}
	// 获取分页参数
	pageOption := tools.NewPageOption(pageReq.PageNum, pageReq.PageSize)

	db := dao.DB
	// 查询更新历史记录，按更新时间降序排序
	err = db.Offset(pageOption.PageNum).Limit(pageOption.PageSize).Order("update_time DESC").Find(&updateHistory).Error
	if err != nil {
		return nil, 0, err
	}

	// 查询总记录数
	err = db.Model(&UpdateMain{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	return updateHistory, total, err
}

// GetUpdateDetail 获取系统更新详情
func (u *UpdateMain) GetUpdateDetail(id string) (updateDetail *UpdateDetails, err error) {
	db := dao.DB
	// 查询更新详情记录
	err = db.Where("main_id = ?", id).Preload("UpdateMain").First(&updateDetail).Error
	if err != nil {
		return nil, err
	}
	return updateDetail, err
}

// GetLatestUpdate 获取最新的系统更新记录
func (u *UpdateMain) GetLatestUpdate() (updateDetail *UpdateMain, err error) {
	db := dao.DB
	// 查询最新的更新记录，按更新时间降序排序
	err = db.Order("update_time DESC").First(&updateDetail).Error
	if err != nil {
		return nil, err
	}
	return updateDetail, err
}
