<template>
  <div>
    <el-card>
      <div class="userheader">
        <div>
          <el-button type="primary" icon="el-icon-plus" @click="showDialog('添加账号')">添加账号</el-button>
          <el-button type="danger" size="mini" :disabled="ids.length==0?true:false" @click="showTableDialog">批量删除</el-button>
        </div>
        <div>
          <el-input placeholder="请输入账号名" style="width: 400px;" v-model="tempSearchObj.username" prefix-icon="el-icon-search" @keyup.enter.native="search">
            <el-button slot="append" icon="el-icon-search" @click="search"></el-button>
          </el-input>
          <el-button icon="el-icon-refresh" circle style="margin-left: 10px;" size="mini" @click="resetSearch"></el-button>
        </div>
      </div>
      <el-table  :data="userlist" style="width: 100%;margin-top: 10px;" border  stripe v-loading="loading" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" align="center" 
        :selectable="row => row.name !== 'datacenter' && row.name !== 'rdcenter'">
        </el-table-column>
        <el-table-column type="index" label="序号" width="80" align="center">
        </el-table-column>
        <el-table-column prop="name" label="账号名" width="width">
        </el-table-column>
        <el-table-column prop="home" label="家目录" width="width">
        </el-table-column>
        <el-table-column prop="passwordExpires" label="密码过期时间" width="width">
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="width">
        </el-table-column>
        <el-table-column label="操作" width="400" align="center">
          <template slot-scope="{row}">
            <el-tooltip class="item" effect="dark" content="修改密码" placement="top" :open-delay="1000" :enterable="false" :visible-arrow="false">
              <el-button type="warning" icon="el-icon-edit" size="mini" @click="showDialog('修改密码',row)"></el-button>
            </el-tooltip>
            <el-tooltip class="item" effect="dark" content="删除账号" placement="top" :open-delay="1000" :enterable="false" :visible-arrow="false">
              <el-popconfirm :title="`确定要删除'${row.name}'吗?`" @confirm="deleteUser(row.id)">
                <el-button type="danger" icon="el-icon-delete" size="mini" slot="reference" style="margin-left: 20px;" 
                :disabled="row.name === 'datacenter' || row.name === 'rdcenter'"></el-button>
              </el-popconfirm>
            </el-tooltip>
            <el-tooltip class="item" effect="dark" content="连接SFTP" placement="top" :open-delay="1000" :enterable="false" :visible-arrow="false">
              <el-button type="success" icon="el-icon-connection" size="mini" style="margin-left: 20px;" @click="showSftpDialog(row)"></el-button>
            </el-tooltip>
            <el-button type="primary" icon="el-icon-message" size="mini" style="margin-left: 20px;" @click="emailDialog = true"></el-button>
          </template>
        </el-table-column>
      </el-table>
      <!-- 分页器 -->
      <el-pagination 
        @current-change="getUserList"
        @size-change="handleSizeChange"
        style="margin-top: 20px;text-align: center;" 
        :current-page="page" 
        :page-sizes="[3, 5, 10]" 
        :page-size="limit" 
        layout="sizes, prev, pager, next, jumper,->,total" 
        :total="total">
      </el-pagination>
    </el-card>
    <!-- 添加或修改账号弹框 -->
    <el-dialog :title="title" :visible.sync="dialogFormVisible" width="30%" @close="close" center :close-on-click-modal="false">
      <el-form :model="userForm" :rules="rules" ref="ruleForm">
        <el-form-item label="账号" label-width="80px" prop="name">
          <el-input v-model="userForm.name" autocomplete="off" :disabled="title=='添加账号'?false:true" style="width: 80%;" v-focus="title=='添加账号'"></el-input>
        </el-form-item>
        <el-form-item label="登录类型" label-width="80px" prop="loginType">
          <el-radio-group v-model="userForm.loginType" size="mini" @input="showEmailForm()">
            <el-radio-button label="Password"></el-radio-button>
            <el-radio-button label="KeyFile"></el-radio-button>
            <el-radio-button label="both"></el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-collapse-transition>
          <div v-if="userForm.loginType != 'KeyFile'">
            <el-form-item label="密码" label-width="80px" prop="password">
              <el-input v-model="userForm.password" autocomplete="off" type="password" show-password style="width: 80%;" v-focus="title=='修改密码'"></el-input>
              <el-tooltip class="item" effect="dark" content="生成随机密码" placement="top">
                <el-button type="primary" style="margin-left: 10px;" @click="generatePassword" icon="el-icon-key" circle=""></el-button>
              </el-tooltip>
            </el-form-item>
            <el-form-item label="确认密码" label-width="80px" prop="checkPass">
              <el-input v-model="userForm.checkPass" autocomplete="off" type="password" show-password style="width: 80%;"></el-input>
            </el-form-item>
          </div>
        </el-collapse-transition>
        <el-form-item label="永不过期" label-width="80px" prop="noExpire">
            <el-switch v-model="userForm.noExpire"></el-switch>
        </el-form-item>
        <el-form-item label="发送邮件" label-width="80px" prop="emailOrNot">
            <el-switch v-model="userForm.emailOrNot" @change="change"></el-switch>
        </el-form-item>
        <el-collapse-transition>
          <div v-if="userForm.emailOrNot">
            <el-form-item label="收件人" label-width="80px"  prop="to">
              <!-- <el-input v-model="userForm.to" autocomplete="off" type="textarea" autosize placeholder="多个收件人用','隔开"></el-input> -->
              <el-select
                style="width: 80%;"
                v-model="userForm.to"
                size="medium"
                multiple
                filterable
                allow-create
                default-first-option
                placeholder="请选择或输入收件人">
                <el-option
                  v-for="item in options"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value">
                </el-option>
              </el-select>
            </el-form-item>
            <el-form-item label="抄送人" label-width="80px" prop="cc">
              <!-- <el-input v-model="userForm.cc" autocomplete="off" type="textarea" autosize placeholder="多个抄送人用','隔开"></el-input> -->
              <el-select
                style="width: 80%;"
                v-model="userForm.cc"
                size="medium"
                multiple
                filterable
                allow-create
                default-first-option
                placeholder="请选择或输入抄送人">
                <el-option
                  v-for="item in options"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value">
                </el-option>
              </el-select>
            </el-form-item>
          </div>
        </el-collapse-transition>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false">取 消</el-button>
        <el-button type="primary" @click="addOrUpdateUser()" :loading="buttonLoading">{{ buttonLoading ? '提交中 ...' : '确 定' }}</el-button>
      </div>
    </el-dialog>
    <!-- 批量删除的弹框 -->
    <el-dialog title="批量删除" :visible.sync="dialogTableVisible">
      <el-table height="250" border :data="multipleSelection">
        <el-table-column type="index" label="序号" width="150" align="center"></el-table-column>
        <el-table-column property="name" label="账号" width="200"></el-table-column>
        <el-table-column property="home" label="家目录"></el-table-column>
      </el-table>
      <el-alert title="注意" description="以上账号将被永久删除" type="error" style="margin-top: 10px;" :closable="false" show-icon></el-alert>
      <div slot="footer" class="dialog-footer">
        <el-button @click="closedialogTable">取 消</el-button>
        <el-button type="primary" @click="deleteUsers">确 定</el-button>
      </div>
    </el-dialog>
    <!-- 发送邮件的抽屉 :before-close="handleClose" -->
    <el-drawer
      title="发送邮件(开发中......)"
      :visible.sync="emailDialog"
      direction="rtl"
      custom-class="demo-drawer"
      ref="drawer"
      size="550px"
      >
      <div class="demo-drawer__content" style="margin: 40px">
        <el-form>
          <el-form-item label="主题" label-width="80px">
            <el-input autocomplete="off"></el-input>
          </el-form-item>
          <el-form-item label="收件人" label-width="80px">
            <el-input autocomplete="off"></el-input>
          </el-form-item>
          <el-form-item label="抄送人" label-width="80px">
            <el-input autocomplete="off"></el-input>
          </el-form-item>
        </el-form>
        <div class="demo-drawer__footer">
          <el-button @click="cancelEmailForm">取 消</el-button>
          <el-button type="primary" @click="$refs.drawer.closeDrawer()" :loading="buttonLoading">{{ buttonLoading ? '提交中 ...' : '确 定' }}</el-button>
        </div>
      </div>
    </el-drawer>
    <!-- sftp登录弹框 -->
    <el-dialog title="SFTP登录" :visible.sync="SftpDialogFormVisible" center width="30%" :close-on-click-modal="false" @close="closeSftpDialogForm()">
      <el-card>
        <el-tabs v-model="activeName">
          <el-tab-pane label="密码登录" name="password">
            <el-form :model="SftpForm">
              <el-form-item label="SFTP账号" :label-width="SftpFormLabelWidth">
                <el-input v-model="SftpForm.username" autocomplete="off" :disabled="true"></el-input>
              </el-form-item>
              <el-form-item label="SFTP密码" :label-width="SftpFormLabelWidth">
                <el-input v-model="SftpForm.password" autocomplete="off" type="password" show-password v-focus="SftpDialogFormVisible" @keyup.enter.native="sftplogin()"></el-input>
              </el-form-item>
            </el-form>
          </el-tab-pane>
          <el-tab-pane label="密钥登录" name="keyfile">
            <el-form :model="SftpForm">
              <el-form-item label="SFTP账号" :label-width="SftpFormLabelWidth">
                <el-input v-model="SftpForm.username" autocomplete="off" :disabled="true"></el-input>
              </el-form-item>
              <el-form-item label="SFTP密钥" :label-width="SftpFormLabelWidth">
                <el-upload
                  class="upload-key"
                  ref="upload"
                  :action="KeyFileUploadUrl"
                  :data="{ username: SftpForm.username }"
                  :headers="uploadHeaders"
                  :file-list="keyFileList"
                  :limit=1
                  :before-upload="beforeUploadKey"
                  :on-change="handleKeyFileChange"
                  :on-success="KeyhandleUploadSuccess"
                  :on-error="KeyhandleUploadError"
                  :auto-upload="false">
                  <el-button slot="trigger" size="small" type="primary">选择密钥文件</el-button>
                </el-upload>
              </el-form-item>
            </el-form>
          </el-tab-pane>
        </el-tabs>
      </el-card>
      <div slot="footer" class="dialog-footer">
        <el-button @click="closeSftpDialogForm()">取 消</el-button>
        <el-button type="primary" @click="sftplogin()" :loading="buttonLoading">{{ buttonLoading ? '提交中 ...' : '确 定' }}</el-button>
      </div>
    </el-dialog>
    <!-- SFTP浏览器 -->
    <el-dialog title="SFTP浏览器" :visible.sync="SftpBrowserVisible" center width="width" top="10vh" :close-on-click-modal="false" @close="closeSftpBrowser()">
      <div class="sftp-browser">
        <!-- 面包屑导航 -->
        <el-breadcrumb separator-class="el-icon-arrow-right">
          <el-breadcrumb-item 
            v-for="(item, index) in breadcrumb" 
            :key="index"
            ><span @click="handleBreadcrumbClick(item, index)" :class="breadcrumb.length - 1 == index ? 'breadcrumbBold' : 'breadcrumb'">{{ item.name }}</span>
          </el-breadcrumb-item>
        </el-breadcrumb>

        <!-- 操作按钮 -->
        <div class="operate">
          <!-- 返回上一级 -->
          <el-page-header @back="goBack" :content="breadcrumb[breadcrumb.length - 1].name" style="margin:20px 0;" class="goBack">
          </el-page-header>
          <div style="width: 220px;">
            <el-progress 
              v-if="showUploadProgress"
              :percentage="uploadPercent" 
              size="mini" 
              :color="customColors" 
              class="upload-progress"
              ></el-progress>  
          </div>
          <div>
            <!-- 上传按钮 -->
            <el-upload 
              class="upload"
              :action="uploadUrl"
              :data="{ path: currentPath }"
              :headers="uploadHeaders"
              :on-success="handleUploadSuccess"
              :on-error="handleUploadError"
              :before-upload="beforeUpload"
              :show-file-list="false"
              multiple
              :on-progress="handleUploadProgress"
            >
              <el-button type="primary" size="mini" icon="el-icon-document-add" round>上传文件</el-button>
            </el-upload>
            <!-- 创建目录按钮 -->
            <el-button type="primary" size="mini" icon="el-icon-folder-add" round @click="showCreateFolderDialog = true">创建目录</el-button>
          </div>
        </div>
        
        <!-- 文件列表 -->
        <el-card shadow="hover">
          <el-empty description="这个目录很穷，什么都没有_(:3 ∠)_" v-if="fileList == null ? true : false" style="height: 500px;"></el-empty>
          <el-table :data="fileList" v-loading="isLoading" height="500" v-if="fileList !== null ? true : false" border>
            <el-table-column prop="name" label="名称">
              <template slot-scope="{row}">
                <div v-if="row.isRenaming">
                  <el-input v-model="row.editName" size="mini" @keyup.enter.native="confirmRename(row)" @blur="cancelRename(row)" v-focus="true"></el-input>
                </div>
                <div v-else @click="handleItemClick(row)" 
                    :class="{'dir-item': row.isDir, 'file-item': !row.isDir}">
                  <i :class="row.isDir ? 'el-icon-folder' : 'el-icon-document'"></i>
                  {{ row.name }}
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="size" label="大小" width="120">
              <template slot-scope="{row}">
                {{ row.isDir ? '-' : formatSize(row.size) }}
              </template>
            </el-table-column>
            <el-table-column prop="modified" label="修改时间" width="200">
              <template slot-scope="{row}">
                {{ formatDate(row.modified) }}
              </template>
            </el-table-column>
            <el-table-column  align="center" label="操作" width="150" fixed="right">
              <template slot-scope="{row}">
                <el-button v-if="!row.isDir" size="mini" type="primary" @click="handleDownload(row)" circle icon="el-icon-download"></el-button>
                <el-button size="mini" circle icon="el-icon-edit" @click="startRename(row)"></el-button>
                <el-button size="mini" type="danger" @click="handleDelete(row)"  circle icon="el-icon-delete"></el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </div>
      <!-- 新建文件夹对话框 -->
      <el-dialog title="新建文件夹" :visible.sync="showCreateFolderDialog" width="30%" append-to-body :close-on-click-modal="false">
        <el-form>
          <el-form-item label="文件夹名称">
            <el-input v-model="newFolderName" autocomplete="off"></el-input>
          </el-form-item>
        </el-form>
        <div slot="footer">
          <el-button @click="showCreateFolderDialog = false">取消</el-button>
          <el-button type="primary" @click="createFolder">确定</el-button>
        </div>
      </el-dialog>
      <!-- 删除确认弹框 -->
      <el-dialog
        title="确认删除"
        :visible.sync="deleteDialogVisible"
        width="30%"
        append-to-body :close-on-click-modal="false"
        >
        <span>确定要删除 {{ deleteTarget.name }} 吗？</span>
        <span slot="footer">
          <el-button @click="deleteDialogVisible = false">取消</el-button>
          <el-button type="danger" @click="confirmDelete">确定</el-button>
        </span>
      </el-dialog>
    </el-dialog>
  </div>
</template>

<script>
export default {
  name:'SftpUser',
  data() {
    // 密码验证器
    var validatePass = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请输入密码'));
      } else if (value.length < 14) {
        callback(new Error('密码长度至少为14位'));
      } else {
        // 检查是否包含大小写字母和符号
        let hasLowercase = /[a-z]/.test(value);
        let hasUppercase = /[A-Z]/.test(value);
        let hasSymbol = /[^\w]/.test(value);
        let hasNumber = /\d/.test(value);
        if (!hasLowercase ||!hasUppercase ||!hasSymbol || !hasNumber) {
          callback(new Error('密码必须包含大小写字母数字和符号'));
        } else {
          if (this.userForm.checkPass!== '') {
            this.$refs.ruleForm.validateField('checkPass');
          }
          callback();
        }
      }
    };
    // 确认密码验证器
    var validatePass2 = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请再次输入密码'));
      } else if (value !== this.userForm.password) {
        callback(new Error('两次输入密码不一致!'));
      } else {
        callback();
      }
    };
    return {
      loading:false, //是否显示加载效果
      userlist:[], //用户列表
      dialogFormVisible:false,
      title:'', //dialog标题
      dialogTableVisible:false,
      emailDialog:false,
      buttonLoading:false,
      showPassword:false,
      userForm:{
        id:'', // 用户id
        name:'', // 账号名
        loginType:'Password', // 登录类型
        noExpire:false, // 永不过期
        emailOrNot:false,  //是否发送邮件
        password:'',   // 密码
        checkPass:'', // 确认密码
        to:'',   //收件人
        cc:'',  //抄送人
      },
      // 邮箱相关
      options:[],
      // SFTP登录
      SftpDialogFormVisible:false,
      SftpFormLabelWidth: '100px',
      SftpForm:{
        username:'',
        password:''
      },
      // 整理的email数据
      email:{
        to:[], 
        cc:[],
        subject:'',
      },
      searchObj:{
        username:''
      },
      tempSearchObj:{
        username:''
      },
      ids:[], // 收集多选的用户id列表
      multipleSelection:[], //多项的用户信息
      // 分页器
      page:1,
      limit:5,
      total:0,
      // 验证规则
      rules:{
        name:[
          {required: true, message: '请输入账号名称', trigger: 'blur'},
          { 
            // pattern: /^[a-zA-Z]([a-zA-Z0-9_]*)$/,
            // message: '账号必须以字母开头且不能有空格和符号',
            pattern: /^[a-zA-Z][a-zA-Z0-9_]*$|^[0-9](?=.*[a-zA-Z])[a-zA-Z0-9_]*$/,
            message: '账号须以字母/数字开头（必须包含字母）仅允许字母、数字和下划线组成',
            trigger: ['blur','change']
          }
        ],
        loginType:[{required: true, message: '请选择登录类型', trigger: 'blur'}],
        password:[
          {required: true, message: '请输入密码', trigger: 'blur'},
          { validator: validatePass, trigger: 'blur' }
        ],
        checkPass: [
          {required: true, message: '请确认密码', trigger: 'blur'},
          { validator: validatePass2, trigger: ['blur','change'] }
        ],
        to:[{required: true, message: '请输入邮箱地址', trigger: 'blur' }],
        cc:[{required: false, message: '请输入邮箱地址', trigger: 'blur' }],
      },
      // sftp浏览器相关数据
      SftpBrowserVisible: false,
      activeName: 'password', // SFTP登录方式
      // 上传进度相关
      showUploadProgress: false, // 是否显示上传进度
      uploadPercent: 0, // 上传进度百分比
      customColors: [
          {color: '#f56c6c', percentage: 20},
          {color: '#e6a23c', percentage: 40},
          {color: '#5cb87a', percentage: 60},
          {color: '#1989fa', percentage: 80},
          {color: '#6f7ad3', percentage: 100}
        ],
      // 密钥登录
      keyFileList: [], // 密钥文件列表
      KeyFileUploadUrl: '/dev-api/sftp/login', // 密钥登录接口
      currentPath: '/', // 当前路径
      fileList: [], // 当前目录下的文件、目录列表
      breadcrumb: [{
        name: '根目录',
        path: '/'
      }],
      isLoading: false,
      uploadUrl: '/dev-api/sftp/upload', // 上传文件的接口地址
      // 添加认证头
      uploadHeaders: {
        'Token': `${this.$store.state.user.token}`,
        'X-SFTP-Token': ''
      },
      showCreateFolderDialog: false, // 是否显示新建文件夹对话框
      newFolderName: '', // 新建文件夹的名称
      deleteDialogVisible: false,
      deleteTarget: {
        name: '',
        path: '',
        isDir: false
      },
      renamingItem:null, //正在重命名的文件或目录
    }
  },
  mounted(){
    this.getUserList()
  },
  methods:{
    // 获取用户列表
    async getUserList(pages = 1){
      this.page = pages
      const {page,limit,searchObj} = this
      try {
        this.loading = true
        const result = await this.$API.sftpuser.reqUserList(page,limit,searchObj)
        if(result.code == 200){
          this.userlist = result.data.users
          this.total = result.data.total
        }
        this.loading = false
      } catch (error) {}
    },
    // 显示对话框,并清除数据
    showDialog(title,row){
      // 设置对话框标题
      this.title = title
      // 展示对话框
      this.dialogFormVisible = true
      // 清除数据
      this.userForm = {
        id:'',
        name:'',
        loginType:'Password',
        noExpire:false,
        emailOrNot:false,
        password:'',
        checkPass:'',
        to:'',
        cc:'',
      }
      // 打开后立即重置验证（确保没有残留提示）
      this.$nextTick(() => {
        if (this.$refs.ruleForm) {
          this.$refs.ruleForm.clearValidate();
          // 如果需要重置字段值（可选）
          // this.$refs.ruleForm.resetFields();
        }
      });
      // 显示准备修改的账号名
      if(title == '修改密码'){
        this.userForm.name = row.name
        this.userForm.id = row.id
      }
    },
    // 添加账号和修改密码的确定按钮的回调
    addOrUpdateUser(){
      this.buttonLoading = true
      this.$refs['ruleForm'].validate(async (valid) =>{
        if (valid) {
          const {id,name,loginType,password,emailOrNot,noExpire,to,cc} = this.userForm
          let message;
          if (loginType === 'KeyFile') {
            message = '密钥修改成功';
          } else if (loginType === 'Password') {
            message = '密码修改成功';
          } else {
            message = '密码和密钥修改成功';
          }
          try {
            if(this.title == '修改密码'){
              const result = await this.$API.sftpuser.reqUpdateUser({id,name,loginType,password,emailOrNot,noExpire})
              if(result.code == 200){
                if(this.userForm.emailOrNot){
                  this.email.to = to
                  this.email.cc = cc
                  this.email.subject = `${name}`
                  await this.$API.sftpuser.reqSendEmail(this.email)
                  this.$message({type:'success',message:`账号:${name} ${message},邮件已发送`})
                }else{
                  this.$message({type:'success',message:`账号:${name} ${message}`})
                }
                this.getUserList(this.page)
                this.dialogFormVisible = false
                this.buttonLoading = false
              }
            }else{
              const result = await this.$API.sftpuser.reqAddUser({name,loginType,password,emailOrNot,noExpire})
              if(result.code == 200){
                if(this.userForm.emailOrNot){
                  this.email.to = to
                  this.email.cc = cc
                  this.email.subject = `${name}`
                  await this.$API.sftpuser.reqSendEmail(this.email)
                  this.$message({type:'success',message:`账号:${name} 添加成功,邮件已发送`})
                }else{
                  this.$message({type:'success',message:`账号:${name} 添加成功`})
                }
                this.getUserList(this.page)
                this.dialogFormVisible = false
                this.buttonLoading = false
              }
            }
          } catch (error) {this.buttonLoading = false}
        }else {
          console.log('error submit!!');
          this.buttonLoading = false
          return false;
        }
      })
    },
    // Dialog 关闭的回调
    close(){
      this.title = ''
      // 清除数据,需要给表单加上prop属性
      this.$refs['ruleForm'].clearValidate();
      this.$refs['ruleForm'].resetFields();
      this.dialogFormVisible = false;
    },
    closedialogTable(){
      this.dialogTableVisible = false
    },
    // 分页器limit改变的回调
    handleSizeChange(limit) {
      // 修改参数
      this.limit = limit;
      this.getUserList(this.page);
    },
    // 搜索的回调
    search(){
      this.searchObj = { ...this.tempSearchObj };
      const {username} = this.searchObj
      if(username == ''){
        this.$message({type:'warning',message:'请输入账号名'})
      }else{
        this.getUserList()
      }
    },
    /* 
    重置输入后搜索
    */
    resetSearch() {
      this.searchObj = {
        username: "",
      };
      this.tempSearchObj = {
        username: "",
      };
      this.getUserList();
    },
    // 删除一个用户
    async deleteUser(id){
      try {
        const {page} = this
        const result = await this.$API.sftpuser.reqDeleteUser(id)
        if(result.code == 200){
          this.getUserList(this.userlist.length==1?page-1:page)
          this.$message({type:"success",message:"删除成功"})
        }
      } catch (error) {}
    },
    // 存储选择的用户id
    handleSelectionChange(val){
      this.ids = []
      this.multipleSelection = val
      val.forEach(item => {
        this.ids.push(item.id)
      });
    },
    showTableDialog(){
      this.dialogTableVisible = true
    },
    async deleteUsers(){
      try {
        const {ids,page} = this
        const result = await this.$API.sftpuser.reqDeleteUsers(ids)
        if(result.code == 200){
          this.dialogTableVisible = false
          this.$message({type:'success',message:'批量删除成功'})
          this.getUserList(this.userlist.length==this.ids.length?page-1:page)
        }
      } catch (error) {}
    },
    // 取消email
    cancelEmailForm(){
      this.buttonLoading = false;
      this.emailDialog = false;
    },
    // 生成密码的回调
    generatePassword(){
      // 生成随机密码,长度14位,至少有一个小写字母,一个大写字母和一个符号和一个数字
      const length = 14;
      // 去掉了易混淆的字符
      const lowercaseChars = 'abcdefghjkmnpqrstuvwxyz';
      const uppercaseChars = 'ABCDEFGHJKLMNPQRSTUVWXYZ';
      const numberChars = '123456789';
      const symbolChars = '!@#$%^&*+~:;?><,.-=';
      let password = '';
      password += lowercaseChars[Math.floor(Math.random() * lowercaseChars.length)]; // 至少一个小写字母
      password += uppercaseChars[Math.floor(Math.random() * uppercaseChars.length)]; // 至少一个大写字母
      password += numberChars[Math.floor(Math.random() * numberChars.length)]; // 至少一个数字
      password += symbolChars[Math.floor(Math.random() * symbolChars.length)]; // 至少一个符号
      // 生成剩余的字符串
      for (let i = 4; i < length; i++) {
        const charSet = lowercaseChars + uppercaseChars + numberChars + symbolChars;
        password += charSet[Math.floor(Math.random() * charSet.length)];
      }
      this.userForm.password = password;
      this.userForm.checkPass = password;
    },
    //  发送邮件变化的事件回调
    async change(value){
      if(value){
        // 获取联系人选项
        try {
          const result = await this.$API.contact.reqContactOptions()
          if(result.code == 200){
            this.options = result.data.options
          }
        } catch (error) {
        }
        // 设置默认抄送人,可根据需求修改
        // this.userForm.cc = 'lw.yang@ho-brostech.com,bingjia.zheng@ho-brostech.com,yinling.chen@ho-brostech.com,it@ho-brostech.com'
      }
      if(!value){
        this.userForm.cc = ''
      }
    },
    // 显示SFTP登录对话框
    showSftpDialog(row){
      this.SftpDialogFormVisible = true
      this.SftpForm.username = row.name
      this.SftpForm.password = ''
    },
    // 发送请求，登录sftp
    async sftplogin(){
      // 根据登录类型发送请求
      if(this.activeName == 'password'){
        // 密码登录
        // 验证密码是否为空
        const {username,password} = this.SftpForm
        if(password == ''){
          this.$message({type:'warning',message:'请输入SFTP密码'})
          return
        }
        // loading
        this.buttonLoading = true
        try {
          const result = await this.$API.sftpuser.reqSftpLogin({username,password})
          if(result.code == 200){
            // 存储SFTP-Token
            localStorage.setItem("sftp_token", result.data.sftp_token)
            this.uploadHeaders['X-SFTP-Token'] = result.data.sftp_token
            this.$message({type:'success',message:'SFTP登录成功'})
            this.fetchFiles(); // 获取文件列表
          }
        }catch (error) {}finally{
          this.SftpDialogFormVisible = false
          this.buttonLoading = false
        }
      }else{
        // 密钥登录
        // 验证是否选择了密钥文件
        if(this.keyFileList.length == 0){
          this.$message({type:'warning',message:'请选择密钥文件'})
          return
        }
        // loading
        this.buttonLoading = true
        // 上传密钥文件
        this.$refs.upload.submit();
      }
    },
    handleKeyFileChange(file, fileList) {
      this.keyFileList = fileList;
    },
    // 上传前检查
    beforeUploadKey(file) {
      // 检查文件大小是否超过100MB
      const isLt100M = file.size / 1024 / 1024 < 100;
      if (!isLt100M) {
        this.$message.error('上传文件大小不能超过 100MB!');
      }
      return isLt100M;
    },
    // 密钥登录上传成功的回调
    KeyhandleUploadSuccess(response, file, fileList) {
      if (response.code === 200) {
        this.$message({ type: 'success', message: 'SFTP登录成功' });
        if (response.data && response.data.sftp_token) {
          localStorage.setItem("sftp_token", response.data.sftp_token);
          this.uploadHeaders['X-SFTP-Token'] = response.data.sftp_token
        }
        this.fetchFiles(); // 获取文件列表
      } else {
        this.$message({ type: 'error', message: response.message || 'SFTP登录失败' });
      }
      this.SftpDialogFormVisible = false;
      this.buttonLoading = false;
      this.keyFileList = []; // 清空文件列表
      this.$refs.upload.clearFiles(); // 清除上传组件的文件列表
    },
    // 密钥登录上传失败的回调
    KeyhandleUploadError(response, file, fileList) {
      this.$message({ type: 'error', message: response.message || 'SFTP登录失败' });
      this.SftpDialogFormVisible = false;
      this.buttonLoading = false;
      this.keyFileList = []; // 清空文件列表
      this.$refs.upload.clearFiles(); // 清除上传组件的文件列表
    },
    // 关闭SFTP登录对话框
    closeSftpDialogForm(){
      this.SftpDialogFormVisible = false
      this.buttonLoading = false
      this.activeName = 'password'
      this.keyFileList = [] // 清空文件列表
      this.$refs.upload.clearFiles() // 清除上传组件的文件列表
    },
    // 发送请求获取文件和目录列表
    async fetchFiles(path = this.currentPath) {
      this.SftpBrowserVisible = true;
      this.isLoading = true;
      try {
        const result = await this.$API.sftpuser.reqSftpFiles({path});
        
        if (result.code === 200) {
          this.fileList = result.data.files;
          this.currentPath = result.data.path;
          this.updateBreadcrumb();
        }
      } catch (e) {
        this.$message.error('获取文件列表失败');
      } finally {
        this.isLoading = false;
      }
    },

    // 更新面包屑导航
    updateBreadcrumb() {
      if (this.currentPath === '/') {
        this.breadcrumb = [{ name: '根目录', path: '/' }];
        return;
      }

      const paths = this.currentPath.split('/').filter(Boolean);

      // 生成面包屑导航
      this.breadcrumb = paths.reduce((acc, cur, index) => {
        // 获取累加数组的最后一个面包屑项
        const lastItem = acc[acc.length - 1];

        // 构建新路径
        const path = lastItem.path == '/'
        ? lastItem.path + cur 
        : lastItem.path + '/' + cur; 

        // 返回新的面包屑项
        return [...acc, { name: cur, path }];
      }, [{ name: '根目录', path: '/' }]);
    },
    // 点击文件/目录
    handleItemClick(item) {
      if (item.isDir) {
        this.fetchFiles(item.path); // 进入目录
      } else {
        this.isFile = true;
        this.currentPath = item.path;
        // 此处可调用文件预览/下载方法
      }
    },

    // 点击面包屑导航
    handleBreadcrumbClick(item, index) {
      if (index === this.breadcrumb.length - 1) return;
      this.fetchFiles(item.path);
    },

    // 返回上一级
    goBack() {
      if (this.currentPath === '/') return;
      
      const paths = this.currentPath.split('/').filter(Boolean);
      paths.pop();
      const newPath = paths.length ? `/${paths.join('/')}` : '/';
      
      this.fetchFiles(newPath);
    },

    // 工具函数：格式化文件大小
    formatSize(bytes) {
      if (bytes === 0) return '0 B';
      const k = 1024;
      const sizes = ['B', 'KB', 'MB', 'GB'];
      
      // const i = Math.floor(Math.log(bytes) / Math.log(k));
      let i = 0;
      let size = bytes;
      
      while (size >= k && i < sizes.length - 1) {
        size /= k;
        i++;
      }
      return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    },

    // 工具函数：格式化日期
    formatDate(dateString) {
      return new Date(dateString).toLocaleString();
    },
    // 关闭 SFTP 浏览器时重置面包屑和文件列表，断开sftp连接
    async closeSftpBrowser() {
      this.breadcrumb = [{ name: '根目录', path: '/' }];
      this.currentPath = '/';
      this.fileList = [];
      this.isLoading = false;
      // 新增：重置上传进度条
      this.resetUploadProgress();
      try {
        const result = await this.$API.sftpuser.reqSftpLogout();
        if (result.code === 200) {
          this.$message({ type: 'success', message: result.message });
        } else {
          this.$message({ type: 'error', message: result.message });
        }
      } catch (error) {}
    },

    // 上传前检查
    beforeUpload(file) {
      const isLt1G = file.size / 1024 / 1024 < 1024;
      if (!isLt1G) {
        this.$message.error('上传文件大小不能超过 1GB!');
      }
      return isLt1G;
    },
    
    // 上传成功处理
    handleUploadSuccess(response, file, fileList) {
      if (response.code === 200) {
        this.$message.success('文件上传成功');
        this.fetchFiles(); // 刷新文件列表
        // 重置进度条（单文件上传直接重置，多文件上传判断是否全部完成）
      if (fileList.every(item => item.status === 'success' || item.status === 'fail')) {
        this.resetUploadProgress();
      }
      } else {
        this.$message.error(response.message || '文件上传失败');
      }
    },
    
    // 上传失败处理
    handleUploadError(err, file, fileList) {
      // 重置进度条
      this.resetUploadProgress();
    },
    // 新增：重置上传进度条的公共方法
    resetUploadProgress() {
      this.showUploadProgress = false;
      this.uploadPercent = 0;
      this.isUploading = false;
    },

    // 创建文件夹
    async createFolder() {
      if (!this.newFolderName) {
        this.$message.error('请输入文件夹名称');
        return;
      }
      
      try {
        const result = await this.$API.sftpuser.reqSftpMkdir({
          path: this.currentPath,
          name: this.newFolderName
        })
        
        if (result.code === 200) {
          this.$message.success('文件夹创建成功');
          this.showCreateFolderDialog = false;
          this.newFolderName = '';
          this.fetchFiles(); // 刷新文件列表
        } else {
          this.$message.error(result.message || '文件夹创建失败');
        }
      } catch (err) {}
    },

    // 下载文件
    async handleDownload(file) {
      // 检查是否为目录，目录不支持下载
      if (file.isDir) {
        this.$message.warning('请选择文件进行下载')
        return
      }

      try {
        // 调用后端下载API，传入文件路径参数
        // 注意：reqSftpDownload方法已配置responseType: 'blob'
        const response = await this.$API.sftpuser.reqSftpDownload({ path: file.path });

        // 创建一个Blob URL用于下载
        // 1. 将响应数据转换为Blob对象
        // 2. 生成临时URL供下载使用
        const url = window.URL.createObjectURL(new Blob([response.data]));

        // 创建隐藏的<a>标签用于触发下载
        const link = document.createElement('a');
        link.href = url;
        
        // 设置下载文件名（从文件对象中获取）
        // 注意：如果文件名包含特殊字符，可能需要额外处理
        link.download = file.name

        // 将链接添加到DOM中（必须添加到文档才能工作）
        document.body.appendChild(link);
        // 触发下载
        link.click();

        // 下载完成后清理资源：
        // 1. 从DOM中移除<a>标签
        // 2. 释放Blob URL（避免内存泄漏）
        document.body.removeChild(link);
        window.URL.revokeObjectURL(url)

      } catch (err) {
        this.$message.error('下载失败: ' + (err.message || '未知错误'))
      }
    },

    // 删除文件或目录按钮的回调
    handleDelete(file) {
      this.deleteTarget = {
        name: file.name,
        path: file.path,
        isDir: file.isDir
      };
      this.deleteDialogVisible = true;
    },

    // 确认删除
    async confirmDelete(){
        try {
        const result = await this.$API.sftpuser.reqSftpDelete({
          path: this.deleteTarget.path
        })
        if (result.code === 200) {
          this.$message.success('删除成功');
          this.fetchFiles(); // 刷新文件列表
        } else {
          this.$message.error(result.message || '删除失败');
        }
      } catch (e) {} 
      finally {
        this.deleteDialogVisible = false;
      }
    },

    // 重命名按钮的回调
    startRename(item){
      // 增加isRenaming属性，显示输入框
      this.$set(item,'isRenaming',true)
      // 增加editName属性，输入框显示当前文件名
      this.$set(item,'editName',item.name)
      this.renamingItem = item
    },

    // 重命名输入框回车，确认重命名
    async confirmRename(item){
      // 名称检查，不能为空
      if (!item.editName || item.editName.trim() === '') {
        this.message.warning('名称不能为空')
        return
      }

      // 不能与原名一样
      if (item.editName === item.name) {
        this.cancelRename(item)
        return
      }

      // 发送请求，重命名
      try {
        const result = await this.$API.sftpuser.reqSftpRename({
          oldPath: item.path,
          newName: item.editName
        })

        if (result.code == 200) {
          this.$message.success('重命名成功')
          // 刷新列表
          this.fetchFiles()
        }else {
          this.$message.error(result.message || '重命名失败')
        }
      } catch (error) {
        this.cancelRename(item)
      }
    },

    // 取消重命名
    cancelRename(item){
      item.isRenaming = false
      item.editName = ''
    },
    // 根据不同登录类型默认显示email表单
    async showEmailForm(){
      if (this.userForm.loginType=='KeyFile' || this.userForm.loginType=='both') {
        this.userForm.emailOrNot=true
        try {
          const result = await this.$API.contact.reqContactOptions()
          if(result.code == 200){
            this.options = result.data.options
          }
        } catch (error) {
        }
      }else{
        this.userForm.emailOrNot=false
        this.userForm.cc = ''
      } 
    },
    // 新增：文件上传进度钩子处理方法
    handleUploadProgress(event, file, fileList) {
      // event.percent 是Element UI返回的上传进度（0-100的数值）
      this.uploadPercent = Math.round(event.percent); // 取整，避免小数抖动
      this.showUploadProgress = true; // 显示进度条
      this.isUploading = true; // 标记正在上传
    },
  }
}
</script>

<style scoped>
.userheader{
  display: flex;
  justify-content: space-between;
}
.demo-drawer__footer{
  display: flex;
  justify-content: flex-end;
}
.breadcrumb {
  cursor: pointer;
}

.breadcrumbBold {
  font-weight: bold;
  color: #409EFF
}

.breadcrumb:hover {
  color: #409EFF;
  font-weight: bold;
}

.goBack:hover {
  color: #409EFF;
  font-weight: bold;
}

.dir-item {
  color: #409EFF;
  cursor: pointer;
}

.file-item {
  color: #606266;
}

.operate {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.upload {
  display: inline-block;
  margin-right: 10px;
}
</style>