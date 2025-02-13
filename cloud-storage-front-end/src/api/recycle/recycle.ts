import { request } from "@/utils/network/axios"
import type * as Recycle from "./types/recycle"
import { RequestEnum } from "@/utils/network/httpEnum"

const prefix = import.meta.env.VITE_APP_BASE_API

// 回收站相关接口
export const recycleApi = {
  // 用户回收站列表
  getRecycleList(data: Recycle.UserRecycleListRequest) {
    return request<Recycle.UserRecycleListResponse>(
      `${prefix}/user_center/v1/user/recycle/list`,
      {
        method: RequestEnum.POST,
        data
      },
      true
    )
  },

  // 用户回收站文件删除
  deleteRecycle(data: Recycle.UserRecycleDeleteRequest) {
    return request<Recycle.UserRecycleDeleteResponse>(`${prefix}/user_center/v1/user/recycle/delete`, {
      method: RequestEnum.POST,
      data
    })
  },

  // 用户回收站文件恢复
  restoreRecycle(data: Recycle.UserRecycleRestoreRequest) {
    return request<Recycle.UserRecycleRestoreResponse>(`${prefix}/user_center/v1/user/recycle/restore`, {
      method: RequestEnum.POST,
      data
    })
  }
}
