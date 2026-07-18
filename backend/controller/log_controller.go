package controller

import (
	"fmt"
	"net/http"
	"os/exec"
	"sftpbackend/config"
	"sftpbackend/models"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ! GetLogList 获取日志列表
func GetLogList(c *gin.Context) {
	// 获取路径参数中的页码和每页记录数
	pageStr := c.Param("page")
	limitStr := c.Param("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "无效的页码参数",
		})
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "无效的每页记录数参数",
		})
		return
	}
	// 获取查询参数中的用户名
	time := c.Query("datetime")
	username := c.Query("username")
	logtype := c.Query("logtype")
	var log models.Log
	if logs, total, err := log.GetLogList(page, limit, time, username, logtype); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": gin.H{
				"logs":  logs,
				"total": total,
			},
		})
	}
}

// ! 获取sftp传输日志
func GetSftpLog(c *gin.Context) {
	// 读取日志文件配置
	configLog := config.GlobalConfig.LogFiles
	configScript := config.GlobalConfig.Script
	dateStr := c.Param("date")
	// 解析日期字符串，按照指定的格式将其转换为time.Time类型
	layout := "20060102"
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		fmt.Println("日期解析错误:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	// 获取当前日期
	now := time.Now()
	// 格式化输出当前日期，按照"20060102"格式输出
	nowDateStr := now.Format(layout)

	// 定义一个命令变量
	var commandFirst *exec.Cmd
	var command *exec.Cmd

	// 如果传入日期等于当前日期，返回不同的命令
	if dateStr == nowDateStr {
		commandFirst = exec.Command("bash", configScript.FixSSHLogScript, configLog.LogFile)

		command = exec.Command("cat", configLog.LogFile)
	} else {
		// 将日期加一天
		nextDay := t.AddDate(0, 0, 1)
		// 格式化输出新的日期，同样按照"20060102"格式输出
		newDateStr := nextDay.Format(layout)
		commandFirst = exec.Command("bash", configScript.FixSSHLogScript, fmt.Sprintf(configLog.DailyLogFile, newDateStr))
		command = exec.Command("cat", fmt.Sprintf(configLog.DailyLogFile, newDateStr))
	}

	// 修复日志
	commandFirst.Run()

	// 执行命令获取sftp日志
	sftplog, err := command.CombinedOutput()
	if err != nil {
		fmt.Println("获取sftp日志失败:", err, "\n命令输出:", string(sftplog))
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "获取sftp日志失败: " + err.Error(),
		})
		return
	}

	// 将sftplog转换为字符串
	sftplogStr := string(sftplog)

	// 以\n换行符分割字符串为结合体切片
	sftploglines := strings.Split(sftplogStr, "\n")
	var sftploglinesNew []models.TransferLog
	for _, line := range sftploglines {
		if line == "" {
			continue
		}
		sftploglinesNew = append(sftploglinesNew, models.TransferLog{Log: line})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"sftplog": sftploglinesNew,
		},
	})
}
