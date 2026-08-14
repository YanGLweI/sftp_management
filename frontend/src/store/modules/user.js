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
  SET_RESULTASYNCROUTES:(state, payload)=>{
    // 从 payload 中解构 { asyncRoutes, filteredConstantRoutes }
    const { asyncRoutes: computedAsync, filteredConstantRoutes: filteredConstant } = payload
    // vuex存储当前用户的异步路由
    state.resultAsyncRoutes = computedAsync
    // 计算出当前用户需要的所有路由:过滤后的常量路由 + 过滤后的异步路由 + 任意路由
    state.resultAllRoutes = filteredConstant.concat(computedAsync, anyRoutes)
    // 用户的全部路由:异步路由,加上任意路由
    state.resultUserRoutes = computedAsync.concat(anyRoutes)
    // 给路由器添加新的路由
    router.addRoutes(state.resultUserRoutes)
  }
}

// ! 定义一个函数:2个数组进行对比,对比出当前用户到底显示哪些路由
const computedAsyncRoutes = (asyncRoutes,routes)=>{
  // 过滤出当前用户需要展示的异步路由
  return asyncRoutes.reduce((result,item)=>{
    if(routes.indexOf(item.name) != -1){
      // 浅拷贝，避免修改原始路由对象的 children
      const copy = { ...item }
      if(copy.children&&copy.children.length){
        copy.children = computedAsyncRoutes(copy.children,routes)
      }
      result.push(copy)
    }
    return result
  },[])
}

// 过滤常量路由：子菜单有权限时保留父级菜单
const filterConstantRoutes = (constantRoutes, routes) => {
  return constantRoutes.reduce((result, item) => {
    // 隐藏路由（/login /404 /file）始终显示
    if (item.hidden) {
      result.push(item)
      return result
    }
    // 无权限列表（异常/兼容）时全部显示
    if (!routes || routes.length === 0) {
      result.push(item)
      return result
    }
    // 浅拷贝，避免修改原始路由对象的 children
    const copy = { ...item }
    // 父级有权限：保留，子菜单递归过滤
    if (routes.indexOf(copy.name) !== -1) {
      if (copy.children && copy.children.length) {
        copy.children = filterConstantRoutes(copy.children, routes)
      }
      result.push(copy)
      return result
    }
    // 父级无权限：检查子菜单，只要有子菜单被授权则保留父级并过滤子菜单
    if (copy.children && copy.children.length) {
      const filteredChildren = filterConstantRoutes(copy.children, routes)
      if (filteredChildren.length > 0) {
        copy.children = filteredChildren
        result.push(copy)
        return result
      }
    }
    return result
  }, [])
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
      // 返回完整响应供登录页判断 must_change_password / password_expired
      return result
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
        // 同时过滤常量路由和异步路由
        const filteredConstant = filterConstantRoutes(constantRoutes, data.routes || [])
        const filteredAsync = computedAsyncRoutes(asyncRoutes, data.routes || [])
        commit('SET_RESULTASYNCROUTES', {
          asyncRoutes: filteredAsync,
          filteredConstantRoutes: filteredConstant
        })
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

