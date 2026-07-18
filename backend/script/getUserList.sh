#!/bin/bash

# SFTP用户信息收集与数据库插入脚本

# 配置信息
DB_NAME="sftp"
DB_TABLE="t_sftp_users"

# SFTP根用户
account="datacenter"

group=""

# 日志文件
LOG_FILE="/var/log/sftp_user_info.log"

# 函数：记录日志
log_message() {
    local message="[$(date '+%Y-%m-%d %H:%M:%S')] $1"
    echo "$message"
    echo "$message" >> "$LOG_FILE"
}

# 检查命令依赖
check_dependency() {
    command -v "$1" >/dev/null 2>&1 || { log_message "错误: 需要安装 $1"; exit 1; }
}

# 检查必要的命令
check_dependency "getent"
check_dependency "chage"
check_dependency "mysql"

# 创建数据库表（如果不存在）
create_table() {
    log_message "创建数据库表..."
    mariadb "$DB_NAME" <<EOF
CREATE TABLE IF NOT EXISTS $DB_TABLE (
    id BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	created_at DATETIME(3) NULL DEFAULT NULL,
	updated_at DATETIME(3) NULL DEFAULT NULL,
	deleted_at DATETIME(3) NULL DEFAULT NULL,
	name LONGTEXT NULL DEFAULT NULL COLLATE 'utf8mb4_general_ci',
	home LONGTEXT NULL DEFAULT NULL COLLATE 'utf8mb4_general_ci',
	password_expires LONGTEXT NULL DEFAULT NULL COLLATE 'utf8mb4_general_ci',
	PRIMARY KEY (id) USING BTREE,
	INDEX idx_t_sftp_users_deleted_at (deleted_at) USING BTREE
)
COLLATE='utf8mb4_general_ci';
EOF
}

# 获取sftpusers组的GID
get_sftp_gid() {
    set_group
    local gid=$(getent group "$group" | cut -d: -f3)
    if [ -z "$gid" ]; then
        log_message "错误: 未找到组 $group"
        exit 1
    fi
    echo "$gid"
}

# 获取sftpusers组成员信息
get_sftp_users() {
    local gid=$1
    getent passwd | awk -F: -v gid="$gid" '$4 == gid {print $1" "$6}'
}

# 获取用户密码到期日期
get_password_expiry() {
    local username=$1
    local expiry=$(chage -li "$username" | awk -F'[[:space:]]*[:：][[:space:]]*' 'NR==2 {print $2}')
    
    # 处理"never"和空值情况
    if [[ "$expiry" == "never" || "$expiry" == "从不" || -z "$expiry" ]]; then
        echo "9999-12-31"  # 表示永不过期
    else
        echo "$expiry"
    fi
}

# 转义特殊字符以防止SQL注入
escape_sql() {
    local input="$1"
    # 替换单引号为两个单引号（SQL转义方式）
    echo "${input//\'/\'\'}"
}

# 插入或更新用户信息到数据库
insert_or_update_user() {
    local username=$(escape_sql "$1")
    local home_dir=$(escape_sql "$2")
    local expires=$(escape_sql "$3")
    local created_at=$(escape_sql "$4")
    
    log_message "处理用户: $1"
    
    # 检查用户是否已存在
    local exists=$(mariadb -h "$DB_HOST" -u "$DB_USER" -p"$DB_PASS" "$DB_NAME" -s -N -e \
        "SELECT 1 FROM $DB_TABLE WHERE name='$username' LIMIT 1")
    
    if [ "$exists" = "1" ]; then
        # 用户存在，执行更新
        mariadb -h "$DB_HOST" -u "$DB_USER" -p"$DB_PASS" "$DB_NAME" <<EOF
UPDATE $DB_TABLE SET 
    home = '$home_dir',
    password_expires = '$expires',
    created_at = '$created_at'
WHERE name = '$username' AND deleted_at IS NULL;
EOF
        log_message "用户 $1 已存在，已更新信息"
    else
        # 用户不存在，执行插入
        mariadb -h "$DB_HOST" -u "$DB_USER" -p"$DB_PASS" "$DB_NAME" <<EOF
INSERT INTO $DB_TABLE (name, home, password_expires,created_at)
VALUES ('$username', '$home_dir', '$expires','$created_at');
EOF
        log_message "用户 $1 不存在，已插入新记录"
    fi

    if [ $? -ne 0 ]; then
        log_message "错误: 无法处理用户 $1 的信息"
    fi
}

set_group() {
    group=$(id -g -n "$account")
    if [ -z "$group" ]; then
        echo "错误: 未找到根用户 $account 的主组"
        exit 1
    fi
}

# 主程序
main() {
    log_message "开始收集SFTP用户信息..."
    
    # 创建数据库表
    create_table
    
    # 获取sftpusers组的GID
    SFTP_GID=$(get_sftp_gid)
    log_message "sftpusers组GID: $SFTP_GID"
    
    # 获取sftpusers组成员
    log_message "获取sftpusers组成员..."
    USERS=$(get_sftp_users "$SFTP_GID")
    
    if [ -z "$USERS" ]; then
        log_message "警告: sftpusers组中没有找到用户"
        exit 0
    fi
    
    # 处理每个用户
    while read -r line; do
        username=$(echo "$line" | awk '{print $1}')
        home_dir=$(echo "$line" | awk '{print $2}')
        created_at=$(stat -c %w "$home_dir" | awk -F. '{print $1}')
        
        # 获取密码到期日期
        expires=$(get_password_expiry "$username")
        
        # 插入或更新数据库
        insert_or_update_user "$username" "$home_dir" "$expires" "$created_at"
    done <<< "$USERS"
    
    log_message "SFTP用户信息收集完成!"
}

# 执行主程序
main
