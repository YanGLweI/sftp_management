<template>
  <div>
    <!-- 上部分卡片 -->
    <el-card style="margin: 20px 0px;">
      <CategorySelect @getCategoryId="getCategoryId" :show="!isShowTable"></CategorySelect>
    </el-card>
    <!-- 下部分卡片 -->
    <el-card>
      <div v-show="isShowTable">
        <!-- :disabled:有3级id的时候才可以点击 -->
        <el-button type="primary" icon="el-icon-plus" :disabled="!category3Id" @click="addAttr">添加属性</el-button>
        <!-- 表格:展示平台属性 -->
        <el-table :data="attrList" style="width: 100%;margin-top: 20px;" border>
            <el-table-column type="index" label="序号" width="80px" align="center">
            </el-table-column>
            <el-table-column prop="attrName" label="属性名称" width="150px">
            </el-table-column>
            <el-table-column prop="prop" label="属性值列表" width="width">
              <template slot-scope="{row}">
                <el-tag type="success" v-for="attrValue in row.attrValueList" :key="attrValue.id" style="margin: 5px 5px;">{{ attrValue.valueName }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="prop" label="操作" width="150px">
              <template slot-scope="{row}">
                <el-button type="warning" icon="el-icon-edit" size="mini" @click="updateAttr(row)" style="margin-right: 10px;"></el-button>
                <!-- 气泡确认框 -->
                <el-popconfirm 
                  :title="`确定删除 ${row.attrName} 吗?`"
                  @onConfirm="deleteAttr(row)"
                  >
                  <el-button type="danger" icon="el-icon-delete" size="mini" slot="reference"></el-button>
                </el-popconfirm>
              </template>
            </el-table-column>
        </el-table>
      </div>
      <!-- 添加属性|修改属性的结构 -->
      <div v-show="!isShowTable">
        <el-form :inline="true" ref="form" :model="attrInfo" label-width="80px">
            <el-form-item label="属性名">
              <el-input placeholder="请输入属性名" v-model="attrInfo.attrName"></el-input>
            </el-form-item>
        </el-form>
        <el-button type="primary" icon="el-icon-plus" size="mini" @click="addAttrValue" :disabled="!attrInfo.attrName">添加属性值</el-button>
        <el-button size="mini" @click="isShowTable = true">取消</el-button>
        <el-table style="width: 100%;margin: 20px 0px;" border :data="attrInfo.attrValueList">
            <el-table-column type="index" label="序号" width="80" align="center">
            </el-table-column>
            <el-table-column header-align="center" label="属性值名称" width="width">
              <template slot-scope="{row,$index}">
                <!-- input和span切换显示 -->
                <el-input v-model="row.valueName" placeholder="请输入属性值名称" size="mini" v-if="row.flag" @blur="toLook(row)" @keyup.native.enter="toLook(row)" :ref="$index"></el-input>
                <span v-else @click="toEdit(row,$index)" style="display: block;">{{ row.valueName }}</span>
              </template>
            </el-table-column>
            <el-table-column header-align="center" label="操作" width="width">
              <template slot-scope="{row,$index}">
                <!-- 气泡确认框 -->
                <el-popconfirm 
                  :title="`确定删除 ${row.valueName} 吗?`"
                  @onConfirm="deleteAttrValue($index)" 
                >
                  <el-button type="danger" icon="el-icon-delete" size="mini" slot="reference"></el-button>
                </el-popconfirm>
              </template>
            </el-table-column>
        </el-table>
        <el-button type="primary" @click="addOrUpdateAttr" :disabled="attrInfo.attrValueList.length < 1">保存</el-button>
        <el-button @click="isShowTable = true">取消</el-button>
      </div>
    </el-card>
  </div>
</template>

<script>
// 按需引入loadash中的深拷贝
import cloneDeep from 'lodash/cloneDeep'
export default {
  name:'Attr',
  data() {
    return {
      // 存储子组件获取的分类id
      category1Id:'',
      category2Id:'',
      category3Id:'',
      // 接收平台属性的数据
      attrList:[],
      // 控制table表格的显示和隐藏
      isShowTable:true,
      // 收集新增属性|修改属性
      attrInfo:{
        attrName: "",       //属性名
        attrValueList: [    //属性值数组
          // { 
          //   attrId: 0,      //属性ID
          //   valueName: ""   //属性值
          // }
        ],
        categoryId: 0,      //category3Id
        categoryLevel: 3,
      },
    }
  },
  methods:{
    // 自定义事件回调
    getCategoryId({categoryId,level}){
      // 通过level区分三级分类id,存储到父组件
      if(level == 1){
        this.category1Id = categoryId
        this.category2Id = ''
        this.category3Id = ''
      }else if(level == 2){
        this.category2Id = categoryId
        this.category3Id = ''
      }else{
        // 当三级分类id有了的时候,可以发请求了
        this.category3Id = categoryId
        // 发请求获取平台属性
        this.getAttrList()
      }
    },
    // 获取平台属性的数据
    async getAttrList(){
      // 获取分类id
      const {category1Id,category2Id,category3Id} = this
      try {
        const result = await this.$API.attr.reqAttrList(category1Id,category2Id,category3Id)
        if(result.code == 200){
          this.attrList = result.data
        }
      } catch (error) {}
    },
    // 添加属性值的回调
    addAttrValue(){
      // 向属性值数组添加一个元素
      // attrId:是属性的id,由服务器生成,所以带给服务器时是undefined
      this.attrInfo.attrValueList.push({
        attrId:this.attrInfo.id,  //对于已有的属性,是有id的,应该直接赋值;新增的属性没有id,就是undefined
        valueName:'',
        flag:true  //给每一个属性添加flag,切换查看和编辑模式,只控制自己的切换
      })
      // 新增属性值时自动聚焦input
      this.$nextTick(()=>{
        this.$refs[this.attrInfo.attrValueList.length-1].focus()
      })
    },
    // 添加属性按钮的回调
    addAttr(){
      // 切换table的显示与隐藏
      this.isShowTable = false
      // 每次点击添加属性,清空上一次输入的属性和属性值数据
      this.attrInfo = {
        attrName: "",       //属性名
        attrValueList: [],    //属性值数组,
        // 点击添加属性按钮的时候,category3Id已经存在,可以直接收集
        categoryId: this.category3Id,      //category3Id
        categoryLevel: 3,
      }
    },
    // 修改某一个属性
    updateAttr(row){
      // 切换table的显示和隐藏
      this.isShowTable = false
      // 将选中的属性赋值给attrInfo
      // 由于双向绑定,在没提前的情况也也会修改页面数据
      // 由于数据结构中存在对象套数组,数组里套对象,因此需要深拷贝来解决这类问题
      this.attrInfo = cloneDeep(row)
      // 在修改某一个元素时,将相应的属性值元素上加上flag
      this.attrInfo.attrValueList.forEach(item =>{
        // 这样添加的属性不是响应式属性,变化时不会触发视图更新
        // item.flag = false
        // $set:第一个参数:对象 第二个参数:添加新的响应式属性 第三个参数:新的属性的值
        this.$set(item,'flag',false)
      })
    },
    // 失去焦点的回调--切换查看模式
    toLook(row){
      // 属性值不能为空
      if(row.valueName.trim()==''){
        this.$message('请输入一个属性值')
        return
      }
      // 新增的属性值不能和已有的重复
      let isRepat = this.attrInfo.attrValueList.some(item => {
        // 需要将当前row在判断的时候去除
        if(row != item){
          return row.valueName == item.valueName
        }
      })
      // 判断是否重复
      if(isRepat){
        this.$message('属性值重复')
        return
      }
      row.flag = false
    },
    // 点击span标签变为编辑模式的回调
    toEdit(row,index){
      row.flag=true
      // 获取input节点实现自动聚焦
      // 点击span时切换为input,页面重绘耗时间,这个时候无法获取到input节点
      // $nextTick,节点重新渲染完成后执行一次
      this.$nextTick(()=>{
        // 获取input节点元素实现聚焦
        this.$refs[index].focus()
      })
    },
    // 删除属性值气泡确认框确认按钮的回调
    deleteAttrValue(index){
      // 当前删除不需要发送请求
      this.attrInfo.attrValueList.splice(index,1)
    },
    // 保存按钮进行添加或者修改属性的操作
    async addOrUpdateAttr(){
      // 整理参数1,如果用户添加了空属性值,不应该提交
      // 提交给服务器的数据不应该有flag属性
      this.attrInfo.attrValueList = this.attrInfo.attrValueList.filter(item =>{
        // 过滤掉不是空的属性值
        if(item.valueName != ''){
          // 删除掉flag属性
          delete item.flag
          return true
        }
      })
      try {
        // 发请求
        await this.$API.attr.reqAddOrUpdateAttr(this.attrInfo)
        this.$message({
          message:'保存成功',
          type:'success'
        })
        // 重新请求平台属性的数据
        this.getAttrList()
        // 切换到属性展示
        this.isShowTable = true
      } catch (error) {}
    },
    // 删除属性气泡确认按钮的回调
    async deleteAttr(row){
      try {
        const result = await this.$API.attr.reqDeleteAttr(row.id)
        this.$message({message:'删除成功',type:'success'})
        this.getAttrList()
      } catch (error) {}
    }
  }
}
</script>

<style>

</style>