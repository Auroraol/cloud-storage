import CacheKey from "@/constants/cache-key"
import { refreshTokenApi } from "@/api/user"
interface TokenInfo {
  accessToken: string
  accessExpire: number
  refreshAfter: number
}

class TokenService {
  private refreshTimeout: NodeJS.Timeout | null = null

  // 获取token
  getToken() {
    return localStorage.getItem(CacheKey.TOKEN)
  }

  // 设置token信息
  setToken(tokenInfo: TokenInfo) {
    const { accessToken, refreshAfter } = tokenInfo
    // console.log("accessToken: ", accessToken)
    // console.log("refreshAfter: ", refreshAfter)
    localStorage.setItem(CacheKey.TOKEN, accessToken)

    // 计算需要等待的时间（毫秒）
    const refreshAfterMs = refreshAfter * 1000
    const currentTimeMs = Date.now()
    const waitTime = refreshAfterMs - currentTimeMs
    // console.log("waitTime: ", waitTime)
    // 清除之前的定时器
    if (this.refreshTimeout) {
      clearTimeout(this.refreshTimeout)
    }

    // 设置新的定时器
    this.refreshTimeout = setTimeout(() => {
      this.refreshToken()
    }, waitTime)
  }

  // 刷新token
  private async refreshToken() {
    try {
      const { data } = await refreshTokenApi()

      const tokenInfo: TokenInfo = data

      // 更新token信息并设置下一次刷新
      this.setToken(tokenInfo)
    } catch (error) {
      console.error("Token refresh failed:", error)
      // 处理刷新失败的情况
    }
  }

  // 清除token和定时器
  clearToken() {
    localStorage.removeItem(CacheKey.TOKEN)
    if (this.refreshTimeout) {
      clearTimeout(this.refreshTimeout)
      this.refreshTimeout = null
    }
  }

  /**
   * 设置用户名
   * @param {String} username
   */
  setUsername(username: string) {
    localStorage.setItem(CacheKey.USERNAME_KEY, username)
  }

  /**
   * 获取用户名
   */
  getUsername() {
    return localStorage.getItem(CacheKey.USERNAME_KEY)
  }

  /**
   * 移除用户名
   */
  removeUsername() {
    localStorage.removeItem(CacheKey.USERNAME_KEY)
  }

  /**
   * 设置密码
   * @param {String} password
   */
  setPassword(password: string) {
    localStorage.setItem(CacheKey.PASSWORD_KEY, password)
  }

  /**
   * 获取密码
   */
  getPassword() {
    return localStorage.getItem(CacheKey.PASSWORD_KEY)
  }

  /**
   * 移除密码
   */
  removePassword() {
    localStorage.removeItem(CacheKey.PASSWORD_KEY)
  }

  /**
   * 设置记住密码
   */
  setRemember(checked: boolean) {
    localStorage.setItem(CacheKey.REMEMBER_KEY, JSON.stringify(checked))
  }

  /**
   * 获取记住密码
   */
  getRemember() {
    const value = localStorage.getItem(CacheKey.REMEMBER_KEY)
    return value ? JSON.parse(value) : null
  }
}

export const tokenService = new TokenService()
