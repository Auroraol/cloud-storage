import { request } from "@/utils/network/axios"
import type * as Share from "./types/share"
import { RequestEnum } from "@/utils/network/httpEnum"

const prefix = import.meta.env.VITE_APP_BASE_API

// 分享相关接口
export const shareApi = {
  // 获取分享列表
  getShareList(params: Share.ShareListRequestData) {
    return request<Share.ShareListResponseData>(`${prefix}/share_service/v1/share/basic/list`, {
      method: RequestEnum.GET,
      params
    })
  },

  // 详情
  getShareDetail(data: Share.DetailRequestData) {
    return request<Share.DetailResponseData>(`${prefix}/share_service/v1/share/basic/detail`, {
      method: RequestEnum.GET,
      data
    })
  },

  // 创建分享
  createShare(data: Share.ShareBasicCreateRequestData) {
    return request<Share.ShareBasicCreateResponseData>(`${prefix}/share_service/v1/share/basic/create`, {
      method: RequestEnum.POST,
      data
    })
  },

  // 保存分享
  saveShare(data: Share.ShareBasicSaveRequestData) {
    return request<Share.ShareBasicSaveResponseData>(`${prefix}/share_service/v1/share/basic/save`, {
      method: RequestEnum.POST,
      data
    })
  },

  // 删除分享
  deleteShare(params: Share.ShareBasicDeleteRequestData) {
    return request<Share.ShareBasicDeleteResponseData>(`${prefix}/share_service/v1/share/basic/delete`, {
      method: RequestEnum.POST,
      params
    })
  }
}
