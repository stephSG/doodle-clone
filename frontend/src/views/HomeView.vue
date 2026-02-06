<template>
  <div class="min-h-screen">
    <!-- Header -->
    <div class="bg-gradient-to-r from-primary/10 to-primary/5 py-16 px-4">
      <div class="max-w-4xl mx-auto">
        <div class="flex justify-between items-center mb-8">
          <div>
            <h1 class="text-4xl md:text-5xl font-black mb-2">Sondages</h1>
            <p class="text-base-content/70">Planifiez vos événements facilement</p>
          </div>
          <router-link
            to="/create"
            class="btn btn-primary btn-lg gap-2 shadow-xl"
          >
            <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            Créer
          </router-link>
        </div>

        <!-- Search -->
        <div class="relative">
          <svg class="w-5 h-5 absolute left-4 top-1/2 -translate-y-1/2 text-base-content/40" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input
            type="text"
            v-model="searchQuery"
            @input="debouncedSearch"
            placeholder="Rechercher un sondage..."
            class="input input-bordered w-full pl-12 h-14 text-lg shadow-sm"
          />
        </div>
      </div>
    </div>

    <!-- Content -->
    <div class="max-w-4xl mx-auto px-4 py-8">
      <!-- Loading -->
      <div v-if="loading" class="flex justify-center py-20">
        <span class="loading loading-spinner loading-lg"></span>
      </div>

      <!-- Polls List -->
      <div v-else class="space-y-6">
        <div class="flex justify-between items-center">
          <h2 class="text-xl font-bold">
            {{ searchQuery ? 'Résultats' : 'Sondages récents' }}
          </h2>
          <span class="text-sm text-base-content/50">{{ polls.length }} sondage{{ polls.length > 1 ? 's' : '' }}</span>
        </div>

        <!-- Empty State -->
        <div v-if="polls.length === 0" class="text-center py-20">
          <div class="inline-block p-6 bg-base-200 rounded-full mb-6">
            <svg class="w-16 h-16 text-base-content/30" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" />
            </svg>
          </div>
          <h3 class="text-xl font-bold mb-2">{{ searchQuery ? 'Aucun résultat' : 'Aucun sondage' }}</h3>
          <p class="text-base-content/50 mb-6">{{ searchQuery ? 'Essayez d autres termes de recherche' : 'Créez votre premier sondage pour commencer' }}</p>
          <router-link to="/create" class="btn btn-primary gap-2">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            Créer un sondage
          </router-link>
        </div>

        <!-- Poll Cards -->
        <div v-else class="space-y-4">
          <router-link
            v-for="poll in polls"
            :key="poll.id"
            :to="`/poll/${poll.id}`"
            class="block group"
          >
            <div class="card bg-base-100 shadow-lg hover:shadow-xl transition-all duration-300 border border-base-200 hover:border-primary/30 overflow-hidden">
              <div class="card-body p-0">
                <div class="p-6">
                  <div class="flex items-start justify-between gap-4">
                    <div class="flex-1 min-w-0">
                      <h3 class="card-title text-lg group-hover:text-primary transition-colors truncate">
                        {{ poll.title }}
                      </h3>
                      <p v-if="poll.description" class="text-base-content/60 text-sm mt-1 line-clamp-2">
                        {{ poll.description }}
                      </p>
                      <div class="flex flex-wrap gap-3 mt-3">
                        <span v-if="poll.location" class="inline-flex items-center gap-1 text-xs text-base-content/60">
                          <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                          </svg>
                          {{ poll.location }}
                        </span>
                        <span class="inline-flex items-center gap-1 text-xs text-base-content/60">
                          <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                          </svg>
                          {{ poll.creator?.name || 'Anonyme' }}
                        </span>
                        <span v-if="poll.final_date" class="badge badge-success badge-sm gap-1">
                          <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7" />
                          </svg>
                          Date fixée
                        </span>
                      </div>
                    </div>
                    <div class="text-center min-w-16">
                      <div class="text-2xl font-bold text-primary">{{ poll.participant_count || 0 }}</div>
                      <div class="text-[10px] uppercase text-base-content/40 font-bold">participants</div>
                    </div>
                  </div>
                </div>
                <!-- Footer with date -->
                <div class="bg-base-50 px-6 py-3 border-t border-base-200">
                  <div class="flex justify-between items-center text-xs text-base-content/50">
                    <span>Créé {{ formatRelativeTime(poll.created_at) }}</span>
                    <svg class="w-4 h-4 group-hover:translate-x-1 transition-transform" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                    </svg>
                  </div>
                </div>
              </div>
            </div>
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { usePollsStore } from '@/stores/polls'

const authStore = useAuthStore()
const pollsStore = usePollsStore()

const polls = ref([])
const loading = ref(false)
const searchQuery = ref('')
let debounceTimer = null

async function loadPolls() {
  loading.value = true
  try {
    const data = await pollsStore.fetchPolls(searchQuery.value ? { search: searchQuery.value } : {})
    polls.value = data.polls || []
  } catch (error) {
    console.error('Error loading polls:', error)
  } finally {
    loading.value = false
  }
}

function debouncedSearch() {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    loadPolls()
  }, 300)
}

function formatRelativeTime(dateStr) {
  const date = new Date(dateStr)
  const now = new Date()
  const seconds = Math.floor((now - date) / 1000)

  if (seconds < 60) return 'à l\'instant'
  if (seconds < 3600) return `il y a ${Math.floor(seconds / 60)} min`
  if (seconds < 86400) return `il y a ${Math.floor(seconds / 3600)} h`
  if (seconds < 604800) return `il y a ${Math.floor(seconds / 86400)} j`
  return date.toLocaleDateString('fr-FR')
}

onMounted(() => {
  loadPolls()
})
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
