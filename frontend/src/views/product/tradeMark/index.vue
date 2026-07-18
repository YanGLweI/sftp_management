<template>
  <div>
    <!-- 按钮 -->
    <el-button type="primary" icon="el-icon-plus" style="margin-bottom: 20px;" @click="showDialog">添加</el-button>
    <!-- 
      表格组件
      data:表格组件将来要展示的数据----数组类型
      border:纵向边框,默认没有边框

      column(列)属性
        label:显示的标题
        width:对应列的宽度
        align:对齐方式left/center/right	默认是left
        type:index 可以展示序号
        prop: 对应列内容的字段名

        注意:elemenUI当中的table组件,展示的数据是以list的长度来决定的
    -->
    <el-table style="width: 100%" border :data="list">
        <el-table-column type="index" label="序号" width="80px" align="center">
        </el-table-column>
        <el-table-column prop="tmName" label="品牌名称" width="width">
        </el-table-column>
        <el-table-column label="品牌LOGO" width="width">
          <template slot-scope="{row}">
            <!-- row:是插槽回传的数据 -->
            <img :src="row.logoUrl" alt="" style="width: 50px;height: 50px;">
          </template>
        </el-table-column>
        <el-table-column label="操作" width="width">
          <template slot-scope="{row}">
            <el-button type="warning" icon="el-icon-edit" size="mini" @click="updateTradeMark(row)">修改</el-button>
            <el-button type="danger" icon="el-icon-delete" size="mini" @click="deleteTradeMark(row)">删除</el-button>
          </template>
        </el-table-column>
    </el-table>
    <!-- 分页器
      当前页第几页、数据总数、每页展示条数、连续页码数
      改变每页展示条数的回调
      @size-change="handleSizeChange" 
      点击页码的回调
      @current-change="handleCurrentChange" 

      layout:布局改变每个功能的位置，用->可以将功能靠右对齐
      current-page：当前第几页
      total：数据总数
      page-size：每页展示条数
      page-sizes：设置每页展示条数
      pager-count:页码按钮的数量 默认7
    -->
    <el-pagination 
      style="margin-top: 20px;text-align: center;"
      :current-page="page" 
      :total="total" 
      :page-size="limit" 
      :pager-count="7"
      :page-sizes="[3, 5, 10]" 
      layout="prev, pager, next, jumper,->, sizes, total" 
      @current-change="getPageList"
      @size-change="handleSizeChange"
      >
    </el-pagination>
    <!-- 对话框
      visible.sync : 控制对话框显示和隐藏
    -->
    <el-dialog :title="tmForm.id?'修改品牌':'添加品牌'" :visible.sync="dialogFormVisible" width="500px">
      <!-- from表单 :model这个属性的作用:把表单的数据收集到哪个对象身上,将来表单验证也需要这个属性-->
      <el-form style="width: 80%;" :model="tmForm" :rules="rules" ref="ruleForm">
        <el-form-item label="品牌名称" label-width="100px" prop="tmName">
          <el-input autocomplete="off" v-model="tmForm.tmName"></el-input>
        </el-form-item>
        <el-form-item label="品牌LOGO" label-width="100px" prop="logoUrl">
          <!-- 图片上传
            收集数据:不能使用v-model,因为不是表单元素
            action:设置图片上传的地址
            :on-success:可以检查到图片上传成功,当上传成功时执行一次
            :before-upload:上传之前,会执行一次
          -->
          <el-upload
            class="avatar-uploader"
            action="/dev-api/admin/product/fileUpload"
            :show-file-list="false"
            :on-success="handleAvatarSuccess"
            :before-upload="beforeAvatarUpload"
          >
            <img v-if="tmForm.logoUrl" :src="tmForm.logoUrl" class="avatar">
            <i v-else class="el-icon-plus avatar-uploader-icon"></i>
            <div slot="tip" class="el-upload__tip">只能上传jpg/png文件，且不超过2MB</div>
          </el-upload>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false">取 消</el-button>
        <el-button type="primary" @click="addOrUpdateTradeMark">确 定</el-button>
      </div>
    </el-dialog>
  </div>
  
  
</template>

<script>
export default {
  name:'tradeMark',
  data() {
    // 自定义校验规则
    var validateTmName = (rule, value, callback) => {
      if(value.length < 2 || value.length > 10){
        callback(new Error('品牌名称2-10位'))
      }else(
        callback()
      )
    };
    return {
      // 代表分页器第几页
      page:1,
      // 当前页数展示数据条数
      limit:3,
      //数据总数
      total:0,
      // 列表展示的数据
      list:[],
      // 对话框显示与隐藏
      dialogFormVisible: false,
      // 对话框title
      title:'',
      // 收集品牌信息
      tmForm:{
        // 属性名根据api要求来写
        // 收集品牌名称
        tmName:'',
        // 收集品牌图片信息
        logoUrl:'',
      },
      // 表单验证规则
      rules:{
        // 品牌名称的验证规则
        tmName: [
          // required:必须要校验的,标签前会带*号
          // message:提示信息
          // trigger:用户行为触发验证的设置 blur失焦 change变化
          { required: true, message: '请输入品牌名称', trigger: 'blur' },
          // 自定义校验规则
          { validator: validateTmName, trigger: 'change' }
        ],
        // 品牌logo的验证规则
        logoUrl: [
          { required: true, message: '请选择品牌图片' }
        ],
      }
    } 
  },
  mounted(){
    // 获取列表数据的方法
    this.getPageList()
  },
  methods:{
    async getPageList(pager = 1){
      this.page = pager
      // 解构参数
      const {page,limit} = this
      const result = await this.$API.trademark.reqTradeMarkList(page,limit)
      if (result.code == 200){
        // 分别是展示数据的总条数和列表
        this.total = result.data.total
        this.list = result.data.records
      }
    },
    /* // 改变当前页,默认收到点击的页码数
    handleCurrentChange(pager){
      this.page = pager
      this.getPageList()
    }, */
    // 改变每页展示的条数,默认收到选择的条数
    handleSizeChange(limit){
      this.limit = limit
      this.getPageList()
    },
    // 显示对话框
    showDialog(){
      this.dialogFormVisible = true
      this.tmForm = {tmName:'',logoUrl:''}
    },
    // 修改某一个品牌
    updateTradeMark(row){
      // row:当前用户选中的品牌信息
      this.dialogFormVisible = true
      // 将品牌信息赋值给tmForm
      // const {id,tmName,logoUrl} = row 
      // this.tmForm = row
      this.tmForm = {...row}
    },
    // 图片上传成功
    handleAvatarSuccess(res, file) {
      // res:是上传成功后返回给前端的数据
      // file:上传成功后,服务器返回的数据
      // 收集品牌图片数据,需要带给服务器
      this.tmForm.logoUrl = res.data
    },
    // 图片上传之前
    beforeAvatarUpload(file) {
      // 判断文件格式
      const isJPG = file.type === 'image/jpeg' || file.type ==='image/png';
      // 判断文件大小
      const isLt2M = file.size / 1024 / 1024 < 2;

      if (!isJPG) {
        this.$message.error('上传头像图片只能是 JPG/PNG 格式!');
      }
      if (!isLt2M) {
        this.$message.error('上传头像图片大小不能超过 2MB!');
      }
      return isJPG && isLt2M;
    },
    // 添加或修改完成的提交按钮
    addOrUpdateTradeMark(){
      // 表单验证全部通过后才去执行业务逻辑
      this.$refs.ruleForm.validate(async (valid) =>{
        // 当表单全部验证通过时valid为true
        if(valid){
          // 关闭对话框
          this.dialogFormVisible = false
          // 发请求(添加或者修改)
          try {
            const result = await this.$API.trademark.reqAddOrUpdateTradeMark(this.tmForm)
            if (result.code == 200){
              // 弹出信息:添加品牌成功或修改品牌成功
              this.$message({
                message: this.tmForm.id ? '修改品牌成功' : '添加品牌成功',
                type: 'success'
              })
              // 重新获取品牌列表进行展示
              // 修改完成留在当前页
              this.getPageList(this.page)
            }
          } catch (error) {}
        }else{
          console.log('error submit!!');
          return false;
        }
      })
    },
    // 删除品牌
    deleteTradeMark(row){
      this.$confirm(`你确定删除 ${row.tmName} 品牌吗?`, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async () => {
        // then:点击确定按钮的时候触发
        // 发送删除请求
        try {
          await this.$API.trademark.reqDeleteTradeMark(row.id)
          this.$message({
            type: 'success',
            message: '删除成功!'
          });
          // 重新请求品牌列表
          // 如果当前页只有一条数据,删除后应该跳到前一页
          this.getPageList(this.list.length>1 ? this.page: this.page-1)
        } catch (error) {}
      }).catch(() => {
        // catch:点击取消按钮时触发
        this.$message({
          type: 'info',
          message: '已取消删除'
        });          
      });
    }
  }
}
</script>

<style>
  .avatar-uploader .el-upload {
    border: 1px dashed #d9d9d9;
    border-radius: 6px;
    cursor: pointer;
    position: relative;
    overflow: hidden;
  }
  .avatar-uploader .el-upload:hover {
    border-color: #409EFF;
  }
  .avatar-uploader-icon {
    font-size: 28px;
    color: #8c939d;
    width: 178px;
    height: 178px;
    line-height: 178px;
    text-align: center;
  }
  .avatar {
    width: 178px;
    height: 178px;
    display: block;
  }
</style>