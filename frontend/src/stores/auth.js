import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import api from '@/utils/api'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(null)
  const token = ref(localStorage.getItem('token') || null)
  const loading = ref(false)

  const isAuthenticated = computed(() => !!token.value && !!user.value)

  // Load user from localStorage if token exists
  const savedUser = localStorage.getItem('user')
  if (savedUser && token.value) {
    try {
      user.value = JSON.parse(savedUser)
    } catch (e) {
      console.error('Failed to parse saved user', e)
    }
  }

  // Sync token to localStorage
  watch(token, (newToken) => {
    if (newToken) {
      localStorage.setItem('token', newToken)
    } else {
      localStorage.removeItem('token')
    }
  })

  // Sync user to localStorage
  watch(user, (newUser) => {
    if (newUser) {
      localStorage.setItem('user', JSON.stringify(newUser))
    } else {
      localStorage.removeItem('user')
    }
  })

  async function login(email, password) {
    loading.value = true
    try {
      const response = await api.post('/api/auth/login', { email, password })
      token.value = response.data.token
      user.value = response.data.user
      localStorage.setItem('token', response.data.token)
      localStorage.setItem('refresh_token', response.data.refresh_token)
      return response.data
    } catch (error) {
      throw error
    } finally {
      loading.value = false
    }
  }

  async function register(name, email, password) {
    loading.value = true
    try {
      const response = await api.post('/api/auth/register', { name, email, password })
      token.value = response.data.token
      user.value = response.data.user
      localStorage.setItem('token', response.data.token)
      localStorage.setItem('refresh_token', response.data.refresh_token)
      return response.data
    } catch (error) {
      throw error
    } finally {
      loading.value = false
    }
  }

  async function logout() {
    try {
      await api.post('/api/auth/logout')
    } catch (error) {
      console.error('Logout error:', error)
    } finally {
      token.value = null
      user.value = null
      localStorage.removeItem('token')
      localStorage.removeItem('refresh_token')
    }
  }

  async function fetchUser() {
    if (!token.value) return

    try {
      const response = await api.get('/api/auth/me')
      user.value = response.data
    } catch (error) {
      // Token might be expired
      await logout()
    }
  }

  async function refreshToken() {
    const refreshToken = localStorage.getItem('refresh_token')
    if (!refreshToken) {
      await logout()
      return false
    }

    try {
      const response = await api.post('/api/auth/refresh')
      token.value = response.data.token
      user.value = response.data.user
      localStorage.setItem('token', response.data.token)
      localStorage.setItem('refresh_token', response.data.refresh_token)
      return true
    } catch (error) {
      await logout()
      return false
    }
  }

  async function updateProfile(name, email) {
    try {
      const response = await api.put('/auth/profile', { name, email })
      user.value = { ...user.value, ...response.data }
      return response.data
    } catch (error) {
      throw error
    }
  }

  async function changePassword(oldPassword, newPassword) {
    try {
      const response = await api.put('/auth/password', {
        old_password: oldPassword,
        new_password: newPassword
      })
      return response.data
    } catch (error) {
      throw error
    }
  }

  function loadFromStorage() {
    const savedToken = localStorage.getItem('token')
    if (savedToken) {
      token.value = savedToken
      fetchUser()
    }
  }

  function googleLogin() {
    window.location.href = `${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/auth/google/login`
  }

  return {
    user,
    token,
    loading,
    isAuthenticated,
    login,
    register,
    logout,
    fetchUser,
    refreshToken,
    updateProfile,
    changePassword,
    loadFromStorage,
    googleLogin
  }
})
