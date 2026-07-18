package models

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sftpbackend/config"
	"sftpbackend/dao"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// User model

type UserList struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	Home            string `json:"home"`
	PasswordExpires string `json:"passwordExpires"`
	CreatedAt       string `json:"createdAt"`
}

type SftpUsers struct {
	gorm.Model
	Name            string `json:"name" binding:"required" form:"name"`
	LoginType       string `json:"loginType" gorm:"-" form:"loginType"`
	Password        string `json:"password" gorm:"-" form:"password"`
	Home            string `json:"home" form:"home"`
	PasswordExpires string `json:"passwordExpires" form:"passwordExpires"`
	EmailOrNot      bool   `json:"emailOrNot" gorm:"-" form:"emailOrNot"` // 是否需要发送邮件
	NoExpire        bool   `json:"noExpire" gorm:"-" form:"noExpire"`     // 是否设置密码不过期
	PublicKey       bool   `json:"publicKey" gorm:"-" form:"publicKey"`   // 是否上传客户提供的公钥

}

type LoginUser struct {
	Name      string `json:"name" binding:"required"`
	Password  string `json:"password" binding:"required"`
	LoginType string `json:"loginType" binding:"required"`
}

// ! User增删改查
// ! 新增一个用户  方法
func (u *SftpUsers) AddUser(c *gin.Context) (err error) {
	// 读取配置脚本路径
	scriptPath := config.GlobalConfig.Script.AdduserScript
	SetupKeyScript := config.GlobalConfig.Script.SetupKeyScript
	// 根据登录类型创建用户
	switch u.LoginType {
	case "Password":
		// 密码登录类型
		var cmd *exec.Cmd
		if u.NoExpire {
			cmd = exec.Command("bash", scriptPath, u.Name, u.Password, "99999")
		} else {
			cmd = exec.Command("bash", scriptPath, u.Name, u.Password)
		}
		// 执行命令并获取输出
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("执行命令失败:", err)
		}
		fmt.Println("命令输出:", string(output))
		return err

	case "KeyFile":
		tmpFilename := "Temp.pub"
		savePath := fmt.Sprintf("%s", tmpFilename)

		//  无论是否上传文件，出错都要删除临时文件
		var fileUploaded bool

		// 只要文件存在就删除
		cleanup := func() {
			if fileUploaded {
				_ = os.Remove(savePath)
			}
		}
		// defer 确保函数退出前一定执行清理（正常/异常退出都生效）
		defer cleanup()

		// 如果使用客户提供的公钥
		if u.PublicKey {
			keyfileHeader, err := c.FormFile("file")
			if err != nil {
				return err
			}
			// 临时保存上传的文件到当前目录

			// 打开文件
			file, err := keyfileHeader.Open()
			if err != nil {
				return err
			}
			defer file.Close()

			fileByte, err := io.ReadAll(file)
			if err != nil {
				return err
			}
			// 将文件内容写入到临时文件
			err = os.WriteFile(savePath, fileByte, 0644)
			if err != nil {
				return err
			}
			// 标记文件已上传
			fileUploaded = true
		}

		// 密钥登录类型
		var cmd *exec.Cmd
		if u.NoExpire {
			cmd = exec.Command("bash", scriptPath, u.Name, u.Password, "99999")
		} else {
			cmd = exec.Command("bash", scriptPath, u.Name, u.Password)
		}

		// 执行命令并获取输出
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("执行命令失败:", err)
			fmt.Println("命令输出:", string(output))
			return err
		}
		fmt.Println("命令输出:", string(output))

		// 创建密钥对
		cmdForKey := exec.Command("bash", SetupKeyScript, u.Name)
		outputKey, err := cmdForKey.CombinedOutput()
		if err != nil {
			fmt.Println("执行命令失败:", err)
			fmt.Println("命令输出:", string(outputKey))
		}
		fmt.Println("命令输出:", string(outputKey))
		return err

	case "both":
		tmpFilename := "Temp.pub"
		savePath := fmt.Sprintf("%s", tmpFilename)

		//  无论是否上传文件，出错都要删除临时文件
		var fileUploaded bool

		// 只要文件存在就删除
		cleanup := func() {
			if fileUploaded {
				_ = os.Remove(savePath)
			}
		}
		// defer 确保函数退出前一定执行清理（正常/异常退出都生效）
		defer cleanup()

		// 如果使用客户提供的公钥
		if u.PublicKey {
			keyfileHeader, err := c.FormFile("file")
			if err != nil {
				return err
			}
			// 临时保存上传的文件到当前目录

			// 打开文件
			file, err := keyfileHeader.Open()
			if err != nil {
				return err
			}
			defer file.Close()

			fileByte, err := io.ReadAll(file)
			if err != nil {
				return err
			}
			// 将文件内容写入到临时文件
			err = os.WriteFile(savePath, fileByte, 0644)
			if err != nil {
				return err
			}
			// 标记文件已上传
			fileUploaded = true
		}

		var cmd *exec.Cmd
		if u.NoExpire {
			cmd = exec.Command("bash", scriptPath, u.Name, u.Password, "99999")
		} else {
			cmd = exec.Command("bash", scriptPath, u.Name, u.Password)
		}
		// 执行命令并获取输出
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("执行命令失败:", err)
			fmt.Println("命令输出:", string(output))
			return err
		}
		fmt.Println("命令输出:", string(output))

		cmdForKey := exec.Command("bash", SetupKeyScript, u.Name)
		outputKey, err := cmdForKey.CombinedOutput()
		if err != nil {
			fmt.Println("执行命令失败:", err)
			fmt.Println("命令输出:", string(outputKey))
		}
		fmt.Println("命令输出:", string(outputKey))

		return err

	default:
		return fmt.Errorf("无效的登录类型: %s", u.LoginType)
	}

}

// ! 查询用户是否存在
func (u *SftpUsers) ExistOrNot(name string) (err error) {
	var user SftpUsers
	err = dao.DB.Where("name = ?", name).First(&user).Error
	return
}

// ! 获取所有用户列表
func (u *SftpUsers) GetUserList(page, limit int, username string) (userlist []UserList, totalCount int64, err error) {
	var users []SftpUsers
	// 根据页码和每页数量计算偏移量
	offset := (page - 1) * limit

	// 构建查询条件
	query := dao.DB.Offset(offset).Limit(limit)
	// 先构建一个用于统计总数的查询条件副本，不添加分页相关设置
	countQuery := dao.DB.Model(&SftpUsers{})

	// 如果有username查询参数，使用LIKE进行模糊查询
	if username != "" {
		likePattern := "%" + username + "%"
		// 带分页的条件
		query = query.Where("name LIKE?", likePattern)
		// 不带分页的条件
		countQuery = countQuery.Where("name LIKE?", likePattern)
	}

	// 查询满足条件的所有用户
	if err = query.Select("id,name,home,password_expires,created_at").Order("created_at ASC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	// 通过新的结构体整理返回的数据
	for _, user := range users {
		// 格式化时间
		userlist = append(userlist, UserList{
			ID:              user.ID,
			Name:            user.Name,
			Home:            user.Home,
			PasswordExpires: user.PasswordExpires,
			CreatedAt:       user.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	// 统计满足条件的用户总数（不分页）
	if err = countQuery.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	return
}

// ! 更新一个用户密码
func (u *SftpUsers) UpdateUser(c *gin.Context) (err error) {
	scriptPath := config.GlobalConfig.Script.UpdateuserScript
	SetupKeyScript := config.GlobalConfig.Script.SetupKeyScript
	// 根据登录类型创建用户
	switch u.LoginType {
	case "Password":
		// 密码登录类型
		var cmd *exec.Cmd
		if u.NoExpire {
			cmd = exec.Command("bash", scriptPath, u.Name, u.Password, "99999")
		} else {
			cmd = exec.Command("bash", scriptPath, u.Name, u.Password)
		}
		// 执行命令并获取输出
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("执行命令失败:", err)
		}
		fmt.Println("命令输出:", string(output))
		return err

	case "KeyFile":
		tmpFilename := "Temp.pub"
		savePath := fmt.Sprintf("%s", tmpFilename)

		//  无论是否上传文件，出错都要删除临时文件
		var fileUploaded bool

		// 只要文件存在就删除
		cleanup := func() {
			if fileUploaded {
				_ = os.Remove(savePath)
			}
		}
		// defer 确保函数退出前一定执行清理（正常/异常退出都生效）
		defer cleanup()

		// 如果使用客户提供的公钥
		if u.PublicKey {
			keyfileHeader, err := c.FormFile("file")
			if err != nil {
				return err
			}
			// 临时保存上传的文件到当前目录

			// 打开文件
			file, err := keyfileHeader.Open()
			if err != nil {
				return err
			}
			defer file.Close()

			fileByte, err := io.ReadAll(file)
			if err != nil {
				return err
			}
			// 将文件内容写入到临时文件
			err = os.WriteFile(savePath, fileByte, 0644)
			if err != nil {
				return err
			}
			// 标记文件已上传
			fileUploaded = true
		}

		// 密钥登录类型

		var cmd *exec.Cmd
		if u.NoExpire {
			cmd = exec.Command("bash", scriptPath, u.Name, u.Password, "99999")
		} else {
			cmd = exec.Command("bash", scriptPath, u.Name, u.Password)
		}

		// 执行命令并获取输出
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("执行命令失败:", err)
			fmt.Println("命令输出:", string(output))
			return err
		}
		fmt.Println("命令输出:", string(output))

		cmdForKey := exec.Command("bash", SetupKeyScript, u.Name)
		outputKey, err := cmdForKey.CombinedOutput()
		if err != nil {
			fmt.Println("执行命令失败:", err)
			fmt.Println("命令输出:", string(outputKey))
		}
		fmt.Println("命令输出:", string(outputKey))
		return err

	case "both":
		tmpFilename := "Temp.pub"
		savePath := fmt.Sprintf("%s", tmpFilename)

		//  无论是否上传文件，出错都要删除临时文件
		var fileUploaded bool

		// 只要文件存在就删除
		cleanup := func() {
			if fileUploaded {
				_ = os.Remove(savePath)
			}
		}
		// defer 确保函数退出前一定执行清理（正常/异常退出都生效）
		defer cleanup()

		// 如果使用客户提供的公钥
		if u.PublicKey {
			keyfileHeader, err := c.FormFile("file")
			if err != nil {
				return err
			}
			// 临时保存上传的文件到当前目录

			// 打开文件
			file, err := keyfileHeader.Open()
			if err != nil {
				return err
			}
			defer file.Close()

			fileByte, err := io.ReadAll(file)
			if err != nil {
				return err
			}
			// 将文件内容写入到临时文件
			err = os.WriteFile(savePath, fileByte, 0644)
			if err != nil {
				return err
			}
			// 标记文件已上传
			fileUploaded = true
		}

		var cmd *exec.Cmd
		if u.NoExpire {
			cmd = exec.Command("bash", scriptPath, u.Name, u.Password, "99999")
		} else {
			cmd = exec.Command("bash", scriptPath, u.Name, u.Password)
		}
		// 执行命令并获取输出
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("执行命令失败:", err)
			fmt.Println("命令输出:", string(output))
			return err
		}
		fmt.Println("命令输出:", string(output))

		cmdForKey := exec.Command("bash", SetupKeyScript, u.Name)
		outputKey, err := cmdForKey.CombinedOutput()
		if err != nil {
			fmt.Println("执行命令失败:", err)
			fmt.Println("命令输出:", string(outputKey))
		}
		fmt.Println("命令输出:", string(outputKey))

		return err

	default:
		return fmt.Errorf("无效的登录类型: %s", u.LoginType)
	}
}

// ! 删除一个用户
func (u *SftpUsers) DeleteUser(username string) (err error) {
	scriptPath := config.GlobalConfig.Script.DeleteuserScript
	// 在shell中执行命令运行脚本删除用户
	cmd := exec.Command("bash", scriptPath, username)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("执行命令失败:", err)
		fmt.Println("命令输出:", string(output))
		return
	}
	fmt.Println("命令输出:", string(output))
	// 删除数据库中的用户记录
	result := dao.DB.Where("Name =?", username).Delete(u)
	if err = result.Error; err != nil {
		return err
	}
	// 获取受影响的行数
	rowsAffected := result.RowsAffected
	if rowsAffected == 0 {
		return fmt.Errorf("没有找到要删除的记录")
	}
	return
}

// ! 批量删除用户
func (u *SftpUsers) DeleteUsers(userNames []string) (err error) {
	// 读取配置脚本路径
	scriptPath := config.GlobalConfig.Script.BatchDeleteScript

	// 创建一个切片来存储参数
	// var args []string
	args := make([]string, 0, len(userNames)+1) // +1是因为第一个参数是脚本路径
	args = append(args, scriptPath)
	args = append(args, userNames...)

	// 在shell中执行命令运行脚本批量删除用户
	cmd := exec.Command("bash", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("执行命令失败:", err)
		fmt.Println("命令输出:", string(output))
		return
	}
	fmt.Println("命令输出:", string(output))

	// 删除数据库中的用户记录
	err = dao.DB.Where("name IN (?)", userNames).Delete(u).Error
	return
}

// ! 通过用户名查询一个用户的信息,并返回用户信息或错误信息
func (u *SftpUsers) GetUserInfo(username string) (user SftpUsers, err error) {
	err = dao.DB.Where("name=?", username).First(&user).Error
	return
}

// ! 通过id查询一个用户的信息,并返回用户信息或错误信息
func (u *SftpUsers) GetUserInfoById(id int) (err error) {
	err = dao.DB.Where("id=?", id).First(u).Error
	return
}

// ! 通过id批量查询用户
func (u *SftpUsers) GetUsersByIds(ids []int64) (users []SftpUsers, err error) {
	err = dao.DB.Where("id IN (?)", ids).Find(&users).Error
	return
}

// ! 获取账号总数
func (u *SftpUsers) GetTotalCount() (count int64, err error) {
	err = dao.DB.Model(u).Count(&count).Error
	return
}

// ! 月新增账号数
func (u *SftpUsers) GetMonthlyNewCount() (count int64, err error) {
	// 获取时间
	now := time.Now()
	// 提取年月
	year, month, _ := now.Date()
	// 构建当前月的开始时间
	startOfCurrentMonth := time.Date(year, month, 1, 0, 0, 0, 0, now.Location())
	err = dao.DB.Model(u).Where("created_at >= ?", startOfCurrentMonth).Count(&count).Error
	return
}
