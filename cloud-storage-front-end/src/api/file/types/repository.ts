// user_repository.api
export interface UserRepositorySaveRequestData {
  parent_id: number
  repository_id: number
  name: string
}

export type UserRepositorySaveResponseData = ApiResponseData<object>

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

export type UserFolderListRequestData = ApiResponseData<{
  id: number
}>

export type UserFolderListResponseData = ApiResponseData<{
  list: UserFolder[]
}>

type UserFolder = {
  id: number
  name: string
}

export type UserFileNameUpdateRequestData = ApiResponseData<{
  id: number
  name: string
}>

export type UserFileNameUpdateResponseData = ApiResponseData<object>

export type UserFolderCreateRequestData = ApiResponseData<{
  parent_id: number
  name: string
}>

export type UserFolderCreateResponseData = ApiResponseData<{
  id: number
}>

export type UserFileDeleteRequestData = ApiResponseData<{
  id: number
}>

export type UserFileDeleteResponseData = ApiResponseData<object>

export type UserFileMoveRequestData = ApiResponseData<{
  id: number
  parent_id: number
}>

export type UserFileMoveResponseData = ApiResponseData<object>
