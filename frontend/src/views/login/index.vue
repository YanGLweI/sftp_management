<template>
  <div class="login-wrapper">
    <div class="login-container">
      <!-- 背景图片容器 -->
      <div class="bg-container"></div>
      <!-- 左上角logo -->
      <div class="logo-container">
        <img :src="logoPath" alt="系统logo" class="logo" @click="goToFilePage" />
      </div>
      <!-- 登录表单 -->
      <el-form 
        ref="loginForm" 
        :model="loginForm" 
        :rules="loginRules" 
        class="login-form" 
        auto-complete="on" 
        label-position="left"
      >
        <div class="title-container">
          <h3 class="title">SFTP管理平台</h3>
        </div>

        <el-form-item prop="username">
          <span class="svg-container">
            <svg-icon icon-class="user" />
          </span>
          <el-input
            ref="username"
            v-model="loginForm.username"
            placeholder="Username"
            name="username"
            type="text"
            tabindex="1"
            autocomplete="off"
          />
        </el-form-item>

        <el-form-item prop="password">
          <span class="svg-container">
            <svg-icon icon-class="password" />
          </span>
          <el-input
            :key="passwordType"
            ref="password"
            v-model="loginForm.password"
            :type="passwordType"
            placeholder="Password"
            name="password"
            tabindex="2"
            autocomplete="new-password"
            @keyup.enter.native="handleLogin"
          />
          <span class="show-pwd" @click="showPwd">
            <svg-icon :icon-class="passwordType === 'password' ? 'eye' : 'eye-open'" />
          </span>
        </el-form-item>

        <el-form-item prop="loginType">
          <span class="svg-container">
            <svg-icon icon-class="earth" />
          </span>
          <el-select
            v-model="loginForm.loginType"
            placeholder="Login Type"
            name="loginType"
            tabindex="3"
            autocomplete="off"
          >
            <el-option label="本地登录" value="local" />
            <el-option label="LDAP登录" value="ldap" />
          </el-select>
        </el-form-item>

        <el-button 
          :loading="loading" 
          type="primary" 
          class="login-btn"
          @click.native.prevent="handleLogin"
        >
          登 录
        </el-button>

        <div class="tips">
          <!-- <span style="margin-right:20px;">username: admin</span> -->
          <!-- <span> password: any</span> -->
        </div>
      </el-form>
    </div>

    <!-- 修改密码弹框（公共组件：自动聚焦新密码、回车提交、标题统一） -->
    <ChangePasswordDialog
      :visible.sync="changePasswordDialogVisible"
      :old-password="loginForm.password"
      :change-token="currentChangeToken"
      @success="handleChangeSuccess"
    />
  </div>
</template>

<script>
import { validUsername } from '@/utils/validate'
import { rsaEncrypt } from '@/utils/encrypt'
import { setToken } from '@/utils/auth'
import { validatePassword } from '@/api/settings'
import ChangePasswordDialog from '@/components/ChangePasswordDialog'

export default {
  name: 'Login',
  components: { ChangePasswordDialog },
  data() {
    const validateUsername = (rule, value, callback) => {
      if (!value) {
        callback(new Error('Please enter the correct user name'))
      } else {
        callback()
      }
    }
    const validatePassword = (rule, value, callback) => {
      if (value.length < 14) {
        callback(new Error('The password can not be less than 14 digits'))
      } else {
        callback()
      }
    }
    // 新增登录类型验证（可选，保证必选）
    const validateLoginType = (rule, value, callback) => {
      if (!value) {
        callback(new Error('Please select login type'))
      } else {
        callback()
      }
    }
    return {
      loginForm: {
        username: '',
        password: '',
        // 初始化时从localStorage读取登录类型，无则为空
        loginType: localStorage.getItem('loginType') || ''
      },
      loginRules: {
        username: [{ required: true, trigger: 'blur', validator: validateUsername }],
        password: [{ required: true, trigger: 'blur', validator: validatePassword }],
        // 新增登录类型必选验证（可选，根据业务需求）
        loginType: [{ required: true, trigger: 'change', validator: validateLoginType }]
      },
      loading: false,
      passwordType: 'password',
      redirect: undefined,
      logoPath: require('@/assets/logo.png'),
      // 修改密码弹框
      changePasswordDialogVisible: false,
      currentChangeToken: '' // 存储受限 Token（需改密场景）
    }
  },
  watch: {
    $route: {
      handler: function(route) {
        this.redirect = route.query && route.query.redirect
      },
      immediate: true
    }
  },
  mounted() {
    this.$nextTick(() => {
      this.$refs.username.focus()
    })
  },
  methods: {
    // 点击logo跳转文件管理页面
    goToFilePage() {
      this.$router.push({ path: '/file' })
    },
    showPwd() {
      this.passwordType = this.passwordType === 'password' ? '' : 'password'
      this.$nextTick(() => {
        this.$refs.password.focus()
      })
    },
    handleLogin() {
      this.$refs.loginForm.validate(valid => {
        if (valid) {
          // 登录验证通过后，将当前选择的登录类型存入localStorage
          localStorage.setItem('loginType', this.loginForm.loginType)
          const {username,password,loginType} = this.loginForm
          const rsaPassword = rsaEncrypt(password)
          this.loading = true
          this.$store.dispatch('user/login', {username,password:rsaPassword,loginType}).then(res => {
            // 检查是否需改密或密码过期（旧密码由公共弹框组件从登录表单预填）
            if (res.data && (res.data.must_change_password || res.data.password_expired)) {
              // 传递受限 Token 给组件（仅限需改密/密码过期场景）
              this.currentChangeToken = res.data.token
              this.changePasswordDialogVisible = true
              this.loading = false
            } else {
              // 登录成功跳转：目标路由由路由守卫根据权限动态处理
              this.$router.push({ path: this.redirect || '/' })
              this.loading = false
            }
          }).catch(() => {
            this.loading = false
          })
        } else {
          console.log('error submit!!')
          return false
        }
      })
    },
    // 修改密码成功：更新 token 并跳转首页
    handleChangeSuccess(token) {
      this.$message.success('密码修改成功')
      if (token) {
        this.$store.commit('user/SET_TOKEN', token)
        setToken(token)
      }
      this.$router.push({ path: this.redirect || '/' })
    }
  }
}
</script>

<style lang="scss">
// 全局样式适配
$primary-color: #64c896; // 核心清新绿
$secondary-color: #74c0fc; // 辅助浅蓝
$light-gray: #f8fafc;
$mid-gray: #94a3b8;
$dark-gray: #2a3b47;
// 登录框背景调整为低透明度，保留样式但可见背景
$glass-bg: rgba(255, 255, 255, 0.2);
$glass-border: rgba(100, 200, 150, 0.2);

// 修复input样式
@supports (-webkit-mask: none) and (not (cater-color: $primary-color)) {
  .login-container .el-input input {
    color: $light-gray;
  }
}

// 重置element-ui样式
.login-container {
  .el-input{
    display: inline-block;
    height: 47px;
    width: 85%;

    input {
      background: transparent;
      border: 0;
      border-radius: 0;
      padding: 12px 5px 12px 15px;
      color: $light-gray; // 改为浅白色适配透明背景
      height: 47px;
      caret-color: $primary-color;
      font-size: 14px;

      &:-webkit-autofill {
        box-shadow: 0 0 0px 1000px rgba(255,255,255,0.1) inset !important;
        -webkit-text-fill-color: $light-gray !important;
      }
      
      &::placeholder {
        color: rgba(255,255,255,0.7) !important; // 浅白占位符
      }
    }
  }
  

  

  .el-input, .el-select {
    display: inline-block;
    height: 47px;
    width: 93%; // 和input保持一致的宽度

    // （如果需要）可以把input的子样式也同步给select的输入框
    &.el-select .el-select__wrapper {
      background: transparent;
      border: 0;
      height: 47px;
    }
    .el-input__suffix {
      right: -24px !important; // 固定箭头位置在右侧
    }
  }

  .el-form-item {
    border: 1px solid $glass-border;
    background: rgba(255, 255, 255, 0.1); // 低透明表单背景
    border-radius: 8px;
    color: $light-gray; // 浅白文字
    margin-bottom: 20px;
    transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);

    &:focus-within {
      border-color: $primary-color;
      box-shadow: 0 0 8px rgba(100, 200, 150, 0.2);
    }
  }

  // 登录按钮改为透明样式
  .el-button--primary {
    background: transparent; // 透明背景
    border: 1px solid rgba(100, 200, 150, 0.6); // 清新绿半透明边框
    color: #ffffff; // 白色文字
    box-shadow: 0 0 8px rgba(100, 200, 150, 0.1); // 轻微阴影
    transition: all 0.3s ease;

    &:hover {
      background: rgba(100, 200, 150, 0.1); // hover时轻微背景色
      border-color: rgba(100, 200, 150, 0.8); // 边框加深
      box-shadow: 0 0 12px rgba(100, 200, 150, 0.2); // 阴影增强
      transform: translateY(-2px); // 轻微上浮
      color: #ffffff; // 保持文字白色
    }

    &:active {
      transform: translateY(0); // 点击还原
      background: rgba(100, 200, 150, 0.05); // 点击时更浅的背景
    }

    // loading状态适配透明样式
    &.is-loading {
      background: transparent;
      border-color: rgba(100, 200, 150, 0.6);
      color: #ffffff;
    }
  }
}
</style>

<style lang="scss" scoped>
// 核心样式变量
$primary-color: #64c896;
$secondary-color: #74c0fc;
$light-gray: #f8fafc;
$mid-gray: rgba(255,255,255,0.7);
$dark-gray: #2a3b47;
$glass-bg: rgba(255, 255, 255, 0.05); // 登录框低透明度背景
$glass-border: rgba(100, 200, 150, 0.5);

// 纯背景图片容器（无任何遮挡）
.bg-container {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background: url(~@/assets/1.png);
  background-size: 100% 100%;
  background-repeat: no-repeat;
  background-attachment: fixed;
  z-index: -2; // 放在粒子下层
}

// logo样式优化
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

// 登录容器
.login-container {
  min-height: 100vh;
  width: 100%;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;

  // 登录表单（样式不变，仅背景透明度调整）
  .login-form {
    position: relative;
    width: 520px;
    max-width: 100%;
    padding: 40px 35px 30px;
    margin: 0 auto;
    overflow: hidden;
    // 低透明度玻璃背景，可见底层图片和粒子
    background: $glass-bg;
    // backdrop-filter: blur(1px);
    // -webkit-backdrop-filter: blur(1px);
    border: 1px solid $glass-border;
    border-radius: 12px;
    box-shadow: 
      0 8px 32px rgba(100, 200, 150, 0.1),
      0 4px 16px rgba(0, 0, 0, 0.08);
    animation: formFadeIn 0.8s cubic-bezier(0.25, 0.8, 0.25, 1);
  }

  // 标题样式（适配透明背景）
  .title-container {
    margin-bottom: 30px;

    .title {
      font-size: 24px;
      color: $light-gray; // 浅白色标题
      margin: 0 auto;
      text-align: center;
      font-weight: 600;
      letter-spacing: 0.5px;
      text-shadow: 0 0 4px rgba(0,0,0,0.2); // 文字阴影增强可读性
    }
  }

  // 图标容器（适配透明背景）
  .svg-container {
    padding: 6px 5px 6px 15px;
    color: $primary-color;
    vertical-align: middle;
    width: 30px;
    display: inline-block;
    font-size: 16px;
  }

  // 提示文字（适配透明背景）
  .tips {
    font-size: 14px;
    color: $mid-gray;
    margin-top: 10px;
    text-align: center;
  }

  // 显示密码按钮（适配透明背景）
  .show-pwd {
    position: absolute;
    right: 10px;
    top: 7px;
    font-size: 16px;
    color: $mid-gray;
    cursor: pointer;
    user-select: none;
    transition: all 0.2s ease;

    &:hover {
      color: $primary-color;
      transform: scale(1.1);
    }
  }

  // 登录按钮尺寸样式（保持不变）
  .login-btn {
    width: 100%;
    height: 48px;
    border-radius: 8px;
    font-size: 16px;
    font-weight: 600;
    letter-spacing: 0.3px;
    margin-bottom: 10px;
  }
}

// 表单入场动画
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
</style>