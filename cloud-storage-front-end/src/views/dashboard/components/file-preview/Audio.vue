<template>
  <div class="preview-audio">
    <div v-if="!props.url" class="error-message">无法获取音频URL</div>
    <AudioPlayer
      v-else
      class="audio-player"
      :option="{
        src: props.url,
        title: getTitle,
        coverImage: 'https://img2.baidu.com/it/u=3895119537,2684520677&fm=253&fmt=auto&app=138&f=JPEG?w=500&h=500'
      }"
    />
  </div>
</template>

<script setup>
import AudioPlayer from "vue3-audio-player"
import "vue3-audio-player/dist/style.css"
import { computed, onMounted } from "vue"

const props = defineProps({
  resource: Object,
  url: String
})

onMounted(() => {
  console.log("Audio组件挂载，URL:", props.url)
  console.log("资源名称:", props.resource?.name)
})

const getTitle = computed({
  get: () => {
    const name = props.resource.name
    if (!name) return ""
    const dot = name.indexOf(".")
    return dot === -1 ? name : name.substr(0, dot)
  }
})
</script>

<script>
export default {
  name: "YAudio"
}
</script>

<style scoped lang="scss">
.preview-audio {
  width: 90%;
  height: 100%;
  margin: auto;
  max-width: 500px;
  display: flex;
  align-items: center;
  flex-direction: column;
  justify-content: center;

  .audio-player {
    height: 170px;
    width: 100%;
    padding: 10px 20px;
    background-color: #fff;
    border-radius: 10px;
    box-sizing: border-box;
  }

  .error-message {
    color: red;
    text-align: center;
  }
}

:deep(.audio__player-play-cont) {
  display: flex;
  justify-content: center;
}

:deep(.audio__player-play-icon) {
  cursor: pointer;
}
</style>
