<template>
  <div class="layout-container">
    <div class="sidebar">
      <div class="sidebar-header">
        <div class="logo">
          <div class="logo-icon">
            <svg width="32" height="32" viewBox="0 0 48 48" fill="none">
              <circle cx="24" cy="24" r="22" stroke="url(#sidebarLogoGradient)" stroke-width="2" fill="none"/>
              <path d="M24 12 L24 36 M12 24 L36 24 M16 16 L32 32 M32 16 L16 32" stroke="url(#sidebarLogoGradient)" stroke-width="2" stroke-linecap="round"/>
              <defs>
                <linearGradient id="sidebarLogoGradient" x1="0%" y1="0%" x2="100%" y2="100%">
                  <stop offset="0%" stop-color="#3b82f6"/>
                  <stop offset="100%" stop-color="#10b981"/>
                </linearGradient>
              </defs>
            </svg>
          </div>
          <div class="logo-text">
            <h1>K8s 管理平台</h1>
            <p>ADMIN</p>
          </div>
        </div>
      </div>

      <nav class="sidebar-nav">
        <router-link to="/" class="nav-item" :class="{ active: route.path === '/' }">
          <el-icon><HomeFilled /></el-icon>
          <span>Dashboard</span>
        </router-link>
        
        <router-link to="/clusters" class="nav-item" :class="{ active: route.path.startsWith('/clusters') }">
          <el-icon><Grid /></el-icon>
          <span>Clusters</span>
        </router-link>
        
        <router-link v-if="isAdmin" to="/users" class="nav-item" :class="{ active: route.path.startsWith('/users') }">
          <el-icon><User /></el-icon>
          <span>用户管理</span>
        </router-link>
      </nav>

      <div class="sidebar-footer">
        <div class="user-info">
          <div class="user-avatar">
            {{ userName.charAt(0).toUpperCase() }}
          </div>
          <div class="user-details">
            <div class="user-name">{{ userName }}</div>
            <div class="user-role">{{ userRole }}</div>
          </div>
        </div>
        
        <el-dropdown @command="handleCommand" trigger="click">
          <el-button class="user-menu-btn" circle>
            <el-icon><MoreFilled /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="change-password">
                <el-icon><Key /></el-icon>
                Change Password
              </el-dropdown-item>
              <el-dropdown-item command="logout" divided>
                <el-icon><SwitchButton /></el-icon>
                Sign Out
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>

    <div class="main-content">
      <router-view />
    </div>

    <ChangePassword v-model="showChangePassword" />
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { HomeFilled, Grid, User, MoreFilled, Key, SwitchButton } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import ChangePassword from './ChangePassword.vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const showChangePassword = ref(false)
const isAdmin = computed(() => authStore.user.value?.role === 'admin')
const userName = computed(() => authStore.user.value?.username || '')
const userRole = computed(() => authStore.user.value?.role === 'admin' ? 'Administrator' : 'User')

const handleCommand = async (command) => {
  if (command === 'change-password') {
    showChangePassword.value = true
  } else if (command === 'logout') {
    try {
      await ElMessageBox.confirm(
        '您确定要退出登录吗？',
        '确认退出',
        {
          confirmButtonText: '退出',
          cancelButtonText: '取消',
          type: 'warning'
        }
      )
      authStore.logout()
      ElMessage.success('退出登录成功')
      router.push('/login')
    } catch {
      // User cancelled
    }
  }
}
</script>

<style scoped>
.layout-container {
  display: flex;
  min-height: 100vh;
  background: var(--bg-primary);
}

.sidebar {
  width: 260px;
  background: var(--bg-card);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  position: fixed;
  top: 0;
  left: 0;
  height: 100vh;
  z-index: 100;
}

.sidebar-header {
  padding: 24px;
  border-bottom: 1px solid var(--border-color);
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-icon {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, var(--primary-color), var(--success-color));
  border-radius: var(--radius-md);
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
}

.logo-text h1 {
  font-size: 20px;
  font-weight: 700;
  color: var(--text-primary);
  background: linear-gradient(135deg, var(--primary-light), var(--success-light));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 2px;
}

.logo-text p {
  font-size: 12px;
  color: var(--text-tertiary);
  letter-spacing: 1px;
}

.sidebar-nav {
  flex: 1;
  padding: 16px 12px;
  overflow-y: auto;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  margin-bottom: 4px;
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  text-decoration: none;
  transition: all var(--transition-fast);
  position: relative;
  overflow: hidden;
}

.nav-item::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 3px;
  background: linear-gradient(180deg, var(--primary-color), var(--success-color));
  opacity: 0;
  transition: opacity var(--transition-fast);
}

.nav-item:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.nav-item.active {
  background: linear-gradient(90deg, rgba(59, 130, 246, 0.1), transparent);
  color: var(--primary-light);
}

.nav-item.active::before {
  opacity: 1;
}

.nav-item .el-icon {
  font-size: 20px;
  flex-shrink: 0;
}

.nav-item span {
  font-size: 14px;
  font-weight: 500;
}

.sidebar-footer {
  padding: 16px;
  border-top: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-info {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
}

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: var(--radius-md);
  background: linear-gradient(135deg, var(--primary-color), var(--info-color));
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  font-weight: 700;
  color: white;
  flex-shrink: 0;
  box-shadow: 0 2px 8px rgba(59, 130, 246, 0.3);
}

.user-details {
  flex: 1;
  min-width: 0;
}

.user-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 2px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.user-role {
  font-size: 11px;
  color: var(--text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.user-menu-btn {
  width: 36px;
  height: 36px;
  border: 1px solid var(--border-color);
  background: transparent;
  color: var(--text-secondary);
  border-radius: var(--radius-md);
  transition: all var(--transition-fast);
}

.user-menu-btn:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
  border-color: var(--border-light);
}

.main-content {
  flex: 1;
  margin-left: 260px;
  min-height: 100vh;
  padding: 32px;
  animation: fadeIn 0.3s ease;
}
</style>
