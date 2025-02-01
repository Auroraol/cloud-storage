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

/** 刷新token */
export function refreshTokenApi() {
  return request<Login.RefreshTokenResponseData>(
    `${prefix}/user_center/v1/refresh/authorization`,
    {
      method: "post"
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
/** 更新用户信息 */
export function updateUserInfoApi(data: Login.UserInfoRequestData) {
  return request<Login.UpdateUserInfoResponseData>(
    `${prefix}/user_center/v1/user/info/update`,
    {
      method: "post",
      data
    },
    true
  )
}

/** 修改密码 */
export function updatePasswordApi(data: Login.UpdatePasswordRequestData) {
  return request<Login.UpdatePasswordResponseData>(
    `${prefix}/user_center/v1/user/password/update`,
    {
      method: "post",
      data
    },
    true
  )
}

/** 更新头像 */
export function updateAvatarApi(data: Login.UpdateAvatarRequestData) {
  return request<Login.UpdateAvatarResponseData>(
    `${prefix}/user_center/v1/user/avatar/update`,
    {
      method: "post",
      data
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
