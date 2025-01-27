/** 日志查询参数 */
export interface LogQueryParams {
  /** 日志文件名 */
  logfile: string
  /** 主机 */
  host: string
  /** 关键字 */
  keyword: string
  /** 页码 */
  page: number
  /** 每页条数 */
  pageSize: number
}

/** 日志行数据 */
export interface LogLine {
  /** 行号 */
  number: number
  /** 内容 */
  content: string
  /** 是否高亮 */
  highlight: boolean
}

/** 日志统计信息 */
export interface LogStats {
  /** 总行数 */
  totalLines: number
  /** 匹配行数 */
  matchLines: number
  /** 当前页 */
  currentPage: number
  /** 总页数 */
  totalPages: number
}

/** 日志读取响应数据 */
export interface LogReadResponse {
  /** 日志行列表 */
  lines: LogLine[]
  /** 统计信息 */
  stats: LogStats
}

/** 图表监控指标 */
export interface ChartMetric {
  /** 指标名称 */
  label: string
  /** 指标值 */
  value: string
}

/** 实时监控参数 */
export interface RealtimeMonitorParams {
  /** 监控指标列表 */
  metrics: string[]
  /** 时间范围 */
  timeRange: string
}

/** 历史分析参数 */
export interface HistoryAnalysisParams {
  /** 时间范围 */
  timeRange: [string, string]
  /** 聚合方式 */
  aggregation: "minute" | "hour" | "day"
}

/** 图表数据系列 */
export interface ChartSeries {
  /** 名称 */
  name: string
  /** 类型 */
  type: string
  /** 数据 */
  data: Array<[number, number]>
}

/** 图表响应数据 */
export interface ChartResponse {
  /** 指标列表 */
  metrics: string[]
  /** 数据系列 */
  series: ChartSeries[]
}
