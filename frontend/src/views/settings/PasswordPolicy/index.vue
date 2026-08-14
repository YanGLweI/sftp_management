<template>
  <div class="password-policy-container">
    <el-card class="policy-config-card" shadow="never">
      <!-- 政策信息头部 -->
      <div slot="header" class="policy-header">
        <div class="policy-header__icon is-primary">
          <i class="el-icon-lock"></i>
        </div>
        <div class="policy-header__info">
          <div class="policy-header__title">
            密码策略
            <el-tag size="mini" type="success" effect="light" class="policy-header__tag">
              全局策略
            </el-tag>
          </div>
          <div class="policy-header__desc">配置系统级的密码复杂度、有效期与安全策略，保障账号安全</div>
        </div>
      </div>

      <!-- 设置内容区 -->
      <div class="policy-content">
        <!-- 区块一：密码复杂度 -->
        <div class="policy-section">
          <div class="policy-section__head">
            <h3 class="policy-section__title">
              <svg width="20" height="20" viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M22.8682 24.2982C25.4105 26.7935 26.4138 30.4526 25.4971 33.8863C24.5805 37.32 21.8844 40.0019 18.4325 40.9137C14.9806 41.8256 11.3022 40.8276 8.79375 38.2986C5.02208 34.4141 5.07602 28.2394 8.91499 24.4206C12.754 20.6019 18.9613 20.5482 22.8664 24.3L22.8682 24.2982Z" stroke="currentColor" stroke-width="4" stroke-linejoin="round"/>
                <path d="M23 24L40 7" stroke="currentColor" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
                <path d="M30.3052 16.9001L35.7337 22.3001L42.0671 16.0001L36.6385 10.6001L30.3052 16.9001Z" stroke="currentColor" stroke-width="4" stroke-linejoin="round"/>
              </svg>
              <span style="margin-left: 8px;">密码复杂度</span>
            </h3>
            <span class="policy-section__hint">设置密码的最小长度和字符组成要求</span>
          </div>

          <div class="complexity-cards">
            <!-- 最小长度 -->
            <div class="config-card">
              <div class="config-card__label">密码最小长度</div>
              <div class="config-card__control">
                <el-input-number 
                  v-model.number="policyForm.minLength" 
                  :min="8" 
                  :max="128" 
                  :step="2"
                  size="small"
                  controls-position="right"
                  class="length-control"
                />
                <span class="config-card__unit">位</span>
              </div>
              <div class="config-card__tips">
                <el-alert
                  title="建议至少 12-14 位以上"
                  type="info"
                  :closable="false"
                  show-icon
                  size="small"
                >
                </el-alert>
              </div>
            </div>

            <!-- 字符要求开关组 -->
            <div class="requirements-grid">
              <div 
                class="req-toggle" 
                :class="{ 'is-active': policyForm.requireUppercase }"
                @click="toggleReq('requireUppercase')"
              >
                <div class="req-toggle__icon">
                  <i :class="policyForm.requireUppercase ? 'el-icon-circle-check' : 'el-icon-plus'"></i>
                </div>
                <div class="req-toggle__text">
                  <div class="req-toggle__title">大写字母 A-Z</div>
                  <div class="req-toggle__status">
                    {{ policyForm.requireUppercase ? '已启用' : '已禁用' }}
                  </div>
                </div>
              </div>

              <div 
                class="req-toggle" 
                :class="{ 'is-active': policyForm.requireLowercase }"
                @click="toggleReq('requireLowercase')"
              >
                <div class="req-toggle__icon">
                  <i :class="policyForm.requireLowercase ? 'el-icon-circle-check' : 'el-icon-plus'"></i>
                </div>
                <div class="req-toggle__text">
                  <div class="req-toggle__title">小写字母 a-z</div>
                  <div class="req-toggle__status">
                    {{ policyForm.requireLowercase ? '已启用' : '已禁用' }}
                  </div>
                </div>
              </div>

              <div 
                class="req-toggle" 
                :class="{ 'is-active': policyForm.requireDigit }"
                @click="toggleReq('requireDigit')"
              >
                <div class="req-toggle__icon">
                  <i :class="policyForm.requireDigit ? 'el-icon-circle-check' : 'el-icon-plus'"></i>
                </div>
                <div class="req-toggle__text">
                  <div class="req-toggle__title">数字 0-9</div>
                  <div class="req-toggle__status">
                    {{ policyForm.requireDigit ? '已启用' : '已禁用' }}
                  </div>
                </div>
              </div>

              <div 
                class="req-toggle" 
                :class="{ 'is-active': policyForm.requireSpecialChar }"
                @click="toggleReq('requireSpecialChar')"
              >
                <div class="req-toggle__icon">
                  <i :class="policyForm.requireSpecialChar ? 'el-icon-circle-check' : 'el-icon-plus'"></i>
                </div>
                <div class="req-toggle__text">
                  <div class="req-toggle__title">特殊符号 !@#$%</div>
                  <div class="req-toggle__status">
                    {{ policyForm.requireSpecialChar ? '已启用' : '已禁用' }}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 区块二：密码过期策略 -->
        <div class="policy-section">
          <div class="policy-section__head">
            <h3 class="policy-section__title">
              <i class="el-icon-time" style="margin-right: 8px;"></i>
              密码过期策略
            </h3>
            <span class="policy-section__hint">控制密码的有效期限，强制定期更换</span>
          </div>

          <div class="expiry-card">
            <div class="expiry-card__label">
              密码有效期
              <el-tooltip content="设置为 0 则永不过期，不建议生产环境使用" placement="top">
                <i class="el-icon-question tooltip-icon"></i>
              </el-tooltip>
            </div>
            <div class="expiry-card__control">
              <el-input-number 
                v-model.number="policyForm.expiryDays" 
                :min="0" 
                :max="3650" 
                :step="7"
                size="small"
                controls-position="right"
                class="expiry-control"
              />
              <span class="expiry-card__unit">天</span>
            </div>
            <div class="expiry-card__options">
              <el-radio-group v-model="expiryMode" size="small">
                <el-radio-button :label="'custom'">自定义天数</el-radio-button>
                <el-radio-button :label="'none'">永不过期</el-radio-button>
                <el-radio-button label="30">30 天</el-radio-button>
                <el-radio-button label="60">60 天</el-radio-button>
                <el-radio-button label="90">90 天</el-radio-button>
                <el-radio-button label="180">180 天</el-radio-button>
              </el-radio-group>
            </div>
          </div>
        </div>

        <!-- 区块三：安全策略 -->
        <div class="policy-section">
          <div class="policy-section__head">
            <h3 class="policy-section__title">
              <i class="el-icon-connection" style="margin-right: 8px;"></i>
              安全策略
            </h3>
            <span class="policy-section__hint">防止密码复用和暴力破解攻击</span>
          </div>

          <div class="security-cards">
            <!-- 历史密码 -->
            <div class="history-card">
              <div class="history-card__icon">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                  <circle cx="12" cy="12" r="9"/>
                </svg>
              </div>
              <div class="history-card__content">
                <div class="history-card__label">禁止使用最近 N 次密码</div>
                <div class="history-card__control">
                  <el-input-number 
                    v-model.number="policyForm.passwordHistory" 
                    :min="0" 
                    :max="50" 
                    :step="1"
                    size="small"
                    controls-position="right"
                    class="history-control"
                  />
                  <span class="history-card__unit">次</span>
                </div>
                <div class="history-card__tips">
                  用户创建新密码时，系统会检查是否在近{{ policyForm.passwordHistory }}次内使用过
                </div>
              </div>
            </div>

            <!-- 登录失败锁定 -->
            <div class="lockout-card">
              <div class="lockout-card__icon">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                  <path d="M7 11V7a5 5 0 0110 0v4"/>
                </svg>
              </div>
              <div class="lockout-card__content">
                <div class="lockout-card__label">最大连续失败登录次数</div>
                <div class="lockout-card__control">
                  <el-input-number 
                    v-model.number="policyForm.maxLoginAttempts" 
                    :min="0" 
                    :max="100" 
                    :step="1"
                    size="small"
                    controls-position="right"
                    class="lockout-control"
                  />
                  <span class="lockout-card__unit">次</span>
                </div>
                <div class="lockout-card__tips">
                  超出限制后账号自动锁定，需管理员解锁或等待锁定时间结束
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 区块四：密码验证测试器（新增） -->
        <div class="policy-section">
          <div class="policy-section__head">
            <h3 class="policy-section__title">
              <i class="el-icon-document-checked" style="margin-right: 8px;"></i>
              密码验证测试
            </h3>
            <span class="policy-section__hint">输入密码实时验证是否符合当前策略</span>
          </div>

          <div class="validator-card">
            <el-form :inline="true" :model="validatorForm">
              <el-form-item label="密码强度测试">
                <el-input 
                  v-model="validatorForm.testPassword" 
                  placeholder="输入要测试的密码"
                  type="password"
                  show-password
                  size="medium"
                  class="validator-input"
                  clearable
                  @input="handleValidateInput"
                ></el-input>
              </el-form-item>
              <el-form-item>
                <el-button 
                  type="primary" 
                  size="medium" 
                  plain
                  icon="el-icon-video-play"
                  @click="handleValidate"
                >
                  验证密码
                </el-button>
                <el-button size="medium" @click="validatorForm.testPassword = ''">清除</el-button>
              </el-form-item>
            </el-form>

            <!-- 验证结果展示 -->
            <div v-if="validationResult" class="validation-result">
              <el-alert
                :type="validationResult.valid ? 'success' : 'error'"
                :title="validationResult.valid ? '密码符合要求' : '密码不符合要求'"
                :description="validationResult.message"
                show-icon
                :closable="false"
                class="validation-alert"
              ></el-alert>
            </div>

            <!-- 密码强度指示器 -->
            <div v-if="validatorForm.testPassword" class="strength-meter">
              <div class="strength-meter__label">密码强度</div>
              <div class="strength-meter__bar">
                <div class="strength-meter__fill" :style="{ width: strengthWidth + '%', background: strengthColor }"></div>
              </div>
              <div class="strength-meter__text" :color="strengthColor">
                {{ strengthText }}
              </div>
            </div>

            <!-- 实时校验规则 -->
            <div class="rules-list" v-if="validatorForm.testPassword">
              <div 
                class="rule-item" 
                :class="{ 'is-passed': validationRules.length }"
                v-for="(rule, index) in validationRules" 
                :key="index"
              >
                <i :class="rule.passed ? 'el-icon-circle-check' : 'el-icon-close'"></i>
                <span>{{ rule.text }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 底部操作栏 -->
      <div class="policy-footer">
        <el-button size="medium" @click="resetPolicy">重置</el-button>
        <el-button
          type="primary"
          size="medium"
          :loading="saving"
          @click="savePolicy"
        >
          {{ saving ? '保存中...' : '保存配置' }}
        </el-button>
      </div>
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
      validatorForm: {
        testPassword: ''
      },
      saving: false,
      validating: false,
      validationResult: null,
      validationRules: [],
      strengthWidth: 0,
      strengthColor: '#ccc',
      strengthText: '',
      expiryMode: 'custom'
    }
  },
  watch: {
    expiryMode(val) {
      if (val === 'none') {
        this.policyForm.expiryDays = 0
      } else if (val !== 'custom') {
        this.policyForm.expiryDays = parseInt(val)
      }
    }
  },
  created() {
    this.fetchPolicy()
  },
  methods: {
    toggleReq(field) {
      this.policyForm[field] = !this.policyForm[field]
    },
    
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
          // 初始化过期模式
          if (p.expiryDays === 0) {
            this.expiryMode = 'none'
          } else if ([30, 60, 90, 180].includes(p.expiryDays)) {
            this.expiryMode = String(p.expiryDays)
          } else {
            this.expiryMode = 'custom'
          }
        }
      } catch (e) {
        console.error(e)
      }
    },
    
    handleValidate() {
      if (!this.validatorForm.testPassword) {
        this.$message.warning('请输入要测试的密码')
        return
      }
      
      this.validating = true
      
      const password = this.validatorForm.testPassword
      const rules = []
      
      // 检查长度
      rules.push({
        text: `长度至少${this.policyForm.minLength}位（当前 ${password.length} 位）`,
        passed: password.length >= this.policyForm.minLength
      })
      
      // 检查大写字母
      if (this.policyForm.requireUppercase) {
        rules.push({
          text: '包含大写字母 A-Z',
          passed: /[A-Z]/.test(password)
        })
      } else {
        rules.push({ text: '不需要大写字母', passed: true })
      }
      
      // 检查小写字母
      if (this.policyForm.requireLowercase) {
        rules.push({
          text: '包含小写字母 a-z',
          passed: /[a-z]/.test(password)
        })
      } else {
        rules.push({ text: '不需要小写字母', passed: true })
      }
      
      // 检查数字
      if (this.policyForm.requireDigit) {
        rules.push({
          text: '包含数字 0-9',
          passed: /[0-9]/.test(password)
        })
      } else {
        rules.push({ text: '不需要数字', passed: true })
      }
      
      // 检查特殊字符
      if (this.policyForm.requireSpecialChar) {
        rules.push({
          text: '包含特殊符号 !@#$%^&* 等',
          passed: /[!@#$%^&*]/.test(password)
        })
      } else {
        rules.push({ text: '不需要特殊字符', passed: true })
      }
      
      this.validationRules = rules
      const allPassed = rules.every(r => r.passed)
      
      if (allPassed) {
        this.validationResult = {
          valid: true,
          message: '✅ 密码符合所有复杂度要求，强度良好！'
        }
        this.strengthWidth = 100
        this.strengthColor = '#67C23A'
        this.strengthText = '强密码'
      } else {
        const failedCount = rules.filter(r => !r.passed).length
        this.validationResult = {
          valid: false,
          message: '❌ 还有' + failedCount + '项不符合要求，请根据提示调整密码'
        }
        this.strengthWidth = Math.max(0, (rules.filter(r => r.passed).length / rules.length) * 100 - 20)
        this.strengthColor = '#F56C6C'
        this.strengthText = '需改进'
      }
      
      this.validating = false
    },
    
    handleValidateInput() {
      if (!this.validatorForm.testPassword) {
        this.strengthWidth = 0
        this.strengthColor = '#ccc'
        this.strengthText = ''
        return
      }
      
      const password = this.validatorForm.testPassword
      let score = 0
      
      if (password.length >= this.policyForm.minLength) score++
      if (password.length >= 16) score++
      if (/[A-Z]/.test(password)) score++
      if (/[a-z]/.test(password)) score++
      if (/[0-9]/.test(password)) score++
      if (/[!@#$%^&*]/.test(password)) score++
      
      if (score >= 6) {
        this.strengthWidth = 100
        this.strengthColor = '#67C23A'
        this.strengthText = '强密码 🟢'
      } else if (score >= 4) {
        this.strengthWidth = 60
        this.strengthColor = '#E6A23C'
        this.strengthText = '中等密码 🟡'
      } else if (score >= 2) {
        this.strengthWidth = 30
        this.strengthColor = '#E6A23C'
        this.strengthText = '弱密码 🔴'
      } else {
        this.strengthWidth = 10
        this.strengthColor = '#F56C6C'
        this.strengthText = '非常弱 🔴'
      }
    },
    
    async savePolicy() {
      this.saving = true
      try {
        const res = await updatePasswordPolicy(this.policyForm)
        if (res.code === 200) {
          this.$message.success('密码策略已更新')
          this.validatorForm.testPassword = ''
          this.validationResult = null
          this.validationRules = []
        } else {
          this.$message.error(res.message || '保存失败')
        }
      } catch (e) {
        console.error('保存密码策略失败:', e)
        this.$message.error('保存密码策略失败')
      } finally {
        this.saving = false
      }
    },
    
    resetPolicy() {
      this.$confirm('确定要重置为上次保存的配置吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.fetchPolicy()
        this.$message.success('已重置到上次保存的配置')
        this.validatorForm.testPassword = ''
        this.validationResult = null
        this.validationRules = []
      }).catch(() => {})
    }
  }
}
</script>

<style scoped>
.password-policy-container {
  padding: 20px;
}

.policy-config-card {
  max-width: 960px;
  margin: 0 auto;
  border-radius: 12px;
}

.policy-config-card >>> .el-card__header {
  padding: 0;
  border-bottom: none;
}

/* ===== 政策信息头部 ===== */
.policy-header {
  display: flex;
  align-items: center;
  padding: 20px 24px;
  background: linear-gradient(135deg, #f0f6ff 0%, #f8fafc 100%);
  border-bottom: 1px solid #e3edfb;
  border-radius: 12px 12px 0 0;
}

.policy-header__icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 52px;
  height: 52px;
  border-radius: 12px;
  margin-right: 16px;
  flex-shrink: 0;
}

.policy-header__icon i {
  font-size: 26px;
  color: #fff;
}

.policy-header__icon.is-primary {
  background: linear-gradient(135deg, #409EFF, #337ecc);
}

.policy-header__info {
  flex: 1;
}

.policy-header__title {
  display: flex;
  align-items: center;
  font-size: 17px;
  font-weight: 600;
  color: #1f2d3d;
  gap: 8px;
}

.policy-header__tag {
  margin-left: 2px;
}

.policy-header__desc {
  margin-top: 6px;
  font-size: 13px;
  color: #7a8ba3;
}

.policy-content {
  padding: 4px 20px;
}

/* ===== 设置区块 ===== */
.policy-section {
  padding: 20px 4px;
  border-bottom: 1px solid #eef1f6;
}

.policy-section:last-child {
  border-bottom: none;
}

.policy-section__head {
  display: flex;
  align-items: baseline;
  margin-bottom: 16px;
}

.policy-section__title svg {
  color: #409EFF;
}

.policy-section__title i {
  margin-right: 8px;
}

.policy-section__title > span {
  margin-left: 8px;
}

.policy-section__hint {
  margin-left: 12px;
  font-size: 12px;
  color: #98a6b8;
}

/* ===== 配置卡片组 ===== */
.complexity-cards {
  display: flex;
  gap: 16px;
  align-items: stretch;
}

.config-card {
  flex: 0 0 200px;
  padding: 16px;
  background: #fafbfd;
  border: 1.5px solid #e4e9f0;
  border-radius: 10px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  align-items: center;
  justify-content: center;
}

.config-card__label {
  font-size: 13px;
  font-weight: 600;
  color: #1f2d3d;
  text-align: center;
}

.config-card__control,
.expiry-card__control {
  display: flex;
  align-items: center;
  gap: 8px;
}

.config-card__unit {
  font-size: 12px;
  color: #98a6b8;
}

.config-card__tips {
  flex: 1;
}

/* ===== 字符要求网格 ===== */
.requirements-grid {
  flex: 1;
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.req-toggle {
  display: flex;
  align-items: center;
  padding: 14px;
  background: #fafbfd;
  border: 1.5px solid #e4e9f0;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s ease;
  min-height: 80px;
}

.req-toggle:hover {
  border-color: #b3d4fc;
  background: #f5f9ff;
}

.req-toggle.is-active {
  border-color: #409EFF;
  background: #f0f7ff;
  box-shadow: 0 0 0 1px #409eff inset;
}

.req-toggle__icon {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  background: #eef3fa;
  color: #5a7ca6;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 12px;
  flex-shrink: 0;
  transition: all 0.2s ease;
}

.req-toggle.is-active .req-toggle__icon {
  background: #409EFF;
  color: #fff;
}

.req-toggle__icon i {
  font-size: 18px;
}

.req-toggle__text {
  flex: 1;
  overflow: hidden;
}

.req-toggle__title {
  font-size: 14px;
  font-weight: 600;
  color: #1f2d3d;
  margin-bottom: 4px;
}

.req-toggle__status {
  font-size: 11px;
  color: #98a6b8;
}

/* ===== 密码过期卡片 ===== */
.expiry-card {
  padding: 16px;
  background: #fafbfd;
  border: 1.5px solid #e4e9f0;
  border-radius: 10px;
  max-width: 500px;
}

.expiry-card__label {
  font-size: 14px;
  font-weight: 600;
  color: #1f2d3d;
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 12px;
}

.tooltip-icon {
  font-size: 13px;
  color: #98a6b8;
  cursor: help;
}

.expiry-card__control {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
}

.expiry-control,
.length-control {
  width: 120px;
}

.expiry-card__unit {
  font-size: 12px;
  color: #98a6b8;
}

.expiry-card__options {
  margin-top: 8px;
}

.expiry-card__options >>> .el-radio-button__original-radio:checked + .el-radio-button__inner {
  border-color: #409EFF;
  background: #409EFF;
  color: #fff;
}

/* ===== 安全卡片 ===== */
.security-cards {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.history-card,
.lockout-card {
  display: flex;
  gap: 12px;
  padding: 16px;
  background: #fafbfd;
  border: 1.5px solid #e4e9f0;
  border-radius: 10px;
  align-items: center;
}

.history-card__icon,
.lockout-card__icon {
  width: 42px;
  height: 42px;
  border-radius: 10px;
  background: linear-gradient(135deg, #6f7ad3, #5560c0);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.history-card__content,
.lockout-card__content {
  flex: 1;
  min-width: 0;
}

.history-card__label,
.lockout-card__label {
  font-size: 13px;
  font-weight: 600;
  color: #1f2d3d;
  margin-bottom: 6px;
}

.history-card__control,
.lockout-card__control {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.history-control,
.lockout-control {
  width: 120px;
}

.history-card__unit,
.lockout-card__unit {
  font-size: 12px;
  color: #98a6b8;
}

.history-card__tips,
.lockout-card__tips {
  font-size: 11px;
  color: #98a6b8;
  line-height: 1.4;
}

/* ===== 验证器卡片 ===== */
.validator-card {
  padding: 16px;
  background: #fafbfd;
  border: 1.5px solid #e4e9f0;
  border-radius: 10px;
}

.validator-card >>> .el-form-item {
  margin-bottom: 12px;
}

.validator-input {
  width: 300px;
}

.validation-result {
  margin-top: 16px;
}

.validation-alert {
  width: 100%;
}

/* ===== 密码强度指示器 ===== */
.strength-meter {
  margin-top: 16px;
  padding: 12px;
  background: #f5f7fa;
  border-radius: 8px;
}

.strength-meter__label {
  font-size: 12px;
  color: #7a8ba3;
  margin-bottom: 8px;
}

.strength-meter__bar {
  height: 8px;
  background: #e3e8f0;
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 4px;
}

.strength-meter__fill {
  height: 100%;
  border-radius: 4px;
  transition: all 0.3s ease;
}

.strength-meter__text {
  font-size: 12px;
  font-weight: 600;
}

/* ===== 规则列表 ===== */
.rules-list {
  margin-top: 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.rule-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: #7a8ba3;
  padding: 6px 8px;
  background: #f5f7fa;
  border-radius: 6px;
}

.rule-item.is-passed {
  background: #f0f9ff;
  color: #409eff;
}

.rule-item.is-passed i {
  color: #67C23A;
}

.rule-item i {
  font-size: 14px;
}

/* ===== 底部操作栏 ===== */
.policy-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px 20px;
  border-top: 1px solid #eef1f6;
  background: #fafbfd;
  border-radius: 0 0 12px 12px;
}

.policy-footer .el-button {
  min-width: 110px;
}
</style>
