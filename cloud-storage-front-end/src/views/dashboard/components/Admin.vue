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
              <div class="icon-container">
                <template v-if="row.fileType === 3 && row.fileCover">
                  <!-- 仅图片类型(3)显示封面 -->
                  <Icon :cover="row.fileCover" :width="32" />
                </template>
                <template v-else>
                  <!-- 其他类型显示对应图标 -->
                  <Icon :file-type="row.fileType" :width="32" />
                </template>
              </div>
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
      </el-table>
    </div>

    <!-- 上传进度 -->
    <div v-if="currentFile" ref="uploadProgressRef" class="upload-progress">
      <div class="file-info">
        <span class="filename">{{ currentFile.name }}</span>
        <span class="filesize">{{ formatFileSize(currentFile.size) }}</span>

        <div class="progress-container">
          <el-progress
            :percentage="uploadProgress"
            :stroke-width="8"
            :show-text="false"
            :color="uploadProgress === 100 ? '#67C23A' : '#FBBC4D'"
          />
          <div class="percentage-text">{{ uploadProgress }}%</div>
        </div>

        <div class="status-text" :class="{ success: uploadProgress === 100 }">
          <el-icon v-if="uploadProgress === 100" class="status-icon"><CircleCheckFilled /></el-icon>
          <el-icon v-else class="status-icon spinning"><Loading /></el-icon>
          {{ uploadProgress === 100 ? "上传完成" : statusText }}
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
          <el-input
            ref="renameInput"
            v-model="renameDialog.form.name"
            @keyup.enter="handleRename(renameDialog.fileId, renameDialog.form.name)"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="renameDialog.visible = false">取消</el-button>
        <el-button type="primary" @click="handleRename(renameDialog.fileId, renameDialog.form.name)"> 确定 </el-button>
      </template>
    </el-dialog>

    <!-- 移动对话框 -->
    <el-dialog v-model="moveDialog.visible" title="移动文件" width="30%">
      <el-form :model="moveDialog.form" label-width="90px">
        <el-form-item label="目标文件夹">
          <el-select v-model="moveDialog.form.targetFolderId" placeholder="选择目标文件夹">
            <el-option v-for="folder in filesFolders" :key="folder.id" :label="folder.name" :value="folder.id" />
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
        <li @click="handleShare">
          <el-icon><Share /></el-icon>
          分享
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
    <el-dialog v-model="previewDialog.visible" title="文件预览" width="80%">
      <Preview :resource="previewResource" style="height: 330px" />
      <template #footer>
        <span>
          <el-button @click="previewDialog.visible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 添加分享结果对话框 -->
    <el-dialog v-model="shareDialog.visible" title="分享成功" width="30%">
      <div class="share-info">
        <p>分享链接已创建，有效期7天</p>
        <div class="share-link">
          <el-input v-model="shareDialog.link" readonly>
            <template #append>
              <el-button @click="copyShareLink">复制</el-button>
            </template>
          </el-input>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted, computed, onUnmounted, nextTick } from "vue"
import { ElMessage } from "element-plus"
import type { UploadRequestOptions, UploadRawFile } from "element-plus"
import Icon from "@/components/FileIcon/Icon.vue"
import { debounce } from "lodash-es"
import Preview from "@/views/dashboard/components/file-preview/Preview.vue"

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
import { formatFileSize } from "@/utils/format/formatFileSize"
import { timestampToDate } from "@/utils/format/formatTime"
import { useUserStore } from "@/store/modules/user"
import { shareApi } from "@/api/share/share"
import { Share, CircleCheckFilled, Loading } from "@element-plus/icons-vue"
import { historyFileApi } from "@/api/file/history"
import type * as History from "@/api/file/types/history"

const userStore = useUserStore()
const capacity = computed(() => userStore.capacity)

// 添加 ref 用于获取输入框元素
const renameInput = ref<HTMLInputElement>()

// 添加上传进度区域的引用
const uploadProgressRef = ref<HTMLElement>()

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

interface UploadHistoryItem {
  filename: string
  size?: number
  status: "success" | "error"
  url?: string
  key?: string
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
const filesFolders = ref<Array<{ id: number; name: string }>>([
  {
    id: 0,
    name: "根目录"
  }
])

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

// 获取文件列表的函数
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

    console.log("文件列表", fileResponse.data)
    // 处理文件和文件夹
    const filesAndFolders = await Promise.all(
      fileResponse.data.list.map(async (item) => {
        // 文件夹的 RepositoryId 为 0
        if (!item.repository_id) {
          // 文件夹
          filesFolders.value.push({
            id: item.id,
            name: item.name
          })
          const folderSize = await getFolderSize(item.id)
          return {
            id: item.id,
            filename: item.name,
            fileType: 0, // 文件夹类型
            fileSize: folderSize, // 文件夹大小
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
            fileSize: item.size,
            updateTime: timestampToDate(item.update_time) || new Date().toISOString(),
            fileCover: item.path,
            isFolder: false,
            ext: item.ext,
            repository_id: item.repository_id
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

// 添加分享对话框状态
const shareDialog = reactive({
  visible: false,
  link: "",
  code: ""
})

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
      console.log("文件大于20M")
      await handleChunkUpload(file)
    } else {
      await handleNormalUpload(file)
    }

    // 上传开始后，等待DOM更新后滚动到进度区域
    nextTick(() => {
      if (uploadProgressRef.value) {
        uploadProgressRef.value.scrollIntoView({ behavior: "smooth", block: "center" })
      }
    })
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

  // 设置初始状态
  uploadProgress.value = 0
  statusText.value = "正在上传..."

  const response = await uploadFileApi.upload(fileUploadRequestData, {
    onUploadProgress: (progressEvent) => {
      if (progressEvent.total) {
        // 限制上传阶段最多到95%，为服务器处理预留空间
        const progress = (progressEvent.loaded / progressEvent.total) * 95
        uploadProgress.value = Math.round(progress)

        // 更新状态文本
        if (progress > 90) {
          statusText.value = "处理中..."
        } else {
          statusText.value = "正在上传..."
        }
      }
    }
  })

  // 上传完成后服务器处理中
  uploadProgress.value = 98
  statusText.value = "即将完成..."

  handleUploadSuccess(response.data, file)
}

// 处理分片上传
const handleChunkUpload = async (file: File) => {
  // 初始化分片上传
  try {
    // 设置初始状态
    uploadProgress.value = 0
    statusText.value = "正在初始化..."

    const initResponse = await uploadFileApi.initiateMultipart({
      file_name: file.name,
      file_size: file.size,
      metadata: JSON.stringify({ type: file.type })
    })

    const { upload_id, key } = initResponse.data
    uploadId.value = upload_id

    // 更新进度到初始化完成
    uploadProgress.value = 5
    statusText.value = "准备分片..."

    // 将文件分片
    const chunkCount = Math.ceil(file.size / CHUNK_SIZE)
    chunks.value = []
    uploadedETags.value = new Array(chunkCount)

    for (let i = 0; i < chunkCount; i++) {
      const start = i * CHUNK_SIZE
      const end = Math.min(start + CHUNK_SIZE, file.size)
      chunks.value.push(file.slice(start, end))
    }

    // 上传所有分片 - 分片上传占总进度的85%
    for (let i = 0; i < chunks.value.length; i++) {
      currentChunkIndex.value = i
      const chunkUploadRequestData: ChunkUploadRequestData = {
        upload_id: uploadId.value,
        chunk_index: i + 1,
        key: key,
        file: chunks.value[i] as unknown as File
      }

      statusText.value = `正在上传第 ${i + 1}/${chunks.value.length} 个分片`

      const chunkResponse = await uploadFileApi.uploadPart(chunkUploadRequestData)
      if (isNotEmpty(chunkResponse.data)) {
        uploadedETags.value[i] = chunkResponse.data.etag
      }

      // 更新进度 - 分片上传占进度的5~90%
      const progress = 5 + ((i + 1) / chunks.value.length) * 85
      uploadProgress.value = Math.round(progress)
    }

    // 更新进度到合并阶段
    uploadProgress.value = 95
    statusText.value = "正在合并分片..."

    // 完成分片上传
    const completeResponse = await uploadFileApi.completeMultipart({
      upload_id: uploadId.value,
      key: key,
      etags: uploadedETags.value
    })

    // 更新进度到即将完成
    uploadProgress.value = 98
    statusText.value = "即将完成..."

    handleUploadSuccess(completeResponse.data, file)
  } catch (error) {
    console.error(error)
    handleUploadError(error, file)
  }
}

// 修改处理上传成功的函数
const handleUploadSuccess = async (data: any, file: File) => {
  try {
    // 保存文件关联信息前显示处理中状态
    uploadProgress.value = 99
    statusText.value = "保存文件信息..."

    const parentId = currentPath.value[currentPath.value.length - 1]
    console.log("parentId: ", parentId)

    // 保存到仓库
    await userFileApi.saveRepository({
      parent_id: Number(parentId),
      repository_id: Number(data.repository_id),
      name: file.name
    })

    // 添加上传历史记录
    await historyFileApi.uploadHistoryFile({
      repository_id: Number(data.repository_id),
      file_name: file.name,
      size: file.size,
      status: 1 // 上传成功
    })

    // 更新进度显示为完成
    uploadStatus.value = "success"
    uploadProgress.value = 100
    statusText.value = "上传成功"

    // 更新容量显示
    capacity.value.now_volume += file.size

    handleRefresh()
    ElMessage.success("文件上传成功")
  } catch (error) {
    console.error("保存文件信息失败:", error)
    ElMessage.error("保存文件信息失败")
  }

  resetUploadState()
}

// 处理上传错误
const handleUploadError = async (error: any, file: File) => {
  uploadStatus.value = "error"
  statusText.value = "上传失败"

  const historyItem: UploadHistoryItem = {
    filename: file.name,
    size: file.size,
    status: "error"
  }
  uploadHistory.value.unshift(historyItem)

  try {
    // 添加上传历史记录
    await historyFileApi.uploadHistoryFile({
      repository_id: -1,
      file_name: file.name,
      size: file.size,
      status: 0 // 上传失败
    })
  } catch (error) {
    console.error("添加上传历史记录失败:", error)
  }

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
    // 文件夹跳转
    currentPath.value.push(row.id)
    pathHistory.value.push({
      id: row.id,
      name: row.filename
    })
    currentPage.value = 1
    loadFileList()
  } else {
    // 文件预览
    previewFile(row)
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
const handleSearch = debounce(async () => {
  const keyword = searchKeyword.value.trim()
  if (!keyword) {
    // 关键词为空时，重新加载原始文件列表
    loadFileList()
    return
  }

  loading.value = true
  try {
    // 调用原始文件列表API后在前端过滤
    // 注意：这是前端实现的临时方案，后续应改为调用后端搜索API
    const parentId = currentPath.value[currentPath.value.length - 1]
    const response = await userFileApi.getFileAndFolderList({
      id: parentId,
      page: 1, // 搜索时获取较多数据
      size: 100 // 获取更多条目便于搜索
    })

    // 在前端过滤符合搜索条件的文件
    const filteredList = response.data.list.filter((item) => item.name.toLowerCase().includes(keyword.toLowerCase()))

    // 处理搜索结果
    const filesAndFolders = await Promise.all(
      filteredList.map(async (item) => {
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
            fileSize: item.size,
            updateTime: timestampToDate(item.update_time) || new Date().toISOString(),
            fileCover: item.path,
            isFolder: false,
            ext: item.ext,
            repository_id: item.repository_id
          }
        }
      })
    )

    fileList.value = filesAndFolders
    total.value = filesAndFolders.length

    // 显示搜索结果提示
    if (filesAndFolders.length === 0) {
      ElMessage.info(`没有找到名称包含"${keyword}"的文件`)
    } else {
      ElMessage.success(`搜索到 ${filesAndFolders.length} 个结果`)
    }
  } catch (error) {
    console.error("搜索文件失败:", error)
    ElMessage.error("搜索文件失败")
  } finally {
    loading.value = false
  }
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

// 修改下载处理函数
const handleDownload = async () => {
  if (selectedFiles.value.length === 0) {
    ElMessage.warning("请选择要下载的文件")
    return
  }

  // 目前只支持单文件下载
  if (selectedFiles.value.length > 1) {
    ElMessage.warning("目前只支持单个文件下载，请只选择一个文件")
    return
  }

  const file = selectedFiles.value[0]

  // 文件夹不能直接下载
  if (file.isFolder) {
    ElMessage.warning("暂不支持文件夹下载")
    return
  }

  try {
    console.log("准备下载文件，repository_id:", file.repository_id)

    // 调用后端获取下载链接
    const response = await uploadFileApi.getDownloadUrl({
      repository_id: file.repository_id
    })

    console.log("获取下载链接响应:", response.code)

    if (response.data && response.data.url) {
      // 创建一个隐藏的a标签用于下载
      const link = document.createElement("a")
      link.href = response.data.url
      link.download = file.filename // 设置下载文件名
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)

      ElMessage.success("下载开始")
    } else {
      console.error("响应中没有有效的下载链接:", response)
      ElMessage.error("获取下载链接失败")
    }
  } catch (error) {
    console.error("下载文件失败:", error)
    ElMessage.error("下载文件失败: " + (error instanceof Error ? error.message : String(error)))
  }
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
    // 更新容量显示
    capacity.value.now_volume -= selectedFiles.value.reduce((total, file) => total + (file.fileSize || 0), 0)
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

// 修改打开重命名对话框的函数
const openRenameDialog = () => {
  if (!contextMenu.row) return
  renameDialog.fileId = contextMenu.row.id
  renameDialog.form.name = contextMenu.row.filename
  renameDialog.visible = true

  // 使用 nextTick 确保对话框完全渲染后再聚焦
  nextTick(() => {
    // 获取 el-input 组件的输入框元素并聚焦
    const input = renameInput.value?.$el.querySelector("input")
    if (input) {
      input.focus()
      input.select() // 可选：自动选中所有文本
    }
  })

  closeContextMenu()
}

// 打开移动对话框
const openMoveDialog = () => {
  if (!contextMenu.row) return
  moveDialog.form.targetFolderId = 0 // 默认选项
  moveDialog.visible = true
  closeContextMenu()
}

// 文件预览相关状态
const previewDialog = ref({ visible: false })
const previewResource = ref({
  id: 0,
  name: "",
  type: "",
  url: "",
  size: 0
})

const previewFile = (file: FileListItem) => {
  // 构建预览资源对象
  previewResource.value = {
    id: file.id,
    name: file.filename,
    // type: getFileExtension(file.filename),
    type: file.ext,
    url: file.fileCover || "", // 这里使用文件的访问路径
    size: file.fileSize || 0
  }
  previewDialog.value.visible = true
}

// 处理分享
const handleShare = async () => {
  if (!contextMenu.row) {
    ElMessage.warning("请选择要分享的文件")
    return
  }

  try {
    // 生成随机提取码(6位数)
    const codes_str = Math.random().toString(36).substring(2, 8)
    // 创建分享
    const response = await shareApi.createShare({
      repository_id: contextMenu.row.repository_id,
      user_repository_id: contextMenu.row.id,
      expired_time: 7 * 24 * 60 * 60, // 7天过期时间(秒)
      code: codes_str
    })

    if (response.data.id === "") {
      ElMessage.warning("文件已分享")
      return
    }

    if (response.data) {
      // console.log("response.data", response.data)
      const baseUrl = import.meta.env.VITE_API_URL_3
      const shareLink = `${baseUrl}/share_service/v1/share/basic/detail?id=${response.data.id}&code=${codes_str}`

      // 显示分享对话框
      shareDialog.link = shareLink
      shareDialog.visible = true
    }

    closeContextMenu()
  } catch (error) {
    console.error("创建分享失败:", error)
    ElMessage.error("创建分享失败")
  }
}

// 复制分享链接
const copyShareLink = async () => {
  try {
    await navigator.clipboard.writeText(shareDialog.link)
    ElMessage.success("链接已复制到剪贴板")
  } catch (err) {
    console.error("复制失败:", err)
    ElMessage.error("复制失败")
  }
}

// 初始化
onMounted(() => {
  loadFileList()
  // loadFolderList()
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

  .upload-progress {
    margin: 16px 0;
    padding: 16px;
    border-radius: 8px;
    background-color: #f8f9fa;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
    transition: all 0.3s ease;

    .file-info {
      display: flex;
      flex-direction: column;
      gap: 10px;

      .filename {
        font-size: 15px;
        font-weight: 500;
        color: #303133;
        display: flex;
        align-items: center;

        &:before {
          content: "";
          display: inline-block;
          width: 20px;
          height: 20px;
          margin-right: 8px;
          background-image: url('data:image/svg+xml;utf8,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="%23FBBC4D"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8l-6-6z"/><path d="M14 3v5h5"/></svg>');
          background-size: contain;
          background-repeat: no-repeat;
        }
      }

      .filesize {
        font-size: 13px;
        color: #606266;
        margin-left: 28px;
      }

      .progress-container {
        margin-top: 5px;
        margin-left: 28px;
        margin-right: 28px;
        position: relative;

        .percentage-text {
          position: absolute;
          right: 0;
          top: -18px;
          font-size: 12px;
          color: #909399;
        }
      }

      .status-text {
        font-size: 13px;
        color: #409eff;
        margin-left: 28px;
        font-weight: 500;
        display: flex;
        align-items: center;

        &.success {
          color: #67c23a;
        }

        .status-icon {
          margin-right: 5px;
          font-size: 16px;

          &.spinning {
            animation: spin 1.5s linear infinite;
          }
        }
      }
    }
  }

  @keyframes spin {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(360deg);
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
    margin-top: 10px;
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

  .share-info {
    padding: 20px;

    .share-link {
      margin-top: 15px;

      :deep(.el-input-group__append) {
        padding: 0;

        .el-button {
          margin: 0;
          border: none;
          height: 100%;
        }
      }
    }
  }
}
</style>
