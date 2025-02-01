import { defineStore } from "pinia"
import { tokenService } from "@/utils/cache/localStorage-storage"

export const useLoginStore = defineStore("login", () => {
  const visible = ref(false) // 是否点击登录
  const username = ref(tokenService.getUsername())
  const password = ref(tokenService.getPassword())

  const changeVisible = (value) => {
    visible.value = value
  }

  const setUsernameValue = (value) => {
    username.value = value
  }

  const setPasswordValue = (value) => {
    password.value = value
  }

  /**
   * 记住用户名和密码
   */
  const setUsernameAndPassword = ({ username: newUsername, password: newPassword }) => {
    username.value = newUsername
    password.value = newPassword
    tokenService.setUsername(newUsername)
    tokenService.setPassword(newPassword)
  }

  /**
   * 清除用户和密码
   */
  const clearUsernameAndPassword = () => {
    username.value = ""
    password.value = ""
    tokenService.removeUsername()
    tokenService.removePassword()
  }

  return {
    visible,
    username,
    password,
    changeVisible,
    setUsernameValue,
    setPasswordValue,
    setUsernameAndPassword,
    clearUsernameAndPassword
  }
})
