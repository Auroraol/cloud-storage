import { request } from "@/utils/network/axios"
import type * as Repository from "./types/repository"
import { RequestEnum } from "@/utils/network/httpEnum"

const prefix = import.meta.env.VITE_APP_BASE_API

// 用户文件相关接口
export const userFileApi = {
  // 用户文件列表
  getFileList(data: Repository.UserFileListRequestData) {
    return request<Repository.UserFileListResponseData>(
      `${prefix}/user_center/v1/user/file/list`,
      {
        method: RequestEnum.POST,
        data
      },
      true
    )
  },

  // 用户文件夹列表
  getFolderList(data: Repository.UserFolderListRequestData) {
    return request<Repository.UserFolderListResponseData>(
      `${prefix}/user_center/v1/user/folder/list`,
      {
        method: RequestEnum.POST,
        data
      },
      true
    )
  },

  // 用户文件夹创建
  createFolder(data: Repository.UserFolderCreateRequestData) {
    return request<Repository.UserFolderCreateResponseData>(
      `${prefix}/user_center/v1/user/folder/create`,
      {
        method: RequestEnum.POST,
        data
      },
      true
    )
  },

  // 用户文件删除
  deleteFile(data: Repository.UserFileDeleteRequestData) {
    return request<Repository.UserFileDeleteResponseData>(
      `${prefix}/user_center/v1/user/file/delete`,
      {
        method: RequestEnum.POST,
        data
      },
      true
    )
  },

  // 用户文件移动
  moveFile(data: Repository.UserFileMoveRequestData) {
    return request<Repository.UserFileMoveResponseData>(
      `${prefix}/user_center/v1/user/file/move`,
      {
        method: RequestEnum.POST,
        data
      },
      true
    )
  },

  // 用户文件名称修改
  updateFileName(data: Repository.UserFileNameUpdateRequestData) {
    return request<Repository.UserFileNameUpdateResponseData>(
      `${prefix}/user_center/v1/user/file/name/update`,
      {
        method: RequestEnum.POST,
        data
      },
      true
    )
  },

  // 用户文件的关联存储
  saveRepository(data: Repository.UserRepositorySaveRequestData) {
    return request<Repository.UserRepositorySaveResponseData>(
      `${prefix}/user_center/v1/user/repository/save`,
      {
        method: RequestEnum.POST,
        data
      },
      true
    )
  },

  // 获取文件夹大小
  getFolderSize(data: Repository.UserFolderSizeRequestData) {
    return request<Repository.UserFolderSizeResponseData>(
      `${prefix}/user_center/v1/user/folder/size`,
      {
        method: RequestEnum.GET,
        params: data
      },
      true
    )
  }
}
