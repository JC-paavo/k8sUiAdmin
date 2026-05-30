import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { guest: true }
  },
  {
    path: '/',
    name: 'Layout',
    component: () => import('@/layout/Layout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue')
      },
      {
        path: 'clusters',
        name: 'ClusterList',
        component: () => import('@/views/cluster/ClusterList.vue')
      },
      {
        path: 'clusters/:id',
        name: 'ClusterDetail',
        component: () => import('@/views/cluster/ClusterDetail.vue')
      },
      {
        path: 'clusters/:id/deployments',
        name: 'DeploymentList',
        component: () => import('@/views/cluster/DeploymentList.vue')
      },
      {
        path: 'clusters/:id/deployments/:namespace/:name',
        name: 'DeploymentDetail',
        component: () => import('@/views/cluster/ResourceDetail.vue')
      },
      {
        path: 'clusters/:id/pods',
        name: 'PodList',
        component: () => import('@/views/cluster/PodList.vue')
      },
      {
        path: 'clusters/:id/pods/:namespace/:name',
        name: 'PodDetail',
        component: () => import('@/views/cluster/ResourceDetail.vue')
      },
      {
        path: 'clusters/:id/services',
        name: 'ServiceList',
        component: () => import('@/views/cluster/ServiceList.vue')
      },
      {
        path: 'clusters/:id/services/:namespace/:name',
        name: 'ServiceDetail',
        component: () => import('@/views/cluster/ResourceDetail.vue')
      },
      {
        path: 'clusters/:id/ingresses',
        name: 'IngressList',
        component: () => import('@/views/cluster/IngressList.vue')
      },
      {
        path: 'clusters/:id/ingresses/:namespace/:name',
        name: 'IngressDetail',
        component: () => import('@/views/cluster/ResourceDetail.vue')
      },
      {
        path: 'clusters/:id/configmaps',
        name: 'ConfigMapList',
        component: () => import('@/views/cluster/ConfigMapList.vue')
      },
      {
        path: 'clusters/:id/configmaps/:namespace/:name',
        name: 'ConfigMapDetail',
        component: () => import('@/views/cluster/ResourceDetail.vue')
      },
      {
        path: 'clusters/:id/secrets',
        name: 'SecretList',
        component: () => import('@/views/cluster/SecretList.vue')
      },
      {
        path: 'clusters/:id/secrets/:namespace/:name',
        name: 'SecretDetail',
        component: () => import('@/views/cluster/ResourceDetail.vue')
      },
      {
        path: 'clusters/:id/statefulsets',
        name: 'StatefulSetList',
        component: () => import('@/views/cluster/StatefulSetList.vue')
      },
      {
        path: 'clusters/:id/statefulsets/:namespace/:name',
        name: 'StatefulSetDetail',
        component: () => import('@/views/cluster/ResourceDetail.vue')
      },
      {
        path: 'clusters/:id/daemonsets',
        name: 'DaemonSetList',
        component: () => import('@/views/cluster/DaemonSetList.vue')
      },
      {
        path: 'clusters/:id/daemonsets/:namespace/:name',
        name: 'DaemonSetDetail',
        component: () => import('@/views/cluster/ResourceDetail.vue')
      },
      {
        path: 'clusters/:id/resources/:type',
        name: 'ResourceList',
        component: () => import('@/views/cluster/ResourceList.vue')
      },
      {
        path: 'clusters/:id/resources/:type/:namespace/:name',
        name: 'ResourceDetail',
        component: () => import('@/views/cluster/ResourceDetail.vue')
      },
      {
        path: 'users',
        name: 'UserList',
        component: () => import('@/views/user/UserList.vue'),
        meta: { requiresAdmin: true }
      },
      {
        path: 'users/create',
        name: 'UserCreate',
        component: () => import('@/views/user/UserCreate.vue'),
        meta: { requiresAdmin: true }
      },
      {
        path: 'users/:id/edit',
        name: 'UserEdit',
        component: () => import('@/views/user/UserEdit.vue'),
        meta: { requiresAdmin: true }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  
  // 如果用户已登录，尝试获取用户信息
  if (authStore.token.value && !authStore.user.value) {
    try {
      const { authAPI } = await import('@/utils/api')
      const response = await authAPI.getUser()
      authStore.setUser(response.data)
    } catch (error) {
      // Token无效，清除登录状态
      authStore.logout()
      if (!to.meta.guest) {
        next('/login')
        return
      }
    }
  }
  
  // 检查是否需要登录
  if (to.meta.requiresAuth) {
    if (!authStore.token.value) {
      next('/login')
      return
    }
  }
  
  // 检查是否需要管理员权限
  if (to.meta.requiresAdmin) {
    if (!authStore.token.value || authStore.user.value?.role !== 'admin') {
      next('/')
      return
    }
  }
  
  // 如果已登录且访问login页面，重定向到首页
  if (to.meta.guest && authStore.token.value) {
    next('/')
    return
  }
  
  next()
})

export default router
