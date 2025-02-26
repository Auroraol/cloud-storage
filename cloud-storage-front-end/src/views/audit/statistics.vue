<template>
  <div class="audit-statistics">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card>
          <template #header>今日操作总数</template>
          <div class="statistics-number">{{ todayTotal }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <template #header>文件上传数</template>
          <div class="statistics-number">{{ uploadCount }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <template #header>文件下载数</template>
          <div class="statistics-number">{{ downloadCount }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <template #header>文件删除数</template>
          <div class="statistics-number">{{ deleteCount }}</div>
        </el-card>
      </el-col>
    </el-row>

    <div class="charts-container">
      <div class="chart" ref="operationChart" />
      <div class="chart" ref="userChart" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import * as echarts from "echarts"
import { auditApi } from "@/api/audit"
import { ElMessage } from "element-plus"

const todayTotal = ref(0)
const uploadCount = ref(0)
const downloadCount = ref(0)
const deleteCount = ref(0)

const operationChart = ref<HTMLElement>()
const userChart = ref<HTMLElement>()

onMounted(() => {
  initOperationChart()
  initUserChart()
  fetchStatisticsData()
})

const initOperationChart = () => {
  const chart = echarts.init(operationChart.value!)
  chart.setOption({
    title: {
      text: "我的操作类型统计"
    },
    tooltip: {},
    xAxis: {
      data: ["上传", "下载", "删除", "恢复", "重命名", "移动", "复制", "创建文件夹", "修改"]
    },
    yAxis: {},
    series: [
      {
        name: "数量",
        type: "bar",
        data: [0, 0, 0, 0, 0, 0, 0, 0, 0]
      }
    ]
  })
}

const initUserChart = () => {
  const chart = echarts.init(userChart.value!)
  chart.setOption({
    title: {
      text: "最近30天操作趋势"
    },
    tooltip: {},
    xAxis: {
      type: "time"
    },
    yAxis: {},
    series: [
      {
        name: "操作次数",
        type: "line",
        data: []
      }
    ]
  })
}

const fetchStatisticsData = async () => {
  try {
    const data = await auditApi.getAuditList({
      flag: -1, // 获取所有类型的操作
      start_time: Math.floor(Date.now() / 1000) - 30 * 24 * 60 * 60, // 30天前
      end_time: Math.floor(Date.now() / 1000), // 现在
      page: 1,
      page_size: 1000
    })

    // 统计各类操作数量
    const operationCounts = data.data.operation_logs.reduce(
      (acc, log) => {
        acc[log.flag] = (acc[log.flag] || 0) + 1
        return acc
      },
      {} as Record<number, number>
    )

    // 更新今日统计数据
    const today = new Date()
    today.setHours(0, 0, 0, 0)
    const todayTimestamp = Math.floor(today.getTime() / 1000)

    const todayLogs = data.data.operation_logs.filter((log) => parseInt(log.created_at) >= todayTimestamp)

    todayTotal.value = todayLogs.length
    uploadCount.value = operationCounts[0] || 0 // 上传
    downloadCount.value = operationCounts[1] || 0 // 下载
    deleteCount.value = operationCounts[2] || 0 // 删除

    // 更新操作类型统计图表
    if (operationChart.value) {
      const chart = echarts.getInstanceByDom(operationChart.value)
      chart?.setOption({
        series: [
          {
            data: [
              operationCounts[0] || 0, // 上传
              operationCounts[1] || 0, // 下载
              operationCounts[2] || 0, // 删除
              operationCounts[3] || 0, // 恢复
              operationCounts[4] || 0, // 重命名
              operationCounts[5] || 0, // 移动
              operationCounts[6] || 0, // 复制
              operationCounts[7] || 0, // 创建文件夹
              operationCounts[8] || 0 // 修改
            ]
          }
        ]
      })
    }

    // 按日期分组统计操作数量
    const dailyStats = data.data.operation_logs.reduce(
      (acc, log) => {
        const date = new Date(parseInt(log.created_at) * 1000).toISOString().split("T")[0]
        acc[date] = (acc[date] || 0) + 1
        return acc
      },
      {} as Record<string, number>
    )

    // 获取最近30天的日期列表
    const dates = []
    for (let i = 0; i < 30; i++) {
      const date = new Date()
      date.setDate(date.getDate() - i)
      dates.unshift(date.toISOString().split("T")[0] as never)
    }

    // 转换为图表所需的数据格式，确保每天都有数据
    const timeSeriesData = dates.map((date) => [date, dailyStats[date] || 0] as [string, number])

    // 更新30天趋势图表
    if (userChart.value) {
      const chart = echarts.getInstanceByDom(userChart.value)
      chart?.setOption({
        series: [
          {
            data: timeSeriesData
          }
        ]
      })
    }
  } catch (error) {
    ElMessage.error("获取统计数据失败")
    console.error("获取统计数据失败:", error)
  }
}
</script>

<style scoped>
.audit-statistics {
  padding: 20px;
}

.statistics-number {
  font-size: 24px;
  font-weight: bold;
  color: #409eff;
  text-align: center;
}

.charts-container {
  margin-top: 20px;
  display: flex;
  gap: 20px;
}

.chart {
  flex: 1;
  height: 400px;
}
</style>
