<template>
  <div>
    <el-card class="role-config-card" shadow="never">
      <!-- 头部信息 -->
      <div slot="header" class="role-header">
        <div class="role-header__icon is-primary">
          <svg width="26" height="26" viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M24 16C28.9706 16 33 12.9706 33 8C33 3.02944 28.9706 0 24 0C19.0294 0 15 3.02944 15 8C15 12.9706 19.0294 16 24 16Z" stroke="currentColor" stroke-width="4"/>
            <path d="M8 40C8 35.5817 13.0294 32 19 32H29C34.9706 32 40 35.5817 40 40V42H8V40Z" stroke="currentColor" stroke-width="4"/>
          </svg>
        </div>
        <div class="role-header__info">
          <div class="role-header__title">
            角色管理
            <el-tag size="mini" type="success" effect="light" class="role-header__tag">
              权限控制
            </el-tag>
          </div>
          <div class="role-header__desc">管理系统角色及其菜单权限与安全组关联配置</div>
        </div>
      </div>

      <!-- 列表内容区 -->
      <div class="role-content">
        <!-- 顶部工具栏 -->
        <div class="role-toolbar">
          <div class="toolbar-left">
            <el-input 
              v-model="searchQuery" 
              placeholder="搜索角色名称" 
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
              新增角色
            </el-button>
          </div>
        </div>

        <!-- 角色列表表格 -->
        <div class="table-container">
          <el-table
            v-loading="loading"
            :data="roleList" 
            border 
            stripe 
            style="width: 100%"
            :header-cell-style="{ background: '#f7f9fa', color: '#6a7b9c', fontWeight: 600 }"
          >
          <el-table-column prop="name" label="角色名称" min-width="180">
            <template slot-scope="{ row }">
              <div class="role-name-cell">
                <span class="role-name">{{ row.name }}</span>
                <el-tag v-if="row.name === '超级管理员'" type="danger" size="mini" effect="plain">默认</el-tag>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column prop="description" label="角色描述" min-width="200">
            <template slot-scope="{ row }">
              <span class="role-desc">{{ row.description || '-' }}</span>
            </template>
          </el-table-column>
          
          <el-table-column label="关联菜单数" width="120" align="center">
            <template slot-scope="{ row }">
              <div class="stat-badge">
                <i class="el-icon-menu"></i>
                <span>{{ row.menus ? row.menus.length : 0 }}</span>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column label="关联安全组" width="140" align="center">
            <template slot-scope="{ row }">
              <div class="stat-badge warning">
                <i class="el-icon-user-solid"></i>
                <span>{{ row.ldapGroupCount || 0 }}</span>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column prop="createdAt" label="创建时间" width="180" align="center">
            <template slot-scope="{ row }">
              <span class="date-text">{{ row.CreatedAt | dateFormat }}</span>
            </template>
          </el-table-column>
          
          <el-table-column label="操作" width="200" fixed="right" align="center">
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
                  type="danger" 
                  size="small" 
                  plain
                  icon="el-icon-delete" 
                  @click="handleDelete(row)"
                  :disabled="row.name === '超级管理员'"
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
    <!-- 弹窗样式 -->
    <el-dialog :title="dialogTitle" :visible.sync="dialogVisible" width="700px" @closed="resetDialog" class="role-dialog">
          
      <!-- 步骤条指示器 -->
      <div class="step-indicator">
        <el-steps :active="currentStep - 1" finish-status="success" align-center>
          <el-step title="基本信息"></el-step>
          <el-step title="菜单权限"></el-step>
          <el-step title="LDAP 安全组"></el-step>
        </el-steps>
      </div>
          
      <el-form ref="roleForm" :model="roleForm" :rules="roleRules" label-width="120px">
            
        <!-- 第一步：基本信息 -->
        <div v-if="currentStep === 1" class="step-content">
          <el-form-item label="角色名称" prop="name">
            <el-input v-model="roleForm.name" :disabled="isSuperAdmin" placeholder="请输入角色名称" />
          </el-form-item>
          <el-form-item label="角色描述" prop="description">
            <el-input v-model="roleForm.description" :disabled="isSuperAdmin" type="textarea" :rows="3" placeholder="请输入角色描述" />
          </el-form-item>
        </div>
            
        <!-- 第二步：菜单权限 -->
        <div v-if="currentStep === 2" class="step-content">
          <div style="margin-bottom: 12px; color: #6a7b9c; font-size: 13px;">
            请勾选该角色可访问的菜单项
          </div>
          <el-tree
            ref="menuTree"
            :data="allMenus"
            show-checkbox
            node-key="routeName"
            :props="{ label: 'menuTitle', children: 'children' }"
            default-expand-all
            :disabled="isSuperAdmin"
            @check="handleMenuCheck"
          />
          <span v-if="isSuperAdmin" style="color: #999; font-size: 12px;">超级管理员菜单权限不可修改</span>
        </div>
            
        <!-- 第三步：LDAP 安全组 -->
        <div v-if="currentStep === 3" class="step-content">
          <div style="margin-bottom: 12px; color: #6a7b9c; font-size: 13px;">
            请填写 LDAP 安全组 DN 和显示名称，用于登录时自动匹配角色
          </div>
          <el-button type="primary" size="mini" icon="el-icon-plus" @click="addLDAPGroup" style="margin-bottom: 12px;">
            添加安全组
          </el-button>
          <div v-for="(group, index) in roleForm.ldapGroups" :key="index" class="ldap-group-row">
            <el-input v-model="group.group_dn" placeholder="安全组 DN，如 CN=IT部，OU=..." style="flex: 2;" />
            <el-input v-model="group.group_name" placeholder="显示名称" style="flex: 1;" />
            <el-button type="danger" icon="el-icon-delete" size="mini" @click="removeLDAPGroup(index)" />
          </div>
        </div>
            
      </el-form>
          
      <!-- 底部按钮 -->
      <span slot="footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button v-if="currentStep > 1" @click="previousStep">上一步</el-button>
        <el-button v-if="currentStep < 3" type="primary" @click="nextStep">下一步</el-button>
        <el-button v-else type="primary" :loading="submitLoading" @click="handleSubmit">保存</el-button>
      </span>
          
    </el-dialog>
    </div>
</template>

<script>
import { getRoleList, createRole, updateRole, deleteRole, getRoleDetail, getMenus, getRoleSelect } from '@/api/settings'

export default {
  name: 'RoleManagement',
  filters: {
    dateFormat(val) {
      return val ? val.slice(0, 19).replace('T', ' ') : ''
    }
  },
  data() {
    return {
      loading: false,
      roleList: [],
      total: 0,
      page: 1,
      limit: 10,
      searchQuery: '',
      dialogVisible: false,
      dialogTitle: '',
      submitLoading: false,
      isEdit: false,
      editId: null,
      currentStep: 1, // 当前步骤：1=基本信息，2=菜单权限，3=LDAP 安全组
      isSuperAdmin: false, // 是否为超级管理员角色
      allMenus: [],
      roleForm: {
        name: '',
        description: '',
        menus: [],
        ldapGroups: []
      },
      roleRules: {
        name: [{ required: true, message: '请输入角色名称', trigger: 'blur' }]
      }
    }
  },
  created() {
    this.fetchRoleList()
    this.fetchMenus()
  },
  methods: {
    async fetchRoleList() {
      this.loading = true
      try {
        const res = await getRoleList({ page: this.page, limit: this.limit, name: this.searchQuery })
        if (res.code === 200) {
          this.roleList = res.data.list
          this.total = res.data.total
        }
      } catch (e) {
        console.error(e)
      }
      this.loading = false
    },
    async fetchMenus() {
      try {
        const res = await getMenus()
        if (res.code === 200) {
          this.allMenus = res.data
        }
      } catch (e) {
        console.error(e)
      }
    },
    handleSearch() {
      this.page = 1
      this.fetchRoleList()
    },
    handleSizeChange(val) {
      this.limit = val
      this.fetchRoleList()
    },
    handlePageChange(val) {
      this.page = val
      this.fetchRoleList()
    },
    showCreateDialog() {
      this.isEdit = false
      this.editId = null
      this.currentStep = 1 // 重置到第一步
      this.dialogTitle = '新增角色'
      this.dialogVisible = true
    },
    async showEditDialog(row) {
      this.isEdit = true
      this.editId = row.ID
      this.isSuperAdmin = row.name === '超级管理员'
      this.currentStep = 1 // 编辑也从第一步开始
      this.dialogTitle = this.isSuperAdmin ? '超级管理员（仅可修改 LDAP 安全组）' : '编辑角色'
      try {
        const res = await getRoleDetail(row.ID)
        if (res.code === 200) {
          const role = res.data
          this.roleForm.name = role.name
          this.roleForm.description = role.description
          this.roleForm.menus = (role.menus || []).map(m => m.routeName)
          this.roleForm.ldapGroups = (role.ldapGroups || []).map(g => ({
            group_dn: g.groupDN,
            group_name: g.groupName
          }))
          this.$nextTick(() => {
            this.restoreMenuChecks()
          })
        }
      } catch (e) {
        console.error(e)
      }
      this.dialogVisible = true
    },
    handleMenuCheck() {
      if (this.$refs.menuTree) {
        // 超级管理员：不响应勾选变更，始终保持全选状态
        if (this.isSuperAdmin) {
          const allMenuKeys = this.collectAllMenuKeys(this.allMenus)
          this.$refs.menuTree.setCheckedKeys(allMenuKeys)
          return
        }
        this.roleForm.menus = this.$refs.menuTree.getCheckedKeys()
      }
    },
    // 递归收集菜单树中所有 routeName
    collectAllMenuKeys(menus) {
      let keys = []
      menus.forEach(item => {
        if (item.routeName) {
          keys.push(item.routeName)
        }
        if (item.children && item.children.length) {
          keys = keys.concat(this.collectAllMenuKeys(item.children))
        }
      })
      return keys
    },
    addLDAPGroup() {
      this.roleForm.ldapGroups.push({ group_dn: '', group_name: '' })
    },
    removeLDAPGroup(index) {
      this.roleForm.ldapGroups.splice(index, 1)
    },
    resetDialog() {
      this.currentStep = 1 // 重置步骤
      this.roleForm = { name: '', description: '', menus: [], ldapGroups: [] }
      this.submitLoading = false
      this.isSuperAdmin = false
    },
    
    // 步骤导航 - 下一步
    nextStep() {
      if (this.currentStep === 1) {
        // 验证基本信息
        this.$refs.roleForm.validate((valid) => {
          if (valid) {
            this.currentStep++
            this.restoreMenuChecks()
            this.scrollToTop()
          }
        })
      } else if (this.currentStep === 2) {
        this.currentStep++
        this.scrollToTop()
      }
    },
    
    // 恢复菜单树的勾选状态（树使用 v-if 按步骤渲染，每次进入第二步需重新设置）
    restoreMenuChecks() {
      this.$nextTick(() => {
        if (this.$refs.menuTree) {
          if (this.isSuperAdmin) {
            const allMenuKeys = this.collectAllMenuKeys(this.allMenus)
            this.$refs.menuTree.setCheckedKeys(allMenuKeys)
          } else {
            this.$refs.menuTree.setCheckedKeys(this.roleForm.menus)
          }
        }
      })
    },
    
    // 步骤导航 - 上一步
    previousStep() {
      if (this.currentStep > 1) {
        this.currentStep--
        // 回到第二步时恢复菜单树勾选状态（树被 v-if 销毁重建，需重新设置）
        if (this.currentStep === 2) {
          this.restoreMenuChecks()
        }
        this.scrollToTop()
      }
    },
    
    // 滚动到内容顶部
    scrollToTop() {
      this.$nextTick(() => {
        const body = this.$el.querySelector('.el-dialog__body')
        if (body) body.scrollTop = 0
      })
    },
    async handleSubmit() {
      this.$refs.roleForm.validate(async valid => {
        if (!valid) return
        this.submitLoading = true
        try {
          const data = {
            name: this.roleForm.name,
            description: this.roleForm.description,
            menus: this.roleForm.menus,
            ldapGroups: this.roleForm.ldapGroups.filter(g => g.group_dn)
          }
          let res
          if (this.isEdit) {
            res = await updateRole(this.editId, data)
          } else {
            res = await createRole(data)
          }
          if (res.code === 200) {
            this.$message.success(res.message)
            this.dialogVisible = false
            this.fetchRoleList()
          } else {
            this.$message.error(res.message)
          }
        } catch (e) {
          console.error(e)
        }
        this.submitLoading = false
      })
    },
    async handleDelete(row) {
      try {
        await this.$confirm(`确定删除角色"${row.name}"？`, '提示', { type: 'warning' })
        const res = await deleteRole(row.ID)
        if (res.code === 200) {
          this.$message.success('删除成功')
          this.fetchRoleList()
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

/* ===== 角色配置卡片 ===== */
.role-config-card {
  max-width: 1200px;
  margin: 0 auto;
  border-radius: 12px;
}

.role-config-card >>> .el-card__header {
  padding: 0;
  border-bottom: none;
}

.role-config-card >>> .el-card__body {
  padding: 0;
}

/* ===== 角色头部信息 ===== */
.role-header {
  display: flex;
  align-items: center;
  padding: 20px 24px;
  background: linear-gradient(135deg, #f0f6ff 0%, #f8fafc 100%);
  border-bottom: 1px solid #e3edfb;
  border-radius: 12px 12px 0 0;
}

.role-header__icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 52px;
  height: 52px;
  border-radius: 12px;
  margin-right: 16px;
  flex-shrink: 0;
}

.role-header__icon svg {
  color: #fff;
  stroke-width: 3;
}

.role-header__icon.is-primary {
  background: linear-gradient(135deg, #409EFF, #337ecc);
}

.role-header__info {
  flex: 1;
}

.role-header__title {
  display: flex;
  align-items: center;
  font-size: 17px;
  font-weight: 600;
  color: #1f2d3d;
  gap: 8px;
}

.role-header__tag {
  margin-left: 2px;
}

.role-header__desc {
  margin-top: 6px;
  font-size: 13px;
  color: #7a8ba3;
}

/* ===== 内容区布局 ===== */
.role-content {
  padding: 4px 20px;
}

/* ===== 工具栏 ===== */
.role-toolbar {
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

.role-name-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.role-name {
  font-weight: 600;
  color: #1f2d3d;
}

.role-desc {
  color: #7a8ba3;
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

.stat-badge.warning {
  color: #E6A23C;
  background: #fdf6ec;
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
  min-width: 60px;
}

/* ===== 对话框样式 ===== */
.role-dialog >>> .el-dialog__header {
  padding: 12px 24px;
  border-bottom: 1px solid #eef1f6;
  display: flex;
  align-items: center;
  gap: 12px;
}

.role-dialog >>> .el-dialog__title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2d3d;
}

.role-dialog >>> .el-dialog__body {
  padding: 24px;
}

.role-dialog >>> .el-form-item {
  margin-bottom: 20px;
}

.role-dialog >>> .el-form-item__label {
  font-weight: 600;
  color: #1f2d3d;
}

.role-dialog >>> .el-input__inner {
  border-radius: 6px;
  border-color: #e4e9f0;
}

.role-dialog >>> .el-input__inner:focus {
  border-color: #409EFF;
}

.role-dialog >>> .el-button--primary {
  border-radius: 6px;
  min-width: 100px;
}

.role-dialog >>> .el-textarea__inner {
  border-radius: 6px;
  border-color: #e4e9f0;
  resize: none;
}

.role-dialog >>> .el-textarea__inner:focus {
  border-color: #409EFF;
}

/* ===== 响应式 ===== */
@media screen and (max-width: 1440px) {
  .role-config-card {
    max-width: 100%;
  }
}

/* ===== 步骤条指示器 ===== */
.role-dialog >>> .step-indicator {
  padding: 20px 24px;
  border-bottom: 1px solid #eef1f6;
  background: linear-gradient(135deg, #f0f6ff 0%, #f8fafc 100%);
  margin-top: -20px;
  margin-left: -24px;
  margin-right: -24px;
  margin-bottom: 10px;
}

.role-dialog >>> .el-step__head.is-process {
  color: #409EFF;
  font-weight: 600;
}

.role-dialog >>> .el-step__title {
  font-size: 14px;
  font-weight: 600;
}

/* ===== 内容区域高度（dialog 总高控制在 590px 内，不滚动） ===== */
.role-dialog >>> .el-dialog__body {
  padding: 16px 24px;
}

/* ===== 步骤内容容器：固定高度，三步共用，保证弹窗高度稳定 ===== */
.step-content {
  height: 450px;
  overflow-y: auto; /* 极端情况兑底，正常三步内容均不触发 */
}

/* ===== LDAP 安全组输入优化 ===== */
.ldap-group-row {
  display: flex;
  gap: 8px;
  margin-top: 8px;
  align-items: center;
}

.ldap-group-row >>> .el-input__inner {
  border-radius: 6px;
  height: 36px;
  font-size: 13px;
}

.ldap-group-row >>> .el-button--mini {
  height: 36px;
  padding: 9px 15px;
  flex-shrink: 0;
}

/* ===== 菜单权限树优化（树内部可滚动，不撑大弹窗） ===== */
.role-dialog >>> .el-tree {
  background: #fff;
  padding: 6px;
  border-radius: 8px;
  border: 1px solid #e4e9f0;
  max-height: 380px;
  overflow-y: auto;
}

.role-dialog >>> .el-tree-node__content {
  height: 26px;
  border-radius: 4px;
}

.role-dialog >>> .el-tree-node__content:hover {
  background-color: #f0f6ff;
}

/* ===== 底部按钮固定 ===== */
.role-dialog >>> .el-dialog__footer {
  position: sticky;
  bottom: 0;
  background: #fff;
  padding-top: 12px;
  padding-bottom: 12px;
  box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.05);
  z-index: 10;
}
</style>