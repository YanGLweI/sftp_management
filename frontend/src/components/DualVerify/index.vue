<template>
  <el-dialog
    title="双控验证"
    :visible.sync="dialogVisible"
    width="500px"
    append-to-body
    custom-class="dual-verify-dialog"
    :close-on-click-modal="false"
    @closed="handleClosed"
  >
    <div class="dual-verify-body">
      <p class="dual-verify-desc">操作「{{ actionDesc }}」需要复核验证!</p>
      <el-form label-position="top" @submit.native.prevent="confirm">
        <el-form-item label="复核账号">
          <el-input
            v-model="form.username"
            placeholder="请输入另一账号"
            autocomplete="off"
            ref="dualUsername"
            @keyup.enter.native="confirm"
          ></el-input>
        </el-form-item>
        <el-form-item label="密码">
          <el-input
            v-model="form.password"
            type="password"
            show-password
            placeholder="请输入密码"
            @keyup.enter.native="confirm"
          ></el-input>
        </el-form-item>
      </el-form>
      <el-alert v-if="errorMsg" :title="errorMsg" type="error" :closable="false" show-icon></el-alert>
    </div>
    <div slot="footer">
      <el-button @click="cancel">取消</el-button>
      <el-button type="primary" :loading="verifying" @click="confirm">{{ verifying ? '验证中...' : '验 证' }}</el-button>
    </div>
  </el-dialog>
</template>

<script>
import { reqSftpDualVerify } from '@/api/sftp/sftpuser'
import { rsaEncrypt } from '@/utils/encrypt'

export default {
  name: 'DualVerify',
  props: {
    // 当前登录的产业部账号（复核账号不得与之相同）
    loginDomainUser: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      dialogVisible: false,
      actionDesc: '', // 当前操作描述
      form: { username: '', password: '' },
      verifying: false,
      errorMsg: '',
      resolver: null // 当前验证的 Promise resolver
    }
  },
  methods: {
    // 打开双控验证弹窗，验证通过后 resolve 双控凭证 token
    verify(actionDesc) {
      this.actionDesc = actionDesc
      this.form = { username: '', password: '' }
      this.errorMsg = ''
      this.dialogVisible = true
      this.$nextTick(() => {
        if (this.$refs.dualUsername) this.$refs.dualUsername.focus()
      })
      return new Promise((resolve, reject) => {
        this.resolver = { resolve, reject }
      })
    },
    async confirm() {
      const { username, password } = this.form
      if (!username || !password) {
        this.errorMsg = '请输入复核账号和密码'
        return
      }
      if (this.loginDomainUser && username.toLowerCase() === this.loginDomainUser.toLowerCase()) {
        this.errorMsg = '复核账号不能与当前登录账号相同'
        return
      }
      this.verifying = true
      this.errorMsg = ''
      try {
        const rsaPassword = await rsaEncrypt(password)
        const res = await reqSftpDualVerify({ username, password: rsaPassword })
        if (res.code === 200) {
          this.dialogVisible = false
          this.resolver && this.resolver.resolve(res.data.dual_token)
        }
      } catch (error) {
        // 后端校验失败（非产业部账号/密码错误等），保留弹窗并展示错误原因
        this.errorMsg = error.message || '双控验证失败'
        this.resolver?.reject(new Error(this.errorMsg))
        this.resolver = null
      } finally {
        this.verifying = false
      }
    },
    cancel() {
      this.dialogVisible = false
    },
    // 弹窗关闭（取消/验证成功后）：未完成的 Promise 需要 reject，调用方中止操作
    handleClosed() {
      if (this.resolver) {
        this.resolver.reject(new Error('双控验证已取消'))
        this.resolver = null
      }
    }
  }
}
</script>

<style>
.dual-verify-dialog .dual-verify-body {
  padding: 0 4px;
}
.dual-verify-dialog .dual-verify-desc {
  margin: 0 0 16px 0;
  color: #606266;
  font-size: 14px;
  line-height: 1.6;
}
</style>
