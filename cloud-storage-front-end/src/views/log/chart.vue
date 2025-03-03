<!-- src/components/Chart.vue -->
<template>
  <div class="app-container">
    <!-- SSH连接对话框 -->
    <el-dialog v-model="sshDialogVisible" title="SSH连接" width="500px">
      <el-form :model="sshForm" label-width="100px">
        <el-form-item label="主机地址">
          <el-input v-model="sshForm.host" placeholder="请输入主机地址" />
        </el-form-item>
        <el-form-item label="端口">
          <el-input-number v-model="sshForm.port" :min="1" :max="65535" />
        </el-form-item>
        <el-form-item label="用户名">
          <el-input v-model="sshForm.user" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="sshForm.password" type="password" placeholder="请输入密码" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="sshDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSSHConnect" :loading="sshConnecting">连接</el-button>
        </span>
      </template>
    </el-dialog>

    <el-card class="chart-container">
      <div class="chart-header">
        <el-select v-model="currentHost" placeholder="选择主机" style="width: 200px" @change="handleHostChange">
          <el-option v-for="host in hosts" :key="host" :label="host" :value="host" />
        </el-select>
        <el-button type="primary" @click="showSSHDialog" style="margin-left: 10px">连接主机</el-button>
      </div>

      <el-tabs v-model="activeTab">
        <el-tab-pane label="实时监控" name="realtime">
          <!-- 实时监控配置 -->
          <div class="chart-config">
            <el-form :inline="true" :model="realtimeConfig">
              <el-form-item label="监控项" style="width: 350px">
                <el-select v-model="realtimeConfig.metrics" multiple placeholder="请选择监控项">
                  <el-option v-for="item in metricOptions" :key="item.value" :label="item.label" :value="item.value" />
                </el-select>
              </el-form-item>
              <el-form-item label="时间范围">
                <el-radio-group v-model="realtimeConfig.timeRange">
                  <el-radio-button label="1h">1小时</el-radio-button>
                  <el-radio-button label="6h">6小时</el-radio-button>
                  <el-radio-button label="12h">12小时</el-radio-button>
                  <el-radio-button label="24h">24小时</el-radio-button>
                </el-radio-group>
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="startMonitor">开始监控</el-button>
                <el-button @click="stopMonitor">停止监控</el-button>
              </el-form-item>
            </el-form>
          </div>

          <!-- 实时图表 -->
          <div ref="realtimeChart" class="chart-view" />
        </el-tab-pane>

        <el-tab-pane label="历史分析" name="history">
          <!-- 历史数据配置 -->
          <div class="chart-config">
            <el-form :inline="true" :model="historyConfig">
              <el-form-item label="时间范围">
                <el-date-picker
                  v-model="historyConfig.timeRange"
                  type="datetimerange"
                  range-separator="至"
                  start-placeholder="开始时间"
                  end-placeholder="结束时间"
                />
              </el-form-item>
              <el-form-item label="聚合方式">
                <el-select v-model="historyConfig.aggregation">
                  <el-option label="按分钟" value="minute" />
                  <el-option label="按小时" value="hour" />
                  <el-option label="按天" value="day" />
                </el-select>
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="queryHistory">查询</el-button>
              </el-form-item>
            </el-form>
          </div>

          <!-- 历史图表 -->
          <div ref="historyChart" class="chart-view" />
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted, onUnmounted, watch } from "vue"
import * as echarts from "echarts"
import { ElMessage } from "element-plus"
import { getRealtimeMetricsApi, getHistoryMetricsApi, connectSSHApi, getHostsApi } from "@/api/log/frontend"
import type { RealtimeMonitorParams, FrontHistoryAnalysisParams, ChartMetric } from "@/api/log/types/frontend"
import { useUserStoreHook } from "@/store/modules/user"

// 用户存储
const userStore = useUserStoreHook()

// SSH连接表单
const sshDialogVisible = ref(false)
const sshConnecting = ref(false)
const currentHost = ref("")
const hosts = ref<string[]>([])
const sshForm = reactive({
  host: "",
  port: 22,
  user: "root",
  password: ""
})

// 显示SSH连接对话框
const showSSHDialog = () => {
  sshDialogVisible.value = true
}

// 处理SSH连接
const handleSSHConnect = async () => {
  if (!sshForm.host) {
    ElMessage.warning("请输入主机地址")
    return
  }

  try {
    sshConnecting.value = true
    await connectSSHApi(sshForm.host, sshForm.port, sshForm.user, sshForm.password)
    ElMessage.success("SSH连接成功")

    // 刷新主机列表
    await refreshHosts()

    // 设置当前主机
    currentHost.value = sshForm.host

    sshDialogVisible.value = false
  } catch (error) {
    ElMessage.error(`SSH连接失败: ${error instanceof Error ? error.message : String(error)}`)
  } finally {
    sshConnecting.value = false
  }
}

// 刷新主机列表
const refreshHosts = async () => {
  try {
    hosts.value = await getHostsApi()
  } catch (error) {
    console.error("获取主机列表失败:", error)
  }
}

// 主机变更处理
const handleHostChange = (host: string) => {
  if (host) {
    userStore.setCurrentSSHHost(host)
  }
}

// 图表实例
const realtimeChart = ref<HTMLElement | null>(null)
const historyChart = ref<HTMLElement | null>(null)
let realtimeChartInstance: echarts.ECharts | null = null
let historyChartInstance: echarts.ECharts | null = null

// 标签页
const activeTab = ref("realtime")

// 监控项选项
const metricOptions: ChartMetric[] = [
  { label: "请求数", value: "requests" },
  { label: "错误数", value: "errors" },
  { label: "响应时间", value: "response_time" }
]

// 实时监控配置
const realtimeConfig = reactive<RealtimeMonitorParams>({
  metrics: ["requests"],
  timeRange: "1h"
})

// 历史分析配置
const historyConfig = reactive<FrontHistoryAnalysisParams>({
  timeRange: [null, null],
  aggregation: "hour"
})

// 定时器
let monitorTimer: number | null = null

// 初始化图表
onMounted(() => {
  if (realtimeChart.value) {
    realtimeChartInstance = echarts.init(realtimeChart.value)
  }
  if (historyChart.value) {
    historyChartInstance = echarts.init(historyChart.value)
  }

  // 监听窗口大小变化
  window.addEventListener("resize", handleResize)

  // 初始化数据
  initData()
})

onUnmounted(() => {
  stopMonitor()
  window.removeEventListener("resize", handleResize)
  realtimeChartInstance?.dispose()
  historyChartInstance?.dispose()
})

// 初始化数据
const initData = async () => {
  await refreshHosts()

  // 如果有当前SSH主机，则自动选择
  if (userStore.currentSSHHost) {
    currentHost.value = userStore.currentSSHHost
  }
}

// 窗口大小变化时重绘图表
const handleResize = () => {
  realtimeChartInstance?.resize()
  historyChartInstance?.resize()
}

// 开始实时监控
const startMonitor = async () => {
  if (!currentHost.value) {
    ElMessage.warning("请先连接主机")
    showSSHDialog()
    return
  }

  if (!realtimeConfig.metrics.length) {
    ElMessage.warning("请选择至少一个监控项")
    return
  }

  stopMonitor()
  await updateRealtimeChart()

  monitorTimer = window.setInterval(async () => {
    await updateRealtimeChart()
  }, 5000) // 每5秒更新一次
}

// 停止监控
const stopMonitor = () => {
  if (monitorTimer) {
    clearInterval(monitorTimer)
    monitorTimer = null
  }
}

// 更新实时图表
const updateRealtimeChart = async () => {
  try {
    const data = await getRealtimeMetricsApi(realtimeConfig)
    renderRealtimeChart(data)
  } catch (error) {
    console.error("更新实时图表失败:", error)
  }
}

// 渲染实时图表
const renderRealtimeChart = (data: any) => {
  if (!realtimeChartInstance) return

  const option = {
    title: {
      text: "实时监控"
    },
    tooltip: {
      trigger: "axis"
    },
    legend: {
      data: realtimeConfig.metrics
    },
    xAxis: {
      type: "time",
      splitLine: {
        show: false
      }
    },
    yAxis: {
      type: "value",
      splitLine: {
        lineStyle: {
          type: "dashed"
        }
      }
    },
    series: data.series
  }

  realtimeChartInstance.setOption(option)
}

// 查询历史数据
const queryHistory = async () => {
  if (!currentHost.value) {
    ElMessage.warning("请先连接主机")
    showSSHDialog()
    return
  }

  if (!historyConfig.timeRange[0] || !historyConfig.timeRange[1]) {
    ElMessage.warning("请选择时间范围")
    return
  }

  try {
    const data = await getHistoryMetricsApi(historyConfig)
    renderHistoryChart(data)
  } catch (error) {
    console.error("查询历史数据失败:", error)
    ElMessage.error("查询失败")
  }
}

// 渲染历史图表
const renderHistoryChart = (data: any) => {
  if (!historyChartInstance) return

  const option = {
    title: {
      text: "历史分析"
    },
    tooltip: {
      trigger: "axis"
    },
    legend: {
      data: data.metrics
    },
    xAxis: {
      type: "time",
      splitLine: {
        show: false
      }
    },
    yAxis: {
      type: "value",
      splitLine: {
        lineStyle: {
          type: "dashed"
        }
      }
    },
    series: data.series
  }

  historyChartInstance.setOption(option)
}

// 监听当前主机变化
watch(currentHost, (newHost) => {
  if (newHost) {
    userStore.setCurrentSSHHost(newHost)
  }
})
</script>

<style scoped>
.chart-container {
  margin-bottom: 20px;
}

.chart-header {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
}

.host-info {
  margin-left: 15px;
  color: #409eff;
  font-weight: bold;
}

.chart-config {
  margin-bottom: 20px;
  padding: 15px;
  background: #f8f9fa;
  border-radius: 4px;
}

.chart-view {
  height: 400px;
  margin-top: 20px;
}
</style>
