package scheduler

import (
	"fmt"
	"sftpbackend/config"
	"sftpbackend/dao"
	"sftpbackend/models"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

var conf = config.GlobalConfig.Script

// 任务函数映射
var taskFuncMap = map[string]func(){
	"CheckKasperskyIsolateZone": CheckKasperskyIsolateZone,
	"SendKasperskyReport":       SendKasperskyReport,
	"SystemUpdate":              SystemUpdate,
	"SystemSecurityCheck":       SystemSecurityCheck,
	"SendUpdateReport":          SendUpdateReport,
	"SendHardeningReport":       SendHardeningReport,
}

var GlobalScheduler *Scheduler

// Scheduler 管理定时任务的调度器
type Scheduler struct {
	c       *cron.Cron
	quit    chan struct{}
	taskIDs map[string]cron.EntryID // 任务名 -> 任务ID，用于定位已注册的任务
	mu      sync.RWMutex            // 保护taskIDs的并发读写安全
}

// ! Run 启动调度器
func Run() {
	// 初始化调度器
	s := newScheduler()

	// 设置全局调度器
	GlobalScheduler = &s

	// 注册定时任务
	s.register()

	s.start()
}

// ! newScheduler 创建新的 Scheduler 实例
func newScheduler() Scheduler {
	return Scheduler{
		c:       cron.New(),
		quit:    make(chan struct{}),
		taskIDs: make(map[string]cron.EntryID),
	}
}

// ! addTask 添加定时任务，基于 cron 表达式
func (s *Scheduler) addTask(name string, cronExpr string, taskFunc func(), immediate bool) {
	if immediate {
		logrus.Infof("%s run immediately...", "Job-"+name)
		taskFunc() // 立即执行
	}

	entryID, err := s.c.AddFunc(cronExpr, func() {
		logrus.Infof("[%s] execute at %s", "Job-"+name, time.Now().Format("2006-01-02 15:04:05"))
		taskFunc()
	})
	if err != nil {
		logrus.Errorf("[%s] add failed, err: %v.", "Job-"+name, err)
	}
	// 加锁更新任务ID映射
	s.mu.Lock()
	s.taskIDs[name] = entryID
	s.mu.Unlock()

	logrus.Infof("[%s] added successfully (cron: %s, entryID: %d)", name, cronExpr, entryID)
}

// ! updateTask 更新指定名称的定时任务
// 参数说明：
//
//	name: 任务名称（需与注册时一致）
//	newCronExpr: 新的cron表达式
//	taskFunc: 任务执行函数（需与原任务一致，避免逻辑不一致）
//	immediate: 更新后是否立即执行一次
//
// 返回值：更新失败时返回错误
func (s *Scheduler) UpdateTask(name string, newCronExpr string, enable bool, immediate bool) error {
	// 加写锁（保证更新操作原子性）
	s.mu.Lock()
	defer s.mu.Unlock()

	// 1. 检查任务是否存在
	entryID, exists := s.taskIDs[name]
	if !exists {
		return fmt.Errorf("[%s] task not found", name)
	}

	// 2. 移除旧任务
	s.c.Remove(entryID)
	logrus.Infof("[%s] old task removed (entryID: %d)", name, entryID)

	// 3. 包装新任务函数（保持日志格式一致）
	wrappedFunc := func() {
		logrus.Infof("[%s] execute at %s", name, time.Now().Format("2006-01-02 15:04:05"))
		taskFuncMap[name]()
	}

	// 4. 添加新任务（使用新的cron表达式）
	if enable {
		newEntryID, err := s.c.AddFunc(newCronExpr, wrappedFunc)
		if err != nil {
			return fmt.Errorf("[%s] add new task failed: %v", name, err)
		}
		s.taskIDs[name] = newEntryID
		logrus.Infof("[%s] updated successfully (new cron: %s, new entryID: %d)", name, newCronExpr, newEntryID)
	}

	// 5. 若需要，立即执行一次新任务
	if immediate {
		logrus.Infof("[%s] run immediately after update...", name)
		wrappedFunc()
	}

	return nil
}

// ! start 启动调度器
func (s *Scheduler) start() {
	s.c.Start()
}

// ! stop 停止调度器
func (s *Scheduler) stop() {
	s.c.Stop()
	close(s.quit)
}

// ! register 注册所有计划任务
func (s *Scheduler) register() {
	// 查询所有计划任务
	var tasks []models.Scheduler
	dao.DB.Find(&tasks)

	// 注册所有计划任务
	for _, task := range tasks {
		// 注册启用的任务
		if task.Enable {
			s.addTask(task.TaskName, task.Cron, func() {
				taskFuncMap[task.TaskName]()
			}, false)
		}
	}
}

// ! Stop 对外暴露的停止函数
func Stop() {
	if GlobalScheduler != nil {
		GlobalScheduler.stop()
		logrus.Info("调度器已优雅停止（等待当前任务执行完成）")
	} else {
		logrus.Warn("调度器未初始化，无需停止")
	}
}
