<template>
  <div class="app-container">
    <el-card v-loading="loading">
      <div class="toolbar-wrapper">
        <el-button :disabled="currentFolderId === 0" type="primary" @click="handleBackToParent">
          <el-icon><Back /></el-icon>返回上级
        </el-button>
        <el-button type="danger" @click="handleBatchDelete" :disabled="selectedRows.length === 0">
          <el-icon><Delete /></el-icon>批量删除
        </el-button>
      </div>
      <div class="table-wrapper">
        <el-table
          ref="multipleTable"
          :data="tableData"
          v-model="selectedRows"
          @selection-change="handleSelectionChange"
          @row-click="handleRowClick"
          style="width: 100%"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column prop="filename" label="文件名" min-width="200">
            <template #default="{ row }">
              <div class="file-item">
                <!-- 文件图标/预览图 -->
                <template v-if="row.fileType === 3 && row.fileCover">
                  <!-- 仅图片类型(3)显示封面 -->
                  <Icon :cover="row.fileCover" :width="32" />
                </template>
                <template v-else>
                  <!-- 其他类型显示对应图标 -->
                  <Icon :file-type="row.fileType" :width="32" />
                </template>
                <span class="filename" :title="row.filename">{{ row.filename }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="fileSize" label="大小" width="120">
            <template #default="{ row }">
              {{ row.fileSize }}
            </template>
          </el-table-column>
          <el-table-column prop="updateTime" label="删除时间" width="180">
            <template #default="{ row }">
              {{ row.updateTime }}
            </template>
          </el-table-column>
          <el-table-column fixed="right" label="操作" width="200">
            <template #default="{ row }">
              <el-button type="primary" @click.stop="handleRestore(row)">
                <el-icon><RefreshRight /></el-icon>还原
              </el-button>
              <el-button type="danger" @click.stop="handleDelete(row)">
                <el-icon><Delete /></el-icon>删除
              </el-button>
            </template>
          </el-table-column>
          <template #empty>
            <div style="text-align: center; padding: 40px 0">
              <svg-icon name="trash" color="green" width="100px" height="80px" />
              <div style="margin-top: -30px; font-size: 16px; font-weight: bold">空空如也 ~</div>
            </div>
          </template>
        </el-table>
      </div>
      <div class="pagination">
        <el-pagination
          v-model:current-page="paginationData.currentPage"
          v-model:page-size="paginationData.pageSize"
          :total="paginationData.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import { ref, watch, reactive, onUnmounted } from "vue"
import { ElMessage, ElMessageBox } from "element-plus"
import { usePagination } from "@/hooks/usePagination"
import Icon from "@/components/FileIcon/Icon.vue"
import { recycleApi } from "@/api/recycle/recycle"
import { formatFileSize } from "@/utils/format/formatFileSize"
import { timestampToDate } from "@/utils/format/formatTime"
import { userFileApi } from "@/api/file/repository"
import { Back } from "@element-plus/icons-vue"

interface FileListItem {
  id: number
  filename: string
  updateTime: string
  fileSize?: number
  fileType?: number
  isFolder: boolean
  fileCover?: string
  ext: string // 扩展名
  repository_id: number // 文件详情id
}
// 判断文件是否可预览
const isPreviewAble = (file: FileListItem): boolean => {
  if (file.isFolder) return false
  return file.fileType === 3 || file.fileType === 4 // 图片(3)和视频(4)可预览
}

// 根据文件扩展名判断文件类型
const getFileTypeNumber = (ext: string): number => {
  if (!ext) return 0 // 文件夹返回0
  const extMap = {
    // 图片
    jpg: 3,
    jpeg: 3,
    png: 3,
    gif: 3,
    // 视频
    mp4: 4,
    avi: 4,
    mov: 4,
    // 音频
    mp3: 5,
    wav: 5,
    // Word
    doc: 6,
    docx: 6,
    // Excel
    xls: 7,
    xlsx: 7,
    // PPT
    ppt: 8,
    pptx: 8,
    // 压缩包
    zip: 9,
    rar: 9,
    "7z": 9,
    // PDF
    pdf: 10,
    // 文本
    txt: 2,
    md: 2
  }
  const extension = ext.toLowerCase().replace(".", "")
  return extMap[extension] || 1 // 默认为普通文件
}

// 获取文件夹大小
const getFolderSize = async (folderId: number): Promise<number> => {
  try {
    // console.log("获取文件夹id:", folderId)
    const response = await userFileApi.getFolderSize({
      id: folderId
    })
    // console.log("文件夹大小:", response.data)
    return response.data.size || 0
  } catch (error) {
    console.error("获取文件夹大小失败:", error)
    return 0
  }
}

const loading = ref<boolean>(false)
const { paginationData, handleCurrentChange, handleSizeChange } = usePagination()

const tableData = ref<FileListItem[]>([])
const selectedRows = ref<FileListItem[]>([])

const currentFolderId = ref<number>(0)

// 在 script 部分添加一个数组来存储文件夹路径
const folderStack = ref<number[]>([])

const fetchRecycleList = async () => {
  loading.value = true
  try {
    const res = await recycleApi.getRecycleList({
      id: currentFolderId.value,
      page: paginationData.currentPage,
      size: paginationData.pageSize
    })
    // console.log(res.data.list[0].update_time)
    tableData.value = await Promise.all(
      res.data.list.map(async (item) => {
        // 文件夹的 RepositoryId 为 0
        if (!item.repository_id) {
          // 文件夹
          const folderSize = await getFolderSize(item.id)
          return {
            id: item.id,
            filename: item.name,
            fileType: 0, // 文件夹类型
            fileSize: formatFileSize(folderSize), // 文件夹大小
            updateTime: timestampToDate(item.update_time) || new Date().toISOString(),
            fileCover: "", // 文件夹没有封面
            isFolder: true,
            ext: "",
            repository_id: 0
          }
        } else {
          // 文件
          return {
            id: item.id,
            filename: item.name,
            fileType: getFileTypeNumber(item.ext),
            fileSize: formatFileSize(item.size),
            updateTime: timestampToDate(item.update_time) || new Date().toISOString(),
            fileCover: item.path,
            isFolder: false,
            ext: item.ext,
            repository_id: item.repository_id
          }
        }
      })
    )
    paginationData.total = res.data.total
  } catch (error) {
    console.error("获取回收站列表失败:", error)
    ElMessage.error("获取回收站列表失败")
  } finally {
    loading.value = false
  }
}

const handleDelete = async (row: FileListItem) => {
  try {
    await ElMessageBox.confirm(`确认删除该项：${row.filename}？`, "提示", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning"
    })

    const res = await recycleApi.deleteRecycle({ id: row.id })
    if (res.data.success) {
      ElMessage.success("删除成功")
      fetchRecycleList()
    }
  } catch (error) {
    if (error !== "cancel") {
      ElMessage.error("删除失败")
    }
  }
}

const handleBatchDelete = async () => {
  if (selectedRows.value.length === 0) {
    ElMessage.warning("请先选择要删除的项")
    return
  }

  try {
    await ElMessageBox.confirm(`确认删除选中的 ${selectedRows.value.length} 项？`, "提示", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning"
    })

    const deletePromises = selectedRows.value.map((row) => recycleApi.deleteRecycle({ id: row.id }))
    await Promise.all(deletePromises)

    ElMessage.success("批量删除成功")
    fetchRecycleList()
    selectedRows.value = []
  } catch (error) {
    if (error !== "cancel") {
      ElMessage.error("批量删除失败")
    }
  }
}

const handleSelectionChange = (val) => {
  selectedRows.value = val
}

const multipleTable = ref()
const handleRowClick = (row: FileListItem) => {
  if (row.isFolder) {
    folderStack.value.push(currentFolderId.value) // 将当前文件夹 ID 推入栈中
    currentFolderId.value = row.id
    paginationData.currentPage = 1
    fetchRecycleList()
  } else {
    selectedRows.value = row
    multipleTable.value.toggleRowSelection(row)
  }
}

watch(
  [() => paginationData.currentPage, () => paginationData.pageSize],
  () => {
    fetchRecycleList()
  },
  { immediate: true }
)

// 添加还原文件功能
const handleRestore = async (row: FileListItem) => {
  try {
    await ElMessageBox.confirm(`确认还原该项：${row.filename}？`, "提示", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning"
    })

    const res = await recycleApi.restoreRecycle({ id: row.id })
    if (res.data.success) {
      ElMessage.success("还原成功")

      fetchRecycleList()
    }
  } catch (error) {
    if (error !== "cancel") {
      ElMessage.error("还原失败")
    }
  }
}

// 修改 handleBackToParent 函数
const handleBackToParent = async () => {
  if (folderStack.value.length === 0) return // 如果栈为空，返回
  currentFolderId.value = folderStack.value.pop() || 0 // 从栈中弹出上一个文件夹 ID
  paginationData.currentPage = 1
  fetchRecycleList()
}

// 在组件卸载时清理事件监听
onUnmounted(() => {
  fetchRecycleList()
})
</script>

<style lang="scss" scoped>
.app-container {
  padding: 20px;

  .file-item {
    display: flex;
    align-items: center;
    gap: 8px;
    cursor: pointer;

    :deep(.icon) {
      width: 32px !important;
      height: 32px !important;
      flex-shrink: 0;
    }
  }

  .filename {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .pagination {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }

  .context-menu {
    position: fixed;
    z-index: 2000;
    background: white;
    border-radius: 4px;
    // box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
    padding: 8px 0;
    min-width: 120px;
    user-select: none;

    ul {
      list-style: none;
      margin: 0;
      padding: 0;
    }

    li {
      padding: 8px 16px;
      cursor: pointer;
      display: flex;
      align-items: center;
      gap: 8px;
      color: #606266;

      &:hover {
        background-color: #f5f7fa;
      }

      .el-icon {
        font-size: 16px;
      }
    }
  }

  .context-menu-mask {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    z-index: 1999;
  }
}
</style>
