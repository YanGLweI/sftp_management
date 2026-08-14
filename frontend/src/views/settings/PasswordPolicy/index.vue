<template>
  <div class="app-container">
    <el-card>
      <div slot="header">
        <span>密码策略配置</span>
      </div>
      <el-form ref="policyForm" :model="policyForm" label-width="180px" label-position="left">
        <el-divider content-position="left">密码复杂度</el-divider>
        <el-form-item label="密码最小长度">
          <el-input-number v-model="policyForm.minLength" :min="6" :max="64" />
          <span style="margin-left: 8px; color: #999;">位字符</span>
        </el-form-item>
        <el-form-item label="需要大写字母">
          <el-switch v-model="policyForm.requireUppercase" />
        </el-form-item>
        <el-form-item label="需要小写字母">
          <el-switch v-model="policyForm.requireLowercase" />
        </el-form-item>
        <el-form-item label="需要数字">
          <el-switch v-model="policyForm.requireDigit" />
        </el-form-item>
        <el-form-item label="需要特殊字符">
          <el-switch v-model="policyForm.requireSpecialChar" />
        </el-form-item>

        <el-divider content-position="left">密码过期</el-divider>
        <el-form-item label="密码有效期">
          <el-input-number v-model="policyForm.expiryDays" :min="0" :max="3650" />
          <span style="margin-left: 8px; color: #999;">天（0=永不过期）</span>
        </el-form-item>

        <el-divider content-position="left">安全策略</el-divider>
        <el-form-item label="禁止使用历史密码次数">
          <el-input-number v-model="policyForm.passwordHistory" :min="0" :max="50" />
          <span style="margin-left: 8px; color: #999;">次（0=不限制）</span>
        </el-form-item>
        <el-form-item label="最大连续失败登录次数">
          <el-input-number v-model="policyForm.maxLoginAttempts" :min="0" :max="100" />
          <span style="margin-left: 8px; color: #999;">次（0=不限制，超出后账号自动锁定）</span>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="saving" @click="handleSave">保存配置</el-button>
          <el-button @click="fetchPolicy">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script>
import { getPasswordPolicy, updatePasswordPolicy } from '@/api/settings'

export default {
  name: 'PasswordPolicy',
  data() {
    return {
      policyForm: {
        minLength: 14,
        requireUppercase: true,
        requireLowercase: true,
        requireDigit: true,
        requireSpecialChar: true,
        expiryDays: 90,
        passwordHistory: 5,
        maxLoginAttempts: 5
      },
      saving: false
    }
  },
  created() {
    this.fetchPolicy()
  },
  methods: {
    async fetchPolicy() {
      try {
        const res = await getPasswordPolicy()
        if (res.code === 200) {
          const p = res.data
          this.policyForm = {
            minLength: p.minLength,
            requireUppercase: p.requireUppercase,
            requireLowercase: p.requireLowercase,
            requireDigit: p.requireDigit,
            requireSpecialChar: p.requireSpecialChar,
            expiryDays: p.expiryDays,
            passwordHistory: p.passwordHistory,
            maxLoginAttempts: p.maxLoginAttempts
          }
        }
      } catch (e) {
        console.error(e)
      }
    },
    async handleSave() {
      this.saving = true
      try {
        const res = await updatePasswordPolicy(this.policyForm)
        if (res.code === 200) {
          this.$message.success('密码策略已更新')
        } else {
          this.$message.error(res.message)
        }
      } catch (e) {
        console.error(e)
      }
      this.saving = false
    }
  }
}
</script>