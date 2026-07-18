package models

// SFTPUser结构体用于存储SFTP用户信息
type SFTPLogin struct {
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
}

// 文件信息结构体
type FileInfo struct {
	Name     string `json:"name"`     // 文件名
	Path     string `json:"path"`     // 完整路径
	IsDir    bool   `json:"isDir"`    // 是否为目录
	Size     int64  `json:"size"`     // 文件大小字节
	Modified string `json:"modified"` // 修改时间
}
