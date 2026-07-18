#!/bin/bash
# 删除SFTP用户脚本（函数式）
# 版本：1.0.0
# 用于删除SFTP用户
# 需要传入用户名作为参数

# 配置信息
SCRIPT_VERSION="1.0.0"  # 脚本版本号
# 全局变量
user=""

# 显示版本信息
show_version() {
    echo "SFTP用户删除脚本 v$SCRIPT_VERSION"
    echo "功能: 删除单个SFTP用户，并同步信息到数据库"
    echo "用法: $0 <用户名> 或 $0 --version"
}

# 检查用户是否已存在
check_user_exists() {
    local haduser=$(grep -w "^$user" /etc/passwd)
    if [ -n "$haduser" ]; then
        echo "用户$user已存在，准备删除"
        return 0  # 用户存在
    else
        echo "用户$user不存在。"
        exit 1  # 用户不存在
    fi

}

# 获取用户家目录
get_user_home() {
    local home=$(getent passwd "$user" | cut -d: -f6)
    if [ -z "$home" ]; then
        echo "错误: 未找到用户 $user 的家目录"
        exit 1
    fi
    echo "$home"
}

delete_sftp_user() {
    # 关闭用户所有进程
    pkill -u "$user" > /dev/null 2>&1
    # 删除用户
    userdel "$user" > /dev/null || {
        echo "删除用户失败"
        exit 1
    }
    rm -rf "$1" || {
        echo "删除用户家目录失败"
        exit 1
    }

    echo "用户$user删除成功"
}

# 主函数
main() {
    if [ "$1" == "--version" ]; then
        show_version
        exit 0
    elif [ $# -ne 1 ]; then
        echo "用法: $0 <用户名> 或 $0 --version"
        exit 1
    fi

    user="$1"

    # 检查用户是否存在
    check_user_exists

    # 获取用户家目录
    local home=$(get_user_home)

    # 删除用户及其家目录
    delete_sftp_user "$home"
}

# 执行主函数
main "$@"