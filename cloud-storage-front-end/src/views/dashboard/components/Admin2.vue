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
        <span class="filesize">{{ currentFile?.size ? formatFileSize(currentFile.size) : "" }}</span>
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
          {{ item.status === "success" ? "上传成功" : "上传失败" }}
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted } from "vue"
import type { UploadRequestOptions, UploadRawFile } from "element-plus"
import { ElMessage, ElMessageBox } from "element-plus"
import { VueGoodTable } from "vue-good-table-next"
import Icon from "@/components/FileIcon/Icon.vue"

import { uploadFileApi } from "@/api/file" // 导入上传API
import {
  type FileUploadRequestData,
  type ChunkUploadRequestData,
  type FileUploadResponseData,
  type ChunkUploadCompleteResponseData
} from "@/api/file/types/upload"

import { isNotEmpty } from "@/utils/isEmpty"
import { userFileApi } from "@/api/file/repository"
import type * as Repository from "@/api/file/types/repository"
import { on } from "events"
import { c } from "vite/dist/node/types.d-aGj9QkWt"

// 定义文件接受类型
const fileAccept = ".jpg,.jpeg,.png,.gif,.zip,.doc,.docx,.pdf"

// 定义上传历史记录项的接口
interface UploadHistoryItem {
  filename: string
  size: number
  status: "success" | "error"
  url?: string
  key?: string
}

// 定义当前文件的类型
const currentFile = ref<UploadRawFile | null>(null)

// 定义上传历史数组的类型
const uploadHistory = ref<UploadHistoryItem[]>([])

// 上传状态
const uploadStatus = ref("")
const uploadProgress = ref(0)
const statusText = ref("")

// 定义分片大小为 5MB
const CHUNK_SIZE = 5 * 1024 * 1024

// 添加分片上传相关的状态
const chunks = ref<Blob[]>([])
const currentChunkIndex = ref(0)
const uploadId = ref("")
const uploadedETags = ref<string[]>([])

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

interface FileListItem {
  id: number
  filename: string
  updateTime: string
  fileSize?: number
  fileType?: number
  status?: number
}

// 测试文件列表数据
const fileList = ref<FileListItem[]>([
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

// // 获取文件列表
// onMounted(async () => {
//   const { data } = await userFileApi.getFileList({
//     id: 0,
//     page: 1,
//     size: 10
//   })

//   const res = data.list
//   for (const item of res) {
//     const fileListItem: FileListItem = {
//       id: item.id,
//       filename: item.name,
//       updateTime: String(item.update_time),
//       fileSize: item.size,
//       fileType: getFileType[item.ext]
//     }
//     fileList.value.push(fileListItem)
//   }
// })

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
  const { file } = options
  try {
    if (!file) {
      throw new Error("No file selected")
    }

    currentFile.value = file as UploadRawFile
    uploadStatus.value = "uploading"
    uploadProgress.value = 0
    statusText.value = "准备上传..."

    // 根据文件大小选择上传方式
    if (file.size > 20 * 1024 * 1024) {
      await handleChunkUpload(file)
    } else {
      await handleNormalUpload(file)
    }
  } catch (err) {
    handleUploadError(err, file)
  }
}

// 处理普通上传
const handleNormalUpload = async (file: File) => {
  const fileUploadRequestData: FileUploadRequestData = {
    file: file,
    metadata: JSON.stringify({
      filename: file.name,
      type: file.type
    })
  }

  const response = await uploadFileApi.upload(fileUploadRequestData, {
    onUploadProgress: (progressEvent) => {
      if (progressEvent.total) {
        const progress = (progressEvent.loaded / progressEvent.total) * 100
        uploadProgress.value = Math.round(progress)
        statusText.value = "正在上传..."
      }
    }
  })

  handleUploadSuccess(response.data, file)
}

// 处理分片上传
const handleChunkUpload = async (file: File) => {
  // 初始化分片上传
  const initResponse = await uploadFileApi.initiateMultipart({
    file_name: file.name,
    file_size: file.size,
    metadata: JSON.stringify({ type: file.type })
  })

  const { upload_id, key } = initResponse.data
  uploadId.value = upload_id

  // 将文件分片
  const chunkCount = Math.ceil(file.size / CHUNK_SIZE)
  chunks.value = []
  uploadedETags.value = new Array(chunkCount)

  for (let i = 0; i < chunkCount; i++) {
    const start = i * CHUNK_SIZE
    const end = Math.min(start + CHUNK_SIZE, file.size)
    chunks.value.push(file.slice(start, end))
  }

  // 上传所有分片
  for (let i = 0; i < chunks.value.length; i++) {
    currentChunkIndex.value = i
    const chunkUploadRequestData: ChunkUploadRequestData = {
      upload_id: uploadId.value,
      chunk_index: i + 1,
      key: key,
      file: chunks.value[i]
    }

    const chunkResponse = await uploadFileApi.uploadPart(chunkUploadRequestData)
    if (isNotEmpty(chunkResponse.data)) {
      uploadedETags.value[i] = chunkResponse.data.etag
    }

    // 更新进度
    const progress = ((i + 1) / chunks.value.length) * 100
    uploadProgress.value = Math.round(progress)
    statusText.value = `正在上传第 ${i + 1}/${chunks.value.length} 个分片`
  }

  // 完成分片上传
  const completeResponse = await uploadFileApi.completeMultipart({
    upload_id: uploadId.value,
    key: key,
    etags: uploadedETags.value
  })

  handleUploadSuccess(completeResponse.data, file)
}

// 修改处理上传成功的函数
const handleUploadSuccess = async (data: any, file: File) => {
  uploadStatus.value = "success"
  uploadProgress.value = 100
  statusText.value = "上传成功"
  // console.log(parentId, data.repository_id, file.name)
  // 保存文件关联信息
  const parentId = currentPath.value[currentPath.value.length - 1] || 0
  await userFileApi.saveRepository({
    parent_id: Number(parentId),
    repository_id: Number(data.repository_id),
    name: ""
    // name: folderName
  })

  const historyItem: UploadHistoryItem = {
    filename: file.name,
    size: file.size,
    status: "success",
    url: data.url,
    key: data.key
  }
  uploadHistory.value.unshift(historyItem)

  handleRefresh()
  ElMessage.success("文件上传成功")

  resetUploadState()
}

// 处理上传错误
const handleUploadError = (error: any, file: File) => {
  uploadStatus.value = "error"
  statusText.value = "上传失败"

  const historyItem: UploadHistoryItem = {
    filename: file.name,
    size: file.size,
    status: "error"
  }
  uploadHistory.value.unshift(historyItem)

  console.error("Upload error:", error)
  ElMessage.error("文件上传失败")

  resetUploadState()
}

// 重置上传状态
const resetUploadState = () => {
  setTimeout(() => {
    currentFile.value = null
    uploadProgress.value = 0
    uploadStatus.value = ""
    statusText.value = ""
    chunks.value = []
    currentChunkIndex.value = 0
    uploadId.value = ""
    uploadedETags.value = []
  }, 2000)
}

const currentPath = ref<string[]>([])
const handleCreateFolder = async () => {
  try {
    const { value: folderName } = await ElMessageBox.prompt("请输入文件夹名称", "新建文件夹", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      inputValidator: (value) => {
        if (!value) {
          return "文件夹名称不能为空"
        }
        return true
      }
    })
    if (folderName) {
      const parentId = currentPath.value[currentPath.value.length - 1] || 0
      await userFileApi.createFolder({
        parent_id: Number(parentId),
        name: folderName
      })
      ElMessage.success("文件夹创建成功")
      handleRefresh()
    }
  } catch {
    // 用户取消操作
  }
}

// 根据文件扩展名判断文件类型
const getFileType = (ext: string): number => {
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

// 刷新文件列表
const handleRefresh = async () => {
  try {
    const parentId = currentPath.value[currentPath.value.length - 1] || 0
    const { data } = await userFileApi.getFileList({
      id: Number(parentId),
      page: 1,
      size: 10
    })
    // console.log(data)
    // 转换后端数据格式为前端所需格式
    fileList.value = data.list.map((file) => ({
      id: file.id,
      filename: file.name,
      fileType: getFileType(file.ext), // 需要添加文件类型判断函数
      fileSize: file.size,
      updateTime: new Date().toISOString(), // 后端未提供,暂时用当前时间
      fileCover: file.path // 如果是图片类型可以用作预览
    }))
    // console.log(fileList.value)
    ElMessage.success("刷新成功")
  } catch (error) {
    console.log(error)
    ElMessage.error("刷新失败")
  }
}

const handleDownload = async () => {
  const selectedRows = tableRows.value.filter((row) => row.selected)
  if (selectedRows.length === 0) {
    ElMessage.warning("请选择要下载的文件")
    return
  }

  try {
    // TODO: 调用下载API
    ElMessage.success("开始下载")
  } catch (error) {
    ElMessage.error("下载失败")
  }
}

const handleRowClick = async (params: any) => {
  const { row } = params
  if (row.fileType === 0) {
    // 文件夹
    currentPath.value.push(row.id)
    await handleRefresh()
  } else {
    // TODO: 处理文件预览或下载
  }
}

const handleBatchDelete = async () => {
  const selectedRows = tableRows.value.filter((row) => row.selected)
  if (selectedRows.length === 0) {
    ElMessage.warning("请选择要删除的文件")
    return
  }

  try {
    await ElMessageBox.confirm("确定要删除选中的文件吗?", "提示", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning"
    })

    // 逐个删除选中的文件
    for (const row of selectedRows) {
      await userFileApi.deleteFile({
        id: row.id
      })
    }

    ElMessage.success("删除成功")
    handleRefresh()
  } catch {
    // 用户取消操作
  }
}

const handleBatchMove = async () => {
  const selectedRows = tableRows.value.filter((row) => row.selected)
  if (selectedRows.length === 0) {
    ElMessage.warning("请选择要移动的文件")
    return
  }

  try {
    const { data } = await userFileApi.getFolderList({
      id: currentPath.value[currentPath.value.length - 1] || 0
    })

    // 显示文件夹选择对话框
    const { value: targetFolderId } = await ElMessageBox.prompt("请选择目标文件夹", "移动到", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      inputType: "select",
      inputValue: "",
      inputPlaceholder: "请选择",
      inputValidator: (value) => {
        if (!value) return "请选择目标文件夹"
        return true
      },
      inputOptions: data.list.map((folder) => ({
        label: folder.name,
        value: folder.id
      }))
    })

    if (targetFolderId) {
      // 逐个移动选中的文件
      for (const row of selectedRows) {
        await userFileApi.moveFile({
          id: row.id,
          parentId: targetFolderId
        })
      }
      ElMessage.success("移动成功")
      handleRefresh()
    }
  } catch (error) {
    ElMessage.error("移动失败")
  }
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
