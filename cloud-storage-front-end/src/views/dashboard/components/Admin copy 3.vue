<template>
  <div class="app-container">
    <div class="dashboard-header">
      <!-- 上传按钮 -->
      <el-upload
        :show-file-list="false"
        :with-credentials="true"
        :multiple="true"
        :http-request="handleFileUpload"
        :accept="fileAccept"
        class="upload-btn"
      >
        <el-button type="primary">
          <el-icon><Upload /></el-icon>上传文件
        </el-button>
      </el-upload>

      <!-- 操作按钮组 -->
      <div class="action-group">
        <el-button @click="handleCreateFolder">
          <el-icon><FolderAdd /></el-icon>新建文件夹
        </el-button>

        <el-button @click="handleRefresh">
          <el-icon><Refresh /></el-icon>刷新
        </el-button>
      </div>

      <!-- 搜索框 -->
      <div class="search-box">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索文件..."
          clearable
          @clear="handleSearch"
          @keyup.enter="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </div>
    </div>

    <!-- 面包屑导航 -->
    <div class="breadcrumb">
      <el-breadcrumb separator="/">
        <el-breadcrumb-item :to="{ path: '/' }">全部文件</el-breadcrumb-item>
        <el-breadcrumb-item v-for="(path, index) in currentPath" :key="index">
          {{ path }}
        </el-breadcrumb-item>
      </el-breadcrumb>
    </div>

    <!-- 文件列表 -->
    <div class="file-table">
      <el-table
        v-loading="loading"
        :data="filteredFileList"
        style="width: 100%"
        @selection-change="handleSelectionChange"
        @row-click="handleRowClick"
      >
        <el-table-column type="selection" width="55" />

        <el-table-column label="文件名" min-width="300">
          <template #default="{ row }">
            <div class="file-item">
              <el-icon :size="20" class="file-icon">
                <component :is="getFileIcon(row)" />
              </el-icon>
              <span class="file-name" :class="{ 'is-folder': row.isFolder }">
                {{ row.filename }}
              </span>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="updateTime" label="修改时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.updateTime) }}
          </template>
        </el-table-column>

        <el-table-column prop="fileSize" label="大小" width="120">
          <template #default="{ row }">
            {{ formatFileSize(row.fileSize) }}
          </template>
        </el-table-column>

        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button-group>
              <el-button v-if="!row.isFolder" link type="primary" @click.stop="handleDownload(row)"> 下载 </el-button>
              <el-button link type="primary" @click.stop="handleShare(row)"> 分享 </el-button>
              <el-dropdown trigger="click" @click.stop>
                <el-button link type="primary">
                  更多<el-icon><ArrowDown /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item @click="handleRename(row)"> 重命名 </el-dropdown-item>
                    <el-dropdown-item @click="handleMove(row)"> 移动 </el-dropdown-item>
                    <el-dropdown-item divided @click="handleDelete(row)">
                      <span style="color: var(--el-color-danger)">删除</span>
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>

    <!-- 批量操作工具栏 -->
    <div class="batch-toolbar" v-show="selectedFiles.length > 0">
      <span class="selected-count">已选择 {{ selectedFiles.length }} 个文件</span>
      <el-button-group>
        <el-button type="primary" @click="handleBatchDownload"> 批量下载 </el-button>
        <el-button type="primary" @click="handleBatchMove"> 批量移动 </el-button>
        <el-button type="danger" @click="handleBatchDelete"> 批量删除 </el-button>
      </el-button-group>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed } from "vue"
import type { UploadRequestOptions } from "element-plus"
import { ElMessage } from "element-plus"
import { VueGoodTable } from "vue-good-table-next"
import Icon from "@/components/FileIcon/Icon.vue"
import { FolderOpened, Document, ZipFile } from "@element-plus/icons-vue"
import { Search } from "@element-plus/icons-vue"

// 定义文件接受类型
const fileAccept = ".jpg,.jpeg,.png,.gif,.zip,.doc,.docx,.pdf"

// 格式化文件大小的函数
const formatFileSize = (size: number): string => {
  if (!size || isNaN(size)) return "0 B"

  const units = ["B", "KB", "MB", "GB", "TB"]
  let index = 0
  let fileSize = size

  while (fileSize >= 1024 && index < units.length - 1) {
    fileSize /= 1024
    index++
  }

  return `${fileSize.toFixed(2)} ${units[index]}`
}

// 格式化日期的函数
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

// 测试文件列表数据
const fileList = ref([
  {
    id: 1,
    filename: "文件1.jpg",
    updateTime: "2023-10-01 14:49:56",
    fileSize: 102400,
    fileType: 3,
    status: 2
  },
  {
    id: 2,
    filename: "文件2",
    updateTime: "2023-10-01 14:48:57",
    fileSize: 209715200,
    fileType: 0,
    status: 2
  },
  {
    id: 3,
    filename: "文件.zip",
    updateTime: "2023-10-01 12:17:57",
    fileSize: 10485760,
    fileType: 9
  },
  {
    id: 1,
    filename: "文件1.jpg",
    updateTime: "2023-10-01 14:49:56",
    fileSize: 102400,
    fileType: 3,
    status: 2
  },
  {
    id: 2,
    filename: "文件2",
    updateTime: "2023-10-01 14:48:57",
    fileSize: 209715200,
    fileType: 0,
    status: 2
  },
  {
    id: 3,
    filename: "文件.zip",
    updateTime: "2023-10-01 12:17:57",
    fileSize: 10485760,
    fileType: 9
  },
  {
    id: 1,
    filename: "文件1.jpg",
    updateTime: "2023-10-01 14:49:56",
    fileSize: 102400,
    fileType: 3,
    status: 2
  },
  {
    id: 2,
    filename: "文件2",
    updateTime: "2023-10-01 14:48:57",
    fileSize: 209715200,
    fileType: 0,
    status: 2
  },
  {
    id: 3,
    filename: "文件.zip",
    updateTime: "2023-10-01 12:17:57",
    fileSize: 10485760,
    fileType: 9
  },
  {
    id: 1,
    filename: "文件1.jpg",
    updateTime: "2023-10-01 14:49:56",
    fileSize: 102400,
    fileType: 3,
    status: 2
  },
  {
    id: 2,
    filename: "文件2",
    updateTime: "2023-10-01 14:48:57",
    fileSize: 209715200,
    fileType: 0,
    status: 2
  },
  {
    id: 3,
    filename: "文件.zip",
    updateTime: "2023-10-01 12:17:57",
    fileSize: 10485760,
    fileType: 9
  },
  {
    id: 1,
    filename: "文件1.jpg",
    updateTime: "2023-10-01 14:49:56",
    fileSize: 102400,
    fileType: 3,
    status: 2
  },
  {
    id: 2,
    filename: "文件2",
    updateTime: "2023-10-01 14:48:57",
    fileSize: 209715200,
    fileType: 0,
    status: 2
  },
  {
    id: 3,
    filename: "文件.zip",
    updateTime: "2023-10-01 12:17:57",
    fileSize: 10485760,
    fileType: 9
  }
])

// 使用 computed 来缓存格式化后的数据
const formattedFileList = computed(() => {
  return fileList.value.map((file) => ({
    ...file,
    formattedSize: formatFileSize(file.fileSize),
    formattedDate: formatDate(file.updateTime)
  }))
})

// 修改表格列定义
const columns = ref([
  {
    label: "文件名",
    field: "filename",
    sortable: true,
    width: "50%"
  },
  {
    label: "修改时间",
    field: "formattedDate", // 使用预先格式化的字段
    sortable: true
  },
  {
    label: "大小",
    field: "formattedSize", // 使用预先格式化的字段
    sortable: true
  }
])

// 修改数据源
const tableRows = computed(() => formattedFileList.value)

// 判断文件是否可预览
const isPreviewable = (file: { fileType: number; status?: number }): boolean => {
  return file.fileType === 3 || (file.fileType === 1 && file.status === 2)
}

// 格式化表格列值
const formatColumnValue = (props) => {
  const { column, formattedRow } = props
  return column.formatFn ? column.formatFn(formattedRow[column.field]) : formattedRow[column.field]
}

// 事件处理方法
const handleFileUpload = async (options: UploadRequestOptions) => {
  try {
    const { file } = options
    // 处理文件上传
    console.log("Uploading file:", file.name)
    ElMessage.success("文件上传成功")
  } catch (err: unknown) {
    console.error("Upload error:", err)
    ElMessage.error("文件上传失败")
  }
}

const currentPath = ref<string[]>([])
const handleCreateFolder = () => {
  // 处理新建文件夹
}

const handleRefresh = () => {
  // 刷新文件列表
}

const handleDownload = () => {
  // 处理下载
}

const handleRowClick = (params) => {
  console.log("Row clicked:", params)
}

const handleSelectionChange = (params) => {
  console.log("Selection changed:", params)
}

const handleBatchDownload = () => {
  // 处理批量下载
}

const handleBatchDelete = () => {
  // 处理批量删除
}

const handleBatchMove = () => {
  // 处理批量移动
}

// 新增的状态和方法
const loading = ref(false)
const searchKeyword = ref("")
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const selectedFiles = ref<any[]>([])

const filteredFileList = computed(() => {
  if (!searchKeyword.value) return fileList.value
  const keyword = searchKeyword.value.toLowerCase()
  return fileList.value.filter((file) => file.filename.toLowerCase().includes(keyword))
})

const getFileIcon = (file: any) => {
  if (file.isFolder) return FolderOpened
  switch (file.fileType) {
    case 3:
      return Document
    case 9:
      return ZipFile
    default:
      return Document
  }
}

const handleSearch = () => {
  currentPage.value = 1
  // 实现搜索逻辑
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  // 重新加载数据
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
  // 重新加载数据
}

const handleShare = (row: any) => {
  // 实现分享逻辑
}

const handleRename = (row: any) => {
  // 实现重命名逻辑
}

const handleMove = (row: any) => {
  // 实现移动逻辑
}

const handleDelete = (row: any) => {
  // 实现删除逻辑
}
</script>

<style lang="scss" scoped>
.dashboard-header {
  display: flex;
  align-items: center;
  padding: 16px 20px;
  background: var(--el-bg-color);
  border-radius: 8px;
  margin-bottom: 16px;
  box-shadow: var(--el-box-shadow-light);

  .action-group {
    display: flex;
    gap: 8px;
    margin: 0 16px;
  }

  .search-box {
    margin-left: auto;
    width: 300px;
  }
}

.breadcrumb {
  margin: 16px 0;
  padding: 0 20px;
}

.file-table {
  background: var(--el-bg-color);
  border-radius: 8px;
  padding: 20px;
  box-shadow: var(--el-box-shadow-light);

  .file-item {
    display: flex;
    align-items: center;
    gap: 8px;

    .file-icon {
      color: var(--el-text-color-secondary);
    }

    .file-name {
      &.is-folder {
        color: var(--el-color-primary);
        cursor: pointer;

        &:hover {
          text-decoration: underline;
        }
      }
    }
  }

  .pagination {
    margin-top: 16px;
    display: flex;
    justify-content: flex-end;
  }
}

.batch-toolbar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 12px 24px;
  background: var(--el-bg-color);
  box-shadow: var(--el-box-shadow);
  display: flex;
  align-items: center;
  justify-content: space-between;
  z-index: 1000;
  transform: translateY(0);
  transition: transform 0.3s ease;

  &.hidden {
    transform: translateY(100%);
  }

  .selected-count {
    color: var(--el-text-color-secondary);
  }
}

// 动画过渡
.el-table {
  .el-table__row {
    transition: background-color 0.3s ease;
  }
}

.el-button {
  transition: all 0.3s ease;
}

// 响应式布局
@media screen and (max-width: 768px) {
  .dashboard-header {
    flex-direction: column;
    gap: 12px;

    .search-box {
      width: 100%;
      margin: 0;
    }
  }

  .file-table {
    .el-table {
      .cell {
        padding: 8px;
      }
    }
  }
}
</style>
