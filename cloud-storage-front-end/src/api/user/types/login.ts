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
  accessToken: string
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

/** 注册请求数据 */
export interface RegisterRequestData {
  name: string
  password: string
}

/** 注册响应数据 */
export type RegisterResponseData = ApiResponseData<object>

/** 更新用户信息请求数据 */
export interface UserInfoRequestData {
  avatar: string
  nickname: string
  mobile: string
  gender: number
  birthday: string
  email: string
  brief: string
}

/** 更新用户信息响应数据  */
export type UpdateUserInfoResponseData = ApiResponseData<object>

/** 更新头像请求数据 */
export interface UpdateAvatarRequestData {
  file: File // 文件
}

/** 更新头像响应数据 */
export type UpdateAvatarResponseData = ApiResponseData<object>

/** 修改密码请求数据 */
export interface UpdatePasswordRequestData {
  oldPassword: string
  newPassword: string
}

/** 修改密码响应数据 */
export type UpdatePasswordResponseData = ApiResponseData<object>

/** 刷新token响应数据 */
export type RefreshTokenResponseData = ApiResponseData<{
  accessToken: string
  accessExpire: number
  refreshAfter: number
}>
