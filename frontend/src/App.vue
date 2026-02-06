<template>
  <div class="min-h-screen flex flex-col bg-slate-50">
    <!-- Premium App Bar -->
    <nav class="sticky top-0 z-50 bg-white/80 backdrop-blur-lg border-b border-slate-100">
      <div class="max-w-md mx-auto px-4 py-3 flex items-center justify-between">
        <!-- Logo -->
        <router-link to="/" class="flex items-center gap-2">
          <div class="w-10 h-10 bg-gradient-to-br from-indigo-500 to-violet-600 rounded-xl flex items-center justify-center text-white shadow-lg">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
              <rect width="18" height="18" x="3" y="4" rx="2" ry="2"/>
              <line x1="16" x2="16" y1="2" y2="6"/>
              <line x1="8" x2="8" y1="2" y2="6"/>
              <line x1="3" x2="21" y1="10" y2="10"/>
            </svg>
          </div>
          <span class="text-xl font-extrabold text-slate-900">Doodle</span>
        </router-link>

        <!-- User Menu / Auth Buttons -->
        <div class="flex items-center gap-2">
          <!-- Dark Mode Toggle -->
          <button
            @click="toggleTheme"
            class="w-10 h-10 bg-slate-100 dark:bg-slate-800 rounded-xl flex items-center justify-center hover:bg-slate-200 dark:hover:bg-slate-700 transition-colors"
            :title="isDark ? 'Mode clair' : 'Mode sombre'"
          >
            <svg v-if="isDark" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-amber-400">
              <circle cx="12" cy="12" r="5"/>
              <path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/>
            </svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-slate-600">
              <path d="M12 3a6 6 0 0 0 9 9 9 9 0 1 1-9-9Z"/>
            </svg>
          </button>

          <template v-if="!authStore.isAuthenticated">
            <router-link to="/login" class="px-4 py-2 text-sm font-bold text-slate-600 hover:text-indigo-600 transition-colors">
              Connexion
            </router-link>
            <router-link to="/register" class="px-4 py-2 bg-gradient-to-r from-indigo-600 to-violet-600 text-white text-sm font-bold rounded-xl shadow-lg shadow-indigo-200 hover:shadow-xl transition-all active:scale-95">
              S'inscrire
            </router-link>
          </template>
          <template v-else>
            <router-link to="/create" class="w-10 h-10 bg-indigo-50 text-indigo-600 rounded-xl flex items-center justify-center hover:bg-indigo-100 transition-colors">
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <path d="M12 5v14"/>
                <path d="M5 12h14"/>
              </svg>
            </router-link>

            <!-- User Dropdown -->
            <div class="dropdown dropdown-end" ref="userDropdown">
              <div tabindex="0" role="button" class="flex items-center gap-2 p-1 pr-3 bg-slate-50 rounded-full hover:bg-slate-100 transition-colors cursor-pointer">
                <div class="w-8 h-8 rounded-full bg-gradient-to-br from-indigo-500 to-violet-600 flex items-center justify-center text-white font-bold text-sm">
                  {{ userInitial }}
                </div>
                <span class="text-sm font-bold text-slate-700 max-w-20 truncate">{{ authStore.user?.name }}</span>
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-slate-400">
                  <path d="m6 9 6 6 6-6"/>
                </svg>
              </div>
              <ul tabindex="0" class="dropdown-content z-[100] menu p-2 shadow-xl bg-white rounded-2xl w-52 border border-slate-100 mt-2">
                <li class="px-2 py-1">
                  <div class="flex items-center gap-3 p-2 bg-slate-50 rounded-xl">
                    <div class="w-10 h-10 rounded-full bg-gradient-to-br from-indigo-500 to-violet-600 flex items-center justify-center text-white font-bold">
                      {{ userInitial }}
                    </div>
                    <div class="flex-1 min-w-0">
                      <p class="font-bold text-slate-800 text-sm truncate">{{ authStore.user?.name }}</p>
                      <p class="text-xs text-slate-400 truncate">{{ authStore.user?.email }}</p>
                    </div>
                  </div>
                </li>
                <li><hr class="my-1 border-slate-100"></li>
                <li>
                  <router-link to="/dashboard" class="flex items-center gap-3 text-slate-700 hover:text-indigo-600 hover:bg-indigo-50 rounded-xl">
                    <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <rect width="7" height="9" x="3" y="3" rx="1"/>
                      <rect width="7" height="5" x="14" y="3" rx="1"/>
                      <rect width="7" height="9" x="14" y="12" rx="1"/>
                      <rect width="7" height="5" x="3" y="16" rx="1"/>
                    </svg>
                    <span class="font-medium">Tableau de bord</span>
                  </router-link>
                </li>
                <li>
                  <router-link to="/profile" class="flex items-center gap-3 text-slate-700 hover:text-indigo-600 hover:bg-indigo-50 rounded-xl">
                    <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"/>
                      <circle cx="12" cy="7" r="4"/>
                    </svg>
                    <span class="font-medium">Profil</span>
                  </router-link>
                </li>
                <li><hr class="my-1 border-slate-100"></li>
                <li>
                  <a @click.prevent="handleLogout" class="flex items-center gap-3 text-rose-600 hover:bg-rose-50 rounded-xl">
                    <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
                      <polyline points="16 17 21 12 16 7"/>
                      <line x1="21" x2="9" y1="12" y2="12"/>
                    </svg>
                    <span class="font-medium">Déconnexion</span>
                  </a>
                </li>
              </ul>
            </div>
          </template>
        </div>
      </div>
    </nav>

    <!-- Main content -->
    <main class="flex-1">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>

    <!-- Toast container -->
    <div class="toast toast-end z-[200]" id="toast-container"></div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useUiStore } from '@/stores/ui'

const router = useRouter()
const authStore = useAuthStore()
const uiStore = useUiStore()

// Dark mode
const isDark = ref(false)

// Load theme preference on mount
onMounted(() => {
  authStore.loadFromStorage()

  // Load theme from localStorage or system preference
  const savedTheme = localStorage.getItem('theme')
  if (savedTheme) {
    isDark.value = savedTheme === 'dark'
  } else {
    // Check system preference
    isDark.value = window.matchMedia('(prefers-color-scheme: dark)').matches
  }
  applyTheme()
})

// Watch for changes and save to localStorage
watch(isDark, () => {
  applyTheme()
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
})

function applyTheme() {
  const html = document.documentElement
  if (isDark.value) {
    html.setAttribute('data-theme', 'dark')
  } else {
    html.setAttribute('data-theme', 'light')
  }
}

function toggleTheme() {
  isDark.value = !isDark.value
}

const userInitial = computed(() => {
  return authStore.user?.name?.charAt(0)?.toUpperCase() || '?'
})

async function handleLogout() {
  await authStore.logout()
  uiStore.success('À bientôt !')
  router.push('/')
}
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

:deep(.dropdown-content li a) {
  border-radius: 0.75rem;
}

:deep(.dropdown-content li hr) {
  margin: 0.25rem 0.5rem;
}
</style>
