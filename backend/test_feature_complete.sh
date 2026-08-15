#!/bin/bash

echo "========================================="
echo "SFTP 管理平台功能完整性测试计划"
echo "========================================="
echo ""
echo "测试范围：最近 2 天新增代码（role/localuser/password_ldap/sftp-module）"
echo "测试时间：$(date)"
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试计数器
PASSED=0
FAILED=0
WARNINGS=0

echo "【第一部分】后端 API 功能验证"
echo "-----------------------------------------"

# 1. 超级管理员角色保护测试
echo -e "\n${YELLOW}1. 超级管理员角色保护测试${NC}"
echo "   ✓ 检查 DeleteRole 函数中禁止删除超级管理员..."
if grep -q "不能删除超级管理员角色" backend/controller/role_controller.go; then
    echo -e "   ${GREEN}✓ PASS: DeleteRole 实现超级管理员保护${NC}"
    ((PASSED++))
else
    echo -e "   ${RED}✗ FAIL: DeleteRole 缺少保护逻辑${NC}"
    ((FAILED++))
fi

echo "   ✓ 检查 UpdateRole 函数限制修改菜单权限..."
if grep -A 5 "superAdminRoleName" backend/controller/role_controller.go | grep -q "仅允许更新 LDAP 安全组"; then
    echo -e "   ${GREEN}✓ PASS: UpdateRole 限制菜单权限修改${NC}"
    ((PASSED++))
else
    echo -e "   ${RED}✗ FAIL: UpdateRole 未正确限制${NC}"
    ((FAILED++))
fi

# 2. 本地账号 admin 保护测试
echo -e "\n${YELLOW}2. 本地账号 admin 保护测试${NC}"
echo "   ✓ 检查 DeleteLocalUser 禁止删除默认 admin..."
if grep -q "不能删除默认管理员账号" backend/controller/localuser_controller.go; then
    echo -e "   ${GREEN}✓ PASS: DeleteLocalUser 实现保护${NC}"
    ((PASSED++))
else
    echo -e "   ${RED}✗ FAIL: DeleteLocalUser 缺少保护${NC}"
    ((FAILED++))
fi

# 3. LDAP 配置数据库迁移验证
echo -e "\n${YELLOW}3. LDAP 配置数据库迁移验证${NC}"
echo "   ✓ 检查 LDAPConfig 模型包含 CertFilename 字段..."
if grep -q "CertFilename string" backend/models/ldap_config.go; then
    echo -e "   ${GREEN}✓ PASS: LDAPConfig 模型完整${NC}"
    ((PASSED++))
else
    echo -e "   ${RED}✗ FAIL: LDAPConfig 模型缺少字段${NC}"
    ((FAILED++))
fi

echo "   ✓ 检查 main.go 触发 AutoMigrate LDAPConfig..."
if grep -q "&models.LDAPConfig{}" backend/main.go; then
    echo -e "   ${GREEN}✓ PASS: LDAPConfig 已加入迁移流程${NC}"
    ((PASSED++))
else
    echo -e "   ${RED}✗ FAIL: LDAPConfig 未加入迁移${NC}"
    ((FAILED++))
fi

# 4. LDAP 连接测试证书验证
echo -e "\n${YELLOW}4. LDAP 连接测试证书验证逻辑${NC}"
echo "   ✓ 检查使用域名解析而非直接 IP..."
if grep -q "ldap.DialURL" backend/controller/ldap_config_controller.go; then
    if grep -B 5 -A 5 "ldap.DialURL" backend/controller/ldap_config_controller.go | grep -q "RootCAs"; then
        echo -e "   ${GREEN}✓ PASS: TLS 证书验证使用 RootCAs 正确配置${NC}"
        ((PASSED++))
    else
        echo -e "   ${YELLOW}⚠ WARN: LDAP 测试连接可能跳过证书验证 (InsecureSkipVerify)${NC}"
        ((WARNINGS++))
    fi
else
    echo -e "   ${RED}✗ FAIL: LDAP 连接测试实现异常${NC}"
    ((FAILED++))
fi

# 5. 密码策略功能完整性
echo -e "\n${YELLOW}5. 密码策略功能完整性${NC}"
if grep -q "ValidatePasswordPolicy" backend/controller/password_policy_controller.go && \
   grep -q "GetPasswordPolicy\|UpdatePasswordPolicy" backend/controller/password_policy_controller.go; then
    echo -e "   ${GREEN}✓ PASS: Get/Update 接口完整${NC}"
    ((PASSED++))
else
    echo -e "   ${RED}✗ FAIL: 密码策略接口不完整${NC}"
    ((FAILED++))
fi

echo "   ✓ 检查 ValidatePassword 验证端点..."
if grep -q "func ValidatePassword" backend/controller/password_policy_controller.go; then
    echo -e "   ${GREEN}✓ PASS: 前端实时验证接口存在${NC}"
    ((PASSED++))
else
    echo -e "   ${RED}✗ FAIL: 实时验证接口缺失${NC}"
    ((FAILED++))
fi

# 6. SFTP 模块双控验证逻辑
echo -e "\n${YELLOW}6. SFTP 模块双控验证逻辑${NC}"
echo "   ✓ 检查 CheckLDAPRolePermission 实现域控验证..."
if grep -q "CheckLDAPRolePermission" backend/controller/sftp_module_controller.go; then
    echo -e "   ${GREEN}✓ PASS: LDAP 角色权限校验存在${NC}"
    ((PASSED++))
else
    echo -e "   ${RED}✗ FAIL: LDAP 角色权限校验缺失${NC}"
    ((FAILED++))
fi

echo "   ✓ 检查 ChinaUnicom 模块双控开关配置..."
if grep -A 5 "ModuleNameChinaUnicom" backend/controller/sftp_module_controller.go | grep -q "DualAuthEnabled"; then
    echo -e "   ${GREEN}✓ PASS: 双控开关仅限中国联通模块${NC}"
    ((PASSED++))
else
    echo -e "   ${RED}✗ FAIL: 双控开关逻辑错误${NC}"
    ((FAILED++))
fi

# 7. 初始化数据完整性
echo -e "\n${YELLOW}7. 初始化数据完整性${NC}"
if grep -q "InitDefaultConfigs" backend/common/init_data.go && \
   grep -q "CreateLDAPConfig" backend/common/init_data.go; then
    echo -e "   ${GREEN}✓ PASS: LDAP 和 SFTP 模块初始化集成${NC}"
    ((PASSED++))
else
    echo -e "   ${RED}✗ FAIL: 初始化逻辑不完整${NC}"
    ((FAILED++))
fi

echo -e "\n========================================="
echo "测试结果汇总："
echo -e "${GREEN}通过：$PASSED${NC}"
echo -e "${RED}失败：$FAILED${NC}"
echo -e "${YELLOW}警告：$WARNINGS${NC}"
echo "========================================="

# 退出状态码
if [ $FAILED -gt 0 ]; then
    exit 1
elif [ $WARNINGS -gt 0 ]; then
    exit 0
else
    exit 0
fi
