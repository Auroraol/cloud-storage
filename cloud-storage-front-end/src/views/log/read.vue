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
        <!-- 添加模式切换 -->
        <el-form-item>
          <el-radio-group v-model="mode" @change="handleModeChange">
            <el-radio-button label="ssh">SSH模式</el-radio-button>
            <el-radio-button label="local">本地模式</el-radio-button>
          </el-radio-group>
        </el-form-item>

        <!-- SSH模式下的主机选择 -->
        <el-form-item label="主机" prop="host" v-if="mode === 'ssh'">
          <el-select
            v-model="queryParams.host"
            placeholder="请选择主机"
            clearable
            style="width: 200px"
            @change="handleHostChange"
          >
            <el-option
              v-for="host in hosts"
              :key="host"
              :label="host"
              :value="host"
              @dblclick="handleHostDblClick(host)"
            />
          </el-select>
          <el-button type="primary" icon="Plus" circle size="small" @click="showSSHDialog" style="margin-left: 10px" />
        </el-form-item>

        <!-- SSH模式下的路径选择 -->
        <el-form-item label="路径" prop="path" v-if="mode === 'ssh'">
          <el-input
            v-model="queryParams.path"
            placeholder="请输入日志文件路径"
            clearable
            style="width: 200px"
            @change="handlePathChange"
          />
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
      <div class="log-toolbar">
        <div class="log-stats-wrapper">
          <div class="stat-item">
            <el-icon><Document /></el-icon>
            <span class="stat-label">总行数:</span>
            <span class="stat-value">{{ stats.totalLines }}</span>
          </div>
          <div class="stat-item">
            <el-icon><Search /></el-icon>
            <span class="stat-label">匹配行数:</span>
            <span class="stat-value">{{ stats.matchLines }}</span>
          </div>
          <div class="stat-item">
            <el-icon><DocumentCopy /></el-icon>
            <span class="stat-label">页码:</span>
            <span class="stat-value">{{ stats.currentPage }}/{{ stats.totalPages }}</span>
          </div>
        </div>
        <div class="log-actions">
          <el-tooltip content="刷新日志" placement="top">
            <el-button type="primary" circle @click="handleQuery">
              <el-icon><Refresh /></el-icon>
            </el-button>
          </el-tooltip>
          <el-tooltip content="下载日志" placement="top">
            <el-button type="success" circle @click="handleDownload">
              <el-icon><Download /></el-icon>
            </el-button>
          </el-tooltip>
          <el-select v-model="refreshRate" placeholder="刷新频率" class="refresh-select">
            <el-option label="5秒刷新" value="5" />
            <el-option label="10秒刷新" value="10" />
            <el-option label="30秒刷新" value="30" />
            <el-option label="手动刷新" value="0" />
          </el-select>
        </div>
      </div>

      <div class="log-viewer" ref="logViewer">
        <div v-for="line in logLines" :key="line.number" class="log-line" :class="getLogLevelClass(line.content)">
          <div class="line-meta">
            <span class="line-number">{{ line.number }}</span>
            <span class="line-level" :class="getLogLevelClass(line.content)">{{ getLogLevel(line.content) }}</span>
          </div>
          <span class="line-content" v-html="highlightContent(line.content)" />
        </div>
        <div v-if="logLines.length === 0" class="empty-logs">
          <el-empty description="暂无日志数据" />
        </div>
      </div>

      <!-- 分页区域 -->
      <div class="pagination-container">
        <div class="pagination-wrapper">
          <!-- 分页信息 -->
          <div class="pagination-info">
            <span
              >共 <strong>{{ stats.totalLines }}</strong> 条记录</span
            >
          </div>

          <!-- 分页控制 -->
          <div class="pagination-controls">
            <!-- 首页按钮 -->
            <button
              class="pagination-btn"
              :class="{ disabled: queryParams.page <= 1 }"
              @click="queryParams.page > 1 && jumpToPage(1)"
            >
              <i class="pagination-icon">«</i>
            </button>

            <!-- 上一页按钮 -->
            <button
              class="pagination-btn"
              :class="{ disabled: queryParams.page <= 1 }"
              @click="queryParams.page > 1 && jumpToPage(queryParams.page - 1)"
            >
              <i class="pagination-icon">‹</i>
            </button>

            <!-- 页码按钮 -->
            <div class="pagination-pages">
              <template v-for="page in displayedPages" :key="page">
                <button
                  v-if="page !== '...'"
                  class="pagination-btn page-num"
                  :class="{ active: queryParams.page === page }"
                  @click="page !== queryParams.page && jumpToPage(Number(page))"
                >
                  {{ page }}
                </button>
                <span v-else class="pagination-ellipsis">...</span>
              </template>
            </div>

            <!-- 下一页按钮 -->
            <button
              class="pagination-btn"
              :class="{ disabled: queryParams.page >= stats.totalPages }"
              @click="queryParams.page < stats.totalPages && jumpToPage(queryParams.page + 1)"
            >
              <i class="pagination-icon">›</i>
            </button>

            <!-- 末页按钮 -->
            <button
              class="pagination-btn"
              :class="{ disabled: queryParams.page >= stats.totalPages }"
              @click="queryParams.page < stats.totalPages && jumpToPage(stats.totalPages)"
            >
              <i class="pagination-icon">»</i>
            </button>
          </div>

          <!-- 跳转到页 -->
          <div class="pagination-goto">
            <span>前往</span>
            <div class="goto-input-wrapper">
              <input
                type="number"
                class="goto-input"
                v-model="currentPageInput"
                min="1"
                :max="stats.totalPages"
                @keyup.enter="handlePageInputChange(currentPageInput)"
              />
            </div>
            <span>页</span>
            <button class="goto-btn" @click="handlePageInputChange(currentPageInput)">确定</button>
          </div>

          <!-- 每页条数 -->
          <div class="page-size-selector">
            <select class="page-size-select" v-model="queryParams.pageSize" @change="(e) => handleSizeChange()">
              <option :value="50">50条/页</option>
              <option :value="100">100条/页</option>
              <option :value="200">200条/页</option>
              <option :value="500">500条/页</option>
            </select>
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted, onUnmounted, watch, computed } from "vue"
import { ElMessage } from "element-plus"
import type { FormInstance } from "element-plus"
import {
  readLogApi,
  downloadLogApi,
  getHostsApi,
  getLogFilesApi,
  getSSHConnectionsApi,
  getLocalLogFilesApi,
  readLocalLogFileApi,
  connectSSHApi
} from "@/api/log/frontend"
import type { LogQueryParams, LogLine, LogStats } from "@/api/log/types/frontend"
import { useUserStore } from "@/store/modules/user"
import { Document, Search, DocumentCopy, Refresh, Download } from "@element-plus/icons-vue"

// 用户存储
const userStore = useUserStore()

// SSH连接表单
const sshDialogVisible = ref(false)
const sshConnecting = ref(false)
const sshForm = reactive({
  host: "",
  port: 22,
  user: "root",
  password: ""
})

// 添加模式状态
const mode = ref<"ssh" | "local">("local")

// 处理模式切换
const handleModeChange = async () => {
  // 重置查询参数
  queryParams.host = ""
  queryParams.path = ""
  queryParams.logfile = ""
  logFiles.value = []
  logLines.value = []
  stats.value = {
    totalLines: 0,
    matchLines: 0,
    currentPage: 1,
    totalPages: 1
  }

  // 根据新模式初始化数据
  await initData()
}

// 获取本地日志文件列表
const getLocalLogFiles = async () => {
  try {
    const response = await getLocalLogFilesApi("")
    console.log("response", response.data.files)
    if (!response || !response.data.files) {
      throw new Error("获取本地日志文件列表失败：返回数据格式错误")
    }
    logFiles.value = response.data.files.filter((file) => !file.isDir).map((file) => file.path)
  } catch (error) {
    console.error("获取本地日志文件列表失败:", error)
    logFiles.value = []
  }
}

// 查询参数
const queryParams = reactive<LogQueryParams>({
  logfile: "",
  host: "",
  keyword: "",
  path: "",
  page: 1,
  pageSize: 50
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

// 定时器
let refreshTimer: number | null = null

// 当前页输入框的值
const currentPageInput = ref(1)

// 查询日志
const handleQuery = async () => {
  if (mode.value === "ssh") {
    await handleSSHQuery()
  } else {
    await handleLocalQuery()
  }
}

// SSH模式查询
const handleSSHQuery = async () => {
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
    const response = await readLogApi(queryParams)
    logLines.value = response.lines
    stats.value = response.stats
  } catch (error) {
    ElMessage.error("获取日志失败")
    console.error(error)
  } finally {
    loading.value = false
  }
}

// 本地模式查询
const handleLocalQuery = async () => {
  if (!queryParams.logfile) {
    ElMessage.warning("请选择日志文件")
    return
  }

  try {
    loading.value = true
    const response = await readLocalLogFileApi({
      path: queryParams.logfile,
      keyword: queryParams.keyword,
      maxResults: queryParams.pageSize
    })

    console.log("本地日志响应:", response)
    console.log("响应数据类型:", typeof response)
    console.log("响应数据结构:", Object.keys(response))
    console.log("响应数据内容:", JSON.stringify(response, null, 2))

    // 转换本地日志响应格式
    if (response && response.data) {
      const logData = response.data
      console.log("日志数据:", logData)

      if (Array.isArray(logData)) {
        logLines.value = logData.map((line, index) => ({
          number: (queryParams.page - 1) * queryParams.pageSize + index + 1,
          content: line,
          highlight: queryParams.keyword ? line.includes(queryParams.keyword) : false
        }))
        stats.value = {
          totalLines: logData.length,
          matchLines: logLines.value.filter((line) => line.highlight).length,
          currentPage: queryParams.page,
          totalPages: Math.ceil(logData.length / queryParams.pageSize)
        }
      } else if (typeof logData === "object" && Array.isArray(logData.entries)) {
        // 处理 LogEntry 数组格式的响应
        logLines.value = logData.entries.map((entry, index) => ({
          number: (queryParams.page - 1) * queryParams.pageSize + index + 1,
          content: entry.content,
          highlight: queryParams.keyword ? entry.content.includes(queryParams.keyword) : false
        }))
        stats.value = {
          totalLines: logData.total || logData.entries.length,
          matchLines: logLines.value.filter((line) => line.highlight).length,
          currentPage: queryParams.page,
          totalPages: Math.ceil((logData.total || logData.entries.length) / queryParams.pageSize)
        }
      } else {
        throw new Error(`返回数据格式错误: ${typeof logData}`)
      }
    } else {
      throw new Error("返回数据格式错误: 缺少 data 字段")
    }
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
    // 获取主机列表
    const hostList = await getHostsApi()

    // 去重处理
    hosts.value = [...new Set(hostList)]
  } catch (error) {
    console.error("获取主机列表失败:", error)
  }
}

// 处理主机变更
const handleHostChange = async (host: string) => {
  if (!host) {
    logFiles.value = []
    return
  }

  try {
    // 设置当前SSH主机
    userStore.setCurrentSSHHost(host)

    // 获取日志文件列表
    logFiles.value = await getLogFilesApi(host, queryParams.path)
  } catch (error) {
    console.log("handleHostChange error", error)
    ElMessage.error("获取日志文件列表失败")
    console.error(error)
  }
}

// 处理路径变更
const handlePathChange = async (path: string) => {
  if (!queryParams.host) {
    ElMessage.warning("请先选择主机")
    return
  }
  try {
    logFiles.value = await getLogFilesApi(queryParams.host, path)
  } catch (error) {
    console.log("handleHostChange error", error)
    ElMessage.error("获取日志文件列表失败")
    console.error(error)
  }
}

// 跳转到指定页
const jumpToPage = (page: number) => {
  if (page < 1) page = 1
  if (page > stats.value.totalPages) page = stats.value.totalPages

  queryParams.page = page
  currentPageInput.value = page
  handleQuery()
}

// 计算要显示的页码
const displayedPages = computed(() => {
  const totalPages = stats.value.totalPages
  const currentPage = queryParams.page
  const pages: (number | string)[] = []

  if (totalPages <= 7) {
    // 总页数少于7，显示所有页码
    for (let i = 1; i <= totalPages; i++) {
      pages.push(i)
    }
  } else {
    // 总页数大于7，显示部分页码
    pages.push(1) // 始终显示第一页

    if (currentPage > 4) {
      pages.push("...") // 当前页大于4，显示前省略号
    }

    // 显示当前页附近的页码
    const startPage = Math.max(2, currentPage - 1)
    const endPage = Math.min(totalPages - 1, currentPage + 1)

    for (let i = startPage; i <= endPage; i++) {
      pages.push(i)
    }

    if (currentPage < totalPages - 3) {
      pages.push("...") // 当前页小于倒数第4页，显示后省略号
    }

    if (totalPages > 1) {
      pages.push(totalPages) // 始终显示最后一页
    }
  }

  return pages
})

// 处理页码输入框变化
const handlePageInputChange = (value: number | string) => {
  const pageNum = Number(value)
  if (isNaN(pageNum)) return

  jumpToPage(Math.floor(pageNum))
}

// 处理每页条数变化
const handleSizeChange = () => {
  queryParams.page = 1
  currentPageInput.value = 1
  handleQuery()
}

// 下载日志
const handleDownload = async () => {
  if (!queryParams.logfile) {
    ElMessage.warning("请选择日志文件")
    return
  }

  try {
    loading.value = true
    let content

    if (mode.value === "ssh") {
      if (!queryParams.host) {
        ElMessage.warning("请选择主机")
        return
      }
      const blob = await downloadLogApi(queryParams)
      content = blob
    } else {
      const response = await readLocalLogFileApi({
        path: queryParams.logfile,
        maxResults: 10000
      })

      if (response && response.data) {
        const logData = response.data
        let lines: string[] = []

        if (Array.isArray(logData)) {
          lines = logData
        } else if (typeof logData === "object" && Array.isArray(logData.entries)) {
          lines = logData.entries.map((entry) => entry.content)
        }

        content = new Blob([lines.join("\n")], { type: "text/plain" })
      } else {
        throw new Error("返回数据格式错误: 缺少 data 字段")
      }
    }

    const url = window.URL.createObjectURL(content)
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
  if (mode.value === "ssh") {
    console.log("initData")
    // 只有在userStore中没有SSH连接信息时才从接口获取
    console.log("userStore.sshConnections", userStore.sshConnections)
    if (!userStore.sshConnections || userStore.sshConnections.length === 0) {
      await getSSHConnectionsApi()
    }

    // 刷新主机列表
    await refreshHosts()

    // 如果有当前SSH主机，则自动选择
    if (userStore.currentSSHHost) {
      queryParams.host = userStore.currentSSHHost
      await handleHostChange(userStore.currentSSHHost)
    }
  } else {
    // 获取本地日志文件列表
    await getLocalLogFiles()
  }

  // 初始化当前页输入框
  currentPageInput.value = queryParams.page
}

onMounted(() => {
  initData()
})

onUnmounted(() => {
  stopAutoRefresh()
})

// 获取日志级别
const getLogLevel = (content: string): string => {
  if (content.includes("ERROR") || content.includes("FATAL")) return "ERROR"
  if (content.includes("WARN") || content.includes("WARNING")) return "WARN"
  if (content.includes("INFO")) return "INFO"
  if (content.includes("DEBUG")) return "DEBUG"
  if (content.includes("TRACE")) return "TRACE"
  return ""
}

// 获取日志级别对应的CSS类名
const getLogLevelClass = (content: string): string => {
  const level = getLogLevel(content)
  if (level) return `log-level-${level.toLowerCase()}`
  return ""
}

// 高亮显示内容
const highlightContent = (content: string): string => {
  if (!queryParams.keyword || !content.includes(queryParams.keyword)) return content

  // 转义HTML特殊字符
  const escapeHtml = (text: string) => {
    return text
      .replace(/&/g, "&amp;")
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;")
      .replace(/"/g, "&quot;")
      .replace(/'/g, "&#039;")
  }

  // 高亮关键字
  const escapedContent = escapeHtml(content)
  const escapedKeyword = escapeHtml(queryParams.keyword)
  return escapedContent.replace(
    new RegExp(escapedKeyword, "gi"),
    (match) => `<span class="keyword-highlight">${match}</span>`
  )
}

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

// 处理主机双击事件
const handleHostDblClick = async (host: string) => {
  if (!host) return

  try {
    // 获取保存的连接信息
    const savedConnection = userStore.getSSHConnection(host)

    if (savedConnection) {
      // 使用保存的连接信息连接
      await connectSSHApi(savedConnection.host, savedConnection.port, savedConnection.user, savedConnection.password)
      ElMessage.success("SSH连接成功")

      // 设置为当前选中的主机
      queryParams.host = host

      // 获取该主机的日志文件列表
      await handleHostChange(host)
    } else {
      ElMessage.warning("未找到该主机的连接信息")
    }
  } catch (error) {
    ElMessage.error(`SSH连接失败: ${error instanceof Error ? error.message : String(error)}`)
  }
}
</script>

<style scoped>
.filter-container {
  margin-bottom: 20px;
  background: #fff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.log-content {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.log-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 20px;
  background: #f9fafc;
  border-bottom: 1px solid #ebeef5;
}

.log-stats-wrapper {
  display: flex;
  gap: 24px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.stat-label {
  font-size: 14px;
  color: #606266;
}

.stat-value {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.log-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.refresh-select {
  width: 120px;
}

.log-viewer {
  height: 600px;
  overflow-y: auto;
  background: #ffffff;
  font-family: "JetBrains Mono", "Fira Code", "Source Code Pro", Consolas, monospace;
  font-size: 13px;
  line-height: 1.5;
  position: relative;
}

.log-line {
  display: flex;
  padding: 4px 0;
  border-bottom: 1px solid #f0f0f0;
  transition: background-color 0.15s ease;
}

.log-line:hover {
  background-color: #f5f7fa;
}

.line-meta {
  display: flex;
  flex-shrink: 0;
  width: 120px;
  padding: 0 12px;
  border-right: 1px solid #ebeef5;
}

.line-number {
  color: #909399;
  font-size: 12px;
  width: 40px;
  text-align: right;
  padding-right: 8px;
  user-select: none;
}

.line-level {
  font-size: 12px;
  font-weight: 600;
  width: 50px;
  text-align: center;
  border-radius: 3px;
  padding: 0 4px;
}

.line-content {
  flex: 1;
  padding: 0 12px;
  white-space: pre-wrap;
  word-break: break-all;
  color: #303133;
}

.keyword-highlight {
  background-color: #ffeb3b;
  color: #000;
  font-weight: bold;
  border-radius: 2px;
  padding: 0 2px;
}

.log-level-error {
  color: #f56c6c;
  background-color: rgba(245, 108, 108, 0.1);
}

.log-level-warn {
  color: #e6a23c;
  background-color: rgba(230, 162, 60, 0.1);
}

.log-level-info {
  color: #409eff;
  background-color: rgba(64, 158, 255, 0.1);
}

.log-level-debug {
  color: #67c23a;
  background-color: rgba(103, 194, 58, 0.1);
}

.log-level-trace {
  color: #909399;
  background-color: rgba(144, 147, 153, 0.1);
}

.pagination-container {
  padding: 16px;
  background: #fff;
  border-top: 1px solid #ebeef5;
}

.pagination-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-wrap: wrap;
  gap: 16px;
}

.pagination-info {
  color: #606266;
  font-size: 14px;
}

.pagination-controls {
  display: flex;
  align-items: center;
}

.pagination-btn {
  display: inline-flex;
  justify-content: center;
  align-items: center;
  min-width: 32px;
  height: 32px;
  padding: 0 4px;
  margin: 0 4px;
  font-size: 13px;
  background-color: #fff;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  color: #606266;
  cursor: pointer;
  transition: all 0.3s;
}

.pagination-btn:hover:not(.disabled):not(.active) {
  color: #409eff;
  border-color: #c6e2ff;
  background-color: #ecf5ff;
}

.pagination-btn.active {
  color: #fff;
  background-color: #409eff;
  border-color: #409eff;
}

.pagination-btn.disabled {
  color: #c0c4cc;
  cursor: not-allowed;
  background-color: #fff;
  border-color: #ebeef5;
}

.pagination-icon {
  font-style: normal;
  font-size: 16px;
}

.pagination-pages {
  display: flex;
  align-items: center;
}

.page-num {
  font-weight: 500;
}

.pagination-ellipsis {
  display: inline-flex;
  justify-content: center;
  align-items: center;
  width: 32px;
  height: 32px;
  margin: 0 4px;
  color: #606266;
}

.pagination-goto {
  display: flex;
  align-items: center;
  font-size: 14px;
  color: #606266;
}

.goto-input-wrapper {
  margin: 0 8px;
  width: 50px;
}

.goto-input {
  width: 100%;
  height: 32px;
  padding: 0 8px;
  text-align: center;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  color: #606266;
  transition: border-color 0.3s;
}

.goto-input:focus {
  outline: none;
  border-color: #409eff;
}

.goto-btn {
  margin-left: 8px;
  padding: 0 12px;
  height: 32px;
  background-color: #f4f4f5;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  color: #606266;
  cursor: pointer;
  transition: all 0.3s;
}

.goto-btn:hover {
  color: #409eff;
  border-color: #c6e2ff;
  background-color: #ecf5ff;
}

.page-size-selector {
  display: flex;
  align-items: center;
}

.page-size-select {
  height: 32px;
  padding: 0 8px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  color: #606266;
  background-color: #fff;
  cursor: pointer;
  transition: border-color 0.3s;
}

.page-size-select:focus {
  outline: none;
  border-color: #409eff;
}

.empty-logs {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  color: #909399;
}
</style>
