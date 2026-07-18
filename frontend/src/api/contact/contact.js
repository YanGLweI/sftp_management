//?通讯录管理模块API
import request from '@/utils/request.js'

// ! 分页获取通讯录列表  /contact/address/{page}/{limit} get
export const reqContactList = (page,limit,searchObj) => request({url:`/contact/address/${page}/${limit}`,method:'get',params:searchObj})

// ! 添加联系人  /contact/address post
export const reqAddContact = (data) => request({url:'/contact/address',method:'post',data})

// ! 更新一个联系人  /contact/address put
export const reqUpdateContact = (data) => request({url:'/contact/address',method:'put',data})

// ! 删除一个联系人  /contact/address/{id} delete
export const reqDeleteContact = (id) => request({url:`/contact/address/${id}`,method:'delete'})

// ! 批量删除联系人  /contact/address delete
export const reqDeleteContacts = (ids) => request({url:'/contact/address',method:'delete',data:ids})

// ! 获取联系人选项列表  /contact/options get
export const reqContactOptions = () => request({url:'/contact/options',method:'get'})