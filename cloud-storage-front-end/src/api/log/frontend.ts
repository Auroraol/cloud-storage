import { logApi } from "./index"
import type { LogQueryParams, LogLine, RealtimeMonitorParams, FrontHistoryAnalysisParams } from "./types/frontend"
import type {
  ReadLogFileReq,
  RealTimeMonitorReq,
  HistoryAnalysisReq,
  LogEntry,
  SSHConnectRequestData
} from "./types/log"
import { useUserStoreHook } from "@/store/modules/user"

// 存储SSH连接状态
let sshConnected = false

// SSH连接
export const connectSSHApi = async (host: string, port: number = 22, user: string = "root", password: string = "") => {
  const userStore = useUserStoreHook()

  // 检查是否已经连接到相同主机
  if (sshConnected && userStore.currentSSHHost === host) {
    return true // 已连接到相同主机
  }

  try {
    // 检查是否有保存的连接信息
    const savedConnection = userStore.getSSHConnection(host)

    // 使用保存的连接信息或新提供的信息
    const data: SSHConnectRequestData = {
      host,
      port: savedConnection?.port || port,
      user: savedConnection?.user || user,
      password: savedConnection?.password || password
    }

    const response = await logApi.sshConnect(data)

    if (response.data.success) {
      sshConnected = true

      // 保存连接信息到store
      userStore.addOrUpdateSSHConnection({
        host,
        port: data.port,
        user: data.user,
        password: data.password
      })

      return true
    } else {
      throw new Error(response.data.message || "SSH连接失败")
    }
  } catch (error) {
    sshConnected = false
    throw error
  }
}

// 确保SSH已连接
const ensureConnected = async (host: string) => {
  const userStore = useUserStoreHook()

  if (!sshConnected || userStore.currentSSHHost !== host) {
    // 获取保存的连接信息
    const savedConnection = userStore.getSSHConnection(host)

    if (savedConnection) {
      // 使用保存的连接信息
      await connectSSHApi(savedConnection.host, savedConnection.port, savedConnection.user, savedConnection.password)
    } else {
      // 没有保存的连接信息，使用默认值
      await connectSSHApi(host)
    }
  }
}

// 读取日志内容
export const readLogApi = async (params: LogQueryParams) => {
  // 确保已连接
  await ensureConnected(params.host)

  // 转换参数格式
  const apiParams: ReadLogFileReq = {
    host: params.host,
    path: params.logfile,
    page: params.page,
    page_size: params.pageSize,
    match: params.keyword
  }

  const response = await logApi.readLogContent(apiParams)

  // 转换响应格式
  const lines: LogLine[] = (response.data.contents || []).map((line, index) => {
    return {
      number: (params.page - 1) * params.pageSize + index + 1,
      content: line,
      highlight: params.keyword ? line.includes(params.keyword) : false
    }
  })

  return {
    lines,
    stats: {
      totalLines: response.data.total_lines,
      matchLines: lines.filter((line) => line.highlight).length,
      currentPage: params.page,
      totalPages: Math.ceil(response.data.total_lines / params.pageSize)
    }
  }
}

// 下载日志
export const downloadLogApi = async (params: LogQueryParams) => {
  // 确保已连接
  await ensureConnected(params.host)

  const apiParams: ReadLogFileReq = {
    host: params.host,
    path: params.logfile,
    page: 1,
    page_size: 10000, // 下载更多行
    match: params.keyword
  }

  const response = await logApi.readLogContent(apiParams)
  const content = (response.data.contents || []).join("\n")
  return new Blob([content], { type: "text/plain" })
}

// 获取主机列表
export const getHostsApi = async () => {
  // 从用户存储中获取已连接的主机列表
  const userStore = useUserStoreHook()
  return userStore.getSSHHosts()
}

// 获取日志文件列表
export const getLogFilesApi = async (host: string, path: string = "") => {
  // 确保已连接
  await ensureConnected(host)

  const response = await logApi.getLogFiles({
    host,
    path: path || "/opt/goTest" // 默认路径，可以根据需要修改
  })

  return response.data.files || []
}

// 获取实时监控数据
export const getRealtimeMetricsApi = async (params: RealtimeMonitorParams) => {
  const userStore = useUserStoreHook()
  const host = userStore.currentSSHHost || "localhost"

  // 确保已连接
  await ensureConnected(host)

  // 转换参数
  const timeRangeMap: Record<string, number> = {
    "1h": 1,
    "6h": 6,
    "12h": 12,
    "24h": 24
  }

  const apiParams: RealTimeMonitorReq = {
    host, // 使用当前连接的主机
    log_file: params.dataFile, // 默认日志文件，可以根据需要修改
    time_range: timeRangeMap[params.timeRange] || 1,
    monitor_items: params.metrics
  }

  const response = await logApi.realTimeMonitor(apiParams)

  // 处理响应数据，转换为echarts需要的格式
  const series = params.metrics.map((metric) => {
    return {
      name: metric,
      type: "line",
      data: (response.data.data || [])
        .filter((item) => item.type === metric)
        .map((item) => [item.timestamp, item.value])
    }
  })

  return { series }
}

// 获取历史分析数据
export const getHistoryMetricsApi = async (params: FrontHistoryAnalysisParams) => {
  const userStore = useUserStoreHook()
  const host = userStore.currentSSHHost || "localhost"

  // 确保已连接
  await ensureConnected(host)

  // 转换参数
  const startTime = params.timeRange[0] ? new Date(params.timeRange[0]).getTime() : 0
  const endTime = params.timeRange[1] ? new Date(params.timeRange[1]).getTime() : 0

  const aggregateByMap: Record<string, string> = {
    minute: "按分钟",
    hour: "按小时",
    day: "按天"
  }

  const apiParams: HistoryAnalysisReq = {
    host, // 使用当前连接的主机
    log_file: params.dataFile, // 默认日志文件，可以根据需要修改
    page: 1,
    page_size: 50,
    start_time: startTime,
    end_time: endTime,
    aggregate_by: aggregateByMap[params.aggregation] || "按小时",
    keywords: "" // 添加必需的 keywords 字段
  }

  const response = await logApi.historyAnalysis(apiParams)
  const { data, total, page, page_size, success } = response.data
  console.log("xxxx", data[0], total)

  // 处理响应数据，转换为echarts需要的格式
  const metrics = Array.from(new Set((data || []).map((item: LogEntry) => item.level)))

  const series = metrics.map((metric) => {
    return {
      name: metric,
      type: "line",
      data: (data || []).filter((item: LogEntry) => item.level === metric).map((item: LogEntry) => [item.timestamp, 1]) // 使用计数为1来统计日志数量
    }
  })

  return { metrics, series, total, page, page_size, success }
}
