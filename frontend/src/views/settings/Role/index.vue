<template>
  <div class="app-container">
    <div class="search-box">
      <el-input v-model="searchQuery" placeholder="搜索角色名称" clearable style="width: 200px;" @keyup.enter.native="handleSearch" />
      <el-button type="primary" icon="el-icon-search" @click="handleSearch">搜索</el-button>
      <el-button type="success" icon="el-icon-plus" @click="showCreateDialog">新增角色</el-button>
    </div>

    <el-table v-loading="loading" :data="roleList" border stripe style="width: 100%">
      <el-table-column prop="name" label="角色名称" min-width="150" />
      <el-table-column prop="description" label="角色描述" min-width="200" />
      <el-table-column label="关联菜单数" width="100">
        <template slot-scope="{ row }">
          <el-tag>{{ row.menus ? row.menus.length : 0 }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="关联安全组数" width="120">
        <template slot-scope="{ row }">
          <el-tag type="warning">{{ row.ldapGroupCount || 0 }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="createdAt" label="创建时间" width="200">
        <template slot-scope="{ row }">{{ row.CreatedAt | dateFormat }}</template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template slot-scope="{ row }">
          <el-button type="primary" size="mini" icon="el-icon-edit" @click="showEditDialog(row)">编辑</el-button>
          <el-button type="danger" size="mini" icon="el-icon-delete" :disabled="row.name === '超级管理员'" @click="handleDelete(row)">删除</el-button>
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

    <el-dialog :title="dialogTitle" :visible.sync="dialogVisible" width="600px" @closed="resetDialog">
      <el-form ref="roleForm" :model="roleForm" :rules="roleRules" label-width="100px">
        <el-form-item label="角色名称" prop="name">
          <el-input v-model="roleForm.name" :disabled="isSuperAdmin" placeholder="请输入角色名称" />
        </el-form-item>
        <el-form-item label="角色描述" prop="description">
          <el-input v-model="roleForm.description" :disabled="isSuperAdmin" type="textarea" :rows="3" placeholder="请输入角色描述" />
        </el-form-item>
        <el-form-item label="菜单权限">
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
        </el-form-item>
        <el-form-item label="LDAP安全组">
          <el-button type="primary" size="mini" icon="el-icon-plus" @click="addLDAPGroup">添加安全组</el-button>
          <div v-for="(group, index) in roleForm.ldapGroups" :key="index" style="margin-top: 8px; display: flex; gap: 8px;">
            <el-input v-model="group.group_dn" placeholder="安全组DN，如 CN=IT部,OU=..." style="flex: 2;" />
            <el-input v-model="group.group_name" placeholder="显示名称" style="flex: 1;" />
            <el-button type="danger" icon="el-icon-delete" size="mini" @click="removeLDAPGroup(index)" />
          </div>
        </el-form-item>
      </el-form>
      <span slot="footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">保存</el-button>
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
      this.dialogTitle = '新增角色'
      this.dialogVisible = true
    },
    async showEditDialog(row) {
      this.isEdit = true
      this.editId = row.ID
      this.isSuperAdmin = row.name === '超级管理员'
      this.dialogTitle = this.isSuperAdmin ? '超级管理员（仅可修改LDAP安全组）' : '编辑角色'
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
            if (this.$refs.menuTree) {
              // 超级管理员：默认勾选所有菜单（不可修改）
              if (this.isSuperAdmin) {
                const allMenuKeys = this.collectAllMenuKeys(this.allMenus)
                this.$refs.menuTree.setCheckedKeys(allMenuKeys)
              } else {
                this.$refs.menuTree.setCheckedKeys(this.roleForm.menus)
              }
            }
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
      this.roleForm = { name: '', description: '', menus: [], ldapGroups: [] }
      this.submitLoading = false
      this.isSuperAdmin = false
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
.search-box {
  margin-bottom: 16px;
  display: flex;
  gap: 8px;
  align-items: center;
}
</style>