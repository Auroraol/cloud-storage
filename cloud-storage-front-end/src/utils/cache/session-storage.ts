/** 统一处理 sessionStorage */

import CacheKey from "@/constants/cache-key"

export const getToken = () => {
  return sessionStorage.getItem(CacheKey.TOKEN)
}
export const setToken = (token: string) => {
  // 临时
  console.log("token: ", token)
  sessionStorage.setItem(CacheKey.TOKEN, token)
}
export const removeToken = () => {
  sessionStorage.removeItem(CacheKey.TOKEN)
}
