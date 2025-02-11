// 详情
export interface DetailRequestData {
  id: string
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
  code: string
}

export type ShareBasicCreateResponseData = ApiResponseData<{
  id: string
}>

// 保存
export interface ShareBasicSaveRequestData {
  repository_id: number
  parent_id: number
}

export type ShareBasicSaveResponseData = ApiResponseData<{
  id: string
}>

// 列表
export interface ShareListRequestData {
  page: number
  page_size: number
}

export type ShareListResponseData = ApiResponseData<{
  list: Share[]
  count: number
}>

type Share = {
  id: string
  repository_id: number
  name: string
  ext: string
  size: number
  path: string
  expired_time: number
  update_time: number
  //
  owner: string
  avatar: string
  //
  click_num: number // 浏览次数
  //
  code: string
}

// 删除
export interface ShareBasicDeleteRequestData {
  id: string
}

export type ShareBasicDeleteResponseData = ApiResponseData<{}>
