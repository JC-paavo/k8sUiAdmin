import { reactive, computed } from 'vue'

const state = reactive({
  user: null,
  token: localStorage.getItem('token') || null
})

export function useAuthStore() {
  const isLoggedIn = computed(() => !!state.token)
  
  const login = (user, token) => {
    state.user = user
    state.token = token
    localStorage.setItem('token', token)
  }
  
  const logout = () => {
    state.user = null
    state.token = null
    localStorage.removeItem('token')
  }
  
  const setUser = (user) => {
    state.user = user
  }
  
  return {
    user: computed(() => state.user),
    token: computed(() => state.token),
    isLoggedIn,
    login,
    logout,
    setUser
  }
}