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
          <h3><span class="title-decorator"></span>SFTP{{ title }}排行 Top6</h3>
          <ul v-if="userRankList && userRankList.length > 0">
            <li class="rank-row" v-for="(user, index) in userRankList" :key="index">
              <div class="rank-bar" :class="{ 'rank-bar--low': index >= 3 }" :style="{ width: barPercent(user.count) }"></div>
              <span class="rank-badge" :class="'rank-badge--' + (index + 1)">{{ user.rank }}</span>
              <span class="rank-name" :title="user.username">{{ user.username }}</span>
              <span class="rank-value">{{ user.count.toLocaleString() }}</span>
            </li>
          </ul>
          <!-- 空状态提示 -->
          <div v-else class="empty-tip">暂无数据</div>
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
import { reqTopTransferUsers } from '@/api/dashboard/dashboard.js'

export default {
  name: "Sale",
  data() {
    return {
      activeName: "sale",
      barCharts: null,
      // 收集日历数据
      date: [],
      userRankList: []
    };
  },
  computed: {
    // 标题
    title() {
      return this.activeName == "sale" ? "传输量" : "访问量";
    },
    ...mapState({
      listState: (state) => state.home.list,
    }),
    // 获取最大数量用于计算条形宽度比例
    maxCount() {
      if (!this.userRankList || this.userRankList.length === 0) return 1;
      const counts = this.userRankList.map(item => item.count);
      return Math.max(...counts, 1);
    }
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
    // 获取传输量排行列表
    async fetchUserRankList() {
      try {
        const res = await reqTopTransferUsers()
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
    
    // 计算条形宽度百分比（保底 8% 避免 0 值行完全无条）
    barPercent(count) {
      const percentage = Math.max((count / this.maxCount) * 100, 8);
      return `${percentage}%`;
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
  font-size: 14px;
}
.title-decorator {
  display: inline-block;
  width: 3px;
  height: 14px;
  background: linear-gradient(to bottom, #409eff, #3a8ee6);
  border-radius: 2px;
  vertical-align: middle;
  margin-right: 6px;
}
ul {
  list-style: none;
  width: 100%;
  height: 250px;
  padding: 0;
  box-sizing: border-box;
}
ul > li {
  position: relative;
  height: 34px;
  margin: 0;
  padding: 0;
  overflow: hidden;
}

/* 排行榜条形图样式 */
.rank-row {
  display: flex;
  align-items: center;
  position: relative;
  padding: 0 10px;
  background: #f5f7fa;
  border-radius: 6px;
  margin-bottom: 6px;
}

.rank-bar {
  position: absolute;
  left: 0;
  top: 2px;
  bottom: 2px;
  width: 0;
  background: linear-gradient(90deg, rgba(64, 158, 255, 0.22), rgba(64, 158, 255, 0.08));
  border-radius: 6px;
  z-index: 0;
  transition: width 0.6s ease-out;
  pointer-events: none;
}

.rank-bar.rank-bar--low {
  background: linear-gradient(90deg, rgba(64, 158, 255, 0.15), rgba(64, 158, 255, 0.05));
}

.rank-badge {
  width: 24px;
  height: 24px;
  border-radius: 12px;
  background: #e4e7ed;
  color: #707d92;
  text-align: center;
  line-height: 24px;
  font-size: 12px;
  font-weight: bold;
  z-index: 1;
  margin-right: 10px;
  flex-shrink: 0;
}

.rank-badge.rank-badge--1 {
  background: linear-gradient(135deg, #F7B500, #FFD700);
  color: white;
  box-shadow: 0 2px 6px rgba(247, 181, 0, 0.3);
}

.rank-badge.rank-badge--2 {
  background: linear-gradient(135deg, #A3ADC2, #C0C8D8);
  color: white;
  box-shadow: 0 2px 6px rgba(163, 173, 194, 0.3);
}

.rank-badge.rank-badge--3 {
  background: linear-gradient(135deg, #D98E5F, #E8AA85);
  color: white;
  box-shadow: 0 2px 6px rgba(217, 142, 95, 0.3);
}

.rank-name {
  flex: 1;
  font-size: 13px;
  color: #303133;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  z-index: 1;
}

.rank-value {
  font-size: 13px;
  color: #606266;
  font-weight: 600;
  min-width: 70px;
  text-align: right;
  z-index: 1;
  font-feature-settings: "tnum";
}

.empty-tip {
  text-align: center;
  color: #909399;
  padding: 60px 0;
  font-size: 13px;
}
</style>
