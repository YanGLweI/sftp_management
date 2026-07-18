<template>
  <div>
    <!-- 上部分看片 -->
    <el-card style="margin: 20px 0px;">
      <!-- 三级联动:全局组件可以直接使用 -->
      <CategorySelect @getCategoryId="getCategoryId" :show="!show"></CategorySelect>
    </el-card>
    <!-- 下部分卡片 -->
    <el-card>
      <!-- 底部这里有三部分结构进行切换显示 -->
      <div v-show="scene==0">
        <!-- 展示SPU列表的结构 -->
        <el-button type="primary" icon="el-icon-plus" :disabled="!category3Id" @click="addSpu">添加SPU</el-button>
        <el-table  style="width: 100%;margin: 20px 0;" border :data="records">
            <el-table-column type="index" label="序号" width="80" align="center">
            </el-table-column>
            <el-table-column prop="spuName" label="SPU名称" width="width">
            </el-table-column>
            <el-table-column prop="description" label="SPU描述" width="width">
            </el-table-column>
            <el-table-column label="操作" width="width">
              <template slot-scope="{row}">
                <!-- 用自己封装的hintButton -->
                <hint-button type="success" icon="el-icon-plus" size="mini" title="添加sku" @click="addSku(row)"></hint-button>
                <hint-button type="warning" icon="el-icon-edit" size="mini" title="修改sku" @click="updateSpu(row)"></hint-button>
                <hint-button type="info" icon="el-icon-info" size="mini" title="查看当前spu全部sku" @click="handler(row)"></hint-button>
                <el-popconfirm :title="`确定删除 ${row.spuName} 吗?`" @onConfirm="deleteSpu(row)">
                  <hint-button type="danger" icon="el-icon-delete" size="mini" title="删除spu" slot="reference"></hint-button>
                </el-popconfirm>
              </template>
            </el-table-column>
        </el-table>
        <!-- 分页器 -->
        <el-pagination 
          :current-page="page" 
          :page-sizes="[3, 5, 10]" 
          :page-size="limit" 
          layout="prev, pager, next, jumper, ->, sizes, total" 
          :total="total"
          @current-change="getSpuList"
          @size-change="handleSizeChange" 
          style="text-align: center;"
          >
        </el-pagination>
      </div>
      <SpuForm v-show="scene==1" @changeScene="changeScene" ref="spu"></SpuForm>
      <SkuForm v-show="scene==2" ref="sku" @changeScenes="changeScenes"></SkuForm> 
    </el-card>
    <el-dialog :title="`${spu.spuName}的sku列表`" :visible.sync="dialogTableVisible" @close="close">
      <!-- table展示sku列表 -->
      <el-table :data="skuList" border v-loading="loading">
          <el-table-column prop="skuName" label="名称" width="width">
          </el-table-column>
          <el-table-column prop="price" label="价格" width="width">
          </el-table-column>
          <el-table-column prop="weight" label="重量" width="width">
          </el-table-column>
          <el-table-column prop="prop" label="默认图片" width="width">
            <template slot-scope="{row}">
              <img :src="row.skuDefaultImg" alt="" style="height: 50px;width: 50px;">
            </template>
          </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script>
import SpuForm from './SpuForm'
import SkuForm from './SkuForm'
export default {
  name:'Spu',
  components:{SpuForm,SkuForm},
  data() {
    return {
      // 存储子组件获取的分类id
      category1Id:'',
      category2Id:'',
      category3Id:'',
      // 控制三级分类联动的可操作性
      show:true,
      page:1, //当前第几页
      limit:3, //每页展示几条数据
      records:[],  //spu列表数据
      total:0,  //分页器展示的数据总数
      scene:0,  //0:展示SPU列表 1:展示添加SPU|修改SPU 2:展示添加SKU
      // 控制对话框的显示与隐藏
      dialogTableVisible:false,
      spu:{},
      skuList:[],
      loading:true
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
        // 获取Spu列表数据
        this.getSpuList()
      }
    },
    // 获取SPU列表的回调
    async getSpuList(pages=1){
      this.page = pages
      const {page,limit,category3Id} = this
      try {
        // 携带三个参数:page\limit\category3Id
        const result = await this.$API.spu.reqSpuList(page,limit,category3Id)
        if(result.code == 200){
          this.records = result.data.records
          this.total = result.data.total
        }
      } catch (error) {}
    },
    // 改变每页显示条数
    handleSizeChange(limit){
      this.limit = limit
      this.getSpuList()
    },
    // 添加Spu按钮的回调
    addSpu(){
      // 切换为 添加Spu结构
      this.scene = 1
      // 三级联动改为 fales不能操作
      this.show = false
      // 通知子组件发请求
      this.$refs.spu.addSpuData(this.category3Id)
    },
    // 修改某一个Spu的回调
    updateSpu(row){
      // 切换为 添加Spu结构
      this.scene = 1
      // 三级联动改为 fales不能操作
      this.show = false
      // 获取子组件SpuForm
      this.$refs.spu.initSpuData(row)
    },
    // 自定义事件的回调(SpuForm)
    changeScene({scene,flag}){
      // flag:区分是修改时的保存,还是添加的保存
      // 切换到spu列表结构
      this.scene = scene
      this.show = true
      if(flag == '修改' || flag == '取消'){
        this.getSpuList(this.page)
      }else{
        this.getSpuList()
      }
    },
    // 删除SPU
    async deleteSpu(row){
      try {
        const result = await this.$API.spu.reqDeleteSpu(row.id)
        if (result.code == 200){
          this.$message({type:'success',message:'删除成功'})
          this.getSpuList(this.records.length>1?this.page:this.page-1)
        }
      } catch (error) {}
    },
    // 添加Sku按钮的回调
    addSku(row){
      // 切换到场景2
      this.scene = 2
      // 父组件调用子组件方法,让子组件发3个请求
      this.$refs.sku.getData(this.category1Id,this.category2Id,row)
    },
    // 自定义事件回调,skuform通知父组件切换场景
    changeScenes(scene){
      // 切换场景为0
      this.scene = scene
    },
    // 查看sku按钮的回调
    async handler(spu){
      // 点击后对话框可见
      this.dialogTableVisible = true
      this.spu = spu
      // 获取sku列表数据展示
      try {
        const result = await this.$API.spu.reqSkuList(spu.id)
        this.skuList = result.data
        this.loading = false
      } catch (error) {}
    },
    // 关闭对话框的回调
    close(){
      // loading属性再次变为true
      this.loading = true
      // 清楚sku列表数据
      this.skuList = []
      // 关闭对话框
    }
  }
}
</script>

<style>

</style>