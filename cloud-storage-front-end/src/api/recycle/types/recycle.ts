// 用户回收站列表
export interface UserRecycleListRequest {
  id: number
  page: number
  size: number
}

export type UserRecycleListResponse = ApiResponseData<{
  list: UserRecycleFile[]
  total: number
}>

export type UserRecycleFile = {
  id: number
  repository_id: number
  name: string
  ext: string
  path: string
  size: number
  update_time: number
}

// 用户回收站文件删除
export interface UserRecycleDeleteRequest {
  id: number
}

export type UserRecycleDeleteResponse = ApiResponseData<{
  success: boolean
}>

// 用户回收站文件恢复
export interface UserRecycleRestoreRequest {
  id: number
}

export type UserRecycleRestoreResponse = ApiResponseData<{
  success: boolean
}>
