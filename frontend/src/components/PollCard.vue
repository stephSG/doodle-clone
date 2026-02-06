<template>
  <div class="poll-card card">
    <div class="card-body">
      <h3 class="card-title">
        <router-link :to="`/poll/${poll.id}`" class="hover:text-primary transition-colors">
          {{ poll.title }}
        </router-link>
      </h3>

      <p v-if="poll.description" class="text-sm text-base-content/70 line-clamp-2">
        {{ poll.description }}
      </p>

      <div class="flex flex-wrap gap-2 mt-2">
        <div v-if="poll.location" class="badge badge-ghost">
          <svg class="w-3 h-3 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
          </svg>
          {{ poll.location }}
        </div>

        <div v-if="poll.expires_at" class="badge badge-outline">
          <svg class="w-3 h-3 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          {{ formatExpiration(poll.expires_at) }}
        </div>

        <div v-if="poll.final_date" class="badge badge-success">
          ✓ Final date selected
        </div>
      </div>

      <div class="flex items-center justify-between mt-4">
        <div class="flex items-center gap-2 text-sm text-base-content/60">
          <div class="avatar placeholder">
            <div class="bg-neutral text-neutral-content rounded-full w-6">
              <span class="text-xs">{{ poll.creator?.name?.charAt(0) || '?' }}</span>
            </div>
          </div>
          <span>{{ poll.creator?.name || 'Unknown' }}</span>
        </div>

        <div v-if="showActions" class="dropdown dropdown-end">
          <div tabindex="0" role="button" class="btn btn-ghost btn-sm btn-circle">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z" />
            </svg>
          </div>
          <ul tabindex="0" class="dropdown-content z-[1] menu p-2 shadow bg-base-100 rounded-box w-32">
            <li><a @click="$emit('edit', poll)">Edit</a></li>
            <li><a @click="$emit('delete', poll)" class="text-error">Delete</a></li>
          </ul>
        </div>

        <router-link v-else :to="`/poll/${poll.id}`" class="btn btn-primary btn-sm">
          View Poll
        </router-link>
      </div>

      <!-- Vote count bar -->
      <div v-if="poll.participant_count !== undefined" class="mt-3">
        <div class="flex justify-between text-xs mb-1">
          <span>{{ poll.participant_count }} participant{{ poll.participant_count !== 1 ? 's' : '' }}</span>
        </div>
        <progress class="progress progress-primary w-full" :value="Math.min(poll.participant_count * 10, 100)" max="100"></progress>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  poll: {
    type: Object,
    required: true
  },
  showActions: {
    type: Boolean,
    default: false
  }
})

defineEmits(['edit', 'delete'])

function formatExpiration(dateStr) {
  const date = new Date(dateStr)
  const now = new Date()
  const diff = date - now

  if (diff < 0) return 'Expired'
  if (diff < 86400000) return 'Expires today'
  if (diff < 172800000) return 'Expires tomorrow'
  return `Expires in ${Math.ceil(diff / 86400000)} days`
}
</script>
