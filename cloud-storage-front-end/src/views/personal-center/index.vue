<template>
  <div class="personal-center">
    <h1>个人信息</h1>
    <div class="profile">
      <div class="info">
        <div>
          <img :src="userInfo.avatar || defaultAvatar" class="avatar-wrapper" />
        </div>
        <p class="username">用户名: {{ userInfo.username }}</p>
        <p class="nickname">昵称: {{ userInfo.nickname }}</p>
        <p class="phone">手机号码: {{ userInfo.mobile }}</p>
        <p class="email">邮箱: {{ userInfo.email }}</p>
        <p class="birthday">生日: {{ userInfo.birthday }}</p>
        <p class="gender">性别: {{ gender }}</p>
        <p class="brief">简介: {{ userInfo.brief }}</p>
      </div>
    </div>
    <div class="settings">
      <h3>设置</h3>
      <ul>
        <li>
          <el-button type="text" @click="openEditProfile">编辑个人资料</el-button>
        </li>
        <li>
          <el-button type="text" @click="openChangePassword">更改密码</el-button>
        </li>
      </ul>
    </div>

    <el-dialog v-model="dialogVisible" title="编辑个人资料" class="edit-profile-dialog" width="680px">
      <ul class="setting-list">
        <li class="item">
          <span class="label">头像</span>
          <div class="avatar-wrapper">
            <img :src="userInfo.avatar || defaultAvatar" class="user-avatar" />
          </div>
          <div class="hint">支持 jpg、png 格式大小 300KB 以内的图片</div>
          <div class="action-box">
            <el-upload
              ref="upload"
              :action="path"
              :file-list="files"
              :multiple="false"
              :limit="1"
              :headers="headers"
              :auto-upload="true"
              :name="'file'"
              :on-exceed="onExceed"
              :before-upload="beforeUpload"
              :on-success="uploadSuccess"
              :on-error="uploadError"
              :show-file-list="false"
            >
              <el-button :loading="loading" type="primary" size="small" class="upload-button">点击上传</el-button>
            </el-upload>
          </div>
        </li>
        <li class="item">
          <span class="label">用户名</span>
          <div class="input-wrapper">
            <el-input v-model="userInfo.username" disabled placeholder="未填写" />
          </div>
          <div class="action-box">
            <el-button v-if="userInfo.username" type="text" disabled class="disabled-button">
              <span><i class="el-icon-lock" />不可修改</span>
            </el-button>
            <el-button v-else type="text" @click="usernamePrompt" class="prompt-button">
              <span><i class="el-icon-lock" />立即填写</span>
            </el-button>
          </div>
        </li>
        <li class="item">
          <span class="label">昵称</span>
          <div class="input-wrapper">
            <el-input
              ref="nicknameIput"
              v-model="form.nickname"
              placeholder="填写你的昵称"
              @focus="nicknameIputFocus"
              class="input-nickname"
            />
          </div>
          <div class="action-box">
            <el-button v-show="!opVisible.nickname" type="text" @click="nicknameIputFocus" class="edit-button">
              <span><i class="el-icon-edit" />修改</span>
            </el-button>
            <span v-show="opVisible.nickname" class="action-buttons">
              <el-button type="text" class="cancel-button" @click="cancelNickname">取消</el-button>
              <el-button type="text" @click="saveNickname">保存</el-button>
            </span>
          </div>
        </li>
        <li class="item">
          <span class="label">手机号</span>
          <div class="input-wrapper">
            <el-input v-model="form.mobile" placeholder="绑定手机号" />
          </div>
          <div class="action-box">
            <router-link v-if="form.mobile" to="/rebind-mobile" class="link-button">
              <span><i class="el-icon-mobile-phone" />更改绑定</span>
            </router-link>
            <router-link v-else to="/bind-mobile" class="link-button">
              <span><i class="el-icon-mobile-phone" />立即绑定</span>
            </router-link>
          </div>
        </li>
        <li class="item">
          <span class="label">邮箱</span>
          <div class="input-wrapper">
            <el-input v-model="form.email" placeholder="未绑定邮箱" />
          </div>
          <div class="action-box">
            <router-link type="text" to="/email-validate" class="link-button">
              <span><i class="el-icon-message" />{{ form.email ? "更改绑定" : "立即绑定" }}</span>
            </router-link>
          </div>
        </li>
        <li class="item">
          <span class="label">生日</span>
          <div class="input-wrapper">
            <el-date-picker
              ref="birthdayIput"
              v-model="form.birthday"
              type="date"
              :prefix-icon="''"
              placeholder="选择日期"
              :clearable="false"
              @focus="birthdayIputFocus"
              class="input-date"
            />
          </div>
          <div class="action-box">
            <el-button v-show="!opVisible.birthday" type="text" @click="birthdayIputFocus" class="edit-button">
              <span><i class="el-icon-date" />修改</span>
            </el-button>
            <span v-show="opVisible.birthday" class="action-buttons">
              <el-button type="text" class="cancel-button" @click="cancelBirthday">取消</el-button>
              <el-button type="text" @click="saveBirthday">保存</el-button>
            </span>
          </div>
        </li>
        <li class="item">
          <span class="label">性别</span>
          <div class="input-wrapper">
            <el-select
              ref="genderIput"
              v-model="form.gender"
              placeholder="请选择"
              @focus="genderIputFocus"
              class="input-gender"
            >
              <el-option v-for="item in options" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </div>
          <div class="action-box">
            <el-button v-show="!opVisible.gender" type="text" @click="genderIputFocus" class="edit-button">
              <span><i class="el-icon-female" />修改</span>
            </el-button>
            <span v-show="opVisible.gender" class="action-buttons">
              <el-button type="text" class="cancel-button" @click="cancelGender">取消</el-button>
              <el-button type="text" @click="saveGender">保存</el-button>
            </span>
          </div>
        </li>
        <li class="item">
          <span class="label">简介</span>
          <div class="input-wrapper">
            <el-input
              ref="briefIput"
              v-model="form.brief"
              type="textarea"
              :rows="2"
              placeholder="填写你的简介"
              @focus="briefIputFocus"
              class="input-brief"
            />
          </div>
          <div class="action-box">
            <el-button v-show="!opVisible.brief" type="text" @click="briefIputFocus" class="edit-button">
              <span><i class="el-icon-edit" />修改</span>
            </el-button>
            <span v-show="opVisible.brief" class="action-buttons">
              <el-button type="text" class="cancel-button" @click="cancelBrief">取消</el-button>
              <el-button type="text" @click="saveBrief">保存</el-button>
            </span>
          </div>
        </li>
      </ul>
    </el-dialog>

    <!-- 添加修改密码对话框 -->
    <el-dialog v-model="passwordDialogVisible" title="修改密码" width="500px">
      <el-form :model="passwordForm" :rules="passwordRules" ref="passwordFormRef" label-width="100px">
        <el-form-item label="原密码" prop="oldPassword">
          <el-input v-model="passwordForm.oldPassword" type="password" show-password placeholder="请输入原密码" />
        </el-form-item>
        <el-form-item label="新密码" prop="newPassword">
          <el-input v-model="passwordForm.newPassword" type="password" show-password placeholder="请输入新密码" />
        </el-form-item>
        <el-form-item label="确认新密码" prop="confirmPassword">
          <el-input
            v-model="passwordForm.confirmPassword"
            type="password"
            show-password
            placeholder="请再次输入新密码"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="passwordDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitChangePassword" :loading="changePasswordLoading"> 确认 </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted } from "vue"
import { ElMessage, FormInstance } from "element-plus"
// import { useUserStore } from "@/store/modules/user"
// import { apiUpdateUser } from "@/api/user" // 需要创建这个API函数

const dialogVisible = ref(false)
const passwordDialogVisible = ref(false)
const changePasswordLoading = ref(false)
const passwordFormRef = ref<FormInstance>()

const openEditProfile = () => {
  dialogVisible.value = true
}

// 响应
const options = ref([
  { value: 1, label: "男" },
  { value: 0, label: "女" }
])

// 性别 userInfo.value.gender是0,1 0是女 1是男
const gender = computed(() => {
  return userInfo.value.gender === 0 ? "女" : "男"
})

const opVisible = ref({
  nickname: false,
  birthday: false,
  gender: false,
  brief: false
})

// 上传头像相关变量
// const path = `${import.meta.env.VITE_APP_BASE_API}/user/avatar/update` // 修复上传路径

const headers = computed(() => {
  return { Authorization: `Bearer ${userStore.token}` } // 修复请求头
})

const files = ref([])
const loading = ref(false)
const form = ref({
  nickname: "",
  mobile: "",
  email: "",
  gender: 0,
  birthday: "",
  brief: ""
})

//  ref获取dom元素
const nicknameIput = ref() // 昵称
const genderIput = ref() // 性别
const briefIput = ref() // 简介
const birthdayIput = ref() // 简介
import { useUserStore } from "@/store/modules/user"
const userStore = useUserStore()
// 默认头像
// const defaultAvatar = computed(() => useSettingsStorePinia.defaultAvatar)

// 计算属性
const userInfo = computed(() => {
  const info = userStore.userInfo
  return Object.keys(info).length === 0 ? null : info
})

// 初始化
onMounted(() => {
  init()
})

const init = () => {
  form.value.nickname = userInfo.value.nickname
  form.value.mobile = sensitiveMobile(userInfo.value.mobile)
  form.value.email = sensitiveEmail(userInfo.value.email)
  form.value.gender = userInfo.value.gender
  form.value.birthday = userInfo.value.birthday
  form.value.brief = userInfo.value.brief
}

// 昵称输入框焦点事件
const nicknameIputFocus = () => {
  nicknameIput.value.focus()
  opVisible.value.nickname = true
}

// 昵称取消保存
const cancelNickname = () => {
  form.value.nickname = userInfo.value.nickname
  opVisible.value.nickname = false
}

// 昵称保存
const saveNickname = () => {
  const data = { nickname: form.value.nickname, userId: userInfo.value.id }
  updateUser(data, 0)
}

// 生日栏聚焦
const birthdayIputFocus = () => {
  birthdayIput.value.focus()
  opVisible.value.birthday = true
}

// 生日取消保存
const cancelBirthday = () => {
  form.value.birthday = userInfo.value.birthday
  opVisible.value.birthday = false
}

// 保存生日
const saveBirthday = () => {
  const data = { birthday: form.value.birthday, userId: userInfo.value.id }
  console.error(data)
  updateUser(data, 1)
}

// 性别栏聚焦
const genderIputFocus = () => {
  genderIput.value.focus()
  opVisible.value.gender = true
}

// 性别取消保存
const cancelGender = () => {
  form.value.gender = userInfo.value.gender
  opVisible.value.gender = false
}

// 性别保存
const saveGender = () => {
  const data = { gender: form.value.gender, userId: userInfo.value.id }
  updateUser(data, 2)
}

// 简介聚焦
const briefIputFocus = () => {
  briefIput.value.focus()
  opVisible.value.brief = true
}

// 简介取消保存
const cancelBrief = () => {
  form.value.brief = userInfo.value.brief
  opVisible.value.brief = false
}

// 保存简介
const saveBrief = () => {
  const data = { brief: form.value.brief, userId: userInfo.value.id }
  updateUser(data, 3)
}

// 更新用户信息
const updateUser = async (data: any, index: number) => {
  try {
    // await apiUpdateUser(data)
    switch (index) {
      case 0:
        opVisible.value.nickname = false
        break
      case 1:
        opVisible.value.birthday = false
        break
      case 2:
        opVisible.value.gender = false
        break
      case 3:
        opVisible.value.brief = false
        break
    }
    ElMessage.success("保存成功")
    // 重写获取用户信息
    // await userStore.getUserInfo()
    // 表单重写初始化
    init()
  } catch (error) {
    console.error(error)
    ElMessage.error("保存失败")
  }
}

// 上传
const onExceed = () => {
  loading.value = false
}

const uploadSuccess = (res: any) => {
  if (res.code !== 200000) {
    ElMessage.error("文件上传失败")
    return
  }
  loading.value = false
  ElMessage.success("上传成功")
  // 重写获取用户信息 + 表单重写初始化
  // userStore.getUserInfo().then(() => init())
}

const uploadError = (err: any) => {
  console.error(err)
  loading.value = false
  ElMessage.error("文件上传失败")
}

const beforeUpload = (file: any) => {
  loading.value = true
  const isImg = file.type === "image/jpeg" || file.type === "image/png" || file.type === "image/jpg"
  const isLt300KB = file.size / 1000 < 300
  if (!isImg) {
    ElMessage.error("文件格式不正确")
    loading.value = false
  }
  if (!isLt300KB) {
    ElMessage.error("文件大小不能大于300KB")
    loading.value = false
  }
  return isImg && isLt300KB
}

const usernamePrompt = () => {
  // $prompt(
  //   "字母开头，允许2-16字节，允许字母数字下划线，并且用户名成功填写后不允许修改。",
  //   "填写用户名",
  //   {
  //     confirmButtonText: "确定",
  //     cancelButtonText: "取消",
  //     inputPattern: /^[a-zA-Z][a-zA-Z0-9_]{1,15}$/,
  //     inputErrorMessage: "用户名格式不正确",
  //   }
  // )
  //   .then(({ value }) => {
  //     const params = { username: value };
  //     bindUsername(params).then(() => {
  //       ElMessage.success("绑定成功");
  //       $store.dispatch("user/getUserInfo").then(() => init());
  //     });
  //   })
  //   .catch(() => {
  //     // TODO
  //   });
}

const sensitiveEmail = (email: string) => {
  return email ? `${email.substr(0, 2)}****${email.substr(email.indexOf("@"))}` : ""
}

const sensitiveMobile = (mobile: string) => {
  const pat = /(\d{3})\d*(\d{4})/
  return mobile ? mobile.toString().replace(pat, "$1****$2") : ""
}

// 修改密码表单
const passwordForm = ref({
  oldPassword: "",
  newPassword: "",
  confirmPassword: ""
})

const validateConfirmPassword = (rule: any, value: string, callback: any) => {
  if (value === "") {
    callback(new Error("请再次输入密码"))
  } else if (value !== passwordForm.value.newPassword) {
    callback(new Error("两次输入密码不一致"))
  } else {
    callback()
  }
}

const passwordRules = {
  oldPassword: [
    { required: true, message: "请输入原密码", trigger: "blur" },
    { min: 6, max: 20, message: "长度在 6 到 20 个字符", trigger: "blur" }
  ],
  newPassword: [
    { required: true, message: "请输入新密码", trigger: "blur" },
    { min: 6, max: 20, message: "长度在 6 到 20 个字符", trigger: "blur" }
  ],
  confirmPassword: [{ required: true, validator: validateConfirmPassword, trigger: "blur" }]
}

const openChangePassword = () => {
  passwordDialogVisible.value = true
  passwordForm.value = {
    oldPassword: "",
    newPassword: "",
    confirmPassword: ""
  }
}

const submitChangePassword = async () => {
  if (!passwordFormRef.value) return

  await passwordFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        changePasswordLoading.value = true
        // 调用修改密码API
        // await apiChangePassword({
        //   oldPassword: passwordForm.value.oldPassword,
        //   newPassword: passwordForm.value.newPassword
        // })
        ElMessage.success("密码修改成功")
        passwordDialogVisible.value = false
      } catch (error: any) {
        ElMessage.error(error.message || "修改失败")
      } finally {
        changePasswordLoading.value = false
      }
    }
  })
}
</script>

<style scoped>
.personal-center {
  padding: 20px;
  background-color: #f9f9f9;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.profile {
  display: flex;
  align-items: center;
  width: 100%;
  margin-bottom: 20px;
}

.avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  margin-right: 20px;
  border: 2px solid #007bff;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.2);
}

.info {
  width: 100%;
  display: flex;
  flex-direction: column;
  background-color: #fff;
  padding: 15px;
  border-radius: 8px;
  box-shadow: 0 1px 5px rgba(0, 0, 0, 0.1);
}

.username,
.nickname,
.phone,
.email,
.birthday,
.gender {
  margin: 5px 0;
  font-size: 16px;
  color: #333;
}

.settings {
  background-color: #fff;
  padding: 15px;
  border-radius: 8px;
  box-shadow: 0 1px 5px rgba(0, 0, 0, 0.1);
}

.settings h3 {
  margin-top: 0;
}

.settings ul {
  list-style: none;
  padding: 0;
}

.settings li {
  margin: 10px 0;
}

.settings a {
  text-decoration: none;
  color: #007bff;
}

.settings a:hover {
  text-decoration: underline;
}

/* 添加对话框相关样式 */
:deep(.edit-profile-dialog) {
  /* .el-dialog__header { */
  /* margin: 0; */
  /* padding: 20px 24px; */
  /* border-bottom: 1px solid #f0f0f0; */
  /* } */

  .el-dialog__title {
    font-size: 16px;
    font-weight: 500;
    color: #333;
  }

  .el-dialog__body {
    padding: 24px;
  }

  .setting-list {
    list-style: none;
    padding: 0;
    margin: 0;
  }

  .item {
    display: flex;
    align-items: flex-start;
    padding: 24px 0;
    border-bottom: 1px solid #f0f0f0;

    &:last-child {
      border-bottom: none;
    }
  }

  .label {
    width: 80px;
    font-size: 14px;
    color: #666;
    line-height: 32px;
  }

  .input-wrapper {
    flex: 1;
    margin: 0 16px;
  }

  .action-box {
    width: 100px;
    text-align: right;
  }

  .hint {
    font-size: 12px;
    color: #999;
    margin-bottom: 8px;
  }

  .user-avatar {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .upload-button {
    margin-top: 8px;
  }

  .action-buttons {
    .el-button {
      padding: 0 4px;
      font-size: 14px;

      &.cancel-button {
        color: #999;
        margin-right: 12px;
      }
    }
  }

  .input-nickname,
  .input-gender,
  .input-date,
  .input-brief {
    width: 100%;
  }

  .input-brief {
    .el-textarea__inner {
      min-height: 80px;
    }
  }
}

.avatar-wrapper {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  overflow: hidden;
  margin: 0 16px;
  border: 1px dashed #d9d9d9;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 修改 Element Plus 默认样式 */
:deep(.el-dialog) {
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

:deep(.el-input.is-disabled .el-input__inner) {
  color: #666;
  -webkit-text-fill-color: #666;
  background-color: #f5f7fa;
}

:deep(.el-button--text) {
  &:not(.disabled-button) {
    color: #1890ff;

    &:hover {
      color: #40a9ff;
    }
  }

  &.disabled-button {
    color: #999;
    cursor: not-allowed;
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
