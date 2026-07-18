<template>
  <div>
    <!-- inline:行内表单,一行可以放置多个表单元素 -->
    <el-form :inline="true" class="demo-form-inline" :model="cForm">
      <el-form-item label="一级分类">
        <el-select  placeholder="请选择" v-model="cForm.category1Id" @change="handler1" :disabled="show">
          <el-option :label="c1.name" :value="c1.id" v-for="c1 in list1" :key="c1.id"></el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="二级分类">
        <el-select  placeholder="请选择"  v-model="cForm.category2Id" @change="handler2" :disabled="show">
          <el-option :label="c2.name" :value="c2.id" v-for="c2 in list2" :key="c2.id"></el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="三级分类">
        <el-select  placeholder="请选择" v-model="cForm.category3Id" @change="handler3" :disabled="show">
          <el-option :label="c3.name" :value="c3.id" v-for="c3 in list3" :key="c3.id"></el-option>
        </el-select>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
export default {
  name: "CategorySelect",
  props:['show'],
  data() {
    return {
      // 一级分类数据
      list1:[],
      // 二级分类数据
      list2:[],
      // 三级分类数据
      list3:[],
      // 收集相应的一二三级分类的id
      cForm:{
        category1Id:'',
        category2Id:'',
        category3Id:'',
      }
    }
  },
  mounted(){
    // 挂载完毕后发送请求获取一级分类
    this.getCategory1List()
  },
  methods:{
    // 获取一级分类
    async getCategory1List(){
      try {
        // 请求
        const result = await this.$API.attr.reqCategory1List()
        if(result.code == 200){
          this.list1 = result.data
        }
      } catch (error) {}
    },
    // 当一级分类的select变化时,获取相应的二级分类数据
    async handler1(){
      // 清除数据
      this.list3 = []
      this.cForm.category2Id = ''
      this.cForm.category3Id = ''
      // 结构出一级分类id
      const {category1Id} = this.cForm
      this.$emit('getCategoryId',{categoryId:category1Id,level:1})
      try {
        // 通过一级分类id获取二级分类
        const result = await this.$API.attr.reqCategory2List(category1Id)
        if(result.code == 200){
          this.list2 = result.data
        }
      } catch (error) {}
    },
    // 当二级分类select变变化时,获取相应的三级分类数据
    async handler2(){
      // 清除数据
      this.cForm.category3Id = ''
      // 解构出二级分类id
      const {category2Id} = this.cForm
      this.$emit('getCategoryId',{categoryId:category2Id,level:2})
      try {
        // 通过二级分类id,获取三级分类数据
        const result = await this.$API.attr.reqCategory3List(category2Id)
        if(result.code == 200){
          this.list3 = result.data
        }
      } catch (error) {}
    },
    handler3(){
      const {category3Id} = this.cForm
      this.$emit('getCategoryId',{categoryId:category3Id,level:3})
    }
  }
};
</script>

<style>
</style>