import { request } from "@/utils/network/axios"
import { RequestEnum } from "@/utils/network/httpEnum"

const prefix = import.meta.env.VITE_APP_BASE_API

import type * as History from "./types/history"
import { pa } from "element-plus/es/locale"

// 历史文件相关接口
export const historyFileApi = {
  // 获取历史文件列表
  getHistoryFileList(data: History.HistoryFileListRequestData) {
    return request<History.HistoryFileListResponseData>(
      `${prefix}/upload_service/v1/file/history/list`,
      {
        method: RequestEnum.POST,
        data
      },
      true
    )
  },

  // 上传历史文件记录
  uploadHistoryFile(data: History.UploadHistoryFileRequestData) {
    return request<History.UploadHistoryFileResponseData>(
      `${prefix}/upload_service/v1/file/history/update`,
      {
        method: RequestEnum.POST,
        data
      },
      true
    )
  },

  // 删除历史文件记录
  deleteHistoryFile(data: History.DeleteHistoryFileRequestData) {
    return request<History.DeleteHistoryFileResponseData>(
      `${prefix}/upload_service/v1/file/history/delete/all`,
      {
        method: RequestEnum.POST,
        data
      },
      true
    )
  }
}

/**
 * 删除指定的历史记录
 * @param ids 要删除的历史记录ID数组
 */
export const deleteHistoryRecords = (ids: number[]) => {
  return request.delete("/api/file/history", { data: { ids } })
}
