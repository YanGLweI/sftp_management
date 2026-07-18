// ! 映入登录\退出登录\获取用户信息的接口函数
import { login, logout, getInfo } from '@/api/user'
// ! 获取token\设置token\删除token的接口函数
import { getToken, setToken, removeToken } from '@/utils/auth'
// ! 路由模块中重置路由的方法
import { asyncRoutes, resetRouter,constantRoutes,anyRoutes } from '@/router'
import router from '@/router'

// ! 箭头函数
const getDefaultState = () => {
  return {
    // 获取token
    token: getToken(),
    // 存储用户名
    name: '',
    // 存储用户id
    id: '',
    // 存储用户头像
    avatar: '',
    // ! 服务器返回的菜单标记,不同角色返回的不一样
    routes:[],
    // ! 角色标记
    roles:[],
    // ! 按钮权限标记
    buttons:[],
    // ! 对比之后:异步路由和服务器返回的路由标记对比
    resultAsyncRoutes:[],
    // ! 用户需要的全部路由:异步路由,加上任意路由
    resultUserRoutes:[],
    // ! 用户全部的路由
    resultAllRoutes:[]
  }
}

const state = getDefaultState()

// ! 唯一修改State的地方
const mutations = {
  // ! 重置state
  RESET_STATE: (state) => {
    Object.assign(state, getDefaultState())
  },
  // ! 存储token
  SET_TOKEN: (state, token) => {
    state.token = token
  },
  // // ! 存储用户名
  // SET_NAME: (state, name) => {
  //   state.name = name
  // },
  // // ! 存储头像
  // SET_AVATAR: (state, avatar) => {
  //   state.avatar = avatar
  // }
  // ! 存储用户信息
  SET_USERINFO:(state,userInfo)=>{
    // 用户名
    state.name = userInfo.name
    // 用户id
    state.id = userInfo.id
    // 用户头像
    state.avatar = userInfo.avatar
    // 菜单权限标记
    state.routes = userInfo.routes
    // 按钮权限标记
    state.buttons = userInfo.buttons
    // 角色信息
    state.roles = userInfo.roles
  },
  // ! 最终计算出来的异步路由
  SET_RESULTASYNCROUTES:(state,asyncRoutes)=>{
    // vuex存储当前用户的异步路由,注意:一个用户需要展示完整的路由:常量\异步\任意路由
    state.resultAsyncRoutes = asyncRoutes
    // 计算出当前用户需要的所有路由
    state.resultAllRoutes = constantRoutes.concat(state.resultAsyncRoutes,anyRoutes)
    // 用户的全部路由:异步路由,加上任意路由
    state.resultUserRoutes = state.resultAsyncRoutes.concat(anyRoutes)
    // ! 给路由器添加新的路由:不包含常量路由,否则会报路由重复
    router.addRoutes(state.resultUserRoutes)
  }
}

// ! 定义一个函数:2个数组进行对比,对比出当前用户到底显示哪些路由
const computedAsyncRoutes = (asyncRoutes,routes)=>{
  // 过滤出当前用户需要展示的异步路由
  return asyncRoutes.filter(item=>{
    // 数组中没有这个元素返回-1,如果有返回一定不是-1
    if(routes.indexOf(item.name) != -1){
      // ! 递归,可能有2级3级...路由
      if(item.children&&item.children.length){
        item.children = computedAsyncRoutes(item.children,routes)
      }
      return true
    }
  })
}

// ! actions
const actions = {
  //! 处理登录业务
  async login({ commit }, userInfo) {
    // 解构用户名和密码
    const { username, password,loginType } = userInfo
    let result = await login({ name: username.trim(), password: password, loginType })
    if (result.code === 20000){
      // ! vuex存储token
      commit('SET_TOKEN', result.data.token)
      // ! 本地持久化存储
      setToken(result.data.token)
      return 'ok'
    }else{
      return Promise.reject(new Error('failed'))
    }
    // return new Promise((resolve, reject) => {
    //   login({ username: username.trim(), password: password }).then(response => {
    //     const { data } = response
    //     commit('SET_TOKEN', data.token)
    //     setToken(data.token)
    //     resolve()
    //   }).catch(error => {
    //     reject(error)
    //   })
    // })
  },

  //! get user info 获取用户信息
  getInfo({ commit, state }) {
    return new Promise((resolve, reject) => {
      getInfo().then(response => {
        // 获取用户信息
        const { data } = response
        
        if (!data) {
          return reject('Verification failed, please Login again.')
        }
        // ! 获取的用户信息包含 用户名name\用户头像avatar\routes[不同的用户应该展示哪些菜单的标记]\角色信息roles\buttons[按钮权限标记]
        // ! vuex存储全部信息
        commit('SET_USERINFO',data)
        commit('SET_RESULTASYNCROUTES',computedAsyncRoutes(asyncRoutes,data.routes))
        resolve(data)
      }).catch(error => {
        reject(error)
      })
    })
  },

  // user logout
  logout({ commit, state }) {
    return new Promise((resolve, reject) => {
      logout(state.token).then(() => {
        removeToken() // must remove  token  first
        resetRouter()
        commit('RESET_STATE')
        resolve()
      }).catch(error => {
        reject(error)
      })
    })
  },

  // remove token
  resetToken({ commit }) {
    return new Promise(resolve => {
      removeToken() // must remove  token  first
      commit('RESET_STATE')
      resolve()
    })
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}

