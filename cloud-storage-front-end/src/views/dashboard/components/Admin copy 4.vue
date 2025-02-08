<template>
  <div class="app-container">
    <div class="main-ui-header">
      <!-- 上传按钮 -->
      <el-upload
        :show-file-list="false"
        :with-credentials="true"
        :multiple="true"
        :http-request="handleFileUpload"
        :accept="fileAccept"
        class="upload-resource"
      >
        <el-button type="primary" color="#FBBC4D">
          <el-icon><Upload /></el-icon>上传
        </el-button>
      </el-upload>

      <!-- 操作按钮组 -->
      <div class="operation-buttons">
        <el-button type="primary" color="#FBBC4D" @click="handleCreateFolder">
          <el-icon><FolderAdd /></el-icon>新建
        </el-button>

        <el-button type="primary" color="#FBBC4D" @click="handleRefresh">
          <el-icon><Refresh /></el-icon>刷新
        </el-button>

        <el-button type="primary" color="#FBBC4D" @click="handleDownload">
          <el-icon><Download /></el-icon>下载
        </el-button>
      </div>
    </div>

    <!-- 文件列表 -->
    <div class="main-ui-table">
      <vue-good-table
        max-height="100%"
        :columns="columns"
        :rows="tableRows"
        :select-options="{
          enabled: true,
          selectOnCheckboxOnly: true,
          selectionText: '选中的行数',
          clearSelectionText: ''
        }"
        :search-options="{ enabled: true, placeholder: '搜索文件...', skipDiacritics: true }"
        :pagination-options="{
          enabled: true,
          mode: 'pages',
          perPage: 10,
          perPageDropdown: [10, 20, 50],
          dropdownAllowAll: false,
          nextLabel: '下一页',
          prevLabel: '上一页',
          allLabel: '全部',
          rowsPerPageLabel: '每页显示',
          pageLabel: '页',
          ofLabel: '/'
        }"
        v-on:row-click="handleRowClick"
        styleClass="vgt-table"
      >
        <!-- 列 -->
        <template #table-row="props">
          <span v-if="props.column.field === 'filename'">
            <div class="file-item">
              <!-- 文件图标/预览图 -->
              <template v-if="isPreviewable(props.row)">
                <Icon :cover="props.row.fileCover" :width="32" />
              </template>
              <template v-else>
                <Icon :file-type="props.row.fileType" :width="32" />
              </template>

              <!-- 文件名 -->
              <span class="filename" :title="props.row.filename">
                {{ props.row.filename }}
              </span>
            </div>
          </span>
          <span v-else>
            <span class="rows-value" :title="props.row.filename">
              {{ formatColumnValue(props) }}
            </span>
          </span>
        </template>
        <!-- 所选行操作槽-->
        <template #selected-row-actions>
          <el-button-group>
            <el-button type="danger" @click="handleBatchDelete">删除</el-button>
            <el-button @click="handleBatchMove">移动到</el-button>
          </el-button-group>
        </template>
      </vue-good-table>
    </div>

    <!-- 上传进度 -->
    <div class="upload-progress">
      <div class="file-info">
        <span class="filename">当前上传文件:</span>
        <span class="filesize">{{ currentFile ? formatFileSize(currentFile.size) : '' }}</span>
      </div>
      <div class="status-text">
        {{ statusText }}
      </div>
    </div>

    <!-- 上传历史 -->
    <div class="upload-history">
      <h3>上传历史</h3>
      <div v-for="(item, index) in uploadHistory" :key="index" class="history-item">
        <div class="file-info">
          {{ item.filename }}
        </div>
        <div class="file-status">
          {{ item.status === 'success' ? '上传成功' : '上传失败' }}
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed } from "vue"
import type { UploadRequestOptions } from "element-plus"
import { ElMessage } from "element-plus"
import { VueGoodTable } from "vue-good-table-next"
import Icon from "@/components/FileIcon/Icon.vue"
import { FolderOpened } from "@element-plus/icons-vue"
import { uploadFileApi } from "@/api/file" // 导入上传API

// 定义文件接受类型
const fileAccept = ".jpg,.jpeg,.png,.gif,.zip,.doc,.docx,.pdf"

// 上传状态
const uploadStatus = ref('')
const uploadProgress = ref(0)
const currentFile = ref(null)
const statusText = ref('')

// 上传历史记录
const uploadHistory = ref([])

// 格式化文件大小
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

// 处理文件上传
const handleFileUpload = async (options: UploadRequestOptions) => {
  try {
    const { file } = options
    currentFile.value = file
    uploadStatus.value = 'uploading'
    uploadProgress.value = 0
    statusText.value = '准备上传...'

    // 创建 FormData
    const formData = new FormData()
    formData.append('file', file)
    formData.append('metadata', JSON.stringify({
      filename: file.name,
      type: file.type
    }))

    // 调用上传API
    const response = await uploadFileApi(formData, {
      onUploadProgress: (progressEvent) => {
        const progress = (progressEvent.loaded / progressEvent.total) * 100
        uploadProgress.value = Math.round(progress)
        statusText.value = '正在上传...'
      }
    })

    // 处理上传成功
    uploadStatus.value = 'success'
    uploadProgress.value = 100
    statusText.value = '上传成功'

    // 添加到上传历史
    uploadHistory.value.unshift({
      filename: file.name,
      size: file.size,
      status: 'success',
      url: response.data.url,
      key: response.data.key
    })

    // 刷新文件列表
    handleRefresh()

    ElMessage.success('文件上传成功')

    // 重置上传状态
    setTimeout(() => {
      currentFile.value = null
      uploadProgress.value = 0
      uploadStatus.value = ''
      statusText.value = ''
    }, 2000)

  } catch (err: unknown) {
    // 处理上传错误
    uploadStatus.value = 'error'
    statusText.value = '上传失败'

    uploadHistory.value.unshift({
      filename: file.name,
      size: file.size,
      status: 'error'
    })

    console.error("Upload error:", err)
    ElMessage.error('文件上传失败')
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
</script>

<style lang="scss">
.main-ui-header {
  display: flex;
  align-items: center;
  // padding: 16px 20px;
  margin-left: 1%;
  margin-top: 1%;
  gap: 16px;

  .breadcrumb-nav {
    flex: 1;
    margin: 0;
    padding: 0;
  }
}

.main-ui-table {
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

  .vgt-table {
    td {
      vertical-align: middle !important;
    }

    td:first-child {
      .file-item {
        min-height: 40px;
      }
    }

    .rows-value {
      display: inline-block;
      line-height: 40px;
    }
  }

  .vgt-wrap {
    .vgt-global-search {
      border: none;
      margin: -37px 0 10px;
      float: right;

      .vgt-global-search__input {
        .input__icon {
          display: none;
        }

        input {
          width: 300px;
          border: 1px solid #dcdfe6;
          border-radius: 4px;
          padding: 8px 12px;
          font-size: 14px;
          transition: all 0.3s;

          &:focus {
            outline: none;
            border-color: #fbbc4d;
            box-shadow: 0 0 0 2px rgba(251, 188, 77, 0.2);
          }

          &::placeholder {
            color: #909399;
          }
        }
      }
    }
  }
}

.breadcrumb-nav {
  .el-breadcrumb {
    font-size: 14px;

    .el-breadcrumb__item {
      .el-icon {
        margin-right: 4px;
        font-size: 16px;
        vertical-align: -0.15em;
      }

      .el-breadcrumb__inner {
        color: #606266;
        display: inline-flex;
        align-items: center;

        &:hover {
          color: #409eff;
        }

        &.is-link {
          font-weight: normal;
        }
      }

      &:last-child {
        .el-breadcrumb__inner {
          color: #303133;
          cursor: default;

          &:hover {
            color: #303133;
          }
        }
      }
    }
  }
}

// 添加上传进度相关样式
.upload-progress {
  margin-top: 10px;
  padding: 10px;
  border-radius: 4px;
  background-color: #f5f7fa;

  .file-info {
    display: flex;
    align-items: center;
    margin-bottom: 10px;

    .filename {
      margin-left: 8px;
      font-weight: bold;
    }

    .filesize {
      margin-left: 8px;
      color: #909399;
    }
  }

  .status-text {
    margin-top: 5px;
    font-size: 12px;
    color: #909399;
  }
}

.upload-history {
  margin-top: 20px;

  .history-item {
    display: flex;
    align-items: center;
    padding: 8px;
    border-bottom: 1px solid #ebeef5;

    .file-info {
      flex: 1;
    }

    .file-status {
      margin: 0 10px;
    }
  }
}
</style>
