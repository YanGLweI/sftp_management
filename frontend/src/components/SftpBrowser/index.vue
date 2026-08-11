<template>
  <el-dialog
    :title="`SFTP浏览器 - ${username}`"
    :visible.sync="innerVisible"
    center
    width="70%"
    top="5vh"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <div class="sftp-browser">
      <!-- 面包屑导航 -->
      <el-breadcrumb separator-class="el-icon-arrow-right">
        <el-breadcrumb-item
          v-for="(item, index) in breadcrumb"
          :key="index"
        >
          <span
            @click="handleBreadcrumbClick(item, index)"
            :class="breadcrumb.length - 1 == index ? 'breadcrumbBold' : 'breadcrumb'"
          >{{ item.name }}</span>
        </el-breadcrumb-item>
      </el-breadcrumb>

      <!-- 操作按钮 -->
      <div class="operate">
        <el-page-header
          @back="goBack"
          :content="breadcrumb[breadcrumb.length - 1].name"
          style="margin:20px 0;width: 300px;"
          class="goBack"
        >
        </el-page-header>

        <div style="width: 45%;">
          <el-progress 
            v-if="showUploadProgress"
            :percentage="uploadPercent" 
            size="mini" 
            :color="customColors"
            class="upload-progress"
          ></el-progress>  
        </div>

        <div style="width: 360px; display: flex; align-items: center; flex-wrap: nowrap; gap: 8px;">
          <el-button icon="el-icon-refresh" circle size="mini" @click="fetchFiles()"></el-button>
          <el-upload
            class="upload"
            :action="uploadUrl"
            :data="{ path: currentPath }"
            :headers="uploadHeaders"
            :on-success="handleUploadSuccess"
            :on-error="handleUploadError"
            :before-upload="beforeUpload"
            :show-file-list="false"
            multiple
            :on-progress="handleUploadProgress"
          >
            <el-button type="primary" size="mini" icon="el-icon-document-add" round>上传文件</el-button>
          </el-upload>

          <el-button
            type="primary"
            size="mini"
            icon="el-icon-folder-add"
            round
            @click="handleCreateDir"
          >创建目录</el-button>

          <el-button
            type="warning"
            size="mini"
            icon="el-icon-search"
            round
            @click="openDeepSearch"
          >搜索</el-button>
        </div>
      </div>
      <!-- 文件列表 -->
      <el-card shadow="hover">
        <div class="table-container" ref="tableContainer" @dragleave.prevent="onTableDragLeave">
          <el-empty
            description="这个目录很穷，什么都没有_(:3 ∠)_"
            v-if="!fileList"
            style="height: calc(100vh - 400px);"
          ></el-empty>

          <el-table
            ref="multipleTable"
            @selection-change="handleSelectionChange"
            :data="filteredFileList"
            v-loading="isLoading"
            :height="searchQuery ? 'calc(100vh - 440px)' : 'calc(100vh - 400px)'"
            border
            v-if="fileList"
          >
            <el-table-column type="selection" width="55" align="center"></el-table-column>
            <el-table-column prop="name" label="名称" sortable show-overflow-tooltip>
              <template slot-scope="{row}">
                <div v-if="row.isRenaming">
                  <el-input
                    v-model="row.editName"
                    size="mini"
                    @keyup.enter.native="confirmRename(row)"
                    @blur="cancelRename(row)"
                    v-focus="true"
                  ></el-input>
                </div>
                <div
                  v-else
                  @click="handleItemClick(row)"
                  :class="{'dir-item': row.isDir, 'file-item': !row.isDir}"
                >
                  <i :class="row.isDir ? 'el-icon-folder' : 'el-icon-document'"></i>
                  {{ row.name }}
                </div>
              </template>
            </el-table-column>

            <el-table-column prop="size" label="大小" width="120" sortable>
              <template slot-scope="{row}">
                {{ row.isDir ? '-' : formatSize(row.size) }}
              </template>
            </el-table-column>

            <el-table-column prop="modified" label="修改时间" width="200" sortable>
              <template slot-scope="{row}">
                {{ formatDate(row.modified) }}
              </template>
            </el-table-column>

            <el-table-column align="center" label="操作" width="150" fixed="right">
              <template slot-scope="{row}">
                <el-button
                  v-if="!row.isDir"
                  size="mini"
                  type="primary"
                  @click="handleDownload(row)"
                  circle
                  icon="el-icon-download"
                ></el-button>
                <el-button
                  v-else
                  size="mini"
                  type="primary"
                  @click="handleDownloadDir(row)"
                  circle
                  icon="el-icon-download"
                ></el-button>
                <el-button
                  size="mini"
                  circle
                  icon="el-icon-edit"
                  @click="startRename(row)"
                ></el-button>
                <el-button
                  size="mini"
                  type="danger"
                  @click="handleDelete(row)"
                  circle
                  icon="el-icon-delete"
                ></el-button>
              </template>
            </el-table-column>
          </el-table>

          <!-- 搜索框 -->
          <div class="search-box-container" v-show="showSearchBox">
            <div class="search-box">
              <el-input
                ref="searchInput"
                v-model="searchQuery"
                size="small"
                placeholder="请输入文件名关键字..."
                prefix-icon="el-icon-search"
                clearable
                @input="applyFilter"
                @keyup.native.esc.native="handleCloseSearchBox"
              ></el-input>
              <el-button
                size="small"
                icon="el-icon-close"
                circle
                @click="handleCloseSearchBox"
              ></el-button>
            </div>
          </div>
          
          <div style="margin-top: 10px;margin-left: 10px; display: flex; justify-content: space-between; align-items: center;">
            <span style="color: #909399;" v-if="selectedFiles.length > 0">{{ selectDescription }}</span>
            <span style="color: #909399;" v-else>{{ description }}</span>
            <div>
              <el-button type="primary" size="mini" icon="el-icon-download" round @click="handleBatchDownload" :disabled="selectedFiles.length === 0">批量下载</el-button>
              <el-button type="danger" size="mini" icon="el-icon-delete" round @click="batchDeleteDialogVisible = true" :disabled="selectedFiles.length === 0">批量删除</el-button>
            </div>
          </div>

          <!-- 拖拽上传遮罩层，默认隐藏，拖拽文件进入时显示 -->
          <div 
            class="drag-overlay" 
            v-show="dragOverlayVisible"
            @dragover.prevent
            @dragenter.prevent="keepOverlayVisible"
            @drop.prevent="handleDrop"
          >
            <div class="drag-overlay-content">
              <i class="el-icon-upload2"></i>
              <span>释放文件以上传</span>
              <span class="drag-sub">支持多文件，最大5GB/文件</span>
            </div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 新建文件夹 -->
    <el-dialog
      title="新建文件夹"
      :visible.sync="showCreateFolderDialog"
      width="30%"
      append-to-body
      :close-on-click-modal="false"
    >
      <el-form @submit.native.prevent>
        <el-form-item label="目录名称">
          <el-input v-model="newFolderName" autocomplete="off" ref="newDirInput" @keydown.enter.native="createFolder"></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer">
        <el-button @click="showCreateFolderDialog = false">取消</el-button>
        <el-button type="primary" @click="createFolder">确定</el-button>
      </div>
    </el-dialog>

    <!-- 删除确认 -->
    <el-dialog
      title="确认删除"
      :visible.sync="deleteDialogVisible"
      width="30%"
      append-to-body
      :close-on-click-modal="false"
    >
      <span>确定要删除 {{ deleteTarget.name }} 吗？</span>
      <span slot="footer">
        <el-button @click="deleteDialogVisible = false">取消</el-button>
        <el-button type="danger" @click="confirmDelete">确定</el-button>
      </span>
    </el-dialog>
    <!-- 递归搜索弹框 -->
    <el-dialog
      title="递归搜索"
      :visible.sync="showDeepSearchDialog"
      width="65%"
      append-to-body
      :close-on-click-modal="false"
    >
      <el-form :inline="true" @submit.native.prevent="doDeepSearch" size="small">
        <el-form-item label="搜索路径">
          <el-input v-model="deepSearchPath" placeholder="输入搜索路径" style="width:280px" clearable></el-input>
        </el-form-item>
        <el-form-item label="关键字">
          <el-input v-model="deepSearchKeyword" placeholder="输入关键字" style="width:220px" clearable
            @keyup.enter.native="doDeepSearch"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="el-icon-search" @click="doDeepSearch" :loading="isSearching">搜索</el-button>
        </el-form-item>
      </el-form>

      <!-- 骨架屏预占位（未搜索时显示） -->
      <div v-if="!deepSearchHasSearched && !isSearching" class="skeleton-table-wrapper">
        <el-skeleton :rows="8" animated>
          <template slot="header">
            <div style="display:flex;gap:16px;align-items:center;padding-bottom:12px;">
              <el-skeleton-item variant="rect" style="width:55px;height:20px;" />
              <el-skeleton-item variant="rect" style="width:200px;height:20px;" />
              <el-skeleton-item variant="rect" style="width:80px;height:20px;" />
              <el-skeleton-item variant="rect" style="width:160px;height:20px;" />
              <el-skeleton-item variant="rect" style="width:120px;height:20px;" />
            </div>
          </template>
        </el-skeleton>
      </div>

      <!-- 搜索结果表格 -->
      <el-table
        :data="deepSearchResults"
        v-loading="isSearching"
        border
        max-height="400"
        v-if="deepSearchResults.length > 0 || isSearching"
      >
        <el-table-column prop="name" label="名称" sortable show-overflow-tooltip>
          <template slot-scope="{row}">
            <div v-if="row.isRenaming">
              <el-input
                v-model="row.editName"
                size="mini"
                @keyup.enter.native="confirmDeepSearchRename(row)"
                @blur="cancelDeepSearchRename(row)"
                v-focus="true"
              ></el-input>
            </div>
            <div v-else :class="{'dir-item': row.isDir, 'file-item': !row.isDir}">
              <i :class="row.isDir ? 'el-icon-folder' : 'el-icon-document'"></i>
              {{ row.name }}
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="parentPath" label="所在路径" show-overflow-tooltip></el-table-column>
        <el-table-column prop="size" label="大小" width="120" sortable>
          <template slot-scope="{row}">
            {{ row.isDir ? '-' : formatSize(row.size) }}
          </template>
        </el-table-column>
        <el-table-column prop="modified" label="修改时间" width="200" sortable>
          <template slot-scope="{row}">
            {{ formatDate(row.modified) }}
          </template>
        </el-table-column>
        <el-table-column align="center" label="操作" width="150" fixed="right">
          <template slot-scope="{row}">
            <el-button
              v-if="!row.isDir"
              size="mini"
              type="primary"
              @click="handleDownload(row)"
              circle
              icon="el-icon-download"
            ></el-button>
            <el-button
              v-else
              size="mini"
              type="primary"
              @click="handleDownloadDir(row)"
              circle
              icon="el-icon-download"
            ></el-button>
            <el-button
              size="mini"
              circle
              icon="el-icon-edit"
              @click="startDeepSearchRename(row)"
            ></el-button>
            <el-button
              size="mini"
              type="danger"
              @click="handleDeepSearchDelete(row)"
              circle
              icon="el-icon-delete"
            ></el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 无结果提示 -->
      <el-empty
        v-if="!isSearching && deepSearchResults.length === 0 && deepSearchHasSearched"
        description="未找到匹配的文件或目录"
        :image-size="80"
      ></el-empty>

      <!-- 搜索结果统计 -->
      <div v-if="!isSearching && deepSearchResults.length > 0" style="margin-top:10px;color:#909399;font-size:13px;">
        共找到 {{ deepSearchResults.length }} 个结果
      </div>
    </el-dialog>

    <!-- 搜索删除确认 -->
    <el-dialog
      title="确认删除"
      :visible.sync="deepSearchDeleteDialogVisible"
      width="30%"
      append-to-body
      :close-on-click-modal="false"
    >
      <span>确定要删除 {{ deepSearchDeleteTarget.name }} 吗？</span>
      <span slot="footer">
        <el-button @click="deepSearchDeleteDialogVisible = false">取消</el-button>
        <el-button type="danger" @click="confirmDeepSearchDelete">确定</el-button>
      </span>
    </el-dialog>

    <!-- 批量删除确认 -->
    <el-dialog
      title="批量删除"
      :visible.sync="batchDeleteDialogVisible"
      width="50%"
      append-to-body
      :close-on-click-modal="false"
    >
      <el-table
        :data="selectedFiles"
        style="width: 100%"
        max-height="350"
        border
        show-overflow-tooltip
      >
        <el-table-column prop="name" label="名称" width>
          <template slot-scope="{row}">
            <div :class="{'dir-item': row.isDir, 'file-item': !row.isDir}">
                <i :class="row.isDir ? 'el-icon-folder' : 'el-icon-document'"></i>
                {{ row.name }}
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="size" label="大小" width="120">
          <template slot-scope="{row}">
            {{ row.isDir ? '-' : formatSize(row.size) }}
          </template>
        </el-table-column>
        <el-table-column prop="modified" label="修改时间" width="200">
          <template slot-scope="{row}">
            {{ formatDate(row.modified) }}
          </template>
        </el-table-column>
      </el-table>
      <span slot="footer">
        <el-button @click="batchDeleteDialogVisible = false">取消</el-button>
        <el-button type="danger" @click="confirmBatchDelete">确定</el-button>
      </span>
    </el-dialog>
  </el-dialog>
</template>

<script>
export default {
  name: 'SftpBrowser',
  props: {
    // 用户名
    username: {
      type: String,
      default: ''
    },
    // 控制显示/隐藏
    visible: {
      type: Boolean,
      default: false
    },
    // 上传请求头（必须传 token）
    uploadHeaders: {
      type: Object,
      required: true
    },
    // 上传接口
    uploadUrl: {
      type: String,
      default: '/dev-api/sftp/upload'
    },
    // 路径
    path: {
      type: String,
      default: ''
    },
  },
  data() {
    return {
      innerVisible: false, // 内部显示状态
      showUploadProgress: false, // 上传进度条显示状态
      uploadPercent: 0, // 上传进度
      // 自定义上传进度条颜色
      customColors: [
        { color: '#f56c6c', percentage: 20 },
        { color: '#e6a23c', percentage: 40 },
        { color: '#5cb87a', percentage: 60 },
        { color: '#1989fa', percentage: 80 },
        { color: '#6f7ad3', percentage: 100 }
      ],
      currentPath: '/', // 当前路径
      fileList: [], // 文件列表
      breadcrumb: [{ name: '根目录', path: '/' }], // 面包屑
      isLoading: false, 
      showCreateFolderDialog: false,
      newFolderName: '',
      deleteDialogVisible: false,
      deleteTarget: { name: '', path: '', isDir: false },
      renamingItem: null,
      // 选中的文件
      selectedFiles: [],
      batchDeleteDialogVisible: false,
      // 拖拽上传相关
      dragOverlayVisible: false,      // 遮罩层显隐
      isUploading: false,              // 防止重复上传
      // 描述
      description: '', // 当前目录描述
      selectDescription: '', // 选中的文件描述
      // 搜索功能相关
      searchQuery: '',           // 搜索关键词
      showSearchBox: false,      // 搜索框显隐控制
      filteredFileList: [],      // 过滤后的文件列表
      originalFileList: [],      // 原始文件列表（用于还原）
      // 递归搜索相关
      showDeepSearchDialog: false,       // 搜索弹框显隐
      deepSearchPath: '',                // 搜索路径
      deepSearchKeyword: '',             // 搜索关键字
      deepSearchResults: [],             // 搜索结果
      isSearching: false,                // 搜索中状态
      deepSearchHasSearched: false,      // 是否已执行过搜索
      deepSearchDeleteDialogVisible: false, // 搜索删除确认弹框
      deepSearchDeleteTarget: { name: '', path: '', isDir: false }, // 搜索中待删除目标
    }
  },
  mounted() {
    // 全局监听 Backspace（8）
    window.addEventListener('keydown', this.handleKey)
    // 监听拖拽进入，显示遮罩层
    window.addEventListener('dragenter', this.onWindowDragEnter)
    // 拖拽结束（释放或取消）隐藏遮罩层
    window.addEventListener('dragend', this.onDragEnd)
  },
  beforeDestroy() {
    window.removeEventListener('keydown', this.handleKey)
    window.removeEventListener('dragenter', this.onWindowDragEnter)
    window.removeEventListener('dragend', this.onDragEnd)
  },
  watch: {
    // sftp浏览器显示时，获取文件列表
    visible(val) {
      this.innerVisible = val
      if (val && this.path === '') {
        this.fetchFiles()
      }else if (val && this.path !== '') {
        this.fetchFiles(this.path)
      } else {
        // 关闭时重置拖拽状态
        this.dragOverlayVisible = false
      }
    }
  },
  methods: {
    // 获取文件列表
    async fetchFiles(path = this.currentPath) {
      this.isLoading = true
      try {
        const res = await this.$API.sftpuser.reqSftpFiles({ path })
        if (res.code === 200) {
          this.fileList = res.data.files
          // 保存原始数据用于搜索还原
          this.originalFileList = JSON.parse(JSON.stringify(res.data.files))
          this.currentPath = res.data.path
          this.description = res.data.description
          this.updateBreadcrumb()
          // 如果已有搜索内容，重新应用过滤
          if (this.searchQuery) {
            this.applyFilter()
          } else {
            this.filteredFileList = this.fileList
          }
        }
      } catch (e) {
        this.$message.error('获取文件列表失败')
      } finally {
        this.isLoading = false
      }
    },
    // 更新面包屑
    updateBreadcrumb() {
      if (this.currentPath === '/') {
        this.breadcrumb = [{ name: '根目录', path: '/' }]
        return
      }
      const paths = this.currentPath.split('/').filter(Boolean)
      this.breadcrumb = paths.reduce((acc, cur) => {
        const last = acc[acc.length - 1]
        const path = last.path === '/' ? `/${cur}` : `${last.path}/${cur}`
        return [...acc, { name: cur, path }]
      }, [{ name: '根目录', path: '/' }])
    },
    // 点击目录/文件
    handleItemClick(item) {
      if (item.isDir) this.fetchFiles(item.path)
      this.selectedFiles = []
    },
    // 面包屑跳转
    handleBreadcrumbClick(item, index) {
      if (index === this.breadcrumb.length - 1) return
      if ( this.path != '' && item.path === '/') return
      this.fetchFiles(item.path)
    },
    // 返回上一级
    goBack() {
      if (this.currentPath === '/') return
      if (this.currentPath === this.path) return
      const paths = this.currentPath.split('/').filter(Boolean)
      paths.pop()
      this.fetchFiles(paths.length ? `/${paths.join('/')}` : '/')
    },
    // 上传前检查
    beforeUpload(file) {
      // 上传文件大小不能超过 5GB
      const isValid = file.size / 1024 / 1024 < 1024 * 5
      if (!isValid) this.$message.error('不能超过 5GB')
      return isValid
    },
    // 上传成功
    handleUploadSuccess() {
      this.$message.success('上传成功')
      this.fetchFiles()
      this.resetUploadProgress()
    },
    // 上传失败
    handleUploadError() {
      this.$message.error('上传失败')
      this.resetUploadProgress()
    },
    // 上传进度
    handleUploadProgress(e) {
      this.uploadPercent = Math.round(e.percent)
      this.showUploadProgress = true
    },
    // 重置进度
    resetUploadProgress() {
      this.showUploadProgress = false
      this.uploadPercent = 0
    },
    // 创建文件夹
    async createFolder() {
      if (!this.newFolderName) return this.$message.warning('请输入名称')
      try {
        const res = await this.$API.sftpuser.reqSftpMkdir({
          path: this.currentPath,
          name: this.newFolderName
        })
        if (res.code === 200) {
          this.$message.success('创建成功')
          this.showCreateFolderDialog = false
          this.newFolderName = ''
          this.fetchFiles()
        }
      } catch {}
    },
    // 下载文件
    handleDownload(file) {
      const token = sessionStorage.getItem('sftp_token')
      const url = `/dev-api/sftp/download?sftp_token=${token}&path=${encodeURIComponent(file.path)}`
      const a = document.createElement('a')
      a.href = url
      a.download = file.name
      a.click()
    },
    // 下载目录
    handleDownloadDir(file) {
      const token = sessionStorage.getItem('sftp_token')
      const url = `/dev-api/sftp/downloaddir?sftp_token=${token}&path=${encodeURIComponent(file.path)}`
      const a = document.createElement('a')
      a.href = url
      a.download = file.name
      a.click()
    },
    // 删除
    handleDelete(file) {
      this.deleteTarget = file
      this.deleteDialogVisible = true
    },
    async confirmDelete() {
      try {
        const res = await this.$API.sftpuser.reqSftpDelete({ path: this.deleteTarget.path })
        if (res.code === 200) {
          this.$message.success('删除成功')
          this.fetchFiles()
        }
      } catch {} finally {
        this.deleteDialogVisible = false
      }
    },
    // 重命名
    startRename(row) {
      this.$set(row, 'isRenaming', true)
      this.$set(row, 'editName', row.name)
    },
    async confirmRename(row) {
      if (!row.editName) return this.$message.warning('名称不能为空')
      try {
        const res = await this.$API.sftpuser.reqSftpRename({
          oldPath: row.path,
          newName: row.editName
        })
        if (res.code === 200) {
          this.$message.success('重命名成功')
          this.fetchFiles()
        }
      } catch {}
    },
    cancelRename(row) {
      row.isRenaming = false
    },
    // 关闭时向外通知
    handleClose() {
      this.breadcrumb = [{ name: '根目录', path: '/' }];
      this.currentPath = '/';
      this.fileList = [];
      this.isLoading = false;
      // 重置拖拽状态
      this.dragOverlayVisible = false;
      // 重置搜索状态
      this.searchQuery = ''
      this.showSearchBox = false
      this.filteredFileList = []
      this.originalFileList = []
      // 重置递归搜索状态
      this.showDeepSearchDialog = false
      this.deepSearchPath = ''
      this.deepSearchKeyword = ''
      this.deepSearchResults = []
      this.isSearching = false
      this.deepSearchHasSearched = false
      // 通知父组件关闭
      this.$emit('close')
    },
    // 工具函数：格式化文件大小
    formatSize(bytes) {
      if (bytes === 0) return '0 B';
      const k = 1024;
      const sizes = ['B', 'KB', 'MB', 'GB'];
      let i = 0;
      let size = bytes;
      while (size >= k && i < sizes.length - 1) {
        size /= k;
        i++;
      }
      return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    },
    // 工具函数：格式化日期
    formatDate(dateString) {
      return new Date(dateString).toLocaleString();
    },
    // 处理选择变化
    handleSelectionChange(selection) {
      this.selectedFiles = selection
      let dirCount = 0
      let fileCount = 0
      selection.forEach((file) => {
        if (file.isDir) {
          dirCount++
        } else {
          fileCount++
        }
      })
      if (dirCount == 0 &&  fileCount > 0) {
        this.selectDescription = `${fileCount} 个文件`
      } else if (fileCount == 0 && dirCount > 0) {
        this.selectDescription = `${dirCount} 个目录`
      } else {
        this.selectDescription = ` ${fileCount} 个文件 和 ${dirCount} 个目录`
      }
    },
    // 应用搜索过滤
    applyFilter() {
      if (!this.searchQuery.trim()) {
        this.filteredFileList = this.originalFileList
        return
      }
      const query = this.searchQuery.toLowerCase()
      this.filteredFileList = this.originalFileList.filter(item => 
        item.name.toLowerCase().includes(query)
      )
    },
    // 切换搜索框显示
    toggleSearchBox() {
      this.showSearchBox = !this.showSearchBox
      if (this.showSearchBox) {
        // 打开时聚焦到搜索输入框
        this.$nextTick(() => {
          if (this.$refs.searchInput) {
            this.$refs.searchInput.focus()
          }
        })
      } else {
        // 关闭时重置
        this.resetSearch()
      }
    },
    // 重置搜索
    resetSearch() {
      this.searchQuery = ''
      this.showSearchBox = false
      this.filteredFileList = this.originalFileList
    },
    // 关闭搜索框（通过按钮或 esc）
    handleCloseSearchBox() {
      this.showSearchBox = false
      this.searchQuery = ''
      this.filteredFileList = this.originalFileList
      // 关闭后聚焦到容器避免焦点丢失
      this.$nextTick(() => {
        document.activeElement.blur()
      })
    },
    // 批量删除
    async confirmBatchDelete(){
      try {
        const res = await this.$API.sftpuser.reqSftpBatchDelete(this.selectedFiles)
        if (res.code === 200) {
          this.$message.success('删除成功')
          this.fetchFiles()
          this.selectedFiles = []
        }
      } catch (error) {} finally {
        this.batchDeleteDialogVisible = false
      }
    },
    // 批量下载
    handleBatchDownload(){
      this.selectedFiles.forEach((file) => {
        if (file.isDir) {
          this.handleDownloadDir(file)
        } else {
          this.handleDownload(file)
        }
      })
    },
    handleKey(e) {
      // 只有在组件可见时才处理按键事件
      if (!this.innerVisible) return;
      
      // Ctrl/Cmd + F 切换搜索框
      if ((e.ctrlKey || e.metaKey) && (e.key === 'f' || e.keyCode === 70)) {
        e.preventDefault()
        e.stopPropagation()
        this.toggleSearchBox()
        return
      }
      
      // 判断是否为真正的文本输入元素（排除复选框、单选按钮、按钮等）
      const tag = e.target.tagName;
      const isTextInput = 
        (tag === 'INPUT' && (e.target.type === 'text' || e.target.type === 'password' || e.target.type === 'search' || e.target.type === 'email' || e.target.type === 'tel' || e.target.type === 'url' || e.target.type === 'number')) ||
        tag === 'TEXTAREA' || (e.target.isContentEditable === true);

      if (isTextInput) return;
      
      // Esc 关闭搜索框
      if (e.key === 'Escape' || e.keyCode === 27) {
        if (this.showSearchBox) {
          e.preventDefault()
          this.handleCloseSearchBox()
          return
        }
      }
      
      // F5 刷新
      if (e.key === 'F5'){
        e.preventDefault() // 阻止默认刷新行为
        e.stopPropagation()  // 阻止事件冒泡
        if (this.innerVisible) {    // 仅当弹窗可见时才刷新列表
          this.fetchFiles()
          this.$message.info('已刷新')
        }
        return
      }
      
      // Backspace 返回上一级
      if (e.keyCode === 8 || e.key === 'Backspace') {
        this.goBack();
        e.preventDefault(); // 阻止浏览器后退
      }
      // Delete 键触发批量删除
      else if (e.keyCode === 46 || e.key === 'Delete') {
        // 只有选中时才弹出对话框
        if (this.selectedFiles.length > 0) {
          this.openBatchDeleteDialog();
          e.preventDefault(); // 阻止可能的默认行为
        }
      }
      
      // Enter 确认批量删除（对话框打开时）
      else if (e.key === 'Enter' || e.keyCode === 13) {
        if (this.batchDeleteDialogVisible) {
          this.confirmBatchDelete();
          e.preventDefault();
        }
        else if (this.deleteDialogVisible){
          this.confirmDelete()
          e.preventDefault()
        }
      }
      // Ctrl+A 全选
      else if (e.ctrlKey && e.key === 'a') {
        this.$refs.multipleTable.toggleAllSelection();
        e.preventDefault();
      }
    },

    // 拖拽上传方法
    // 全局 dragenter 处理: 只有当拖拽进入组件内部且有文件时显示遮罩层
    onWindowDragEnter(e) {
      // 只有在组件可见时才处理拖拽事件
      if (!this.innerVisible) return;
      // 检查拖拽内容是否包含文件
      if (!e.dataTransfer || !e.dataTransfer.types) return
      const hasFiles = Array.from(e.dataTransfer.types).includes('Files')
      if (!hasFiles) return

      // 判断拖拽目标是否在当前组件容器内
      const container = this.$refs.tableContainer
      if (!container || !container.contains(e.target)) return

      // 阻止默认避免异常
      e.preventDefault()
      // 显示遮罩层
      this.dragOverlayVisible = true
    },

    // 全局 dragend: 拖拽结束隐藏遮罩层
    onDragEnd() {
      // 只有在组件可见时才处理拖拽事件
      if (!this.innerVisible) return;
      this.dragOverlayVisible = false
    },

    // 保持遮罩层在自身区域不隐藏（防止闪烁，也可不处理，全局dragend会确保最终隐藏）
    keepOverlayVisible(e) {
      // 什么都不做，只是为了覆盖默认行为并维持显示
      e.preventDefault()
    },

    // 处理文件释放（拖拽上传）
    async handleDrop(e) {
      // 立即隐藏遮罩层，避免重复触发
      this.dragOverlayVisible = false
      
      const files = Array.from(e.dataTransfer.files)
      if (files.length === 0) return

      if (this.isUploading) {
        this.$message.warning('正在上传中，请稍后再试')
        return
      }

      // 文件大小预检（不能超过5GB）
      const invalidFile = files.find(f => f.size / 1024 / 1024 >= 1024 * 5)
      if (invalidFile) {
        this.$message.error(`文件 ${invalidFile.name} 超过 5GB，上传终止`)
        return
      }

      // 开始上传
      this.isUploading = true
      this.showUploadProgress = true
      let successCount = 0
      let failCount = 0

      // 串行上传，并更新整体进度
      for (let i = 0; i < files.length; i++) {
        const file = files[i]
        // 更新整体百分比（基于完成数量）
        this.uploadPercent = Math.floor((i / files.length) * 100)
        try {
          await this.uploadSingleFile(file)
          successCount++
        } catch (err) {
          console.error(err)
          failCount++
        }
      }
      this.uploadPercent = 100
      this.isUploading = false
      // 延迟重置进度条，避免一闪而过
      setTimeout(() => {
        this.resetUploadProgress()
      }, 1000)
      
      if (successCount > 0) {
        this.$message.success(`成功上传 ${successCount} 个文件${failCount ? `，失败 ${failCount} 个` : ''}`)
        this.fetchFiles() // 刷新列表
      } else if (failCount > 0) {
        this.$message.error('所有文件上传失败')
      }
    },

    // 单文件上传（使用与 el-upload 相同的接口和参数）
    uploadSingleFile(file) {
      return new Promise((resolve, reject) => {
        const formData = new FormData()
        formData.append('file', file)
        formData.append('path', this.currentPath)

        const xhr = new XMLHttpRequest()
        xhr.open('POST', this.uploadUrl, true)
        // 设置请求头
        Object.keys(this.uploadHeaders).forEach(key => {
          xhr.setRequestHeader(key, this.uploadHeaders[key])
        })
        xhr.onload = () => {
          if (xhr.status === 200) {
            try {
              const res = JSON.parse(xhr.responseText)
              if (res.code === 200) {
                resolve()
              } else {
                reject(new Error(res.message || '上传失败'))
              }
            } catch (e) {
              reject(e)
            }
          } else {
            reject(new Error(`HTTP ${xhr.status}`))
          }
        }
        xhr.onerror = () => {
          reject(new Error('网络错误'))
        }
        xhr.send(formData)
      })
    },
    onTableDragLeave(e){
      // 如果 relatedTarget 仍然在容器内，说明只是离开了某个子元素，不隐藏
      if (this.$refs.tableContainer && this.$refs.tableContainer.contains(e.relatedTarget)) {
        return
      }
      this.dragOverlayVisible = false 
    },
    // 打开创建目录对话框
    handleCreateDir(){
      this.showCreateFolderDialog = true
      this.$nextTick(() => {
        this.$refs.newDirInput.focus()
      })
    },
    // 打开递归搜索弹框
    openDeepSearch() {
      this.deepSearchPath = this.currentPath
      this.deepSearchKeyword = ''
      this.deepSearchResults = []
      this.deepSearchHasSearched = false
      this.showDeepSearchDialog = true
    },
    // 执行递归搜索
    async doDeepSearch() {
      if (!this.deepSearchKeyword.trim()) {
        return this.$message.warning('请输入搜索关键字')
      }
      if (!this.deepSearchPath.trim()) {
        return this.$message.warning('请输入搜索路径')
      }
      this.isSearching = true
      this.deepSearchHasSearched = true
      try {
        const res = await this.$API.sftpuser.reqSftpSearch({
          path: this.deepSearchPath,
          keyword: this.deepSearchKeyword
        })
        if (res.code === 200) {
          this.deepSearchResults = res.data.results || []
        }
      } catch (e) {
        this.$message.error('搜索失败')
      } finally {
        this.isSearching = false
      }
    },
    // 搜索结果中开始重命名
    startDeepSearchRename(row) {
      this.$set(row, 'isRenaming', true)
      this.$set(row, 'editName', row.name)
    },
    // 搜索结果中确认重命名
    async confirmDeepSearchRename(row) {
      if (!row.editName) return this.$message.warning('名称不能为空')
      try {
        const res = await this.$API.sftpuser.reqSftpRename({
          oldPath: row.path,
          newName: row.editName
        })
        if (res.code === 200) {
          this.$message.success('重命名成功')
          this.doDeepSearch()
        }
      } catch {} 
    },
    // 搜索结果中取消重命名
    cancelDeepSearchRename(row) {
      row.isRenaming = false
    },
    // 搜索结果中删除
    handleDeepSearchDelete(row) {
      this.deepSearchDeleteTarget = row
      this.deepSearchDeleteDialogVisible = true
    },
    // 确认搜索结果中的删除
    async confirmDeepSearchDelete() {
      try {
        const res = await this.$API.sftpuser.reqSftpDelete({ path: this.deepSearchDeleteTarget.path })
        if (res.code === 200) {
          this.$message.success('删除成功')
          this.doDeepSearch()
        }
      } catch {} finally {
        this.deepSearchDeleteDialogVisible = false
      }
    },
    // 打开批量删除确认对话框
    openBatchDeleteDialog() {
      if (this.selectedFiles.length === 0) {
        return;
      }
      this.batchDeleteDialogVisible = true;
    },
  }
}
</script>

<style scoped>
.breadcrumb { cursor: pointer; }
.breadcrumbBold { font-weight: bold; color: #409EFF; }
.breadcrumb:hover { color: #409EFF; font-weight: bold; }
.goBack:hover { color: #409EFF; font-weight: bold; }
.dir-item { color: #409EFF; cursor: pointer; }
.file-item { color: #606266; }
.operate { display: flex; justify-content: space-between; align-items: center; }
.upload { display: inline-block; }

/* 表格容器相对定位，用于拖拽遮罩层绝对定位 */
.table-container {
  position: relative;
}

/* 拖拽上传遮罩层样式：覆盖整个表格区域，半透明背景，中心提示 */
.drag-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: #64c89631;
  backdrop-filter: blur(2px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10;
  border-radius: 4px;
  pointer-events: auto;  /* 仅显示时接收拖拽事件，不影响表格本身点击（因隐藏时不可见） */
}

.drag-overlay-content {
  text-align: center;
  color: white;
  font-size: 24px;
  background: rgba(0, 0, 0, 0.6);
  padding: 20px 40px;
  border-radius: 12px;
  pointer-events: none;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.drag-overlay-content i {
  font-size: 48px;
}

.drag-sub {
  font-size: 14px;
  opacity: 0.9;
}

/* 搜索框样式 */
.search-box-container {
  padding: 12px 15px;
  background-color: #f5f7fa;
  border-top: 1px solid #ebeef5;
  border-bottom: 1px solid #ebeef5;
}

.search-box {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
}

.search-box .el-input {
  flex: 1;
}

/* 骨架屏表格占位 */
.skeleton-table-wrapper {
  border: 1px solid #ebeef5;
  border-radius: 4px;
  padding: 16px;
  min-height: 200px;
  background-color: #fff;
}
</style>