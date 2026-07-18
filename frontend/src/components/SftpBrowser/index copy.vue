<template>
  <el-dialog
    title="SFTP浏览器"
    :visible.sync="innerVisible"
    center
    width="70%"
    top="10vh"
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

        <div style="width: 210px;">
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
            @click="showCreateFolderDialog = true"
          >创建目录</el-button>
        </div>
      </div>
      <!-- 文件列表 -->
      <el-card shadow="hover">
        <el-empty
          description="这个目录很穷，什么都没有_(:3 ∠)_"
          v-if="!fileList"
          style="height: 500px;"
        ></el-empty>

        <el-table
          ref="multipleTable"
          @selection-change="handleSelectionChange"
          :data="fileList"
          v-loading="isLoading"
          height="calc(100vh - 400px)"
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
        <el-collapse-transition>
          <div v-show="selectedFiles.length > 0" style="margin-top: 10px; display: flex; justify-content: flex-end;">
            <el-button type="primary" size="mini" icon="el-icon-download" round @click="handleBatchDownload">批量下载</el-button>
            <el-button type="danger" size="mini" icon="el-icon-delete" round @click="batchDeleteDialogVisible = true">批量删除</el-button>
          </div>
        </el-collapse-transition>
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
      <el-form>
        <el-form-item label="文件夹名称">
          <el-input v-model="newFolderName" autocomplete="off"></el-input>
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
    }
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
    }
  },
  mounted() {
    // 全局监听 Backspace（8）
    window.addEventListener('keydown', this.handleKey)
  },
  beforeDestroy() {
    window.removeEventListener('keydown', this.handleKey)
  },
  watch: {
    // sftp浏览器显示时，获取文件列表
    visible(val) {
      this.innerVisible = val
      if (val) {
        this.fetchFiles()
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
          // console.log(this.fileList)
          this.currentPath = res.data.path
          this.updateBreadcrumb()
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
      this.fetchFiles(item.path)
    },
    // 返回上一级
    goBack() {
      if (this.currentPath === '/') return
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
      const token = localStorage.getItem('sftp_token')
      const url = `/dev-api/sftp/download?sftp_token=${token}&path=${encodeURIComponent(file.path)}`
      const a = document.createElement('a')
      a.href = url
      a.download = file.name
      a.click()
    },
    // 下载目录
    handleDownloadDir(file) {
      const token = localStorage.getItem('sftp_token')
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
      // 避免在输入框里按删除也返回
      const tag = e.target.tagName
      if (tag === 'INPUT' || tag === 'TEXTAREA') return

      if (e.keyCode === 8 ) {
        this.goBack()
      }
    }
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
.upload { display: inline-block; margin-right: 10px; }
</style>