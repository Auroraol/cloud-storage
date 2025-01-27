import { request } from "@/utils/service"
import type * as Login from "./types/login"

/** 获取登录验证码 */
export function getLoginCodeApi() {
  return request<Login.LoginCodeResponseData>({
    url: "login/code",
    method: "get"
  })
}

/** 账号密码登录 */
export function loginApi(data: Login.LoginRequestData) {
  return request<Login.LoginResponseData>({
    url: "users/login",
    method: "post",
    data
  })
}

/** 手机验证码登录 */
export function phoneLoginApi(data: Login.PhoneLoginRequestData) {
  return request<Login.LoginResponseData>({
    url: "users/phone-login",
    method: "post",
    data
  })
}

/** 发送手机验证码 */
export function sendSmsCodeApi(mobile: string) {
  return request<Login.SmsCodeResponseData>({
    url: "users/send-code",
    method: "post",
    data: { mobile }
  })
}

/** 获取用户详情 */
export function getUserInfoApi() {
  return request<Login.UserInfoResponseData>({
    url: "users/info",
    method: "get"
  })
}
