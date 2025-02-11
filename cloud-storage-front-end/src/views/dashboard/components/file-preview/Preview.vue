<template>
  <div class="preview">
    <Transition name="preview">
      <div class="preview-content" ref="previewRef">
        <div class="preview-title">
          {{ props.resource.name }}
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
import { computed, watch, ref } from "vue"
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

    .preview-title {
      height: 30px;
      display: none;
      font-size: 18px;
      font-weight: 700;
      text-align: center;
    }
    .preview-body {
      width: 100%;
      height: 100%;
    }
  }
  .preview-content.fixed {
    position: fixed;
    top: 50% !important;
    left: 50% !important;
    transform: translate(-50%, -50%) !important;
    z-index: 999;
    width: 80vw;
    height: calc(80vh + 30px);
    padding: 10px 20px;
    background-color: #fff;
    border-radius: 10px;

    .preview-title {
      display: block;
    }
    .preview-body {
      height: calc(100% - 30px);
      overflow-y: auto;
    }
  }

  .overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background-color: rgba(0, 0, 0, 0.5);
    z-index: 998;
  }

  .maximize-btn {
    position: absolute;
    bottom: 0;
    right: 0;
    padding: 5px;
  }
}
// .preview-enter-from {
//   top: 0;
// }

// .preview-leave-to {
//   top: 0;
// }

// .preview-enter-from ,
// .preview-leave-to {
//   top: unset;
//   bottom: 0;
//   left: 0;
// }
</style>
