export interface LoginRequestData {
  /** admin 或 editor */
  name: string
  /** 密码 */
  password: string
}

export interface PhoneLoginRequestData {
  /** 手机号 */
  mobile: string
  /** 验证码 */
  code: string
}

export type LoginCodeResponseData = ApiResponseData<string>

export type LoginResponseData = ApiResponseData<{
  accesssToken: string
  accessExpire: number
  refreshAfter: number
}>

export type UserInfoResponseData = ApiResponseData<{ username: string; roles: string[] }>

/** 获取短信验证码的响应数据 */
export type SmsCodeResponseData = ApiResponseData<string>
