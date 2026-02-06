<template>
  <div class="max-w-2xl mx-auto space-y-6">
    <h1 class="text-3xl font-bold">Profile</h1>

    <!-- Profile Card -->
    <div class="card bg-base-100 shadow-xl">
      <div class="card-body">
        <div class="flex items-center gap-4 mb-6">
          <div class="avatar">
            <div class="w-24 rounded-full">
              <img :src="authStore.user?.avatar || `https://api.dicebear.com/7.x/avataaars/svg?seed=${authStore.user?.name}`" />
            </div>
          </div>
          <div>
            <h2 class="text-2xl font-bold">{{ authStore.user?.name }}</h2>
            <p class="text-base-content/70">{{ authStore.user?.email }}</p>
            <div class="badge badge-ghost mt-1">{{ authStore.user?.provider || 'email' }}</div>
          </div>
        </div>

        <form @submit.prevent="updateProfile" class="space-y-4">
          <div class="form-control">
            <label class="label">
              <span class="label-text">Name</span>
            </label>
            <input
              v-model="profileForm.name"
              type="text"
              class="input input-bordered"
            />
          </div>

          <div class="form-control">
            <label class="label">
              <span class="label-text">Email</span>
            </label>
            <input
              v-model="profileForm.email"
              type="email"
              class="input input-bordered"
            />
          </div>

          <button type="submit" class="btn btn-primary" :disabled="loading">
            <span v-if="loading" class="loading loading-spinner"></span>
            Save Changes
          </button>
        </form>
      </div>
    </div>

    <!-- Change Password -->
    <div v-if="authStore.user?.provider === 'email'" class="card bg-base-100 shadow-xl">
      <div class="card-body">
        <h2 class="card-title">Change Password</h2>

        <form @submit.prevent="changePassword" class="space-y-4">
          <div class="form-control">
            <label class="label">
              <span class="label-text">Current Password</span>
            </label>
            <input
              v-model="passwordForm.old_password"
              type="password"
              class="input input-bordered"
              required
            />
          </div>

          <div class="form-control">
            <label class="label">
              <span class="label-text">New Password</span>
            </label>
            <input
              v-model="passwordForm.new_password"
              type="password"
              class="input input-bordered"
              minlength="8"
              required
            />
          </div>

          <div v-if="passwordError" class="alert alert-error">
            <span>{{ passwordError }}</span>
          </div>

          <button type="submit" class="btn btn-warning" :disabled="loading">
            <span v-if="loading" class="loading loading-spinner"></span>
            Change Password
          </button>
        </form>
      </div>
    </div>

    <!-- Danger Zone -->
    <div class="card bg-base-100 shadow-xl border border-error/20">
      <div class="card-body">
        <h2 class="card-title text-error">Danger Zone</h2>
        <p class="text-sm text-base-content/70">Once you delete your account, there is no going back. Please be certain.</p>

        <button @click="deleteAccount" class="btn btn-error btn-outline">
          Delete Account
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useUiStore } from '@/stores/ui'

const router = useRouter()
const authStore = useAuthStore()
const uiStore = useUiStore()

const loading = ref(false)
const passwordError = ref('')

const profileForm = ref({
  name: '',
  email: ''
})

const passwordForm = ref({
  old_password: '',
  new_password: ''
})

onMounted(() => {
  profileForm.value = {
    name: authStore.user?.name || '',
    email: authStore.user?.email || ''
  }
})

async function updateProfile() {
  loading.value = true
  try {
    await authStore.updateProfile(profileForm.value.name, profileForm.value.email)
    uiStore.success('Profile updated successfully!')
  } catch (error) {
    uiStore.error(error.response?.data?.error || 'Failed to update profile')
  } finally {
    loading.value = false
  }
}

async function changePassword() {
  loading.value = true
  passwordError.value = ''

  try {
    await authStore.changePassword(passwordForm.value.old_password, passwordForm.value.new_password)
    uiStore.success('Password changed successfully!')
    passwordForm.value = { old_password: '', new_password: '' }
  } catch (error) {
    passwordError.value = error.response?.data?.error || 'Failed to change password'
  } finally {
    loading.value = false
  }
}

function deleteAccount() {
  if (confirm('Are you sure you want to delete your account? This action cannot be undone.')) {
    // TODO: Implement delete account API
    uiStore.info('Account deletion will be implemented soon')
  }
}
</script>
