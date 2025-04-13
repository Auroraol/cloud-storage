// SSH连接参数
export interface SSHConnectRequestData {
  host: string // 主机地址
  port: number // 端口
  user: string // 用户名
  password: string // 密码
  private_key_path?: string // 私钥路径
}

// SSH连接响应
export type SSHConnectResponseData = ApiResponseData<{
  success: boolean // 是否成功
  message: string // 消息
}>

// 实时监控请求
export interface RealTimeMonitorReq {
  host?: string // 主机地址
  log_file: string // 日志文件名
  monitor_items: string[] // 监控项（请求数、错误数、响应时间）
  time_range: number // 时间范围（1小时、6小时、12小时、24小时）
}

// 实时监控响应
export type RealTimeMonitorRes = ApiResponseData<{
  data: MonitorData[] // 监控数据
  total: number // 总数
  success: boolean // 是否成功
}>

// 监控数据
export interface MonitorData {
  timestamp: number // 时间戳
  value: number // 值
  type: string // 类型（请求数、错误数、响应时间）
}

// 历史分析请求
export interface HistoryAnalysisReq {
  host?: string // 主机地址
  log_file: string // 日志文件名
  start_time: number // 开始时间
  end_time: number // 结束时间
  keywords?: string // 关键字
}

// 历史分析响应
export type HistoryAnalysisRes = ApiResponseData<{
  data: LogEntry[] // 日志条目
  total: number // 总数
  success: boolean // 是否成功
}>

// 日志条目
export interface LogEntry {
  timestamp: number // 时间戳
  content: string // 内容
  level: string // 级别
  source: string // 来源
  value: number // 数量
}

// 获取日志文件列表请求
export interface GetLogFilesReq {
  host: string // 主机地址
  path: string // 日志路径
}

// 获取日志文件列表响应
export type GetLogFilesRes = ApiResponseData<{
  files: string[] // 文件列表
  success: boolean // 是否成功
}>

// 读取日志文件请求
export interface ReadLogFileReq {
  host: string // 主机地址
  path: string // 日志路径
  match: string // 匹配字符串
  page: number // 页码
  page_size: number // 每页大小
}

// 读取日志文件响应
export type ReadLogFileRes = ApiResponseData<{
  contents: string[] // 内容
  total_lines: number // 总行数
  page: number // 页码
  page_size: number // 每页大小
  success: boolean // 是否成功
}>

// 删除SSH连接信息请求
export interface DeleteSSHConnectReq {
  ssh_id: number // SSH记录ID
}

// 删除SSH连接信息响应
export type DeleteSSHConnectRes = ApiResponseData<{
  success: boolean // 是否成功
  message: string // 消息
  ssh_id: number // SSH记录ID
}>

// SSH连接详细信息
export interface SshInfoDetailResp {
  user_id: number // 关联用户ID
  ssh_id: number // SSH记录ID
  host: string // 主机地址
  port: number // 端口号
  user: string // 用户名
  password: string // 密码
}

// SSH连接列表响应
export type SshInfoListResp = ApiResponseData<{
  items: SshInfoDetailResp[] // SSH记录列表
}>

// 通用API响应数据类型
interface ApiResponseData<T> {
  code: number
  data: T
  msg: string
}
