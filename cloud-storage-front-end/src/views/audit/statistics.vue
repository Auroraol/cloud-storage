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
import { getAuditStatisticsApi } from "@/api/audit"
// import { useUserStore } from "@/stores/user"
import { ElMessage } from "element-plus"

// const userStore = useUserStore()

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
      data: ["上传", "下载", "删除", "修改"]
    },
    yAxis: {},
    series: [
      {
        name: "数量",
        type: "bar",
        data: [0, 0, 0, 0]
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
    const data = await getAuditStatisticsApi(userStore.userId)
    todayTotal.value = data.todayTotal
    uploadCount.value = data.uploadCount
    downloadCount.value = data.downloadCount
    deleteCount.value = data.deleteCount

    // 更新图表数据
    if (operationChart.value) {
      const chart = echarts.getInstanceByDom(operationChart.value)
      chart?.setOption({
        series: [
          {
            data: [data.uploadCount, data.downloadCount, data.deleteCount, data.modifyCount]
          }
        ]
      })
    }

    if (userChart.value) {
      const chart = echarts.getInstanceByDom(userChart.value)
      chart?.setOption({
        series: [
          {
            data: data.dailyTrend
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
