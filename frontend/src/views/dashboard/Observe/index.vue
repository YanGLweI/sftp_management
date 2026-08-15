<template>
  <div>
    <el-row :gutter="10" class="observe-row">
      <el-col :span="12" class="observe-col">
        <Search></Search>
      </el-col>
      <el-col :span="12" class="observe-col">
        <Category></Category>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import Search from "./Search";
import Category from "./Category";
export default {
  name: "Observe",
  components: { Search, Category },
};
</script>

<style>
/* 左右两列弹性等高：由内容更高的卡片（右侧图表）决定高度 */
.observe-row {
  display: flex;
}
.observe-col {
  display: flex;
}
.observe-col > div {
  width: 100%;
  display: flex;
  flex-direction: column;
}
.observe-col .el-card {
  flex: 1;
  width: 100%;
  display: flex;
  flex-direction: column;
}
.observe-col .el-card__body {
  flex: 1;
  display: flex;
  flex-direction: column;
}
</style>