<template>
  <div class="app-container center">
    <el-empty description="欢迎来到 admin 角色专属首页" />
  </div>
</template>

<script lang="ts" setup>
import { nextTick, reactive, ref } from "vue"

const emit = defineEmits(["addFile"])
const dataTableRef = ref()
const editNameRef = ref()
const editing = ref(false)
const currentFolder = ref({ id: "0" })
const params = reactive({
  page: 1 as number,
  pageSize: 20 as number,
  filePid: "0" as string,
  category: "all" as string
})

const fileList = ref([
  { name: "文件夹1", updateTime: "2023-10-01 14:49:56", size: "-" },
  { name: "文件2", updateTime: "2023-10-01 14:48:57", size: "-" },
  { name: "文件3", updateTime: "2023-10-01 12:17:57", size: "-" },
  { name: "文件夹4", updateTime: "2023-10-01 11:11:11", size: "-" },
  { name: "文件5", updateTime: "2023-10-01 14:48:06", size: "0 字节" },
  { name: "文件6", updateTime: "2023-10-01 12:10:48", size: "4.8 MB" }
])

const reloadTable = () => {
  dataTableRef.value?.reload()
}

const newFolder = () => {
  if (editing.value) return
  dataTableRef.value.setRowFieldValue("showEdit", false)
  editing.value = true
  dataTableRef.value.unshiftRow({
    showEdit: true,
    folderType: 1,
    id: "",
    filePid: currentFolder.value.id,
    filenameReal: "",
    updateTime: new Date().toLocaleString()
  })
  nextTick(() => {
    editNameRef.value.focus()
  })
}

defineExpose({ reloadTable })
</script>
<style lang="scss" scoped>
.center {
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}
.main {
  padding-top: 20px;
}
.main-header {
  height: 40px;
  padding: 4px 24px;
  .toolbar {
    display: flex;
    align-items: center;
    .toolbar-action {
      width: 100%;
      .toolbar-group {
        display: flex;
        color: #fff;
        .el-icon {
          margin-right: 6px;
        }
        .action-list {
          :deep(button) {
            background-color: #f0faff;
            padding: 0 12px;
          }
        }
      }
      .action-main {
        margin-right: 16px;
      }
    }
    .toolbar-customize {
      .search {
        width: 270px;
        :deep(.el-input__wrapper) {
          border-radius: 18px;
          padding: 1px 15px;
        }
        .text {
          cursor: pointer;
        }
      }
    }
  }
}

.main-body {
  height: calc(100% - 40px);
  .nav {
    width: 100%;
    padding: 0 0 0 24px;
    box-sizing: border-box;
    height: 40px;
    .nav-list {
      height: 40px;
      line-height: 40px;
      position: relative;
      border-radius: 4px 4px 0 0;
      overflow: hidden;
      .nav-left {
        float: left;
        color: #03081a;
      }
    }
  }
  .file-list {
    padding-left: 8px;
    height: calc(100% - 40px);
  }
  .table-empty {
    height: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
    overflow: auto;
  }
}
</style>
