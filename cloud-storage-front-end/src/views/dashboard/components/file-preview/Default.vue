<template>
  <div>
    <n-empty description="当前文件类型暂不支持预览">
      <template #icon>
        <c-icon name="user" />
      </template>
    </n-empty>
  </div>
</template>

<script setup>
const _props = defineProps({
  resource: Object
})
</script>

<script>
export default {
  name: "Default"
}
</script>
