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
          <el-radio-group v-model="userForm.loginType" size="mini">
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
        <el-form-item label="导入公钥" label-width="80px" prop="publicKey" v-if="userForm.loginType != 'Password'">
          <el-switch v-model="userForm.publicKey"></el-switch>
          <el-upload
            v-if="userForm.publicKey"
            class="upload-key"
            ref="uploadKey"
            action="https://127.0.0.1"
            :file-list="publicKeyFile"
            :limit=1
            :on-change="handleChange"
            :before-upload="beforeUploadKey"
            :auto-upload="false">
            <el-button  size="small" type="primary">选择公钥文件</el-button>
          </el-upload>
        </el-form-item>
        <el-form-item label="下载私钥" label-width="80px" prop="downloadKey" v-if="userForm.loginType === 'KeyFile' || userForm.loginType === 'both'">
          <el-switch v-model="userForm.downloadKey"></el-switch>
        </el-form-item>
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
    <!-- 公共 SFTP 浏览器组件 -->
    <SftpBrowser
      :username="SftpForm.username"
      :visible="SftpBrowserVisible"
      :upload-headers="uploadHeaders"
      upload-url="/dev-api/sftp/upload"
      @close="closeSftpBrowser"
    />
  </div>
</template>

<script>
import SftpBrowser from '@/components/SftpBrowser/index.vue'
import { rsaEncrypt } from '@/utils/encrypt'

export default {
  name:'SftpUser',
  components: {
    SftpBrowser
  },
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
        id:'', // 用户 id
        name:'', // 账号名
        loginType:'Password', // 登录类型
        noExpire:false, // 永不过期
        emailOrNot:false,  //是否发送邮件
        password:'',   // 密码
        checkPass:'', // 确认密码
        to:'',   //收件人
        cc:'',  //抄送人
        downloadKey:false, // 是否下载私钥
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
        downloadKey: [
          {
            validator: (rule, value, callback) => {
              // 如果选择了下载私钥，必须是 KeyFile 或 both 类型
              if (value && this.userForm.loginType !== 'KeyFile' && this.userForm.loginType !== 'both') {
                callback(new Error('下载私钥仅支持 KeyFile 和 both 登录类型'))
              } else {
                callback()
              }
            },
            trigger: 'change'
          }
        ],
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
      publicKeyFile: [], // 公钥文件列表
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
        publicKey:false,
        password:'',
        checkPass:'',
        to:'',
        cc:'',
        downloadKey:false,
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
          const {id,name,loginType,password,emailOrNot,noExpire,to,cc,downloadKey} = this.userForm
          let message;
          if (loginType === 'KeyFile') {
            message = '密钥修改成功';
          } else if (loginType === 'Password') {
            message = '密码修改成功';
          } else {
            message = '密码和密钥修改成功';
          }
          try {
            const formData = new FormData()
              for (let key of Object.keys(this.userForm)){
                if (key == 'password'){
                  // 对密码加密
                  const rsaPassword = rsaEncrypt(this.userForm[key])
                  formData.append(key, rsaPassword)
                  continue
                }else if ( key == 'checkPass'){
                  continue
                }
                formData.append(key, this.userForm[key])
              }
              // 如果选择了导入公钥文件，则添加到 formData 中
              if (this.userForm.publicKey && this.publicKeyFile.length > 0) {
                this.publicKeyFile.forEach(item => {
                  formData.append('file', item.raw);
                })
              }
            if(this.title == '修改密码'){
              const result = await this.$API.sftpuser.reqUpdateUser(formData)
              if(result.code == 200){
                // 处理私钥下载
                if (downloadKey && (loginType === 'KeyFile' || loginType === 'both')) {
                  try {
                    const downloadResult = await this.$API.sftpuser.reqDownloadKey(name)
                    if (downloadResult && downloadResult.data instanceof Blob) {
                      // 触发浏览器下载
                      const blob = downloadResult.data
                      const url = window.URL.createObjectURL(blob)
                      const link = document.createElement('a')
                      link.href = url
                      link.download = `${name}_sftp_rsa_key`
                      document.body.appendChild(link)
                      link.click()
                      document.body.removeChild(link)
                      window.URL.revokeObjectURL(url)
                          
                      // 更新消息提示
                      message += ', 私钥已下载'
                    }
                  } catch (error) {
                    console.error('私钥下载失败:', error)
                    // 不阻止主流程继续
                  }
                }
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
              const result = await this.$API.sftpuser.reqAddUser(formData)
              if(result.code == 200){
                // 处理私钥下载
                if (downloadKey && (loginType === 'KeyFile' || loginType === 'both')) {
                  try {
                    const downloadResult = await this.$API.sftpuser.reqDownloadKey(name)
                    if (downloadResult && downloadResult.data instanceof Blob) {
                      // 触发浏览器下载
                      const blob = downloadResult.data
                      const url = window.URL.createObjectURL(blob)
                      const link = document.createElement('a')
                      link.href = url
                      link.download = `${name}_sftp_rsa_key`
                      document.body.appendChild(link)
                      link.click()
                      document.body.removeChild(link)
                      window.URL.revokeObjectURL(url)
                                  
                      // 更新消息提示
                      message += ', 私钥已下载'
                    }
                  } catch (error) {
                    console.error('私钥下载失败:', error)
                    // 不阻止主流程继续
                  }
                }
                if(this.userForm.emailOrNot){
                  this.email.to = to
                  this.email.cc = cc
                  this.email.subject = `${name}`
                  await this.$API.sftpuser.reqSendEmail(this.email)
                  this.$message({type:'success',message:`账号:${name} 添加成功${message ? ',' + message : ''},邮件已发送`})
                }else{
                  this.$message({type:'success',message:`账号:${name} 添加成功${message ? ':' + message : ''}`})
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
      // 清除公钥文件列表
      this.publicKeyFile = []
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
        const rsaPassword = rsaEncrypt(password)
        if(password == ''){
          this.$message({type:'warning',message:'请输入SFTP密码'})
          return
        }
        // loading
        this.buttonLoading = true
        try {
          const result = await this.$API.sftpuser.reqSftpLogin({username,password:rsaPassword})
          if(result.code == 200){
            // 存储SFTP-Token
            sessionStorage.setItem("sftp_token", result.data.sftp_token)
            this.uploadHeaders['X-SFTP-Token'] = result.data.sftp_token
            this.$message({type:'success',message:'SFTP登录成功'})
            this.SftpBrowserVisible = true
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
          sessionStorage.setItem("sftp_token", response.data.sftp_token);
          this.uploadHeaders['X-SFTP-Token'] = response.data.sftp_token
        }
        this.SftpBrowserVisible = true
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
    // 关闭 SFTP 浏览器时断开sftp连接
    async closeSftpBrowser() {
      this.SftpBrowserVisible = false
      this.SftpForm.password = ''
      try {
        const result = await this.$API.sftpuser.reqSftpLogout();
        if (result.code === 200) {
          this.$message({ type: 'success', message: result.message });
        } else {
          this.$message({ type: 'error', message: result.message });
        }
      } catch (error) {}
    },
    handleChange(file, fileList) {
      this.publicKeyFile = fileList
    }
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