<template>
  <div>
    <div class="header">
      <span class="search-header">搜索用户数</span><i class="el-icon-info"></i>
    </div>
    <div class="main">
      <span class="main-title">12321 </span>
      <span class="main-conten">17.1</span>
      <svg t="1730193703168" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="5432" width="16" height="16"><path d="M239.616 634.88 507.904 346.112 776.192 634.88Z" fill="#d81e06" p-id="5433"></path></svg>
      <svg t="1730194239993" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="6730" width="16" height="16"><path d="M193 320.667 515.149 640.851 834 321Z" fill="#1afa29" p-id="6731"></path></svg>
    </div>
    <div class="footer">
      <div class="charts" ref="charts"></div>
    </div>
  </div>
</template>

<script>
// 引入Echarts
import echarts from "echarts";
export default {
  name:'',
  mounted() {
    // 初始化echarts
    let lineCharts = echarts.init(this.$refs.charts);
    // 配置数据
    lineCharts.setOption({
      xAxis: {
        // 隐藏X轴
        show: false,
        type: "category",
      },
      yAxis: {
        // 隐藏Y轴
        show: false,
      },
      // 系列
      series: [
        {
          type: "line",
          data: [10, 22, 8, 64, 16, 35,44,9,71,54,57,100],
          // 折线拐点标志的样式
          itemStyle: {
            opacity: 0,
          },
          // 线条样式
          lineStyle: {
            color: "skyblue",
          },
          // 区域填充样式
          areaStyle: {
            // 渐变颜色
            color: {
              type: "linear",
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [
                {
                  offset: 0,
                  color: "skyblue", // 0% 处的颜色
                },
                {
                  offset: 1,
                  color: "#fff", // 100% 处的颜色
                },
              ],
              global: false, // 缺省为 false
            },
          },
          // 是否平滑曲线显示
          smooth:true
        },
      ],
      // 布局调试
      grid: {
        left: 0,
        top: 0,
        right: 0,
        bottom: 0,
      },
    });
  },
}
</script>

<style scoped>
.header{
  display: flex;
  align-items: center;
}
.search-header{
  margin-right: 20px;
}
.main{
  display: flex;
  align-items: center;
  margin: 10px 0;
}
.main-title{
  margin-right: 30px;
}
.main-conten{
  margin-right: 5px;
}
.charts{
  width: 100%;
  height: 50px;
  margin-bottom: 10px;
}
</style>