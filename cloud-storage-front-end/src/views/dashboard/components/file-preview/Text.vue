<template>
  <div>
    <div v-show="props.maximize" style="width: 100%; height: calc(100% - 35px)">
      <TextEditor
        ref="editorRef"
        style="border: 1px #ececec solid; box-sizing: border-box"
        height="100%"
        content=""
        read-only
      />
      <div class="operateor">
        <n-button v-if="edit == false" size="small" color="#ff69b4" @click="toggleEdit(true)"> 编辑 </n-button>
        <n-button v-if="edit" size="small" @click="toggleEdit(false)"> 取消 </n-button>
        <n-button v-if="edit" size="small" type="primary" @click="saveContent"> 保存 </n-button>
      </div>
    </div>
    <div v-show="!props.maximize" class="text-content">
      {{ content }}
    </div>
  </div>
</template>

<script setup>
import { defineExpose, onMounted, ref, watch } from "vue"
import TextEditor from "./TextEditor.vue"
import axios from "axios"
// import { saveContent as saveContentHttp } from "@/http/Explore"

const props = defineProps({
  resource: Object,
  url: String,
  maximize: Boolean
})

watch(
  () => props.maximize,
  (val) => {
    console.log("最大化状态变更:", val)
    setContent(content.value)
  }
)

const edit = ref(false)
const resource = ref(props.resource)
const editorRef = ref(null)
const content = ref("")

onMounted(() => {
  console.log("Text组件挂载，URL:", props.url)
  refreshContent()
})

const refreshContent = function () {
  if (!props.url) {
    console.error("无有效的文本URL")
    setContent("无法获取文件内容")
    return
  }

  console.log("尝试获取文本内容:", props.url)
  axios
    .get(props.url)
    .then((response) => {
      console.log("获取文本成功")
      setContent(response.data)
    })
    .catch((err) => {
      console.error("获取文本失败:", err)
      setContent("获取文件内容失败: " + err.message)
    })
}

const setContent = function (text) {
  if (props.maximize && editorRef.value) {
    editorRef.value.setValue(text)
  }
  content.value = text
}

const setReadOnly = function (flag = false) {
  if (editorRef.value) {
    editorRef.value.setReadOnly(flag)
  }
}

const toggleEdit = function (flag) {
  edit.value = flag
  if (editorRef.value) {
    editorRef.value.setReadOnly(!flag)
  }
}

/**
 * 保存文件内容
 * @param resource
 */
const saveContent = function () {
  if (!resource.value) return
  console.log("保存文本内容功能暂未实现")
  // return saveContentHttp(resource.value.id, editorRef.value.getValue())
  //   .then(() => {
  //     window.$message.success("保存成功")
  //     toggleEdit(false)
  //   })
  //   .catch((err) => {
  //     return Promise.reject(err)
  //   })
}

defineExpose({
  setReadOnly,
  saveContent
})
</script>

<script>
export default {
  name: "YText"
}
</script>

<style scoped lang="scss">
.code-preview {
  width: 100%;
  height: 100%;
}
.operateor {
  height: 30px;
  padding-top: 5px;
  box-sizing: border-box;
  text-align: right;

  button {
    margin: 0 5px;
  }
}
.text-content {
  white-space: pre-wrap;
  max-height: 100%;
  overflow-y: auto;
  padding: 10px;
  background-color: #f8f8f8;
  border-radius: 5px;
  font-family: monospace;
}
#editor {
  margin-top: unset;
}
</style>
@/http/Axios
