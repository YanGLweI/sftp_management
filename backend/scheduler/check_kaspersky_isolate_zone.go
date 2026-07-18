package scheduler

import (
	"fmt"
	"os"
	"sftpbackend/kaspersky"
	"sftpbackend/utils"
	"strings"

	"github.com/sirupsen/logrus"
)

func CheckKasperskyIsolateZone() {
	html := kaspersky.HtmlHead

	htmlBody := kaspersky.HtmlBodyStart
	threat := kaspersky.ThreatReport()
	htmlBody += threat + kaspersky.HtmlBodyEnd

	html += htmlBody

	if strings.Contains(threat, "正常") {
		return
	}

	hostname, _ := os.Hostname()

	email := utils.NewEmail(utils.Email{
		Subject: fmt.Sprintf("【%s】DataBuffer Kaspersky 威胁报告", hostname),
		Body:    html,
	})

	if err := email.Send(); err != nil {
		logrus.Errorf("邮件发送失败: %v", err)
	}
}
