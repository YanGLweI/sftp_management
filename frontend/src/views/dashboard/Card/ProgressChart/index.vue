<template>
  <div class="charts" ref="charts"></div>
</template>

<script>
// 引入Echarts
import echarts from "echarts";
export default {
  name:'ProgressChart',
  mounted(){
    let barCharts = echarts.init(this.$refs.charts);
    barCharts.setOption({
      // 提示信息
      xAxis:{
        show:false,
        // 最小值与最大值的设置
        min:0,
        max:100,
      },
      yAxis:{
        show:false,
        // 均分
        type:'category',
      },
      series:[
        {
          // 柱状图
          type:'bar',
          data:[78],
          color:'yellowgreen',
          barWidth: 10,
          // 背景颜色
          showBackground:true,
          backgroundStyle:{
            color:'#eee'
          },
          // 文本标签
          label:{
            show:true,
            // 改变文本内容
            formatter:'|',
            // 标签位置
            position:'right'
          }
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