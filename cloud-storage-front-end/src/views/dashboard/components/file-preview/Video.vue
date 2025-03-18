<template>
  <div class="preview-video" :class="{ maximized: maximize }">
    <div v-if="!url" class="error-message">无法获取视频URL</div>
    <div v-else-if="loading" class="loading-container">
      <div class="loading-spinner" />
      <div class="loading-text">视频加载中...</div>
    </div>
    <div v-else-if="error" class="error-container">
      <div class="error-icon">
        <svg viewBox="0 0 1024 1024" width="48" height="48">
          <path
            d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64zm0 820c-205.4 0-372-166.6-372-372s166.6-372 372-372 372 166.6 372 372-166.6 372-372 372z"
            fill="#ff4d4f"
          />
          <path
            d="M464 688a48 48 0 1 0 96 0 48 48 0 1 0-96 0zm24-112h48c4.4 0 8-3.6 8-8V296c0-4.4-3.6-8-8-8h-48c-4.4 0-8 3.6-8 8v272c0 4.4 3.6 8 8 8z"
            fill="#ff4d4f"
          />
        </svg>
      </div>
      <div class="error-title">视频加载失败</div>
      <div class="error-message-details">{{ errorMessage }}</div>
      <div class="error-actions">
        <button class="error-btn retry-btn" @click="retryLoadVideo">
          <span>重试加载</span>
        </button>
        <button class="error-btn alt-player-btn" @click="useNativePlayer" v-if="isCorsProblem">
          <span>使用系统播放器</span>
        </button>
      </div>
      <div id="dplayer" style="display: none" />
    </div>
    <div v-else-if="nativePlayer" class="native-player-container">
      <div id="native-video-container">
        <video id="native-video" controls preload="auto" />
      </div>
      <div class="native-video-tools">
        <a :href="props.url" download :title="props.resource?.name || '下载视频'" class="download-btn">
          <svg viewBox="0 0 1024 1024" width="16" height="16">
            <path
              d="M505.7 661c3.2 4.1 9.4 4.1 12.6 0l112-141.7c4.1-5.2 0.4-12.9-6.3-12.9h-74.1V168c0-4.4-3.6-8-8-8h-60c-4.4 0-8 3.6-8 8v338.3H400c-6.7 0-10.4 7.7-6.3 12.9l112 141.8z"
              fill="currentColor"
            />
            <path
              d="M878 626h-60c-4.4 0-8 3.6-8 8v154H214V634c0-4.4-3.6-8-8-8h-60c-4.4 0-8 3.6-8 8v198c0 17.7 14.3 32 32 32h684c17.7 0 32-14.3 32-32V634c0-4.4-3.6-8-8-8z"
              fill="currentColor"
            />
          </svg>
          <span>下载视频</span>
        </a>
      </div>
    </div>
    <div v-else class="video-container">
      <div id="dplayer" />
      <div class="video-tools">
        <div class="tool-group">
          <button class="tool-btn" @click="rotate(-90)" title="向左旋转">
            <svg viewBox="0 0 1024 1024" width="16" height="16">
              <path
                d="M672 418H144c-17.7 0-32 14.3-32 32v414c0 17.7 14.3 32 32 32h528c17.7 0 32-14.3 32-32V450c0-17.7-14.3-32-32-32z m-44 402H188V494h440v326z"
                fill="currentColor"
              />
              <path
                d="M819.3 328.5c-78.8-100.7-196-153.6-314.6-154.2l-0.2-64c0-6.5-7.6-10.1-12.6-6.1l-128 101c-4 3.1-4 9.1 0 12.3l128 101c5 3.9 12.6 0.4 12.6-6.1v-64.2c89.1 0.4 178.6 36.3 240.4 108.3 65.6 76.1 86.2 185.5 53.8 280.9-2.3 6.5 1.2 13.5 7.9 15.3 25.2 6.8 50.2 10.2 75 10.2 12.9 0 25.8-0.8 38.6-2.6 6.6-0.9 11-7.2 9.4-13.7-32.2-128.4-1.1-277.4 80.7-318.1z"
                fill="currentColor"
              />
            </svg>
          </button>
          <button class="tool-btn" @click="rotate(90)" title="向右旋转">
            <svg viewBox="0 0 1024 1024" width="16" height="16">
              <path
                d="M672 418H144c-17.7 0-32 14.3-32 32v414c0 17.7 14.3 32 32 32h528c17.7 0 32-14.3 32-32V450c0-17.7-14.3-32-32-32z m-44 402H188V494h440v326z"
                fill="currentColor"
              />
              <path
                d="M819.3 328.5c-78.8-100.7-196-153.6-314.6-154.2l-0.2-64c0-6.5-7.6-10.1-12.6-6.1l-128 101c-4 3.1-4 9.1 0 12.3l128 101c5 3.9 12.6 0.4 12.6-6.1v-64.2c89.1 0.4 178.6 36.3 240.4 108.3 65.6 76.1 86.2 185.5 53.8 280.9-2.3 6.5 1.2 13.5 7.9 15.3 25.2 6.8 50.2 10.2 75 10.2 12.9 0 25.8-0.8 38.6-2.6 6.6-0.9 11-7.2 9.4-13.7-32.2-128.4-1.1-277.4 80.7-318.1z"
                fill="currentColor"
                transform="scale(-1, 1) translate(-1024, 0)"
              />
            </svg>
          </button>
        </div>
        <div class="tool-group">
          <button
            v-for="speed in playbackSpeeds"
            :key="speed"
            class="tool-btn"
            :class="{ active: currentSpeed === speed }"
            @click="setPlaybackRate(speed)"
            :title="`${speed}x 播放速度`"
          >
            {{ speed }}x
          </button>
        </div>
      </div>
    </div>

    <div ref="backupPlayerContainer" class="backup-player-container" style="display: none">
      <div id="backup-dplayer" />
    </div>
  </div>
</template>

<script setup>
import { onMounted, watch, ref, onBeforeUnmount, nextTick } from "vue"
import DPlayer from "dplayer"

const props = defineProps({
  resource: Object,
  url: String,
  maximize: Boolean
})

const loading = ref(true)
const error = ref(false)
const errorMessage = ref("")
const isCorsProblem = ref(false)
const showControls = ref(false)
const currentRotation = ref(0)
const currentSpeed = ref(1.0)
const playbackSpeeds = [0.5, 0.75, 1.0, 1.25, 1.5, 2.0]
const retryCount = ref(0)
const backupPlayerContainer = ref(null)
const nativePlayer = ref(false)
const showRetryButton = ref(false)
let player = null

// 添加一个调试日志函数，便于跟踪问题
function debugLog(message, data = null) {
  const timestamp = new Date().toISOString().substr(11, 12)
  const logPrefix = `[VideoPlayer ${timestamp}]`

  if (data) {
    console.log(logPrefix, message, data)
  } else {
    console.log(logPrefix, message)
  }
}

onMounted(() => {
  debugLog("Video组件挂载，URL:", props.url)

  // 默认显示控制栏
  showControls.value = true

  // 使用备份容器ID作为回退
  if (backupPlayerContainer.value) {
    const backupDplayer = backupPlayerContainer.value.querySelector("#backup-dplayer")
    if (backupDplayer) {
      debugLog("备用播放器容器已准备就绪")
    }
  }

  // 记录初始DOM状态，帮助诊断问题
  debugLog("初始DOM状态:", {
    mainContainer: !!document.getElementById("dplayer"),
    backupContainer: !!(backupPlayerContainer.value && backupPlayerContainer.value.querySelector("#backup-dplayer")),
    videoContainerExists: !!document.querySelector(".video-container")
  })

  if (props.url) {
    initPlayer()
  }
})

watch(
  () => props.url,
  (newUrl) => {
    debugLog("视频URL变更:", newUrl)
    if (newUrl) {
      loading.value = true
      error.value = false
      errorMessage.value = ""
      retryCount.value = 0
      initPlayer()
    }
  }
)

watch(
  () => props.maximize,
  () => {
    if (player) {
      setTimeout(() => {
        player.resize()
      }, 300)
    }
  }
)

async function checkVideoAvailability(url) {
  try {
    const isAliyunOSS = url.includes("aliyuncs.com")

    if (isAliyunOSS) {
      const response = await fetch(url, {
        method: "HEAD",
        mode: "no-cors",
        cache: "no-cache"
      })
      return true
    } else {
      const response = await fetch(url, {
        method: "HEAD",
        cache: "no-cache"
      })
      return response.ok
    }
  } catch (e) {
    console.error("视频URL检查失败:", e)
    return true
  }
}

async function initPlayer() {
  debugLog("开始初始化播放器")
  try {
    // 如果已存在播放器实例，先销毁
    if (player) {
      debugLog("销毁现有播放器实例")
      player.destroy()
      player = null
    }

    // 重置原生播放器状态
    nativePlayer.value = false

    // 创建一个变量来跟踪DOM检查的重试次数
    const maxDomRetries = 3
    let domRetryCount = 0
    let container = null

    // 使用一个函数检查DOM元素
    const checkDomElement = async () => {
      // 等待DOM渲染完成
      await nextTick()

      debugLog("检查DOM元素状态", {
        attempt: domRetryCount + 1,
        maxAttempts: maxDomRetries
      })

      // 首先检查主播放器容器
      container = document.getElementById("dplayer")

      // 如果主容器不存在，尝试使用备用容器
      if (!container && backupPlayerContainer.value) {
        debugLog("主播放器容器不存在，尝试使用备用容器")
        const backupContainer = backupPlayerContainer.value.querySelector("#backup-dplayer")
        if (backupContainer) {
          debugLog("使用备用播放器容器")
          // 显示备用容器
          backupPlayerContainer.value.style.display = "block"
          return backupContainer
        }
      }

      if (!container) {
        domRetryCount++
        if (domRetryCount <= maxDomRetries) {
          debugLog(`视频容器元素不存在，将在300ms后重试 (${domRetryCount}/${maxDomRetries})`)

          return new Promise((resolve) => {
            setTimeout(async () => {
              if (!error.value) {
                resolve(await checkDomElement())
              } else {
                resolve(null)
              }
            }, 300)
          })
        } else {
          debugLog(`视频容器元素检查已达到最大重试次数(${maxDomRetries})，尝试使用系统播放器`)
          return null
        }
      }

      debugLog("找到播放器容器", container)
      return container
    }

    // 尝试检查和获取DOM元素
    container = await checkDomElement()

    // 如果容器元素仍然不存在，使用系统播放器或显示错误
    if (!container) {
      debugLog("无法找到视频容器元素，将使用系统播放器")
      isCorsProblem.value = true
      useNativePlayer()
      return
    }

    // 先检查视频URL是否可访问
    debugLog("检查视频URL可访问性")
    const isURLValid = await checkVideoAvailability(props.url)
    debugLog("视频URL可访问性检查结果:", isURLValid)

    // 处理URL (如果是阿里云OSS，可能需要特殊处理)
    const videoUrl = props.url
    const isOssUrl = videoUrl.includes("aliyuncs.com")

    if (isOssUrl) {
      debugLog("检测到阿里云OSS链接")

      // 尝试直接使用系统播放器，避免CORS问题
      debugLog("阿里云视频直接使用系统原生播放器")
      isCorsProblem.value = true
      useNativePlayer()
      return
    }

    currentRotation.value = 0
    currentSpeed.value = 1.0

    // 根据视频类型设置更多参数
    const videoType = _getType(getFileExtension(props.resource?.name || ""))
    debugLog("视频MIME类型:", videoType)

    debugLog("初始化DPlayer播放器", {
      container: container,
      url: videoUrl,
      type: videoType
    })

    try {
      // 配置DPlayer
      player = new DPlayer({
        container: container,
        autoplay: false,
        theme: "#42b983",
        loop: false,
        lang: "zh-cn",
        screenshot: true,
        hotkey: true,
        preload: "metadata",
        volume: 0.7,
        mutex: true,
        video: {
          url: videoUrl,
          type: videoType,
          pic: props.resource?.fileCover || "",
          thumbnails: props.resource?.fileCover || "",
          defaultQuality: 0,
          quality: [
            {
              name: "标准",
              url: videoUrl,
              type: videoType
            }
          ]
        },
        subtitle: {
          url: props.resource?.subtitleUrl || "",
          type: "webvtt",
          fontSize: "20px",
          bottom: "10%",
          color: "#fff"
        }
      })

      debugLog("DPlayer实例创建成功")
    } catch (initError) {
      debugLog("DPlayer初始化失败:", initError)
      // 如果DPlayer初始化失败，尝试使用系统播放器
      handleVideoError(initError)
      return
    }

    // 绑定事件
    player.on("loadedmetadata", () => {
      debugLog("视频元数据加载完成")
      loading.value = false
      error.value = false
      showControls.value = true
    })

    player.on("error", (e) => {
      debugLog("视频加载失败", e)
      handleVideoError(e)
    })

    // 添加加载超时处理
    const loadTimeout = setTimeout(() => {
      if (loading.value) {
        debugLog("视频加载超时")
        handleVideoError({ message: "视频加载超时，请检查网络连接或视频格式。" })
      }
    }, 15000)

    player.on("play", () => {
      debugLog("视频开始播放")
      clearTimeout(loadTimeout)
    })

    player.on("pause", () => {
      debugLog("视频暂停")
    })

    player.on("waiting", () => {
      debugLog("视频加载中...")
    })

    debugLog("视频播放器初始化成功")
  } catch (error) {
    debugLog("视频播放器初始化失败:", error)
    handleVideoError(error)
  }
}

function handleVideoError(error) {
  debugLog("处理视频错误:", error)
  loading.value = false
  error.value = true

  // 根据错误类型设置用户友好的错误信息
  if (error && error.type === "NetworkError") {
    errorMessage.value = "网络错误：无法加载视频。请检查您的网络连接。"
  } else if (error && error.type === "MediaError") {
    errorMessage.value = "媒体错误：视频格式不受支持或文件已损坏。"
  } else if (error && error.message) {
    errorMessage.value = `视频加载失败: ${error.message}`
  } else {
    errorMessage.value = "视频加载失败，请稍后重试或尝试使用其他浏览器。"
  }

  // 自动重试逻辑
  if (retryCount.value < 2) {
    retryCount.value++
    debugLog(`自动重试加载视频 (${retryCount.value}/2)`)

    // 在重试之前先销毁现有播放器
    if (player) {
      debugLog("重试前销毁现有播放器")
      player.destroy()
      player = null
    }

    setTimeout(() => {
      if (!error.value) return // 如果用户已取消错误状态，则不重试

      debugLog("尝试改用系统播放器")
      useNativePlayer()
    }, 1000)
  } else {
    debugLog("达到最大重试次数，提示用户手动重试")
    showRetryButton.value = true
  }
}

function retryLoadVideo() {
  debugLog("用户手动重试加载视频")
  retryCount.value = 0
  error.value = false
  loading.value = true
  errorMessage.value = ""
  isCorsProblem.value = false
  nativePlayer.value = false // 重置原生播放器状态
  showRetryButton.value = false

  // 如果已存在播放器实例，先销毁
  if (player) {
    debugLog("重试前销毁现有播放器实例")
    player.destroy()
    player = null
  }

  setTimeout(() => {
    initPlayer()
  }, 500)
}

function useNativePlayer() {
  debugLog("启用系统原生播放器模式")
  nativePlayer.value = true
  loading.value = false
  error.value = false
  errorMessage.value = ""

  // 强制使用video标签的原生控件
  nextTick(() => {
    const videoElement = document.getElementById("native-video")

    if (!videoElement) {
      debugLog("无法找到原生视频元素", {
        elementExists: !!document.getElementById("native-video"),
        nativePlayerValue: nativePlayer.value,
        nativeContainer: !!document.getElementById("native-video-container")
      })
      error.value = true
      nativePlayer.value = false // 如果失败则重置状态
      errorMessage.value = "无法初始化视频播放器，请刷新页面或使用其他浏览器尝试。"
      return
    }

    debugLog("原生视频元素初始化成功", videoElement)

    // 添加事件监听
    videoElement.addEventListener("loadeddata", () => {
      debugLog("原生播放器视频数据加载完成")
      loading.value = false
    })

    videoElement.addEventListener("error", (e) => {
      debugLog("原生播放器加载视频失败", e)
      error.value = true
      nativePlayer.value = false // 如果失败则重置状态
      errorMessage.value = "视频加载失败，请检查网络连接或文件格式。"
    })

    videoElement.addEventListener("play", () => {
      debugLog("原生播放器开始播放")
    })

    // 设置源
    videoElement.src = props.url
    debugLog("原生播放器URL已设置", props.url)

    // 加载视频
    try {
      videoElement.load()
      debugLog("原生播放器开始加载视频")
    } catch (e) {
      debugLog("原生播放器加载视频出错", e)
      error.value = true
      nativePlayer.value = false // 如果失败则重置状态
      errorMessage.value = `加载视频时出错: ${e.message}`
    }
  })

  // 添加必要的样式
  if (!document.getElementById("native-video-styles")) {
    const style = document.createElement("style")
    style.id = "native-video-styles"
    style.textContent = `
      #native-video-container {
        width: 100%;
        height: 100%;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
      }
      
      #native-video {
        max-width: 100%;
        max-height: 100%;
        border-radius: 4px;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
      }
      
      .video-controls {
        margin-top: 15px;
        display: flex;
        gap: 10px;
      }
      
      .download-link {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        padding: 5px 15px;
        background-color: #42b983;
        color: white;
        text-decoration: none;
        border-radius: 4px;
        font-size: 14px;
        transition: background-color 0.2s;
      }
    `
    document.head.appendChild(style)
  }

  setTimeout(() => {
    const videoElement = document.getElementById("native-video")
    if (videoElement) {
      videoElement.play().catch((err) => {
        debugLog("自动播放失败:", err)
      })
    }
  }, 500)
}

function rotate(degrees) {
  if (!player) return

  currentRotation.value = (currentRotation.value + degrees) % 360
  const video = player.video.container.querySelector("video")
  if (video) {
    video.style.transform = `rotate(${currentRotation.value}deg)`

    if (currentRotation.value % 180 !== 0) {
      const aspectRatio = video.videoHeight / video.videoWidth
      video.style.transformOrigin = "center center"

      if (window.innerWidth < window.innerHeight) {
        video.style.width = `${video.parentElement.clientHeight * aspectRatio}px`
        video.style.height = `${video.parentElement.clientWidth / aspectRatio}px`
      } else {
        video.style.width = `${video.parentElement.clientHeight}px`
        video.style.maxHeight = `${video.parentElement.clientWidth}px`
      }
    } else {
      video.style.width = "100%"
      video.style.height = "auto"
      video.style.maxHeight = "100%"
    }
  }
}

function setPlaybackRate(rate) {
  if (!player) return
  player.speed(rate)
  currentSpeed.value = rate
}

function getFileExtension(filename) {
  if (!filename) return ""
  return filename.slice(((filename.lastIndexOf(".") - 1) >>> 0) + 2).toLowerCase()
}

function _getType(type) {
  switch (type) {
    case "mp4":
      return "video/mp4"
    case "webm":
      return "video/webm"
    case "ogg":
      return "video/ogg"
    case "avi":
      return "video/webm"
    case "mov":
      return "video/mp4"
    default:
      return "video/mp4"
  }
}

onBeforeUnmount(() => {
  if (player) {
    player.destroy()
    player = null
  }
})
</script>

<script>
export default {
  name: "YVideo"
}
</script>

<style scoped lang="scss">
.preview-video {
  height: 100%;
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  position: relative;
  flex-direction: column;

  #dplayer {
    width: 100%;
    height: auto;
    max-height: 100%;
    transition: all 0.3s ease;
    border-radius: 6px;
    overflow: hidden;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }

  &.maximized #dplayer {
    width: 100%;
    height: calc(100% - 50px);
    max-width: 100% !important;
  }

  .native-player-container {
    width: 100%;
    height: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
    flex-direction: column;

    #native-video-container {
      width: 100%;
      max-width: 900px;
      padding: 20px;
    }

    #native-video {
      width: 100%;
      border-radius: 6px;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    }

    .native-video-tools {
      display: flex;
      justify-content: center;
      margin-top: 15px;
      gap: 10px;

      .download-btn {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        padding: 6px 15px;
        background-color: #42b983;
        color: white;
        text-decoration: none;
        border-radius: 4px;
        font-size: 14px;
        transition: all 0.2s ease;

        &:hover {
          background-color: #3da976;
        }

        svg {
          margin-right: 5px;
        }
      }
    }
  }

  .error-message {
    color: red;
    text-align: center;
    padding: 20px;
    background: rgba(255, 0, 0, 0.05);
    border-radius: 6px;
    font-weight: 500;
  }

  .error-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 30px;
    background: rgba(255, 240, 240, 0.2);
    border-radius: 8px;
    border: 1px solid rgba(255, 77, 79, 0.2);
    width: 90%;
    max-width: 400px;

    .error-icon {
      margin-bottom: 15px;
    }

    .error-title {
      font-size: 18px;
      font-weight: 600;
      color: #ff4d4f;
      margin-bottom: 12px;
    }

    .error-message-details {
      font-size: 14px;
      color: #666;
      text-align: center;
      margin-bottom: 20px;
    }

    .error-actions {
      display: flex;
      gap: 15px;

      .error-btn {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        padding: 6px 16px;
        border-radius: 4px;
        border: none;
        font-size: 14px;
        cursor: pointer;
        transition: all 0.2s;

        &.retry-btn {
          background-color: #42b983;
          color: white;

          &:hover {
            background-color: #3da976;
          }
        }

        &.alt-player-btn {
          background-color: #f0f0f0;
          color: #666;

          &:hover {
            background-color: #e0e0e0;
            color: #333;
          }
        }
      }
    }
  }

  .loading-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 200px;
    width: 100%;

    .loading-spinner {
      width: 40px;
      height: 40px;
      border: 4px solid rgba(66, 185, 131, 0.1);
      border-radius: 50%;
      border-top-color: #42b983;
      animation: spin 1s infinite linear;
    }

    .loading-text {
      margin-top: 12px;
      color: #666;
      font-size: 14px;
    }
  }

  .video-tools {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 12px;
    padding: 8px 12px;
    background-color: #f8f8f8;
    border-radius: 6px;
    width: 100%;
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.05);

    .tool-group {
      display: flex;
      gap: 5px;
    }

    .tool-btn {
      background: none;
      border: 1px solid #ddd;
      border-radius: 4px;
      padding: 5px 8px;
      cursor: pointer;
      color: #555;
      font-size: 13px;
      transition: all 0.2s ease;

      &:hover {
        background-color: #eee;
        color: #333;
        border-color: #ccc;
      }

      &.active {
        background-color: #42b983;
        color: white;
        border-color: #42b983;
      }

      svg {
        vertical-align: middle;
      }
    }
  }
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

@media screen and (max-width: 768px) {
  .preview-video {
    #dplayer {
      width: 100%;
      max-width: 100%;
    }

    .video-tools {
      flex-wrap: wrap;
      gap: 5px;

      .tool-group {
        flex-wrap: wrap;
      }
    }
  }
}

:deep(.dplayer-controller) {
  .dplayer-bar-wrap {
    height: 8px !important;

    &:hover .dplayer-bar-time {
      opacity: 1;
    }
  }

  .dplayer-mobile-play {
    width: 50px;
    height: 50px;
  }
}

:deep(.dplayer-video-wrap) {
  video {
    transition: transform 0.3s ease;
    transform-origin: center center;
  }
}
</style>
