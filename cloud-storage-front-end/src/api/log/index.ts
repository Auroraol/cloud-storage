import { request } from "@/utils/service"
import type {
  LogQueryParams,
  LogReadResponse,
  RealtimeMonitorParams,
  HistoryAnalysisParams,
  ChartResponse
} from "./types"

const prefix = "/api/logs"

/**
 * 读取日志内容
 * @param params 查询参数
 */
export function readLogApi(params: LogQueryParams) {
  return request<LogReadResponse>({
    url: `${prefix}/read`,
    method: "post",
    data: params
  })
}

/**
 * 下载日志文件
 * @param params 查询参数
 */
export function downloadLogApi(params: LogQueryParams) {
  return request({
    url: `${prefix}/download`,
    method: "post",
    data: params,
    responseType: "blob"
  })
}

/**
 * 获取实时监控数据
 * @param params 监控参数
 */
export function getRealtimeMetricsApi(params: RealtimeMonitorParams) {
  return request<ChartResponse>({
    url: `${prefix}/metrics/realtime`,
    method: "post",
    data: params
  })
}

/**
 * 获取历史分析数据
 * @param params 分析参数
 */
export function getHistoryMetricsApi(params: HistoryAnalysisParams) {
  return request<ChartResponse>({
    url: `${prefix}/metrics/history`,
    method: "post",
    data: params
  })
}

/**
 * 获取可用主机列表
 */
export function getHostsApi() {
  return request<string[]>({
    url: `${prefix}/hosts`,
    method: "get"
  })
}

/**
 * 获取日志文件列表
 * @param host 主机名
 */
export function getLogFilesApi(host: string) {
  return request<string[]>({
    url: `${prefix}/files`,
    method: "get",
    params: { host }
  })
}
