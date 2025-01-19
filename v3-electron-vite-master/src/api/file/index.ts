import { request, post, get } from "@/utils/network/axios"
import { ContentTypeEnum, RequestEnum } from "@/utils/network/httpEnum"

const prefix = import.meta.env.VITE_APP_BASE_API

/**
 * @description: 文件列表
 * @param {Object} params
 */
export function loadDataListApi(params) {
  return get(`${prefix}/loadDataList/`, params, true)
}

/**
 * @description: 上传文件
 * @param {Object} formData
 * @param {Object} config
 */
export function uploadFileApi(formData, config) {
  return request(
    `${prefix}/uploadFile`,
    {
      method: RequestEnum.POST,
      headers: {
        "Content-Type": ContentTypeEnum.FORM_DATA
      },
      data: formData,
      ...config
    },
    true
  )
}

/**
 * @description: 新建目录
 * @param {Object} data
 */
export function newFolderApi(data) {
  return post(`${prefix}/newFolder`, data, true)
}

/**
 * @description: 获取目录信息
 * @param {Object} params
 */
export function getFolderInfoApi(params) {
  return get(`${prefix}/getFolderInfo`, params, true)
}

/**
 * @description: 修改文件名
 * @param {Object} data
 */
export function renameApi(data) {
  return request(
    `${prefix}/rename`,
    {
      method: RequestEnum.PUT,
      data
    },
    true
  )
}

/**
 * @description: 获取所有目录
 * @param {Object} params
 */
export function loadAllFolderApi(params) {
  return get(`${prefix}/loadAllFolder`, params, true)
}

/**
 * @description: 移动文件
 * @param {Object} data
 */
export function changeFileFolderApi(data) {
  return request(
    `${prefix}/changeFileFolder`,
    {
      method: RequestEnum.PUT,
      data
    },
    true
  )
}

/**
 * @description: 回收文件
 * @param {Object} params
 */
export function delFileApi(params) {
  return request(
    `${prefix}/delFile/${params}`,
    {
      method: RequestEnum.DELETE
    },
    true
  )
}

/**
 * @description: 活动文件流
 * @param {string} url
 */
export function getFileBolbApi(url) {
  return get(
    url,
    {
      responseType: "arraybuffer" // 因为是流文件，所以要指定blob类型
    },
    true
  )
}

/**
 * @description: 创建下载链接
 * @param {string} url
 */
export function createDownloadUrlApi(url) {
  return get(url, {}, true)
}
