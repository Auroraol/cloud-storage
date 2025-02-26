<template>
  <div class="real-time-monitor">
    <div class="monitor-header">
      <div class="tabs">
        <div class="tab" :class="{ active: activeTab === 'realtime' }" @click="activeTab = 'realtime'">实时监控</div>
        <div class="tab" :class="{ active: activeTab === 'history' }" @click="activeTab = 'history'">历史分析</div>
      </div>
    </div>

    <!-- SSH连接表单 -->
    <div v-if="!isConnected" class="ssh-connect-form">
      <h3>连接主机</h3>
      <el-form :model="sshForm" label-width="100px">
        <el-form-item label="主机地址">
          <el-input v-model="sshForm.host" placeholder="请输入主机地址"></el-input>
        </el-form-item>
        <el-form-item label="用户名">
          <el-input v-model="sshForm.user" placeholder="请输入用户名"></el-input>
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="sshForm.password" type="password" placeholder="请输入密码"></el-input>
        </el-form-item>
        <el-form-item label="私钥路径">
          <el-input v-model="sshForm.privateKeyPath" placeholder="请输入私钥路径（可选）"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="connectSSH">连接</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div v-else>
      <div v-if="activeTab === 'realtime'" class="monitor-content">
        <div class="connection-info">
          <span>已连接到主机: {{ sshForm.host }}</span>
          <el-button size="small" @click="disconnectSSH">断开连接</el-button>
        </div>

        <div class="filter-bar">
          <div class="log-file-select">
            <span>日志文件：</span>
            <el-select v-model="logFile" placeholder="请选择日志文件" @change="fetchLogFiles">
              <el-option
                v-for="file in logFiles"
                :key="file"
                :label="file.split('/').pop()"
                :value="file"
              ></el-option>
            </el-select>
            <el-button size="small" @click="fetchLogFiles">刷新</el-button>
          </div>

          <div class="monitor-items">
            <span>监控项：</span>
            <el-checkbox-group v-model="monitorItems">
              <el-checkbox label="请求数"></el-checkbox>
              <el-checkbox label="错误数"></el-checkbox>
              <el-checkbox label="响应时间"></el-checkbox>
            </el-checkbox-group>
          </div>

          <div class="time-range">
            <span>时间范围：</span>
            <el-radio-group v-model="timeRange">
              <el-radio-button :label="1">1小时</el-radio-button>
              <el-radio-button :label="6">6小时</el-radio-button>
              <el-radio-button :label="12">12小时</el-radio-button>
              <el-radio-button :label="24">24小时</el-radio-button>
            </el-radio-group>
          </div>
        </div>

        <div class="action-buttons">
          <el-button type="primary" @click="startMonitoring">开始监控</el-button>
          <el-button @click="stopMonitoring">停止监控</el-button>
        </div>

        <div class="chart-container" v-if="chartData.length > 0">
          <div ref="chartRef" class="chart"></div>
        </div>
        <div v-else class="no-data">
          <p>暂无监控数据，请选择监控项并点击"开始监控"</p>
        </div>
      </div>

      <div v-else class="history-content">
        <HistoryAnalysis :host="sshForm.host" :logFiles="logFiles" />
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted, watch } from 'vue';
import * as echarts from 'echarts';
import HistoryAnalysis from './HistoryAnalysis.vue';
import axios from 'axios';

export default {
  name: 'RealTimeMonitor',
  components: {
    HistoryAnalysis
  },
  setup() {
    const activeTab = ref('realtime');
    const isConnected = ref(false);
    const sshForm = ref({
      host: '',
      user: '',
      password: '',
      privateKeyPath: ''
    });
    const logFiles = ref([]);
    const logFile = ref('');
    const monitorItems = ref(['请求数']);
    const timeRange = ref(1);
    const chartData = ref([]);
    const chartRef = ref(null);
    const chart = ref(null);
    const monitoringInterval = ref(null);

    // 连接SSH
    const connectSSH = async () => {
      try {
        const response = await axios.post('/api/log_service/v1/ssh/connect', {
          host: sshForm.value.host,
          user: sshForm.value.user,
          password: sshForm.value.password,
          private_key_path: sshForm.value.privateKeyPath
        });

        if (response.data.success) {
          isConnected.value = true;
          fetchLogFiles();
        } else {
          alert('连接失败: ' + response.data.message);
        }
      } catch (error) {
        console.error('连接失败:', error);
        alert('连接失败: ' + error.message);
      }
    };

    // 断开SSH连接
    const disconnectSSH = () => {
      isConnected.value = false;
      stopMonitoring();
      logFiles.value = [];
      logFile.value = '';
      chartData.value = [];
    };

    // 获取日志文件列表
    const fetchLogFiles = async () => {
      try {
        const response = await axios.post('/api/log_service/v1/ssh/logfiles', {
          host: sshForm.value.host,
          path: '/var/log' // 默认日志路径
        });

        if (response.data.success) {
          logFiles.value = response.data.files;
          if (logFiles.value.length > 0 && !logFile.value) {
            logFile.value = logFiles.value[0];
          }
        } else {
          console.error('获取日志文件列表失败');
        }
      } catch (error) {
        console.error('获取日志文件列表失败:', error);
      }
    };

    // 获取监控数据
    const fetchMonitorData = async () => {
      if (!logFile.value) {
        alert('请选择日志文件');
        return;
      }

      try {
        const response = await axios.post('/api/log_service/v1/monitor/realtime', {
          host: sshForm.value.host,
          log_file: logFile.value,
          monitor_items: monitorItems.value,
          time_range: timeRange.value
        });

        if (response.data.success) {
          chartData.value = response.data.data;
          updateChart();
        }
      } catch (error) {
        console.error('获取监控数据失败:', error);
      }
    };

    // 开始监控
    const startMonitoring = () => {
      if (monitorItems.value.length === 0) {
        alert('请至少选择一个监控项');
        return;
      }

      fetchMonitorData();
      
      // 每30秒更新一次数据
      stopMonitoring();
      monitoringInterval.value = setInterval(fetchMonitorData, 30000);
    };

    // 停止监控
    const stopMonitoring = () => {
      if (monitoringInterval.value) {
        clearInterval(monitoringInterval.value);
        monitoringInterval.value = null;
      }
    };

    // 初始化图表
    const initChart = () => {
      if (chartRef.value) {
        chart.value = echarts.init(chartRef.value);
      }
    };

    // 更新图表
    const updateChart = () => {
      if (!chart.value) return;

      const series = [];
      const xAxisData = [];
      const dataMap = {};

      // 按类型分组数据
      chartData.value.forEach(item => {
        if (!dataMap[item.type]) {
          dataMap[item.type] = [];
        }
        dataMap[item.type].push({
          value: item.value,
          timestamp: item.timestamp
        });
      });

      // 为每种类型创建一个系列
      Object.keys(dataMap).forEach(type => {
        // 按时间戳排序
        const sortedData = dataMap[type].sort((a, b) => a.timestamp - b.timestamp);
        
        // 提取时间戳作为X轴数据（仅对第一个系列）
        if (series.length === 0) {
          xAxisData.push(...sortedData.map(item => {
            const date = new Date(item.timestamp * 1000);
            return `${date.getHours()}:${date.getMinutes().toString().padStart(2, '0')}`;
          }));
        }

        // 创建系列
        series.push({
          name: type,
          type: 'line',
          smooth: true,
          data: sortedData.map(item => item.value)
        });
      });

      // 设置图表选项
      const option = {
        title: {
          text: '实时监控数据'
        },
        tooltip: {
          trigger: 'axis'
        },
        legend: {
          data: Object.keys(dataMap)
        },
        grid: {
          left: '3%',
          right: '4%',
          bottom: '3%',
          containLabel: true
        },
        xAxis: {
          type: 'category',
          boundaryGap: false,
          data: xAxisData
        },
        yAxis: {
          type: 'value'
        },
        series: series
      };

      chart.value.setOption(option);
    };

    // 监听时间范围变化
    watch(timeRange, () => {
      if (monitoringInterval.value) {
        fetchMonitorData();
      }
    });

    // 组件挂载时初始化图表
    onMounted(() => {
      initChart();
      window.addEventListener('resize', () => {
        chart.value && chart.value.resize();
      });
    });

    // 组件卸载时清除定时器和图表
    onUnmounted(() => {
      stopMonitoring();
      if (chart.value) {
        chart.value.dispose();
        chart.value = null;
      }
    });

    return {
      activeTab,
      isConnected,
      sshForm,
      logFiles,
      logFile,
      monitorItems,
      timeRange,
      chartData,
      chartRef,
      connectSSH,
      disconnectSSH,
      fetchLogFiles,
      startMonitoring,
      stopMonitoring
    };
  }
};
</script>

<style scoped>
.real-time-monitor {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.monitor-header {
  padding: 10px 0;
  border-bottom: 1px solid #eee;
}

.tabs {
  display: flex;
}

.tab {
  padding: 8px 16px;
  cursor: pointer;
  border-bottom: 2px solid transparent;
}

.tab.active {
  border-bottom: 2px solid #409eff;
  color: #409eff;
}

.ssh-connect-form {
  max-width: 500px;
  margin: 20px auto;
  padding: 20px;
  border: 1px solid #eee;
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.connection-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px;
  margin-bottom: 10px;
  background-color: #f0f9eb;
  border-radius: 4px;
}

.monitor-content, .history-content {
  flex: 1;
  padding: 20px;
  overflow: auto;
}

.filter-bar {
  display: flex;
  flex-direction: column;
  gap: 15px;
  margin-bottom: 20px;
  padding: 15px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.log-file-select, .monitor-items, .time-range {
  display: flex;
  align-items: center;
  gap: 10px;
}

.action-buttons {
  margin-bottom: 20px;
}

.chart-container {
  width: 100%;
  height: 400px;
  margin-top: 20px;
}

.chart {
  width: 100%;
  height: 100%;
}

.no-data {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 400px;
  background-color: #f5f7fa;
  border-radius: 4px;
}
</style> 