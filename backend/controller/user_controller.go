package controller

import (
	"fmt"
	"net/http"
	"os"
	"sftpbackend/config"
	"sftpbackend/models"
	"sftpbackend/tools"
	"sftpbackend/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// ! 添加用户,需要字段{"name":"string","password":"string"}
func AddUser(c *gin.Context) {
	var user models.SftpUsers
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	// 解密密码
	decryptedPassword, err := tools.DecryptPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	user.Password = decryptedPassword
	// 保存原始密码,用于附件发送
	originalPassword := user.Password
	// ! 查询用户是否已经存在
	if err := user.ExistOrNot(user.Name); err == nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "账号已存在",
		})
		return
	} else {
		if err := user.AddUser(c); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		} else {
			// 判断是否需要创建附件
			if user.EmailOrNot {
				// 创建附件
				if err := utils.CreateAttachment(user.Name, originalPassword, user.LoginType); err != nil {
					fmt.Println("附件创建失败:", err)
					c.JSON(http.StatusOK, gin.H{
						"code":    200,
						"data":    "",
						"message": "账号添加成功,但附件创建失败",
					})
					return
				}
			}
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"data":    "",
				"message": "账号添加成功",
			})
			// 创建操作日志
			log := models.Log{
				Username: c.MustGet("username").(string),
				IP:       c.ClientIP(),
				Action:   "Add",
				Message:  "添加用户: " + user.Name,
			}
			if err := log.CreateLog(); err != nil {
				fmt.Println("日志创建失败:", err)
			}
		}
	}
}

// ! 获取所有用户列表
func GetUserList(c *gin.Context) {
	// 获取路径参数中的页码和每页记录数
	pageStr := c.Param("page")
	limitStr := c.Param("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "无效的页码参数",
		})
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "无效的每页记录数参数",
		})
		return
	}

	// 获取查询参数中的用户名
	username := c.Query("username")
	var user models.SftpUsers
	if users, total, err := user.GetUserList(page, limit, username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": gin.H{
				"users": users,
				"total": total,
			},
		})
	}
}

// ! 更新一个用户密码
func UpdateUser(c *gin.Context) {
	var user models.SftpUsers
	c.ShouldBind(&user)
	// 解密密码
	decryptedPassword, err := tools.DecryptPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	user.Password = decryptedPassword
	// 保存原始密码,用于附件发送
	originalPassword := user.Password
	// ! 查询用户是否已经存在
	if err := user.ExistOrNot(user.Name); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "账号不存在",
		})
		return
	} else {
		if err := user.UpdateUser(c); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"message": err.Error(),
			})
			return
		} else {
			// 判断是否需要创建附件
			if user.EmailOrNot {
				// 创建附件
				if err := utils.CreateAttachment(user.Name, originalPassword, user.LoginType); err != nil {
					fmt.Println("附件创建失败:", err)
					c.JSON(http.StatusOK, gin.H{
						"code":    200,
						"data":    "",
						"message": "密码修改成功,但附件创建失败",
					})
					return
				}
			}
			var message string
			switch user.LoginType {
			case "Password":
				message = "密码修改成功"
			case "KeyFile":
				message = "密钥修改成功"
			case "both":
				message = "密码和密钥修改成功"
			}

			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"data":    "",
				"message": message,
			})
			// 创建操作日志
			log := models.Log{
				Username: c.MustGet("username").(string),
				IP:       c.ClientIP(),
				Action:   "Update",
				Message:  user.Name + " " + message,
			}
			if err := log.CreateLog(); err != nil {
				fmt.Println("日志创建失败:", err)
			}
		}
	}
}

// ! 删除一个用户
func DeleteUser(c *gin.Context) {
	// 获取路径上的id参数
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "无效的ID",
		})
		return
	}

	var user models.SftpUsers
	// 获取将被删除的用户信息,因为日志需要记录被删除的用户信息
	if err := user.GetUserInfoById(id); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
		})
		return
	}
	// 删除用户
	if err := user.DeleteUser(user.Name); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
		})
		return
	} else {
		// 返回成功响应
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"data":    gin.H{},
			"message": "删除成功",
		})
		// 创建操作日志
		log := models.Log{
			Username: c.MustGet("username").(string),
			IP:       c.ClientIP(),
			Action:   "Delete",
			Message:  "删除用户: " + user.Name,
		}
		if err := log.CreateLog(); err != nil {
			fmt.Println("日志创建失败:", err)
		}
	}
}

// ! 批量删除用户
func DeleteUsers(c *gin.Context) {
	var userIdList []int64
	if err := c.BindJSON(&userIdList); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	var user models.SftpUsers
	// 获取将被删除的用户信息,因为日志需要记录被删除的用户信息
	users, err := user.GetUsersByIds(userIdList)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	// 循环users获取全部的用户名
	var userNames []string
	for _, user := range users {
		userNames = append(userNames, user.Name)
	}

	// 删除用户
	if err := user.DeleteUsers(userNames); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"data":    gin.H{},
			"message": "删除成功",
		})

		// 将切片转换为字符串
		userNamesStr := strings.Join(userNames, ",")
		// 创建操作日志
		log := models.Log{
			Username: c.MustGet("username").(string),
			IP:       c.ClientIP(),
			Action:   "Delete",
			Message:  "批量删除: " + userNamesStr,
		}
		if err := log.CreateLog(); err != nil {
			fmt.Println("日志创建失败:", err)
		}
	}
}

// ! 邮件发送
func SendEmail(c *gin.Context) {
	var email models.EmailConfig

	if err := c.BindJSON(&email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	// 密钥附件
	keyAttachmentPath := email.Subject + "_sftp_rsa_key"
	// 邮件主题
	email.Subject = config.GlobalConfig.Email.Subject + "-" + email.Subject
	// 账号信息附件
	email.Path = email.Subject + ".txt"
	// 附件列表
	Attachments := []string{email.Path, keyAttachmentPath}
	if err := email.SendEmail(Attachments); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"data":    gin.H{},
			"message": "发送成功",
		})
		// 邮件发送完成后删除附件
		os.Remove(email.Path)
		os.Remove(keyAttachmentPath)
	}
}

// ! 下载私钥
type DownloadKeyResponse struct {
	Message string `json:"message"`
}

func DownloadPrivateKey(c *gin.Context) {
	// 获取用户名参数
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "用户名不能为空",
		})
		return
	}

	// 私钥文件路径 (当前工作目录下的 {用户名}_sftp_rsa_key)
	privateKeyPath := username + "_sftp_rsa_key"

	// 检查私钥文件是否存在
	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "私钥文件不存在，请先创建或修改用户并选择导出私钥",
		})
		return
	}

	// 读取私钥文件内容
	fileContent, err := os.ReadFile(privateKeyPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "读取私钥文件失败：" + err.Error(),
		})
		return
	}

	// 以附件形式返回文件供下载
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s_sftp_rsa_key", username))
	
	// 在下载完成后自动清理私钥文件，提高安全性
	defer os.Remove(privateKeyPath)
	c.Data(http.StatusOK, "application/octet-stream", fileContent)
}
