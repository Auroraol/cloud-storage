import { request, post, get } from "@/utils/network/axios"
import type * as Login from "./types/login"
const prefix = import.meta.env.VITE_APP_BASE_API

/**
 * 用户注册
 * @param {Login.RegisterRequestData} data - 注册信息
 */
export function registerApi(data: Login.RegisterRequestData) {
  return request<Login.RegisterResponseData>(
    `${prefix}/user_center/v1/oauth/register`,
    {
      method: "post",
      data
    },
    true
  )
}

/**
 * 账号登录
 * @param {Login.LoginRequestData} data - 登录信息
 */
export function accountLoginApi(data: Login.LoginRequestData) {
  return request<Login.LoginResponseData>(
    `${prefix}/user_center/v1/oauth/login`,
    {
      method: "post",
      data
    },
    true
  )
}

/** 获取用户详情 */
export function getUserInfoApi() {
  return request<Login.UserInfoResponseData>(
    `${prefix}/user_center/v1/user/detail`,
    {
      method: "post"
    },
    true
  )
}

/**
 * @description: 获取用户网盘空间
 * @return {Promise<SpaceInfo>} 空间使用信息
 */
export function getUseSpaceApi() {
  return get<SpaceInfo>(`${prefix}/getUseSpace`, {}, true)
}
