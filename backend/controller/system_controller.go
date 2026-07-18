package controller

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sftpbackend/config"
	"sftpbackend/models"
	"sftpbackend/scheduler"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true // 允许所有跨域请求，生产环境应限制
// 	},
// }

// KasperskyItem 定义返回给前端的结构体
type KasperskyItem struct {
	Item    string `json:"item"`    // 对应前部分（键）
	Content string `json:"content"` // 对应后部分（值）
}

// 计划任务信息
type KasperskySchedule struct {
	RuleType            string `json:"ruleType"`            // 任务启动计划:once Monthly Weekly Daily Hourly Minutely
	RunMissedStartRules string `json:"runMissedStartRules"` // 是否运行错过的任务计划 yes no
	Date                string `json:"date"`                // 日期
	Time                string `json:"time"`                // 时间
	Interval            string `json:"interval"`            // 任务执行间隔（分钟、小时、天、周、月）
	RandomInterval      string `json:"randomInterval"`      // 随机延迟时间（分钟）
}
type LastRunDate struct {
	Date string `json:"Date"` // 上次运行日期
}

var conf = config.GlobalConfig.Script

// ! 系统更新模块

// 立即更新
func SystemUpdate(c *gin.Context) {
	// 获取Websocket连接
	ws := models.NewWebSocketManager()
	conn, err := ws.UpgradeConnection(c.Writer, c.Request)
	if err != nil {
		log.Printf("连接WebSocket失败: %v", err)
		return
	}
	defer ws.CloseConnection(conn)

	// 创建命令
	cmd := exec.Command("bash", conf.SystemUpdateScript)

	// 获取标准输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("获取标准输出管道失败: %v", err)
		conn.WriteMessage(1, []byte("获取标准输出管道失败: "+err.Error()))
		return
	}
	// 获取标准错误管道
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("获取标准错误管道失败: %v", err)
		conn.WriteMessage(1, []byte("获取标准错误管道失败: "+err.Error()))
		return
	}

	// 执行命令
	if err := cmd.Start(); err != nil {
		log.Printf("命令启动失败: %v", err)
		conn.WriteMessage(1, []byte("命令启动失败: "+err.Error()))
		return
	}

	// 使用WaitGroup等待输出处理完成
	var wg sync.WaitGroup
	wg.Add(2)

	// 读取标准输出
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			// 发送到客户端
			if err := conn.WriteMessage(1, []byte(line)); err != nil {
				log.Printf("Error writing message: %v", err)
				return
			}
		}
	}()

	// 读取标准错误
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			// 发送到客户端
			if err := conn.WriteMessage(1, []byte(line)); err != nil {
				log.Printf("Error writing message: %v", err)
				return
			}
		}
	}()

	// 等待命令执行完成
	wg.Wait()
	if err := cmd.Wait(); err != nil {
		conn.WriteMessage(1, []byte("[ERROR] 命令执行失败: "+err.Error()))
	} else {
		conn.WriteMessage(1, []byte("[SUCCESS] 命令执行成功"))
	}

}

// 获取更新历史
func GetUpdateHistory(c *gin.Context) {

	var u models.UpdateMain
	// var updateHistory []models.UpdateMain
	// var total int64
	updateHistory, total, err := u.GetUpdateHistory(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}
	var updateHistoryResponse []models.UpdateMainResponse
	for _, item := range *updateHistory {
		updateHistoryResponse = append(updateHistoryResponse, models.UpdateMainResponse{
			ID:          item.ID,
			Hostname:    item.Hostname,
			UpdateTime:  item.UpdateTime.Format("2006-01-02 15:04:05"),
			Status:      item.Status,
			Duration:    item.Duration,
			UpdateBrief: item.UpdateBrief,
		})
	}
	// 返回结果
	c.JSON(200, gin.H{
		"code":    200,
		"data":    gin.H{"updateHistory": updateHistoryResponse, "total": total},
		"message": "success",
	})
}

// 获取系统更新详情
func GetUpdateDetail(c *gin.Context) {
	var u models.UpdateMain
	id := c.Param("id")
	updateDetail, err := u.GetUpdateDetail(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}
	// 返回结果
	c.JSON(200, gin.H{
		"code":    200,
		"data":    updateDetail,
		"message": "success",
	})
}

// 获取系统更新计划
func GetUpdateSchedule(c *gin.Context) {
	var s models.Scheduler
	schedule, err := s.GetUpdateSchedule()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}
	cronExpr := schedule.Cron
	var resp models.SchedulerResp

	// 解析表达式
	resp, err = resp.CronToSchedulerResp(cronExpr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}

	resp.Enable = schedule.Enable

	// 返回结果
	c.JSON(200, gin.H{
		"code":    200,
		"data":    resp,
		"message": "success",
	})
}

// 设置系统更新计划
func SetUpdateSchedule(c *gin.Context) {
	var schedule models.Scheduler
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}

	// 更新计划任务
	if err := scheduler.GlobalScheduler.UpdateTask("SystemUpdate", schedule.Cron, schedule.Enable, false); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}

	// 更新数据库记录
	if err := schedule.UpdateSchedule(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}

	// 返回结果
	c.JSON(200, gin.H{
		"code":    200,
		"data":    nil,
		"message": "success",
	})
}

// 获取更新报告计划
func GetUpdateReportSchedule(c *gin.Context) {
	var s models.Scheduler
	schedule, err := s.GetUpdateReportSchedule()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}
	cronExpr := schedule.Cron
	var resp models.SchedulerResp

	// 解析表达式
	resp, err = resp.CronToSchedulerResp(cronExpr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}

	resp.Enable = schedule.Enable

	// 返回结果
	c.JSON(200, gin.H{
		"code":    200,
		"data":    resp,
		"message": "success",
	})
}

// 设置系统更新报告计划
func SetUpdateReportSchedule(c *gin.Context) {
	var schedule models.Scheduler
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}

	// 更新计划任务
	if err := scheduler.GlobalScheduler.UpdateTask("SendUpdateReport", schedule.Cron, schedule.Enable, false); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}

	// 更新数据库记录
	if err := schedule.UpdateReportSchedule(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}

	// 返回结果
	c.JSON(200, gin.H{
		"code":    200,
		"data":    nil,
		"message": "success",
	})
}

// ! 卡巴斯基模块
// 获取卡巴斯基应用信息
func Antivirus(c *gin.Context) {
	// 定义第一个命令 `kesl-control --app-info`
	cmd1 := exec.Command("kesl-control", "--app-info") // 参数拆分为切片，每个参数一个元素

	// 定义第二个命令 `sed -n '1p;2p;5p;6p;12p;14p;15p'`
	// sed的参数拆分：-n 是第一个参数，过滤规则是第二个参数（整体作为字符串）
	cmd2 := exec.Command("sed", "-n", "1p;2p;5p;6p;12p;14p;15p")

	// 关键：将cmd1的标准输出接到cmd2的标准输入
	// 相当于：kesl-control --app-info | sed -n '1p;2p;5p;6p;12p;14p;15p'
	cmd2.Stdin, _ = cmd1.StdoutPipe() // 获取cmd1的stdout管道，作为cmd2的stdin

	// 定义缓冲区，用于捕获cmd2的最终输出（也可以直接设置为os.Stdout，输出到控制台）
	var outputBuf bytes.Buffer
	cmd2.Stdout = &outputBuf

	// 定义第三个命令，查询Scan_My_Computer 任务的上次结束日期
	query := "TaskName == 'Scan_My_Computer' and TaskState == 'stopped'"
	cmd3 := exec.Command("kesl-control", "-E", "--query", query, "-n", "1", "--reverse")
	output3, err := cmd3.CombinedOutput()
	if err != nil {
		log.Printf("任务的上次结束日期：%v，%s", err, string(output3))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error() + string(output3),
		})
		return
	}
	// 解析输出，获取结束日期,查找Date开头的行，等于号后面提取日期
	lines3 := strings.Split(string(output3), "\n")
	var lastRunDate string
	for _, line := range lines3 {
		if strings.HasPrefix(line, "Date") {
			lastRunDate = strings.TrimSpace(line[len("Date="):])
			break
		}
	}
	// var lastRunDates []LastRunDate
	// output3Str := string(output3)
	// // 解析JSON
	// if err := json.Unmarshal([]byte(output3Str), &lastRunDates); err != nil {
	// 	log.Printf("解析JSON失败：%v，%s", err, output3Str)
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"code":    200,
	// 		"data":    nil,
	// 		"message": err.Error() + output3Str,
	// 	})
	// 	return
	// }
	// // 提取上次运行日期
	// var lastRunDate string
	// if len(lastRunDates) > 0 {
	// 	lastRunDate = lastRunDates[0].Date
	// }

	// 定义第4个命令，查询当前Scan_My_Computer 任务的状态
	cmd4 := exec.Command("kesl-control", "--get-task-state", "2")
	output4, err := cmd4.CombinedOutput()
	if err != nil {
		log.Printf("查询任务状态失败：%v，%s", err, string(output4))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error() + string(output4),
		})
		return
	}
	// 解析输出，获取最后一行:后的状态
	lines4 := strings.Split(string(output4), "\n")
	var taskState string
	// 倒序找最后一个非空行
	for i := len(lines4) - 1; i >= 0; i-- {
		trimmed := strings.TrimSpace(lines4[i])
		if trimmed != "" {
			colonIdx := strings.Index(trimmed, ":")
			if colonIdx != -1 {
				taskState = strings.TrimSpace(trimmed[colonIdx+1:])
			}
			break
		}
	}

	// 启动第二个命令（先启动cmd2，再启动cmd1，避免管道阻塞）
	if err := cmd2.Start(); err != nil {
		log.Printf("启动sed命令失败：%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}

	// 启动第一个命令
	if err := cmd1.Run(); err != nil {
		log.Printf("启动kesl-control命令失败：%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}

	// 等待第二个命令执行完成
	if err := cmd2.Wait(); err != nil {
		log.Printf("sed命令执行失败：%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}
	// 获取最终输出结果
	output := outputBuf.String()
	var appInfo []KasperskyItem

	// 按行分割字符串
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		// 去除每行首尾空白字符
		trimmedLine := strings.TrimSpace(line)
		// 跳过空行
		if trimmedLine == "" {
			continue
		}

		// 查找第一个冒号的位置（分割键和值）
		colonIdx := strings.Index(trimmedLine, ":")
		if colonIdx == -1 {
			continue // 无冒号的行跳过
		}

		// 提取键和值（去除首尾空白）
		item := strings.TrimSpace(trimmedLine[:colonIdx])
		content := strings.TrimSpace(trimmedLine[colonIdx+1:])

		// 添加到结果切片
		appInfo = append(appInfo, KasperskyItem{
			Item:    item,
			Content: content,
		})
	}
	appInfo = append(appInfo, KasperskyItem{
		Item:    "Scan_My_Computer 任务的上次结束日期",
		Content: lastRunDate,
	})
	appInfo = append(appInfo, KasperskyItem{
		Item:    "Scan_My_Computer 任务的当前状态",
		Content: taskState,
	})

	// 返回结果
	c.JSON(200, gin.H{
		"code":    200,
		"data":    appInfo,
		"message": "success",
	})
}

// 获取隔离区状态
func IsolationZone(c *gin.Context) {
	// 查询是否检查到威胁
	// query := "EventType == 'ThreatDetected' and TaskType == 'ODS'"

	// 查询今天是否检查到威胁
	today := time.Now().Format("2006-01-02")
	dateQuery := today + " 00:00:00"
	query := fmt.Sprintf("Date > '%s' and EventType == 'ThreatDetected'", dateQuery)
	cmd5 := exec.Command("kesl-control", "-E", "--query", query)
	output5, err := cmd5.CombinedOutput()
	if err != nil {
		log.Printf("查询威胁失败：%v，%s", err, string(output5))
		c.JSON(200, gin.H{
			"code":    200,
			"data":    "",
			"message": string(output5),
		})
		return
	}
	// 返回结果
	c.JSON(200, gin.H{
		"code":    200,
		"data":    string(output5),
		"message": "success",
	})
}

// 获取计划任务信息
func Schedule(c *gin.Context) {
	// 通过命令kesl-control --get-schedule 2 获取任务2的计划任务信息
	cmd := exec.Command("kesl-control", "--get-schedule", "2")
	outputBytes, err := cmd.Output()
	if err != nil {
		log.Printf("获取计划任务信息失败：%v", err)
		c.JSON(200, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}
	output := string(outputBytes)
	var schedule KasperskySchedule

	// 按行分割字符串
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		// 去除每行首尾空白字符
		trimmedLine := strings.TrimSpace(line)
		// 跳过空行
		if trimmedLine == "" {
			continue
		}

		// 查找第一个冒号的位置（分割键和值）
		colonIdx := strings.Index(trimmedLine, "=")
		if colonIdx == -1 {
			continue // 无冒号的行跳过
		}

		// 提取键和值（去除首尾空白）
		item := strings.TrimSpace(trimmedLine[:colonIdx])
		content := strings.TrimSpace(trimmedLine[colonIdx+1:])

		// 添加到结果切片
		switch item {
		case "RuleType":
			schedule.RuleType = content
		case "RunMissedStartRules":
			schedule.RunMissedStartRules = content
		case "RandomInterval":
			schedule.RandomInterval = content
		case "StartTime":
			// 判断是否有空格
			if strings.Contains(content, " ") {
				parts := strings.Split(content, " ")
				schedule.Date = parts[0]

				// 替换日期中月份英文为数字
				schedule.Date = strings.ReplaceAll(schedule.Date, "Jan", "01")
				schedule.Date = strings.ReplaceAll(schedule.Date, "Feb", "02")
				schedule.Date = strings.ReplaceAll(schedule.Date, "Mar", "03")
				schedule.Date = strings.ReplaceAll(schedule.Date, "Apr", "04")
				schedule.Date = strings.ReplaceAll(schedule.Date, "May", "05")
				schedule.Date = strings.ReplaceAll(schedule.Date, "Jun", "06")
				schedule.Date = strings.ReplaceAll(schedule.Date, "Jul", "07")
				schedule.Date = strings.ReplaceAll(schedule.Date, "Aug", "08")
				schedule.Date = strings.ReplaceAll(schedule.Date, "Sep", "09")
				schedule.Date = strings.ReplaceAll(schedule.Date, "Oct", "10")
				schedule.Date = strings.ReplaceAll(schedule.Date, "Nov", "11")
				schedule.Date = strings.ReplaceAll(schedule.Date, "Dec", "12")

				if strings.Contains(parts[1], ";") {
					parts := strings.Split(parts[1], ";")
					schedule.Time = parts[0]
					schedule.Interval = parts[1]
				} else {
					schedule.Time = parts[1]
					schedule.Interval = ""
				}
			} else {
				schedule.Date = ""
				if strings.Contains(content, ";") {
					parts := strings.Split(content, ";")
					schedule.Time = parts[0]
					schedule.Interval = parts[1]
				} else {
					schedule.Time = content
					schedule.Interval = ""
				}
			}
		}
	}
	// 返回结果
	c.JSON(200, gin.H{
		"code":    200,
		"data":    schedule,
		"message": "success",
	})
}

// 设置卡巴斯基计划任务
func SetSchedule(c *gin.Context) {
	// 获取请求体
	var schedule KasperskySchedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}

	var startTime string
	if schedule.Date == "" {
		startTime = fmt.Sprintf("%s;%s", schedule.Time, schedule.Interval)
	} else {
		if schedule.Interval == "" {
			startTime = fmt.Sprintf("%s %s", schedule.Date, schedule.Time)
		} else {
			startTime = fmt.Sprintf("%s %s;%s", schedule.Date, schedule.Time, schedule.Interval)
		}
	}
	// 构建计划任务文件schedule.ini
	scheduleIni := fmt.Sprintf(`RuleType=%s
RunMissedStartRules=%s
RandomInterval=%s
StartTime=%s
`, schedule.RuleType, schedule.RunMissedStartRules, schedule.RandomInterval, startTime)
	// 写入计划任务文件
	if err := os.WriteFile("schedule.ini", []byte(scheduleIni), 0644); err != nil {
		log.Printf("写入计划任务文件失败：%v", err)
		c.JSON(200, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}
	// 执行命令kesl-control --set-schedule 2 --file schedule.ini 设置任务2的计划任务
	cmd := exec.Command("kesl-control", "--set-schedule", "2", "--file", "schedule.ini")
	outputBytes, err := cmd.Output()
	fmt.Println(string(outputBytes))
	if err != nil {
		log.Printf("设置计划任务失败：%v", err)
		c.JSON(200, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}
	// 返回结果
	c.JSON(200, gin.H{
		"code":    200,
		"data":    nil,
		"message": "success",
	})
}

// 启动卡巴斯基扫描
func StartScan(c *gin.Context) {
	// 执行命令kesl-control --scan --type full 启动全量扫描
	cmd := exec.Command("kesl-control", "--start-task", "2")
	outputBytes, err := cmd.CombinedOutput()
	// 获取错误输出

	fmt.Println(string(outputBytes))
	if err != nil {
		log.Printf("启动扫描失败：%v，%s", err, string(outputBytes))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error() + "：" + string(outputBytes),
		})
		return
	}
	// 返回结果
	c.JSON(200, gin.H{
		"code":    200,
		"data":    nil,
		"message": string(outputBytes),
	})
}

// 获取卡巴斯基报告计划任务
func GetKReportSchedule(c *gin.Context) {
	var s models.Scheduler
	schedule, err := s.GetKReportSchedule()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}
	cronExpr := schedule.Cron
	var resp models.SchedulerResp

	// 解析表达式
	resp, err = resp.CronToSchedulerResp(cronExpr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}

	resp.Enable = schedule.Enable

	// 返回结果
	c.JSON(200, gin.H{
		"code":    200,
		"data":    resp,
		"message": "success",
	})
}

// 设置卡巴斯基报告计划
func SetKReportSchedule(c *gin.Context) {
	var schedule models.Scheduler
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}

	// 更新计划任务
	if err := scheduler.GlobalScheduler.UpdateTask("SendKasperskyReport", schedule.Cron, schedule.Enable, false); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}

	// 更新数据库记录
	if err := schedule.UpdateKReportSchedule(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}

	// 返回结果
	c.JSON(200, gin.H{
		"code":    200,
		"data":    nil,
		"message": "success",
	})
}

// ! 系统加固检查模块
// 分页获取系统加固检查列表
func GetSystemSecurityCheck(c *gin.Context) {
	var s models.SystemSecurity
	checklist, total, err := s.GetSystemSecurityCheck(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}
	// 返回结果
	c.JSON(200, gin.H{
		"code":    200,
		"data":    gin.H{"checklist": checklist, "total": total},
		"message": "success",
	})
}

// 立即启动系统加固任务
func StartSystemSecurityCheck(c *gin.Context) {
	// 执行系统加固检查脚本
	scheduler.SystemSecurityCheck()
	// 返回结果
	c.JSON(200, gin.H{
		"code":    200,
		"data":    nil,
		"message": "success",
	})
}

// 获取系统加固计划
func GetSystemSecuritySchedule(c *gin.Context) {
	var s models.Scheduler
	schedule, err := s.GetSystemSecuritySchedule()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}
	cronExpr := schedule.Cron
	var resp models.SchedulerResp

	// 解析表达式
	resp, err = resp.CronToSchedulerResp(cronExpr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}

	resp.Enable = schedule.Enable

	// 返回结果
	c.JSON(200, gin.H{
		"code":    200,
		"data":    resp,
		"message": "success",
	})
}

// 设置系统加固计划
func SetSystemSecuritySchedule(c *gin.Context) {
	var schedule models.Scheduler
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}

	// 更新计划任务
	if err := scheduler.GlobalScheduler.UpdateTask("SystemSecurityCheck", schedule.Cron, schedule.Enable, false); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}

	// 更新数据库记录
	if err := schedule.UpdateSystemSecuritySchedule(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}

	// 返回结果
	c.JSON(200, gin.H{
		"code":    200,
		"data":    nil,
		"message": "success",
	})
}

// 获取系统加固报告计划任务
func GetSystemSecurityReportSchedule(c *gin.Context) {
	var s models.Scheduler
	schedule, err := s.GetSystemSecurityReportSchedule()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}
	cronExpr := schedule.Cron
	var resp models.SchedulerResp

	// 解析表达式
	resp, err = resp.CronToSchedulerResp(cronExpr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}

	resp.Enable = schedule.Enable

	// 返回结果
	c.JSON(200, gin.H{
		"code":    200,
		"data":    resp,
		"message": "success",
	})
}

// 设置系统加固报告计划
func SetSystemSecurityReportSchedule(c *gin.Context) {
	var schedule models.Scheduler
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}

	// 更新计划任务
	if err := scheduler.GlobalScheduler.UpdateTask("SendHardeningReport", schedule.Cron, schedule.Enable, false); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}

	// 更新数据库记录
	if err := schedule.UpdateSystemSecurityReportSchedule(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    200,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}

	// 返回结果
	c.JSON(200, gin.H{
		"code":    200,
		"data":    nil,
		"message": "success",
	})
}
