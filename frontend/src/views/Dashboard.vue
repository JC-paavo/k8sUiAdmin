<template>
  <div class="dashboard">
    <div class="dashboard-header">
      <div class="welcome-section">
        <h1>欢迎回来，{{ user?.username }}</h1>
        <p>Kubernetes 集群管理控制台</p>
      </div>
      <div class="header-actions">
        <el-button v-if="isAdmin" type="primary" @click="goToClusterList">
          <el-icon><Plus /></el-icon>
          添加集群
        </el-button>
      </div>
    </div>
    
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon blue">
          <el-icon><Link /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ clusters.length }}</div>
          <div class="stat-label">集群总数</div>
        </div>
      </div>
      
      <div class="stat-card success">
        <div class="stat-icon green">
          <el-icon><CircleCheck /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ connectedClusters }}</div>
          <div class="stat-label">Connected</div>
        </div>
        <div class="stat-indicator success"></div>
      </div>
      
      <div class="stat-card warning">
        <div class="stat-icon orange">
          <el-icon><Clock /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ pendingClusters }}</div>
          <div class="stat-label">Pending</div>
        </div>
        <div class="stat-indicator warning"></div>
      </div>
      
      <div class="stat-card danger">
        <div class="stat-icon red">
          <el-icon><Warning /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ errorClusters }}</div>
          <div class="stat-label">失败</div>
        </div>
        <div class="stat-indicator error"></div>
      </div>
    </div>
    
    <div class="section">
      <div class="section-header">
        <h2>集群状态</h2>
        <div class="section-actions">
          <el-button size="small" @click="fetchClusters" :loading="loading">
            <el-icon><Refresh /></el-icon>
            刷新全部
          </el-button>
        </div>
      </div>
      
      <div v-if="loading" class="loading-state">
        <el-icon class="is-loading" :size="32"><Loading /></el-icon>
        <p>加载集群中...</p>
      </div>
      
      <div v-else-if="clusters.length > 0" class="cluster-list">
        <div 
          v-for="cluster in clusters" 
          :key="cluster?.id ?? Math.random()"
          class="cluster-item"
          :class="cluster?.status"
          @click="goToCluster(cluster?.id)"
        >
          <div class="cluster-status-indicator" :class="cluster?.status"></div>
          
          <div class="cluster-main">
            <div class="cluster-header">
              <h3>{{ cluster?.alias || cluster?.name || '未知集群' }}</h3>
              <el-tag size="small" :type="getStatusTagType(cluster?.status)">
                {{ getStatusText(cluster?.status) }}
              </el-tag>
            </div>
            
            <div class="cluster-details">
              <div class="detail-item">
                <el-icon><Link /></el-icon>
                <span>{{ cluster?.server || '-' }}</span>
              </div>
              <div class="detail-item">
                <el-icon><Setting /></el-icon>
                <span>K8s {{ cluster?.version || 'Unknown' }}</span>
              </div>
            </div>
          </div>
          
          <div class="cluster-actions">
            <el-button size="small" type="primary" plain @click.stop="goToCluster(cluster?.id)">
              <el-icon><View /></el-icon>
              管理
            </el-button>
            <el-button size="small" @click.stop="refreshCluster(cluster?.id)">
              <el-icon><Refresh /></el-icon>
            </el-button>
          </div>
        </div>
      </div>
      
      <div v-else class="empty-state">
        <div class="empty-icon">
          <el-icon :size="64"><Link /></el-icon>
        </div>
        <h3>{{ isAdmin ? '暂无集群' : '暂无授权集群' }}</h3>
        <p>{{ isAdmin ? '添加您的第一个 Kubernetes 集群以开始使用' : '请联系管理员为您分配集群访问权限' }}</p>
        <el-button v-if="isAdmin" type="primary" @click="goToClusterList">
          <el-icon><Plus /></el-icon>
          添加第一个集群
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { 
  Plus, Link, CircleCheck, Clock, Warning, Refresh, 
  Loading, Setting, View 
} from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { useClusterStore } from '@/stores/cluster'
import { clusterAPI } from '@/utils/api'

const router = useRouter()
const authStore = useAuthStore()
const clusterStore = useClusterStore()

const clusters = ref([])
const loading = ref(false)

const user = computed(() => authStore.user.value)
const isAdmin = computed(() => authStore.user.value?.role === 'admin')
const connectedClusters = computed(() => clusters.value.filter(c => c.status === 'connected').length)
const pendingClusters = computed(() => clusters.value.filter(c => c.status === 'pending').length)
const errorClusters = computed(() => clusters.value.filter(c => c.status === 'error').length)

const getStatusText = (status) => {
  const texts = {
    connected: 'Connected',
    pending: 'Pending',
    error: 'Failed'
  }
  return texts[status] || status
}

const getStatusTagType = (status) => {
  const types = {
    connected: 'success',
    pending: 'warning',
    error: 'danger'
  }
  return types[status] || 'info'
}

const goToCluster = (id) => {
  router.push(`/clusters/${id}`)
}

const goToClusterList = () => {
  router.push('/clusters')
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

const fetchClusters = async () => {
  if (!authStore.token.value) {
    router.push('/login')
    return
  }
  
  loading.value = true
  try {
    const response = await clusterAPI.list()
    console.log('获取到的集群数据:', response.data)
    console.log('集群数量:', response.data?.length)
    if (response.data?.length > 0) {
      console.log('第一个集群:', response.data[0])
      console.log('第一个集群的ID:', response.data[0].id)
    }
    clusters.value = response.data
    clusterStore.setClusters(response.data)
  } catch (error) {
    console.error('获取集群列表失败:', error)
    if (error.response?.status === 401) {
      ElMessage.error('登录已过期，请重新登录')
      authStore.logout()
      router.push('/login')
    }
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchClusters()
})
</script>

<style scoped>
.dashboard {
  max-width: 1200px;
  margin: 0 auto;
  padding: 32px;
  animation: fadeIn 0.4s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
  padding: 24px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
}

.welcome-section h1 {
  font-size: 28px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 8px;
  background: linear-gradient(135deg, var(--primary-light), var(--success-light));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.welcome-section p {
  color: var(--text-tertiary);
  font-size: 14px;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  margin-bottom: 32px;
}

.stat-card {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 16px;
  position: relative;
  overflow: hidden;
  transition: all var(--transition-base);
  animation: slideUp 0.5s ease backwards;
}

.stat-card:nth-child(1) { animation-delay: 0.1s; }
.stat-card:nth-child(2) { animation-delay: 0.2s; }
.stat-card:nth-child(3) { animation-delay: 0.3s; }
.stat-card:nth-child(4) { animation-delay: 0.4s; }

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

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-lg);
}

.stat-card.success {
  border-color: var(--success-color);
  background: linear-gradient(135deg, rgba(16, 185, 129, 0.05), transparent);
}

.stat-card.warning {
  border-color: var(--warning-color);
  background: linear-gradient(135deg, rgba(245, 158, 11, 0.05), transparent);
}

.stat-card.danger {
  border-color: var(--error-color);
  background: linear-gradient(135deg, rgba(239, 68, 68, 0.05), transparent);
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: var(--radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  flex-shrink: 0;
}

.stat-icon.blue {
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.2), rgba(59, 130, 246, 0.1));
  color: var(--primary-color);
}

.stat-icon.green {
  background: linear-gradient(135deg, rgba(16, 185, 129, 0.2), rgba(16, 185, 129, 0.1));
  color: var(--success-color);
}

.stat-icon.orange {
  background: linear-gradient(135deg, rgba(245, 158, 11, 0.2), rgba(245, 158, 11, 0.1));
  color: var(--warning-color);
}

.stat-icon.red {
  background: linear-gradient(135deg, rgba(239, 68, 68, 0.2), rgba(239, 68, 68, 0.1));
  color: var(--error-color);
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 32px;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1;
  margin-bottom: 8px;
}

.stat-label {
  color: var(--text-tertiary);
  font-size: 13px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.stat-indicator {
  position: absolute;
  top: 0;
  right: 0;
  width: 4px;
  height: 100%;
}

.stat-indicator.success {
  background: var(--success-color);
}

.stat-indicator.warning {
  background: var(--warning-color);
}

.stat-indicator.error {
  background: var(--error-color);
}

.section {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  padding: 24px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.section-header h2 {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
}

.loading-state {
  text-align: center;
  padding: 60px 20px;
  color: var(--text-tertiary);
}

.loading-state .el-icon {
  margin-bottom: 16px;
  color: var(--primary-color);
}

.cluster-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.cluster-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all var(--transition-fast);
  position: relative;
  overflow: hidden;
}

.cluster-item:hover {
  background: var(--bg-hover);
  border-color: var(--border-light);
  transform: translateX(4px);
}

.cluster-item.connected {
  border-left: 3px solid var(--success-color);
}

.cluster-item.pending {
  border-left: 3px solid var(--warning-color);
}

.cluster-item.error {
  border-left: 3px solid var(--error-color);
}

.cluster-status-indicator {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  flex-shrink: 0;
  animation: pulse 2s ease-in-out infinite;
}

.cluster-status-indicator.connected {
  background: var(--success-color);
  box-shadow: 0 0 8px var(--success-glow);
}

.cluster-status-indicator.pending {
  background: var(--warning-color);
  box-shadow: 0 0 8px var(--warning-glow);
}

.cluster-status-indicator.error {
  background: var(--error-color);
  box-shadow: 0 0 8px var(--error-glow);
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.cluster-main {
  flex: 1;
  min-width: 0;
}

.cluster-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.cluster-header h3 {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.cluster-details {
  display: flex;
  gap: 24px;
}

.detail-item {
  display: flex;
  align-items: center;
  gap: 6px;
  color: var(--text-tertiary);
  font-size: 13px;
}

.detail-item .el-icon {
  font-size: 14px;
}

.cluster-actions {
  display: flex;
  gap: 8px;
  opacity: 0;
  transition: opacity var(--transition-fast);
}

.cluster-item:hover .cluster-actions {
  opacity: 1;
}

.empty-state {
  text-align: center;
  padding: 80px 20px;
}

.empty-icon {
  width: 120px;
  height: 120px;
  margin: 0 auto 24px;
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.1), rgba(16, 185, 129, 0.1));
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-tertiary);
  opacity: 0.5;
}

.empty-state h3 {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.empty-state p {
  color: var(--text-tertiary);
  margin-bottom: 24px;
}

/* Responsive */
@media (max-width: 1024px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
  
  .dashboard-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }
  
  .cluster-details {
    flex-direction: column;
    gap: 8px;
  }
  
  .cluster-actions {
    opacity: 1;
  }
}
</style>
