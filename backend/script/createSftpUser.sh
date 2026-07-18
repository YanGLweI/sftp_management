#!/bin/bash

# 添加SFTP用户脚本（函数式）
# 版本: 2.0.0-20250908

# 配置信息
DB_NAME="sftp"
DB_TABLE="t_sftp_users"
SCRIPT_VERSION="2.0.0"  # 脚本版本号
DateTime=$(date +"%Y-%m-%d %H:%M:%S")



# SFTP根用户
account="datacenter"

# 全局变量
localip=""
path=""
group=""
user=""
pass=""
expireTime=""

# 显示版本信息
show_version() {
    echo "SFTP用户创建脚本 v$SCRIPT_VERSION"
    echo "功能: 创建和管理SFTP用户，配置目录权限，并同步信息到数据库"
    echo "用法: $0 <用户名> <密码> 或 $0 <用户名> <密码> <过期时间,例:10> 或 $0 --version"
}

# 获取本机IP
get_local_ip() {
    localip=$(ip a | grep global | awk -F'[ /]' '{print $6}' | sed -n '1p')
}

# 根据本地IP设置主机信息
set_host_info() {
    path=$(getent passwd "$account" | cut -d: -f6)
    if [ -z "$path" ]; then
        echo "错误: 未找到根用户 $account 的家目录"
        exit 1
    fi
    group=$(id -g -n "$account")
    if [ -z "$group" ]; then
        echo "错误: 未找到根用户 $account 的主组"
        exit 1
    fi

}

# 检查用户是否已存在
check_user_exists() {
    local haduser=$(grep -w "^$user" /etc/passwd)
    if [ -n "$haduser" ]; then
        echo "已经存在$user用户"
        return 0  # 用户存在
    fi
    return 1  # 用户不存在
}

# 创建SFTP用户
create_sftp_user() {
    # 创建用户，禁止系统登录，不自动创建home目录
    useradd "$user" -g "$group" -s /bin/false -M || {
        echo "创建用户$user失败"
        exit 1
    }

    # 创建用户目录结构
    mkdir -p "$path"/"$user"/{input,output} || {
        echo "创建目录结构失败"
        exit 1
    }

    # 设置用户主目录
    usermod "$user" -d "$path"/"$user" || {
        echo "设置用户主目录失败"
        exit 1
    }

    # 设置目录权限
    chown "$user":"$group" "$path"/"$user"/* || {
        echo "设置目录所有者失败"
        exit 1
    }

    # 修改目录权限
    chmod -R 775 "$path"/"$user"/* || {
        echo "修改目录权限失败"
        exit 1
    }

    # 设置权限继承
    setfacl -m d:g:"$group":rwx -R "$path"/"$user"/* || {
        echo "设置权限继承失败"
        exit 1
    }

    if [ "$pass" != "" ]; then
        # 设置用户密码
        echo "$pass" | passwd --stdin "$user" > /dev/null || {
            echo "设置用户密码失败"
            exit 1
        }
    fi

    # 设置密码永不过期,expires参数非空时设置密码永不过期
    if [ -n "$expireTime" ]; then
        chage -M $expireTime "$user" || {
        echo "设置密码永不过期失败"
        exit 1
        }
    fi
    echo "用户$user创建成功"
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

# 插入或更新用户信息到数据库
insert_or_update_user() {
    local home_dir="$path/$user"
    local expires=$1
    
    # 使用SQL预处理语句防止注入
    mariadb "$DB_NAME" <<EOF
INSERT INTO $DB_TABLE (name, home, password_expires,created_at)
VALUES ('$user', '$home_dir', '$expires','$DateTime')
ON DUPLICATE KEY UPDATE 
    home = VALUES(home),
    password_expires = VALUES(password_expires);
EOF

    if [ $? -ne 0 ]; then
        echo "错误: 无法插入/更新用户 $username 的信息"
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
    get_local_ip
    set_host_info
    
    if ! check_user_exists; then
        create_sftp_user
    fi

    # 插入数据库
    local expiry=$(get_password_expiry "$user")
    insert_or_update_user "$expiry"
}

# 执行主函数
main "$@"
