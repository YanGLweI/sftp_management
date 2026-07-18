package report

import (
	"sftpbackend/models"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var updateInfo string

// HtmlHead HTML页面的头部内容，包含DOCTYPE声明、元标签、样式定义等
var HtmlHead = `<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>系统更新报告</title>
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

        .status-normal {
            color: #2ecc71; /* 绿色 */
        }
					
		.status-abnormal {
            color: #e74c3c;
        }

        .detail-card {
            flex: 1;
            background: #ffffff;
            border-radius: 10px;
            margin: 0 10px;
            padding: 20px;
            text-align: left;
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
            border-left: 6px solid #2ecc71;
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
						text-align: center;
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

// UpdateInfo 生成系统更新信息的HTML内容
func UpdateInfo() string {
	var u models.UpdateMain
	latestUpdate, err := u.GetLatestUpdate()
	if err != nil {
		logrus.Printf("GetLatestUpdate failed: %v", err)
		return ""
	}

	// 系统更新状态
	var updateStatus string
	var updateStatusClass string
	// 更新任务状态
	var updateTaskStatus string
	var updateTaskStatusClass string

	// 判断最近更新的时间是否超过7天
	if time.Since(latestUpdate.UpdateTime) > 7*24*time.Hour {
		updateTaskStatus = "异常"
		updateTaskStatusClass = "status-abnormal"
		updateStatus = "异常"
		updateStatusClass = "status-abnormal"
	} else {
		if latestUpdate.Status == "成功" || latestUpdate.Status == "暂无可用更新" {
			updateStatus = "正常"
			updateTaskStatus = "正常"
			updateStatusClass = "status-normal"
			updateTaskStatusClass = "status-normal"
		} else {
			updateStatus = "异常"
			updateTaskStatus = "异常"
			updateTaskStatusClass = "status-abnormal"
			updateStatusClass = "status-abnormal"
		}
	}
	// 生成HTML内容
	updateInfo = `
		<h2>DataBuffer 系统更新报告</h2>
		<div class="status-container">
			<div class="status-card">
				<div class="status-title">更新任务状态</div>
				<div class="status-value"><span class="` + updateTaskStatusClass + `">` + updateTaskStatus + `</span></div>
			</div>
			<div class="status-card">
				<div class="status-title">系统更新状态</div>
				<div class="status-value"><span class="` + updateStatusClass + `">` + updateStatus + `</span></div>
			</div>
		</div>
		<!-- 系统信息表格 -->
        <div class="table-container">
            <table>
                <tr><th>更新时间</th><th>主机名</th><th>更新状态</th><th>更新摘要</th></tr>
                <tr><td>` + latestUpdate.UpdateTime.Format("2006-01-02 15:04:05") + `</td><td>` + latestUpdate.Hostname + `</td><td>` + latestUpdate.Status + `</td><td>` + latestUpdate.UpdateBrief + `</td></tr>
            </table>
        </div>
	`

	if latestUpdate.Status == "成功" {
		// 获取详情
		id := strconv.FormatUint(uint64(latestUpdate.ID), 10)
		updateDetail, err := u.GetUpdateDetail(id)
		if err != nil {
			logrus.Printf("GetDetail failed: %v", err)
			return ""
		}
		detailStr := strings.ReplaceAll(updateDetail.UpdateDetails, "\n", "<br>")
		// updateDetail是长文本，需要进行HTML转义
		updateInfo += `
        <h2>更新详情</h2>
				<div class="detail-card">
        	<p>` + detailStr + `</p>
				</div>
		`
	}

	return updateInfo
}
