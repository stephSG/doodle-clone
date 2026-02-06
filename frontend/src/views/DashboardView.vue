<template>
  <div class="space-y-8">
    <h1 class="text-3xl font-bold">Dashboard</h1>

    <!-- Stats -->
    <div class="stats stats-vertical lg:stats-horizontal shadow w-full">
      <div class="stat">
        <div class="stat-figure text-primary">
          <svg class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
          </svg>
        </div>
        <div class="stat-title">Polls Created</div>
        <div class="stat-value text-primary">{{ myPolls.length }}</div>
      </div>

      <div class="stat">
        <div class="stat-figure text-secondary">
          <svg class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <div class="stat-title">Votes Cast</div>
        <div class="stat-value text-secondary">{{ myVotes.length }}</div>
      </div>
    </div>

    <!-- Tabs -->
    <div role="tablist" class="tabs tabs-bordered">
      <input type="radio" name="dashboard_tabs" role="tab" class="tab" aria-label="My Polls" :checked="activeTab === 'polls'" @change="activeTab = 'polls'" />
      <div role="tabpanel" class="tab-content pt-4">
        <div v-if="loading" class="flex justify-center">
          <span class="loading loading-spinner loading-lg"></span>
        </div>

        <div v-else-if="myPolls.length === 0" class="text-center py-12">
          <p class="text-lg mb-4">You haven't created any polls yet</p>
          <router-link to="/create" class="btn btn-primary">Create Your First Poll</router-link>
        </div>

        <div v-else class="grid gap-4 md:grid-cols-2">
          <PollCard
            v-for="poll in myPolls"
            :key="poll.id"
            :poll="poll"
            :show-actions="true"
            @edit="editPoll"
            @delete="deletePoll"
          />
        </div>
      </div>

      <input type="radio" name="dashboard_tabs" role="tab" class="tab" aria-label="My Votes" :checked="activeTab === 'votes'" @change="activeTab = 'votes'" />
      <div role="tabpanel" class="tab-content pt-4">
        <div v-if="loading" class="flex justify-center">
          <span class="loading loading-spinner loading-lg"></span>
        </div>

        <div v-else-if="myVotes.length === 0" class="text-center py-12">
          <p class="text-lg">You haven't voted on any polls yet</p>
        </div>

        <div v-else class="overflow-x-auto">
          <table class="table">
            <thead>
              <tr>
                <th>Poll</th>
                <th>Date</th>
                <th>Your Vote</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="vote in myVotes" :key="vote.id">
                <td>
                  <router-link :to="`/poll/${vote.poll_id}`" class="link link-primary">
                    {{ vote.poll_title }}
                  </router-link>
                  <div v-if="vote.poll_location" class="text-sm text-base-content/70">{{ vote.poll_location }}</div>
                </td>
                <td>{{ formatDateTime(vote.start_time) }}</td>
                <td>
                  <span :class="{
                    'badge badge-success': vote.response === 'yes',
                    'badge badge-error': vote.response === 'no',
                    'badge badge-warning': vote.response === 'maybe'
                  }">
                    {{ vote.response.toUpperCase() }}
                  </span>
                </td>
                <td>
                  <router-link :to="`/poll/${vote.poll_id}`" class="btn btn-ghost btn-sm">View</router-link>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { usePollsStore } from '@/stores/polls'
import PollCard from '@/components/PollCard.vue'

const router = useRouter()
const pollsStore = usePollsStore()

const activeTab = ref('polls')
const loading = ref(false)
const myPolls = ref([])
const myVotes = ref([])

async function loadMyPolls() {
  loading.value = true
  try {
    const data = await pollsStore.fetchUserPolls()
    myPolls.value = data.polls || []
  } catch (error) {
    console.error('Error loading polls:', error)
  } finally {
    loading.value = false
  }
}

async function loadMyVotes() {
  loading.value = true
  try {
    const data = await pollsStore.fetchUserVotes()
    myVotes.value = data.votes || []
  } catch (error) {
    console.error('Error loading votes:', error)
  } finally {
    loading.value = false
  }
}

function editPoll(poll) {
  router.push(`/poll/${poll.id}/edit`)
}

async function deletePoll(poll) {
  if (!confirm(`Are you sure you want to delete "${poll.title}"?`)) return

  try {
    await pollsStore.deletePoll(poll.id)
    myPolls.value = myPolls.value.filter(p => p.id !== poll.id)
  } catch (error) {
    console.error('Error deleting poll:', error)
  }
}

function formatDateTime(dateStr) {
  return new Date(dateStr).toLocaleString()
}

onMounted(() => {
  loadMyPolls()
  loadMyVotes()
})
</script>
