/**
 * 获取用户有权限的第一个可见路由路径
 *
 * 从用户过滤后的路由列表中，找到第一个非 hidden、非登录、非404的可见路由路径。
 * 用于：登录成功后的重定向、路由权限校验失败时的兜底跳转、404页面"返回首页"按钮等场景。
 *
 * @param {Array} routes - 用户过滤后的路由列表（store.state.user.resultAllRoutes）
 * @returns {string} 第一个可用路由的 path，兜底返回 '/login'
 */
export function getFirstAccessibleRoutePath(routes) {
  if (!routes || routes.length === 0) return '/login'
  for (const route of routes) {
    if (route.hidden) continue
    // 过滤 404 和 login
    if (route.path === '*' || route.path === '/404' || route.path === '/login') continue
    // 有子路由：找第一个可见的子路由
    if (route.children && route.children.length > 0) {
      for (const child of route.children) {
        if (child.hidden) continue
        if (child.path === '*' || child.path === '/404' || child.path === '/login') continue
        const base = route.path.replace(/\/+$/, '')
        const childPath = child.path.startsWith('/') ? child.path : '/' + child.path
        return base + childPath
      }
    }
    return route.path
  }
  return '/login'
}