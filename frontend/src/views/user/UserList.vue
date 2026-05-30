<template>
  <div class="user-list-page">
    <div class="page-header">
      <div class="header-left">
        <h1>用户管理</h1>
        <p class="header-sub">管理系统用户账号与集群权限</p>
      </div>
      <el-button type="primary" size="large" @click="goToCreate">
        <el-icon><Plus /></el-icon>
        创建用户
      </el-button>
    </div>

    <div class="stats-row">
      <div class="stat-card">
        <div class="stat-icon blue">
          <el-icon><User /></el-icon>
        </div>
        <div class="stat-body">
          <span class="stat-value">{{ total }}</span>
          <span class="stat-label">用户总数</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon orange">
          <el-icon><UserFilled /></el-icon>
        </div>
        <div class="stat-body">
          <span class="stat-value">{{ adminCount }}</span>
          <span class="stat-label">管理员</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon green">
          <el-icon><CircleCheck /></el-icon>
        </div>
        <div class="stat-body">
          <span class="stat-value">{{ activeCount }}</span>
          <span class="stat-label">已启用</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon gray">
          <el-icon><Remove /></el-icon>
        </div>
        <div class="stat-body">
          <span class="stat-value">{{ inactiveCount }}</span>
          <span class="stat-label">已禁用</span>
        </div>
      </div>
    </div>

    <div class="content-card">
      <div class="toolbar">
        <div class="search-box">
          <el-icon class="search-icon"><Search /></el-icon>
          <input 
            v-model="keyword"
            placeholder="搜索用户名或邮箱..."
            @input="onSearch"
            class="search-input"
          />
          <el-icon 
            v-if="keyword" 
            class="clear-icon" 
            @click="clearSearch"
          >
            <Close />
          </el-icon>
        </div>
        <div class="toolbar-right">
          <span class="result-count" v-if="keyword">找到 {{ total }} 条结果</span>
        </div>
      </div>

      <el-table 
        :data="users" 
        stripe
        style="width:100%"
        :header-cell-style="{ background:'#fafbfc', color:'#475569', fontWeight:'600', fontSize:'13px' }"
        row-class-name="user-row"
      >
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column label="用户名" min-width="140">
          <template #default="scope">
            <div class="user-cell">
              <div class="user-avatar-sm" :class="{ admin: scope.row.role === 'admin' }">
                {{ (scope.row.username || '?').charAt(0).toUpperCase() }}
              </div>
              <span class="user-name">{{ scope.row.username }}</span>
              <el-tag v-if="scope.row.role === 'admin'" size="small" type="warning" effect="light">管理员</el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="email" label="邮箱" min-width="180">
          <template #default="scope">
            <span class="email-text">{{ scope.row.email || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="角色" width="100" align="center">
          <template #default="scope">
            <span class="role-badge" :class="scope.row.role">
              {{ scope.row.role === 'admin' ? '管理员' : '普通用户' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="90" align="center">
          <template #default="scope">
            <span class="status-dot" :class="{ active: scope.row.status }" />
            <span :class="scope.row.status ? 'text-success' : 'text-muted'">
              {{ scope.row.status ? '启用' : '禁用' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="170" align="center">
          <template #default="scope">
            <span class="time-text">{{ formatTime(scope.row.created_at) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right" align="center">
          <template #default="scope">
            <el-button size="small" plain @click="goToEdit(scope.row.id)">编辑</el-button>
            <el-button 
              size="small"
              plain
              :type="scope.row.status ? 'warning' : 'success'"
              :disabled="scope.row.id === 1"
              @click="toggleStatus(scope.row)"
            >
              {{ scope.row.id === 1 ? '受保护' : (scope.row.status ? '禁用' : '启用') }}
            </el-button>
            <el-button
              size="small"
              plain
              type="danger"
              :disabled="scope.row.id === 1 || scope.row.id === currentUserId"
              @click="deleteUser(scope.row.id)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrap">
        <el-pagination
          v-model:current-page="page"
          :page-size="pageSize"
          :total="total"
          @current-change="handlePageChange"
          layout="prev, pager, next, jumper, ->, total"
          background
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, User, UserFilled, CircleCheck, Remove, Search, Close } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { userAPI } from '@/utils/api'

const router = useRouter()
const authStore = useAuthStore()

const users = ref([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const keyword = ref('')
let searchTimer = null

const currentUserId = computed(() => authStore.user.value?.id)
const adminCount = computed(() => users.value.filter(u => u.role === 'admin').length)
const activeCount = computed(() => users.value.filter(u => u.status).length)
const inactiveCount = computed(() => users.value.filter(u => !u.status).length)

const formatTime = (ts) => {
  if (!ts) return '-'
  return ts.replace('T', ' ').substring(0, 19)
}

const goToCreate = () => {
  router.push('/users/create')
}

const goToEdit = (id) => {
  router.push(`/users/${id}/edit`)
}

const toggleStatus = async (user) => {
  if (user.id === 1) {
    ElMessage.warning('默认管理员账号不可禁用')
    return
  }
  try {
    await userAPI.update(user.id, { ...user, status: !user.status })
    await fetchUsers()
    ElMessage.success(`用户已${user.status ? '禁用' : '启用'}`)
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '操作失败')
  }
}

const deleteUser = async (id) => {
  try {
    await ElMessageBox.confirm('确定要删除这个用户吗？该操作不可撤销。', '删除确认', {
      confirmButtonText: '确定删除',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await userAPI.delete(id)
    await fetchUsers()
    ElMessage.success('用户已删除')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.error || '删除失败')
    }
  }
}

const onSearch = () => {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    page.value = 1
    fetchUsers()
  }, 300)
}

const clearSearch = () => {
  keyword.value = ''
  page.value = 1
  fetchUsers()
}

const handlePageChange = (newPage) => {
  page.value = newPage
  fetchUsers()
}

const fetchUsers = async () => {
  try {
    const response = await userAPI.list(page.value, pageSize.value, keyword.value)
    users.value = response.data.list || []
    total.value = response.data.total || 0
  } catch (error) {
    console.error('获取用户列表失败:', error)
  }
}

onMounted(() => {
  fetchUsers()
})
</script>

<style scoped>
.user-list-page {
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 28px;
}

.header-left h1 {
  font-size: 26px;
  font-weight: 700;
  color: #1e293b;
  margin: 0 0 4px 0;
  letter-spacing: -0.5px;
}

.header-sub {
  color: #94a3b8;
  font-size: 14px;
  margin: 0;
}

/* 统计卡片 */
.stats-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.06);
  transition: box-shadow 0.2s;
}

.stat-card:hover {
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

.stat-icon {
  width: 44px;
  height: 44px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  flex-shrink: 0;
}

.stat-icon.blue {
  background: #eff6ff;
  color: #3b82f6;
}

.stat-icon.orange {
  background: #fff7ed;
  color: #f97316;
}

.stat-icon.green {
  background: #f0fdf4;
  color: #22c55e;
}

.stat-icon.gray {
  background: #f8fafc;
  color: #94a3b8;
}

.stat-body {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #1e293b;
  line-height: 1;
}

.stat-label {
  font-size: 13px;
  color: #94a3b8;
}

/* 内容卡片 */
.content-card {
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.06);
  overflow: hidden;
}

/* 搜索栏 */
.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 24px;
  border-bottom: 1px solid #f1f5f9;
}

.search-box {
  position: relative;
  display: flex;
  align-items: center;
  width: 320px;
}

.search-icon {
  position: absolute;
  left: 12px;
  color: #94a3b8;
  font-size: 16px;
  pointer-events: none;
}

.search-input {
  width: 100%;
  height: 38px;
  padding: 0 36px 0 36px;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  font-size: 14px;
  color: #334155;
  background: #f8fafc;
  outline: none;
  transition: all 0.2s;
}

.search-input:focus {
  border-color: #3b82f6;
  background: #fff;
  box-shadow: 0 0 0 3px rgba(59,130,246,0.1);
}

.search-input::placeholder {
  color: #cbd5e1;
}

.clear-icon {
  position: absolute;
  right: 12px;
  cursor: pointer;
  color: #94a3b8;
  font-size: 14px;
  transition: color 0.15s;
}

.clear-icon:hover {
  color: #64748b;
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.result-count {
  font-size: 13px;
  color: #94a3b8;
}

/* 表格样式 */
:deep(.user-row td) {
  padding: 14px 0;
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-avatar-sm {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  background: #f1f5f9;
  color: #64748b;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 600;
  flex-shrink: 0;
}

.user-avatar-sm.admin {
  background: #fef3c7;
  color: #f59e0b;
}

.user-name {
  font-size: 14px;
  font-weight: 500;
  color: #1e293b;
}

.email-text {
  color: #64748b;
  font-size: 13px;
}

.role-badge {
  display: inline-block;
  padding: 2px 10px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
}

.role-badge.admin {
  background: #fef3c7;
  color: #92400e;
}

.role-badge.user {
  background: #f1f5f9;
  color: #475569;
}

.status-dot {
  display: inline-block;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  margin-right: 6px;
  vertical-align: middle;
  background: #cbd5e1;
}

.status-dot.active {
  background: #22c55e;
}

.text-success {
  color: #16a34a;
  font-size: 13px;
}

.text-muted {
  color: #94a3b8;
  font-size: 13px;
}

.time-text {
  color: #94a3b8;
  font-size: 13px;
}

.pagination-wrap {
  padding: 20px 24px;
  display: flex;
  justify-content: center;
}
</style>
