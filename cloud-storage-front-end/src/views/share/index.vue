<template>
  <div class="app-container">
    <div class="share-header">
      <div class="top">
        <el-checkbox class="all-selected" v-model="allSelected" @change="toggleSelectAll" />
        <span class="title"> 全部文件 </span>
        <span class="count">已加载{{ tableDataNum }}个 </span>
        <el-icon class="refresh-icon" @click="refreshList"><RefreshIcon /></el-icon>
      </div>
      <div class="actions">
        <button @click="copyLink">复制链接</button>
        <button @click="cancelShare">取消分享</button>
        <button @click="exportLink">导出链接</button>
      </div>
    </div>
    <div class="share-body">
      <div class="left-container">
        <div class="file-list-container" @scroll="handleScroll">
          <div class="file-list">
            <el-card v-for="item in fileList" 
                     :key="item.id" 
                     class="file-card"
                     :class="{ 'is-selected': isSelected(item) }"
                     shadow="hover"
                     @click="selectFile(item)">
              <div class="file-info">
                <div class="file-header">
                  <el-avatar :src="item.avatar" size="small" />
                  <span class="owner-name">{{ item.owner }}</span>
                </div>
                <div class="file-name">{{ item.name }}</div>
                <div class="file-meta"> 
                  <span>大小: {{ formatFileSize(item.size) }} </span>
                  <span>浏览次数: {{ item.click_num }}</span>
                </div>
              </div>
            </el-card>
          </div>
          <div v-if="loading" class="loading">加载中...</div>
          <div v-if="noMore" class="no-more">没有更多文件了</div>
        </div>
      </div>
      <div class="right-detail">
        <ShareDetail :selectFileProps="selectFiles" />
      </div>
    </div>
    <FileDialog
      :show="dialogConfig.show"
      :title="dialogConfig.title"
      :buttons="dialogConfig.buttons"
      :show-cancel="true"
      width="400px"
      @close="dialogConfig.show = false"
    >
      <div class="delete-info">
        <div class="info">
          <span class="text">取消分享后，该条记录将被删除，好友将无法访问此链接。您确认要取消分享吗？</span>
        </div>
      </div>
    </FileDialog>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, computed, onMounted, watch, onUnmounted } from "vue"
import { ElMessage } from "element-plus"
import FileDialog from "@/components/FileDialog/index.vue"
import useClipboard from "vue-clipboard3"
import ShareDetail from "./detail.vue"
import { shareApi } from "@/api/share/share"
import type * as Share from "@/api/share/types/share"
import { formatFileSize } from "@/utils/format/formatFileSize"
import { formatTime } from "@/utils/format/formatTime"
import dayjs from 'dayjs'
import { time } from "console"
import { Refresh as RefreshIcon } from '@element-plus/icons-vue'

interface ShareItem {
  id: number
  repository_id: number
  name: string
  ext: string
  size: number
  path: string
  expired_time: number
  update_time: number
  //
  owner: string
  avatar: string
  //
  click_num: number // 浏览次数
  //
  code: string
  selected?: boolean
}

const selectFiles = ref<ShareItem[]>([])
const allSelected = ref(false)
const shareIds = ref<string>("")
const tableDataNum = ref<number>(0)
const loading = ref(false)
const noMore = ref(false)
const page = ref(1)
const pageSize = ref(5) // 每页加载的文件数量

const dialogConfig = reactive({
  show: false,
  title: "确认取消分享",
  showCancel: true,
  buttons: [
    {
      type: "primary",
      text: "确认",
      click: () => {
        deleteShare()
      }
    }
  ]
})

const fileList = ref<ShareItem[]>([])

// const paginatedFileList = computed(() => {
//   const start = (page.value - 1) * pageSize.value
//   return fileList.value.slice(start, start + pageSize.value)
// })

// const tableRows = computed(() => {
//   return paginatedFileList.value.map((file) => ({
//     ...file,
//   }))
// })

// const handleSelectionChange = (val) => {
//   selectFiles.value = val
// }

// const multipleTable = ref()
// const handleRowClick = (row) => {
//   selectFiles.value = [row]
//   //通过ref绑定来操作bom元素
//   multipleTable.value.toggleRowSelection(row)
// }

const selectedFileIds = computed(() => {
  return selectFiles.value.map((file) => file.id).join(",")
})

const cancelShare = () => {
  console.log("selectedFileIds", selectedFileIds.value)
  shareIds.value = selectedFileIds.value // 使用计算属性
  dialogConfig.show = true
}

const deleteShare = async () => {
  try {
    await shareApi.deleteShare({ id: shareIds.value})

    ElMessage.success("取消分享成功")
    dialogConfig.show = false
    refreshList()
  } catch (error) {
    ElMessage.error("取消分享失败")
  }
}

// 复制链接
const copyLink = () => {
  if (selectFiles.value.length === 0) {
    ElMessage.warning("请先选择文件")
    return
  }
  const baseUrl = import.meta.env.VITE_API_URL_3
  const links = selectFiles.value.map((file) => `链接: ${baseUrl}/share_service/v1/share/basic/detail?id=${file.id}&code=${file.code}`).join("\n")

  navigator.clipboard
    .writeText(links)
    .then(() => {
      ElMessage.success("链接已复制到剪贴板")
    })
    .catch(() => {
      ElMessage.error("复制链接失败")
    })
}

// 导出链接
const exportLink = () => {
  if (selectFiles.value.length === 0) {
    ElMessage.warning("请先选择文件")
    return
  }

  const csvContent =
    "data:text/csv;charset=utf-8," +
    selectFiles.value
      .map(
        (file) =>
          `${file.name},${formatTime(file.update_time)},${formatExpiredTime(file.update_time,file.expired_time)},${file.click_num},${file.code}`
      )
      .join("\n")

  const encodedUri = encodeURI(csvContent)
  const link = document.createElement("a")
  link.setAttribute("href", encodedUri)
  link.setAttribute("download", "exported_files.csv")
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  ElMessage.success("文件已导出为CSV格式")
}

// 计算过期时间
const formatExpiredTime = (time,expiredTime) => {
  if (!expiredTime) return '-'
  const expirationDate = dayjs(time).add(expiredTime, 'second') // 将当前时间加上过期时间（秒）
  return expirationDate.format('YYYY-MM-DD HH:mm:ss') 
}
const fileListContainer = ref<HTMLElement | null>(null)

// 处理滚动事件
const handleScroll = (e: Event) => {
  const target = e.target as HTMLElement
  const scrollHeight = target.scrollHeight
  const scrollTop = target.scrollTop
  const clientHeight = target.clientHeight
  
  // 当滚动到距离底部100px时加载更多
  if (scrollHeight - scrollTop - clientHeight < 100 && !loading.value && !noMore.value) {
    console.log('Loading more files...')
    loadFiles()
  }
}

// 加载文件
const loadFiles = async () => {
  if (loading.value || noMore.value) return
  loading.value = true
  console.log('Loading more files...') 

  try {
    const response = await shareApi.getShareList({ 
      page: page.value, 
      page_size: pageSize.value 
    })
    const { list } = response.data

    if (list && list.length > 0) {
      fileList.value.push(...list)
      tableDataNum.value = fileList.value.length
      page.value += 1
    } else {
      noMore.value = true
    }
  } catch (error) {
    console.error('Failed to load files:', error)
    ElMessage.error("加载文件失败")
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  // 初始加载
  loadFiles()
})

// 保存分享
const saveShare = async (repositoryId: number, parentId: number) => {
  try {
    const response = await shareApi.saveShare({
      repository_id: repositoryId,
      parent_id: parentId
    })

    if (response.data) {
      ElMessage.success("保存分享成功")
    }
  } catch (error) {
    ElMessage.error("保存分享失败")
  }
}

watch(fileList, (newVal) => {
  console.log('fileList changed:', newVal)
}, { deep: true })

// 修改选择文件相关的方法
const isSelected = (file: ShareItem) => {
  return selectFiles.value.some(item => item.id === file.id)
}

// 选择某一行文件
const selectFile = (file: ShareItem) => {
  const index = selectFiles.value.findIndex(item => item.id === file.id)
  if (index === -1) {
    console.log("selectFile", file)
    selectFiles.value = [file] // 只保留当前选中的文件
  } else {
    selectFiles.value = [] // 取消选择
  }
  // 更新全选状态
  allSelected.value = selectFiles.value.length === fileList.value.length
}

// 监控文件列表变化，更新全选状态
watch(fileList, (newVal) => {
  allSelected.value = newVal.length > 0 && newVal.length === selectFiles.value.length
})

// 修改 toggleSelectAll 方法
const toggleSelectAll = () => {
  if (allSelected.value) {
    selectFiles.value = [...fileList.value]
  } else {
    selectFiles.value = []
  }
}

// 刷新列表的函数
const refreshList = () => {
  page.value = 1; // 重置页码
  fileList.value = []; // 清空当前文件列表
  loadFiles(); // 重新加载文件
}
</script>

<style lang="scss">
.app-container {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.share-header {
  height: 60px;
  padding: 10px 20px;
  background-color: #ffffff;
  border-bottom: 1px solid #dcdcdc;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  border-radius: 8px;

  .top {
    display: flex;
    align-items: center;
    line-height: 36px;
    color: #333;
    font-weight: bold;
    font-size: 14px;
  }

  .title {
    margin-right: 10px;
  }

  .all-selected {
    margin-right: 10px;
  }

  .refresh-icon {
    cursor: pointer;
    margin-left: 2px;
  }

  .actions {
    display: flex;
    align-items: center;

    button {
      margin-left: 10px;
      padding: 8px 12px;
      border: none;
      border-radius: 4px;
      background-color: #007bff;
      color: #ffffff;
      cursor: pointer;
      transition: background-color 0.3s;
      font-size: 14px;

      &:hover {
        background-color: #0056b3;
      }
    }
  }
}

.share-body {
  display: flex;
  height: calc(100% - 60px);
  background-color: #ffffff;

  .left-container {
    flex: 1;
    display: flex;
    flex-direction: column;

    .file-list-container {
      flex: 1;
      overflow-y: auto;
      height: calc(100vh - 180px);
      padding: 20px;
      background-color: #f9f9f9;
      border: 1px solid #dcdcdc;
      border-radius: 8px;

      .file-list {
        display: flex;
        flex-direction: column; /* 每个文件占一行 */
        gap: 10px; /* 文件之间的间距 */
      }
    }
  }
}

.right-detail {
  width: 240px;
  overflow: hidden;
  display: inline-block;
  position: relative;
  height: 100%;
  border-radius: 8px;
  font-size: 14px;
  border: 1px solid #dcdcdc;
}

.delete-info {
  padding: 10px;
  color: #333;
  text-align: center;
}

.file-card {
  cursor: pointer;
  transition: all 0.3s;
  border: 2px solid transparent;

  &.is-selected {
    border-color: #409EFF;
    background-color: #ecf5ff;
  }
}

.file-card:hover {
  transform: translateY(-5px);
}

.file-info {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.file-header {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 0;

  .el-avatar {
    flex-shrink: 0; // 防止头像被压缩
    border: 1px solid #eee;
  }
}

.owner-name {
  font-weight: 500;
  flex: 1; // 让名字占据剩余空间
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-name {
  font-size: 16px;
  word-break: break-all;
  line-height: 1.4;
  margin: 4px 0;
}

.file-meta {
  display: flex;
  justify-content: space-between;
  color: #666;
  font-size: 14px;
  align-items: center; // 垂直居中对齐
  padding: 4px 0;

  span {
    display: flex;
    align-items: center;
    gap: 4px;
  }
}

.loading {
  text-align: center;
  color: #666;
  padding: 20px;
}

.no-more {
  text-align: center;
  color: #999;
  padding: 20px;
}


//滚动条样式
::-webkit-scrollbar {
  width: 0.5rem;
  height: 0.5rem;
  background: rgba(255, 255, 255, 0.6);
}

::-webkit-scrollbar-track {
  border-radius: 0;
}

::-webkit-scrollbar-thumb {
  border-radius: 0;
  background-color: rgb(218, 218, 218);
  transition: all 0.2s;
  border-radius: 8px;

  &:hover {
    background-color: rgb(172, 172, 172);
  }
}
</style>
