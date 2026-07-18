//?SFTP管理模块API
import request from '@/utils/request.js'

// ! 获取sftp用户列表  /user/account/{page}/{limit} get
export const reqUserList = (page,limit,searchObj) => request({url:`/user/account/${page}/${limit}`,method:'get',params:searchObj})

// ! 添加用户的接口  /user/account  post
export const reqAddUser = (data) => request({url:'/user/account',method:'post',data})

// ! 修改用户密码的接口 /user/account/ put
export const reqUpdateUser = (data) => request({url:'/user/account',method:'put',data})

// ! 删除一个用户 /user/account/{id}  delete
export const reqDeleteUser = (id) => request({url:`/user/account/${id}`,method:'delete'})

// ! 批量删除一个用户 /user/account delete
export const reqDeleteUsers = (ids) => request({url:'/user/account',method:'delete',data:ids})

// ! 发送邮件 /user/email post
export const reqSendEmail = (data) => request({url:'/user/email',method:'post',data})

// ! 获取日志列表  /user/log/{page}/{limit} get
export const reqLogList = (page,limit,searchObj) => request({url:`/user/log/${page}/${limit}`,method:'get',params:searchObj})

// ! 获取传输日志 /user/log/sftplog/{date} get
export const reqSftpLog = (date) => request({url:`/user/log/sftplog/${date}`,method:'get'})

// ? SFTP传输模块
// ! 登录sftp /sftp/login post
export const reqSftpLogin = (data) => request({url:'/sftp/login',method:'post',data})

// ! 读取sftp文件和目录 /sftp/files get
export const reqSftpFiles = (path) => request({url:'/sftp/files',method:'get',params:path})

// ! 登出sftp /sftp/logout get
export const reqSftpLogout = () => request({url:'/sftp/logout',method:'get'})

// ! 创建目录 /sftp/mkdir post
export const reqSftpMkdir = (data) => request({url:'/sftp/mkdir',method:'post',data})

// ! 下载文件 /sftp/download get（未使用reqSftpDownload）
export const reqSftpDownload = (path) => request({url:'/sftp/download',method:'get',params:path,responseType: 'blob'})

// ! 删除文件或目录 /sftp/delete post
export const reqSftpDelete = (data) => request({url:'/sftp/delete',method:'post',data})

// ! 重命名文件或目录 /sftp/rename post
export const reqSftpRename = (data) => request({url:'/sftp/rename',method:'post',data})

// ! 下载目录  /sftp/downloaddir get（未使用reqSftpDownloadDir）
export const reqSftpDownloadDir = (path) => request({url:'/sftp/downloaddir',method:'get',params:path,responseType: 'blob'})

// ! 批量删除文件或目录 /sftp/batchdelete post
export const reqSftpBatchDelete = (data) => request({url:'/sftp/batchdelete',method:'post',data})