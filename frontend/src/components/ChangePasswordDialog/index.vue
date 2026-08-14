<template>
  <el-dialog
    title="请修改密码"
    :visible.sync="visible"
    width="450px"
    :close-on-click-modal="false"
    :show-close="false"
    @opened="focusNewPasswordInput"
  >
    <el-form label-width="100px">
      <el-form-item label="旧密码">
        <el-input v-model="form.oldPassword" type="password" disabled />
        <span style="font-size: 12px; color: #999;">（登录已验证）</span>
      </el-form-item>
      <el-form-item label="新密码" required>
        <el-input ref="newPasswordInput" v-model="form.newPassword" type="password" placeholder="请输入新密码" show-password />
      </el-form-item>
      <el-form-item label="确认密码" required>
        <el-input v-model="form.confirmPassword" type="password" placeholder="请再次输入新密码" show-password @keyup.enter.native="handleChange" />
      </el-form-item>
    </el-form>
    <span slot="footer">
      <el-button type="primary" :loading="loading" @click="handleChange">确认修改</el-button>
    </span>
  </el-dialog>
</template>

<script>
import axios from 'axios'
import { rsaEncrypt } from '@/utils/encrypt'
import { changePassword } from '@/api/user'

/**
 * ChangePasswordDialog 修改密码公共弹框组件
 * 平台登录（/login）与 SFTP 模块登录（/file）共用
 * props:
 *   - visible: 是否显示（配合 .sync）
 *   - oldPassword: 登录时输入的密码（预填旧密码框）
 *   - changeToken: 受限改密凭证（/file 使用）；平台登录传空，由请求拦截器自动携带平台 token
 * 事件:
 *   - success(token): 修改成功，token 为后端返回的新完整 token（平台登录使用；/file 忽略）
 *   - update:visible: 关闭弹框
 */
export default {
  name: 'ChangePasswordDialog',
  props: {
    visible: {
      type: Boolean,
      default: false
    },
    oldPassword: {
      type: String,
      default: ''
    },
    changeToken: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      form: {
        oldPassword: '',
        newPassword: '',
        confirmPassword: ''
      },
      loading: false
    }
  },
  watch: {
    visible(val) {
      if (val) {
        // 打开时用登录密码预填旧密码框，清空新密码/确认密码
        this.form = {
          oldPassword: this.oldPassword || '',
          newPassword: '',
          confirmPassword: ''
        }
      }
    }
  },
  methods: {
    // 弹框打开后自动聚焦到新密码输入框
    focusNewPasswordInput() {
      this.$nextTick(() => {
        if (this.$refs.newPasswordInput) {
          this.$refs.newPasswordInput.focus()
        }
      })
    },
    // 提交修改密码
    async handleChange() {
      const { newPassword, confirmPassword } = this.form
      if (!newPassword) {
        this.$message.warning('请输入新密码')
        return
      }
      if (newPassword !== confirmPassword) {
        this.$message.warning('两次输入的密码不一致')
        return
      }
      this.loading = true
      try {
        const rsaOldPwd = rsaEncrypt(this.form.oldPassword)
        const rsaNewPwd = rsaEncrypt(newPassword)
        let res
        if (this.changeToken) {
          // /file 公共页面（可能无平台 token）：使用登录时签发的受限改密凭证
          const resp = await axios.post('/dev-api/user/change-password', {
            oldPassword: rsaOldPwd,
            newPassword: rsaNewPwd
          }, { headers: { token: this.changeToken } })
          res = resp.data || {}
        } else {
          // 平台登录：请求拦截器自动携带平台 token
          res = await changePassword({
            oldPassword: rsaOldPwd,
            newPassword: rsaNewPwd
          })
        }
        if (res.code === 20000) {
          this.$emit('success', res.data ? res.data.token : '')
          this.$emit('update:visible', false)
        } else {
          this.$message.error(res.message || '密码修改失败')
        }
      } catch (error) {
        this.$message.error('密码修改失败，请重试')
        console.error('修改密码失败:', error)
      } finally {
        this.loading = false
      }
    }
  }
}
</script>
