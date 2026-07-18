// package kaspersky 用于处理卡巴斯基Endpoint Security相关信息的采集与报告生成
package kaspersky

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// appInfo 全局变量，存储卡巴斯基应用信息的每行内容
var appInfo []string

// HtmlHead HTML页面的头部内容，包含DOCTYPE声明、元标签、样式定义等
var HtmlHead = `<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Databuffer Kaspersky 保护状态报告</title>
    <style>
        body {
            font-family: "Segoe UI", Arial, sans-serif;
            margin: 0;
            padding: 30px;
            background: #f4f6f9;
            color: #333;
        }

        h2 {
            text-align: center;
            margin-bottom: 25px;
            color: #2c3e50;
        }

        /* 状态卡片容器 */
        .status-container {
            display: flex;
            justify-content: space-between;
            margin-bottom: 30px;
            gap: 20px;
        }

        /* 单个状态卡片 */
        .status-card {
            flex: 1;
            background: #ffffff;
            border-radius: 10px;
            margin: 0 10px;
            padding: 20px;
            text-align: center;
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
            border-left: 6px solid #2ecc71;
        }

        .status-title {
            font-size: 17px;
            color: #555;
            margin-bottom: 10px;
        }

        .status-value {
            font-size: 22px;
            font-weight: bold;
            color: #2ecc71;
        }

        /* 表格容器 */
        .table-container {
            padding: 0 10px;
        }

        /* 表格样式 */
        table {
            width: 100%;
            border-collapse: collapse;
            background: #fff;
            border-radius: 8px;
            overflow: hidden;
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
        }

        th {
            background: #2c7be5;
            color: #fff;
            padding: 12px;
            font-size: 15px;
            letter-spacing: 1px;
        }

        td {
            padding: 12px;
            border-bottom: 1px solid #e5e7eb;
        }

        tr:last-child td {
            border-bottom: none;
        }

        tr:nth-child(even) td {
            background: #f9fbfc;
        }

        tr:hover td {
            background: #eef5ff;
        }
    </style>
</head>`

// HtmlBodyStart HTML主体的开始标签
var HtmlBodyStart = `<body>`

// HtmlBodyEnd HTML主体的结束标签（包含整个HTML文档的闭合）
var HtmlBodyEnd = `</body></html>`

// Summary 生成卡巴斯基系统信息摘要的HTML内容
// 包含应用基本信息、数据库状态、扫描任务状态等表格展示
func Summary() string {
	// 获取卡巴斯基应用详细信息
	appInfo = getAppInfo()

	// 反病毒库更新状态
	var updateStatus string
	// 扫描任务状态
	var scanStatus string

	for _, line := range appInfo {
		if strings.Contains(line, "数据库的上次发布日期:") {
			lastPublishDate := strings.Split(line, ":  ")[1] // 分割获取数据库发布日期
			// 判断更新时间是否在最近7天内
			layout := "2006-01-02 15:04:05"
			publishTime, _ := time.Parse(layout, strings.Trim(lastPublishDate, " "))
			if time.Since(publishTime) <= 7*24*time.Hour {
				updateStatus = "正常"
			} else {
				updateStatus = "异常"
			}
		}
		if strings.Contains(line, "Scan_My_Computer 任务的上次运行日期:") {
			lastScanDate := strings.Split(line, ":  ")[1] // 分割获取扫描任务日期
			// 判断扫描时间是否在最近24小时内
			layout := "2006-01-02 15:04:05"
			scanTime, _ := time.Parse(layout, strings.Trim(lastScanDate, " "))
			if time.Since(scanTime) <= 24*time.Hour {
				scanStatus = "正常"
			} else {
				scanStatus = "异常"
			}
		}
	}

	// 初始化摘要HTML内容，包含标题、状态卡片和表格结构
	summary := fmt.Sprintf(`
    <h2>卡巴斯基 Endpoint Security for Linux 系统信息</h2>

    <!-- 状态卡片模块 -->
    <div class="status-container">
        <div class="status-card">
            <div class="status-title">反病毒库更新状态</div>
            <div class="status-value">%s</div>
        </div>
        <div class="status-card">
            <div class="status-title">扫描任务状态</div>
            <div class="status-value">%s</div>
        </div>
    </div>
    <!-- 系统信息表格 -->
    <div class="table-container">
        <table>
            <tr><th>项目</th><th>内容</th></tr>`, updateStatus, scanStatus)
	// 添加计算机名称信息行
	hostname, _ := os.Hostname()
	summary += "\n  <tr><td>计算机名称</td><td>" + hostname + "</td></tr>\n"

	// 遍历应用信息，提取关键字段并拼接表格行
	for _, line := range appInfo {
		switch {
		case strings.Contains(line, "名称:"):
			name := strings.Split(line, ":  ")[1] // 分割获取名称值
			summary += fmt.Sprintf("  <tr><td>%s</td><td>%s</td></tr>\n", "应用名称", strings.Trim(name, " "))
		case strings.Contains(line, "版本:"):
			version := strings.Split(line, ":  ")[1] // 分割获取版本值
			summary += fmt.Sprintf("  <tr><td>%s</td><td>%s</td></tr>\n", "应用版本", strings.Trim(version, " "))
		case strings.Contains(line, "数据库的上次发布日期:"):
			lastPublishDate := strings.Split(line, ":  ")[1] // 分割获取数据库发布日期
			summary += fmt.Sprintf("  <tr><td>%s</td><td>%s</td></tr>\n", "数据库的上次发布日期", strings.Trim(lastPublishDate, " "))
		case strings.Contains(line, "应用程序数据库已加载:"):
			loaded := strings.Split(line, ":  ")[1] // 分割获取数据库加载状态
			summary += fmt.Sprintf("  <tr><td>%s</td><td>%s</td></tr>\n", "应用程序数据库已加载", strings.Trim(loaded, " "))
		case strings.Contains(line, "Scan_My_Computer 任务的上次运行日期:"):
			lastRunDate := strings.Split(line, ":  ")[1] // 分割获取扫描任务上次运行日期
			summary += fmt.Sprintf("  <tr><td>%s</td><td>%s</td></tr>", "Scan_My_Computer 任务的上次运行日期", strings.Trim(lastRunDate, " "))
		}
	}

	// 闭合表格标签
	summary += "</table></div>"
	return summary
}

// ThreatReport 生成卡巴斯基威胁报告的HTML内容
// 根据隔离区状态判断是否存在威胁，并展示对应状态卡片
func ThreatReport() string {
	// 获取隔离区（威胁报告）状态
	status := getThreatReport()
	var report string

	// 判断隔离区是否有威胁文件：无威胁则显示正常，否则显示异常
	if status == "" {
		report = `<h2>卡巴斯基 Endpoint Security for Linux 威胁报告</h2>
    <!-- 状态卡片模块 -->
    <div class="status-container">
        <div class="status-card">
            <div class="status-title">隔离区状态</div>
            <div class="status-value">正常</div>
        </div>
    </div>
    `
	} else {
		report = fmt.Sprintf(`<h2>卡巴斯基 Endpoint Security for Linux 威胁报告</h2>
    <!-- 状态卡片模块 -->
    <div class="status-container">
        <div class="status-card">
            <div class="status-title">隔离区状态</div>
            <div class="status-value">异常</div>
        </div>
        <div class="status-card">
            <div class="status-title">威胁详情</div>
            <div style="font-size: 14px;color: #000000ff;">%s</div>
        </div>
    </div>
    `, status)
	}

	return report
}

// getAppInfo 执行kesl-control --app-info命令获取卡巴斯基应用信息
// 并将信息保存到data目录下的日志文件，返回按行分割的信息切片
func getAppInfo() []string {
	// 执行kesl-control命令获取应用信息（kesl-control是卡巴斯基的命令行工具）
	cmd := exec.Command("kesl-control", "--app-info")
	output_bytes, err := cmd.Output()
	if err != nil {
		return []string{"N/A"} // 命令执行失败返回N/A
	}
	output := string(output_bytes)

	// 创建data目录（用于存储日志文件），不存在则创建
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.Mkdir("data", 0777) // 0777表示所有用户可读可写可执行
	}

	// 生成日志文件名（包含当前日期）
	filename := "data/kaspersky_app_info_" + time.Now().Format("2006-01-02") + ".txt"
	file, err := os.Create(filename) // 创建文件
	if err != nil {
		fmt.Println(err) // 打印文件创建错误
	}
	defer file.Close() // 函数结束时关闭文件

	// 将应用信息写入文件
	file.WriteString(output)

	// 按换行符分割信息，返回每行内容的切片
	return strings.Split(output, "\n")
}

// getThreatReport 执行kesl-control -Q --query命令获取隔离区（威胁）报告
// 返回命令执行的输出结果字符串
func getThreatReport() string {
	// 查询今天是否检查到威胁
	today := time.Now().Format("2006-01-02")
	dateQuery := today + " 00:00:00"
	query := fmt.Sprintf("Date > '%s' and EventType == 'ThreatDetected'", dateQuery)
	cmd5 := exec.Command("kesl-control", "-E", "--query", query)
	output5, err := cmd5.CombinedOutput()
	if err != nil {
		return string(output5) // 命令执行失败返回输出内容
	}
	return string(output5)
}
