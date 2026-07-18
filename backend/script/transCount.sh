#!/bin/bash
LOG_DIR="/var/log"

dates=()
total_transfer=()
visit_counts=()

for i in {6..0}; do
    day=$(date -d "$i days ago" +%Y%m%d)
    day_cn=$(date -d "$i days ago -1 day" +%-m月%-d日)
    log_file="${LOG_DIR}/sftp.log-${day}"

    # 当天日志直接使用 sftp.log
    [ $i = 0 ] && log_file="${LOG_DIR}/sftp.log"

    if [ -f "$log_file" ]; then
        # 上传：write > 0
        upload=$(grep -E ' close ".+" bytes read [0-9]+ written [1-9]' "$log_file" | wc -l)
        
        # 下载：read > 0
        download=$(grep -E ' close ".+" bytes read [1-9]+ written 0' "$log_file" | wc -l)
        
        # 总传输 = 上传 + 下载
        total=$((upload + download))
    else
        total=0
    fi

    dates+=("$day_cn")
    total_transfer+=("$total")
done

# 输出数组格式（直接给前端/图表用）
echo "${dates[*]}" | sed 's/ /,/g'
echo "${total_transfer[*]}" | sed 's/ /,/g'
