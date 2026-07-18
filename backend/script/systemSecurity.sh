#!/bin/bash
source /etc/profile
#定义变量
day=$(date +%Y/%m/%d_%H:%M:%S)
local_ip=$(ip a |grep global |awk -F'[ /]' '{print $6}' |sed -n '1p')   #IP地址
name=$(hostname)                                                        #主机名
GID=$(grep "^root:" /etc/passwd | cut -f4 -d:)                          #根账户GID
kernel=$(uname -r)                                                      #系统内核
redhat=$(cat /etc/redhat-release)                                       #系统版本
gpgcheck=`grep "gpgcheck=" /etc/dnf/dnf.conf`
PASS=`grep "^PASS" /etc/login.defs | tr -s " " | awk  '{print $2}'`
ipass=`echo $PASS | tr -d " "`
INACTIVE=`useradd -D | grep INACTIVE | cut -d= -f2`
NTPServer="ntp.aliyun.com"   
#-----------------------------------------------------------------------------------------------------------#
echo -e "                          Linux系统加固检查表"
echo
echo "一、服务器信息"
echo
echo -e "Hostname=$name \nSystem=$redhat \nKernel=$kernel "
echo -e "IP=$local_ip" 
echo "-----------------------------------------------------------------------------"
echo "二、系统安全更新"
echo
echo "   全球激活gpgcheck状态"
case $gpgcheck in #判断/etc/dnf/dnf.conf中gpgcheck的值，修改并输出。
    gpgcheck=1)
    echo -e "$gpgcheck";;
    gpgcheck=0)
    sed -i 's/gpgcheck=0/gpgcheck=1/g' /etc/dnf/dnf.conf && echo "gpgcheck已修改为`grep "gpgcheck=" /etc/dnf/dnf.conf`";;
    *)
        echo "没有找到gpgcheck的值，请检查配置文件。"
esac

echo "   repo文件gpgcheck状态"
if [ -e /etc/yum.repos.d/redhat.repo ];then #判断repo文件是否存在
    case `grep gpgcheck /etc/yum.repos.d/redhat.repo | awk -F' = ' '{print $2}' |sed -n '1p'` in
        1)
        echo -e "`grep gpgcheck /etc/yum.repos.d/redhat.repo | awk -F' ' '{print $1$2$3}' |sed -n '1p'`";;
        *)
            sed -i 's/gpgcheck = 0/gpgcheck = 1/g' /etc/yum.repos.d/redhat.repo &&  echo -e "`grep gpgcheck /etc/yum.repos.d/redhat.repo | awk -F' ' '{print $1$2$3}' |sed -n '1p'`"
    esac
else
    echo -e "没有找到repo文件!"
fi
echo "----------------------------------------------"
echo "三、用户账户和环境"
echo
echo "   用户账户密码时效策略"
case $ipass in
    '301147')
        echo -e "PASS_MAX_DAYS=`echo $PASS | awk '{print $1}'`"
        echo -e "PASS_MIN_DAYS=`echo $PASS | awk '{print $2}'`"
        echo -e "PASS_MIN_LEN=`echo $PASS | awk '{print $3}'`"
    echo -e "PASS_WARN_AGE=`echo $PASS | awk '{print $4}'`";;
    *)
        sed -i.bak -e 's/^\(PASS_MAX_DAYS\).*/\1   30/' /etc/login.defs
        sed -i.bak -e 's/^\(PASS_MIN_DAYS\).*/\1   1/' /etc/login.defs
        # 检查是否存在含有PASS_MIN_LEN的行
        if grep -q "^PASS_MIN_LEN" /etc/login.defs; then
            # 如果存在，则修改该行为新的值
            sed -i.bak -e 's/^\(PASS_MIN_LEN\).*/\1   14/' /etc/login.defs
        else
            # 如果不存在，则在PASS_MIN_DAYS行的下一行插入新的行
            sed -i.bak -e '/^PASS_MIN_DAYS/a\PASS_MIN_LEN    14' /etc/login.defs
        fi
        #sed -i.bak -e "/^PASS_MIN_DAYS/a\PASS_MIN_LEN    14" /etc/login.defs
        sed -i.bak -e 's/^\(PASS_WARN_AGE\).*/\1   7/' /etc/login.defs
        PASS=`grep "^PASS" /etc/login.defs | tr -s " " | awk  '{print $2}'`
        echo -e "已修改PASS_MAX_DAYS=`echo $PASS | awk '{print $1}'`"
        echo -e "已修改PASS_MIN_DAYS=`echo $PASS | awk '{print $2}'`"
        echo -e "已修改PASS_MIN_LEN=`echo $PASS | awk '{print $3}'`"
    echo -e "已修改PASS_WARN_AGE=`echo $PASS | awk '{print $4}'`";;
esac

case $INACTIVE in
    30)
    echo -e "`useradd -D | grep INACTIVE`";;
    *)
        useradd -D -f 30
    echo -e "已修改`useradd -D | grep INACTIVE`";;
esac
echo -e "根账户的默认组GID=$GID"
echo "   用户shell超时时间"

if [ -z "$TMOUT" ];then
    printf '%s\n' "# Set TMOUT to 180 seconds" "typeset -xr TMOUT=180" > /etc/profile.d/50-tmout.sh
    typeset -xr TMOUT=180
    echo -e "已设置超时：`echo $TMOUT`"
else
    case $TMOUT in
        "180")
        echo -e "`echo $TMOUT`";;
        *)
            printf '%s\n' "# Set TMOUT to 180 seconds" "typeset -xr TMOUT=180" > /etc/profile.d/50-tmout.sh
        echo -e "已修改`grep "=" /etc/profile.d/50-tmout.sh | awk -F= '{print $2}'`,重新登录后生效。";;
    esac
fi
echo "-------------------------------------------"
echo "四、计划任务"
echo
systemctl is-enabled crond &> /dev/null || systemctl enable crond &> /dev/null
cron=$(systemctl is-enabled crond)   #对Cron赋值
echo "   Cron守护进程启用状态"
echo -e "Cron=$cron"
echo "   Cron权限配置情况"
A=`stat  -c %a%u%g /etc/crontab`
B=`stat  -c %a%u%g /etc/cron.hourly`
C=`stat  -c %a%u%g //etc/cron.daily`
D=`stat  -c %a%u%g /etc/cron.weekly`
E=`stat  -c %a%u%g /etc/cron.monthly`
if [ $A != 60000 ];then
    chmod 600 /etc/crontab && chown root:root /etc/crontab
    echo -e "crontab=`stat  /etc/crontab | sed -n '4p' |tr -s " "`"
else
    echo -e "crontab=`stat  /etc/crontab | sed -n '4p' |tr -s " "`"
fi
if [ $B != 70000 ];then
    chmod 700 /etc/cron.hourly && chown root:root /etc/cron.hourly
    echo -e "cron.hourly=`stat  /etc/cron.hourly | sed -n '4p' |tr -s " "`"
else
    echo -e "cron.hourly=`stat  /etc/cron.hourly | sed -n '4p' |tr -s " "`"
fi
if [ $C != 70000 ];then
    chmod 700 /etc/cron.daily && chown root:root /etc/cron.daily
    echo -e "cron.daily=`stat  /etc/cron.daily | sed -n '4p' |tr -s " "`"
else
    echo -e "cron.daily=`stat  /etc/cron.daily | sed -n '4p' |tr -s " "`"
fi
if [ $D != 70000 ];then
    chmod 700 /etc/cron.weekly && chown root:root /etc/cron.weekly
    echo -e "cron.weekly=`stat  /etc/cron.weekly | sed -n '4p' |tr -s " "`"
else
    echo -e "cron.weekly=`stat  /etc/cron.weekly | sed -n '4p' |tr -s " "`"
fi
if [ $E != 70000 ];then
    chmod 700 /etc/cron.monthly && chown root:root /etc/cron.monthly
    echo -e "cron.monthly=`stat  /etc/cron.monthly | sed -n '4p' |tr -s " "`"
else
    echo -e "cron.monthly=`stat  /etc/cron.monthly | sed -n '4p' |tr -s " "`"
fi
echo "   Cron仅限于授权用户"

for i in cron.deny at.deny cron.allow at.allow #使用循环查询不同配置文件的权限情况，不符合要求的就修改。
do
    if [ -e /etc/$i ];then
        if [ `stat  -c %a /etc/$i` == 600 ];then
            echo -e "$i=`stat /etc/$i | sed -n '4p' | awk '{print $1}'`"
        else
            chmod 600 /etc/$i
            echo -e "已修改$i=`stat /etc/$i | sed -n '4p' | awk '{print $1}'`"
        fi
    else
        echo -e "$i=No such file or directory"
    fi
done
echo "----------------------------------------------"
echo "五、SSH服务器配置"
echo
echo "   sshd_config权限配置情况"
F=`stat  -c %a%u%g /etc/ssh/sshd_config`
if [ $F != 60000 ];then
    chmod 600 /etc/cron.monthly && chown root:root /etc/ssh/sshd_config
    echo -e "sshd_config=`stat  /etc/ssh/sshd_config | sed -n '4p' |tr -s " "`"
else
    echo -e "sshd_config=`stat  /etc/ssh/sshd_config | sed -n '4p' |tr -s " "`"
fi
echo "   sshd_config参数配置情况"
echo -e "SSH protocol=v.2"
case `sshd -T | grep loglevel` in
    "loglevel INFO")
    echo -e "`grep "LogLevel" /etc/ssh/sshd_config | sed 's/ /=/'`";;
    *)
        sed -i '/LogLevel/c LogLevel INFO' /etc/ssh/sshd_config
    echo -e "已修改`grep "LogLevel" /etc/ssh/sshd_config | sed 's/ /=/'`";;
esac
case `grep "X11Forwarding" /etc/ssh/sshd_config | sed -n '1p'` in
    "X11Forwarding no")
    echo -e "`grep "^X11Forwarding" /etc/ssh/sshd_config | sed 's/ /=/'`";;
    *)
        sed -i '/^X11Forwarding/c X11Forwarding no' /etc/ssh/sshd_config
        sed -i '/^#X11Forwarding/c X11Forwarding no' /etc/ssh/sshd_config
    echo -e "已修改`grep "^X11Forwarding" /etc/ssh/sshd_config | sed 's/ /=/'`";;
esac
case `grep "MaxAuthTries" /etc/ssh/sshd_config` in
    "MaxAuthTries 4")
    echo -e "`grep "MaxAuthTries" /etc/ssh/sshd_config | sed 's/ /=/'`";;
    *)
        sed -i '/MaxAuthTries/c MaxAuthTries 4' /etc/ssh/sshd_config
    echo -e "已修改`grep "MaxAuthTries" /etc/ssh/sshd_config | sed 's/ /=/'`";;
esac
case `grep "IgnoreRhosts" /etc/ssh/sshd_config` in
    "IgnoreRhosts yes")
    echo -e "`grep "IgnoreRhosts" /etc/ssh/sshd_config | sed 's/ /=/'`";;
    *)
        sed -i '/IgnoreRhosts/c IgnoreRhosts yes' /etc/ssh/sshd_config
    echo -e "已修改`grep "IgnoreRhosts" /etc/ssh/sshd_config | sed 's/ /=/'`";;
esac
case `grep HostbasedAuthentication /etc/ssh/sshd_config | sed -n '1p'` in
    "HostbasedAuthentication no")
    echo -e "`grep HostbasedAuthentication /etc/ssh/sshd_config | sed -n '1p' | sed 's/ /=/'`";;
    *)
        sed -i '/^#HostbasedAuthentication/c HostbasedAuthentication no' /etc/ssh/sshd_config
    echo -e "已修改`grep HostbasedAuthentication /etc/ssh/sshd_config | sed -n '1p' | sed 's/ /=/'`";;
esac
case `grep "PermitRootLogin" /etc/ssh/sshd_config | sed -n '1p'` in
    "PermitRootLogin no")
    echo -e "`grep "^PermitRootLogin" /etc/ssh/sshd_config | sed 's/ /=/'`";;
    *)
        sed -i '/^PermitRootLogin/c PermitRootLogin no' /etc/ssh/sshd_config
        sed -i '/^#PermitRootLogin/c PermitRootLogin no' /etc/ssh/sshd_config
    echo -e "已修改 `grep "PermitRootLogin" /etc/ssh/sshd_config | sed 's/ /=/'`";;
esac
case `grep "PermitEmptyPasswords" /etc/ssh/sshd_config` in
    "PermitEmptyPasswords no")
    echo -e "`grep "PermitEmptyPasswords" /etc/ssh/sshd_config | sed 's/ /=/'`";;
    *)
        sed -i '/PermitEmptyPasswords/c PermitEmptyPasswords no' /etc/ssh/sshd_config
    echo -e "已修改`grep "PermitEmptyPasswords" /etc/ssh/sshd_config | sed 's/ /=/'`";;
esac
case `grep "PermitUserEnvironment" /etc/ssh/sshd_config` in
    "PermitUserEnvironment no")
    echo -e "`grep "PermitUserEnvironment" /etc/ssh/sshd_config | sed 's/ /=/'`";;
    *)
        sed -i '/PermitUserEnvironment/c PermitUserEnvironment no' /etc/ssh/sshd_config
    echo -e "已修改`grep "PermitUserEnvironment" /etc/ssh/sshd_config | sed 's/ /=/'`";;
esac
case `grep "ClientAliveInterval" /etc/ssh/sshd_config` in
    "ClientAliveInterval 60")
    echo -e "`grep "ClientAliveInterval" /etc/ssh/sshd_config | sed 's/ /=/'`";;
    *)
        sed -i '/ClientAliveInterval/c ClientAliveInterval 60' /etc/ssh/sshd_config
    echo -e "已修改`grep "ClientAliveInterval" /etc/ssh/sshd_config | sed 's/ /=/'`";;
esac
case `grep "ClientAliveCountMax" /etc/ssh/sshd_config` in
    "ClientAliveCountMax 3")
    echo -e "`grep "ClientAliveCountMax" /etc/ssh/sshd_config | sed 's/ /=/'`";;
    *)
        sed -i '/ClientAliveCountMax/c ClientAliveCountMax 3' /etc/ssh/sshd_config
    echo -e "已修改`grep "ClientAliveCountMax" /etc/ssh/sshd_config | sed 's/ /=/'`";;
esac
case `grep "LoginGraceTime" /etc/ssh/sshd_config` in
    "LoginGraceTime 60")
    echo -e "`grep "LoginGraceTime" /etc/ssh/sshd_config | sed 's/ /=/'`";;
    *)
        sed -i '/LoginGraceTime/c LoginGraceTime 60' /etc/ssh/sshd_config
    echo -e "已修改`grep "LoginGraceTime" /etc/ssh/sshd_config | sed 's/ /=/'`";;
esac
echo "---------------------------------------------"
echo "六、PAM配置"
echo
echo "   密码复杂度配置情况"
case `grep -Psi -- '^\h*minlen\h*=\h*(1[4-9]|[2-9][0-9]|[1-9][0-9]{2,})\b'  /etc/security/pwquality.conf.d/50-pwlength.conf` in
    "minlen = 14")
    echo -e "`grep -Psi -- '^\h*minlen\h*=\h*(1[4-9]|[2-9][0-9]|[1-9][0-9]{2,})\b'  /etc/security/pwquality.conf.d/50-pwlength.conf`";;
    *)
        sed -ri 's/^\s*minlen\s*=/# &/' /etc/security/pwquality.conf
        printf '# 密码长度\n%s' "minlen = 14" > /etc/security/pwquality.conf.d/50-pwlength.conf
    echo -e "已修改`grep -Psi -- '^\h*minlen\h*=\h*(1[4-9]|[2-9][0-9]|[1-9][0-9]{2,})\b'  /etc/security/pwquality.conf.d/50-pwlength.conf`";;
esac

if [ ! -e "/etc/security/pwquality.conf.d/50-pwcomplexity.conf" ];then
    # 注释掉/etc/security/pwquality.conf配置文件里的修改参数
    sed -ri 's/^\s*minclass\s*=/# &/' /etc/security/pwquality.conf
    sed -ri 's/^\s*[dulo]credit\s*=/# &/' /etc/security/pwquality.conf
    # 创建配置文件，写入相关配置
    printf '#密码复杂度\n%s\n' "minclass = 4" > /etc/security/pwquality.conf.d/50-pwcomplexity.conf
    printf '\n%s\n' "dcredit = -1" "ucredit = -1" "ocredit = -1" "lcredit = -1" >> /etc/security/pwquality.conf.d/50-pwcomplexity.conf
    # 查询设置好的配置
    grep -Psi -- '^\h*(minclass|[dulo]credit)\b' /etc/security/pwquality.conf.d/50-pwcomplexity.conf
else
    for credit in minclass dcredit ucredit lcredit ocredit
    do
        case `grep "$credit" /etc/security/pwquality.conf.d/50-pwcomplexity.conf` in
            "$credit = -1")
                if [ $credit == "minclass" ];then
                    sed -i "/$credit/c $credit = 4" /etc/security/pwquality.conf.d/50-pwcomplexity.conf
                    echo -e "已修改`grep "$credit" /etc/security/pwquality.conf.d/50-pwcomplexity.conf`"
                else
                    echo -e "`grep "$credit" /etc/security/pwquality.conf.d/50-pwcomplexity.conf`"
                fi
                ;;
            "$credit = 4")
                if [ $credit == "minclass" ];then
                    echo -e "`grep "$credit" /etc/security/pwquality.conf.d/50-pwcomplexity.conf`"
                else
                    sed -i "/$credit/c $credit = -1" /etc/security/pwquality.conf.d/50-pwcomplexity.conf
                    echo -e "已修改`grep "$credit" /etc/security/pwquality.conf.d/50-pwcomplexity.conf`"
                fi
                ;;
            "")
                if [ $credit == "minclass" ];then
                    printf '\n%s\n' "$credit = 4" >> /etc/security/pwquality.conf.d/50-pwcomplexity.conf
                    echo -e "已配置`grep "$credit" /etc/security/pwquality.conf.d/50-pwcomplexity.conf`"
                else
                    printf '\n%s\n' "$credit = -1" >> /etc/security/pwquality.conf.d/50-pwcomplexity.conf
                    echo -e "已配置`grep "$credit" /etc/security/pwquality.conf.d/50-pwcomplexity.conf`"
                fi
                ;;
            *)
                if [ $credit == "minclass" ];then
                    sed -i "/$credit/c $credit = 4" /etc/security/pwquality.conf.d/50-pwcomplexity.conf
                    echo -e "已修改`grep "$credit" /etc/security/pwquality.conf.d/50-pwcomplexity.conf`"
                else
                    sed -i "/$credit/c $credit = -1" /etc/security/pwquality.conf.d/50-pwcomplexity.conf
                    echo -e "已修改`grep "$credit" /etc/security/pwquality.conf.d/50-pwcomplexity.conf`"
                fi
                ;;
        esac
    done
fi

echo "   密码重复受到限制"

case `grep -Pi -- '^\h*remember\h*=\h*(2[4-9]|[3-9][0-9]|[1-9][0-9]{2,})\b' /etc/security/pwhistory.conf` in
    "remember = 24")
    echo -e "`grep -Pi -- '^\h*remember\h*=\h*(2[4-9]|[3-9][0-9]|[1-9][0-9]{2,})\b' /etc/security/pwhistory.conf`";;
    *)
        sed -i 's/^.*remember = .*$/remember = 24/' /etc/security/pwhistory.conf
    echo -e "`grep -Pi -- '^\h*remember\h*=\h*(2[4-9]|[3-9][0-9]|[1-9][0-9]{2,})\b' /etc/security/pwhistory.conf`";;
esac
echo "---------------------------------------------"
echo "七、文件系统权限"
echo
echo "   系统文件权限"
for i in passwd passwd- group group- #使用循环查询不同配置文件的权限情况，不符合要求的就修改。
do
    if [ `stat  -c %a /etc/$i` == 644 ];then
        echo -e "$i=`stat /etc/$i | sed -n '4p' | tr -s " "`"
    else
        chmod 644 /etc/$i
        echo -e "已修改$i=`stat /etc/$i | sed -n '4p' | tr -s " "`"
    fi
done
for i in shadow shadow- gshadow gshadow- #使用循环查询不同配置文件的权限情况，不符合要求的就修改。
do
    if [ `stat  -c %a /etc/$i` == 0 ];then
        echo -e "$i=`stat /etc/$i | sed -n '4p' | tr -s " "`"
    else
        chmod 0000 /etc/$i
        echo -e "已修改$i=`stat /etc/$i | sed -n '4p' | tr -s " "`"
    fi
done
echo "---------------------------------------------"
echo "八、强密码套件"
echo
# 需要检查的策略列表
REQUIRED_POLICIES=(
    "NO-SHA1"
    "NO-WEAKMAC"
    "NO-SSHCBC"
    "NO-SSHCHACHA20"
    "NO-SSHETM"
    "NO-SSHWEAKCIPHERS"
    "NO-SSHWEAKMACS"
)

# 获取当前策略
CURRENT_POLICY=$(update-crypto-policies --show)

# 检查每个策略
ALL_PRESENT=true
for policy in "${REQUIRED_POLICIES[@]}"; do
    if ! grep -qw "$policy" <<< "$CURRENT_POLICY"; then
        echo "未应用策略: $policy"
        # 判断是哪个策略缺失，并添加相应的配置
        case $policy in
            "NO-SHA1")
                # 不存在就创建 NO-SHA1.pmod 文件
                if [ ! -f /etc/crypto-policies/policies/modules/NO-SHA1.pmod ]; then
                    # 创建 NO-SHA1.pmod 文件
                    printf '%s\n' "# This is a subpolicy dropping the SHA1 hash and signature \
                    support" "hash = -SHA1" "sign = -*-SHA1" "sha1_in_certs = 0" > /etc/crypto-policies/policies/modules/NO-SHA1.pmod
                fi
                echo "已添加 NO-SHA1 策略"
                ;;
            "NO-WEAKMAC")
                # 不存在就创建 NO-WEAKMAC.pmod 文件
                if [ ! -f /etc/crypto-policies/policies/modules/NO-WEAKMAC.pmod ]; then
                    # 创建 NO-WEAKMAC.pmod 文件
                    printf '%s\n' "# This is a subpolicy to disable weak macs" "mac = -*-64" > /etc/crypto-policies/policies/modules/NO-WEAKMAC.pmod
                fi
                echo "已添加 NO-WEAKMAC 策略"
                ;;
            "NO-SSHCBC")
                # 不存在就创建 NO-SSHCBC.pmod 文件
                if [ ! -f /etc/crypto-policies/policies/modules/NO-SSHCBC.pmod ]; then
                    # 创建 NO-SSHCBC.pmod 文件
                    printf '%s\n' "# This is a subpolicy to disable all CBC mode ciphers" "# for the SSH protocol \
                    (libssh and OpenSSH)" "cipher@SSH = -*-CBC" > /etc/crypto-policies/policies/modules/NO-SSHCBC.pmod
                fi
                echo "已添加 NO-SSHCBC 策略"
                ;;
            "NO-SSHCHACHA20")
                # 不存在就创建 NO-SSHCHACHA20.pmod 文件
                if [ ! -f /etc/crypto-policies/policies/modules/NO-SSHCHACHA20.pmod ]; then
                    # 创建 NO-SSHCHACHA20.pmod 文件
                    printf '%s\n' "# This is a subpolicy to disable the chacha20-poly1305 ciphers" "# for the SSH protocol \
                    (libssh and OpenSSH)" "cipher@SSH = -CHACHA20-POLY1305" > /etc/crypto-policies/policies/modules/NO-SSHCHACHA20.pmod
                fi
                echo "已添加 NO-SSHCHACHA20 策略"
                ;;
            "NO-SSHETM")
                # 不存在就创建 NO-SSHETM.pmod 文件
                if [ ! -f /etc/crypto-policies/policies/modules/NO-SSHETM.pmod ]; then
                    # 创建 NO-SSHETM.pmod 文件
                    printf '%s\n' "# This is a subpolicy to disable Encrypt then MAC" "# for the SSH protocol \
                    (libssh and OpenSSH)" "etm@SSH = DISABLE_ETM" > /etc/crypto-policies/policies/modules/NO-SSHETM.pmod
                fi
                echo "已添加 NO-SSHETM 策略"
                ;;
            "NO-SSHWEAKCIPHERS")
                # 不存在就创建 NO-SSHWEAKCIPHERS.pmod 文件
                if [ ! -f /etc/crypto-policies/policies/modules/NO-SSHWEAKCIPHERS.pmod ]; then
                    # 创建 NO-SSHWEAKCIPHERS.pmod 文件
                    printf '%s\n' "# This is a subpolicy to disable weak ciphers" "# for the SSH protocol \
                    (libssh and OpenSSH)" "cipher@SSH = -3DES-CBC -AES-128-CBC -AES-192-CBC -AES-256-CBC \
                    -CHACHA20-POLY1305" > /etc/crypto-policies/policies/modules/NO-SSHWEAKCIPHERS.pmod
                fi
                echo "已添加 NO-SSHWEAKCIPHERS 策略"
                ;;
            "NO-SSHWEAKMACS")
                # 不存在就创建 NO-SSHWEAKMACS.pmod 文件
                if [ ! -f /etc/crypto-policies/policies/modules/NO-SSHWEAKMACS.pmod ]; then
                    # 创建 NO-SSHWEAKMACS.pmod 文件
                    printf '%s\n' "# This is a subpolicy to disable weak MACs" "# for the SSH protocol \
                    (libssh and OpenSSH)" "mac@SSH = -HMAC-MD5* -UMAC-64* -UMAC-128* -HMAC-SHA1" > /etc/crypto-policies/policies/modules/NO-SSHWEAKMACS.pmod
                fi
                echo "已添加 NO-SSHWEAKMACS 策略"
                ;;
        esac
        ALL_PRESENT=false
    fi
done

if $ALL_PRESENT; then
    echo "所有加密策略均已应用。"
    update-crypto-policies --show
else
    # 如果有缺失的策略，应用默认策略并输出
    update-crypto-policies --set DEFAULT:NO-SHA1:NO-WEAKMAC:NO-SSHCBC:NO-SSHCHACHA20:NO-SSHETM:NO-SSHWEAKCIPHERS:NO-SSHWEAKMACS
    echo "已应用缺失的加密策略。"
    update-crypto-policies --show
fi

echo "---------------------------------------------"
echo "九、时间同步"
echo
if grep -Prs -- '^\h*(server|pool)\h+[^#\n\r]+' /etc/chrony.conf > /dev/null; then
    case `grep -Prs -- '^\h*(server|pool)\h+[^#\n\r]+' /etc/chrony.conf | sed -n '1p' | awk '{print $2}'` in
        "$NTPServer")
            echo -e "NTP服务器:`grep -Prs -- '^\h*(server|pool)\h+[^#\n\r]+' /etc/chrony.conf | sed -n '1p'`"
            # 重启chronyd服务
            systemctl restart chronyd;;
        *)
            sed -i "/^server /c server $NTPServer iburst" /etc/chrony.conf
            sed -i "/^pool /c server $NTPServer iburst" /etc/chrony.conf
            echo -e "已修改NTP服务器:`grep -Prs -- '^\h*(server|pool)\h+[^#\n\r]+' /etc/chrony.conf | sed -n '1p'`"
            # 重启chronyd服务
            systemctl restart chronyd;;
    esac 
else
    echo "未配置时间同步服务器，请检查 chrony.conf 文件。"
fi
echo
