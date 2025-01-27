<template>
  <div class="audit-log">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="操作类型">
          <el-select v-model="searchForm.operationType" clearable placeholder="请选择操作类型">
            <el-option label="上传" value="upload" />
            <el-option label="下载" value="download" />
            <el-option label="删除" value="delete" />
            <el-option label="修改" value="modify" />
          </el-select>
        </el-form-item>
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="searchForm.timeRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="resetForm">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table :data="logList" border style="width: 100%">
      <el-table-column prop="operationType" label="操作类型">
        <template #default="{ row }">
          <el-tag :type="getOperationTypeTag(row.operationType)">
            {{ getOperationTypeLabel(row.operationType) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="fileName" label="文件名" />
      <el-table-column prop="fileSize" label="文件大小">
        <template #default="{ row }">
          {{ formatFileSize(row.fileSize) }}
        </template>
      </el-table-column>
      <el-table-column prop="operationTime" label="操作时间">
        <template #default="{ row }">
          {{ formatDateTime(row.operationTime) }}
        </template>
      </el-table-column>
      <el-table-column prop="ip" label="IP地址" />
      <el-table-column label="操作" width="120">
        <template #default="{ row }">
          <el-button type="text" @click="viewDetail(row)">查看详情</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      :current-page="page"
      :page-size="pageSize"
      :total="total"
      layout="total, prev, pager, next"
      @current-change="handlePageChange"
    />

    <!-- 详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="操作详情" width="500px">
      <div v-if="currentDetail">
        <p><strong>文件名：</strong>{{ currentDetail.fileName }}</p>
        <p><strong>操作类型：</strong>{{ getOperationTypeLabel(currentDetail.operationType) }}</p>
        <p><strong>操作时间：</strong>{{ formatDateTime(currentDetail.operationTime) }}</p>
        <p><strong>文件大小：</strong>{{ formatFileSize(currentDetail.fileSize) }}</p>
        <p><strong>IP地址：</strong>{{ currentDetail.ip }}</p>
        <p><strong>操作结果：</strong>{{ currentDetail.result ? "成功" : "失败" }}</p>
        <p v-if="currentDetail.errorMessage"><strong>错误信息：</strong>{{ currentDetail.errorMessage }}</p>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from "vue"
import { ElMessage } from "element-plus"
import { getAuditLogsApi } from "@/api/audit"
import type { LogRecord } from "@/api/audit/types"
// import { useUserStore } from "@/stores/user"

interface LogRecord {
  id: string
  operationType: "upload" | "download" | "delete" | "modify"
  fileName: string
  fileSize: number
  operationTime: string
  ip: string
  result: boolean
  errorMessage?: string
}

// const userStore = useUserStore()

const searchForm = reactive({
  operationType: "",
  timeRange: [] as Date[]
})

const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const logList = ref<LogRecord[]>([])
const detailDialogVisible = ref(false)
const currentDetail = ref<LogRecord | null>(null)

const handleSearch = async () => {
  try {
    const res = await getAuditLogsApi({
      userId: userStore.userId,
      operationType: searchForm.operationType,
      startTime: searchForm.timeRange[0]?.toISOString(),
      endTime: searchForm.timeRange[1]?.toISOString(),
      page: page.value,
      pageSize: pageSize.value
    })
    logList.value = res.list
    total.value = res.total
  } catch (error) {
    ElMessage.error("获取日志列表失败")
    console.error("获取日志列表失败:", error)
  }
}

const resetForm = () => {
  searchForm.operationType = ""
  searchForm.timeRange = []
  handleSearch()
}

const handlePageChange = (newPage: number) => {
  page.value = newPage
  handleSearch()
}

const viewDetail = (row: LogRecord) => {
  currentDetail.value = row
  detailDialogVisible.value = true
}

const getOperationTypeLabel = (type: string) => {
  const types = {
    upload: "上传",
    download: "下载",
    delete: "删除",
    modify: "修改"
  }
  return types[type] || type
}

const getOperationTypeTag = (type: string) => {
  const types = {
    upload: "success",
    download: "info",
    delete: "danger",
    modify: "warning"
  }
  return types[type] || ""
}

const formatFileSize = (size: number) => {
  if (size < 1024) return size + " B"
  if (size < 1024 * 1024) return (size / 1024).toFixed(2) + " KB"
  if (size < 1024 * 1024 * 1024) return (size / 1024 / 1024).toFixed(2) + " MB"
  return (size / 1024 / 1024 / 1024).toFixed(2) + " GB"
}

const formatDateTime = (time: string) => {
  return new Date(time).toLocaleString()
}
</script>

<style scoped>
.audit-log {
  padding: 20px;
}

.search-bar {
  margin-bottom: 20px;
}

.el-tag {
  width: 60px;
  text-align: center;
}
</style>
