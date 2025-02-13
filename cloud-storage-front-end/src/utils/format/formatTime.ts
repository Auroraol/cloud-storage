import dayjs from "dayjs"
// 添加时间格式化函数
export const formatTime = (time) => {
  if (!time) return "-"
  return dayjs(time).format("YYYY-MM-DD HH:mm:ss")
}

// 毫秒级时间戳转时间
export const timestampToDateMs = (timestamp: number) => {
  return dayjs(timestamp).format("YYYY-MM-DD HH:mm:ss")
}

// 时间戳转日期
export const timestampToDate = (timestamp: number): string => {
  return new Date(timestamp * 1000).toLocaleString()
}
