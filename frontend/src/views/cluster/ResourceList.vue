<template>
  <div class="resource-list-page">
    <div v-if="!hasReadPermission" class="no-permission">
      <el-icon name="Lock" size="64" color="#909399" />
      <p>您没有权限访问该集群</p>
    </div>
    
    <div v-else class="page-container">
      <!-- 页面头部信息 -->
      <div class="page-header-hero">
        <div class="header-content">
          <div class="header-title-section">
            <div class="resource-icon">
              <el-icon :size="48"><component :is="getResourceIcon" /></el-icon>
            </div>
            <div class="title-text">
              <h1>{{ resourceTypeLabel }} 资源管理</h1>
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
              <div class="stat-label">资源总数</div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 工具栏 -->
      <div class="toolbar-section">
        <div class="toolbar-left">
          <el-select 
            v-model="selectedNamespace" 
            placeholder="选择命名空间"
            class="namespace-select"
            clearable
          >
            <el-option label="所有命名空间" value="" />
            <el-option 
              v-for="ns in namespaces" 
              :key="ns.name" 
              :label="ns.name" 
              :value="ns.name" 
            />
          </el-select>
          
          <el-input
            v-model="searchQuery"
            placeholder="搜索资源名称..."
            class="search-input"
            clearable
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </div>
        
        <div class="toolbar-right">
          <el-button 
            type="primary" 
            @click="showCreateModal = true" 
            v-if="hasWritePermission"
            class="create-btn"
          >
            <el-icon><Plus /></el-icon>
            创建资源
          </el-button>
          <el-button @click="fetchResources" class="refresh-btn">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </div>
      
      <!-- 资源表格 -->
      <div class="resource-table-container">
        <el-table 
          :data="paginatedResources" 
          border 
          stripe
          class="resource-table"
          :header-cell-style="{ background: 'var(--bg-secondary)', color: 'var(--text-primary)', textAlign: 'center' }"
        >
          <template v-if="resourceType === 'deployments'">
            <el-table-column prop="metadata.name" label="名称" min-width="200">
              <template #default="scope">
                <div class="resource-name-cell">
                  <el-icon class="resource-icon-small"><Box /></el-icon>
                  <span>{{ scope.row.metadata.name }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="metadata.namespace" label="命名空间" width="150" />
            <el-table-column label="副本数" width="140">
              <template #default="scope">
                <div class="replicas-info">
                  <span class="ready">{{ scope.row.status.readyReplicas || 0 }}</span>
                  <span class="separator">/</span>
                  <span class="total">{{ scope.row.status.replicas || 0 }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="状态" width="120">
              <template #default="scope">
                <el-tag :type="getDeploymentStatusTagType(scope.row)" size="small">
                  {{ getDeploymentStatus(scope.row) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="镜像" min-width="250">
              <template #default="scope">
                <el-tooltip 
                  v-if="scope.row.spec.template.spec.containers[0]" 
                  :content="scope.row.spec.template.spec.containers[0].image"
                  placement="top"
                >
                  <span class="image-text">
                    {{ truncateImage(scope.row.spec.template.spec.containers[0].image) }}
                  </span>
                </el-tooltip>
              </template>
            </el-table-column>
            <el-table-column prop="metadata.creationTimestamp" label="创建时间" width="180">
              <template #default="scope">
                {{ formatTime(scope.row.metadata.creationTimestamp) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="280" fixed="right" :cell-style="{ textAlign: 'center' }">
              <template #default="scope">
                <div class="action-buttons">
                  <el-button size="small" @click="viewResource(scope.row)" class="action-btn">
                    <el-icon><View /></el-icon>
                    详情
                  </el-button>
                  <el-button size="small" @click="viewEvents(scope.row)" class="action-btn">
                    <el-icon><Document /></el-icon>
                    事件
                  </el-button>
                  <el-button size="small" @click="handleScale(scope.row)" v-if="hasWritePermission" class="action-btn scale-btn">
                    <el-icon><Setting /></el-icon>
                    扩缩
                  </el-button>
                  <el-button size="small" danger @click="handleDelete(scope.row)" v-if="hasWritePermission" class="action-btn delete-btn">
                    <el-icon><Delete /></el-icon>
                  </el-button>
                </div>
              </template>
            </el-table-column>
          </template>
        
        <template v-else-if="resourceType === 'statefulsets'">
          <el-table-column prop="metadata.name" label="名称" />
          <el-table-column prop="metadata.namespace" label="命名空间" />
          <el-table-column label="副本数">
            <template #default="scope">
              {{ scope.row.spec.replicas || 0 }} / {{ scope.row.status.replicas || 0 }}
            </template>
          </el-table-column>
          <el-table-column prop="metadata.creationTimestamp" label="创建时间" />
          <el-table-column label="操作" :cell-style="{ textAlign: 'center' }">
            <template #default="scope">
              <el-button size="small" @click="viewResource(scope.row)">查看</el-button>
              <el-button size="small" danger @click="handleDelete(scope.row)" v-if="hasWritePermission">删除</el-button>
            </template>
          </el-table-column>
        </template>
        
        <template v-else-if="resourceType === 'services'">
          <el-table-column prop="metadata.name" label="名称" />
          <el-table-column prop="metadata.namespace" label="命名空间" />
          <el-table-column label="类型">
            <template #default="scope">{{ scope.row.spec.type }}</template>
          </el-table-column>
          <el-table-column label="端口">
            <template #default="scope">{{ formatPorts(scope.row.spec.ports) }}</template>
          </el-table-column>
          <el-table-column prop="metadata.creationTimestamp" label="创建时间" />
          <el-table-column label="操作" :cell-style="{ textAlign: 'center' }">
            <template #default="scope">
              <el-button size="small" @click="viewResource(scope.row)">查看</el-button>
              <el-button size="small" danger @click="handleDelete(scope.row)" v-if="hasWritePermission">删除</el-button>
            </template>
          </el-table-column>
        </template>
        
        <template v-else-if="resourceType === 'ingresses'">
          <el-table-column prop="metadata.name" label="名称" />
          <el-table-column prop="metadata.namespace" label="命名空间" />
          <el-table-column label="主机">
            <template #default="scope">{{ formatHosts(scope.row.spec.rules) }}</template>
          </el-table-column>
          <el-table-column prop="metadata.creationTimestamp" label="创建时间" />
          <el-table-column label="操作" :cell-style="{ textAlign: 'center' }">
            <template #default="scope">
              <el-button size="small" @click="viewResource(scope.row)">查看</el-button>
              <el-button size="small" danger @click="handleDelete(scope.row)" v-if="hasWritePermission">删除</el-button>
            </template>
          </el-table-column>
        </template>
        
        <template v-else-if="resourceType === 'configmaps'">
          <el-table-column prop="metadata.name" label="名称" />
          <el-table-column prop="metadata.namespace" label="命名空间" />
          <el-table-column label="数据项数">
            <template #default="scope">{{ Object.keys(scope.row.data || {}).length }}</template>
          </el-table-column>
          <el-table-column prop="metadata.creationTimestamp" label="创建时间" />
          <el-table-column label="操作" :cell-style="{ textAlign: 'center' }">
            <template #default="scope">
              <el-button size="small" @click="viewResource(scope.row)">查看</el-button>
              <el-button size="small" danger @click="handleDelete(scope.row)" v-if="hasWritePermission">删除</el-button>
            </template>
          </el-table-column>
        </template>
        
        <template v-else-if="resourceType === 'secrets'">
          <el-table-column prop="metadata.name" label="名称" />
          <el-table-column prop="metadata.namespace" label="命名空间" />
          <el-table-column prop="type" label="类型" />
          <el-table-column prop="metadata.creationTimestamp" label="创建时间" />
          <el-table-column label="操作" :cell-style="{ textAlign: 'center' }">
            <template #default="scope">
              <el-button size="small" @click="viewResource(scope.row)">查看</el-button>
              <el-button size="small" danger @click="handleDelete(scope.row)" v-if="hasWritePermission">删除</el-button>
            </template>
          </el-table-column>
        </template>
        
        <template v-else-if="resourceType === 'pods'">
          <el-table-column prop="metadata.name" label="名称" />
          <el-table-column prop="metadata.namespace" label="命名空间" />
          <el-table-column label="状态">
            <template #default="scope">
              <span :class="getPodStatusClass(scope.row)">
                {{ getPodStatus(scope.row) }}
              </span>
            </template>
          </el-table-column>
          <el-table-column label="节点">
            <template #default="scope">{{ scope.row.spec.nodeName || '-' }}</template>
          </el-table-column>
          <el-table-column prop="metadata.creationTimestamp" label="创建时间" />
          <el-table-column label="操作" :cell-style="{ textAlign: 'center' }">
            <template #default="scope">
              <el-button size="small" @click="viewResource(scope.row)">查看</el-button>
              <el-button size="small" @click="viewLogs(scope.row)">查看日志</el-button>
              <el-button size="small" @click="viewEvents(scope.row)">查看事件</el-button>
              <el-button size="small" danger @click="handleDelete(scope.row)" v-if="hasWritePermission">删除</el-button>
            </template>
          </el-table-column>
        </template>
      </el-table>
    </div>
    
    <el-dialog title="创建资源" v-model="showCreateModal" width="800px">
      <el-form :model="createForm">
        <el-form-item label="名称">
          <el-input v-model="createForm.name" placeholder="请输入名称" />
        </el-form-item>
        <el-form-item label="命名空间">
          <el-select v-model="createForm.namespace">
            <el-option v-for="ns in namespaces" :key="ns.name" :label="ns.name" :value="ns.name" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateModal = false">取消</el-button>
        <el-button type="primary" @click="handleCreate">确定</el-button>
      </template>
    </el-dialog>
    
    <el-dialog title="查看日志" v-model="showLogsModal" width="800px" :close-on-click-modal="false">
      <div class="logs-container">
        <div class="logs-header">
          <span>{{ currentPod?.metadata.namespace }} / {{ currentPod?.metadata.name }}</span>
          <el-select v-model="selectedContainer" placeholder="选择容器">
            <el-option 
              v-for="container in podContainers" 
              :key="container.name" 
              :label="container.name" 
              :value="container.name" 
            />
          </el-select>
        </div>
        <pre class="logs-content">{{ logsContent }}</pre>
      </div>
      <template #footer>
        <el-button @click="showLogsModal = false">关闭</el-button>
        <el-button type="primary" @click="fetchLogs">刷新</el-button>
      </template>
    </el-dialog>
    
    <el-dialog title="缩扩容" v-model="showScaleModal" width="400px">
      <el-form :model="scaleForm">
        <el-form-item label="副本数">
          <el-input-number v-model="scaleForm.replicas" :min="0" :max="100" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showScaleModal = false">取消</el-button>
        <el-button type="primary" @click="handleScaleConfirm">确定</el-button>
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
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useAuthStore } from '@/stores/auth'
import { k8sAPI, authAPI, clusterAPI } from '@/utils/api'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const resources = ref([])
const namespaces = ref([])
const selectedNamespace = ref('')
const showCreateModal = ref(false)
const showLogsModal = ref(false)
const showScaleModal = ref(false)
const showEventsModal = ref(false)
const logsContent = ref('')
const eventsContent = ref([])
const currentPod = ref(null)
const currentResource = ref(null)
const selectedContainer = ref('')
const scalingResource = ref(null)
const searchQuery = ref('')
const currentPage = ref(1)
const pageSize = ref(20)
const cluster = ref(null)

const createForm = reactive({
  name: '',
  namespace: ''
})

const scaleForm = reactive({
  replicas: 1
})

const resourceType = computed(() => route.params.type)
const clusterId = computed(() => route.params.id)
const clusterName = ref('')

// 通用函数：获取有效的集群 ID
const getValidClusterId = () => {
  const id = route.params.id
  if (!id) return null
  
  const numId = Number(id)
  if (isNaN(numId) || numId <= 0) return null
  
  return numId
}

const resourceTypeLabel = computed(() => {
  const labels = {
    deployments: 'Deployment',
    statefulsets: 'StatefulSet',
    services: 'Service',
    ingresses: 'Ingress',
    configmaps: 'ConfigMap',
    secrets: 'Secret',
    pods: 'Pod'
  }
  return labels[resourceType.value] || resourceType.value
})

const hasReadPermission = ref(false)
const hasWritePermission = ref(false)

// 计算属性：过滤后的资源
const filteredResources = computed(() => {
  if (!searchQuery.value) {
    return resources.value
  }
  const query = searchQuery.value.toLowerCase()
  return resources.value.filter(resource => 
    resource.metadata?.name?.toLowerCase().includes(query)
  )
})

// 计算属性：分页后的资源
const paginatedResources = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredResources.value.slice(start, end)
})

// 计算属性：根据资源类型返回图标
const getResourceIcon = computed(() => {
  const icons = {
    deployments: 'Monitor',
    statefulsets: 'Box',
    services: 'Connection',
    ingresses: 'Link',
    configmaps: 'Document',
    secrets: 'Lock',
    pods: 'Box'
  }
  return icons[resourceType.value] || 'Box'
})

const checkPermissions = async () => {
  let id = clusterId.value
  if (!id) {
    console.error('集群ID为空，无法检查权限')
    hasReadPermission.value = false
    hasWritePermission.value = false
    return
  }
  
  id = Number(id)
  if (isNaN(id) || id <= 0) {
    console.error('集群ID无效，无法检查权限:', clusterId.value)
    hasReadPermission.value = false
    hasWritePermission.value = false
    return
  }
  
  try {
    console.log('检查集群权限, clusterId:', id)
    const response = await authAPI.checkClusterPermission(id)
    hasReadPermission.value = response.data.has_read
    hasWritePermission.value = response.data.has_write
  } catch (error) {
    console.error('检查权限失败:', error)
    hasReadPermission.value = false
    hasWritePermission.value = false
  }
}

const podContainers = computed(() => {
  if (!currentPod.value?.spec?.containers) return []
  return currentPod.value.spec.containers
})

const getDeploymentStatus = (row) => {
  const ready = row.status.readyReplicas || 0
  const total = row.status.replicas || 0
  if (ready === total && total > 0) return '运行中'
  if (row.status.updatedReplicas > 0 && ready < row.status.updatedReplicas) return '更新中'
  return '异常'
}

const getDeploymentStatusClass = (row) => {
  const status = getDeploymentStatus(row)
  if (status === '运行中') return 'status-success'
  if (status === '更新中') return 'status-warning'
  return 'status-error'
}

const getPodStatus = (row) => {
  const phase = row.status.phase
  const containerStatuses = row.status.containerStatuses || []
  const totalContainers = containerStatuses.length
  const runningContainers = containerStatuses.filter(c => c.state.running).length
  
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

const formatPorts = (ports) => {
  if (!ports) return '-'
  return ports.map(p => `${p.port}/${p.protocol || 'TCP'}`).join(', ')
}

const formatHosts = (rules) => {
  if (!rules) return '-'
  return rules.map(r => r.host).join(', ')
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

const truncateImage = (image) => {
  if (!image) return '-'
  if (image.length > 60) {
    return image.substring(0, 60) + '...'
  }
  return image
}

const getDeploymentStatusTagType = (row) => {
  const ready = row.status.readyReplicas || 0
  const total = row.status.replicas || 0
  if (ready === total && total > 0) return 'success'
  if (row.status.updatedReplicas > 0 && ready < row.status.updatedReplicas) return 'warning'
  return 'danger'
}

const viewResource = (resource) => {
  const id = getValidClusterId()
  if (!id) {
    ElMessage.error('集群ID无效')
    return
  }
  
  const namespace = resource.metadata.namespace
  const name = resource.metadata.name
  window.open(`/clusters/${id}/resources/${resourceType.value}/${namespace}/${name}`, '_blank')
}

const viewLogs = (pod) => {
  currentPod.value = pod
  if (pod.spec.containers?.length > 0) {
    selectedContainer.value = pod.spec.containers[0].name
  }
  showLogsModal.value = true
  fetchLogs()
}

const fetchLogs = async () => {
  const id = getValidClusterId()
  if (!id) {
    logsContent.value = '获取日志失败：集群ID无效'
    return
  }
  
  try {
    const response = await k8sAPI.getPodLogs(id, currentPod.value.metadata.namespace, currentPod.value.metadata.name, selectedContainer.value)
    logsContent.value = response.data.logs
  } catch (error) {
    logsContent.value = '获取日志失败: ' + (error.response?.data?.error || error.message)
  }
}

const viewEvents = async (resource) => {
  currentResource.value = resource
  showEventsModal.value = true
  await fetchEvents()
}

const fetchEvents = async () => {
  const id = getValidClusterId()
  if (!id) {
    eventsContent.value = []
    ElMessage.error('集群ID无效')
    return
  }
  
  try {
    let response
    if (resourceType.value === 'pods') {
      response = await k8sAPI.getPodEvents(id, currentResource.value.metadata.namespace, currentResource.value.metadata.name)
    } else if (resourceType.value === 'deployments') {
      response = await k8sAPI.getDeploymentEvents(id, currentResource.value.metadata.namespace, currentResource.value.metadata.name)
    }
    eventsContent.value = response.data || []
  } catch (error) {
    eventsContent.value = []
    ElMessage.error('获取事件失败: ' + (error.response?.data?.error || error.message))
  }
}

const handleScale = (resource) => {
  scalingResource.value = resource
  scaleForm.replicas = resource.spec.replicas || 1
  showScaleModal.value = true
}

const handleScaleConfirm = async () => {
  const id = getValidClusterId()
  if (!id) {
    ElMessage.error('集群ID无效')
    return
  }
  
  try {
    await k8sAPI.scaleDeployment(id, scalingResource.value.metadata.namespace, scalingResource.value.metadata.name, scaleForm.replicas)
    showScaleModal.value = false
    await fetchResources()
    ElMessage.success('缩扩容成功')
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '缩扩容失败')
  }
}

const handleDelete = async (resource) => {
  const id = getValidClusterId()
  if (!id) {
    ElMessage.error('集群ID无效')
    return
  }
  
  try {
    await ElMessageBox.confirm(
      `确定要删除 ${resource.metadata.name} 吗？`,
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    const deleteFn = {
      deployments: k8sAPI.deleteDeployment,
      statefulsets: () => Promise.resolve({ data: { message: 'deleted' } }),
      services: k8sAPI.deleteService,
      ingresses: k8sAPI.deleteIngress,
      configmaps: k8sAPI.deleteConfigMap,
      secrets: k8sAPI.deleteSecret,
      pods: k8sAPI.deletePod
    }[resourceType.value]
    
    if (deleteFn) {
      await deleteFn(id, resource.metadata.namespace, resource.metadata.name)
    }
    
    await fetchResources()
    ElMessage.success('删除成功')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleCreate = async () => {
  const id = getValidClusterId()
  if (!id) {
    ElMessage.error('集群ID无效')
    return
  }
  
  try {
    await k8sAPI[`create${resourceTypeLabel.value}`](id, createForm)
    showCreateModal.value = false
    createForm.name = ''
    createForm.namespace = ''
    await fetchResources()
    ElMessage.success('创建成功')
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '创建失败')
  }
}

const fetchNamespaces = async () => {
  const id = getValidClusterId()
  if (!id) {
    console.error('无法获取命名空间：集群ID无效')
    return
  }
  
  try {
    console.log('获取命名空间, clusterId:', id)
    const response = await k8sAPI.listNamespaces(id)
    namespaces.value = response.data
  } catch (error) {
    console.error('Failed to fetch namespaces:', error)
  }
}

const fetchResources = async () => {
  const id = getValidClusterId()
  if (!id) {
    console.error('无法获取资源：集群ID无效')
    return
  }
  
  try {
    console.log('获取资源, clusterId:', id, 'resourceType:', resourceType.value, 'namespace:', selectedNamespace.value)
    const listFn = {
      deployments: k8sAPI.listDeployments,
      statefulsets: k8sAPI.listStatefulSets,
      services: k8sAPI.listServices,
      ingresses: k8sAPI.listIngresses,
      configmaps: k8sAPI.listConfigMaps,
      secrets: k8sAPI.listSecrets,
      pods: k8sAPI.listPods
    }[resourceType.value]
    
    if (listFn) {
      const response = await listFn(id, selectedNamespace.value)
      resources.value = response.data
    }
  } catch (error) {
    console.error('Failed to fetch resources:', error)
  }
}

watch(selectedNamespace, () => {
  currentPage.value = 1
  fetchResources()
})

watch(searchQuery, () => {
  currentPage.value = 1
})

onMounted(async () => {
  console.log('ResourceList onMounted, route.params:', route.params)
  
  let id = route.params.id
  if (!id) {
    console.error('集群ID为空')
    ElMessage.error('集群ID无效')
    router.push('/clusters')
    return
  }
  
  id = Number(id)
  if (isNaN(id) || id <= 0) {
    console.error('集群ID无效:', route.params.id, '转换后:', id)
    ElMessage.error('集群ID无效')
    router.push('/clusters')
    return
  }
  
  console.log('资源列表页面加载, clusterId:', id, 'resourceType:', route.params.type)
  
  // 获取集群信息
  try {
    const response = await clusterAPI.get(id)
    cluster.value = response.data
    clusterName.value = response.data.alias || response.data.name
    console.log('获取到集群信息:', cluster.value)
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

.no-permission {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  color: var(--text-tertiary);
}

.no-permission .el-icon {
  font-size: 64px;
  margin-bottom: 24px;
  opacity: 0.5;
}

.no-permission p {
  margin-top: 16px;
  font-size: 16px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
  padding: 24px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  animation: slideDown 0.4s ease;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.page-header h1 {
  font-size: 28px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 8px;
  background: linear-gradient(135deg, var(--primary-light), var(--success-light));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.page-header p {
  color: var(--text-tertiary);
  font-size: 14px;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.header-actions .el-select {
  min-width: 200px;
}

.resource-table {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  overflow: hidden;
  animation: fadeIn 0.5s ease;
}

.status-success {
  color: var(--success-color);
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 6px;
}

.status-warning {
  color: var(--warning-color);
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 6px;
}

.status-error {
  color: var(--error-color);
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 6px;
}

.logs-container {
  height: 400px;
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

.logs-content {
  flex: 1;
  background: #0d1117;
  color: #c9d1d9;
  padding: 16px;
  overflow-y: auto;
  font-family: var(--font-mono);
  font-size: 13px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-all;
  border-radius: var(--radius-md);
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
  border-radius: var(--radius-sm);
}

.no-events {
  text-align: center;
  padding: 60px 20px;
  color: var(--text-tertiary);
  font-size: 14px;
  background: var(--bg-secondary);
  border-radius: var(--radius-md);
  margin: 20px 0;
}

/* Element Plus overrides for dark theme */
:deep(.el-table) {
  background: transparent;
  color: var(--text-primary);
}

:deep(.el-table th) {
  background: var(--bg-secondary);
  color: var(--text-secondary);
  font-weight: 600;
  border-bottom: 1px solid var(--border-color);
}

:deep(.el-table tr) {
  background: transparent;
  transition: background var(--transition-fast);
}

:deep(.el-table tr:hover > td) {
  background: var(--bg-hover);
}

:deep(.el-table td) {
  border-bottom: 1px solid var(--border-color);
  color: var(--text-primary);
}

:deep(.el-button) {
  border-radius: var(--radius-md);
  font-weight: 500;
  transition: all var(--transition-fast);
}

:deep(.el-button--primary) {
  background: linear-gradient(135deg, var(--primary-color), var(--primary-dark));
  border: none;
}

:deep(.el-button--primary:hover) {
  background: linear-gradient(135deg, var(--primary-light), var(--primary-color));
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.4);
}

:deep(.el-button--danger) {
  background: linear-gradient(135deg, var(--error-color), var(--error-dark));
  border: none;
}

:deep(.el-button--danger:hover) {
  background: linear-gradient(135deg, var(--error-light), var(--error-color));
  box-shadow: 0 4px 12px rgba(239, 68, 68, 0.4);
}

:deep(.el-input__wrapper) {
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  box-shadow: none;
}

:deep(.el-input__wrapper:hover) {
  border-color: var(--border-light);
}

:deep(.el-input__wrapper.is-focus) {
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px var(--primary-glow);
}

:deep(.el-select .el-input__wrapper) {
  background: var(--bg-secondary);
}

:deep(.el-dialog) {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-xl);
}

:deep(.el-dialog__header) {
  border-bottom: 1px solid var(--border-color);
  padding: 20px 24px;
}

:deep(.el-dialog__title) {
  color: var(--text-primary);
  font-weight: 600;
}

:deep(.el-dialog__body) {
  padding: 24px;
  color: var(--text-secondary);
}

:deep(.el-dialog__footer) {
  border-top: 1px solid var(--border-color);
  padding: 16px 24px;
}

/* 新的界面样式 */
.page-container {
  animation: fadeIn 0.5s ease;
}

.page-header-hero {
  background: linear-gradient(135deg, var(--primary-color) 0%, var(--primary-dark) 100%);
  border-radius: var(--radius-xl);
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

.resource-icon {
  width: 80px;
  height: 80px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: var(--radius-lg);
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

.cluster-info .el-icon {
  font-size: 16px;
}

.header-stats {
  display: flex;
  gap: 16px;
}

.stat-card {
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(10px);
  border-radius: var(--radius-lg);
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
  border-radius: var(--radius-lg);
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

.create-btn {
  background: linear-gradient(135deg, var(--primary-color), var(--primary-dark));
  border: none;
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
}

.create-btn:hover {
  box-shadow: 0 6px 16px rgba(59, 130, 246, 0.4);
  transform: translateY(-1px);
}

.refresh-btn {
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  color: var(--text-secondary);
}

.refresh-btn:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.resource-table-container {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
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

.resource-table {
  border-radius: var(--radius-lg);
}

.resource-name-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 500;
  color: var(--text-primary);
}

.resource-icon-small {
  color: var(--primary-color);
  font-size: 18px;
}

.replicas-info {
  display: flex;
  align-items: center;
  gap: 4px;
  font-weight: 600;
  font-family: 'JetBrains Mono', 'Courier New', monospace;
}

.replicas-info .ready {
  color: var(--success-color);
  font-size: 16px;
}

.replicas-info .separator {
  color: var(--text-tertiary);
}

.replicas-info .total {
  color: var(--text-secondary);
  font-size: 16px;
}

.image-text {
  font-family: 'JetBrains Mono', 'Courier New', monospace;
  font-size: 12px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: color 0.2s;
}

.image-text:hover {
  color: var(--primary-color);
}

.action-buttons {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.action-btn {
  border-radius: var(--radius-md);
  font-size: 12px;
  padding: 4px 8px;
}

.action-btn :deep(.el-icon) {
  margin-right: 4px;
}

.scale-btn {
  background: linear-gradient(135deg, var(--warning-color), var(--warning-dark));
  border: none;
  color: white;
}

.scale-btn:hover {
  box-shadow: 0 4px 12px rgba(245, 158, 11, 0.4);
}

.delete-btn {
  background: linear-gradient(135deg, var(--error-color), var(--error-dark));
  border: none;
  color: white;
}

.delete-btn:hover {
  box-shadow: 0 4px 12px rgba(239, 68, 68, 0.4);
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
  background: var(--success-color);
  box-shadow: 0 0 8px var(--success-color);
}

.status-dot.status-warning {
  background: var(--warning-color);
  box-shadow: 0 0 8px var(--warning-color);
}

.status-dot.status-error {
  background: var(--error-color);
  box-shadow: 0 0 8px var(--error-color);
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

.pagination-container {
  padding: 20px;
  border-top: 1px solid var(--border-color);
  display: flex;
  justify-content: flex-end;
}

.custom-pagination {
  margin: 0;
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

.empty-state p {
  color: var(--text-tertiary);
  margin-bottom: 24px;
}
</style>