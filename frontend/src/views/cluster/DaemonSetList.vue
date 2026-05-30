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
            <el-button @click="goBack" class="back-btn">
              <el-icon><ArrowLeft /></el-icon>
              <span>返回</span>
            </el-button>
            <div class="resource-icon">
              <IconDaemonSet :size="48" />
            </div>
            <div class="title-text">
              <h1>DaemonSet 资源管理</h1>
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
              <div class="stat-label">DaemonSet 总数</div>
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
            placeholder="搜索 DaemonSet 名称..."
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
            创建 DaemonSet
          </el-button>
          <el-button @click="handleRefresh" class="refresh-btn">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </div>
      
      <!-- 资源列表 -->
      <div class="resource-table-container">
        <el-table 
          :data="paginatedResources" 
          border 
          stripe
          class="resource-table"
          :header-cell-style="{ background: 'var(--bg-secondary)', color: 'var(--text-primary)', textAlign: 'center' }"
        >
          <el-table-column prop="metadata.name" label="名称" min-width="200">
            <template #default="scope">
              <span class="clickable-name" @click="viewResource(scope.row)">
                <el-icon class="resource-icon-small"><Collection /></el-icon>
                {{ scope.row.metadata.name }}
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="metadata.namespace" label="命名空间" width="150" />
          <el-table-column label="状态" width="120">
            <template #default="scope">
              <el-tag :type="getDaemonSetStatusTagType(scope.row)" size="small">
                {{ getDaemonSetStatus(scope.row) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="就绪" width="120">
            <template #default="scope">
              <span class="replicas-info">
                <span class="ready">{{ scope.row.status.numberReady || 0 }}</span>
                <span class="separator">/</span>
                <span class="total">{{ scope.row.status.desiredNumberScheduled || 0 }}</span>
              </span>
            </template>
          </el-table-column>
          <el-table-column label="镜像" min-width="250">
            <template #default="scope">
              <span v-if="scope.row?.spec?.template?.spec?.containers?.length > 0">
                <el-tooltip :content="scope.row.spec.template.spec.containers[0].image" placement="top">
                  <span class="image-text">{{ truncateImage(scope.row.spec.template.spec.containers[0].image) }}</span>
                </el-tooltip>
              </span>
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column prop="metadata.creationTimestamp" label="创建时间" width="180">
            <template #default="scope">
              {{ formatTime(scope.row.metadata.creationTimestamp) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="240" fixed="right" :cell-style="{ textAlign: 'center' }">
            <template #default="scope">
              <div class="action-buttons">
                <el-button size="small" @click="viewEvents(scope.row)" class="action-btn">
                  <el-icon><Document /></el-icon>
                  事件
                </el-button>
                <el-button size="small" @click="handleEdit(scope.row)" v-if="hasWritePermission" class="action-btn edit-btn">
                  <el-icon><Edit /></el-icon>
                  编辑
                </el-button>
                <el-button size="small" danger @click="handleDelete(scope.row)" v-if="hasWritePermission" class="action-btn delete-btn">
                  <el-icon><Delete /></el-icon>
                  删除
                </el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>
        
        <!-- 分页 -->
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
        
        <!-- 空状态 -->
        <div v-if="filteredResources.length === 0" class="empty-state">
          <el-icon :size="80" color="var(--text-tertiary)"><DocumentDelete /></el-icon>
          <h3>暂无 DaemonSet</h3>
          <p v-if="searchQuery">未找到匹配 "{{ searchQuery }}" 的 DaemonSet</p>
          <p v-else>该命名空间下暂无 DaemonSet 资源</p>
          <el-button 
            type="primary" 
            @click="showCreateModal = true" 
            v-if="hasWritePermission && !searchQuery"
          >
            创建第一个 DaemonSet
          </el-button>
        </div>
      </div>
    </div>
    
    <!-- 创建 DaemonSet 对话框 -->
    <el-dialog title="创建 DaemonSet (YAML)" v-model="showCreateModal" width="900px" :close-on-click-modal="false">
      <div class="yaml-editor-container">
        <el-input
          v-model="createYaml"
          type="textarea"
          :rows="25"
          placeholder="apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: my-daemonset
  namespace: default
spec:
  selector:
    matchLabels:
      app: my-app
  template:
    metadata:
      labels:
        app: my-app
    spec:
      containers:
      - name: nginx
        image: nginx:latest
        ports:
        - containerPort: 80"
          class="yaml-editor"
        />
      </div>
      <template #footer>
        <el-button @click="showCreateModal = false">取消</el-button>
        <el-button type="primary" @click="handleCreate" :loading="loadingYaml">创建</el-button>
      </template>
    </el-dialog>
    
    <!-- 编辑对话框 -->
    <el-dialog title="编辑 DaemonSet (YAML)" v-model="showEditModal" width="900px" :close-on-click-modal="false">
      <div class="yaml-editor-container">
        <el-input
          v-model="editYaml"
          type="textarea"
          :rows="25"
          class="yaml-editor"
          placeholder="请输入 YAML 内容"
        />
      </div>
      <template #footer>
        <el-button @click="showEditModal = false">取消</el-button>
        <el-button type="primary" @click="handleUpdate" :loading="updatingYaml">保存</el-button>
      </template>
    </el-dialog>
    
    <!-- 查看事件对话框 -->
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
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import IconDaemonSet from '@/assets/icons/IconDaemonSet.vue'
import { k8sAPI, authAPI, clusterAPI } from '@/utils/api'
import jsYaml from 'js-yaml'

const route = useRoute()
const router = useRouter()

const resources = ref([])
const namespaces = ref([])
const selectedNamespace = ref('default')
const showCreateModal = ref(false)
const showEventsModal = ref(false)
const showEditModal = ref(false)
const eventsContent = ref([])
const currentResource = ref(null)
const searchQuery = ref('')
const currentPage = ref(1)
const pageSize = ref(20)
const cluster = ref(null)
const createYaml = ref('')
const editYaml = ref('')
const loadingYaml = ref(false)
const updatingYaml = ref(false)

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

const getDaemonSetStatus = (row) => {
  const ready = row.status.numberReady || 0
  const total = row.status.desiredNumberScheduled || 0
  if (ready === total && total > 0) return '运行中'
  if (row.status.updatedNumberScheduled > 0 && ready < row.status.updatedNumberScheduled) return '更新中'
  return '异常'
}

const getDaemonSetStatusTagType = (row) => {
  const ready = row.status.numberReady || 0
  const total = row.status.desiredNumberScheduled || 0
  if (ready === total && total > 0) return 'success'
  if (row.status.updatedNumberScheduled > 0 && ready < row.status.updatedNumberScheduled) return 'warning'
  return 'danger'
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
  router.push(`/clusters/${id}/daemonsets/${namespace}/${name}`)
}

const handleEdit = async (resource) => {
  currentResource.value = resource
  try {
    const id = clusterId.value
    const response = await k8sAPI.getDaemonSet(id, resource.metadata.namespace, resource.metadata.name)
    try {
      editYaml.value = jsYaml.dump(response.data)
    } catch {
      editYaml.value = JSON.stringify(response.data, null, 2)
    }
    showEditModal.value = true
  } catch (error) {
    ElMessage.error('获取资源信息失败')
  }
}

const handleUpdate = async () => {
  if (!editYaml.value.trim()) {
    ElMessage.error('请输入 YAML 内容')
    return
  }
  
  updatingYaml.value = true
  try {
    const updatedResource = jsYaml.load(editYaml.value)
    const id = clusterId.value
    await k8sAPI.updateDaemonSet(id, currentResource.value.metadata.namespace, currentResource.value.metadata.name, updatedResource)
    showEditModal.value = false
    await fetchResources()
    ElMessage.success('更新成功')
  } catch (error) {
    ElMessage.error(error.response?.data?.error || error.message || '更新失败')
  } finally {
    updatingYaml.value = false
  }
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
    const response = await k8sAPI.getDaemonSetEvents(
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

const handleDelete = async (resource) => {
  const id = clusterId.value
  if (!id) {
    ElMessage.error('集群ID无效')
    return
  }
  
  try {
    await ElMessageBox.confirm(
      `确定要删除 DaemonSet ${resource.metadata.name} 吗？`,
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await k8sAPI.deleteDaemonSet(
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

const handleCreate = async () => {
  const id = clusterId.value
  if (!id) {
    ElMessage.error('集群ID无效')
    return
  }
  
  if (!createYaml.value.trim()) {
    ElMessage.error('请输入 YAML 内容')
    return
  }
  
  loadingYaml.value = true
  try {
    const daemonSet = jsYaml.load(createYaml.value)
    
    await k8sAPI.createDaemonSet(id, daemonSet)
    showCreateModal.value = false
    createYaml.value = ''
    await fetchResources()
    ElMessage.success('创建成功')
  } catch (error) {
    ElMessage.error(error.response?.data?.error || error.message || '创建失败')
  } finally {
    loadingYaml.value = false
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
    const response = await k8sAPI.listDaemonSets(id, selectedNamespace.value)
    resources.value = response.data
  } catch (error) {
    console.error('Failed to fetch daemonsets:', error)
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
  console.log('DaemonSetList onMounted, route.params:', route.params)
  
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

.create-btn {
  background: linear-gradient(135deg, #f093fb, #f5576c);
  border: none;
  box-shadow: 0 4px 12px rgba(240, 147, 251, 0.3);
}

.create-btn:hover {
  box-shadow: 0 6px 16px rgba(240, 147, 251, 0.4);
  transform: translateY(-1px);
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

.image-text {
  font-family: 'Courier New', monospace;
  font-size: 13px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: color 0.2s;
}

.image-text:hover {
  color: #f5576c;
}

.replicas-info {
  display: flex;
  align-items: center;
  gap: 4px;
  font-weight: 600;
  font-family: 'Courier New', monospace;
}

.replicas-info .ready {
  color: #67c23a;
  font-size: 16px;
}

.replicas-info .separator {
  color: var(--text-tertiary);
}

.replicas-info .total {
  color: var(--text-secondary);
  font-size: 16px;
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

.edit-btn {
  background: linear-gradient(135deg, #409eff, #337ecc);
  border: none;
  color: white;
}

.edit-btn:hover {
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.4);
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
  padding: 24px 0;
  display: flex;
  justify-content: center;
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
</style>