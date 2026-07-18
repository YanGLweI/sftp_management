package utils

import (
	"fmt"
	"os"
	"sftpbackend/config"
)

func CreateAttachment(name, password, loginType string) error {
	// 读取配置文件中的内容模板
	config := config.GlobalConfig.Email
	contentTemplate := config.Content

	switch loginType {
	case "Password":
		loginType = "Normal"
	case "KeyFile":
		password = "Use the private key file to login"
	case "both":
		loginType = "Password and KeyFile"
	}

	// 创建一个txt文件,写入账号和密码信息
	content := fmt.Sprintf(contentTemplate, loginType, name, password)
	filePath := fmt.Sprintf("%s-%s.txt", config.Subject, name)

	// 创建一个txt文件
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("创建文件失败:", err)
	}
	defer file.Close()

	// 将填充后的内容写入文件
	_, err = file.Write([]byte(content))
	if err != nil {
		fmt.Println("写入文件内容失败:", err)
	}
	return err
}
