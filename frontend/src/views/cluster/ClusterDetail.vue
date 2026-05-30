<template>
  <div class="cluster-detail-page">
    <div v-if="cluster" class="page-content">
      <div class="page-header">
        <div>
          <h1>{{ cluster.alias || cluster.name }}</h1>
          <p>{{ cluster.server }}</p>
        </div>
        <div class="header-actions">
          <span class="status-badge" :class="cluster.status">
            {{ getStatusText(cluster.status) }}
          </span>
          <el-button size="small" @click="refreshStatus">刷新状态</el-button>
        </div>
      </div>
      
      <div class="cluster-info">
        <div class="info-section">
          <h3>基本信息</h3>
          <div class="info-grid">
            <div class="info-item">
              <label>集群名称</label>
              <span>{{ cluster.name }}</span>
            </div>
            <div class="info-item">
              <label>集群别名</label>
              <span>{{ cluster.alias || '-' }}</span>
            </div>
            <div class="info-item">
              <label>API Server</label>
              <span>{{ cluster.server }}</span>
            </div>
            <div class="info-item">
              <label>Kubernetes 版本</label>
              <span>{{ cluster.version || '未知' }}</span>
            </div>
            <div class="info-item full">
              <label>描述</label>
              <span>{{ cluster.description || '-' }}</span>
            </div>
          </div>
        </div>
        
        <div class="info-section">
          <h3>资源管理</h3>
          <div class="resource-grid">
            <div 
              v-for="resource in resourceTypes" 
              :key="resource.type"
              class="resource-card"
              @click="goToResource(resource.type)"
            >
              <div class="resource-icon">
                <component :is="resource.icon" :size="22" />
              </div>
              <div class="resource-info">
                <h4>{{ resource.label }}</h4>
                <p>{{ resource.description }}</p>
              </div>
              <el-icon name="ChevronRight" class="arrow-icon" />
            </div>
          </div>
        </div>
        
        <div class="info-section" v-if="isAdmin">
          <h3>用户权限</h3>
          <div class="permissions-header">
            <span>已授权用户</span>
            <el-button size="small" type="primary" @click="showPermissionModal = true">添加权限</el-button>
          </div>
          <div class="permissions-list">
            <div 
              v-for="perm in permissions" 
              :key="perm.user_id"
              class="permission-item"
            >
              <div class="perm-info">
                <span class="username">{{ perm.user?.username }}</span>
                <span class="perm-badge" :class="perm.permission">
                  {{ perm.permission === 'write' ? '完全权限' : '只读权限' }}
                </span>
              </div>
              <el-button size="small" @click="removePermission(perm.user_id)" danger>移除</el-button>
            </div>
          </div>
          <div v-if="permissions.length === 0" class="empty-permissions">
            <p>暂无授权用户</p>
          </div>
        </div>
      </div>
    </div>
    
    <div v-else class="loading-state">
      <el-loading text="加载中..." />
    </div>
    
    <el-dialog title="添加用户权限" v-model="showPermissionModal" width="400px" @open="onPermDialogOpen">
      <el-form ref="permFormRef" :model="permissionForm" :rules="permissionRules" label-width="80px">
        <el-form-item label="用户" prop="userId">
          <el-select v-model="permissionForm.userId" placeholder="请选择用户" filterable>
            <el-option 
              v-for="user in availableUsers" 
              :key="user.id" 
              :label="user.username" 
              :value="user.id" 
            />
          </el-select>
        </el-form-item>
        <el-form-item label="权限" prop="permission">
          <el-select v-model="permissionForm.permission" placeholder="请选择权限">
            <el-option label="只读" value="read" />
            <el-option label="完全权限" value="write" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showPermissionModal = false">取消</el-button>
        <el-button type="primary" :loading="permSubmitting" @click="handleAddPermission">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useAuthStore } from '@/stores/auth'
import { clusterAPI, userAPI } from '@/utils/api'
import { k8sIcons } from '@/assets/icons/index.js'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const cluster = ref(null)
const permissions = ref([])
const availableUsers = ref([])
const showPermissionModal = ref(false)
const loading = ref(false)
const permSubmitting = ref(false)
const permFormRef = ref(null)

const permissionForm = reactive({
  userId: '',
  permission: 'read'
})

const permissionRules = {
  userId: [{ required: true, message: '请选择用户', trigger: 'blur' }],
  permission: [{ required: true, message: '请选择权限', trigger: 'blur' }]
}

const isAdmin = computed(() => authStore.user.value?.role === 'admin')

const resourceTypes = [
  { type: 'deployments', label: 'Deployment', description: '管理无状态应用', icon: k8sIcons.deployments },
  { type: 'statefulsets', label: 'StatefulSet', description: '管理有状态应用', icon: k8sIcons.statefulsets },
  { type: 'daemonsets', label: 'DaemonSet', description: '管理节点守护进程', icon: k8sIcons.daemonsets },
  { type: 'services', label: 'Service', description: '管理网络服务', icon: k8sIcons.services },
  { type: 'ingresses', label: 'Ingress', description: '管理入口流量', icon: k8sIcons.ingresses },
  { type: 'configmaps', label: 'ConfigMap', description: '管理配置数据', icon: k8sIcons.configmaps },
  { type: 'secrets', label: 'Secret', description: '管理密钥凭证', icon: k8sIcons.secrets },
  { type: 'pods', label: 'Pod', description: '管理容器实例', icon: k8sIcons.pods }
]

const getStatusText = (status) => {
  const texts = {
    connected: '已连接',
    pending: '待连接',
    error: '连接失败'
  }
  return texts[status] || status
}

const goToResource = (type) => {
  console.log('goToResource 被调用, type:', type, '当前路由参数:', route.params)
  
  let clusterId = route.params.id
  
  if (!clusterId) {
    ElMessage.error('集群ID无效')
    return
  }
  
  clusterId = Number(clusterId)
  if (isNaN(clusterId) || clusterId <= 0) {
    ElMessage.error('集群ID无效')
    return
  }
  
  console.log('准备跳转到资源页面, clusterId:', clusterId, 'type:', type)
  
  // 跳转到独立的资源列表页面
  const routes = {
    deployments: `/clusters/${clusterId}/deployments`,
    statefulsets: `/clusters/${clusterId}/statefulsets`,
    daemonsets: `/clusters/${clusterId}/daemonsets`,
    services: `/clusters/${clusterId}/services`,
    ingresses: `/clusters/${clusterId}/ingresses`,
    configmaps: `/clusters/${clusterId}/configmaps`,
    secrets: `/clusters/${clusterId}/secrets`,
    pods: `/clusters/${clusterId}/pods`
  }
  
  const targetRoute = routes[type]
  if (!targetRoute) {
    ElMessage.error('未知的资源类型')
    return
  }
  
  try {
    router.push(targetRoute)
  } catch (error) {
    console.error('跳转失败:', error)
    ElMessage.error('页面跳转失败')
  }
}

const fetchCluster = async () => {
  let clusterId = route.params.id
  
  // 强制类型转换，确保 ID 是有效的数字
  if (clusterId) {
    clusterId = Number(clusterId)
    if (isNaN(clusterId) || clusterId <= 0) {
      console.error('集群ID无效:', route.params.id, '转换后:', clusterId)
      ElMessage.error('集群ID无效')
      return
    }
  }
  
  console.log('fetchCluster - clusterId:', clusterId, '原始值:', route.params.id, '类型:', typeof route.params.id)
  
  if (!clusterId) {
    console.error('集群ID为空')
    ElMessage.error('集群ID无效')
    return
  }
  
  loading.value = true
  try {
    const response = await clusterAPI.get(clusterId)
    console.log('获取到的集群详情:', response.data)
    cluster.value = response.data
  } catch (error) {
    console.error('获取集群信息失败:', error)
    ElMessage.error('获取集群信息失败')
  } finally {
    loading.value = false
  }
}

const fetchPermissions = async () => {
  if (!isAdmin.value) return
  
  let clusterId = route.params.id
  if (!clusterId) return
  
  clusterId = Number(clusterId)
  if (isNaN(clusterId) || clusterId <= 0) return
  
  try {
    const response = await clusterAPI.getPermissions(clusterId)
    permissions.value = response.data
  } catch (error) {
    console.error('获取权限列表失败:', error)
  }
}

const fetchUsers = async () => {
  if (!isAdmin.value) return
  
  let clusterId = route.params.id
  if (!clusterId) return
  
  clusterId = Number(clusterId)
  if (isNaN(clusterId) || clusterId <= 0) return
  
  try {
    const response = await userAPI.list(1, 100)
    const adminId = authStore.user.value?.id
    availableUsers.value = response.data.list.filter(u => u.id !== adminId)
  } catch (error) {
    console.error('获取用户列表失败:', error)
  }
}

const handleAddPermission = async () => {
  const valid = await permFormRef.value?.validate().catch(() => false)
  if (!valid) return

  let clusterId = route.params.id
  if (!clusterId) {
    ElMessage.error('集群ID无效')
    return
  }
  
  clusterId = Number(clusterId)
  if (isNaN(clusterId) || clusterId <= 0) {
    ElMessage.error('集群ID无效')
    return
  }
  
  if (!permissionForm.userId) {
    ElMessage.error('请选择用户')
    return
  }
  
  try {
    permSubmitting.value = true
    await clusterAPI.addPermission(clusterId, permissionForm.userId, permissionForm.permission)
    ElMessage.success('权限添加成功')
    showPermissionModal.value = false
    permissionForm.userId = ''
    permissionForm.permission = 'read'
    await fetchPermissions()
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '添加失败')
  } finally {
    permSubmitting.value = false
  }
}

const onPermDialogOpen = () => {
  fetchUsers()
}

const removePermission = async (userId) => {
  let clusterId = route.params.id
  if (!clusterId) {
    ElMessage.error('集群ID无效')
    return
  }
  
  clusterId = Number(clusterId)
  if (isNaN(clusterId) || clusterId <= 0) {
    ElMessage.error('集群ID无效')
    return
  }
  
  try {
    await ElMessageBox.confirm(
      '确定要移除该用户的权限吗？',
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await clusterAPI.removePermission(clusterId, userId)
    ElMessage.success('权限已移除')
    await fetchPermissions()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('移除失败')
    }
  }
}

const refreshStatus = async () => {
  let clusterId = route.params.id
  
  if (clusterId) {
    clusterId = Number(clusterId)
    if (isNaN(clusterId) || clusterId <= 0) {
      ElMessage.error('集群ID无效')
      return
    }
  }
  
  if (!clusterId) {
    ElMessage.error('集群ID无效')
    return
  }
  
  try {
    await clusterAPI.refreshStatus(clusterId)
    await fetchCluster()
    ElMessage.success('状态已刷新')
  } catch (error) {
    ElMessage.error('刷新失败')
  }
}

// 监听路由参数变化
watch(
  () => route.params.id,
  (newId) => {
    if (newId && newId !== 'undefined' && newId !== '' && !Number.isNaN(Number(newId))) {
      fetchCluster()
      fetchPermissions()
      fetchUsers()
    }
  },
  { immediate: true }
)

 onMounted(() => {
  setTimeout(() => {
    const clusterId = route.params.id
    if (clusterId && clusterId !== 'undefined' && clusterId !== '' && !Number.isNaN(Number(clusterId))) {
      fetchCluster()
      fetchPermissions()
      fetchUsers()
    }
  }, 200)
})
</script>

<style scoped>
.cluster-detail-page {
  max-width: 1000px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.page-header h1 {
  font-size: 24px;
  color: #303133;
  margin-bottom: 4px;
}

.page-header p {
  color: #909399;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.status-badge {
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 13px;
  background: #f3f4f6;
  color: #64748b;
}

.status-badge.connected {
  background: #dcfce7;
  color: #16a34a;
}

.status-badge.error {
  background: #fee2e2;
  color: #dc2626;
}

.cluster-info {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.info-section {
  background: white;
  padding: 24px;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.info-section h3 {
  font-size: 16px;
  color: #303133;
  margin-bottom: 16px;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-item.full {
  grid-column: span 2;
}

.info-item label {
  font-size: 13px;
  color: #909399;
}

.info-item span {
  font-size: 14px;
  color: #303133;
}

.resource-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.resource-card {
  display: flex;
  align-items: center;
  padding: 16px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.resource-card:hover {
  border-color: #3b82f6;
  background: #f8fafc;
}

.resource-icon {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  background: #e8f4fd;
  color: #3b82f6;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  margin-right: 12px;
}

.resource-info h4 {
  font-size: 14px;
  color: #303133;
  margin-bottom: 2px;
}

.resource-info p {
  font-size: 12px;
  color: #909399;
}

.arrow-icon {
  margin-left: auto;
  color: #909399;
}

.permissions-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.permissions-header span {
  font-size: 14px;
  color: #64748b;
}

.permissions-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.permission-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #f9fafb;
  border-radius: 8px;
}

.perm-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.username {
  font-size: 14px;
  color: #303133;
}

.perm-badge {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  background: #e0e7ff;
  color: #6366f1;
}

.perm-badge.write {
  background: #dcfce7;
  color: #16a34a;
}

.empty-permissions {
  padding: 20px;
  text-align: center;
  color: #909399;
}

.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 400px;
}
</style>