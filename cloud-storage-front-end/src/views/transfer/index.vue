<template>
  <div class="upload-view">
    <div class="header">
      <el-button type="primary" @click="clear">清空记录</el-button>
      <el-button type="primary" @click="loadHistoryFiles">刷新记录</el-button>
      <div class="tabs">
        <span class="tab" :class="{ active: activeTab === 'uploading' }" @click="activeTab = 'uploading'">
          上传中({{ uploadingCount }})
        </span>
        <span class="tab" :class="{ active: activeTab === 'completed' }" @click="activeTab = 'completed'">
          已完成({{ completedCount }})
        </span>
      </div>
    </div>

    <div class="file-list">
      <!-- 上传中的表格 -->
      <template v-if="activeTab === 'uploading'">
        <div class="list-header uploading-header">
          <span>文件</span>
          <span>大小</span>
          <span>状态</span>
          <span>操作</span>
        </div>

        <div class="list-items">
          <div v-for="file in displayFiles" :key="file.id" class="list-item">
            <div class="item-main uploading-item">
              <span class="file-name">{{ file.name }}</span>
              <span class="file-size">{{ file.size }}</span>
              <span class="file-status">{{ file.status }}</span>
              <span class="file-action">
                <el-button type="danger" @click="deleteFile(file.id)">删除</el-button>
              </span>
            </div>
          </div>
        </div>
      </template>

      <!-- 已完成的表格 -->
      <template v-else>
        <div class="list-header uploading-item">
          <el-checkbox v-model="selectAll" @change="handleSelectAll">文件</el-checkbox>
          <span>大小</span>
          <span>完成时间</span>
        </div>

        <div class="list-items">
          <div
            v-for="file in displayFiles"
            :key="file.id"
            class="list-item"
            :class="{ 'item-selected': file.selected }"
          >
            <div class="item-main completed-item">
              <el-checkbox v-model="file.selected" />
              <span class="file-name">{{ file.name }}</span>
              <span class="file-size">{{ file.size }}</span>
              <span class="file-date">{{ file.date }}</span>
            </div>
          </div>
        </div>
      </template>
    </div>

    <!-- 添加分页 -->
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
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted } from "vue"
import { historyFileApi } from "@/api/file/history"
import { formatFileSize } from "@/utils/format/formatFileSize"
import { formatTime } from "@/utils/format/formatTime"

interface UploadFile {
  id: string
  name: string
  size: number
  date: string
  status: "uploading" | "completed"
  selected: boolean
  repository_id: number
}

// 实际数据
const uploadedFiles = ref<UploadFile[]>([])
const loading = ref(false)

// 添加分页相关
const currentPage = ref(0)
const pageSize = ref(10)
const total = ref(0)

// 加载历史记录
const loadHistoryFiles = async () => {
  loading.value = true
  try {
    const response = await historyFileApi.getHistoryFileList({
      page: currentPage.value,
      size: pageSize.value
    })

    console.log("获取历史记录列表:", response.data.history_list)
    uploadedFiles.value = response.data.history_list.map((item) => ({
      id: item.id,
      name: item.file_name,
      size: formatFileSize(item.size),
      date: formatTime(item.update_time),
      status: item.status === 1 ? "completed" : "uploading",
      selected: false,
      repository_id: item.repository_id
    }))

    total.value = response.data.total
  } catch (error) {
    console.error("获取历史记录列表失败:", error)
  } finally {
    loading.value = false
  }
}

// 添加分页处理
const handleSizeChange = (val: number) => {
  pageSize.value = val
  loadHistoryFiles()
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
  loadHistoryFiles()
}

onMounted(() => {
  loadHistoryFiles()
})

const activeTab = ref<"uploading" | "completed">("uploading")
// const uploadedFiles = ref<UploadFile[]>([])
const selectAll = ref(false)

const displayFiles = computed(() => {
  return uploadedFiles.value.filter((file) => file.status === activeTab.value)
})

const uploadingCount = computed(() => {
  return uploadedFiles.value.filter((file) => file.status === "uploading").length
})

const completedCount = computed(() => {
  return uploadedFiles.value.filter((file) => file.status === "completed").length
})

const handleSelectAll = (val: boolean) => {
  displayFiles.value.forEach((file) => {
    file.selected = val
  })
}

const clear = async () => {
  try {
    // 删除勾选的文件
    if (uploadedFiles.value.length > 0) {
      console.log("删除所有文件", uploadedFiles.value)
      const ids = uploadedFiles.value.filter((file) => file.selected).map((file) => file.id.toString())
      await historyFileApi.deleteHistoryFile({ ids })
      uploadedFiles.value = uploadedFiles.value.filter((file) => !file.selected) // 清空本地记录
      await loadHistoryFiles() // 重新加载历史记录
    }
  } catch (error) {
    console.error("删除历史记录失败:", error)
  }
}

// 添加删除单行文件的函数
const deleteFile = async (fileId: number | string) => {
  try {
    console.log("删除文件:", fileId)
    await historyFileApi.deleteHistoryFile({ ids: [fileId.toString()] })
    uploadedFiles.value = uploadedFiles.value.filter((file) => file.id !== fileId.toString()) // 从本地记录中删除
  } catch (error) {
    console.error("删除历史记录失败:", error)
  }
}
</script>

<style scoped>
.upload-view {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: flex-start;
  align-items: center;
  margin-bottom: 20px;
  gap: 0;
}

.tabs {
  display: flex;
  gap: 20px;
  margin-left: auto;
}

.tab {
  cursor: pointer;
  padding: 8px 16px;
  border-bottom: 2px solid transparent;
}

.tab.active {
  border-bottom-color: #409eff;
  color: #409eff;
}

.file-list {
  background: white;
  border-radius: 8px;
  padding: 16px;
}

.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid #ebeef5;
  color: #909399;
}

.list-items {
  margin-top: 8px;
}

.list-item {
  padding: 12px 16px;
  border-radius: 4px;
  list-style: none;
}

.list-item:hover {
  background-color: #f5f7fa;
}

.item-selected {
  background-color: #f0f7ff;
}

.item-main {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.file-info {
  /* display: flex;
  align-items: center;
  gap: 8px;
  overflow: hidden; */
}

.file-icon {
  color: #909399;
  margin-right: 8px;
}

.file-name {
  /* overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1; */
}

.file-size,
.file-date,
.file-ai,
.file-status {
  text-align: left;
}

.uploading-header {
  display: grid;
  grid-template-columns: 1fr 120px 120px 120px;
  padding: 12px 16px;
  border-bottom: 1px solid #ebeef5;
  color: #909399;
  align-items: center;
}

.completed-header {
  display: grid;
  grid-template-columns: 1fr 120px 120px 120px;
  padding: 12px 16px;
  border-bottom: 1px solid #ebeef5;
  color: #909399;
  align-items: center;
}

.uploading-item {
  display: grid;
  grid-template-columns: 1fr 120px 120px 120px;
  gap: 8px;
  align-items: center;
  width: 100%;
}

.completed-item {
  display: grid;
  grid-template-columns: 24px 1fr 120px 120px;
  gap: 8px;
  align-items: center;
}

.file-action {
  text-align: left;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
