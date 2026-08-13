package models

// SFTPUser结构体用于存储SFTP用户信息
type SFTPLogin struct {
	Username  string `json:"username"`  // 用户名
	Password  string `json:"password"`  // 密码
	LoginType string `json:"loginType"` // 登录模块标识，值为 "hotlabel"/"chinaunicom" 时走域控验证流程（公共账号登录并绑定模块根路径）
}

// 文件信息结构体
type FileInfo struct {
	Name     string `json:"name"`     // 文件名
	Path     string `json:"path"`     // 完整路径
	IsDir    bool   `json:"isDir"`    // 是否为目录
	Size     int64  `json:"size"`     // 文件大小字节
	Modified string `json:"modified"` // 修改时间
}

// 搜索结果文件信息结构体（包含父目录路径）
type SearchFileInfo struct {
	FileInfo
	ParentPath string `json:"parentPath"` // 父目录路径
}
