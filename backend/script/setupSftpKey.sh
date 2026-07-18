#!/bin/bash

# 显示用法信息
usage() {
    echo "用法: $0 <用户名>"
    echo "示例: $0 sftpuser"
}

# 检查用户是否存在并获取用户主组
check_user_and_get_group() {
    local username="$1"
    if ! id -u "$username" >/dev/null 2>&1; then
        echo "错误: 用户 '$username' 不存在，请先创建该用户。"
        return 1
    fi
    
    # 获取用户主组名
    local user_group=$(id -g -n "$username")
    echo "$user_group"
    return 0
}

# 获取用户主目录
get_user_home() {
    local username="$1"
    echo "$(eval echo ~"$username")"
}

# 创建.ssh目录并设置权限
create_ssh_directory() {
    local ssh_dir="$1"
    local username="$2"
    local user_group="$3"
    
    if [ ! -d "$ssh_dir" ]; then
        echo "创建.ssh目录..."
        mkdir -p "$ssh_dir"
        chown "$username:$user_group" "$ssh_dir"
        chmod 700 "$ssh_dir"
    else
        echo ".ssh目录已存在，检查权限..."
        # 确保目录权限正确
        chown "$username:$user_group" "$ssh_dir"
        chmod 700 "$ssh_dir"
    fi
}

# 生成SSH密钥对
generate_ssh_keys() {
    local private_key="$1"
    local key_type="$2"
    local key_bits="$3"
    local username="$4"
    
    if [ -f "$private_key" ] && [ -f "${private_key}.pub" ]; then
        echo "警告: 用户 '$username' 已存在SSH密钥对，将直接覆盖..."
    fi
    
    # 生成新的密钥对（无密码）
    echo "生成新的${key_type}密钥对..."
    yes | sudo -u "$username" ssh-keygen -t "$key_type" -b "$key_bits" -N "" -f "$private_key" -q
}

# 配置authorized_keys文件
configure_authorized_keys() {
    local authorized_keys="$1"
    local public_key="$2"
    local username="$3"
    local user_group="$4"
    
    if [ ! -f "$authorized_keys" ]; then
        echo "创建authorized_keys文件..."
        touch "$authorized_keys"
        chown "$username:$user_group" "$authorized_keys"
        chmod 600 "$authorized_keys"
    else
        echo "检查authorized_keys权限..."
        # 确保文件权限正确
        chown "$username:$user_group" "$authorized_keys"
        chmod 600 "$authorized_keys"
    fi
    
    # 检查公钥是否已在authorized_keys中
    if ! grep -qxF "$(cat "$public_key")" "$authorized_keys"; then
        echo "将公钥添加到authorized_keys..."
        cat "$public_key" >> "$authorized_keys"
    else
        echo "公钥已在authorized_keys中"
    fi
}

# 导出私钥到当前目录
export_private_key() {
    local source_key="$1"
    local output_key="$2"
    local current_user=$(whoami)
    local current_group=$(id -g -n "$current_user")
    
    if [ -f "$source_key" ]; then
        cp "$source_key" "$output_key"
        chown "$current_user:$current_group" "$output_key"
        chmod 600 "$output_key"
        echo "私钥已导出到: $output_key"
        return 0
    else
        echo "错误: 无法找到私钥文件 $source_key"
        return 1
    fi
}

# 主函数
main() {
    # 检查参数
    if [ $# -ne 1 ]; then
        usage
        exit 1
    fi
    
    local username="$1"
    local key_type="rsa"
    local key_bits=4096
    local output_key_file="./${username}_sftp_${key_type}_key"
    
    # 检查用户是否存在并获取主组
    local user_group=$(check_user_and_get_group "$username")
    if [ $? -ne 0 ]; then
        exit 1
    fi
    echo "用户 '$username' 的主组为: $user_group"
    
    # 获取用户目录信息
    local user_home=$(get_user_home "$username")
    local ssh_dir="$user_home/.ssh"
    local authorized_keys="$ssh_dir/authorized_keys"
    local private_key="$ssh_dir/id_${key_type}"
    local public_key="${private_key}.pub"
    local temp_pub_key="Temp.pub"
    
    # 执行主要操作
    create_ssh_directory "$ssh_dir" "$username" "$user_group"
    # 如果没有temp.pub,则生成密钥对
    if [ ! -f "$temp_pub_key" ]; then
        generate_ssh_keys "$private_key" "$key_type" "$key_bits" "$username"
        configure_authorized_keys "$authorized_keys" "$public_key" "$username" "$user_group"
        # 导出私钥到当前目录
        if ! export_private_key "$private_key" "$output_key_file"; then
            exit 1
        fi
    else
        # 如果有temp.pub,则使用temp.pub中的公钥
        configure_authorized_keys "$authorized_keys" "$temp_pub_key" "$username" "$user_group"
    fi
    
    echo "操作完成。请妥善保管私钥文件，它是登录SFTP的凭证。"
}

# 启动主函数
main "$@"
    