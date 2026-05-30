<template>
  <div class="user-form-page">
    <div class="page-header">
      <h1>编辑用户</h1>
    </div>
    
    <div v-if="user" class="form-container">
      <el-form ref="userFormRef" :model="userForm" :rules="userRules" label-width="100px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="userForm.username" placeholder="请输入用户名" :disabled="isAdmin" />
        </el-form-item>
        
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="userForm.email" placeholder="请输入邮箱（可选）" />
        </el-form-item>
        
        <el-form-item label="角色" prop="role">
          <el-select v-model="userForm.role" placeholder="请选择角色" :disabled="isAdmin">
            <el-option label="普通用户" value="user" />
            <el-option label="管理员" value="admin" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="状态">
          <el-switch v-model="userForm.status" :disabled="isAdmin" />
        </el-form-item>

        <el-form-item label="新密码">
          <el-input 
            v-model="userForm.password" 
            placeholder="留空则不修改密码" 
            type="password" 
            show-password 
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSubmit">保存</el-button>
          <el-button @click="goBack">取消</el-button>
        </el-form-item>
      </el-form>
    </div>
    
    <div v-else class="loading-state">
      <el-loading text="加载中..." />
    </div>
  </div>
</template>

<script setup>
import { reactive, ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { userAPI } from '@/utils/api'

const route = useRoute()
const router = useRouter()

const user = ref(null)
const userFormRef = ref(null)
const isAdmin = computed(() => route.params.id === '1')

const userForm = reactive({
  username: '',
  email: '',
  role: 'user',
  status: true,
  password: ''
})

const userRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在3-20个字符', trigger: 'blur' }
  ],
  email: [
    { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' }
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'blur' }
  ]
}

const handleSubmit = async () => {
  const valid = await userFormRef.value.validate().catch(() => false)
  if (!valid) return

  try {
    const data = {
      username: userForm.username,
      email: userForm.email,
      role: userForm.role,
      status: userForm.status
    }
    if (userForm.password) {
      data.password = userForm.password
    }
    
    await userAPI.update(route.params.id, data)
    ElMessage.success('用户更新成功')
    router.push('/users')
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '更新失败')
  }
}

const goBack = () => {
  router.push('/users')
}

const fetchUser = async () => {
  try {
    const response = await userAPI.get(route.params.id)
    user.value = response.data
    userForm.username = response.data.username
    userForm.email = response.data.email || ''
    userForm.role = response.data.role
    userForm.status = response.data.status
  } catch (error) {
    console.error('获取用户信息失败:', error)
  }
}

onMounted(() => {
  fetchUser()
})
</script>

<style scoped>
.user-form-page {
  max-width: 500px;
  margin: 0 auto;
}

.page-header h1 {
  font-size: 24px;
  color: #303133;
  margin-bottom: 24px;
}

.form-container {
  background: white;
  padding: 32px;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 300px;
}
</style>