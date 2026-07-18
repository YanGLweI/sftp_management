package models

import (
	"fmt"
	"sftpbackend/sshutils"
)

type TransferLog struct {
	Log string `json:"log"`
}

func GetSSHCommandResult(command string) (string, error) {
	// test ssh connect
	session, err := sshutils.NewSSHSession("root", "XH.svr.passw0rd", "10.7.254.188")
	if err != nil {
		fmt.Println("连接失败", err.Error())
		return "", err
	}
	defer session.Close()

	// 执行命令
	cmd := command

	output, err := session.RunCommand(cmd)
	if err != nil {
		fmt.Println("执行命令失败", err.Error())
		return "", err
	}
	// fmt.Println(output)
	return output, nil
}
