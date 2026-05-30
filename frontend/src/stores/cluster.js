import { reactive, computed } from 'vue'

const state = reactive({
  clusters: [],
  currentCluster: null
})

export function useClusterStore() {
  const clusters = computed(() => state.clusters)
  const currentCluster = computed(() => state.currentCluster)
  
  const setClusters = (clusters) => {
    state.clusters = clusters
  }
  
  const setCurrentCluster = (cluster) => {
    state.currentCluster = cluster
  }
  
  const addCluster = (cluster) => {
    state.clusters.push(cluster)
  }
  
  const updateCluster = (updatedCluster) => {
    const index = state.clusters.findIndex(c => c.id === updatedCluster.id)
    if (index !== -1) {
      state.clusters[index] = updatedCluster
      if (state.currentCluster?.id === updatedCluster.id) {
        state.currentCluster = updatedCluster
      }
    }
  }
  
  const removeCluster = (clusterId) => {
    state.clusters = state.clusters.filter(c => c.id !== clusterId)
    if (state.currentCluster?.id === clusterId) {
      state.currentCluster = null
    }
  }
  
  return {
    clusters,
    currentCluster,
    setClusters,
    setCurrentCluster,
    addCluster,
    updateCluster,
    removeCluster
  }
}