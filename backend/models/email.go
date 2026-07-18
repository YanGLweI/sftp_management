package models

import (
	"os"
	"sftpbackend/config"

	"gopkg.in/gomail.v2"
)

// EmailConfig 邮件配置结构体
type EmailConfig struct {
	To      []string `json:"to" binding:"required"`      // 收件人邮箱列表
	Cc      []string `json:"cc" binding:"required"`      // 抄送人邮箱列表,可以为空列表
	Path    string   `json:"path"`                       // 附件路径
	Subject string   `json:"subject" binding:"required"` // 邮件主题
}

func (e *EmailConfig) SendEmail(Attachments []string) (err error) {
	// 获取邮件配置信息
	config := config.GlobalConfig.Email

	// 创建一个新的邮件对象
	m := gomail.NewMessage()
	// 设置发件人
	m.SetHeader("From", config.From)
	// 设置收件人
	m.SetHeader("To", e.To...)
	// 设置抄送人，有抄送人才设置
	if len(e.Cc) > 0 {
		m.SetHeader("Cc", e.Cc...)
	}
	// 设置邮件主题
	m.SetHeader("Subject", e.Subject)
	// 设置邮件正文
	m.SetBody("text/html", config.Body)
	// 如果有附件，添加附件到邮件
	// if e.Path != "" {
	// 	file, err := os.Open(e.Path)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	defer file.Close()

	// 	// 添加附件到邮件
	// 	m.Attach(e.Path)
	// }
	// 循环添加所有附件
	if len(Attachments) > 0 {
		for _, attachmentPath := range Attachments {
			if attachmentPath == "" {
				continue // 跳过空路径
			}
			// 检查文件是否存在
			if _, err := os.Stat(attachmentPath); os.IsNotExist(err) {
				continue // 文件不存在，跳过
			}
			// 添加附件到邮件，Attach方法会处理文件打开和读取
			m.Attach(attachmentPath)
		}
	}
	// 创建一个SMTP连接对象
	d := gomail.NewDialer(config.Host, config.Port, config.From, config.Password)

	// 发送邮件
	err = d.DialAndSend(m)

	return
}
