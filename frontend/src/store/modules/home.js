import mockRequest from '@/utils/mockRequest.js'
import { reqAccountCount } from '@/api/dashboard/dashboard.js'

const state = {
  list:{}
}
const mutations = {
  GETDATA(state,list){
    state.list = list
  }
}
const actions = {
  // 发请求获取首页的模拟数据
  async getData({commit}){
    // let result = await mockRequest.get('/home/list')
    let result = await reqAccountCount()
    // if(result.code == 20000){
    if(result.code == 200){
      commit('GETDATA',result.data)
    }
  }
}
const getters ={}
export default {
  state,
  mutations,
  actions,
  getters
}