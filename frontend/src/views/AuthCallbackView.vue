<template>
  <div class="flex flex-col items-center justify-center min-h-64">
    <span class="loading loading-spinner loading-lg"></span>
    <p class="mt-4">Signing you in...</p>
    <p v-if="error" class="mt-4 text-error">{{ error }}</p>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import axios from 'axios'

const router = useRouter()
const route = useRoute()
const error = ref(null)

onMounted(async () => {
  const token = route.query.token

  console.log('Token received:', token ? 'YES' : 'NO')

  if (!token) {
    error.value = 'No token received'
    setTimeout(() => router.push('/login'), 2000)
    return
  }

  try {
    // Store the token in localStorage first
    localStorage.setItem('token', token)
    console.log('Token stored in localStorage')

    // Fetch user data directly with axios
    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    console.log('Fetching user...')

    const response = await axios.get(apiUrl + '/api/auth/me', {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })

    console.log('User data received:', response.data)

    // Store user data in localStorage
    localStorage.setItem('user', JSON.stringify(response.data))

    // Redirect to dashboard - the app will load the auth state from localStorage
    window.location.href = '/dashboard'
  } catch (err) {
    console.error('Auth callback error:', err)
    error.value = err.response?.data?.error || err.message || 'Failed to sign in'
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    setTimeout(() => router.push('/login'), 3000)
  }
})
</script>
