<template>
  <div class="app-container">
    <div class="setting-box">
      <h1>个人信息</h1>
      <ul class="setting-list">
        <li class="item">
          <span class="label">头像</span>
          <div class="avatar-wrapper">
            <!-- <img :src="userInfo.avatar || defaultAvatar" class="user-avatar" /> -->
          </div>
          <div class="action-box">
            <div class="hint">支持 jpg、png 格式大小 300KB 以内的图片</div>
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
              <el-button :loading="loading" type="primary" size="mini" class="upload-button">点击上传</el-button>
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
            <el-input v-model="form.mobile" disabled placeholder="绑定手机号" />
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
            <el-input v-model="form.email" disabled placeholder="未绑定邮箱" />
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
      <el-button @click="$emit('close')">关闭</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue"
// import { getAccessToken } from "/@/utils/auth"
// import { updateUser as apiUpdateUser, bindUsername } from "/@/api/user/user"
// import { useSettingsStore, useUserStore } from "/@/store/index"
// import { useGetters } from "/@/store/getters"

// pinia
// const useUserStorePinia = useUserStore()
// const useGettersPinia = useGetters()
// const useSettingsStorePinia = useSettingsStore()

// 响应
const options = ref([
  { value: 1, label: "男" },
  { value: 0, label: "女" }
])

const opVisible = ref({
  nickname: false,
  birthday: false,
  gender: false,
  brief: false
})

// 上传头像相关变量
// const path = import.meta.env.VITE_APP_BASE_API + "/user/avatar/update" // 后端服务器api接口
const headers = computed(() => {
  // return { Authorization: "Bearer " + getAccessToken() } //上传文件请求头
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
    await apiUpdateUser(data)
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
    await useUserStorePinia.getUserInfo()
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
  console.error(res)

  if (res.code !== 200000) {
    ElMessage.error("文件上传失败")
    return
  }
  loading.value = false
  ElMessage.success("上传成功")
  // 重写获取用户信息  + 表单重写初始化
  useUserStorePinia.getUserInfo().then(() => init())
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
</script>

<style lang="scss" scoped>
.app-container {
  padding: 20px;
  background: #f5f5f5;
  min-height: calc(100vh - 50px);
  position: relative;

  .setting-box {
    margin: 20px auto;
    width: 90%;
    max-width: 600px;
    background: #fff;
    border-radius: 12px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
    padding: 30px;
    transition: box-shadow 0.3s ease;

    h1 {
      font-size: 26px;
      font-weight: bold;
      margin-bottom: 20px;
      color: #333;
      text-align: center;
    }

    .setting-list {
      list-style: none;
      padding: 0;

      .item {
        border-bottom: 1px solid #eaeaea;
        padding: 20px 0;
        display: flex;
        align-items: center;

        .label {
          min-width: 100px;
          font-weight: 600;
          color: #555;
          font-size: 16px;
        }

        .input-wrapper {
          flex: 1;
        }

        .action-box {
          display: flex;
          align-items: center;

          .hint {
            color: #888;
            font-size: 12px;
            margin-right: 10px;
          }

          .upload-button {
            padding: 6px 12px;
            background-color: #007bff;
            color: white;
            border-radius: 4px;
            transition: background-color 0.3s;

            &:hover {
              background-color: #0056b3;
            }
          }

          .disabled-button {
            color: #c0c4cc;
          }

          .prompt-button,
          .edit-button,
          .cancel-button {
            color: #007bff;
            font-size: 14px;
            margin-left: 10px;
            cursor: pointer;
            transition: color 0.3s;

            &:hover {
              color: #0056b3;
            }
          }

          .link-button {
            color: #007bff;
            font-size: 14px;
            margin-left: 10px;
            cursor: pointer;
            transition: color 0.3s;

            &:hover {
              color: #0056b3;
            }
          }

          .action-buttons {
            display: flex;
            align-items: center;
          }
        }
      }
    }
  }
}
</style>
