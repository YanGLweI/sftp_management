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
