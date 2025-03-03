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
  timeRange: [Date | null, Date | null]
  aggregation: string
  dataFile: string
}

// 图表指标选项
export interface ChartMetric {
  label: string
  value: string
}
