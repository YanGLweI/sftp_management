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
    redirect: '/sftp/sftpuser',
    children: [
      {
        path: 'sftpuser',
        name: 'SftpUser',
        component: () => import('@/views/sftp/SftpUser'),
        meta: { title: '账号管理' },
      },
      {
        path: 'sftplog',
        name: 'SftpLog',
        component: () => import('@/views/sftp/Sftplog'),
        meta: { title: '日志管理' },
      },
      {
        path: 'contacts',
        name: 'Contacts',
        component: () => import('@/views/sftp/Contacts'),
        meta: { title: '通讯邮箱' },
      },
    ],
  },
  // ! 系统安全
  {
    path: '/system',
    component: Layout,
    name: 'System',
    meta: { title: '系统安全', icon: 'el-icon-lock' },
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
]

// ! 异步路由:不同的用户(角色),需要过滤筛选出的路由
export const asyncRoutes = [
  // ! 权限管理
  /* {
    name: 'Acl',
    path: '/acl',
    component: Layout,
    redirect: '/acl/user/list',
    meta: {
      title: '权限管理',
      icon: 'el-icon-lock'
    },
    children: [
      {
        name: 'User',
        path: 'user/list',
        component: () => import('@/views/acl/user/list'),
        meta: {
          title: '用户管理',
        },
      },
      {
        name: 'Role',
        path: 'role/list',
        component: () => import('@/views/acl/role/list'),
        meta: {
          title: '角色管理',
        },
      },
      {
        name: 'RoleAuth',
        path: 'role/auth/:id',
        component: () => import('@/views/acl/role/roleAuth'),
        meta: {
          activeMenu: '/acl/role/list',
          title: '角色授权',
        },
        hidden: true,
      },
      {
        name: 'Permission',
        path: 'permission/list',
        component: () => import('@/views/acl/permission/list'),
        meta: {
          title: '菜单管理',
        },
      },
    ]
  }, */
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



