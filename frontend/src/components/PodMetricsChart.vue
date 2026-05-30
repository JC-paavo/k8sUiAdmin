<template>
  <div class="pod-metrics-charts" v-if="hasData">
    <div class="charts-row">
      <div class="chart-card">
        <h4 class="chart-title">CPU 使用量 (cores)</h4>
        <div ref="cpuChartRef" style="width:100%;height:280px"></div>
      </div>
      <div class="chart-card">
        <h4 class="chart-title">内存使用量 (MiB)</h4>
        <div ref="memChartRef" style="width:100%;height:280px"></div>
      </div>
    </div>
  </div>
  <div v-else class="metrics-loading" v-loading="loading">
    <span v-if="!loading && !error">暂无监控数据（等待采集...）</span>
    <span v-if="error" style="color:#f56c6c">{{ error }}</span>
  </div>
</template>

<script setup>
import { ref, watch, onMounted, onBeforeUnmount, computed } from 'vue'
import * as echarts from 'echarts'
import { k8sAPI } from '@/utils/api'

const props = defineProps({
  clusterId: { type: [Number, String], required: true },
  namespace: { type: String, required: true },
  podName: { type: String, required: true }
})

const cpuChartRef = ref(null)
const memChartRef = ref(null)
const loading = ref(false)
const error = ref('')
const metricsData = ref([])

let cpuChart = null
let memChart = null
let timer = null

const hasData = computed(() => metricsData.value.length > 0)

const formatTime = (ts) => {
  if (!ts) return ''
  const d = new Date(ts)
  return d.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit', second: '2-digit' })
}

const makeOption = (title, data, field, unit) => ({
  tooltip: {
    trigger: 'axis',
    formatter: (params) => {
      const p = params[0]
      return `${p.axisValue}<br/>${title}: <b>${p.value.toFixed(2)} ${unit}</b>`
    }
  },
  grid: { top: 10, right: 30, bottom: 30, left: 50 },
  xAxis: {
    type: 'category',
    data: data.map(d => formatTime(d.time)),
    axisLabel: { fontSize: 10, rotate: 30 }
  },
  yAxis: {
    type: 'value',
    axisLabel: { fontSize: 10, formatter: (v) => v.toFixed(1) }
  },
  series: [{
    type: 'line',
    data: data.map(d => d[field]),
    smooth: true,
    symbol: 'none',
    lineStyle: { width: 2 },
    areaStyle: { opacity: 0.1 },
    itemStyle: { color: field === 'cpu' ? '#409eff' : '#67c23a' }
  }]
})

const fetchMetrics = async () => {
  loading.value = true
  error.value = ''
  try {
    const res = await k8sAPI.getPodMetrics(props.clusterId, props.namespace, props.podName)
    metricsData.value = res.data?.data || []
    renderCharts()
  } catch (e) {
    error.value = e.response?.data?.error || '获取监控数据失败'
  } finally {
    loading.value = false
  }
}

const renderCharts = () => {
  if (!cpuChartRef.value || !memChartRef.value) return

  if (!cpuChart) {
    cpuChart = echarts.init(cpuChartRef.value)
    memChart = echarts.init(memChartRef.value)
  }

  cpuChart.setOption(makeOption('CPU', metricsData.value, 'cpu', 'cores'), true)
  memChart.setOption(makeOption('内存', metricsData.value, 'memory', 'MiB'), true)
}

const handleResize = () => {
  cpuChart?.resize()
  memChart?.resize()
}

onMounted(() => {
  fetchMetrics()
  timer = setInterval(fetchMetrics, 60000)
  window.addEventListener('resize', handleResize)
})

onBeforeUnmount(() => {
  clearInterval(timer)
  window.removeEventListener('resize', handleResize)
  cpuChart?.dispose()
  memChart?.dispose()
})

watch(() => [props.podName, props.namespace], () => {
  fetchMetrics()
})
</script>

<style scoped>
.pod-metrics-charts {
  margin-top: 20px;
}

.charts-row {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.chart-card {
  flex: 1;
  min-width: 340px;
  background: #fff;
  border: 1px solid #ebeef5;
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 1px 4px rgba(0,0,0,0.04);
}

.chart-title {
  margin: 0 0 8px 0;
  font-size: 14px;
  font-weight: 600;
  color: #303133;
}

.metrics-loading {
  margin-top: 20px;
  padding: 40px;
  text-align: center;
  color: #909399;
  font-size: 14px;
}
</style>
