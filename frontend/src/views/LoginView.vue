<template>
  <div class="min-h-screen bg-gradient-to-br from-primary/5 via-base-100 to-primary/10 flex items-center justify-center px-4 py-12">
    <div class="w-full max-w-md">
      <!-- Logo/Branding -->
      <div class="text-center mb-8">
        <div class="inline-block p-4 bg-primary rounded-2xl shadow-lg shadow-primary/30 mb-4">
          <svg class="w-12 h-12 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
        </div>
        <h1 class="text-3xl font-black">Connexion</h1>
        <p class="text-base-content/60 mt-2">Connectez-vous pour accéder à vos sondages</p>
      </div>

      <!-- Login Card -->
      <div class="card bg-base-100 shadow-2xl border border-base-200">
        <div class="card-body p-8">
          <!-- Google OAuth Button -->
          <button
            @click="googleLogin"
            class="btn btn-outline w-full h-14 gap-3 font-semibold hover:border-primary/50 hover:bg-primary/5 transition-all"
          >
            <svg class="w-5 h-5" viewBox="0 0 24 24">
              <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
              <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
              <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
              <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
            </svg>
            Continuer avec Google
          </button>

          <div class="divider my-6 text-sm text-base-content/50">ou</div>

          <!-- Email/Password form -->
          <form @submit.prevent="handleLogin" class="space-y-5">
            <div class="form-control">
              <label class="label">
                <span class="label-text font-medium">Email</span>
              </label>
              <input
                v-model="email"
                type="email"
                placeholder="vous@exemple.com"
                class="input input-bordered h-12"
                required
              />
            </div>

            <div class="form-control">
              <label class="label">
                <span class="label-text font-medium">Mot de passe</span>
              </label>
              <input
                v-model="password"
                type="password"
                placeholder="••••••••"
                class="input input-bordered h-12"
                required
              />
            </div>

            <div v-if="errorMessage" class="alert alert-error alert-sm">
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <span class="text-sm">{{ errorMessage }}</span>
            </div>

            <button
              type="submit"
              class="btn btn-primary w-full h-12 text-base font-semibold"
              :disabled="loading"
            >
              <span v-if="loading" class="loading loading-spinner"></span>
              <span v-if="!loading">Se connecter</span>
              <span v-else>Connexion...</span>
            </button>
          </form>

          <!-- Sign up link -->
          <p class="text-center mt-6 text-sm text-base-content/60">
            Pas encore de compte ?
            <router-link to="/register" class="link link-primary font-medium">Créer un compte</router-link>
          </p>
        </div>
      </div>

      <!-- Back to home -->
      <div class="text-center mt-6">
        <router-link to="/" class="link link-hover flex items-center justify-center gap-1 text-sm">
          <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
          Retour à l'accueil
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useUiStore } from '@/stores/ui'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const uiStore = useUiStore()

const email = ref('')
const password = ref('')
const loading = ref(false)
const errorMessage = ref('')

async function handleLogin() {
  errorMessage.value = ''
  loading.value = true

  try {
    await authStore.login(email.value, password.value)
    uiStore.success('Bienvenue !')
    const redirect = route.query.redirect || '/dashboard'
    router.push(redirect)
  } catch (error) {
    errorMessage.value = error.response?.data?.error || 'Identifiants incorrects'
  } finally {
    loading.value = false
  }
}

function googleLogin() {
  authStore.googleLogin()
}
</script>
