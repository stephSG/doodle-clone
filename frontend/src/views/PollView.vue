<template>
  <div v-if="loading" class="flex justify-center items-center min-h-screen">
    <span class="loading loading-spinner loading-lg"></span>
  </div>

  <div v-else-if="poll" class="max-w-5xl mx-auto px-4 py-8 space-y-8">
    <!-- Poll Header -->
    <div class="card bg-gradient-to-r from-primary/10 to-primary/5 shadow-xl">
      <div class="card-body">
        <div class="flex flex-col md:flex-row justify-between items-start gap-4">
          <div class="flex-1">
            <div class="flex items-center gap-3 mb-3">
              <h1 class="text-3xl font-bold">{{ poll.title }}</h1>
              <div v-if="poll.final_date" class="badge badge-success gap-1">
                <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7" />
                </svg>
                Final Date Set
              </div>
            </div>
            <p v-if="poll.description" class="text-base-content/70 mb-4">{{ poll.description }}</p>
            <div class="flex flex-wrap gap-4 text-sm">
              <div v-if="poll.location" class="flex items-center gap-2 btn btn-ghost btn-sm">
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                </svg>
                {{ poll.location }}
              </div>
              <div class="flex items-center gap-2 btn btn-ghost btn-sm">
                <div class="avatar placeholder">
                  <div class="bg-primary text-primary-content rounded-full w-6">
                    <span class="text-xs">{{ poll.creator?.name?.charAt(0) || '?' }}</span>
                  </div>
                </div>
                {{ poll.creator?.name }}
              </div>
              <div v-if="poll.expires_at" class="flex items-center gap-2 btn btn-ghost btn-sm">
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                {{ formatDate(poll.expires_at) }}
              </div>
            </div>
          </div>

          <div class="dropdown dropdown-end">
            <div tabindex="0" role="button" class="btn btn-ghost btn-circle">
              <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z" />
              </svg>
            </div>
            <ul tabindex="0" class="dropdown-content z-[1] menu p-2 shadow bg-base-100 rounded-box w-52">
              <li><a @click="copyLink" class="flex gap-2">
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                </svg>
                Copy Link
              </a></li>
              <li v-if="isCreator"><router-link :to="`/poll/${poll.id}/edit`" class="flex gap-2">
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                </svg>
                Edit Poll
              </router-link></li>
              <li><a @click="exportPDF" class="flex gap-2">
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
                Export PDF
              </a></li>
              <li><a @click="exportICS" class="flex gap-2">
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
                Add to Calendar
              </a></li>
            </ul>
          </div>
        </div>
      </div>
    </div>

    <!-- Voting Section -->
    <div class="card bg-base-100 shadow-xl">
      <div class="card-body">
        <h2 class="card-title flex items-center gap-2">
          <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
          </svg>
          Choose Your Availability
        </h2>

        <!-- User's Voting Row -->
        <div v-if="authStore.isAuthenticated" class="bg-primary/5 rounded-lg p-4 mb-6">
          <div class="flex items-center gap-2 mb-4">
            <div class="avatar placeholder">
              <div class="bg-primary text-primary-content rounded-full w-8">
                <span>{{ authStore.user?.name?.charAt(0) || 'U' }}</span>
              </div>
            </div>
            <span class="font-medium">Your votes</span>
          </div>
          <div class="space-y-3">
            <div v-for="option in dateOptions" :key="option.id" class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 p-3 bg-base-100 rounded-lg">
              <div class="flex items-center gap-3">
                <div class="text-center">
                  <div class="text-2xl font-bold">{{ formatDateShort(option.start_time) }}</div>
                  <div class="text-xs text-base-content/50">{{ formatTime(option.start_time) }}</div>
                </div>
                <div v-if="poll.final_date === option.id" class="badge badge-success">
                  Final
                </div>
              </div>
              <div class="join w-full sm:w-auto">
                <button
                  @click="vote(option.id, 'yes')"
                  class="btn btn-sm join-item flex-1 sm:flex-none"
                  :class="getUserVoteClass(option.id, 'yes')"
                >
                  <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                  </svg>
                  OK
                </button>
                <button
                  v-if="poll.allow_maybe"
                  @click="vote(option.id, 'maybe')"
                  class="btn btn-sm join-item flex-1 sm:flex-none"
                  :class="getUserVoteClass(option.id, 'maybe')"
                >
                  <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  Maybe
                </button>
                <button
                  @click="vote(option.id, 'no')"
                  class="btn btn-sm join-item flex-1 sm:flex-none"
                  :class="getUserVoteClass(option.id, 'no')"
                >
                  <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                  No
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Login Prompt -->
        <div v-else class="alert alert-info">
          <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
          </svg>
          <div>
            <h3 class="font-bold">Login to vote!</h3>
            <div class="text-xs">Sign in with Google or email to participate in this poll.</div>
          </div>
          <router-link to="/login" class="btn btn-sm btn-primary">Login</router-link>
        </div>

        <!-- Results Summary -->
        <div class="mt-8">
          <h3 class="font-semibold mb-4 flex items-center gap-2">
            <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
            Results
          </h3>
          <div class="overflow-x-auto">
            <table class="table table-zebra">
              <thead>
                <tr>
                  <th>Participant</th>
                  <th v-for="option in dateOptions" :key="option.id" class="text-center">
                    <div class="text-xs">{{ formatDateShort(option.start_time) }}</div>
                    <div class="text-xs text-base-content/50">{{ formatTime(option.start_time) }}</div>
                  </th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(votes, userId) in votesByUser" :key="userId">
                  <td class="font-medium">{{ getUserDisplayName(votes[0]) }}</td>
                  <td v-for="option in dateOptions" :key="option.id" class="text-center">
                    <span
                      v-if="getVoteResponse(option.id, votes)"
                      :class="['vote-badge', getVoteResponse(option.id, votes)]"
                    >
                      {{ getVoteIcon(getVoteResponse(option.id, votes)) }}
                    </span>
                    <span v-else class="text-base-content/20">—</span>
                  </td>
                </tr>
                <tr v-for="(vote, idx) in anonymousVotes" :key="'anon-' + idx">
                  <td class="font-medium text-base-content/70">{{ vote.user_name }}</td>
                  <td v-for="option in dateOptions" :key="option.id" class="text-center">
                    <span
                      v-if="vote.date_option_id === option.id"
                      :class="['vote-badge', vote.response]"
                    >
                      {{ getVoteIcon(vote.response) }}
                    </span>
                    <span v-else class="text-base-content/20">—</span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Summary Cards -->
        <div class="mt-8">
          <h3 class="font-semibold mb-4">Availability Summary</h3>
          <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            <div
              v-for="option in dateOptions"
              :key="option.id"
              class="card bg-base-200 cursor-pointer hover:bg-base-300 transition-colors"
              :class="{ 'ring-2 ring-success': poll.final_date === option.id }"
            >
              <div class="card-body p-4">
                <div class="flex justify-between items-start mb-2">
                  <div class="text-sm font-medium">{{ formatDateTime(option.start_time) }}</div>
                  <div v-if="poll.final_date === option.id" class="badge badge-success badge-sm">Final</div>
                </div>
                <div class="flex gap-3 mt-2">
                  <div class="text-center flex-1">
                    <div class="text-2xl font-bold text-success">{{ option.yes_count }}</div>
                    <div class="text-xs text-base-content/50">OK</div>
                  </div>
                  <div v-if="poll.allow_maybe" class="text-center flex-1">
                    <div class="text-2xl font-bold text-warning">{{ option.maybe_count }}</div>
                    <div class="text-xs text-base-content/50">Maybe</div>
                  </div>
                  <div class="text-center flex-1">
                    <div class="text-2xl font-bold text-error">{{ option.no_count }}</div>
                    <div class="text-xs text-base-content/50">No</div>
                  </div>
                </div>
                <progress
                  class="progress progress-success w-full mt-3"
                  :value="option.yes_count"
                  :max="Math.max(option.total_votes, 1)"
                ></progress>
              </div>
            </div>
          </div>
        </div>

        <!-- Set Final Date (creator only) -->
        <div v-if="isCreator && !poll.final_date" class="mt-8 p-4 bg-success/10 rounded-lg">
          <h3 class="font-semibold mb-3 flex items-center gap-2">
            <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
            Set Final Date
          </h3>
          <p class="text-sm text-base-content/70 mb-3">Once you set a final date, participants will be notified and voting will close.</p>
          <div class="flex gap-2">
            <select v-model="selectedFinalDate" class="select select-bordered flex-1">
              <option value="">Select a date...</option>
              <option v-for="option in dateOptions" :key="option.id" :value="option.id">
                {{ formatDateTime(option.start_time) }}
              </option>
            </select>
            <button
              @click="setFinalDate"
              class="btn btn-success"
              :disabled="!selectedFinalDate || loadingVote"
            >
              Confirm
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Comments -->
    <div class="card bg-base-100 shadow-xl">
      <div class="card-body">
        <h2 class="card-title flex items-center gap-2">
          <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
          </svg>
          Comments ({{ comments.length }})
        </h2>

        <div v-if="authStore.isAuthenticated" class="form-control">
          <textarea
            v-model="newComment"
            class="textarea textarea-bordered"
            placeholder="Write a comment..."
            rows="2"
          ></textarea>
          <button @click="addComment" class="btn btn-primary mt-2" :disabled="!newComment || loadingVote">
            Send Comment
          </button>
        </div>

        <div v-else class="alert alert-soft">
          <router-link to="/login" class="btn btn-sm btn-primary">Login</router-link>
          <span>to post a comment</span>
        </div>

        <div v-if="comments.length > 0" class="space-y-4 mt-4">
          <div v-for="comment in comments" :key="comment.id" class="flex gap-3 p-3 bg-base-200 rounded-lg">
            <div class="avatar placeholder">
              <div class="bg-neutral text-neutral-content rounded-full w-10">
                <span>{{ comment.user?.name?.charAt(0) || '?' }}</span>
              </div>
            </div>
            <div class="flex-1">
              <div class="flex items-center gap-2 mb-1">
                <span class="font-medium">{{ comment.user?.name }}</span>
                <span class="text-xs text-base-content/50">{{ formatRelativeTime(comment.created_at) }}</span>
              </div>
              <p class="text-sm">{{ comment.content }}</p>
            </div>
          </div>
        </div>
        <div v-else class="text-center py-8 text-base-content/50">
          No comments yet. Be the first to comment!
        </div>
      </div>
    </div>
  </div>

  <div v-else class="flex justify-center items-center min-h-screen">
    <div class="alert alert-error max-w-md">
      <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
      </svg>
      <span>Poll not found or may have been deleted.</span>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { usePollsStore } from '@/stores/polls'
import { useUiStore } from '@/stores/ui'

const route = useRoute()
const authStore = useAuthStore()
const pollsStore = usePollsStore()
const uiStore = useUiStore()

const poll = ref(null)
const dateOptions = ref([])
const comments = ref([])
const votesByUser = ref({})
const anonymousVotes = ref([])
const loading = ref(false)
const loadingVote = ref(false)
const newComment = ref('')
const selectedFinalDate = ref('')

const isCreator = computed(() => {
  return authStore.isAuthenticated && poll.value?.creator_id === authStore.user?.id
})

async function loadPoll() {
  loading.value = true
  try {
    const data = await pollsStore.fetchPoll(route.params.id)
    poll.value = data.poll
    dateOptions.value = data.date_options || []
    comments.value = data.comments || []

    // Group votes by user
    const grouped = {}
    const anon = []

    data.votes?.forEach(vote => {
      if (vote.user_id) {
        if (!grouped[vote.user_id]) {
          grouped[vote.user_id] = []
        }
        grouped[vote.user_id].push(vote)
      } else {
        anon.push(vote)
      }
    })

    votesByUser.value = grouped
    anonymousVotes.value = anon
  } catch (error) {
    uiStore.error('Failed to load poll')
  } finally {
    loading.value = false
  }
}

function vote(dateOptionId, response) {
  loadingVote.value = true
  const votes = [{ date_option_id: dateOptionId, response }]
  pollsStore.vote(poll.value.id, votes)
    .then(() => {
      uiStore.success('Vote recorded!')
      loadPoll()
    })
    .catch(error => {
      uiStore.error(error.response?.data?.error || 'Failed to vote')
    })
    .finally(() => {
      loadingVote.value = false
    })
}

function addComment() {
  if (!newComment.value) return

  loadingVote.value = true
  pollsStore.addComment(poll.value.id, newComment.value)
    .then(() => {
      newComment.value = ''
      uiStore.success('Comment added!')
      loadPoll()
    })
    .catch(error => {
      uiStore.error(error.response?.data?.error || 'Failed to add comment')
    })
    .finally(() => {
      loadingVote.value = false
    })
}

function setFinalDate() {
  if (!selectedFinalDate.value) return

  loadingVote.value = true
  pollsStore.setFinalDate(poll.value.id, selectedFinalDate.value)
    .then(() => {
      uiStore.success('Final date set!')
      loadPoll()
    })
    .catch(error => {
      uiStore.error(error.response?.data?.error || 'Failed to set final date')
    })
    .finally(() => {
      loadingVote.value = false
    })
}

function copyLink() {
  navigator.clipboard.writeText(window.location.href)
  uiStore.success('Link copied!')
}

function exportPDF() {
  window.open(`/api/polls/${poll.value.id}/export/pdf`, '_blank')
}

function exportICS() {
  window.open(`/api/polls/${poll.value.id}/export/ics`, '_blank')
}

function formatDate(dateStr) {
  return new Date(dateStr).toLocaleDateString('en-US', {
    weekday: 'short',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

function formatDateShort(dateStr) {
  return new Date(dateStr).toLocaleDateString('en-US', {
    weekday: 'short',
    month: 'short',
    day: 'numeric'
  })
}

function formatTime(dateStr) {
  return new Date(dateStr).toLocaleTimeString('en-US', {
    hour: '2-digit',
    minute: '2-digit'
  })
}

function formatDateTime(dateStr) {
  return new Date(dateStr).toLocaleString('en-US', {
    weekday: 'short',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

function formatRelativeTime(dateStr) {
  const date = new Date(dateStr)
  const now = new Date()
  const seconds = Math.floor((now - date) / 1000)

  if (seconds < 60) return 'just now'
  if (seconds < 3600) return `${Math.floor(seconds / 60)}m ago`
  if (seconds < 86400) return `${Math.floor(seconds / 3600)}h ago`
  return `${Math.floor(seconds / 86400)}d ago`
}

function getUserDisplayName(votes) {
  return votes[0]?.user?.name || votes[0]?.user_name || 'Anonymous'
}

function getVoteResponse(dateOptionId, votes) {
  const vote = votes.find(v => v.date_option_id === dateOptionId)
  return vote?.response
}

function getVoteIcon(response) {
  const icons = {
    yes: '✓',
    no: '✗',
    maybe: '?'
  }
  return icons[response] || ''
}

function getUserVoteClass(dateOptionId, response) {
  const userVotes = votesByUser.value[authStore.user?.id] || []
  const existingVote = userVotes.find(v => v.date_option_id === dateOptionId)

  if (existingVote?.response === response) {
    return { yes: 'btn-success', no: 'btn-error', maybe: 'btn-warning' }[response]
  }
  return 'btn-outline'
}

onMounted(() => {
  loadPoll()
})
</script>

<style scoped>
.vote-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  font-weight: bold;
}

.vote-badge.yes {
  background-color: hsl(var(--su));
  color: white;
}

.vote-badge.no {
  background-color: hsl(var(--er));
  color: white;
}

.vote-badge.maybe {
  background-color: hsl(var(--wa));
  color: white;
}

.vote-matrix :deep(th) {
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}
</style>
