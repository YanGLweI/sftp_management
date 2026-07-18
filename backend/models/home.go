package models

import (
	"fmt"
	"os/exec"
	"sftpbackend/config"
	"strings"
)

// 获取最近7天传输统计，用于首页展示柱状图
func GetTransCount() ([]string, []string) {
	script := config.GlobalConfig.Script.TransCountScript
	cmd := exec.Command(script)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("执行脚本失败:", err)
		return nil, nil
	}

	outputLines := strings.Split(string(output), "\n")

	transXaxis := strings.Split(outputLines[0], ",")
	transFullDay := strings.Split(outputLines[1], ",")

	return transXaxis, transFullDay
}

// 获取最近7天访问统计，用于首页展示柱状图
func GetAccessCount() ([]string, []string) {
	script := config.GlobalConfig.Script.AccessCountScript
	cmd := exec.Command(script)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("执行脚本失败:", err)
		return nil, nil
	}

	outputLines := strings.Split(string(output), "\n")

	accessXaxis := strings.Split(outputLines[0], ",")
	accessFullDay := strings.Split(outputLines[1], ",")

	return accessXaxis, accessFullDay
}
