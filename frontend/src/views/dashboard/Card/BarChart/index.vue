<template>
  <div class="charts" ref="charts"></div>
</template>

<script>
// 引入 Echarts
import echarts from "echarts";
export default {
  name:'BarChart',
  // 接收外部传入的柱状图数据
  props: {
    data: {
      type: Array,
      default: () => []
    }
  },
  mounted(){
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
    initChart(){
      let barCharts = echarts.init(this.$refs.charts);
      this.chartInstance = barCharts;
      barCharts.setOption({
        // 提示信息
        tooltip: {},
        xAxis:{
          show:false,
          // 均分
          type:'category'
        },
        yAxis:{
          show:false
        },
        series:[
          {
            // 柱状图
            type:'bar',
            data: this.data && this.data.length > 0 ? this.data : [0],
            color:'skyblue'
          }
        ],
        // 布局调试
        grid:{
          left: 0,
          top: 0,
          right: 0,
          bottom: 0,
        }
      })
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
.charts{
  width: 100%;
  height: 100%;
}
</style>