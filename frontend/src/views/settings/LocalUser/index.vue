<template>
  <div>
    <el-card class="user-config-card" shadow="never">
      <!-- 头部信息 -->
      <div slot="header" class="user-header">
        <div class="user-header__icon is-primary">
          <svg width="26" height="26" viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M24 16C28.9706 16 33 12.9706 33 8C33 3.02944 28.9706 0 24 0C19.0294 0 15 3.02944 15 8C15 12.9706 19.0294 16 24 16Z" stroke="currentColor" stroke-width="4"/>
            <path d="M8 40C8 35.5817 13.0294 32 19 32H29C34.9706 32 40 35.5817 40 40V42H8V40Z" stroke="currentColor" stroke-width="4"/>
          </svg>
        </div>
        <div class="user-header__info">
          <div class="user-header__title">
            本地账号管理
            <el-tag size="mini" type="success" effect="light" class="user-header__tag">
              系统账号
            </el-tag>
          </div>
          <div class="user-header__desc">管理平台本地用户账号及其权限、密码策略配置</div>
        </div>
      </div>

      <!-- 列表内容区 -->
      <div class="user-content">
        <!-- 顶部工具栏 -->
        <div class="user-toolbar">
          <div class="toolbar-left">
            <el-input 
              v-model="searchQuery" 
              placeholder="搜索用户名" 
              clearable
              prefix-icon="el-icon-search"
              style="width: 280px;"
              @keyup.enter.native="handleSearch"
            />
          </div>
          <div class="toolbar-right">
            <el-button 
              type="primary" 
              icon="el-icon-plus" 
              @click="showCreateDialog"
              class="add-btn"
            >
              新增账号
            </el-button>
          </div>
        </div>

        <!-- 账号列表表格 -->
        <div class="table-container">
          <el-table
            v-loading="loading"
            :data="userList" 
            border 
            stripe 
            style="width: 100%"
            :header-cell-style="{ background: '#f7f9fa', color: '#6a7b9c', fontWeight: 600 }"
          >
            <el-table-column prop="username" label="用户名" min-width="150">
              <template slot-scope="{ row }">
                <div class="username-cell">
                  <span class="username">{{ row.username }}</span>
                  <el-tag v-if="row.username === 'admin'" type="danger" size="mini" effect="plain">默认</el-tag>
                </div>
              </template>
            </el-table-column>
            
            <el-table-column prop="roleName" label="角色" min-width="120">
              <template slot-scope="{ row }">
                <span class="role-name">{{ row.roleName || '-' }}</span>
              </template>
            </el-table-column>
            
            <el-table-column label="状态" width="100" align="center">
              <template slot-scope="{ row }">
                <div class="status-badge" :class="row.enabled ? 'success' : 'danger'">
                  <i :class="row.enabled ? 'el-icon-circle-check' : 'el-icon-remove'"></i>
                  <span>{{ row.enabled ? '启用' : '禁用' }}</span>
                </div>
              </template>
            </el-table-column>
            
            <el-table-column label="需改密" width="100" align="center">
              <template slot-scope="{ row }">
                <div class="flag-badge warning" v-if="row.mustChangePassword">
                  <i class="el-icon-key"></i>
                  <span>是</span>
                </div>
                <span v-else class="text-muted">否</span>
              </template>
            </el-table-column>
            
            <el-table-column label="永不过期" width="110" align="center">
              <template slot-scope="{ row }">
                <div class="flag-badge info" v-if="row.passwordNeverExpires">
                  <i class="el-icon-time"></i>
                  <span>是</span>
                </div>
                <span v-else class="text-muted">否</span>
              </template>
            </el-table-column>
            
            <el-table-column label="失败次数" width="100" align="center">
              <template slot-scope="{ row }">
                <div class="stat-badge" :class="{ danger: row.failedAttempts >= 5 }">
                  <i class="el-icon-warning"></i>
                  <span>{{ row.failedAttempts || 0 }}</span>
                </div>
              </template>
            </el-table-column>
            
            <el-table-column prop="lastLoginAt" label="最后登录" width="180" align="center">
              <template slot-scope="{ row }">
                <span class="date-text">{{ row.lastLoginAt | dateFormat }}</span>
              </template>
            </el-table-column>
            
            <el-table-column prop="CreatedAt" label="创建时间" width="180" align="center">
              <template slot-scope="{ row }">
                <span class="date-text">{{ row.CreatedAt | dateFormat }}</span>
              </template>
            </el-table-column>
            
            <el-table-column label="操作" width="320" fixed="right" align="center">
              <template slot-scope="{ row }">
                <div class="action-buttons">
                  <el-button 
                    type="primary" 
                    size="small" 
                    plain
                    icon="el-icon-edit" 
                    @click="showEditDialog(row)"
                  >
                    编辑
                  </el-button>
                  <el-button 
                    type="warning" 
                    size="small" 
                    plain
                    icon="el-icon-key" 
                    @click="showResetPasswordDialog(row)"
                  >
                    重置密码
                  </el-button>
                  <el-button 
                    type="danger" 
                    size="small" 
                    plain
                    icon="el-icon-delete" 
                    @click="handleDelete(row)"
                    :disabled="row.username === 'admin'"
                  >
                    删除
                  </el-button>
                </div>
              </template>
            </el-table-column>
          </el-table>

          <!-- 分页 -->
          <el-pagination
            v-show="total > 0"
            :current-page="page"
            :page-sizes="[10, 20, 50]"
            :page-size="limit"
            :total="total"
            layout="total, sizes, prev, pager, next, jumper"
            style="margin-top: 16px;"
            @size-change="handleSizeChange"
            @current-change="handlePageChange"
          />
        </div>
      </div>
    </el-card>

    <!-- 新增/编辑对话框 -->
    <el-dialog :title="dialogTitle" :visible.sync="dialogVisible" width="600px" @closed="resetDialog" class="user-dialog">
      <el-form ref="userForm" :model="userForm" :rules="userRules" label-width="120px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="userForm.username" :disabled="isEdit" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item v-if="!isEdit" label="密码" prop="password">
          <el-input v-model="userForm.password" type="password" placeholder="请输入密码" show-password />
        </el-form-item>
        <el-form-item v-if="!isEdit" label="确认密码" prop="checkPass">
          <el-input v-model="userForm.checkPass" type="password" placeholder="请再次输入密码" show-password />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="userForm.roleId" placeholder="请选择角色" style="width: 100%;" clearable>
            <el-option v-for="role in roleOptions" :key="role.ID" :label="role.name" :value="role.ID" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="userForm.enabled" active-text="启用" inactive-text="禁用" />
        </el-form-item>
        <el-form-item label="需修改密码">
          <el-switch v-model="userForm.mustChangePassword" />
          <span style="margin-left: 8px; color: #98a6b8; font-size: 12px;">用户首次登录时需修改密码</span>
        </el-form-item>
        <el-form-item label="密码永不过期">
          <el-switch v-model="userForm.passwordNeverExpires" />
          <span style="margin-left: 8px; color: #98a6b8; font-size: 12px;">密码不受有效期限制</span>
        </el-form-item>
      </el-form>
      <span slot="footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">{{ isEdit ? '保存' : '创建' }}</el-button>
      </span>
    </el-dialog>

    <!-- 重置密码对话框 -->
    <el-dialog title="重置密码" :visible.sync="resetDialogVisible" width="500px" @closed="resetPassword = ''" class="reset-dialog">
      <el-form label-width="120px">
        <el-form-item label="用户名">
          <span class="reset-username">{{ resetUsername }}</span>
        </el-form-item>
        <el-form-item label="新密码" :required="true">
          <el-input v-model="resetPassword" type="password" placeholder="请输入新密码" show-password />
          <div class="password-hint">
            <i class="el-icon-info"></i>
            <span>密码需符合系统密码策略要求</span>
          </div>
        </el-form-item>
      </el-form>
      <span slot="footer">
        <el-button @click="resetDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="resetLoading" @click="handleResetPassword">确认重置</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import { getLocalUserList, createLocalUser, updateLocalUser, deleteLocalUser, resetLocalUserPassword } from '@/api/settings'
import { getRoleSelect } from '@/api/settings'
import { rsaEncrypt } from '@/utils/encrypt'

export default {
  name: 'LocalUserManagement',
  filters: {
    dateFormat(val) {
      return val ? val.slice(0, 19).replace('T', ' ') : ''
    }
  },
  data() {
    return {
      loading: false,
      userList: [],
      total: 0,
      page: 1,
      limit: 10,
      searchQuery: '',
      roleOptions: [],
      dialogVisible: false,
      dialogTitle: '',
      submitLoading: false,
      isEdit: false,
      editId: null,
      userForm: {
        username: '',
        password: '',
        checkPass: '',
        roleId: null,
        enabled: true,
        mustChangePassword: false,
        passwordNeverExpires: false
      },
      userRules: {
        username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
        password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
        checkPass: [
          { required: true, message: '请再次输入密码', trigger: 'blur' },
          { validator: this.validateCheckPass, trigger: 'blur' }
        ]
      },
      resetDialogVisible: false,
      resetUserId: null,
      resetUsername: '',
      resetPassword: '',
      resetLoading: false
    }
  },
  created() {
    this.fetchUserList()
    this.fetchRoleOptions()
  },
  methods: {
    // 确认密码校验
    validateCheckPass(rule, value, callback) {
      if (value !== this.userForm.password) {
        callback(new Error('两次输入的密码不一致'))
      } else {
        callback()
      }
    },
    async fetchUserList() {
      this.loading = true
      try {
        const res = await getLocalUserList({ page: this.page, limit: this.limit, username: this.searchQuery })
        if (res.code === 200) {
          this.userList = res.data.list
          this.total = res.data.total
        }
      } catch (e) {
        console.error(e)
      }
      this.loading = false
    },
    async fetchRoleOptions() {
      try {
        const res = await getRoleSelect()
        if (res.code === 200) {
          this.roleOptions = res.data
        }
      } catch (e) {
        console.error(e)
      }
    },
    handleSearch() {
      this.page = 1
      this.fetchUserList()
    },
    handleSizeChange(val) {
      this.limit = val
      this.fetchUserList()
    },
    handlePageChange(val) {
      this.page = val
      this.fetchUserList()
    },
    showCreateDialog() {
      this.isEdit = false
      this.editId = null
      this.dialogTitle = '新增本地账号'
      this.userForm.enabled = true
      this.userForm.mustChangePassword = false
      this.userForm.passwordNeverExpires = false
      this.dialogVisible = true
    },
    showEditDialog(row) {
      this.isEdit = true
      this.editId = row.ID
      this.dialogTitle = '编辑账号'
      this.userForm.username = row.username
      this.userForm.password = ''
      this.userForm.roleId = row.roleId || null
      this.userForm.enabled = row.enabled
      this.userForm.mustChangePassword = !!row.mustChangePassword
      this.userForm.passwordNeverExpires = !!row.passwordNeverExpires
      this.dialogVisible = true
    },
    resetDialog() {
      this.userForm = { username: '', password: '', checkPass: '', roleId: null, enabled: true, mustChangePassword: false, passwordNeverExpires: false }
      this.submitLoading = false
    },
    async handleSubmit() {
      this.$refs.userForm.validate(async valid => {
        if (!valid) return
        this.submitLoading = true
        try {
          let res
          if (this.isEdit) {
            res = await updateLocalUser(this.editId, {
              roleId: this.userForm.roleId,
              enabled: this.userForm.enabled,
              mustChangePassword: this.userForm.mustChangePassword,
              passwordNeverExpires: this.userForm.passwordNeverExpires
            })
          } else {
            const rsaPassword = await rsaEncrypt(this.userForm.password)
            res = await createLocalUser({
              username: this.userForm.username,
              password: rsaPassword,
              roleId: this.userForm.roleId,
              enabled: this.userForm.enabled,
              mustChangePassword: this.userForm.mustChangePassword,
              passwordNeverExpires: this.userForm.passwordNeverExpires
            })
          }
          if (res.code === 200) {
            this.$message.success(res.message)
            this.dialogVisible = false
            this.fetchUserList()
          } else {
            this.$message.error(res.message)
          }
        } catch (e) {
          console.error(e)
        }
        this.submitLoading = false
      })
    },
    showResetPasswordDialog(row) {
      this.resetUserId = row.ID
      this.resetUsername = row.username
      this.resetPassword = ''
      this.resetDialogVisible = true
    },
    async handleResetPassword() {
      if (!this.resetPassword) {
        this.$message.warning('请输入新密码')
        return
      }
      this.resetLoading = true
      try {
        const rsaPassword = await rsaEncrypt(this.resetPassword)
        const res = await resetLocalUserPassword(this.resetUserId, { password: rsaPassword })
        if (res.code === 200) {
          this.$message.success(res.message)
          this.resetDialogVisible = false
        } else {
          this.$message.error(res.message)
        }
      } catch (e) {
        console.error('重置密码失败:', e)
        this.$message.error('重置密码失败，请重试')
      }
      this.resetLoading = false
    },
    async handleDelete(row) {
      try {
        await this.$confirm(`确定删除账号"${row.username}"？`, '提示', { type: 'warning' })
        const res = await deleteLocalUser(row.ID)
        if (res.code === 200) {
          this.$message.success('删除成功')
          this.fetchUserList()
        } else {
          this.$message.error(res.message)
        }
      } catch (e) {
        if (e !== 'cancel') console.error(e)
      }
    }
  }
}
</script>

<style scoped>

/* ===== 用户配置卡片 ===== */
.user-config-card {
  max-width: 1500px;
  margin: 0 auto;
  border-radius: 12px;
}

.user-config-card >>> .el-card__header {
  padding: 0;
  border-bottom: none;
}

.user-config-card >>> .el-card__body {
  padding: 0;
}

/* ===== 用户头部信息 ===== */
.user-header {
  display: flex;
  align-items: center;
  padding: 20px 24px;
  background: linear-gradient(135deg, #f0f6ff 0%, #f8fafc 100%);
  border-bottom: 1px solid #e3edfb;
  border-radius: 12px 12px 0 0;
}

.user-header__icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 52px;
  height: 52px;
  border-radius: 12px;
  margin-right: 16px;
  flex-shrink: 0;
}

.user-header__icon svg {
  color: #fff;
  stroke-width: 3;
}

.user-header__icon.is-primary {
  background: linear-gradient(135deg, #409EFF, #337ecc);
}

.user-header__info {
  flex: 1;
}

.user-header__title {
  display: flex;
  align-items: center;
  font-size: 17px;
  font-weight: 600;
  color: #1f2d3d;
  gap: 8px;
}

.user-header__tag {
  margin-left: 2px;
}

.user-header__desc {
  margin-top: 6px;
  font-size: 13px;
  color: #7a8ba3;
}

/* ===== 内容区布局 ===== */
.user-content {
  padding: 4px 20px;
}

/* ===== 工具栏 ===== */
.user-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 8px;
  gap: 16px;
}

.toolbar-left {
  flex: 1;
}

.toolbar-right {
  white-space: nowrap;
}

.add-btn {
  min-width: 110px;
}

/* ===== 表格容器 ===== */
.table-container {
  position: relative;
}

.username-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.username {
  font-weight: 600;
  color: #1f2d3d;
}

.role-name {
  color: #7a8ba3;
}

.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 600;
}

.status-badge.success {
  color: #67C23A;
  background: #f0f9eb;
}

.status-badge.danger {
  color: #F56C6C;
  background: #fef0f0;
}

.flag-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 600;
}

.flag-badge.warning {
  color: #E6A23C;
  background: #fdf6ec;
}

.flag-badge.info {
  color: #909399;
  background: #f4f4f5;
}

.text-muted {
  color: #c0c4cc;
  font-size: 13px;
}

.stat-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  background: #f7f9fa;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 600;
  color: #409EFF;
}

.stat-badge.danger {
  color: #F56C6C;
  background: #fef0f0;
}

.date-text {
  color: #98a6b8;
  font-size: 13px;
}

.action-buttons {
  display: flex;
  gap: 8px;
  justify-content: center;
}

.action-buttons .el-button {
  min-width: 70px;
}

/* ===== 对话框样式 ===== */
.user-dialog >>> .el-dialog__header,
.reset-dialog >>> .el-dialog__header {
  padding: 20px 24px;
  border-bottom: 1px solid #eef1f6;
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-dialog >>> .el-dialog__title,
.reset-dialog >>> .el-dialog__title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2d3d;
}

.user-dialog >>> .el-dialog__body,
.reset-dialog >>> .el-dialog__body {
  padding: 24px;
}

.user-dialog >>> .el-form-item {
  margin-bottom: 20px;
}

.user-dialog >>> .el-form-item__label {
  font-weight: 600;
  color: #1f2d3d;
}

.user-dialog >>> .el-input__inner {
  border-radius: 6px;
  border-color: #e4e9f0;
}

.user-dialog >>> .el-input__inner:focus {
  border-color: #409EFF;
}

.user-dialog >>> .el-button--primary {
  border-radius: 6px;
  min-width: 100px;
}

/* ===== 重置密码对话框特殊样式 ===== */
.reset-username {
  font-weight: 600;
  color: #1f2d3d;
  font-size: 14px;
}

.password-hint {
  margin-top: 8px;
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  background: #f0f7ff;
  border-radius: 6px;
  font-size: 12px;
  color: #409EFF;
}

.password-hint i {
  font-size: 14px;
}

/* ===== 响应式 ===== */
@media screen and (max-width: 1600px) {
  .user-config-card {
    max-width: 100%;
  }
}
</style>