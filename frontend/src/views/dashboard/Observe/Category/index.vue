<template>
  <div>
    <el-card>
      <div slot="header" class="header">
        <div class="category-header">
          <span>认证方式分布</span>
          <el-radio-group size="mini" v-model="radio1">
            <el-radio-button label="全部渠道"></el-radio-button>
            <el-radio-button label="内部"></el-radio-button>
            <el-radio-button label="外部"></el-radio-button>
          </el-radio-group>
        </div>
      </div>
      <div class="charts" ref="charts"></div>
    </el-card>
  </div>
</template>

<script>
import echarts from "echarts";
import { reqAuthDistribution } from '@/api/dashboard/dashboard.js'

export default {
  name: "Category",
  data() {
    return {
      radio1: "全部渠道",
      authData: [],  // 认证方式分布数据
      chartInstance: null  // 图表实例引用
    };
  },
  mounted() {
    this.initCharts();
    this.fetchAuthDistribution();
  },
  methods: {
    initCharts() {
      let myCharts = echarts.init(this.$refs.charts);
      this.chartInstance = myCharts;
      
      myCharts.setOption({
        title: {
          left: "center",
          top: "center",
        },
        tooltip: {
          trigger: "item",
          formatter: '{b}: {c} ({d}%)'  // 显示数值和百分比
        },
        // 底部图例
        legend: {
          bottom: 0,
          itemWidth: 12,
          itemHeight: 12,
          textStyle: {
            fontSize: 12
          }
        },
        series: [
          {
            name: "认证方式",
            type: "pie",
            radius: ["45%", "62%"],  // 单环样式
            // 扇区调色板：与平台 element-ui 主色系协调的柔和蓝/绿
            color: ['#79BBFF', '#95D475', '#F0C78A', '#B0A8E8'],
            avoidLabelOverlap: false,
            padAngle: 3,  // 减小扇区间间隔
            itemStyle: {
              borderRadius: 5,
            },
            label: {
              // 外标签显示：类别名和次数
              show: true,
              position: "outside",
              formatter: function(params) {
                return params.name + '\n' + params.value + '次';
              },
              fontSize: 12,
              color: "#333",
            },
            labelLine: {
              // 连接线样式：加长连接线
              show: true,
              length: 25,   // 第一段线长
              length2: 15,  // 第二段水平线长
              lineStyle: {
                width: 1.5,
                color: "#999"
              }
            },
            data: []  // 动态填充
          }
        ],
      });
      
      // 绑定事件
      myCharts.on("mouseover", (params) => {
        const {name,value} = params
        myCharts.setOption({
          title:{
            text:name,
            subtext:value.toString()
          }
        })
      });
    },
    
    async fetchAuthDistribution() {
      try {
        const res = await reqAuthDistribution()
        if (res.code === 200) {
          // 将 Map 转换为 ECharts 需要的数组格式
          this.authData = Object.entries(res.data).map(([name, value]) => ({
            value: value,
            name: name
          }))
          
          // 更新图表数据
          if (this.chartInstance) {
            this.chartInstance.setOption({
              series: [{
                data: this.authData
              }]
            })
          }
        } else {
          this.$message.warning('加载认证分布数据失败')
        }
      } catch (error) {
        console.error('加载认证分布数据失败:', error)
        this.$message.warning('加载认证分布数据失败')
      }
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
.header {
  border-bottom: 1px solid #eee;
  padding: 5px 0;
}
.category-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.charts {
  width: 100%;
  height: 300px;
}
</style>