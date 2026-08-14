<template>
  <el-card class="module-config-card" shadow="never">
    <!-- 模块信息头部（card header） -->
    <div slot="header" class="config-header">
      <div class="config-header__icon" :class="'is-' + moduleName">
        <i :class="moduleName === 'chinaunicom' ? 'el-icon-connection' : 'el-icon-collection-tag'"></i>
      </div>
      <div class="config-header__info">
        <div class="config-header__title">
          {{ moduleTitle }}
          <el-tag
            size="mini"
            :type="currentConfig.loginType === 'local' ? 'success' : 'primary'"
            effect="light"
            class="config-header__tag"
          >
            {{ currentConfig.loginType === 'local' ? '本地登录' : 'LDAP 登录' }}
          </el-tag>
          <el-tag
            v-if="showDualAuth"
            size="mini"
            :type="currentConfig.dualAuthEnabled ? 'warning' : 'info'"
            effect="light"
          >
            双控{{ currentConfig.dualAuthEnabled ? '已开启' : '已关闭' }}
          </el-tag>
        </div>
        <div class="config-header__desc">{{ moduleDesc }}</div>
      </div>
    </div>

    <!-- 设置区块容器（card body） -->
    <div class="config-card-body">
    <!-- 设置区块：登录方式 -->
    <div class="config-section">
      <div class="config-section__head">
        <h3 class="config-section__title">登录方式</h3>
        <span class="config-section__hint">选择该模块在 SFTP 登录页使用的认证方式</span>
      </div>
      <div class="login-type-cards">
        <div
          class="login-type-card"
          :class="{ 'is-active': currentConfig.loginType === 'local' }"
          @click="currentConfig.loginType = 'local'"
        >
          <div class="login-type-card__icon"><i class="el-icon-user"></i></div>
          <div class="login-type-card__body">
            <div class="login-type-card__title">本地登录</div>
            <div class="login-type-card__desc">使用平台本地账号与密码登录</div>
          </div>
          <div class="login-type-card__check">
            <i class="el-icon-check" v-if="currentConfig.loginType === 'local'"></i>
          </div>
        </div>
        <div
          class="login-type-card"
          :class="{ 'is-active': currentConfig.loginType === 'ldap' }"
          @click="currentConfig.loginType = 'ldap'"
        >
          <div class="login-type-card__icon">
            <!-- LDAP 目录结构图标（stroke 使用 currentColor 跟随选中态） -->
            <svg class="ldap-icon" viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
              <circle cx="14" cy="29" r="5" stroke="currentColor" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
              <circle cx="34" cy="29" r="5" stroke="currentColor" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
              <circle cx="24" cy="9" r="5" stroke="currentColor" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M24 44C24 38.4772 19.5228 34 14 34C8.47715 34 4 38.4772 4 44" stroke="currentColor" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M44 44C44 38.4772 39.5228 34 34 34C28.4772 34 24 38.4772 24 44" stroke="currentColor" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M34 24C34 18.4772 29.5228 14 24 14C18.4772 14 14 18.4772 14 24" stroke="currentColor" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </div>
          <div class="login-type-card__body">
            <div class="login-type-card__title">LDAP 登录</div>
            <div class="login-type-card__desc">使用域账号（LDAP 安全组）登录</div>
          </div>
          <div class="login-type-card__check">
            <i class="el-icon-check" v-if="currentConfig.loginType === 'ldap'"></i>
          </div>
        </div>
      </div>
    </div>

    <!-- 设置区块：可登录角色 -->
    <div class="config-section">
      <div class="config-section__head">
        <h3 class="config-section__title">可登录角色</h3>
        <span class="config-section__hint">仅列表中角色可登录该模块，未配置则全部拒绝</span>
      </div>
      <div class="role-select-wrap">
        <el-select
          v-model="currentConfig.enabledRolesArray"
          multiple
          filterable
          collapse-tags
          placeholder="请选择允许登录的角色"
          class="role-select"
          size="medium"
        >
          <el-option
            v-for="role in roleList"
            :key="role.ID"
            :label="role.name"
            :value="role.ID"
          >
            <span>{{ role.name }}</span>
            <span class="role-option-desc">{{ role.description }}</span>
          </el-option>
        </el-select>
        <div class="role-select__tools">
          <el-button type="text" size="small" @click="selectAllRoles">全选</el-button>
          <el-button type="text" size="small" @click="currentConfig.enabledRolesArray = []">清空</el-button>
          <span class="role-select__count" v-if="currentConfig.enabledRolesArray.length">
            已选 {{ currentConfig.enabledRolesArray.length }} 个角色
          </span>
        </div>
        <div class="role-tags" v-if="currentConfig.enabledRolesArray.length">
          <el-tag
            v-for="roleId in currentConfig.enabledRolesArray"
            :key="roleId"
            closable
            size="small"
            effect="plain"
            @close="removeRole(roleId)"
          >
            {{ roleName(roleId) }}
          </el-tag>
        </div>
        <div class="role-empty-tip" v-else>
          <i class="el-icon-warning-outline"></i>
          未配置角色，该模块当前不允许任何人登录
        </div>
      </div>
    </div>

    <!-- 设置区块：双控开关（仅中国联通） -->
    <div class="config-section" v-if="showDualAuth">
      <div class="config-section__head">
        <h3 class="config-section__title">双控验证</h3>
        <span class="config-section__hint">写操作（上传/删除/重命名等）需另一账号复核</span>
      </div>
      <div class="dual-auth-row">
        <div class="dual-auth-row__info">
          <div class="dual-auth-row__title">启用双控</div>
          <div class="dual-auth-row__desc">
            {{ currentConfig.dualAuthEnabled ? '开启后，写操作需另一账号进行双控验证' : '关闭后，写操作无需双控验证' }}
          </div>
        </div>
        <el-switch
          v-model="currentConfig.dualAuthEnabled"
          active-color="#409EFF"
          inactive-color="#C0C4CC"
        ></el-switch>
      </div>
      <el-alert
        v-if="currentConfig.dualAuthEnabled"
        type="warning"
        :closable="false"
        show-icon
        class="dual-auth-alert"
        title="双控已开启：所有写操作将要求输入复核账号（不得与当前登录账号相同）"
      ></el-alert>
    </div>

    <!-- 保存操作栏 -->
    <div class="config-footer">
      <el-button size="medium" @click="resetConfig">重置</el-button>
      <el-button
        type="primary"
        size="medium"
        :loading="loading"
        @click="saveConfig"
      >
        {{ loading ? '保存中...' : '保存配置' }}
      </el-button>
    </div>
    </div>
  </el-card>
</template>

<script>
import sftpModulesApi from '@/api/admin/sftpModules'
import { getRoleSelect } from '@/api/settings'

/**
 * SftpModuleConfigForm 模块配置公共表单组件（设置面板风格）
 * 由各模块配置页面（标签上传/中国联通）复用，避免重复代码
 * props:
 *   - moduleName: 模块标识（hotlabel / chinaunicom）
 *   - showDualAuth: 是否显示双控开关区块（仅中国联通 true）
 */
export default {
  name: 'SftpModuleConfigForm',
  props: {
    moduleName: {
      type: String,
      required: true,
      validator: value => ['hotlabel', 'chinaunicom'].indexOf(value) !== -1
    },
    showDualAuth: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      currentConfig: {
        loginType: 'ldap',
        enabledRolesArray: [],
        dualAuthEnabled: false
      },
      roleList: [],
      loading: false,
      savedSnapshot: null // 记录保存时的快照，用于重置
    }
  },
  computed: {
    moduleTitle() {
      return this.moduleName === 'chinaunicom' ? '中国联通' : '标签上传'
    },
    moduleDesc() {
      return this.moduleName === 'chinaunicom'
        ? '中国联通数据传输模块，支持双控验证保障写操作安全'
        : '标签上传模块，用于标签数据的安全传输与归档'
    }
  },
  created() {
    this.fetchConfigs()
    this.fetchRoleList()
  },
  methods: {
    // 获取当前模块配置
    async fetchConfigs() {
      try {
        const res = await sftpModulesApi.getAllConfigs()
        if (res.code === 200) {
          const config = res.data.find(item => item.moduleName === this.moduleName)
          if (config) {
            this.currentConfig = {
              loginType: config.loginType || 'ldap',
              enabledRolesArray: JSON.parse(config.enabledRoles || '[]'),
              dualAuthEnabled: !!config.dualAuthEnabled
            }
            this.savedSnapshot = JSON.parse(JSON.stringify(this.currentConfig))
          }
        } else {
          this.$message.error('获取配置失败：' + res.message)
        }
      } catch (error) {
        console.error('获取配置失败:', error)
        this.$message.error('获取配置失败')
      }
    },

    // 获取角色列表
    async fetchRoleList() {
      try {
        const res = await getRoleSelect()
        if (res.code === 200) {
          this.roleList = res.data
        }
      } catch (error) {
        console.error('获取角色列表失败:', error)
      }
    },

    // 角色名映射
    roleName(roleId) {
      const role = this.roleList.find(item => item.ID === roleId)
      return role ? role.name : '角色#' + roleId
    },

    // 全选角色
    selectAllRoles() {
      this.currentConfig.enabledRolesArray = this.roleList.map(role => role.ID)
    },

    // 移除单个角色
    removeRole(roleId) {
      this.currentConfig.enabledRolesArray = this.currentConfig.enabledRolesArray.filter(id => id !== roleId)
    },

    // 重置为已保存的配置
    resetConfig() {
      if (this.savedSnapshot) {
        this.currentConfig = JSON.parse(JSON.stringify(this.savedSnapshot))
        this.$message.info('已恢复为上次保存的配置')
      }
    },

    // 保存配置
    async saveConfig() {
      try {
        this.loading = true
        const data = {
          loginType: this.currentConfig.loginType,
          enabledRoles: this.currentConfig.enabledRolesArray,
          dualAuthEnabled: this.currentConfig.dualAuthEnabled
        }

        const res = await sftpModulesApi.updateModuleConfig(this.moduleName, data)
        if (res.code === 200) {
          this.savedSnapshot = JSON.parse(JSON.stringify(this.currentConfig))
          this.$message.success('配置保存成功')
        } else {
          this.$message.error('配置保存失败：' + res.message)
        }
      } catch (error) {
        console.error('保存配置失败:', error)
        this.$message.error('配置保存失败')
      } finally {
        this.loading = false
      }
    }
  }
}
</script>

<style scoped>
.module-config-card {
  max-width: 860px;
  border-radius: 12px;
}
.module-config-card >>> .el-card__header {
  padding: 0;
  border-bottom: none;
}

/* ===== 模块信息头部 ===== */
.config-header {
  display: flex;
  align-items: center;
  padding: 20px 24px;
  background: linear-gradient(135deg, #f0f6ff 0%, #f8fafc 100%);
  border-bottom: 1px solid #e3edfb;
  border-radius: 12px 12px 0 0;
}
.config-header__icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 52px;
  height: 52px;
  border-radius: 12px;
  margin-right: 16px;
  flex-shrink: 0;
}
.config-header__icon i {
  font-size: 26px;
  color: #fff;
}
.config-header__icon.is-hotlabel {
  background: linear-gradient(135deg, #409EFF, #337ecc);
}
.config-header__icon.is-chinaunicom {
  background: linear-gradient(135deg, #6f7ad3, #5560c0);
}
.config-header__info {
  flex: 1;
}
.config-header__title {
  display: flex;
  align-items: center;
  font-size: 17px;
  font-weight: 600;
  color: #1f2d3d;
  gap: 8px;
}
.config-header__tag {
  margin-left: 2px;
}
.config-header__desc {
  margin-top: 6px;
  font-size: 13px;
  color: #7a8ba3;
}

/* ===== 设置区块（card 内部，分隔线划分） ===== */
.config-section {
  padding: 20px 4px;
  border-bottom: 1px solid #eef1f6;
}
.config-section:last-child {
  border-bottom: none;
}
.config-section__head {
  display: flex;
  align-items: baseline;
  margin-bottom: 16px;
}
.config-section__title {
  margin: 0;
  font-size: 15px;
  font-weight: 600;
  color: #1f2d3d;
}
.config-section__hint {
  margin-left: 12px;
  font-size: 12px;
  color: #98a6b8;
}

/* ===== 登录方式卡片 ===== */
.login-type-cards {
  display: flex;
  gap: 16px;
}
.login-type-card {
  position: relative;
  display: flex;
  align-items: center;
  flex: 1;
  padding: 16px 18px;
  background: #fafbfd;
  border: 1.5px solid #e4e9f0;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s ease;
}
.login-type-card:hover {
  border-color: #b3d4fc;
  background: #f5f9ff;
}
.login-type-card.is-active {
  border-color: #409EFF;
  background: #f0f7ff;
  box-shadow: 0 0 0 1px #409eff inset;
}
.login-type-card__icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 42px;
  height: 42px;
  border-radius: 10px;
  margin-right: 14px;
  background: #eef3fa;
  color: #5a7ca6;
  flex-shrink: 0;
  transition: all 0.2s ease;
}
.login-type-card.is-active .login-type-card__icon {
  background: #409EFF;
  color: #fff;
}
.login-type-card__icon i {
  font-size: 22px;
}
.ldap-icon {
  width: 22px;
  height: 22px;
  color: inherit;
  display: block;
}
.login-type-card__title {
  font-size: 15px;
  font-weight: 600;
  color: #1f2d3d;
}
.login-type-card__desc {
  margin-top: 4px;
  font-size: 12px;
  color: #98a6b8;
}
.login-type-card__check {
  position: absolute;
  top: 12px;
  right: 12px;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: #409EFF;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
}

/* ===== 可登录角色 ===== */
.role-select-wrap {
  padding: 4px 0;
}
.role-select {
  width: 100%;
  max-width: 480px;
}
.role-select__tools {
  display: flex;
  align-items: center;
  margin-top: 8px;
  gap: 4px;
}
.role-select__tools .el-button {
  padding: 0 6px;
}
.role-select__count {
  margin-left: 8px;
  font-size: 12px;
  color: #7a8ba3;
}
.role-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 12px;
}
.role-empty-tip {
  display: flex;
  align-items: center;
  margin-top: 12px;
  font-size: 12px;
  color: #e6a23c;
  gap: 4px;
}
.role-option-desc {
  float: right;
  margin-left: 16px;
  color: #98a6b8;
  font-size: 12px;
}

/* ===== 双控开关 ===== */
.dual-auth-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 4px 0;
}
.dual-auth-row__title {
  font-size: 14px;
  font-weight: 500;
  color: #1f2d3d;
}
.dual-auth-row__desc {
  margin-top: 4px;
  font-size: 12px;
  color: #98a6b8;
}
.dual-auth-alert {
  margin-top: 14px;
}

/* ===== 底部操作栏 ===== */
.config-card-body {
  padding: 4px 20px;
}
.config-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px 20px;
  border-top: 1px solid #eef1f6;
  background: #fafbfd;
  border-radius: 0 0 12px 12px;
}
.config-footer .el-button {
  min-width: 110px;
}
</style>
