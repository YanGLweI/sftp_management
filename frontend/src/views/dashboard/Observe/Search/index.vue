<template>
  <div>
    <el-card>
      <div slot="header" class="header">
        <div class="search-header">
          <span>线上热门搜索</span>
          <el-dropdown>
            <span class="el-dropdown-link">
              <i class="el-icon-more"></i>
            </span>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item>mia</el-dropdown-item>
              <el-dropdown-item>liumin</el-dropdown-item>
              <el-dropdown-item>ylw</el-dropdown-item>
              <el-dropdown-item>longliuming</el-dropdown-item>
              <el-dropdown-item>wangpeng</el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
        </div>
      </div>
      <div>
        <el-row :gutter="10">
          <el-col :span="12">
            <LineCharts></LineCharts>
          </el-col>
          <el-col :span="12">
            <LineCharts></LineCharts>
          </el-col>
        </el-row>
        <!-- table表格 -->
        <el-table :data="tableData" border style="width: 100%">
          <el-table-column type="index" label="排名" width="80" align="center">
          </el-table-column>
          <el-table-column label="搜索关键字" width="180">
          </el-table-column>
          <el-table-column label="用户数" sortable>
          </el-table-column>
          <el-table-column label="周涨幅" sortable>
          </el-table-column>
        </el-table>
        <!-- 分页器 -->
        <el-pagination
          align="center"
          layout="prev, pager, next"
          :total="1000">
        </el-pagination>
      </div>
    </el-card>
  </div>
</template>

<script>
import LineCharts from "./lineCharts";
export default {
  name: "Search",
  data() {
    return {
      tableData:[{}]
    }
  },
  components: { LineCharts },
};
</script>

<style scoped>
.header {
  border-bottom: 1px solid #eee;
  padding: 5px 0;
}
.search-header {
  display: flex;
  justify-content: space-between;
}
.el-dropdown-link {
  cursor: pointer;
}
</style>