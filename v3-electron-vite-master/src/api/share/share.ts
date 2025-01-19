import { request, post, get } from "@/utils/network/axios"
import { RequestEnum } from "@/utils/network/httpEnum"

const prefix = import.meta.env.VITE_APP_BASE_API

/**
 * @description: 回收站文件列表
 * @param {Object} params
 */
export function loadShareListApi(params) {
  return get(`${prefix}/loadShareList/`, params, true)
}

/**
 * @description: 分享文件
 * @param {Object} data
 * @param {boolean} isPublic - 是否公开分享
 */
export function shareFileApi(data, isPublic = false) {
  return post(`${prefix}/shareFile`, { ...data, isPublic }, true) // 直接使用 post 函数
}

/**
 * @description: 取消分享
 * @param {string} ids
 */
export function cancelShareApi(ids) {
  return request(
    `${prefix}/cancelShare/${ids}`,
    {
      method: RequestEnum.DELETE
    },
    true // isNeedToken
  )
}
