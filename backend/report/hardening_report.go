package report

import (
	"fmt"
	"sftpbackend/models"
	"time"

	"github.com/sirupsen/logrus"
)

var hardeningInfo string

// HtmlHead HTML页面的头部内容，包含DOCTYPE声明、元标签、样式定义等
var HardeningHtmlHead = `<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>系统加固报告</title>
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

				h3 {
            text-align: center;
            margin-bottom: 15px;
            color: #3c6489ff;
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

        /* 加固详情表单样式 */
        .hardening-form {
            background: #fff;
						margin: 0 10px;
            border-radius: 10px;
            padding: 25px;
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
            display: flex;
            flex-wrap: wrap;
            gap: 20px;
        }

        .form-column {
            flex: 1;
            min-width: 400px;
        }

        .form-group {
            margin-bottom: 18px;
        }

        .form-label {
            display: block;
            font-weight: 600;
            color: #2c3e50;
            margin-bottom: 8px;
            font-size: 14px;
        }

        .form-value {
            width: 100%;
            padding: 10px;
            border: 1px solid #e5e7eb;
            border-radius: 6px;
            background: #f9fbfc;
            font-size: 14px;
            color: #333;
            min-height: 40px;
            box-sizing: border-box;
            white-space: pre-wrap;
            word-wrap: break-word;
        }
    </style>
</head>`

// HtmlBodyStart HTML主体的开始标签
var HardeningHtmlBodyStart = `<body>`

// HtmlBodyEnd HTML主体的结束标签（包含整个HTML文档的闭合）
var HardeningHtmlBodyEnd = `</body></html>`

func HardeningInfo() string {
	var h models.SystemSecurity
	latestHardening, err := h.FindLatest()
	if err != nil {
		logrus.Printf("获取最新的系统加固结果失败：%v", err)
		return ""
	}

	// 系统加固状态
	var hardeningStatus string
	// 系统加固状态类
	var hardeningStatusClass string

	// 加固任务状态
	var taskStatus string
	var taskStatusClass string

	if latestHardening.Result == "正常" {
		hardeningStatus = "正常"
		hardeningStatusClass = "status-normal"
	} else {
		hardeningStatus = "异常"
		hardeningStatusClass = "status-abnormal"
	}

	// 判断日期，超过1天，任务状态为异常
	if time.Since(latestHardening.Date) > 24*time.Hour {
		taskStatus = "异常"
		taskStatusClass = "status-abnormal"
	} else {
		taskStatus = "正常"
		taskStatusClass = "status-normal"
	}

	hardeningInfo = fmt.Sprintf(`
		<h2>DataBuffer 系统加固报告</h2>
		<div class="status-container">
			<div class="status-card">
					<div class="status-title">系统加固状态</div>
					<div class="status-value"><span class="%s">%s</span></div>
			</div>
			<div class="status-card">
					<div class="status-title">加固任务状态</div>
					<div class="status-value"><span class="%s">%s</span></div>
			</div>
    </div>
		<!-- 系统加固表格 -->
		<div class="table-container">
			<table>
				<tr><th>主机名</th><th>IP地址</th><th>操作系统</th><th>加固状态</th><th>日期</th></tr>
				<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>
			</table>
		</div>
		`, hardeningStatusClass, hardeningStatus, taskStatusClass, taskStatus, latestHardening.Hostname,
		latestHardening.IP, latestHardening.Operasystem, latestHardening.Result,
		latestHardening.Date.Format("2006-01-02 15:04:05"))

	// 构建加固详情表单
	hardeningInfo += fmt.Sprintf(`
		<!-- 加固详情卡片 -->
		<h2>系统加固详情</h2>
		<h3>DNF/Repo 配置</h3>
		<div class="hardening-form">
			
			<!-- 第一列 -->
			<div class="form-column">

				<div class="form-group">
					<label class="form-label">内核版本</label>
					<div class="form-value">%s</div>
				</div>
				
				<div class="form-group">
					<label class="form-label">redhat.repo中gpgcheck配置值</label>
					<div class="form-value">%s</div>
				</div>
			</div>

			<!-- 第二列 -->
			<div class="form-column">
				<div class="form-group">
					<label class="form-label">dnf.conf中gpgcheck配置值</label>
					<div class="form-value">%s</div>
				</div>
			</div>
		</div>`,
		latestHardening.Kernel,
		latestHardening.RedhatRepoGpgcheck,
		latestHardening.DnfConfGpgcheck)

	hardeningInfo += fmt.Sprintf(`
		<h3>密码策略</h3>
		<div class="hardening-form">
			<!-- 第一列 -->
			<div class="form-column">
				<div class="form-group">
					<label class="form-label">密码最大有效期（天）</label>
					<div class="form-value">%s</div>
				</div>
				<div class="form-group">
					<label class="form-label">密码最小长度</label>
					<div class="form-value">%s</div>
				</div>
				<div class="form-group">
					<label class="form-label">账户非活动锁定天数</label>
					<div class="form-value">%s</div>
				</div>
				<div class="form-group">
					<label class="form-label">终端自动超时时间（秒）</label>
					<div class="form-value">%s</div>
				</div>
			</div>

			<!-- 第二列 -->
			<div class="form-column">
				<div class="form-group">
					<label class="form-label">密码最小修改间隔（天）</label>
					<div class="form-value">%s</div>
				</div>
				<div class="form-group">
					<label class="form-label">密码过期警告天数</label>
					<div class="form-value">%s</div>
				</div>
				<div class="form-group">
					<label class="form-label">默认GID</label>
					<div class="form-value">%s</div>
				</div>
			</div>

		</div>`,
		latestHardening.PASSMAXDAYS,
		latestHardening.PASSMINLEN,
		latestHardening.INACTIVE,
		latestHardening.TMOUT,
		latestHardening.PASSMINDAYS,
		latestHardening.PASSWARNAGE,
		latestHardening.GID,
	)

	hardeningInfo += fmt.Sprintf(`
		<h3>Cron/At 任务配置</h3>
		<div class="hardening-form">
			

			<!-- 第一列 -->
			<div class="form-column">
				<div class="form-group">
					<label class="form-label">cron启用</label>
					<div class="form-value ">%s</div>
				</div>

				<div class="form-group">
					<label class="form-label">cron.hourly目录</label>
					<div class="form-value">%s</div>
				</div>

				<div class="form-group">
					<label class="form-label">cron.weekly目录</label>
					<div class="form-value">%s</div>
				</div>

				<div class="form-group">
					<label class="form-label">cron.deny文件</label>
					<div class="form-value">%s</div>
				</div>

				<div class="form-group">
					<label class="form-label">cron.allow文件</label>
					<div class="form-value">%s</div>
				</div>
			</div>

			<!-- 第二列 -->
			<div class="form-column">
				<div class="form-group">
					<label class="form-label">crontab文件</label>
					<div class="form-value ">%s</div>
				</div>
				<div class="form-group">
					<label class="form-label">cron.daily目录</label>
					<div class="form-value">%s</div>
				</div>
				<div class="form-group">
					<label class="form-label">cron.monthly目录</label>
					<div class="form-value">%s</div>
				</div>
				<div class="form-group">
					<label class="form-label">at.deny文件</label>
					<div class="form-value">%s</div>
				</div>
				<div class="form-group">
					<label class="form-label">at.allow文件</label>
					<div class="form-value">%s</div>
				</div>
			</div>
		</div>`,
		latestHardening.Cron,
		latestHardening.CronHourly,
		latestHardening.CronWeekly,
		latestHardening.CronDeny,
		latestHardening.CronAllow,
		latestHardening.Crontab,
		latestHardening.CronDaily,
		latestHardening.CronMonthly,
		latestHardening.AtDeny,
		latestHardening.AtAllow,
	)

	hardeningInfo += fmt.Sprintf(`
		<h3>SSHD 配置</h3>
		<div class="hardening-form">
			

			<!-- 第一列 -->
			<div class="form-column">
				<div class="form-group">
					<label class="form-label">SSHD配置文件</label>
					<div class="form-value ">%s</div>
				</div>

				<div class="form-group">
					<label class="form-label">X11转发开关（yes/no）</label>
					<div class="form-value">%s</div>
				</div>

				<div class="form-group">
					<label class="form-label">是否忽略rhosts（yes/no）</label>
					<div class="form-value">%s</div>
				</div>

				<div class="form-group">
					<label class="form-label">是否允许root登录SSH</label>
					<div class="form-value">%s</div>
				</div>

				<div class="form-group">
					<label class="form-label">是否允许用户自定义环境</label>
					<div class="form-value">%s</div>
				</div>

				<div class="form-group">
					<label class="form-label">SSH客户端存活最大次数</label>
					<div class="form-value">%s</div>
				</div>
			</div>

			<!-- 第二列 -->
			<div class="form-column">
				<div class="form-group">
					<label class="form-label">SSH日志级别</label>
					<div class="form-value">%s</div>
				</div>

				<div class="form-group">
					<label class="form-label">SSH最大认证尝试次数</label>
					<div class="form-value">%s</div>
				</div>

				<div class="form-group">
					<label class="form-label">基于主机的认证开关</label>
					<div class="form-value">%s</div>
				</div>

				<div class="form-group">
					<label class="form-label">是否允许空密码登录</label>
					<div class="form-value">%s</div>
				</div>

				<div class="form-group">
					<label class="form-label">SSH客户端存活间隔（秒）</label>
					<div class="form-value">%s</div>
				</div>

				<div class="form-group">
					<label class="form-label">SSH登录宽限时间（秒）</label>
					<div class="form-value">%s</div>
				</div>

			</div>
		</div>`,
		latestHardening.SshdConfig,
		latestHardening.X11Forwarding,
		latestHardening.IgnoreRhosts,
		latestHardening.PermitRootLogin,
		latestHardening.PermitUserEnvironment,
		latestHardening.ClientAliveCountMax,
		latestHardening.LogLevel,
		latestHardening.MaxAuthTries,
		latestHardening.HostbasedAuthentication,
		latestHardening.PermitEmptyPasswords,
		latestHardening.ClientAliveInterval,
		latestHardening.LoginGraceTime,
	)

	hardeningInfo += fmt.Sprintf(`
		<h3>密码复杂度</h3>
		<div class="hardening-form">
			

			<!-- 第一列 -->
			<div class="form-column">
				<div class="form-group">
					<label class="form-label">密码最小长度</label>
					<div class="form-value">%s</div>
				</div>

				<div class="form-group">
					<label class="form-label">数字字符信用值</label>
					<div class="form-value">%s</div>
				</div>

				<div class="form-group">
					<label class="form-label">小写字符信用值</label>
					<div class="form-value">%s</div>
				</div>

				<div class="form-group">
					<label class="form-label">密码历史</label>
					<div class="form-value">%s</div>
				</div>
			</div>

			<!-- 第二列 -->
			<div class="form-column">
				<div class="form-group">
					<label class="form-label">密码字符种类</label>
					<div class="form-value">%s</div>
				</div>

				<div class="form-group">
					<label class="form-label">大写字符信用值</label>
					<div class="form-value">%s</div>
				</div>

				<div class="form-group">
					<label class="form-label">特殊字符信用值</label>
					<div class="form-value">%s</div>
				</div>

			</div>
		</div>`,
		latestHardening.Minlen,
		latestHardening.Dcredit,
		latestHardening.Lcredit,
		latestHardening.PasswordRemember,
		latestHardening.Minclass,
		latestHardening.Ucredit,
		latestHardening.Ocredit,
	)

	hardeningInfo += fmt.Sprintf(`
		<h3>系统文件内容</h3>
		<div class="hardening-form">
			

			<!-- 第一列 -->
			<div class="form-column">
				<div class="form-group">
					<label class="form-label">passwd文件内容</label>
					<div class="form-value ">%s</div>
				</div>

				<div class="form-group">
					<label class="form-label">group文件内容</label>
					<div class="form-value ">%s</div>
				</div>
				
				<div class="form-group">
					<label class="form-label">shadow文件内容</label>
					<div class="form-value ">%s</div>
				</div>
			
				<div class="form-group">
					<label class="form-label">gshadow文件内容</label>
					<div class="form-value ">%s</div>
				</div>
			</div>

			<!-- 第二列 -->
			<div class="form-column">
				<div class="form-group">
					<label class="form-label">passwd-备份文件内容</label>
					<div class="form-value ">%s</div>
				</div>

				<div class="form-group">
					<label class="form-label">group-备份文件内容</label>
					<div class="form-value ">%s</div>
				</div>

				<div class="form-group">
					<label class="form-label">shadow-备份文件内容</label>
					<div class="form-value ">%s</div>
				</div>
				
				<div class="form-group">
					<label class="form-label">gshadow-备份文件内容</label>
					<div class="form-value ">%s</div>
				</div>
			</div>
		</div>`,
		latestHardening.Passwd,
		latestHardening.Group,
		latestHardening.Shadow,
		latestHardening.Gshadow,
		latestHardening.PasswdDash,
		latestHardening.GroupDash,
		latestHardening.ShadowDash,
		latestHardening.GshadowDash,
	)

	hardeningInfo += fmt.Sprintf(`
		<h3>加密/时间 策略</h3>
		<div class="hardening-form">
			
			
			<!-- 第一列 -->
			<div class="form-column">
				<div class="form-group">
					<label class="form-label">加密策略</label>
					<div class="form-value">%s</div>
				</div>
			</div>

			<!-- 第二列 -->
			<div class="form-column">
				<div class="form-group">
					<label class="form-label">NTP</label>
					<div class="form-value">%s</div>
				</div>
			</div>
		</div>`,
		latestHardening.CryptoPolicies,
		latestHardening.NtpServer,
	)

	return hardeningInfo
}
