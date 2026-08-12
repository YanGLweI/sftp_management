package utils

import (
	"fmt"
	"io"
	"math/rand"
	"path/filepath"
	"sftpbackend/config"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// SFTPConnection 封装单个SFTP连接的所有资源
type SFTPConnection struct {
	SftpClient   *sftp.Client
	SSHClient    *ssh.Client
	CreateTime   time.Time // 连接创建时间，用于过期清理
	Username     string    // 关联的用户名，可选
	LastUsedTime time.Time // 最后使用时间
}

// 连接管理器（全局）
type sftpConnManager struct {
	connMap map[string]*SFTPConnection // key: token, value: 连接实例
	mu      sync.RWMutex               // 读写锁，保证并发安全
}

// 全局连接管理器实例
var SFTPConnManager = &sftpConnManager{
	connMap: make(map[string]*SFTPConnection),
}

// 生成唯一Token（UUID）
func generateConnToken() string {
	return uuid.New().String()
}

// 添加连接到管理器
func (m *sftpConnManager) AddConn(conn *SFTPConnection) string {
	m.mu.Lock()
	defer m.mu.Unlock()

	token := generateConnToken()
	m.connMap[token] = conn
	return token
}

// 通过Token获取连接
func (m *sftpConnManager) GetConn(token string) (*SFTPConnection, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	conn, exists := m.connMap[token]
	if !exists {
		return nil, fmt.Errorf("SFTP Token 已失效，请重新登录")
	}
	// 更新最后使用时间
	conn.LastUsedTime = time.Now()
	return conn, nil
}

// 删除并关闭连接
func (m *sftpConnManager) RemoveConn(token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	conn, exists := m.connMap[token]
	if !exists {
		// return fmt.Errorf("连接不存在")
		return nil // 已经不存在了，算是成功删除
	}

	// 关闭连接
	if err := conn.Close(); err != nil {
		return err
	}

	// 从映射中删除
	delete(m.connMap, token)
	return nil
}

// 定期清理过期连接（比如30分钟未使用）
func (m *sftpConnManager) CleanExpiredConns(expireTime time.Duration) {
	ticker := time.NewTicker(5 * time.Minute) // 每5分钟检查一次
	defer ticker.Stop()

	for range ticker.C {
		m.mu.Lock()
		now := time.Now()
		for token, conn := range m.connMap {
			if now.Sub(conn.LastUsedTime) > expireTime {
				// 关闭并删除过期连接
				conn.Close()
				delete(m.connMap, token)
			}
		}
		m.mu.Unlock()
	}
}

// ! 初始化SFTP连接（密码登录）
func NewSFTPConnection(user, password string) (*SFTPConnection, error) {
	conf := config.GlobalConfig.SFTP
	sshConfig := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshClient, err := ssh.Dial("tcp", conf.Host+":"+fmt.Sprint(conf.Port), sshConfig)
	if err != nil {
		return nil, fmt.Errorf("SSH连接失败: %v", err)
	}

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		sshClient.Close()
		return nil, fmt.Errorf("SFTP客户端创建失败: %v", err)
	}

	return &SFTPConnection{
		SftpClient:   sftpClient,
		SSHClient:    sshClient,
		CreateTime:   time.Now(),
		Username:     user,
		LastUsedTime: time.Now(),
	}, nil
}

// ! 初始化SFTP连接（密钥登录）
func NewSFTPConnectionByKey(user string, key []byte) (*SFTPConnection, error) {
	conf := config.GlobalConfig.SFTP
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("私钥解析失败: %v", err)
	}

	sshConfig := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshClient, err := ssh.Dial("tcp", conf.Host+":"+fmt.Sprint(conf.Port), sshConfig)
	if err != nil {
		return nil, fmt.Errorf("SSH连接失败: %v", err)
	}

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		sshClient.Close()
		return nil, fmt.Errorf("SFTP客户端创建失败: %v", err)
	}

	return &SFTPConnection{
		SftpClient:   sftpClient,
		SSHClient:    sshClient,
		CreateTime:   time.Now(),
		Username:     user,
		LastUsedTime: time.Now(),
	}, nil
}

// 关闭连接
func (conn *SFTPConnection) Close() error {
	if conn == nil {
		return nil
	}

	var err error
	if conn.SftpClient != nil {
		err = conn.SftpClient.Close()
		conn.SftpClient = nil
	}

	if conn.SSHClient != nil {
		if sshErr := conn.SSHClient.Close(); sshErr != nil && err == nil {
			err = sshErr
		}
		conn.SSHClient = nil
	}

	return err
}

// 创建上传文件（srcFile 支持任意 io.Reader，便于 multipart part 流式直写）
func (conn *SFTPConnection) CreateUploadFile(dstPath string, srcFile io.Reader) error {
	if conn == nil || conn.SftpClient == nil {
		return fmt.Errorf("SFTP连接未初始化")
	}

	dstFile, err := conn.SftpClient.Create(dstPath)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %v", err)
	}
	defer dstFile.Close()

	if _, err := dstFile.ReadFrom(srcFile); err != nil {
		return fmt.Errorf("复制文件内容失败: %v", err)
	}

	return nil
}

// 创建目录
func (conn *SFTPConnection) CreateFolder(path, folderName string) error {
	if conn == nil || conn.SftpClient == nil {
		return fmt.Errorf("SFTP连接未初始化")
	}
	// 先检查目录是否已经存在

	fullPath := filepath.Join(path, folderName)
	// 先检查目录是否已经存在
	_, err := conn.SftpClient.Stat(fullPath)
	if err == nil {
		return fmt.Errorf("目录已存在")
	}

	if err := conn.SftpClient.Mkdir(fullPath); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	return nil
}

// 检测目录是否可删除（通过重命名测试）
func (conn *SFTPConnection) isDeletableDirectory(path string) error {
	// 生成临时名称，避免冲突
	tmpPath := path + "__delete_test_" + fmt.Sprintf("%d", rand.Int63())

	// 尝试重命名
	err := conn.SftpClient.Rename(path, tmpPath)
	if err != nil {
		return fmt.Errorf("无法删除目录 %s（权限不足或其他错误）: %v", path, err)
	}
	// 立即改回
	if err := conn.SftpClient.Rename(tmpPath, path); err != nil {
		// 理论上重命名成功后改回不应该失败，若失败则严重错误
		return fmt.Errorf("严重错误：重命名测试后无法恢复目录名: %v", err)
	}
	return nil
}

// 递归删除目录
func (conn *SFTPConnection) DeleteDirectory(path string) error {
	if conn == nil || conn.SftpClient == nil {
		return fmt.Errorf("SFTP连接未初始化")
	}

	// 检查目录是否可删除
	if err := conn.isDeletableDirectory(path); err != nil {
		return err
	}

	files, err := conn.SftpClient.ReadDir(path)
	if err != nil {
		return fmt.Errorf("读取目录失败: %v", err)
	}

	for _, file := range files {
		filePath := filepath.Join(path, file.Name())
		if file.IsDir() {
			if err := conn.DeleteDirectory(filePath); err != nil {
				return err
			}
		} else {
			if err := conn.SftpClient.Remove(filePath); err != nil {
				return err
			}
		}
	}

	return conn.SftpClient.RemoveDirectory(path)
}
