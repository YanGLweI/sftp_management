#!/bin/bash
source /etc/profile
#执行加固脚本并输出日志到指定文件夹。
# mkdir -p /data/System-Check-logs
# /script/System_Check-1.2.sh &> /data/System-Check-logs/`ip a |grep global |awk -F'[ /]' '{print $6}' |sed -n '1p'`_systemcheck.`date '+%Y_%m_%d'`.log
# #安装MariaDB客户端，用于连接远程数据库。
# #判断是否存在mysql-insert.sh的计划任务，不存在就创建mysql-insert.sh计划任务每天0点执行一次。
# cron=`grep "mysql-insert.sh" /var/spool/cron/root` #赋值
# if  [ -z "$cron" ];then
#     echo "0 * * * * /script/mysql-insert.sh" >> /var/spool/cron/root
#     echo "将mysql-insert.sh添加到计划任务成功！"
# fi
./script/systemSecurity.sh
######################################################################################################
# DB_USER=it            #远程数据库用户
# DB_PASSWD=a*999999    #远程数据库密码
DB_NAME=sftp            #远程数据库名
TABLE="t_system_securities"    #远程数据库表
# DB_IP="10.60.254.127" #远程数据库IP地址
#以下为赋值加固后的参数，以便于插入数据库表中。
day=$(date +%Y/%m/%d_%H:%M:%S)
name=$(hostname)
centos=$(cat /etc/redhat-release)
kernel=$(uname -r)
local_ip=$(ip a |grep global |awk -F'[ /]' '{print $6}' |sed -n '1p')
values5=`grep "gpgcheck=" /etc/yum.conf |awk -F= '{print $2}'`
values6=`grep gpgcheck /etc/yum.repos.d/redhat.repo | awk -F' = ' '{print $2}' |sed -n '1p'`
values7=`grep "^PASS_MAX_DAYS" /etc/login.defs |awk  '{print $2}' |sed -n '1p'`
values8=`grep "^PASS_MIN_DAYS" /etc/login.defs |awk  '{print $2}' |sed -n '1p'`
values9=`grep "^PASS_MIN_LEN" /etc/login.defs |awk  '{print $2}' |sed -n '1p'`
values10=`grep "^PASS_WARN_AGE" /etc/login.defs |awk  '{print $2}' |sed -n '1p'`
values11=`useradd -D | grep INACTIVE |awk -F= '{print $2}'`
values12=`grep "^root:" /etc/passwd | cut -f4 -d:`
values13=`echo $TMOUT`
values14=`systemctl is-enabled crond`
values15=`stat  /etc/crontab | sed -n '4p' |tr -s " " | sed 's/^.../Access:/'`
values16=`stat  /etc/cron.hourly | sed -n '4p' |tr -s " " | sed 's/^.../Access:/'`
values17=`stat  /etc/cron.daily | sed -n '4p' |tr -s " " | sed 's/^.../Access:/'`
values18=`stat  /etc/cron.weekly | sed -n '4p' |tr -s " " | sed 's/^.../Access:/'`
values19=`stat  /etc/cron.monthly | sed -n '4p' |tr -s " " | sed 's/^.../Access:/'`
values20=`stat /etc/cron.deny >/dev/null 2>&1 && stat /etc/cron.deny | sed -n '4p' |tr -s " " | sed 's/^.../Access:/' || echo "No such file or directory"`
values21=`stat /etc/at.deny >/dev/null 2>&1 && stat /etc/at.deny | sed -n '4p' |tr -s " " | sed 's/^.../Access:/' || echo "No such file or directory"`
values22=`stat /etc/cron.allow >/dev/null 2>&1 && stat /etc/cron.allow | sed -n '4p' |tr -s " " | sed 's/^.../Access:/' || echo "No such file or directory"`
values23=`stat /etc/at.allow >/dev/null 2>&1 && stat /etc/at.allow | sed -n '4p' |tr -s " " | sed 's/^.../Access:/' || echo "No such file or directory"`
values24=`stat  /etc/ssh/sshd_config | sed -n '4p' |tr -s " "  | sed 's/^.../Access:/' `
values26=`grep "^LogLevel" /etc/ssh/sshd_config |sed -n '1p' | awk '{print $2}'`
values27=`grep "^X11Forwarding" /etc/ssh/sshd_config | awk '{print $2}'`
values28=`grep "MaxAuthTries" /etc/ssh/sshd_config |awk '{print $2}'`
values29=`grep "IgnoreRhosts" /etc/ssh/sshd_config |awk '{print $2}'`
values30=`grep HostbasedAuthentication /etc/ssh/sshd_config | sed -n '1p' |awk '{print $2}'`
values31=`grep "^PermitRootLogin" /etc/ssh/sshd_config | awk '{print $2}'`
values32=`grep "PermitEmptyPasswords" /etc/ssh/sshd_config | awk '{print $2}'`
values33=`grep "PermitUserEnvironment" /etc/ssh/sshd_config | awk '{print $2}'`
values34=`grep "ClientAliveInterval" /etc/ssh/sshd_config | awk '{print $2}'`
values35=`grep "ClientAliveCountMax" /etc/ssh/sshd_config | awk '{print $2}'`
values36=`grep "LoginGraceTime" /etc/ssh/sshd_config | awk '{print $2}'`
values37=`grep "minlen" /etc/security/pwquality.conf.d/50-pwlength.conf | awk -F' ' '{print $3}'`
values52=`grep "minclass" /etc/security/pwquality.conf.d/50-pwcomplexity.conf | awk -F' ' '{print $3}'`
values38=`grep "dcredit" /etc/security/pwquality.conf.d/50-pwcomplexity.conf | awk -F' ' '{print $3}'`
values39=`grep "ucredit" /etc/security/pwquality.conf.d/50-pwcomplexity.conf | awk -F' ' '{print $3}'`
values40=`grep "lcredit" /etc/security/pwquality.conf.d/50-pwcomplexity.conf | awk -F' ' '{print $3}'`
values41=`grep "ocredit" /etc/security/pwquality.conf.d/50-pwcomplexity.conf | awk -F' ' '{print $3}'`
values42=`grep -oP '^\h*remember\h*=\h*\K(2[4-9]|[3-9][0-9]|[1-9][0-9]{2,})\b' /etc/security/pwhistory.conf`
values44=`stat /etc/passwd | sed -n '4p' | tr -s " " | sed 's/^.../Access:/' `
values45=`stat /etc/passwd- | sed -n '4p' | tr -s " " | sed 's/^.../Access:/' `
values46=`stat /etc/group | sed -n '4p' | tr -s " " | sed 's/^.../Access:/' `
values47=`stat /etc/group- | sed -n '4p' | tr -s " " | sed 's/^.../Access:/' `
values48=`stat /etc/shadow | sed -n '4p' | tr -s " " | sed 's/^.../Access:/' `
values49=`stat /etc/shadow- | sed -n '4p' | tr -s " " | sed 's/^.../Access:/' `
values50=`stat /etc/gshadow | sed -n '4p' | tr -s " " | sed 's/^.../Access:/' `
values51=`stat /etc/gshadow- | sed -n '4p' | tr -s " " | sed 's/^.../Access:/' `
values53=`update-crypto-policies --show`
values54=`grep -Prs -- '^\h*(server|pool)\h+[^#\n\r]+' /etc/chrony.conf | sed -n '1p'`
#############################################################################

# 检测是否支持 --skip-ssl 参数
if mariadb -h $DB_IP --user=$DB_USER --password=$DB_PASSWD --skip-ssl --connect-timeout=5 -e "SELECT 1;" $DB_NAME &>/dev/null; then
    SKIP_SSL="--skip-ssl"
else
    SKIP_SSL=""
fi

#连接数据库，插入加固检查的参数到表中。
# mysql -h $DB_IP --user=$DB_USER  --password=$DB_PASSWD $SKIP_SSL $DB_NAME << EOF
mariadb $SKIP_SSL $DB_NAME << EOF
INSERT INTO $TABLE (
    Date,
    Hostname,
    operasystem,
    Kernel,
    IP,
    \`dnf_conf_gpgcheck\`,
    \`redhat_repo_gpgcheck\`,
    PASS_MAX_DAYS,
    PASS_MIN_DAYS,
    PASS_MIN_LEN,
    PASS_WARN_AGE,
    INACTIVE,
    GID,
    TMOUT,
    Cron,
    crontab,
    \`cron_hourly\`,
    \`cron_daily\`,
    \`cron_weekly\`,
    \`cron_monthly\`,
    \`cron_deny\`,
    \`at_deny\`,
    \`cron_allow\`,
    \`at_allow\`,
    sshd_config,
    LogLevel,
    X11Forwarding,
    MaxAuthTries,
    IgnoreRhosts,
    HostbasedAuthentication,
    PermitRootLogin,
    PermitEmptyPasswords,
    PermitUserEnvironment,
    ClientAliveInterval,
    ClientAliveCountMax,
    LoginGraceTime,
    minlen,
    minclass,
    dcredit,
    ucredit,
    lcredit,
    ocredit,
    password_remember,
    passwd,
    \`passwd-\`,
    \`group\`,
    \`group-\`,
    shadow,
    \`shadow-\`,
    gshadow,
    \`gshadow-\`,
    crypto_policies,
    ntp_server
)
VALUES (
    "$day",
    "$name",
    "$centos",
    "$kernel",
    "$local_ip",
    "$values5",
    "$values6",
    "$values7",
    "$values8",
    "$values9",
    "$values10",
    "$values11",
    "$values12",
    "$values13",
    "$values14",
    "$values15",
    "$values16",
    "$values17",
    "$values18",
    "$values19",
    "$values20",
    "$values21",
    "$values22",
    "$values23",
    "$values24",
    "$values26",
    "$values27",
    "$values28",
    "$values29",
    "$values30",
    "$values31",
    "$values32",
    "$values33",
    "$values34",
    "$values35",
    "$values36",
    "$values37",
    "$values52",
    "$values38",
    "$values39",
    "$values40",
    "$values41",
    "$values42",
    "$values44",
    "$values45",
    "$values46",
    "$values47",
    "$values48",
    "$values49",
    "$values50",
    "$values51",
    "$values53",
    "$values54"
);
EOF

echo "数据库插入完成！"
