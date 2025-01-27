<template>
  <div class="app-container">
    <div class="recycle">
      <div class="recycle-header">
        <div class="actions">
          <el-button type="primary" round class="clear-button">
            <el-icon><ShoppingCart /></el-icon>
            清空回收站
          </el-button>
          <div class="action-more" v-show="selectFiles.length > 0">
            <FileAction :value="selectFiles" :actions="actions" />
          </div>
        </div>
      </div>
      <div class="nav">
        <span class="txt">回收站</span>
        <span class="count">已加载{{ tableDataNum }}个</span>
      </div>
      <div class="recycle-body">
        <div class="file-list">
          <vue-good-table
            ref="recycleTableRef"
            :columns="canEditColumns"
            :rows="tableRows"
            :select-options="{
              enabled: true,
              selectOnCheckboxOnly: true,
              selectionText: '选中的行数',
              clearSelectionText: ''
            }"
            :pagination-options="{
              enabled: true,
              mode: 'pages',
              perPage: 10,
              perPageDropdown: [10, 20, 50],
              dropdownAllowAll: false,
              nextLabel: '下一页',
              prevLabel: '上一页',
              rowsPerPageLabel: '每页显示',
              pageLabel: '页',
              ofLabel: '/'
            }"
            v-on:row-click="handleTableRowSelect"
            styleClass="vgt-table condensed"
          >
            <template #table-row="props">
              <span v-if="props.column.field === 'filename'">
                <div class="file-item">
                  <Icon
                    v-if="(props.row.fileType === 3 || props.row.fileType === 1) && props.row.status === 2"
                    :cover="props.row.fileCover"
                    :width="32"
                  />
                  <Icon v-else-if="props.row.folderType === 0" :file-type="props.row.fileType" />
                  <Icon v-else icon-name="folder_2" :width="40" />
                  <span class="filename">{{ props.row.filename }}</span>
                </div>
              </span>
            </template>
            <template #emptystate>
              <el-empty description="空空如也" />
            </template>
          </vue-good-table>
        </div>
      </div>
      <FileDialog
        :show="dialogConfig.show"
        :title="dialogTitle"
        :buttons="dialogConfig.buttons"
        :show-cancel="true"
        width="400px"
        @close="dialogConfig.show = false"
      >
        <div class="delete-info">
          <div class="info">
            <span class="text">{{ dialogContent }}</span>
          </div>
        </div>
      </FileDialog>
    </div>
  </div>
</template>

<script lang="ts">
export default {
  name: "Recycle"
}
</script>

<script setup lang="ts">
import { computed, ref, reactive } from "vue"
import { columns } from "./columns"
import { FileAction } from "@/components/FileAction/index"
import Icon from "@/components/FileIcon/Icon.vue"
import { isArray } from "@/utils/is"
import { ElMessage } from "element-plus"
import FileDialog from "@/components/FileDialog/index.vue"
import { VueGoodTable } from "vue-good-table-next"

const recycleTableRef = ref()
const selectFiles = ref([])
const canEditColumns = ref(columns)
const tableDataNum = ref(0)
const flag = ref(1) // 0恢复，1彻底删除
const ids = ref()

interface RecycleItem {
  id: string
  filePid: string
  fileSize: number
  filename: string
  fileCover: string
  folderType: number
  fileCategory: number
  fileType: number
  status: number
  updateTime: Date
  recoveryTime: Date
}

const dialogTitle = computed(() => {
  return flag.value ? "彻底删除" : "确认还原"
})

const dialogContent = computed(() => {
  return flag.value ? "文件删除后将无法恢复，您确认要彻底删除所选文件吗？" : "确认还原选中的文件？"
})

const actions = [
  {
    title: "还原",
    label: "copyLink",
    showLabel: true,
    icon: "recovery",
    onClick: (data) => {
      showDialog(data, 0)
    }
  },
  {
    title: "删除",
    label: "delete",
    showLabel: true,
    icon: "delete",
    onClick: (data) => {
      showDialog(data, 1)
    }
  },
  {
    title: "批量删除",
    label: "batchDelete",
    showLabel: true,
    icon: "delete",
    onClick: (data) => {
      showBatchDeleteDialog(data)
    }
  }
]

const dialogConfig = reactive({
  show: false,
  showCancel: true,
  buttons: [
    {
      type: "primary",
      text: "确认",
      click: () => {
        confirmDothing()
      }
    }
  ]
})

function showDialog(data, f) {
  console.log(data)
  flag.value = f
  if (!isArray(data)) {
    ids.value = data.id
  } else {
    ids.value = data.map((item) => item.id).join(",")
  }
  dialogConfig.show = true
}

function reloadTable() {
  recycleTableRef.value?.reload()
}

function fetchSuccess({ totalData }) {
  tableDataNum.value = totalData
}

function handleTableRowSelect(params) {
  selectFiles.value = [params.row] // 将选中的行作为数组传递
}

// // 计算属性来处理选中的文件
// const selectedFileIds = computed(() => {
//   return selectFiles.value.map((file) => file.fileId).join(",")
// })

const confirmDothing = async () => {
  if (flag.value) {
    try {
      await delFileApi(ids.value)
      ElMessage.success("文件已彻底删除")
    } catch (error) {
      ElMessage.error("操作失败")
    }
  } else {
    try {
      await recoveryFileApi(ids.value)
      ElMessage.success("文件已恢复")
    } catch (error) {
      ElMessage.error("操作失败")
    }
  }
  dialogConfig.show = false
  reloadTable()
}

const RecycleList = ref<RecycleItem[]>([
  {
    id: "1",
    filePid: "1",
    fileSize: 1024,
    filename: "测试文件1",
    fileCover: "",
    folderType: 0,
    fileCategory: 1,
    fileType: 1,
    status: 2,
    updateTime: new Date("2023-01-01"),
    recoveryTime: new Date("2023-01-01")
  },
  {
    id: "2",
    filePid: "2",
    fileSize: 2048,
    filename: "测试文件2",
    fileCover: "",
    folderType: 0,
    fileCategory: 2,
    fileType: 2,
    status: 1,
    updateTime: new Date("2023-01-01"),
    recoveryTime: new Date("2023-01-01")
  }
])

const loadDataTable = (res: any) => {
  // return loadShareList(res)
}

const formatDate = (date: string): string => {
  if (!date) return "-"
  try {
    const dateObj = new Date(date)
    const year = dateObj.getFullYear()
    const month = String(dateObj.getMonth() + 1).padStart(2, "0")
    const day = String(dateObj.getDate()).padStart(2, "0")
    const hours = String(dateObj.getHours()).padStart(2, "0")
    const minutes = String(dateObj.getMinutes()).padStart(2, "0")
    const seconds = String(dateObj.getSeconds()).padStart(2, "0")

    return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
  } catch {
    return "-"
  }
}

const tableRows = computed(() => {
  return RecycleList.value.map((file) => ({
    ...file,
    updateTime: formatDate(file.updateTime.toISOString()),
    recoveryTime: formatDate(file.recoveryTime.toISOString())
  }))
})

function showBatchDeleteDialog(data) {
  console.log("批量删除的文件:", data)
  // 这里可以调用删除 API
}
</script>

<style lang="scss" scoped>
.recycle {
  height: 100%;
  width: 100%;
  // min-width: 750px;

  .recycle-header {
    height: 50px;
    padding-left: 4px;
    font-size: 12px;
    display: flex;
    justify-content: space-between;
    align-items: center;

    .top {
      display: inline-block;
      line-height: 36px;
    }

    .actions {
      display: flex;
      align-items: center;

      .clear-button {
        // margin-right: 16px;
        font-weight: bold;
      }

      .action-more {
        // margin-left: 16px;
      }
    }
  }

  .nav {
    font-size: 12px;
    color: #03081a;
    padding-left: 4px;
    display: flex;
    justify-content: space-between;

    .txt {
      color: #03081a;
      font-weight: 600;
      font-size: 18px;
      margin-bottom: 15px;
    }
    .count {
      font-size: 18px;
      color: #afb3bf;
    }
  }

  .recycle-body {
    height: calc(100% - 75px);
    // background-color: #fff;
    // padding: 10px;

    .file-list {
      // padding: 0 20px;
    }
  }
}
.vgt-table {
  // min-height: 450px;

  tr {
    height: 50px !important; // 强制设置行高
  }

  td {
    padding: 10px !important; // 强制设置单元格内边距
    vertical-align: middle !important;
    text-align: center; // 设置水平对齐方式为居中
  }

  .file-item {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .filename {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}

// .left-align {
//   text-align: left; // 左对齐
// }

// .right-align {
//   text-align: right; // 右对齐
// }
</style>
