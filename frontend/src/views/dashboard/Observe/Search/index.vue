<template>
  <div>
    <el-card>
      <div slot="header" class="header">
        <div class="search-header">
          <span>活跃用户 Top6</span>
        </div>
      </div>
      <div>
        <!-- 排行榜列表展示 -->
        <div class="rank-list">
          <div v-for="(user, index) in activeUsers" :key="index" class="rank-item">
            <span class="rank-index" :class="'rank-' + (index + 1)">
              {{ index + 1 }}
            </span>
            <span class="rank-name">{{ user.username }}</span>
            <span class="rank-value">{{ user.count }}</span>
          </div>
          <!-- 空状态提示 -->
          <div v-if="!activeUsers || activeUsers.length === 0" class="empty-tip">
            暂无数据
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script>
import { reqActiveUsersTop6 } from '@/api/dashboard/dashboard.js'

export default {
  name: "Search",
  data() {
    return {
      activeUsers: []  // 活跃用户列表
    }
  },
  mounted() {
    this.fetchActiveUsers()
  },
  methods: {
    async fetchActiveUsers() {
      try {
        const res = await reqActiveUsersTop6()
        if (res.code === 200) {
          this.activeUsers = res.data.map(item => ({
            username: item.username,
            count: item.count
          }))
        }
      } catch (error) {
        console.error('加载活跃用户失败:', error)
        this.activeUsers = []
      }
    }
  }
}
</script>

<style scoped>
.header {
  border-bottom: 1px solid #eee;
  padding: 5px 0;
}

.search-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

/* 新增排行榜样式 */
.rank-list {
  margin-top: 10px;
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: space-evenly;  /* 条目均匀分布，填满卡片高度 */
}

.rank-item {
  display: flex;
  align-items: center;
  padding: 10px;
  border-bottom: 1px solid #eee;
}

.rank-index {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  background: #ddd;
  color: #333;
  text-align: center;
  line-height: 30px;
  margin-right: 15px;
  font-weight: bold;
}

.rank-index.rank-1 { background: #ff4757; color: white; }
.rank-index.rank-2 { background: #ffa502; color: white; }
.rank-index.rank-3 { background: #2ed573; color: white; }

.rank-name {
  flex: 1;
  font-size: 14px;
}

.rank-value {
  font-size: 14px;
  color: #666;
}

.empty-tip {
  text-align: center;
  color: #999;
  padding: 30px 0;
  font-size: 14px;
}
</style>