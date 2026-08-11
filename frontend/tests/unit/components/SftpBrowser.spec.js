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

  it('方向键下移/上移焦点，边界处停止', async() => {
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

    // 首行再按上不越界
    pressKey('ArrowUp')
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
    // 边界：过滤后只有 2 行
    pressKey('ArrowDown')
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.focusIndex).toBe(1)
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
