<template>
  <div class="preview">
    <Transition name="preview-fade" mode="out-in">
      <div class="preview-content" ref="previewRef" v-show="showPreview">
        <div class="preview-title">
          {{ props.resource.name }}
          <div class="preview-actions">
            <button class="action-btn download-btn" @click="downloadFile" title="下载文件">
              <svg class="icon" viewBox="0 0 1024 1024" width="16" height="16">
                <path
                  d="M505.7 661c3.2 4.1 9.4 4.1 12.6 0l112-141.7c4.1-5.2 0.4-12.9-6.3-12.9h-74.1V168c0-4.4-3.6-8-8-8h-60c-4.4 0-8 3.6-8 8v338.3H400c-6.7 0-10.4 7.7-6.3 12.9l112 141.8z"
                  fill="currentColor"
                />
                <path
                  d="M878 626h-60c-4.4 0-8 3.6-8 8v154H214V634c0-4.4-3.6-8-8-8h-60c-4.4 0-8 3.6-8 8v198c0 17.7 14.3 32 32 32h684c17.7 0 32-14.3 32-32V634c0-4.4-3.6-8-8-8z"
                  fill="currentColor"
                />
              </svg>
            </button>
            <button class="action-btn maximize-btn" @click="toggleMaximize" title="最大化">
              <svg v-if="!maximize" class="icon" viewBox="0 0 1024 1024" width="16" height="16">
                <path
                  d="M853.333333 170.666667H170.666667v682.666666h682.666666V170.666667z m-42.666666 640H213.333333V213.333333h597.333334v597.333334z"
                  fill="currentColor"
                />
                <path
                  d="M682.666667 426.666667H341.333333V341.333333h341.333334v85.333334zM384 682.666667h-42.666667v-42.666667h42.666667v42.666667z m85.333333 0h-42.666666v-42.666667h42.666666v42.666667z m85.333334 0h-42.666667v-42.666667h42.666667v42.666667z m85.333333 0h-42.666667v-42.666667H640v42.666667z"
                  fill="currentColor"
                />
              </svg>
              <svg v-else class="icon" viewBox="0 0 1024 1024" width="16" height="16">
                <path
                  d="M563.2 353.6l107.52-107.52L711.68 287.36l107.52-107.52v194.56H563.2v-20.8zM460.8 670.4l-107.52 107.52-40.96-40.96L204.8 844.48V649.6h256v20.8z"
                  fill="currentColor"
                />
                <path
                  d="M204.8 241.28L312.32 348.8l-40.96 40.96 107.52 107.52H204.8V241.28z m501.12 542.08l107.52-107.52 40.96 40.96 107.52-107.52V844.8h-256v-61.44z"
                  fill="currentColor"
                />
              </svg>
            </button>
          </div>
        </div>
        <div class="preview-body">
          <component
            v-if="url"
            style="width: 100%; height: 100%"
            :is="viewType"
            :resource="props.resource"
            :maximize="maximize"
            :url="url"
          />
        </div>
      </div>
    </Transition>
  </div>
</template>
<script lang="ts">
export default {
  name: "Preview"
}
</script>
<script setup lang="ts">
import { computed, ref } from "vue"
// 正确导入组件
import Default from "./Default.vue"
import YImage from "./Image.vue"
import YVideo from "./Video.vue"
import YAudio from "./Audio.vue"
import YOffice from "./Office.vue"
import YText from "./Text.vue"
// import { preview } from "@/http/Explore"
// import store from "@/store/temp"

const previewRef = ref(null)
const maximize = ref(false)
const showPreview = ref(true)

// 添加一个函数来切换最大化状态
function toggleMaximize() {
  maximize.value = !maximize.value

  if (maximize.value && previewRef.value) {
    // 进入全屏模式
    const previewElement = previewRef.value as HTMLElement
    previewElement.classList.add("fixed")

    // 创建遮罩层
    const overlay = document.createElement("div")
    overlay.classList.add("overlay")

    // 点击遮罩层关闭全屏预览
    overlay.addEventListener("click", (e) => {
      if (e.target === overlay) {
        toggleMaximize()
      }
    })

    // 添加ESC键关闭全屏预览
    const handleEscKey = (e: KeyboardEvent) => {
      if (e.key === "Escape" && maximize.value) {
        toggleMaximize()
      }
    }
    document.addEventListener("keydown", handleEscKey)

    // 存储引用，用于移除
    overlay.dataset.escKeyHandler = "true"
    document.body.appendChild(overlay)

    // 添加过渡动画类
    setTimeout(() => {
      previewElement.classList.add("animation-done")
    }, 50)

    // 在移动设备上，禁止背景滚动
    document.body.style.overflow = "hidden"
  } else if (previewRef.value) {
    // 退出全屏模式
    const previewElement = previewRef.value as HTMLElement
    previewElement.classList.remove("animation-done")

    // 使用setTimeout确保动画效果完成后再移除fixed类
    setTimeout(() => {
      previewElement.classList.remove("fixed")

      // 移除遮罩层
      const overlay = document.querySelector(".overlay")
      if (overlay) overlay.remove()

      // 移除ESC键事件监听器
      document.removeEventListener("keydown", (e: KeyboardEvent) => {
        if (e.key === "Escape") {
          toggleMaximize()
        }
      })

      // 恢复背景滚动
      document.body.style.overflow = ""
    }, 200)
  }
}

// 添加下载文件功能
function downloadFile() {
  if (!url.value) return

  // 创建一个临时链接
  const a = document.createElement("a")
  a.href = url.value
  a.download = props.resource.name || `文件.${props.resource.type}`
  // 不要设置 target="_blank" 以避免在新标签页打开
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
}

const props = defineProps({
  resource: {
    type: Object,
    required: true
  }
})

const fileTypes = {
  image: ["jpg", "jpeg", "png", "gif", "bmp", "webp"],
  video: ["mp4", "webm", "ogg"],
  audio: ["mp3", "wav", "ogg"],
  office: ["doc", "docx", "xls", "xlsx", "ppt", "pptx", "pdf"],
  text: ["txt", "md", "json", "js", "css", "html"]
}

// 计算属性
const url = computed(() => props.resource.url || "")
const viewType = computed({
  get: () => {
    let type = props.resource.type
    if (type.indexOf(".") !== -1) {
      type = type.substring(type.indexOf(".") + 1)
    }
    let v = getView(type)
    if (v === Default) {
      const dot = props.resource.name.indexOf(".")
      if (dot === -1) return v
      type = props.resource.name.substring(dot + 1)
      v = getView(type)
    }
    return v
  },
  set: () => {}
})

// 获取相应的预览组件
function getView(type: string) {
  console.log("type: ", type)
  if (fileTypes.image.includes(type)) {
    return YImage
  } else if (fileTypes.video.includes(type)) {
    return YVideo
  } else if (fileTypes.audio.includes(type)) {
    return YAudio
  } else if (fileTypes.office.includes(type)) {
    return YOffice
  } else if (fileTypes.text.includes(type)) {
    return YText
  } else {
    return Default
  }
}

// 导出功能
defineExpose({
  toggleMaximize,
  downloadFile
})
</script>

<style scoped lang="scss">
.preview {
  width: 100%;
  height: 100%;
  padding: 10px;
  background-color: #f7f7f7;
  display: flex;
  box-sizing: border-box;
  color: #555;
  font-weight: 700;
  justify-content: center;
  align-items: center;
  position: relative;
  overflow: hidden;

  .preview-content {
    width: 100%;
    height: 100%;
    padding: 0 25px;
    position: relative;
    top: 0;
    left: 0;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    transform: translateY(0);
    opacity: 1;

    .preview-title {
      height: 30px;
      display: none;
      font-size: 18px;
      font-weight: 700;
      text-align: center;
      position: relative;
      margin-bottom: 10px;
      transition: opacity 0.3s ease;

      .preview-actions {
        position: absolute;
        right: 10px;
        top: 0;

        .action-btn {
          background: none;
          border: none;
          cursor: pointer;
          padding: 5px;
          margin-left: 10px;
          color: #666;
          font-size: 16px;
          transition:
            color 0.2s ease,
            transform 0.2s ease;

          &:hover {
            color: #000;
            transform: scale(1.1);
          }

          .icon {
            display: inline-block;
            vertical-align: middle;
          }
        }
      }
    }
    .preview-body {
      width: 100%;
      height: 100%;
      transition: height 0.3s ease;
    }
  }
  .preview-content.fixed {
    position: fixed;
    top: 50% !important;
    left: 50% !important;
    transform: translate(-50%, -50%) scale(0.95);
    z-index: 999;
    width: 90vw;
    height: 90vh;
    padding: 15px 20px;
    background-color: #fff;
    border-radius: 10px;
    box-shadow: 0 5px 25px rgba(0, 0, 0, 0.2);
    opacity: 0.95;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

    .preview-title {
      display: block;
      animation: fadeIn 0.3s ease forwards;
    }
    .preview-body {
      height: calc(100% - 40px);
      overflow-y: auto;
    }
  }

  .preview-content.fixed.animation-done {
    transform: translate(-50%, -50%) scale(1);
    opacity: 1;
  }

  .overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background-color: rgba(0, 0, 0, 0.6);
    z-index: 998;
    backdrop-filter: blur(2px);
    animation: fadeIn 0.3s ease forwards;
    cursor: pointer;
  }
}

/* 淡入淡出过渡效果 */
.preview-fade-enter-active,
.preview-fade-leave-active {
  transition: opacity 0.3s ease;
}

.preview-fade-enter-from,
.preview-fade-leave-to {
  opacity: 0;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes scaleIn {
  from {
    transform: translate(-50%, -50%) scale(0.9);
  }
  to {
    transform: translate(-50%, -50%) scale(1);
  }
}

/* 响应式调整 */
@media screen and (max-width: 768px) {
  .preview {
    .preview-content.fixed {
      width: 100vw;
      height: 100vh;
      padding: 10px;
      border-radius: 0;
    }
  }
}
</style>
