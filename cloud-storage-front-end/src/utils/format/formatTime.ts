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

// 2025-02-11T16:00:00.000Z 转成时间戳
export const dateToTimestamp = (date: string): number => {
  if (!date) return 0
  return dayjs(date).unix()
}
