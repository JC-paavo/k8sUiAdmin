import axios from 'axios'
import { useAuthStore } from '@/stores/auth'

const api = axios.create({
  baseURL: '/api',
  timeout: 30000
})

api.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore()
    if (authStore.token.value) {
      config.headers.Authorization = `Bearer ${authStore.token.value}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      const authStore = useAuthStore()
      authStore.logout()
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export const authAPI = {
	login: (username, password) => api.post('/login', { username, password }),
	getUser: () => api.get('/user'),
	changePassword: (oldPassword, newPassword) => 
		api.post('/user/change-password', { old_password: oldPassword, new_password: newPassword }),
	checkClusterPermission: (clusterId) => api.get(`/user/clusters/${clusterId}/permission`)
}

export const userAPI = {
  list: (page, pageSize) => api.get('/admin/users', { params: { page, pageSize } }),
  create: (user) => api.post('/admin/users', user),
  get: (id) => api.get(`/admin/users/${id}`),
  update: (id, user) => api.put(`/admin/users/${id}`, user),
  delete: (id) => api.delete(`/admin/users/${id}`)
}

export const clusterAPI = {
  list: () => api.get('/clusters'),
  get: (id) => api.get(`/clusters/${id}`),
  create: (cluster) => api.post('/admin/clusters', cluster),
  update: (id, cluster) => api.put(`/admin/clusters/${id}`, cluster),
  delete: (id) => api.delete(`/admin/clusters/${id}`),
  testConnection: (id) => api.post(`/clusters/test/${id}`),
  refreshStatus: (id) => api.post(`/clusters/refresh/${id}`),
  getPermissions: (id) => api.get(`/admin/clusters/${id}/permissions`),
  addPermission: (id, userId, permission) => 
    api.post(`/admin/clusters/${id}/permissions`, { user_id: userId, permission }),
  removePermission: (id, userId) => 
    api.delete(`/admin/clusters/${id}/permissions/${userId}`)
}

export const k8sAPI = {
  listNamespaces: (clusterId) => api.get(`/k8s/${clusterId}/namespaces`),
  getClusterResourceUsage: (clusterId) => api.get(`/k8s/${clusterId}/cluster/resource-usage`),
  listDeployments: (clusterId, namespace) => 
    api.get(`/k8s/${clusterId}/deployments`, { params: { namespace } }),
  getDeployment: (clusterId, namespace, name) => 
    api.get(`/k8s/${clusterId}/deployments/${namespace}/${name}`),
  createDeployment: (clusterId, data) => 
    api.post(`/k8s/${clusterId}/deployments`, data),
  updateDeployment: (clusterId, namespace, name, data) => 
    api.put(`/k8s/${clusterId}/deployments/${namespace}/${name}`, data),
  deleteDeployment: (clusterId, namespace, name) => 
    api.delete(`/k8s/${clusterId}/deployments/${namespace}/${name}`),
  scaleDeployment: (clusterId, namespace, name, replicas) => 
    api.post(`/k8s/${clusterId}/deployments/${namespace}/${name}/scale`, {}, { params: { replicas } }),
  getDeploymentHistory: (clusterId, namespace, name) => 
    api.get(`/k8s/${clusterId}/deployments/${namespace}/${name}/history`),
  rollbackDeployment: (clusterId, namespace, name, revision) => 
    api.post(`/k8s/${clusterId}/deployments/${namespace}/${name}/rollback`, {}, { params: { revision } }),
  
  listStatefulSets: (clusterId, namespace) => 
    api.get(`/k8s/${clusterId}/statefulsets`, { params: { namespace } }),
  getStatefulSet: (clusterId, namespace, name) => 
    api.get(`/k8s/${clusterId}/statefulsets/${namespace}/${name}`),
  createStatefulSet: (clusterId, data) => 
    api.post(`/k8s/${clusterId}/statefulsets`, data),
  updateStatefulSet: (clusterId, namespace, name, data) => 
    api.put(`/k8s/${clusterId}/statefulsets/${namespace}/${name}`, data),
  deleteStatefulSet: (clusterId, namespace, name) => 
    api.delete(`/k8s/${clusterId}/statefulsets/${namespace}/${name}`),
  scaleStatefulSet: (clusterId, namespace, name, replicas) => 
    api.post(`/k8s/${clusterId}/statefulsets/${namespace}/${name}/scale`, {}, { params: { replicas } }),
  getStatefulSetEvents: (clusterId, namespace, name) => 
    api.get(`/k8s/${clusterId}/statefulsets/${namespace}/${name}/events`),
  
  listServices: (clusterId, namespace) => 
    api.get(`/k8s/${clusterId}/services`, { params: { namespace } }),
  getService: (clusterId, namespace, name) => 
    api.get(`/k8s/${clusterId}/services/${namespace}/${name}`),
  createService: (clusterId, data) => 
    api.post(`/k8s/${clusterId}/services`, data),
  updateService: (clusterId, namespace, name, data) => 
    api.put(`/k8s/${clusterId}/services/${namespace}/${name}`, data),
  deleteService: (clusterId, namespace, name) => 
    api.delete(`/k8s/${clusterId}/services/${namespace}/${name}`),
  
  listIngresses: (clusterId, namespace) => 
    api.get(`/k8s/${clusterId}/ingresses`, { params: { namespace } }),
  getIngress: (clusterId, namespace, name) => 
    api.get(`/k8s/${clusterId}/ingresses/${namespace}/${name}`),
  createIngress: (clusterId, data) => 
    api.post(`/k8s/${clusterId}/ingresses`, data),
  updateIngress: (clusterId, namespace, name, data) => 
    api.put(`/k8s/${clusterId}/ingresses/${namespace}/${name}`, data),
  deleteIngress: (clusterId, namespace, name) => 
    api.delete(`/k8s/${clusterId}/ingresses/${namespace}/${name}`),
  
  listConfigMaps: (clusterId, namespace) => 
    api.get(`/k8s/${clusterId}/configmaps`, { params: { namespace } }),
  getConfigMap: (clusterId, namespace, name) => 
    api.get(`/k8s/${clusterId}/configmaps/${namespace}/${name}`),
  createConfigMap: (clusterId, data) => 
    api.post(`/k8s/${clusterId}/configmaps`, data),
  updateConfigMap: (clusterId, namespace, name, data) => 
    api.put(`/k8s/${clusterId}/configmaps/${namespace}/${name}`, data),
  deleteConfigMap: (clusterId, namespace, name) => 
    api.delete(`/k8s/${clusterId}/configmaps/${namespace}/${name}`),
  
  listSecrets: (clusterId, namespace) => 
    api.get(`/k8s/${clusterId}/secrets`, { params: { namespace } }),
  getSecret: (clusterId, namespace, name) => 
    api.get(`/k8s/${clusterId}/secrets/${namespace}/${name}`),
  createSecret: (clusterId, data) => 
    api.post(`/k8s/${clusterId}/secrets`, data),
  updateSecret: (clusterId, namespace, name, data) => 
    api.put(`/k8s/${clusterId}/secrets/${namespace}/${name}`, data),
  deleteSecret: (clusterId, namespace, name) => 
    api.delete(`/k8s/${clusterId}/secrets/${namespace}/${name}`),
  
  listPods: (clusterId, namespace) => 
    api.get(`/k8s/${clusterId}/pods`, { params: { namespace } }),
  getPod: (clusterId, namespace, name) => 
    api.get(`/k8s/${clusterId}/pods/${namespace}/${name}`),
  deletePod: (clusterId, namespace, name) => 
    api.delete(`/k8s/${clusterId}/pods/${namespace}/${name}`),
  getPodLogs: (clusterId, namespace, name, container) => 
    api.get(`/k8s/${clusterId}/pods/${namespace}/${name}/logs`, { params: { container } }),
  getPodEvents: (clusterId, namespace, name) => 
    api.get(`/k8s/${clusterId}/pods/${namespace}/${name}/events`),
  getDeploymentEvents: (clusterId, namespace, name) => 
    api.get(`/k8s/${clusterId}/deployments/${namespace}/${name}/events`),
  getServiceEvents: (clusterId, namespace, name) => 
    api.get(`/k8s/${clusterId}/services/${namespace}/${name}/events`),
  getIngressEvents: (clusterId, namespace, name) => 
    api.get(`/k8s/${clusterId}/ingresses/${namespace}/${name}/events`),
  getConfigMapEvents: (clusterId, namespace, name) => 
    api.get(`/k8s/${clusterId}/configmaps/${namespace}/${name}/events`),
  getSecretEvents: (clusterId, namespace, name) => 
    api.get(`/k8s/${clusterId}/secrets/${namespace}/${name}/events`),
  getStatefulSetEvents: (clusterId, namespace, name) => 
    api.get(`/k8s/${clusterId}/statefulsets/${namespace}/${name}/events`),
  getDaemonSetEvents: (clusterId, namespace, name) => 
    api.get(`/k8s/${clusterId}/daemonsets/${namespace}/${name}/events`),
  
  listDaemonSets: (clusterId, namespace) => 
    api.get(`/k8s/${clusterId}/daemonsets`, { params: { namespace } }),
  getDaemonSet: (clusterId, namespace, name) => 
    api.get(`/k8s/${clusterId}/daemonsets/${namespace}/${name}`),
  createDaemonSet: (clusterId, data) => 
    api.post(`/k8s/${clusterId}/daemonsets`, data),
  updateDaemonSet: (clusterId, namespace, name, data) => 
    api.put(`/k8s/${clusterId}/daemonsets/${namespace}/${name}`, data),
  deleteDaemonSet: (clusterId, namespace, name) => 
    api.delete(`/k8s/${clusterId}/daemonsets/${namespace}/${name}`),
  getDaemonSetEvents: (clusterId, namespace, name) => 
    api.get(`/k8s/${clusterId}/daemonsets/${namespace}/${name}/events`),
  
  listEvents: (clusterId, namespace) => 
    api.get(`/k8s/${clusterId}/events`, { params: { namespace } })
}

export default api