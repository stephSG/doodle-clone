import axios from 'axios'
import { useAuthStore } from '@/stores/auth'

const api = axios.create({
  baseURL: '/',  // Use relative path to leverage Vite proxy
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor to add auth token
api.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore()
    if (authStore.token) {
      config.headers.Authorization = `Bearer ${authStore.token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor to handle token refresh
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const authStore = useAuthStore()

    // Don't retry on logout or refresh endpoints
    if (error.config?.url?.includes('/auth/logout') ||
        error.config?.url?.includes('/auth/refresh')) {
      return Promise.reject(error)
    }

    // If 401 and not already retrying
    if (error.response?.status === 401 && !error.config._retry) {
      error.config._retry = true

      try {
        await authStore.refreshToken()
        // Retry the original request with new token
        error.config.headers.Authorization = `Bearer ${authStore.token}`
        return api.request(error.config)
      } catch (refreshError) {
        // Refresh failed, redirect to login
        authStore.logout()
        window.location.href = '/login'
        return Promise.reject(refreshError)
      }
    }

    return Promise.reject(error)
  }
)

export default api
