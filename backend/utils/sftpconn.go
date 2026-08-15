package utils

import (
	"fmt"
	"io"
	"math/rand"
	"path/filepath"
	"sftpbackend/config"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/sftp"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

// SFTPConnection 封装单个SFTP连接的所有资源
type SFTPConnection struct {
	SftpClient   *sftp.Client
	SSHClient    *ssh.Client
	CreateTime   time.Time // 连接创建时间，用于过期清理
	Username     string    // 关联的SFTP用户名（专用账号或用户账号）
	LoginType    string    // 登录方式：password/keyfile/hotlabel/chinaunicom
	DomainUser   string    // 域账号（hotlabel/chinaunicom 登录时的实际操作者，其他登录方式为空）
	HomePath     string    // 连接允许访问的根路径（空表示不限制）
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

// ! 初始化SFTP连接（密码登录，不限制访问路径）
func NewSFTPConnection(user, password string) (*SFTPConnection, error) {
	return NewSFTPConnectionWithHome(user, password, "")
}

// ! 初始化SFTP连接（密码登录，可指定允许访问的根路径）
// homePath 为空表示不限制路径
func NewSFTPConnectionWithHome(user, password, homePath string) (*SFTPConnection, error) {
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
		HomePath:     homePath,
		LastUsedTime: time.Now(),
	}, nil
}

// NewSFTPConnectionForModule 初始化SFTP连接（模块专用登录：域控验证通过后用公共账号登录，绑定根路径）
// loginType: hotlabel/chinaunicom 等模块标识；domainUser: 通过域控验证的操作者账号
func NewSFTPConnectionForModule(user, password, homePath, loginType, domainUser string) (*SFTPConnection, error) {
	conn, err := NewSFTPConnectionWithHome(user, password, homePath)
	if err != nil {
		return nil, err
	}
	conn.LoginType = loginType
	conn.DomainUser = domainUser
	return conn, nil
}

// ResolvePath 校验并规范化请求路径，确保不超出连接允许的根路径
// filepath.Clean 处理 ".." 穿越与重复斜杠（如 /hotlabel/../.. Clean 后为 /，前缀校验拒绝）
// 边界处理：当根路径为 "/" 时，允许所有绝对路径
// 越界时记录警告日志（便于审计可疑路径尝试），并返回错误
func (conn *SFTPConnection) ResolvePath(requestPath string) (string, error) {
	cleaned := filepath.Clean(requestPath)
	if conn.HomePath == "" {
		return cleaned, nil // 普通连接不限制
	}
	home := filepath.Clean(conn.HomePath)
	if cleaned == home {
		return cleaned, nil
	}
	// home 为 "/" 时，所有绝对路径都在允许范围内（len(cleaned) > len(home) 且首位为 /）
	if home == "/" {
		if strings.HasPrefix(cleaned, "/") && len(cleaned) > 1 {
			return cleaned, nil
		}
		logrus.WithFields(logrus.Fields{
			"user":       conn.Username,
			"loginType":  conn.LoginType,
			"domainUser": conn.DomainUser,
			"homePath":   conn.HomePath,
			"request":    requestPath,
		}).Warn("可疑路径尝试：路径超出允许范围")
		return "", fmt.Errorf("路径超出允许范围: %s", requestPath)
	}
	// 前缀校验：确保 cleaned 在 home 目录内
	if strings.HasPrefix(cleaned, home) && len(cleaned) > len(home) && cleaned[len(home)] == '/' {
		return cleaned, nil
	}
	logrus.WithFields(logrus.Fields{
		"user":       conn.Username,
		"loginType":  conn.LoginType,
		"domainUser": conn.DomainUser,
		"homePath":   conn.HomePath,
		"request":    requestPath,
	}).Warn("可疑路径尝试：路径超出允许范围")
	return "", fmt.Errorf("路径超出允许范围: %s", requestPath)
}

// ValidateFileName 校验文件名/目录名合法性，防止路径穿越
// 拒绝：空名称、"."、".."、包含路径分隔符（/ 或 \）的名称
func ValidateFileName(name string) error {
	if name == "" {
		return fmt.Errorf("名称不能为空")
	}
	if name == "." || name == ".." {
		return fmt.Errorf("名称不能为 . 或 ..")
	}
	if strings.ContainsAny(name, "/\\") {
		return fmt.Errorf("名称不能包含路径分隔符")
	}
	return nil
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
