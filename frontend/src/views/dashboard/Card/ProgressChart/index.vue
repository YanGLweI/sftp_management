<template>
  <div class="charts" ref="charts"></div>
</template>

<script>
// 引入 Echarts
import echarts from "echarts";
export default {
  name:'ProgressChart',
  // 接收外部传入的进度值（0-100），默认 0
  props: {
    value: {
      type: Number,
      default: 0
    }
  },
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
          data:[this.value],
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
  },
  watch: {
    // 数据变化时更新图表
    value(newVal) {
      if (this.$refs.charts) {
        let barCharts = echarts.getInstanceByDom(this.$refs.charts)
        if (barCharts) {
          barCharts.setOption({
            series: [{
              data: [newVal]
            }]
          })
        }
      }
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