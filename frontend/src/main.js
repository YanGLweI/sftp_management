import Vue from 'vue'

import 'normalize.css/normalize.css' // A modern alternative to CSS resets

import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
import locale from 'element-ui/lib/locale/lang/en' // lang i18n

import '@/styles/index.scss' // global css

import App from './App'
import store from './store'
import router from './router'

import '@/icons' // icon
import '@/permission' // permission control

/**
 * If you don't want to use mock-server
 * you want to use MockJs for mock api
 * you can execute: mockXHR()
 *
 * Currently MockJs will be used in the production environment,
 * please remove it before going online ! ! !
 */
// if (process.env.NODE_ENV === 'production') {
//   const { mockXHR } = require('../mock')
//   mockXHR()
// }

// set ElementUI lang to EN
// Vue.use(ElementUI, { locale })
// 如果想要中文版 element-ui，按如下方式声明
Vue.use(ElementUI)

// 引入相关API请求接口
import API from '@/api'
//! 组件实例的原型的原型指向Vue.prototype
//! 任意组件都能使用API相关接口
Vue.prototype.$API = API

import CategorySelect from '@/components/CategorySelect'
Vue.component(CategorySelect.name,CategorySelect)

import HintButton from '@/components/HintButton'
Vue.component(HintButton.name,HintButton)

Vue.config.productionTip = false

// 全局自定义指令v-focus,自动聚焦输入框
Vue.directive('focus', {
  // 当元素插入DOM时
  inserted: function (el, binding) {
    // 使用setTimeout确保DOM更新后执行
    setTimeout(() => {
      if (binding.value && el.offsetParent !== null) {
        const input = el.querySelector('input');
        if (input) {
          input.focus();
        }
      }
    }, 50);
  },
  // 当元素更新时执行
  update: function (el, binding) {
    // 只有当binding.value为true且之前为false时才执行
    if (binding.value && binding.oldValue !== binding.value) {
      // 添加防抖机制
      clearTimeout(el._focusTimer);
      el._focusTimer = setTimeout(() => {
        if (el.offsetParent !== null) {
          const input = el.querySelector('input');
          if (input) {
            input.focus();
          }
        }
      }, 100);
    }
  },
  // 当指令与元素解绑时清理定时器
  unbind: function (el) {
    clearTimeout(el._focusTimer);
  }
});

import idleTimeout from '@/utils/idleTimeout'

// 空闲超时挂载到Vue原型
Vue.prototype.$idleTimeout = idleTimeout



new Vue({
  el: '#app',
  router,
  store,
  render: h => h(App),
  mounted(){
    // 初始化超时监听
    this.$idleTimeout.init()
  }
})
