import axios from "axios"
import type { AxiosInstance, AxiosRequestConfig, AxiosResponse } from "axios"
import { ElMessage } from "element-plus"

// 创建 axios 实例
const service: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_URL_4,
  timeout: 15000
})

// 请求拦截器
service.interceptors.request.use(
  (config) => {
    return config
  },
  (error) => {
    console.error("请求错误:", error)
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  (response: AxiosResponse) => {
    const res = response.data
    if (res.code !== 200) {
      ElMessage.error(res.message || "请求失败")
      return Promise.reject(new Error(res.message || "请求失败"))
    }
    return res
  },
  (error) => {
    console.error("响应错误:", error)
    ElMessage.error(error.message || "请求失败")
    return Promise.reject(error)
  }
)

// 封装请求方法
const request = <T = any>(config: AxiosRequestConfig): Promise<T> => {
  return service.request(config)
}

export default request
