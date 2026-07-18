package scheduler

import (
	"fmt"
	"log"
	"os/exec"
	"reflect"
	"sftpbackend/models"
)

// SystemSecurityCheck 执行系统加固检查
func SystemSecurityCheck() {
	// 系统加固检查
	// 执行系统加固检查脚本
	cmd := exec.Command(conf.SystemSecurityCheckScript)
	if err := cmd.Run(); err != nil {
		log.Printf("SystemSecurityCheck failed: %v", err)
	}

	// 查询数据库获取刚刚插入的系统加固检查结果
	// 检查结果是否符合预期，更新result字段
	var latest models.SystemSecurity
	latestResult, err := latest.FindLatest()
	if err != nil {
		log.Printf("FindLatest failed: %v", err)
	}

	var standard models.SystemSecurityStandard
	standardResult, err := standard.FindStandard()
	if err != nil {
		log.Printf("FindStandard failed: %v", err)
	}

	if compareSecurityFields(latestResult, standardResult) {
		latestResult.Result = "正常"
	} else {
		latestResult.Result = "异常"
	}

	// 更新数据库记录
	if err := latestResult.Update(); err != nil {
		log.Printf("Update failed: %v", err)
	}
}

// compareSecurityFields 对比两个结构体的字段（仅对比SystemSecurityStandard包含的字段，排除ID）
// 返回 true: 所有字段一致 | false: 存在不一致字段
func compareSecurityFields(actual *models.SystemSecurity, standard *models.SystemSecurityStandard) bool {
	if actual == nil || standard == nil {
		log.Println("对比失败：实际结果/标准配置为空")
		return false
	}

	// 获取结构体的反射值（Elem() 解析指针类型）
	stdVal := reflect.ValueOf(standard).Elem()
	actVal := reflect.ValueOf(actual).Elem()

	// 遍历标准配置的所有字段
	for i := 0; i < stdVal.NumField(); i++ {
		// 获取字段元信息
		stdField := stdVal.Type().Field(i)
		fieldName := stdField.Name

		// 排除ID字段
		if fieldName == "ID" {
			continue
		}

		// 获取实际结果中对应的字段
		actField := actVal.FieldByName(fieldName)
		if !actField.IsValid() {
			log.Printf("对比失败：实际结果中不存在字段 %s", fieldName)
			return false
		}

		// 统一转为字符串对比（兼容所有类型，此处实际都是string）
		stdValStr := fmt.Sprintf("%v", stdVal.Field(i).Interface())
		actValStr := fmt.Sprintf("%v", actField.Interface())

		if fieldName == "NtpServer" && actValStr != "" {
			continue // NtpServer允许为空，不为空则视为符合标准
		}

		// 字段值不一致则返回false
		if stdValStr != actValStr {
			// 更新actual这个字段值为 ：标准值: %s | 实际值: %s
			actField.SetString(fmt.Sprintf("当前值: %s | 标准值: %s", actValStr, stdValStr))

			log.Printf("字段 %s 不匹配 | 标准值: %s | 实际值: %s", fieldName, stdValStr, actValStr)
			return false
		}
	}

	return true
}
