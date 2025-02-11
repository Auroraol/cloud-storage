// 详情
export interface DetailRequestData {
  id: number
}

export type DetailResponseData = ApiResponseData<{
  repository_id: number
  name: string
  ext: string
  size: number
  path: string
}>

// 创建
export interface ShareBasicCreateRequestData {
  user_repository_id: number
  repository_id: number
  expired_time: number
}

export type ShareBasicCreateResponseData = ApiResponseData<{
  id: number
}>

// 保存
export interface ShareBasicSaveRequestData {
  repository_id: number
  parent_id: number
}

export type ShareBasicSaveResponseData = ApiResponseData<{
  id: number
}>

// 列表
export interface ShareListRequestData {
  page: number
  size: number
}

export type ShareListResponseData = ApiResponseData<{
  list: Share[]
  count: number
}>

type Share = {
  id: number
  file_name: string
  expired_time: number
  update_time: number
  browse_count: number
}
