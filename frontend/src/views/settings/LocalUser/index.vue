<template>
  <div class="app-container">
    <div class="search-box">
      <el-input v-model="searchQuery" placeholder="搜索用户名" clearable style="width: 200px;" @keyup.enter.native="handleSearch" />
      <el-button type="primary" icon="el-icon-search" @click="handleSearch">搜索</el-button>
      <el-button type="success" icon="el-icon-plus" @click="showCreateDialog">新增账号</el-button>
    </div>

    <el-table v-loading="loading" :data="userList" border stripe style="width: 100%">
      <el-table-column prop="username" label="用户名" min-width="120" />
      <el-table-column prop="roleName" label="角色" min-width="120" />
      <el-table-column label="状态" width="100">
        <template slot-scope="{ row }">
          <el-tag :type="row.enabled ? 'success' : 'danger'">{{ row.enabled ? '启用' : '禁用' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="需改密" width="80">
        <template slot-scope="{ row }">
          <el-tag v-if="row.mustChangePassword" type="warning" size="mini">是</el-tag>
          <span v-else>否</span>
        </template>
      </el-table-column>
      <el-table-column label="永不过期" width="80">
        <template slot-scope="{ row }">
          <el-tag v-if="row.passwordNeverExpires" type="info" size="mini">是</el-tag>
          <span v-else>否</span>
        </template>
      </el-table-column>
      <el-table-column label="失败次数" width="80" align="center">
        <template slot-scope="{ row }">
          <span>{{ row.failedAttempts || 0 }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="lastLoginAt" label="最后登录" width="200">
        <template slot-scope="{ row }">{{ row.lastLoginAt | dateFormat }}</template>
      </el-table-column>
      <el-table-column prop="CreatedAt" label="创建时间" width="200">
        <template slot-scope="{ row }">{{ row.CreatedAt | dateFormat }}</template>
      </el-table-column>
      <el-table-column label="操作" width="320" fixed="right">
        <template slot-scope="{ row }">
          <el-button type="primary" size="mini" icon="el-icon-edit" @click="showEditDialog(row)">编辑</el-button>
          <el-button type="warning" size="mini" icon="el-icon-key" @click="showResetPasswordDialog(row)">重置密码</el-button>
          <el-button type="danger" size="mini" icon="el-icon-delete" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-show="total > 0"
      :current-page="page"
      :page-sizes="[10, 20, 50]"
      :page-size="limit"
      :total="total"
      layout="total, sizes, prev, pager, next, jumper"
      style="margin-top: 20px;"
      @size-change="handleSizeChange"
      @current-change="handlePageChange"
    />

    <el-dialog :title="dialogTitle" :visible.sync="dialogVisible" width="500px" @closed="resetDialog">
      <el-form ref="userForm" :model="userForm" :rules="userRules" label-width="100px">
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
        </el-form-item>
        <el-form-item label="密码永不过期">
          <el-switch v-model="userForm.passwordNeverExpires" />
        </el-form-item>
      </el-form>
      <span slot="footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">{{ isEdit ? '保存' : '创建' }}</el-button>
      </span>
    </el-dialog>

    <el-dialog title="重置密码" :visible.sync="resetDialogVisible" width="450px" @closed="resetPassword = ''">
      <el-form label-width="100px">
        <el-form-item label="用户名">
          <span>{{ resetUsername }}</span>
        </el-form-item>
        <el-form-item label="新密码" :required="true">
          <el-input v-model="resetPassword" type="password" placeholder="请输入新密码" show-password />
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
            const rsaPassword = rsaEncrypt(this.userForm.password)
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
        const rsaPassword = rsaEncrypt(this.resetPassword)
        const res = await resetLocalUserPassword(this.resetUserId, { password: rsaPassword })
        if (res.code === 200) {
          this.$message.success(res.message)
          this.resetDialogVisible = false
        } else {
          this.$message.error(res.message)
        }
      } catch (e) {
        console.error(e)
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
.search-box {
  margin-bottom: 16px;
  display: flex;
  gap: 8px;
  align-items: center;
}
</style>