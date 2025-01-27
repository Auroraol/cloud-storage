<!-- src/components/Chart.vue -->
<template>
  <div class="app-container">
    <div class="chart-container">
      <el-tabs :stretch="true">
        <el-tab-pane label="`连续时间" class="toLogin">
          <div class="continuous-time-form">
            <div class="panel panel-default">
              <div class="panel-body">
                <form role="form" class="form-horizontal tab-content" @submit.prevent="fetchChartData">
                  <!-- 修改日志查询条件区域 -->
                  <div class="search-conditions">
                    <div class="form-group">
                      <div class="col-sm-3">
                        <input type="text" v-model="logfile" class="form-control" placeholder="LoggroveLog" />
                      </div>
                      <div class="col-sm-3">
                        <select v-model="selectedHosts" class="form-control">
                          <option value="localhost">localhost</option>
                        </select>
                      </div>
                      <div class="col-sm-3">
                        <select v-model="selectedMonitorItems" class="form-control">
                          <option value="Login, Total">Login, Total</option>
                        </select>
                      </div>
                      <div class="col-sm-3">
                        <div class="btn-group">
                          <button type="button" class="btn btn-default" @click="removeItemGroup">
                            <i class="fa fa-minus" />
                          </button>
                          <button type="button" class="btn btn-default" @click="addItemGroup">
                            <i class="fa fa-plus" />
                          </button>
                        </div>
                      </div>
                    </div>
                  </div>

                  <!-- 修改时间选择区域 -->
                  <div class="time-selector">
                    <div class="quick-select">
                      <button
                        v-for="(label, time) in timeOptions"
                        :key="time"
                        type="button"
                        class="btn btn-default"
                        @click="selectTimeRange(time)"
                      >
                        {{ label }}
                      </button>
                    </div>

                    <div class="custom-time">
                      <div class="col-sm-5">
                        <input type="text" v-model="beginTime" class="form-control" placeholder="开始时间" />
                      </div>
                      <div class="col-sm-5">
                        <input type="text" v-model="endTime" class="form-control" placeholder="结束时间" />
                      </div>
                      <div class="col-sm-2">
                        <button type="button" class="btn btn-primary" @click="showIntervalChart">查看</button>
                      </div>
                    </div>
                  </div>
                </form>
              </div>
            </div>
            <!-- 图表展示区域 -->
            <div class="chart-display" v-if="chartVisible">
              <div class="panel panel-default">
                <div class="panel-body" id="log_chart">
                  <!-- 图表将在这里渲染 -->
                </div>
              </div>
            </div>
          </div>
        </el-tab-pane>
        <el-tab-pane label="日期对比" class="toLogin" />
      </el-tabs>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from "vue"

const logfile = ref("")
const selectedHosts = ref<string[]>([])
const selectedMonitorItems = ref<string[]>([])
const logfileError = ref("")
const beginTime = ref("")
const endTime = ref("")
const beginTimeError = ref("")
const endTimeError = ref("")
const hosts = ref<string[]>([]) // 从 API 获取的主机列表
const monitorItems = ref<string[]>([]) // 从 API 获取的监控项列表
const chartVisible = ref(false)
const yesterday = new Date(new Date().setDate(new Date().getDate() - 1)).toISOString().split("T")[0]
const lastWeekSameDay = new Date(new Date().setDate(new Date().getDate() - 7)).toISOString().split("T")[0]
const lastMonthSameDay = new Date(new Date().setMonth(new Date().getMonth() - 1)).toISOString().split("T")[0]
const lastYearSameDay = new Date(new Date().setFullYear(new Date().getFullYear() - 1)).toISOString().split("T")[0]
const contrastDate = ref("")
const dateError = ref("")

const timeOptions = {
  "1h": "最近1小时",
  "6h": "最近6小时",
  "12h": "最近12小时",
  "24h": "最近24小时",
  "2d": "最近2天",
  "7d": "最近7天"
}

const selectTimeRange = (range: string) => {
  // 处理时间范围选择逻辑
  showIntervalChart(range)
}

const fetchChartData = async () => {
  // 这里调用 API 获取图表数据
  // 处理逻辑
}

const showIntervalChart = (timeFrame: string) => {
  // 这里调用 API 显示时间区间图表
}

const showContrastChart = (date: string) => {
  // 这里调用 API 显示对比图表
}

// 其他方法和逻辑
</script>

<style scoped>
.continuous-time-form {
  background: #fff;
  border-radius: 4px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
  margin-bottom: 20px;
}

.form-header {
  padding: 15px 15px 0;
}

.nav-tabs {
  border-bottom: 1px solid #ddd;
}

.nav-tabs > li > a {
  color: #666;
  padding: 8px 15px;
}

.search-conditions {
  padding: 20px 15px;
  border-bottom: 1px solid #eee;
}

.time-selector {
  padding: 20px 15px;
}

.quick-select {
  margin-bottom: 15px;
}

.quick-select .btn {
  margin-right: 8px;
  margin-bottom: 8px;
  padding: 6px 12px;
  border-radius: 4px;
  background: #f5f5f5;
  border: 1px solid #ddd;
}

.custom-time {
  display: flex;
  align-items: center;
  gap: 10px;
}

.form-control {
  height: 34px;
  border: 1px solid #ddd;
  border-radius: 4px;
  padding: 6px 12px;
}

.btn-primary {
  background-color: #1890ff;
  border-color: #1890ff;
  color: white;
}

.btn-primary:hover {
  background-color: #40a9ff;
  border-color: #40a9ff;
}

.chart-display {
  background: #fff;
  border-radius: 4px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
  padding: 20px;
}

.error_text {
  color: #ff4d4f;
  font-size: 12px;
  margin-top: 4px;
}
</style>
