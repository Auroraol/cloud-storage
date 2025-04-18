import { logApi } from "./index"
import type {
  LogQueryParams,
  LogLine,
  RealtimeMonitorParams,
  FrontHistoryAnalysisParams,
  LocalLogQueryParams,
  LocalFileContentParams
} from "./types/frontend"
import type {
  ReadLogFileReq,
  RealTimeMonitorReq,
  HistoryAnalysisReq,
  LogEntry,
  SSHConnectRequestData,
  DeleteSSHConnectReq
} from "./types/log"
import { useUserStoreHook } from "@/store/modules/user"
import request from "@/utils/request"

// 存储SSH连接状态
let sshConnected = false

// 获取SSH连接列表
export const getSSHConnectionsApi = async () => {
  try {
    const response = await logApi.getSSHConnect()
    const userStore = useUserStoreHook()

    if (response.data.items && response.data.items.length > 0) {
      // 清空现有连接信息
      userStore.clearSSHConnections()

      // 将接口返回的连接信息转换为store中的格式
      const connections = response.data.items.map((item) => ({
        host: item.host,
        port: item.port,
        user: item.user,
        password: item.password
      }))

      // 设置连接信息到store
      userStore.setSSHConnections(connections)
    }

    return response.data.items || []
  } catch (error) {
    console.error("获取SSH连接列表失败:", error)
    return []
  }
}

// 删除SSH连接
export const deleteSSHConnectionApi = async (sshId: number) => {
  try {
    const data: DeleteSSHConnectReq = {
      ssh_id: sshId
    }

    const response = await logApi.deleteSSHConnect(data)
    return response.data
  } catch (error) {
    console.error("删除SSH连接失败:", error)
    throw error
  }
}

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

// 从用户存储中获取已连接的主机列表
export const getHostsApi = async () => {
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
  // console.log("response", response)
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
  let host = params.host // 默认主机

  // 确保已连接
  if (host !== "") {
    host = userStore.currentSSHHost
    await ensureConnected(host)
  }

  // 转换参数
  const startTime = params.timeRange[0] ? new Date(params.timeRange[0]).getTime() / 1000 : 0
  const endTime = params.timeRange[1] ? new Date(params.timeRange[1]).getTime() / 1000 : 0
  const apiParams: HistoryAnalysisReq = {
    host, // 使用当前连接的主机
    log_file: params.dataFile, // 默认日志文件，可以根据需要修改
    start_time: startTime,
    end_time: endTime,
    keywords: "" // 添加必需的 keywords 字段
  }

  const response = await logApi.historyAnalysis(apiParams)
  const { data, total, page, page_size, success } = response.data

  // 处理响应数据，转换为echarts需要的格式
  const metrics = Array.from(new Set((data || []).map((item: LogEntry) => item.level)))

  const series = metrics.map((metric) => {
    return {
      name: metric,
      type: "line",
      data: (data || [])
        .filter((item: LogEntry) => item.level === metric)
        .map((item: LogEntry) => [item.timestamp, item.value]) // 使用计数为1来统计日志数量
    }
  })

  return { metrics, series, total, page, page_size, success }
}

// 获取本地日志文件列表
export const getLocalLogFilesApi = async (path: string) => {
  try {
    console.log("正在获取本地日志文件列表，路径:", path)
    const response = await logApi.getLocalLogFiles({ path })
    console.log("获取本地日志文件列表响应:", response)

    // if (!response) {
    //   throw new Error("获取本地日志文件列表失败：响应为空")
    // }

    // if (!response.files) {
    //   throw new Error("获取本地日志文件列表失败：响应数据格式错误")
    // }

    return response
  } catch (error) {
    console.error("获取本地日志文件列表失败:", error)
    if (error instanceof Error) {
      throw new Error(`获取本地日志文件列表失败: ${error.message}`)
    }
    throw new Error("获取本地日志文件列表失败: 未知错误")
  }
}

// 读取本地日志文件
export const readLocalLogFileApi = async (params: LocalLogQueryParams) => {
  try {
    const response = await logApi.readLocalLogFile(params)
    return response
  } catch (error) {
    console.error("读取本地日志文件失败:", error)
    throw error
  }
}

// 获取本地文件内容
export const getLocalFileContentApi = async (params: LocalFileContentParams) => {
  try {
    const response = await logApi.getLocalFileContent(params)
    return response
  } catch (error) {
    console.error("获取本地文件内容失败:", error)
    throw error
  }
}

// 获取本地日志指标
export const getLocalLogMetricsApi = async (params: { dataFile: string; metrics: string[]; timeRange: string }) => {
  try {
    console.log("正在获取本地日志指标，参数:", params)
    // 转换参数
    const timeRangeMap: Record<string, number> = {
      "1h": 1,
      "6h": 6,
      "12h": 12,
      "24h": 24
    }

    const apiParams: RealTimeMonitorReq = {
      log_file: params.dataFile, // 默认日志文件，可以根据需要修改
      time_range: timeRangeMap[params.timeRange] || 1,
      monitor_items: params.metrics
    }
    const response = await logApi.getLocalLogMetrics(apiParams)
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
  } catch (error) {
    console.error("获取本地日志指标失败:", error)
    throw error
  }
}
