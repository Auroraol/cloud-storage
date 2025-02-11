// 格式化文件大小
export const formatFileSize = (size: number): string => {
  if (!size || isNaN(size)) return "0 B"
  const units = ["B", "KB", "MB", "GB", "TB"]
  let index = 0
  let fileSize = size
  while (fileSize >= 1024 && index < units.length - 1) {
    fileSize /= 1024
    index++
  }
  return `${fileSize.toFixed(2)} ${units[index]}`
}
