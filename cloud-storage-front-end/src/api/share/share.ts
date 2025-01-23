import type { AxiosResponse } from "axios"
// import { request } from "@/utils/request"

// export interface ShareParams {
//   page?: number
//   pageSize?: number
//   keyword?: string
// }

// export interface ShareData {
//   id: string
//   filename: string
//   fileType: number
//   folderType: number
//   status: number
//   fileCover?: string
//   createTime: string
//   expireTime: string
//   code?: string
//   browseCount: number
//   saveCount: number
//   downloadCount: number
// }

// /**
//  * 获取分享列表
//  */
// export function getShareList(params: ShareParams): Promise<AxiosResponse> {
//   return request({
//     url: "/share/list",
//     method: "GET",
//     params
//   })
// }

// /**
//  * 创建分享
//  */
// export function createShare(data: {
//   fileIds: string[]
//   isPublic?: boolean
//   expireTime?: number
// }): Promise<AxiosResponse> {
//   return request({
//     url: "/share/create",
//     method: "POST",
//     data
//   })
// }

// /**
//  * 取消分享
//  */
// export function cancelShare(shareIds: string[]): Promise<AxiosResponse> {
//   return request({
//     url: "/share/cancel",
//     method: "POST",
//     data: { shareIds }
//   })
// }
