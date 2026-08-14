import Vue from 'vue'
import router from './router'
import store from './store'
import { Message } from 'element-ui'
import NProgress from 'nprogress' // progress bar
import 'nprogress/nprogress.css' // progress bar style
import { getToken } from '@/utils/auth' // get token from cookie
import getPageTitle from '@/utils/get-page-title'
import { getFirstAccessibleRoutePath } from '@/utils/get-first-accessible-route'

NProgress.configure({ showSpinner: false }) // NProgress Configuration

const whiteList = ['/login', '/file', '/404'] // no redirect whitelist

// /file 为对外独立页面：自身及其子路径均不重定向到主登录页
const isFileRoute = path => path === '/file' || path.startsWith('/file/')

// 校验当前目标路由是否有权限访问
// routes 为空时放行（兼容旧Token/异常态，避免锁死）
const hasRoutePermission = (to, routes) => {
  if (!routes || routes.length === 0) return true
  if (whiteList.indexOf(to.path) !== -1) return true
  // 取 matched 中最后一个非 hidden 的路由名作为校验目标
  const matched = to.matched || []
  for (let i = matched.length - 1; i >= 0; i--) {
    const record = matched[i]
    if (record.meta && record.meta.title && !record.hidden) {
      return routes.indexOf(record.name) !== -1
    }
  }
  // 没有匹配到具体路由则放行
  return true
}

// 默认首页（/ → /dashboard）无权限时：静默跳转到第一个可访问路由，
// 避免先进入 /dashboard 触发“无权访问”警告再跳走，提升登录体验
const redirectDefaultHomeIfNoPermission = (to, next) => {
  if (to.path !== '/') return false
  const routes = store.state.user.routes || []
  if (routes.length > 0 && routes.indexOf('Dashboard') === -1) {
    const accessiblePath = getFirstAccessibleRoutePath(store.state.user.resultAllRoutes)
    if (accessiblePath && accessiblePath !== '/login' && accessiblePath !== '/404') {
      next(accessiblePath)
      NProgress.done()
      return true
    }
  }
  return false
}

router.beforeEach(async(to, from, next) => {
  // start progress bar
  NProgress.start()

  // 不存在的 /file 子路由：直接回 /file 页面，不进入登录重定向逻辑
  if (to.path.startsWith('/file/')) {
    next('/file')
    NProgress.done()
    return
  }

  // set page title
  document.title = getPageTitle(to.meta.title)

  // 更新目标路由路径到超时管理工具
  Vue.prototype.$idleTimeout.updateTargetPath(to.path)

  // determine whether the user has logged in
  const hasToken = getToken()

  if (hasToken) {
    if (to.path === '/login') {
      // if is logged in, redirect to the home page
      next({ path: '/' })
      NProgress.done()
    } else {
      const hasGetUserInfo = store.getters.name
      if (hasGetUserInfo) {
        // 默认首页无权限：静默跳转到第一个可访问路由（覆盖手动输入 / 的场景）
        if (redirectDefaultHomeIfNoPermission(to, next)) return
        // 已获取用户信息：校验路由权限
        if (!hasRoutePermission(to, store.state.user.routes)) {
          Message.warning('无权访问该页面')
          const accessiblePath = getFirstAccessibleRoutePath(store.state.user.resultAllRoutes)
          next(accessiblePath === '/login' ? '/404' : accessiblePath)
          NProgress.done()
          return
        }
        Vue.prototype.$idleTimeout.checkAndStart()
        next()
      } else {
        try {
          // get user info
          await store.dispatch('user/getInfo')
          // 登录后默认首页无权限：直接跳转到第一个可访问路由，避免中间产生警告
          if (redirectDefaultHomeIfNoPermission(to, next)) return
          // 获取用户权限后校验路由权限
          if (!hasRoutePermission(to, store.state.user.routes)) {
            // 登录页重定向链（刚登录）无权限时静默跳转，非主动访问不弹警告
            if (from.path !== '/login') {
              Message.warning('无权访问该页面')
            }
            const accessiblePath = getFirstAccessibleRoutePath(store.state.user.resultAllRoutes)
            next(accessiblePath === '/login' ? '/404' : accessiblePath)
            NProgress.done()
            return
          }
          // 重置定时器
          Vue.prototype.$idleTimeout.checkAndStart()
          next({...to})
        } catch (error) {
          // remove token and go to login page to re-login
          await store.dispatch('user/resetToken')
          Message.error(error || 'Has Error')
          Vue.prototype.$idleTimeout.stop()
          next(isFileRoute(to.path) ? '/file' : `/login?redirect=${to.path}`)
          NProgress.done()
        }
      }
    }
  } else {
    /* has no token*/
    //  未登录时，清除定时器
    Vue.prototype.$idleTimeout.stop()

    if (whiteList.indexOf(to.path) !== -1) {
      // in the free login whitelist, go directly
      next()
    } else {
      // other pages that do not have permission to access are redirected to the login page.
      next(`/login?redirect=${to.path}`)
      NProgress.done()
    }
  }
})

router.afterEach(() => {
  // finish progress bar
  NProgress.done()
  // 路由跳转完成后清除目标路径缓存
  Vue.prototype.$idleTimeout.updateTargetPath(null)
  Vue.prototype.$idleTimeout.checkAndStart()
})