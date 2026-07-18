import router from '@/router'
import store from '@/store'
import { Message } from 'element-ui'
import { getToken } from '@/utils/auth'

let timeout = 180000; // 180秒
let warningTime = 30000; // 提前30秒警告
let timer = null;
let warningTimer = null;
let eventListener = null;
let currentTargetPath = null;
let isRunning = false; // 新增状态标记，减少重复判断

// 更新目标路由路径
function updateTargetPath(path) {
  currentTargetPath = path;
}

// 判断是否需要启动计时器
function shouldStartTimer() {
  const hasToken = getToken();
  const isLoginPage = currentTargetPath 
    ? currentTargetPath === '/login' 
    : router.currentRoute.path === '/login';
    
  return hasToken && !isLoginPage;
}

// 重置定时器
function resetTimer() {
  if (!isRunning) return; // 未运行状态直接返回，减少判断
  
  clearTimeout(warningTimer);
  clearTimeout(timer);
  warningTimer = setTimeout(() => {
    Message.warning({
      message: '您即将在30秒后超时，请尽快操作',
      duration: 30000,
      showClose: true
    });
  }, timeout - warningTime);
  
  timer = setTimeout(() => {
    logout();
  }, timeout);
}

// 清除定时器
function clearTimer() {
  clearTimeout(warningTimer);
  clearTimeout(timer);
  warningTimer = null;
  timer = null;
}

// 执行登出操作
function logout() {
  if (!isRunning) return; // 未运行状态直接返回
  
  store.dispatch('user/logout').then(() => {
    Message.success('由于长时间未操作，已自动登出');
    router.push(`/login?redirect=${router.currentRoute.fullPath}`);
  }).catch(err => {
    console.error('登出失败', err);
  });
}

// 注册事件监听
function registerListeners() {
  if (eventListener || !isRunning) return; // 已存在监听器或未运行时直接返回
  
  eventListener = function() {
    resetTimer();
  };
  
  const events = ['click', 'mousemove', 'keydown', 'scroll', 'touchstart'];
  events.forEach(event => {
    document.addEventListener(event, eventListener, true);
  });
}

// 移除事件监听
function removeListeners() {
  if (eventListener) {
    const events = ['click', 'mousemove', 'keydown', 'scroll', 'touchstart'];
    events.forEach(event => {
      document.removeEventListener(event, eventListener, true);
    });
    eventListener = null;
  }
}

// 启动计时器和事件监听
function start() {
  // 只在需要启动且当前未运行时执行
  if (shouldStartTimer() && !isRunning) {
    isRunning = true;
    registerListeners();
    resetTimer();
  }
}

// 停止计时器和事件监听
function stop() {
  if (isRunning) { // 只在运行状态时执行
    isRunning = false;
    removeListeners();
    clearTimer();
  }
}

export default {
  init() {
    start(); // 直接调用start，内部已包含判断
  },
  // start,
  stop,
  // resetTimer,
  checkAndStart() {
    start(); // 直接复用start方法
  },
  updateTargetPath
}
    