<template>
  <div class="audit-log">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="操作类型">
          <el-select v-model="searchForm.operationType" clearable placeholder="请选择操作类型" style="width: 160px">
            <el-option
              v-for="option in operationTypeOptions"
              :key="option.value"
              :label="option.label"
              :value="option.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="searchForm.timeRange"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            :shortcuts="dateShortcuts"
            :editable="false"
            :clearable="true"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
            :popper-options="{
              modifiers: [{ name: 'arrow', options: { padding: 10 } }]
            }"
            popper-class="custom-date-picker"
            :teleported="true"
            confirm
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
          <el-tag :type="getOperationTypeTag(row.flag)">
            {{ getOperationTypeLabel(row.flag) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="fileName" label="文件名">
        <template #default="{ row }">
          {{ row.file_name }}
        </template>
      </el-table-column>
      <el-table-column prop="operationTime" label="操作时间">
        <template #default="{ row }">
          {{ timestampToDate(row.created_at) }}
        </template>
      </el-table-column>
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
        <p><strong>文件名：</strong>{{ currentDetail.file_name }}</p>
        <p><strong>文件ID / 所在文件夹ID：</strong>{{ currentDetail.file_id }}</p>
        <p><strong>操作类型：</strong>{{ getOperationTypeLabel(currentDetail.flag) }}</p>
        <p><strong>操作时间：</strong>{{ timestampToDate(Number(currentDetail.created_at)) }}</p>
        <p><strong>操作内容：</strong>{{ currentDetail.content }}</p>
        <p v-if="currentDetail.errorMessage"><strong>错误信息：</strong>{{ currentDetail.errorMessage }}</p>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from "vue"
import { ElMessage } from "element-plus"
import { auditApi } from "@/api/audit"
import { timestampToDate, dateToTimestamp } from "@/utils/format/formatTime"
interface LogRecord {
  content: string
  flag: number
  file_size: number
  errorMessage?: string
  created_at: string
  file_name: string
  file_id: string
}

const searchForm = reactive({
  operationType: -1,
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
    console.log("searchForm", searchForm.timeRange)
    console.log("searchForm.timeRange[0]", searchForm.timeRange[0])
    console.log("searchForm.timeRange[1]", searchForm.timeRange[1])
    console.log(
      "dateToTimestamp(searchForm.timeRange[0]?.toISOString())",
      dateToTimestamp(searchForm.timeRange[0]?.toISOString())
    )
    const res = await auditApi.getAuditList({
      flag: Number(searchForm.operationType),
      start_time: dateToTimestamp(searchForm.timeRange[0]?.toISOString()) || 0,
      end_time: dateToTimestamp(searchForm.timeRange[1]?.toISOString()) || 0,
      page: page.value,
      page_size: pageSize.value
    })

    logList.value = res.data.operation_logs.map((item) => ({
      file_name: item.file_name,
      file_id: item.file_id,
      content: item.content,
      flag: item.flag,
      file_size: item.file_size,
      created_at: item.created_at
    }))
    total.value = res.data.total
  } catch (error) {
    ElMessage.error("获取日志列表失败")
    console.error("获取日志列表失败:", error)
  }
}

const resetForm = () => {
  searchForm.operationType = -1
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

//操作类型，0：上传，1：下载，2：删除，3.恢复 4：重命名，5：移动，6：复制，7：创建文件夹，8：修改文件, -1: 全部
const getOperationTypeLabel = (type: number) => {
  // 定义操作类型与对应标签的映射关系
  const types = {
    0: "上传",
    1: "下载",
    2: "删除",
    3: "恢复",
    4: "重命名",
    5: "移动",
    6: "复制",
    7: "创建文件夹",
    8: "修改文件"
  }
  return types[type] || type
}

const getOperationTypeTag = (type: number) => {
  // 定义操作类型与对应标签的映射关系
  const types = {
    0: "success",
    1: "info",
    2: "danger",
    3: "warning",
    4: "warning",
    5: "warning",
    6: "warning",
    7: "warning",
    8: "warning"
  }
  return types[type] || ""
}

// 添加日期快捷选项
const dateShortcuts = [
  {
    text: "最近一周",
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
      return [start, end]
    }
  },
  {
    text: "最近一个月",
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 30)
      return [start, end]
    }
  },
  {
    text: "最近三个月",
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 90)
      return [start, end]
    }
  }
]

const operationTypeOptions = [
  { label: "全部", value: -1 },
  { label: "上传", value: 0 },
  { label: "下载", value: 1 },
  { label: "删除", value: 2 },
  { label: "修改", value: 8 },
  { label: "恢复", value: 3 },
  { label: "重命名", value: 4 },
  { label: "移动", value: 5 },
  { label: "复制", value: 6 },
  { label: "创建文件夹", value: 7 }
]

onMounted(() => {
  handleSearch()
})
</script>

<style scoped>
.audit-log {
  padding: 20px;
}

.search-bar {
  margin-bottom: 10px;
}

.el-tag {
  width: 60px;
  text-align: center;
}

/* 新增样式 */
.el-select .el-option {
  font-size: 14px; /* 调整字体大小 */
}

/* 自定义日期选择器样式 */
:deep(.el-date-editor) {
  width: 280px;
}

:deep(.custom-date-picker) {
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

:deep(.custom-date-picker .el-picker-panel__footer) {
  text-align: right;
  padding: 8px 12px;
  border-top: 1px solid #e4e7ed;
}

:deep(.custom-date-picker .el-picker-panel__footer .el-button) {
  padding: 8px 20px;
  font-size: 14px;
  border-radius: 4px;
}

:deep(.custom-date-picker .el-button--default) {
  margin-right: 8px;
}

:deep(.custom-date-picker .el-button--primary) {
  background-color: #409eff;
  border-color: #409eff;
}

:deep(.custom-date-picker .el-date-picker__time-header) {
  padding: 8px 12px;
  border-bottom: 1px solid #e4e7ed;
}

:deep(.custom-date-picker .el-date-picker__header) {
  margin: 8px 12px;
}

:deep(.custom-date-picker .el-picker-panel__content) {
  margin: 0 12px;
}

:deep(.custom-date-picker .el-time-spinner__wrapper) {
  width: 100%;
}

:deep(.custom-date-picker .el-date-table th) {
  padding: 8px 0;
  color: #606266;
}

:deep(.custom-date-picker .el-date-table td) {
  padding: 4px 0;
}

:deep(.custom-date-picker .el-time-panel__content) {
  padding: 8px 12px;
}

:deep(.custom-date-picker .el-time-panel__footer) {
  padding: 8px 12px;
  border-top: 1px solid #e4e7ed;
}

:deep(.custom-date-picker .el-picker-panel__shortcut) {
  padding: 8px 12px;
  font-size: 14px;
}
</style>
