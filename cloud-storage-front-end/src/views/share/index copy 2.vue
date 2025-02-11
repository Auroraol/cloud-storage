<template>
  <div class="app-container">
    <div class="share-header">
      <div class="top">
        <span class="title">全部文件</span>
        <span class="count">已加载{{ tableDataNum }}个</span>
      </div>
      <div class="actions">
        <button @click="copyLink">复制链接</button>
        <button @click="cancelShare">取消分享</button>
        <button @click="exportLink">导出链接</button>
      </div>
    </div>
    <div class="share-body">
      <div class="left-container">
        <div class="file-list" v-infinite-scroll="loadFiles" infinite-scroll-disabled="loading">
          <el-table
            ref="multipleTable"
            :data="tableRows"
            @selection-change="handleSelectionChange"
            @row-click="handleRowClick"
            class="custom-table"
          >
            <el-table-column type="selection" width="55" />
            <el-table-column prop="filename" label="文件名" />
            <el-table-column prop="createTime" label="分享时间" />
            <el-table-column prop="expireTime" label="状态" />
            <el-table-column prop="browseCount" label="浏览次数" />
            <template #empty>
              <div style="text-align: center; margin-top: 20px">
                <div style="margin-top: -30px; font-size: 16px; font-weight: bold">您的分享为空 ~</div>
              </div>
            </template>
          </el-table>
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
import { ref, reactive, computed, onMounted, watch } from "vue"
import { ElMessage } from "element-plus"
import FileDialog from "@/components/FileDialog/index.vue"
import useClipboard from "vue-clipboard3"
import ShareDetail from "./detail.vue"
import { shareApi } from "@/api/share/share"
import type * as Share from "@/api/share/types/share"

interface ShareItem {
  id: number
  repository_id: number
  name: string
  ext: string
  size: number
  path: string
  code: string
  browseCount: number
  saveCount: number
  downloadCount: number
  expireTime: Date
  createTime: Date
}

const shareTableRef = ref()
const selectFiles = ref<ShareItem[]>([])
const shareIds = ref<string>("")
const tableDataNum = ref<number>(0)
const loading = ref(false)
const noMore = ref(false)
const page = ref(1)
const pageSize = ref(10) // 每页加载的文件数量
const shareUrl = ref<string>(document.location.origin + "/share/")

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

const paginatedFileList = computed(() => {
  const start = (page.value - 1) * pageSize.value
  return fileList.value.slice(start, start + pageSize.value)
})

const tableRows = computed(() => {
  return paginatedFileList.value.map((file) => ({
    ...file,
    createTime: formatDate(file.createTime.toISOString()),
    expireTime: formatDate(file.expireTime.toISOString())
  }))
})

const handleSelectionChange = (val) => {
  selectFiles.value = val
}

const multipleTable = ref()
const handleRowClick = (row) => {
  selectFiles.value = [row]
  //通过ref绑定来操作bom元素
  multipleTable.value.toggleRowSelection(row)
}

const selectedFileIds = computed(() => {
  return selectFiles.value.map((file) => file.fileId).join(",")
})

const cancelShare = () => {
  console.log("selectedFileIds", selectedFileIds.value)
  shareIds.value = selectedFileIds.value // 使用计算属性
  dialogConfig.show = true
}

const deleteShare = async () => {
  try {
    // 这里需要后端提供删除分享的 API
    // await shareApi.deleteShare({ id: shareIds.value })

    ElMessage.success("取消分享成功")
    dialogConfig.show = false
    reloadTable()
  } catch (error) {
    ElMessage.error("取消分享失败")
  }
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

const reloadTable = () => {
  shareTableRef.value?.reload()
}

// 复制链接
const copyLink = () => {
  if (selectFiles.value.length === 0) {
    ElMessage.warning("请先选择文件")
    return
  }

  // 链接: https://pan.baidu.com/s/1OtdlYyBNGL0NPpIiLgyG3A 提取码: 93vd
  const links = selectFiles.value.map((file) => `链接: ${shareUrl.value}${file.fileId} 提取码: ${file.code}`).join("\n")
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
          `${file.filename},${file.createTime.toISOString()},${file.expireTime.toISOString()},${file.browseCount},${file.saveCount},${file.downloadCount},${file.code}`
      )
      .join("\n")

  const encodedUri = encodeURI(csvContent)
  const link = document.createElement("a")
  link.setAttribute("href", encodedUri)
  link.setAttribute("download", "exported_files.csv")
  document.body.appendChild(link) // Required for FF

  link.click()
  document.body.removeChild(link)
  ElMessage.success("文件已导出为CSV格式")
}

const loadFiles = async () => {
  if (loading.value || noMore.value) return
  loading.value = true

  try {
    const response = await shareApi.getShareDetail({
      id: page.value // 使用页码作为 id 参数
    })

    if (response.data) {
      const newItem = {
        ...response.data,
        browseCount: 0,
        saveCount: 0,
        downloadCount: 0,
        code: "XXXX", // 从响应中获取
        expireTime: new Date(),
        createTime: new Date()
      }

      fileList.value.push(newItem)
      tableDataNum.value += 1
      page.value += 1
    } else {
      noMore.value = true
      ElMessage.warning("没有更多文件了")
    }
  } catch (error) {
    ElMessage.error("加载文件失败")
  } finally {
    loading.value = false
  }
}

const observer = new IntersectionObserver((entries) => {
  if (entries[0].isIntersecting) {
    loadFiles()
  }
})

onMounted(() => {
  loadFiles() // 初始加载
  const loadingElement = document.querySelector(".loading")
  if (loadingElement) {
    observer.observe(loadingElement)
  }
})

watch(tableRows, (newVal) => {
  if (newVal.length === 0) {
    ElMessage.warning("没有文件")
  }
})

// 创建新分享
const createShare = async (repositoryId: number) => {
  try {
    const response = await shareApi.createShare({
      user_repository_id: repositoryId,
      expired_time: 7 * 24 * 60 * 60 // 7天过期时间
    })

    if (response.data) {
      ElMessage.success("创建分享成功")
      reloadTable()
    }
  } catch (error) {
    ElMessage.error("创建分享失败")
  }
}

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
    display: inline-block;
    line-height: 36px;
    color: #333;
    font-weight: bold;
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
    // padding: 20px;

    .file-list {
      flex: 1;
      overflow-y: auto;
      min-height: 300px;
      border: 1px solid #dcdcdc;
      border-radius: 8px;
      background-color: #f9f9f9;
      padding: 10px;

      .custom-table {
        border: none;
        .el-table__header {
          background-color: #f0f4f8;
          border-bottom: 2px solid #007bff;
        }
        .el-table__body {
          tr {
            &:hover {
              background-color: #e6f7ff;
            }
            td {
              padding: 12px;
              border-bottom: 1px solid #dcdcdc;
            }
          }
        }
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
  // margin-left: 20px;
  // background-color: #f9f9f9;
  // background-color: #f0f4f8;
  border-radius: 8px;
  // padding: 16px;
  font-size: 14px;
  border: 1px solid #dcdcdc;
  // box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.delete-info {
  padding: 10px;
  color: #333;
  text-align: center;
}
</style>
