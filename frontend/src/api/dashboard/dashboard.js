// ? 看板模块 API
import request from '@/utils/request.js'

// ! 获取看板卡片账号数
export const reqAccountCount = () => request({url:'/dashboard/user/total',method:'get'})

// ! 获取累计访问量（登录总次数）
export const reqTotalAccessCount = () => request({url:'/dashboard/access/total',method:'get'})

// ! 获取今日访问量
export const reqTodayAccessCount = () => request({url:'/dashboard/access/today',method:'get'})

// ! 获取访问量增长率（今日 vs 昨日）
export const reqAccessGrowth = () => request({url:'/dashboard/access/growth',method:'get'})

// ! 获取累计传输数（上传 + 下载总次数）
export const reqTotalTransferCount = () => request({url:'/dashboard/transfer/total',method:'get'})

// ! 获取今日传输量
export const reqTodayTransferCount = () => request({url:'/dashboard/transfer/today',method:'get'})

// ! 获取传输量增长率（今日 vs 昨日）
export const reqTransferGrowth = () => request({url:'/dashboard/transfer/growth',method:'get'})

// ! 获取认证方式分布
export const reqAuthDistribution = () => request({url:'/dashboard/auth/distribution',method:'get'})

// ! 获取活跃用户 Top6（按登录次数排行）
export const reqActiveUsersTop6 = () => request({url:'/dashboard/users/active-top6',method:'get'})

// ! 获取传输量排行 Top10（按上传+下载次数排序）
export const reqTopTransferUsers = () => request({url:'/dashboard/users/top-transfers',method:'get'})