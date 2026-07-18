<template>
  <div>
    <el-form ref="form" label-width="80px">
      <el-form-item label="SPU名称"> {{ spu.spuName }} </el-form-item>
      <el-form-item label="SKU名称">
        <el-input placeholder="Sku名称" v-model="skuInfo.skuName"></el-input>
      </el-form-item>
      <el-form-item label="价格(元)">
        <el-input placeholder="价格(元)" type="number" v-model="skuInfo.price"></el-input>
      </el-form-item>
      <el-form-item label="重量(千克)">
        <el-input placeholder="重量(千克)" v-model="skuInfo.weight"></el-input>
      </el-form-item>
      <el-form-item label="规格描述">
        <el-input placeholder="规格描述" type="textarea" rows="3" v-model="skuInfo.skuDesc"></el-input>
      </el-form-item>
      <el-form-item label="平台属性">
        <el-form :inline="true" ref="form1" label-width="80px">
          <el-form-item :label="attr.attrName" v-for="(attr) in attrInfoList" :key="attr.id">
            <el-select v-model="attr.attrIdAndValueId" placeholder="请选择">
              <el-option :label="attrValue.valueName" :value="`${attr.id}:${attrValue.id}`" v-for="(attrValue) in attr.attrValueList" :key="attrValue.id"></el-option>
            </el-select>
          </el-form-item>
        </el-form>
      </el-form-item>
      <el-form-item label="销售属性">
        <el-form :inline="true" ref="form1" label-width="80px">
          <el-form-item :label="saleAttr.saleAttrName" v-for="saleAttr in spuSaleAttrList" :key="saleAttr.id">
            <el-select v-model="saleAttr.attrIdAndValueId" placeholder="请选择">
              <el-option :label="asleAttrValue.saleAttrValueName" :value="`${saleAttr.id}:${asleAttrValue.id}`" v-for="asleAttrValue in saleAttr.spuSaleAttrValueList" :key="asleAttrValue.id"></el-option>
            </el-select>
          </el-form-item>
        </el-form>
      </el-form-item>
      <el-form-item label="图片列表">
        <el-table style="width: 100%" border :data="spuImageList" @selection-change="handleSelectionChange">
          <el-table-column type="selection" width="80"> </el-table-column>
          <el-table-column prop="prop" label="图片" width="width">
            <template slot-scope="{ row }">
              <img :src="row.imgUrl" alt="" style="height: 50px; width: 50px" />
            </template>
          </el-table-column>
          <el-table-column prop="imgName" label="名称" width="width">
          </el-table-column>
          <el-table-column prop="prop" label="操作" width="width">
            <template slot-scope="{row}">
              <el-button type="primary" v-if="row.isDefault==0" @click="changeDefault(row)">设置默认</el-button>
              <el-button v-else>默认</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-form-item>
      <el-form-item label="">
        <el-button type="primary" @click="save">保存</el-button>
        <el-button @click="cannel">取消</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
export default {
  name: "SkuForm",
  data() {
    return {
      // 存储图片信息
      spuImageList: [],
      // 存储销售属性信息
      spuSaleAttrList: [],
      // 平台属性信息
      attrInfoList: [],
      // 收集sku数据的字段
      skuInfo: {
        // 第一类收集的数据:父组件给的数据
        category3Id: 0,
        spuId: 0,
        tmId: 0,
        // 第二类:需要通过数据双向绑定v-model收集
        skuName: "",
        price: 0,
        weight: "",
        skuDesc: "",
        // 第三类:需要自己书写代码收集
        // 默认图片
        skuDefaultImg: "",
        skuImageList: [
          // {
          //   id: 0,
          //   imgName: "string",
          //   imgUrl: "string",
          //   isDefault: "string",
          //   skuId: 0,
          //   spuImgId: 0,
          // },
        ],
        // 平台属性
        skuAttrValueList: [
          // {
          //   attrId: 0,
          //   valueId: 0,
          // },
        ],
        // 销售属性
        skuSaleAttrValueList: [
          // {
          //   saleAttrId: 0,
          //   saleAttrValueId: 0,
          // },
        ],
      },
      spu:{},
      // 收集图片的数据字段:收集的数据缺少isDefault,提交给服务器的时候需要整理参数
      imageList:[]
    };
  },
  methods: {
    // 获取skuForm数据
    async getData(category1Id, category2Id, spu) {
      // 收集父组件给予的数据
      this.skuInfo.category3Id = spu.category3Id
      this.skuInfo.spuId = spu.id
      this.skuInfo.tmId = spu.tmId
      this.spu = spu
      // 获取图片的数据
      try {
        const result = await this.$API.spu.reqSpuImageList(spu.id);
        if (result.code == 200) {
          let list = result.data;
          list.forEach(item =>{
            item.isDefault = 0
          })
          this.spuImageList = list
        }
      } catch (error) {}
      // 获取销售属性的数据
      try {
        const result1 = await this.$API.spu.reqSpuSaleAttrList(spu.id);
        if (result1.code == 200) {
          this.spuSaleAttrList = result1.data;
        }
      } catch (error) {}
      // 获取平台属性的数据
      try {
        const result2 = await this.$API.spu.reqAttrInfoList(
          category1Id,
          category2Id,
          spu.category3Id
        );
        if (result2.code == 200) {
          this.attrInfoList = result2.data;
        }
      } catch (error) {}
    },
    // table表格复选框按钮的事件
    handleSelectionChange(selection){
      // 获取到用户选中图片的信息数据,但是当前收集的数据缺失isDefault字段
      this.imageList = selection
    },
    // 排他操作
    changeDefault(row){
      // 图片列表的数据的isDefault字段都变为0
      this.spuImageList.forEach(item => {
        item.isDefault = 0
      })
      // 点击的那个isDefault变为1
      row.isDefault = 1
      // 收集默认图片的地址
      this.skuInfo.skuDefaultImg = row.imgUrl
    },
    cannel(){
      // 触发自定义事件,让父组件切换场景为0
      this.$emit('changeScenes',0)
      // 清楚数据
      Object.assign(this._data,this.$options.data())
    },
    // 保存按钮的事件
    async save(){
      // 整理参数
      // 整理平台属性
      const {attrInfoList,skuInfo,spuSaleAttrList,imageList} = this
      // 将整理好的数据赋值给skuInfo.skuAttrValueList
      skuInfo.skuAttrValueList = attrInfoList.reduce((prev,item) => {
        if(item.attrIdAndValueId){
          const [attrId,valueId] =  item.attrIdAndValueId.split(':')
          prev.push({attrId,valueId})
        }
        return prev
      },[])
      // 整理销售属性
      skuInfo.skuSaleAttrValueList = spuSaleAttrList.reduce((prev,item) => {
        if(item.attrIdAndValueId){
          const [saleAttrId,saleAttrValueId] = item.attrIdAndValueId.split(':')
          prev.push({saleAttrId,saleAttrValueId})
        }
        return prev
      },[])
      // 整理图片的数据
      skuInfo.skuImageList = imageList.map(item =>{
        return {
          imgName:item.imgName,
          imgUrl:item.imgUrl,
          isDefault:item.isDefault,
          spuImgId:item.id,
        }
      })
      // 发请求
      try {
        const result = await this.$API.spu.reqAddsku(skuInfo)
        if (result.code == 200){
          this.$message({type:'success',message:'添加SKU成功'})
          // 触发自定义事件,让父组件切换场景为0
          this.$emit('changeScenes',0)
          // 清楚数据
          Object.assign(this._data,this.$options.data())
        }
      } catch (error) {}
    }
  },
};
</script>

<style>
</style>