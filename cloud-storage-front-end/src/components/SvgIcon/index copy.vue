<!--
  自定义svg图片组件
  使用方法：在阿里云iconfont下载需要图标并命名文件名为英文名，将命名好的文件放到icons文件夹下svg文件夹下，
  如logo图标使用：
  <svg-icon name="logo"  color="red" width="100px" height="100px" />
  这里的name即为logo文件名, 可自定义color, width, height

  <div class="svg-icon-container">
    <svg-icon name="fasong"> </svg-icon>
    <span>发送</span>
  </div>
  .svg-icon-container {
  display: flex; /* 使用 Flex 布局 */
  align-items: center; /* 垂直居中对齐 */
}

-->
<!-- src\components\SvgIcon.vue -->
<script lang="ts" setup>
defineProps({
  //xlink:href属性值的前缀
  prefix: {
    type: String,
    default: "#icon-"
  },
  // 接收父组件传来的svg矢量图的名字
  name: String,
  // 接收父组件传来的图标颜色
  color: {
    type: String,
    default: ""
  },
  // 接收父组件传来的图标宽度
  width: {
    type: String,
    default: "1em"
  },
  // 接收父组件传来的图标高度
  height: {
    type: String,
    default: "1em"
  }
})
</script>

<template>
  <div>
    <svg :style="{ width, height }">
      <use :xlink:href="prefix + name" :fill="color" />
    </svg>
  </div>
</template>

<!-- <style lang="scss" scoped>
.svg-icon {
  width: 1em;
  height: 1em;
  overflow: hidden;
} -->
<!-- </style> -->
