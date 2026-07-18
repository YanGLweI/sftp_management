// ? 看板模块API
import request from '@/utils/request.js'

// ! 获取看板卡片账号数
export const reqAccountCount = () => request({url:'/dashboard/user/total',method:'get'})