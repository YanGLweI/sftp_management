package models

import (
	"fmt"
	"sftpbackend/dao"
	"strings"

	"gorm.io/gorm"
)

type Scheduler struct {
	gorm.Model
	TaskName string `gorm:"column:task_name;type:varchar(100);not null;comment:任务名称"`
	Cron     string `gorm:"column:cron;type:varchar(100);not null;comment:定时任务表达式" json:"cron"`
	Enable   bool   `gorm:"column:enable;type:bool;not null;default:false;comment:是否启用" json:"enable"`
}

type SchedulerResp struct {
	RuleType string `json:"ruleType"` // Weekly or Daily
	Time     string `json:"time"`     // HH:mm
	Enable   bool   `json:"enable"`   // 是否启用
	Interval string `json:"interval"` // 如果是Weekly ： Mon,Tue,Wed,Thu,Fri,Sat,Sun,否则为空
}

// cronWeekMap Cron周数字（0=周日，1=周一...6=周六）到星期缩写的映射
var cronWeekMap = map[string]string{
	"0": "Sun",
	"1": "Mon",
	"2": "Tue",
	"3": "Wed",
	"4": "Thu",
	"5": "Fri",
	"6": "Sat",
	// 兼容部分系统用7表示周日的情况
	"7": "Sun",
}

// 获取系统更新计划
func (s *Scheduler) GetUpdateSchedule() (Scheduler, error) {
	var schedule Scheduler
	if err := dao.DB.Where("task_name = ?", "SystemUpdate").First(&schedule).Error; err != nil {
		return schedule, err
	}
	return schedule, nil
}

// 更新数据库记录
func (s *Scheduler) UpdateSchedule() error {
	// s.Enable 可能为false,直接updates GORM会不处理false值,导致遗漏
	updateData := map[string]interface{}{
		"cron":   s.Cron,
		"enable": s.Enable,
	}
	if err := dao.DB.
		Model(&Scheduler{}).
		Where("task_name = ?", "SystemUpdate").
		Updates(updateData).Error; err != nil {
		return err
	}
	return nil
}

// 获取系统更新报告计划
func (s *Scheduler) GetUpdateReportSchedule() (Scheduler, error) {
	var schedule Scheduler
	if err := dao.DB.Where("task_name = ?", "SendUpdateReport").First(&schedule).Error; err != nil {
		return schedule, err
	}
	return schedule, nil
}

// 更新系统更新报告计划数据库记录
func (s *Scheduler) UpdateReportSchedule() error {
	// s.Enable 可能为false,直接updates GORM会不处理false值,导致遗漏
	updateData := map[string]interface{}{
		"cron":   s.Cron,
		"enable": s.Enable,
	}
	if err := dao.DB.
		Model(&Scheduler{}).
		Where("task_name = ?", "SendUpdateReport").
		Updates(updateData).Error; err != nil {
		return err
	}
	return nil
}

// 获取卡巴斯基报告计划
func (s *Scheduler) GetKReportSchedule() (Scheduler, error) {
	var schedule Scheduler
	if err := dao.DB.Where("task_name = ?", "SendKasperskyReport").First(&schedule).Error; err != nil {
		return schedule, err
	}
	return schedule, nil
}

// 更新卡巴斯基报告计划数据库记录
func (s *Scheduler) UpdateKReportSchedule() error {
	// s.Enable 可能为false,直接updates GORM会不处理false值,导致遗漏
	updateData := map[string]interface{}{
		"cron":   s.Cron,
		"enable": s.Enable,
	}
	if err := dao.DB.
		Model(&Scheduler{}).
		Where("task_name = ?", "SendKasperskyReport").
		Updates(updateData).Error; err != nil {
		return err
	}
	return nil
}

// 获取系统加固报告计划
func (s *Scheduler) GetSystemSecurityReportSchedule() (Scheduler, error) {
	var schedule Scheduler
	if err := dao.DB.Where("task_name = ?", "SendHardeningReport").First(&schedule).Error; err != nil {
		return schedule, err
	}
	return schedule, nil
}

// 更新系统加固报告计划数据库记录
func (s *Scheduler) UpdateSystemSecurityReportSchedule() error {
	// s.Enable 可能为false,直接updates GORM会不处理false值,导致遗漏
	updateData := map[string]interface{}{
		"cron":   s.Cron,
		"enable": s.Enable,
	}
	if err := dao.DB.
		Model(&Scheduler{}).
		Where("task_name = ?", "SendHardeningReport").
		Updates(updateData).Error; err != nil {
		return err
	}
	return nil
}

// 将Cron表达式转换为SchedulerResp
func (sr *SchedulerResp) CronToSchedulerResp(cronExpr string) (SchedulerResp, error) {
	// 分割Cron表达式（分 时 日 月 周）
	parts := strings.Fields(cronExpr)
	if len(parts) != 5 {
		return SchedulerResp{}, fmt.Errorf("无效的cron表达式：%s", cronExpr)
	}

	minute, hour, _, _, week := parts[0], parts[1], parts[2], parts[3], parts[4]

	// 处理时间（补前导零，确保HH:mm格式）
	formatHour := fmt.Sprintf("%02s", hour)
	formatMinute := fmt.Sprintf("%02s", minute)
	timeStr := fmt.Sprintf("%s:%s", formatHour, formatMinute)

	resp := SchedulerResp{
		Time: timeStr,
	}

	// 判断规则类型并处理Interval
	if week == "*" {
		// 周为通配符 → 每天执行
		resp.RuleType = "Daily"
		resp.Interval = ""
	} else {
		// 周指定具体天 → 每周执行，转换为星期缩写
		resp.RuleType = "Weekly"
		weekParts := strings.Split(week, ",") // 分割周字段，例如"1,2,3"
		var weekNames []string
		for _, w := range weekParts {
			name, ok := cronWeekMap[w]
			if !ok {
				return SchedulerResp{}, fmt.Errorf("无效的星期数字：%s", w)
			}
			weekNames = append(weekNames, name)
		}
		resp.Interval = strings.Join(weekNames, ",")
	}

	return resp, nil
}

// 获取系统加固计划
func (s *Scheduler) GetSystemSecuritySchedule() (Scheduler, error) {
	var schedule Scheduler
	if err := dao.DB.Where("task_name = ?", "SystemSecurityCheck").First(&schedule).Error; err != nil {
		return schedule, err
	}
	return schedule, nil
}

// 更新系统加固计划数据库记录
func (s *Scheduler) UpdateSystemSecuritySchedule() error {
	// s.Enable 可能为false,直接updates GORM会不处理false值,导致遗漏
	updateData := map[string]interface{}{
		"cron":   s.Cron,
		"enable": s.Enable,
	}
	if err := dao.DB.
		Model(&Scheduler{}).
		Where("task_name = ?", "SystemSecurityCheck").
		Updates(updateData).Error; err != nil {
		return err
	}
	return nil
}
