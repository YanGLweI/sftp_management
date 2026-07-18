#!/bin/bash

# 日志目录
LOG_DIR="/data/sftplogs"

# 声明数组
dates=()
counts=()

# 倒序循环：6天前 → 今天（最早在前）
for i in {6..0}; do
    # 文件名用：当天日期（例如 20260602）
    day=$(date -d "$i days ago" +%Y%m%d)
    
    # 真实日期 = 文件名日期 -1天（日志内容是前一天的数据）
    day_cn=$(date -d "$i days ago -1 day" +%-m月%-d日)
    
    # 日志文件名
    log_file="${LOG_DIR}/sftp.log-${day}"

    # 统计当天登录总次数
    if [ -f "$log_file" ]; then
        total=$(grep 'session opened for user' "$log_file" 2>/dev/null \
        | awk '{print $NF}' \
        | sort \
        | uniq -c \
        | awk '{sum += $1} END {print sum+0}')
    else
        total=0
    fi

    # 存入数组
    dates+=("$day_cn")
    counts+=("$total")
done

# 输出最终格式
echo "${dates[*]}" | sed 's/ /,/g'
echo "${counts[*]}" | sed 's/ /,/g'
