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

export type UserInfoResponseData = ApiResponseData<{
  id: number // 用户ID
  username: string // 用户名
  mobile: string // 手机号
  nickname: string // 昵称
  gender: number // 性别，1：男，0：女，默认为1
  avatar: string // 用户头像
  birthday: string // 生日
  email: string // 电子邮箱
  brief: string // 简介|个性签名
  info: string // 新增信息
  roles: string[] // 角色
}>

/** 获取短信验证码的响应数据 */
export type SmsCodeResponseData = ApiResponseData<string>
