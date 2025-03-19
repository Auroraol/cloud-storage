<template>
  <div class="office-preview">
    <div v-if="!processedUrl" class="error-message">无法获取文档URL</div>
    <div v-else-if="errorMsg" class="error-message">
      {{ errorMsg }}
      <button v-if="processedUrl" @click="downloadFile" class="download-btn">下载文件</button>
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
      <div class="viewer-toolbar">
        <div class="toolbar-info">当前使用: {{ currentViewerName }}</div>
        <div class="toolbar-actions">
          <button @click="toggleViewerService" class="viewer-switch-btn">
            {{ useGoogleViewer ? "切换到Microsoft预览" : "切换到Google预览" }}
          </button>
          <button @click="tryDirectEmbed" class="viewer-switch-btn">直接嵌入预览</button>
          <button @click="downloadFile" class="download-btn">下载文件</button>
        </div>
      </div>
      <iframe
        v-if="currentViewerUrl"
        :src="currentViewerUrl"
        width="100%"
        height="100%"
        frameborder="0"
        allowfullscreen
        :allow="useDirectEmbed ? '*' : 'autoplay; fullscreen *'"
        @load="handleIframeLoad"
        @error="handleIframeError"
        referrerpolicy="no-referrer"
        :sandbox="
          !useDirectEmbed
            ? 'allow-scripts allow-same-origin allow-forms allow-popups allow-top-navigation allow-popups-to-escape-sandbox allow-downloads'
            : undefined
        "
      />
      <div v-else class="iframe-placeholder">
        <p>无法生成预览URL，检查文件URL是否有效</p>
        <p>当前文件URL: {{ processedUrl }}</p>
        <button v-if="processedUrl" @click="forceCreateUrl" class="alt-btn">尝试强制生成预览URL</button>
      </div>
    </div>
    <vue-office-docx
      v-else-if="useNativeViewer && isDocx"
      :src="processedUrl"
      @rendered="rendered"
      @error="handleError"
    />
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
let loadingTimer = null // 改为let声明，而不是const
const loadFailCount = ref(0) // 加载失败次数计数器
const useGoogleViewer = ref(false) // 是否使用Google文档预览器
const useDirectEmbed = ref(false) // 是否直接嵌入文件URL

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

// 添加处理后的URL响应式变量
const processedUrl = ref(props.url || "")

// 创建预览URL
const viewerUrl = computed(() => {
  if (!processedUrl.value) return ""

  // 使用Microsoft Office在线查看器 - 对于https链接效果较好
  // 优先使用embed.aspx而不是view.aspx，因为嵌入模式通常在iframe中更可靠
  if (processedUrl.value.startsWith("https://")) {
    console.log("使用Microsoft Office在线查看器", processedUrl.value)
    // 检查URL中是否有特殊字符，如果有则对其进行处理
    const cleanUrl = encodeURIComponent(processedUrl.value)
    return `https://view.officeapps.live.com/op/embed.aspx?src=${cleanUrl}`
  }

  // 使用Google Docs查看器 - 适用于各种文档, 但有文件大小限制
  // URL格式: https://docs.google.com/viewer?url=YOUR_URL_HERE&embedded=true
  console.log("使用Google Docs查看器", processedUrl.value)
  return `https://docs.google.com/viewer?url=${encodeURIComponent(processedUrl.value)}&embedded=true`
})

// 当前使用的预览器名称
const currentViewerName = computed(() => {
  if (useDirectEmbed.value) {
    return "直接嵌入模式"
  } else if (useGoogleViewer.value) {
    return "Google Docs预览"
  } else {
    return "Microsoft Office预览"
  }
})

// 当前使用的预览URL
const currentViewerUrl = computed(() => {
  console.log("计算currentViewerUrl")
  console.log("viewerUrl:", viewerUrl.value)

  // 如果是直接嵌入模式，返回原始URL
  if (useDirectEmbed.value) {
    return processedUrl.value
  }

  // 根据预览模式选择使用哪个预览服务
  if (useGoogleViewer.value) {
    return `https://docs.google.com/viewer?url=${encodeURIComponent(processedUrl.value)}&embedded=true`
  }
  return viewerUrl.value
})

// 启动加载计时器函数
const startLoadingTimer = () => {
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
        if (isDoc.value) {
          loadingTextEl.textContent = ".doc文件加载较慢，请耐心等待或考虑下载到本地..."
        } else {
          loadingTextEl.textContent = "文档较大，加载可能需要较长时间，请耐心等待..."
        }
      }
    }
  }, 1000)
}

onMounted(() => {
  console.log("Office组件挂载，URL:", props.url)

  // 确保processedUrl被正确赋值
  processedUrl.value = props.url || ""

  console.log("处理后的URL:", processedUrl.value)
  console.log("文档类型:", props.resource?.type)
  console.log("文件扩展名:", fileExtension.value)

  // 检查预览URL是否正确计算
  console.log("预览URL计算值:", viewerUrl.value)
  console.log("当前使用的预览URL:", currentViewerUrl.value)

  // 启动加载时间计时器和进度条
  startLoadingTimer()

  // 将等待时间从3秒增加到15秒
  setTimeout(() => {
    if (!iframeLoaded.value && isOfficeFile.value) {
      if (isDoc.value) {
        errorMsg.value = ".doc格式文件通常需要转换，可能无法在线预览。您可以下载后查看"
      } else {
        errorMsg.value = "文档正在加载中，可能需要较长时间。您也可以选择下载后查看"
      }
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
  console.log("加载的URL是:", currentViewerUrl.value)

  // 尝试获取iframe内容，检查是否真正加载了内容
  try {
    const iframe = document.querySelector(".viewer-container iframe")
    if (iframe) {
      console.log("iframe元素存在")

      // 尝试检测iframe是否有内容，但由于跨域限制可能会失败
      // 在这种情况下，我们依然认为加载成功，不再执行检测逻辑
      // iframe加载完成即可视为成功

      // 如果加载次数超过一定限制并且使用的是Microsoft服务，尝试切换到Google预览
      if (loadFailCount.value > 0 && currentViewerUrl.value.includes("officeapps.live.com")) {
        console.log("Office预览可能不稳定，尝试切换到Google预览")
        tryAlternativeViewer()
        return
      }
    } else {
      console.log("未找到iframe元素")
    }
  } catch (e) {
    console.error("检查iframe时出错:", e)
  }

  isLoading.value = false
  iframeLoaded.value = true
  loadingProgress.value = 100
  clearInterval(loadingTimer)
}

const handleIframeError = function (error) {
  console.error("iframe加载失败", error)
  loadFailCount.value++

  // 如果失败，尝试使用替代预览方式
  if (loadFailCount.value <= 2) {
    console.log(`尝试切换预览方式，当前尝试次数: ${loadFailCount.value}`)
    tryAlternativeViewer()
    return
  }

  // 显示错误信息
  if (isDoc.value) {
    errorMsg.value = ".doc格式文件可能无法在线预览，建议下载后使用Microsoft Word或WPS等软件打开"
  } else {
    errorMsg.value = "在线预览服务暂时不可用，请尝试下载文件查看"
  }
  isLoading.value = false
  clearInterval(loadingTimer)
}

// 添加下载文件的方法
const downloadFile = () => {
  if (!processedUrl.value) return

  // 创建一个临时的a元素
  const a = document.createElement("a")
  a.href = processedUrl.value
  a.download = props.resource?.name || `文件.${fileExtension.value}`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
}

const forceCreateUrl = () => {
  // 实现强制生成预览URL的逻辑
  console.log("强制生成预览URL")
  if (!processedUrl.value) {
    console.error("没有可用的文件URL")
    return
  }

  // 强制刷新processedUrl以触发计算属性重新计算
  const tempUrl = processedUrl.value
  processedUrl.value = ""
  setTimeout(() => {
    processedUrl.value = tempUrl

    // 根据URL类型选择合适的预览服务
    let previewUrl = ""
    if (processedUrl.value.startsWith("https://")) {
      // 尝试使用嵌入模式，有时候这种模式兼容性更好
      previewUrl = `https://view.officeapps.live.com/op/embed.aspx?src=${encodeURIComponent(processedUrl.value)}`
    } else {
      // 使用Google Docs预览
      previewUrl = `https://docs.google.com/viewer?url=${encodeURIComponent(processedUrl.value)}&embedded=true`
    }

    // 直接将iframe的src设置为新URL
    const iframe = document.querySelector(".viewer-container iframe")
    if (iframe) {
      console.log("设置iframe URL为:", previewUrl)
      iframe.src = previewUrl

      // 添加必要的属性以解决可能的跨域问题
      iframe.setAttribute("referrerpolicy", "no-referrer")
      iframe.setAttribute("allow", "autoplay; fullscreen *")
    } else {
      // 如果找不到iframe，创建一个新的
      const container = document.querySelector(".viewer-container")
      if (container) {
        console.log("创建新iframe，URL为:", previewUrl)
        const newIframe = document.createElement("iframe")
        newIframe.src = previewUrl
        newIframe.width = "100%"
        newIframe.height = "100%"
        newIframe.frameBorder = "0"
        newIframe.allowFullscreen = true
        newIframe.setAttribute("referrerpolicy", "no-referrer")
        newIframe.setAttribute("allow", "autoplay; fullscreen *")
        newIframe.addEventListener("load", handleIframeLoad)
        newIframe.addEventListener("error", handleIframeError)
        container.appendChild(newIframe)
      }
    }
  }, 100)
}

// 添加可以直接使用embed.aspx的方法，有时这种方式更可靠
const tryAlternativeViewer = () => {
  if (!processedUrl.value) return

  console.log("尝试使用替代预览方式")

  // 尝试使用Google Docs预览服务，通常对于跨域场景有更好的兼容性
  const googleDocsUrl = `https://docs.google.com/viewer?url=${encodeURIComponent(processedUrl.value)}&embedded=true`

  const iframe = document.querySelector(".viewer-container iframe")
  if (iframe) {
    console.log("设置iframe URL为Google Docs方式:", googleDocsUrl)
    iframe.src = googleDocsUrl
  } else {
    console.log("没有找到iframe元素")
  }
}

// 切换预览服务
const toggleViewerService = () => {
  useGoogleViewer.value = !useGoogleViewer.value
  useDirectEmbed.value = false // 关闭直接嵌入模式
  console.log(`切换到${useGoogleViewer.value ? "Google" : "Microsoft"}预览服务`)
  loadFailCount.value = 0 // 重置失败计数

  // 刷新iframe
  const iframe = document.querySelector(".viewer-container iframe")
  if (iframe) {
    // 先设置为加载状态
    isLoading.value = true
    iframeLoaded.value = false
    startLoadingTimer()

    // 设置新的URL
    iframe.src = currentViewerUrl.value

    // 恢复sandbox属性
    iframe.setAttribute(
      "sandbox",
      "allow-scripts allow-same-origin allow-forms allow-popups allow-top-navigation allow-popups-to-escape-sandbox allow-downloads"
    )
  }
}

// 尝试直接嵌入预览
const tryDirectEmbed = () => {
  console.log("尝试直接嵌入预览")
  useDirectEmbed.value = true
  useGoogleViewer.value = false
  loadFailCount.value = 0 // 重置失败计数

  // 重置加载状态
  isLoading.value = true
  iframeLoaded.value = false
  startLoadingTimer()

  // 刷新iframe
  const iframe = document.querySelector(".viewer-container iframe")
  if (iframe) {
    iframe.src = processedUrl.value

    // 移除sandbox属性以允许直接访问
    iframe.removeAttribute("sandbox")
    iframe.setAttribute("allow", "*")
  }
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
  display: flex;
  flex-direction: column;
}

.viewer-container iframe {
  flex: 1;
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

.alt-btn {
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
  margin-left: 10px;
}

.alt-btn:hover {
  background-color: #3a7bd5;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.2);
}

.debug-info {
  position: absolute;
  top: 0;
  left: 0;
  background-color: rgba(0, 0, 0, 0.7);
  color: white;
  padding: 10px;
  z-index: 10;
  font-size: 12px;
  max-width: 80%;
  overflow-wrap: break-word;
}

.iframe-placeholder {
  text-align: center;
  padding: 20px;
}

.viewer-toolbar {
  width: 100%;
  padding: 8px;
  background-color: #f5f5f5;
  border-bottom: 1px solid #e0e0e0;
  display: flex;
  justify-content: flex-end;
  align-items: center;
}

.toolbar-info {
  flex: 1;
  text-align: left;
  font-size: 12px;
  color: #666;
  margin-left: 10px;
}

.toolbar-actions {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 8px;
}

.viewer-switch-btn {
  padding: 5px 10px;
  background-color: #4a89dc;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  transition: background-color 0.3s;
}

.viewer-switch-btn:hover {
  background-color: #3a7bd5;
}
</style>
