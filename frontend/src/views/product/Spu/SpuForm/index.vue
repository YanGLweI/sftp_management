<template>
  <div>
    <el-form ref="form" label-width="80px" :model="spu">
      <el-form-item label="SPU名称">
        <el-input placeholder="SPU名称" v-model="spu.spuName"></el-input>
      </el-form-item>
      <el-form-item label="品牌">
        <el-select placeholder="请选择品牌" v-model="spu.tmId">
          <el-option :label="tm.tmName" :value="tm.id" v-for="tm in tradeMarkList" :key="tm.id"></el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="SPU描述">
        <!-- type="textarea" 可以设置input为文本区, rows="4" 可以设置文本区大小 -->
        <el-input placeholder="SPU描述" type="textarea" rows="4" v-model="spu.description"></el-input>
      </el-form-item>
      <el-form-item label="SPU图片">
        <!-- 照片墙  list-type:文件列表的类型 on-preview:预览是触发 on-remove:移除时触发 file-list:上传的文件列表,数组的元素需要有name和url字段
        on-preview:图片预览 on-remove:删除图片时触发 -->
        <el-upload
          action="/dev-api/admin/product/fileUpload"
          list-type="picture-card"
          :on-preview="handlePictureCardPreview"
          :on-remove="handleRemove"
          :on-success="handleSuccess"
          :file-list="spuImageList"
        >
          <i class="el-icon-plus"></i>
        </el-upload>
        <el-dialog :visible.sync="dialogVisible">
          <img width="100%" :src="dialogImageUrl" alt="" />
        </el-dialog>
      </el-form-item>
      <el-form-item label="销售属性">
        <el-select v-model="attrIdAndAttrName" :placeholder="`还有${unSelectSaleAttr.length}个未选择`">
          <el-option :label="unselect.name" :value="`${unselect.id}:${unselect.name}`" v-for="unselect in unSelectSaleAttr" :key="unselect.id"></el-option>
        </el-select>
        <el-button type="primary" icon="el-icon-plus" style="margin-left: 5px" :disabled="!attrIdAndAttrName" @click="addSaleAttr">添加销售属性</el-button>
        <!-- 展示当前SPU属于自己的销售属性 -->
        <el-table style="width: 100%; margin-top: 5px" border :data="spu.spuSaleAttrList">
          <el-table-column type="index" label="序号" width="80" align="center">
          </el-table-column>
          <el-table-column prop="saleAttrName" label="属性名" width="width">
          </el-table-column>
          <el-table-column prop="prop" label="属性值名称列表" width="width">
            <template slot-scope="{row}">
              <!-- @close="handleClose(tag)" -->
              <!-- 展示已有的属性值列表 -->
              <el-tag :key="tag.id" v-for="(tag,index) in row.spuSaleAttrValueList" closable :disable-transitions="false" @close="handleClose(row,index)">{{tag.saleAttrValueName}}</el-tag>
              <!-- input和Button切换显示 -->
              <!-- @keyup.enter.native="handleInputConfirm"  @blur="handleInputConfirm" -->
              <el-input class="input-new-tag" v-if="row.inputVisible" v-model="row.inputValue" ref="saveTagInput" size="small" @keyup.enter.native="$event.target.blur"  @blur="handleInputConfirm(row)"></el-input>
              <!-- @click="showInput" -->
              <el-button v-else class="button-new-tag" size="small" @click="addSaleAttrValue(row)">添加</el-button>
            </template>
          </el-table-column>
          <el-table-column prop="prop" label="操作" width="width">
            <template slot-scope="{$index}">
              <el-button type="danger" icon="el-icon-delete" size="mini" @click="spu.spuSaleAttrList.splice($index,1)"></el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="addOrUpdateSpu">保存</el-button>
        <el-button @click="cancel">取消</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
export default {
  name: "SpuForm",
  data() {
    return {
      dialogImageUrl: "",
      dialogVisible: false,
      // 存储SPU信息
      spu: {
        category3Id: 0,   //三级分类id
        description: "",  //描述
        tmId: '',          //品牌id
        spuName: "",      //Spu名称
        // SPU图片信息
        spuImageList: [
          // {
          //   id: 0,
          //   imgName: "",
          //   imgUrl: "",
          //   spuId: 0,
          // },
        ],
        // 平台属性与属性值的收集
        spuSaleAttrList: [
          // {
          //   baseSaleAttrId: 0,
          //   id: 0,
          //   saleAttrName: "string",
          //   spuId: 0,
          //   spuSaleAttrValueList: [
          //     {
          //       baseSaleAttrId: 0,
          //       id: 0,
          //       isChecked: "string",
          //       saleAttrName: "string",
          //       saleAttrValueName: "string",
          //       spuId: 0,
          //     },
          //   ],
          // },
        ],
      },
      // 存储品牌列表
      tradeMarkList: [],
      // 存储SPU图片列表
      spuImageList: [],
      // 存储平台销售属性
      saleAttrList: [],
      attrIdAndAttrName:'' ,//收集未选择的销售属性id和属性名
    };
  },
  methods: {
    // 上传成功的回调
    handleSuccess(response, file, fileList){
      // 收集图片信息
      this.spuImageList = fileList
    },
    // 照片墙删除某一个图片时触发
    handleRemove(file, fileList) {
      // file:删除的那张图片,fileList:剩余的图片
      // console.log(file, fileList);
      // 收集照片墙图片数据
      // 已有的数据含有name和url,服务器不需要,将来要删除
      this.spuImageList = fileList
    },
    // 照片墙的预览回调
    handlePictureCardPreview(file) {
      // 将地址赋值给这个属性,对话框要用
      this.dialogImageUrl = file.url;
      // 显示对话框
      this.dialogVisible = true;
    },
    // 初始化SpuForm数据
    async initSpuData(spu) {
      // 请求数据
      try {
        // 获取Spu信息的数据
        const spuResult = await this.$API.spu.reqSpu(spu.id);
        if (spuResult.code == 200) {
          this.spu = spuResult.data;
        }
      } catch (error) {}
      try {
        // 获取品牌信息
        const tradeMarkResult = await this.$API.spu.reqTradeMarkList();
        if (tradeMarkResult.code == 200) {
          this.tradeMarkList = tradeMarkResult.data;
        }
      } catch (error) {}
      try {
        // 获取spu图片
        const SpuImageResult = await this.$API.spu.reqSpuImageList(spu.id);
        if (SpuImageResult.code == 200) {
          let listArr = SpuImageResult.data;
          // 由于照片墙显示图片,数组的元素需要有name和url字段
          listArr.forEach(item=>{
            item.name = item.imgName
            item.url = item.imgUrl
          })
          // 把整理好的数据赋值给spuImageList
          this.spuImageList = listArr
        }
      } catch (error) {}
      try {
        // 获取平台全部的销售属性
        const saleResult = await this.$API.spu.reqBaseSaleAttrList();
        if (saleResult.code == 200) {
          this.saleAttrList = saleResult.data;
        }
      } catch (error) {}
    },
    // 添加新的销售属性按钮回调
    addSaleAttr(){
      // 已经收集到了attrIdAndAttrName中,但是格式不对
      // 把收集到的数据分割
      const [baseSaleAttrId,saleAttrName] = this.attrIdAndAttrName.split(':')
      // 向spu对象的spuSaleAttrList中push新的销售属性数据
      let newSaleAttr = {baseSaleAttrId,saleAttrName,spuSaleAttrValueList:[]}
      this.spu.spuSaleAttrList.push(newSaleAttr)
      // 清空数据
      this.attrIdAndAttrName = ''
    },
    // 添加属性值按钮的回调
    addSaleAttrValue(row){
      // 在row上添加inputVisible,控制input和Button的切换
      this.$set(row,'inputVisible',true)
      this.$set(row,'inputValue','')
      this.$nextTick(_=>{
        this.$refs.saveTagInput.focus()
      })
    },
    // el-input的回车和失焦事件
    handleInputConfirm(row){
      // 结构出需要的数据
      const {baseSaleAttrId,inputValue} = row
      // 新增的属性值,不能为空或者重复
      if(inputValue.trim() == ''){
        this.$message('属性值不能为空')
        // 切换为Button
        row.inputVisible = false
        // 清空数据
        row.inputValue = ''
        return
      }
      // 属性值不能重复
      let result = row.spuSaleAttrValueList.every((item)=>item.saleAttrValueName != inputValue)
      if(!result) {
        this.$message('属性值不能重复')
        // 切换为Button
        row.inputVisible = false
        // 清空数据
        row.inputValue = ''
        return
      }
      let newSaleAttrValue = {baseSaleAttrId,saleAttrValueName:inputValue}
      // 新增
      row.spuSaleAttrValueList.push(newSaleAttrValue);
      // 切换为Button
      row.inputVisible = false
      // 清空数据
      row.inputValue = ''
    },
    // 删除属性值
    handleClose(row,index){
      row.spuSaleAttrValueList.splice(index,1)
    },
    // 保存按钮的回调
    async addOrUpdateSpu(){
      // 整理参数:需要整理照片墙的数据
      // 携带的参数:对于图片,imgName imgUrl
      // 用map返回新数组
      this.spu.spuImageList = this.spuImageList.map(item=>{
        return {
          imgName:item.name,
          imgUrl:(item.response&&item.response.data)||item.url
        }
      })
      // 发请求
      try {
        const result = await this.$API.spu.reqAddOrUpdateSpu(this.spu)
        if(result.code==200){
          this.$message({type:'success',message:'保存成功'})
          // 返回spu列表
          this.$emit('changeScene', {scene:0,flag:this.spu.id?'修改':'添加'})
        }
      } catch (error) {}
      // 清除数据
      Object.assign(this._data,this.$options.data())
    },
    // 点击添加spu按钮发请求获取品牌和全部销售属性
    async addSpuData(category3Id){
      this.spu.category3Id = category3Id
      try {
        // 获取品牌信息
        const tradeMarkResult = await this.$API.spu.reqTradeMarkList();
        if (tradeMarkResult.code == 200) {
          this.tradeMarkList = tradeMarkResult.data;
        }
      } catch (error) {}
      try {
        // 获取平台全部的销售属性
        const saleResult = await this.$API.spu.reqBaseSaleAttrList();
        if (saleResult.code == 200) {
          this.saleAttrList = saleResult.data;
        }
      } catch (error) {}
    },
    // 取消按钮
    cancel(){
      // 切换场景
      this.$emit('changeScene', {scene:0,flag:'取消'})
      // 清除数据
      // Object.assign:es6中新增的方法可以合并对象 
      // 组件实例的this._data,可以操作data中的响应式数据
      // this.$options:可以获取配置对象,配置对象的data函数执行,返回初始的响应式数据
      Object.assign(this._data,this.$options.data())
    }
  },
  computed:{
    // 计算出还未选择的销售属性
    unSelectSaleAttr(){
      // 整个平台一共三个销售属性:颜色\尺寸\版本 ---- saleAttrList
      // spu中有属于自己销售属性 ----- spu.spuSaleAttrList
      let newList = this.saleAttrList.filter((item)=>{
        // every方法会返回一个布尔值
        return this.spu.spuSaleAttrList.every((item1)=>{
          return item1.saleAttrName != item.name
        })
      })
      return newList
    },
  },
};
</script>

<style scoped>
  .el-tag + .el-tag {
    margin-left: 10px;
  }
  .button-new-tag {
    margin-left: 10px;
    height: 32px;
    line-height: 30px;
    padding-top: 0;
    padding-bottom: 0;
  }
  .input-new-tag {
    width: 90px;
    margin-left: 10px;
    vertical-align: bottom;
  }
</style>