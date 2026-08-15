<template>
  <div>
    <div class="header">
      <span class="search-header">{{ title }}</span><i class="el-icon-info"></i>
    </div>
    <div class="main">
      <span class="main-title">{{ value }} </span>
      <span class="main-conten">{{ rate }}</span>
      <svg t="1730193703168" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="5432" width="16" height="16"><path d="M239.616 634.88 507.904 346.112 776.192 634.88Z" fill="#d81e06" p-id="5433"></path></svg>
      <svg t="1730194239993" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="6730" width="16" height="16"><path d="M193 320.667 515.149 640.851 834 321Z" fill="#1afa29" p-id="6731"></path></svg>
    </div>
    <div class="footer">
      <div class="charts" ref="charts"></div>
    </div>
  </div>
</template>

<script>
// 引入 Echarts
import echarts from "echarts";
export default {
  name:'SearchLineChart',
  props: {
    // 标题，如：今日访问 / 今日传输
    title: {
      type: String,
      default: "统计"
    },
    // 主数字
    value: {
      type: [Number, String],
      default: 0
    },
    // 增长率百分比
    rate: {
      type: [Number, String],
      default: "0"
    },
    // 折线图数据（7 天趋势）
    chartData: {
      type: Array,
      default: () => []
    }
  },
  mounted() {
    this.initChart();
  },
  watch: {
    // 数据变化时更新图表
    chartData: {
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
            data: this.chartData && this.chartData.length > 0 ? this.chartData : [0],
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