import { ref } from "vue"
import { pinia } from "@/store"
import { defineStore } from "pinia"
import { useTagsViewStore } from "./tags-view"
import { useSettingsStore } from "./settings"
import { getToken, removeToken, setToken } from "@/utils/cache/session-storage"
import { resetRouter } from "@/router"
import routeSettings from "@/config/route"
// api
import { accountLoginApi, getUserInfoApi, phoneLoginApi } from "@/api/user"
import { type LoginRequestData } from "@/api/user/types/login"
import { da } from "element-plus/es/locale"
import { tokenService } from "@/utils/cache/localStorage-storage"

// SSH连接信息接口
export interface SSHConnectionInfo {
  host: string
  port: number
  user: string
  password: string
}

export const useUserStore = defineStore("user", () => {
  // 全局变量
  // const token = ref<string>(getToken() || "")
  const token = ref<string>(tokenService.getToken() || "")
  const roles = ref<string[]>([])
  const username = ref<string>("")
  const avatar = ref<string>("")
  const userInfo = ref<any>({})
  const capacity = ref<{
    now_volume: number
    total_volume: number
  }>({
    now_volume: 0,
    total_volume: 0
  })

  // SSH连接信息
  const sshConnections = ref<SSHConnectionInfo[]>([])
  const currentSSHHost = ref<string>("")

  const tagsViewStore = useTagsViewStore()
  const settingsStore = useSettingsStore()

  /** 登录 */
  const login = async (req: LoginRequestData) => {
    let dataResult: any = {}
    let messageResult = ""
    // 判断是否为手机号
    if (req.isPhone) {
      const { data, message } = await phoneLoginApi({ mobile: req.mobile, code: req.code })
      dataResult = data
      messageResult = message
    } else {
      const { data, message } = await accountLoginApi({ name: req.name, password: req.password })
      dataResult = data
      messageResult = message
    }
    // console.log(data)
    // 方式1
    // setToken(data.accessToken) //保存在sessionStorage, 也可以使用cookie
    // 方式2 tokenService
    tokenService.setToken(dataResult)
    token.value = dataResult.accessToken
    return messageResult // 添加返回值
  }

  /** 获取用户详情(动态路由) */
  const getInfo = async () => {
    const { data } = await getUserInfoApi()
    console.log(data)
    username.value = data.nickname
    // 验证返回的 roles 是否为一个非空数组，否则塞入一个没有任何作用的默认角色，防止路由守卫逻辑进入无限循环
    roles.value = data.roles?.length > 0 ? data.roles : routeSettings.defaultRoles
    // 头像
    avatar.value = data.avatar
    // 用户详情
    userInfo.value = data
    // 容量
    capacity.value.now_volume = data.now_volume
    capacity.value.total_volume = data.total_volume
  }

  /** 模拟角色变化 */
  const changeRoles = async (role: string) => {
    const newToken = "token-" + role
    token.value = newToken
    setToken(newToken)
    // 用刷新页面代替重新登录
    window.location.reload()
  }

  /** 登出 */
  const logout = () => {
    // 方式1
    // removeToken()
    // 方式2 tokenService
    tokenService.clearToken()
    token.value = ""
    roles.value = []
    avatar.value = ""
    resetRouter()
    _resetTagsView()
  }

  /** 重置 Token */
  const resetToken = () => {
    // 方式1
    // removeToken()
    // 方式2 tokenService
    tokenService.clearToken()
    token.value = ""
    roles.value = []
    avatar.value = ""
  }
  /** 重置 Visited Views 和 Cached Views */
  const _resetTagsView = () => {
    if (!settingsStore.cacheTagsView) {
      tagsViewStore.delAllVisitedViews()
      tagsViewStore.delAllCachedViews()
    }
  }

  /** 添加或更新SSH连接信息 */
  const addOrUpdateSSHConnection = (connection: SSHConnectionInfo) => {
    const index = sshConnections.value.findIndex((conn) => conn.host === connection.host)
    if (index !== -1) {
      // 更新现有连接
      sshConnections.value[index] = connection
    } else {
      // 添加新连接
      sshConnections.value.push(connection)
    }
    // 设置当前连接的主机
    currentSSHHost.value = connection.host
  }

  /** 获取SSH连接信息 */
  const getSSHConnection = (host: string): SSHConnectionInfo | undefined => {
    return sshConnections.value.find((conn) => conn.host === host)
  }

  /** 设置当前SSH主机 */
  const setCurrentSSHHost = (host: string) => {
    currentSSHHost.value = host
  }

  /** 获取所有SSH连接主机 */
  const getSSHHosts = (): string[] => {
    return sshConnections.value.map((conn) => conn.host)
  }

  /** 清空SSH连接列表 */
  const clearSSHConnections = () => {
    sshConnections.value = []
  }

  /** 设置SSH连接列表 */
  const setSSHConnections = (connections: SSHConnectionInfo[]) => {
    sshConnections.value = connections
  }

  return {
    token,
    roles,
    username,
    avatar,
    userInfo,
    capacity,
    sshConnections,
    currentSSHHost,
    login,
    getInfo,
    changeRoles,
    logout,
    resetToken,
    addOrUpdateSSHConnection,
    getSSHConnection,
    setCurrentSSHHost,
    getSSHHosts,
    clearSSHConnections,
    setSSHConnections
  }
})

/**
 * 在 SPA 应用中可用于在 pinia 实例被激活前使用 store
 * 在 SSR 应用中可用于在 setup 外使用 store
 */
export function useUserStoreHook() {
  return useUserStore(pinia)
}
