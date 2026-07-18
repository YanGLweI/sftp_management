#!/bin/bash
set -euo pipefail

# ===================== 配置项（保留你的原始配置）=====================
TEMP_LOG_FILE="/tmp/dnf_update_$(date +%Y%m%d%H%M%S).log"
HOSTNAME=$(hostname)
DB_MAIN_TABLE="t_update_mains"    # 主表名
DB_DETAIL_TABLE="t_update_details"  # 详情表名
DB_NAME="sftp"        # 数据库名

# ===================== 核心逻辑 =====================
# 1. 记录时间
START_TIME=$(date +%s)
UPDATE_TIME=$(date +"%Y-%m-%d %H:%M:%S")

# 2. 执行dnf更新并捕获输出（保留原始输出，用于失败/无更新场景）
echo "[$(date)] 开始执行系统更新..."
dnf -y update --security 2>&1 | tee "${TEMP_LOG_FILE}" || true  # 继续执行，即使更新失败,避免中断脚本

# 3. 计算耗时
END_TIME=$(date +%s)
DURATION=$(echo "scale=2; ${END_TIME} - ${START_TIME}" | bc)

# 4. 先读取原始更新输出（用于失败/暂无更新场景）
RAW_DETAILS=$(cat "${TEMP_LOG_FILE}" | sed -e 's|\\|\\\\|g' -e 's|"|\\"|g' -e "s|'|\\\\'|g")

# 5. 初始化摘要变量
UPDATE_BRIEF=""

# 6. 判定更新状态 + 赋值摘要/详情（简化包名提取逻辑）
DNF_EXIT_CODE=$?

if [ ${DNF_EXIT_CODE} -ne 0 ]; then
    STATUS="失败"
    DETAILS="${RAW_DETAILS}"
    UPDATE_BRIEF="更新失败请查看详情"  # 失败场景摘要
elif grep -qiE "Nothing to do|无需任何处理" "${TEMP_LOG_FILE}"; then  # 兼容中英文无更新提示
    STATUS="暂无可用更新"
    DETAILS="${RAW_DETAILS}"
    UPDATE_BRIEF="当前已是最新版本"    # 无更新场景摘要
# 没有Complete或者完毕，说明失败
elif ! grep -qiE "Complete|完毕" "${TEMP_LOG_FILE}"; then
    STATUS="失败"
    DETAILS="${RAW_DETAILS}"
    UPDATE_BRIEF="更新失败请查看详情"  # 失败场景摘要
else
    STATUS="成功"
    # 成功时：通过dnf history获取结构化详情（简化包名提取）
    echo "[$(date)] 更新成功，获取结构化更新历史..."
    # 步骤1：获取最新update事务ID（匹配history list中的"更新"关键词）
    HISTORY_ID=$(dnf history list | awk 'NR>=2 && NR<=20' | grep -iE "update" | head -n 1 | awk '{print $1}' || true)
    if [ -n "${HISTORY_ID}" ] && [[ "${HISTORY_ID}" =~ ^[0-9]+$ ]]; then
        # 步骤2：获取结构化详情（保留原始输出用于提取）
        HISTORY_DETAILS_RAW=$(dnf history info "${HISTORY_ID}" 2>&1)
        # 步骤3：【简化核心】仅取"已改变的包："/"Packages Altered"的下一行，提取包名
        FIRST_PACKAGE=$(echo "${HISTORY_DETAILS_RAW}" | \
            grep -A1 -E "已改变的包：|Packages Altered" |  # 匹配目标行并取其后1行
            tail -n1 |                                    # 只保留目标行的下一行
            grep -v "^$" |                                # 过滤空行
            awk '{print $2}' || true)                     # 提取包名（第二个字段）
        
        # 容错：如果提取不到包名，设置默认摘要
        if [ -n "${FIRST_PACKAGE}" ]; then
            UPDATE_BRIEF="更新成功：${FIRST_PACKAGE} 等..."  # 成功场景摘要（带包名）
        else
            UPDATE_BRIEF="更新成功（无具体包信息）"
        fi
        
        # 转义特殊字符后存入DETAILS
        DETAILS=$(echo "${HISTORY_DETAILS_RAW}" | sed -e 's|\\|\\\\|g' -e 's|"|\\"|g' -e "s|'|\\\\'|g")
    else
        # 容错：如果获取history ID失败，仍用原始输出
        echo "[$(date)] 警告：获取dnf history ID失败，使用原始更新输出"
        DETAILS="${RAW_DETAILS}"
        UPDATE_BRIEF="更新成功（日志异常）"
    fi
fi

# 7. 插入主表 + 同一会话获取自增ID（加入update_brief字段）
MAIN_ID=$(mariadb -D ${DB_NAME} -N -e "
INSERT INTO ${DB_MAIN_TABLE} (hostname, update_time, status, duration, update_brief)
VALUES ('${HOSTNAME}', '${UPDATE_TIME}', '${STATUS}', ${DURATION}, '${UPDATE_BRIEF}');
SELECT LAST_INSERT_ID();
")

# 8. 校验ID有效性，防止空值/0
if [ -z "${MAIN_ID}" ] || [ "${MAIN_ID}" -eq 0 ]; then
    echo "[$(date)] 错误：获取主表自增ID失败！ID=${MAIN_ID}"
    rm -f "${TEMP_LOG_FILE}"  # 清理临时文件
    exit 1  # 终止脚本，避免插入详情表时关联错误ID
fi

# 9. 插入详情表（关联主表ID，逻辑不变）
mariadb -D ${DB_NAME} -e "
INSERT INTO ${DB_DETAIL_TABLE} (main_id, update_details)
VALUES (${MAIN_ID}, '${DETAILS}');
"

# 10. 清理临时文件
rm -f "${TEMP_LOG_FILE}"

# 11. 输出结果（新增摘要展示）
echo "=================================================="
echo "更新完成！结果如下："
echo "主机名：${HOSTNAME}"
echo "更新时间：${UPDATE_TIME}"
echo "更新状态：${STATUS}"
echo "更新摘要：${UPDATE_BRIEF}"
echo "耗时：${DURATION} 秒"
echo "主记录ID：${MAIN_ID}"
echo "=================================================="