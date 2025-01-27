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

interface ShareItem {
  filename: string
  fileType: number
  status: number
  fileId: string // 文件id
  userId: string // 用户id
  validType: boolean // 有效期类型 0:1天 1:7天 2:30天 3:永久有效
  code: string // 提取码
  browseCount: number // 浏览次数
  saveCount: number // 收藏次数
  downloadCount: number // 下载次数
  expireTime: Date // 过期时间
  createTime: Date // 更新时间
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

// 模拟数据
const fileList = ref<ShareItem[]>([
  {
    filename: "文件1",
    fileType: 0,
    status: 2,
    fileId: "1",
    userId: "2",
    validType: true,
    code: "ABC123",
    browseCount: 10,
    saveCount: 5,
    downloadCount: 2,
    expireTime: new Date("2024-12-31"),
    createTime: new Date("2023-01-01")
  },
  {
    filename: "文件2",
    fileType: 9,
    status: 1,
    fileId: "2",
    userId: "1",
    validType: false,
    code: "XYZ789",
    browseCount: 20,
    saveCount: 10,
    downloadCount: 5,
    expireTime: new Date("2023-06-30"),
    createTime: new Date("2023-02-01")
  }
])

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

const deleteShare = () => {
  ElMessage.success("取消外链分享成功")
  dialogConfig.show = false
  reloadTable()
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
  if (loading.value || noMore.value) return // 防止重复加载
  loading.value = true

  // 模拟 API 请求
  const response = await fetch(`/api/files?page=${page.value}&size=${pageSize.value}`)
  const data = await response.json()

  if (data && data.length) {
    fileList.value.push(...data)
    tableDataNum.value += data.length
    page.value += 1 // 增加页码
  } else {
    noMore.value = true // 没有更多数据
    ElMessage.warning("没有更多文件了")
  }
  loading.value = false
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
