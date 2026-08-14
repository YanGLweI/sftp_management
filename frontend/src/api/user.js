//! 引入axios(axios进行了二次封装)
import request from '@/utils/request'
// import mockRequest from '@/utils/mockRequest.js'

//! 对外暴露登录接口函数
export function login(data) {
  return request({
    // url: '/admin/acl/index/login',
    // url: '/vue-admin-template/user/login',
    url: '/login',
    method: 'post',
    data
  })
}

//! 对外暴露获取用户信息函数
export function getInfo() {
  return request({
    // url: '/admin/acl/index/info',
    // url: '/vue-admin-template/user/info',
    url: '/user/info',
    method: 'get',
    // params: { token }
  })
}

// !对外暴露退出登录的函数
export function logout() {
  return request({
    // url: '/admin/acl/index/logout',
    // url: '/vue-admin-template/user/logout',
    url: '/user/logout',
    method: 'get'
  })
}

// 修改密码
export function changePassword(data) {
  return request({
    url: '/user/change-password',
    method: 'post',
    data
  })
}
