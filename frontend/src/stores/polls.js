import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/utils/api'

export const usePollsStore = defineStore('polls', () => {
  const polls = ref([])
  const currentPoll = ref(null)
  const userPolls = ref([])
  const userVotes = ref([])
  const loading = ref(false)

  async function fetchPolls(params = {}) {
    loading.value = true
    try {
      const response = await api.get('/api/polls', { params })
      polls.value = response.data.polls || []
      return response.data
    } catch (error) {
      console.error('Error fetching polls:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  async function fetchPoll(id) {
    loading.value = true
    try {
      const response = await api.get(`/api/polls/${id}`)
      currentPoll.value = response.data.poll
      return response.data
    } catch (error) {
      console.error('Error fetching poll:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  async function createPoll(pollData) {
    loading.value = true
    try {
      const response = await api.post('/api/polls', pollData)
      polls.value.unshift(response.data)
      return response.data
    } catch (error) {
      console.error('Error creating poll:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  async function updatePoll(id, pollData) {
    loading.value = true
    try {
      const response = await api.put(`/api/polls/${id}`, pollData)
      if (currentPoll.value?.id === id) {
        currentPoll.value = { ...currentPoll.value, ...response.data }
      }
      return response.data
    } catch (error) {
      console.error('Error updating poll:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  async function deletePoll(id) {
    loading.value = true
    try {
      await api.delete(`/api/polls/${id}`)
      polls.value = polls.value.filter(p => p.id !== id)
      currentPoll.value = null
      return true
    } catch (error) {
      console.error('Error deleting poll:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  async function setFinalDate(pollId, dateOptionId) {
    loading.value = true
    try {
      const response = await api.post(`/api/polls/${pollId}/final`, { date_option_id: dateOptionId })
      if (currentPoll.value?.id === pollId) {
        currentPoll.value = { ...currentPoll.value, final_date: dateOptionId }
      }
      return response.data
    } catch (error) {
      console.error('Error setting final date:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  async function vote(pollId, votes, userName = null) {
    loading.value = true
    try {
      const payload = userName ? { votes, user_name: userName } : { votes }
      // Use /vote endpoint (singular) which supports optional auth for anonymous voting
      const response = await api.post(`/api/polls/${pollId}/vote`, payload)
      await fetchPoll(pollId) // Refresh poll data
      return response.data
    } catch (error) {
      console.error('Error voting:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  async function updateVote(pollId, voteId, response) {
    loading.value = true
    try {
      const res = await api.put(`/api/polls/${pollId}/votes/${voteId}`, { response })
      await fetchPoll(pollId)
      return res.data
    } catch (error) {
      console.error('Error updating vote:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  async function deleteVote(pollId, voteId) {
    loading.value = true
    try {
      await api.delete(`/api/polls/${pollId}/votes/${voteId}`)
      await fetchPoll(pollId)
      return true
    } catch (error) {
      console.error('Error deleting vote:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  async function addComment(pollId, content) {
    loading.value = true
    try {
      const response = await api.post(`/api/polls/${pollId}/comments`, { content })
      await fetchPoll(pollId)
      return response.data
    } catch (error) {
      console.error('Error adding comment:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  async function deleteComment(pollId, commentId) {
    loading.value = true
    try {
      await api.delete(`/api/polls/${pollId}/comments/${commentId}`)
      await fetchPoll(pollId)
      return true
    } catch (error) {
      console.error('Error deleting comment:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  async function fetchUserPolls() {
    loading.value = true
    try {
      const response = await api.get('/api/user/polls')
      userPolls.value = response.data.polls || []
      return response.data
    } catch (error) {
      console.error('Error fetching user polls:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  async function fetchUserVotes() {
    loading.value = true
    try {
      const response = await api.get('/api/user/votes')
      userVotes.value = response.data.votes || []
      return response.data
    } catch (error) {
      console.error('Error fetching user votes:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  function clearCurrentPoll() {
    currentPoll.value = null
  }

  return {
    polls,
    currentPoll,
    userPolls,
    userVotes,
    loading,
    fetchPolls,
    fetchPoll,
    createPoll,
    updatePoll,
    deletePoll,
    setFinalDate,
    vote,
    updateVote,
    deleteVote,
    addComment,
    deleteComment,
    fetchUserPolls,
    fetchUserVotes,
    clearCurrentPoll
  }
})
