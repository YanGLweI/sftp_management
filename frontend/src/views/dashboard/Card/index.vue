<template>
  <div>
    <el-row :gutter="10">
        <el-col :span="6">
          <el-card>
            <!-- 第一个Detail -->
            <Detail title="账号总数" :count="listState.accountCount">
              <template slot="charts">
                <div style="display: flex;align-items: center;">
                  <span style="font-size: 12px;">周同比 56.67% </span>
                  <svg t="1730193703168" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="5432" width="16" height="16">
                    <path d="M239.616 634.88 507.904 346.112 776.192 634.88Z" fill="#d81e06" p-id="5433"></path>
                  </svg>
                  <span style="font-size: 12px;padding-left: 10px;">日同比 19.96% </span>
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
            <!-- 第二个Detail -->
            <Detail title="访问量" count="88460">
              <template slot="charts">
                <LineChart></LineChart>
              </template>
              <template slot="footer">
                <span>日访问量 1234</span>
              </template>
            </Detail>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <!-- 第三个Detail -->
            <Detail title="传输数" count="88460">
              <template slot="charts">
                <BarChart></BarChart>
              </template>
              <template slot="footer">
                <span>转化率 65%</span>
              </template>
            </Detail>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <!-- 第四个Detail -->
            <Detail title="访问量增长效果" count="78%">
              <template slot="charts">
              <ProgressChart></ProgressChart>
              </template>
              <template slot="footer">
                <div style="display: flex; align-items: center;">
                  <span style="font-size: 16px;">周 56.67%</span>
                  <svg t="1730193703168" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="5432" width="16" height="16"><path d="M239.616 634.88 507.904 346.112 776.192 634.88Z" fill="#d81e06" p-id="5433"></path></svg>
                  <span style="font-size: 16px;padding-left: 10px;">日 19.96% </span>
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
export default {
  name:'Card',
  components:{Detail,LineChart,BarChart,ProgressChart},
  data(){
    return{
      cardData:''
    }
  },
  mounted(){
    // this.getAccountCount()
  },
  computed: {
    ...mapState({
      listState: (state) => state.home.list,
    }),
  },
  methods:{
    // 获取账号总数
    async getAccountCount(){
      try {
        const result = await this.$API.dashboard.reqAccountCount()
        if (result.code == 200){
          this.cardData = result.data
        }
      } catch (error) {
        this.$message.error('获取账号总数失败')
      }
      
    }
  }
}
</script>

<style>

</style>