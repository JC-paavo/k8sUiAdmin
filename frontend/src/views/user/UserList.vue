<template>
  <div class="user-list-page">
    <div class="page-header">
      <div>
        <h1>用户管理</h1>
        <p>管理系统用户和权限</p>
      </div>
      <el-button type="primary" @click="goToCreate">
        <el-icon name="Plus" />
        创建用户
      </el-button>
    </div>
    
    <div class="user-table">
      <el-table :data="users" border>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" />
        <el-table-column prop="email" label="邮箱" />
        <el-table-column label="角色">
          <template #default="scope">
            <span :class="scope.row.role === 'admin' ? 'role-admin' : 'role-user'">
              {{ scope.row.role === 'admin' ? '管理员' : '普通用户' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="状态">
          <template #default="scope">
            <span :class="scope.row.status ? 'status-active' : 'status-inactive'">
              {{ scope.row.status ? '启用' : '禁用' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" />
        <el-table-column label="操作">
          <template #default="scope">
            <el-button size="small" @click="goToEdit(scope.row.id)">编辑</el-button>
            <el-button 
              size="small" 
              :type="scope.row.status ? 'warning' : 'success'"
              @click="toggleStatus(scope.row)"
            >
              {{ scope.row.status ? '禁用' : '启用' }}
            </el-button>
            <el-button
              size="small"
              danger
              @click="deleteUser(scope.row.id)"
              :disabled="scope.row.id === 1 || scope.row.id === currentUserId"
            >
              {{ scope.row.id === 1 ? '不能删除' : '删除' }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <div class="pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :total="total"
          @current-change="handlePageChange"
          layout="prev, pager, next, jumper, ->, total"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useAuthStore } from '@/stores/auth'
import { userAPI } from '@/utils/api'

const router = useRouter()
const authStore = useAuthStore()

const users = ref([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const currentUserId = computed(() => authStore.user.value?.id)

const goToCreate = () => {
  router.push('/users/create')
}

const goToEdit = (id) => {
  router.push(`/users/${id}/edit`)
}

const toggleStatus = async (user) => {
  try {
    await userAPI.update(user.id, { ...user, status: !user.status })
    await fetchUsers()
    ElMessage.success(`用户已${user.status ? '禁用' : '启用'}`)
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const deleteUser = async (id) => {
  try {
    await ElMessageBox.confirm(
      '确定要删除这个用户吗？',
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await userAPI.delete(id)
    await fetchUsers()
    ElMessage.success('删除成功')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handlePageChange = (newPage) => {
  page.value = newPage
  fetchUsers()
}

const fetchUsers = async () => {
  try {
    const response = await userAPI.list(page.value, pageSize.value)
    users.value = response.data.list
    total.value = response.data.total
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
  align-items: center;
  margin-bottom: 24px;
}

.page-header h1 {
  font-size: 24px;
  color: #303133;
  margin-bottom: 4px;
}

.page-header p {
  color: #909399;
}

.user-table {
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.role-admin {
  color: #f97316;
}

.role-user {
  color: #64748b;
}

.status-active {
  color: #22c55e;
}

.status-inactive {
  color: #909399;
}

.pagination {
  padding: 20px;
  display: flex;
  justify-content: center;
}
</style>