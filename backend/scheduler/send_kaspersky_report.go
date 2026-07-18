package scheduler

import (
	"fmt"
	"os"
	"sftpbackend/kaspersky"
	"sftpbackend/utils"
	"time"

	"github.com/sirupsen/logrus"
)

func SendKasperskyReport() {
	html := kaspersky.HtmlHead

	htmlBody := kaspersky.HtmlBodyStart
	summary := kaspersky.Summary()
	threat := kaspersky.ThreatReport()
	htmlBody += summary + threat + kaspersky.HtmlBodyEnd

	html += htmlBody

	hostname, _ := os.Hostname()

	email := utils.NewEmail(utils.Email{
		Subject: fmt.Sprintf("【%s】DataBuffer Kaspersky 保护状态报告", hostname),
		Body:    html,
		Attach:  "data/kaspersky_app_info_" + time.Now().Format("2006-01-02") + ".txt",
	})

	if err := email.Send(); err != nil {
		logrus.Errorf("邮件发送失败: %v", err)
	}
}
