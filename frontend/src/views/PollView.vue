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

    <!-- Error State -->
    <div v-else-if="!poll" class="flex justify-center items-center min-h-screen p-6">
      <div class="bg-rose-50 text-rose-600 p-6 rounded-3xl text-center">
        <svg class="w-12 h-12 mx-auto mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
        <p class="font-bold">Sondage non trouvé</p>
      </div>
    </div>

    <!-- Poll Detail -->
    <template v-else>
      <!-- Premium Header -->
      <div class="bg-indigo-600 p-8 pt-12 pb-10 rounded-b-[48px] text-white shadow-2xl relative overflow-hidden">
        <div class="absolute -top-10 -right-10 w-48 h-48 bg-white opacity-5 rounded-full blur-3xl"></div>
        <div class="flex justify-between items-center mb-8 relative z-10">
          <button @click="$router.back()" class="w-11 h-11 bg-white/10 rounded-xl backdrop-blur-md flex items-center justify-center active:scale-90 transition-all">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
              <path d="m15 18-6-6 6-6"/>
            </svg>
          </button>
          <button @click="isShareSheetOpen = true" class="px-5 py-2.5 bg-white/20 rounded-2xl backdrop-blur-md flex items-center gap-2 text-xs font-black border border-white/20 active:scale-90 transition-all uppercase tracking-wider">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
              <path d="M4 12v8a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-8"/>
              <polyline points="16 6 12 2 8 6"/>
              <line x1="12" x2="12" y1="2" y2="15"/>
            </svg>
            Partager
          </button>
        </div>
        <h1 class="text-3xl font-black leading-tight mb-4">{{ poll.title }}</h1>
        <div class="flex flex-wrap gap-3 items-center">
          <div v-if="poll.location" class="flex items-center gap-2 text-indigo-50 text-[11px] font-bold bg-white/10 px-3 py-1.5 rounded-full backdrop-blur-sm">
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
              <path d="M20 10c0 6-8 12-8 12s-8-6-8-12a8 8 0 0 1 16 0Z"/>
              <circle cx="12" cy="10" r="3"/>
            </svg>
            {{ poll.location }}
          </div>
          <div class="flex items-center gap-2 text-indigo-50 text-[11px] font-bold bg-white/10 px-3 py-1.5 rounded-full backdrop-blur-sm">
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
              <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/>
              <circle cx="9" cy="7" r="4"/>
              <path d="M22 21v-2a4 4 0 0 0-3-3.87"/>
              <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
            </svg>
            {{ participantCount }} participants
          </div>
          <div v-if="poll.final_date" class="flex items-center gap-2 text-emerald-50 text-[11px] font-bold bg-emerald-500/30 px-3 py-1.5 rounded-full backdrop-blur-sm">
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round">
              <path d="M5 13l4 4L19 7"/>
            </svg>
            Date fixée
          </div>
        </div>
      </div>

      <!-- Tab Switcher -->
      <div class="mx-6 -mt-8 bg-white/90 backdrop-blur-xl p-1.5 rounded-[28px] flex shadow-xl border border-indigo-50/50 relative z-20">
        <button @click="eventViewMode = 'vote'" :class="['flex-1 py-4 rounded-[22px] text-sm font-black transition-all duration-300', eventViewMode === 'vote' ? 'bg-indigo-600 text-white shadow-lg' : 'text-slate-400']">
          Voter
        </button>
        <button @click="eventViewMode = 'results'" :class="['flex-1 py-4 rounded-[22px] text-sm font-black transition-all duration-300', eventViewMode === 'results' ? 'bg-indigo-600 text-white shadow-lg' : 'text-slate-400']">
          Résultats
        </button>
      </div>

      <div class="px-6 mt-10 pb-32">
        <!-- VOTE MODE -->
        <div v-if="eventViewMode === 'vote'" class="space-y-6">
          <div class="flex justify-between items-end mb-2">
            <div>
              <h2 class="font-black text-2xl" :class="isDark ? 'text-white' : 'text-slate-800'">C'est à vous !</h2>
              <p class="text-sm font-medium" :class="isDark ? 'text-slate-400' : 'text-slate-400'">Sélectionnez vos disponibilités.</p>
            </div>
            <div class="bg-indigo-50 text-indigo-600 text-[10px] font-black px-3 py-1.5 rounded-full uppercase tracking-tighter border border-indigo-100">
              {{ selectedSlots.length }} choix
            </div>
          </div>

          <!-- Anonymous Name Input (only in main content, not in fixed bar) -->
          <div v-if="!authStore.isAuthenticated" class="p-5 rounded-3xl border mb-6" :class="isDark ? 'bg-slate-800 border-slate-700' : 'bg-slate-50 border-slate-100'">
            <label class="text-[11px] font-black text-indigo-500 uppercase tracking-widest mb-2 block">Votre nom</label>
            <input
              v-model="anonymousName"
              type="text"
              placeholder="Entrez votre nom..."
              class="w-full h-14 px-5 rounded-2xl outline-none transition-all font-bold border-2 border-transparent"
              :class="[
                isDark
                  ? 'bg-slate-700 text-white placeholder:text-slate-500 focus:bg-slate-600 focus:border-indigo-500'
                  : 'bg-white text-slate-700 placeholder:text-slate-300 focus:border-indigo-100'
              ]"
            />
          </div>

          <!-- Date Options -->
          <div class="grid gap-3">
            <div
              v-for="option in dateOptions"
              :key="option.id"
              @click="toggleSlot(option.id)"
              :class="['group p-5 rounded-[32px] border-2 transition-all flex items-center gap-5 cursor-pointer relative overflow-hidden active:scale-[0.97]', selectedSlots.includes(option.id) ? 'border-indigo-500 bg-indigo-500/20' : isDark ? 'border-slate-700 bg-slate-800/50' : 'border-slate-50 bg-slate-50/50']"
            >
              <div :class="['w-16 h-16 rounded-[22px] flex items-center justify-center shrink-0 transition-all', selectedSlots.includes(option.id) ? 'bg-indigo-600 text-white shadow-lg' : isDark ? 'bg-slate-700 text-slate-400 shadow-sm' : 'bg-white text-slate-300 shadow-sm']">
                <svg xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <rect width="18" height="18" x="3" y="4" rx="2" ry="2"/>
                  <line x1="16" x2="16" y1="2" y2="6"/>
                  <line x1="8" x2="8" y1="2" y2="6"/>
                  <line x1="3" x2="21" y1="10" y2="10"/>
                </svg>
              </div>
              <div class="flex-1">
                <p class="font-black text-lg leading-none transition-colors" :class="[selectedSlots.includes(option.id) ? 'text-indigo-500' : '', isDark ? 'text-slate-200' : 'text-slate-700']">{{ formatDayOfWeek(option.start_time) }}</p>
                <p class="text-sm font-bold mt-2" :class="isDark ? 'text-slate-400' : 'text-slate-400'">{{ formatTime(option.start_time) }}</p>
              </div>
              <div :class="['w-10 h-10 rounded-2xl border-2 flex items-center justify-center shrink-0 transition-all', selectedSlots.includes(option.id) ? 'bg-indigo-600 border-indigo-600 text-white scale-110' : isDark ? 'border-slate-600 bg-slate-700' : 'border-slate-200 bg-white']">
                <svg v-if="selectedSlots.includes(option.id)" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="4" stroke-linecap="round" stroke-linejoin="round">
                  <polyline points="20 6 9 17 4 12"/>
                </svg>
                <div v-else class="w-1.5 h-1.5 bg-slate-200 rounded-full"></div>
              </div>
              <div v-if="poll.final_date === option.id" class="absolute top-0 right-0 bg-emerald-400 text-white text-[8px] font-black px-2 py-1 rounded-br-2xl rounded-tl-lg uppercase tracking-widest shadow-sm">
                Final
              </div>
            </div>
          </div>

          <!-- Submit Vote Button -->
          <div class="fixed bottom-0 left-0 right-0 p-6 z-40" :class="isDark ? 'bg-slate-900/90 backdrop-blur-md' : 'bg-gradient-to-t from-white via-white to-transparent'">
            <div class="max-w-md mx-auto flex gap-3">
              <input
                v-if="!authStore.isAuthenticated"
                type="text"
                v-model="anonymousName"
                placeholder="Votre nom"
                :class="['flex-1 h-16 px-6 rounded-2xl shadow-2xl outline-none font-bold transition-all', isDark ? 'bg-slate-800 border-slate-700 text-white placeholder:text-slate-500 focus:border-indigo-500' : 'bg-white border-2 border-slate-100 text-slate-700 placeholder:text-slate-300 focus:border-indigo-100']"
              />
              <button
                @click="submitVote"
                :disabled="!canSubmitVote"
                :class="['flex-1 h-16 px-8 rounded-2xl font-black shadow-2xl disabled:opacity-50 active:scale-[0.97] transition-all flex items-center justify-center gap-2', canSubmitVote ? 'bg-indigo-600 text-white hover:bg-indigo-700' : 'bg-slate-300 text-slate-500']"
              >
                <span v-if="loadingVote" class="loading loading-spinner loading-sm"></span>
                <span v-else>{{ authStore.isAuthenticated ? 'Confirmer' : 'Voter' }}</span>
                <svg v-if="!loadingVote" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                  <polyline points="20 6 9 17 4 12"/>
                </svg>
              </button>
            </div>
          </div>
        </div>

        <!-- RESULTS MODE -->
        <div v-if="eventViewMode === 'results'" class="space-y-10 pb-10">
          <div>
            <h2 class="font-black text-2xl flex items-center gap-3" :class="isDark ? 'text-white' : 'text-slate-800'">
              Tendances
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-amber-400">
                <path d="M6 9H4.5a2.5 2.5 0 0 1 0-5H6"/>
                <path d="M18 9h1.5a2.5 2.5 0 0 0 0-5H18"/>
                <path d="M4 22h16"/>
                <path d="M10 14.66V17c0 .55-.47.98-.97 1.21C7.85 18.75 7 20.24 7 22"/>
                <path d="M14 14.66V17c0 .55.47.98.97 1.21C16.15 18.75 17 20.24 17 22"/>
                <path d="M18 2H6v7a6 6 0 0 0 12 0V2Z"/>
              </svg>
            </h2>
            <p class="text-slate-400 text-sm font-medium mt-1">L'option favorite se dessine.</p>
          </div>

          <!-- Results Bars -->
          <div class="space-y-4">
            <div
              v-for="option in sortedOptions"
              :key="option.id"
              class="p-6 rounded-[32px] shadow-sm border transition-all"
              :class="[
                isDark ? 'bg-slate-800' : 'bg-white',
                option.isWinner ? 'border-indigo-500 bg-indigo-500/20' : (isDark ? 'border-slate-700' : 'border-slate-50')
              ]"
            >
              <div class="flex justify-between items-start mb-4">
                <div>
                  <p class="font-black text-lg" :class="[option.isWinner ? 'text-indigo-400' : '', isDark ? 'text-white' : 'text-slate-800']">{{ formatDayOfWeek(option.start_time) }}</p>
                  <p class="text-xs font-bold uppercase tracking-widest mt-1" :class="isDark ? 'text-slate-500' : 'text-slate-400'">{{ formatTime(option.start_time) }}</p>
                </div>
                <div class="text-right">
                  <span :class="['text-2xl font-black', option.isWinner ? 'text-indigo-600' : 'text-slate-400']">{{ option.yesCount }}</span>
                  <p class="text-[9px] font-black uppercase text-slate-300 tracking-tighter">VOTES OUI</p>
                </div>
              </div>
              <div class="h-4 w-full bg-slate-50 rounded-full overflow-hidden border border-slate-100/50">
                <div :class="['h-full transition-all duration-1000 ease-out', option.isWinner ? 'bg-gradient-to-r from-indigo-500 to-violet-500 shadow-lg' : 'bg-slate-200']" :style="{ width: option.percentage + '%' }"></div>
              </div>
            </div>
          </div>

          <!-- Participants List -->
          <div class="mt-12">
            <h3 class="font-black text-lg mb-6 flex items-center gap-2" :class="isDark ? 'text-white' : 'text-slate-800'">
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" class="text-indigo-500">
                <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/>
                <circle cx="9" cy="7" r="4"/>
                <path d="M22 21v-2a4 4 0 0 0-3-3.87"/>
                <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
              </svg>
              Participants ({{ participantCount }})
            </h3>
            <div class="grid grid-cols-1 gap-3">
              <div
                v-for="(votes, userId) in votesByUser"
                :key="userId"
                class="flex items-center gap-4 p-4 rounded-3xl border"
                :class="isDark ? 'bg-slate-800 border-slate-700' : 'bg-slate-50 border-slate-100/50'"
              >
                <div class="w-11 h-11 rounded-2xl bg-indigo-100 text-indigo-600 flex items-center justify-center font-black text-lg">
                  {{ getUserDisplayName(votes)?.charAt(0) || '?' }}
                </div>
                <div class="flex-1">
                  <p class="font-black" :class="isDark ? 'text-white' : 'text-slate-800'">{{ getUserDisplayName(votes) }}</p>
                  <p class="text-[10px] font-bold uppercase tracking-widest" :class="isDark ? 'text-slate-500' : 'text-slate-400'">{{ votes.filter(v => v.response === 'yes').length }} choix</p>
                </div>
                <div class="flex -space-x-1">
                  <div
                    v-for="v in votes.filter(v => v.response === 'yes')"
                    :key="v.id"
                    class="w-2.5 h-2.5 rounded-full bg-indigo-400 border-2 border-slate-50"
                  ></div>
                </div>
              </div>
              <div
                v-for="(anonVote, idx) in anonymousVotes"
                :key="'anon-' + idx"
                class="flex items-center gap-4 bg-slate-50 p-4 rounded-3xl border border-slate-100/50"
              >
                <div class="w-11 h-11 rounded-2xl bg-slate-200 text-slate-500 flex items-center justify-center font-black text-lg">
                  {{ anonVote.user_name?.charAt(0) || '?' }}
                </div>
                <div class="flex-1">
                  <p class="font-black text-slate-800">{{ anonVote.user_name }}</p>
                  <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest">Anonyme</p>
                </div>
              </div>
              <div v-if="participantCount === 0" class="text-center py-10 rounded-[32px] border-2 border-dashed" :class="isDark ? 'bg-slate-800 border-slate-700 text-slate-500' : 'bg-slate-50 border-slate-100 text-slate-300'">
                <p class="font-bold text-sm italic">Aucun vote enregistré.</p>
              </div>
            </div>
          </div>

          <!-- Creator Actions -->
          <div v-if="isCreator" class="mt-8 space-y-4">
            <div class="p-6 rounded-[32px] border" :class="isDark ? 'bg-slate-800 border-slate-700' : 'bg-slate-50 border-slate-100'">
              <h3 class="font-black mb-4 flex items-center gap-2" :class="isDark ? 'text-white' : 'text-slate-800'">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" class="text-emerald-500">
                  <path d="M5 13l4 4L19 7"/>
                </svg>
                Fixer la date finale
              </h3>
              <p v-if="!poll.final_date" class="text-sm text-slate-400 mb-4">Une fois la date fixée, les participants seront notifiés.</p>
              <div v-if="!poll.final_date" class="flex gap-2">
                <select v-model="selectedFinalDate" class="flex-1 bg-white border-2 border-slate-100 h-14 px-5 rounded-2xl outline-none font-bold text-slate-700 focus:border-indigo-100">
                  <option value="">Choisir une date...</option>
                  <option v-for="option in dateOptions" :key="option.id" :value="option.id">
                    {{ formatFullDateTime(option.start_time) }}
                  </option>
                </select>
                <button
                  @click="setFinalDate"
                  :disabled="!selectedFinalDate || loadingVote"
                  class="h-14 px-6 bg-emerald-500 text-white rounded-2xl font-bold shadow-lg disabled:bg-slate-200 active:scale-95 transition-all"
                >
                  Confirmer
                </button>
              </div>
              <div v-else class="flex items-center gap-2 text-emerald-600 font-bold">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M5 13l4 4L19 7"/>
                </svg>
                Date fixée : {{ formatFullDateTime(dateOptions.find(o => o.id === poll.final_date)?.start_time) }}
              </div>
            </div>

            <router-link :to="`/poll/${poll.id}/edit`" class="block bg-slate-50 p-6 rounded-[32px] border border-slate-100 text-center">
              <span class="font-black text-indigo-600 flex items-center justify-center gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                  <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                </svg>
                Modifier le sondage
              </span>
            </router-link>

            <div class="flex gap-3">
              <button @click="exportPDF" class="flex-1 bg-slate-50 p-4 rounded-2xl border border-slate-100 font-bold text-slate-600 active:scale-95 transition-all flex items-center justify-center gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                  <polyline points="7 10 12 15 17 10"/>
                  <line x1="12" x2="12" y1="15" y2="3"/>
                </svg>
                PDF
              </button>
              <button @click="exportICS" class="flex-1 bg-slate-50 p-4 rounded-2xl border border-slate-100 font-bold text-slate-600 active:scale-95 transition-all flex items-center justify-center gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                </svg>
                Calendrier
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- SHARE SHEET MODAL -->
      <transition name="fade">
        <div v-if="isShareSheetOpen" class="fixed inset-0 z-[100] flex items-end justify-center">
          <div class="absolute inset-0 bg-slate-900/60 backdrop-blur-md" @click="isShareSheetOpen = false"></div>
          <div class="relative bg-white w-full max-w-md rounded-t-[56px] p-10 shadow-2xl animate-in slide-in-from-bottom">
            <div class="w-16 h-1.5 bg-slate-100 rounded-full mx-auto mb-10"></div>
            <div class="flex justify-between items-center mb-10">
              <div>
                <h3 class="text-2xl font-black text-slate-900">Invitez vos amis</h3>
                <p class="text-slate-400 font-bold text-sm">Collectez les réponses en un clic.</p>
              </div>
              <button @click="isShareSheetOpen = false" class="w-12 h-12 bg-slate-50 rounded-2xl flex items-center justify-center text-slate-400 active:scale-90 transition-all">
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M18 6 6 18"/>
                  <path d="m6 6 12 12"/>
                </svg>
              </button>
            </div>
            <div class="bg-slate-50 p-6 rounded-[32px] border-2 border-slate-100 flex flex-col gap-4 mb-10">
              <div class="truncate text-xs text-slate-400 font-mono font-bold">{{ shareUrl }}</div>
              <button @click="copyLink" :class="['flex items-center justify-center gap-3 h-14 rounded-2xl font-black text-sm transition-all shadow-xl', copyFeedback ? 'bg-emerald-500 text-white' : 'bg-indigo-600 text-white']">
                <svg v-if="copyFeedback" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                  <polyline points="20 6 9 17 4 12"/>
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                  <rect width="14" height="14" x="8" y="8" rx="2" ry="2"/>
                  <path d="M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2"/>
                </svg>
                {{ copyFeedback ? 'COPIÉ !' : 'COPIER LE LIEN' }}
              </button>
            </div>
            <div class="grid grid-cols-4 gap-6">
              <button @click="shareWhatsApp" class="flex flex-col items-center gap-3 cursor-pointer group">
                <div class="w-16 h-16 rounded-[22px] bg-emerald-500 flex items-center justify-center shadow-lg active:scale-90 transition-all text-white group-hover:-translate-y-1">
                  <svg xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M7.9 20A9 9 0 1 0 4 16.1L2 22Z"/>
                    <path d="M8 12h.01"/>
                    <path d="M12 12h.01"/>
                    <path d="M16 12h.01"/>
                  </svg>
                </div>
                <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest">WhatsApp</span>
              </button>
              <button @click="shareSMS" class="flex flex-col items-center gap-3 cursor-pointer group">
                <div class="w-16 h-16 rounded-[22px] bg-sky-500 flex items-center justify-center shadow-lg active:scale-90 transition-all text-white group-hover:-translate-y-1">
                  <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="m22 2-7 20-4-9-9-4Z"/>
                    <path d="M22 2 11 13"/>
                  </svg>
                </div>
                <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest">SMS</span>
              </button>
              <button @click="shareEmail" class="flex flex-col items-center gap-3 cursor-pointer group">
                <div class="w-16 h-16 rounded-[22px] bg-slate-800 flex items-center justify-center shadow-lg active:scale-90 transition-all text-white group-hover:-translate-y-1">
                  <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <rect width="20" height="16" x="2" y="4" rx="2"/>
                    <path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7"/>
                  </svg>
                </div>
                <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Email</span>
              </button>
              <button @click="openLink" class="flex flex-col items-center gap-3 cursor-pointer group">
                <div class="w-16 h-16 rounded-[22px] bg-indigo-600 flex items-center justify-center shadow-lg active:scale-90 transition-all text-white group-hover:-translate-y-1">
                  <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"/>
                    <polyline points="15 3 21 3 21 9"/>
                    <line x1="10" x2="21" y1="14" y2="3"/>
                  </svg>
                </div>
                <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Ouvrir</span>
              </button>
            </div>
          </div>
        </div>
      </transition>
    </template>
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

// Dark mode
const isDark = computed(() => document.documentElement.getAttribute('data-theme') === 'dark')

const poll = ref(null)
const dateOptions = ref([])
const votesByUser = ref({})
const anonymousVotes = ref([])
const loading = ref(false)
const loadingVote = ref(false)
const anonymousName = ref('')
const selectedFinalDate = ref('')
const eventViewMode = ref('vote')
const selectedSlots = ref([])
const isShareSheetOpen = ref(false)
const copyFeedback = ref(false)

const isCreator = computed(() => {
  return authStore.isAuthenticated && poll.value?.creator_id === authStore.user?.id
})

const participantCount = computed(() => {
  const authUsers = Object.keys(votesByUser.value).length
  const anonUsers = anonymousVotes.value.length
  return authUsers + anonUsers
})

const canSubmitVote = computed(() => {
  if (authStore.isAuthenticated) {
    return selectedSlots.value.length > 0
  }
  return selectedSlots.value.length > 0 && anonymousName.value.trim()
})

const shareUrl = computed(() => {
  if (!poll.value?.access_code) return window.location.href
  // Use access_code for sharing - shorter and more user-friendly
  const protocol = window.location.protocol === 'https:' ? 'https:' : 'http:'
  return `${protocol}//${window.location.host}/poll/${poll.value.access_code}`
})

const sortedOptions = computed(() => {
  if (!dateOptions.value.length) return []

  const maxYes = Math.max(...dateOptions.value.map(o => o.yes_count || 0))

  return dateOptions.value.map(option => {
    const yesCount = option.yes_count || 0
    const total = Math.max(participantCount.value, 1)
    return {
      ...option,
      yesCount,
      percentage: total > 0 ? (yesCount / total) * 100 : 0,
      isWinner: yesCount > 0 && yesCount === maxYes
    }
  }).sort((a, b) => b.yesCount - a.yesCount)
})

async function loadPoll() {
  loading.value = true
  try {
    const data = await pollsStore.fetchPoll(route.params.id)
    poll.value = data.poll
    dateOptions.value = data.date_options || []

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

    // Pre-select user's existing votes
    if (authStore.isAuthenticated && authStore.user?.id) {
      const userVotes = grouped[authStore.user.id] || []
      selectedSlots.value = userVotes.filter(v => v.response === 'yes').map(v => v.date_option_id)
    }
  } catch (error) {
    uiStore.error('Impossible de charger le sondage')
  } finally {
    loading.value = false
  }
}

function toggleSlot(optionId) {
  if (selectedSlots.value.includes(optionId)) {
    selectedSlots.value = selectedSlots.value.filter(s => s !== optionId)
  } else {
    selectedSlots.value.push(optionId)
  }
}

async function submitVote() {
  if (!canSubmitVote.value) return

  loadingVote.value = true

  const votes = selectedSlots.value.map(dateOptionId => ({
    date_option_id: dateOptionId,
    response: 'yes'
  }))

  try {
    await pollsStore.vote(poll.value.id, votes, anonymousName.value)
    uiStore.success('Vote enregistré !')
    eventViewMode.value = 'results'
    await loadPoll()
  } catch (error) {
    uiStore.error(error.response?.data?.error || 'Erreur lors du vote')
  } finally {
    loadingVote.value = false
  }
}

function setFinalDate() {
  if (!selectedFinalDate.value) return

  loadingVote.value = true
  pollsStore.setFinalDate(poll.value.id, selectedFinalDate.value)
    .then(() => {
      uiStore.success('Date finale fixée !')
      loadPoll()
    })
    .catch(error => {
      uiStore.error(error.response?.data?.error || 'Erreur')
    })
    .finally(() => {
      loadingVote.value = false
    })
}

function copyLink() {
  navigator.clipboard.writeText(shareUrl.value)
  copyFeedback.value = true
  setTimeout(() => copyFeedback.value = false, 2000)
}

function shareWhatsApp() {
  const text = `Disponibilités pour "${poll.value?.title}" : ${shareUrl.value}`
  window.open(`https://wa.me/?text=${encodeURIComponent(text)}`, '_blank')
}

function shareSMS() {
  const text = `Disponibilités pour "${poll.value?.title}" : ${shareUrl.value}`
  window.location.href = `sms:?body=${encodeURIComponent(text)}`
}

function shareEmail() {
  const subject = encodeURIComponent(`Sondage: ${poll.value?.title}`)
  const body = encodeURIComponent(`Votez pour le sondage "${poll.value?.title}" :\n\n${shareUrl.value}`)
  window.location.href = `mailto:?subject=${subject}&body=${body}`
}

function openLink() {
  window.open(shareUrl.value, '_blank')
}

function exportPDF() {
  window.open(`/api/polls/${poll.value.id}/export/pdf`, '_blank')
}

function exportICS() {
  window.open(`/api/polls/${poll.value.id}/export/ics`, '_blank')
}

function formatDayOfWeek(dateStr) {
  if (!dateStr) return '--'
  return new Date(dateStr).toLocaleDateString('fr-FR', {
    weekday: 'long',
    day: 'numeric',
    month: 'long'
  })
}

function formatTime(dateStr) {
  if (!dateStr) return '--'
  return new Date(dateStr).toLocaleTimeString('fr-FR', {
    hour: '2-digit',
    minute: '2-digit'
  })
}

function formatFullDateTime(dateStr) {
  if (!dateStr) return '--'
  return new Date(dateStr).toLocaleString('fr-FR', {
    weekday: 'short',
    day: 'numeric',
    month: 'short',
    hour: '2-digit',
    minute: '2-digit'
  })
}

function getUserDisplayName(votes) {
  return votes[0]?.user?.name || votes[0]?.user_name || 'Anonyme'
}

onMounted(() => {
  loadPoll()
})
</script>

<style scoped>
:deep(*) {
  -webkit-tap-highlight-color: transparent;
}
</style>
