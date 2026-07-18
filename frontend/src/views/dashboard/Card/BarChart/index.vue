<template>
  <div class="charts" ref="charts"></div>
</template>

<script>
// 引入Echarts
import echarts from "echarts";
export default {
  name:'BarChart',
  mounted(){
    let barCharts = echarts.init(this.$refs.charts);
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
          data:[5, 20, 36, 46, 10, 20, 4, 41, 20, 36, 46, 10, 20,],
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
}
</script>

<style scoped>
.charts{
  width: 100%;
  height: 100%;
}
</style>