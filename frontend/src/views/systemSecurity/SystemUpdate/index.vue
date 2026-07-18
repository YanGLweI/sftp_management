<template>
  <div class="system-update-container">
    <!-- 系统更新标题 -->
    <div class="page-title">系统更新管理
      <div class="but-right">
        <el-tooltip effect="dark" content="更新计划" placement="top">
          <el-button circle icon="el-icon-alarm-clock" class="icon-btn" @click="handlePlanClick"></el-button>
        </el-tooltip>
      </div>
    </div>

    <!-- 状态卡片区域 -->
    <el-row :gutter="20" class="status-row">
      <el-col :span="12">
        <el-card class="status-card">
          <div class="status-label">系统更新状态</div>
          <div 
            class="status-value" 
            :class="{ 'status-error': status.systemUpdate.status === '异常' || status.systemUpdate.status === '获取中' }"
          >
            {{ status.systemUpdate.status }}
          </div>
          <div class="status-desc">{{ status.systemUpdate.desc }}</div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card class="status-card">
          <div class="status-label">更新任务状态</div>
          <div 
            class="status-value" 
            :class="{ 'status-error': status.updateTask.status === '异常' || status.updateTask.status === '获取中' }"
          >
            {{ status.updateTask.status }}
          </div>
          <div class="status-desc">{{ status.updateTask.desc }}</div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 命令执行控制区 -->
    <el-card class="control-card">
      <div style="display: flex; justify-content: space-between;">
        <el-button 
          type="primary" 
          :loading="isExecuting" 
          @click="startExecution"
          class="icon-exec-btn"
        >
          {{ isExecuting ? '执行中...' : '立即更新' }}
        </el-button>
        <div>
          <el-tooltip effect="dark" content="清空输出" placement="top">
            <el-button circle icon="el-icon-delete" class="icon-btn" @click="clearOutput" :disabled="!output"></el-button>
          </el-tooltip>
          <el-tooltip effect="dark" content="停止更新" placement="top">
            <el-button circle icon="el-icon-circle-close" class="icon-btn" :disabled="!isExecuting" @click="stopExecution"></el-button>
          </el-tooltip>
        </div>
      </div>
      <!-- 错误提示 -->
      <el-alert
        v-if="errorMessage"
        :title="errorMessage"
        type="error"
        show-icon
        closable
        @close="errorMessage = ''"
        style="margin-top: 20px;"
      />
      <el-collapse-transition>
        <!-- 命令输出区域 -->
        <div class="output-container" v-if="output">
          <div class="output-label">实时状态:</div>
          <pre ref="output" class="output-content">{{ output }}</pre>
        </div>
      </el-collapse-transition>
      
    </el-card>

    <!-- 更新历史表格 -->
    <div class="page-title" style="margin-top: 30px;">更新历史记录</div>
    <el-table
      :data="updateHistory"
      v-loading="tableLoading"
      border
      stripe
      style="width: 100%; margin-top: 20px;border-radius: 10px;"
      :header-cell-style="{ textAlign: 'center' }"
    >
      <el-table-column
        type="index"
        label="序号"
        align="center"
        width="80"
      />
      <el-table-column
        prop="update_time"
        label="更新时间"
        align="center"
        width="200"
      />
      <el-table-column  align="center" prop="hostname" label="主机名" width="120">
      </el-table-column>
      <el-table-column
        prop="status"
        label="更新状态"
        align="center"
        width="120"
      >
        <template slot-scope="{row}">
          <el-tag :type="row.status !== '失败' ? 'success' : 'danger'" effect="plain">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column
        prop="duration"
        label="耗时"
        align="center"
        width="100"
      >
        <template slot-scope="{row}">
          <el-tag type="primary" effect="plain">{{ formatDuration(row.duration) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column
        prop="update_brief"
        label="更新摘要"
        show-overflow-tooltip
        align="left"
      >
        <template slot-scope="{row}">
          <div style="display: flex; justify-content: space-between;">
            <span>{{ row.update_brief }}</span>
            <el-tooltip effect="dark" content="查看详情" placement="top">
              <el-button class="icon-btn-table" icon="el-icon-more-outline" @click="handleDetailClick(row.id)"></el-button>
            </el-tooltip>
          </div>
        </template>
      </el-table-column>
    </el-table>
    <!-- 分页  -->
    <el-pagination 
      @current-change="getUpdateHistory"
      @size-change="handleSizeChange"
      style="margin-top: 20px;text-align: center;" 
      :current-page="params.pageNum" 
      :page-sizes="[3, 5, 10]" 
      :page-size="params.pageSize" 
      layout="sizes, prev, pager, next, jumper,->,total" 
      :total="total">
    </el-pagination>

    <!-- 更新详情弹窗 -->
    <el-dialog
      title="更新详情"
      :visible.sync="dialogVisible"
      width="60%"
      :close-on-click-modal="false"
    >
      <el-card>
        <div class="detail-content">{{ updateDetail }}</div>
      </el-card>
    </el-dialog>

    <!-- 更新计划弹窗 -->
    <el-dialog
      title="系统更新计划"
      :visible.sync="dialogPlanVisible"
      width="480px"
      :close-on-click-modal="false"
      @close="closePlanDialog"
    >
      <!-- 弹窗内容 -->
      <el-tabs v-model="activeName" @tab-click="handleClick">
          <el-tab-pane label="更新任务" name="updatePlan">
            <el-card shadow="never">
              <el-form :model="form" :rules="rules" ref="formRef" label-width="140px" v-loading="formloading">
                <!-- 扫描间隔设置 -->
                <el-form-item label="计划" prop="ruleType">
                  <!-- <el-input v-model="form.ruleType" placeholder="请输入计划" /> -->
                  <el-select v-model="form.ruleType" placeholder="请选择任务启动计划" @change="handleRuleTypeChange">
                    <el-option label="每天" value="Daily"></el-option>
                    <el-option label="每周" value="Weekly"></el-option>
                  </el-select>
                </el-form-item>
                <!-- 日期时间设置 -->
                <el-form-item label="开始于" prop="time">
                  <el-time-select
                    v-model="form.time"
                    :picker-options="{
                      start: '00:00',
                      step: '00:30',
                      end: '23:30'
                    }"
                    placeholder="选择时间">
                  </el-time-select>
                </el-form-item>
                <!-- 任务执行间隔设置（每周） -->
                <el-form-item label="运行间隔" prop="interval" v-if="form.ruleType === 'Weekly'">
                  <el-select v-model="form.interval" placeholder="请选择任务执行间隔">
                    <el-option label="周一" value="Mon"></el-option>
                    <el-option label="周二" value="Tue"></el-option>
                    <el-option label="周三" value="Wed"></el-option>
                    <el-option label="周四" value="Thu"></el-option>
                    <el-option label="周五" value="Fri"></el-option>
                    <el-option label="周六" value="Sat"></el-option>
                    <el-option label="周日" value="Sun"></el-option>
                  </el-select>
                </el-form-item>
                <!-- 是否启用计划任务 -->
                <el-form-item label="是否启用计划任务" prop="enable">
                  <el-radio-group v-model="form.enable">
                    <el-radio :label="true">是</el-radio>
                    <el-radio :label="false">否</el-radio>
                  </el-radio-group>
                </el-form-item>
              </el-form>
            </el-card>
          </el-tab-pane>
          <el-tab-pane label="报告任务" name="reportPlan">
            <el-card shadow="never">
              <el-form :model="form" :rules="rules" ref="formRef" label-width="140px" v-loading="formloading">
                <!-- 扫描间隔设置 -->
                <el-form-item label="计划" prop="ruleType">
                  <!-- <el-input v-model="form.ruleType" placeholder="请输入计划" /> -->
                  <el-select v-model="form.ruleType" placeholder="请选择任务启动计划" @change="handleRuleTypeChange">
                    <el-option label="每天" value="Daily"></el-option>
                    <el-option label="每周" value="Weekly"></el-option>
                  </el-select>
                </el-form-item>
                <!-- 日期时间设置 -->
                <el-form-item label="开始于" prop="time">
                  <el-time-select
                    v-model="form.time"
                    :picker-options="{
                      start: '00:00',
                      step: '00:30',
                      end: '23:30'
                    }"
                    placeholder="选择时间">
                  </el-time-select>
                </el-form-item>
                <!-- 任务执行间隔设置（每周） -->
                <el-form-item label="运行间隔" prop="interval" v-if="form.ruleType === 'Weekly'">
                  <el-select v-model="form.interval" placeholder="请选择任务执行间隔">
                    <el-option label="周一" value="Mon"></el-option>
                    <el-option label="周二" value="Tue"></el-option>
                    <el-option label="周三" value="Wed"></el-option>
                    <el-option label="周四" value="Thu"></el-option>
                    <el-option label="周五" value="Fri"></el-option>
                    <el-option label="周六" value="Sat"></el-option>
                    <el-option label="周日" value="Sun"></el-option>
                  </el-select>
                </el-form-item>
                <!-- 是否启用计划任务 -->
                <el-form-item label="是否启用计划任务" prop="enable">
                  <el-radio-group v-model="form.enable">
                    <el-radio :label="true">是</el-radio>
                    <el-radio :label="false">否</el-radio>
                  </el-radio-group>
                </el-form-item>
              </el-form>
            </el-card>
          </el-tab-pane>
      </el-tabs>
      <!-- 弹窗底部：确认/取消按钮 -->
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogPlanVisible = false">取 消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="butloading">确 定</el-button>
      </div>
    </el-dialog>
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
      ws: null,
      tableLoading: false,

      updateHistory: [],
      latestUpdate:{},
      params: {
        pageNum: 1,
        pageSize: 5
      },
      total: 0,
      dialogVisible: false,
      updateDetail: '',
      dialogPlanVisible: false,
      form: {
        ruleType: '',
        time: '',
        interval: '',
        enable: '',
      },
      reqForm:{
        cron:'',
        enable:''
      },
      rules: {
        ruleType: [{ required: true, message: '请选择计划', trigger: 'blur' }],
        time: [{ required: true, message: '请选择时间', trigger: 'blur' }],
        interval: [{ required: true, message: '请选择任务执行间隔', trigger: 'blur' }],
        enable: [{ required: true, message: '请选择是否启用计划任务', trigger: 'change' }],
      },
      formloading: false,
      butloading: false,
      activeName: 'updatePlan',
    }
  },
  created() {
    // 初始化获取更新历史
    this.getUpdateHistory();
  },
  computed: {
    // 计算更新状态信息
    status() {
      return {
        // 系统更新状态（模拟：超过3天未更新则异常）
        systemUpdate: this.getSystemUpdateStatus(),
        // 更新任务状态（模拟：根据当前执行状态和历史记录判断）
        updateTask: this.getUpdateTaskStatus()
      };
    },
  },
  methods: {
    // 计算系统更新状态
    getSystemUpdateStatus() {
      // 列表中第一个的状态为成功或暂无可用更新时，系统更新状态为正常
      const latestUpdate = this.latestUpdate || {};
      
      if (!latestUpdate) {
        return { status: "异常", desc: "未检测到历史更新记录" };
      }

      // 最新记录状态不是成功/暂无可用更新 → 直接判定异常
      if (!['成功', '暂无可用更新'].includes(latestUpdate.status)) {
        return { 
          status: "异常", 
          desc: `最新更新任务失败（${latestUpdate.update_time}）` 
        };
      }
      
      try {
        const lastUpdateTime = new Date(latestUpdate.update_time);
        const now = new Date();
        const timeDiff = now - lastUpdateTime;
        const dayDiff = Math.floor(timeDiff / (1000 * 60 * 60 * 24));
        
        if (dayDiff > 7) {
          return { 
            status: "异常", 
            desc: `系统已${dayDiff}天未更新（超过7天）` 
          };
        } else {
          return { 
            status: "正常", 
            desc: `系统${dayDiff}天前更新（≤7天）` 
          };
        }
      } catch (e) {
        return { status: "异常", desc: "更新时间解析失败" };
      }
    },
    
    // 计算更新任务状态
    getUpdateTaskStatus() {
      // 执行中状态优先
      if (this.isExecuting) {
        return { status: "正常", desc: "更新任务正在执行中" };
      }
      
      // 检查最后一次更新结果
      const lastTask = this.latestUpdate || {};
      if (!lastTask) {
        return { status: "获取中", desc: "未获取到更新任务记录" };
      }
      
      if (lastTask.status === '失败') {
        return { status: "异常", desc: `最后一次更新失败：${lastTask.update_time}` };
      } else {
        return { status: "正常", desc: `最后一次更新成功（${lastTask.update_time}）` };
      }
    },
    
    // 启动更新执行
    startExecution() {
      // 重置状态
      this.isExecuting = true;
      this.output = '';
      this.errorMessage = '';
      
      try {
        // 确定WebSocket协议
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/dev-api/system/ws/update`;

        // 创建WebSocket连接
        this.ws = new WebSocket(wsUrl);
        
        // 连接建立
        this.ws.onopen = () => {
          this.output += '✅ 连接到更新服务器...\n';
          this.output += '🔄 开始执行系统更新流程...\n';
        };
        
        // 接收消息
        this.ws.onmessage = (event) => {
          this.output += event.data + '\n';
          // 自动滚动到底部
          this.$nextTick(() => {
            this.scrollToBottom();
          });
        };
        
        // 连接错误
        this.ws.onerror = (error) => {
          console.error('WebSocket error:', error);
          this.errorMessage = 'WebSocket连接发生错误';
          this.isExecuting = false;
        };
        
        // 连接关闭
        this.ws.onclose = () => {
          this.output += '\n🔌 连接已关闭，更新流程结束';
          this.isExecuting = false;
          this.getUpdateHistory();
        };
      } catch (error) {
        this.errorMessage = '创建WebSocket连接失败: ' + error.message;
        this.isExecuting = false;
      }
    },
    
    // 停止更新
    stopExecution() {
      if (this.ws) {
        this.ws.close();
        this.output += '\n🛑 用户手动终止更新流程\n';
      }
      this.isExecuting = false;
    },
    
    // 清空输出
    clearOutput() {
      this.output = '';
    },
    
    // 滚动到底部
    scrollToBottom() {
      const outputEl = this.$refs.output;
      if (outputEl) {
        outputEl.scrollTop = outputEl.scrollHeight;
      }
    },
    // 获取更新历史
    async getUpdateHistory(pages = 1){
      // 设置分页参数
      this.params.pageNum = pages;
      
      try {
        this.tableLoading = true;
        const res = await this.$API.system.reqUpdateHistory(this.params);
        if (res.code === 200) {
          this.updateHistory = res.data.updateHistory || [];
          this.total = res.data.total || 0;
          // 如果是第一页，记录最新更新记录
          if (pages === 1) {
            this.latestUpdate = res.data.updateHistory[0] || {};
          }
        } else {
          this.$message.error(res.msg || '获取更新历史失败');
        }
      } catch (error) {
        this.$message.error('获取更新历史失败: ' + error.message);
      } finally {
        this.tableLoading = false;
      }
    },
    // 处理页大小变化
    handleSizeChange(val) {
      this.params.pageSize = val;
      this.getUpdateHistory();
    },
    // 处理详情点击
    async handleDetailClick(id){
      // 清空详情内容
      this.updateDetail = '';
      try {
        const res = await this.$API.system.reqUpdateDetail(id);
        if (res.code === 200) {
          this.updateDetail = res.data.UpdateDetails || '';
          this.dialogVisible = true;
        } else {
          this.$message.error(res.msg || '获取更新详情失败');
        }
      } catch (error) {
        this.$message.error('获取更新详情失败: ' + error.message);
      }
    },
    // 计算耗时，超过60秒，显示1m,例如62秒，显示1m2s
    formatDuration(seconds) {
      if (seconds <= 60) {
        return `${seconds}s`;
      }
      const minutes = Math.floor(seconds / 60);
      const remainingSeconds = seconds % 60;
      return `${minutes}m${remainingSeconds}s`;
    },
    // 处理规则类型变化
    handleRuleTypeChange(val) {
      // 清理时间和间隔
      this.form.time = '00:00';
      this.form.interval = '';
      this.form.enable = '';
    },
    // 关闭计划弹窗时，重置表单
    resetForm(){
      this.$refs.formRef.resetFields();
      this.$refs.formRef.clearValidate();
      this.form = {
        ruleType: '',
        time: '',
        interval: '',
        enable: '',
      };
      this.activeName = 'updatePlan';
    },
    closePlanDialog() {
      this.resetForm();
      this.dialogPlanVisible = false;
    },
    // 提交计划任务表单
    handleSubmit(){
      this.$refs.formRef.validate(async (valid) => {
        if (valid) {
          try {
            this.butloading = true;
            // 转换表达式
            const cronExpression = this.convertToCron(this.form);
            // console.log('转换后的Cron表达式:', cronExpression);
            this.reqForm.cron = cronExpression;
            this.reqForm.enable = this.form.enable;
            // 提交计划任务
            let res
            if(this.activeName === 'updatePlan'){
              res = await this.$API.system.reqSetUpdateSchedule(this.reqForm);
            } else if(this.activeName === 'reportPlan'){
              res = await this.$API.system.reqSetUpdateReportSchedule(this.reqForm);
            }
            if (res.code === 200) {
              this.$message.success('计划任务更新成功');
              this.closePlanDialog();
            } 
            this.closePlanDialog();
          } catch (error) {
            this.$message.error('计划任务提交失败: ' + error.message);
          } finally {
            this.butloading = false;
          }
          
        } else {
          this.$message.error('请填写完整计划任务信息');
        }
      });
    },

    // 将Daily/Weekly规则转换为Linux Cron表达式
    convertToCron(rule) {
      // 1. 基础校验
      if (!rule || !rule.ruleType || !rule.time) {
          console.error('缺少必要参数：ruleType或time');
          return '';
      }

      // 2. 解析时间（HH:MM → 小时、分钟）
      const timeParts = rule.time.split(':');
      if (timeParts.length !== 2) {
          console.error('时间格式错误，需为HH:MM（如00:00）');
          return '';
      }
      const [hour, minute] = timeParts.map(Number);
      // 校验时分合法性
      if (isNaN(hour) || isNaN(minute) || hour < 0 || hour > 23 || minute < 0 || minute > 59) {
          console.error('时间值超出范围：小时0-23，分钟0-59');
          return '';
      }

      // 3. 星期映射（Cron标准：0/7=周日，1=周一，2=周二...6=周六）
      const weekMap = {Sun: 0,Mon: 1,Tue: 2,Wed: 3,Thu: 4,Fri: 5,Sat: 6};

      // 4. 按规则类型生成Cron
      switch (rule.ruleType) {
        case 'Daily':
            // 每天 时分 * * *
            return `${minute} ${hour} * * *`;
        
        case 'Weekly':
            // 校验星期参数
            if (!rule.interval) {
                console.error('Weekly规则缺少interval（星期）参数');
                return '';
            }
            const weekKey = rule.interval.trim();
            const weekValue = weekMap[weekKey];
            if (weekValue === undefined) {
                console.error('星期格式错误，支持：Sun/Mon/Tue/Wed/Thu/Fri/Sat');
                return '';
            }
            // 每周指定星期 时分 * * 星期值
            return `${minute} ${hour} * * ${weekValue}`;
          
        default:
            console.error(`不支持的规则类型：${rule.ruleType}，仅支持Daily/Weekly`);
            return '';
      }
    },
    handlePlanClick(){
      this.dialogPlanVisible = true;
      this.getUpdateSchedule();
    },

    // 获取更新计划
    async getUpdateSchedule(){
      try {
        this.dialogPlanVisible = true;
        this.formloading = true;
        const res = await this.$API.system.reqUpdateSchedule();
        if (res.code === 200) {
          this.form = res.data || {};
          
        } else {
          this.$message.error(res.msg || '获取更新计划失败');
        }
      } catch (error) {
        this.$message.error('获取更新计划失败: ' + error.message);
      } finally {
        this.formloading = false;
      }
    },

    // 获取更新报告任务计划
    async getUpdateReportSchedule(){
      try {
        this.formloading = true;
        const res = await this.$API.system.reqUpdateReportSchedule();
        if (res.code === 200) {
          this.form = res.data || {};
          
        } else {
          this.$message.error(res.msg || '获取更新报告任务计划失败');
        }
      } catch (error) {
        this.$message.error('获取更新报告任务计划失败: ' + error.message);
      } finally {
        this.formloading = false;
      }
    },

    // 点击tab切换时，获取不同的任务计划
    handleClick(tab){
      if(tab.name === 'updatePlan'){
        this.getUpdateSchedule();
      } else if(tab.name === 'reportPlan'){
        this.getUpdateReportSchedule();
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
/* 整体容器样式，和目标页面保持一致 */
.system-update-container {
  padding: 24px;
  border-radius: 12px;
  background: #f5f7fa;
  border: 1px solid rgba(100, 200, 150, 0.15);
  box-shadow: 
    0 4px 20px rgba(100, 200, 150, 0.08),
    0 8px 30px rgba(0, 0, 0, 0.05);
}

/* 页面标题样式复用 */
.page-title {
  text-align: center;
  font-size: 18px;
  font-weight: 600;
  color: #2a3b47;
  margin-bottom: 24px;
  position: relative;
  /* padding-right: 60px; */
  letter-spacing: 0.3px;
}

/* 右侧按钮组 */
.but-right {
  position: absolute;
  top: 50%;
  right: 0;
  transform: translateY(-50%);
  display: flex;
  gap: 16px;
}

/* 图标按钮样式，和目标页面一致 */
.icon-btn {
  background: transparent !important;
  border: none !important;
  padding: 8px !important;
  width: auto !important;
  height: auto !important;
  font-size: 30px;
  color: #94a3b8;
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
  
  --el-button-hover-bg-color: transparent !important;
  --el-button-active-bg-color: transparent !important;
  --el-button-hover-border-color: transparent !important;
  --el-button-active-border-color: transparent !important;
}

.icon-btn:hover {
  color: #64c896 !important;
  transform: scale(1.1);
  text-shadow: 0 0 6px rgba(100, 200, 150, 0.2);
}

.icon-btn:disabled {
  color: #cbd5e1 !important;
  transform: none !important;
  cursor: not-allowed;
}

.icon-btn-table {
  background: transparent !important;
  border: none !important;
  padding: 0px !important;
  width: auto !important;
  height: auto !important;
  font-size: 24px;
  color: #94a3b8;
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
  
  --el-button-hover-bg-color: transparent !important;
  --el-button-active-bg-color: transparent !important;
  --el-button-hover-border-color: transparent !important;
  --el-button-active-border-color: transparent !important;
}

.icon-btn-table:hover {
  color: #64c896 !important;
  transform: scale(1.1);
  text-shadow: 0 0 6px rgba(100, 200, 150, 0.2);
}

.icon-btn-table:disabled {
  color: #cbd5e1 !important;
  transform: none !important;
  cursor: not-allowed;
}

/* 状态卡片行间距 */
.status-row {
  margin-bottom: 24px;
}

/* 状态卡片样式，完全复用目标页面 */
.status-card {
  text-align: center;
  padding: 24px 16px;
  border-radius: 10px;
  border-left: 6px solid #64c896 !important;
  background: linear-gradient(90deg, #f9fffb, #f0f9ff);
  border: 1px solid rgba(100, 200, 150, 0.2);
  box-shadow: 0 2px 12px rgba(100, 200, 150, 0.08);
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
  min-height: 140px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

/* 异常状态卡片样式 */
.status-card:has(.status-error) {
  border-left-color: #e57373 !important;
  background: linear-gradient(90deg, #fff8f8, #fef0f0);
  border-color: rgba(229, 115, 115, 0.2);
}

.status-error {
  color: #e57373 !important;
  font-weight: 600;
}

.status-label {
  font-size: 16px;
  color: #475569;
  margin-bottom: 12px;
  letter-spacing: 0.2px;
  height: 22px;
  line-height: 22px;
}

.status-value {
  font-size: 18px;
  color: #64c896;
  font-weight: 600;
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
  line-height: 24px;
  display: inline-block;
  transform-origin: center;
  height: 24px;
  margin: 0 auto;
}

.status-card:hover .status-value {
  transform: scale(1.2);
  text-shadow: 0 0 8px rgba(100, 200, 150, 0.2);
}

.status-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(100, 200, 150, 0.12);
}

.status-desc {
  font-size: 14px;
  color: #94a3b8;
  margin-top: 8px;
  line-height: 1.6;
  min-height: 22px;
}

/* 控制卡片样式 */
.control-card {
  border-radius: 10px;
  border: 1px solid rgba(100, 200, 150, 0.2);
  border-left: 6px solid #64c896 !important;
  /* background: #ffffff; */
  background: linear-gradient(90deg, #f9fffb, #f0f9ff);
  box-shadow: 0 2px 12px rgba(100, 200, 150, 0.08);
  padding: 24px;
  margin-bottom: 24px;
}

/* 立即更新按钮样式（对齐目标页面确认按钮风格） */
.icon-exec-btn {
  background: transparent !important;
  border: 1px solid #64c896 !important;
  color: #64c896 !important;
  padding: 8px 24px !important;
  font-size: 14px !important;
  border-radius: 8px !important;
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1) !important;
  
  /* 覆盖Element默认样式 */
  --el-button-primary-bg-color: transparent !important;
  --el-button-primary-border-color: #64c896 !important;
  --el-button-primary-text-color: #64c896 !important;
  --el-button-primary-hover-bg-color: rgba(100, 200, 150, 0.1) !important;
  --el-button-primary-hover-border-color: #55b785 !important;
  --el-button-primary-hover-text-color: #55b785 !important;
  --el-button-primary-active-bg-color: rgba(100, 200, 150, 0.2) !important;
  --el-button-primary-active-border-color: #4aa876 !important;
  --el-button-primary-active-text-color: #4aa876 !important;
  --el-button-loading-text-color: #64c896 !important;
}

.icon-exec-btn:hover {
  transform: scale(1.05);
  box-shadow: 0 4px 12px rgba(100, 200, 150, 0.15);
}

.icon-exec-btn:disabled,
.icon-exec-btn.is-loading {
  transform: none !important;
  box-shadow: none !important;
  opacity: 0.8;
}

/* 输出容器样式 */
.output-container {
  margin-top: 20px;
  border-radius: 10px;
  background: #f9fffb;
  border: 1px solid rgba(100, 200, 150, 0.15);
  padding: 16px;
}

.output-label {
  font-size: 16px;
  color: #475569;
  margin-bottom: 12px;
  font-weight: 500;
}

.output-content {
  height: 300px;
  overflow-y: auto;
  background-color: #ffffff;
  padding: 16px;
  border-radius: 8px;
  font-family: "Consolas", "Monaco", monospace;
  white-space: pre-wrap;
  line-height: 1.6;
  color: #2a3b47;
  font-size: 12px;
  border: 1px solid rgba(100, 200, 150, 0.1);
}

/* 表格样式适配 */
:deep(.el-table) {
  --el-table-header-text-color: #2a3b47;
  --el-table-row-hover-bg-color: #f9fffb;
  --el-table-border-color: rgba(100, 200, 150, 0.2);
}

:deep(.el-table th) {
  background-color: #f0f9ff !important;
}

:deep(.el-table--striped .el-table__row--striped td) {
  background-color: #f9fffb !important;
}
.detail-content{
  /* 关键样式：解析\n为换行，同时自动合并多余空格 */
  white-space: pre-line;
  /* 可选：增加行高，提升可读性 */
  line-height: 1.6;
  height: 500px;
  /* 内容超出部分滚动 */
  overflow-y: auto;
  /* 内容超出部分滚动 */
  overflow-x: auto;
  /* 文字大小 */
  font-size: 10px;
}
</style>