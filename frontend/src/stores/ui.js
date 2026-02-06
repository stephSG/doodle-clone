import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUiStore = defineStore('ui', () => {
  const notifications = ref([])
  const theme = ref(localStorage.getItem('theme') || 'light')

  function toast(message, type = 'info') {
    const id = Date.now()
    const notification = { id, message, type }
    notifications.value.push(notification)

    // Auto-remove after 3 seconds
    setTimeout(() => {
      removeNotification(id)
    }, 3000)

    // Also show DaisyUI toast
    showToastElement(message, type)
  }

  function removeNotification(id) {
    const index = notifications.value.findIndex(n => n.id === id)
    if (index > -1) {
      notifications.value.splice(index, 1)
    }
  }

  function showToastElement(message, type) {
    const container = document.getElementById('toast-container')
    if (!container) return

    const alertClass = {
      success: 'alert-success',
      error: 'alert-error',
      warning: 'alert-warning',
      info: 'alert-info'
    }[type] || 'alert-info'

    const toast = document.createElement('div')
    toast.className = `alert ${alertClass} shadow-lg`
    toast.innerHTML = `<span>${message}</span>`

    container.appendChild(toast)

    setTimeout(() => {
      toast.remove()
    }, 3000)
  }

  function success(message) {
    toast(message, 'success')
  }

  function error(message) {
    toast(message, 'error')
  }

  function warning(message) {
    toast(message, 'warning')
  }

  function info(message) {
    toast(message, 'info')
  }

  function setTheme(newTheme) {
    theme.value = newTheme
    localStorage.setItem('theme', newTheme)
    document.documentElement.setAttribute('data-theme', newTheme)
  }

  function toggleTheme() {
    const newTheme = theme.value === 'light' ? 'dark' : 'light'
    setTheme(newTheme)
  }

  return {
    notifications,
    theme,
    toast,
    success,
    error,
    warning,
    info,
    removeNotification,
    setTheme,
    toggleTheme
  }
})
