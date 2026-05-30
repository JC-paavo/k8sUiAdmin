<template>
  <div class="resource-list-page">
    <div v-if="!hasReadPermission" class="no-permission">
      <el-icon name="Lock" size="64" color="#909399" />
      <p>您没有权限访问该集群</p>
    </div>
    
    <div v-else class="page-container">
      <div class="page-header-hero">
        <div class="header-content">
          <div class="header-title-section">
            <el-button @click="goBack" class="back-btn">
              <el-icon><ArrowLeft /></el-icon>
              <span>返回</span>
            </el-button>
            <div class="resource-icon">
              <IconPod :size="48" />
            </div>
            <div class="title-text">
              <h1>Pod 资源管理</h1>
              <div class="cluster-info">
                <el-icon><Location /></el-icon>
                <span>{{ clusterName || '加载中...' }}</span>
                <el-tag size="small" :type="cluster?.status === 'connected' ? 'success' : 'warning'">
                  {{ cluster?.status === 'connected' ? '已连接' : '连接中' }}
                </el-tag>
              </div>
            </div>
          </div>
          <div class="header-stats">
            <div class="stat-card">
              <div class="stat-value">{{ resources.length }}</div>
              <div class="stat-label">Pod 总数</div>
            </div>
            <div class="stat-card">
              <div class="stat-value">{{ runningPods }}</div>
              <div class="stat-label">运行中</div>
            </div>
          </div>
        </div>
      </div>
      
      <div class="toolbar-section">
        <div class="toolbar-left">
          <el-select 
            v-model="selectedNamespace" 
            placeholder="选择命名空间"
            class="namespace-select"
            clearable
            @change="handleNamespaceChange"
          >
            <el-option label="所有命名空间" value="" />
            <el-option 
              v-for="ns in namespaces" 
              :key="ns.name || ns.metadata?.name" 
              :label="ns.name || ns.metadata?.name" 
              :value="ns.name || ns.metadata?.name" 
            />
          </el-select>
          
          <el-input
            v-model="searchQuery"
            placeholder="搜索 Pod 名称..."
            class="search-input"
            clearable
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </div>
        
        <div class="toolbar-right">
          <el-button @click="handleRefresh" class="refresh-btn">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </div>
      
      <div class="resource-table-container">
        <el-table 
          :data="paginatedResources" 
          border 
          stripe
          class="resource-table"
          :header-cell-style="{ background: 'var(--bg-secondary)', color: 'var(--text-primary)', textAlign: 'center' }"
        >
          <el-table-column prop="metadata.name" label="名称" min-width="250">
            <template #default="scope">
              <span class="clickable-name" @click="viewResource(scope.row)">
                <IconPod :size="16" />
                {{ scope.row.metadata.name }}
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="metadata.namespace" label="命名空间" width="150" />
          <el-table-column label="状态" width="180">
            <template #default="scope">
              <div class="pod-status">
                <span :class="['status-dot', getPodStatusClass(scope.row)]"></span>
                <span>{{ getPodStatus(scope.row) }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="节点" width="150">
            <template #default="scope">
              <el-tag size="small" type="info">{{ scope.row.spec.nodeName || '-' }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="Pod IP" width="140">
            <template #default="scope">
              <code class="pod-ip">{{ scope.row.status?.podIP || '-' }}</code>
            </template>
          </el-table-column>
          <el-table-column label="容器数" width="100">
            <template #default="scope">
              <span>{{ scope.row.spec.containers?.length || 0 }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="metadata.creationTimestamp" label="创建时间" width="180">
            <template #default="scope">
              {{ formatTime(scope.row.metadata.creationTimestamp) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="320" fixed="right" :cell-style="{ textAlign: 'center' }">
            <template #default="scope">
              <div class="action-buttons">
                <el-button size="small" @click="openExec(scope.row)" class="action-btn exec-btn" :disabled="!isPodRunning(scope.row)">
                  <el-icon><Monitor /></el-icon>
                  终端
                </el-button>
                <el-button size="small" @click="viewLogs(scope.row)" class="action-btn log-btn">
                  <el-icon><Document /></el-icon>
                  日志
                </el-button>
                <el-button size="small" @click="viewEvents(scope.row)" class="action-btn event-btn">
                  <el-icon><Document /></el-icon>
                  事件
                </el-button>
                <el-button size="small" danger @click="handleDelete(scope.row)" v-if="hasWritePermission" class="action-btn delete-btn">
                  <el-icon><Delete /></el-icon>
                  删除
                </el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>
        
        <div class="pagination-container" v-if="filteredResources.length > 0">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="filteredResources.length"
            layout="total, sizes, prev, pager, next, jumper"
            class="custom-pagination"
          />
        </div>
        
        <div v-if="filteredResources.length === 0" class="empty-state">
          <el-icon :size="80" color="var(--text-tertiary)"><DocumentDelete /></el-icon>
          <h3>暂无 Pod</h3>
          <p v-if="searchQuery">未找到匹配 "{{ searchQuery }}" 的 Pod</p>
          <p v-else>该命名空间下暂无 Pod 资源</p>
        </div>
      </div>
    </div>
    
    <el-dialog title="查看日志" v-model="showLogsModal" width="960px" :close-on-click-modal="false">
      <div class="logs-container">
        <div class="logs-header">
          <div class="logs-path">{{ currentPod?.metadata.namespace }} / {{ currentPod?.metadata.name }}</div>
          <div class="logs-controls">
            <el-select v-model="selectedContainer" placeholder="选择容器" class="container-select">
              <el-option-group label="普通容器">
                <el-option 
                  v-for="container in normalContainers" 
                  :key="container.name" 
                  :label="container.name" 
                  :value="container.name" 
                />
              </el-option-group>
              <el-option-group v-if="initContainers.length > 0" label="Init 容器">
                <el-option 
                  v-for="container in initContainers" 
                  :key="container.name" 
                  :label="container.name" 
                  :value="container.name" 
                />
              </el-option-group>
            </el-select>
            <el-button @click="fetchLogs" size="small" type="primary">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
            <el-button @click="downloadLogs" size="small">
              <el-icon><Download /></el-icon>
              下载
            </el-button>
          </div>
        </div>
        <div class="logs-toolbar">
          <span class="container-tag" v-if="selectedContainer">
            <span v-if="initContainerNames.includes(selectedContainer)" class="init-tag">Init</span>
            <span v-else class="normal-tag">Container</span>
            {{ selectedContainer }}
          </span>
        </div>
        <pre class="logs-content">{{ logsContent }}</pre>
      </div>
      <template #footer>
        <el-button @click="showLogsModal = false">关闭</el-button>
      </template>
    </el-dialog>
    
    <el-dialog title="查看事件" v-model="showEventsModal" width="900px">
      <div class="events-container">
        <div class="events-header">
          <span>{{ currentResource?.metadata.namespace }} / {{ currentResource?.metadata.name }}</span>
          <span class="events-count">共 {{ eventsContent.length }} 个事件</span>
        </div>
        <el-table :data="eventsContent" border max-height="400">
          <el-table-column prop="lastTimestamp" label="时间" width="180" />
          <el-table-column prop="type" label="类型" width="100">
            <template #default="scope">
              <el-tag :type="scope.row.type === 'Normal' ? 'success' : 'danger'" size="small">
                {{ scope.row.type === 'Normal' ? '正常' : '警告' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="reason" label="原因" width="150" />
          <el-table-column prop="message" label="消息" />
          <el-table-column prop="count" label="次数" width="80" />
        </el-table>
        <div v-if="eventsContent.length === 0" class="no-events">
          暂无事件
        </div>
      </div>
      <template #footer>
        <el-button @click="showEventsModal = false">关闭</el-button>
        <el-button type="primary" @click="fetchEvents">刷新</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showTerminal" title="Pod 终端" width="80%" :close-on-click-modal="false" @closed="handleTerminalClosed" destroy-on-close>
      <div class="terminal-header">
        <div class="terminal-info">
          <span class="terminal-path">{{ execPod?.metadata?.namespace }} / {{ execPod?.metadata?.name }}</span>
          <span class="terminal-path" v-if="execContainers.length > 1">选择容器：</span>
          <el-select v-if="execContainers.length > 1" v-model="execContainer" size="small" @change="reconnectTerminal">
            <el-option v-for="c in execContainers" :key="c" :label="c" :value="c" />
          </el-select>
          <span v-else class="terminal-container-badge">{{ execContainer }}</span>
        </div>
      </div>
      <WebTerminal 
        v-if="showTerminal"
        :key="terminalKey"
        :cluster-id="clusterId"
        :namespace="execPod?.metadata?.namespace"
        :pod-name="execPod?.metadata?.name"
        :container="execContainer"
      />
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Delete, Document, Monitor } from '@element-plus/icons-vue'
import { k8sAPI, authAPI, clusterAPI } from '@/utils/api'
import IconPod from '@/assets/icons/IconPod.vue'
import WebTerminal from '@/components/WebTerminal.vue'

const route = useRoute()
const router = useRouter()

const resources = ref([])
const namespaces = ref([])
const selectedNamespace = ref('default')
const showLogsModal = ref(false)
const showEventsModal = ref(false)
const logsContent = ref('')
const eventsContent = ref([])
const currentPod = ref(null)
const currentResource = ref(null)
const selectedContainer = ref('')
const searchQuery = ref('')
const currentPage = ref(1)
const pageSize = ref(20)
const cluster = ref(null)
const showTerminal = ref(false)
const execPod = ref(null)
const execContainer = ref('')
const execContainers = ref([])
const terminalKey = ref(0)

const clusterId = computed(() => route.params.id)
const clusterName = ref('')

const hasReadPermission = ref(false)
const hasWritePermission = ref(false)

const filteredResources = computed(() => {
  if (!searchQuery.value) {
    return resources.value
  }
  const query = searchQuery.value.toLowerCase()
  return resources.value.filter(resource => 
    resource.metadata?.name?.toLowerCase().includes(query)
  )
})

const paginatedResources = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredResources.value.slice(start, end)
})

const runningPods = computed(() => {
  return resources.value.filter(pod => pod.status?.phase === 'Running').length
})

const podContainers = computed(() => {
  if (!currentPod.value?.spec?.containers) return []
  return currentPod.value.spec.containers
})

const normalContainers = computed(() => {
  if (!currentPod.value?.spec?.containers) return []
  return currentPod.value.spec.containers
})

const initContainers = computed(() => {
  if (!currentPod.value?.spec?.initContainers) return []
  return currentPod.value.spec.initContainers
})

const initContainerNames = computed(() => {
  return initContainers.value.map(c => c.name)
})

const allContainers = computed(() => {
  return [...normalContainers.value, ...initContainers.value]
})

const checkPermissions = async () => {
  let id = clusterId.value
  if (!id) return
  
  id = Number(id)
  if (isNaN(id) || id <= 0) return
  
  try {
    const response = await authAPI.checkClusterPermission(id)
    hasReadPermission.value = response.data.has_read
    hasWritePermission.value = response.data.has_write
  } catch (error) {
    console.error('检查权限失败:', error)
    hasReadPermission.value = false
    hasWritePermission.value = false
  }
}

const getPodStatus = (row) => {
  const phase = row.status.phase
  const containerStatuses = row.status.containerStatuses || []
  const totalContainers = containerStatuses.length
  const runningContainers = containerStatuses.filter(c => c.state?.running).length
  
  let statusText = ''
  switch (phase) {
    case 'Running':
      const ready = row.status.conditions?.find(c => c.type === 'Ready')
      if (ready?.status === 'True') {
        statusText = `运行中 (${runningContainers}/${totalContainers})`
      } else {
        statusText = `启动中 (${runningContainers}/${totalContainers})`
      }
      break
    case 'Pending':
      statusText = '等待中'
      break
    case 'Succeeded':
      statusText = '已完成'
      break
    case 'Failed':
      statusText = '失败'
      break
    case 'Unknown':
      statusText = '未知'
      break
    default:
      statusText = phase
  }
  
  return statusText
}

const getPodStatusClass = (row) => {
  const phase = row.status.phase
  if (phase === 'Running') {
    const ready = row.status.conditions?.find(c => c.type === 'Ready')
    if (ready?.status === 'True') return 'status-success'
    return 'status-warning'
  }
  if (phase === 'Succeeded') return 'status-success'
  if (phase === 'Pending') return 'status-warning'
  return 'status-error'
}

const formatTime = (timestamp) => {
  if (!timestamp) return '-'
  const date = new Date(timestamp)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const goBack = () => {
  router.push(`/clusters/${clusterId.value}`)
}

const viewResource = (resource) => {
  const id = clusterId.value
  if (!id) {
    ElMessage.error('集群ID无效')
    return
  }
  const namespace = resource.metadata.namespace
  const name = resource.metadata.name
  router.push(`/clusters/${id}/pods/${namespace}/${name}`)
}

const viewLogs = (pod) => {
  currentPod.value = pod
  if (normalContainers.value.length > 0) {
    selectedContainer.value = normalContainers.value[0].name
  } else if (initContainers.value.length > 0) {
    selectedContainer.value = initContainers.value[0].name
  }
  showLogsModal.value = true
  fetchLogs()
}

const fetchLogs = async () => {
  const id = clusterId.value
  if (!id) {
    logsContent.value = '获取日志失败：集群ID无效'
    return
  }
  
  try {
    const response = await k8sAPI.getPodLogs(
      id, 
      currentPod.value.metadata.namespace, 
      currentPod.value.metadata.name, 
      selectedContainer.value
    )
    logsContent.value = response.data.logs || '暂无日志'
  } catch (error) {
    logsContent.value = '获取日志失败: ' + (error.response?.data?.error || error.message)
  }
}

const downloadLogs = () => {
  if (!logsContent.value || logsContent.value === '暂无日志') {
    ElMessage.warning('暂无日志可下载')
    return
  }
  const podName = currentPod.value?.metadata?.name || 'pod'
  const containerName = selectedContainer.value || 'unknown'
  const blob = new Blob([logsContent.value], { type: 'text/plain;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${podName}_${containerName}_logs.txt`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}

const viewEvents = async (resource) => {
  currentResource.value = resource
  showEventsModal.value = true
  await fetchEvents()
}

const fetchEvents = async () => {
  const id = clusterId.value
  if (!id) {
    eventsContent.value = []
    return
  }
  
  try {
    const response = await k8sAPI.getPodEvents(
      id, 
      currentResource.value.metadata.namespace, 
      currentResource.value.metadata.name
    )
    eventsContent.value = response.data || []
  } catch (error) {
    eventsContent.value = []
    ElMessage.error('获取事件失败: ' + (error.response?.data?.error || error.message))
  }
}

const isPodRunning = (pod) => {
  return pod?.status?.phase === 'Running'
}

const openExec = (pod) => {
  const containers = pod?.spec?.containers || []
  if (containers.length === 0) return
  execPod.value = pod
  execContainers.value = containers.map(c => c.name)
  execContainer.value = containers[0].name
  showTerminal.value = true
}

const handleTerminalClosed = () => {
  execPod.value = null
  execContainer.value = ''
  execContainers.value = []
}

const reconnectTerminal = () => {
  terminalKey.value++
}

const handleDelete = async (resource) => {
  const id = clusterId.value
  if (!id) {
    ElMessage.error('集群ID无效')
    return
  }
  
  try {
    await ElMessageBox.confirm(
      `确定要删除 Pod ${resource.metadata.name} 吗？`,
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await k8sAPI.deletePod(
      id, 
      resource.metadata.namespace, 
      resource.metadata.name
    )
    
    await fetchResources()
    ElMessage.success('删除成功')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const fetchNamespaces = async () => {
  const id = clusterId.value
  if (!id) return
  
  try {
    const response = await k8sAPI.listNamespaces(id)
    namespaces.value = response.data
  } catch (error) {
    console.error('Failed to fetch namespaces:', error)
  }
}

const fetchResources = async () => {
  const id = clusterId.value
  if (!id) return
  
  try {
    const response = await k8sAPI.listPods(id, selectedNamespace.value)
    resources.value = response.data
  } catch (error) {
    console.error('Failed to fetch pods:', error)
    ElMessage.error('获取资源失败')
  }
}

const handleRefresh = async () => {
  await fetchResources()
  await fetchNamespaces()
  ElMessage.success('刷新成功')
}

const handleNamespaceChange = () => {
  currentPage.value = 1
  fetchResources()
}

watch(selectedNamespace, () => {
  currentPage.value = 1
  fetchResources()
})

watch(searchQuery, () => {
  currentPage.value = 1
})

onMounted(async () => {
  console.log('PodList onMounted, route.params:', route.params)
  
  const id = clusterId.value
  if (!id) {
    ElMessage.error('集群ID无效')
    router.push('/clusters')
    return
  }
  
  try {
    const response = await clusterAPI.get(id)
    cluster.value = response.data
    clusterName.value = response.data.alias || response.data.name
  } catch (error) {
    console.error('获取集群信息失败:', error)
  }
  
  await checkPermissions()
  if (hasReadPermission.value) {
    fetchNamespaces()
    fetchResources()
  }
})
</script>

<style scoped>
.resource-list-page {
  max-width: 1400px;
  margin: 0 auto;
  padding: 32px;
}

.page-container {
  animation: fadeIn 0.5s ease;
}

.page-header-hero {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  border-radius: 16px;
  padding: 32px;
  margin-bottom: 24px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.3);
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-title-section {
  display: flex;
  align-items: center;
  gap: 20px;
}

.back-btn {
  background: rgba(255, 255, 255, 0.8);
  border: 2px solid rgba(0, 0, 0, 0.2);
  color: black;
  width: auto;
  height: 48px;
  border-radius: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
  gap: 8px;
  padding: 0 16px;
  font-weight: 500;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.back-btn:hover {
  background: rgba(255, 255, 255, 1);
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.2);
}

.resource-icon {
  width: 80px;
  height: 80px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  backdrop-filter: blur(10px);
}

.title-text h1 {
  color: white;
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 8px;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

.cluster-info {
  display: flex;
  align-items: center;
  gap: 12px;
  color: rgba(255, 255, 255, 0.9);
  font-size: 14px;
}

.header-stats {
  display: flex;
  gap: 16px;
}

.stat-card {
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  padding: 20px 32px;
  text-align: center;
  border: 1px solid rgba(255, 255, 255, 0.2);
  transition: all 0.3s ease;
}

.stat-card:hover {
  background: rgba(255, 255, 255, 0.25);
  transform: translateY(-2px);
}

.stat-value {
  font-size: 32px;
  font-weight: 700;
  color: white;
  line-height: 1;
  margin-bottom: 8px;
}

.stat-label {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.8);
  text-transform: uppercase;
  letter-spacing: 1px;
}

.toolbar-section {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 16px;
  padding: 20px;
  margin-bottom: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.toolbar-left {
  display: flex;
  gap: 12px;
  flex: 1;
}

.namespace-select {
  width: 250px;
}

.search-input {
  flex: 1;
  max-width: 400px;
}

.toolbar-right {
  display: flex;
  gap: 12px;
}

.refresh-btn {
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  color: var(--text-secondary);
}

.resource-table-container {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 16px;
  overflow: hidden;
  animation: slideUp 0.5s ease;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.resource-name-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 500;
  color: var(--text-primary);
}

.clickable-name {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  color: var(--text-primary);
  font-weight: 500;
  transition: all 0.2s ease;
}

.clickable-name:hover {
  color: #f5576c;
  text-decoration: underline;
}

.resource-icon-small {
  color: #f5576c;
  font-size: 18px;
}

.pod-status {
  display: flex;
  align-items: center;
  gap: 8px;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  animation: pulse 2s infinite;
}

.status-dot.status-success {
  background: #67c23a;
  box-shadow: 0 0 8px #67c23a;
}

.status-dot.status-warning {
  background: #e6a23c;
  box-shadow: 0 0 8px #e6a23c;
}

.status-dot.status-error {
  background: #f56c6c;
  box-shadow: 0 0 8px #f56c6c;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

.action-buttons {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.action-btn {
  border-radius: 8px;
  font-size: 12px;
  padding: 4px 8px;
}

.log-btn {
  background: linear-gradient(135deg, #409eff, #337ecc);
  border: none;
  color: white;
}

.log-btn:hover {
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.4);
}

.event-btn {
  background: linear-gradient(135deg, #e6a23c, #f56c6c);
  border: none;
  color: white;
}

.event-btn:hover {
  box-shadow: 0 4px 12px rgba(245, 158, 11, 0.4);
}

.delete-btn {
  background: linear-gradient(135deg, #f56c6c, #e04040);
  border: none;
  color: white;
}

.delete-btn:hover {
  box-shadow: 0 4px 12px rgba(239, 68, 68, 0.4);
}

.pagination-container {
  padding: 20px;
  border-top: 1px solid var(--border-color);
  display: flex;
  justify-content: flex-end;
}

.empty-state {
  padding: 80px 20px;
  text-align: center;
}

.empty-state .el-icon {
  font-size: 80px;
  color: var(--text-tertiary);
  margin-bottom: 24px;
  opacity: 0.5;
}

.empty-state h3 {
  font-size: 20px;
  color: var(--text-primary);
  margin-bottom: 12px;
}

.logs-container {
  height: 500px;
  display: flex;
  flex-direction: column;
}

.logs-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--border-color);
  font-size: 14px;
  color: var(--text-secondary);
}

.logs-controls {
  display: flex;
  gap: 12px;
  align-items: center;
}

.container-select {
  width: 200px;
}

.logs-content {
  flex: 1;
  background: #0d1117;
  color: #c9d1d9;
  padding: 16px;
  overflow-y: auto;
  font-family: 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-all;
  border-radius: 8px;
  border: 1px solid var(--border-color);
}

.events-container {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.events-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid var(--border-color);
  font-size: 14px;
  color: var(--text-secondary);
}

.events-count {
  color: var(--text-tertiary);
  font-size: 12px;
  padding: 4px 12px;
  background: var(--bg-secondary);
  border-radius: 4px;
}

.no-events {
  text-align: center;
  padding: 60px 20px;
  color: var(--text-tertiary);
  font-size: 14px;
  background: var(--bg-secondary);
  border-radius: 8px;
  margin: 20px 0;
}

.no-permission {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  color: var(--text-tertiary);
}

.pod-ip {
  font-family: 'Courier New', monospace;
  font-size: 13px;
  font-weight: 600;
  color: #67c23a;
  background: #f0f9eb;
  padding: 2px 8px;
  border-radius: 4px;
}

.logs-path {
  font-weight: 500;
  color: var(--text-primary);
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.exec-btn {
  background: linear-gradient(135deg, #67c23a, #529b2e);
  border: none;
  color: white;
}

.exec-btn:hover {
  box-shadow: 0 4px 12px rgba(103, 194, 58, 0.4);
}

.terminal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding: 8px 16px;
  background: var(--bg-secondary);
  border-radius: 8px;
}

.terminal-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.terminal-path {
  font-weight: 600;
  color: var(--text-primary);
  font-size: 14px;
}

.terminal-container-badge {
  font-size: 12px;
  color: #67c23a;
  background: #f0f9eb;
  padding: 2px 10px;
  border-radius: 4px;
  font-weight: 500;
}

.logs-toolbar {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
}

.container-tag {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: var(--text-primary);
}

.normal-tag {
  font-size: 10px;
  font-weight: 600;
  padding: 1px 6px;
  border-radius: 3px;
  background: #e8f4fd;
  color: #3b82f6;
  text-transform: uppercase;
}

.init-tag {
  font-size: 10px;
  font-weight: 600;
  padding: 1px 6px;
  border-radius: 3px;
  background: #fef3c7;
  color: #d97706;
  text-transform: uppercase;
}
</style>
