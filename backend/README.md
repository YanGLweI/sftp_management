# SFTP Backend

基于 Go 语言开发的 SFTP 账号管理系统后端，提供 SFTP 用户管理、文件操作、系统安全加固检查、自动化报告等功能。

## 功能特性

- **SFTP 用户管理**：支持密码、密钥、混合三种认证方式的用户增删改查，支持批量删除
- **SFTP 文件管理**：通过 Web 接口实现文件浏览、上传、下载（含目录打包）、重命名、删除等操作
- **LDAP 集成认证**：对接企业 LDAP/AD 域，实现统一身份认证
- **JWT 鉴权**：基于 JWT Token 的 API 认证机制
- **WebSocket 实时通信**：支持实时日志推送等场景
- **计划任务调度**：
  - 卡巴斯基安全报告自动发送
  - 卡巴斯基隔离区定时检查
  - 系统自动更新
  - 系统安全加固检查
  - 更新报告 / 加固报告定时发送
- **操作日志审计**：记录用户登录及所有操作行为
- **邮件通知**：账号创建后自动发送 SFTP 连接信息到用户邮箱
- **RSA 加密**：敏感密码传输采用 RSA 非对称加密
- **优雅启停**：支持 HTTP 服务和调度器的优雅关闭

## 技术栈

| 组件 | 技术 |
|------|------|
| Web 框架 | [Gin](https://github.com/gin-gonic/gin) |
| ORM | [GORM](https://gorm.io/) (MariaDB/MySQL) |
| 认证 | JWT (`dgrijalva/jwt-go`) |
| LDAP | `go-ldap/ldap/v3` |
| SFTP 客户端 | `pkg/sftp` |
| WebSocket | `gorilla/websocket` |
| 定时任务 | `robfig/cron/v3` |
| 邮件 | `gomail.v2` |
| 配置 | `gopkg.in/yaml.v3` |
| 日志 | `sirupsen/logrus` |

## 项目结构

```
.
├── main.go              # 程序入口
├── config.yml           # 配置文件
├── config/              # 配置加载
├── controller/          # 控制器层（API 处理）
│   ├── login_controller.go
│   ├── user_controller.go
│   ├── sftp_controller.go
│   ├── dashboard_controller.go
│   ├── log_controller.go
│   ├── contact_controller.go
│   └── system_controller.go
├── models/              # 数据模型层（业务逻辑）
├── dao/                 # 数据库连接与初始化
├── routers/             # 路由注册
├── middleware/          # JWT 认证中间件
├── scheduler/           # 计划任务（定时调度）
├── script/              # Shell 脚本（用户管理、系统检查等）
├── utils/               # 工具包（SFTP 连接池、邮件发送）
├── sshutils/            # SSH 工具函数
├── jwt/                 # JWT 生成与解析
├── tools/               # 通用工具（RSA 加密、分页）
├── graceful/            # 优雅关闭
├── kaspersky/           # 卡巴斯基安全报告解析
├── report/              # 报告生成（加固报告、更新报告）
├── common/              # 初始化数据
├── key/                 # RSA 密钥对
└── certificate/         # CA 证书
```

## 快速开始

### 环境要求

- Go 1.23+
- MariaDB / MySQL
- Linux 系统（依赖 Shell 脚本管理 SFTP 用户）

### 安装与运行

```bash
# 克隆项目
git clone https://github.com/YanGLweI/sftp_backend.git
cd sftp_backend

# 安装依赖
go mod tidy

# 修改配置文件
vim config.yml

# 编译运行
go build -o sftpbackend .
./sftpbackend
```

### 配置说明

编辑 `config.yml`，主要配置项：

| 配置项 | 说明 |
|--------|------|
| `system.port` | HTTP 服务监听端口（默认 8888） |
| `database.*` | MariaDB/MySQL 数据库连接配置 |
| `ldap.*` | LDAP/AD 域认证配置 |
| `email.*` | SMTP 邮件发送配置 |
| `jwt.*` | JWT 密钥及过期时间 |
| `script.*` | Shell 脚本路径 |
| `scheduler.*` | 计划任务 Cron 表达式 |

## API 模块

| 模块 | 路径前缀 | 说明 |
|------|----------|------|
| 登录认证 | `/api/login` | 用户登录、LDAP 认证 |
| 用户管理 | `/api/user` | SFTP 账号增删改查 |
| SFTP 操作 | `/api/sftp` | 文件浏览、上传下载、目录管理 |
| 看板 | `/api/dashboard` | 数据统计概览 |
| 日志 | `/api/log` | 操作日志查询 |
| 通讯录 | `/api/contact` | 联系人管理 |
| 系统管理 | `/api/system` | 安全加固、系统更新、计划任务管理 |

## License

MIT
