import { request, post, get } from "@/utils/network/axios"
import { ContentTypeEnum, RequestEnum } from "@/utils/network/httpEnum"
import type {
  FileUploadRequestData,
  FileUploadResponseData,
  ChunkUploadInitResponseData,
  ChunkUploadInitRequestData,
  ChunkUploadResponseData,
  ChunkUploadRequestData,
  ChunkUploadCompleteResponseData,
  ChunkUploadCompleteRequestData,
  ChunkUploadStatusRequestData,
  ChunkUploadStatusResponseData,
  DownloadUrlRequestData,
  DownloadUrlResponseData
} from "./types/upload"

const prefix = import.meta.env.VITE_APP_BASE_API

// 文件上传相关接口
export const uploadFileApi = {
  // 普通上传
  upload(data: FileUploadRequestData, config?: any) {
    return request<FileUploadResponseData>(
      `${prefix}/upload_service/v1/file/upload`,
      {
        method: RequestEnum.POST,
        headers: {
          "Content-Type": ContentTypeEnum.FORM_DATA
        },
        data,
        ...config
      },
      true
    )
  },

  // 获取下载链接
  getDownloadUrl(data: DownloadUrlRequestData) {
    return request<DownloadUrlResponseData>(
      `${prefix}/upload_service/v1/file/download/url`,
      {
        method: RequestEnum.POST,
        data
      },
      true
    )
  },

  // 初始化分片上传
  initiateMultipart(data: ChunkUploadInitRequestData) {
    return request<ChunkUploadInitResponseData>(
      `${prefix}/upload_service/v1/file/multipart/init`,
      {
        method: RequestEnum.POST,
        params: data
      },
      true
    )
  },

  // 上传分片
  uploadPart(formData: ChunkUploadRequestData) {
    return request<ChunkUploadResponseData>(
      `${prefix}/upload_service/v1/file/multipart/upload`,
      {
        method: RequestEnum.POST,
        headers: {
          "Content-Type": ContentTypeEnum.FORM_DATA
        },
        data: formData
      },
      true
    )
  },

  // 完成分片上传
  completeMultipart(data: ChunkUploadCompleteRequestData) {
    return request<ChunkUploadCompleteResponseData>(
      `${prefix}/upload_service/v1/file/multipart/complete`,
      {
        method: RequestEnum.POST,
        data
      },
      true
    )
  },

  // 查询分片上传状态
  getUploadStatus(params: ChunkUploadStatusRequestData) {
    return request<ChunkUploadStatusResponseData>(
      `${prefix}/upload_service/v1/file/multipart/status`,
      {
        method: RequestEnum.GET,
        params
      },
      true
    )
  }
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
