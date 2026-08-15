# LDAP Config CertFilename 字段补充说明

## 📝 修改内容

### 1. **模型定义** (`models/ldap_config.go`)
```go
type LDAPConfig struct {
    // ... 其他字段 ...
    CertFilename string `json:"cert_filename" gorm:"size:255;"`  // CA 证书文件名
}
```

### 2. **数据库自动迁移** (`main.go`)
```go
err = dao.DB.AutoMigrate(
    // ...
    &models.LDAPConfig{},  // 自动建表 / 自动添加缺失列
)
```
GORM 自动检测模型变化并同步表结构，**无需任何手动 SQL**！

### 3. **初始化数据** (`common/init_data.go`)
```go
// InitData() 函数开头（确保在任何 early return 之前执行）
var ldapConfigCount int64
dao.DB.Model(&models.LDAPConfig{}).Count(&ldapConfigCount)
if ldapConfigCount == 0 {
    if err := models.CreateLDAPConfig("", ""); err != nil {
        logrus.Printf("InitData Create LDAPConfig failed: %v", err)
    } else {
        logrus.Println("InitData Create LDAPConfig success")
    }
}
```

### 4. **Controller 处理** (`controller/ldap_config_controller.go`)
- ✅ 接收 `cert_filename` JSON 参数
- ✅ 保存到 `LDAPConfig.CertFilename` 字段

---

## 🔧 部署步骤

**无需执行任何 SQL 脚本！** 只需：

1. **重启后端服务**（AutoMigrate 自动建表/加列 + InitData 自动初始化数据）
   ```bash
   cd /src/sftp_management/backend
   go build -o sftpbackend main.go
   ./sftpbackend
   ```

2. **查看日志确认**
   ```
   [213.348ms] [rows:0] CREATE TABLE `t_ldap_config` ...   ← AutoMigrate 建表
   InitData Create LDAPConfig success                      ← 初始化数据
   Gin服务启动成功，监听地址: :8888                         ← 服务正常启动
   ```

3. **重新编译前端**
   ```bash
   cd /src/sftp_management/frontend
   npm run build  # 或 npm run dev
   ```

---

## ✅ 验证方法

### **1. 检查表结构和初始化数据**
```sql
DESC t_ldap_config;
-- cert_filename varchar(255) 字段应存在

SELECT * FROM t_ldap_config;
-- 应有一条默认记录：
-- id=1, server='', base_dn='', use_tls=0, insecure=1,
-- user_filter='(sAMAccountName=%s)', cert_filename='', updated_by=0
```

### **2. 测试功能流程**
1. 上传 `dcpm.crt` → 立即显示预览（蓝色 info 框）✅
2. 点击"保存配置" → 调用 `saveLDAPConfig(cert_filename)` ✅
3. **刷新页面** → 显示 `"当前已上传 CA 证书"` + `"dcpm.crt"` ✅

---

## 🎯 设计原则（严禁 SQL）

| 事项 | 处理方式 |
|------|----------|
| **建表 / 加列** | ✅ GORM `AutoMigrate` 自动处理 |
| **初始化数据** | ✅ `InitData()` Go 代码处理（表为空时插入） |
| **默认值** | ✅ Go 代码中设置（如 `UserFilter: "(sAMAccountName=%s)"`） |
| **手动 SQL** | ❌ **严禁使用**（不创建 .sql 脚本，不手动执行） |

> ⚠️ 注意：GORM struct tag 中的 `default` 无法正确处理含括号/特殊字符的默认值（如 `(sAMAccountName=%s)` 会报 SQL 语法错误 1064），因此默认值一律在 Go 代码中设置。

---

## ⚠️ 注意事项

1. **`InitData()` 中 LDAP 初始化必须放在函数开头**
   - 该函数中间存在 `return`（admin 账号已存在时提前返回）
   - 放在开头可确保 LDAP 初始化**始终执行**

2. **`CreateLDAPConfig` 已设置默认值**
   - `UserFilter: "(sAMAccountName=%s)"` ✅
   - `UseTLS: false`、`Insecure: true` ✅
   - `CertFilename` 默认为空字符串（无证书状态）✅

3. **兼容性处理**
   - 旧数据没有 `cert_filename` 记录 → 默认值 `""`
   - 前端回退逻辑：`config.cert_filename || 'ca.crt'`
