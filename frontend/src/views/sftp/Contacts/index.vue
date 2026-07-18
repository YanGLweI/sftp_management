<template>
  <div>
    <el-card>
      <div class="userheader">
        <div>
          <el-button type="primary" icon="el-icon-plus" @click="showDialog('添加联系人')">添加联系人</el-button>
          <el-button type="danger" size="mini" :disabled="ids.length==0?true:false" @click="showTableDialog">批量删除</el-button>
        </div>
        <div>
          <el-input placeholder="请输入联系人姓名" style="width: 400px;" v-model="tempSearchObj.name" prefix-icon="el-icon-search" @keyup.enter.native="search"> 
            <el-button slot="append" icon="el-icon-search" @click="search"></el-button>
          </el-input>
          <el-button icon="el-icon-refresh" circle style="margin-left: 10px;" size="mini" @click="resetSearch"></el-button>
        </div>
      </div>
      <el-table  :data="contactlist" style="width: 100%;margin-top: 10px;" border  stripe v-loading="loading" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" align="center">
        </el-table-column>
        <el-table-column type="index" label="序号" width="80" align="center">
        </el-table-column>
        <el-table-column prop="name" label="姓名" width="width">
        </el-table-column>
        <el-table-column prop="email" label="邮箱" width="width">
        </el-table-column>
        <el-table-column prop="CreatedAt" label="创建时间" width="width" :formatter="formatDateTime">
        </el-table-column>
        <el-table-column prop="UpdatedAt" label="更新时间" width="width" :formatter="formatDateTime">
        </el-table-column>
        <el-table-column label="操作" width="400" align="center">
          <template slot-scope="{row}">
            <el-tooltip class="item" effect="dark" content="修改联系人" placement="top" :open-delay="1000" :enterable="false" :visible-arrow="false">
              <el-button type="warning" icon="el-icon-edit" size="medium" @click="showDialog('修改联系人',row)" circle></el-button>
            </el-tooltip>
            <el-tooltip class="item" effect="dark" content="删除联系人" placement="top" :open-delay="1000" :enterable="false" :visible-arrow="false">
              <el-popconfirm :title="`确定要删除'${row.name}'吗?`" @confirm="deleteContact(row.ID)">
                <el-button type="danger" icon="el-icon-delete" size="medium" slot="reference" style="margin-left: 20px;" circle></el-button>
              </el-popconfirm>
            </el-tooltip>
          </template>
        </el-table-column>
      </el-table>
      <!-- 分页器 -->
      <el-pagination 
        @current-change="getContactList"
        @size-change="handleSizeChange"
        style="margin-top: 20px;text-align: center;" 
        :current-page="page" 
        :page-sizes="[3, 5, 10]" 
        :page-size="limit" 
        layout="sizes, prev, pager, next, jumper,->,total" 
        :total="total">
      </el-pagination>
    </el-card>
    <!-- 添加或修改联系人弹框 -->
    <el-dialog :title="title" :visible.sync="dialogFormVisible" width="30%" @close="close" center :close-on-click-modal="false">
      <el-card>
        <el-form :model="contactForm" :rules="rules" ref="ruleForm">
          <el-form-item label="姓名" label-width="80px" prop="name">
            <el-input v-model="contactForm.name" autocomplete="off" style="width: 80%;" v-focus="title=='添加联系人'"></el-input>
          </el-form-item>
          <el-form-item label="邮箱" label-width="80px" prop="email">
            <el-input v-model="contactForm.email" autocomplete="off" style="width: 80%;" v-focus="title=='修改联系人'"></el-input>
          </el-form-item>
        </el-form>
      </el-card>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false">取 消</el-button>
        <el-button type="primary" @click="addOrUpdateContact()" :loading="buttonLoading">{{ buttonLoading ? '提交中 ...' : '确 定' }}</el-button>
      </div>
    </el-dialog>
    <!-- 批量删除的弹框 -->
    <el-dialog title="批量删除" :visible.sync="dialogTableVisible">
      <el-card>
        <el-table height="250" border :data="multipleSelection">
          <el-table-column type="index" label="序号" width="55" align="center"></el-table-column>
          <el-table-column property="name" label="姓名" width="200"></el-table-column>
          <el-table-column property="email" label="邮箱"></el-table-column> 
        </el-table>
      </el-card>
      <el-alert title="注意" description="以上联系人将被永久删除" type="error" style="margin-top: 10px;" :closable="false" show-icon></el-alert>
      <div slot="footer" class="dialog-footer">
        <el-button @click="closedialogTable">取 消</el-button>
        <el-button type="primary" @click="deleteContacts">确 定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
export default {
  name: "Contacts",
  data() {
    return {
      contactlist:[],
      loading: false,
      searchObj:{
        name:''
      },
      tempSearchObj:{
        name:''
      },
      // 选择的联系人id
      ids:[],
      multipleSelection:[], //多项的联系人信息
      // 分页器
      page:1,
      limit:10,
      total:0,
      // 添加或修改联系人弹框
      dialogFormVisible: false,
      title:'',
      buttonLoading:false,
      contactForm:{
        ID:"",
        name:'',
        email:'',
      },
      rules: {
        name: [
          { required: true, message: '请输入联系人姓名', trigger: 'blur' },
          { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
        ],
        email: [
          { required: true, message: '请输入联系人邮箱', trigger: 'blur' },
          { type: 'email', message: '请输入正确的邮箱地址', trigger: ['blur', 'change'] }
        ]
      },
      dialogTableVisible: false, // 批量删除弹框
    };
  },
  mounted() {
    this.getContactList();
  },
  methods:{
    // 获取联系人列表
    async getContactList(pages = 1){
      this.page = pages
      const {page,limit,searchObj} = this
      try {
        this.loading = true
        const result = await this.$API.contact.reqContactList(page,limit,searchObj)
        if(result.code === 200){
          this.contactlist = result.data.contacts
          this.total = result.data.total
        }
        this.loading = false
      } catch (error) {
      }
    },
     // 搜索的回调
    search(){
      this.searchObj = { ...this.tempSearchObj };
      const {name} = this.searchObj
      if(name == ''){
        this.$message({type:'warning',message:'请输入联系人姓名'})
      }else{
        this.getContactList()
      }
    },
    /* 
    重置输入后搜索
    */
    resetSearch() {
      this.searchObj = {
        name: "",
      };
      this.tempSearchObj = {
        name: "",
      };
      this.getContactList();
    },
        // 存储选择的用户id
    handleSelectionChange(val){
      this.ids = []
      this.multipleSelection = val
      val.forEach(item => {
        this.ids.push(item.ID)
      });
    },
    // 分页器limit改变的回调
    handleSizeChange(limit) {
      // 修改参数
      this.limit = limit;
      this.getContactList(this.page);
    },
    // 时间格式化方法
    formatDateTime(row, column, cellValue, index) {
      if (!cellValue) return '';
      
      // 转换为 Date 对象
      const date = new Date(cellValue);
      
      // 获取年、月、日、时、分、秒
      const year = date.getFullYear();
      const month = this.padZero(date.getMonth() + 1);
      const day = this.padZero(date.getDate());
      const hours = this.padZero(date.getHours());
      const minutes = this.padZero(date.getMinutes());
      const seconds = this.padZero(date.getSeconds());
      
      // 拼接成 "YYYY-MM-DD HH:MM:SS" 格式
      return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
    },
    
    // 补零方法，确保数字为两位数
    padZero(num) {
      return num < 10 ? '0' + num : num;
    },
    // 显示添加或修改联系人弹框
    showDialog(title,row){
      // 重置表单
      this.contactForm = {
        ID:"",
        name:'',
        email:'',
      }
      // 设置标题
      this.title = title
      // 如果是修改联系人，填充表单数据
      this.dialogFormVisible = true
      // 打开后立即重置验证（确保没有残留提示）
      this.$nextTick(() => {
        if (this.$refs.ruleForm) {
          this.$refs.ruleForm.clearValidate();
          // 如果需要重置字段值（可选）
          // this.$refs.ruleForm.resetFields();
        }
      });
      // 显示准备修改的账号名
      if(title == '修改联系人'){
        this.contactForm.name = row.name
        this.contactForm.email = row.email
        this.contactForm.ID = row.ID
      }
    },
    // 显示添加或修改联系人弹框 关闭的回调
    close(){
      this.title = ''
      // 清除数据,需要给表单加上prop属性
      this.$refs['ruleForm'].clearValidate();
      this.$refs['ruleForm'].resetFields();
      this.dialogFormVisible = false;
    },
    // 添加或修改联系人
    addOrUpdateContact(){
      this.$refs.ruleForm.validate(async (valid) => {
        if (valid) {
          this.buttonLoading = true
          const {ID,name,email} = this.contactForm
          try {
            if(this.title == '添加联系人'){
              const result = await this.$API.contact.reqAddContact({name,email})
              if(result.code === 200){
                this.$message({type:'success',message:'添加成功'})
                this.close()
                this.getContactList(this.page)
              }
            }else{
              const result = await this.$API.contact.reqUpdateContact(this.contactForm)
              if(result.code === 200){
                this.$message({type:'success',message:'修改成功'})
                this.close()
                this.getContactList(this.page)
              }
            }
          } catch (error) {
            this.$message({type:'error',message:'操作失败，请稍后再试'})
            this.buttonLoading = false
          }finally{
            this.buttonLoading = false
          }
        }
      })
    },
    // 删除联系人
    async deleteContact(id){
      try {
        const result = await this.$API.contact.reqDeleteContact(id)
        if(result.code === 200){
          this.$message({type:'success',message:'删除成功'})
          this.getContactList(this.page)
        }
      } catch (error) {
        this.$message({type:'error',message:'删除失败，请稍后再试'})
      }
    },
    // 显示批量删除弹框
    showTableDialog(){
      this.dialogTableVisible = true
    },
    // 批量删除联系人
    async deleteContacts(){
      try {
        const {ids,page} = this
        const result = await this.$API.contact.reqDeleteContacts(ids)
        if(result.code === 200){
          this.dialogTableVisible = false
          this.$message({type:'success',message:'批量删除成功'})
          this.getContactList(this.contactlist.length==1?page-1:page)
        }
      } catch (error) {
      }
    },
    // 关闭批量删除弹框
    closedialogTable(){
      this.dialogTableVisible = false
    }
  }
}
</script>

<style scoped>
.userheader{
  display: flex;
  justify-content: space-between;
}
</style>