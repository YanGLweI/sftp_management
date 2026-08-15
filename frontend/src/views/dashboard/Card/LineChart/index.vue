<template>
  <!-- 容器 -->
  <div class="charts" ref="charts"></div>
</template>

<script>
// 引入 Echarts
import echarts from "echarts";
export default {
  name: "LineChart",
  // 接收外部传入的折线数据
  props: {
    data: {
      type: Array,
      default: () => []
    }
  },
  mounted() {
    this.initChart();
  },
  watch: {
    // 数据变化时更新图表
    data: {
      deep: true,
      handler(newVal) {
        if (this.$refs.charts && this.chartInstance) {
          this.chartInstance.setOption({
            series: [{
              data: newVal
            }]
          })
        } else if (this.$refs.charts) {
          this.initChart()
        }
      }
    }
  },
  methods: {
    initChart() {
      // 初始化 echarts
      let lineCharts = echarts.init(this.$refs.charts);
      this.chartInstance = lineCharts;
      // 配置数据
      lineCharts.setOption({
        xAxis: {
          // 隐藏 X 轴
          show: false,
          type: "category",
        },
        yAxis: {
          // 隐藏 Y 轴
          show: false,
        },
        // 系列
        series: [
          {
            type: "line",
            data: this.data && this.data.length > 0 ? this.data : [0],
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
            smooth: true
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
    }
  },
  beforeDestroy() {
    if (this.chartInstance) {
      this.chartInstance.dispose()
    }
  }
};
</script>

<style scoped>
.charts {
  width: 100%;
  height: 100%;
}
</style>