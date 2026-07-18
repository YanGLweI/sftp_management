// ? 系统安全模块API
import request from '@/utils/request.js'

//? 卡巴斯基模块
// ! 获取卡巴斯基信息
export const reqKasperskyInfo = () => request({url:'/system/antivirus',method:'get'})

// !获取卡巴斯基计划
export const reqKasperskySchedule = () => request({url:'/system/antivirus/schedule',method:'get'})

// ! 设置卡巴斯基计划
export const reqSetKasperskySchedule = (data) => request({url:'/system/antivirus/schedule',method:'post',data})

// ! 立即扫描
export const reqKasperskyScan = () => request({url:'/system/antivirus/scan',method:'get'})

// ! 获取卡巴斯基隔离区
export const reqKasperskyIsolationZone = () => request({url:'/system/antivirus/isolationzone',method:'get'})

// ! 获取卡巴斯基报告计划
export const reqKasperskyReportSchedule = () => request({url:'/system/antivirus/report/schedule',method:'get'})

// ! 设置卡巴斯基报告计划
export const reqSetKasperskyReportSchedule = (data) => request({url:'/system/antivirus/report/schedule',method:'post',data})

// ? 更新模块
// ! 获取更新历史
export const reqUpdateHistory = (params) => request({url:'/system/update/history',method:'get',params})

// ! 获取更新详情 /system/update/detail/:id
export const reqUpdateDetail = (id) => request({url:`/system/update/detail/${id}`,method:'get'})

// ! 获取更新计划
export const reqUpdateSchedule = () => request({url:'/system/update/schedule',method:'get'})

// ! 设置更新计划
export const reqSetUpdateSchedule = (data) => request({url:'/system/update/schedule',method:'post',data})

// ! 获取更新报告任务计划
export const reqUpdateReportSchedule = () => request({url:'/system/update/report/schedule',method:'get'})

// ! 设置更新报告任务计划
export const reqSetUpdateReportSchedule = (data) => request({url:'/system/update/report/schedule',method:'post',data})

// ? 系统加固模块
// ! 获取系统加固列表
export const reqSystemHardeningList = (params) => request({url:'/system/security/checklist',method:'get',params})

// ! 立即执行系统加固
export const reqSystemHardening = () => request({url:'/system/security/start',method:'get'})

// ! 获取系统加固计划
export const reqSystemHardeningSchedule = () => request({url:'/system/security/schedule',method:'get'})

// ! 设置系统加固计划
export const reqSetSystemHardeningSchedule = (data) => request({url:'/system/security/schedule',method:'post',data})

// ! 获取系统加固报告计划
export const reqSystemHardeningReportSchedule = () => request({url:'/system/security/report/schedule',method:'get'})

// ! 设置系统加固报告计划
export const reqSetSystemHardeningReportSchedule = (data) => request({url:'/system/security/report/schedule',method:'post',data})