<template>
  <div class="echarts-test-container">
    <h2>ECharts 测试页面</h2>

    <div class="chart-actions">
      <el-button type="primary" @click="renderTestChart">渲染测试图表</el-button>
      <el-button type="success" @click="renderActualDataChart">渲染实际数据图表</el-button>
      <el-button type="warning" @click="clearChart">清空图表</el-button>
    </div>

    <div ref="chartContainer" class="chart-container" />

    <div class="chart-info" v-if="chartInfo">
      <h3>图表信息</h3>
      <pre>{{ chartInfo }}</pre>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue"
import * as echarts from "echarts"

const chartContainer = ref<HTMLElement | null>(null)
let chartInstance: echarts.ECharts | null = null
const chartInfo = ref("")

// 测试数据
const testData = [
  {
    name: "测试数据",
    type: "line",
    data: [
      [Date.now() - 86400000, 10],
      [Date.now() - 43200000, 20],
      [Date.now(), 30]
    ]
  }
]

// 模拟实际数据
const actualData = [
  {
    name: "INFO",
    type: "line",
    smooth: true,
    showSymbol: true,
    symbolSize: 8,
    lineStyle: { width: 2 },
    data: [
      [1741028346000, 1],
      [1741055909000, 1],
      [1741028346000, 1],
      [1741055909000, 1],
      [1741025346000, 1],
      [1741055909000, 1],
      [1733280449000, 8]
    ]
  },
  {
    name: "ERROR",
    type: "line",
    smooth: true,
    showSymbol: true,
    symbolSize: 8,
    lineStyle: { width: 2 },
    data: [
      [1741025946000, 5],
      [1696170900000, 1]
    ]
  }
]

onMounted(() => {
  initChart()
})

onUnmounted(() => {
  disposeChart()
})

// 初始化图表
const initChart = () => {
  if (chartContainer.value) {
    // 确保容器有尺寸
    chartContainer.value.style.width = "100%"
    chartContainer.value.style.height = "400px"

    // 初始化图表
    chartInstance = echarts.init(chartContainer.value)

    // 记录图表信息
    chartInfo.value = `
ECharts 版本: ${echarts.version}
容器尺寸: ${chartContainer.value.offsetWidth}x${chartContainer.value.offsetHeight}
图表实例: ${chartInstance ? "已创建" : "未创建"}
    `
  }
}

// 销毁图表
const disposeChart = () => {
  if (chartInstance) {
    chartInstance.dispose()
    chartInstance = null
  }
}

// 渲染测试图表
const renderTestChart = () => {
  if (!chartInstance) {
    initChart()
  }

  if (chartInstance) {
    try {
      chartInstance.clear()

      const option = {
        title: {
          text: "测试图表"
        },
        tooltip: {
          trigger: "axis"
        },
        xAxis: {
          type: "time"
        },
        yAxis: {
          type: "value"
        },
        series: testData
      }

      chartInstance.setOption(option)
      chartInfo.value += "\n测试图表渲染成功"
    } catch (error) {
      console.error("渲染测试图表失败:", error)
      chartInfo.value += `\n测试图表渲染失败: ${error}`
    }
  }
}

// 渲染实际数据图表
const renderActualDataChart = () => {
  if (!chartInstance) {
    initChart()
  }

  if (chartInstance) {
    try {
      chartInstance.clear()

      const option = {
        title: {
          text: "实际数据图表"
        },
        tooltip: {
          trigger: "axis",
          formatter: function (params: any) {
            if (!params || !params.length) return ""

            const date = new Date(params[0].value[0])
            const formattedDate = `${date.getFullYear()}-${(date.getMonth() + 1).toString().padStart(2, "0")}-${date.getDate().toString().padStart(2, "0")} ${date.getHours().toString().padStart(2, "0")}:${date.getMinutes().toString().padStart(2, "0")}:${date.getSeconds().toString().padStart(2, "0")}`

            let result = `${formattedDate}<br/>`
            params.forEach((param: any) => {
              const value = Array.isArray(param.value) && param.value.length > 1 ? param.value[1] : param.value
              result += `${param.seriesName}: ${value}<br/>`
            })

            return result
          }
        },
        legend: {
          data: ["INFO", "ERROR"],
          type: "scroll",
          orient: "horizontal",
          bottom: 0
        },
        grid: {
          left: "3%",
          right: "4%",
          bottom: "10%",
          top: "10%",
          containLabel: true
        },
        xAxis: {
          type: "time",
          splitLine: {
            show: false
          }
        },
        yAxis: {
          type: "value",
          splitLine: {
            lineStyle: {
              type: "dashed"
            }
          }
        },
        series: actualData
      }

      chartInstance.setOption(option)

      // 检查数据范围
      if (actualData.length > 0 && actualData[0].data.length > 0) {
        const timestamps = actualData[0].data.map((point: any) => point[0])
        const minTime = Math.min(...timestamps)
        const maxTime = Math.max(...timestamps)
        chartInfo.value += `\n数据时间范围: ${new Date(minTime)} 至 ${new Date(maxTime)}`
        chartInfo.value += `\n时间跨度(天): ${(maxTime - minTime) / (1000 * 60 * 60 * 24)}`
      }

      chartInfo.value += "\n实际数据图表渲染成功"
    } catch (error) {
      console.error("渲染实际数据图表失败:", error)
      chartInfo.value += `\n实际数据图表渲染失败: ${error}`
    }
  }
}

// 清空图表
const clearChart = () => {
  if (chartInstance) {
    chartInstance.clear()
    chartInfo.value += "\n图表已清空"
  }
}
</script>

<style scoped>
.echarts-test-container {
  padding: 20px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.chart-actions {
  margin-bottom: 20px;
}

.chart-container {
  height: 400px;
  border: 1px solid #ebeef5;
  border-radius: 4px;
  margin-bottom: 20px;
}

.chart-info {
  padding: 15px;
  background-color: #f8f9fa;
  border-radius: 4px;
  margin-top: 20px;
}

.chart-info pre {
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
