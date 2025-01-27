<template>
  <div class="app-container">
    <div class="form-container">
      <form @submit.prevent="logfileRead" class="log-form">
        <div class="log-form-container">
          <div class="left-element">
            <div class="form-group">
              <label>日志文件 *</label>
              <input type="text" v-model="logfile" class="form-control" placeholder="请输入日志文件路径" />
              <span class="error_text">{{ logfileError }}</span>
            </div>
          </div>
          <div class="right-element">
            <div class="form-group">
              <label>主机 *</label>
              <select v-model="selectedHost" class="form-control" @change="fetchPaths">
                <option v-for="host in hosts" :key="host" :value="host">{{ host }}</option>
              </select>
              <span class="error_text">{{ hostError }}</span>
            </div>
          </div>
        </div>

        <div class="log-form-container">
          <div class="left-element">
            <div class="form-group">
              <label>路径 *</label>
              <select v-model="selectedPath" class="form-control">
                <option v-for="path in paths" :key="path" :value="path">{{ path }}</option>
              </select>
              <span class="error_text">{{ pathError }}</span>
            </div>
          </div>
          <div class="right-element">
            <div class="form-group">
              <label>匹配</label>
              <input type="text" v-model="match" class="form-control" placeholder="请输入匹配内容" />
              <span class="error_text">{{ matchError }}</span>
            </div>
          </div>
        </div>
        <div class="form-group">
          <label>
            <input type="checkbox" v-model="filterSearchLine" />
            仅查看和输出匹配行
          </label>
          <button type="submit" class="btn btn-primary">查看</button>
        </div>
      </form>
    </div>

    <div
      class="log-content"
      v-if="logContentVisible"
      v-infinite-scroll="loadMoreLogs"
      infinite-scroll-disabled="loading"
      infinite-scroll-distance="10"
    >
      <el-table :data="logContent" style="width: 100%">
        <el-table-column prop="totalLines" label="文件总行数" />
        <el-table-column prop="matchLines" label="匹配总行数" />
        <el-table-column prop="lines" label="窗口行数" />
        <el-table-column prop="highlightLines" label="高亮行数" />
        <el-table-column prop="page" label="页码" />
      </el-table>
      <div v-if="loading">加载中...</div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from "vue"

// 定义日志条目的接口
interface LogEntry {
  totalLines: number
  matchLines: number
  lines: number
  highlightLines: number
  page: number
}

// 添加接口类型定义
interface LogResponse {
  totalLines: number
  matchLines: number
  lines: string[]
  highlightLines: number[]
  page: number
}

const logfile = ref("")
const selectedHost = ref("")
const selectedPath = ref("")
const match = ref("")
const filterSearchLine = ref(true)
const logContent = ref<LogEntry[]>([])
const logContentVisible = ref(false)

const logfileError = ref("")
const hostError = ref("")
const pathError = ref("")
const matchError = ref("")

const hosts = ref<string[]>(["localhost", "remote-host"]) // 示例主机列表
const paths = ref<string[]>(["/tmp/loggrove.log"]) // 示例路径列表

const loading = ref(false) // 添加加载状态

const fetchHosts = async () => {
  // 这里调用 API 获取主机列表
}

const fetchPaths = async () => {
  // 这里调用 API 获取路径列表
}

const logfileRead = async () => {
  try {
    // 验证必填项
    if (!logfile.value) {
      logfileError.value = "请输入日志文件"
      return
    }
    if (!selectedHost.value) {
      hostError.value = "请选择主机"
      return
    }
    if (!selectedPath.value) {
      pathError.value = "请选择路径"
      return
    }

    // 清除错误提示
    logfileError.value = ""
    hostError.value = ""
    pathError.value = ""

    loading.value = true
    const params = {
      logfile: logfile.value,
      host: selectedHost.value,
      path: selectedPath.value,
      match: match.value,
      filterSearchLine: filterSearchLine.value,
      page: 1,
      pageSize: 100
    }

    const res = await fetch('/api/log/read', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(params)
    })

    const data = await res.json()
    logContent.value = data.logs
    logContentVisible.value = true
  } catch (error) {
    console.error('读取日志失败:', error)
  } finally {
    loading.value = false
  }
}

// 实现加载更多日志的方法
const fetchMoreLogs = async () => {
  if (loading.value) return []

  try {
    const currentPage = Math.ceil(logContent.value.length / 100) + 1
    const params = {
      logfile: logfile.value,
      host: selectedHost.value,
      path: selectedPath.value,
      match: match.value,
      filterSearchLine: filterSearchLine.value,
      page: currentPage,
      pageSize: 100
    }

    const res = await fetch('/api/log/read', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(params)
    })

    const data = await res.json()
    return data.logs
  } catch (error) {
    console.error('加载更多日志失败:', error)
    return []
  }
}

// 在组件挂载时获取主机列表
fetchHosts()
logfileRead()
</script>

<style scoped>
.form-container {
  margin-bottom: 20px;
  background-color: #f8f9fa; /* 更柔和的背景色 */
  padding: 20px;
  border-radius: 8px; /* 更圆的边角 */
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1); /* 更明显的阴影效果 */
}

.panel {
  border: 1px solid #dee2e6;
  border-radius: 0.25rem;
}

.panel-body {
  padding: 15px;
}

.log-form-container {
  display: flex;
  align-items: center; /* 垂直居中 */
  justify-content: space-between; /* 水平分布 */
  margin-bottom: 15px; /* 增加间距 */
}

.left-element,
.right-element {
  flex: 1; /* 使左右元素占据相同的空间 */
  margin-right: 10px; /* 调整间距 */
}

.right-element {
  margin-right: 0; /* 右侧元素不需要右边距 */
}

.form-group {
  margin-bottom: 15px;
}

.form-control {
  border: 1px solid #ced4da;
  border-radius: 4px; /* 更圆的边角 */
  padding: 10px;
  transition: border-color 0.3s; /* 添加过渡效果 */
}

.form-control:focus {
  border-color: #007bff; /* 聚焦时的边框颜色 */
  box-shadow: 0 0 5px rgba(0, 123, 255, 0.5); /* 聚焦时的阴影效果 */
}

.btn-primary {
  background-color: #007bff;
  border-color: #007bff;
  border-radius: 4px; /* 更圆的边角 */
  padding: 10px 15px; /* 增加内边距 */
  transition: background-color 0.3s; /* 添加过渡效果 */
}

.btn-primary:hover {
  background-color: #0056b3; /* 悬停时的背景颜色 */
}

.error_text {
  color: #dc3545; /* 更鲜明的错误提示颜色 */
  font-size: 0.875em;
}

.log-content {
  max-height: 500px;
  overflow: auto;
  line-height: 20px;
  background-color: #ffffff;
  padding: 20px;
  border-radius: 5px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
}
</style>
