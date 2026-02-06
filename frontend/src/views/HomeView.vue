<template>
  <div class="min-h-screen max-w-md mx-auto shadow-[0_0_80px_rgba(0,0,0,0.03)] overflow-x-hidden relative transition-colors duration-300" :class="isDark ? 'bg-slate-900' : 'bg-white'">
    <!-- Loading State -->
    <div v-if="loading" class="flex flex-col items-center justify-center min-h-screen">
      <div class="relative w-16 h-16">
        <div class="absolute inset-0 border-[3px] border-indigo-50 rounded-full"></div>
        <div class="absolute inset-0 border-[3px] border-indigo-600 rounded-full border-t-transparent animate-spin"></div>
      </div>
      <p class="mt-6 text-indigo-600 font-bold tracking-wide animate-pulse">Chargement...</p>
    </div>

    <!-- Home View -->
    <div v-else class="p-6 pb-32 animate-in">
      <!-- Header -->
      <header class="flex justify-between items-center mb-10 mt-4">
        <div>
          <h1 class="text-4xl font-extrabold tracking-tight" :class="isDark ? 'text-white' : 'text-slate-900'">Doodle</h1>
          <p class="font-medium mt-1" :class="isDark ? 'text-slate-400' : 'text-slate-400'">Vos événements, simplifiés.</p>
        </div>
        <router-link
          to="/create"
          class="w-14 h-14 bg-gradient-to-br from-indigo-500 to-violet-600 rounded-2xl flex items-center justify-center text-white shadow-xl shadow-indigo-100 transition-transform active:scale-90"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="m12 3-1.912 5.813a2 2 0 0 1-1.275 1.275L3 12l5.813 1.912a2 2 0 0 1 1.275 1.275L12 21l1.912-5.813a2 2 0 0 1 1.275-1.275L21 12l-5.813-1.912a2 2 0 0 1-1.275-1.275L12 3Z"/>
            <path d="M5 3v4"/>
            <path d="M19 17v4"/>
            <path d="M3 5h4"/>
            <path d="M17 19h4"/>
          </svg>
        </router-link>
      </header>

      <!-- Search -->
      <div class="relative mb-6">
        <svg class="w-5 h-5 absolute left-5 top-1/2 -translate-y-1/2" fill="none" viewBox="0 0 24 24" stroke="currentColor" :class="isDark ? 'text-slate-500' : 'text-slate-300'">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
        <input
          type="text"
          v-model="searchQuery"
          @input="debouncedSearch"
          placeholder="Rechercher un sondage..."
          class="w-full h-14 pl-14 pr-6 rounded-2xl outline-none transition-all font-medium border-2 border-transparent"
          :class="[
            isDark
              ? 'bg-slate-800 text-slate-200 placeholder:text-slate-500 focus:bg-slate-700 focus:border-indigo-500'
              : 'bg-slate-50 text-slate-600 placeholder:text-slate-300 focus:bg-white focus:border-indigo-100'
          ]"
        />
      </div>

      <!-- Polls List -->
      <div class="space-y-4">
        <div
          v-for="poll in polls"
          :key="poll.id"
          @click="goToPoll(poll.id)"
          class="group p-5 rounded-[32px] border transition-all cursor-pointer relative overflow-hidden flex items-center gap-5 active:scale-[0.98]"
          :class="[
            isDark
              ? 'bg-slate-800 border-slate-700 hover:border-indigo-500'
              : 'bg-white border-slate-100 hover:border-indigo-100 hover:shadow-md'
          ]"
        >
          <!-- Date Badge -->
          <div class="w-16 h-16 rounded-2xl bg-indigo-50 text-indigo-600 flex flex-col items-center justify-center font-bold shrink-0 transition-transform group-hover:rotate-2">
            <span class="text-[10px] uppercase tracking-widest opacity-70 mb-0.5">{{ formatMonth(poll.date_options?.[0]?.start_time) }}</span>
            <span class="text-xl">{{ formatDay(poll.date_options?.[0]?.start_time) }}</span>
          </div>

          <!-- Content -->
          <div class="flex-1 min-w-0">
            <h3 class="font-bold text-lg leading-tight truncate group-hover:text-indigo-600 transition-colors" :class="isDark ? 'text-white' : 'text-slate-800'">
              {{ poll.title }}
            </h3>
            <p v-if="poll.description" class="text-sm mt-1 line-clamp-2 truncate" :class="isDark ? 'text-slate-400' : 'text-slate-400'">
              {{ poll.description }}
            </p>
            <div class="flex items-center gap-4 text-xs mt-2 font-semibold" :class="isDark ? 'text-slate-500' : 'text-slate-400'">
              <span class="flex items-center gap-1.5">
                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" class="text-indigo-400">
                  <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/>
                  <circle cx="9" cy="7" r="4"/>
                  <path d="M22 21v-2a4 4 0 0 0-3-3.87"/>
                  <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
                </svg>
                {{ poll.participant_count || 0 }}
              </span>
              <span v-if="poll.location" class="truncate flex items-center gap-1.5">
                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" class="text-indigo-400">
                  <path d="M20 10c0 6-8 12-8 12s-8-6-8-12a8 8 0 0 1 16 0Z"/>
                  <circle cx="12" cy="10" r="3"/>
                </svg>
                {{ poll.location }}
              </span>
            </div>
          </div>

          <!-- Arrow -->
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" class="text-slate-200 group-hover:text-indigo-300 transition-all shrink-0">
            <path d="m9 18 6-6-6-6"/>
          </svg>

          <!-- Final Date Badge -->
          <div v-if="poll.final_date" class="absolute top-0 right-0 bg-emerald-400 text-white text-[8px] font-black px-3 py-1 rounded-bl-xl uppercase tracking-widest shadow-sm flex items-center gap-1">
            <svg xmlns="http://www.w3.org/2000/svg" width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round">
              <path d="M5 13l4 4L19 7"/>
            </svg>
            Fixé
          </div>
        </div>

        <!-- Empty State -->
        <div v-if="polls.length === 0" class="text-center py-20">
          <div class="inline-block p-6 bg-slate-50 rounded-full mb-6">
            <svg class="w-16 h-16 text-slate-200" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" />
            </svg>
          </div>
          <h3 class="text-xl font-bold mb-2 text-slate-800">{{ searchQuery ? 'Aucun résultat' : 'Aucun sondage' }}</h3>
          <p class="text-slate-400 mb-6 text-sm">{{ searchQuery ? 'Essayez d\'autres termes de recherche' : 'Créez votre premier sondage pour commencer' }}</p>
          <router-link to="/create" class="inline-flex items-center gap-2 bg-indigo-600 text-white h-14 px-8 rounded-2xl font-bold shadow-xl shadow-indigo-100 active:scale-95 transition-all">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
              <path d="M12 5v14"/>
              <path d="M5 12h14"/>
            </svg>
            Créer un sondage
          </router-link>
        </div>
      </div>

      <!-- Create Button -->
      <div class="fixed bottom-8 left-6 right-6 z-50 max-w-md mx-auto">
        <router-link
          to="/create"
          class="w-full bg-slate-900 text-white h-16 rounded-2xl font-bold shadow-2xl flex items-center justify-center gap-3 active:scale-95 transition-all group relative overflow-hidden block"
        >
          <div class="absolute inset-0 bg-indigo-600 translate-y-full group-hover:translate-y-0 transition-transform duration-300"></div>
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" class="relative z-10">
            <path d="M12 5v14"/>
            <path d="M5 12h14"/>
          </svg>
          <span class="relative z-10">Créer un sondage</span>
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { usePollsStore } from '@/stores/polls'

const router = useRouter()
const pollsStore = usePollsStore()

// Dark mode
const isDark = computed(() => document.documentElement.getAttribute('data-theme') === 'dark')

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

function goToPoll(id) {
  router.push(`/poll/${id}`)
}

function formatMonth(dateStr) {
  if (!dateStr) return '--'
  const date = new Date(dateStr)
  return date.toLocaleDateString('fr-FR', { month: 'short' }).replace('.', '')
}

function formatDay(dateStr) {
  if (!dateStr) return '--'
  const date = new Date(dateStr)
  return date.getDate()
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

/* Hide scrollbar for mobile */
:deep(*) {
  -webkit-tap-highlight-color: transparent;
}
</style>
