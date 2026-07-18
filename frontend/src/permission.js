import Vue from 'vue'
import router from './router'
import store from './store'
import { Message } from 'element-ui'
import NProgress from 'nprogress' // progress bar
import 'nprogress/nprogress.css' // progress bar style
import { getToken } from '@/utils/auth' // get token from cookie
import getPageTitle from '@/utils/get-page-title'

NProgress.configure({ showSpinner: false }) // NProgress Configuration

const whiteList = ['/login','/file'] // no redirect whitelist

router.beforeEach(async(to, from, next) => {
  // start progress bar
  NProgress.start()

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
        Vue.prototype.$idleTimeout.checkAndStart()
        next()
      } else {
        try {
          // get user info
          await store.dispatch('user/getInfo')
          // 重置定时器
          Vue.prototype.$idleTimeout.checkAndStart()
          next({...to})
        } catch (error) {
          // remove token and go to login page to re-login
          await store.dispatch('user/resetToken')
          Message.error(error || 'Has Error')
          Vue.prototype.$idleTimeout.stop()
          next(`/login?redirect=${to.path}`)
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
