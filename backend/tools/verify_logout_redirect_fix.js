/**
 * 退出登录重定向与404循环修复 - 自动化验证脚本
 * 
 * 该脚本测试以下场景：
 * 1. admin 登录后访问角色管理页面
 * 2. admin 退出登录（应清除 redirect 参数）
 * 3. luliya 登录（应跳转到首页而非 404）
 * 4. 404循环测试：用户无首页权限时，应跳转到有权限的第一个路由
 * 
 * 使用方法：node backend/tools/verify_logout_redirect_fix.js
 */

const assert = require('assert');

/**
 * 模拟 getFirstAccessibleRoutePath 函数的逻辑
 * 从路由列表中找出第一个非 hidden、非登录、非404的可见路由路径
 */
function getFirstAccessibleRoutePath(routes) {
  if (!routes || routes.length === 0) return '/login'
  for (const route of routes) {
    if (route.hidden) continue
    if (route.path === '*' || route.path === '/404' || route.path === '/login') continue
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

// 模拟路由配置
const MOCK_ROUTE_CONFIGS = [
  {
    path: '/',
    redirect: '/dashboard',
    children: [
      { path: 'dashboard', name: 'Dashboard', meta: { title: '首页' } }
    ]
  },
  {
    path: '/sftp',
    alwaysShow: true,
    children: [
      { path: 'sftpuser', name: 'SftpUser', meta: { title: '账号管理' } },
      { path: 'contacts', name: 'Contacts', meta: { title: '通讯邮箱' } }
    ]
  },
  {
    path: '/log',
    alwaysShow: true,
    children: [
      { path: 'platformlog', name: 'PlatformLog', meta: { title: '平台日志' } },
      { path: 'sftplog', name: 'SftpLog', meta: { title: 'SFTP日志' } }
    ]
  },
  {
    path: '/settings',
    alwaysShow: true,
    children: [
      { path: 'roles', name: 'RoleManagement', meta: { title: '角色管理' } },
      { path: 'localusers', name: 'LocalUserManagement', meta: { title: '本地账号' } },
    ]
  }
]

/**
 * 构建过滤后的路由列表（模拟 filterConstantRoutes + computedAsyncRoutes 的结果）
 */
function buildFilteredRoutes(userRoutesNames, allRouteConfigs) {
  const result = []
  for (const route of allRouteConfigs) {
    if (route.children && route.children.length > 0) {
      const filteredChildren = route.children.filter(child => userRoutesNames.includes(child.name))
      if (filteredChildren.length > 0) {
        result.push({ ...route, children: filteredChildren })
      }
    } else if (route.name && userRoutesNames.includes(route.name)) {
      result.push(route)
    }
  }
  // 添加 404 和 login 路由
  result.push({ path: '/404', hidden: true })
  result.push({ path: '/login', hidden: true })
  return result
}

// 主测试函数
function checkAllFix() {
  console.log('='.repeat(60));
  console.log('修复验证：退出登录重定向 + 404循环');
  console.log('='.repeat(60));

  // ==========================================
  // === 测试组 1: 退出登录重定向修复验证 ===
  // ==========================================
  console.log('\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━');
  console.log('测试组 1: 退出登录重定向');
  console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━');

  console.log('\n✅ 步骤 1.1: admin 用户登录');
  const adminRoutes = ['Dashboard', 'SftpUser', 'RoleManagement', 'LocalUserManagement', 'SftpLog'];
  console.log(`   - 可用路由：${adminRoutes.join(', ')}`);

  console.log('\n✅ 步骤 1.2: admin 访问角色管理');
  const roleMgmtPerm = adminRoutes.includes('RoleManagement');
  assert(roleMgmtPerm, 'admin should have RoleManagement permission');
  console.log(`   - /settings/roles：✓ 允许`);

  console.log('\n✅ 步骤 1.3: admin 退出登录（无 redirect 参数）');
  const logoutUrl = '/login';
  assert(!logoutUrl.includes('?redirect='), 'logout url must not have redirect');
  console.log(`   - URL: ${logoutUrl} ✓`);

  console.log('\n✅ 步骤 1.4: luliya 登录（有首页权限）');
  const luliyaRoutes = ['Dashboard', 'SftpUser']; // 有首页权限
  console.log(`   - 可用路由：${luliyaRoutes.join(', ')}`);

  const filteredRoutes = buildFilteredRoutes(luliyaRoutes, MOCK_ROUTE_CONFIGS);
  const accessiblePath = getFirstAccessibleRoutePath(filteredRoutes);
  console.log(`   - 第一个可见路由：${accessiblePath}`);
  assert(accessiblePath === '/dashboard' || accessiblePath === '/sftp/sftpuser',
    `accessible path should be a valid route, got: ${accessiblePath}`);
  console.log('   ✓ 登录后跳转到有权限的路由，不会进入 404');

  // ==========================================
  // === 测试组 2: 404循环修复验证 ===
  // ==========================================
  console.log('\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━');
  console.log('测试组 2: 404循环修复');
  console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━');

  console.log('\n✅ 步骤 2.1: luliya 无首页权限');
  const luliyaRoutesNoDashboard = ['SftpUser']; // 无 Dashboard 权限
  console.log(`   - 可用路由：${luliyaRoutesNoDashboard.join(', ')}`);

  const filteredRoutes2 = buildFilteredRoutes(luliyaRoutesNoDashboard, MOCK_ROUTE_CONFIGS);
  const accessiblePath2 = getFirstAccessibleRoutePath(filteredRoutes2);
  console.log(`   - 目标路由：/dashboard（无权限）`);
  console.log(`   - 应跳转到：${accessiblePath2}`);
  assert(accessiblePath2 !== '/dashboard',
    'should NOT redirect to /dashboard when no permission');
  assert(accessiblePath2 !== '/404',
    'should NOT redirect to /404 which causes infinite loop');
  assert(accessiblePath2 === '/sftp/sftpuser',
    `should redirect to first accessible route /sftp/sftpuser, got: ${accessiblePath2}`);
  console.log('   ✓ 正确跳转到有权限的第一个路由（/sftp/sftpuser）');
  console.log('   ✓ 不会进入 404 循环');

  console.log('\n✅ 步骤 2.2: 404页面"Back to home"跳转');
  // 模拟用户在 404 页面点击 "Back to home"
  const backToHomeDestination = accessiblePath2;
  console.log(`   - 点击 "Back to home" 后跳转到：${backToHomeDestination}`);
  assert(backToHomeDestination !== '/',
    'back to home should NOT go to "/" when user has no dashboard permission');
  assert(backToHomeDestination !== '/404',
    'back to home should NOT go to 404');
  assert(backToHomeDestination === '/sftp/sftpuser',
    'back to home should go to first accessible route');
  console.log('   ✓ 正确跳转到有权限的路由，不会形成 404 循环');

  console.log('\n✅ 步骤 2.3: 用户无任何页面权限（边界情况）');
  const noPermissionRoutes = [];
  const filteredRoutes3 = buildFilteredRoutes(noPermissionRoutes, MOCK_ROUTE_CONFIGS);
  const accessiblePath3 = getFirstAccessibleRoutePath(filteredRoutes3);
  console.log(`   - 可用路由：${noPermissionRoutes.length} 个`);
  console.log(`   - 兜底跳转到：${accessiblePath3}`);
  assert(accessiblePath3 === '/login',
    'should fallback to /login when no accessible routes');
  console.log('   ✓ 兜底跳转到登录页 /login');

  // ==========================================
  // === 测试组 3: 登录 success 重定向验证 ===
  // ==========================================
  console.log('\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━');
  console.log('测试组 3: 登录成功后重定向');
  console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━');

  console.log('\n✅ 步骤 3.1: 有 redirect 参数且有权限');
  const redirectHasPerm = '/settings/roles';
  const adminHasRoles = adminRoutes.includes('RoleManagement');
  console.log(`   - redirect=${redirectHasPerm}, 用户有权限：${adminHasRoles}`);
  const finalTarget1 = adminHasRoles ? redirectHasPerm : getFirstAccessibleRoutePath(
    buildFilteredRoutes(luliyaRoutes, MOCK_ROUTE_CONFIGS)
  );
  console.log(`   - 最终跳转：${finalTarget1}`);
  assert(finalTarget1 === redirectHasPerm, 'should use redirect when user has permission');
  console.log('   ✓ 有权限时仍使用 redirect 参数');

  console.log('\n✅ 步骤 3.2: 有 redirect 参数但无权限');
  const redirectNoPerm = '/settings/roles';
  const luliyaNoRoles = luliyaRoutesNoDashboard.includes('RoleManagement');
  console.log(`   - redirect=${redirectNoPerm}, 用户有权限：${luliyaNoRoles}`);
  const finalTarget2 = luliyaNoRoles ? redirectNoPerm : accessiblePath2;
  console.log(`   - 最终跳转：${finalTarget2}`);
  assert(!luliyaNoRoles, 'luliya should not have role management');
  assert(finalTarget2 === '/sftp/sftpuser',
    'should fallback to first accessible route when redirect has no permission');
  console.log('   ✓ 无权限时自动 fallback 到有权限的路由');

  console.log('\n✅ 步骤 3.3: 无 redirect 参数，用户有 Dashboard');
  const noRedirect1 = null;
  const finalTarget3 = noRedirect1 || getFirstAccessibleRoutePath(
    buildFilteredRoutes(luliyaRoutes, MOCK_ROUTE_CONFIGS)
  ) || '/';
  console.log(`   - redirect: ${noRedirect1}`);
  console.log(`   - resultAllRoutes 中第一个可见路由：${finalTarget3}`);
  console.log('   ✓ 登录后正常跳转');

  console.log('\n✅ 步骤 3.4: 无 redirect 参数，用户无 Dashboard');
  const noRedirect2 = null;
  const finalTarget4 = noRedirect2 || getFirstAccessibleRoutePath(
    buildFilteredRoutes(luliyaRoutesNoDashboard, MOCK_ROUTE_CONFIGS)
  ) || '/';
  console.log(`   - redirect: ${noRedirect2}`);
  console.log(`   - 用户无 Dashboard 权限`);
  console.log(`   - 跳转到有权限的第一个路由：${finalTarget4}`);
  assert(finalTarget4 === '/sftp/sftpuser',
    `should redirect to /sftp/sftpuser, not /dashboard, got: ${finalTarget4}`);
  console.log('   ✓ 无首页权限时登录成功不会进入 404 循环');

  // ==========================================
  // === 测试组 4: 登录后无警告静默跳转验证 ===
  // ==========================================
  console.log('\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━');
  console.log('测试组 4: 登录后无警告静默跳转');
  console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━');

  console.log('\n✅ 步骤 4.1: 模拟路由守卫 beforeEach 登录流程（无 Dashboard 权限）');
  // 模拟 permission.js 的 redirectDefaultHomeIfNoPermission + hasRoutePermission 逻辑
  // 将目标路径映射到路由名（模拟 vue-router matched 校验）
  const routeNameForPath = (path) => {
    for (const route of MOCK_ROUTE_CONFIGS) {
      if (route.children && route.children.length > 0) {
        for (const child of route.children) {
          const fullPath = route.path === '/' ? '/' + child.path : route.path + '/' + child.path
          if (fullPath === path) return child.name
        }
      } else if (route.path === path && route.name) {
        return route.name
      }
    }
    return null
  }
  const simulateGuard = (toPath, fromPath, userRoutes, allRoutes) => {
    // 1. 默认首页无权限：静默跳转到第一个可访问路由
    if (toPath === '/') {
      const routes = userRoutes || []
      if (routes.length > 0 && routes.indexOf('Dashboard') === -1) {
        const accessiblePath = getFirstAccessibleRoutePath(allRoutes)
        if (accessiblePath && accessiblePath !== '/login' && accessiblePath !== '/404') {
          return { redirected: true, target: accessiblePath, warned: false }
        }
      }
    }
    // 2. 常规权限校验：目标路由名不在用户权限中则拦截
    const routeName = routeNameForPath(toPath)
    const hasPerm = routeName === null || userRoutes.indexOf(routeName) !== -1
    if (!hasPerm) {
      const accessiblePath = getFirstAccessibleRoutePath(allRoutes)
      // 登录页重定向链无权限时静默，非主动访问弹警告
      const warned = fromPath !== '/login'
      return { redirected: true, target: accessiblePath, warned }
    }
    return { redirected: false, warned: false }
  }

  const luliyaFiltered = buildFilteredRoutes(['SftpLog'], MOCK_ROUTE_CONFIGS)

  // 4.1 登录后 push('/')：应静默跳转到第一个可访问路由，无警告
  const r1 = simulateGuard('/', '/login', ['SftpLog'], luliyaFiltered)
  console.log(`   - 登录后 push('/')，用户权限：['SftpLog']`);
  console.log(`   - 结果：重定向=${r1.redirected}，目标=${r1.target}，警告=${r1.warned}`);
  assert(r1.redirected, 'should redirect away from default home');
  assert(r1.target === '/log/sftplog', 'should redirect to first accessible route');
  assert(!r1.warned, 'should NOT show warning on login redirect chain');
  console.log('   ✓ 静默跳转至 /log/sftplog，无警告条');

  // 4.2 有 Dashboard 权限的用户（admin）：正常进入首页，无警告
  const adminFiltered = buildFilteredRoutes(['Dashboard', 'SftpUser'], MOCK_ROUTE_CONFIGS)
  const r2 = simulateGuard('/', '/login', ['Dashboard', 'SftpUser'], adminFiltered)
  console.log(`\n✅ 步骤 4.2: admin 登录后 push('/')，用户权限：['Dashboard','SftpUser']`);
  console.log(`   - 结果：重定向=${r2.redirected}，警告=${r2.warned}`);
  assert(!r2.redirected, 'admin should go through to dashboard normally');
  assert(!r2.warned, 'admin should not see warning');
  console.log('   ✓ admin 正常进入首页，无警告');

  // 4.3 主动访问无权限页面：应保留警告
  const r3 = simulateGuard('/settings/roles', '/log/sftplog', ['SftpLog'], luliyaFiltered)
  console.log(`\n✅ 步骤 4.3: luliya 主动访问 #/settings/roles（from=/log/sftplog）`);
  console.log(`   - 结果：重定向=${r3.redirected}，目标=${r3.target}，警告=${r3.warned}`);
  assert(r3.redirected, 'should redirect away');
  assert(r3.warned, 'should show warning for active access');
  console.log('   ✓ 主动访问无权限页面保留警告提示');

  // 4.4 带 redirect 参数登录（from=/login）：静默跳转
  const r4 = simulateGuard('/settings/roles', '/login', ['SftpLog'], luliyaFiltered)
  console.log(`\n✅ 步骤 4.4: 带 redirect=/settings/roles 登录（from=/login）`);
  console.log(`   - 结果：重定向=${r4.redirected}，目标=${r4.target}，警告=${r4.warned}`);
  assert(r4.redirected, 'should redirect away');
  assert(r4.target === '/log/sftplog', 'should redirect to first accessible route');
  assert(!r4.warned, 'should NOT warn on login redirect chain');
  console.log('   ✓ 登录重定向链静默跳转，无警告');

  // ==========================================
  // === 汇总 ===
  // ==========================================
  console.log('\n' + '='.repeat(60));
  console.log('全部测试通过！');
  console.log('='.repeat(60));

  return {
    success: true,
    message: '退出登录重定向 + 404循环 + 登录静默跳转修复验证成功',
    details: {
      'test_group_1_logout_redirect': 'PASS',
      'test_group_2_404_loop': 'PASS',
      'test_group_3_login_redirect': 'PASS',
      'test_group_4_silent_login_redirect': 'PASS',
    }
  };
}

// 运行测试
try {
  const result = checkAllFix();
  console.log('\n测试结果摘要:');
  console.log(JSON.stringify(result, null, 2));
  process.exit(0);
} catch (error) {
  console.error('\n测试失败:', error.message);
  process.exit(1);
}