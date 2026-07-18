#!/bin/bash
# 批量删除SFTP用户脚本（函数式）
# 版本：1.0.0
# 用于批量删除SFTP用户
# 需要传入用户名作为参数

# 配置信息
SCRIPT_VERSION="1.0.0"  # 脚本版本号
# 全局变量
users=""

# 显示版本信息
show_version() {
    echo "SFTP用户批量删除脚本 v$SCRIPT_VERSION"
    echo "功能: 批量删除多个SFTP用户，并同步信息到数据库"
    echo "用法: $0 <用户名1 用户名2 ···> 或 $0 --version"
}

# 检查用户是否已存在
check_users_exist() {
	for user in "${users[@]}"; do
		local haduser=$(grep -w "^$user" /etc/passwd)
		if [ -n "$haduser" ]; then
			echo "用户$user已存在，准备删除"
		else
			echo "用户$user不存在。"
			exit 1  # 用户不存在
		fi
	done
}

# 获取用户家目录
get_user_home() {
	local home=$(getent passwd "$1" | cut -d: -f6)
	if [ -z "$home" ]; then
		echo "错误: 未找到用户 $1 的家目录"
		continue
	fi
	echo "$home"
}

delete_sftp_users() {
	for user in "${users[@]}"; do
		local home=$(get_user_home "$user")
		# 如果用户登陆过有进程，关闭用户所有进程
    pkill -u "$user" > /dev/null 2>&1
		# 删除用户
		userdel "$user" > /dev/null || {
			echo "删除用户$user失败"
			exit 1
		}
		rm -rf "$home" || {
			echo "删除用户家目录$home失败"
			exit 1
		}
		echo "用户$user删除成功"
	done
}

# 主函数
main() {
	if [ "$1" == "--version" ]; then
		show_version
		exit 0
	fi
	if [ $# -eq 0 ]; then
		echo "错误: 需要至少一个用户名作为参数"
		show_version
		exit 1
	fi
	# 将传入的参数转换为数组
	users=("$@")
	# 检查用户是否存在
	check_users_exist
	# 删除用户
	delete_sftp_users
}

# 调用主函数
main "$@"