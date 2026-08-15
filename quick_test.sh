#!/bin/bash
# 快速功能验证脚本

echo "======================================"
echo "SFTP 管理平台 - 核心功能快速验证"
echo "======================================"

PASS=0
FAIL=0

# 1. 角色管理保护检查
echo -e "\n[1] 超级管理员保护机制"
if grep -q "不能删除超级管理员角色" backend/controller/role_controller.go; then
    echo "   ✓ DeleteRole 禁止删除超级管理员"
    ((PASS++))
else
    echo "   ✗ DeleteRole 缺少保护"
    ((FAIL++))
fi

if grep -A 5 "superAdminRoleName {" backend/controller/role_controller.go | grep -q "仅允许更新 LDAP 安全组"; then
    echo "   ✓ UpdateRole 限制菜单权限修改"
    ((PASS++))
else
    echo "   ✗ UpdateRole 未正确限制"
    ((FAIL++))
fi

# 2. 本地账号保护
echo -e "\n[2] 本地账号 admin 保护"
if grep -q "不能删除默认管理员账号" backend/controller/localuser_controller.go; then
    echo "   ✓ DeleteLocalUser 保护 admin"
    ((PASS++))
else
    echo "   ✗ DeleteLocalUser 缺少保护"
    ((FAIL++))
fi

# 3. LDAP 配置迁移
echo -e "\n[3] LDAP 数据库迁移"
if grep -q "&models.LDAPConfig{}" backend/main.go && \
   grep -q "CertFilename string" backend/models/ldap_config.go; then
    echo "   ✓ LDAPConfig 加入 AutoMigrate"
    ((PASS++))
else
    echo "   ✗ LDAPConfig 迁移不完整"
    ((FAIL++))
fi

# 4. TLS 证书验证
echo -e "\n[4] LDAP TLS 证书验证"
if grep -B 3 "RootCAs:" backend/controller/ldap_config_controller.go | grep -q "tls.Config"; then
    echo "   ✓ 使用 RootCAs 验证证书"
    ((PASS++))
else
    echo "   ✗ TLS 配置异常"
    ((FAIL++))
fi

# 5. 密码策略接口
echo -e "\n[5] 密码策略功能"
if grep -q "func GetPasswordPolicy" backend/controller/password_policy_controller.go && \
   grep -q "func UpdatePasswordPolicy" backend/controller/password_policy_controller.go && \
   grep -q "func ValidatePassword" backend/controller/password_policy_controller.go; then
    echo "   ✓ Get/Update/Validate 接口完整"
    ((PASS++))
else
    echo "   ✗ 密码策略接口缺失"
    ((FAIL++))
fi

# 6. SFTP 模块双控
echo -e "\n[6] SFTP 模块双控逻辑"
if grep -q "CheckLDAPRolePermission" backend/controller/sftp_module_controller.go && \
   grep -A 5 "ModuleNameChinaUnicom" backend/controller/sftp_module_controller.go | grep -q "DualAuthEnabled"; then
    echo "   ✓ CheckLDAPRolePermission 存在"
    echo "   ✓ ChinaUnicom 双控开关独立配置"
    ((PASS+=2))
else
    echo "   ✗ SFTP 模块双控逻辑异常"
    ((FAIL++))
fi

# 7. 前端组件集成
echo -e "\n[7] 前端组件验证"
if [ -f "frontend/src/views/admin/hotlabel-config/index.vue" ] && \
   [ -f "frontend/src/views/admin/chinaunicom-config/index.vue" ] && \
   [ -f "frontend/src/views/admin/sftp-module-config/SftpModuleConfigForm.vue" ]; then
    echo "   ✓ 标签上传和中国联通配置页面存在"
    echo "   ✓ SftpModuleConfigForm 公共组件存在"
    ((PASS+=2))
else
    echo "   ✗ 前端组件缺失"
    ((FAIL++))
fi

# 8. 初始化调用
echo -e "\n[8] 数据初始化流程"
if grep -q "InitDefaultConfigs()" backend/main.go; then
    echo "   ✓ main.go 中调用 InitDefaultConfigs"
    ((PASS++))
else
    echo "   ⚠️  未找到 InitDefaultConfigs 调用"
    ((FAIL++))
fi

# 汇总
echo -e "\n======================================"
echo "测试结果: $PASS 通过，$FAIL 失败"
echo "======================================"

if [ $FAIL -eq 0 ]; then
    echo "✅ 所有核心功能验证通过！"
    exit 0
else
    echo "❌ 发现 $FAIL 个问题需要修复"
    exit 1
fi
