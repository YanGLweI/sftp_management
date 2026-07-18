#!/bin/bash
# 更新SFTP用户脚本（函数式）
# 版本：1.0.0
# 用于更新SFTP用户的密码
# 需要传入用户名和密码作为参数

# 配置信息
DB_NAME="sftp"
DB_TABLE="t_sftp_users"
SCRIPT_VERSION="1.0.0"  # 脚本版本号
DateTime=$(date +"%Y-%m-%d %H:%M:%S")



# 全局变量
user=""
pass=""
expireTime=""

# 显示版本信息
show_version() {
    echo "SFTP用户创建脚本 v$SCRIPT_VERSION"
    echo "功能: 更新或修改SFTP用户密码，并同步信息到数据库"
    echo "用法: $0 <用户名> <密码> 或 $0 <用户名> <密码> <过期时间,例:10> 或 $0 --version"
}

# 检查用户是否已存在
check_user_exists() {
    local haduser=$(grep -w "^$user" /etc/passwd)
    if [ -n "$haduser" ]; then
        echo "用户$user已存在，准备更新密码"
        return 0  # 用户存在
    else
        echo "用户$user不存在，请先创建用户"
        exit 1  # 用户不存在
    fi

}

update_sftp_user() {
    # 设置密码永不过期,expires参数非空时设置密码永不过期
    if [ -n "$expireTime" ]; then
        chage -M $expireTime "$user" || {
        echo "设置密码永不过期失败"
        exit 1
        }
    else
        # 如果没有设置过期时间
        chage -M 30 "$user" || {
            echo "设置密码过期时间失败"
            exit 1
        }
        chage -m 1 "$user" || {
            echo "设置密码最小间隔失败"
            exit 1
        }
        chage -I 30 "$user" || {
            echo "设置密码失效时间失败"
            exit 1
        }
    fi
    if [ "$pass" != "" ]; then
        # 设置用户密码
        echo "$pass" | passwd --stdin "$user" > /dev/null || {
            echo "设置用户密码失败"
            exit 1
        }
    fi

    echo "用户$user密码更新成功"
}

# 获取用户密码到期日期
get_password_expiry() {
    local expiry=$(chage -li "$user" | awk -F'[[:space:]]*[:：][[:space:]]*' 'NR==2 {print $2}')

    # 处理"never"和空值情况
    if [[ "$expiry" == "never" || "$expiry" == "从不" || -z "$expiry" ]]; then
        echo "9999-12-31"  # 表示永不过期
    else
        echo "$expiry"
    fi
}

# 修改密码后更新数据库表中，对应用户的过期时间
update_user_expiry() {
    local expiry=$1

    mariadb "$DB_NAME" <<EOF
UPDATE $DB_TABLE
SET password_expires = '$expiry',
    updated_at = '$DateTime'
WHERE name = '$user';
EOF
    if [ $? -ne 0 ]; then
        echo "错误: 无法更新用户 $user 的密码过期时间"
        exit 1
    fi
}

# 主函数
main() {
    # 检查是否需要显示版本信息
    if [ $# -eq 1 ] && [ "$1" = "--version" ]; then
        show_version
        exit 0
    fi
    # 检查参数不等于2也不等于3
    if [ $# -ne 2 ] && [ $# -ne 3 ]; then
        echo "用法: $0 <用户名> <密码> 或 $0 <用户名> <密码> <过期时间,例:10> 或 $0 --version"
        exit 1
    fi

    user="$1"
    pass="$2"
    # 如果有第三个参数，则设置过期时间
    if [ $# -eq 3 ]; then   
        expireTime="$3"
    fi

    # 执行主要流程
    # 用户存在则更新用户密码
    if check_user_exists; then
        update_sftp_user
    fi

    # 更新数据库
    local expiry=$(get_password_expiry "$user")
    update_user_expiry "$expiry"
}

# 执行主函数
main "$@"