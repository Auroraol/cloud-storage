<template>
  <el-image
    style="width: 100px; height: 100px"
    :src="url"
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

// 计算属性
const url = computed(() => props.url || "")
const srcList = computed(() => [url.value])
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
