<template>
  <span :class="customClass" :style="{ width: width + 'px', height: computedHeight }">
    <img :src="getImage()" :style="{ 'object-fit': fit }" />
  </span>
</template>

<!--
使用:<template>
    写一下参数说明
    fileType: 文件类型(0 目录 1 视频 2 音频 3 图片 4 pdf 5 doc 6 excel 7 纯文本 8 代码 9 压缩包 10 其他文件)
    cover: 封面
    width: 宽度
    height: 高度
    fit: 图片适应方式
    customClass: 自定义类名
-->

<script>
export default {
  name: "Icon"
}
</script>
<script setup>
import { useGlobalSetting } from "@/settings/setting"
import fileSetting from "@/settings/fileSetting"
import { isString } from "@/utils/is"
import { computed } from "vue"

const props = defineProps({
  fileType: {
    type: Number
  },
  iconName: {
    type: String
  },
  cover: {
    type: String
  },
  width: {
    type: Number,
    default: 32
  },
  height: {
    type: [Number, String],
    default: null
  },
  fit: {
    type: String,
    default: "cover"
  },
  customClass: {
    type: String,
    default: "icon"
  }
})

const globalSetting = useGlobalSetting()
const imageUrl = globalSetting.imgUrl

const computedHeight = computed(() => {
  if (props.height) {
    if (isString(props.height)) {
      return props.height
    } else {
      return props.height + "px"
    }
  } else {
    return props.width + "px"
  }
})
const { fileTypeMap } = fileSetting
const getImage = () => {
  if (props.cover) {
    // 视频，图片封面
    return imageUrl + props.cover
  }
  let icon = "unknow_icon"
  if (props.iconName) {
    icon = props.iconName
  } else {
    const iconMap = fileTypeMap[props.fileType]
    if (iconMap != undefined) {
      icon = iconMap["icon"]
    }
  }
  return new URL(`/src/assets/icon-image/${icon}.png`, import.meta.url).href
}
</script>
<style lang="scss" scoped>
.icon {
  text-align: center;
  display: inline-block;
  border-radius: 3px;
  // overflow: hidden;
}
img {
  display: block;
  width: 100%;
  height: 100%;
}
</style>
