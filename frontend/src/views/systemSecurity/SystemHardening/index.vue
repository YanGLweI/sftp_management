<template>
  <div class="system-hardening-container">
    <!-- 系统加固标题 -->
    <div class="page-title">系统加固管理
      <div class="but-right">
        <el-tooltip effect="dark" content="立即加固" placement="top">
          <el-button circle icon="el-icon-video-play" class="icon-btn" @click="startHardening"></el-button>
        </el-tooltip>
        <el-tooltip effect="dark" content="计划" placement="top">
          <el-button circle icon="el-icon-setting" class="icon-btn" @click="handlePlanClick"></el-button> 
        </el-tooltip>
      </div>
    </div>

    <!-- 状态卡片区域 -->
    <el-row :gutter="20" class="status-row">
      <el-col :span="12">
        <el-card class="status-card">
          <div class="status-label">系统加固状态</div>
          <div 
            class="status-value" 
            :class="{ 'status-error': status.hardening.status === '异常' || status.hardening.status === '获取中' }"
          >
            {{ status.hardening.status }}
          </div>
          <div class="status-desc">{{ status.hardening.desc }}</div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card class="status-card">
          <div class="status-label">加固任务状态</div>
          <div 
            class="status-value" 
            :class="{ 'status-error': status.task.status === '异常' || status.task.status === '获取中' }"
          >
            {{ status.task.status }}
          </div>
          <div class="status-desc">{{ status.task.desc }}</div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 加固信息表格 -->
    <el-table
      :data="hardeningTableData"
      v-loading="loading"
      border
      stripe
      style="width: 100%; margin-top: 20px;border-radius: 10px;"
      :header-cell-style="{ textAlign: 'center' }"
    >
      <el-table-column type="expand">
        <template slot-scope="props">
            <el-form label-position="left" inline class="demo-table-expand" label-width="180px" style="margin: 20px;">
              <el-divider>DNF/Repo 配置</el-divider>
              <el-form-item label="Kernel">
                <span>{{ props.row.kernel }}</span>
              </el-form-item>
              <el-form-item label="dnf_conf_gpgcheck">
                <span>{{ props.row.dnf_conf_gpgcheck }}</span>
              </el-form-item>
              <el-form-item label="redhat_repo_gpgcheck">
                <span>{{ props.row.redhat_repo_gpgcheck }}</span>
              </el-form-item>
              <el-divider>密码策略</el-divider>
              <el-form-item label="pass_max_days">
                <span>{{ props.row.pass_max_days }}</span>
              </el-form-item>
              <el-form-item label="pass_min_days">
                <span>{{ props.row.pass_min_days }}</span>
              </el-form-item>
              <el-form-item label="pass_min_len">
                <span>{{ props.row.pass_min_len }}</span>
              </el-form-item>
              <el-form-item label="pass_warn_age">
                <span>{{ props.row.pass_warn_age }}</span>
              </el-form-item>
              <el-form-item label="inactive">
                <span>{{ props.row.inactive }}</span>
              </el-form-item>
              <el-form-item label="gid">
                <span>{{ props.row.gid }}</span>
              </el-form-item>
              <el-form-item label="tmout">
                <span>{{ props.row.tmout }}</span>
              </el-form-item>
              <el-divider>Cron/At 任务配置</el-divider>
              <el-form-item label="cron">
                <span>{{ props.row.cron }}</span>
              </el-form-item>
              <el-form-item label="crontab">
                <span>{{ props.row.crontab }}</span>
              </el-form-item>
              <el-form-item label="cron_hourly">
                <span>{{ props.row.cron_hourly }}</span>
              </el-form-item>
              <el-form-item label="cron_daily">
                <span>{{ props.row.cron_daily }}</span>
              </el-form-item>
              <el-form-item label="cron_weekly">
                <span>{{ props.row.cron_weekly }}</span>
              </el-form-item>
              <el-form-item label="cron_monthly">
                <span>{{ props.row.cron_monthly }}</span>
              </el-form-item>
              <el-form-item label="cron_deny">
                <span>{{ props.row.cron_deny }}</span>
              </el-form-item>
              <el-form-item label="at_deny">
                <span>{{ props.row.at_deny }}</span>
              </el-form-item>
              <el-form-item label="cron_allow">
                <span>{{ props.row.cron_allow }}</span>
              </el-form-item>
              <el-form-item label="at_allow">
                <span>{{ props.row.at_allow }}</span>
              </el-form-item>
              <el-divider>SSHD 配置</el-divider>
              <el-form-item label="sshd_config">
                <span>{{ props.row.sshd_config }}</span>
              </el-form-item>
              <el-form-item label="log_level">
                <span>{{ props.row.log_level }}</span>
              </el-form-item>
              <el-form-item label="x11_forwarding">
                <span>{{ props.row.x11_forwarding }}</span>
              </el-form-item>
              <el-form-item label="max_auth_tries">
                <span>{{ props.row.max_auth_tries }}</span>
              </el-form-item>
              <el-form-item label="ignore_rhosts">
                <span>{{ props.row.ignore_rhosts }}</span>
              </el-form-item>
              <el-form-item label="hostbased_authentication">
                <span>{{ props.row.hostbased_authentication }}</span>
              </el-form-item>
              <el-form-item label="permit_root_login">
                <span>{{ props.row.permit_root_login }}</span>
              </el-form-item>
              <el-form-item label="permit_empty_passwords">
                <span>{{ props.row.permit_empty_passwords }}</span>
              </el-form-item>
              <el-form-item label="permit_user_environment">
                <span>{{ props.row.permit_user_environment }}</span>
              </el-form-item>
              <el-form-item label="client_alive_interval">
                <span>{{ props.row.client_alive_interval }}</span>
              </el-form-item>
              <el-form-item label="client_alive_count_max">
                <span>{{ props.row.client_alive_count_max }}</span>
              </el-form-item>
              <el-form-item label="login_grace_time">
                <span>{{ props.row.login_grace_time }}</span>
              </el-form-item>
              <el-divider>密码复杂度</el-divider>
              <el-form-item label="minlen">
                <span>{{ props.row.minlen }}</span>
              </el-form-item>
              <el-form-item label="minclass">
                <span>{{ props.row.minclass }}</span>
              </el-form-item>
              <el-form-item label="dcredit">
                <span>{{ props.row.dcredit }}</span>
              </el-form-item>
              <el-form-item label="ucredit">
                <span>{{ props.row.ucredit }}</span>
              </el-form-item>
              <el-form-item label="lcredit">
                <span>{{ props.row.lcredit }}</span>
              </el-form-item>
              <el-form-item label="ocredit">
                <span>{{ props.row.ocredit }}</span>
              </el-form-item>
              <el-form-item label="password_remember">
                <span>{{ props.row.password_remember }}</span>
              </el-form-item>
              <el-divider>系统文件内容</el-divider>
              <el-form-item label="passwd">
                <span>{{ props.row.passwd }}</span>
              </el-form-item>
              <el-form-item label="passwd-">
                <span>{{ props.row.passwd_dash }}</span>
              </el-form-item>
              <el-form-item label="group">
                <span>{{ props.row.group }}</span>
              </el-form-item>
              <el-form-item label="group-">
                <span>{{ props.row.group_dash }}</span>
              </el-form-item>
              <el-form-item label="shadow">
                <span>{{ props.row.shadow }}</span>
              </el-form-item>
              <el-form-item label="shadow-">
                <span>{{ props.row.shadow_dash }}</span>
              </el-form-item>
              <el-form-item label="gshadow">
                <span>{{ props.row.gshadow }}</span>
              </el-form-item>
              <el-form-item label="gshadow-">
                <span>{{ props.row.gshadow_dash }}</span>
              </el-form-item>
              <el-divider>加密/时间 策略</el-divider>
              <el-form-item label="crypto_policies">
                <span>{{ props.row.crypto_policies }}</span>
              </el-form-item>
              <el-form-item label="ntp_server">
                <span>{{ props.row.ntp_server }}</span>
              </el-form-item>
            </el-form>
        </template>
      </el-table-column>
      <el-table-column
        type="index"
        label="序号"
        align="center"
        width="80"
      />
      <el-table-column
        prop="hostname"
        label="主机名"
        align="center"
        min-width="120"
      />
      <el-table-column
        prop="ip"
        label="IP地址"
        align="center"
        min-width="150"
      />
      <el-table-column
        prop="operasystem"
        label="操作系统"
        align="center"
        min-width="180"
      />
      <el-table-column
        prop="result"
        label="加固状态"
        align="center"
        min-width="120"
      >
        <template slot-scope="scope">
          <el-tag 
            :type="scope.row.result === '正常' ? 'success' : (scope.row.result === '异常' ? 'danger' : 'warning')"
          >
            {{ scope.row.result }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column
        prop="date"
        label="日期"
        align="center"
        min-width="180"
      />
    </el-table>

    <!-- 分页组件 -->
    <el-pagination
      style="margin-top: 20px; text-align: center;"
      @size-change="handleSizeChange"
      @current-change="getHardeningList"
      :current-page="params.pageNum"
      :page-sizes="[1, 3, 5, 10]"
      :page-size="params.pageSize"
      layout="sizes, prev, pager, next, jumper,->,total"
      :total="total"
    >
    </el-pagination>

    <!-- 加固计划弹窗 -->
    <el-dialog
      title="系统加固计划"
      :visible.sync="dialogPlanVisible"
      width="480px"
      :close-on-click-modal="false"
      @close="closePlanDialog"
    >
      <!-- 弹窗内容 -->
      <el-tabs v-model="activeName" @tab-click="handleClick">
          <el-tab-pane label="加固任务" name="hardeningPlan">
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
  name: "SystemHardening",
  data() {
    return {
      // 表格加载状态
      loading: false,
      // 分页参数
      params: {
        pageNum: 1,
        pageSize: 5
      },
      total: 0,
      hardeningTableData: [],
      // 最新加固记录
      latestHardening: {},
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
      activeName: 'hardeningPlan',
    };
  },
  mounted() {
    // 初始化加载表格数据
    this.getHardeningList();
  },
  computed: {
    // 状态计算（模拟数据，后续可对接后端）
    status() {
      return {
        // 系统加固状态
        hardening: this.getHardeningStatus(),
        // 加固任务状态
        task: this.getHardeningPlanStatus()
      };
    }
  },
  methods: {
    // 获取加固列表数据（模拟后端请求）
    async getHardeningList(pages = 1) {
      // 设置分页参数
      this.params.pageNum = pages;

      try {
        this.loading = true;
        const res = await this.$API.system.reqSystemHardeningList(this.params);
        this.hardeningTableData = res.data.checklist || [];
        this.hardeningTableData.forEach(item => {
          item.date = item.date.replace(/T(.*?)(\+.*)?$/, ' $1')
        })
        this.total = res.data.total || 0;
        // 如果是第一页，记录最新加固记录
        if (pages === 1 && res.data.total > 0) {
          this.latestHardening = res.data.checklist[0] || {};
          this.latestHardening.date = this.latestHardening.date.replace(/T(.*?)(\+.*)?$/, ' $1')
        }
      } catch (error) {
        
      } finally {
        this.loading = false;
      }
    },
    // 分页-每页条数变化
    handleSizeChange(val) {
      this.params.pageSize = val;
      this.getHardeningList();
    },
    // // 分页-当前页变化
    // handleCurrentChange(val) {
    //   this.params.pageNum = val;
    //   this.getHardeningList();
    // },
    // 立即加固
    startHardening() {
      this.$confirm("确定要执行立即加固操作吗？", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning"
      }).then(async () => {
        try {
          // 模拟接口请求
          this.loading = true
          await this.$API.system.reqSystemHardening();
          this.$message.success("系统加固已启动");
          // 刷新列表
          this.getHardeningList();
        } catch (error) {
          this.$message.error("系统加固启动失败");
        } finally {
          this.loading = false;
        }
      }).catch(() => {
        this.$message.info("已取消加固操作");
      });
    },
    // 打开计划设置弹窗
    openScheduleDialog() {
      this.dialogVisible = true;
    },

    // 获取加固状态
    getHardeningStatus() {

      const latestHardening = this.latestHardening || {};
      
      if (this.total === 0) {
        return { status: "异常", desc: "未检测到加固记录" };
      }

      if (!['正常'].includes(latestHardening.result)) {
        return { 
          status: "异常", 
          desc: `最新加固检查发现异常（${latestHardening.date}）` 
        };
      }
      
      return { 
        status: "正常", 
        desc: `核心系统配置已加固` 
      };
    },
    getHardeningPlanStatus() {
      const latestHardening = this.latestHardening || {};
      if (this.total === 0) {
        return { status: "异常", desc: "未检测到加固任务记录" };
      }

      // 判断最新任务时间，距离现在是否超过1天，超过一天则为异常
      const now = new Date();
      const lastDate = new Date(latestHardening.date);
      const diffTime = now - lastDate;
      const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24));

      if (diffDays > 1) {
        return { 
          status: "异常", 
          desc: `加固任务异常（${latestHardening.date}）（距离现在${diffDays}天）` 
        };
      }
      return { 
        status: "正常", 
        desc: `加固任务正常（${latestHardening.date}）` 
      };
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
      this.activeName = 'hardeningPlan';
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
            if ( this.activeName === 'hardeningPlan'){
              res = await this.$API.system.reqSetSystemHardeningSchedule(this.reqForm);
            } else if ( this.activeName === 'reportPlan'){
              res = await this.$API.system.reqSetSystemHardeningReportSchedule(this.reqForm);
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
      this.getHardeningSchedule();
    },

    // 获取系统加固计划
    async getHardeningSchedule(){
      try {
        this.dialogPlanVisible = true;
        this.formloading = true;
        const res = await this.$API.system.reqSystemHardeningSchedule();
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

    // 获取系统加固报告计划
    async getReportSchedule(){
      try {
        this.formloading = true;
        const res = await this.$API.system.reqSystemHardeningReportSchedule();
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
    // 处理选项卡点击
    handleClick(tab) {
      if (tab.name === 'hardeningPlan') {
        this.getHardeningSchedule();
      } else if (tab.name === 'reportPlan') {
        this.getReportSchedule();
      }
    },
  }
};
</script>

<style scoped>
/* 复用目标页面的容器样式 */
.system-hardening-container {
  padding: 24px;
  border-radius: 12px;
  background: #f5f7fa;
  border: 1px solid rgba(100, 200, 150, 0.15);
  box-shadow: 
    0 4px 20px rgba(100, 200, 150, 0.08),
    0 8px 30px rgba(0, 0, 0, 0.05);
}

/* 标题样式（和目标页面一致） */
.page-title {
  text-align: center;
  font-size: 18px;
  font-weight: 600;
  color: #2a3b47;
  margin-bottom: 24px;
  position: relative;
  letter-spacing: 0.3px;
}

/* 右侧按钮容器（和目标页面一致） */
.but-right {
  position: absolute;
  top: 50%;
  right: 0;
  transform: translateY(-50%);
  display: flex;
  gap: 16px;
}

/* 图标按钮样式（和目标页面一致） */
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

/* 状态行样式 */
.status-row {
  margin-bottom: 24px;
}

/* 状态卡片样式（和目标页面一致） */
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

/* 状态文本样式 */
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

/* 弹窗空提示样式 */
.dialog-empty-tip {
  text-align: center;
  padding: 40px 0;
  color: #94a3b8;
  font-size: 14px;
}

/* 弹窗样式适配 */
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

.demo-table-expand {
  font-size: 0;
}
.demo-table-expand label {
  width: 90px !important;
  color: #99a9bf !important;
}
.demo-table-expand .el-form-item {
  margin-right: 0 !important;
  margin-bottom: 0 !important;
  width: 50% !important;
}
</style>