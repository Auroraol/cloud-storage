export interface HistoryFileListRequestData {
  page: number
  size: number
}

export type HistoryFileListResponseData = ApiResponseData<{
  history_list: History[]
  total: number
}>

interface History {
  id: string
  file_name: string
  size: number
  status: number //上传状态，0：上传中，1：上传成功，2：上传失败
  update_time: string
  repository_id: number
}

// 上传历史文件记录
export interface UploadHistoryFileRequestData {
  repository_id: number
  file_name: string
  size: number
  status: number
}

export type UploadHistoryFileResponseData = ApiResponseData<object>

// 删除历史文件记录
export interface DeleteHistoryFileRequestData {
  ids: string[]
}

export type DeleteHistoryFileResponseData = ApiResponseData<object>
