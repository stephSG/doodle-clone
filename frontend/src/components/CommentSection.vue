<template>
  <div class="space-y-4">
    <!-- Comment List -->
    <div class="space-y-3">
      <div v-for="comment in comments" :key="comment.id" class="comment">
        <div class="flex items-start gap-3">
          <div class="avatar placeholder">
            <div class="bg-neutral text-neutral-content rounded-full w-10">
              <span>{{ comment.user?.name?.charAt(0) || '?' }}</span>
            </div>
          </div>
          <div class="flex-1">
            <div class="flex items-center gap-2">
              <span class="font-medium">{{ comment.user?.name || 'Unknown' }}</span>
              <span class="text-xs text-base-content/50">{{ formatRelativeTime(comment.created_at) }}</span>
              <button
                v-if="canDelete(comment)"
                @click="$emit('delete', comment)"
                class="btn btn-ghost btn-xs text-error"
              >
                Delete
              </button>
            </div>
            <p class="mt-1">{{ comment.content }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Add Comment Form -->
    <div class="card bg-base-200">
      <div class="card-body p-4">
        <textarea
          v-model="newComment"
          class="textarea textarea-bordered textarea-sm"
          placeholder="Write a comment..."
          rows="2"
          @keydown.ctrl.enter="submitComment"
        ></textarea>
        <div class="flex justify-between items-center mt-2">
          <span class="text-xs text-base-content/50">Ctrl+Enter to submit</span>
          <button
            @click="submitComment"
            class="btn btn-primary btn-sm"
            :disabled="!newComment.trim()"
          >
            Post Comment
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'

const props = defineProps({
  comments: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['add', 'delete'])

const authStore = useAuthStore()
const newComment = ref('')

function submitComment() {
  if (newComment.value.trim()) {
    emit('add', newComment.value)
    newComment.value = ''
  }
}

function canDelete(comment) {
  return authStore.isAuthenticated && (
    comment.user_id === authStore.user?.id ||
    authStore.user?.id === comment.poll?.creator_id
  )
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
</script>
