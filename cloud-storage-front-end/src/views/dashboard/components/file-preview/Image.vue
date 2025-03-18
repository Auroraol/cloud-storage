<template>
  <el-image
    style="width: 100px; height: 100px"
    :src="imageUrl"
    :alt="props.resource.name"
    :preview-src-list="srcList"
    :zoom-rate="1.2"
    :max-scale="7"
    :min-scale="0.2"
    show-progress
    :initial-index="4"
    fit="cover"
  />
</template>

<script setup>
import { computed, onMounted } from "vue"

const props = defineProps({
  resource: {
    type: Object,
    required: true
  },
  url: {
    type: String,
    required: true
  }
})

onMounted(() => {
  console.log("Image组件挂载，URL:", props.url)
})

// 计算属性
const imageUrl = computed(() => props.url || "")
const srcList = computed(() => [imageUrl.value])
</script>

<script>
export default {
  name: "YImage"
}
</script>

<style scoped lang="scss">
.demo-image__error .image-slot {
  font-size: 30px;
}
.demo-image__error .image-slot .el-icon {
  font-size: 30px;
}
.demo-image__error .el-image {
  width: 100%;
  height: 200px;
}
</style>
