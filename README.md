# SFTP Management System

> 一站式 SFTP 账号管理与服务器运维平台，面向运维与安全团队，提供用户全生命周期管理、文件操作、安全审计、自动化巡检及报告生成等能力。

## 功能概览

| 模块 | 功能 |
|------|------|
| **SFTP 用户管理** | 创建/更新/删除/批量删除 SFTP 账号，支持密码、密钥、混合三种认证方式，可选密码过期策略 |
| **SFTP 文件管理** | 在线浏览目录、上传/下载文件，基于 `pkg/sftp` 实现远程文件操作 |
| **LDAP/AD 域认证** | 集成 LDAP/Active Directory 登录，支持 TLS 加密连接与安全组权限控制 |
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
│   │   ├── login_controller.go     #   登录认证
│   │   ├── user_controller.go      #   SFTP 用户管理
│   │   ├── sftp_controller.go      #   SFTP 文件操作
│   │   ├── dashboard_controller.go #   数据看板
│   │   ├── log_controller.go       #   日志查询
│   │   ├── contact_controller.go   #   通讯录
│   │   └── system_controller.go    #   系统安全/更新/调度
│   ├── models/                     # 数据模型与业务逻辑
│   ├── routers/                    # 路由注册
│   ├── middleware/                 # JWT 鉴权中间件
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
        │   ├── login/              # 登录页
        │   ├── dashboard/          # 数据看板（统计卡片 + 图表）
        │   ├── sftp/               # SFTP 管理
        │   │   ├── SftpUser/       #   用户管理
        │   │   ├── Sftplog/        #   操作日志
        │   │   └── Contacts/       #   通讯录
        │   ├── file/               # 文件管理
        │   ├── systemSecurity/     # 系统安全
        │   │   ├── Antivirus/      #   卡巴斯基监控
        │   │   ├── SystemHardening/#   系统加固
        │   │   └── SystemUpdate/   #   系统更新
        │   └── acl/                # 权限管理
        └── ...
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

`backend/config.yml` 包含以下配置段：

| 配置段 | 说明 |
|--------|------|
| `system` | 监听端口、运行模式、RSA 私钥路径 |
| `ldap` | LDAP 服务器地址、Base DN、TLS、安全组 |
| `database` | MariaDB 连接信息 |
| `email` | SMTP 邮件发送配置 |
| `jwt` | Token 密钥、过期时间 |
| `script` | Shell 脚本路径 |
| `logfiles` | SFTP 日志路径 |
| `scheduler` | 计划任务初始 Cron 表达式 |

## 系统截图

> 待补充

## License

MIT
