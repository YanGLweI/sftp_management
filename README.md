# SFTP Management System

SFTP 账号管理系统，包含前端和后端。

## 项目结构

```
sftp_management/
├── backend/    # Go 后端（Gin + GORM）
└── frontend/   # Vue 前端（Vue 2 + Element UI）
```

## 后端（backend）

基于 Go + Gin 框架开发的 SFTP 账号管理后端服务。

### 主要功能

- SFTP 用户管理（密码/密钥/混合认证）
- SFTP 文件管理（上传/下载/目录操作）
- LDAP/AD 域集成认证
- 计划任务调度（安全报告、系统加固检查、自动更新）
- 操作日志审计
- 邮件通知
- WebSocket 实时通信

### 快速开始

```bash
cd backend
go mod tidy
vim config.yml   # 修改配置
go build -o sftpbackend .
./sftpbackend
```

详细文档请参考 [backend/README.md](./backend/README.md)

## 前端（frontend）

基于 Vue 2 + Element UI 开发的 SFTP 管理系统前端。

### 快速开始

```bash
cd frontend
npm install
npm run dev
```

## 技术栈

| 端 | 技术 |
|----|------|
| 后端 | Go 1.23、Gin、GORM、MariaDB、JWT、LDAP |
| 前端 | Vue 2、Element UI、Axios |

## License

MIT
