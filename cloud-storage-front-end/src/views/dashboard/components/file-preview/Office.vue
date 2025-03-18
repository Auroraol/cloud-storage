<template>
  <div class="office-preview">
    <div v-if="!props.url" class="error-message">无法获取文档URL</div>
    <div v-else-if="errorMsg" class="error-message">
      {{ errorMsg }}
      <button v-if="props.url" @click="downloadFile" class="download-btn">下载文件</button>
    </div>
    <div v-else-if="isLoading" class="loading-message">
      <div class="spinner" />
      <div class="loading-text">文档加载中，请稍候...</div>
      <div class="progress-container">
        <div class="progress-bar" :style="{ width: loadingProgress + '%' }" />
      </div>
      <div class="loading-time">已加载 {{ loadingTime }}秒</div>
    </div>
    <div v-else-if="isOfficeFile" class="viewer-container">
      <iframe
        :src="viewerUrl"
        width="100%"
        height="100%"
        frameborder="0"
        allowfullscreen
        @load="handleIframeLoad"
        @error="handleIframeError"
      />
    </div>
    <vue-office-docx v-else-if="useNativeViewer && isDocx" :src="props.url" @rendered="rendered" @error="handleError" />
    <div v-else class="fallback-message">
      <p>当前文件类型 (.{{ fileExtension }}) 不支持在线预览</p>
      <p>您可以下载后在本地查看</p>
      <button @click="downloadFile" class="download-btn">下载文件</button>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref, computed, onBeforeUnmount } from "vue"
import VueOfficeDocx from "@vue-office/docx"
// 引入相关样式
import "@vue-office/docx/lib/index.css"

const props = defineProps({
  resource: Object,
  url: String
})

const errorMsg = ref("")
const isLoading = ref(true)
const iframeLoaded = ref(false)
const useNativeViewer = ref(false) // 是否使用内置的docx预览器
const loadingProgress = ref(0)
const loadingTime = ref(0)
let loadingTimer = null

// 获取文件扩展名
const fileExtension = computed(() => {
  if (!props.resource?.name) return ""
  const dot = props.resource.name.lastIndexOf(".")
  return dot !== -1 ? props.resource.name.substring(dot + 1).toLowerCase() : ""
})

// 判断是否为docx文件
const isDocx = computed(() => {
  return fileExtension.value === "docx"
})

// 判断是否为doc文件
const isDoc = computed(() => {
  return fileExtension.value === "doc"
})

// 判断是否为PDF文件
const isPdf = computed(() => {
  return fileExtension.value === "pdf"
})

// 判断是否为Excel文件
const isExcel = computed(() => {
  return ["xls", "xlsx"].includes(fileExtension.value)
})

// 判断是否为PPT文件
const isPpt = computed(() => {
  return ["ppt", "pptx"].includes(fileExtension.value)
})

// 判断是否为Office文件
const isOfficeFile = computed(() => {
  return isDoc.value || isDocx.value || isPdf.value || isExcel.value || isPpt.value
})

// 创建预览URL
const viewerUrl = computed(() => {
  if (!props.url) return ""

  // 使用Microsoft Office在线查看器 - 对于https链接效果较好
  // URL格式: https://view.officeapps.live.com/op/view.aspx?src=YOUR_URL_HERE
  if (props.url.startsWith("https")) {
    return `https://view.officeapps.live.com/op/view.aspx?src=${encodeURIComponent(props.url)}`
  }

  // 使用Google Docs查看器 - 适用于各种文档, 但有文件大小限制
  // URL格式: https://docs.google.com/viewer?url=YOUR_URL_HERE&embedded=true
  return `https://docs.google.com/viewer?url=${encodeURIComponent(props.url)}&embedded=true`
})

onMounted(() => {
  console.log("Office组件挂载，URL:", props.url)
  console.log("文档类型:", props.resource?.type)
  console.log("文件扩展名:", fileExtension.value)

  // 启动加载时间计时器和进度条
  loadingTimer = setInterval(() => {
    loadingTime.value += 1

    // 让进度条看起来是在逐渐加载
    if (loadingProgress.value < 90) {
      if (loadingTime.value <= 3) {
        loadingProgress.value += 10
      } else if (loadingTime.value <= 8) {
        loadingProgress.value += 5
      } else {
        loadingProgress.value += 1
      }
    }

    // 在进度条达到一定程度时更新提示文字
    if (loadingTime.value === 5) {
      const loadingTextEl = document.querySelector(".loading-text")
      if (loadingTextEl) {
        loadingTextEl.textContent = "文档较大，加载可能需要较长时间，请耐心等待..."
      }
    }
  }, 1000)

  // 将等待时间从3秒增加到15秒
  setTimeout(() => {
    if (!iframeLoaded.value && isOfficeFile.value) {
      errorMsg.value = "文档正在加载中，可能需要较长时间。您也可以选择下载后查看"
    }
    isLoading.value = false
    clearInterval(loadingTimer)
  }, 15000) // 从10000改为15000毫秒
})

onBeforeUnmount(() => {
  // 清除定时器
  if (loadingTimer) {
    clearInterval(loadingTimer)
  }
})

const rendered = function () {
  console.log("本地文档渲染成功")
  errorMsg.value = ""
  isLoading.value = false
  iframeLoaded.value = true
  clearInterval(loadingTimer)
}

const handleError = function (error) {
  console.error("本地文档渲染失败:", error)
  errorMsg.value = "文档渲染失败，将尝试使用在线查看器"
  isLoading.value = false

  // 尝试使用在线预览服务
  useNativeViewer.value = false
  clearInterval(loadingTimer)
}

const handleIframeLoad = function () {
  console.log("iframe加载完成")
  isLoading.value = false
  iframeLoaded.value = true
  loadingProgress.value = 100
  clearInterval(loadingTimer)
}

const handleIframeError = function () {
  console.error("iframe加载失败")
  errorMsg.value = "在线预览服务暂时不可用，请尝试下载文件"
  isLoading.value = false
  clearInterval(loadingTimer)
}

// 添加下载文件的方法
const downloadFile = () => {
  if (!props.url) return

  // 创建一个临时的a元素
  const a = document.createElement("a")
  a.href = props.url
  a.download = props.resource?.name || `文件.${fileExtension.value}`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
}
</script>

<script>
export default {
  name: "YOffice"
}
</script>

<style scoped>
.office-preview {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
}

.viewer-container {
  width: 100%;
  height: 100%;
  position: relative;
}

.spinner {
  border: 4px solid rgba(0, 0, 0, 0.1);
  width: 36px;
  height: 36px;
  border-radius: 50%;
  border-left-color: #4a89dc;
  animation: spin 1s linear infinite;
  margin-bottom: 10px;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

.loading-message {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #333;
  text-align: center;
  padding: 20px;
  max-width: 80%;
}

.loading-text {
  margin-top: 10px;
  margin-bottom: 15px;
}

.loading-time {
  margin-top: 8px;
  font-size: 12px;
  color: #666;
}

.progress-container {
  width: 200px;
  height: 6px;
  background-color: #f1f1f1;
  border-radius: 3px;
  overflow: hidden;
  margin-top: 10px;
}

.progress-bar {
  height: 100%;
  background-color: #4a89dc;
  width: 0;
  transition: width 0.3s ease;
}

.error-message {
  color: red;
  text-align: center;
  padding: 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.fallback-message {
  text-align: center;
  padding: 20px;
}

.download-btn {
  display: inline-block;
  margin-top: 15px;
  padding: 10px 20px;
  background-color: #4a89dc;
  color: white;
  text-decoration: none;
  border-radius: 4px;
  font-weight: bold;
  transition: background-color 0.3s;
  border: none;
  cursor: pointer;
}

.download-btn:hover {
  background-color: #3a7bd5;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.2);
}
</style>
