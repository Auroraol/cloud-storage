/** 文件上传请求数据 */
export interface FileUploadRequestData {
  // 文件表单字段名为 "file"
  file: File
  // 可选的元数据，以JSON字符串形式传递
  metadata?: string
}

// 定义上传响应类型
export type FileUploadResponseData = ApiResponseData<{
  url: string
  key: string
  size: number
  repository_id: number
}>

// 分片上传初始化请求
export interface ChunkUploadInitRequestData {
  file_name: string
  file_size: number
  metadata?: string
}

// 定义分片上传响应类型
export type ChunkUploadInitResponseData = ApiResponseData<{
  upload_id: string
  key: string
}>

// 分片上传请求
export interface ChunkUploadRequestData {
  upload_id: string
  chunk_index: number
  key: string
  file: File
}

// 定义分片上传响应类型
export type ChunkUploadResponseData = ApiResponseData<{
  etag: string
}>

// 完成分片上传请求
export interface ChunkUploadCompleteRequestData {
  upload_id: string
  key: string
  etags: string[]
}

// 定义完成分片上传响应类型
export type ChunkUploadCompleteResponseData = ApiResponseData<{
  url: string
  size: number
}>

// 分片上传状态请求
export interface ChunkUploadStatusRequestData {
  upload_id: string
  key: string
}

// 定义分片上传状态响应类型
interface PartInfo {
  part_number: number // 分片编号
  size: number // 分片大小
  etag: string // 分片ETag
}

// 列出分片上传的分片信息响应
export type ChunkUploadStatusResponseData = ApiResponseData<{
  parts: PartInfo[] // 已上传的分片信息列表
  total_parts: number // 总分片数
  file_size: number // 文件总大小
}>
