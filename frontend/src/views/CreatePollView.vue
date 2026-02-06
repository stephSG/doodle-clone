<template>
  <div class="max-w-2xl mx-auto px-4 py-8">
    <!-- Header -->
    <div class="flex items-center gap-4 mb-8">
      <button @click="$router.back()" class="btn btn-ghost btn-circle">
        <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
      </button>
      <div>
        <h1 class="text-2xl font-bold">{{ isEdit ? 'Modifier' : 'Créer' }} un sondage</h1>
        <p class="text-sm text-base-content/60">Proposez des dates et laissez les participants voter</p>
      </div>
    </div>

    <form @submit.prevent="handleSubmit" class="space-y-6">
      <!-- Basic Info -->
      <div class="card bg-base-100 shadow-lg">
        <div class="card-body">
          <h2 class="card-title text-lg">Informations de base</h2>

          <div class="form-control">
            <label class="label">
              <span class="label-text font-medium">Titre *</span>
            </label>
            <input
              v-model="form.title"
              type="text"
              placeholder="Ex: Réunion d'équipe, Dîner entre amis..."
              class="input input-bordered"
              required
            />
          </div>

          <div class="form-control">
            <label class="label">
              <span class="label-text font-medium">Description</span>
            </label>
            <textarea
              v-model="form.description"
              class="textarea textarea-bordered"
              placeholder="Ajoutez plus de détails..."
              rows="2"
            ></textarea>
          </div>

          <div class="form-control">
            <label class="label">
              <span class="label-text font-medium">Lieu</span>
            </label>
            <div class="relative">
              <svg class="w-5 h-5 absolute left-3 top-1/2 -translate-y-1/2 text-base-content/40" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
              </svg>
              <input
                v-model="form.location"
                type="text"
                placeholder="Ex: Paris, Zoom, En présentiel..."
                class="input input-bordered pl-10"
              />
            </div>
          </div>
        </div>
      </div>

      <!-- Options -->
      <div class="card bg-base-100 shadow-lg">
        <div class="card-body">
          <h2 class="card-title text-lg">Options</h2>

          <div class="space-y-3">
            <label class="label cursor-pointer justify-start gap-3 p-3 bg-base-50 rounded-lg">
              <input type="checkbox" v-model="form.allow_multiple" class="checkbox checkbox-primary" />
              <div>
                <span class="label-text font-medium">Votes multiples</span>
                <span class="label-text-alt">Les participants peuvent voter pour plusieurs dates</span>
              </div>
            </label>

            <label class="label cursor-pointer justify-start gap-3 p-3 bg-base-50 rounded-lg">
              <input type="checkbox" v-model="form.allow_maybe" class="checkbox checkbox-primary" />
              <div>
                <span class="label-text font-medium">Option "Peut-être"</span>
                <span class="label-text-alt">Permet les réponses incertaines</span>
              </div>
            </label>

            <label class="label cursor-pointer justify-start gap-3 p-3 bg-base-50 rounded-lg">
              <input type="checkbox" v-model="form.anonymous" class="checkbox checkbox-primary" />
              <div>
                <span class="label-text font-medium">Votes anonymes</span>
                <span class="label-text-alt">Pas besoin de compte pour voter</span>
              </div>
            </label>
          </div>
        </div>
      </div>

      <!-- Dates -->
      <div class="card bg-base-100 shadow-lg">
        <div class="card-body">
          <h2 class="card-title text-lg">Dates proposées</h2>

          <div class="space-y-4">
            <div
              v-for="(date, index) in form.dates"
              :key="index"
              class="p-4 bg-base-50 rounded-lg space-y-4"
            >
              <div class="flex justify-between items-start">
                <span class="text-sm font-medium">Date #{{ index + 1 }}</span>
                <button
                  v-if="form.dates.length > 1"
                  type="button"
                  @click="removeDate(index)"
                  class="btn btn-sm btn-ghost btn-circle text-error"
                >
                  <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>

              <!-- Date & Time Picker combined -->
              <div class="form-control">
                <label class="label">
                  <span class="label-text text-xs">Date et heure *</span>
                </label>
                <VueDatePicker
                  v-model="form.dates[index].datetime"
                  :format="formatDateTime"
                  :min-date="new Date()"
                  :teleport="true"
                  :enable-time-picker="true"
                  :is-24="true"
                  :minutes-increment="15"
                  :auto-apply="true"
                  :month-change-on-scroll="false"
                  placeholder="Sélectionnez une date et heure"
                  class="dp-input-bordered"
                />
              </div>
            </div>
          </div>

          <button
            type="button"
            @click="addDate"
            class="btn btn-outline w-full gap-2 mt-4"
          >
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            Ajouter une date
          </button>

          <!-- Expiration -->
          <div class="form-control mt-6">
            <label class="label">
              <span class="label-text font-medium">Date limite (optionnel)</span>
            </label>
            <VueDatePicker
              v-model="expiresAt"
              :format="formatDateTime"
              :min-date="new Date()"
              :teleport="true"
              :enable-time-picker="true"
              :is-24="true"
              :minutes-increment="15"
              :auto-apply="true"
              :month-change-on-scroll="false"
              placeholder="Sélectionnez la date limite"
              class="dp-input-bordered"
            />
            <span class="label-text-alt">Après cette date, les votes ne seront plus acceptés</span>
          </div>
        </div>
      </div>

      <!-- Submit -->
      <div class="flex gap-3">
        <button
          type="button"
          @click="$router.back()"
          class="btn btn-ghost flex-1"
        >
          Annuler
        </button>
        <button
          type="submit"
          class="btn btn-primary flex-1"
          :disabled="loading || !isFormValid"
        >
          <span v-if="loading" class="loading loading-spinner"></span>
          {{ isEdit ? 'Modifier' : 'Créer' }} le sondage
        </button>
      </div>
    </form>
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

function formatDateTime(date) {
  if (!date) return ''
  return new Intl.DateTimeFormat('fr-FR', {
    weekday: 'short',
    day: 'numeric',
    month: 'short',
    year: 'numeric',
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
      router.push('/dashboard')
    }
  }
})
</script>

<style>
.dp__main {
  border: none;
}

.dp__input {
  width: 100%;
  padding: 0.75rem 1rem;
  border-radius: 0.5rem;
  border: 1px solid var(--fallback-bc, oklch(var(--bc) / 0.2));
  background-color: var(--fallback-b1, oklch(var(--b1) / 1));
  color: var(--fallback-bc, oklch(var(--bc) / 1));
  font-size: 0.875rem;
  transition: border-color 0.2s;
}

.dp__input:focus {
  outline: none;
  border-color: var(--fallback-p, oklch(var(--p) / 1));
  box-shadow: 0 0 0 2px var(--fallback-p, oklch(var(--p) / 0.2));
}

.dp__input::placeholder {
  color: var(--fallback-bc, oklch(var(--bc) / 0.5));
}

.dp__menu {
  border: 1px solid var(--fallback-bc, oklch(var(--bc) / 0.2));
  border-radius: 0.75rem;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.15);
  font-family: inherit;
}

.dp__active_date {
  background-color: var(--fallback-p, oklch(var(--p) / 1)) !important;
  color: white !important;
}

.dp__today {
  border-color: var(--fallback-p, oklch(var(--p) / 1)) !important;
}

.dp__range_end, .dp__range_start, .dp__range_between {
  background-color: var(--fallback-p, oklch(var(--p) / 0.9)) !important;
}

.dp__action_select {
  background-color: var(--fallback-p, oklch(var(--p) / 1)) !important;
  color: white !important;
}

.dp__action_select:hover {
  background-color: var(--fallback-p, oklch(var(--p) / 0.9)) !important;
}

.dp__overlay_cell_active {
  background-color: var(--fallback-p, oklch(var(--p) / 1)) !important;
}

.dp__month_year_select {
  font-weight: 600;
}

.dp__calendar_header_item {
  font-weight: 500;
  color: var(--fallback-bc, oklch(var(--bc) / 0.7));
}

.dp__time_display {
  font-size: 1.1rem;
}

/* Mobile optimizations */
@media (max-width: 640px) {
  .dp__menu {
    width: calc(100vw - 2rem) !important;
    max-width: none !important;
  }
}
</style>
