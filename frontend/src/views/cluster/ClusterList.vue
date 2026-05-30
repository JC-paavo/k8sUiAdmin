<template>
  <div class="cluster-list-page">
    <div class="page-header">
      <div>
        <h1>集群管理</h1>
        <p>管理您的 Kubernetes 集群</p>
      </div>
      <el-button type="primary" @click="showCreateModal = true" v-if="isAdmin">
        <el-icon name="Plus" />
        添加集群
      </el-button>
    </div>
    
    <!-- 搜索栏 -->
    <div class="search-section">
      <el-input
        v-model="searchQuery"
        placeholder="搜索集群名称..."
        class="search-input"
        clearable
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
    </div>
    
    <div class="cluster-grid">
      <div 
        v-for="cluster in paginatedClusters" 
        :key="cluster?.id ?? Math.random()"
        class="cluster-card"
        :class="cluster?.status"
      >
        <div class="card-top">
          <div class="status-badge" :class="cluster?.status">
            <span class="status-dot"></span>
            <span>{{ getStatusText(cluster?.status) }}</span>
          </div>
          <el-dropdown v-if="isAdmin" trigger="click">
            <el-icon class="more-icon"><MoreFilled /></el-icon>
            <el-dropdown-menu>
              <el-dropdown-item @click="editCluster(cluster)">编辑</el-dropdown-item>
              <el-dropdown-item @click="deleteCluster(cluster?.id)" divided>删除</el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
        </div>
        
        <div class="card-body">
          <div class="cluster-info">
            <h3 class="cluster-name">{{ cluster?.alias || cluster?.name || '未知集群' }}</h3>
            <div class="cluster-meta">
              <span class="meta-item">
                <el-icon><Server /></el-icon>
                {{ cluster?.server || '-' }}
              </span>
              <span class="meta-divider"></span>
              <span class="meta-item">
                <el-icon><Coin /></el-icon>
                Kubernetes {{ cluster?.version || '未知' }}
              </span>
              <span class="meta-divider"></span>
              <span class="meta-item">
                <el-icon><FolderOpened /></el-icon>
                {{ namespaceCounts[cluster?.id] !== undefined ? '命名空间: ' + namespaceCounts[cluster?.id] : '命名空间: ...' }}
              </span>
            </div>
            <p class="description" v-if="cluster?.description">{{ cluster.description }}</p>
          </div>

          <!-- 资源利用率 -->
          <div v-if="resourceUsage[cluster?.id]" class="usage-section">
            <div class="usage-row">
              <div class="usage-bar-group">
                <span class="usage-label">
                  <el-icon><Cpu /></el-icon> CPU
                </span>
                <div class="usage-bar-wrapper">
                  <div class="usage-bar-bg">
                    <div 
                      class="usage-bar-fill cpu" 
                      :style="{ width: Math.min(resourceUsage[cluster?.id].cpu_percent, 100) + '%' }"
                    ></div>
                  </div>
                  <span class="usage-value">{{ resourceUsage[cluster?.id].cpu_percent.toFixed(1) }}%</span>
                </div>
                <span class="usage-detail">{{ resourceUsage[cluster?.id].allocated_cpu_cores.toFixed(1) }} / {{ resourceUsage[cluster?.id].total_cpu_cores.toFixed(1) }} 核</span>
              </div>
              <div class="usage-bar-group">
                <span class="usage-label">
                  <el-icon><Memo /></el-icon> 内存
                </span>
                <div class="usage-bar-wrapper">
                  <div class="usage-bar-bg">
                    <div 
                      class="usage-bar-fill memory" 
                      :style="{ width: Math.min(resourceUsage[cluster?.id].memory_percent, 100) + '%' }"
                    ></div>
                  </div>
                  <span class="usage-value">{{ resourceUsage[cluster?.id].memory_percent.toFixed(1) }}%</span>
                </div>
                <span class="usage-detail">{{ (resourceUsage[cluster?.id].allocated_memory_mi / 1024).toFixed(1) }} / {{ (resourceUsage[cluster?.id].total_memory_mi / 1024).toFixed(1) }} Gi</span>
              </div>
            </div>
          </div>
          <div v-else-if="cluster?.status === 'connected'" class="usage-section loading-usage">
            <el-icon class="is-loading"><Loading /></el-icon> 加载资源使用率...
          </div>
        </div>
        
        <div class="card-actions">
          <el-button size="small" @click="goToDetail(cluster?.id)">查看详情</el-button>
          <el-button size="small" plain @click="refreshCluster(cluster?.id)">刷新状态</el-button>
          <el-button v-if="isAdmin" size="small" text type="danger" @click="deleteCluster(cluster?.id)">
            删除
          </el-button>
        </div>
      </div>
    </div>
    
    <!-- 分页 -->
    <div class="pagination-section" v-if="totalClusters > pageSize">
      <el-pagination
        v-model:current-page="currentPage"
        :page-size="pageSize"
        :total="totalClusters"
        layout="prev, pager, next"
        background
      />
    </div>
    
    <div v-if="clusters.length === 0" class="empty-state">
      <el-icon name="Server" class="empty-icon" />
      <p>暂无集群</p>
      <el-button type="primary" @click="showCreateModal = true" v-if="isAdmin">添加集群</el-button>
    </div>
    
    <el-dialog title="创建集群" v-model="showCreateModal" width="700px" :close-on-click-modal="false">
      <el-form ref="clusterFormRef" :model="clusterForm" :rules="clusterRules" label-width="120px">
        <el-form-item label="集群名称" prop="name">
          <el-input v-model="clusterForm.name" placeholder="请输入集群名称" />
        </el-form-item>
        <el-form-item label="集群别名" prop="alias">
          <el-input v-model="clusterForm.alias" placeholder="请输入集群别名（可选）" />
        </el-form-item>
        <template v-if="clusterForm.authType === 'token'">
          <el-form-item label="API Server" prop="server" :required="true">
            <el-input
              v-model="clusterForm.server"
              placeholder="https://your-k8s-server:6443"
            />
          </el-form-item>
        </template>
        
        <el-divider content-position="left">认证方式</el-divider>
        
        <el-form-item label="认证方式">
          <el-radio-group v-model="clusterForm.authType">
            <el-radio label="kubeconfig">Kubeconfig</el-radio>
            <el-radio label="token">Token</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <!-- Kubeconfig 模式 -->
        <template v-if="clusterForm.authType === 'kubeconfig'">
          <el-form-item label="Kubeconfig" required>
            <el-input
              v-model="clusterForm.kubeconfig"
              type="textarea"
              placeholder="粘贴 kubeconfig 文件内容"
              :rows="8"
            />
          </el-form-item>
          <el-alert
            title="提示"
            type="info"
            :closable="false"
            show-icon
            style="margin-bottom: 20px;"
          >
            Kubeconfig 将自动提取服务器地址和认证信息
          </el-alert>
        </template>
        
        <!-- Token 模式 -->
        <template v-else>
          <el-alert
            title="提示"
            type="info"
            :closable="false"
            show-icon
            style="margin-bottom: 15px;"
          >
            请手动输入 API Server 地址
          </el-alert>
          <el-form-item label="Token" required>
            <el-input
              v-model="clusterForm.token"
              type="textarea"
              placeholder="粘贴 ServiceAccount 的 Token"
              :rows="4"
            />
          </el-form-item>
          <el-form-item label="CA证书">
            <el-input
              v-model="clusterForm.caCert"
              type="textarea"
              placeholder="base64编码的CA证书（可选，不填则跳过证书验证）"
              :rows="3"
            />
          </el-form-item>
        </template>
        
        <el-form-item label="描述">
          <el-input
            v-model="clusterForm.description"
            type="textarea"
            placeholder="集群描述（可选）"
            :rows="2"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="handleCancelCreate">取消</el-button>
        <el-button type="primary" @click="handleCreate" :loading="creating">确定</el-button>
      </template>
    </el-dialog>
    
    <el-dialog title="编辑集群" v-model="showEditModal" width="700px">
      <el-form ref="editClusterFormRef" :model="editClusterForm" label-width="120px">
        <el-form-item label="集群名称">
          <el-input v-model="editClusterForm.name" placeholder="请输入集群名称" />
        </el-form-item>
        <el-form-item label="集群别名">
          <el-input v-model="editClusterForm.alias" placeholder="请输入集群别名（可选）" />
        </el-form-item>
        <el-form-item label="API Server">
          <el-input v-model="editClusterForm.server" placeholder="https://your-k8s-server:6443" />
        </el-form-item>
        
        <el-divider content-position="left">认证信息</el-divider>
        
        <el-form-item label="认证类型">
          <el-radio-group v-model="editClusterForm.authType" @change="handleAuthTypeChange">
            <el-radio label="kubeconfig">Kubeconfig</el-radio>
            <el-radio label="token">Token</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <template v-if="editClusterForm.authType === 'kubeconfig'">
          <el-form-item label="Kubeconfig">
            <el-input
              v-model="editClusterForm.kubeconfig"
              type="textarea"
              placeholder="粘贴 kubeconfig 文件内容（可选，留空则保持原配置）"
              :rows="6"
            />
          </el-form-item>
        </template>
        
        <template v-else>
          <el-form-item label="Token">
            <el-input
              v-model="editClusterForm.token"
              type="textarea"
              placeholder="粘贴 ServiceAccount 的 Token（可选，留空则保持原配置）"
              :rows="4"
            />
          </el-form-item>
          <el-form-item label="CA证书">
            <el-input
              v-model="editClusterForm.caCert"
              type="textarea"
              placeholder="base64编码的CA证书（可选）"
              :rows="3"
            />
          </el-form-item>
        </template>
        
        <el-form-item label="描述">
          <el-input
            v-model="editClusterForm.description"
            type="textarea"
            placeholder="集群描述（可选）"
            :rows="2"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="handleCancelEdit">取消</el-button>
        <el-button type="primary" @click="handleUpdate" :loading="updating">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useAuthStore } from '@/stores/auth'
import { useClusterStore } from '@/stores/cluster'
import { clusterAPI, k8sAPI } from '@/utils/api'

const router = useRouter()
const authStore = useAuthStore()
const clusterStore = useClusterStore()

const clusters = ref([])
const namespaceCounts = ref({})
const resourceUsage = ref({})
const searchQuery = ref('')
const currentPage = ref(1)
const pageSize = ref(20)
const showCreateModal = ref(false)
const showEditModal = ref(false)
const creating = ref(false)
const updating = ref(false)
const editingCluster = ref(null)

const clusterForm = reactive({
  name: '',
  alias: '',
  server: '',
  authType: 'kubeconfig',
  kubeconfig: '',
  caCert: '',
  token: '',
  description: ''
})

const editClusterForm = reactive({
  name: '',
  alias: '',
  server: '',
  authType: 'kubeconfig',
  kubeconfig: '',
  caCert: '',
  token: '',
  description: ''
})

const clusterRules = computed(() => {
  const rules = {
    name: [{ required: true, message: '请输入集群名称', trigger: 'blur' }]
  }
  
  // Token 模式下才验证 server
  if (clusterForm.authType === 'token') {
    rules.server = [
      { required: true, message: '请输入API Server地址', trigger: 'blur' },
      { pattern: /^https?:\/\/.+/, message: 'API Server地址必须是有效的URL（以http://或https://开头）', trigger: 'blur' }
    ]
  }
  
  return rules
})

const validateAuth = (rule, value, callback) => {
  if (clusterForm.authType === 'kubeconfig' && !clusterForm.kubeconfig.trim()) {
    callback(new Error('请输入Kubeconfig内容'))
  } else if (clusterForm.authType === 'token' && !clusterForm.token.trim()) {
    callback(new Error('请输入Token'))
  } else {
    callback()
  }
}

const isAdmin = computed(() => authStore.user.value?.role === 'admin')

const filteredClusters = computed(() => {
  if (!searchQuery.value) return clusters.value
  const q = searchQuery.value.toLowerCase()
  return clusters.value.filter(c =>
    (c.alias || c.name || '').toLowerCase().includes(q)
  )
})

const totalClusters = computed(() => filteredClusters.value.length)

const paginatedClusters = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredClusters.value.slice(start, end)
})

const usageColor = (percent) => {
  if (percent >= 80) return '#f56c6c'
  if (percent >= 50) return '#e6a23c'
  return '#67c23a'
}

const getStatusText = (status) => {
  const texts = {
    connected: '已连接',
    pending: '待连接',
    error: '连接失败'
  }
  return texts[status] || status
}

const goToDetail = (id) => {
  router.push(`/clusters/${id}`)
}

const refreshCluster = async (id) => {
  console.log('刷新集群 - ID:', id, '类型:', typeof id)
  if (!id || id === 'undefined' || id === '') {
    ElMessage.error('集群ID无效')
    return
  }
  try {
    await clusterAPI.refreshStatus(id)
    await fetchClusters()
    ElMessage.success('状态已刷新')
  } catch (error) {
    console.error('刷新失败:', error)
    ElMessage.error('刷新失败')
  }
}

const editCluster = (cluster) => {
  editingCluster.value = cluster
  editClusterForm.name = cluster.name
  editClusterForm.alias = cluster.alias || ''
  editClusterForm.server = cluster.server || ''
  editClusterForm.kubeconfig = cluster.kubeconfig || ''
  editClusterForm.caCert = cluster.caCert || ''
  editClusterForm.token = ''
  editClusterForm.description = cluster.description || ''
  
  // 根据是否有 kubeconfig 判断认证类型
  editClusterForm.authType = cluster.kubeconfig ? 'kubeconfig' : 'token'
  showEditModal.value = true
}

const handleAuthTypeChange = () => {
  // 切换认证类型时清空相应字段
  if (editClusterForm.authType === 'kubeconfig') {
    editClusterForm.token = ''
    editClusterForm.caCert = ''
  } else {
    editClusterForm.kubeconfig = ''
  }
}

const handleCancelEdit = () => {
  showEditModal.value = false
  editClusterForm.name = ''
  editClusterForm.alias = ''
  editClusterForm.server = ''
  editClusterForm.authType = 'kubeconfig'
  editClusterForm.kubeconfig = ''
  editClusterForm.caCert = ''
  editClusterForm.token = ''
  editClusterForm.description = ''
}

const deleteCluster = async (id) => {
  try {
    await ElMessageBox.confirm(
      '确定要删除这个集群吗？',
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await clusterAPI.delete(id)
    await fetchClusters()
    ElMessage.success('删除成功')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleCancelCreate = () => {
  showCreateModal.value = false
  clusterForm.name = ''
  clusterForm.alias = ''
  clusterForm.server = ''
  clusterForm.authType = 'kubeconfig'
  clusterForm.kubeconfig = ''
  clusterForm.caCert = ''
  clusterForm.token = ''
  clusterForm.description = ''
  // 重置表单验证
  if (clusterFormRef.value) {
    clusterFormRef.value.clearValidate()
  }
}

const handleCreate = async () => {
  // 验证必填字段
  if (!clusterForm.name.trim()) {
    ElMessage.error('请输入集群名称')
    return
  }
  
  if (!clusterForm.kubeconfig.trim() && clusterForm.authType === 'kubeconfig') {
    ElMessage.error('请输入Kubeconfig内容')
    return
  }
  
  if (clusterForm.authType === 'token') {
    if (!clusterForm.server.trim()) {
      ElMessage.error('请输入API Server地址')
      return
    }
    if (!clusterForm.token.trim()) {
      ElMessage.error('请输入Token')
      return
    }
  }
  
  creating.value = true
  try {
    const data = {
      name: clusterForm.name,
      alias: clusterForm.alias,
      description: clusterForm.description
    }
    
    if (clusterForm.authType === 'kubeconfig') {
      data.kubeconfig = clusterForm.kubeconfig
    } else {
      data.server = clusterForm.server
      data.token = clusterForm.token
      data.caCert = clusterForm.caCert
    }
    
    await clusterAPI.create(data)
    await fetchClusters()
    handleCancelCreate()
    ElMessage.success('创建成功')
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '创建失败')
  } finally {
    creating.value = false
  }
}

const handleUpdate = async () => {
  updating.value = true
  try {
    const data = {
      name: editClusterForm.name,
      alias: editClusterForm.alias,
      description: editClusterForm.description
    }
    
    if (editClusterForm.authType === 'kubeconfig') {
      if (editClusterForm.kubeconfig.trim()) {
        data.kubeconfig = editClusterForm.kubeconfig
      }
    } else {
      data.server = editClusterForm.server
      if (editClusterForm.token.trim()) {
        data.token = editClusterForm.token
      }
      if (editClusterForm.caCert.trim()) {
        data.caCert = editClusterForm.caCert
      }
    }
    
    await clusterAPI.update(editingCluster.value.id, data)
    await fetchClusters()
    handleCancelEdit()
    ElMessage.success('更新成功')
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '更新失败')
  } finally {
    updating.value = false
  }
}

const fetchClusters = async () => {
  try {
    const response = await clusterAPI.list()
    clusters.value = response.data
    clusterStore.setClusters(response.data)
    
    // 获取每个已连接集群的命名空间数量和资源使用率
    await Promise.all([
      fetchNamespaceCounts(),
      fetchResourceUsage()
    ])
  } catch (error) {
    console.error('获取集群列表失败:', error)
  }
}

const fetchNamespaceCounts = async () => {
  const counts = {}
  const promises = clusters.value.map(async (cluster) => {
    if (cluster.status !== 'connected') {
      counts[cluster.id] = 0
      return
    }
    try {
      const nsResponse = await k8sAPI.listNamespaces(cluster.id)
      counts[cluster.id] = nsResponse.data?.length || 0
    } catch {
      counts[cluster.id] = 0
    }
  })
  await Promise.all(promises)
  namespaceCounts.value = counts
}

const fetchResourceUsage = async () => {
  const usage = {}
  const promises = clusters.value.map(async (cluster) => {
    if (cluster.status !== 'connected') return
    try {
      const resp = await k8sAPI.getClusterResourceUsage(cluster.id)
      usage[cluster.id] = resp.data
    } catch {
    }
  })
  await Promise.all(promises)
  resourceUsage.value = usage
}

watch(searchQuery, () => {
  currentPage.value = 1
})

// 监听认证类型变化，自动清空不必要的字段
watch(() => clusterForm.authType, (newType) => {
  if (newType === 'kubeconfig') {
    clusterForm.server = ''
    clusterForm.token = ''
    clusterForm.caCert = ''
  } else {
    clusterForm.kubeconfig = ''
  }
})

onMounted(() => {
  fetchClusters()
})
</script>

<style scoped>
.cluster-list-page {
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 28px;
}

.page-header h1 {
  font-size: 22px;
  font-weight: 600;
  color: #1a1a2e;
  margin-bottom: 4px;
}

.page-header p {
  color: #8b8fa3;
  font-size: 14px;
}

.search-section {
  margin-bottom: 24px;
}

.search-section .search-input {
  max-width: 380px;
}

.cluster-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(380px, 1fr));
  gap: 20px;
}

.cluster-card {
  background: #ffffff;
  border-radius: 14px;
  padding: 0;
  border: 1px solid #eaecf0;
  transition: all 0.25s ease;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.cluster-card:hover {
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.06);
  border-color: #d0d5dd;
  transform: translateY(-2px);
}

.cluster-card.connected {
  border-left: 4px solid #22c55e;
}

.cluster-card.error {
  border-left: 4px solid #ef4444;
}

.cluster-card.disconnected {
  border-left: 4px solid #94a3b8;
}

.card-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 18px 22px 0;
}

.status-badge {
  display: flex;
  align-items: center;
  gap: 7px;
  font-size: 12px;
  font-weight: 500;
  padding: 4px 12px;
  border-radius: 20px;
  background: #f8f9fc;
  color: #64748b;
}

.status-badge.connected {
  background: #f0fdf4;
  color: #16a34a;
}

.status-badge.error {
  background: #fef2f2;
  color: #dc2626;
}

.status-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #94a3b8;
}

.cluster-card.connected .status-dot {
  background: #22c55e;
}

.cluster-card.error .status-dot {
  background: #ef4444;
}

.more-icon {
  cursor: pointer;
  color: #98a2b3;
  font-size: 18px;
  padding: 4px;
  border-radius: 6px;
  transition: all 0.15s;
}

.more-icon:hover {
  background: #f2f4f7;
  color: #475467;
}

.card-body {
  padding: 16px 22px;
  flex: 1;
}

.cluster-info {
  margin-bottom: 4px;
}

.cluster-name {
  font-size: 17px;
  font-weight: 600;
  color: #1a1a2e;
  margin-bottom: 10px;
}

.cluster-meta {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 6px;
}

.meta-item {
  font-size: 12px;
  color: #667085;
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.meta-item .el-icon {
  font-size: 13px;
}

.meta-divider {
  width: 3px;
  height: 3px;
  border-radius: 50%;
  background: #d0d5dd;
  display: inline-block;
}

.description {
  font-size: 13px;
  color: #667085;
  margin-top: 10px;
  line-height: 1.5;
}

.usage-section {
  margin-top: 16px;
  padding: 14px 16px;
  background: #f8fafc;
  border-radius: 10px;
  border: 1px solid #f1f5f9;
}

.usage-row {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.usage-bar-group {
  display: flex;
  align-items: center;
  gap: 10px;
}

.usage-label {
  font-size: 12px;
  font-weight: 500;
  color: #475467;
  min-width: 50px;
  display: flex;
  align-items: center;
  gap: 4px;
}

.usage-bar-wrapper {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 8px;
}

.usage-bar-bg {
  flex: 1;
  height: 8px;
  background: #e4e7ec;
  border-radius: 100px;
  overflow: hidden;
}

.usage-bar-fill {
  height: 100%;
  border-radius: 100px;
  transition: width 0.6s ease;
}

.usage-bar-fill.cpu {
  background: linear-gradient(90deg, #3b82f6, #2563eb);
}

.usage-bar-fill.memory {
  background: linear-gradient(90deg, #8b5cf6, #7c3aed);
}

.usage-value {
  font-size: 12px;
  font-weight: 600;
  color: #1a1a2e;
  min-width: 44px;
  text-align: right;
}

.usage-detail {
  font-size: 11px;
  color: #98a2b3;
  min-width: 85px;
  text-align: right;
}

.loading-usage {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #98a2b3;
  justify-content: center;
  padding: 18px;
}

.card-actions {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 14px 22px;
  border-top: 1px solid #f2f4f7;
  background: #fafbfc;
}

.empty-state {
  text-align: center;
  padding: 80px 20px;
  background: #ffffff;
  border-radius: 14px;
  border: 1px solid #eaecf0;
}

.empty-state p {
  color: #98a2b3;
  margin-bottom: 16px;
}

.pagination-section {
  display: flex;
  justify-content: center;
  margin-top: 28px;
  padding-top: 24px;
  border-top: 1px solid #eaecf0;
}
</style>