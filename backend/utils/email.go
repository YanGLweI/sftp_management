package utils

import (
	"sftpbackend/config"

	"gopkg.in/gomail.v2"
)

type Email struct {
	Subject string // 邮件主题
	Body    string // 邮件内容
	Attach  string // 附件路径
}

func NewEmail(email Email) *Email {
	return &email
}

func (e *Email) Send() error {
	// 获取配置
	host := config.GlobalConfig.Email.Host
	port := config.GlobalConfig.Email.Port
	password := config.GlobalConfig.Email.Password
	from := config.GlobalConfig.Email.From
	tos := config.GlobalConfig.Email.Tos

	// 组装邮件
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", tos...)
	m.SetHeader("Subject", e.Subject)
	m.SetBody("text/html", e.Body)
	if e.Attach != "" {
		m.Attach(e.Attach)
	}
	d := gomail.NewDialer(host, port, from, password)

	return d.DialAndSend(m)
}

func (e *Email) toMessage() string {
	from := config.GlobalConfig.Email.From
	to := "lw.yang@ho-brostech.com"

	// 组装邮件内容
	message := "From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + e.Subject + "\r\n" +
		"\r\n" + e.Body

	return message
}
