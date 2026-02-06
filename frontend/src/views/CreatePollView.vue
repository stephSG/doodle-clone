<template>
  <div class="min-h-screen max-w-md mx-auto shadow-[0_0_80px_rgba(0,0,0,0.03)] overflow-x-hidden relative transition-colors duration-300" :class="isDark ? 'bg-slate-900' : 'bg-white'">
    <div class="p-6 pb-32 animate-in">
      <!-- Header -->
      <header class="flex justify-between items-center mb-10 mt-4">
        <button @click="$router.back()" class="w-12 h-12 flex items-center justify-center rounded-2xl active:scale-90 transition-all border" :class="isDark ? 'bg-slate-800 text-slate-300 border-slate-700' : 'bg-slate-50 text-slate-600 border-slate-100'">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="m15 18-6-6 6-6"/>
          </svg>
        </button>
        <h2 class="text-xl font-extrabold" :class="isDark ? 'text-white' : 'text-slate-900'">{{ isEdit ? 'Modifier' : 'Nouvel Event' }}</h2>
        <div class="w-12 h-12"></div>
      </header>

      <form @submit.prevent="handleSubmit" class="space-y-6">
        <!-- Title -->
        <div class="space-y-2">
          <label class="text-[11px] font-black text-indigo-500 uppercase tracking-widest ml-1">Titre de l'événement</label>
          <input
            v-model="form.title"
            type="text"
            placeholder="Ex: Dîner de fin d'année"
            class="w-full h-16 px-6 rounded-2xl outline-none transition-all font-bold border-2 border-transparent"
            :class="[
              isDark
                ? 'bg-slate-800 text-white placeholder:text-slate-500 focus:bg-slate-700 focus:border-indigo-500'
                : 'bg-slate-50 text-slate-700 placeholder:text-slate-300 focus:bg-white focus:border-indigo-100'
            ]"
            required
          />
        </div>

        <!-- Location -->
        <div class="space-y-2">
          <label class="text-[11px] font-black text-indigo-500 uppercase tracking-widest ml-1">Lieu</label>
          <div class="relative">
            <svg class="w-5 h-5 absolute left-5 top-1/2 -translate-y-1/2" fill="none" viewBox="0 0 24 24" stroke="currentColor" :class="isDark ? 'text-slate-500' : 'text-slate-300'">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
            </svg>
            <input
              v-model="form.location"
              type="text"
              placeholder="Ex: Paris ou Discord"
              class="w-full h-16 pl-14 pr-6 rounded-2xl outline-none transition-all font-bold border-2 border-transparent"
              :class="[
                isDark
                  ? 'bg-slate-800 text-white placeholder:text-slate-500 focus:bg-slate-700 focus:border-indigo-500'
                  : 'bg-slate-50 text-slate-700 placeholder:text-slate-300 focus:bg-white focus:border-indigo-100'
              ]"
            />
          </div>
        </div>

        <!-- Description -->
        <div class="space-y-2">
          <label class="text-[11px] font-black text-indigo-500 uppercase tracking-widest ml-1">Description</label>
          <textarea
            v-model="form.description"
            class="w-full min-h-24 px-6 py-4 rounded-2xl outline-none transition-all font-medium border-2 border-transparent resize-none"
            :class="[
              isDark
                ? 'bg-slate-800 text-white placeholder:text-slate-500 focus:bg-slate-700 focus:border-indigo-500'
                : 'bg-slate-50 text-slate-700 placeholder:text-slate-300 focus:bg-white focus:border-indigo-100'
            ]"
            placeholder="Ajoutez plus de détails..."
          ></textarea>
        </div>

        <!-- Options -->
        <div class="space-y-3">
          <label class="text-[11px] font-black text-indigo-500 uppercase tracking-widest ml-1">Options</label>
          <div class="space-y-3">
            <label class="flex items-center gap-3 p-4 bg-slate-50 rounded-2xl cursor-pointer hover:bg-indigo-50 transition-colors">
              <input type="checkbox" v-model="form.allow_multiple" class="w-5 h-5 rounded border-2 border-slate-200 text-indigo-600 focus:ring-indigo-500 focus:ring-offset-0" />
              <div class="flex-1">
                <span class="block font-bold text-slate-800 text-sm">Votes multiples</span>
                <span class="block text-xs text-slate-400">Les participants peuvent voter pour plusieurs dates</span>
              </div>
            </label>

            <label class="flex items-center gap-3 p-4 bg-slate-50 rounded-2xl cursor-pointer hover:bg-indigo-50 transition-colors">
              <input type="checkbox" v-model="form.allow_maybe" class="w-5 h-5 rounded border-2 border-slate-200 text-indigo-600 focus:ring-indigo-500 focus:ring-offset-0" />
              <div class="flex-1">
                <span class="block font-bold text-slate-800 text-sm">Option "Peut-être"</span>
                <span class="block text-xs text-slate-400">Permet les réponses incertaines</span>
              </div>
            </label>

            <label class="flex items-center gap-3 p-4 bg-slate-50 rounded-2xl cursor-pointer hover:bg-indigo-50 transition-colors">
              <input type="checkbox" v-model="form.anonymous" class="w-5 h-5 rounded border-2 border-slate-200 text-indigo-600 focus:ring-indigo-500 focus:ring-offset-0" />
              <div class="flex-1">
                <span class="block font-bold text-slate-800 text-sm">Votes anonymes</span>
                <span class="block text-xs text-slate-400">Pas besoin de compte pour voter</span>
              </div>
            </label>
          </div>
        </div>

        <!-- Dates & Times -->
        <div class="space-y-3">
          <label class="text-[11px] font-black text-indigo-500 uppercase tracking-widest ml-1">Dates & Heures</label>
          <div
            v-for="(date, index) in form.dates"
            :key="index"
            class="flex gap-2 group animate-in"
          >
            <VueDatePicker
              v-model="date.datetime"
              :format="formatDateTimeDisplay"
              :min-date="new Date()"
              :teleport="true"
              :enable-time-picker="true"
              :is-24="true"
              :minutes-increment="15"
              :auto-apply="true"
              :month-change-on-scroll="false"
              :placeholder="index === 0 ? 'Sélectionnez date et heure' : ''"
              class="dp-input-bordered flex-1"
            >
              <template #input-icon>
                <svg class="w-5 h-5 text-slate-300 ml-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
              </template>
            </VueDatePicker>
            <button
              v-if="form.dates.length > 1"
              type="button"
              @click="removeDate(index)"
              class="w-16 h-16 flex items-center justify-center bg-rose-50 text-rose-500 rounded-2xl active:scale-90 transition-all shrink-0"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <path d="M18 6 6 18"/>
                <path d="m6 6 12 12"/>
              </svg>
            </button>
          </div>
          <button
            type="button"
            @click="addDate"
            class="w-full h-14 border-2 border-dashed border-slate-200 text-slate-400 rounded-2xl font-bold flex items-center justify-center gap-2 mt-2 hover:border-indigo-300 hover:text-indigo-500 transition-all"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
              <path d="M12 5v14"/>
              <path d="M5 12h14"/>
            </svg>
            Ajouter une date
          </button>

          <!-- Expiration -->
          <div class="mt-4 space-y-2">
            <label class="text-[11px] font-black text-slate-400 uppercase tracking-widest ml-1">Date limite (optionnel)</label>
            <VueDatePicker
              v-model="expiresAt"
              :format="formatDateTimeDisplay"
              :min-date="new Date()"
              :teleport="true"
              :enable-time-picker="true"
              :is-24="true"
              :minutes-increment="15"
              :auto-apply="true"
              :month-change-on-scroll="false"
              placeholder="Sélectionnez la date limite"
              class="dp-input-bordered w-full"
            >
              <template #input-icon>
                <svg class="w-5 h-5 text-slate-300 ml-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </template>
            </VueDatePicker>
          </div>
        </div>
      </form>

      <!-- Submit Button -->
      <div class="fixed bottom-8 left-6 right-6 z-50 max-w-md mx-auto">
        <button
          @click="handleSubmit"
          :disabled="loading || !isFormValid"
          class="w-full bg-indigo-600 text-white h-16 rounded-2xl font-bold shadow-xl shadow-indigo-100 disabled:bg-slate-200 disabled:text-slate-400 disabled:shadow-none transition-all flex items-center justify-center gap-3 active:scale-95"
        >
          <span v-if="loading" class="loading loading-spinner loading-sm"></span>
          <span v-else>{{ isEdit ? 'Modifier' : 'Publier le' }} sondage</span>
          <svg v-if="!loading" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
            <path d="m9 18 6-6-6-6"/>
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { usePollsStore } from '@/stores/polls'
import { useUiStore } from '@/stores/ui'
import { VueDatePicker } from '@vuepic/vue-datepicker'
import '@vuepic/vue-datepicker/dist/main.css'

const router = useRouter()
const route = useRoute()
const pollsStore = usePollsStore()
const uiStore = useUiStore()

// Dark mode
const isDark = computed(() => document.documentElement.getAttribute('data-theme') === 'dark')

const isEdit = ref(false)
const loading = ref(false)

const form = ref({
  title: '',
  description: '',
  location: '',
  allow_multiple: false,
  allow_maybe: true,
  anonymous: false,
  dates: [
    { datetime: null }
  ]
})

const expiresAt = ref(null)

const isFormValid = computed(() => {
  return form.value.title.trim() &&
         form.value.dates.some(d => d.datetime)
})

function formatDateTimeDisplay(date) {
  if (!date) return ''
  return new Intl.DateTimeFormat('fr-FR', {
    weekday: 'short',
    day: 'numeric',
    month: 'short',
    hour: '2-digit',
    minute: '2-digit'
  }).format(date)
}

function addDate() {
  form.value.dates.push({ datetime: null })
}

function removeDate(index) {
  if (form.value.dates.length > 1) {
    form.value.dates.splice(index, 1)
  }
}

async function handleSubmit() {
  if (!isFormValid.value) {
    uiStore.error('Veuillez remplir le titre et au moins une date')
    return
  }

  loading.value = true

  // Convert datetime objects to ISO strings
  const datesWithDateTime = form.value.dates
    .filter(d => d.datetime)
    .map(d => {
      return {
        start_time: d.datetime.toISOString(),
        end_time: null
      }
    })

  // Handle expiration
  const expiresAtValue = expiresAt.value ? expiresAt.value.toISOString() : null

  const payload = {
    title: form.value.title,
    description: form.value.description || null,
    location: form.value.location || null,
    allow_multiple: form.value.allow_multiple,
    allow_maybe: form.value.allow_maybe,
    anonymous: form.value.anonymous,
    dates: datesWithDateTime,
    expires_at: expiresAtValue
  }

  try {
    if (isEdit.value) {
      await pollsStore.updatePoll(route.params.id, payload)
      uiStore.success('Sondage modifié !')
      router.push(`/poll/${route.params.id}`)
    } else {
      const result = await pollsStore.createPoll(payload)
      uiStore.success('Sondage créé avec succès !')
      router.push(`/poll/${result.id || result.poll?.id}`)
    }
  } catch (error) {
    console.error('Error creating poll:', error)
    uiStore.error(error.response?.data?.error || 'Erreur lors de la création du sondage')
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  if (route.params.id) {
    isEdit.value = true
    try {
      const data = await pollsStore.fetchPoll(route.params.id)
      const poll = data.poll

      // Format dates for the form
      const dates = data.date_options?.map(d => {
        return {
          datetime: new Date(d.start_time)
        }
      }) || [{ datetime: null }]

      // Format expiration
      if (poll.expires_at) {
        expiresAt.value = new Date(poll.expires_at)
      }

      form.value = {
        title: poll.title,
        description: poll.description || '',
        location: poll.location || '',
        allow_multiple: poll.allow_multiple,
        allow_maybe: poll.allow_maybe,
        anonymous: poll.anonymous,
        dates: dates.length > 0 ? dates : [{ datetime: null }]
      }
    } catch (error) {
      uiStore.error('Erreur lors du chargement du sondage')
      router.push('/')
    }
  }
})
</script>

<style scoped>
.dp__main {
  border: none;
}

:deep(.dp__input) {
  width: 100%;
  height: 4rem;
  padding: 0 1.5rem 0 3rem;
  border-radius: 1rem;
  border: 2px solid transparent;
  background-color: #f8fafc;
  color: #334155;
  font-size: 1rem;
  font-weight: 600;
  font-family: 'Plus Jakarta Sans', sans-serif;
  transition: all 0.2s;
}

:deep(.dp__input:focus) {
  outline: none;
  border-color: #e0e7ff;
  background-color: white;
}

:deep(.dp__input::placeholder) {
  color: #cbd5e1;
  font-weight: 500;
}

:deep(.dp__menu) {
  border: 1px solid #e2e8f0;
  border-radius: 1rem;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.15);
  font-family: 'Plus Jakarta Sans', sans-serif;
}

:deep(.dp__active_date) {
  background-color: #4f46e5 !important;
  color: white !important;
}

:deep(.dp__today) {
  border-color: #4f46e5 !important;
}

:deep(.dp__range_end), :deep(.dp__range_start), :deep(.dp__range_between) {
  background-color: #4f46e5 !important;
}

:deep(.dp__action_select) {
  background-color: #4f46e5 !important;
  color: white !important;
}

:deep(.dp__action_select:hover) {
  background-color: #4338ca !important;
}

:deep(.dp__overlay_cell_active) {
  background-color: #4f46e5 !important;
}

:deep(.dp__month_year_select) {
  font-weight: 600;
}

:deep(.dp__calendar_header_item) {
  font-weight: 500;
  color: #94a3b8;
}

:deep(.dp__time_display) {
  font-size: 1.1rem;
}

:deep(.dp__input_wrap) {
  position: relative;
}

:deep(.dp__input_icon) {
  left: 0;
}

/* Mobile optimizations */
@media (max-width: 640px) {
  :deep(.dp__menu) {
    width: calc(100vw - 2rem) !important;
    max-width: none !important;
  }
}

:deep(*) {
  -webkit-tap-highlight-color: transparent;
}
</style>
