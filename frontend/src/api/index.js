// !将四个模块请求的接口函数统一暴露
import * as trademark from './product/tradeMark'
import * as attr from './product/attr'
import * as spu from './product/spu'
import * as sku from './product/sku'
// !引入 SFTP 相关的接口文件
import * as sftpuser from './sftp/sftpuser'
// ! 引入通讯录相关接口文件
import * as contact from './contact/contact'
// ! 引入系统安全相关接口文件
import * as system from './system/system'
// ! 引入看板相关接口文件
import * as dashboard from './dashboard/dashboard'



export default {
  trademark,
  attr,
  spu,
  sku,
  sftpuser,
  contact,
  system,
  dashboard
}
