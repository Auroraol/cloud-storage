<script lang="ts" setup>
import { reactive, ref, onMounted } from "vue"
import { useRouter } from "vue-router"
import { useUserStore } from "@/store/modules/user"
import { type FormInstance, type FormRules } from "element-plus"
import { User, Lock } from "@element-plus/icons-vue"
import ThemeSwitch from "@/components/ThemeSwitch/index.vue"
import Owl from "./components/Owl.vue"
import { useFocus } from "./hooks/useFocus"
import SvgIcon from "@/components/SvgIcon/index.vue"

import { registerApi } from "@/api/user"
import { type LoginRequestData, type RegisterRequestData } from "@/api/user/types/login"

import { ElMessage } from "element-plus"

const router = useRouter()
const { isFocus, handleBlur, handleFocus } = useFocus()

// pinia
const useUserPinia = useUserStore()

/** 表单元素的引用 */
const loginFormRef = ref<FormInstance | null>(null)
const phoneFormRef = ref<FormInstance | null>(null)
const registerFormRef = ref<FormInstance | null>(null)
const registerPhoneFormRef = ref<FormInstance | null>(null)

/** 当前激活的标签页 */
const activeTab = ref("account")

/** 登录按钮 Loading */
const loading = ref(false)
/** 发送验证码按钮状态 */
const sendCodeLoading = ref(false)
const countdown = ref(0)

/** 账号密码登录表单数据 */
const loginFormData: LoginRequestData = reactive({
  name: "lfj",
  password: "123456"
})

/** 手机号登录表单数据 */
const phoneFormData = reactive({
  mobile: "",
  code: ""
})

/** 账号密码登录表单校验规则 */
const loginFormRules: FormRules = {
  name: [{ required: true, message: "请输入用户名", trigger: "blur" }],
  password: [
    { required: true, message: "请输入密码", trigger: "blur" },
    { min: 6, max: 16, message: "长度在 6 到 16 个字符", trigger: "blur" }
  ]
}

/** 手机号登录表单校验规则 */
const phoneFormRules: FormRules = {
  mobile: [
    { required: true, message: "请输入手机号", trigger: "blur" },
    { pattern: /^1[3-9]\d{9}$/, message: "请输入正确的手机号", trigger: "blur" }
  ],
  code: [
    { required: true, message: "请输入验证码", trigger: "blur" },
    { len: 6, message: "验证码长度应为6位", trigger: "blur" }
  ]
}

/** 记住密码选项 */
const rememberPassword = ref(false)
/** 记住手机号选项 */
const rememberPhone = ref(false)

/** 登录逻辑 */
const handleLogin = async () => {
  const valid = await loginFormRef.value?.validate()
  if (valid) {
    loading.value = true
    try {
      // Pinia出来登录
      await useUserPinia.login(loginFormData)
      // 如果选择记住密码，可以在这里保存登录信息
      // if (rememberPassword.value) {
      //   localStorage.setItem("rememberedUsername", loginFormData.username)
      //   // 注意：在实际应用中，不建议直接存储密码，应该使用更安全的方式
      // }
      router.push({ path: "/" })
    } catch {
      loginFormData.password = ""
    } finally {
      loading.value = false
    }
  } else {
    console.error("表单校验不通过")
  }
}

/** 手机号登录逻辑 */
const handlePhoneLogin = async () => {
  const valid = await phoneFormRef.value?.validate()
  if (valid) {
    loading.value = true
    try {
      // Pinia
      useUserPinia.login(loginFormData)
      await phoneLoginApi(phoneFormData)
      // 如果选择记住手机号，可以在这里保存手机号
      if (rememberPhone.value) {
        localStorage.setItem("rememberedPhone", phoneFormData.mobile)
      }
      router.push({ path: "/" })
    } catch {
      phoneFormData.code = ""
    } finally {
      loading.value = false
    }
  } else {
    console.error("表单校验不通过")
  }
}

/** 发送验证码 */
const handleSendCode = async () => {
  try {
    await phoneFormRef.value?.validateField("mobile")
    sendCodeLoading.value = true
    await sendSmsCodeApi(phoneFormData.mobile)
    countdown.value = 60
    const timer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0) {
        clearInterval(timer)
      }
    }, 1000)
  } catch (error) {
    console.error("发送验证码失败", error)
  } finally {
    sendCodeLoading.value = false
  }
}

const handleForgetPassword = () => {
  router.push({ path: "/forget-password" })
}

// 在组件挂载时读取保存的登录信息
onMounted(() => {
  const rememberedUsername = localStorage.getItem("rememberedUsername")
  const rememberedPhone = localStorage.getItem("rememberedPhone")

  if (rememberedUsername) {
    loginFormData.name = rememberedUsername
    rememberPassword.value = true
  }

  if (rememberedPhone) {
    phoneFormData.mobile = rememberedPhone
    rememberPhone.value = true
  }
})

// 添加注册相关的响应式变量
const registerDialogVisible = ref(false)
const registerActiveTab = ref("account")
const registerLoading = ref(false)

// 注册表单数据
const registerFormData = reactive({
  name: "",
  password: "",
  confirmPassword: ""
})

// 手机注册表单数据
const registerPhoneFormData = reactive({
  mobile: "",
  code: "",
  password: ""
})

// 注册表单校验规则
const registerFormRules: FormRules = {
  name: [
    { required: true, message: "请输入用户名", trigger: "blur" },
    { min: 2, max: 20, message: "长度在 2 到 20 个字符", trigger: "blur" }
  ],
  password: [
    { required: true, message: "请输入密码", trigger: "blur" },
    { min: 6, max: 16, message: "长度在 6 到 16 个字符", trigger: "blur" }
  ],
  confirmPassword: [
    { required: true, message: "请确认密码", trigger: "blur" },
    {
      validator: (rule: any, value: string, callback: Function) => {
        if (value !== registerFormData.password) {
          callback(new Error("两次输入密码不一致"))
        } else {
          callback()
        }
      },
      trigger: "blur"
    }
  ]
}

// 手机注册表单校验规则
const registerPhoneFormRules: FormRules = {
  mobile: [
    { required: true, message: "请输入手机号", trigger: "blur" },
    { pattern: /^1[3-9]\d{9}$/, message: "请输入正确的手机号", trigger: "blur" }
  ],
  code: [
    { required: true, message: "请输入验证码", trigger: "blur" },
    { len: 6, message: "验证码长度应为6位", trigger: "blur" }
  ],
  password: [
    { required: true, message: "请输入密码", trigger: "blur" },
    { min: 6, max: 16, message: "长度在 6 到 16 个字符", trigger: "blur" }
  ]
}

// 处理注册
const handleRegister = async () => {
  if (registerActiveTab.value === "account") {
    const valid = await registerFormRef.value?.validate()
    if (valid) {
      registerLoading.value = true
      try {
        console.log("registerFormData", registerFormData)
        const data: RegisterRequestData = {
          name: registerFormData.name,
          password: registerFormData.password
        }
        await registerApi(data)
        ElMessage.success("注册成功")
        // 注册成功后，直接跳转登录
        const loginData: LoginRequestData = {
          name: registerFormData.name,
          password: registerFormData.password
        }
        await useUserPinia.login(loginData)
        router.push({ path: "/" })
        registerDialogVisible.value = false
      } catch (error) {
        console.error("注册失败:", error)
      } finally {
        registerLoading.value = false
      }
    }
  } else {
    const valid = await registerPhoneFormRef.value?.validate()
    if (valid) {
      registerLoading.value = true
      try {
        console.log(registerPhoneFormData)
        const data: RegisterRequestData = {
          name: registerPhoneFormData.mobile,
          password: registerPhoneFormData.password
        }
        await registerApi(data)
        ElMessage.success("注册成功")
        registerDialogVisible.value = false
      } catch (error) {
        console.error("注册失败:", error)
      } finally {
        registerLoading.value = false
      }
    }
  }
}
</script>

<template>
  <div class="login-container">
    <ThemeSwitch class="theme-switch" />
    <Owl :close-eyes="isFocus" />
    <div class="login-card">
      <div class="content">
        <el-tabs v-model="activeTab" stretch>
          <el-tab-pane label="账号密码登录" name="account">
            <el-form ref="loginFormRef" :model="loginFormData" :rules="loginFormRules" @keyup.enter="handleLogin">
              <el-form-item prop="username">
                <el-input
                  v-model.trim="loginFormData.name"
                  placeholder="用户名"
                  type="text"
                  tabindex="1"
                  :prefix-icon="User"
                  size="large"
                />
              </el-form-item>
              <el-form-item prop="password">
                <el-input
                  v-model.trim="loginFormData.password"
                  placeholder="密码"
                  type="password"
                  tabindex="2"
                  :prefix-icon="Lock"
                  size="large"
                  show-password
                  @blur="handleBlur"
                  @focus="handleFocus"
                />
              </el-form-item>
              <div class="form-footer">
                <el-checkbox v-model="rememberPassword">记住密码</el-checkbox>
                <el-link type="primary" @click="handleForgetPassword" :underline="false">忘记密码？</el-link>
              </div>
              <el-button :loading="loading" type="primary" size="large" @click.prevent="handleLogin">登 录</el-button>
            </el-form>
          </el-tab-pane>

          <el-tab-pane label="手机号登录" name="phone">
            <el-form ref="phoneFormRef" :model="phoneFormData" :rules="phoneFormRules" @keyup.enter="handlePhoneLogin">
              <el-form-item prop="mobile">
                <el-input
                  v-model.trim="phoneFormData.mobile"
                  placeholder="请输入手机号"
                  type="text"
                  tabindex="1"
                  size="large"
                >
                  <template #prepend>+86</template>
                </el-input>
              </el-form-item>
              <el-form-item prop="code">
                <el-input
                  v-model.trim="phoneFormData.code"
                  placeholder="请输入验证码"
                  type="text"
                  tabindex="2"
                  size="large"
                >
                  <template #append>
                    <el-button :loading="sendCodeLoading" :disabled="countdown > 0" @click="handleSendCode">
                      {{ countdown > 0 ? `${countdown}s后重试` : "获取验证码" }}
                    </el-button>
                  </template>
                </el-input>
              </el-form-item>
              <div class="form-footer">
                <el-checkbox v-model="rememberPhone">记住手机号</el-checkbox>
              </div>
              <el-button :loading="loading" type="primary" size="large" @click.prevent="handlePhoneLogin"
                >登 录</el-button
              >
            </el-form>
          </el-tab-pane>
        </el-tabs>
        <div class="other-login">
          <div class="divider">
            <span class="line" />
            <span class="text">其他登录方式</span>
            <span class="line" />
          </div>
          <div class="icons">
            <a href="" class="icon-item" title="QQ登录">
              <svg-icon name="qq-login" color="green" width="26" height="26" />
            </a>
            <a href="" class="icon-item" title="GitHub登录">
              <svg-icon name="github-login" width="24" height="24" />
            </a>
            <a
              href="https://gitee.com/oauth/authorize?client_id=62d3af1f2058b5facec5316ab9d18b8d3602406fcb704adde7f5b38de381996b&redirect_uri=http%3A%2F%2Flocalhost%3A9000%2Foauth&response_type=code&state=3"
              class="icon-item"
              title="Gitee登录"
            >
              <svg-icon name="gitee-login" width="26" height="26" />
            </a>
            <a href="javascript:;" class="icon-item" title="注册" @click="registerDialogVisible = true">
              <svg-icon name="register" width="26" height="26" />
            </a>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- 注册对话框 -->
  <el-dialog v-model="registerDialogVisible" title="用户注册" width="400px" :close-on-click-modal="false">
    <el-tabs v-model="registerActiveTab" stretch>
      <el-tab-pane label="账号注册" name="account">
        <el-form ref="registerFormRef" :model="registerFormData" :rules="registerFormRules" label-width="80px">
          <el-form-item label="用户名" prop="name">
            <el-input v-model="registerFormData.name" placeholder="请输入用户名" />
          </el-form-item>
          <el-form-item label="密码" prop="password">
            <el-input v-model="registerFormData.password" type="password" placeholder="请输入密码" show-password />
          </el-form-item>
          <el-form-item label="确认密码" prop="confirmPassword">
            <el-input
              v-model="registerFormData.confirmPassword"
              type="password"
              placeholder="请确认密码"
              show-password
            />
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="手机注册" name="phone">
        <el-form
          ref="registerPhoneFormRef"
          :model="registerPhoneFormData"
          :rules="registerPhoneFormRules"
          label-width="80px"
        >
          <el-form-item label="手机号" prop="mobile">
            <el-input v-model="registerPhoneFormData.mobile" placeholder="请输入手机号">
              <template #prepend>+86</template>
            </el-input>
          </el-form-item>
          <el-form-item label="验证码" prop="code">
            <el-input v-model="registerPhoneFormData.code" placeholder="请输入验证码">
              <template #append>
                <el-button :loading="sendCodeLoading" :disabled="countdown > 0" @click="handleSendCode">
                  {{ countdown > 0 ? `${countdown}s后重试` : "获取验证码" }}
                </el-button>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item label="密码" prop="password">
            <el-input v-model="registerPhoneFormData.password" type="password" placeholder="请输入密码" show-password />
          </el-form-item>
        </el-form>
      </el-tab-pane>
    </el-tabs>

    <template #footer>
      <span class="dialog-footer">
        <el-button @click="registerDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="registerLoading" @click="handleRegister"> 注册 </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<style lang="scss" scoped>
.login-container {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  width: 100%;
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);

  .theme-switch {
    position: fixed;
    top: 5%;
    right: 5%;
    cursor: pointer;
    transition: transform 0.3s ease;

    &:hover {
      transform: rotate(30deg);
    }
  }

  .login-card {
    width: 480px;
    max-width: 90%;
    border-radius: 20px;
    background-color: var(--el-bg-color);
    box-shadow: 0 8px 30px rgba(0, 0, 0, 0.1);
    overflow: hidden;
    transition:
      transform 0.3s ease,
      box-shadow 0.3s ease;

    &:hover {
      transform: translateY(-5px);
      box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
    }

    .title {
      display: flex;
      justify-content: center;
      align-items: center;
      height: 150px;
      padding: 20px 0;

      img {
        height: 100%;
        transition: transform 0.3s ease;

        &:hover {
          transform: scale(1.05);
        }
      }
    }

    .content {
      padding: 20px 50px 40px;

      :deep(.el-tabs__nav) {
        width: 100%;

        .el-tabs__item {
          width: 50%;
          height: 44px;
          line-height: 44px;
          text-align: center;
          font-size: 16px;
          transition: all 0.3s ease;

          &.is-active {
            font-weight: 600;
          }
        }
      }

      .el-form-item {
        margin-bottom: 25px;

        :deep(.el-input__wrapper) {
          border-radius: 8px;
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
        }

        .el-input {
          --el-input-height: 44px;

          &:focus-within {
            :deep(.el-input__wrapper) {
              box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
            }
          }
        }
      }

      .form-footer {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 20px;
      }

      .el-button {
        width: 100%;
        height: 44px;
        border-radius: 8px;
        font-size: 16px;
        font-weight: 600;
        transition: all 0.3s ease;

        &:hover {
          transform: translateY(-1px);
          box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
        }
      }

      .other-login {
        margin-top: 30px;

        .divider {
          display: flex;
          align-items: center;
          margin: 20px 0;

          .line {
            flex: 1;
            height: 1px;
            background-color: var(--el-border-color-lighter);
          }

          .text {
            padding: 0 15px;
            font-size: 14px;
            color: var(--el-text-color-secondary);
          }
        }

        .icons {
          display: flex;
          justify-content: center;
          gap: 30px;
          margin-top: 20px;

          .iconfont {
            font-size: 24px;
            color: var(--el-text-color-secondary);
            cursor: pointer;
            transition: all 0.3s ease;

            &:hover {
              transform: scale(1.2);
              color: var(--el-color-primary);
            }
          }
        }
      }
    }
  }
}
</style>
