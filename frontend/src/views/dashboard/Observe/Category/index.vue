<template>
  <div>
    <el-card>
      <div slot="header" class="header">
        <div class="category-header">
          <span>访问渠道占比</span>
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
export default {
  name: "Category",
  data() {
    return {
      radio1: "全部渠道",
    };
  },
  mounted() {
    let myCharts = echarts.init(this.$refs.charts);
    myCharts.setOption({
      title: {
        left: "center",
        top: "center",
      },
      tooltip: {
        trigger: "item",
      },
      series: [
        {
          name: "Access From",
          type: "pie",
          radius: ["40%", "70%"],
          avoidLabelOverlap: false,
          padAngle: 5,
          itemStyle: {
            borderRadius: 10,
          },
          label: {
            show: true,
            position: "outside",
          },
          labelLine: {
            show: true,
          },
          data: [
            { value: 1048, name: "密钥" },
            { value: 735, name: "密码" },
            { value: 580, name: "PIN" },
            { value: 484, name: "ID" },
            { value: 300, name: "code" },
          ],
        },
      ],
    });
    // 绑定事件
    myCharts.on("mouseover", function (params) {
      // 结构出name和value
      const {name,value} = params
      // 鼠标移上后,修改标题和子标题
      myCharts.setOption({
        title:{
          text:name,
          subtext:value
        }
      })
    });
  },
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