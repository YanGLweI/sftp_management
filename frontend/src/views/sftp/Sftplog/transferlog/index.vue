<template>
  <div>
    <el-card shadow="hover">
      <div slot="header">
        <el-date-picker 
          v-model="date" 
          value-format="yyyyMMdd"
          type="date" 
          placeholder="选择日期"
          :picker-options="pickerOptions"
          @change="getTransferLog"
        ></el-date-picker>
        <div style="float: right;">
          <!-- 搜索 -->
          <el-input placeholder="请输入搜索内容" style="width: 400px;" prefix-icon="el-icon-search" v-model="input" @keyup.enter.native="search">
            <el-button slot="append" icon="el-icon-search" @click="search"></el-button>
          </el-input>
          <!-- 重置搜索 -->
          <el-button icon="el-icon-refresh" circle style="margin-left: 10px;" size="mini" @click="reset"></el-button>
        </div>
      </div>
      <div>
        <el-empty v-loading="loading" :description="description" v-if="logList.length > 0 ? false : true" style="height: 590px;"></el-empty> 
        <el-table ref="tableRef" v-loading="loading" v-if="logList.length > 0 ? true : false" :data="logList" style="width: 100%" height="590" border :show-header="false" stripe>
          <el-table-column prop="log" label="label" width="width">
          </el-table-column>
        </el-table>
      </div>
    </el-card>
  </div>
</template>

<script>
export default {
  name:'TransferLog',
  data() {
    return {
      date:'',
      input:'',
      pickerOptions: {
        disabledDate(time) {
          return time.getTime() > Date.now();
        }
      },
      logList:[],
      description:'请选择日期',
      loading:false
    }
  },
  watch: {
    // 监听date,为null时修改description
    date(newVal) {
      if (!newVal) {
        this.description = '请选择日期'
      }
    },
  },
  methods: {
    // 获取传输日志,日期变化时触发
    async getTransferLog() {
      if (!this.date) {
        this.logList = []
        return
      }
      try {
        this.loading = true
        const result = await this.$API.sftpuser.reqSftpLog(this.date)
        if (result.code == 200){
          this.logList = result.data.sftplog === null ? [] : result.data.sftplog
          this.description = this.logList.length == 0 ? '暂无数据' : ''
        }
        // 核心：数据加载完成后，滚动条回到顶部
        this.$nextTick(() => {
          // 确保表格Ref存在，且表格已渲染
          if (this.$refs.tableRef) {
            // 找到表格的滚动容器
            const scrollWrapper = this.$refs.tableRef.$el.querySelector('.el-table__body-wrapper')
            if (scrollWrapper) {
              scrollWrapper.scrollTop = 0 // 重置滚动位置
            }
          }
        })
      } catch (error) {
        this.$message.error('获取传输日志失败')
      } finally {
        this.loading = false
      }
      
    },
    // 搜索
    search() {
      if (!this.input) {
        this.$message.error('请输入搜索内容')
        return
      }
      const newlist = this.logList.filter(item => item.log.includes(this.input))
      if (newlist.length == 0) {
        this.$message.error('未找到相关内容')
        return
      }
      this.logList = newlist
      // 搜索后重置滚动条
      this.$nextTick(() => {
        if (this.$refs.tableRef) {
          const scrollWrapper = this.$refs.tableRef.$el.querySelector('.el-table__body-wrapper')
          if (scrollWrapper) scrollWrapper.scrollTop = 0
        }
      })
    },
    // 重置搜索
    reset() {
      this.input = ''
      this.getTransferLog(this.date)
    }
  }
}
</script>

<style>

</style>