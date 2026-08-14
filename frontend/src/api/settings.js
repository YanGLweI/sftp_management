import request from '@/utils/request'

// 获取所有菜单树
export function getMenus() {
  return request({
    url: '/settings/menus',
    method: 'get'
  })
}

// 获取角色列表
export function getRoleList(params) {
  return request({
    url: '/settings/roles',
    method: 'get',
    params
  })
}

// 创建角色
export function createRole(data) {
  return request({
    url: '/settings/role',
    method: 'post',
    data
  })
}

// 更新角色
export function updateRole(id, data) {
  return request({
    url: `/settings/role/${id}`,
    method: 'put',
    data
  })
}

// 删除角色
export function deleteRole(id) {
  return request({
    url: `/settings/role/${id}`,
    method: 'delete'
  })
}

// 获取角色详情
export function getRoleDetail(id) {
  return request({
    url: `/settings/role/${id}`,
    method: 'get'
  })
}

// 获取角色下拉列表
export function getRoleSelect() {
  return request({
    url: '/settings/role/select',
    method: 'get'
  })
}

// 获取本地账号列表
export function getLocalUserList(params) {
  return request({
    url: '/settings/localusers',
    method: 'get',
    params
  })
}

// 创建本地账号
export function createLocalUser(data) {
  return request({
    url: '/settings/localuser',
    method: 'post',
    data
  })
}

// 更新本地账号
export function updateLocalUser(id, data) {
  return request({
    url: `/settings/localuser/${id}`,
    method: 'put',
    data
  })
}

// 重置密码
export function resetLocalUserPassword(id, data) {
  return request({
    url: `/settings/localuser/${id}/reset-password`,
    method: 'put',
    data
  })
}

// 删除本地账号
export function deleteLocalUser(id) {
  return request({
    url: `/settings/localuser/${id}`,
    method: 'delete'
  })
}

// 获取密码策略
export function getPasswordPolicy() {
  return request({
    url: '/settings/password-policy',
    method: 'get'
  })
}

// 更新密码策略
export function updatePasswordPolicy(data) {
  return request({
    url: '/settings/password-policy',
    method: 'put',
    data
  })
}

// 验证密码
export function validatePassword(data) {
  return request({
    url: '/settings/password-policy/validate',
    method: 'post',
    data
  })
}