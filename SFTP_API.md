# SFTP 平台 API 文档

本文档仅涵盖 SFTP 登录和 SFTP 浏览器相关接口。

## 目录
- [认证与登录](#认证与登录)
- [SFTP 浏览器接口](#sftp 浏览器接口)
- [数据格式](#数据格式)
- [错误码说明](#错误码说明)

---

## 认证与登录

### 1. 获取模块配置

获取可用的 SFTP 登录模块配置（公共接口，无需认证）。

**请求**
```http
GET /sftp/module-configs
```

**响应**
```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "modules": [
      {
        "name": "HOTLABEL",
        "enabled": true,
        "login_type": "ldap",
        "root_path": "/hotlabel"
      },
      {
        "name": "CHINAUNICOM",
        "enabled": true,
        "login_type": "ldap",
        "root_path": "/chinaunicom"
      }
    ]
  }
}
```

---

### 2. SFTP 登录

支持密码登录和密钥文件登录两种方式。

#### 2.1 密码登录

**请求**
```http
POST /sftp/login
Content-Type: application/json

{
  "username": "string",
  "password": "string",       // RSA 加密后的密码
  "loginType": "normal|hotlabel|chinaunicom"
}
```

**参数说明**
- `username`: 用户名
- `password`: 密码（需使用 RSA 公钥加密后传输）
- `loginType`: 
  - `normal`: 普通密码登录
  - `hotlabel`: HOTLABEL 模块（走 LDAP 域控验证 + 角色权限校验）
  - `chinaunicom`: 中国联通模块（走 LDAP 域控验证 + 角色权限校验）

**响应（成功）**
```json
{
  "code": 200,
  "message": "SFTP 登录成功",
  "data": {
    "sftp_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expire_tips": "连接有效期 8 小时，超时需重新登录"
  }
}
```

**响应（需修改密码）**
```json
{
  "code": 200,
  "message": "该账号需先修改密码",
  "data": {
    "must_change_password": true,
    "change_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**响应（失败）**
```json
{
  "code": 400,
  "message": "<密码错误> invalid user"
}
```

#### 2.2 密钥文件登录

**请求**
```http
POST /sftp/login
Content-Type: multipart/form-data

username: string
file: <RSA 私钥文件>
```

**响应（成功）**
```json
{
  "code": 200,
  "message": "SFTP 登录成功",
  "data": {
    "sftp_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expire_tips": "连接有效期 8 小时，超时需重新登录"
  }
}
```

---

### 3. 双控验证

仅中国联通登录的连接需要双控验证。验证通过后签发短期双控凭证。

**请求**
```http
POST /sftp/dualverify
Content-Type: application/json
X-SFTP-Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

{
  "username": "string",
  "password": "string"        // RSA 加密后的密码
}
```

**参数说明**
- `X-SFTP-Token`: 登录返回的 Token
- `username`: 双控复核人用户名
- `password`: 双控复核人密码（RSA 加密）

**响应（成功）**
```json
{
  "code": 200,
  "message": "双控验证通过",
  "data": {
    "dual_token": "dual_eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**响应（当前连接无需双控验证）**
```json
{
  "code": 400,
  "message": "当前连接无需双控验证"
}
```

**响应（双控验证失败）**
```json
{
  "code": 400,
  "message": "双控验证失败：invalid user"
}
```

---

### 4. SFTP 登出

删除 Token 和连接。

**请求**
```http
GET /sftp/logout?sftp_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

或使用 Header
```http
GET /sftp/logout
X-SFTP-Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应**
```json
{
  "code": 200,
  "message": "SFTP 登出成功"
}
```

---

## SFTP 浏览器接口

所有 SFTP 浏览器接口都需要携带 Token（优先使用 Header，其次 Query/Body）。

### Token 传递方式

**推荐方式（Header）**
```http
X-SFTP-Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**备选方式（Query）**
```http
GET /sftp/files?sftp_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**备选方式（Form）**
```http
POST /sftp/upload
X-SFTP-Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

### 5. 获取文件列表

获取指定路径下的文件和目录列表。

**请求**
```http
GET /sftp/files?path=/current/dir
```

**参数说明**
- `path`: 远程路径，默认为 `/`

**响应**
```json
{
  "code": 200,
  "message": "文件列表获取成功",
  "data": {
    "path": "/current/dir",
    "description": "3 个文件 和 2 个目录",
    "files": [
      {
        "name": "documents",
        "path": "/current/dir/documents",
        "isDir": true,
        "size": 0,
        "modified": "2026-08-25T10:30:00+08:00"
      },
      {
        "name": "report.pdf",
        "path": "/current/dir/report.pdf",
        "isDir": false,
        "size": 102400,
        "modified": "2026-08-25T09:15:00+08:00"
      }
    ]
  }
}
```

---

### 6. 上传文件

支持单文件和多文件批量上传，流式传输，进度条反映真实速度。

**请求**
```http
POST /sftp/upload
X-SFTP-Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: multipart/form-data

path: /remote/destination/path
file: <file1.txt>
file: <file2.txt>
```

**参数说明**
- `path`: 目标远程路径（必填）
- `file`: 上传的文件（支持多个，必填）

**安全特性**
- 文件名合法性校验（防止 `../` 路径穿越）
- 路径边界检查（防止超出允许的根路径）
- 流式直写 SFTP（避免整体缓冲到内存）

**响应（成功）**
```json
{
  "code": 200,
  "message": "文件上传成功"
}
```

**响应（失败）**
```json
{
  "code": 400,
  "message": "上传文件名非法：文件名包含非法字符"
}
```

---

### 7. 下载文件

下载单个文件，流式传输。

**请求**
```http
GET /sftp/download?sftp_token=XXX&path=/remote/file.txt
```

**参数说明**
- `path`: 远程文件完整路径（必填）

**响应头**
```http
Content-Type: application/octet-stream
Content-Disposition: attachment; filename="file.txt"; filename*=UTF-8''file.txt
Content-Length: 102400
Cache-Control: no-cache, no-store, must-revalidate
Pragma: no-cache
Expires: 0
```

**响应（成功）**
- 二进制文件流

**响应（失败）**
```json
{
  "code": 404,
  "message": "文件不存在：stat /remote/file.txt: no such file"
}
```

---

### 8. 下载目录

将目录打包为 ZIP 后下载。

**请求**
```http
GET /sftp/downloaddir?sftp_token=XXX&path=/remote/folder
```

**参数说明**
- `path`: 远程目录完整路径（必填）

**响应头**
```http
Content-Type: application/zip
Content-Disposition: attachment; filename="folder.zip"; filename*=UTF-8''folder.zip
Cache-Control: no-cache, no-store, must-revalidate
Pragma: no-cache
Expires: 0
```

**响应（成功）**
- ZIP 压缩文件流

**响应（失败）**
```json
{
  "code": 400,
  "message": "该路径不是目录，请使用文件下载接口"
}
```

---

### 9. 创建目录

创建新目录。

**请求**
```http
POST /sftp/mkdir
X-SFTP-Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
X-Dual-Token: dual_eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...  // 仅中国联通模块需要
Content-Type: application/json

{
  "path": "/parent/directory",
  "name": "new_folder"
}
```

**参数说明**
- `path`: 父目录路径（必填）
- `name`: 新目录名称（必填）

**响应（成功）**
```json
{
  "code": 200,
  "message": "目录创建成功"
}
```

**响应（失败）**
```json
{
  "code": 400,
  "message": "目录名非法：目录名包含非法字符"
}
```

---

### 10. 删除文件或目录

删除指定的文件或目录（支持递归删除目录）。

**请求**
```http
POST /sftp/delete
X-SFTP-Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
X-Dual-Token: dual_eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...  // 仅中国联通模块需要
Content-Type: application/json

{
  "path": "/path/to/delete"
}
```

**参数说明**
- `path`: 待删除文件或目录的完整路径（必填）

**响应（成功）**
```json
{
  "code": 200,
  "message": "删除成功"
}
```

**响应（失败）**
```json
{
  "code": 404,
  "message": "路径不存在：stat /path/to/delete: no such file"
}
```

---

### 11. 批量删除

批量删除多个文件或目录，容错设计（部分失败不中断）。

**请求**
```http
POST /sftp/batchdelete
X-SFTP-Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
X-Dual-Token: dual_eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...  // 仅中国联通模块需要
Content-Type: application/json

[
  {
    "path": "/path/to/delete1",
    "isDir": false,
    "name": "file1.txt"
  },
  {
    "path": "/path/to/delete2",
    "isDir": true,
    "name": "folder2"
  }
]
```

**参数说明**
- 数组元素包含：
  - `path`: 完整路径（必填）
  - `isDir`: 是否为目录（可选）
  - `name`: 文件名或目录名（可选）

**响应（全部成功）**
```json
{
  "code": 200,
  "message": "批量删除成功"
}
```

**响应（部分成功，部分失败）**
```json
{
  "code": 206,
  "message": "部分删除成功:删除失败：权限拒绝，删除失败：路径不存在",
  "errors": [
    "删除失败：权限拒绝",
    "删除失败：路径不存在"
  ]
}
```

**响应（全部失败）**
```json
{
  "code": 500,
  "message": "批量删除失败：删除失败：权限拒绝，删除失败：路径不存在"
}
```

---

### 12. 重命名文件或目录

修改文件或目录的名称。

**请求**
```http
POST /sftp/rename
X-SFTP-Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
X-Dual-Token: dual_eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...  // 仅中国联通模块需要
Content-Type: application/json

{
  "oldPath": "/old/path/file.txt",
  "newName": "new_name.txt"
}
```

**参数说明**
- `oldPath`: 原文件/目录完整路径（必填）
- `newName`: 新名称（仅文件名，不含路径）

**响应（成功）**
```json
{
  "code": 200,
  "message": "重命名成功"
}
```

**响应（失败）**
```json
{
  "code": 400,
  "message": "名称包含非法字符：名称不能包含特殊字符"
}
```

---

### 13. 搜索文件或目录

在当前路径下递归搜索（不区分大小写模糊匹配）。

**请求**
```http
GET /sftp/search?sftp_token=XXX&path=/start/dir&keyword=filename
```

**参数说明**
- `path`: 搜索起始路径，默认为 `/`
- `keyword`: 搜索关键字（必填）

**响应**
```json
{
  "code": 200,
  "message": "搜索完成",
  "data": {
    "total": 5,
    "results": [
      {
        "name": "REPORT.PDF",
        "path": "/start/dir/subdir/REPORT.PDF",
        "isDir": false,
        "size": 204800,
        "modified": "2026-08-24T15:20:00+08:00",
        "parentPath": "/start/dir/subdir"
      },
      {
        "name": "report_backup.pdf",
        "path": "/start/dir/report_backup.pdf",
        "isDir": false,
        "size": 102400,
        "modified": "2026-08-23T10:10:00+08:00",
        "parentPath": "/start/dir"
      }
    ]
  }
}
```

**参数说明**
- `total`: 搜索结果总数
- `results`: 结果数组，每个结果包含：
  - `FileInfo`: 文件基本信息
  - `parentPath`: 父目录路径

---

## 数据格式

### 通用响应结构

```json
{
  "code": 200,          // 状态码
  "message": "string",  // 响应消息
  "data": {}            // 响应数据（可选）
}
```

### 文件信息结构

```json
{
  "name": "string",             // 文件名或目录名
  "path": "string",             // 完整路径
  "isDir": false,               // 是否为目录
  "size": 1024,                 // 文件大小（字节），目录为 0
  "modified": "2026-08-25T10:30:00+08:00"  // 最后修改时间（RFC3339 格式）
}
```

### 搜索结果结构

```json
{
  "name": "string",
  "path": "string",
  "isDir": false,
  "size": 1024,
  "modified": "2026-08-25T10:30:00+08:00",
  "parentPath": "/parent/directory"
}
```

---

## 错误码说明

### HTTP 状态码

| 代码 | 说明 |
|------|------|
| 200 | 成功 |
| 206 | 部分成功（批量操作时使用） |
| 400 | 请求参数错误 |
| 403 | 禁止访问（越界路径） |
| 404 | 资源不存在 |
| 429 | 请求过于频繁（双控验证达到并发限制） |
| 500 | 服务器内部错误 |

### 业务错误码

| 代码 | 说明 |
|------|------|
| 400 | 通用业务错误 |
| 403 | 无权限访问 |
| 428 | 需要双控验证（未提供有效 Dual-Token） |
| 50014 | SFTP 连接失效，请重新登录 |

### 常见错误消息

| 消息 | 原因 | 解决方案 |
|------|------|----------|
| `<密码错误> invalid user` | 用户名或密码错误 | 检查用户名和密码 |
| `<密钥错误> private key file load failed` | 密钥文件格式错误或密码不正确 | 检查私钥文件格式 |
| SFTP 连接失效 | Token 过期或连接已断开 | 重新登录获取新 Token |
| 您的角色无权登录该模块 | 用户角色不在白名单内 | 联系管理员添加角色权限 |
| 密码已过期，请先在平台修改密码 | 账号密码已过期 | 前往账号管理平台修改密码 |
| 该账号需先修改密码 | MustChangePassword 标记为 true | 使用 change_token 修改密码 |
| 双控验证失败：无效账号 | 双控复核人账号不存在 | 检查复核人账号是否正确 |
| 路径包含非法字符 | 文件名或目录名包含 `..` 等危险字符 | 检查输入内容是否合法 |
| 目标路径超出连接允许的根路径 | 试图访问模块根路径之外的目录 | 确认目标路径在允许的范围内 |

---

## 安全特性

### 1. 身份认证
- 密码登录：RSA 加密传输
- 密钥登录：私钥文件验证
- 模块化登录：支持普通、HOTLABEL、CHINAUNICOM 三种模式

### 2. 授权控制
- Token 机制：登录签发 Token，8 小时有效期
- 双控验证：中国联通模块强制双人复核
- 角色白名单：按模块配置角色权限
- 路径隔离：每个连接限制在指定根路径内

### 3. 数据安全
- 路径边界校验：防止越权访问
- 文件名合法性校验：防止路径穿越攻击
- 日志审计：记录登录、操作、删除等全链路日志

### 4. 操作审计
- 自动记录：登录、登出、上传、下载、删除等操作
- 双控追溯：记录复核人身份
- 来源标识：区分 normal/hotlabel/chinaunicom 来源

---

## 使用示例

### 流程 1：普通用户密码登录并上传文件

```javascript
// 1. 登录
const loginResp = await fetch('/sftp/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    username: 'user1',
    password: encryptedPassword,  // RSA 加密
    loginType: 'normal'
  })
});
const { data } = await loginResp.json();
const token = data.sftp_token;

// 2. 上传文件
const formData = new FormData();
formData.append('path', '/upload/dir');
formData.append('file', fileInput.files[0]);

const uploadResp = await fetch('/sftp/upload', {
  method: 'POST',
  headers: { 'X-SFTP-Token': token },
  body: formData
});
```

### 流程 2：中国联通模块登录 + 双控验证 + 批量上传

```javascript
// 1. 中国联通模块登录
const loginResp = await fetch('/sftp/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    username: 'domain_user',
    password: encryptedPassword,
    loginType: 'chinaunicom'
  })
});
const { data: loginData } = await loginResp.json();
const sftpToken = loginData.sftp_token;

// 2. 双控验证
const verifyResp = await fetch('/sftp/dualverify', {
  method: 'POST',
  headers: { 
    'Content-Type': 'application/json',
    'X-SFTP-Token': sftpToken
  },
  body: JSON.stringify({
    username: 'reviewer',
    password: encryptedReviewerPassword
  })
});
const { data: verifyData } = await verifyResp.json();
const dualToken = verifyData.dual_token;

// 3. 批量上传（带双控 Token）
for (const file of files) {
  const formData = new FormData();
  formData.append('path', '/chinaunicom/uploads');
  formData.append('file', file);

  await fetch('/sftp/upload', {
    method: 'POST',
    headers: { 
      'X-SFTP-Token': sftpToken,
      'X-Dual-Token': dualToken  // 复用同一凭证
    },
    body: formData
  });
}
```

---

## 备注

1. **Token 时效性**: SFTP Token 有效期 8 小时，超时需重新登录
2. **双控 Token**: 60 秒有效，可复用（同一批次上传共享凭证）
3. **流式传输**: 上传/下载采用流式处理，避免内存溢出
4. **隐藏文件**: 文件列表和搜索自动过滤以`.`开头的隐藏文件
5. **大小写敏感**: 文件名搜索不区分大小写
6. **编码规范**: 中文文件名使用 UTF-8 URL 编码兼容低版本浏览器

---

*文档生成时间：2026-08-25*
