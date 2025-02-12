<template>
  <div class="upload-view">
    <div class="header">
      <el-button type="primary" @click="clearAll">清空全部记录</el-button>
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
        </div>

        <div class="list-items">
          <div v-for="file in displayFiles" :key="file.id" class="list-item">
            <div class="item-main uploading-item">
              <span class="file-name">{{ file.name }}</span>
              <span class="file-size">{{ file.size }}</span>
              <span class="file-status">{{ file.status }}</span>
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
              <div class="file-info">
                <span class="file-name">{{ file.name }}</span>
              </div>
              <span class="file-size">{{ file.size }}</span>
              <span class="file-date">{{ file.date }}</span>
            </div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed } from "vue"

interface UploadFile {
  id: number
  name: string
  size: string
  date: string
  status: "uploading" | "completed"
  selected: boolean
}

// 模拟数据
const uploadedFiles = ref<UploadFile[]>([
  {
    id: 1,
    name: "文件1.txt",
    size: "150 KB",
    date: "2023-10-01",
    status: "uploading",
    selected: false
  },
  {
    id: 2,
    name: "文件2.jpg",
    size: "2.5 MB",
    date: "2023-10-02",
    status: "completed",
    selected: false
  },
  {
    id: 3,
    name: "文件3.pdf",
    size: "1.2 MB",
    date: "2023-10-03",
    status: "uploading",
    selected: false
  },
  {
    id: 4,
    name: "文件4.docx",
    size: "300 KB",
    date: "2023-10-04",
    status: "completed",
    selected: false
  }
])

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

const clearAll = () => {
  uploadedFiles.value = []
}

const handleFileUpload = (event: Event) => {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (file) {
    const newFile: UploadFile = {
      id: uploadedFiles.value.length + 1,
      name: file.name,
      size: `${(file.size / 1024).toFixed(2)} KB`,
      date: new Date().toLocaleDateString(),
      status: "completed",
      selected: false,
      aiStatus: "未编辑"
    }
    uploadedFiles.value.push(newFile)
  }
}
</script>

<style scoped>
.upload-view {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.tabs {
  display: flex;
  gap: 20px;
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
  display: grid;
  grid-template-columns: 1fr 120px 120px 120px;
  padding: 12px 16px;
  border-bottom: 1px solid #ebeef5;
  color: #909399;
  align-items: center;
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
  display: grid;
  grid-template-columns: 24px 1fr 120px 120px 120px;
  gap: 8px;
  align-items: center;
}

.file-info {
  display: flex;
  align-items: center;
  gap: 8px;
  overflow: hidden;
}

.file-icon {
  color: #909399;
  margin-right: 8px;
}

.file-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
}

.file-size,
.file-date,
.file-ai,
.file-status {
  text-align: left;
}

.uploading-header {
  display: grid;
  grid-template-columns: 1fr 120px 120px;
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
  grid-template-columns: 1fr 120px 120px;
  gap: 8px;
  align-items: center;
}

.completed-item {
  display: grid;
  grid-template-columns: 24px 1fr 120px 120px;
  gap: 8px;
  align-items: center;
}
</style>
