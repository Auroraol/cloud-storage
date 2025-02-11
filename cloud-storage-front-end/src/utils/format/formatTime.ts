import dayjs from 'dayjs'
// 添加时间格式化函数
export const formatTime = (time) => {
  if (!time) return '-'
  return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
}
