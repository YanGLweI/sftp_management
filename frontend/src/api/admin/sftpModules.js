/* 
SFTP 模块配置相关的 API 请求函数
*/
import request from '@/utils/request'

const api_name = '/settings/sftp-modules'

export default {
  /* 
  获取所有模块的公共配置（无需登录，/file SFTP 登录页渲染表单用）
  */
  getPublicConfigs() {
    return request({
      url: '/sftp/module-configs',
      method: 'get'
    })
  },

  /* 
  获取所有模块配置
  */
  getAllConfigs() {
    return request({
      url: `${api_name}/all`,
      method: 'get'
    })
  },

  /* 
  获取指定模块的配置
  */
  getModuleConfig(name) {
    return request({
      url: `${api_name}/${name}`,
      method: 'get'
    })
  },

  /* 
  更新模块配置
  */
  updateModuleConfig(name, data) {
    return request({
      url: `${api_name}/${name}`,
      method: 'put',
      data
    })
  }
}
