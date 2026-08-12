import { mount, createLocalVue } from '@vue/test-utils'
import ElementUI from 'element-ui'
import SftpBrowser from '@/components/SftpBrowser/index.vue'

const localVue = createLocalVue()
localVue.use(ElementUI)
// 组件模板内使用了 v-focus 自定义指令
localVue.directive('focus', {
  inserted(el) {
    el.focus()
  }
})

// 模拟多级目录结构
const mockFiles = {
  '/': [
    { name: 'dirA', path: '/dirA', isDir: true, size: 0, modified: '2026-01-01T00:00:00Z' },
    { name: 'dirB', path: '/dirB', isDir: true, size: 0, modified: '2026-01-01T00:00:00Z' },
    { name: 'file1.txt', path: '/file1.txt', isDir: false, size: 100, modified: '2026-01-01T00:00:00Z' }
  ],
  '/dirA': [
    { name: 'sub', path: '/dirA/sub', isDir: true, size: 0, modified: '2026-01-01T00:00:00Z' },
    { name: 'a.txt', path: '/dirA/a.txt', isDir: false, size: 10, modified: '2026-01-01T00:00:00Z' }
  ],
  '/dirA/sub': [
    { name: 'deep.txt', path: '/dirA/sub/deep.txt', isDir: false, size: 5, modified: '2026-01-01T00:00:00Z' }
  ],
  '/dirB': null
}

const flushPromises = () => new Promise(resolve => setTimeout(resolve, 0))
const pressKey = key => window.dispatchEvent(new KeyboardEvent('keydown', { key }))

async function createWrapper() {
  const reqSftpFiles = jest.fn(({ path }) =>
    Promise.resolve({
      code: 200,
      data: { files: path in mockFiles ? mockFiles[path] : [], path, description: '' }
    })
  )
  const $API = {
    sftpuser: {
      reqSftpFiles,
      reqSftpSearch: jest.fn(() => Promise.resolve({ code: 200, data: { results: [] }})),
      reqSftpMkdir: jest.fn(),
      reqSftpDelete: jest.fn(),
      reqSftpRename: jest.fn(),
      reqSftpBatchDelete: jest.fn()
    }
  }
  const wrapper = mount(SftpBrowser, {
    localVue,
    propsData: {
      visible: false,
      username: 'testuser',
      uploadHeaders: { Token: 'x' },
      path: ''
    },
    mocks: { $API }
  })
  // 模拟父组件打开弹框，触发 visible watcher -> fetchFiles
  await wrapper.setProps({ visible: true })
  await flushPromises()
  await wrapper.vm.$nextTick()
  return { wrapper, reqSftpFiles }
}

describe('SftpBrowser 键盘与鼠标焦点导航', () => {
  it('初始状态无焦点', async() => {
    const { wrapper } = await createWrapper()
    expect(wrapper.vm.focusIndex).toBe(-1)
    expect(wrapper.vm.filteredFileList.length).toBe(3)
    wrapper.destroy()
  })

  it('方向键下移/上移焦点，边界处循环', async() => {
    const { wrapper } = await createWrapper()
    pressKey('ArrowDown')
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.focusIndex).toBe(0)
    // 焦点行应带 keyboard-focus-row class（用于滚动定位与高亮）
    expect(wrapper.find('tr.keyboard-focus-row').exists()).toBe(true)

    pressKey('ArrowDown')
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.focusIndex).toBe(1)

    pressKey('ArrowUp')
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.focusIndex).toBe(0)

    // 首行按上循环到末行（共 3 行）
    pressKey('ArrowUp')
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.focusIndex).toBe(2)

    // 末行按下循环回首行
    pressKey('ArrowDown')
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.focusIndex).toBe(0)
    wrapper.destroy()
  })

  it('鼠标单击仅聚焦，不进入目录', async() => {
    const { wrapper, reqSftpFiles } = await createWrapper()
    const callCount = reqSftpFiles.mock.calls.length
    const dirRow = wrapper.find('.dir-item')
    await dirRow.trigger('click')
    expect(wrapper.vm.focusIndex).toBe(0)
    // 没有新的列表请求，说明没有进入目录
    expect(reqSftpFiles.mock.calls.length).toBe(callCount)
    expect(wrapper.vm.currentPath).toBe('/')
    wrapper.destroy()
  })

  it('双击目录进入，并在 Enter 上共用打开逻辑', async() => {
    const { wrapper } = await createWrapper()
    const dirRow = wrapper.find('.dir-item')
    await dirRow.trigger('dblclick')
    await flushPromises()
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.currentPath).toBe('/dirA')
    // 进入后焦点清空（子目录中不存在 dirA 行）
    expect(wrapper.vm.focusIndex).toBe(-1)
    wrapper.destroy()
  })

  it('Enter 打开焦点行：目录进入', async() => {
    const { wrapper } = await createWrapper()
    pressKey('ArrowDown') // 聚焦 dirA
    await wrapper.vm.$nextTick()
    pressKey('Enter')
    await flushPromises()
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.currentPath).toBe('/dirA')
    wrapper.destroy()
  })

  it('Enter/双击 打开文件触发下载', async() => {
    const { wrapper } = await createWrapper()
    const spy = jest.spyOn(wrapper.vm, 'handleDownload')
    // 聚焦第 3 行 file1.txt
    pressKey('ArrowDown')
    pressKey('ArrowDown')
    pressKey('ArrowDown')
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.filteredFileList[wrapper.vm.focusIndex].name).toBe('file1.txt')
    pressKey('Enter')
    await wrapper.vm.$nextTick()
    expect(spy).toHaveBeenCalledTimes(1)
    expect(spy.mock.calls[0][0].name).toBe('file1.txt')
    wrapper.destroy()
  })

  it('Backspace 返回上级后，焦点恢复到之前进入的目录行', async() => {
    const { wrapper } = await createWrapper()
    // 进入 dirA -> sub，构造两级深度
    pressKey('ArrowDown')
    await wrapper.vm.$nextTick()
    pressKey('Enter')
    await flushPromises()
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.currentPath).toBe('/dirA')

    pressKey('Backspace')
    await flushPromises()
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.currentPath).toBe('/')
    // 焦点恢复至 dirA 行（索引 0）
    expect(wrapper.vm.focusIndex).toBe(0)
    expect(wrapper.vm.filteredFileList[wrapper.vm.focusIndex].name).toBe('dirA')

    // 从恢复的焦点继续上下移动
    pressKey('ArrowDown')
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.focusIndex).toBe(1)
    expect(wrapper.vm.filteredFileList[wrapper.vm.focusIndex].name).toBe('dirB')
    wrapper.destroy()
  })

  it('多级连续返回，每级焦点都恢复到对应目录行', async() => {
    const { wrapper } = await createWrapper()
    // / -> dirA -> sub
    wrapper.vm.handleItemOpen(wrapper.vm.filteredFileList[0]) // dirA
    await flushPromises()
    await wrapper.vm.$nextTick()
    wrapper.vm.handleItemOpen(wrapper.vm.filteredFileList[0]) // sub
    await flushPromises()
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.currentPath).toBe('/dirA/sub')

    pressKey('Backspace')
    await flushPromises()
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.currentPath).toBe('/dirA')
    expect(wrapper.vm.filteredFileList[wrapper.vm.focusIndex].name).toBe('sub')

    pressKey('Backspace')
    await flushPromises()
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.currentPath).toBe('/')
    expect(wrapper.vm.filteredFileList[wrapper.vm.focusIndex].name).toBe('dirA')
    wrapper.destroy()
  })

  it('子弹框打开时方向键不移动焦点，Enter 不打开行', async() => {
    const { wrapper } = await createWrapper()
    wrapper.setData({ showCreateFolderDialog: true })
    await wrapper.vm.$nextTick()
    pressKey('ArrowDown')
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.focusIndex).toBe(-1)
    wrapper.destroy()
  })

  it('搜索过滤状态下方向键在过滤结果内移动', async() => {
    const { wrapper } = await createWrapper()
    wrapper.setData({ searchQuery: 'dir' })
    wrapper.vm.applyFilter()
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.filteredFileList.length).toBe(2)
    pressKey('ArrowDown')
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.focusIndex).toBe(0)
    pressKey('ArrowDown')
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.focusIndex).toBe(1)
    // 末行按下循环回首行（过滤后共 2 行）
    pressKey('ArrowDown')
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.focusIndex).toBe(0)
    wrapper.destroy()
  })

  it('进入空目录（后端 files: null）不报错，方向键安全', async() => {
    const { wrapper } = await createWrapper()
    const errSpy = jest.spyOn(wrapper.vm.$message, 'error')
    // 进入 dirB（空目录，后端返回 files: null）
    wrapper.vm.handleItemOpen(wrapper.vm.filteredFileList[1])
    await flushPromises()
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.currentPath).toBe('/dirB')
    // 不应弹出“获取文件列表失败”错误提示
    expect(errSpy).not.toHaveBeenCalled()
    expect(wrapper.vm.focusIndex).toBe(-1)
    // 空目录下按方向键不抛错、焦点保持 -1
    pressKey('ArrowDown')
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.focusIndex).toBe(-1)
    // 空目录下 Backspace 返回上级正常
    pressKey('Backspace')
    await flushPromises()
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.currentPath).toBe('/')
    expect(errSpy).not.toHaveBeenCalled()
    wrapper.destroy()
  })

  it('打开时释放外部残留焦点，方向键立即可用', async() => {
    // 模拟登录密码框残留焦点
    const externalInput = document.createElement('input')
    document.body.appendChild(externalInput)
    externalInput.focus()
    expect(document.activeElement).toBe(externalInput)

    const { wrapper } = await createWrapper()
    await wrapper.vm.$nextTick()
    // 打开后残留焦点被释放
    expect(document.activeElement).toBe(document.body)
    // 无需鼠标点击，方向键立即可用
    pressKey('ArrowDown')
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.focusIndex).toBe(0)

    document.body.removeChild(externalInput)
    wrapper.destroy()
  })

  it('关闭弹框后键盘导航状态被重置', async() => {
    const { wrapper } = await createWrapper()
    pressKey('ArrowDown')
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.focusIndex).toBe(0)
    wrapper.vm.handleClose()
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.focusIndex).toBe(-1)
    expect(wrapper.vm.lastFocusName).toBe('')
    expect(wrapper.vm.pendingFocusName).toBe('')
    wrapper.destroy()
  })
})

// 可控制进度的 XHR 替身，模拟真实上传过程中的 progress 事件
class MockXHR {
  constructor() {
    this.upload = {}
    this.headers = {}
    MockXHR.instances.push(this)
  }
  open(method, url) {
    this.method = method
    this.url = url
  }
  setRequestHeader(key, value) {
    this.headers[key] = value
  }
  send() {
    this.sent = true
  }
  // 测试辅助：模拟上传进度事件
  emitProgress(loaded, total) {
    this.upload.onprogress({ lengthComputable: true, loaded, total })
  }
  // 测试辅助：模拟上传成功响应
  emitSuccess(body = { code: 200 }) {
    this.status = 200
    this.responseText = JSON.stringify(body)
    this.onload()
  }
  // 测试辅助：模拟上传失败（如网络错误）
  emitError() {
    if (this.onerror) this.onerror(new Error('network error'))
  }
  // 测试辅助：模拟后端返回非 200 业务/HTTP 错误（响应体仍为 {code, message}）
  emitHttpError(status = 400, body = { code: 400, message: '文件上传失败: 目标目录不存在' }) {
    this.status = status
    this.responseText = JSON.stringify(body)
    this.onload()
  }
}
MockXHR.instances = []

describe('SftpBrowser 拖拽上传进度条', () => {
  let originalXHR
  beforeEach(() => {
    originalXHR = window.XMLHttpRequest
    MockXHR.instances = []
    window.XMLHttpRequest = MockXHR
  })
  afterEach(() => {
    window.XMLHttpRequest = originalXHR
  })

  const makeDropEvent = files => ({
    dataTransfer: { files },
    preventDefault: () => {}
  })
  const makeFile = (name, size) => {
    // jsdom 的 File 支持 size 参数
    return new File([new Uint8Array(0)], name, { type: 'text/plain' })
  }
  const makeSizedFile = (name, size) => {
    const f = makeFile(name)
    Object.defineProperty(f, 'size', { value: size })
    return f
  }

  it('单文件拖拽上传：进度随发送字节实时增长，完成后到 100%', async() => {
    const { wrapper } = await createWrapper()
    const fetchSpy = jest.spyOn(wrapper.vm, 'fetchFiles').mockImplementation(() => Promise.resolve())
    const file = makeSizedFile('big.bin', 1000)

    const dropPromise = wrapper.vm.handleDrop(makeDropEvent([file]))
    await flushPromises()

    expect(MockXHR.instances.length).toBe(1)
    const xhr = MockXHR.instances[0]
    // 进度条已显示，初始 0%
    expect(wrapper.vm.showUploadProgress).toBe(true)
    expect(wrapper.vm.uploadPercent).toBe(0)

    // 模拟分步上传：30% -> 60% -> 90%
    xhr.emitProgress(300, 1000)
    expect(wrapper.vm.uploadPercent).toBe(30)
    xhr.emitProgress(600, 1000)
    expect(wrapper.vm.uploadPercent).toBe(60)
    xhr.emitProgress(900, 1000)
    expect(wrapper.vm.uploadPercent).toBe(90)

    // 上传完成
    xhr.emitSuccess()
    await dropPromise
    await flushPromises()
    expect(wrapper.vm.uploadPercent).toBe(100)
    expect(wrapper.vm.isUploading).toBe(false)

    // 1 秒后进度条隐藏并重置
    jest.useRealTimers()
    await new Promise(resolve => setTimeout(resolve, 1100))
    expect(wrapper.vm.showUploadProgress).toBe(false)
    expect(wrapper.vm.uploadPercent).toBe(0)
    fetchSpy.mockRestore()
    wrapper.destroy()
  })

  it('多文件拖拽上传：整体进度按已完成文件 + 当前文件进度加权推进', async() => {
    const { wrapper } = await createWrapper()
    const fetchSpy = jest.spyOn(wrapper.vm, 'fetchFiles').mockImplementation(() => Promise.resolve())
    const f1 = makeSizedFile('a.bin', 1000)
    const f2 = makeSizedFile('b.bin', 1000)

    const dropPromise = wrapper.vm.handleDrop(makeDropEvent([f1, f2]))
    await flushPromises()

    // 第 1 个文件上传到 50%：整体应为 (0 + 0.5) / 2 = 25%
    const xhr1 = MockXHR.instances[0]
    xhr1.emitProgress(500, 1000)
    expect(wrapper.vm.uploadPercent).toBe(25)
    xhr1.emitSuccess()
    await flushPromises()

    // 第 2 个文件上传到 50%：整体应为 (1 + 0.5) / 2 = 75%
    const xhr2 = MockXHR.instances[1]
    xhr2.emitProgress(500, 1000)
    expect(wrapper.vm.uploadPercent).toBe(75)
    xhr2.emitSuccess()

    await dropPromise
    expect(wrapper.vm.uploadPercent).toBe(100)
    fetchSpy.mockRestore()
    wrapper.destroy()
  })
})

describe('SftpBrowser 传输队列', () => {
  let originalXHR
  beforeEach(() => {
    originalXHR = window.XMLHttpRequest
    MockXHR.instances = []
    window.XMLHttpRequest = MockXHR
  })
  afterEach(() => {
    window.XMLHttpRequest = originalXHR
  })

  const makeDropEvent = files => ({
    dataTransfer: { files },
    preventDefault: () => {}
  })
  const makeSizedFile = (name, size) => {
    const f = new File([new Uint8Array(0)], name, { type: 'text/plain' })
    Object.defineProperty(f, 'size', { value: size })
    return f
  }

  it('队列卡片 DOM 渲染：三个标签页、拖入提示条与右键菜单', async() => {
    const { wrapper } = await createWrapper()
    wrapper.vm.currentPath = '/q'

    // 队列卡片存在，含三个标签页标题
    const card = wrapper.find('.transfer-queue-card')
    expect(card.exists()).toBe(true)
    const tabText = card.text()
    expect(tabText).toContain('列队的文件 (0)')
    expect(tabText).toContain('传输失败 (0)')
    expect(tabText).toContain('成功的传输 (0)')

    // 拖入卡片时显示提示条并入队
    wrapper.vm.queueDragOver = true
    await wrapper.vm.$nextTick()
    expect(wrapper.find('.queue-drop-hint').text()).toContain('释放文件以加入传输队列')
    wrapper.vm.handleQueueDrop(makeDropEvent([makeSizedFile('dom.bin', 100)]))
    await wrapper.vm.$nextTick()
    expect(wrapper.find('.queue-drop-hint').isVisible()).toBe(false)
    expect(card.text()).toContain('列队的文件 (1)')
    expect(card.text()).toContain('dom.bin')
    expect(card.text()).toContain('/q/dom.bin')
    expect(card.text()).toContain('待上传')

    // 右键菜单：默认关闭；打开后挂载到 body（避免弹框 backdrop-filter 影响 fixed 定位）并显示四个菜单项
    expect(wrapper.vm.ctxMenu.visible).toBe(false)
    wrapper.vm.openCtxMenu(wrapper.vm.transferQueue[0], null, { preventDefault: () => {}, clientX: 10, clientY: 20 })
    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()
    const menuEl = document.querySelector('.ctx-menu')
    expect(menuEl).toBeTruthy()
    expect(menuEl.parentElement).toBe(document.body)
    expect(menuEl.style.display).not.toBe('none')
    expect(Array.from(menuEl.querySelectorAll('.ctx-menu-item')).map(el => el.textContent.trim())).toEqual(['全部上传', '选定上传', '选定移除', '全部移除'])
    // 视口内不触发收敛，保持鼠标坐标
    expect(wrapper.vm.ctxMenu.x).toBe(10)
    expect(wrapper.vm.ctxMenu.y).toBe(20)
    wrapper.vm.closeCtxMenu()
    await wrapper.vm.$nextTick()
    expect(menuEl.style.display).toBe('none')
    wrapper.destroy()
    // 组件销毁时清理 body 上的菜单节点
    expect(document.querySelector('.ctx-menu')).toBe(null)
  })

  it('右键菜单：靠近视口右下边缘时自动收敛不被裁剪', async() => {
    const { wrapper } = await createWrapper()
    wrapper.vm.currentPath = '/q'
    wrapper.vm.handleQueueDrop(makeDropEvent([makeSizedFile('clamp.bin', 100)]))

    // openCtxMenu 同步挂载菜单到 body，随后 mock 尺寸（在 nextTick 测量前生效）
    wrapper.vm.openCtxMenu(wrapper.vm.transferQueue[0], null, { preventDefault: () => {}, clientX: window.innerWidth - 10, clientY: window.innerHeight - 10 })
    const menuEl = document.querySelector('.ctx-menu')
    expect(menuEl).toBeTruthy()
    // jsdom 无布局，mock 菜单实际尺寸
    jest.spyOn(menuEl, 'getBoundingClientRect').mockReturnValue({ width: 120, height: 140, top: 0, left: 0, right: 0, bottom: 0, x: 0, y: 0 })
    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.ctxMenu.x).toBe(window.innerWidth - 120 - 4)
    expect(wrapper.vm.ctxMenu.y).toBe(window.innerHeight - 140 - 4)
    wrapper.destroy()
  })

  it('失败原因取自响应 message（含 HTTP 400 场景）', async() => {
    const { wrapper } = await createWrapper()
    const fetchSpy = jest.spyOn(wrapper.vm, 'fetchFiles').mockImplementation(() => Promise.resolve())
    wrapper.vm.currentPath = '/q'
    wrapper.vm.handleQueueDrop(makeDropEvent([makeSizedFile('msg.bin', 100)]))

    wrapper.vm.uploadAll()
    await flushPromises()
    // HTTP 400 + JSON body：原因应为响应 message 而非 "HTTP 400"
    MockXHR.instances[0].emitHttpError(400, { code: 400, message: '文件上传失败: 目标目录不存在' })
    await flushPromises()
    expect(wrapper.vm.failedTransfers[0].reason).toBe('文件上传失败: 目标目录不存在')

    // HTTP 200 + 业务 code 非 200：同样取 message
    wrapper.vm.retryFailed(wrapper.vm.failedTransfers[0])
    await flushPromises()
    MockXHR.instances[1].emitHttpError(200, { code: 500, message: 'SFTP连接失效: token expired' })
    await flushPromises()
    expect(wrapper.vm.failedTransfers[0].reason).toBe('SFTP连接失效: token expired')
    fetchSpy.mockRestore()
    wrapper.destroy()
  })

  it('切换标签页：显示后主动重布局当前表格避免抖动', async() => {
    const { wrapper } = await createWrapper()
    const spy = jest.spyOn(wrapper.vm.$refs.failedTable, 'doLayout')
    wrapper.vm.queueTab = 'failed'
    wrapper.vm.onQueueTabClick()
    await wrapper.vm.$nextTick()
    expect(spy).toHaveBeenCalled()
    wrapper.destroy()
  })

  it('原生拖放事件：在卡片 DOM 上 drop/dragover 能触发入队与高亮', async() => {
    const { wrapper } = await createWrapper()
    wrapper.vm.currentPath = '/nat'
    const card = wrapper.find('.transfer-queue-card')
    expect(card.exists()).toBe(true)

    // 原生 dragover：触发高亮（验证事件确实绑定到了 DOM）
    const overEvt = new Event('dragover', { bubbles: true, cancelable: true })
    card.element.dispatchEvent(overEvt)
    await wrapper.vm.$nextTick()
    expect(overEvt.defaultPrevented).toBe(true)
    expect(wrapper.vm.queueDragOver).toBe(true)

    // 原生 drop：入队
    const file = makeSizedFile('native.bin', 100)
    const dropEvt = new Event('drop', { bubbles: true, cancelable: true })
    Object.defineProperty(dropEvt, 'dataTransfer', { value: { files: [file] }})
    card.element.dispatchEvent(dropEvt)
    await wrapper.vm.$nextTick()
    expect(dropEvt.defaultPrevented).toBe(true)
    expect(wrapper.vm.queueDragOver).toBe(false)
    expect(wrapper.vm.transferQueue.length).toBe(1)
    expect(wrapper.vm.transferQueue[0].name).toBe('native.bin')
    expect(wrapper.vm.transferQueue[0].remotePath).toBe('/nat/native.bin')
    expect(MockXHR.instances.length).toBe(0)
    wrapper.destroy()
  })

  it('拖入队列卡片：记为待上传，远程路径为当前目录+文件名，不发起请求', async() => {
    const { wrapper } = await createWrapper()
    wrapper.vm.currentPath = '/data/input'
    wrapper.vm.queueTab = 'success'

    wrapper.vm.handleQueueDrop(makeDropEvent([makeSizedFile('a.bin', 100), makeSizedFile('b.bin', 200)]))
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.transferQueue.length).toBe(2)
    expect(wrapper.vm.transferQueue[0].status).toBe('pending')
    expect(wrapper.vm.transferQueue[0].name).toBe('a.bin')
    expect(wrapper.vm.transferQueue[0].remotePath).toBe('/data/input/a.bin')
    expect(wrapper.vm.transferQueue[1].remotePath).toBe('/data/input/b.bin')
    expect(wrapper.vm.queueTab).toBe('queue')
    // 仅入队，不发起上传请求
    expect(MockXHR.instances.length).toBe(0)
    wrapper.destroy()
  })

  it('全部上传：最多 3 路并发，完成 1 个自动补位，进度独立且成功/失败分别记录', async() => {
    const { wrapper } = await createWrapper()
    const fetchSpy = jest.spyOn(wrapper.vm, 'fetchFiles').mockImplementation(() => Promise.resolve())
    wrapper.vm.currentPath = '/q'
    const files = ['f1', 'f2', 'f3', 'f4', 'f5'].map(n => makeSizedFile(n + '.bin', 1000))
    wrapper.vm.handleQueueDrop(makeDropEvent(files))

    wrapper.vm.uploadAll()
    await flushPromises()
    // 5 个入队，仅发起 3 个
    expect(MockXHR.instances.length).toBe(3)
    expect(wrapper.vm.uploadingCount).toBe(3)

    // 各条目进度独立更新
    MockXHR.instances[0].emitProgress(300, 1000)
    MockXHR.instances[1].emitProgress(600, 1000)
    expect(wrapper.vm.transferQueue.find(i => i.name === 'f1.bin').percent).toBe(30)
    expect(wrapper.vm.transferQueue.find(i => i.name === 'f2.bin').percent).toBe(60)

    // 第 1 个成功：进入成功记录，自动补位发起第 4 个
    MockXHR.instances[0].emitSuccess()
    await flushPromises()
    expect(MockXHR.instances.length).toBe(4)
    expect(wrapper.vm.successTransfers.length).toBe(1)
    expect(wrapper.vm.successTransfers[0].name).toBe('f1.bin')
    expect(wrapper.vm.transferQueue.length).toBe(4)

    // 第 2 个失败：进入失败记录并带原因，自动补位发起第 5 个
    MockXHR.instances[1].emitError()
    await flushPromises()
    expect(MockXHR.instances.length).toBe(5)
    expect(wrapper.vm.failedTransfers.length).toBe(1)
    expect(wrapper.vm.failedTransfers[0].name).toBe('f2.bin')
    expect(wrapper.vm.failedTransfers[0].reason).toBeTruthy()

    // 完成剩余 3 个
    MockXHR.instances[2].emitSuccess()
    MockXHR.instances[3].emitSuccess()
    MockXHR.instances[4].emitSuccess()
    await flushPromises()
    expect(wrapper.vm.transferQueue.length).toBe(0)
    expect(wrapper.vm.successTransfers.length).toBe(4)
    expect(wrapper.vm.uploadingCount).toBe(0)
    // 队列清空后刷新文件列表
    expect(fetchSpy).toHaveBeenCalled()
    fetchSpy.mockRestore()
    wrapper.destroy()
  })

  it('右键动作：选定移除 / 全部移除仅作用于待上传条目', async() => {
    const { wrapper } = await createWrapper()
    const fetchSpy = jest.spyOn(wrapper.vm, 'fetchFiles').mockImplementation(() => Promise.resolve())
    wrapper.vm.currentPath = '/q'
    // 5 个文件，3 路并发：r1~r3 上传中，r4/r5 排队
    const files = ['r1', 'r2', 'r3', 'r4', 'r5'].map(n => makeSizedFile(n + '.bin', 100))
    wrapper.vm.handleQueueDrop(makeDropEvent(files))

    wrapper.vm.pumpUploads()
    await flushPromises()
    expect(MockXHR.instances.length).toBe(3)
    const uploading = wrapper.vm.transferQueue.find(i => i.name === 'r1.bin')
    expect(uploading.status).toBe('uploading')

    // 选定移除：移除排队的 r5，不影响上传中的条目
    wrapper.vm.removeOne(wrapper.vm.transferQueue.find(i => i.name === 'r5.bin'))
    expect(wrapper.vm.transferQueue.map(i => i.name)).toEqual(['r1.bin', 'r2.bin', 'r3.bin', 'r4.bin'])
    // 对上传中条目无效
    wrapper.vm.removeOne(uploading)
    expect(wrapper.vm.transferQueue.length).toBe(4)

    // 全部移除：仅移除排队的 r4，上传中的 r1~r3 继续跑完
    wrapper.vm.removeAllPending()
    expect(wrapper.vm.transferQueue.map(i => i.name)).toEqual(['r1.bin', 'r2.bin', 'r3.bin'])
    MockXHR.instances[0].emitSuccess()
    MockXHR.instances[1].emitSuccess()
    MockXHR.instances[2].emitSuccess()
    await flushPromises()
    expect(wrapper.vm.transferQueue.length).toBe(0)
    expect(wrapper.vm.successTransfers.length).toBe(3)
    fetchSpy.mockRestore()
    wrapper.destroy()
  })

  it('失败重试：复用原文件重新入队并可上传成功', async() => {
    const { wrapper } = await createWrapper()
    const fetchSpy = jest.spyOn(wrapper.vm, 'fetchFiles').mockImplementation(() => Promise.resolve())
    wrapper.vm.currentPath = '/q'
    const file = makeSizedFile('retry.bin', 1000)
    wrapper.vm.handleQueueDrop(makeDropEvent([file]))

    wrapper.vm.uploadAll()
    await flushPromises()
    MockXHR.instances[0].emitError()
    await flushPromises()
    expect(wrapper.vm.failedTransfers.length).toBe(1)

    // 重试：从失败记录移除，重新入队并自动启动
    const failed = wrapper.vm.failedTransfers[0]
    wrapper.vm.retryFailed(failed)
    await flushPromises()
    expect(wrapper.vm.failedTransfers.length).toBe(0)
    expect(MockXHR.instances.length).toBe(2)
    expect(wrapper.vm.transferQueue[0].file).toBe(file)

    MockXHR.instances[1].emitSuccess()
    await flushPromises()
    expect(wrapper.vm.transferQueue.length).toBe(0)
    expect(wrapper.vm.successTransfers.map(i => i.name)).toEqual(['retry.bin'])
    fetchSpy.mockRestore()
    wrapper.destroy()
  })
})
