import { request } from "@/utils/network/axios"
import { RequestEnum } from "@/utils/network/httpEnum"
import type * as Audit from "./types/audit"
const prefix = import.meta.env.VITE_APP_BASE_API

// 审计相关接口
export const auditApi = {
  // 获取审计列表
  getAuditList(data: Audit.GetOperationLogReq) {
    return request<Audit.GetOperationLogRes>(`${prefix}/log_service/v1/operation`, {
      method: RequestEnum.POST,
      data: data
    })
  },
  // 获取审计列表
  getAuditStatistics(data) {
    return request<Audit.StatisticsData>(`${prefix}/log_service/v1/operation`, {
      method: RequestEnum.POST,
      data: data
    })
  }
}
