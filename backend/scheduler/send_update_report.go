package scheduler

import (
	"fmt"
	"os"
	"sftpbackend/report"
	"sftpbackend/utils"

	"github.com/sirupsen/logrus"
)

func SendUpdateReport() {
	html := report.HtmlHead
	htmlBody := report.HtmlBodyStart
	summary := report.UpdateInfo()
	htmlBody += summary + report.HtmlBodyEnd
	html += htmlBody
	hostname, _ := os.Hostname()

	email := utils.NewEmail(utils.Email{
		Subject: fmt.Sprintf("【%s】DataBuffer 更新报告", hostname),
		Body:    html,
	})

	if err := email.Send(); err != nil {
		logrus.Errorf("邮件发送失败: %v", err)
	}
}
