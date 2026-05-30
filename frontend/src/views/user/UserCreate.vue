<template>
  <div class="user-form-page">
    <div class="page-header">
      <h1>创建用户</h1>
    </div>
    
    <div class="form-container">
      <el-form ref="userFormRef" :model="userForm" :rules="userRules" label-width="100px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="userForm.username" placeholder="请输入用户名" />
        </el-form-item>
        
        <el-form-item label="密码" prop="password">
          <el-input v-model="userForm.password" type="password" placeholder="请输入密码" />
        </el-form-item>
        
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input v-model="userForm.confirmPassword" type="password" placeholder="请确认密码" />
        </el-form-item>
        
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="userForm.email" placeholder="请输入邮箱（可选）" />
        </el-form-item>
        
        <el-form-item label="角色" prop="role">
          <el-select v-model="userForm.role" placeholder="请选择角色">
            <el-option label="普通用户" value="user" />
            <el-option label="管理员" value="admin" />
          </el-select>
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="handleSubmit">创建</el-button>
          <el-button @click="goBack">取消</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { userAPI } from '@/utils/api'

const router = useRouter()

const userFormRef = ref(null)

const userForm = reactive({
  username: '',
  password: '',
  confirmPassword: '',
  email: '',
  role: 'user'
})

const userRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在3-20个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少6个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    { 
      validator: (rule, value, callback) => {
        if (value !== userForm.password) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  email: [],
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
      password: userForm.password,
      email: userForm.email,
      role: userForm.role,
      status: true
    }
    
    await userAPI.create(data)
    ElMessage.success('用户创建成功')
    router.push('/users')
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '创建失败')
  }
}

const goBack = () => {
  router.push('/users')
}
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
</style>