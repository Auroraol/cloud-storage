<template>
  <!-- <div class="app-container center">
    <el-table :data="fileList" style="width: 100%">
      <el-table-column prop="name" label="名称" width="180" />
      <el-table-column prop="updateTime" label="更新时间" width="180" />
      <el-table-column prop="size" label="大小" width="180" />
    </el-table>
  </div> -->
  <!-- <div class="app-container center"> -->
  <!-- {{ canEditColumns }} -->
  <!--     :dataSource="loadDataTable" -->
  <BasicTable
    ref="dataTableRef"
    :columns="canEditColumns"
    :row-key="(row) => row.id"
    :dataSource="loadDataTable"
    :loading="true"
    @select-change="handleTableRowSelect"
  >
    <template #filename="{ index, row }">
      <div class="file-item">
        <template v-if="(row.fileType == 3 || row.fileType == 1) && row.status == 2">
          <!-- 视频，图片，且转码成功才有封面 -->
          <!-- <Icon :cover="row.fileCover" :width="32" /> -->
        </template>
        <template v-else>
          <!-- <Icon v-if="row.folderType == 0" :file-type="row.fileType" />
          文件夹 -->
          <!--<Icon v-else="row.folderType == 1" :file-type="0" /> -->
        </template>
        <span class="file-name" v-if="!row.showEdit">
          <span class="text" @click.stop="openFile(row)">{{ row.filename }}</span>
        </span>
        <!-- 编辑框 -->
        <div class="edit-panel" v-if="row.showEdit" @click.stop="">
          <el-input
            v-model.trim="row.filenameReal"
            ref="editNameRef"
            size="small"
            :maxLength="190"
            @keyup.enter="saveNameEdit(index)"
          >
            <template #suffix>
              {{ row.fileSuffix }}
            </template>
          </el-input>
          <div class="edit-action is-confirm">
            <span class="iconfont icon-check" @click="saveNameEdit(index)" />
          </div>
          <div class="edit-action is-cancel">
            <span class="iconfont icon-close" @click="cancelNameEdit(index)" />
          </div>
        </div>
      </div>
    </template>
    <template #RowAction="{ index, row }">
      <FileAction mode="row" :value="row" field="updateTime" :key="row.showEdit" :actions="actions" :offset="-130" />
    </template>
    <template #fileSize="{ index, row }">
      <div class="file-size">
        <!-- <span v-if="row.fileSize">{{ size2Str(row.fileSize) }}</span>
        <span v-else> - </span> -->
      </div>
    </template>
    <template #empty>
      <div class="table-empty">
        <div class="main-empty">
          <div class="u-empty">
            <div class="empty-image">
              <Icon icon-name="empty" :width="120" height="auto" />
            </div>
            <div class="empty-description">
              <span class="text">当前列表为空，上传你的第一个文件吧</span>
            </div>
            <div class="empty-bottom">
              <div class="empty-action">
                <!-- <el-upload
                  :show-file-list="false"
                  :with-credentials="true"
                  :multiple="true"
                  :http-request="addFile"
                  :accept="fileAccept"
                > -->
                <el-upload :show-file-list="false" :with-credentials="true" :multiple="true">
                  <div class="action-item">
                    <div class="action-item-image">
                      <Icon icon-name="empty_upload" :width="42" />
                    </div>
                    <div class="action-item-text">上传文件</div>
                  </div>
                </el-upload>

                <!-- <div class="action-item" @click="newFolder" v-show="params.category == 'all'">
                  <div class="action-item-image">
                    <Icon icon-name="empty_newFolder" :width="42" />
                  </div>
                  <div class="action-item-text">新建文件夹</div>
                </div> -->
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>
  </BasicTable>
  <!-- </div> -->
</template>

<script lang="ts" setup>
import { nextTick, reactive, ref } from "vue"
import { ElMessage } from "element-plus"
import { BasicTable } from "@/components/Table/index"
// import Icon from "@/components/Icon/Icon.vue"
import FileAction from "@/components/FileAction/FileAction.vue"
import { createActions } from "./actions"

const columns = [
  {
    label: "文件名",
    prop: "filename",
    scopedSlots: "filename",
    ellipsis: true,
    minWidth: "40%"
  },
  {
    label: "修改时间",
    prop: "updateTime",
    minWidth: "25%",
    scopedSlots: "RowAction"
  },
  {
    label: "大小",
    prop: "fileSize",
    scopedSlots: "fileSize",
    minWidth: "23%"
  }
]

// 表格数据的操作
const actions = createActions({
  // deleteFile: (data) => {
  //   // 批量删除
  //   deleteFileRef.value.show(data)
  // },
  // renameFile: (data) => {
  //   renameFile(data)
  // },
  // moveFile: (data) => {
  //   moveFileRef.value.show(data)
  // },
  // downloadFile: async (data) => {
  //   handleDownloadFile(data)
  // },
  // shareFile: (data) => {
  //   shareFileRef.value.show(data)
  // }
})

const canEditColumns = ref(columns) // 可编辑列
const selectedRow = ref([]) // 选中行
const navigationRef = ref() // 导航
const previewRef = ref() // 预览
const dataTableRef = ref() // 表格
const editNameRef = ref() // 编辑框

const editing = ref(false)

// const deleteFileRef = ref() // 删除文件
// const moveFileRef = ref() // 移动文件
// const shareFileRef = ref() // 分享文件

// const loadDataTable = async () => {
const loadDataTable = () => {
  // const param = { ...res, ...params }
  // if (param.category !== "all") {
  //   delete param.filePid
  // }
  // try {
  // return loadDataList(param)
  // 返回模拟数据
  return [
    {
      id: "1",
      filePid: "0",
      fileSize: 100 * 1024, // 转换为字节
      filename: "文件1",
      fileCover: "",
      folderType: 0,
      fileCategory: 5,
      fileType: 10,
      status: 2,
      updateTime: "2023-10-01 14:49:56",
      recoveryTime: null
    },
    {
      id: "2",
      filePid: "0",
      fileSize: 200 * 1024, // 转换为字节
      filename: "文件2",
      fileCover: "",
      folderType: 0,
      fileCategory: 5,
      fileType: 10,
      status: 2,
      updateTime: "2023-10-01 14:48:57",
      recoveryTime: null
    },
    {
      id: "3",
      filePid: "0",
      fileSize: 300 * 1024, // 转换为字节
      filename: "文件3",
      fileCover: "",
      folderType: 0,
      fileCategory: 5,
      fileType: 10,
      status: 2,
      updateTime: "2023-10-01 12:17:57",
      recoveryTime: null
    }
  ]
  // } catch (error) {
  // return []
  // }
}

const handleTableRowSelect = (rows) => {
  if (rows.length > 0) {
    canEditColumns.value[0].label = `已选中${rows.length}个文件/文件夹`
  } else {
    canEditColumns.value[0].label = "文件名"
  }
  selectedRow.value = rows
}

// 打开文件夹或者预览文件
const openFile = (row) => {
  if (row.folderType == 1) {
    navigationRef.value.openFolder(row)
    return
  }
  // 文件
  if (row.status != 2) {
    ElMessage.warning("文件未完成转码，无法预览")
    return
  }
  previewRef.value.showPreview(row, 0)
}

// 搜索文件
const searchFile = () => {
  // params.filename = filenameFuzzy.value
  // reloadTable()
}

// 新建文件夹
const newFolder = () => {
  //   if (editing.value) {
  //     return
  //   }
  //   dataTableRef.value.setRowFieldValue("showEdit", false)
  //   editing.value = true
  //   dataTableRef.value.unshiftRow({
  //     showEdit: true,
  //     folderType: 1,
  //     id: "",
  //     filePid: currentFolder.value.id,
  //     filenameReal: "",
  //     updateTime: dateFormat("YY-mm-dd HH:MM", new Date())
  //   })
  //   nextTick(() => {
  //     editNameRef.value.focus()
  //   })
  // }
  // // 重命名文件
  // const renameFile = (row) => {
  //   if (isArray(row)) return
  //   if (dataTableRef.value.tableElRef.data[0].id == "") {
  //     dataTableRef.value.tableElRef.data.splice(0, 1)
  //   }
  //   dataTableRef.value.setRowFieldValue("showEdit", false)
  //   row.showEdit = true
  //   editing.value = true
  //   // 编辑文件
  //   if (row.folderType == 0) {
  //     // 拿到文件名
  //     row.filenameReal = row.filename?.substring(0, row.filename.lastIndexOf("."))
  //     row.fileSuffix = row.filename?.substring(row.filename.lastIndexOf("."))
  //   } else {
  //     row.filenameReal = row.filename
  //     row.fileSuffix = ""
  //   }
  //   nextTick(() => {
  //     editNameRef.value = true
  //   })
}

// 取消重命名
const cancelNameEdit = (index) => {
  const fileData = dataTableRef.value.tableElRef.data[index]
  if (fileData.id) {
    fileData.showEdit = false
  } else {
    dataTableRef.value.tableElRef.data.splice(index, 1)
  }
  editing.value = false
}

// 保存文件名
const saveNameEdit = async (index) => {
  // dataTableRef.value.setLoading(true)
  // const { id, filePid, folderType, filenameReal } = dataTableRef.value.tableElRef.data[index]
  // if (filenameReal == "" || filenameReal.indexOf("/") != -1) {
  //   ElMessage.warning(`文件${folderType == 1 ? "夹" : ""}名不能为空，且不能含有 '/'`)
  //   return
  // }
  // if (id == "") {
  //   // 新建文件
  //   try {
  //     // const res = await newFolderApi({ filePid, filename: filenameReal })
  //     // dataTableRef.value.setRowValue(res, 0)
  //   } catch (error) {
  //     cancelNameEdit(0)
  //   }
  // } else {
  //   // 重命名
  //   try {
  //     // const res = await renameApi({ id, filename: filenameReal })
  //     // dataTableRef.value.setRowValue(res, index)
  //   } catch (error) {
  //     cancelNameEdit(0)
  //   }
  // }
  // editing.value = false
  // dataTableRef.value.setLoading(false)
}

// const emit = defineEmits(["addFile"])
// const dataTableRef = ref()
// const editNameRef = ref()
// const editing = ref(false)
// const currentFolder = ref({ id: "0" })
// const params = reactive({
//   page: 1 as number,
//   pageSize: 20 as number,
//   filePid: "0" as string,
//   category: "all" as string
// })

// const fileList = ref([
//   { name: "文件夹1", updateTime: "2023-10-01 14:49:56", size: "-" },
//   { name: "文件2", updateTime: "2023-10-01 14:48:57", size: "-" },
//   { name: "文件3", updateTime: "2023-10-01 12:17:57", size: "-" },
//   { name: "文件夹4", updateTime: "2023-10-01 11:11:11", size: "-" },
//   { name: "文件5", updateTime: "2023-10-01 14:48:06", size: "0 字节" },
//   { name: "文件6", updateTime: "2023-10-01 12:10:48", size: "4.8 MB" }
// ])

// const reloadTable = () => {
//   dataTableRef.value?.reload()
// }

// const newFolder = () => {
//   if (editing.value) return
//   dataTableRef.value.setRowFieldValue("showEdit", false)
//   editing.value = true
//   dataTableRef.value.unshiftRow({
//     showEdit: true,
//     folderType: 1,
//     id: "",
//     filePid: currentFolder.value.id,
//     filenameReal: "",
//     updateTime: new Date().toLocaleString()
//   })
//   nextTick(() => {
//     editNameRef.value.focus()
//   })
// }

// defineExpose({ reloadTable })
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
