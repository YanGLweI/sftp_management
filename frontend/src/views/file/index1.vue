<template>
  <div class="page-container">
    <!-- 背景图片容器 -->
    <div class="bg-container"></div>
    <!-- 左上角logo -->
    <div class="logo-container">
      <img :src="logoPath" alt="系统logo" class="logo" @click="goToLoginPage" />
    </div>
    <!-- sftp登录-->
    <el-card class="sftp-login-fade" style="width: 480px;">
      <el-card>
        <div slot="header">
          <span>SFTP登录</span>
        </div>
        <el-tabs v-model="activeName" @tab-click="handleClick">
          <el-tab-pane label="密码登录" name="password">
            <el-form :model="SftpForm">
              <el-form-item label="SFTP账号" :label-width="SftpFormLabelWidth">
                <el-input v-model="SftpForm.username" autocomplete="off" ref="username"></el-input>
              </el-form-item>
              <el-form-item label="SFTP密码" :label-width="SftpFormLabelWidth">
                <el-input
                  v-model="SftpForm.password"
                  autocomplete="off"
                  type="password"
                  show-password
                  v-focus="SftpDialogFormVisible"
                  @keyup.enter.native="sftplogin()"
                ></el-input>
              </el-form-item>
            </el-form>
          </el-tab-pane>
          <el-tab-pane label="密钥登录" name="keyfile">
            <el-form :model="SftpForm">
              <el-form-item label="SFTP账号" :label-width="SftpFormLabelWidth">
                <el-input v-model="SftpForm.username" autocomplete="off" ref="keyUsername"></el-input>
              </el-form-item>
              <el-form-item label="SFTP密钥" :label-width="SftpFormLabelWidth">
                <el-upload
                  class="upload-key"
                  ref="upload"
                  :action="KeyFileUploadUrl"
                  :data="{ username: SftpForm.username }"
                  :headers="uploadHeaders"
                  :file-list="keyFileList"
                  :limit="1"
                  :before-upload="beforeUploadKey"
                  :on-change="handleKeyFileChange"
                  :on-success="KeyhandleUploadSuccess"
                  :on-error="KeyhandleUploadError"
                  :auto-upload="false"
                >
                  <el-button slot="trigger" size="small" type="primary">选择密钥文件</el-button>
                </el-upload>
              </el-form-item>
            </el-form>
          </el-tab-pane>
        </el-tabs>
      </el-card>
      <div class="dialog-footer">
        <el-button @click="closeSftpDialogForm()">取 消</el-button>
        <el-button
          type="primary"
          @click="sftplogin()"
          :loading="buttonLoading"
        >{{ buttonLoading ? '提交中 ...' : '确 定' }}</el-button>
      </div>
    </el-card>

    <!-- SFTP浏览器 -->
    <el-dialog
      title="SFTP浏览器"
      :visible.sync="SftpBrowserVisible"
      center
      width="width"
      top="10vh"
      :close-on-click-modal="false"
      @close="closeSftpBrowser()"
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
          <!-- 返回上一级 -->
          <el-page-header
            @back="goBack"
            :content="breadcrumb[breadcrumb.length - 1].name"
            style="margin:20px 0;width: 200px;"
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
            <!-- 上传按钮 -->
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
            <!-- 创建目录按钮 -->
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
            v-if="fileList == null ? true : false"
            style="height: 500px;"
          ></el-empty>
          <el-table
            :data="fileList"
            v-loading="isLoading"
            height="500"
            border
            v-if="fileList !== null ? true : false"
          >
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
        </el-card>
      </div>
      <!-- 新建文件夹对话框 -->
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
      <!-- 删除确认弹框 -->
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
    </el-dialog>
  </div>
</template>

<script>
export default {
  name: 'File',
  data() {
    // 密码验证器
    var validatePass = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请输入密码'));
      } else if (value.length < 14) {
        callback(new Error('密码长度至少为14位'));
      } else {
        // 检查是否包含大小写字母和符号
        let hasLowercase = /[a-z]/.test(value);
        let hasUppercase = /[A-Z]/.test(value);
        let hasSymbol = /[^\w]/.test(value);
        let hasNumber = /\d/.test(value);
        if (!hasLowercase || !hasUppercase || !hasSymbol || !hasNumber) {
          callback(new Error('密码必须包含大小写字母数字和符号'));
        } else {
          if (this.userForm.checkPass !== '') {
            this.$refs.ruleForm.validateField('checkPass');
          }
          callback();
        }
      }
    };
    // 确认密码验证器
    var validatePass2 = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请再次输入密码'));
      } else if (value !== this.userForm.password) {
        callback(new Error('两次输入密码不一致!'));
      } else {
        callback();
      }
    };
    return {
      logoPath: require('@/assets/logo.png'),
      loading: false, //是否显示加载效果
      userlist: [], //用户列表
      dialogFormVisible: false,
      title: '', //dialog标题
      dialogTableVisible: false,
      emailDialog: false,
      buttonLoading: false,
      showPassword: false,
      userForm: {
        id: '', // 用户id
        name: '', // 账号名
        loginType: 'Password', // 登录类型
        noExpire: false, // 永不过期
        emailOrNot: false,  //是否发送邮件
        password: '',   // 密码
        checkPass: '', // 确认密码
        to: '',   //收件人
        cc: '',  //抄送人
      },
      // 邮箱相关
      options: [],
      // SFTP登录
      SftpDialogFormVisible: false,
      SftpFormLabelWidth: '100px',
      SftpForm: {
        username: '',
        password: ''
      },
      // 整理的email数据
      email: {
        to: [],
        cc: [],
        subject: '',
      },
      searchObj: {
        username: ''
      },
      tempSearchObj: {
        username: ''
      },
      ids: [], // 收集多选的用户id列表
      multipleSelection: [], //多项的用户信息
      // 分页器
      page: 1,
      limit: 5,
      total: 0,
      // 验证规则
      rules: {
        name: [
          { required: true, message: '请输入账号名称', trigger: 'blur' },
          {
            pattern: /^[a-zA-Z][a-zA-Z0-9_]*$|^[0-9](?=.*[a-zA-Z])[a-zA-Z0-9_]*$/,
            message: '账号须以字母/数字开头（必须包含字母）仅允许字母、数字和下划线组成',
            trigger: ['blur', 'change']
          }
        ],
        loginType: [{ required: true, message: '请选择登录类型', trigger: 'blur' }],
        password: [
          { required: true, message: '请输入密码', trigger: 'blur' },
          { validator: validatePass, trigger: 'blur' }
        ],
        checkPass: [
          { required: true, message: '请确认密码', trigger: 'blur' },
          { validator: validatePass2, trigger: ['blur', 'change'] }
        ],
        to: [{ required: true, message: '请输入邮箱地址', trigger: 'blur' }],
        cc: [{ required: false, message: '请输入邮箱地址', trigger: 'blur' }],
      },
      // sftp浏览器相关数据
      SftpBrowserVisible: false,
      activeName: 'password', // SFTP登录方式
      // 上传进度相关
      showUploadProgress: false, // 是否显示上传进度
      uploadPercent: 0, // 上传进度百分比
      customColors: [
          {color: '#f56c6c', percentage: 20},
          {color: '#e6a23c', percentage: 40},
          {color: '#5cb87a', percentage: 60},
          {color: '#1989fa', percentage: 80},
          {color: '#6f7ad3', percentage: 100}
        ],
      // 密钥登录
      keyFileList: [], // 密钥文件列表
      KeyFileUploadUrl: '/dev-api/sftp/login', // 密钥登录接口
      currentPath: '/', // 当前路径
      fileList: [], // 当前目录下的文件、目录列表
      breadcrumb: [{
        name: '根目录',
        path: '/'
      }],
      isLoading: false,
      uploadUrl: '/dev-api/sftp/upload', // 上传文件的接口地址
      // 添加认证头
      uploadHeaders: {
        'Token': `${this.$store.state.user.token}`,
        'X-SFTP-Token': ''
      },
      showCreateFolderDialog: false, // 是否显示新建文件夹对话框
      newFolderName: '', // 新建文件夹的名称
      deleteDialogVisible: false,
      deleteTarget: {
        name: '',
        path: '',
        isDir: false
      },
      renamingItem: null, //正在重命名的文件或目录
    }
  },
  mounted() {
    // this.getUserList()
    this.$nextTick(() => {
      this.$refs.username.focus()
    })
  },
  methods: {
    goToLoginPage() {
      this.$router.push({ path: '/login' })
    },
    // 获取用户列表
    async getUserList(pages = 1) {
      this.page = pages
      const { page, limit, searchObj } = this
      try {
        this.loading = true
        const result = await this.$API.sftpuser.reqUserList(page, limit, searchObj)
        if (result.code == 200) {
          this.userlist = result.data.users
          this.total = result.data.total
        }
        this.loading = false
      } catch (error) { }
    },
    // 显示对话框,并清除数据
    showDialog(title, row) {
      this.title = title
      this.dialogFormVisible = true
      this.userForm = {
        id: '',
        name: '',
        loginType: 'Password',
        noExpire: false,
        emailOrNot: false,
        password: '',
        checkPass: '',
        to: '',
        cc: '',
      }
      this.$nextTick(() => {
        if (this.$refs.ruleForm) {
          this.$refs.ruleForm.clearValidate();
        }
      });
      if (title == '修改密码') {
        this.userForm.name = row.name
        this.userForm.id = row.id
      }
    },
    // 添加账号和修改密码的确定按钮的回调
    addOrUpdateUser() {
      this.buttonLoading = true
      this.$refs['ruleForm'].validate(async (valid) => {
        if (valid) {
          const { id, name, loginType, password, emailOrNot, noExpire, to, cc } = this.userForm
          let message;
          if (loginType === 'KeyFile') {
            message = '密钥修改成功';
          } else if (loginType === 'Password') {
            message = '密码修改成功';
          } else {
            message = '密码和密钥修改成功';
          }
          try {
            if (this.title == '修改密码') {
              const result = await this.$API.sftpuser.reqUpdateUser({ id, name, loginType, password, emailOrNot, noExpire })
              if (result.code == 200) {
                if (this.userForm.emailOrNot) {
                  this.email.to = to
                  this.email.cc = cc
                  this.email.subject = `${name}`
                  await this.$API.sftpuser.reqSendEmail(this.email)
                  this.$message({ type: 'success', message: `账号:${name} ${message},邮件已发送` })
                } else {
                  this.$message({ type: 'success', message: `账号:${name} ${message}` })
                }
                this.getUserList(this.page)
                this.dialogFormVisible = false
                this.buttonLoading = false
              }
            } else {
              const result = await this.$API.sftpuser.reqAddUser({ name, loginType, password, emailOrNot, noExpire })
              if (result.code == 200) {
                if (this.userForm.emailOrNot) {
                  this.email.to = to
                  this.email.cc = cc
                  this.email.subject = `${name}`
                  await this.$API.sftpuser.reqSendEmail(this.email)
                  this.$message({ type: 'success', message: `账号:${name} 添加成功,邮件已发送` })
                } else {
                  this.$message({ type: 'success', message: `账号:${name} 添加成功` })
                }
                this.getUserList(this.page)
                this.dialogFormVisible = false
                this.buttonLoading = false
              }
            }
          } catch (error) { this.buttonLoading = false }
        } else {
          console.log('error submit!!');
          this.buttonLoading = false
          return false;
        }
      })
    },
    // Dialog 关闭的回调
    close() {
      this.title = ''
      this.$refs['ruleForm'].clearValidate();
      this.$refs['ruleForm'].resetFields();
      this.dialogFormVisible = false;
    },
    closedialogTable() {
      this.dialogTableVisible = false
    },
    // 分页器limit改变的回调
    handleSizeChange(limit) {
      this.limit = limit;
      this.getUserList(this.page);
    },
    // 搜索的回调
    search() {
      this.searchObj = { ...this.tempSearchObj };
      const { username } = this.searchObj
      if (username == '') {
        this.$message({ type: 'warning', message: '请输入账号名' })
      } else {
        this.getUserList()
      }
    },
    resetSearch() {
      this.searchObj = {
        username: "",
      };
      this.tempSearchObj = {
        username: "",
      };
      this.getUserList();
    },
    // 删除一个用户
    async deleteUser(id) {
      try {
        const { page } = this
        const result = await this.$API.sftpuser.reqDeleteUser(id)
        if (result.code == 200) {
          this.getUserList(this.userlist.length == 1 ? page - 1 : page)
          this.$message({ type: "success", message: "删除成功" })
        }
      } catch (error) { }
    },
    // 存储选择的用户id
    handleSelectionChange(val) {
      this.ids = []
      this.multipleSelection = val
      val.forEach(item => {
        this.ids.push(item.id)
      });
    },
    showTableDialog() {
      this.dialogTableVisible = true
    },
    async deleteUsers() {
      try {
        const { ids, page } = this
        const result = await this.$API.sftpuser.reqDeleteUsers(ids)
        if (result.code == 200) {
          this.dialogTableVisible = false
          this.$message({ type: 'success', message: '批量删除成功' })
          this.getUserList(this.userlist.length == this.ids.length ? page - 1 : page)
        }
      } catch (error) { }
    },
    // 取消email
    cancelEmailForm() {
      this.buttonLoading = false;
      this.emailDialog = false;
    },
    // 生成密码的回调
    generatePassword() {
      const length = 14;
      const lowercaseChars = 'abcdefghjkmnpqrstuvwxyz';
      const uppercaseChars = 'ABCDEFGHJKLMNPQRSTUVWXYZ';
      const numberChars = '123456789';
      const symbolChars = '!@#$%^&*+~:;?><,.-=';
      let password = '';
      password += lowercaseChars[Math.floor(Math.random() * lowercaseChars.length)];
      password += uppercaseChars[Math.floor(Math.random() * uppercaseChars.length)];
      password += numberChars[Math.floor(Math.random() * numberChars.length)];
      password += symbolChars[Math.floor(Math.random() * symbolChars.length)];
      for (let i = 4; i < length; i++) {
        const charSet = lowercaseChars + uppercaseChars + numberChars + symbolChars;
        password += charSet[Math.floor(Math.random() * charSet.length)];
      }
      this.userForm.password = password;
      this.userForm.checkPass = password;
    },
    //  发送邮件变化的事件回调
    async change(value) {
      if (value) {
        try {
          const result = await this.$API.contact.reqContactOptions()
          if (result.code == 200) {
            this.options = result.data.options
          }
        } catch (error) {
        }
      }
      if (!value) {
        this.userForm.cc = ''
      }
    },
    // 显示SFTP登录对话框
    showSftpDialog(row) {
      this.SftpDialogFormVisible = true
      this.SftpForm.username = row.name
      this.SftpForm.password = ''
    },
    // 发送请求，登录sftp
    async sftplogin() {
      if (this.activeName == 'password') {
        const { username, password } = this.SftpForm

        if (password == '' || username == '') {
          this.$message({ type: 'warning', message: '请输入SFTP账号和密码' })
          return
        }
        this.buttonLoading = true
        try {
          const result = await this.$API.sftpuser.reqSftpLogin({ username, password })
          if (result.code == 200) {
            localStorage.setItem("sftp_token", result.data.sftp_token)
            this.uploadHeaders['X-SFTP-Token'] = result.data.sftp_token
            this.$message({ type: 'success', message: 'SFTP登录成功' })
            this.fetchFiles();
          }
        } catch (error) { } finally {
          this.SftpDialogFormVisible = false
          this.buttonLoading = false
        }
      } else {
        const { username } = this.SftpForm
        if (username == '') {
          this.$message({ type: 'warning', message: '请输入SFTP账号' })
          return
        }
        if (this.keyFileList.length == 0) {
          this.$message({ type: 'warning', message: '请选择密钥文件' })
          return
        }
        this.buttonLoading = true
        this.$refs.upload.submit();
      }
    },
    handleKeyFileChange(file, fileList) {
      this.keyFileList = fileList;
    },
    // 上传前检查
    beforeUploadKey(file) {
      const isLt100M = file.size / 1024 / 1024 < 100;
      if (!isLt100M) {
        this.$message.error('上传文件大小不能超过 100MB!');
      }
      return isLt100M;
    },
    // 密钥登录上传成功的回调
    KeyhandleUploadSuccess(response, file, fileList) {
      if (response.code === 200) {
        this.$message({ type: 'success', message: 'SFTP登录成功' });
        if (response.data && response.data.sftp_token) {
          localStorage.setItem("sftp_token", response.data.sftp_token);
          this.uploadHeaders['X-SFTP-Token'] = response.data.sftp_token
        }
        this.fetchFiles();
      } else {
        this.$message({ type: 'error', message: response.message || 'SFTP登录失败' });
      }
      this.SftpDialogFormVisible = false;
      this.buttonLoading = false;
      this.keyFileList = [];
      this.$refs.upload.clearFiles();
    },
    // 密钥登录上传失败的回调
    KeyhandleUploadError(response, file, fileList) {
      this.$message({ type: 'error', message: response.message || 'SFTP登录失败' });
      this.SftpDialogFormVisible = false;
      this.buttonLoading = false;
      this.keyFileList = [];
      this.$refs.upload.clearFiles();
    },
    // 关闭SFTP登录对话框
    closeSftpDialogForm() {
      this.SftpDialogFormVisible = false
      this.buttonLoading = false
      this.activeName = 'password'
      this.keyFileList = []
      this.$refs.upload.clearFiles()
      this.SftpForm = {
        username: '',
        password: '',
      }
    },
    // 发送请求获取文件和目录列表
    async fetchFiles(path = this.currentPath) {
      this.SftpBrowserVisible = true;
      this.isLoading = true;
      try {
        const result = await this.$API.sftpuser.reqSftpFiles({ path });
        if (result.code === 200) {
          this.fileList = result.data.files;
          this.currentPath = result.data.path;
          this.updateBreadcrumb();
        }
      } catch (e) {
        this.$message.error('获取文件列表失败');
      } finally {
        this.isLoading = false;
      }
    },
    // 更新面包屑导航
    updateBreadcrumb() {
      if (this.currentPath === '/') {
        this.breadcrumb = [{ name: '根目录', path: '/' }];
        return;
      }
      const paths = this.currentPath.split('/').filter(Boolean);
      this.breadcrumb = paths.reduce((acc, cur, index) => {
        const lastItem = acc[acc.length - 1];
        const path = lastItem.path == '/'
          ? lastItem.path + cur
          : lastItem.path + '/' + cur;
        return [...acc, { name: cur, path }];
      }, [{ name: '根目录', path: '/' }]);
    },
    // 点击文件/目录
    handleItemClick(item) {
      if (item.isDir) {
        this.fetchFiles(item.path);
      } else {
        this.isFile = true;
        this.currentPath = item.path;
      }
    },
    // 点击面包屑导航
    handleBreadcrumbClick(item, index) {
      if (index === this.breadcrumb.length - 1) return;
      this.fetchFiles(item.path);
    },
    // 返回上一级
    goBack() {
      if (this.currentPath === '/') return;
      const paths = this.currentPath.split('/').filter(Boolean);
      paths.pop();
      const newPath = paths.length ? `/${paths.join('/')}` : '/';
      this.fetchFiles(newPath);
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
    // 关闭 SFTP 浏览器时重置面包屑和文件列表，断开sftp连接
    async closeSftpBrowser() {
      this.breadcrumb = [{ name: '根目录', path: '/' }];
      this.currentPath = '/';
      this.fileList = [];
      this.isLoading = false;
      this.SftpForm = {
        username: '',
        password: '',
      }
      // 新增：重置上传进度条
      this.resetUploadProgress();
      try {
        const result = await this.$API.sftpuser.reqSftpLogout();
        if (result.code === 200) {
          this.$message({ type: 'success', message: result.message });
        } else {
          this.$message({ type: 'error', message: result.message });
        }
      } catch (error) { }
    },
    // 上传前检查
    beforeUpload(file) {
      const isLt1G = file.size / 1024 / 1024 < 1024;
      if (!isLt1G) {
        this.$message.error('上传文件大小不能超过 1GB!');
      }
      return isLt1G;
    },
    // 上传成功处理
    handleUploadSuccess(response, file, fileList) {
      if (response.code === 200) {
        this.$message.success('文件上传成功');
        this.fetchFiles();
      } else {
        this.$message.error(response.message || '文件上传失败');
      }
      // 重置进度条（单文件上传直接重置，多文件上传判断是否全部完成）
      if (fileList.every(item => item.status === 'success' || item.status === 'fail')) {
        this.resetUploadProgress();
      }
    },
    // 上传失败处理
    handleUploadError(err, file, fileList) {
      this.$message.error('文件上传失败: ' + err.message);
      // 重置进度条
      this.resetUploadProgress();
    },
    // 新增：重置上传进度条的公共方法
    resetUploadProgress() {
      this.showUploadProgress = false;
      this.uploadPercent = 0;
      this.isUploading = false;
    },
    // 创建文件夹
    async createFolder() {
      if (!this.newFolderName) {
        this.$message.error('请输入文件夹名称');
        return;
      }
      try {
        const result = await this.$API.sftpuser.reqSftpMkdir({
          path: this.currentPath,
          name: this.newFolderName
        })
        if (result.code === 200) {
          this.$message.success('文件夹创建成功');
          this.showCreateFolderDialog = false;
          this.newFolderName = '';
          this.fetchFiles();
        } else {
          this.$message.error(result.message || '文件夹创建失败');
        }
      } catch (err) { }
    },
    // 下载文件
    // async handleDownload(file) {
    //   if (file.isDir) {
    //     this.$message.warning('请选择文件进行下载')
    //     return
    //   }
    //   try {
    //     const response = await this.$API.sftpuser.reqSftpDownload({ path: file.path });
    //     const url = window.URL.createObjectURL(new Blob([response.data]));
    //     const link = document.createElement('a');
    //     link.href = url;
    //     link.download = file.name
    //     document.body.appendChild(link);
    //     link.click();
    //     document.body.removeChild(link);
    //     window.URL.revokeObjectURL(url)
    //   } catch (err) {
    //     this.$message.error('下载失败: ' + (err.message || '未知错误'))
    //   }
    // },
    // 下载文件（简单版：原生a标签流式下载，无错误捕获）
    async handleDownload(file) {
      if (file.isDir) {
        this.$message.warning('请选择文件进行下载')
        return
      }
      try {
        // 1. 拼接下载接口URL，对path参数编码（避免特殊字符/中文导致请求错误）
        const baseUrl = '/dev-api/sftp/download'
        const params = new URLSearchParams()
        params.append('path', file.path)
        const sftpToken = localStorage.getItem('sftp_token')
        const downloadUrl = `${baseUrl}?sftp_token=${sftpToken}&${params.toString()}`

        // 2. 创建原生a标签，触发浏览器原生下载（流式传输，不占前端内存）
        const link = document.createElement('a')
        link.href = downloadUrl
        link.download = file.name // 下载后的文件名（可自定义）
        link.style.display = 'none'
        document.body.appendChild(link)
        link.click() // 触发下载

        // 3. 清理DOM
        document.body.removeChild(link)
        window.URL.revokeObjectURL(link.href)
      } catch (err) {
        this.$message.error('下载失败: ' + (err.message || '未知错误'))
      }
    },
    // 下载目录（简单版：原生a标签流式下载，无错误捕获）
    async handleDownloadDir(file) {
      if (!file.isDir) {
        this.$message.warning('请选择目录进行下载')
        return
      }
      try {
        // 1. 拼接下载接口URL，对path参数编码（避免特殊字符/中文导致请求错误）
        const baseUrl = '/dev-api/sftp/downloaddir'
        const params = new URLSearchParams()
        params.append('path', file.path)
        const sftpToken = localStorage.getItem('sftp_token')
        const downloadUrl = `${baseUrl}?sftp_token=${sftpToken}&${params.toString()}`

        // 2. 创建原生a标签，触发浏览器原生下载（流式传输，不占前端内存）
        const link = document.createElement('a')
        link.href = downloadUrl
        link.download = file.name // 下载后的文件名（可自定义）
        link.style.display = 'none'
        document.body.appendChild(link)
        link.click() // 触发下载

        // 3. 清理DOM
        document.body.removeChild(link)
        window.URL.revokeObjectURL(link.href)
      } catch (err) {
        this.$message.error('下载失败: ' + (err.message || '未知错误'))
      }
    },
    // 删除文件或目录按钮的回调
    handleDelete(file) {
      this.deleteTarget = {
        name: file.name,
        path: file.path,
        isDir: file.isDir
      };
      this.deleteDialogVisible = true;
    },
    // 确认删除
    async confirmDelete() {
      try {
        const result = await this.$API.sftpuser.reqSftpDelete({
          path: this.deleteTarget.path
        })
        if (result.code === 200) {
          this.$message.success('删除成功');
          this.fetchFiles();
        } else {
          this.$message.error(result.message || '删除失败');
        }
      } catch (e) { }
      finally {
        this.deleteDialogVisible = false;
      }
    },
    // 重命名按钮的回调
    startRename(item) {
      this.$set(item, 'isRenaming', true)
      this.$set(item, 'editName', item.name)
      this.renamingItem = item
    },
    // 重命名输入框回车，确认重命名
    async confirmRename(item) {
      if (!item.editName || item.editName.trim() === '') {
        this.$message.warning('名称不能为空')
        return
      }
      if (item.editName === item.name) {
        this.cancelRename(item)
        return
      }
      try {
        const result = await this.$API.sftpuser.reqSftpRename({
          oldPath: item.path,
          newName: item.editName
        })
        if (result.code == 200) {
          this.$message.success('重命名成功')
          this.fetchFiles()
        } else {
          this.$message.error(result.message || '重命名失败')
        }
      } catch (error) {
        this.cancelRename(item)
      }
    },
    // 取消重命名
    cancelRename(item) {
      item.isRenaming = false
      item.editName = ''
    },
    // 根据不同登录类型默认显示email表单
    async showEmailForm() {
      if (this.userForm.loginType == 'KeyFile' || this.userForm.loginType == 'both') {
        this.userForm.emailOrNot = true
        try {
          const result = await this.$API.contact.reqContactOptions()
          if (result.code == 200) {
            this.options = result.data.options
          }
        } catch (error) {
        }
      } else {
        this.userForm.emailOrNot = false
        this.userForm.cc = ''
      }
    },
    // 切换登录类型时，默认显示email表单
    handleClick(tab) {
      if (tab.name === 'password') {
        this.$nextTick(() => {
          this.$refs.username.focus()
        })
      }
      else {
        this.$nextTick(() => {
          this.$refs.keyUsername.focus()
        })
      }
    },
    // 新增：文件上传进度钩子处理方法
    handleUploadProgress(event, file, fileList) {
      // event.percent 是Element UI返回的上传进度（0-100的数值）
      this.uploadPercent = Math.round(event.percent); // 取整，避免小数抖动
      this.showUploadProgress = true; // 显示进度条
      this.isUploading = true; // 标记正在上传
    },
  }
}
</script>

<style scoped>
.userheader {
  display: flex;
  justify-content: space-between;
}

.demo-drawer__footer {
  display: flex;
  justify-content: flex-end;
}

.breadcrumb {
  cursor: pointer;
}

.breadcrumbBold {
  font-weight: bold;
  color: #409EFF
}

.breadcrumb:hover {
  color: #409EFF;
  font-weight: bold;
}

.goBack:hover {
  color: #409EFF;
  font-weight: bold;
}

.dir-item {
  color: #409EFF;
  cursor: pointer;
}

.file-item {
  color: #606266;
}

.operate {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.upload {
  display: inline-block;
  margin-right: 10px;
}

.dialog-footer {
  display: flex;
  margin-top: 20px;
  justify-content: center;
}
</style>

<style scoped lang="scss">
// 根容器核心布局：实现登录框垂直水平居中
.page-container {
  position: relative; // 让logo绝对定位基于根容器，不跑偏
  display: flex;
  justify-content: center; // 水平居中
  align-items: center;     // 垂直居中
  width: 100vw;           // 占满视口宽度
  min-height: 100vh;      // 占满视口高度
  overflow: hidden;       // 禁止容器内部滚动
  padding: 20px;          // 小屏幕留边距，避免贴边
  box-sizing: border-box; // 内边距不占用宽度，防止溢出
}

// 背景图片容器（优化适配，无拉伸）
.bg-container {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background: url(~@/assets/1.png);
  background-size: 100% 100%; // 替代100%100%，避免图片拉伸变形
  background-repeat: no-repeat;
  background-attachment: fixed;
  background-position: center; // 图片居中显示
  z-index: -2; // 放在最下层，不遮挡其他元素
}

// logo样式（原有样式保留，无修改）
.logo-container {
  position: absolute;
  top: 20px;
  left: 20px;
  z-index: 10;
  transition: all 0.3s ease;

  &:hover {
    transform: scale(1.05);
  }
}

.logo {
  width: 60px;
  height: 60px;
  object-fit: contain;
  filter: drop-shadow(0 2px 8px rgba(100, 200, 150, 0.15));
}

// 入场动画（原有样式保留，无修改）
@keyframes formFadeIn {
  0% {
    opacity: 0;
    transform: translateY(20px) scale(0.78);
  }
  100% {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

// SFTP登录框：调整布局，确保居中生效
.sftp-login-fade {
  animation: formFadeIn 0.8s cubic-bezier(0.25, 0.8, 0.25, 1);
  z-index: 5; // 确保在背景上层
  margin: 0 !important; // 覆盖原有margin，让flex居中生效
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12); // 轻微阴影，提升视觉层次
}

// 全局禁止body滚动条，彻底消除页面滚动
:deep(body) {
  margin: 0;
  padding: 0;
  overflow: hidden !important;
}
</style>