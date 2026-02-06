<template>
  <div class="vote-matrix">
    <table class="table table-zebra table-sm">
      <thead>
        <tr>
          <th class="w-48">Participant</th>
          <th v-for="option in options" :key="option.id" class="text-center min-w-20">
            <div class="text-xs font-normal">{{ formatTime(option.start_time) }}</div>
            <div v-if="showStats" class="text-xs text-base-content/60 mt-1">
              ✓ {{ option.yes_count }} ? {{ option.maybe_count }}
            </div>
          </th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(participant, idx) in participants" :key="idx">
          <td class="font-medium">{{ participant.name }}</td>
          <td v-for="option in options" :key="option.id" class="text-center p-1">
            <button
              :class="['w-10 h-10 rounded flex items-center justify-center transition-all', getCellClass(option.id, participant.id)]"
              :disabled="!interactive"
              @click="$emit('vote', option.id, participant.id)"
            >
              {{ getCellIcon(option.id, participant.id) }}
            </button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup>
defineProps({
  options: {
    type: Array,
    required: true
  },
  participants: {
    type: Array,
    required: true
  },
  votes: {
    type: Object,
    default: () => ({})
  },
  interactive: {
    type: Boolean,
    default: true
  },
  showStats: {
    type: Boolean,
    default: true
  }
})

defineEmits(['vote'])

function formatTime(dateStr) {
  const date = new Date(dateStr)
  return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

function getCellClass(optionId, participantId) {
  const response = getVoteResponse(optionId, participantId)
  switch (response) {
    case 'yes': return 'bg-success text-success-content hover:bg-success/80'
    case 'no': return 'bg-error text-error-content hover:bg-error/80'
    case 'maybe': return 'bg-warning text-warning-content hover:bg-warning/80'
    default: return 'bg-base-300 hover:bg-base-200'
  }
}

function getCellIcon(optionId, participantId) {
  const response = getVoteResponse(optionId, participantId)
  switch (response) {
    case 'yes': return '✓'
    case 'no': return '✗'
    case 'maybe': return '?'
    default: return ''
  }
}

function getVoteResponse(optionId, participantId) {
  const participantVotes = props.votes[participantId]
  if (!participantVotes) return null

  const vote = participantVotes.find(v => v.date_option_id === optionId)
  return vote?.response || null
}
</script>
