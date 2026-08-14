package controller

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"sftpbackend/config"
	jwtpkg "sftpbackend/jwt"
	"sftpbackend/models"
	"sftpbackend/tools"
	"sftpbackend/utils"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 从请求中获取SFTP连接Token（优先从Header获取，也可从Query/Body获取）
func getSFTPToken(c *gin.Context) string {
	// 1. 优先从Header获取
	token := c.GetHeader("X-SFTP-Token")
	if token != "" {
		return token
	}

	// 2. 备选：从Query参数获取
	token = c.Query("sftp_token")
	if token != "" {
		return token
	}

	// 3. 备选：从Form/Body获取（根据前端传参方式调整）
	token = c.PostForm("sftp_token")
	return token
}

// ! LoginSFTP 处理SFTP登录（返回Token）
func LoginSFTP(c *gin.Context) {
	var conn *utils.SFTPConnection
	var err error

	// （密钥登录）获取上传的密钥文件
	file, fileErr := c.FormFile("file")
	if fileErr == nil && file != nil {
		srcFile, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "无法打开上传的密钥文件: " + err.Error(),
			})
			return
		}
		defer srcFile.Close()

		username := c.PostForm("username")
		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "用户名不能为空",
			})
			return
		}

		keyBytes, err := io.ReadAll(srcFile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "无法读取密钥文件内容: " + err.Error(),
			})
			return
		}

		// 创建连接实例
		conn, err = utils.NewSFTPConnectionByKey(username, keyBytes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "<密钥错误> " + err.Error(),
			})
			return
		}
	} else {
		// （密码登录）
		var sftpLogin models.SFTPLogin
		if err := c.ShouldBindJSON(&sftpLogin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请求数据格式错误",
			})
			return
		}

		// 解密密码
		decryptedPassword, decerr := tools.DecryptPassword(sftpLogin.Password)
		if decerr != nil {
			recordSftpLog(c, sftpLogin.Username, "Login", "SFTP登录失败: 密码解密失败", "", "")
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "<密码错误> " + decerr.Error(),
			})
			return
		}

		// 标签上传/中国联通模块：根据数据库配置决定登录方式（本地/LDAP），并校验角色权限
		if sftpLogin.LoginType == "hotlabel" || sftpLogin.LoginType == "chinaunicom" {
			// 获取该模块的配置（登录方式 + 可登录角色）
			moduleConfig, cfgErr := models.GetSFTPModuleConfig(sftpLogin.LoginType)
			loginType := models.LoginTypeLDAP // 配置不存在时默认 LDAP（向后兼容）
			if cfgErr == nil && moduleConfig != nil {
				loginType = moduleConfig.LoginType
			}

			if loginType == models.LoginTypeLDAP {
				// LDAP 域控验证（安全组 ldap.sftp_security_group_dn）通过后，读取公共SFTP账号登录并绑定模块根路径
				// 1. LDAP 验证域账号 + 安全组成员身份
				ldapUserInfo, statusCode, ldapErr := models.AuthenticateLDAPWithGroup(
					sftpLogin.Username, decryptedPassword, config.GlobalConfig.LDAP.SftpSecurityGroupDN)
				if ldapErr != nil {
					recordSftpLog(c, sftpLogin.Username, "Login", "SFTP登录失败: 域控验证失败: "+ldapErr.Error(), "", "")
					c.JSON(http.StatusOK, gin.H{
						"code":    statusCode,
						"message": "域控验证失败: " + ldapErr.Error(),
					})
					return
				}

				// 1.5 角色白名单校验：域账号的安全组（memberOf）匹配角色，检查角色是否在模块 enabled_roles 中
				if moduleConfig != nil && !CheckLDAPRolePermission(ldapUserInfo["memberOf"], sftpLogin.LoginType) {
					recordSftpLog(c, sftpLogin.Username, "Login", "SFTP登录失败: 角色无权限", "", "")
					c.JSON(http.StatusForbidden, gin.H{
						"code":    403,
						"message": "您的角色无权登录该模块",
					})
					return
				}

				// 2. 读取公共SFTP账号建立连接，绑定模块根路径限制
				account := config.GlobalConfig.SftpAccount
				var rootPath string
				if sftpLogin.LoginType == "hotlabel" {
					rootPath = config.GlobalConfig.HotLabel.RootPath
				} else {
					rootPath = config.GlobalConfig.ChinaUnicom.RootPath
				}
				conn, err = utils.NewSFTPConnectionForModule(account.SFTPUsername, account.SFTPPassword, rootPath, sftpLogin.LoginType, sftpLogin.Username)
				if err != nil {
					recordSftpLog(c, sftpLogin.Username, "Login", "SFTP登录失败: 专用账号连接失败: "+err.Error(), "", "")
					c.JSON(http.StatusInternalServerError, gin.H{
						"code":    500,
						"message": "专用账号连接失败: " + err.Error(),
					})
					return
				}
			} else {
				// 本地登录：平台本地账号验证 → 公共SFTP账号连接并绑定模块根路径
				// 1. 平台本地账号验证（启用状态、bcrypt密码、失败锁定、密码过期检查）
				localUser, expired, localErr := models.AuthenticateLocal(sftpLogin.Username, decryptedPassword)
				if localErr != nil {
					recordSftpLog(c, sftpLogin.Username, "Login", "SFTP登录失败: 本地账号验证失败: "+localErr.Error(), "", "")
					c.JSON(http.StatusOK, gin.H{
						"code":    400,
						"message": "本地账号验证失败: " + localErr.Error(),
					})
					return
				}
				if expired {
					recordSftpLog(c, sftpLogin.Username, "Login", "SFTP登录失败: 密码已过期", "", "")
					c.JSON(http.StatusOK, gin.H{
						"code":    400,
						"message": "密码已过期，请先在平台修改密码",
					})
					return
				}

				// 1.5 需修改密码检查：签发受限改密 token，由前端弹出修改密码弹框（与平台登录一致）
				if localUser.MustChangePassword {
					changeToken, tokenErr := jwtpkg.GenerateLimitedToken(sftpLogin.Username)
					if tokenErr != nil {
						c.JSON(http.StatusInternalServerError, gin.H{
							"code":    500,
							"message": "生成改密凭证失败",
						})
						return
					}
					recordSftpLog(c, sftpLogin.Username, "Login", "SFTP登录需先修改密码", "", "")
					c.JSON(http.StatusOK, gin.H{
						"code":    200,
						"message": "该账号需先修改密码",
						"data": gin.H{
							"must_change_password": true,
							"change_token":         changeToken,
						},
					})
					return
				}

				// 2. 角色白名单校验：平台本地账号的角色必须在模块 enabled_roles 中
				if moduleConfig != nil && (localUser.RoleID == nil || !CheckRolePermission(*localUser.RoleID, sftpLogin.LoginType)) {
					recordSftpLog(c, sftpLogin.Username, "Login", "SFTP登录失败: 角色无权限", "", "")
					c.JSON(http.StatusForbidden, gin.H{
						"code":    403,
						"message": "您的角色无权登录该模块",
					})
					return
				}

				// 3. 公共SFTP账号连接，绑定模块根路径，DomainUser = 平台账号
				account := config.GlobalConfig.SftpAccount
				var rootPath string
				if sftpLogin.LoginType == "hotlabel" {
					rootPath = config.GlobalConfig.HotLabel.RootPath
				} else {
					rootPath = config.GlobalConfig.ChinaUnicom.RootPath
				}
				conn, err = utils.NewSFTPConnectionForModule(account.SFTPUsername, account.SFTPPassword, rootPath, sftpLogin.LoginType, sftpLogin.Username)
				if err != nil {
					recordSftpLog(c, sftpLogin.Username, "Login", "SFTP登录失败: 专用账号连接失败: "+err.Error(), "", "")
					c.JSON(http.StatusInternalServerError, gin.H{
						"code":    500,
						"message": "专用账号连接失败: " + err.Error(),
					})
					return
				}
			}
		} else {
			// 普通密码登录：创建连接实例
			conn, err = utils.NewSFTPConnection(sftpLogin.Username, decryptedPassword)
			if err != nil {
				recordSftpLog(c, sftpLogin.Username, "Login", "SFTP登录失败: "+err.Error(), "", "")
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "<密码错误> " + err.Error(),
				})
				return
			}
		}
	}

	// 将连接添加到管理器，获取唯一Token
	token := utils.SFTPConnManager.AddConn(conn)

	// 登录成功，返回Token给前端
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "SFTP登录成功",
		"data": gin.H{
			"sftp_token":  token, // 前端需要保存这个Token
			"expire_tips": "连接有效期8小时，超时需重新登录",
		},
	})

	// 记录SFTP登录成功日志（独立逻辑，不影响主流程）
	recordSftpLog(c, sftpLogUsername(conn), "Login", "SFTP登录成功", conn.HomePath, "")
}

// 获取SFTP日志的操作者用户名（域账号优先，其次SFTP账号）
func sftpLogUsername(conn *utils.SFTPConnection) string {
	if conn == nil {
		return ""
	}
	if conn.DomainUser != "" {
		return conn.DomainUser
	}
	return conn.Username
}

// 从请求中获取双控复核人账号
func getReviewer(c *gin.Context) string {
	return utils.DualAuthManager.GetReviewer(c.GetHeader("X-Dual-Token"))
}

// 记录SFTP登录与操作日志（独立逻辑，不影响主流程）
func recordSftpLog(c *gin.Context, username, action, message, path, reviewer string) {
	if username == "" {
		return
	}
	log := models.SftpLog{
		Username: username,
		Reviewer: reviewer,
		IP:       c.ClientIP(),
		Action:   action,
		Message:  message,
		Path:     path,
	}
	if err := log.CreateSftpLog(); err != nil {
		fmt.Println("SFTP日志创建失败:", err)
	}
}

// ! DualVerify 双控验证：验证另一个产业部账号，通过后签发短期双控凭证
// 仅中国联通登录的连接需要双控验证；验证账号须属于产业部安全组且与当前登录账号不同
func DualVerify(c *gin.Context) {
	// 1. 获取Token与连接
	token := getSFTPToken(c)
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "SFTP Token不能为空",
		})
		return
	}
	conn, err := utils.SFTPConnManager.GetConn(token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    50014,
			"message": "SFTP连接失效: " + err.Error(),
		})
		return
	}

	// 2. 仅中国联通连接需要双控验证
	if conn.LoginType != "chinaunicom" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "当前连接无需双控验证",
		})
		return
	}

	// 3. 绑定请求参数
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求数据格式错误",
		})
		return
	}
	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "请输入双控验证账号和密码",
		})
		return
	}

	// 4. 解密密码
	decryptedPassword, decerr := tools.DecryptPassword(req.Password)
	if decerr != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "密码解密失败: " + decerr.Error(),
		})
		return
	}

	// 5. 双控账号不得与当前登录账号相同（大小写不敏感）
	if strings.EqualFold(req.Username, conn.DomainUser) {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "双控验证账号不能与当前登录账号相同",
		})
		return
	}

	// 6. LDAP 验证双控账号（须属于产业部安全组）
	_, statusCode, ldapErr := models.AuthenticateLDAPWithGroup(
		req.Username, decryptedPassword, config.GlobalConfig.LDAP.SftpSecurityGroupDN)
	if ldapErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    statusCode,
			"message": "双控验证失败: " + ldapErr.Error(),
		})
		return
	}

	// 7. 签发双控凭证（60秒有效，可复用），记录复核人
	dualToken := utils.DualAuthManager.IssueToken(token, req.Username)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "双控验证通过",
		"data": gin.H{
			"dual_token": dualToken,
		},
	})
}

// ! LogoutSFTP 处理SFTP登出（删除Token和连接）
func LogoutSFTP(c *gin.Context) {
	token := getSFTPToken(c)
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "SFTP Token不能为空",
		})
		return
	}

	// 登出前获取连接信息（用于记录日志）
	var username string
	if conn, err := utils.SFTPConnManager.GetConn(token); err == nil {
		username = sftpLogUsername(conn)
	}

	// 从管理器中删除并关闭连接
	if err := utils.SFTPConnManager.RemoveConn(token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "登出失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "SFTP登出成功",
	})

	// 记录SFTP登出日志（独立逻辑，不影响主流程）
	recordSftpLog(c, username, "Logout", "SFTP登出", "", "")
}

// ! 获取目录下的文件列表（通过Token获取连接）
func GetFiles(c *gin.Context) {
	// 1. 获取Token
	token := getSFTPToken(c)
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "SFTP Token不能为空",
		})
		return
	}

	// 2. 通过Token获取对应的连接实例
	conn, err := utils.SFTPConnManager.GetConn(token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    50014,
			"message": "SFTP连接失效: " + err.Error(),
		})
		return
	}

	// 3. 执行具体操作（使用当前连接实例）
	path := c.Query("path")
	if path == "" {
		path = "/"
	}

	// 校验路径不超出连接允许的根路径
	path, err = conn.ResolvePath(path)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": err.Error(),
		})
		return
	}

	files, err := conn.SftpClient.ReadDir(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "无法读取目录: " + err.Error(),
		})
		return
	}

	var dirs []models.FileInfo
	var fileList []models.FileInfo
	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			continue // 跳过隐藏项，不加入结果集
		}
		fileInfo := models.FileInfo{
			Name:     file.Name(),
			Path:     filepath.Join(path, file.Name()),
			IsDir:    file.IsDir(),
			Size:     file.Size(),
			Modified: file.ModTime().Format(time.RFC3339),
		}
		if file.IsDir() {
			dirs = append(dirs, fileInfo)

		} else {
			fileList = append(fileList, fileInfo)
		}
	}
	// 按照Modified排序
	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i].Modified > dirs[j].Modified
	})
	sort.Slice(fileList, func(i, j int) bool {
		return fileList[i].Modified > fileList[j].Modified
	})

	dirsCount := len(dirs)
	filesCount := len(fileList)
	var description string
	if dirsCount == 0 && filesCount == 0 {
		description = "空目录"
	} else if dirsCount == 0 && filesCount > 0 {
		description = fmt.Sprintf("%d 个文件", filesCount)
	} else if dirsCount > 0 && filesCount == 0 {
		description = fmt.Sprintf("%d 个目录", dirsCount)
	} else {
		description = fmt.Sprintf("%d 个文件 和 %d 个目录", filesCount, dirsCount)
	}

	sortedFiles := append(dirs, fileList...)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "文件列表获取成功",
		"data": gin.H{
			"path":        path,
			"description": description,
			"files":       sortedFiles,
		},
	})
}

// ! UploadFile 上传文件（适配Token，流式直写SFTP，避免整体缓冲导致进度虚高）
func UploadFile(c *gin.Context) {
	// 1. 获取Token（仅从 Header/Query 获取，避免解析请求体破坏流式读取）
	token := c.GetHeader("X-SFTP-Token")
	if token == "" {
		token = c.Query("sftp_token")
	}
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "SFTP Token不能为空",
		})
		return
	}

	// 2. 获取连接实例
	conn, err := utils.SFTPConnManager.GetConn(token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    50014,
			"message": "SFTP连接失效: " + err.Error(),
		})
		return
	}

	// 3. 流式读取 multipart 请求体（不能用 FormFile/PostForm，否则会先把整个文件缓冲到本地）
	mr, err := c.Request.MultipartReader()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "文件上传失败: " + err.Error(),
		})
		return
	}

	var targetPath string
	var uploadedFile string // 记录最后上传的文件完整路径（用于日志）
	uploaded := false
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "文件上传失败: " + err.Error(),
			})
			return
		}

		switch part.FormName() {
		case "path":
			data, err := io.ReadAll(part)
			part.Close()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    400,
					"message": "读取目标路径失败: " + err.Error(),
				})
				return
			}
			targetPath = string(data)
			// 校验目标路径不超出连接允许的根路径
			targetPath, err = conn.ResolvePath(targetPath)
			if err != nil {
				c.JSON(http.StatusForbidden, gin.H{
					"code":    403,
					"message": err.Error(),
				})
				return
			}
		case "file":
			if targetPath == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    400,
					"message": "目标路径不能为空",
				})
				return
			}
			filename := part.FileName()
			if filename == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    400,
					"message": "上传文件名为空",
				})
				return
			}
			dstPath := filepath.Join(targetPath, filename)
			// 边接收边写入SFTP，SFTP写入速度通过TCP背压传导到浏览器，进度条反映真实进度
			if err := conn.CreateUploadFile(dstPath, part); err != nil {
				part.Close()
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "文件上传失败: " + err.Error(),
				})
				return
			}
			part.Close()
			uploadedFile = dstPath
			uploaded = true
		default:
			// 忽略其他字段，读完丢弃以保持 multipart 流推进
			io.Copy(io.Discard, part)
			part.Close()
		}
	}

	if !uploaded {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "未接收到上传文件",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "文件上传成功",
	})

	// 记录SFTP上传日志（独立逻辑，不影响主流程）
	recordSftpLog(c, sftpLogUsername(conn), "Upload", "上传文件", uploadedFile, getReviewer(c))
}

// ! CreateFolder 创建目录（适配Token）
func CreateFolder(c *gin.Context) {
	// 1. 获取Token
	token := getSFTPToken(c)
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "SFTP Token不能为空",
		})
		return
	}

	// 2. 获取连接实例
	conn, err := utils.SFTPConnManager.GetConn(token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    50014,
			"message": "SFTP连接失效: " + err.Error(),
		})
		return
	}

	// 3. 执行创建目录
	var req struct {
		Path string `json:"path"`
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求数据格式错误",
		})
		return
	}

	if req.Path == "" || req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "路径和名称不能为空",
		})
		return
	}

	// 校验路径不超出连接允许的根路径
	req.Path, err = conn.ResolvePath(req.Path)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": err.Error(),
		})
		return
	}

	if err := conn.CreateFolder(req.Path, req.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建目录失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "目录创建成功",
	})

	// 记录SFTP创建目录日志（独立逻辑，不影响主流程）
	recordSftpLog(c, sftpLogUsername(conn), "Mkdir", "创建目录: "+req.Name, filepath.Join(req.Path, req.Name), getReviewer(c))
}

// ! DownloadFile 下载文件（适配Token）
func DownloadFile(c *gin.Context) {
	// 1. 获取Token
	token := getSFTPToken(c)
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "SFTP Token不能为空",
		})
		return
	}

	// 2. 获取连接实例
	conn, err := utils.SFTPConnManager.GetConn(token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    50014,
			"message": "SFTP连接失效: " + err.Error(),
		})
		return
	}

	// 3. 执行下载操作
	filePath := c.Query("path")
	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "文件路径不能为空",
		})
		return
	}

	// 校验路径不超出连接允许的根路径
	filePath, err = conn.ResolvePath(filePath)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": err.Error(),
		})
		return
	}

	// 检查文件是否存在
	fileInfo, err := conn.SftpClient.Stat(filePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "文件不存在: " + err.Error(),
		})
		return
	}

	// 4. 检查是否为目录
	if fileInfo.IsDir() {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "路径是目录不是文件",
		})
		return
	}

	// 5. 打开文件并传输
	file, err := conn.SftpClient.Open(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "打开文件失败: " + err.Error(),
		})
		return
	}
	defer file.Close()

	// fileName := filepath.Base(filePath)
	fileName := fileInfo.Name() // 获取SFTP服务器上的文件名
	// 二进制流文件，让浏览器识别为下载文件
	c.Header("Content-Type", "application/octet-stream")
	// c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"; filename*=UTF-8''%s",
		escapeFileName(fileName),  // 兼容低版本浏览器的文件名转义
		url.QueryEscape(fileName), // 标准UTF-8转义，解决中文文件名乱码
	))
	// c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	c.Header("Content-Length", strconv.FormatInt(fileInfo.Size(), 10)) // 文件大小，让浏览器显示下载进度
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")   // 禁止缓存，避免大文件缓存问题
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.Status(http.StatusOK) // 设置HTTP状态码为200

	// if _, err := io.Copy(c.Writer, file); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"code":    500,
	// 		"message": "文件传输失败: " + err.Error(),
	// 	})
	// 	return
	// }
	// 流式传输（SFTP流 → HTTP响应流）
	// io.Copy会自动分块复制数据（默认缓冲区32KB），不会将整个文件加载到后端内存
	// 直接将SFTP读取的字节流式写入gin的ResponseWriter，前端浏览器实时接收
	_, err = io.Copy(c.Writer, file)
	if err != nil {
		// 流式传输过程中出错，记录日志即可（浏览器会自动中断下载）
		log.Printf("SFTP文件流式下载失败，路径：%s，错误：%v", filePath, err)
		return
	}

	// 记录SFTP下载日志（独立逻辑，不影响主流程）
	recordSftpLog(c, sftpLogUsername(conn), "Download", "下载文件", filePath, getReviewer(c))
}

// 辅助函数：转义文件名，兼容IE/Edge等低版本浏览器（解决中文乱码）
func escapeFileName(name string) string {
	// 替换特殊字符，避免响应头解析错误
	replacer := strings.NewReplacer(`"`, `\"`, `\`, `\\`, `/`, `\/`)
	return replacer.Replace(name)
}

// ! DeletePath 删除文件/目录（适配Token）
func DeletePath(c *gin.Context) {
	// 1. 获取Token
	token := getSFTPToken(c)
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "SFTP Token不能为空",
		})
		return
	}

	// 2. 获取连接实例
	conn, err := utils.SFTPConnManager.GetConn(token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    50014,
			"message": "SFTP连接失效: " + err.Error(),
		})
		return
	}

	// 3. 执行删除操作
	var req struct {
		Path string `json:"path"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求数据格式错误",
		})
		return
	}

	if req.Path == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "路径不能为空",
		})
		return
	}

	// 校验路径不超出连接允许的根路径
	req.Path, err = conn.ResolvePath(req.Path)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": err.Error(),
		})
		return
	}

	fileInfo, err := conn.SftpClient.Stat(req.Path)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "路径不存在: " + err.Error(),
		})
		return
	}

	var errDelete error
	if fileInfo.IsDir() {
		errDelete = conn.DeleteDirectory(req.Path)
	} else {
		errDelete = conn.SftpClient.Remove(req.Path)
	}

	if errDelete != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除失败: " + errDelete.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})

	// 记录SFTP删除日志（独立逻辑，不影响主流程）
	recordSftpLog(c, sftpLogUsername(conn), "Delete", "删除", req.Path, getReviewer(c))
}

// ! 批量删除目录和文件
func BatchDelete(c *gin.Context) {
	// 获取Token
	token := getSFTPToken(c)
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "SFTP Token不能为空",
		})
		return
	}

	// 获取连接实例
	conn, err := utils.SFTPConnManager.GetConn(token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    50014,
			"message": "SFTP连接失效: " + err.Error(),
		})
		return
	}

	var req []models.FileInfo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求数据格式错误",
		})
		return
	}

	if len(req) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "路径不能为空",
		})
		return
	}

	// 校验所有路径是否存在
	type pathCheckResult struct {
		Path     string // 路径
		IsDir    bool   // 是否目录
		NotExist bool   // 是否不存在
		Error    string // 错误信息
	}
	var checkResults []pathCheckResult

	for _, file := range req {
		path := file.Path
		// 校验路径不超出连接允许的根路径，任一越界整体拒绝
		path, err = conn.ResolvePath(path)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": err.Error(),
			})
			return
		}
		fileInfo, err := conn.SftpClient.Stat(path)

		// 路径不存在：记录错误，继续下一个
		if err != nil {
			checkResults = append(checkResults, pathCheckResult{
				Path:     path,
				NotExist: true,
				Error:    "路径不存在: " + err.Error(),
			})
			continue
		}

		// 路径存在：记录类型（文件/目录）
		checkResults = append(checkResults, pathCheckResult{
			Path:  path,
			IsDir: fileInfo.IsDir(),
		})
	}

	// 收集所有校验失败的路径
	var validateErrors []string
	for _, res := range checkResults {
		if res.NotExist {
			validateErrors = append(validateErrors, res.Error)
		}
	}

	// 如果有校验失败，直接返回所有错误，不执行删除
	if len(validateErrors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "部分路径校验失败",
			"errors":  validateErrors,
		})
		return
	}

	// 所有路径都存在，执行批量删除（容错不中断）
	var successCount int      // 成功删除数量
	var deleteErrors []string // 删除失败的错误信息

	for _, res := range checkResults {
		path := res.Path
		var errDel error

		// 目录 / 文件分别删除
		if res.IsDir {
			errDel = conn.DeleteDirectory(path)
		} else {
			errDel = conn.SftpClient.Remove(path)
		}

		// 删除失败：记录错误，继续下一个
		if errDel != nil {
			deleteErrors = append(deleteErrors, "删除失败: "+errDel.Error())
			continue
		}

		successCount++
	}

	if len(deleteErrors) == 0 {
		// 全部删除成功
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "批量删除成功",
		})
	} else if successCount == 0 {
		// 全部删除失败（重点补充）
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "批量删除失败:" + strings.Join(deleteErrors, ", "),
		})
	} else {
		// 部分删除成功，部分失败
		c.JSON(http.StatusPartialContent, gin.H{
			"code":    206,
			"message": "部分删除成功:" + strings.Join(deleteErrors, ", "),
		})
	}

	// 记录SFTP批量删除日志（独立逻辑，不影响主流程；有成功删除时记录）
	if successCount > 0 {
		var paths []string
		for _, res := range checkResults {
			paths = append(paths, res.Path)
		}
		recordSftpLog(c, sftpLogUsername(conn), "BatchDelete", "批量删除"+fmt.Sprintf("（成功%d条）", successCount), strings.Join(paths, ","), getReviewer(c))
	}
}

// ! RenamePath 重命名（适配Token）
func RenamePath(c *gin.Context) {
	// 1. 获取Token
	token := getSFTPToken(c)
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "SFTP Token不能为空",
		})
		return
	}

	// 2. 获取连接实例
	conn, err := utils.SFTPConnManager.GetConn(token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    50014,
			"message": "SFTP连接失效: " + err.Error(),
		})
		return
	}

	// 3. 执行重命名
	var req struct {
		OldPath string `json:"oldPath"`
		NewName string `json:"newName"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求数据格式错误",
		})
		return
	}

	if req.OldPath == "" || req.NewName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "原路径和新名称不能为空",
		})
		return
	}

	if strings.Contains(req.NewName, "/") || req.NewName == "." || req.NewName == ".." {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "名称包含非法字符",
		})
		return
	}

	// 校验原路径不超出连接允许的根路径
	req.OldPath, err = conn.ResolvePath(req.OldPath)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": err.Error(),
		})
		return
	}

	oldDir := filepath.Dir(req.OldPath)
	newPath := filepath.Join(oldDir, req.NewName)

	if err := conn.SftpClient.Rename(req.OldPath, newPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "重命名失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "重命名成功",
	})

	// 记录SFTP重命名日志（独立逻辑，不影响主流程）
	recordSftpLog(c, sftpLogUsername(conn), "Rename", "重命名为: "+req.NewName, req.OldPath, getReviewer(c))
}

// ! DownloadDirectory 下载目录
func DownloadDirectory(c *gin.Context) {
	// 获取Token
	token := getSFTPToken(c)
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "SFTP Token不能为空",
		})
		return
	}

	// 获取连接实例
	conn, err := utils.SFTPConnManager.GetConn(token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    50014,
			"message": "SFTP连接失效: " + err.Error(),
		})
		return
	}

	// 获取要下载的远程目录路径
	remoteDir := c.Query("path")
	if remoteDir == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "目录路径不能为空",
		})
		return
	}

	// 校验路径不超出连接允许的根路径
	remoteDir, err = conn.ResolvePath(remoteDir)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": err.Error(),
		})
		return
	}

	// 校验路径是否存在且是目录
	fileInfo, err := conn.SftpClient.Stat(remoteDir)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "目录不存在: " + err.Error(),
		})
		return
	}
	if !fileInfo.IsDir() {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "该路径不是目录，请使用文件下载接口",
		})
		return
	}

	// 设置下载响应头
	zipName := filepath.Base(remoteDir) + ".zip"
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", fmt.Sprintf(
		"attachment; filename=\"%s\"; filename*=UTF-8''%s",
		escapeFileName(zipName),
		url.QueryEscape(zipName),
	))
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.Status(http.StatusOK)

	// 创建 zip 流式写入器
	zipWriter := zip.NewWriter(c.Writer)
	defer zipWriter.Close()

	// 递归遍历并压缩目录
	err = addDirToZip(zipWriter, conn, remoteDir, remoteDir)
	if err != nil {
		log.Printf("目录打包下载失败：%s，错误：%v", remoteDir, err)
		return
	}

	// 记录SFTP下载目录日志（独立逻辑，不影响主流程）
	recordSftpLog(c, sftpLogUsername(conn), "Download", "下载目录", remoteDir, getReviewer(c))

	log.Printf("目录打包下载完成：%s", remoteDir)
}

// --- DownloadDirectory辅助函数：递归把SFTP目录写入zip ---
func addDirToZip(zipW *zip.Writer, conn *utils.SFTPConnection, remoteRoot, currentDir string) error {
	// 读取当前目录内容
	files, err := conn.SftpClient.ReadDir(currentDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		// 跳过隐藏文件（和你文件列表逻辑一致）
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}

		remotePath := filepath.Join(currentDir, file.Name())
		// zip 内部路径（保持目录结构）
		zipPath, _ := filepath.Rel(remoteRoot, remotePath)

		if file.IsDir() {
			// 递归子目录
			if err := addDirToZip(zipW, conn, remoteRoot, remotePath); err != nil {
				return err
			}
		} else {
			// 写入文件到 zip
			if err := writeFileToZip(zipW, conn, remotePath, zipPath); err != nil {
				return err
			}
		}
	}
	return nil
}

// 写入单个文件到 zip
func writeFileToZip(zipW *zip.Writer, conn *utils.SFTPConnection, remotePath, zipPath string) error {
	// 打开远程文件
	srcFile, err := conn.SftpClient.Open(remotePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 创建 zip 中的文件
	zipFile, err := zipW.Create(zipPath)
	if err != nil {
		return err
	}

	// 流式复制（不占内存）
	_, err = io.Copy(zipFile, srcFile)
	return err
}

// ! SearchFiles 递归搜索文件和目录（不区分大小写模糊匹配）
func SearchFiles(c *gin.Context) {
	// 1. 获取Token
	token := getSFTPToken(c)
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "SFTP Token不能为空",
		})
		return
	}

	// 2. 通过Token获取对应的连接实例
	conn, err := utils.SFTPConnManager.GetConn(token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    50014,
			"message": "SFTP连接失效: " + err.Error(),
		})
		return
	}

	// 3. 获取搜索参数
	searchPath := c.Query("path")
	if searchPath == "" {
		searchPath = "/"
	}
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "搜索关键字不能为空",
		})
		return
	}

	// 校验路径不超出连接允许的根路径
	searchPath, err = conn.ResolvePath(searchPath)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": err.Error(),
		})
		return
	}

	// 4. 验证搜索路径是否存在
	fileInfo, err := conn.SftpClient.Stat(searchPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "搜索路径不存在: " + err.Error(),
		})
		return
	}
	if !fileInfo.IsDir() {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "搜索路径不是目录",
		})
		return
	}

	// 5. 递归搜索
	var results []models.SearchFileInfo
	recursiveSearch(conn, searchPath, keyword, &results)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "搜索完成",
		"data": gin.H{
			"total":   len(results),
			"results": results,
		},
	})
}

// recursiveSearch 递归搜索目录，匹配文件名（不区分大小写）
func recursiveSearch(conn *utils.SFTPConnection, dirPath, keyword string, results *[]models.SearchFileInfo) error {
	files, err := conn.SftpClient.ReadDir(dirPath)
	if err != nil {
		// 跳过无权限或无法读取的目录
		return nil
	}

	lowerKeyword := strings.ToLower(keyword)

	for _, f := range files {
		// 跳过隐藏文件/目录
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}

		fullPath := filepath.Join(dirPath, f.Name())

		// 模糊匹配文件名（不区分大小写）
		if strings.Contains(strings.ToLower(f.Name()), lowerKeyword) {
			*results = append(*results, models.SearchFileInfo{
				FileInfo: models.FileInfo{
					Name:     f.Name(),
					Path:     fullPath,
					IsDir:    f.IsDir(),
					Size:     f.Size(),
					Modified: f.ModTime().Format(time.RFC3339),
				},
				ParentPath: dirPath,
			})
		}

		// 递归进入子目录
		if f.IsDir() {
			recursiveSearch(conn, fullPath, keyword, results)
		}
	}
	return nil
}
