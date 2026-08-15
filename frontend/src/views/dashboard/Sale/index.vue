<template>
  <el-card class="box-card" style="margin: 10px 0">
    <div slot="header" class="clearfix">
      <!-- v-model="activeName" @tab-click="handleClick" -->
      <!-- 头部左侧内容 -->
      <el-tabs class="tab" v-model="activeName">
        <el-tab-pane label="传输量" name="sale"></el-tab-pane>
        <el-tab-pane label="访问量" name="visits"></el-tab-pane>
      </el-tabs>
      <!-- 头部右侧内容 -->
      <div class="right">
        <span @click="setDay">今日</span>
        <span @click="setWeek">本周</span>
        <span @click="setMonth">本月</span>
        <span @click="setYear">本年</span>
        <!-- v-model="value1" -->
        <el-date-picker
          v-model="date"
          value-format="yyyy-MM-dd HH:mm:ss"
          class="date"
          size="mini"
          type="datetimerange"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
        >
        </el-date-picker>
      </div>
    </div>
    <div>
      <el-row :gutter="10">
        <el-col :span="18">
          <!-- 容器 -->
          <div class="charts" ref="charts"></div>
        </el-col>
        <el-col :span="6" class="right1">
          <h3>SFTP{{ title }}排行 Top10</h3>
          <ul>
            <li v-for="(user, index) in userRankList" :key="index">
              <span :class="index < 3 ? 'rindex' : 'rindex1'">{{ user.rank }}</span>
              <span>{{ user.username }}</span>
              <span class="rvalue">{{ user.count.toLocaleString() }}</span>
            </li>
            <!-- 兼容处理：如果数据未加载完成或为空 -->
            <template v-if="!userRankList || userRankList.length === 0">
              <li v-for="i in 7" :key="i">
                <span :class="i <= 3 ? 'rindex' : 'rindex1'">{{ i }}</span>
                <span>暂无数据</span>
                <span class="rvalue">0</span>
              </li>
            </template>
          </ul>
        </el-col>
      </el-row>
    </div>
  </el-card>
</template>

<script>
// 引入 Echarts
import echarts from "echarts";
// 引入 dayjs
import dayjs from "dayjs";

import { mapState } from "vuex";
import { reqActiveUsersTop10 } from '@/api/dashboard/dashboard.js'

export default {
  name: "Sale",
  data() {
    return {
      activeName: "sale",
      barCharts: null,
      // 收集日历数据
      date: [],
      userRankList: []  // 新增：用户排行榜数据
    };
  },
  mounted() {
    // 初始化 echarts 实例
    this.barCharts = echarts.init(this.$refs.charts);
    // 配置数据
    this.barCharts.setOption({
      title: {
        text: this.title + "趋势",
      },
      tooltip: {
        trigger: "axis",
        axisPointer: {
          type: "shadow",
        },
      },
      grid: {
        left: "3%",
        right: "4%",
        bottom: "3%",
        containLabel: true,
      },
      xAxis: [
        {
          type: "category",
          data: [],
          axisTick: {
            alignWithLabel: true,
          },
        },
      ],
      yAxis: [
        {
          type: "value",
        },
      ],
      series: [
        {
          name: "Direct",
          type: "bar",
          barWidth: "60%",
          data: [],
          color: "orange",
          barWidth: "25px",
        },
      ],
    });
    
    // 加载用户排行榜
    this.fetchUserRankList();
    
    // 防御性渲染：如果 store 已有数据（getData 早于本组件完成），立即绘制趋势图
    // 避免 watch listState 因数据未发生"变化"而不触发，导致图表空白
    this.$nextTick(() => {
      if (this.listState && this.listState.transXaxis && this.listState.transXaxis.length > 0) {
        this.barCharts.setOption({
          xAxis: [
            {
              data: this.listState.transXaxis,
            },
          ],
          series: [
            {
              data: this.listState.transFullDay,
            },
          ],
        });
      }
    });
  },
  computed: {
    // 标题
    title() {
      return this.activeName == "sale" ? "传输量" : "访问量";
    },
    ...mapState({
      listState: (state) => state.home.list,
    }),
  },
  // 监听
  watch: {
    // 监听 title 变化。改变数据
    title() {
      // 重新修改图表配置数据
      // 图表的配置数据可以再次修改，新数值会替换旧数值
      this.barCharts.setOption({
        title: {
          text: this.title + "趋势",
        },
        xAxis: {
          data:
            this.title == "传输量"
              ? this.listState.transXaxis
              : this.listState.accessXaxis,
        },
        series: [
          {
            data:
              this.title == "传输量"
                ? this.listState.transFullDay
                : this.listState.accessFullDay,
          },
        ],
      });
    },
    // 刚打开首页是的数据
    listState() {
      // 配置数据
      this.barCharts.setOption({
        xAxis: [
          {
            data: this.listState.transXaxis,
          },
        ],
        series: [
          {
            data: this.listState.transFullDay,
          },
        ],
      });
    },
  },
  methods: {
    // 获取活跃用户列表
    async fetchUserRankList() {
      try {
        const res = await reqActiveUsersTop10()
        if (res.code === 200) {
          this.userRankList = res.data.map((item, index) => ({
            rank: index + 1,
            username: item.username,
            count: item.count
          }))
        }
      } catch (error) {
        console.error('加载用户排行榜失败:', error)
        // 如果请求失败，保留空数组，不显示默认假数据
        this.userRankList = []
      }
    },
    
    // 今日
    setDay() {
      dayjs();
      const dayStart = dayjs().format("YYYY-MM-DD 00:00:00");
      const dayEnd = dayjs().format("YYYY-MM-DD 23:59:59");
      this.date = [dayStart, dayEnd];
    },
    // 本周
    setWeek() {
      // 获取本周开始日期（周日）
      const startOfWeek = dayjs().startOf("week").format("YYYY-MM-DD HH:mm:ss");
      // 获取本周结束日期（周六）
      const endOfWeek = dayjs().endOf("week").format("YYYY-MM-DD HH:mm:ss");
      this.date = [startOfWeek, endOfWeek];
    },
    // 本月
    setMonth() {
      const startOfMonth = dayjs()
        .startOf("Month")
        .format("YYYY-MM-DD HH:mm:ss");
      const endOfMonth = dayjs().endOf("Month").format("YYYY-MM-DD HH:mm:ss");
      this.date = [startOfMonth, endOfMonth];
    },
    // 本年
    setYear() {
      const startOfYear = dayjs().startOf("year").format("YYYY-MM-DD HH:mm:ss");
      const endOfYear = dayjs().endOf("year").format("YYYY-MM-DD HH:mm:ss");
      this.date = [startOfYear, endOfYear];
    },
  },
};
</script>

<style scoped>
.clearfix {
  position: relative;
  display: flex;
  justify-content: space-between;
}
.tab {
  width: 100%;
}
.right {
  position: absolute;
  right: 0px;
}
.date {
  width: 350px;
  margin: 0 20px;
}
.right span {
  margin: 0 10px;
}
.charts {
  width: 100%;
  height: 300px;
}
.right1 > h3 {
  margin: 0;
  padding: 0;
}
ul {
  list-style: none;
  width: 100%;
  height: 250px;
  padding: 0;
}
ul > li {
  height: 8%;
  margin: 15px 0;
}
.rindex {
  float: left;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: black;
  color: white;
  text-align: center;
  margin-right: 20px;
}
.rindex1 {
  float: left;
  width: 20px;
  height: 20px;
  text-align: center;
  margin-right: 20px;
}
.rvalue {
  float: right;
}
</style>