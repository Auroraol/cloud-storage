<template>
  <div class="app-container">
    <!-- SSH连接对话框 -->
    <el-dialog v-model="sshDialogVisible" title="SSH连接" width="500px">
      <el-form :model="sshForm" label-width="100px">
        <el-form-item label="主机地址">
          <el-input v-model="sshForm.host" placeholder="请输入主机地址" />
        </el-form-item>
        <el-form-item label="端口">
          <el-input-number v-model="sshForm.port" :min="1" :max="65535" />
        </el-form-item>
        <el-form-item label="用户名">
          <el-input v-model="sshForm.user" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="sshForm.password" type="password" placeholder="请输入密码" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="sshDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSSHConnect" :loading="sshConnecting">连接</el-button>
        </span>
      </template>
    </el-dialog>

    <div class="filter-container">
      <!-- 搜索条件区域 -->
      <el-form :model="queryParams" ref="queryForm" :inline="true">
        <el-form-item label="主机" prop="host">
          <el-select
            v-model="queryParams.host"
            placeholder="请选择主机"
            clearable
            style="width: 200px"
            @change="handleHostChange"
          >
            <el-option v-for="host in hosts" :key="host" :label="host" :value="host" />
          </el-select>
          <el-button type="primary" icon="Plus" circle size="small" @click="showSSHDialog" style="margin-left: 10px" />
        </el-form-item>
        <el-form-item label="日志文件" prop="logfile">
          <el-select v-model="queryParams.logfile" placeholder="请选择日志文件" clearable style="width: 200px">
            <el-option v-for="file in logFiles" :key="file" :label="file" :value="file" />
          </el-select>
        </el-form-item>
        <el-form-item label="关键字" prop="keyword">
          <el-input v-model="queryParams.keyword" placeholder="请输入关键字" clearable style="width: 200px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleQuery">查询</el-button>
          <el-button @click="resetQuery(queryForm)">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 日志内容展示区域 -->
    <el-card class="log-content" v-loading="loading">
      <div class="log-stats">
        <el-row :gutter="20">
          <el-col :span="6">
            <div class="stat-item">
              <div class="stat-title">总行数</div>
              <div class="stat-value">{{ stats.totalLines }}</div>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="stat-item">
              <div class="stat-title">匹配行数</div>
              <div class="stat-value">{{ stats.matchLines }}</div>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="stat-item">
              <div class="stat-title">当前页</div>
              <div class="stat-value">{{ stats.currentPage }}/{{ stats.totalPages }}</div>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="stat-item">
              <div class="stat-title">刷新频率</div>
              <el-select v-model="refreshRate" style="width: 120px">
                <el-option label="5秒" value="5" />
                <el-option label="10秒" value="10" />
                <el-option label="30秒" value="30" />
                <el-option label="手动" value="0" />
              </el-select>
            </div>
          </el-col>
        </el-row>
      </div>

      <div class="log-viewer" ref="logViewer">
        <div v-for="line in logLines" :key="line.number" class="log-line" :class="{ highlight: line.highlight }">
          <span class="line-number">{{ line.number }}</span>
          <span class="line-content" v-html="line.content" />
        </div>
      </div>

      <div class="log-actions">
        <el-button-group>
          <el-button type="primary" @click="loadPrevPage" :disabled="isFirstPage">上一页</el-button>
          <el-button type="primary" @click="loadNextPage" :disabled="isLastPage">下一页</el-button>
        </el-button-group>
        <el-button type="success" @click="handleQuery">刷新</el-button>
        <el-button type="warning" @click="handleDownload">下载</el-button>
      </div>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted, onUnmounted, computed, watch } from "vue"
import { ElMessage } from "element-plus"
import type { FormInstance } from "element-plus"
import { readLogApi, downloadLogApi, getHostsApi, getLogFilesApi, connectSSHApi } from "@/api/log/frontend"
import type { LogQueryParams, LogLine, LogStats } from "@/api/log/types/frontend"
import { useUserStoreHook } from "@/store/modules/user"

// 用户存储
const userStore = useUserStoreHook()

// SSH连接表单
const sshDialogVisible = ref(false)
const sshConnecting = ref(false)
const sshForm = reactive({
  host: "",
  port: 22,
  user: "root",
  password: ""
})

// 显示SSH连接对话框
const showSSHDialog = () => {
  sshDialogVisible.value = true
}

// 处理SSH连接
const handleSSHConnect = async () => {
  if (!sshForm.host) {
    ElMessage.warning("请输入主机地址")
    return
  }

  try {
    sshConnecting.value = true
    await connectSSHApi(sshForm.host, sshForm.port, sshForm.user, sshForm.password)
    ElMessage.success("SSH连接成功")

    // 刷新主机列表
    await refreshHosts()

    // 设置为当前选中的主机
    queryParams.host = sshForm.host

    // 获取该主机的日志文件列表
    await handleHostChange(sshForm.host)

    sshDialogVisible.value = false
  } catch (error) {
    ElMessage.error(`SSH连接失败: ${error instanceof Error ? error.message : String(error)}`)
  } finally {
    sshConnecting.value = false
  }
}

// 查询参数
const queryParams = reactive<LogQueryParams>({
  logfile: "",
  host: "",
  keyword: "",
  page: 1,
  pageSize: 100
})

// 状态数据
const loading = ref(false)
const logLines = ref<LogLine[]>([])
const stats = ref<LogStats>({
  totalLines: 0,
  matchLines: 0,
  currentPage: 1,
  totalPages: 1
})
const refreshRate = ref("0")
const hosts = ref<string[]>([])
const logFiles = ref<string[]>([])
const queryForm = ref<FormInstance>()

// 计算属性
const isFirstPage = computed(() => stats.value.currentPage <= 1)
const isLastPage = computed(() => stats.value.currentPage >= stats.value.totalPages)

// 定时器
let refreshTimer: number | null = null

// 查询日志
const handleQuery = async () => {
  if (!queryParams.host) {
    ElMessage.warning("请选择主机")
    return
  }

  if (!queryParams.logfile) {
    ElMessage.warning("请选择日志文件")
    return
  }

  try {
    loading.value = true
    const res = await readLogApi(queryParams)
    logLines.value = res.lines
    stats.value = res.stats
  } catch (error) {
    ElMessage.error("获取日志失败")
    console.error(error)
  } finally {
    loading.value = false
  }
}

// 重置查询
const resetQuery = (formEl: FormInstance | undefined) => {
  if (!formEl) return
  formEl.resetFields()
  queryParams.page = 1
  handleQuery()
}

// 刷新主机列表
const refreshHosts = async () => {
  try {
    hosts.value = await getHostsApi()
  } catch (error) {
    console.error("获取主机列表失败:", error)
  }
}

// 主机变更
const handleHostChange = async (host: string) => {
  if (!host) {
    logFiles.value = []
    return
  }

  try {
    // 设置当前SSH主机
    userStore.setCurrentSSHHost(host)

    // 获取日志文件列表
    logFiles.value = await getLogFilesApi(host, "")
  } catch (error) {
    ElMessage.error("获取日志文件列表失败")
    console.error(error)
  }
}

// 翻页操作
const loadPrevPage = () => {
  if (!isFirstPage.value) {
    queryParams.page--
    handleQuery()
  }
}

const loadNextPage = () => {
  if (!isLastPage.value) {
    queryParams.page++
    handleQuery()
  }
}

// 下载日志
const handleDownload = async () => {
  if (!queryParams.host || !queryParams.logfile) {
    ElMessage.warning("请选择主机和日志文件")
    return
  }

  try {
    loading.value = true
    const blob = await downloadLogApi(queryParams)
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement("a")
    link.href = url
    link.download = `log-${new Date().getTime()}.txt`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
  } catch (error) {
    ElMessage.error("下载日志失败")
    console.error(error)
  } finally {
    loading.value = false
  }
}

// 自动刷新
const startAutoRefresh = () => {
  stopAutoRefresh()
  if (refreshRate.value !== "0") {
    refreshTimer = window.setInterval(
      () => {
        handleQuery()
      },
      parseInt(refreshRate.value) * 1000
    )
  }
}

const stopAutoRefresh = () => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
}

// 监听刷新频率变化
watch(refreshRate, () => {
  startAutoRefresh()
})

// 初始化数据
const initData = async () => {
  await refreshHosts()

  // 如果有当前SSH主机，则自动选择
  if (userStore.currentSSHHost) {
    queryParams.host = userStore.currentSSHHost
    await handleHostChange(userStore.currentSSHHost)
  }
}

onMounted(() => {
  initData()
})

onUnmounted(() => {
  stopAutoRefresh()
})
</script>

<style scoped>
.filter-container {
  margin-bottom: 20px;
}

.log-content {
  background: #fff;
}

.log-stats {
  padding: 15px;
  border-bottom: 1px solid #eee;
}

.stat-item {
  text-align: center;
}

.stat-title {
  color: #666;
  font-size: 14px;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #409eff;
}

.log-viewer {
  height: 600px;
  overflow-y: auto;
  padding: 10px;
  font-family: monospace;
  background: #f8f9fa;
}

.log-line {
  padding: 2px 5px;
  border-bottom: 1px solid #eee;
  white-space: pre-wrap;
}

.log-line:hover {
  background: #f0f0f0;
}

.line-number {
  color: #999;
  padding-right: 10px;
  user-select: none;
}

.line-content {
  color: #333;
}

.highlight {
  background-color: #ff980020;
}

.log-actions {
  padding: 15px;
  text-align: right;
  border-top: 1px solid #eee;
}
</style>
