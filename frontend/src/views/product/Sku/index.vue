<template>
  <div>
    <!-- 表格 -->
    <el-table style="width: 100%" border :data="records">
      <el-table-column type="index" label="序号" width="80" align="center">
      </el-table-column>
      <el-table-column prop="skuName" label="名称" width="width">
      </el-table-column>
      <el-table-column prop="skuDesc" label="描述" width="width">
      </el-table-column>
      <el-table-column label="默认图片" width="110">
        <template slot-scope="{ row }">
          <img
            :src="row.skuDefaultImg"
            alt=""
            style="width: 50px; height: 50px"
          />
        </template>
      </el-table-column>
      <el-table-column prop="weight" label="重量(KG)" width="80">
      </el-table-column>
      <el-table-column prop="price" label="价格(元)" width="80">
      </el-table-column>
      <el-table-column prop="prop" label="操作" width="width">
        <template slot-scope="{ row }">
          <el-tooltip
            class="item"
            effect="dark"
            content="下架"
            placement="top"
            v-if="row.isSale == 1"
          >
            <el-button
              type="info"
              icon="el-icon-bottom"
              size="mini"
              @click="cancel(row)"
            ></el-button>
          </el-tooltip>
          <el-tooltip
            class="item"
            effect="dark"
            content="上架"
            placement="top"
            v-else
          >
            <el-button
              type="success"
              icon="el-icon-top"
              size="mini"
              @click="sale(row)"
            ></el-button>
          </el-tooltip>
          <el-button
            type="primary"
            icon="el-icon-edit"
            size="mini"
            @click="edit"
          ></el-button>
          <el-tooltip
            class="item"
            effect="dark"
            content="查看信息"
            placement="top"
          >
            <el-button
              type="info"
              icon="el-icon-info"
              size="mini"
              @click="getSkuInfo(row)"
            ></el-button>
          </el-tooltip>
          <el-button
            type="danger"
            icon="el-icon-delete"
            size="mini"
          ></el-button>
        </template>
      </el-table-column>
    </el-table>
    <!-- 分页器       @size-change="handleSizeChange"
      @current-change="handleCurrentChange"-->
    <el-pagination
      @current-change="getSkuList"
      @size-change="handleSizeChange"
      style="text-align: center"
      :current-page="page"
      :page-sizes="[3, 5, 10]"
      :page-size="limit"
      layout="prev, pager, next, jumper, ->,sizes,total "
      :total="total"
    >
    </el-pagination>
    <!-- 抽屉 -->
    <el-drawer
      :visible.sync="show"
      :show-close="false"
      size="50%"
      @close="close"
    >
      <el-row>
        <el-col :span="5">名称</el-col>
        <el-col :span="16">{{ skuInfo.skuName }}</el-col>
      </el-row>
      <el-row>
        <el-col :span="5">描述</el-col>
        <el-col :span="16">{{ skuInfo.skuDesc }}</el-col>
      </el-row>
      <el-row>
        <el-col :span="5">价格</el-col>
        <el-col :span="16">{{ skuInfo.price }} 元</el-col>
      </el-row>
      <el-row>
        <el-col :span="5">平台属性</el-col>
        <el-col :span="16">
          <el-tag
            type="success"
            v-for="attr in skuInfo.skuAttrValueList"
            :key="attr.id"
            style="margin-right: 5px"
            >{{ attr.attrId }}+{{ attr.valueId }}</el-tag
          >
        </el-col>
      </el-row>
      <el-row>
        <el-col :span="5">商品图片</el-col>
        <el-col :span="16">
          <el-carousel height="400px">
            <el-carousel-item
              v-for="item in skuInfo.skuImageList"
              :key="item.id"
            >
              <img :src="item.imgUrl" alt="" style="height: 400px;margin-left: 31px;" />
            </el-carousel-item>
          </el-carousel>
        </el-col>
      </el-row>
    </el-drawer>
  </div>
</template>

<script>
export default {
  name: "Sku",
  data() {
    return {
      page: 1, //当前第几页
      limit: 10, //当前页面有几条数据
      records: [], //存储SKU列表的数据
      total: 0, //存储分页器一共展示的数据
      skuInfo: {}, //存储SKU信息
      show: false, //控制抽屉显示
    };
  },
  mounted() {
    // 获取sku列表
    this.getSkuList();
  },
  methods: {
    async getSkuList(pages = 1) {
      this.page = pages;
      // 结构出默认参数
      const { page, limit } = this;
      try {
        let result = await this.$API.sku.reqSkuList(page, limit);
        if (result.code == 200) {
          this.total = result.data.total;
          this.records = result.data.records;
        }
      } catch (error) {}
    },
    handleSizeChange(limit) {
      // 修改参数
      this.limit = limit;
      this.getSkuList();
    },
    // 下架按钮的回调
    async cancel(row) {
      try {
        const result = await this.$API.sku.reqCancel(row.id);
        if (result.code == 200) {
          this.getSkuList(this.page);
          this.$message({ type: "success", message: "下架成功" });
        }
      } catch (error) {}
    },
    // 上架
    async sale(row) {
      try {
        const result = await this.$API.sku.reqSale(row.id);
        if (result.code == 200) {
          this.getSkuList(this.page);
          this.$message({ type: "success", message: "上架成功" });
        }
      } catch (error) {}
    },
    edit() {
      this.$message("功能正在开发中...");
    },
    // 查看sku信息按钮的回调
    async getSkuInfo(row) {
      // 展示抽屉
      this.show = true;
      // 获取sku数据
      try {
        const result = await this.$API.sku.reqSkuById(row.id);
        if (result.code == 200) {
          this.skuInfo = result.data;
        }
      } catch (error) {}
    },
    // 抽到关闭的回调
    close() {
      this.skuInfo = {};
    },
  },
};
</script>

<style scoped>
.el-row .el-col-5 {
  font-size: 18px;
  font-weight: bold;
  text-align: right;
}
.el-col {
  margin: 10px;
}
</style>

<style>
.el-carousel__button{
  background: red;
  width: 10px;
  height: 10px;
  border-radius: 50%;
}
</style>