import { request, post, get } from "@/utils/network/axios"
import { RequestEnum } from "@/utils/network/httpEnum" // 引入请求枚举

const prefix = import.meta.env.VITE_APP_BASE_API

/**
 * @description: 回收站文件列表
 * @param {Object} params
 */
export function loadDataListApi(params) {
  return get(`${prefix}/loadRecycleList/`, params, true) // 修改为使用 get 方法
}

/**
 * @description: 恢复文件
 * @param {Object} params
 */
export function recoveryFileApi(params) {
  return request(
    `${prefix}/recoverFile/${params}`,
    {
      method: RequestEnum.PUT
    },
    true
  )
}

/**
 * @description: 删除文件
 * @param {Object} params
 */
export function delFileApi(params) {
  return post(`${prefix}/delFile/${params}`, true)
}
