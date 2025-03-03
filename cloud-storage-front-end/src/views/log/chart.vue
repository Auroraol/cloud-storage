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

    <!-- 文件选择对话框 -->
    <el-dialog v-model="fileDialogVisible" title="选择文件" width="600px">
      <div class="file-browser">
        <div class="file-path-bar">
          <el-input v-model="fileBrowser.currentPath" placeholder="当前路径">
            <template #append>
              <el-button @click="loadFiles">刷新</el-button>
            </template>
          </el-input>
        </div>
        <div class="file-list" v-loading="fileBrowser.loading">
          <el-table
            :data="fileBrowser.files"
            style="width: 100%"
            @row-click="handleFileRowClick"
            highlight-current-row
            :row-class-name="getFileRowClass"
          >
            <el-table-column label="名称" prop="name">
              <template #default="{ row }">
                <el-icon v-if="!row.includes('.')"><Folder /></el-icon>
                <el-icon v-else><Document /></el-icon>
                <span style="margin-left: 8px">{{ row }}</span>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="120">
              <template #default="{ row }">
                <el-button v-if="!isDirectory(row)" type="primary" size="small" @click.stop="selectFile(row)">
                  选择
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
        <div class="file-pagination">
          <el-pagination
            v-model:current-page="fileBrowser.currentPage"
            v-model:page-size="fileBrowser.pageSize"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            :total="fileBrowser.total"
            @size-change="handleFileSizeChange"
            @current-change="handleFilePageChange"
          />
        </div>
        <div class="selected-file-info" v-if="fileBrowser.selectedFile">
          <span class="selected-label">已选择文件:</span>
          <span class="selected-value">{{ fileBrowser.selectedFile }}</span>
        </div>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="fileDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="confirmSelectFile" :disabled="!fileBrowser.selectedFile">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <div class="filter-container">
      <el-form :model="queryParams" ref="queryForm" :inline="true">
        <el-form-item label="主机" prop="host">
          <el-select
            v-model="currentHost"
            placeholder="选择主机"
            clearable
            style="width: 200px"
            @change="handleHostChange"
          >
            <el-option v-for="host in hosts" :key="host" :label="host" :value="host" />
          </el-select>
          <el-button type="primary" icon="Plus" circle size="small" @click="showSSHDialog" style="margin-left: 10px" />
        </el-form-item>
        <el-form-item label="数据文件" prop="dataFile">
          <el-input v-model="queryParams.dataFile" placeholder="选择数据文件" readonly style="width: 300px">
            <template #append>
              <el-button @click="showFileDialog">浏览</el-button>
            </template>
          </el-input>
        </el-form-item>
      </el-form>
    </div>

    <el-card class="chart-container">
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

          <!-- 分页区域 -->
          <div class="pagination-container" v-if="historyData.length > 0">
            <el-pagination
              v-model:current-page="historyPage.current"
              v-model:page-size="historyPage.size"
              :page-sizes="[10, 20, 50, 100]"
              layout="total, sizes, prev, pager, next, jumper"
              :total="historyPage.total"
              @size-change="handleHistorySizeChange"
              @current-change="handleHistoryPageChange"
            />
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted, onUnmounted, watch } from "vue"
import * as echarts from "echarts"
import { ElMessage } from "element-plus"
import {
  getRealtimeMetricsApi,
  getHistoryMetricsApi,
  connectSSHApi,
  getHostsApi,
  getLogFilesApi
} from "@/api/log/frontend"
import type { RealtimeMonitorParams, FrontHistoryAnalysisParams, ChartMetric } from "@/api/log/types/frontend"
import { useUserStoreHook } from "@/store/modules/user"
import { Document, Folder } from "@element-plus/icons-vue"

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

// 文件浏览器
const fileDialogVisible = ref(false)
const fileBrowser = reactive({
  currentPath: "/opt",
  files: [] as string[],
  loading: false,
  selectedFile: "",
  currentPage: 1,
  pageSize: 20,
  total: 0
})

// 查询参数
const queryParams = reactive({
  dataFile: "",
  host: ""
})

// 历史数据分页
const historyData = ref<any[]>([])
const historyPage = reactive({
  current: 1,
  size: 20,
  total: 0
})

// 显示SSH连接对话框
const showSSHDialog = () => {
  sshDialogVisible.value = true
}

// 显示文件选择对话框
const showFileDialog = () => {
  if (!currentHost.value) {
    ElMessage.warning("请先选择主机")
    return
  }

  fileDialogVisible.value = true
  loadFiles()
}

// 加载文件列表
const loadFiles = async () => {
  if (!currentHost.value) return

  try {
    fileBrowser.loading = true
    const response = await getLogFilesApi(currentHost.value, fileBrowser.currentPath)

    // 处理分页
    const startIndex = (fileBrowser.currentPage - 1) * fileBrowser.pageSize
    const endIndex = startIndex + fileBrowser.pageSize

    fileBrowser.files = response.slice(startIndex, endIndex)
    fileBrowser.total = response.length
  } catch (error) {
    ElMessage.error(`获取文件列表失败: ${error instanceof Error ? error.message : String(error)}`)
  } finally {
    fileBrowser.loading = false
  }
}

// 判断是否是目录
const isDirectory = (file: string) => {
  return !file.includes(".")
}

// 获取文件行的类名
const getFileRowClass = (row: { row: string; rowIndex: number }) => {
  const fullPath = fileBrowser.currentPath + row.row
  return fullPath === fileBrowser.selectedFile ? "selected-file-row" : ""
}

// 选择文件
const selectFile = (file: string) => {
  fileBrowser.selectedFile = file
}

// 处理文件点击
const handleFileClick = (file: string) => {
  // 判断是否是目录
  if (isDirectory(file)) {
    // 如果是目录，进入该目录
    fileBrowser.currentPath = fileBrowser.currentPath + file
    fileBrowser.currentPage = 1
    loadFiles()
  } else {
    // 如果是文件，选中该文件
    selectFile(file)
  }
}

// 确认选择文件
const confirmSelectFile = () => {
  if (fileBrowser.selectedFile) {
    queryParams.dataFile = fileBrowser.selectedFile
    // 同步更新配置中的dataFile
    realtimeConfig.dataFile = fileBrowser.selectedFile
    historyConfig.dataFile = fileBrowser.selectedFile
    fileDialogVisible.value = false
    ElMessage.success(`已选择文件: ${fileBrowser.selectedFile}`)
  } else {
    ElMessage.warning("请选择一个文件")
  }
}

// 处理文件分页大小变化
const handleFileSizeChange = (size: number) => {
  fileBrowser.pageSize = size
  loadFiles()
}

// 处理文件分页页码变化
const handleFilePageChange = (page: number) => {
  fileBrowser.currentPage = page
  loadFiles()
}

// 处理历史数据分页大小变化
const handleHistorySizeChange = (size: number) => {
  historyPage.size = size
  displayHistoryData()
}

// 处理历史数据分页页码变化
const handleHistoryPageChange = (page: number) => {
  historyPage.current = page
  displayHistoryData()
}

// 显示历史数据（分页）
const displayHistoryData = () => {
  if (!historyChartInstance || !historyData.value.length) return

  try {
    // 根据当前页和页大小计算要显示的数据
    const start = (historyPage.current - 1) * historyPage.size
    const end = Math.min(start + historyPage.size, historyData.value.length)

    // 获取当前页的数据
    const currentPageData = historyData.value.slice(start, end)

    // 重新渲染图表
    const newSeries = currentPageData.map((item: any) => {
      return {
        name: item.name,
        type: "line",
        smooth: true,
        showSymbol: false,
        data: item.data
      }
    })

    // 更新图表
    historyChartInstance.setOption({
      series: newSeries
    })
  } catch (error) {
    console.error("显示历史数据失败:", error)
    ElMessage.error("显示历史数据失败")
  }
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
    queryParams.host = host
    // 清空已选文件
    queryParams.dataFile = ""
    fileBrowser.selectedFile = ""
    // 同步清空配置中的dataFile
    realtimeConfig.dataFile = ""
    historyConfig.dataFile = ""
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
  timeRange: "1h",
  dataFile: ""
})

// 历史分析配置
const historyConfig = reactive<FrontHistoryAnalysisParams>({
  timeRange: [null, null],
  aggregation: "hour",
  dataFile: ""
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
    queryParams.host = userStore.currentSSHHost
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

  if (!queryParams.dataFile) {
    ElMessage.warning("请选择数据文件")
    showFileDialog()
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
    // 确保dataFile是最新的
    realtimeConfig.dataFile = queryParams.dataFile
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

  if (!queryParams.dataFile) {
    ElMessage.warning("请选择数据文件")
    showFileDialog()
    return
  }

  if (!historyConfig.timeRange[0] || !historyConfig.timeRange[1]) {
    ElMessage.warning("请选择时间范围")
    return
  }

  try {
    // 确保dataFile是最新的
    historyConfig.dataFile = queryParams.dataFile
    const data = await getHistoryMetricsApi(historyConfig)

    // 检查数据结构
    if (!data || !data.series || !Array.isArray(data.series) || data.series.length === 0) {
      ElMessage.warning("未获取到有效的历史数据")
      return
    }

    // 保存原始数据
    historyData.value = data.series || []
    historyPage.total = data.total || historyData.value.length
    historyPage.current = 1

    // 渲染图表
    renderHistoryChart(data)
  } catch (error) {
    console.error("查询历史数据失败:", error)
    ElMessage.error(`查询失败: ${error instanceof Error ? error.message : String(error)}`)
  }
}

// 渲染历史图表
const renderHistoryChart = (data: any) => {
  if (!historyChartInstance) return

  // 确保数据格式正确
  if (!data || !data.series || !Array.isArray(data.series) || data.series.length === 0) {
    console.error("历史图表数据格式不正确:", data)
    ElMessage.warning("历史数据格式不正确，无法显示图表")
    return
  }

  // 处理图表数据
  const metrics = data.metrics || []
  const series = data.series.map((item: any) => {
    // 确保数据点格式正确
    const formattedData = Array.isArray(item.data)
      ? item.data.map((point: any) => {
          // 如果数据点是数组格式 [timestamp, value]
          if (Array.isArray(point)) {
            return point
          }
          // 如果数据点是对象格式 {time: timestamp, value: value}
          else if (point && typeof point === "object" && "time" in point && "value" in point) {
            return [point.time, point.value]
          }
          return point
        })
      : []

    return {
      name: item.name,
      type: "line",
      smooth: true,
      showSymbol: false,
      symbolSize: 6,
      lineStyle: {
        width: 2
      },
      data: formattedData
    }
  })

  const option = {
    title: {
      text: "历史分析"
    },
    tooltip: {
      trigger: "axis",
      formatter: function (params: any) {
        if (!params || !params.length) return ""

        const date = new Date(params[0].value[0])
        const formattedDate = `${date.getFullYear()}-${(date.getMonth() + 1).toString().padStart(2, "0")}-${date.getDate().toString().padStart(2, "0")} ${date.getHours().toString().padStart(2, "0")}:${date.getMinutes().toString().padStart(2, "0")}:${date.getSeconds().toString().padStart(2, "0")}`

        let result = `${formattedDate}<br/>`
        params.forEach((param: any) => {
          const value = Array.isArray(param.value) && param.value.length > 1 ? param.value[1] : param.value
          result += `${param.seriesName}: ${value}<br/>`
        })

        return result
      }
    },
    legend: {
      data: metrics,
      type: "scroll",
      orient: "horizontal",
      bottom: 0
    },
    grid: {
      left: "3%",
      right: "4%",
      bottom: "10%",
      top: "10%",
      containLabel: true
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
    series: series
  }

  // 完全重置图表
  historyChartInstance.clear()
  historyChartInstance.setOption(option, true)

  // 如果有数据，显示第一页
  if (historyData.value.length > 0) {
    displayHistoryData()
  }
}

// 监听当前主机变化
watch(currentHost, (newHost) => {
  if (newHost) {
    userStore.setCurrentSSHHost(newHost)
    queryParams.host = newHost
  }
})

// 处理表格行点击
const handleFileRowClick = (row: string) => {
  handleFileClick(row)
}
</script>

<style scoped>
.filter-container {
  margin-bottom: 20px;
  background: #fff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

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

.file-browser {
  display: flex;
  flex-direction: column;
  height: 400px;
}

.file-path-bar {
  margin-bottom: 10px;
}

.file-list {
  flex: 1;
  overflow-y: auto;
  border: 1px solid #ebeef5;
  border-radius: 4px;
}

.file-pagination {
  margin-top: 10px;
  display: flex;
  justify-content: flex-end;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.selected-file-info {
  margin-top: 10px;
  padding: 8px;
  background-color: #f0f9eb;
  border-radius: 4px;
  display: flex;
  align-items: center;
}

.selected-label {
  font-weight: bold;
  margin-right: 8px;
  color: #606266;
}

.selected-value {
  color: #67c23a;
  word-break: break-all;
}

:deep(.selected-file-row) {
  background-color: #f0f9eb;
}

:deep(.selected-file-row td) {
  background-color: #f0f9eb !important;
}
</style>
