<template>
  <div>
    <el-card class="policy-config-card" shadow="never">
      <!-- 政策信息头部 -->
      <div slot="header" class="policy-header">
        <div class="policy-header__icon is-primary">
          <svg viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg" style="width: 26px; height: 26px;">
            <path d="M19 20C22.866 20 26 16.866 26 13C26 9.13401 22.866 6 19 6C15.134 6 12 9.13401 12 13C12 16.866 15.134 20 19 20Z" fill="white" fill-opacity="0.2" stroke="currentColor" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M32.6077 7C34.6405 8.2249 36.0001 10.4537 36.0001 13C36.0001 15.5463 34.6405 17.7751 32.6077 19" stroke="currentColor" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M4 40.8V42H34V40.8C34 36.3196 34 34.0794 33.1281 32.3681C32.3611 30.8628 31.1372 29.6389 29.6319 28.8719C27.9206 28 25.6804 28 21.2 28H16.8C12.3196 28 10.0794 28 8.36808 28.8719C6.86278 29.6389 5.63893 30.8628 4.87195 32.3681C4 34.0794 4 36.3196 4 40.8Z" fill="white" fill-opacity="0.2" stroke="currentColor" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M43.9999 42.0001V40.8001C43.9999 36.3197 43.9999 34.0795 43.128 32.3682C42.361 30.8629 41.1371 29.6391 39.6318 28.8721" stroke="currentColor" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </div>
        <div class="policy-header__info">
          <div class="policy-header__title">LDAP 管理</div>
          <div class="policy-header__desc">配置 LDAP 服务器连接参数，实现统一身份认证</div>
        </div>
      </div>

      <!-- 设置内容区 -->
      <div class="policy-content">
        
        <!-- 区块一：基本连接配置 -->
        <div class="policy-section">
          <div class="policy-section__head">
            <h3 class="policy-section__title">
              <svg viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg" style="width: 18px; height: 18px; color: #409EFF;">
                <path d="M41.5 10H35.5" stroke="currentColor" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
                <path d="M27.5 6V14" stroke="currentColor" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
                <path d="M27.5 10L5.5 10" stroke="currentColor" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
                <path d="M13.5 24H5.5" stroke="currentColor" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
                <path d="M21.5 20V28" stroke="currentColor" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
                <path d="M43.5 24H21.5" stroke="currentColor" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
                <path d="M41.5 38H35.5" stroke="currentColor" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
                <path d="M27.5 34V42" stroke="currentColor" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
                <path d="M27.5 38H5.5" stroke="currentColor" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
              <span style="margin-left: 8px;">基本连接配置</span>
            </h3>
            <span class="policy-section__hint">填写 LDAP 服务器的基础连接信息</span>
          </div>

          <div class="config-grid">
            <div class="config-item">
              <label class="config-label">LDAP 服务器地址</label>
              <el-input 
                v-model="ldapForm.server" 
                placeholder="ldaps://10.60.254.252:636"
                clearable
                :disabled="loading"
              ></el-input>
            </div>
            
            <div class="config-item">
              <label class="config-label">Base DN</label>
              <el-input 
                v-model="ldapForm.base_dn" 
                placeholder="dc=hot,dc=local"
                clearable
                :disabled="loading"
              ></el-input>
            </div>
            
            <div class="config-item full-width">
              <label class="config-label">用户过滤模板</label>
              <el-input 
                v-model="ldapForm.user_filter" 
                placeholder="(sAMAccountName=%s)"
                clearable
                :disabled="loading"
              >
                <template #append>
                  <el-tooltip content="%s 会被替换为实际用户名" placement="top">
                    <i class="el-icon-question" style="cursor: pointer; color: #909399;"></i>
                  </el-tooltip>
                </template>
              </el-input>
            </div>
          </div>
        </div>

        <!-- 区块二：安全连接选项 -->
        <div class="policy-section">
          <div class="policy-section__head">
            <h3 class="policy-section__title">
              <i class="el-icon-lock" style="margin-right: 8px;"></i>
              <span>安全连接选项</span>
            </h3>
            <span class="policy-section__hint">配置 TLS 证书验证与加密传输</span>
          </div>

          <div class="security-config">
            <el-checkbox v-model="ldapForm.use_tls" :disabled="loading" class="tls-checkbox">
              使用 TLS 加密连接（LDAPS）
            </el-checkbox>

            <el-checkbox v-model="ldapForm.insecure" :disabled="!ldapForm.use_tls || loading" class="insecure-checkbox">
              跳过证书验证（仅用于测试环境）
            </el-checkbox>
          </div>

          <!-- CA 证书文件上传 -->
          <div class="cert-upload-section" v-if="ldapForm.use_tls">
            <label class="upload-label">CA 证书文件</label>
            <el-upload
              ref="certificateUpload"
              action="#"
              :auto-upload="false"
              :limit="1"
              accept=".crt,.cer,.pem"
              :on-change="handleCertificateChange"
              :file-list="certificateList"
              :disabled="loading"
            >
              <el-button size="small">选择证书文件</el-button>
            </el-upload>
            
            <!-- 新选择的证书预览（已读取但未保存） -->
            <el-alert 
              v-if="ldapFile && isCertModified"
              title="证书预览" 
              type="info" 
              :closable="false" 
              show-icon
              class="alert-info"
            >
              <pre style="max-height: 200px; overflow-y: auto;">{{ ldapFile.content }}</pre>
            </el-alert>
            
            <!-- 已保存的证书提示 -->
            <el-alert 
              v-else-if="ldapForm.cert_base64 && !isCertModified" 
              title="当前已上传 CA 证书" 
              type="success" 
              :closable="false" 
              show-icon
              class="alert-success"
            >
              <span slot="content">
                证书已准备就绪，请检查其他配置项后保存或测试连接
              </span>
            </el-alert>
          </div>
        </div>

        <!-- 区块三：绑定账户配置 -->
        <div class="policy-section">
          <div class="policy-section__head">
            <h3 class="policy-section__title">
              <i class="el-icon-user" style="margin-right: 8px;"></i>
              <span>绑定账户配置</span>
            </h3>
            <span class="policy-section__hint">LDAP 管理员账号（用于搜索用户的绑定账户）</span>
          </div>

          <div class="config-grid">
            <div class="config-item">
              <label class="config-label">绑定 DN (Username)</label>
              <el-input 
                v-model="ldapForm.username" 
                placeholder="ylw@hot.local"
                clearable
                :disabled="loading"
              ></el-input>
            </div>
            
            <div class="config-item">
              <label class="config-label">密码</label>
              <el-input 
                v-model="ldapForm.password" 
                type="password" 
                placeholder="请输入管理员密码（留空表示不修改）"
                show-password
                clearable
                :disabled="loading"
              ></el-input>
            </div>
          </div>
        </div>

        <!-- 底部操作栏 -->
        <div class="policy-footer">
          <el-button size="medium" @click="handleReset">重置</el-button>
          <el-button 
            type="primary" 
            size="medium" 
            :loading="testLoading"
            :disabled="loading || testDisabled"
            @click="handleTest"
          >
            {{ testLoading ? '测试中...' : '测试连接' }}
          </el-button>
          <el-button
            type="primary"
            size="medium"
            :loading="submitLoading"
            @click="handleSubmit"
          >
            {{ submitLoading ? '保存中...' : '保存配置' }}
          </el-button>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script>
import { getLDAPConfig, saveLDAPConfig, testLDAPConnection } from '@/api/settings'
import { rsaEncrypt } from '@/utils/encrypt'

export default {
  name: 'LDAPManagement',
  data() {
    return {
      // 加载中状态
      loading: false,
      submitLoading: false,
      testLoading: false,
      
      // 表单数据
      ldapForm: {
        server: '',
        base_dn: '',
        use_tls: false,
        insecure: true,
        user_filter: '(sAMAccountName=%s)',
        username: '',
        password: '',
        cert_base64: ''
      },
      
      // 证书相关文件
      ldapFile: null,
      certificateList: [],
      isCertModified: false,
      
      // 校验规则
      rules: {
        server: [
          { required: true, message: '请输入 LDAP 服务器地址', trigger: 'blur' },
          { pattern: /^ldaps?:\/\/.+/, message: '格式应为 ldaps://host:port 或 ldap://host:port', trigger: 'blur' }
        ],
        base_dn: [
          { required: true, message: '请输入 Base DN', trigger: 'blur' }
        ],
        user_filter: [
          { required: true, message: '请输入用户过滤模板', trigger: 'blur' }
        ]
      },
      
      // 测试连接是否被调用过
      testCalled: false
    }
  },
  computed: {
    // 计算是否已有表单提交过
    isDirty() {
      return this.testCalled || this.isCertModified || this.ldapForm.password !== ''
    },
    // 判断测试按钮是否禁用
    testDisabled() {
      return !this.ldapForm.server || !this.ldapForm.base_dn || (this.ldapForm.use_tls && !this.ldapForm.cert_base64)
    }
  },
  mounted() {
    this.fetchData()
  },
  methods: {
    // 获取 LDAP 配置数据
    async fetchData() {
      this.loading = true
      try {
        const response = await getLDAPConfig()
        if (response.code === 200) {
          const config = response.data
          this.ldapForm = {
            server: config.server || '',
            base_dn: config.base_dn || '',
            use_tls: config.use_tls || false,
            // 注意：不能用 || true，否则 false 会被默认值覆盖
            insecure: config.insecure !== undefined ? config.insecure : true,
            user_filter: config.user_filter || '(sAMAccountName=%s)',
            username: response.username || '', // 取顶层解密后的明文，data 内为加密原值
            password: '', // 出于安全考虑，返回时不显示密码
            cert_base64: config.cert_base64 || ''
          }
          
          // 如果有证书，保留原文件名（如果后端返回了的话）
          if (config.cert_base64) {
            // 这里可以扩展：如果后端 API 返回了 cert_filename，优先使用原文件名
            this.certificateList = [{ name: config.cert_filename || 'ca.crt', url: '#' }]
          }
        } else {
          this.$message.error(response.message || '获取 LDAP 配置失败')
        }
      } catch (error) {
        // 错误提示已由响应拦截器统一弹出，避免重复提示
        console.error('加载 LDAP 配置失败:', error)
      } finally {
        this.loading = false
      }
    },
    
    // 处理证书文件变化
    handleCertificateChange(file) {
      const rawFile = file.raw
      
      if (!rawFile) return
      
      // 检查文件类型
      const fileName = rawFile.name.toLowerCase()
      const isValidType = fileName.endsWith('.crt') || fileName.endsWith('.cer') || fileName.endsWith('.pem')
      
      if (!isValidType) {
        this.$message.warning('仅支持 .crt, .cer, .pem 格式的证书文件')
        return
      }
      
      // 以二进制方式读取文件内容（避免 DER 二进制证书按文本解码产生乱码/损坏）
      const reader = new FileReader()
      reader.onload = (e) => {
        const bytes = new Uint8Array(e.target.result)

        // 分块构建原始字节的 binary string（避免大文件栈溢出）
        let binary = ''
        const CHUNK = 0x8000
        for (let i = 0; i < bytes.length; i += CHUNK) {
          binary += String.fromCharCode.apply(null, bytes.subarray(i, i + CHUNK))
        }
        // 与原始字节完全一致的 base64
        const rawBase64 = btoa(binary)

        // 格式判定：PEM 文本直接预览；DER 二进制转换为 PEM 文本预览
        let previewText
        if (binary.startsWith('-----BEGIN')) {
          previewText = new TextDecoder().decode(bytes)
        } else {
          previewText = '-----BEGIN CERTIFICATE-----\n' +
            rawBase64.match(/.{1,64}/g).join('\n') +
            '\n-----END CERTIFICATE-----'
        }

        this.ldapFile = { name: rawFile.name, content: previewText }
        this.ldapForm.cert_base64 = previewText // 保存与预览一致的明文 PEM 文本（后端兼容解析）
        this.isCertModified = true
        
        // 添加到上传列表
        this.certificateList = [{ name: rawFile.name, url: '#' }]
      }
      reader.readAsArrayBuffer(rawFile)
    },
    
    // 提交表单
    async handleSubmit() {
      // 如果未测试过，先确认是否需要测试提示
      if (!this.testCalled) {
        const needConfirm = await this.confirmRequired()
        if (!needConfirm) return
      }
      
      // 手动校验表单字段
      if (!this.ldapForm.server || !this.ldapForm.base_dn) {
        this.$message.warning('请填写必要的配置项后再保存')
        return
      }
      
      this.submitLoading = true
      try {
      // 密码使用与登录平台一致的 RSA 加密传输，留空表示不修改
      let encryptedPwd = ''
      if (this.ldapForm.password) {
        encryptedPwd = rsaEncrypt(this.ldapForm.password)
        if (!encryptedPwd) {
          this.$message.error('密码加密失败，请重试')
          this.submitLoading = false
          return
        }
      }
      const response = await saveLDAPConfig({
          server: this.ldapForm.server,
          base_dn: this.ldapForm.base_dn,
          use_tls: this.ldapForm.use_tls,
          insecure: this.ldapForm.insecure,
          user_filter: this.ldapForm.user_filter,
          username: this.ldapForm.username,
          password: encryptedPwd,
          cert_base64: this.ldapForm.cert_base64,
          cert_filename: this.certificateList.length > 0 ? this.certificateList[0].name : '' // 发送文件名
        })
        
        if (response.code === 200) {
          this.$message.success('LDAP 配置保存成功')
          this.fetchData()
          this.testCalled = true
        } else {
          this.$message.error(response.message || '保存失败')
        }
      } catch (error) {
        // 错误提示已由响应拦截器统一弹出，避免重复提示
        console.error('保存 LDAP 配置失败:', error)
      } finally {
        this.submitLoading = false
      }
    },
    
    // 确认是否需要提示
    confirmRequired() {
      return new Promise((resolve) => {
        if (this.testCalled) {
          resolve(true)
          return
        }
        
        this.$confirm(
          '尚未进行连接测试，直接保存可能导致配置无法使用。是否继续？',
          '未测试连接提示',
          {
            confirmButtonText: '确定保存',
            cancelButtonText: '取消测试',
            type: 'warning'
          }
        ).then(() => {
          resolve(true)
        }).catch(() => {
          resolve(false)
        })
      })
    },
    
    // 测试连接
    async handleTest() {
      if (this.testDisabled) {
        this.$message.warning('请填写必要的配置项后再测试')
        return
      }
      
      this.testLoading = true
      try {
        // 密码使用与登录平台一致的 RSA 加密传输
        let encryptedPwd = ''
        if (this.ldapForm.password) {
          encryptedPwd = rsaEncrypt(this.ldapForm.password)
          if (!encryptedPwd) {
            this.$message.error('密码加密失败，请重试')
            this.testLoading = false
            return
          }
        }
        const response = await testLDAPConnection({
          server: this.ldapForm.server,
          base_dn: this.ldapForm.base_dn,
          use_tls: this.ldapForm.use_tls,
          insecure: this.ldapForm.insecure,
          username: this.ldapForm.username,
          password: encryptedPwd,
          cert_base64: this.ldapForm.cert_base64,
          user_filter: this.ldapForm.user_filter
        })
        
        if (response.code === 200) {
          this.$message.success('✓ ' + response.message)
          this.testCalled = true
        } else {
          this.$message.error('✗ ' + response.message)
        }
      } catch (error) {
        // 错误提示已由响应拦截器统一弹出，避免重复提示
        console.error('测试连接失败:', error)
      } finally {
        this.testLoading = false
      }
    },
    
    // 重置表单
    handleReset() {
      if (this.isDirty) {
        this.$confirm('重置将放弃所有更改，是否继续？', '重置确认', {
          confirmButtonText: '确定重置',
          cancelButtonText: '取消',
          type: 'warning'
        }).then(() => {
          this.fetchData()
        }).catch(() => {})
      } else {
        this.fetchData()
      }
    }
  }
}
</script>

<style scoped>
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

.policy-header__icon svg {
  color: #fff;
  width: 26px;
  height: 26px;
}

.policy-header__icon.is-primary {
  background: linear-gradient(135deg, #409EFF, #337ecc);
}

.policy-header__info {
  flex: 1;
}

.policy-header__title {
  font-size: 17px;
  font-weight: 600;
  color: #1f2d3d;
  margin-bottom: 6px;
}

.policy-header__desc {
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

/* ===== 配置网格布局 ===== */
.config-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.config-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.config-item.full-width {
  grid-column: 1 / -1;
}

.config-label {
  font-size: 14px;
  font-weight: 600;
  color: #1f2d3d;
}

/* ===== 安全配置选项 ===== */
.security-config {
  padding: 16px;
  background: #fafbfd;
  border: 1.5px solid #e4e9f0;
  border-radius: 10px;
}

.tls-checkbox,
.insecure-checkbox {
  font-size: 15px;
  font-weight: 500;
  margin-bottom: 12px;
}

/* ===== 证书上传区域 ===== */
.cert-upload-section {
  padding: 16px;
  background: #fafbfd;
  border: 1.5px solid #e4e9f0;
  border-radius: 10px;
  margin-top: 12px;
}

.upload-label {
  font-size: 14px;
  font-weight: 600;
  color: #1f2d3d;
  margin-bottom: 12px;
  display: block;
}

.alert-success,
.alert-info {
  margin-top: 12px;
}

/* ===== 底部操作栏 ===== */
.policy-footer {
  padding: 20px;
  text-align: right;
  border-top: 1px solid #eef1f6;
}
</style>
