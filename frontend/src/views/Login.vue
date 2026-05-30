<template>
  <div class="login-container">
    <div class="login-background">
      <div class="grid-bg"></div>
      <div class="glow-orb orb-1"></div>
      <div class="glow-orb orb-2"></div>
      <div class="glow-orb orb-3"></div>
    </div>

    <div class="login-content">
      <div class="login-card">
        <div class="login-header">
          <div class="logo-section">
            <div class="logo-icon">
              <svg width="48" height="48" viewBox="0 0 48 48" fill="none">
                <circle cx="24" cy="24" r="22" stroke="url(#logoGradient)" stroke-width="2" fill="none"/>
                <path d="M24 12 L24 36 M12 24 L36 24 M16 16 L32 32 M32 16 L16 32" stroke="url(#logoGradient)" stroke-width="2" stroke-linecap="round"/>
                <defs>
                  <linearGradient id="logoGradient" x1="0%" y1="0%" x2="100%" y2="100%">
                    <stop offset="0%" stop-color="#3b82f6"/>
                    <stop offset="100%" stop-color="#10b981"/>
                  </linearGradient>
                </defs>
              </svg>
            </div>
            <div class="logo-text">
              <h1>K8s 管理平台</h1>
              <p>Kubernetes 集群管理控制台</p>
            </div>
          </div>
        </div>

        <div class="login-form-section">
          <h2>登录</h2>
          <p class="form-subtitle">访问您的集群管理控制台</p>

          <el-form ref="loginForm" :model="form" :rules="rules" @submit.prevent="handleLogin">
            <div class="form-group">
              <label for="username">用户名</label>
              <el-input
                id="username"
                v-model="form.username"
                placeholder="请输入用户名"
                size="large"
                :prefix-icon="User"
                @keyup.enter="handleLogin"
              />
            </div>

            <div class="form-group">
              <label for="password">密码</label>
              <el-input
                id="password"
                v-model="form.password"
                type="password"
                placeholder="请输入密码"
                size="large"
                :prefix-icon="Lock"
                show-password
                @keyup.enter="handleLogin"
              />
            </div>

            <el-button
              type="primary"
              class="login-button"
              size="large"
              :loading="loading"
              @click="handleLogin"
            >
              <span v-if="!loading">登录</span>
              <span v-else>登录中...</span>
            </el-button>

            <div class="login-tip">
              <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
                <path d="M8 1a7 7 0 100 14A7 7 0 008 1zm0 12.5a5.5 5.5 0 110-11 5.5 5.5 0 010 11zM8 4a.75.75 0 00-.75.75v3.5a.75.75 0 001.5 0v-3.5A.75.75 0 008 4zm0 8a1 1 0 100-2 1 1 0 000 2z"/>
              </svg>
              <span>默认账号：<strong>admin / admin</strong></span>
            </div>
          </el-form>
        </div>

        <div class="login-footer">
          <div class="feature-badges">
            <span class="badge">安全</span>
            <span class="badge">快速</span>
            <span class="badge">可靠</span>
          </div>
        </div>
      </div>

      <div class="decoration">
        <div class="deco-line"></div>
        <div class="deco-dots">
          <span></span><span></span><span></span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import { User, Lock } from '@element-plus/icons-vue'
import { authAPI } from '@/utils/api'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()
const loading = ref(false)

const form = reactive({
  username: '',
  password: ''
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  if (!form.username || !form.password) {
    ElMessage.warning('请输入用户名和密码')
    return
  }

  loading.value = true
  try {
    const response = await authAPI.login(form.username, form.password)
    const { user, token } = response.data
    authStore.login(user, token)
    ElMessage.success('欢迎回来！')
    router.push('/')
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '登录失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-primary);
  position: relative;
  overflow: hidden;
}

.login-background {
  position: absolute;
  inset: 0;
  pointer-events: none;
}

.grid-bg {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(59, 130, 246, 0.03) 1px, transparent 1px),
    linear-gradient(90deg, rgba(59, 130, 246, 0.03) 1px, transparent 1px);
  background-size: 50px 50px;
  animation: gridMove 20s linear infinite;
}

@keyframes gridMove {
  0% { transform: perspective(500px) rotateX(60deg) translateY(0); }
  100% { transform: perspective(500px) rotateX(60deg) translateY(50px); }
}

.glow-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(100px);
  opacity: 0.4;
  animation: orbFloat 15s ease-in-out infinite;
}

.orb-1 {
  width: 600px;
  height: 600px;
  background: radial-gradient(circle, rgba(59, 130, 246, 0.3), transparent 70%);
  top: -200px;
  right: -200px;
}

.orb-2 {
  width: 400px;
  height: 400px;
  background: radial-gradient(circle, rgba(16, 185, 129, 0.3), transparent 70%);
  bottom: -100px;
  left: -100px;
  animation-delay: -5s;
}

.orb-3 {
  width: 300px;
  height: 300px;
  background: radial-gradient(circle, rgba(99, 102, 241, 0.3), transparent 70%);
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  animation-delay: -10s;
}

@keyframes orbFloat {
  0%, 100% { transform: translate(0, 0) scale(1); }
  33% { transform: translate(30px, -30px) scale(1.1); }
  66% { transform: translate(-20px, 20px) scale(0.9); }
}

.login-content {
  position: relative;
  z-index: 10;
  animation: fadeInUp 0.8s ease;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.login-card {
  width: 440px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-xl);
  padding: 48px;
  box-shadow: var(--shadow-xl);
  backdrop-filter: blur(20px);
}

.login-header {
  margin-bottom: 40px;
}

.logo-section {
  display: flex;
  align-items: center;
  gap: 16px;
}

.logo-icon {
  width: 64px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, var(--primary-color), var(--success-color));
  border-radius: var(--radius-lg);
  animation: logoPulse 3s ease-in-out infinite;
}

@keyframes logoPulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.05); }
}

.logo-text h1 {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 4px;
  background: linear-gradient(135deg, var(--primary-light), var(--success-light));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.logo-text p {
  font-size: 12px;
  color: var(--text-tertiary);
  letter-spacing: 0.5px;
}

.login-form-section h2 {
  font-size: 28px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.form-subtitle {
  font-size: 14px;
  color: var(--text-tertiary);
  margin-bottom: 32px;
}

.form-group {
  margin-bottom: 24px;
}

.form-group label {
  display: block;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-secondary);
  margin-bottom: 8px;
}

.login-button {
  width: 100%;
  height: 48px;
  font-size: 16px;
  font-weight: 600;
  background: linear-gradient(135deg, var(--primary-color), var(--primary-dark));
  border: none;
  border-radius: var(--radius-md);
  margin-top: 8px;
  transition: all var(--transition-base);
  position: relative;
  overflow: hidden;
}

.login-button::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, var(--primary-light), var(--success-color));
  opacity: 0;
  transition: opacity var(--transition-fast);
}

.login-button:hover::before {
  opacity: 1;
}

.login-button span {
  position: relative;
  z-index: 1;
}

.login-tip {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 24px;
  padding: 12px 16px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  font-size: 13px;
  color: var(--text-tertiary);
}

.login-tip svg {
  flex-shrink: 0;
  color: var(--primary-color);
}

.login-tip strong {
  color: var(--primary-light);
}

.login-footer {
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid var(--border-color);
}

.feature-badges {
  display: flex;
  justify-content: center;
  gap: 12px;
}

.badge {
  padding: 4px 12px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  font-size: 11px;
  font-weight: 600;
  color: var(--text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.decoration {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 40px;
  animation: fadeIn 1s ease 0.3s backwards;
}

.deco-line {
  width: 60px;
  height: 2px;
  background: linear-gradient(90deg, transparent, var(--primary-color), transparent);
}

.deco-dots {
  display: flex;
  gap: 8px;
  margin-left: 16px;
}

.deco-dots span {
  width: 6px;
  height: 6px;
  background: var(--primary-color);
  border-radius: 50%;
  animation: dotPulse 2s ease-in-out infinite;
}

.deco-dots span:nth-child(2) {
  animation-delay: 0.3s;
}

.deco-dots span:nth-child(3) {
  animation-delay: 0.6s;
}

@keyframes dotPulse {
  0%, 100% { opacity: 0.3; transform: scale(1); }
  50% { opacity: 1; transform: scale(1.2); }
}

/* Responsive */
@media (max-width: 768px) {
  .login-card {
    width: 100%;
    margin: 0 20px;
    padding: 32px 24px;
  }

  .logo-section {
    flex-direction: column;
    text-align: center;
  }
}
</style>
