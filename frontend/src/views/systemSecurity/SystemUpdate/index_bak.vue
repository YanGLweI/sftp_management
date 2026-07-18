<template>
  <div class="command-executor">
    <el-card>
      <div slot="header" class="clearfix">
        <span>命令执行器</span>
        <el-button 
          style="float: right; padding: 3px 0" 
          type="text" 
          @click="clearOutput">
          清空输出
        </el-button>
      </div>
      
      <el-button 
        type="primary" 
        :loading="isExecuting" 
        @click="startExecution">
        {{ isExecuting ? '执行中...' : '立即更新' }}
      </el-button>
      <el-button 
        type="danger" 
        circle 
        icon="el-icon-circle-close" 
        size="mini" 
        :disabled="!isExecuting"
        @click="stopExecution">
      </el-button>
      
      <el-alert
        v-if="errorMessage"
        :title="errorMessage"
        type="error"
        show-icon
        closable
        @close="errorMessage = ''"
        style="margin-top: 20px;"
      />
      
      <div class="output-container">
        <h3>命令输出:</h3>
        <pre ref="output" class="output-content">{{ output }}</pre>
      </div>
    </el-card>
  </div>
</template>

<script>
export default {
  name: 'SystemUpdate',
  data() {
    return {
      isExecuting: false,
      output: '',
      errorMessage: '',
      ws: null
    }
  },
  methods: {
    startExecution() {
      // 重置状态
      this.isExecuting = true;
      this.output = '';
      this.errorMessage = '';
      
      try {
        // 确定WebSocket协议，根据当前页面协议自动选择
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        // 使用相对路径，通过Vue代理转发
        const wsUrl = `${protocol}//${window.location.host}/dev-api/system/ws/update`;

        // 创建WebSocket连接
        this.ws = new WebSocket(wsUrl);
        
        // 连接建立时的处理
        this.ws.onopen = () => {
          this.output += '连接到服务器...\n';
        };
        
        // 处理接收到的消息
        this.ws.onmessage = (event) => {
          this.output += event.data + '\n';
          // 自动滚动到底部
          this.$nextTick(() => {
            this.scrollToBottom();
          });
        };
        
        // 处理错误
        this.ws.onerror = (error) => {
          console.error('WebSocket error:', error);
          this.errorMessage = 'WebSocket连接发生错误';
          this.isExecuting = false;
        };
        
        // 连接关闭时的处理
        this.ws.onclose = () => {
          this.output += '连接已关闭\n';
          this.isExecuting = false;
        };
      } catch (error) {
        this.errorMessage = '创建WebSocket连接失败: ' + error.message;
        this.isExecuting = false;
      }
    },
    stopExecution() {
      if (this.ws) {
        this.ws.close();
      }
      this.isExecuting = false;
    },
    clearOutput() {
      this.output = '';
    },
    scrollToBottom() {
      const outputEl = this.$refs.output;
      if (outputEl) {
        outputEl.scrollTop = outputEl.scrollHeight;
      }
    }
  },
  beforeDestroy() {
    // 组件销毁前关闭连接
    if (this.ws) {
      this.ws.close();
    }
  }
}
</script>

<style scoped>
.command-executor {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.output-container {
  margin-top: 20px;
  border: 1px solid #ebeef5;
  border-radius: 4px;
  padding: 15px;
}

.output-content {
  height: 400px;
  overflow-y: auto;
  background-color: #f5f7fa;
  padding: 10px;
  border-radius: 4px;
  font-family: monospace;
  white-space: pre-wrap;
}

.clearfix:before,
.clearfix:after {
  display: table;
  content: "";
}
.clearfix:after {
  clear: both;
}
</style>