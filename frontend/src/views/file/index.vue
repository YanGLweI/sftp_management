<template>
  <div class="page-container">
    <!-- 背景图片容器 -->
    <div class="bg-container"></div>
    <!-- 左上角logo -->
    <div class="logo-container">
      <img :src="logoPath" alt="系统logo" class="logo" @click="goToLoginPage" />
    </div>
    <!-- sftp登录-->
    <el-card class="sftp-login-fade apple-login-card">
      <div class="apple-card-header">
        <h1 class="apple-title">SFTP登录</h1>
      </div>
      <el-tabs class="apple-segment" v-model="activeName" @tab-click="handleClick">
        <el-tab-pane label="密码登录" name="password">
          <el-form :model="SftpForm" label-position="top">
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
          <el-form :model="SftpForm" label-position="top">
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
        <!-- 标签上传 Tab：仅在 enabled 时为真时显示 -->
        <el-tab-pane 
          label="标签上传" 
          name="hotlabel" 
          v-if="moduleConfigs.hotlabel && moduleConfigs.hotlabel.enabled">
          <el-form :model="SftpForm" label-position="top">
            <template v-if="moduleConfigs.hotlabel.loginType === 'local'">
              <el-form-item label="账号" :label-width="SftpFormLabelWidth">
                <el-input v-model="SftpForm.username" autocomplete="off"></el-input>
              </el-form-item>
              <el-form-item label="密码" :label-width="SftpFormLabelWidth">
                <el-input
                  ref="labelLocal"
                  v-model="SftpForm.password"
                  autocomplete="off"
                  type="password"
                  show-password
                  v-focus="SftpDialogFormVisible"
                  @keyup.enter.native="sftplogin()"
                ></el-input>
              </el-form-item>
            </template>
            <template v-else>
              <el-form-item label="账号" :label-width="SftpFormLabelWidth">
                <el-input v-model="SftpForm.username" autocomplete="off"></el-input>
              </el-form-item>
              <el-form-item label="密码" :label-width="SftpFormLabelWidth">
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
            </template>
          </el-form>
        </el-tab-pane>
        <!-- 中国联通 Tab：仅在 enabled 时为真时显示 -->
        <el-tab-pane 
          label="中国联通" 
          name="chinaunicom" 
          v-if="moduleConfigs.chinaunicom && moduleConfigs.chinaunicom.enabled">
          <el-form :model="SftpForm" label-position="top">
            <template v-if="moduleConfigs.chinaunicom.loginType === 'local'">
              <el-form-item label="账号" :label-width="SftpFormLabelWidth">
                <el-input v-model="SftpForm.username" autocomplete="off"></el-input>
              </el-form-item>
              <el-form-item label="密码" :label-width="SftpFormLabelWidth">
                <el-input
                  ref="unicomLabelLocal"
                  v-model="SftpForm.password"
                  autocomplete="off"
                  type="password"
                  show-password
                  v-focus="SftpDialogFormVisible"
                  @keyup.enter.native="sftplogin()"
                ></el-input>
              </el-form-item>
            </template>
            <template v-else>
              <el-form-item label="账号" :label-width="SftpFormLabelWidth">
                <el-input v-model="SftpForm.username" autocomplete="off"></el-input>
              </el-form-item>
              <el-form-item label="密码" :label-width="SftpFormLabelWidth">
                <el-input
                  ref="unicomLabel"
                  v-model="SftpForm.password"
                  autocomplete="off"
                  type="password"
                  show-password
                  v-focus="SftpDialogFormVisible"
                  @keyup.enter.native="sftplogin()"
                ></el-input>
              </el-form-item>
            </template>
          </el-form>
        </el-tab-pane>
      </el-tabs>
      <div class="dialog-footer apple-footer">
        <el-button @click="closeSftpDialogForm()">取 消</el-button>
        <el-button
          type="primary"
          @click="sftplogin()"
          :loading="buttonLoading"
        >{{ buttonLoading ? '提交中 ...' : '确 定' }}</el-button>
      </div>

      <div class="back-to-platform-link">
        <span class="link-text" @click="goToLoginPage">
          返回 SFTP 管理平台
        </span>
      </div>
    </el-card>

    <!-- 公共 SFTP 浏览器组件 -->
    <SftpBrowser
      :path="path"
      :username="SftpForm.username"
      :visible="SftpBrowserVisible"
      :upload-headers="uploadHeaders"
      :dual-verify-enabled="dualAuthEnabledMap[currentLoginType] || false"
      :login-domain-user="moduleConfigs[currentLoginType] === 'local' ? '' : SftpForm.username"
      upload-url="/dev-api/sftp/upload"
      @close="closeSftpBrowser"
    />

    <!-- 修改密码弹框（公共组件：自动聚焦新密码、回车提交、标题统一） -->
    <ChangePasswordDialog
      :visible.sync="changePasswordDialogVisible"
      :old-password="changeOldPassword"
      :change-token="changeToken"
      @success="handleChangeSuccess"
    />
  </div>
</template>

<script>
import SftpBrowser from '@/components/SftpBrowser'
import ChangePasswordDialog from '@/components/ChangePasswordDialog'
import sftpModulesApi from '@/api/admin/sftpModules'
import { rsaEncrypt } from '@/utils/encrypt'

export default {
  name: 'File',
  components: { SftpBrowser, ChangePasswordDialog },
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
      currentLoginType: 'password', // 当前登录的模块类型（password/hotlabel/chinaunicom），用于决定双控开关
      keyFileList: [],
      KeyFileUploadUrl: '/dev-api/sftp/login',
      uploadHeaders: {
        'Token': `${this.$store.state.user.token}`,
        'X-SFTP-Token': ''
      },
      path: '', //指定登录sftp后的路径
      // SFTP 模块配置 (登录方式 + 启用状态),默认 ldap 兼容未配置场景
      moduleConfigs: {
        hotlabel: {
          loginType: 'ldap',
          enabled: false  // 新增：默认不启用
        },
        chinaunicom: {
          loginType: 'ldap',
          enabled: false  // 新增：默认不启用
        }
      },
      // 双控开关配置（仅中国联通生效）
      dualAuthEnabledMap: {
        hotlabel: false,
        chinaunicom: true
      },
      // 修改密码弹框（账号需修改密码时弹出，公共组件）
      changePasswordDialogVisible: false,
      changeToken: '',
      changeOldPassword: '',
    }
  },
  mounted() {
    this.fetchModuleConfigs()
    this.$nextTick(() => {
      this.$refs.username.focus()
    })
  },
  methods: {
    goToLoginPage() {
      this.$router.push({ path: '/login' })
    },
    // 获取 SFTP 模块配置（登录方式 + 双控开关），用于动态渲染登录表单
    // 使用公共接口 /sftp/module-configs：/file 是公共页面（无需平台 token），与 /login 平级
    async fetchModuleConfigs() {
      try {
        const res = await sftpModulesApi.getPublicConfigs()
        if (res.code === 200 && Array.isArray(res.data)) {
          const configMap = {}
          res.data.forEach(item => {
            if (item.moduleName === 'hotlabel' || item.moduleName === 'chinaunicom') {
              configMap[item.moduleName] = {
                loginType: item.loginType === 'local' ? 'local' : 'ldap',
                enabled: !!item.enabled,  // 新增：解析 enabled 字段
              }
              this.dualAuthEnabledMap[item.moduleName] = !!item.dualAuthEnabled
            }
          })
          this.moduleConfigs = { ...this.moduleConfigs, ...configMap }
        }
      } catch (error) {
        console.error('获取模块配置失败:', error)
      }
    },
    // 修改密码成功：提示重新登录并清空登录表单
    handleChangeSuccess() {
      this.$message.success('密码修改成功，请重新登录')
      this.changeToken = ''
      this.changeOldPassword = ''
      this.SftpForm = { username: '', password: '' }
    },
    async sftplogin() {
      if (this.activeName == 'password') {
        const { username, password } = this.SftpForm
        if (!username || !password) return this.$message.warning('请输入账号密码')
        this.buttonLoading = true
        try {
          const rsaPassword = await rsaEncrypt(password)
          const res = await this.$API.sftpuser.reqSftpLogin({ username, password: rsaPassword })
          if (res.code == 200) {
            this.currentLoginType = this.activeName
            sessionStorage.setItem("sftp_token", res.data.sftp_token)
            this.uploadHeaders['X-SFTP-Token'] = res.data.sftp_token
            this.$message.success('SFTP 登录成功')
            this.SftpBrowserVisible = true
          }
        } catch {} finally {
          this.SftpDialogFormVisible = false
          this.buttonLoading = false
        }
      } else if (this.activeName == 'hotlabel' || this.activeName == 'chinaunicom') {
        // 标签上传/中国联通：提交账号密码，由后端根据模块配置决定本地或 LDAP 验证
        const { username, password } = this.SftpForm
        const isLocal = this.moduleConfigs[this.activeName] === 'local'
        if (!username || !password) return this.$message.warning(isLocal ? '请输入本地账号密码' : '请输入域账号密码')
        this.buttonLoading = true
        try {
          const rsaPassword = await rsaEncrypt(password)
          const res = await this.$API.sftpuser.reqSftpLogin({ username, password: rsaPassword, loginType: this.activeName })
          if (res.code == 200) {
            if (res.data && res.data.must_change_password) {
              // 账号需修改密码：弹出修改密码弹框（公共组件，旧密码由 changeOldPassword 预填）
              this.changeOldPassword = password
              this.changeToken = res.data.change_token
              this.changePasswordDialogVisible = true
              return
            }
            this.currentLoginType = this.activeName
            sessionStorage.setItem("sftp_token", res.data.sftp_token)
            this.uploadHeaders['X-SFTP-Token'] = res.data.sftp_token
            this.path = this.activeName == 'hotlabel' ? '/hotlabel' : '/ChinaUnicom'
            this.$message.success('SFTP登录成功')
            this.SftpBrowserVisible = true
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
        this.currentLoginType = this.activeName
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
      if (this.activeName == 'hotlabel' || this.activeName == 'chinaunicom') {
        if (this.activeName == 'hotlabel') {
          if (this.moduleConfigs.hotlabel?.loginType === 'local') {
            this.$refs.labelLocal ? this.$refs.labelLocal.focus() : this.$refs.username.focus()
          } else {
            this.$refs.label.focus()
          }
        } else {
          if (this.moduleConfigs.chinaunicom?.loginType === 'local') {
            this.$refs.unicomLabelLocal ? this.$refs.unicomLabelLocal.focus() : this.$refs.username.focus()
          } else {
            this.$refs.unicomLabel.focus()
          }
        }
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
      // 关闭浏览器后将焦点还给当前标签页的登录输入框
      this.$nextTick(() => {
        if (this.activeName === 'keyfile') {
          this.$refs.keyUsername.focus()
        } else if (this.activeName === 'hotlabel') {
          this.$refs.label.focus()
        } else if (this.activeName === 'chinaunicom') {
          this.$refs.unicomLabel.focus()
        } else {
          this.$refs.username.focus()
        }
      })
    },
    handleClick(tab) {
      this.$nextTick(() => {
        if (tab.name === 'hotlabel') {
          if (this.moduleConfigs.hotlabel?.loginType === 'local') {
            this.$refs.labelLocal ? this.$refs.labelLocal.focus() : this.$refs.username.focus()
          } else {
            this.$refs.label.focus()
          }
        } else if (tab.name === 'chinaunicom') {
          if (this.moduleConfigs.chinaunicom?.loginType === 'local') {
            this.$refs.unicomLabelLocal ? this.$refs.unicomLabelLocal.focus() : this.$refs.username.focus()
          } else {
            this.$refs.unicomLabel.focus()
          }
        } else {
          this.SftpForm = {
            username: '',
            password: '',
          }
          tab.name === 'password' ? this.$refs.username.focus() : this.$refs.keyUsername.focus()
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
@import '@/styles/variables.scss';

// ===== Apple 圆润风格设计 tokens =====
$apple-blue: #0071e3;
$apple-blue-hover: #0077ed;
$apple-text: #1d1d1f;
$apple-text-secondary: #6e6e73;
$apple-text-disabled: #8e8e93;
$apple-fill: rgba(120, 120, 128, 0.08);
$apple-fill-strong: rgba(120, 120, 128, 0.12);

// 颜色变量
$light-gray: #f8fafc;
$primary-color: #64c896;

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
  font-family: $font-family-base;
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
  0% { opacity: 0; transform: translateY(20px) scale(0.92); }
  100% { opacity: 1; transform: translateY(0) scale(1); }
}
.sftp-login-fade {
  animation: formFadeIn 0.5s cubic-bezier(0.25, 0.8, 0.25, 1);
  z-index: 5;
  margin: 0 !important;
}

// ===== 磨砂玻璃圆润卡片 =====
.apple-login-card {
  width: 480px;
  max-width: 100%;
  border: 1px solid rgba(255, 255, 255, 0.45);
  border-radius: 24px;
  background: rgba(255, 255, 255, 0.38);
  backdrop-filter: blur(14px) saturate(160%);
  -webkit-backdrop-filter: blur(14px) saturate(160%);
  box-shadow: 0 24px 48px rgba(0, 0, 0, 0.16), 0 4px 12px rgba(0, 0, 0, 0.08);

  ::v-deep .el-card__body {
    padding: 32px;
  }

  // 按钮通用（footer + 密钥上传）
  ::v-deep .el-button {
    height: 44px;
    padding: 0 24px;
    border: none;
    border-radius: 12px;
    font-family: $font-family-base;
    font-size: 15px;
    font-weight: 600;
    transition: background-color 0.2s ease, box-shadow 0.2s ease, transform 0.15s ease;
    &:active { transform: scale(0.98); }
  }
  ::v-deep .el-button--primary {
    background: $apple-blue;
    color: #fff;
    &:hover, &:focus { background: $apple-blue-hover; color: #fff; }
  }
  ::v-deep .el-button:not(.el-button--primary) {
    background: rgba(255, 255, 255, 0.14);
    color: #fff;
    &:hover, &:focus { background: rgba(255, 255, 255, 0.22); color: #fff; }
  }
}

.apple-card-header {
  margin-bottom: 24px;
  text-align: center;
}
.apple-title {
  margin: 0;
  font-size: 22px;
  font-weight: 700;
  letter-spacing: 0.5px;
  color: #fff;
}

// ===== iOS 分段控件式标签页 =====
.apple-segment {
  ::v-deep .el-tabs__header {
    margin: 0 0 24px;
  }
  ::v-deep .el-tabs__nav-wrap::after {
    display: none;
  }
  ::v-deep .el-tabs__nav {
    position: relative;
    display: flex;
    width: 100%;
    padding: 2px;
    background: rgba(255, 255, 255, 0.14);
    border-radius: 10px;
  }
  // 白色滑块：复用 element 的 active-bar（自带 transform 过渡），在标签间平滑滑动
  ::v-deep .el-tabs__active-bar {
    top: 2px;
    bottom: auto;
    height: 32px;
    border-radius: 8px;
    background: #fff;
    box-shadow: 0 1px 4px rgba(0, 0, 0, 0.12);
    z-index: 0;
  }
  ::v-deep .el-tabs__item {
    position: relative;
    z-index: 1;
    flex: 1;
    height: 32px;
    line-height: 32px;
    padding: 0 !important;
    text-align: center;
    font-size: 14px;
    color: rgba(255, 255, 255, 0.7);
    transition: color 0.2s ease;
    &:hover:not(.is-active) { color: #fff; }
    &.is-active {
      color: $apple-text;
      font-weight: 600;
    }
  }
  ::v-deep .el-tabs__content {
    padding: 0;
  }

  // 表单
  ::v-deep .el-form-item {
    margin-bottom: 20px;
  }
  ::v-deep .el-form-item__label {
    padding: 0 0 6px;
    line-height: 1.2;
    font-size: 13px;
    font-weight: 500;
    color: rgba(255, 255, 255, 0.75);
  }
  ::v-deep .el-input__inner {
    height: 44px;
    line-height: 44px;
    border: 1px solid transparent;
    border-radius: 12px;
    background: rgba(0, 0, 0, 0.10);
    color: #fff;
    transition: background-color 0.2s ease, border-color 0.2s ease, box-shadow 0.2s ease;
    &::placeholder { color: rgba(255, 255, 255, 0.6); }
    &:focus {
      background: rgba(0, 0, 0, 0.16);
      border-color: transparent;
      box-shadow: none;
    }
  }
  ::v-deep .el-input__icon {
    color: rgba(255, 255, 255, 0.7);
    line-height: 44px;
  }
  ::v-deep .el-input.is-disabled .el-input__inner {
    background: rgba(0, 0, 0, 0.06);
    border-color: transparent;
    color: rgba(255, 255, 255, 0.5);
  }

  // 密钥上传按钮改为次级灰底风格（高度与输入框一致，避免切 tab 抖动）
  ::v-deep .upload-key .el-button {
    height: 44px;
    padding: 0 20px;
    background: rgba(255, 255, 255, 0.14);
    color: #fff;
    font-size: 14px;
    &:hover, &:focus { background: rgba(255, 255, 255, 0.22); color: #fff; }
  }
  ::v-deep .el-upload-list__item {
    border-radius: 8px;
  }
}

// ===== 底部按钮等宽 =====
.apple-footer {
  display: flex;
  gap: 12px;
  justify-content: center;
  margin-top: 28px;

  ::v-deep .el-button {
    flex: 1;
  }
}

// 返回平台链接
.back-to-platform-link {
  text-align: center;
  margin-top: 15px;

  .link-text {
    display: inline-block;
    color: rgba($light-gray, 0.9);
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.3s ease;
    padding: 4px 0;

    &:hover {
      color: rgba($light-gray, 1);
      transform: scale(1.08);
    }
  }
}

:deep(body) {
  margin: 0; padding: 0; overflow: hidden !important;
}

// 无障碍：偏好减少动效时关闭动画
@media (prefers-reduced-motion: reduce) {
  .sftp-login-fade {
    animation: none;
  }
  .apple-login-card,
  .apple-segment,
  .apple-footer {
    ::v-deep * {
      transition: none !important;
    }
  }
}
</style>