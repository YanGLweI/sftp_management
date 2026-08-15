<template>
  <div>
    <el-card shadow="hover">
      <div slot="header">
        <div style="float: left;">
          <el-date-picker
            v-model="tempSearchObj.datetime"
            value-format="yyyy-MM-dd"
            align="right"
            type="date"
            placeholder="选择日期"
            :picker-options="pickerOptions"
            @change="search"
          >
          </el-date-picker>
        </div>
        <div style="float: right;">
          <!-- 搜索 -->
          <el-input placeholder="请输入用户名" style="width: 400px;" prefix-icon="el-icon-search" v-model="tempSearchObj.username" @keyup.enter.native="search">
            <el-button slot="append" icon="el-icon-search" @click="search"></el-button>
          </el-input>
          <!-- 重置搜索 -->
          <el-button icon="el-icon-refresh" circle style="margin-left: 10px;" size="mini" @click="resetSearch"></el-button>
        </div>
      </div>
      <div>
        <!-- 日志列表 -->
        <el-table :data="logList" style="width: 100%;margin: 20px 0;" border stripe>
          <el-table-column type="index" label="序号" width="80" align="center">
          </el-table-column>
          <el-table-column prop="created_at" label="时间" width="200">
          </el-table-column>
          <el-table-column prop="username" label="用户名" width="140">
          </el-table-column>
          <el-table-column prop="reviewer" label="复核人" width="140">
            <template slot-scope="{row}">
              {{ row.reviewer || '-' }}
            </template>
          </el-table-column>
          <el-table-column prop="ip" label="IP" width="140">
          </el-table-column>
          <el-table-column prop="action" label="动作" width="140" align="center">
            <template slot-scope="{row}">
              <el-tag :type="getStatusType(row.action)">{{ row.action }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="path" label="路径" show-overflow-tooltip>
          </el-table-column>
          <el-table-column prop="message" label="消息" width="400" show-overflow-tooltip>
          </el-table-column>
        </el-table>
        <!-- 分页器 -->
        <el-pagination
          @current-change="getLogList"
          @size-change="handleSizeChange"
          style="margin-top: 20px;text-align: center;"
          :current-page="page"
          :page-sizes="[3, 5, 10]"
          :page-size="limit"
          layout="sizes, prev, pager, next, jumper,->,total"
          :total="total">
        </el-pagination>
      </div>
    </el-card>
  </div>
</template>

<script>
export default {
  name: "SftpLogList",
  data() {
    return {
      logList: [],
      // 分页器
      page: 1,
      limit: 10,
      total: 0,
      // 日期选择
      tempSearchObj: {
        datetime: '',
        username: '',
      },
      pickerOptions: {
        disabledDate(time) {
          return time.getTime() > Date.now();
        }
      },
      // 时间搜索
      searchObj: {
        datetime: '',
        username: '',
      }
    }
  },
  mounted() {
    this.getLogList();
  },
  methods: {
    async getLogList(pages = 1) {
      // 修改参数
      this.page = pages;
      const { page, limit, searchObj } = this
      try {
        const result = await this.$API.sftpuser.reqSftpLogList(page, limit, searchObj)
        if (result.code == 200) {
          this.logList = result.data.logs
          this.total = result.data.total
        }
      } catch (error) {}
    },
    handleSizeChange(limit) {
      // 修改参数
      this.limit = limit;
      this.getLogList(this.page);
    },
    // 根据不同的 action 返回不同的tag类型
    getStatusType(action) {
      switch (action) {
        case "Login":
          return "primary";
        case "Logout":
          return "info";
        case "Upload":
        case "Mkdir":
          return "success";
        case "Delete":
        case "BatchDelete":
          return "danger";
        case "Rename":
        case "Download":
          return "warning";
        default:
          return "success"; // 如果有其他未知情况，可以返回默认值，或者根据实际需求处理
      }
    },
    // 搜索和日期变化的回调
    search() {
      this.searchObj = { ...this.tempSearchObj };
      this.getLogList();
    },
    // 重置搜索
    resetSearch() {
      this.searchObj = {
        datetime: '',
        username: ''
      };
      this.tempSearchObj = {
        datetime: '',
        username: '',
      }
      this.getLogList();
    },
  }
};
</script>

<style>
</style>
