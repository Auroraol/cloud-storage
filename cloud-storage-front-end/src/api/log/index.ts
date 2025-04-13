import { request } from "@/utils/network/axios"
import type * as log from "./types/log"
import { RequestEnum } from "@/utils/network/httpEnum"
import type * as frontend from "./types/frontend"

const prefix = import.meta.env.VITE_APP_BASE_API

// 日志相关接口
export const logApi = {
  // SSH连接
  sshConnect(data: log.SSHConnectRequestData) {
    return request<log.SSHConnectResponseData>(
      `${prefix}/log_service/v1/ssh/connect`,
      {
        method: RequestEnum.POST,
        data
      },
      true
    )
  },

  // 获取SSH连接信息列表
  getSSHConnect() {
    return request<log.SshInfoListResp>(
      `${prefix}/log_service/v1/ssh/get`,
      {
        method: RequestEnum.POST,
        data: {}
      },
      true
    )
  },

  // 删除SSH连接信息
  deleteSSHConnect(data: log.DeleteSSHConnectReq) {
    return request<log.DeleteSSHConnectRes>(
      `${prefix}/log_service/v1/ssh/delete`,
      {
        method: RequestEnum.POST,
        data
      },
      true
    )
  },

  // 获取日志文件列表
  getLogFiles(data: log.GetLogFilesReq) {
    return request<log.GetLogFilesRes>(
      `${prefix}/log_service/v1/ssh/readlog`,
      {
        method: RequestEnum.POST,
        data
      },
      true
    )
  },

  // 读取日志文件内容
  readLogContent(data: log.ReadLogFileReq) {
    return request<log.ReadLogFileRes>(
      `${prefix}/log_service/v1/ssh/logfiles`,
      {
        method: RequestEnum.POST,
        data
      },
      true
    )
  },

  // 实时监控(ssh)
  realTimeMonitor(data: log.RealTimeMonitorReq) {
    return request<log.RealTimeMonitorRes>(
      `${prefix}/log_service/v1/monitor/realtime`,
      {
        method: RequestEnum.POST,
        data
      },
      true
    )
  },

  // 实时监控(本地)
  getLocalLogMetrics(data: log.RealTimeMonitorReq) {
    return request<log.RealTimeMonitorRes>(
      `${prefix}/log_service/v1/local/monitor`,
      {
        method: RequestEnum.POST,
        data
      },
      true
    )
  },

  // 历史分析
  historyAnalysis(data: log.HistoryAnalysisReq) {
    return request<log.HistoryAnalysisRes>(
      `${prefix}/log_service/v1/monitor/history`,
      {
        method: RequestEnum.POST,
        data
      },
      true
    )
  },

  // 获取本地日志文件列表
  getLocalLogFiles(data: { path: string }) {
    return request<frontend.LocalLogFilesRes>(
      `${prefix}/log_service/v1/local/files`,
      {
        method: RequestEnum.POST,
        data
      },
      true
    )
  },

  // 读取本地日志文件
  readLocalLogFile(data: frontend.LocalLogQueryParams) {
    return request<frontend.LocalLogReadRes>(
      `${prefix}/log_service/v1/local/read`,
      {
        method: RequestEnum.POST,
        data
      },
      true
    )
  },

  // 获取本地文件内容
  getLocalFileContent(data: frontend.LocalFileContentParams) {
    return request<{
      content: string
      length: number
    }>(
      `${prefix}/log_service/v1/local/content`,
      {
        method: RequestEnum.POST,
        data
      },
      true
    )
  }
}
