<template>
  <div>
    <el-row :gutter="10">
        <el-col :span="6">
          <el-card>
            <!-- 第一个Detail -->
            <Detail title="账号总数" :count="listState.accountCount">
              <template slot="charts">
                <div style="display: flex;align-items: center;">
                  <span style="font-size: 12px;">本月新增 {{ listState.monthlyNewCount }} 人 </span>
                  <svg t="1730193703168" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="5432" width="16" height="16">
                    <path d="M239.616 634.88 507.904 346.112 776.192 634.88Z" fill="#d81e06" p-id="5433"></path>
                  </svg>
                  <span style="font-size: 12px;padding-left: 10px;">总账号 {{ listState.accountCount }} 个 </span>
                  <svg t="1730194239993" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="6730" width="16" height="16">
                    <path d="M193 320.667 515.149 640.851 834 321Z" fill="#1afa29" p-id="6731"></path>
                  </svg>
                </div>
              </template>
              <template slot="footer">
                <span>月新增 {{ listState.monthlyNewCount }}</span>
              </template>
            </Detail>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <!-- 第二个 Detail -->
            <Detail title="访问量" :count="totalAccess">
              <template slot="charts">
                <LineChart :data="accessTrend"></LineChart>
              </template>
              <template slot="footer">
                <span>今日访问 {{ visitCount }} 次 (同比增长 {{ growthRate }}%)</span>
              </template>
            </Detail>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <!-- 第三个 Detail -->
            <Detail title="传输数" :count="totalTransfer">
              <template slot="charts">
                <BarChart :data="transferTrend"></BarChart>
              </template>
              <template slot="footer">
                <span>今日传输 {{ transferCount }} 次</span>
              </template>
            </Detail>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <!-- 第四个 Detail -->
            <Detail title="访问量增长效果" :count="growthRate + '%'">
              <template slot="charts">
                <ProgressChart :value="Math.min(Math.abs(Number(growthRate) || 0), 100)"></ProgressChart>
              </template>
              <template slot="footer">
                <div style="display: flex; align-items: center;">
                  <span style="font-size: 16px;">今日 {{ visitCount }} 次</span>
                  <svg t="1730193703168" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="5432" width="16" height="16"><path d="M239.616 634.88 507.904 346.112 776.192 634.88Z" fill="#d81e06" p-id="5433"></path></svg>
                  <span style="font-size: 16px;padding-left: 10px;">日增 {{ growthRate }}% </span>
                  <svg t="1730194239993" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="6730" width="16" height="16"><path d="M193 320.667 515.149 640.851 834 321Z" fill="#1afa29" p-id="6731"></path></svg>
                </div>
              </template>
            </Detail>
          </el-card>
        </el-col>
    </el-row>
  </div>
</template>

<script>
import Detail from './Detail'
import LineChart from './LineChart'
import BarChart from './BarChart'
import ProgressChart from './ProgressChart'
import { mapState } from "vuex";
import { reqTotalAccessCount, reqTodayAccessCount, reqTotalTransferCount, reqTodayTransferCount, reqAccessGrowth } from '@/api/dashboard/dashboard.js'

export default {
  name:'Card',
  components:{Detail,LineChart,BarChart,ProgressChart},
  data(){
    return{
      cardData:'',
      totalAccess: 0,     // 累计访问量（总数）
      visitCount: 0,      // 今日访问量
      totalTransfer: 0,   // 累计传输数（总数）
      transferCount: 0,   // 今日传输量
      growthRate: 0,      // 增长率（百分比数值部分）
      accessTrend: [],    // 访问量 7 天趋势
      transferTrend: []   // 传输量 7 天趋势
    }
  },
  mounted(){
    this.fetchDashboardStats()
  },
  computed: {
    ...mapState({
      listState: (state) => state.home.list,
    }),
  },
  watch: {
    // 监听 store 数据到达（getData 异步完成），更新迷你图 7 天趋势
    listState: {
      deep: true,
      handler() {
        if (this.listState && this.listState.accessFullDay) {
          this.accessTrend = this.listState.accessFullDay.map(Number)
        }
        if (this.listState && this.listState.transFullDay) {
          this.transferTrend = this.listState.transFullDay.map(Number)
        }
      }
    }
  },
  methods:{
    // 获取 Dashboard 统计数据
    async fetchDashboardStats(){
      try {
        const [totalAccessRes, accessRes, totalTransferRes, transferRes, growthRes] = await Promise.all([
          reqTotalAccessCount(),
          reqTodayAccessCount(),
          reqTotalTransferCount(),
          reqTodayTransferCount(),
          reqAccessGrowth()
        ])
        
        if (totalAccessRes.code === 200) {
          this.totalAccess = totalAccessRes.data
        }
        
        if (accessRes.code === 200) {
          this.visitCount = accessRes.data
        }
        
        if (totalTransferRes.code === 200) {
          this.totalTransfer = totalTransferRes.data
        }
        
        if (transferRes.code === 200) {
          this.transferCount = transferRes.data
        }
        
        // 使用后端计算的真实增长率（今日 vs 昨日）
        if (growthRes.code === 200 && growthRes.data) {
          this.growthRate = growthRes.data.growthRate.toFixed(2)
        }
        
        // 从 vuex 获取 7 天趋势数据（字符串数组转数字）
        if (this.listState && this.listState.accessFullDay) {
          this.accessTrend = this.listState.accessFullDay.map(Number)
        }
        if (this.listState && this.listState.transFullDay) {
          this.transferTrend = this.listState.transFullDay.map(Number)
        }
        
      } catch (error) {
        console.error('加载统计数据失败:', error)
        this.$message.warning('加载统计数据失败，显示占位数据')
      }
    }
  }
}
</script>

<style>

</style>