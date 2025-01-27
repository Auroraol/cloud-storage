<template>
  <div class="app-container">
    <el-card v-loading="loading">
      <div class="toolbar-wrapper">
        <el-button type="danger" @click="handleBatchDelete" :disabled="selectedRows.length === 0">批量删除</el-button>
      </div>
      <div class="table-wrapper">
        <el-table
          ref="multipleTable"
          :data="tableData"
          v-model="selectedRows"
          @selection-change="handleSelectionChange"
          @row-click="handleRowClick"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column prop="filename" label="文件">
            <template #default="{ row }">
              <div style="display: flex; align-items: center">
                <Icon
                  v-if="(row.fileType === 3 || row.fileType === 1) && row.status === 2"
                  :cover="row.fileCover"
                  :width="32"
                />
                <Icon v-else-if="row.folderType === 0" :file-type="row.fileType" />
                <Icon v-else icon-name="folder_2" :width="40" />
                <span class="filename" style="margin-left: 8px">{{ row.filename }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="deletedAt" label="删除时间" />
          <el-table-column fixed="right" label="操作" width="150">
            <template #default="{ row }">
              <el-button type="danger" @click="handleDelete(row)">删除</el-button>
            </template>
          </el-table-column>
          <template #empty>
            <div style="text-align: center; margin-top: 20px">
              <svg-icon name="trash" color="green" width="100px" height="80px" />
              <div style="margin-top: -30px; font-size: 16px; font-weight: bold">您的回收站为空 ~</div>
            </div>
          </template>
        </el-table>
      </div>
      <el-pagination
        :total="paginationData.total"
        :page-size="paginationData.pageSize"
        :current-page="paginationData.currentPage"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import { ref, watch } from "vue"
import { ElMessage, ElMessageBox } from "element-plus"
import { usePagination } from "@/hooks/usePagination"
import Icon from "@/components/FileIcon/Icon.vue"
import SvgIcon from "@/components/SvgIcon/index.vue"
// 定义模拟数据
const mockData = [
  // { id: "1", filename: "测试文件1.txt", deletedAt: "2023-01-01 10:00:00" },
  // { id: "2", filename: "测试文件2.jpg", deletedAt: "2023-01-02 11:00:00" },
  // { id: "3", filename: "测试文件3.pdf", deletedAt: "2023-01-03 12:00:00" }
]

defineOptions({
  name: "Recycle"
})

const loading = ref<boolean>(false)
const { paginationData, handleCurrentChange, handleSizeChange } = usePagination()

const tableData = ref(mockData)
const selectedRows = ref([])

const handleBatchDelete = () => {
  console.log(selectedRows.value)
  if (selectedRows.value.length === 0) {
    ElMessage.warning("请先选择要删除的项")
    return
  }
  ElMessageBox.confirm(`确认删除选中的 ${selectedRows.value.length} 项？`, "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning"
  }).then(() => {
    selectedRows.value.forEach((row) => {
      tableData.value = tableData.value.filter((item) => item.id !== row.id)
    })
    selectedRows.value = []
    ElMessage.success("删除成功")
  })
}

const handleDelete = (row) => {
  ElMessageBox.confirm(`确认删除该项：${row.filename}？`, "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning"
  }).then(() => {
    ElMessage.success("删除成功")
    tableData.value = tableData.value.filter((item) => item.id !== row.id)
  })
}

const handleSelectionChange = (val) => {
  selectedRows.value = val
}

const multipleTable = ref()
const handleRowClick = (row) => {
  selectedRows.value = row
  //通过ref绑定来操作bom元素
  multipleTable.value.toggleRowSelection(row)
}

watch(
  [() => paginationData.currentPage, () => paginationData.pageSize],
  () => {
    // 这里可以添加获取数据的逻辑
  },
  { immediate: true }
)
</script>

<style scoped>
.toolbar-wrapper {
  margin-bottom: 20px;
}
.table-wrapper {
  margin-bottom: 20px;
}
</style>
