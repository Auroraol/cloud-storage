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
        <!-- 添加模式切换 -->
        <el-form-item>
          <el-radio-group v-model="mode" @change="handleModeChange">
            <el-radio-button label="ssh">SSH模式</el-radio-button>
            <el-radio-button label="local">本地模式</el-radio-button>
          </el-radio-group>
        </el-form-item>

        <!-- SSH模式下的主机选择 -->
        <el-form-item label="主机" prop="host" v-if="mode === 'ssh'">
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

    <!-- 图像 -->
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
                <el-button type="primary" @click="startMonitor" :disabled="isMonitoring">开始监控</el-button>
                <el-button @click="stopMonitor" :type="isMonitoring ? 'primary' : 'default'" :disabled="!isMonitoring"
                  >停止监控</el-button
                >
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
import {
  getRealtimeMetricsApi,
  getHistoryMetricsApi,
  connectSSHApi,
  getHostsApi,
  getLogFilesApi,
  getSSHConnectionsApi,
  getLocalLogFilesApi,
  getLocalLogMetricsApi
} from "@/api/log/frontend"
import type { RealtimeMonitorParams, FrontHistoryAnalysisParams, ChartMetric } from "@/api/log/types/frontend"
import { useUserStore } from "@/store/modules/user"
import { Document, Folder } from "@element-plus/icons-vue"
import dayjs from "dayjs"

// 用户存储
const userStore = useUserStore()

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

// 历史数据
const historyData = ref<any[]>([])

// 添加模式状态
const mode = ref<"ssh" | "local">("local")

// 添加监控状态
const isMonitoring = ref(false)

// 处理模式切换
const handleModeChange = async () => {
  // 重置查询参数
  currentHost.value = ""
  queryParams.dataFile = ""
  fileBrowser.files = []
  fileBrowser.selectedFile = ""
  fileBrowser.currentPath = mode.value === "local" ? "" : "/opt"

  // 根据新模式初始化数据
  await initData()
}

// 初始化数据
const initData = async () => {
  if (mode.value === "ssh") {
    // 只有在userStore中没有SSH连接信息时才从接口获取
    if (!userStore.sshConnections || userStore.sshConnections.length === 0) {
      await getSSHConnectionsApi()
    }

    // 刷新主机列表
    await refreshHosts()

    // 如果有当前SSH主机，则自动选择
    if (userStore.currentSSHHost) {
      currentHost.value = userStore.currentSSHHost
      await handleHostChange(userStore.currentSSHHost)
    }
  } else {
    // 获取本地日志文件列表
    await loadLocalFiles()
  }
}

// 获取本地日志文件列表
const loadLocalFiles = async () => {
  try {
    fileBrowser.loading = true
    const response = await getLocalLogFilesApi(fileBrowser.currentPath)
    if (!response || !response.data.files) {
      throw new Error("获取本地日志文件列表失败：返回数据格式错误")
    }
    fileBrowser.files = response.data.files.filter((file) => !file.isDir).map((file) => file.path)
    fileBrowser.total = fileBrowser.files.length
  } catch (error) {
    console.error("获取本地日志文件列表失败:", error)
    fileBrowser.files = []
  } finally {
    fileBrowser.loading = false
  }
}

// 显示SSH连接对话框
const showSSHDialog = () => {
  sshDialogVisible.value = true
}

// 显示文件选择对话框
const showFileDialog = () => {
  if (mode.value === "ssh") {
    if (!currentHost.value) {
      ElMessage.warning("请先选择主机")
      return
    }
  }

  fileDialogVisible.value = true
  if (mode.value === "local") {
    // 本地模式下重置当前路径
    fileBrowser.currentPath = ""
    loadLocalFiles()
  } else {
    // SSH模式下使用默认路径
    fileBrowser.currentPath = "/opt"
    loadFiles()
  }
}

// 加载文件列表
const loadFiles = async () => {
  if (mode.value === "local") {
    await loadLocalFiles()
    return
  }

  if (!currentHost.value) {
    ElMessage.warning("请先选择主机")
    return
  }

  try {
    fileBrowser.loading = true
    const files = await getLogFilesApi(currentHost.value, fileBrowser.currentPath)
    fileBrowser.files = files
    fileBrowser.total = files.length
  } catch (error) {
    console.error("获取文件列表失败:", error)
    ElMessage.error("获取文件列表失败")
  } finally {
    fileBrowser.loading = false
  }
}

// 判断是否为目录
const isDirectory = (name: string) => {
  return !name.includes(".")
}

// 获取文件行样式
const getFileRowClass = ({ row }: { row: string }) => {
  return isDirectory(row) ? "directory-row" : ""
}

// 处理文件行点击
const handleFileRowClick = (row: string) => {
  if (isDirectory(row)) {
    fileBrowser.currentPath = `${fileBrowser.currentPath}/${row}`
    loadFiles()
  } else {
    handleFileClick(row)
  }
}

// 处理文件点击
const handleFileClick = (file: string) => {
  fileBrowser.selectedFile = file
  queryParams.dataFile = file
  realtimeConfig.dataFile = file
  historyConfig.dataFile = file
}

// 选择文件
const selectFile = (file: string) => {
  handleFileClick(file)
}

// 确认选择文件
const confirmSelectFile = () => {
  if (fileBrowser.selectedFile) {
    fileDialogVisible.value = false
  }
}

// 处理文件大小变化
const handleFileSizeChange = (size: number) => {
  fileBrowser.pageSize = size
  loadFiles()
}

// 处理文件页码变化
const handleFilePageChange = (page: number) => {
  fileBrowser.currentPage = page
  loadFiles()
}

// 处理主机变更
const handleHostChange = async (host: string) => {
  if (!host) {
    fileBrowser.files = []
    return
  }

  try {
    // 设置当前SSH主机
    userStore.setCurrentSSHHost(host)

    // 获取日志文件列表
    fileBrowser.files = await getLogFilesApi(host, fileBrowser.currentPath)
    fileBrowser.total = fileBrowser.files.length
  } catch (error) {
    console.error("获取日志文件列表失败:", error)
    ElMessage.error("获取日志文件列表失败")
  }
}

// 刷新主机列表
const refreshHosts = async () => {
  try {
    // 获取主机列表
    const hostList = await getHostsApi()

    // 去重处理
    hosts.value = [...new Set(hostList)]
  } catch (error) {
    console.error("获取主机列表失败:", error)
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

    // 设置为当前选中的主机
    currentHost.value = sshForm.host

    // 获取该主机的日志文件列表
    await handleHostChange(sshForm.host)

    sshDialogVisible.value = false
  } catch (error) {
    ElMessage.error(`SSH连接失败: ${error instanceof Error ? error.message : String(error)}`)
  } finally {
    sshConnecting.value = false
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
  { label: "debug_logs", value: "debug_logs" },
  { label: "info_logs", value: "info_logs" },
  { label: "warn_logs", value: "warn_logs" }
]

// 实时监控配置
const realtimeConfig = reactive<RealtimeMonitorParams>({
  metrics: ["requests"],
  timeRange: "1h",
  dataFile: ""
})

// 历史分析配置
const historyConfig = reactive<FrontHistoryAnalysisParams>({
  host: "",
  dataFile: "",
  timeRange: [null, null]
})

// 定时器
let monitorTimer: number | null = null

// 初始化图表
onMounted(() => {
  // 延迟初始化，确保DOM已经渲染完成
  setTimeout(() => {
    if (realtimeChart.value) {
      realtimeChartInstance = echarts.init(realtimeChart.value)
      console.log("实时图表实例初始化完成")
    }
    if (historyChart.value) {
      historyChartInstance = echarts.init(historyChart.value)
      console.log("历史图表实例初始化完成")
    }

    // 监听窗口大小变化
    window.addEventListener("resize", handleResize)

    // 初始化数据
    initData()
  }, 200)
})

onUnmounted(() => {
  stopMonitor()
  window.removeEventListener("resize", handleResize)
  realtimeChartInstance?.dispose()
  historyChartInstance?.dispose()
})

// 窗口大小变化时重绘图表
const handleResize = () => {
  realtimeChartInstance?.resize()
  historyChartInstance?.resize()
}

// 开始实时监控
const startMonitor = async () => {
  if (mode.value === "ssh") {
    if (!currentHost.value) {
      showSSHDialog()
      return
    }
  }

  if (!queryParams.dataFile) {
    showFileDialog()
    return
  }

  if (!realtimeConfig.metrics.length) {
    ElMessage.warning("请选择至少一个监控项")
    return
  }

  stopMonitor()
  await updateRealtimeChart()
  isMonitoring.value = true

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
  isMonitoring.value = false
}

// 更新实时图表
const updateRealtimeChart = async () => {
  try {
    // 确保dataFile是最新的
    realtimeConfig.dataFile = queryParams.dataFile

    // 检查必要的参数
    if (!realtimeConfig.dataFile) {
      console.error("缺少数据文件路径")
      ElMessage.warning("请先选择数据文件")
      return
    }

    if (!realtimeConfig.metrics || realtimeConfig.metrics.length === 0) {
      console.error("缺少监控指标")
      ElMessage.warning("请选择至少一个监控项")
      return
    }

    console.log("请求实时监控数据，参数:", JSON.stringify(realtimeConfig))
    let data

    if (mode.value === "local") {
      // 本地模式使用本地日志分析API
      data = await getLocalLogMetricsApi({
        dataFile: realtimeConfig.dataFile,
        metrics: realtimeConfig.metrics,
        timeRange: realtimeConfig.timeRange
      })
    } else {
      // SSH模式使用远程日志分析API
      data = await getRealtimeMetricsApi(realtimeConfig)
    }

    console.log("获取到实时监控数据:", data)

    // 检查返回的数据
    if (!data || !data.series) {
      console.error("实时监控API返回的数据格式不正确:", data)
      ElMessage.warning("获取实时监控数据失败，返回数据格式不正确")
      return
    }

    // 检查是否有数据
    if (data.series.length === 0) {
      console.warn("实时监控API返回的数据为空")
      ElMessage.info("当前没有实时监控数据")
      return
    }

    // 渲染图表
    renderRealtimeChart(data)
  } catch (error) {
    console.error("更新实时图表失败:", error)
    ElMessage.error(`更新实时图表失败: ${error instanceof Error ? error.message : String(error)}`)
  }
}

// 渲染实时图表
const renderRealtimeChart = (data: any) => {
  if (!realtimeChartInstance) {
    console.error("实时图表实例未初始化")
    return
  }

  try {
    // 验证数据格式
    if (!data || !data.series || !Array.isArray(data.series)) {
      console.error("图表数据格式不正确:", data)
      ElMessage.warning("图表数据格式不正确")
      return
    }

    // 处理数据点，确保时间戳格式正确
    const processedSeries = data.series.map((series: any, index: number) => {
      console.log("series", series)

      console.log("series.data", series.data)

      if (!series.data || !Array.isArray(series.data)) {
        return { ...series, data: [] }
      }

      // 过滤和处理数据点
      const processedData = series.data
        .filter((point: any) => {
          // 过滤无效的时间戳
          if (!Array.isArray(point) || point.length < 2) return false
          const timestamp = point[0]
          return timestamp && timestamp > 1000000 // 过滤掉无效的时间戳
        })
        .map((point: any) => {
          let timestamp = point[0]
          // 检查时间戳是否需要转换（秒转毫秒）
          if (timestamp < 10000000000) {
            timestamp = timestamp * 1000
            console.log(`转换时间戳从秒到毫秒: ${point[0]} -> ${timestamp}`)
          }
          return [timestamp, point[1]]
        })

      // 为错误日志设置固定的红色渐变
      const isErrorLog = series.name.toLowerCase().includes("error")
      const colorList = isErrorLog
        ? [["#f56c6c", "#ffa39e"]] // 错误日志使用红色渐变
        : [
            ["#1890FF", "#91d5ff"],
            ["#2FC25B", "#8effa0"],
            ["#FACC14", "#ffe58f"],
            ["#722ED1", "#d3adf7"],
            ["#F5222D", "#ffa39e"]
          ]
      const colorIndex = isErrorLog ? 0 : index % colorList.length
      const gradientColor = new echarts.graphic.LinearGradient(0, 0, 0, 1, [
        { offset: 0, color: colorList[colorIndex][0] },
        { offset: 1, color: colorList[colorIndex][1] }
      ])

      return {
        ...series,
        data: processedData,
        // 添加渐变色和样式
        itemStyle: {
          color: colorList[colorIndex][0]
        },
        lineStyle: {
          width: 3,
          color: gradientColor
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: colorList[colorIndex][0] + "50" }, // 添加透明度
            { offset: 1, color: colorList[colorIndex][1] + "00" }
          ])
        }
      }
    })

    // 根据选择的时间范围设置X轴范围
    const now = Date.now()
    const timeRangeHours = parseInt(realtimeConfig.timeRange) || 1
    const rangeMs = timeRangeHours * 60 * 60 * 1000
    const minTime = now - rangeMs
    const maxTime = now

    // 格式化时间的函数
    const formatTime = (timestamp: number) => {
      return dayjs(timestamp).format("YYYY-MM-DD HH:mm:ss")
    }

    // 根据时间范围选择合适的时间格式
    const getTimeFormat = () => {
      const timeRangeHours = parseInt(realtimeConfig.timeRange) || 1
      if (timeRangeHours <= 1) {
        return "HH:mm:ss" // 1小时内显示时:分:秒
      } else if (timeRangeHours <= 12) {
        return "MM-DD HH:mm" // 12小时内显示月-日 时:分
      } else {
        return "MM-DD HH:mm" // 24小时显示月-日 时:分
      }
    }

    // 设置图表选项
    const option = {
      title: {
        text: "实时监控",
        left: "center"
      },
      tooltip: {
        trigger: "axis",
        formatter: function (params: any) {
          if (!params || params.length === 0) return ""
          const time = formatTime(params[0].value[0])
          let result = `${time}<br/>`
          params.forEach((param: any) => {
            result += `${param.seriesName}: ${param.value[1]}<br/>`
          })
          return result
        }
      },
      legend: {
        data: processedSeries.map((s: any) => s.name),
        top: 30,
        type: "scroll",
        orient: "horizontal"
      },
      grid: {
        left: "3%",
        right: "4%",
        bottom: "10%",
        containLabel: true
      },
      dataZoom: [
        {
          type: "slider",
          show: true,
          start: 0,
          end: 100,
          bottom: 10
        },
        {
          type: "inside",
          start: 0,
          end: 100
        }
      ],
      xAxis: {
        type: "time",
        boundaryGap: false,
        min: minTime,
        max: maxTime,
        axisLabel: {
          formatter: function (value: number) {
            return dayjs(value).format(getTimeFormat())
          }
        }
      },
      yAxis: {
        type: "value",
        minInterval: 1, // 确保Y轴刻度为整数
        splitLine: {
          show: true
        },
        axisLabel: {
          formatter: function (value: number) {
            return Math.floor(value) // 确保显示为整数
          }
        }
      },
      series: processedSeries.map((series) => ({
        ...series,
        smooth: true,
        showSymbol: false, // 默认不显示数据点
        symbolSize: 8,
        emphasis: {
          focus: "series",
          itemStyle: {
            shadowBlur: 10,
            shadowColor: "rgba(0, 0, 0, 0.3)"
          }
        }
      }))
    }

    // 设置图表选项
    realtimeChartInstance.setOption(option, true)
    console.log("实时图表渲染完成")
  } catch (error) {
    console.error("渲染实时图表失败:", error)
    ElMessage.error(`渲染实时图表失败: ${error instanceof Error ? error.message : String(error)}`)
  }
}

// 查询历史数据
const queryHistory = async () => {
  if (mode.value === "ssh") {
    if (!currentHost.value) {
      ElMessage.warning("请先连接主机")
      showSSHDialog()
      return
    }
  }

  if (!queryParams.dataFile) {
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

    // 构建请求参数
    const params = {
      ...historyConfig,
      host: mode.value === "ssh" ? currentHost.value || "" : ""
    }

    const data = await getHistoryMetricsApi(params)
    console.log("historyData", data)

    // 检查数据结构
    if (!data || !data.series || !Array.isArray(data.series) || data.series.length === 0) {
      return
    }

    // 预处理数据，确保时间戳有效
    if (data.series && Array.isArray(data.series)) {
      data.series = data.series.map((series) => {
        if (series.data && Array.isArray(series.data)) {
          // 过滤掉无效的时间戳数据点
          series.data = series.data.filter((point) => {
            if (Array.isArray(point)) {
              const timestamp = Number(point[0])
              return timestamp !== 0 && timestamp > 1000000 // 过滤掉0或很小的时间戳
            } else if (point && typeof point === "object" && "time" in (point as Record<string, any>)) {
              const timestamp = Number((point as Record<string, any>).time)
              return timestamp !== 0 && timestamp > 1000000 // 过滤掉0或很小的时间戳
            }
            return true
          })
        }
        return series
      })
    }

    // 保存原始数据
    historyData.value = data.series || []

    // 渲染图表
    renderHistoryChart(data)
  } catch (error) {
    console.error("查询历史数据失败:", error)
    ElMessage.error(`查询失败: ${error instanceof Error ? error.message : String(error)}`)
  }
}

// 渲染历史图表
const renderHistoryChart = (data: any) => {
  // 确保图表实例存在
  if (!historyChartInstance) {
    console.error("图表实例未初始化")

    // 尝试重新初始化
    if (historyChart.value) {
      console.log("尝试重新初始化历史图表实例")
      historyChartInstance = echarts.init(historyChart.value)
    } else {
      console.error("历史图表容器不存在")
      return
    }
  }

  console.log("historyData", data.series[0].data)
  // 确保数据格式正确
  if (!data || !data.series || !Array.isArray(data.series) || data.series.length === 0) {
    console.error("历史图表数据格式不正确:", data)
    ElMessage.warning("历史数据格式不正确，无法显示图表")
    return
  }

  console.log("historyData", data)

  try {
    // 检查图表容器
    if (historyChart.value) {
      console.log("图表容器尺寸:", historyChart.value.offsetWidth, historyChart.value.offsetHeight)
      console.log(
        "图表容器可见性:",
        window.getComputedStyle(historyChart.value).display,
        window.getComputedStyle(historyChart.value).visibility
      )

      // 如果容器尺寸为0，设置明确的尺寸
      if (!historyChart.value.offsetWidth || !historyChart.value.offsetHeight) {
        historyChart.value.style.width = "100%"
        // historyChart.value.style.height = "400px"
        // 重新初始化图表
        if (historyChartInstance) {
          historyChartInstance.dispose()
        }
        historyChartInstance = echarts.init(historyChart.value)
      }
    }

    // 处理图表数据
    const metrics = data.metrics || []
    const series = data.series.map((item: any, index: number) => {
      // 确保数据点格式正确
      console.log("item-data", JSON.stringify(item.data))

      const formattedData = Array.isArray(item.data)
        ? item.data
            .map((point: any) => {
              // 如果数据点是数组格式 [timestamp, value]
              if (Array.isArray(point)) {
                // 确保时间戳是毫秒格式
                const timestamp = Number(point[0])
                console.log("原始时间戳:", timestamp, "日期:", new Date(timestamp))

                // 检查时间戳是否为0或接近0的值
                if (timestamp === 0 || timestamp < 1000000) {
                  console.error("检测到无效时间戳:", timestamp)
                  return null // 返回null以便后续过滤
                }

                // 检查时间戳是否是秒格式（10位数字）
                const timestampMs = timestamp < 10000000000 ? timestamp * 1000 : timestamp
                console.log("转换后时间戳:", timestampMs, "日期:", new Date(timestampMs))
                return [timestampMs, Number(point[1])]
              }
              // 如果数据点是对象格式 {time: timestamp, value: value}
              else if (point && typeof point === "object" && "time" in point && "value" in point) {
                // 确保时间戳是毫秒格式
                const timestamp = Number(point.time)
                console.log("原始时间戳(对象):", timestamp, "日期:", new Date(timestamp))

                // 检查时间戳是否为0或接近0的值
                if (timestamp === 0 || timestamp < 1000000) {
                  console.error("检测到无效时间戳(对象):", timestamp)
                  return null // 返回null以便后续过滤
                }

                // 检查时间戳是否是秒格式（10位数字）
                const timestampMs = timestamp < 10000000000 ? timestamp * 1000 : timestamp
                console.log("转换后时间戳(对象):", timestampMs, "日期:", new Date(timestampMs))
                return [timestampMs, Number(point.value)]
              }
              return point
            })
            .filter((point) => point !== null) // 过滤掉无效的时间戳
        : []

      // 为错误日志设置固定的红色渐变
      const isErrorLog = item.name.toLowerCase().includes("error")
      const colorList = isErrorLog
        ? [["#f56c6c", "#ffa39e"]] // 错误日志使用红色渐变
        : [
            ["#1890FF", "#91d5ff"],
            ["#2FC25B", "#8effa0"],
            ["#FACC14", "#ffe58f"],
            ["#722ED1", "#d3adf7"],
            ["#F5222D", "#ffa39e"]
          ]
      const colorIndex = isErrorLog ? 0 : index % colorList.length
      const gradientColor = new echarts.graphic.LinearGradient(0, 0, 0, 1, [
        { offset: 0, color: colorList[colorIndex][0] },
        { offset: 1, color: colorList[colorIndex][1] }
      ])

      return {
        name: item.name,
        type: "line",
        smooth: true,
        showSymbol: true, // 显示数据点
        symbolSize: 8, // 增大数据点大小
        lineStyle: {
          width: 3,
          color: gradientColor
        },
        itemStyle: {
          color: colorList[colorIndex][0]
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: colorList[colorIndex][0] + "50" }, // 添加透明度
            { offset: 1, color: colorList[colorIndex][1] + "00" }
          ])
        },
        emphasis: {
          focus: "series",
          itemStyle: {
            shadowBlur: 10,
            shadowColor: "rgba(0, 0, 0, 0.3)"
          }
        },
        data: formattedData
      }
    })

    console.log("series", JSON.stringify(series))

    // 检查数据范围
    if (series.length > 0 && series[0].data.length > 0) {
      const timestamps = series[0].data.map((point) => point[0])
      const minTime = Math.min(...timestamps)
      const maxTime = Math.max(...timestamps)
      console.log("数据时间范围:", new Date(minTime), new Date(maxTime))
      console.log("时间跨度(天):", (maxTime - minTime) / (1000 * 60 * 60 * 24))
    }

    // 根据时间范围选择合适的时间格式
    const getTimeFormat = () => {
      if (!historyConfig.timeRange[0] || !historyConfig.timeRange[1]) return "MM-DD HH:mm"

      const startTime = new Date(historyConfig.timeRange[0]).getTime()
      const endTime = new Date(historyConfig.timeRange[1]).getTime()
      const diffHours = (endTime - startTime) / (1000 * 60 * 60)

      if (diffHours <= 1) {
        return "HH:mm:ss" // 1小时内显示时:分:秒
      } else if (diffHours <= 24) {
        return "HH:mm" // 24小时内显示时:分
      } else if (diffHours <= 24 * 7) {
        return "MM-DD HH:mm" // 一周内显示月-日 时:分
      } else {
        return "YYYY-MM-DD" // 超过一周显示年-月-日
      }
    }

    setTimeout(() => {
      console.log("尝试渲染实际数据")
      const option = {
        title: {
          text: "历史分析",
          left: "center"
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
          top: 30
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
          },
          axisLabel: {
            formatter: function (value: number) {
              return dayjs(value).format(getTimeFormat())
            }
          },
          min: function (value) {
            // 如果最小值是1970年，则尝试使用数据中的最小时间
            if (new Date(value.min).getFullYear() <= 1970) {
              console.log("检测到无效的x轴最小值:", value.min, new Date(value.min))
              // 尝试从数据中找出有效的最小时间
              if (series.length > 0 && series[0].data.length > 0) {
                const timestamps = series[0].data
                  .map((point) => point[0])
                  .filter((ts) => new Date(ts).getFullYear() > 1970)
                if (timestamps.length > 0) {
                  const minTime = Math.min(...timestamps)
                  console.log("使用数据中的最小时间:", minTime, new Date(minTime))
                  return minTime
                }
              }

              // 如果没有找到有效的时间戳，使用当前时间减去一天作为默认值
              const defaultTime = Date.now() - 24 * 60 * 60 * 1000
              console.log("使用默认最小时间:", defaultTime, new Date(defaultTime))
              return defaultTime
            }
            return value.min
          },
          max: function (value) {
            // 如果最大值是1970年，则尝试使用数据中的最大时间
            if (new Date(value.max).getFullYear() <= 1970) {
              console.log("检测到无效的x轴最大值:", value.max, new Date(value.max))
              // 尝试从数据中找出有效的最大时间
              if (series.length > 0 && series[0].data.length > 0) {
                const timestamps = series[0].data
                  .map((point) => point[0])
                  .filter((ts) => new Date(ts).getFullYear() > 1970)
                if (timestamps.length > 0) {
                  const maxTime = Math.max(...timestamps)
                  console.log("使用数据中的最大时间:", maxTime, new Date(maxTime))
                  return maxTime
                }
              }

              // 如果没有找到有效的时间戳，使用当前时间作为默认值
              const defaultTime = Date.now()
              console.log("使用默认最大时间:", defaultTime, new Date(defaultTime))
              return defaultTime
            }
            return value.max
          }
        },
        yAxis: {
          type: "value",
          minInterval: 1, // 确保Y轴刻度为整数
          splitLine: {
            lineStyle: {
              type: "dashed"
            }
          },
          axisLabel: {
            formatter: function (value: number) {
              return Math.floor(value) // 确保显示为整数
            }
          }
        },
        series: series,
        dataZoom: [
          {
            type: "slider",
            show: true,
            xAxisIndex: [0],
            start: 0,
            end: 100
          },
          {
            type: "inside",
            xAxisIndex: [0],
            start: 0,
            end: 100
          }
        ]
      }

      try {
        // 完全重置图表并设置选项
        if (historyChartInstance) {
          historyChartInstance.clear()
          historyChartInstance.setOption(option, true)
          console.log("历史图表设置选项成功")
        } else {
          console.error("历史图表实例不存在，无法设置选项")
        }

        // 如果有数据，显示第一页
        if (historyData.value.length > 0) {
          // displayHistoryData()
        }
      } catch (error) {
        console.error("历史图表设置选项失败:", error)
      }
    }, 1000)
  } catch (error) {
    console.error("渲染图表时发生错误:", error)
    ElMessage.error("渲染图表失败")
  }
}

// 监听标签页切换
watch(activeTab, (newTab) => {
  console.log("标签页切换到:", newTab)

  // 延迟执行，确保DOM已更新
  setTimeout(() => {
    if (newTab === "history" && historyChart.value) {
      console.log("历史分析标签页被激活，检查图表实例")

      // 如果实例不存在或已被销毁，重新创建
      if (!historyChartInstance) {
        console.log("重新创建历史图表实例")
        historyChartInstance = echarts.init(historyChart.value)
      }

      // 调整大小以适应容器
      historyChartInstance.resize()

      // 如果有数据，重新渲染
      if (historyData.value.length > 0) {
        console.log("重新渲染历史数据")
        const data = {
          metrics: Array.from(new Set(historyData.value.map((item) => item.name))),
          series: historyData.value
        }
        renderHistoryChart(data)
      }
    }

    if (newTab === "realtime" && realtimeChart.value) {
      console.log("实时监控标签页被激活，检查图表实例")

      // 如果实例不存在或已被销毁，重新创建
      if (!realtimeChartInstance) {
        console.log("重新创建实时图表实例")
        realtimeChartInstance = echarts.init(realtimeChart.value)
      }

      // 调整大小以适应容器
      realtimeChartInstance.resize()
    }
  }, 200)
})

// 监听当前主机变化
watch(currentHost, (newHost) => {
  if (newHost) {
    userStore.setCurrentSSHHost(newHost)
    queryParams.host = newHost
  }
})
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
  min-height: 400px;
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
  /* //滚动条隐藏样式 */
  scrollbar-width: none;
  -ms-overflow-style: none;
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

.log-level-error {
  color: #f56c6c !important;
  background-color: rgba(245, 108, 108, 0.1) !important;
}

.log-level-warn {
  color: #e6a23c !important;
  background-color: rgba(230, 162, 60, 0.1) !important;
}

.log-level-info {
  color: #409eff !important;
  background-color: rgba(64, 158, 255, 0.1) !important;
}

.log-level-debug {
  color: #67c23a !important;
  background-color: rgba(103, 194, 58, 0.1) !important;
}

.log-level-trace {
  color: #909399 !important;
  background-color: rgba(144, 147, 153, 0.1) !important;
}
</style>
