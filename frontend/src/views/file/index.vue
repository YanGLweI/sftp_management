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
          <el-tab-pane label="标签上传" name="hotlabel">
            <el-form :model="SftpForm">
              <el-form-item label="SFTP账号" :label-width="SftpFormLabelWidth">
                <el-input v-model="SftpForm.username" autocomplete="off"  disabled></el-input>
              </el-form-item>
              <el-form-item label="SFTP密码" :label-width="SftpFormLabelWidth">
                <el-input
                  ref="label"
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

    <!-- 公共 SFTP 浏览器组件 -->
    <SftpBrowser
      :path="path"
      :username="SftpForm.username"
      :visible="SftpBrowserVisible"
      :upload-headers="uploadHeaders"
      upload-url="/dev-api/sftp/upload"
      @close="closeSftpBrowser"
    />
  </div>
</template>

<script>
import SftpBrowser from '@/components/SftpBrowser'
import { rsaEncrypt } from '@/utils/encrypt'

export default {
  name: 'File',
  components: { SftpBrowser },
  data() {
    // 密码验证器
    var validatePass = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请输入密码'));
      } else if (value.length < 14) {
        callback(new Error('密码长度至少为14位'));
      } else {
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
      loading: false,
      userlist: [],
      dialogFormVisible: false,
      title: '',
      dialogTableVisible: false,
      emailDialog: false,
      buttonLoading: false,
      showPassword: false,
      userForm: {
        id: '',
        name: '',
        loginType: 'Password',
        noExpire: false,
        emailOrNot: false,
        password: '',
        checkPass: '',
        to: '',
        cc: '',
      },
      options: [],
      SftpDialogFormVisible: false,
      SftpFormLabelWidth: '100px',
      SftpForm: {
        username: '',
        password: ''
      },
      email: { to: [], cc: [], subject: '' },
      searchObj: { username: '' },
      tempSearchObj: { username: '' },
      ids: [],
      multipleSelection: [],
      page: 1, limit: 5, total: 0,
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
      SftpBrowserVisible: false,
      activeName: 'password',
      keyFileList: [],
      KeyFileUploadUrl: '/dev-api/sftp/login',
      uploadHeaders: {
        'Token': `${this.$store.state.user.token}`,
        'X-SFTP-Token': ''
      },
      path: '', //指定登录sftp后的路径
    }
  },
  mounted() {
    this.$nextTick(() => {
      this.$refs.username.focus()
    })
  },
  methods: {
    goToLoginPage() {
      this.$router.push({ path: '/login' })
    },
    async sftplogin() {
      if (this.activeName == 'password') {
        const { username, password } = this.SftpForm
        const rsaPassword = rsaEncrypt(password)
        if (!username || !password) return this.$message.warning('请输入账号密码')
        this.buttonLoading = true
        try {
          const res = await this.$API.sftpuser.reqSftpLogin({ username, password: rsaPassword })
          if (res.code == 200) {
            sessionStorage.setItem("sftp_token", res.data.sftp_token)
            this.uploadHeaders['X-SFTP-Token'] = res.data.sftp_token
            this.$message.success('SFTP登录成功')
            this.SftpBrowserVisible = true
          }
        } catch {} finally {
          this.SftpDialogFormVisible = false
          this.buttonLoading = false
        }
      } else if (this.activeName == 'hotlabel') {
        const { username, password } = this.SftpForm
        const rsaPassword = rsaEncrypt(password)
        if (!username || !password) return this.$message.warning('请输入账号密码')
        this.buttonLoading = true
        try {
          const res = await this.$API.sftpuser.reqSftpLogin({ username, password: rsaPassword })
          if (res.code == 200) {
            this.SftpForm.username = 'HotLabel'
            const {VUE_APP_HotLabel_Username,VUE_APP_HotLabel_Password} = process.env
            const rsaPassword = rsaEncrypt(VUE_APP_HotLabel_Password)
            const result = await this.$API.sftpuser.reqSftpLogin({ username: VUE_APP_HotLabel_Username, password: rsaPassword })
            if (result.code == 200) {
              sessionStorage.setItem("sftp_token", result.data.sftp_token)
              this.uploadHeaders['X-SFTP-Token'] = result.data.sftp_token
              this.path = "/hotlabel"
              this.$message.success('SFTP登录成功')
              this.SftpBrowserVisible = true
            }
          }
        } catch {} finally {
          this.SftpDialogFormVisible = false
          this.buttonLoading = false
        }
      }else {
        const { username } = this.SftpForm
        if (!username) return this.$message.warning('请输入SFTP账号')
        if (!this.keyFileList.length) return this.$message.warning('请选择密钥文件')
        this.buttonLoading = true
        this.$refs.upload.submit();
      }
    },
    handleKeyFileChange(f, fileList) { this.keyFileList = fileList },
    beforeUploadKey(file) {
      const ok = file.size / 1024 / 1024 < 100
      if (!ok) this.$message.error('不能超过100MB')
      return ok
    },
    KeyhandleUploadSuccess(res) {
      if (res.code === 200) {
        sessionStorage.setItem("sftp_token", res.data.sftp_token);
        this.uploadHeaders['X-SFTP-Token'] = res.data.sftp_token
        this.$message.success('SFTP登录成功');
        this.SftpBrowserVisible = true
      } else {
        this.$message.error(res.message || '登录失败')
      }
      this.SftpDialogFormVisible = false
      this.buttonLoading = false
      this.keyFileList = []
      this.$refs.upload.clearFiles()
    },
    KeyhandleUploadError() {
      this.$message.error('登录失败');
      this.SftpDialogFormVisible = false
      this.buttonLoading = false
      this.keyFileList = []
      this.$refs.upload.clearFiles()
    },
    closeSftpDialogForm() {
      this.SftpDialogFormVisible = false
      this.buttonLoading = false
      this.activeName = 'password'
      this.keyFileList = []
      this.$refs.upload.clearFiles()
      this.SftpForm = { username: '', password: '' }
    },
    async closeSftpBrowser() {
      this.SftpBrowserVisible = false
      this.SftpForm = {
        username: '',
        password: '',
      }
      this.path = ''
      if (this.activeName == 'hotlabel') {
        this.SftpForm.username = 'HotData'
        this.$refs.label.focus()
        return
      }
      try {
        const result = await this.$API.sftpuser.reqSftpLogout();
        if (result.code === 200) {
          this.$message({ type: 'success', message: result.message });
        } else {
          this.$message({ type: 'error', message: result.message });
        }
      } catch {}
    },
    handleClick(tab) {
      this.$nextTick(() => {
        tab.name === 'password' ? this.$refs.username.focus() : this.$refs.keyUsername.focus()
        if (tab.name === 'hotlabel') {
          this.SftpForm.username = 'HotData'
          this.$refs.label.focus()
        }else{
          this.SftpForm = {
            username: '',
            password: '',
          }
        }
      })
    },
  }
}
</script>

<style scoped>
.userheader { display: flex; justify-content: space-between; }
.demo-drawer__footer { display: flex; justify-content: flex-end; }
.breadcrumb { cursor: pointer; }
.breadcrumbBold { font-weight: bold; color: #409EFF }
.breadcrumb:hover { color: #409EFF; font-weight: bold; }
.goBack:hover { color: #409EFF; font-weight: bold; }
.dir-item { color: #409EFF; cursor: pointer; }
.file-item { color: #606266; }
.operate { display: flex; justify-content: space-between; align-items: center; }
.upload { display: inline-block; margin-right: 10px; }
.dialog-footer { display: flex; margin-top: 20px; justify-content: center; }
</style>

<style scoped lang="scss">
.page-container {
  position: relative;
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100vw;
  min-height: 100vh;
  overflow: hidden;
  padding: 20px;
  box-sizing: border-box;
}
.bg-container {
  position: fixed; top: 0; left: 0;
  width: 100vw; height: 100vh;
  background: url(~@/assets/1.png);
  background-size: cover;
  background-repeat: no-repeat;
  background-position: center;
  z-index: -2;
}
.logo-container {
  position: absolute; top: 20px; left: 20px;
  z-index: 10;
  transition: all 0.3s ease;
  &:hover { transform: scale(1.05); }
}
.logo {
  width: 60px; height: 60px;
  object-fit: contain;
}
@keyframes formFadeIn {
  0% { opacity: 0; transform: translateY(20px) scale(0.78); }
  100% { opacity: 1; transform: translateY(0) scale(1); }
}
.sftp-login-fade {
  animation: formFadeIn 0.8s cubic-bezier(0.25, 0.8, 0.25, 1);
  z-index: 5;
  margin: 0 !important;
  box-shadow: 0 8px 24px rgba(0,0,0,0.12);
}
:deep(body) {
  margin: 0; padding: 0; overflow: hidden !important;
}
</style>