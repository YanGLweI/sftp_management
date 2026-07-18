<template>
  <div class="antivirus-container">
    <!-- 系统信息标题 -->
    <div class="page-title">卡巴斯基 Endpoint Security for Linux 系统信息
      <div class="but-right">
        <el-tooltip effect="dark" content="立即扫描" placement="top">
          <el-button circle icon="el-icon-video-play" class="icon-btn" @click="startScan"></el-button>
        </el-tooltip>
        <el-tooltip effect="dark" content="计划" placement="top">
          <el-button circle icon="el-icon-setting" class="icon-btn" @click="openScheduleDialog"></el-button> 
        </el-tooltip>
      </div>
      
    </div>
    

    <!-- 状态卡片区域 -->
    <el-row :gutter="20" class="status-row">
      <el-col :span="12">
        <el-card class="status-card">
          <div class="status-label">反病毒库更新状态</div>
          <!-- 动态绑定状态值和样式 -->
          <div 
            class="status-value" 
            :class="{ 'status-error': status.antivirusDb.status === '异常' || status.antivirusDb.status === '获取中' }"
          >
            {{ status.antivirusDb.status }}
          </div>
          <div class="status-desc">{{ status.antivirusDb.desc }}</div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card class="status-card">
          <div class="status-label">扫描任务状态</div>
          <!-- 动态绑定扫描任务状态和样式 -->
          <div 
            class="status-value" 
            :class="{ 'status-error': status.scanTask.status === '异常' || status.scanTask.status === '获取中' }"
          >
            {{ status.scanTask.status }}
          </div>
          <div class="status-desc">{{ status.scanTask.desc }}</div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 信息表格 -->
    <el-table
      :data="systemInfo"
      v-loading="loading"
      border
      stripe
      style="width: 100%; margin-top: 20px;border-radius: 10px;"
      :header-cell-style="{ 
        textAlign: 'center'
      }"
    >
      <el-table-column
        prop="item"
        label="项目"
        align="left"
        sortable
      />
      <el-table-column
        prop="content"
        label="内容"
        align="left"
        sortable
      />
    </el-table>

    <!-- 威胁报告区域 -->
    <div class="page-title" style="margin-top: 30px;">卡巴斯基 Endpoint Security for Linux 威胁报告</div>
    <el-card class="threat-card" style="margin-top: 10px;">
      <div class="status-label">隔离区状态</div>
      <div class="status-value">{{ isolationZoneStatus }}</div>
    </el-card>

    <!-- 设置按钮的弹窗 -->
    <el-dialog
      title="计划"
      :visible.sync="dialogVisible"
      width="480px"
      :close-on-click-modal="false"
      @close="handleClose"
    >
      <!-- 弹窗内容 -->
      <el-tabs v-model="activeName" @tab-click="handleClick">
        <el-tab-pane label="病毒扫描任务" name="scanPlan">
          <el-card shadow="never">
            <el-form :model="form" :rules="rules" ref="formRef" label-width="140px" v-loading="formloading">
              <!-- 扫描间隔设置 -->
              <el-form-item label="计划" prop="ruleType">
                <!-- <el-input v-model="form.ruleType" placeholder="请输入计划" /> -->
                <el-select v-model="form.ruleType" placeholder="请选择任务启动计划" @change="handleRuleTypeChange">
                  <el-option label="仅运行一次" value="Once"></el-option>
                  <el-option label="每分钟" value="Minutely"></el-option>
                  <el-option label="每小时" value="Hourly"></el-option>
                  <el-option label="每天" value="Daily"></el-option>
                  <el-option label="每周" value="Weekly"></el-option>
                  <el-option label="每月" value="Monthly"></el-option>
                </el-select>
              </el-form-item>
              <!-- 日期时间设置 -->
              <el-form-item label="日期" prop="date" 
              v-if="form.ruleType === 'once' || form.ruleType === 'Hourly'">
                <el-date-picker
                  v-model="form.date"
                  type="date"
                  value-format="yyyy/MM/dd"
                  placeholder="选择时间">
                </el-date-picker>
              </el-form-item>
              <el-form-item label="开始于" prop="time">
                <el-time-picker
                  v-model="form.time"
                  value-format="HH:mm:ss"
                  placeholder="选择时间">
                </el-time-picker>
              </el-form-item>
              <!-- 任务执行间隔设置（分钟、小时） -->
              <el-form-item 
              :label="form.ruleType === 'Minutely' ? '运行间隔(分钟)' : '运行间隔(小时)'" 
              v-if="form.ruleType === 'Minutely' || form.ruleType === 'Hourly'"
              prop="intervalTmp">
                <el-input-number controls-position="right" :min="1" :max="999" 
                v-model="form.intervalTmp" placeholder="请输入任务执行间隔" />
              </el-form-item>
              <!-- 任务执行间隔设置（每天） -->
              <el-form-item label="运行间隔(天)" prop="intervalTmp" v-if="form.ruleType === 'Daily'">
                <el-input-number controls-position="right" :min="1" :max="365" 
                v-model="form.intervalTmp" placeholder="请输入任务执行间隔" />
              </el-form-item>
              <!-- 任务执行间隔设置（每月） -->
              <el-form-item label="运行间隔(月)" prop="intervalTmp" v-if="form.ruleType === 'Monthly'">
                <el-input-number controls-position="right" :min="1" :max="12"
                v-model="form.intervalTmp" placeholder="请输入任务执行间隔" />
              </el-form-item>
              <!-- 任务执行间隔设置（每周） -->
              <el-form-item label="运行间隔(周)" prop="intervalTmp" v-if="form.ruleType === 'Weekly'">
                <el-select v-model="form.intervalTmp" placeholder="请选择任务执行间隔">
                  <el-option label="周一" value="Mon"></el-option>
                  <el-option label="周二" value="Tue"></el-option>
                  <el-option label="周三" value="Wed"></el-option>
                  <el-option label="周四" value="Thu"></el-option>
                  <el-option label="周五" value="Fri"></el-option>
                  <el-option label="周六" value="Sat"></el-option>
                  <el-option label="周日" value="Sun"></el-option>
                </el-select>
              </el-form-item>
              <!-- 是否运行错过的任务计划 -->
              <el-form-item label="运行错过的任务" prop="runMissedStartRules">
                <el-radio-group v-model="form.runMissedStartRules">
                  <el-radio label="Yes">是</el-radio>
                  <el-radio label="No">否</el-radio>
                </el-radio-group>
              </el-form-item>
            </el-form>
          </el-card>
        </el-tab-pane>
        <el-tab-pane label="保护状态报告" name="reportPlan">
          <el-card shadow="never">
            <el-form :model="reportForm" :rules="rules" ref="formRef" label-width="140px" v-loading="formloading">
              <!-- 扫描间隔设置 -->
              <el-form-item label="计划" prop="ruleType">
                <!-- <el-input v-model="form.ruleType" placeholder="请输入计划" /> -->
                <el-select v-model="reportForm.ruleType" placeholder="请选择任务启动计划" @change="handleRuleTypeChange">
                  <el-option label="每天" value="Daily"></el-option>
                  <el-option label="每周" value="Weekly"></el-option>
                </el-select>
              </el-form-item>
              <!-- 日期时间设置 -->
              <el-form-item label="开始于" prop="time">
                <el-time-select
                  v-model="reportForm.time"
                  :picker-options="{
                    start: '00:00',
                    step: '00:30',
                    end: '23:30'
                  }"
                  placeholder="选择时间">
                </el-time-select>
              </el-form-item>
              <!-- 任务执行间隔设置（每周） -->
              <el-form-item label="运行间隔" prop="interval" v-if="reportForm.ruleType === 'Weekly'">
                <el-select v-model="reportForm.interval" placeholder="请选择任务执行间隔">
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
                <el-radio-group v-model="reportForm.enable">
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
        <el-button @click="dialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="butloading">确 定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
export default {
  name: "Antivirus",
  data() {
    return {
      // 对应图中的系统信息数据
      systemInfo: [],
      // 弹窗是否可见
      dialogVisible: false,
      // loading状态
      loading: false,
      // 表单loading状态
      formloading: false,
      // 按钮loading状态
      butloading: false,
      // 表单数据模型
      form: {
        ruleType: '', //任务启动计划:once Monthly Weekly Daily Hourly Minutely
        runMissedStartRules: '' , //是否运行错过的任务计划 Yes No
        date: '', //日期 
        time: '', //时间 
        intervalTmp: '', // 临时任务执行间隔（分钟、小时、天、周、月）
        interval: '', // 任务执行间隔（分钟、小时、天、周、月）
        startTime: '', //任务启动时间 [<年>/<月>/<日>] [时]:[分]:[秒]；[<一个月中的第几天>|<一周中的第几天>]；[<启动周期>]。
        randomInterval: "0"
      },
      reportForm: {
        ruleType: '',
        time: '',
        interval: '',
        enable: '',
      },
      reqForm:{
        cron:'',
        enable:''
      },
      // 表单验证规则
      rules: {
        ruleType: [
          { required: true, message: '请选择任务启动计划', trigger: 'change' }
        ],
        date: [
          { required: true, message: '请选择日期', trigger: 'change' }
        ],
        time: [
          { required: true, message: '请选择开始时间', trigger: 'change' }
        ],
        intervalTmp: [
          { required: true, message: '请选择任务间隔', trigger: 'blur' }
        ],
        runMissedStartRules: [
          { required: true, message: '请选择是否运行错过的任务计划', trigger: 'change' }
        ],
        interval: [{ required: true, message: '请选择任务执行间隔', trigger: 'blur' }],
        enable: [{ required: true, message: '请选择是否启用计划任务', trigger: 'change' }],
      },
      isolationZoneStatus: '',
      activeName: 'scanPlan'
    };
  },
  mounted() {
    this.getAntivirusInfo();
    this.getIsolationZone();
  },
  // 计算属性：处理状态判断
  computed: {
    status() {
      return {
        // 反病毒库更新状态（数据库发布时间超过7天则异常）
        antivirusDb: this.getAntivirusDbStatus(),
        // 扫描任务状态（从未运行则异常）
        scanTask: this.getScanTaskStatus()
      };
    }
  },
  methods: {
    // 调用获取卡巴斯基信息的接口
    async getAntivirusInfo() {
      try {
        this.loading = true;
        const res = await this.$API.system.reqKasperskyInfo();
        if (res.code === 200) {
          // 处理返回的系统信息数据
          this.systemInfo = res.data;
        }
      } catch (error) {
        console.error("获取卡巴斯基信息失败：", error);
        // 接口异常时给默认值，避免页面空白
        this.systemInfo = [];
      } finally {
        this.loading = false;
      }
    },

    // 计算反病毒库更新状态
    getAntivirusDbStatus() {
      // 1. 找到数据库上次发布日期的条目
      const dbItem = this.systemInfo.find(item => item.item === "数据库的上次发布日期");
      
      // 未找到数据时默认异常
      if (!dbItem || !dbItem.content) {
        return { status: "获取中", desc: "未获取到数据库发布时间" };
      }

      try {
        // 2. 解析时间字符串（兼容 "2025-11-20 13:58:00" 格式）
        const dbTimeStr = dbItem.content;
        // 替换空格为T，适配Date.parse
        const parseTimeStr = dbTimeStr.replace(/\s+/, 'T');
        const dbTime = new Date(parseTimeStr);
        const now = new Date();

        // 3. 计算时间差（毫秒转天数，取整）
        const timeDiff = now - dbTime;
        const dayDiff = Math.floor(timeDiff / (1000 * 60 * 60 * 24));

        // 4. 判断是否超过7天
        if (dayDiff > 7) {
          return { 
            status: "异常", 
            desc: `数据库已${dayDiff}天未更新（超过7天）` 
          };
        } else {
          if (dayDiff === 0) {
            return { 
              status: "正常", 
              desc: `反病毒数据库版本：${dbTimeStr}（≤7天）` 
            };
          }
          return { 
            status: "正常", 
            desc: `数据库${dayDiff}天前更新（≤7天）` 
          };
        }
      } catch (e) {
        console.error("解析数据库时间失败：", e);
        return { status: "异常", desc: "时间格式解析失败" };
      }
    },

    // 计算扫描任务状态
    getScanTaskStatus() {
      // 1. 找到扫描任务上次运行日期的条目
      const scanItem = this.systemInfo.find(item => item.item === "Scan_My_Computer 任务的上次运行日期");
      
      // 未找到数据时默认“获取中”
      if (!scanItem || !scanItem.content) {
        return { status: "获取中", desc: "未获取到扫描任务状态" };
      }

      // 2. 判断是否从未运行
      if (scanItem.content === "从未运行") {
        return { status: "异常", desc: "扫描任务从未运行" };
      }

      try {
        // 3. 解析扫描任务上次运行时间
        const lastRunTimeStr = scanItem.content.replace(/\s+/, 'T');
        const lastRunTime = new Date(lastRunTimeStr);

        // 校验时间解析是否有效（Invalid Date 判定）
        if (isNaN(lastRunTime.getTime())) {
          return { status: "异常", desc: "扫描任务时间格式解析失败" };
        }

        // 4. 构造「当天中午12点」的时间对象（关键修正点）
        const today = new Date();
        const todayNoon = new Date(
          today.getFullYear(),
          today.getMonth(),
          today.getDate(),
          12, // 中午12点
          0, 0, 0 // 分、秒、毫秒置0
        );
        // console.log("今日：", today);

        // 5. 计算「上次运行时间」与「当天12点」的时间差（小时）
        // 时间差 = 绝对值(12点时间 - 上次运行时间) / 3600000（毫秒转小时）
        // const timeDiffMs = Math.abs(todayNoon - lastRunTime);
        const timeDiffMs = Math.abs(today - lastRunTime);
        const diffHours = timeDiffMs / (1000 * 60 * 60);
        // console.log("111111111", timeDiffMs);

        // 6. 判断是否超过12小时（核心逻辑）
        if (diffHours > 24) {
          return { 
            status: "异常", 
            desc: `上次运行时间在${diffHours.toFixed(2)}小时前（超过24小时）` 
          };
        } else {
          return { 
            status: "正常", 
            desc: `${diffHours.toFixed(2)}小时前运行了扫描任务（≤24小时）` 
          };
        }
      } catch (e) {
        console.error("扫描任务时间计算失败：", e);
        return { status: "异常", desc: "扫描任务时间计算异常" };
      }
    },
    // 打开设置计划弹窗
    openScheduleDialog() {
      this.dialogVisible = true;
      // 请求获取当前计划
      this.getAntivirusSchedule();
      this.getReportSchedule();
    },
    // 处理计划类型变化
    handleRuleTypeChange() {
      // 重置时间选择器
      this.form.time = '';
      this.form.intervalTmp = '';
      this.form.date = '';
    },
    // 重置表单
    resetForm() {
      this.$refs.formRef.resetFields();
      this.$refs.formRef.clearValidate();
      this.form = {
        ruleType: '', //任务启动计划:once Monthly Weekly Daily Hourly Minutely
        runMissedStartRules: '' , //是否运行错过的任务计划 yes no
        date: '', //日期 
        time: '', //时间 
        intervalTmp: '', // 任务执行间隔（分钟、小时、天、周、月）
        interval: '', // 任务执行间隔（分钟、小时、天、周、月）
        startTime: '', //任务启动时间 [<年>/<月>/<日>] [时]:[分]:[秒]；[<一个月中的第几天>|<一周中的第几天>]；[<启动周期>]。
        randomInterval: "0"
      },
      this.reportForm = {
        ruleType: '',
        time: '',
        interval: '',
        enable: '',
      };
      this.activeName = 'scanPlan';
    },
    // 关闭弹窗时重置表单
    handleClose() {
      this.resetForm();
      // // 重置验证
      // this.$nextTick(() => {
      //   this.$refs.formRef.resetFields();
      // });
      // this.$refs.formRef.clearValidate();
      this.dialogVisible = false;
    },
    // 获取当前的反病毒计划设置
    async getAntivirusSchedule() {
      try {
        this.formloading = true;
        const res = await this.$API.system.reqKasperskySchedule();
        if (res.code === 200) {
          const scheduleData = res.data;
          this.form.ruleType = scheduleData.ruleType;
          this.form.runMissedStartRules = scheduleData.runMissedStartRules;
          this.form.date = scheduleData.date;
          this.form.time = scheduleData.time;
          this.form.intervalTmp = scheduleData.interval;
          // this.form.randomInterval = scheduleData.randomInterval || "0";
        } else {
          this.$message.error(res.msg || "获取计划失败");
        }
      } catch (e) {} finally {
        this.formloading = false;
      }
    },
    // 提交计划设置
    handleSubmit() {
      // 校验表单
      this.form.interval = this.form.intervalTmp.toString();
      this.$refs['formRef'].validate(async (valid) => {
        if (valid) {
          try {
            // 按钮loading状态
            this.butloading = true;
            // 如果有日期，装换月份为英文月份，如12月转为December
            if (this.form.date) {
              this.form.date = this.convertMonthToEnglish(this.form.date);
              console.log(this.form.date);
            }

            let res
            if (this.activeName === 'scanPlan') {
              res = await this.$API.system.reqSetKasperskySchedule(this.form);
            }else if (this.activeName === 'reportPlan') {
              // 转换表达式
              const cronExpression = this.convertToCron(this.reportForm);
              // console.log('转换后的Cron表达式:', cronExpression);
              this.reqForm.cron = cronExpression;
              this.reqForm.enable = this.reportForm.enable;
              res = await this.$API.system.reqSetKasperskyReportSchedule(this.reqForm);
            }
            if (res.code === 200) {
              this.$message.success(res.msg || "计划设置成功");
              this.$refs.formRef.clearValidate();
            } else {
              this.$message.error(res.msg || "计划设置失败");
            }
            this.handleClose();
          } catch (e) {} finally {
            // 按钮loading状态
            this.butloading = false;
          }
        }else{
          this.$message.error("请填写完整计划设置");
        }
      });
    },
    /**
     * 将日期字符串（YYYY/MM/DD）中的数字月份替换为英文月份
     * @param {string} dateStr - 原始日期字符串（如 "2025/12/09"）
     * @returns {string} 替换后的日期字符串（如 "2025/December/09"）
    */
    convertMonthToEnglish(dateStr) {
      // 定义月份英文映射（索引 1-12 对应 1-12 月，索引 0 占位）
      const monthNames = [
        '', // 占位，索引 0 不使用
        'January', 'February', 'March', 'April', 'May', 'June',
        'July', 'August', 'September', 'October', 'November', 'December'
      ];

      // 按 / 分割日期为 年、月、日 三部分
      const [year, monthNum, day] = dateStr.split('/');

      // 校验格式：必须是 3 部分，且月份是 1-12 的数字
      if (!year || !monthNum || !day) {
        console.warn('日期格式错误，需为 YYYY/MM/DD 格式');
        return dateStr; // 格式错误则返回原字符串
      }
      const monthIndex = parseInt(monthNum, 10);
      if (isNaN(monthIndex) || monthIndex < 1 || monthIndex > 12) {
        console.warn('月份无效，需为 1-12 之间的数字');
        return dateStr;
      }

      // 替换月份并拼接结果
      const monthEn = monthNames[monthIndex];
      return `${year}/${monthEn}/${day}`;
    },
    // 立即扫描
    startScan() {
      this.$confirm("确定要执行立即扫描操作吗？", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning"
      }).then(async () => {
        try {
          const res = await this.$API.system.reqKasperskyScan();
          if (res.code === 200) {
            this.$message.success(res.message || "扫描启动成功");
          } else {
            this.$message.error(res.message || "扫描启动失败");
          }
        } catch (e) {}finally{
          this.getAntivirusInfo();
        }
      }).catch(() => {
        // 取消操作
        this.$message.info("已取消扫描操作");
      });
    },
    // 获取卡巴斯基隔离区
    async getIsolationZone() {
      try {
        const res = await this.$API.system.reqKasperskyIsolationZone();
        if (res.code === 200 && res.data === "") {
          this.isolationZoneStatus = "正常";
        } else {
          this.isolationZoneStatus = "异常";
        }
      } catch (e){}
    },
    handleClick(tab){
      if(tab.name === 'scanPlan'){
        this.getAntivirusSchedule();
      } else if(tab.name === 'reportPlan'){
        this.getReportSchedule();
      }
    },
    // 处理规则类型变化
    handleRuleTypeChange(val) {
      // 清理时间和间隔
      this.reportForm.time = '00:00';
      this.reportForm.interval = '';
      this.reportForm.enable = '';
    },

    // 获取报告计划
    async getReportSchedule() {
      try {
        this.formloading = true;
        const res = await this.$API.system.reqKasperskyReportSchedule();
        if (res.code === 200) {
          this.reportForm = res.data;
        } else {
          this.$message.error(res.msg || "获取报告计划失败");
        }
      } catch (e) {} finally {
        this.formloading = false;
      }
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
  }
}
</script>

<style scoped>
.antivirus-container {
  padding: 24px;
  border-radius: 12px;
  background: #f5f7fa;
  border: 1px solid rgba(100, 200, 150, 0.15);
  box-shadow: 
    0 4px 20px rgba(100, 200, 150, 0.08),
    0 8px 30px rgba(0, 0, 0, 0.05);
}

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

.but-right {
  position: absolute;
  top: 50%;
  right: 0;
  transform: translateY(-50%);
  display: flex;
  gap: 16px;
}

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

.status-row {
  margin-bottom: 24px;
}

/* 核心：卡片设置固定高度相关，防止缩放撑大 */
.status-card {
  text-align: center;
  padding: 24px 16px;
  border-radius: 10px;
  border-left: 6px solid #64c896 !important;
  background: linear-gradient(90deg, #f9fffb, #f0f9ff);
  border: 1px solid rgba(100, 200, 150, 0.2);
  box-shadow: 0 2px 12px rgba(100, 200, 150, 0.08);
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
  /* 关键：设置最小高度，固定卡片尺寸 */
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
  /* 固定占位，不随子元素变化 */
  height: 22px;
  line-height: 22px;
}

/* 核心优化：用transform缩放代替font-size修改，不改变占位 */
.status-value {
  font-size: 18px;
  color: #64c896;
  font-weight: 600;
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
  line-height: 24px;
  /* 关键：固定行高+inline-block，transform缩放不影响布局 */
  display: inline-block;
  transform-origin: center;
  height: 24px;
  margin: 0 auto;
}

/* hover时用scale缩放，不改变元素占位，不会撑大卡片 */
.status-card:hover .status-value,
.threat-card:hover .status-value {
  transform: scale(1.2); /* 18*1.2≈21.6，接近原22px效果 */
  text-shadow: 0 0 8px rgba(100, 200, 150, 0.2);
}

/* 卡片hover上浮效果 */
.status-card:hover,
.threat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(100, 200, 150, 0.12);
}

.status-desc {
  font-size: 14px;
  color: #94a3b8;
  margin-top: 8px;
  line-height: 1.6;
  /* 固定占位 */
  min-height: 22px;
}

/* 威胁卡片同步优化 */
.threat-card {
  text-align: center;
  padding: 24px 16px;
  border-radius: 10px;
  border-left: 6px solid #64c896 !important;
  background: linear-gradient(90deg, #f9fffb, #f0f9ff);
  border: 1px solid rgba(100, 200, 150, 0.2);
  box-shadow: 0 2px 12px rgba(100, 200, 150, 0.08);
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
  margin-top: 16px !important;
  /* 同步固定最小高度 */
  min-height: 140px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

/* Dialog内表单样式适配 */
:deep(.el-dialog .el-card) {
  background: transparent;
  border: none;
  box-shadow: none;
}

:deep(.el-form-item__label) {
  color: #2a3b47;
  font-weight: 500;
}

:deep(.el-form-item__content) {
  color: #475569;
}
</style>