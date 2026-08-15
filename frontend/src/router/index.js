// 引入Vue|Vue-router
import Vue from 'vue'
import Router from 'vue-router'

//! 使用路由插件
Vue.use(Router)

// !映入 最外层骨架的一级路由组件
/* Layout */
import Layout from '@/layout'

/**
 * Note: sub-menu only appear when route children.length >= 1
 * Detail see: https://panjiachen.github.io/vue-element-admin-site/guide/essentials/router-and-nav.html
 *
 * hidden: true                   if set true, item will not show in the sidebar(default is false)
 * alwaysShow: true               if set true, will always show the root menu
 *                                if not set alwaysShow, when item has more than one children route,
 *                                it will becomes nested mode, otherwise not show the root menu
 * redirect: noRedirect           if set noRedirect will no redirect in the breadcrumb
 * name:'router-name'             the name is used by <keep-alive> (must set!!!)
 * meta : {
    roles: ['admin','editor']    control the page roles (you can set multiple roles)
    title: 'title'               the name show in sidebar and breadcrumb (recommend set)
    icon: 'svg-name'/'el-icon-x' the icon show in the sidebar
    breadcrumb: false            if set false, the item will hidden in breadcrumb(default is true)
    activeMenu: '/example/list'  if set path, the sidebar will highlight the path you set
  }
 */

/**
 * constantRoutes
 * a base page that does not have permission requirements
 * all roles can be accessed
 */
// ! 路由配置
// ! 需要把项目中的路由拆分

// ! 常量路由:所有用户都能看见的路由
// ! 所有角色:登录\404\首页
export const constantRoutes = [
  // ! 登录页
  {
    path: '/login',
    component: () => import('@/views/login/index'),
    hidden: true
  },
  // ! 文件管理
  {
    path: '/file',
    component: () => import('@/views/file/index'),
    hidden: true
  },
  // ! 404
  {
    path: '/404',
    component: () => import('@/views/404'),
    hidden: true
  },
  // ! 首页
  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [{
      path: 'dashboard',
      name: 'Dashboard',
      component: () => import('@/views/dashboard/index'),
      meta: { title: '首页', icon: 'dashboard' }
    },
    ]
  },
  // // ! 传输管理
  // {
  //   path: '/sftp',
  //   component: Layout,
  //   children: [{
  //     path: 'index',
  //     name: 'Sftp',
  //     component: () => import('@/views/sftp'),
  //     meta: { title: '传输管理', icon: 'el-icon-sort' }
  //   }],
  // },
  // ! 传输管理
  {
    path: '/sftp',
    component: Layout,
    name: 'Sftp',
    meta: { title: '传输管理', icon: 'el-icon-sort' },
    alwaysShow: true,
    redirect: '/sftp/sftpuser',
    children: [
      {
        path: 'sftpuser',
        name: 'SftpUser',
        component: () => import('@/views/sftp/SftpUser'),
        meta: { title: '账号管理' },
      },
      {
        path: 'contacts',
        name: 'Contacts',
        component: () => import('@/views/sftp/Contacts'),
        meta: { title: '通讯邮箱' },
      },
    ],
  },
  // ! 日志管理（1级菜单：平台日志 + SFTP日志）
  {
    path: '/log',
    component: Layout,
    name: 'Log',
    meta: { title: '日志管理', icon: 'el-icon-tickets' },
    alwaysShow: true,
    redirect: '/log/platformlog',
    children: [
      {
        path: 'platformlog',
        name: 'PlatformLog',
        component: () => import('@/views/sftp/Sftplog'),
        meta: { title: '平台日志' },
      },
      {
        path: 'sftplog',
        name: 'SftpLog',
        component: () => import('@/views/log/SftpLog'),
        meta: { title: 'SFTP日志' },
      },
    ],
  },
  // ! 系统安全
  {
    path: '/system',
    component: Layout,
    name: 'System',
    meta: { title: '系统安全', icon: 'el-icon-lock' },
    alwaysShow: true,
    redirect: '/system/update',
    children: [
      {
        path: 'update',
        name: 'SystemUpdate',
        component: () => import('@/views/systemSecurity/SystemUpdate'),
        meta: { title: '系统更新' },
      },
      {
        path: 'antivirus',
        name: 'Antivirus',
        component: () => import('@/views/systemSecurity/Antivirus'),
        meta: { title: '病毒管理' },
      },
      {
        path: 'systemHardening',
        name: 'SystemHardening',
        component: () => import('@/views/systemSecurity/SystemHardening'),
        meta: { title: '系统加固' },
      },
    ],
  },
  // ! 平台设置
  {
    path: '/settings',
    component: Layout,
    name: 'Settings',
    meta: { title: '平台设置', icon: 'el-icon-setting' },
    alwaysShow: true,
    redirect: '/settings/roles',
    children: [
      {
        path: 'roles',
        name: 'RoleManagement',
        component: () => import('@/views/settings/Role'),
        meta: { title: '角色管理' },
      },
      {
        path: 'localusers',
        name: 'LocalUserManagement',
        component: () => import('@/views/settings/LocalUser'),
        meta: { title: '本地账号' },
      },
      {
        path: 'password-policy',
        name: 'PasswordPolicy',
        component: () => import('@/views/settings/PasswordPolicy'),
        meta: { title: '密码策略' },
      },
      {
        path: 'ldap',
        name: 'LDAPManagement',
        component: () => import('@/views/settings/LDAPManagement'),
        meta: { title: 'LDAP 管理' },
      },
    ],
  },
  // ! SFTP 管理（新增）
  {
    path: '/sftp-module',
    component: Layout,
    name: 'SftpModuleManagement',
    meta: { title: 'SFTP 管理', icon: 'el-icon-cpu' },
    alwaysShow: true,
    redirect: '/sftp-module/hotlabel-config',
    children: [
      {
        path: 'hotlabel-config',
        name: 'HotLabelConfig',
        component: () => import('@/views/admin/hotlabel-config/index.vue'),
        meta: { title: '标签上传配置' },
      },
      {
        path: 'chinaunicom-config',
        name: 'ChinaUnicomConfig',
        component: () => import('@/views/admin/chinaunicom-config/index.vue'),
        meta: { title: '中国联通配置' },
      },
    ],
  },
]

// ! 异步路由：不同的用户 (角色), 需要过滤筛选出的路由
export const asyncRoutes = [
  // ! 商品管理
  /* {
    path: '/product',
    component: Layout,
    name: 'Product',
    meta: { title: '商品管理', icon: 'el-icon-goods' },
    children: [
      {
        path: 'trademark',
        name: 'Trademark',
        component: () => import('@/views/product/tradeMark'),
        meta: { title: '品牌管理' }
      },
      {
        path: 'attr',
        name: 'Attr',
        component: () => import('@/views/product/Attr'),
        meta: { title: '平台属性管理' }
      },
      {
        path: 'spu',
        name: 'Spu',
        component: () => import('@/views/product/Spu'),
        meta: { title: 'Spu管理' }
      },
      {
        path: 'sku',
        name: 'Sku',
        component: () => import('@/views/product/Sku'),
        meta: { title: 'Sku管理' }
      },
    ]
  }, */
]

// ! 任意路由
// ! 路径出现错误的时候重定向到404
export const anyRoutes = { path: '*', redirect: '/404', hidden: true }


const createRouter = () => new Router({
  // mode: 'history', // require service support
  scrollBehavior: () => ({ y: 0 }),
  routes: constantRoutes
})

const router = createRouter()

// Detail see: https://github.com/vuejs/vue-router/issues/1234#issuecomment-357941465
export function resetRouter() {
  const newRouter = createRouter()
  router.matcher = newRouter.matcher // reset router
}


export default router



