// user_repository.api
export interface UserRepositorySaveRequestData {
  parent_id: number
  repository_id: number
  name: string
}

export type UserRepositorySaveResponseData = ApiResponseData<object>

export interface UserFileAndFolderListRequestData {
  id: number //查询的文件夹id
  page: number //查询的第几页
  size: number //每页页数
}

export type UserFileAndFolderListResponseData = ApiResponseData<{
  list: UserFile[]
  count: number
}>

export interface UserFileListRequestData {
  id: number //查询的文件夹id
  page: number //查询的第几页
  size: number //每页页数
}

export type UserFileListResponseData = ApiResponseData<{
  list: UserFile[]
  count: number
}>

type UserFile = {
  id: number
  repository_id: number
  name: string
  ext: string
  path: string
  size: number
  update_time: number
}

export interface UserFolderListRequestData {
  id: number
}

export type UserFolderListResponseData = ApiResponseData<{
  list: UserFolder[]
}>

type UserFolder = {
  id: number
  name: string
  update_time: number
}

export interface UserFileNameUpdateRequestData {
  id: number
  name: string
}

export type UserFileNameUpdateResponseData = ApiResponseData<object>

export interface UserFolderCreateRequestData {
  parent_id: number
  name: string
}

export type UserFolderCreateResponseData = ApiResponseData<{
  id: number
}>

export interface UserFileDeleteRequestData {
  id: number
}

export type UserFileDeleteResponseData = ApiResponseData<object>

export interface UserFileMoveRequestData {
  id: number
  parent_id: number // 父文件夹ID
}

export type UserFileMoveResponseData = ApiResponseData<object>

export interface UserFolderSizeRequestData {
  id: number
}

export type UserFolderSizeResponseData = ApiResponseData<{
  size: number
}>

export interface GetFileListParams {
  id: number
  page: number
  size: number
  keyword?: string
}

export interface GetFolderListParams {
  id: number
  keyword?: string
}
