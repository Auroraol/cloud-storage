import { post, get } from "@/utils/network/axios"

const prefix = import.meta.env.VITE_APP_BASE_API

/**
 * 用户注册
 * @param {Object} data
 */
export function register(data) {
  return post(`${prefix}/user/register`, data)
}

/**
 * 账号登录
 * @param {Object} params
 */
export function accountLogin(data) {
  return post(`${prefix}/account/login`, data)
}

/**
 * @description: 获取用户网盘空间
 */
export function getUseSpaceApi() {
  return get(`${prefix}/getUseSpace`, {}, true)
}
