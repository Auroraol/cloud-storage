import { request } from "@/utils/service"
import type { LogQueryParams, LogQueryResponse, StatisticsData } from "./types"

/**
 * 获取审计日志列表
 * @param params 查询参数
 */
export function getAuditLogsApi(params: Omit<LogQueryParams, "userId">) {
  return request<LogQueryResponse>({
    url: "/audit/logs",
    method: "post",
    data: params
  })
}

/**
 * 获取用户审计统计数据
 */
export function getAuditStatisticsApi() {
  return request<StatisticsData>({
    url: "/audit/statistics",
    method: "get"
  })
}
