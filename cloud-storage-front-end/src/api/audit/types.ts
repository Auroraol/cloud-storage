/** 操作类型 */
export type OperationType = "upload" | "download" | "delete" | "modify"

/** 日志记录 */
export interface LogRecord {
  /** ID */
  id: string
  /** 操作类型 */
  operationType: OperationType
  /** 文件名 */
  fileName: string
  /** 文件大小 */
  fileSize: number
  /** 操作时间 */
  operationTime: string
  /** IP地址 */
  ip: string
  /** 操作结果 */
  result: boolean
  /** 错误信息 */
  errorMessage?: string
}

/** 日志查询参数 */
export interface LogQueryParams {
  /** 用户ID */
  userId: string
  /** 操作类型 */
  operationType?: OperationType
  /** 开始时间 */
  startTime?: string
  /** 结束时间 */
  endTime?: string
  /** 页码 */
  page: number
  /** 每页条数 */
  pageSize: number
}

/** 日志查询响应 */
export interface LogQueryResponse {
  /** 日志列表 */
  list: LogRecord[]
  /** 总数 */
  total: number
}

/** 统计数据 */
export interface StatisticsData {
  /** 今日总数 */
  todayTotal: number
  /** 上传数 */
  uploadCount: number
  /** 下载数 */
  downloadCount: number
  /** 删除数 */
  deleteCount: number
  /** 修改数 */
  modifyCount: number
  /** 每日趋势 */
  dailyTrend: Array<[string, number]>
}
