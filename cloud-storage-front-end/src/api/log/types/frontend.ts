// 前端组件使用的类型定义

// 日志查询参数
export interface LogQueryParams {
  host: string
  logfile: string
  page: number
  pageSize: number
  keyword: string
  path: string
}

// 日志行
export interface LogLine {
  number: number
  content: string
  highlight: boolean
}

// 日志统计信息
export interface LogStats {
  totalLines: number
  matchLines: number
  currentPage: number
  totalPages: number
}

// 实时监控参数
export interface RealtimeMonitorParams {
  timeRange: string
  metrics: string[]
  dataFile: string
}

// 历史分析参数
export interface FrontHistoryAnalysisParams {
  host?: string
  timeRange: [Date | null, Date | null]
  dataFile: string
}

// 图表指标选项
export interface ChartMetric {
  label: string
  value: string
}

// 本地日志文件信息
export interface LocalFileInfo {
  path: string
  name: string
  size: number
  isDir: boolean
  modTime: string
  extension: string
}

// 本地文件统计信息
export interface LocalFileStat {
  totalFiles: number
  totalDirs: number
  totalSize: number
  logFileCount: number
  recentModified: number
}

export type LocalLogFilesRes = ApiResponseData<{
  files: LocalFileInfo[]
  stat: LocalFileStat
}>

// 本地日志查询参数
export interface LocalLogQueryParams {
  path: string
  startTime?: string
  endTime?: string
  level?: string
  keyword?: string
  maxResults?: number
}

// 本地文件内容查询参数
export interface LocalFileContentParams {
  path: string
  offset?: number
  limit?: number
}

export interface LogEntry {
  timestamp: string
  level: string
  content: string
  source: string
  lineNum: number
}

export type LocalLogReadRes = ApiResponseData<{
  entries: LogEntry[]
  total: number
}>

// 本地文件尾部查询参数
export interface LocalFileTailParams {
  filePath: string
  lines?: number
}
