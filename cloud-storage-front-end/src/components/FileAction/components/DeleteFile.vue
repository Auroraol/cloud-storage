<script lang="ts" setup>
import { useRoute, useRouter } from "vue-router"

const route = useRoute()
const router = useRouter()

router.replace({ path: "/" + route.params.path, query: route.query })
</script>

<template>
  <div />
</template>

<!-- <template>
  <Dialog
    :show="dialogConfig.show"
    :title="dialogConfig.title"
    :buttons="dialogConfig.buttons"
    :show-cancel="true"
    width="400px"
    @close="dialogConfig.show = false"
  >
    <div class="delete-info">
      <Icon icon-name="warning" :width="40" />
      <div class="info">
        <span class="text">确定删除所选的文件吗？</span>
        <span class="text">删除的文件可在 10天 内通过回收站还原</span>
      </div>
    </div>
  </Dialog>
</template>

<script setup>
import Dialog from "@/components/FileAction/components/Dialog.vue"
import Icon from "@/components/SvgIcon/index.vue"
import { ElMessage } from "element-plus"
import { isArray } from "@/utils/is"
import { reactive } from "vue"
import { useFileStore } from "@/stores/file"

const emit = defineEmits(["delete-success"])
const fileStore = useFileStore()
let delIds = ""

const dialogConfig = reactive({
  show: false,
  title: "确定删除",
  showCancel: true,
  buttons: [
    {
      type: "primary",
      text: "删除",
      click: (e) => {
        deleteFile()
      }
    }
  ]
})

function show(data) {
  if (isArray(data)) {
    delIds = data.map((item) => item.id).join(",")
  } else {
    delIds = data.id
  }
  dialogConfig.show = true
}

const deleteFile = async () => {
  try {
    await fileStore.deleteFile(delIds)
    emit("delete-success")
    ElMessage.success("删除成功")
  } catch (error) {
    ElMessage.error("删除失败")
  }
  dialogConfig.show = false
}

defineExpose({
  show
})
</script>
<style lang="scss" scoped>
.delete-info {
  display: flex;
  align-items: center;
  .info {
    display: flex;
    flex-direction: column;
    margin-left: 16px;
  }
}
</style> -->
