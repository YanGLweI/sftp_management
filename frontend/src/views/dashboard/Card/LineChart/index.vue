<template>
  <!-- 容器 -->
  <div class="charts" ref="charts"></div>
</template>

<script>
// 引入Echarts
import echarts from "echarts";
export default {
  name: "LineChart",
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
            color: "purple",
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
                  color: "purple", // 0% 处的颜色
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
};
</script>

<style scoped>
.charts {
  width: 100%;
  height: 100%;
}
</style>