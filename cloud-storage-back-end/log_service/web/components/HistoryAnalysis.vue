<template>
  <div class="history-analysis">
    <div class="filter-section">
      <div class="time-filter">
        <span class="filter-label">时间范围：</span>
        <div class="time-picker-container">
          <el-date-picker
            v-model="timeRange"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="X"
            :default-time="['00:00:00', '23:59:59']"
          />
        </div>
        <div class="aggregate-by">
          <span class="filter-label">聚合方式：</span>
          <el-select v-model="aggregateBy" placeholder="请选择">
            <el-option label="按分钟" value="按分钟" />
            <el-option label="按小时" value="按小时" />
            <el-option label="按天" value="按天" />
          </el-select>
        </div>
      </div>

      <div class="search-filter">
        <div class="log-file-select">
          <span class="filter-label">日志文件：</span>
          <el-select v-model="logFile" placeholder="请选择日志文件">
            <el-option
              v-for="file in logFiles"
              :key="file"
              :label="file.split('/').pop()"
              :value="file"
            ></el-option>
          </el-select>
        </div>
        <div class="keyword-search">
          <span class="filter-label">关键字：</span>
          <el-input v-model="keywords" placeholder="请输入关键字" />
        </div>
      </div>

      <div class="action-buttons">
        <el-button type="primary" @click="searchLogs">查询</el-button>
        <el-button @click="resetFilters">重置</el-button>
      </div>
    </div>

    <div class="log-table-section">
      <el-table
        v-loading="loading"
        :data="logEntries"
        style="width: 100%"
        border
        stripe
      >
        <el-table-column prop="timestamp" label="时间" width="180">
          <template #default="scope">
            {{ formatTimestamp(scope.row.timestamp) }}
          </template>
        </el-table-column>
        <el-table-column prop="level" label="级别" width="100">
          <template #default="scope">
            <el-tag
              :type="getTagType(scope.row.level)"
              size="small"
            >
              {{ scope.row.level }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="source" label="来源" width="150" />
        <el-table-column prop="content" label="内容" />
      </el-table>

      <div class="pagination-container">
        <el-pagination
          background
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          :page-size="pageSize"
          :current-page="currentPage"
          :page-sizes="[10, 20, 50, 100]"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, watch } from 'vue';
import axios from 'axios';

export default {
  name: 'HistoryAnalysis',
  props: {
    host: {
      type: String,
      required: true
    },
    logFiles: {
      type: Array,
      default: () => []
    }
  },
  setup(props) {
    const timeRange = ref([]);
    const aggregateBy = ref('按小时');
    const logFile = ref('');
    const keywords = ref('');
    const logEntries = ref([]);
    const loading = ref(false);
    const total = ref(0);
    const currentPage = ref(1);
    const pageSize = ref(10);

    // 监听logFiles变化
    watch(() => props.logFiles, (newFiles) => {
      if (newFiles.length > 0 && !logFile.value) {
        logFile.value = newFiles[0];
      }
    }, { immediate: true });

    // 格式化时间戳
    const formatTimestamp = (timestamp) => {
      const date = new Date(timestamp * 1000);
      return date.toLocaleString();
    };

    // 根据日志级别获取标签类型
    const getTagType = (level) => {
      switch (level) {
        case 'ERROR':
          return 'danger';
        case 'WARN':
          return 'warning';
        case 'INFO':
          return 'info';
        case 'DEBUG':
          return 'success';
        default:
          return '';
      }
    };

    // 查询日志
    const searchLogs = async () => {
      if (!props.host) {
        alert('请先连接主机');
        return;
      }

      if (!logFile.value) {
        alert('请选择日志文件');
        return;
      }

      loading.value = true;
      try {
        const [startTime, endTime] = timeRange.value || [0, 0];
        
        const response = await axios.post('/api/log_service/v1/monitor/history', {
          host: props.host,
          log_file: logFile.value,
          start_time: startTime || 0,
          end_time: endTime || 0,
          keywords: keywords.value,
          page: currentPage.value,
          page_size: pageSize.value,
          aggregate_by: aggregateBy.value
        });

        if (response.data.success) {
          logEntries.value = response.data.data;
          total.value = response.data.total;
        } else {
          console.error('查询日志失败:', response.data);
        }
      } catch (error) {
        console.error('查询日志失败:', error);
      } finally {
        loading.value = false;
      }
    };

    // 重置筛选条件
    const resetFilters = () => {
      timeRange.value = [];
      aggregateBy.value = '按小时';
      keywords.value = '';
      currentPage.value = 1;
      searchLogs();
    };

    // 处理每页显示数量变化
    const handleSizeChange = (size) => {
      pageSize.value = size;
      searchLogs();
    };

    // 处理页码变化
    const handleCurrentChange = (page) => {
      currentPage.value = page;
      searchLogs();
    };

    // 组件挂载时加载数据
    onMounted(() => {
      // 设置默认时间范围为最近24小时
      const now = Math.floor(Date.now() / 1000);
      timeRange.value = [now - 86400, now];
      
      // 如果已连接主机且有日志文件，则查询日志
      if (props.host && props.logFiles.length > 0) {
        logFile.value = props.logFiles[0];
        searchLogs();
      }
    });

    return {
      timeRange,
      aggregateBy,
      logFile,
      keywords,
      logEntries,
      loading,
      total,
      currentPage,
      pageSize,
      formatTimestamp,
      getTagType,
      searchLogs,
      resetFilters,
      handleSizeChange,
      handleCurrentChange
    };
  }
};
</script>

<style scoped>
.history-analysis {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.filter-section {
  padding: 20px;
  background-color: #f5f7fa;
  border-radius: 4px;
  margin-bottom: 20px;
}

.time-filter, .search-filter {
  display: flex;
  flex-wrap: wrap;
  margin-bottom: 15px;
  gap: 20px;
}

.filter-label {
  font-weight: bold;
  margin-right: 10px;
}

.time-picker-container {
  width: 380px;
}

.log-file-select, .host-select, .keyword-search, .aggregate-by {
  display: flex;
  align-items: center;
}

.action-buttons {
  margin-top: 15px;
}

.log-table-section {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style> 