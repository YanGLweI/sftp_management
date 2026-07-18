#!/bin/bash

# convert_octal_utf8.sh
# 将日志文件中的八进制UTF-8转义序列转换为中文字符（直接修改原文件）

# 如果没有提供文件名，则显示用法
if [ $# -eq 0 ]; then
    echo "用法: $0 <文件名>"
    echo "示例: $0 sshd.log"
    echo "或者处理多个文件: $0 *.log"
    echo "警告：此脚本会直接修改原文件，建议先备份！"
    exit 1
fi

# 处理每个输入文件
for input_file in "$@"; do
    # 检查文件是否存在
    if [ ! -f "$input_file" ]; then
        echo "错误: 文件 '$input_file' 不存在"
        continue
    fi
    
    echo "正在处理: $input_file"
    
    # 备份原文件
   # backup_file="${input_file}.backup.$(date +%Y%m%d_%H%M%S)"
   # cp "$input_file" "$backup_file"
   # echo "已创建备份: $backup_file"
    
    # 方法1: 使用perl直接修改原文件（创建.bak备份）
    perl -i -pe 's/\\\\([0-7]{3})/chr(oct($1))/ge' "$input_file"
    
    echo "转换完成: $input_file (已修改)"
    #echo "转换前后的差异:"
    #diff -u "$backup_file" "$input_file" | head -20
    #echo ""
    systemctl restart rsyslog.service
done
