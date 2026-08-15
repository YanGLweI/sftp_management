# SFTP Management System

> 一站式 SFTP 账号管理与服务器运维平台，面向运维与安全团队，提供用户全生命周期管理、文件操作、安全审计、自动化巡检及报告生成等能力。

## 功能概览

| 模块 | 功能 |
|------|------|
| **SFTP 用户管理** | 创建/更新/删除/批量删除 SFTP 账号，支持密码、密钥、混合三种认证方式，可选密码过期策略 |
| **RBAC 权限管理** | 🌟 基于角色的访问控制：角色列表、菜单权限树、LDAP 安全组绑定、超级管理员保护机制 |
| **本地账号认证** | 集成本地用户管理系统，支持 PAM 认证验证、Shell 合法性校验、与 LDAP 双因子登录 |
| **密码策略管理** | 强密码规则强制实施：最小长度、大小写字母、数字、特殊字符、历史密码防复用、过期提醒 |
| **文件传输队列** | 🚀 专业级文件操作：拖拽上传、批量下载、递归搜索、传输进度跟踪、失败重试机制 |
| **SFTP 浏览器** | 🌟 专业级文件传输体验：双面板布局、拖拽上传、传输队列、递归搜索、键盘导航、右键菜单、批量操作等 |
| **LDAP/AD 域认证** | 集成 LDAP/Active Directory 登录，支持 TLS 加密连接与安全组权限控制（数据库配置管理）|
| **本地 PAM 认证** | 支持 RHEL 系统本地账号 PAM 认证，校验 Shell 合法性 |
| **操作日志审计** | 记录登录日志与操作日志（增/删/改），支持按时间、用户名模糊检索与分页 |
| **数据看板** | 首页展示近 7 天文件传输量与访问次数统计（ECharts 柱状图），账号总数与月度新增统计 |
| **通讯录管理** | 维护邮件联系人，支持增删改查与批量操作，用于报告发送 |
| **卡巴斯基安全监控** | 采集 Kaspersky Endpoint Security 状态，生成反病毒库/扫描任务/威胁报告（HTML 邮件） |
| **系统加固检查** | 自动巡检 SSH 配置、密码策略、Cron/At 任务、DNF GPG 校验、系统加密策略等 50+ 项安全基线 |
| **系统更新管理** | 自动执行 DNF 更新，记录更新状态/耗时/详情，支持历史查看 |
| **计划任务调度** | 基于 Cron 的可视化任务调度，支持动态修改执行周期、启用/禁用、立即执行 |
| **安全报告推送** | 定时生成加固报告、更新报告、卡巴斯基报告，以 HTML 邮件形式推送给指定收件人 |
| **WebSocket 实时推送** | 后端任务执行状态通过 WebSocket 实时推送到前端页面 |
| **JWT 认证** | 基于 JWT 的 API 鉴权，支持 Token 过期控制 |
| **RSA 加密传输** | 前端密码经 RSA 加密后传输，后端解密，防止明文泄露 |
| **优雅关闭** | 支持信号监听，服务优雅停止，等待在途请求处理完毕 |

## 技术栈

### 后端

| 技术 | 版本 | 用途 |
|------|------|------|
| Go | 1.23 | 编程语言 |
| Gin | v1.10 | HTTP Web 框架 |
| GORM | v1.25 | ORM 框架 |
| MariaDB | - | 关系型数据库 |
| gorilla/websocket | v1.5 | WebSocket 通信 |
| go-ldap/ldap | v3 | LDAP/AD 集成 |
| pkg/sftp | v1.13 | SFTP 文件操作 |
| robfig/cron | v3 | 定时任务调度 |
| golang-asm/pam | v1.2 | PAM 本地认证 |
| dgrijalva/jwt-go | v3 | JWT 鉴权 |
| gomail | v2 | 邮件发送 |
| logrus | v1.9 | 日志框架 |
| gopkg.in/yaml | v3 | 配置文件解析 |

### 前端

| 技术 | 版本 | 用途 |
|------|------|------|
| Vue | 2.6 | 前端框架 |
| Element UI | 2.15 | UI 组件库 |
| Vue Router | 3.0 | 路由管理 |
| Vuex | 3.1 | 状态管理 |
| Axios | 0.18 | HTTP 请求 |
| ECharts | 4.9 | 数据可视化图表 |
| JSEncrypt | 3.5 | RSA 加密 |
| dayjs | 1.11 | 日期处理 |
| Sass | - | CSS 预处理 |

## 项目结构

```
sftp_management/
├── backend/                        # Go 后端
│   ├── main.go                     # 入口：数据库连接、迁移、调度器、HTTP 服务
│   ├── config.yml                  # 全局配置文件
│   ├── config/                     # 配置加载与解析
│   ├── dao/                        # 数据库连接（GORM + MariaDB）
│   ├── common/                     # 初始化数据
│   ├── controller/                 # 控制器层（请求处理）
│   │   ├── login_controller.go     #   登录认证（支持 LDAP+ 本地账号 + 双控凭证）
│   │   ├── user_controller.go      #   SFTP 用户管理
│   │   ├── localuser_controller.go #   本地用户管理
│   │   ├── role_controller.go      #   RBAC 角色权限管理
│   │   ├── password_policy_controler.go # 密码策略管理
│   │   ├── ldap_config_controller.go   # LDAP 配置管理（数据库持久化）
│   │   ├── sftp_module_controller.go   # SFTP 模块配置（标签上传/中国联通）
│   │   ├── sftp_controller.go      #   SFTP 文件操作
│   │   ├── dashboard_controller.go #   数据看板
│   │   ├── log_controller.go       #   平台日志查询
│   │   ├── contact_controller.go   #   通讯录
│   │   └── system_controller.go    #   系统安全/更新/调度
│   ├── models/                     # 数据模型（业务逻辑 + ORM 映射）
│   │   ├── user.go                 # SFTP 用户
│   │   ├── localuser.go            # 本地账号
│   │   ├── role.go                 # 角色与菜单/LDAP 安全组关联
│   │   ├── password_policy.go      # 密码策略
│   │   ├── ldap_config.go          # LDAP 配置表（含证书文件名持久化）
│   │   ├── sftp_module_config.go   # SFTP 模块动态配置表
│   │   ├── sftp.go                 # SFTP 文件操作相关
│   │   ├── sftplog.go              # SFTP 操作日志
│   │   ├── contact.go              # 通讯录
│   │   ├── log.go                  # 平台操作日志
│   │   ├── scheduler.go            # 定时任务配置
│   │   ├── systemcheck.go          # 系统安全加固标准
│   │   ├── update.go               # 系统更新记录
│   ├── routers/                    # 路由注册
│   ├── middleware/                 # JWT 鉴权中间件 + 双控凭证校验
│   │   ├── jwtAuth.go              # JWT Token 解析与验签
│   │   └── dualAuth.go             # 双控账号验证逻辑（临时 Token 管理）
│   ├── scheduler/                  # 定时任务实现
│   ├── kaspersky/                  # 卡巴斯基安全监控
│   ├── report/                     # 报告生成（加固/更新）
│   ├── script/                     # Shell 脚本（用户管理/系统检查）
│   ├── jwt/                        # JWT 工具
│   ├── tools/                      # 通用工具（分页、RSA）
│   ├── utils/                      # 工具函数（邮件、SFTP 连接池）
│   ├── sshutils/                   # SSH 工具
│   ├── graceful/                   # 优雅关闭
│   ├── key/                        # RSA 密钥对
│   └── certificate/                # CA 证书
│
└── frontend/                       # Vue 前端
    └── src/
        ├── views/
        │   ├── login/              # 登录页（RSA 加密 + 双因子）
        │   ├── dashboard/          # 数据看板（统计卡片 + ECharts 图表）
        │   ├── sftp/               # SFTP 管理
        │   │   ├── SftpUser/       #   用户管理（增删改查 + 批量删除）
        │   │   ├── Sftplog/        #   平台操作日志
        │   │   └── Contacts/       #   通讯邮箱管理
        │   ├── file/               # 文件管理（临时文件预览）
        │   ├── log/                # SFTP 日志（按日切割文件查询）
        │   ├── systemSecurity/     # 系统安全
        │   │   ├── Antivirus/      #   卡巴斯基监控（威胁报告/隔离区检查）
        │   │   ├── SystemHardening/#   系统加固检查（50+ 项安全基线）
        │   │   └── SystemUpdate/   #   系统自动更新（DNF 包管理器）
        │   └── settings/           # 平台设置（RBAC + LDAP + 策略）
        │       ├── Role/           #   角色管理（菜单权限/SFTP 模块权限/LDAP 安全组）
        │       ├── LocalUser/      #   本地账号（PAM 认证/Shell 校验）
        │       ├── PasswordPolicy/ #   密码策略（强度规则/历史记录）
        │       └── LDAPManagement/ #   LDAP 配置（证书上传/绑定 DN/安全组过滤）
        └── components/
            ├── SftpBrowser/        # SFTP 浏览器核心组件（双面板布局/拖拽/队列）
            ├── DualVerify/         # 双控账号验证组件（Token 临时登录）
            └── ChangePasswordDialog/ # 改密弹窗（旧密码验证/新密码加密提交）
```

## 快速开始

### 环境要求

- Go >= 1.23
- Node.js >= 16
- MariaDB / MySQL
- Linux 操作系统（SFTP 用户管理依赖 Shell 脚本）

### 后端启动

```bash
cd backend

# 安装依赖
go mod tidy

# 修改配置文件
vim config.yml

# 编译运行
go build -o sftpbackend .
./sftpbackend
```

后端默认监听 `:8888`，可在 `config.yml` 中修改。

### 前端启动

```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

前端默认监听 `:9529`。

### 配置说明

`backend/config.yml` 包含以下配置段:

| 配置段 | 说明 |
|--------|------|
| `system` | 监听端口、运行模式、RSA 私钥路径、加密私钥路径 |
| `sftp_account` | 专用 SFTP 账号（标签上传/中国联通等模块共用） |
| `hotlabel` | 标签上传允许访问的根路径限制 |
| `chinaunicom` | 中国联通允许访问的根路径限制 |
| `ldap` | LDAP 服务器地址、Base DN、TLS、安全组 |
| `database` | MariaDB 连接信息、表前缀、编码方式 |
| `email` | SMTP 邮件发送配置（含 HTML 正文模板） |
| `jwt` | Token 密钥、过期时间、签发人 |
| `script` | Shell 脚本路径（用户管理/系统检查/文件统计） |
| `logfiles` | SFTP 日志路径（总日志 + 每日切割文件） |
| `scheduler` | 计划任务初始 Cron 表达式（卡巴斯基/系统更新/加固检查） |

## SFTP 浏览器功能说明

🌟 **SFTP Browser** 是本平台的核心文件传输组件，提供专业级的文件管理体验。

### 核心特性

#### 🖥️ 双面板布局
- **左侧面板**: 文件列表浏览区（支持动态高度、磨砂玻璃风格）
- **右侧面板**: 传输队列管理区（列队/失败/成功三个标签页）
- **拖拽分隔线**: 可自由调整左右面板宽度，折叠/恢复交互流畅

#### 🚀 上传与下载
- **批量上传**: 支持多文件同时选择上传，显示实时进度条
- **拖拽上传**: 拖入文件至对应卡片自动加入传输队列
- **分栏上传**: 右键菜单支持选定单个文件或全部上传
- **批量下载**: 多选文件打包压缩下载
- **单文件下载**: 右键菜单快速下载

#### 🔍 搜索功能
- **即时过滤**: 输入框实时过滤当前目录文件
- **递归搜索**: 深度扫描指定路径下的所有匹配文件
- **骨架屏预占位**: 搜索前显示骨架屏避免页面跳动

#### ⌨️ 键盘导航
- **方向键**: ↑↓ 循环聚焦文件列表项
- **回车/双击**: 打开文件或进入目录
- **Delete**: 删除选中文件（支持批量删除）
- **Backspace**: 返回上级目录
- **Tab**: 焦点在行项目间循环切换

#### 🎯 右键菜单
- **文件操作**: 下载、重命名、删除、复制文件名
- **队列操作**: 选定上传、移除、全部上传、清空
- **失败重试**: 单个/全部重试、清空失败记录
- **成功清理**: 清空成功传输记录

#### 📋 传输队列管理
- **待上传队列**: 显示已添加待上传文件，支持右键控制
- **失败记录**: 显示失败原因、文件大小、时间戳
- **成功记录**: 显示成功传输的文件及完成时间

### 技术亮点

- **动态表格高度**: 根据可用空间自动计算表高，自适应屏幕变化
- **磨砂玻璃效果**: UI 采用半透明模糊背景，提升视觉质感
- **防抖动优化**: 隐藏标签页主动触发 `doLayout` 避免表格抖动
- **平滑动效**: Element UI 标签指示器滑动过渡
- **上下文聚焦**: 当前聚焦行高亮显示，视觉引导清晰
- **事件阻止**: 完善的拖拽事件处理，防止默认行为干扰

### 使用场景

1. **运维人员日常维护**: 高效管理服务器文件结构
2. **安全团队审计**: 批量上传/下载配置文件与日志
3. **开发部署流程**: 快速部署应用包与资源文件
4. **故障排查**: 递归搜索关键文件定位问题

---

## 系统截图

> 待补充

## License

MIT
