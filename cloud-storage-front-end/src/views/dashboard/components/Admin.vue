<template>
  <div class="app-container">
    <!-- 顶部工具栏 -->
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

      <!-- 添加搜索框 -->
      <el-input v-model="searchKeyword" placeholder="搜索文件" class="search-input" @input="handleSearch">
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>

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

        <el-button type="danger" @click="handleDelete">
          <el-icon><Delete /></el-icon>删除
        </el-button>
      </div>
    </div>

    <!-- 面包屑导航 -->
    <div class="breadcrumb">
      <el-breadcrumb separator="/">
        <el-breadcrumb-item @click="handleRootClick">
          <el-icon><FolderOpened /></el-icon>全部文件
        </el-breadcrumb-item>
        <el-breadcrumb-item v-for="(path, index) in pathHistory" :key="index" @click="handlePathClick(index)">
          {{ path.name }}
        </el-breadcrumb-item>
      </el-breadcrumb>
    </div>

    <!-- 文件列表 -->
    <div class="main-ui-table">
      <el-table
        v-loading="loading"
        :data="fileList"
        style="width: 100%"
        @selection-change="handleSelectionChange"
        @row-click="handleRowClick"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column label="文件名" min-width="200">
          <template #default="{ row }">
            <div class="file-item" @contextmenu.prevent="handleContextMenu(row, $event)">
              <!-- 文件图标/预览图 -->
              <template v-if="isPreviewAble(row)">
                <Icon :cover="row.fileCover" :width="32" />
              </template>
              <template v-else>
                <Icon :file-type="row.fileType" :width="32" />
              </template>
              <!-- 文件名 -->
              <span class="filename" :title="row.filename">{{ row.filename }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="size" label="大小" width="120">
          <template #default="{ row }">
            {{ formatFileSize(row.fileSize) }}
          </template>
        </el-table-column>
        <el-table-column prop="updateTime" label="修改时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.updateTime) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120">
          <template #default="{ row }">
            <el-button @click="previewFile(row)">预览</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 上传进度 -->
    <div v-if="currentFile" class="upload-progress">
      <div class="file-info">
        <span class="filename">当前上传文件: {{ currentFile.name }}</span>
        <span class="filesize">{{ formatFileSize(currentFile.size) }}</span>
      </div>
      <div class="status-text">
        {{ statusText }}
      </div>
    </div>

    <!-- 上传历史 -->
    <div class="upload-history">
      <div v-for="(item, index) in uploadHistory" :key="index" class="history-item">
        <div class="file-info">
          {{ item.filename }}
        </div>
        <div class="file-status">
          {{ item.status === "success" ? "上传成功" : "上传失败" }}
        </div>
      </div>
    </div>

    <!-- 分页 -->
    <div class="pagination">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    <!-- 新建文件夹对话框 -->
    <el-dialog v-model="folderDialog.visible" :title="folderDialog.title" width="30%">
      <el-form :model="folderDialog.form" label-width="80px">
        <el-form-item label="文件夹名">
          <el-input v-model="folderDialog.form.name" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="folderDialog.visible = false">取消</el-button>
        <el-button type="primary" @click="confirmCreateFolder">确定</el-button>
      </template>
    </el-dialog>

    <!-- 重命名对话框 -->
    <el-dialog v-model="renameDialog.visible" title="重命名" width="30%">
      <el-form :model="renameDialog.form" label-width="80px">
        <el-form-item label="新名称">
          <el-input v-model="renameDialog.form.name" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="renameDialog.visible = false">取消</el-button>
        <el-button type="primary" @click="handleRename(renameDialog.fileId, renameDialog.form.name)"> 确定 </el-button>
      </template>
    </el-dialog>

    <!-- 移动对话框 -->
    <el-dialog v-model="moveDialog.visible" title="移动文件" width="30%">
      <el-form :model="moveDialog.form" label-width="80px">
        <el-form-item label="目标文件夹">
          <el-select v-model="moveDialog.form.targetFolderId" placeholder="选择目标文件夹">
            <el-option v-for="folder in folderList" :key="folder.id" :label="folder.name" :value="folder.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="moveDialog.visible = false">取消</el-button>
        <el-button type="primary" @click="handleMove()">确定</el-button>
      </template>
    </el-dialog>

    <!-- 右键菜单 -->
    <div
      v-show="contextMenu.visible"
      class="context-menu"
      :style="{
        left: contextMenu.x + 'px',
        top: contextMenu.y + 'px'
      }"
    >
      <ul>
        <li @click="openRenameDialog">
          <el-icon><Edit /></el-icon>
          重命名
        </li>
        <li @click="openMoveDialog">
          <el-icon><FolderAdd /></el-icon>
          移动到
        </li>
        <li @click="handleDelete">
          <el-icon><Delete /></el-icon>
          删除
        </li>
      </ul>
    </div>

    <!-- 添加点击其他区域关闭右键菜单 -->
    <div v-show="contextMenu.visible" class="context-menu-mask" @click="closeContextMenu" @contextmenu.prevent />

    <!-- 文件预览对话框 -->
    <el-dialog v-model:visible="previewDialog.visible" title="文件预览" width="80%">
      <component :is="previewComponent" :url="previewUrl" :resource="previewResource" />
      <template v-slot:footer>
        <span>
          <el-button @click="previewDialog.visible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted, computed, onUnmounted } from "vue"
import { ElMessage } from "element-plus"
import type { UploadRequestOptions, UploadRawFile } from "element-plus"
import Icon from "@/components/FileIcon/Icon.vue"
import { debounce } from "lodash-es"
import Audio from "@/views/dashboard/components/file-preview/Audio.vue"
import Image from "@/views/dashboard/components/file-preview/Image.vue"
import Default from "@/views/dashboard/components/file-preview/Default.vue"
import Office from "@/views/dashboard/components/file-preview/Office.vue"

// api接口
import { userFileApi } from "@/api/file/repository"
import type * as Repository from "@/api/file/types/repository"
import { uploadFileApi } from "@/api/file"
import {
  type FileUploadRequestData,
  type ChunkUploadRequestData
  // type FileUploadResponseData,
  // type ChunkUploadCompleteResponseData
} from "@/api/file/types/upload"

import { isNotEmpty } from "@/utils/isEmpty"
import { el } from "element-plus/es/locale"

interface FileListItem {
  id: number
  filename: string
  updateTime: string
  fileSize?: number
  fileType?: number
  isFolder: boolean
  fileCover?: string
}

interface UploadHistoryItem {
  filename: string
  size?: number
  status: "success" | "error"
  url?: string
  key?: string
}

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

// 文件接受类型
const fileAccept = ".jpg,.jpeg,.png,.gif,.zip,.doc,.docx,.pdf"

// 状态变量
const loading = ref(false)
const searchKeyword = ref("")
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const fileList = ref<FileListItem[]>([])
const selectedFiles = ref<FileListItem[]>([])

// 上传相关状态
const currentFile = ref<UploadRawFile | null>(null)
const statusText = ref("")
const uploadHistory = ref<UploadHistoryItem[]>([])

// 路径历史
const pathHistory = ref<Array<{ id: number; name: string }>>([])
const currentPath = ref<number[]>([0]) // 0表示根目录

// 新建文件夹对话框
const folderDialog = reactive({
  visible: false,
  title: "新建文件夹",
  form: {
    name: ""
  }
})

// 重命名对话框
const renameDialog = reactive({
  visible: false,
  fileId: 0,
  form: {
    name: ""
  }
})

// 右键菜单相关状态
const contextMenu = reactive({
  visible: false,
  x: 0,
  y: 0,
  row: null as FileListItem | null
})

// 移动对话框
const moveDialog = reactive({
  visible: false,
  form: {
    targetFolderId: -1
  }
})

// 添加文件夹列表
const folderList = ref<Array<{ id: number; name: string }>>([])

// 判断文件是否可预览
const isPreviewAble = (file: FileListItem): boolean => {
  if (file.isFolder) return false
  return file.fileType === 3 || file.fileType === 1
}

// 根据文件扩展名判断文件类型
const getFileType = (ext: string): number => {
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

// 修改获取文件列表的函数
const loadFileList = async () => {
  loading.value = true
  try {
    const parentId = currentPath.value[currentPath.value.length - 1]

    // 获取文件列表
    const fileResponse = await userFileApi.getFileAndFolderList({
      id: parentId,
      page: currentPage.value,
      size: pageSize.value
    })

    // 处理文件和文件夹
    const filesAndFolders = await Promise.all(
      fileResponse.data.list.map(async (item) => {
        // 文件夹的 RepositoryId 为 0
        if (!item.repository_id) {
          // 文件夹
          const folderSize = await getFolderSize(item.id)
          return {
            id: item.id,
            filename: item.name,
            fileType: 0, // 文件夹类型
            fileSize: folderSize, // 文件夹大小
            updateTime: timestampToDate(item.update_time) || new Date().toISOString(),
            fileCover: "", // 文件夹没有封面
            isFolder: true
          }
        } else {
          // 文件
          return {
            id: item.id,
            filename: item.name,
            fileType: getFileType(item.ext),
            fileSize: item.size,
            updateTime: timestampToDate(item.update_time) || new Date().toISOString(),
            fileCover: item.path,
            isFolder: false
          }
        }
      })
    )

    fileList.value = filesAndFolders
    total.value = fileResponse.data.count
  } catch (error) {
    console.error("获取文件列表失败:", error)
    ElMessage.error("获取文件列表失败")
  }
  loading.value = false
}

// 时间戳转日期
const timestampToDate = (timestamp: number): string => {
  return new Date(timestamp * 1000).toLocaleString()
}

// 上传状态
const uploadStatus = ref("")
const uploadProgress = ref(0)

// 定义分片大小为 5MB
const CHUNK_SIZE = 5 * 1024 * 1024

// 添加分片上传相关的状态
const chunks = ref<Blob[]>([])
const currentChunkIndex = ref(0)
const uploadId = ref("")
const uploadedETags = ref<string[]>([])

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
  console.log("parentId: ", parentId)
  await userFileApi.saveRepository({
    parent_id: Number(parentId),
    repository_id: Number(data.repository_id),
    name: file.name
  })

  // const historyItem: UploadHistoryItem = {
  //   filename: file.name,
  //   size: file.size,
  //   status: "success",
  //   url: data.url,
  //   key: data.key
  // }
  // uploadHistory.value.unshift(historyItem)

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

// 新建文件夹
const handleCreateFolder = () => {
  folderDialog.form.name = ""
  folderDialog.visible = true
}

const confirmCreateFolder = async () => {
  if (!folderDialog.form.name) {
    ElMessage.warning("请输入文件夹名称")
    return
  }

  try {
    const parentId = currentPath.value[currentPath.value.length - 1]
    await userFileApi.createFolder({
      parent_id: parentId,
      name: folderDialog.form.name
    })

    folderDialog.visible = false
    ElMessage.success("创建成功")
    loadFileList()
  } catch (error) {
    console.log("创建失败:", error)
    ElMessage.error("创建失败")
  }
}

// 处理文件/文件夹点击
const handleRowClick = (row: FileListItem) => {
  // console.log("点击的行:", row)
  if (row.isFolder) {
    currentPath.value.push(row.id)
    pathHistory.value.push({
      id: row.id,
      name: row.filename
    })
    currentPage.value = 1
    loadFileList()
  }
}

// 处理面包屑点击
const handleRootClick = () => {
  pathHistory.value = []
  currentPath.value = [0]
  currentPage.value = 1
  loadFileList()
}

const handlePathClick = (index: number) => {
  pathHistory.value = pathHistory.value.slice(0, index + 1)
  currentPath.value = currentPath.value.slice(0, index + 2)
  currentPage.value = 1
  loadFileList()
}

// 其他事件处理器
const handleRefresh = () => {
  loadFileList()
}

// 处理搜索
const handleSearch = debounce(() => {
  // 调用搜索api接口
  const keyword = searchKeyword.value.trim()
}, 300)

const handleSelectionChange = (selection: FileListItem[]) => {
  selectedFiles.value = selection
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  loadFileList()
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
  loadFileList()
}

// 添加文件操作相关的方法
const handleDownload = () => {
  if (selectedFiles.value.length === 0) {
    ElMessage.warning("请选择要下载的文件")
    return
  }
  // TODO: 实现文件下载逻辑
}

const handleDelete = async () => {
  if (selectedFiles.value.length === 0) {
    ElMessage.warning("请选择要删除的文件")
    return
  }

  try {
    for (const file of selectedFiles.value) {
      await userFileApi.deleteFile({
        id: file.id
      })
    }
    ElMessage.success("删除成功")
    loadFileList()
  } catch (error) {
    console.error("删除失败:", error)
    ElMessage.error("删除失败")
  }
}

const handleMove = async () => {
  if (moveDialog.form.targetFolderId === -1) {
    ElMessage.warning("请选择目标文件夹")
    return
  }

  try {
    for (const file of selectedFiles.value) {
      console.log("移动id:%d-->%d", file.id, moveDialog.form.targetFolderId)
      await userFileApi.moveFile({
        id: file.id,
        parent_id: moveDialog.form.targetFolderId
      })
    }
    ElMessage.success("移动成功")
    moveDialog.visible = false
    loadFileList()
  } catch (error) {
    console.error("移动失败:", error)
    ElMessage.error("移动失败")
  }
}

const handleRename = async (fileId: number, newName: string) => {
  try {
    await userFileApi.updateFileName({
      id: fileId,
      name: newName
    })
    ElMessage.success("重命名成功")
    renameDialog.visible = false
    loadFileList()
  } catch (error) {
    console.error("重命名失败:", error)
    ElMessage.error("重命名失败")
  }
}

// 处理右键点击事件
const handleContextMenu = (row: FileListItem, event: MouseEvent) => {
  console.log("右键点击事件", row, event)
  event.preventDefault()
  event.stopPropagation()

  // 选中该行
  selectedFiles.value = [row]

  // 如果菜单已经显示，先关闭它
  if (contextMenu.visible) {
    closeContextMenu()
  }

  // 设置新的菜单位置和数据
  contextMenu.x = event.pageX
  contextMenu.y = event.pageY
  contextMenu.row = row
  contextMenu.visible = true

  // 添加全局点击事件监听
  document.addEventListener("click", closeContextMenu, { once: true })
}

// 修改关闭菜单函数
const closeContextMenu = () => {
  contextMenu.visible = false
  contextMenu.row = null
  // 移除全局点击事件监听
  document.removeEventListener("click", closeContextMenu)
}

// 打开重命名对话框
const openRenameDialog = () => {
  if (!contextMenu.row) return
  renameDialog.fileId = contextMenu.row.id
  renameDialog.form.name = contextMenu.row.filename
  renameDialog.visible = true
  closeContextMenu()
}

// 打开移动对话框
const openMoveDialog = () => {
  if (!contextMenu.row) return
  moveDialog.form.targetFolderId = 0 // 默认选项
  moveDialog.visible = true
  closeContextMenu()
}

// 获取文件夹列表
const loadFolderList = async () => {
  try {
    // 获取文件夹列表
    const response = await userFileApi.getFolderList({
      id: 0
    })
    folderList.value = [
      {
        id: 0,
        name: "/"
      },
      ...response.data.list
    ]
  } catch (error) {
    console.error("获取文件夹列表失败:", error)
  }
}

// 文件预览相关状态
const previewDialog = ref({ visible: false })
const previewUrl = ref("")
const previewResource = ref(null)
const previewComponent = ref(Default)

const previewFile = (file) => {
  previewUrl.value = file.url // 假设文件对象有一个 url 属性
  previewResource.value = file // 传递文件资源
  if (file.type === "audio") {
    previewComponent.value = Audio
  } else if (file.type === "image") {
    previewComponent.value = Image
  } else if (file.type === "office") {
    previewComponent.value = Office
  } else {
    previewComponent.value = Default
  }
  previewDialog.value.visible = true
}

// 初始化
onMounted(() => {
  loadFileList()
  loadFolderList()
})

// 在组件卸载时清理事件监听
onUnmounted(() => {
  closeContextMenu()
})
</script>

<style lang="scss" scoped>
.app-container {
  padding: 20px;

  .main-ui-header {
    display: flex;
    align-items: center;
    margin-bottom: 20px;
    gap: 16px;

    .search-input {
      width: 300px;
      margin-right: auto;

      :deep(.el-input__wrapper) {
        background-color: #f5f7fa;
      }

      :deep(.el-input__inner) {
        &::placeholder {
          color: #909399;
        }
      }
    }
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

  .breadcrumb {
    margin-bottom: 20px;

    :deep(.el-breadcrumb__item) {
      .el-breadcrumb__inner {
        display: flex;
        align-items: center;
        gap: 4px;
        cursor: pointer;

        .el-icon {
          margin-right: 4px;
        }
      }
    }
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
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
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
