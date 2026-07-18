package sshutils

import (
	"bytes"
	"fmt"

	"golang.org/x/crypto/ssh"
)

// SSHSession 用于表示SSH会话的结构体
type SSHSession struct {
	session *ssh.Session
	client  *ssh.Client
}

// NewSSHSession 创建一个新的SSH会话
func NewSSHSession(user, password, serverAddress string) (*SSHSession, error) {
	// 设置SSH连接配置
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(password)},
		// 设置接受所有主机密钥的回调函数
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// 连接到SSH服务器
	client, err := ssh.Dial("tcp", serverAddress+":22", sshConfig)
	if err != nil {
		return nil, err
	}

	// 创建一个会话
	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, err
	}

	return &SSHSession{
		session: session,
		client:  client,
	}, nil
}

// RunCommand 在SSH会话中执行命令并返回输出结果
func (s *SSHSession) RunCommand(command string) (string, error) {
	var outputBuf bytes.Buffer

	// 将命令输出重定向到outputBuf
	s.session.Stdout = &outputBuf

	// 执行命令
	err := s.session.Run(command)
	if err != nil {
		return "", err
	}

	// 获取命令输出结果
	output := outputBuf.String()

	return output, nil
}

// Close 关闭SSH会话和连接
func (s *SSHSession) Close() {
	s.session.Close()
	s.client.Close()
	fmt.Println("SSH Close Session")
}
