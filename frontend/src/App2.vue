<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Doodle Mobile - Vue.js Premium</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://unpkg.com/vue@3/dist/vue.global.js"></script>
    <style>
        @import url('https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700;800&display=swap');
        
        body {
            font-family: 'Plus Jakarta Sans', sans-serif;
            -webkit-tap-highlight-color: transparent;
            background-color: #f8fafc;
        }

        [v-cloak] { display: none; }

        /* Animations fluides */
        .fade-enter-active, .fade-leave-active { transition: opacity 0.3s ease; }
        .fade-enter-from, .fade-leave-to { opacity: 0; }

        .slide-in-right-enter-active { transition: all 0.4s cubic-bezier(0.16, 1, 0.3, 1); }
        .slide-in-right-enter-from { transform: translateX(30px); opacity: 0; }

        .animate-in {
            animation: fadeIn 0.6s cubic-bezier(0.16, 1, 0.3, 1) forwards;
        }

        @keyframes fadeIn {
            from { opacity: 0; transform: translateY(15px); }
            to { opacity: 1; transform: translateY(0); }
        }

        /* Personnalisation Scrollbar */
        ::-webkit-scrollbar { width: 0px; }
    </style>
</head>
<body class="bg-slate-50">
    <div id="app" v-cloak>
        <div class="min-h-screen bg-white max-w-md mx-auto shadow-[0_0_80px_rgba(0,0,0,0.03)] overflow-x-hidden relative">
            
            <!-- ÉCRAN DE CHARGEMENT -->
            <div v-if="loading" class="flex flex-col items-center justify-center min-h-screen">
                <div class="relative w-16 h-16">
                    <div class="absolute inset-0 border-[3px] border-indigo-50 rounded-full"></div>
                    <div class="absolute inset-0 border-[3px] border-indigo-600 rounded-full border-t-transparent animate-spin"></div>
                </div>
                <p class="mt-6 text-indigo-600 font-bold tracking-wide animate-pulse">Initialisation...</p>
            </div>

            <template v-else>
                <!-- VUE : ACCUEIL -->
                <div v-if="view === 'home'" class="p-6 pb-32 animate-in">
                    <header class="flex justify-between items-center mb-10 mt-4">
                        <div>
                            <h1 class="text-4xl font-extrabold text-slate-900 tracking-tight">Doodle</h1>
                            <p class="text-slate-400 font-medium mt-1">Vos événements, simplifiés.</p>
                        </div>
                        <div class="w-14 h-14 bg-gradient-to-br from-indigo-500 to-violet-600 rounded-2xl flex items-center justify-center text-white shadow-xl shadow-indigo-100 transition-transform active:scale-90">
                            <svg xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m12 3-1.912 5.813a2 2 0 0 1-1.275 1.275L3 12l5.813 1.912a2 2 0 0 1 1.275 1.275L12 21l1.912-5.813a2 2 0 0 1 1.275-1.275L21 12l-5.813-1.912a2 2 0 0 1-1.275-1.275L12 3Z"/><path d="M5 3v4"/><path d="M19 17v4"/><path d="M3 5h4"/><path d="M17 19h4"/></svg>
                        </div>
                    </header>

                    <div class="space-y-4">
                        <div 
                            v-for="(event, index) in events" 
                            :key="event.id"
                            @click="openEvent(event)"
                            class="group bg-white p-5 rounded-[32px] border border-slate-100 hover:border-indigo-100 hover:shadow-md transition-all cursor-pointer relative overflow-hidden flex items-center gap-5 active:scale-[0.98]"
                        >
                            <div :class="['w-16 h-16 rounded-2xl flex flex-col items-center justify-center font-bold shrink-0 transition-transform group-hover:rotate-2', event.id.toString().startsWith('mock-') ? 'bg-amber-50 text-amber-600' : 'bg-indigo-50 text-indigo-600']">
                                <span class="text-[10px] uppercase tracking-widest opacity-70 mb-0.5">{{ event.options[0]?.date.split(' ')[0] }}</span>
                                <span class="text-xl">{{ event.options[0]?.date.split(' ')[1] }}</span>
                            </div>
                            
                            <div class="flex-1 min-w-0">
                                <h3 class="font-bold text-slate-800 text-lg leading-tight truncate group-hover:text-indigo-600 transition-colors">{{ event.title }}</h3>
                                <div class="flex items-center gap-4 text-xs text-slate-400 mt-2 font-semibold">
                                    <span class="flex items-center gap-1.5"><svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" class="text-indigo-400"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M22 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg> {{ event.responses?.length || 0 }}</span>
                                    <span class="truncate flex items-center gap-1.5"><svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" class="text-indigo-400"><path d="M20 10c0 6-8 12-8 12s-8-6-8-12a8 8 0 0 1 16 0Z"/><circle cx="12" cy="10" r="3"/></svg> {{ event.location }}</span>
                                </div>
                            </div>
                            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" class="text-slate-200 group-hover:text-indigo-300 transition-all"><path d="m9 18 6-6-6-6"/></svg>
                            
                            <div v-if="event.id.toString().startsWith('mock-')" class="absolute top-0 right-0 bg-amber-400 text-white text-[8px] font-black px-3 py-1 rounded-bl-xl uppercase tracking-widest shadow-sm">Démo</div>
                        </div>
                    </div>

                    <div class="fixed bottom-8 left-6 right-6 z-50">
                        <button @click="view = 'create'" class="w-full bg-slate-900 text-white h-16 rounded-2xl font-bold shadow-2xl flex items-center justify-center gap-3 active:scale-95 transition-all group relative overflow-hidden">
                            <div class="absolute inset-0 bg-indigo-600 translate-y-full group-hover:translate-y-0 transition-transform duration-300"></div>
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" class="relative z-10"><path d="M12 5v14"/><path d="M5 12h14"/></svg>
                            <span class="relative z-10">Créer un sondage</span>
                        </button>
                    </div>
                </div>

                <!-- VUE : CRÉATION -->
                <div v-if="view === 'create'" class="p-6 pb-32 animate-in">
                    <header class="flex justify-between items-center mb-10 mt-4">
                        <button @click="view = 'home'" class="w-12 h-12 flex items-center justify-center bg-slate-50 rounded-2xl text-slate-600 active:scale-90 transition-all border border-slate-100">
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m15 18-6-6 6-6"/></svg>
                        </button>
                        <h2 class="text-xl font-extrabold text-slate-900">Nouvel Event</h2>
                        <div class="w-12 h-12"></div> <!-- Spacer -->
                    </header>
                    
                    <div class="space-y-6">
                        <div class="space-y-2">
                            <label class="text-[11px] font-black text-indigo-500 uppercase tracking-widest ml-1">Titre de l'événement</label>
                            <input type="text" v-model="newTitle" placeholder="Ex: Dîner de fin d'année" class="w-full bg-slate-50 border-2 border-transparent focus:border-indigo-100 focus:bg-white h-16 px-6 rounded-2xl outline-none transition-all font-bold text-slate-700 placeholder:text-slate-300" />
                        </div>

                        <div class="space-y-2">
                            <label class="text-[11px] font-black text-indigo-500 uppercase tracking-widest ml-1">Lieu</label>
                            <input type="text" v-model="newLocation" placeholder="Ex: Paris ou Discord" class="w-full bg-slate-50 border-2 border-transparent focus:border-indigo-100 focus:bg-white h-16 px-6 rounded-2xl outline-none transition-all font-bold text-slate-700 placeholder:text-slate-300" />
                        </div>

                        <div class="space-y-3">
                            <label class="text-[11px] font-black text-indigo-500 uppercase tracking-widest ml-1">Dates & Heures</label>
                            <div v-for="(opt, idx) in newOptions" :key="idx" class="flex gap-2 group animate-in">
                                <input type="text" v-model="opt.date" placeholder="Lun 12 Fév" class="flex-[2] bg-slate-50 border-2 border-transparent focus:border-indigo-100 focus:bg-white h-16 px-5 rounded-2xl outline-none transition-all font-bold text-slate-700" />
                                <input type="text" v-model="opt.time" placeholder="19:30" class="flex-1 bg-slate-50 border-2 border-transparent focus:border-indigo-100 focus:bg-white h-16 px-4 rounded-2xl outline-none transition-all font-bold text-slate-700" />
                                <button v-if="newOptions.length > 1" @click="newOptions.splice(idx, 1)" class="w-16 h-16 flex items-center justify-center bg-rose-50 text-rose-500 rounded-2xl active:scale-90 transition-all">
                                    <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M18 6 6 18"/><path d="m6 6 12 12"/></svg>
                                </button>
                            </div>
                            <button @click="newOptions.push({ date: '', time: '' })" class="w-full h-14 border-2 border-dashed border-slate-200 text-slate-400 rounded-2xl font-bold flex items-center justify-center gap-2 mt-2 hover:border-indigo-300 hover:text-indigo-500 transition-all">
                                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M12 5v14"/><path d="M5 12h14"/></svg> Ajouter une date
                            </button>
                        </div>
                    </div>

                    <div class="fixed bottom-8 left-6 right-6 z-50">
                        <button @click="createEvent" :disabled="!newTitle" class="w-full bg-indigo-600 text-white h-16 rounded-2xl font-bold shadow-xl shadow-indigo-100 disabled:bg-slate-200 disabled:text-slate-400 disabled:shadow-none transition-all flex items-center justify-center gap-3 active:scale-95">
                            Publier le sondage <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="m9 18 6-6-6-6"/></svg>
                        </button>
                    </div>
                </div>

                <!-- VUE : DÉTAIL ÉVÉNEMENT -->
                <div v-if="view === 'event'" class="pb-32 animate-in">
                    <!-- HEADER PREMIUM -->
                    <div class="bg-indigo-600 p-8 pt-12 rounded-b-[48px] text-white shadow-2xl relative overflow-hidden">
                        <div class="absolute -top-10 -right-10 w-48 h-48 bg-white opacity-5 rounded-full blur-3xl"></div>
                        <div class="flex justify-between items-center mb-8 relative z-10">
                            <button @click="view = 'home'" class="w-11 h-11 bg-white/10 rounded-xl backdrop-blur-md flex items-center justify-center active:scale-90 transition-all">
                                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="m15 18-6-6 6-6"/></svg>
                            </button>
                            <button @click="isShareSheetOpen = true" class="px-5 py-2.5 bg-white/20 rounded-2xl backdrop-blur-md flex items-center gap-2 text-xs font-black border border-white/20 active:scale-90 transition-all uppercase tracking-wider">
                                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M4 12v8a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-8"/><polyline points="16 6 12 2 8 6"/><line x1="12" x2="12" y1="2" y2="15"/></svg> Partager
                            </button>
                        </div>
                        <h1 class="text-3xl font-black leading-tight mb-4">{{ activeEvent.title }}</h1>
                        <div class="flex flex-wrap gap-3 items-center">
                            <div class="flex items-center gap-2 text-indigo-50 text-[11px] font-bold bg-white/10 px-3 py-1.5 rounded-full backdrop-blur-sm">
                                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M20 10c0 6-8 12-8 12s-8-6-8-12a8 8 0 0 1 16 0Z"/><circle cx="12" cy="10" r="3"/></svg> {{ activeEvent.location }}
                            </div>
                            <div class="flex items-center gap-2 text-indigo-50 text-[11px] font-bold bg-white/10 px-3 py-1.5 rounded-full backdrop-blur-sm">
                                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M22 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg> {{ activeEvent.responses?.length || 0 }} participants
                            </div>
                        </div>
                    </div>

                    <!-- TAB SWITCHER -->
                    <div class="mx-6 -mt-8 bg-white/90 backdrop-blur-xl p-1.5 rounded-[28px] flex shadow-xl border border-indigo-50/50 relative z-20">
                        <button @click="eventViewMode = 'vote'" :class="['flex-1 py-4 rounded-[22px] text-sm font-black transition-all duration-300 flex items-center justify-center gap-2', eventViewMode === 'vote' ? 'bg-indigo-600 text-white shadow-lg' : 'text-slate-400']">
                            Voter
                        </button>
                        <button @click="eventViewMode = 'results'" :class="['flex-1 py-4 rounded-[22px] text-sm font-black transition-all duration-300 flex items-center justify-center gap-2', eventViewMode === 'results' ? 'bg-indigo-600 text-white shadow-lg' : 'text-slate-400']">
                            Résultats
                        </button>
                    </div>

                    <div class="px-6 mt-10">
                        <!-- MODE : VOTE -->
                        <div v-if="eventViewMode === 'vote'" class="space-y-6">
                            <div class="flex justify-between items-end mb-2">
                                <div>
                                    <h2 class="font-black text-2xl text-slate-800">C'est à vous !</h2>
                                    <p class="text-slate-400 text-sm font-medium">Sélectionnez vos disponibilités.</p>
                                </div>
                                <div class="bg-indigo-50 text-indigo-600 text-[10px] font-black px-3 py-1.5 rounded-full uppercase tracking-tighter border border-indigo-100">
                                    {{ selectedSlots.length }} choix
                                </div>
                            </div>

                            <div class="grid gap-3">
                                <div 
                                    v-for="option in activeEvent.options" 
                                    :key="option.id"
                                    @click="toggleSlot(option.id)"
                                    :class="['group p-5 rounded-[32px] border-2 transition-all flex items-center gap-5 cursor-pointer relative overflow-hidden active:scale-[0.97]', selectedSlots.includes(option.id) ? 'border-indigo-500 bg-indigo-50/40' : 'border-slate-50 bg-slate-50/50']"
                                >
                                    <div :class="['w-16 h-16 rounded-[22px] flex items-center justify-center shrink-0 transition-all', selectedSlots.includes(option.id) ? 'bg-indigo-600 text-white shadow-lg' : 'bg-white text-slate-300 shadow-sm']">
                                        <svg xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="18" height="18" x="3" y="4" rx="2" ry="2"/><line x1="16" x2="16" y1="2" y2="6"/><line x1="8" x2="8" y1="2" y2="6"/><line x1="3" x2="21" y1="10" y2="10"/></svg>
                                    </div>
                                    <div class="flex-1">
                                        <p :class="['font-black text-lg leading-none transition-colors', selectedSlots.includes(option.id) ? 'text-indigo-700' : 'text-slate-700']">{{ option.date }}</p>
                                        <p class="text-sm text-slate-400 font-bold mt-2">{{ option.time }}</p>
                                    </div>
                                    <div :class="['w-10 h-10 rounded-2xl border-2 flex items-center justify-center shrink-0 transition-all', selectedSlots.includes(option.id) ? 'bg-indigo-600 border-indigo-600 text-white scale-110' : 'border-slate-200 bg-white']">
                                        <svg v-if="selectedSlots.includes(option.id)" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>
                                        <div v-else class="w-1.5 h-1.5 bg-slate-200 rounded-full"></div>
                                    </div>
                                </div>
                            </div>

                            <div class="fixed bottom-0 left-0 right-0 p-6 bg-gradient-to-t from-white via-white to-transparent z-40">
                                <div class="max-w-md mx-auto flex gap-3">
                                    <input type="text" v-model="userName" placeholder="Votre nom" class="flex-1 bg-white border-2 border-slate-100 h-16 px-6 rounded-2xl shadow-2xl outline-none font-black text-slate-700 placeholder:text-slate-300 focus:border-indigo-100 transition-all" />
                                    <button @click="submitVote" :disabled="!userName || selectedSlots.length === 0" class="h-16 px-8 bg-indigo-600 text-white rounded-2xl font-black shadow-2xl disabled:bg-slate-200 active:scale-95 transition-all flex items-center gap-2">
                                        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>
                                    </button>
                                </div>
                            </div>
                        </div>

                        <!-- MODE : RÉSULTATS -->
                        <div v-if="eventViewMode === 'results'" class="space-y-10 pb-10">
                            <div>
                                <h2 class="font-black text-2xl text-slate-800 flex items-center gap-3">Tendances <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-amber-400"><path d="M6 9H4.5a2.5 2.5 0 0 1 0-5H6"/><path d="M18 9h1.5a2.5 2.5 0 0 0 0-5H18"/><path d="M4 22h16"/><path d="M10 14.66V17c0 .55-.47.98-.97 1.21C7.85 18.75 7 20.24 7 22"/><path d="M14 14.66V17c0 .55.47.98.97 1.21C16.15 18.75 17 20.24 17 22"/><path d="M18 2H6v7a6 6 0 0 0 12 0V2Z"/></svg></h2>
                                <p class="text-slate-400 text-sm font-medium mt-1">L'option favorite se dessine.</p>
                            </div>

                            <div class="space-y-4">
                                <div v-for="option in sortedOptions" :key="option.id" :class="['bg-white p-6 rounded-[32px] shadow-sm border transition-all', option.isWinner ? 'border-indigo-100 bg-indigo-50/20' : 'border-slate-50']">
                                    <div class="flex justify-between items-start mb-4">
                                        <div>
                                            <p :class="['font-black text-lg', option.isWinner ? 'text-indigo-600' : 'text-slate-800']">{{ option.date }}</p>
                                            <p class="text-xs text-slate-400 font-bold uppercase tracking-widest mt-1">{{ option.time }}</p>
                                        </div>
                                        <div class="text-right">
                                            <span :class="['text-2xl font-black', option.isWinner ? 'text-indigo-600' : 'text-slate-400']">{{ option.votes }}</span>
                                            <p class="text-[9px] font-black uppercase text-slate-300 tracking-tighter">VOTES</p>
                                        </div>
                                    </div>
                                    <div class="h-4 w-full bg-slate-50 rounded-full overflow-hidden border border-slate-100/50">
                                        <div :class="['h-full transition-all duration-1000 ease-out', option.isWinner ? 'bg-gradient-to-r from-indigo-500 to-violet-500 shadow-lg' : 'bg-slate-200']" :style="{ width: option.percentage + '%' }"></div>
                                    </div>
                                </div>
                            </div>
                            
                            <div class="mt-12">
                                <h3 class="font-black text-slate-800 text-lg mb-6 flex items-center gap-2">
                                    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" class="text-indigo-500"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M22 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>
                                    Participants ({{ activeEvent.responses?.length || 0 }})
                                </h3>
                                <div class="grid grid-cols-1 gap-3">
                                    <div v-for="resp in activeEvent.responses" class="flex items-center gap-4 bg-slate-50 p-4 rounded-3xl border border-slate-100/50">
                                        <img :src="resp.avatar" class="w-11 h-11 rounded-2xl bg-white shadow-sm" alt="Avatar">
                                        <div class="flex-1">
                                            <p class="font-black text-slate-800">{{ resp.name }}</p>
                                            <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest">{{ resp.votes.length }} choix</p>
                                        </div>
                                        <div class="flex -space-x-1">
                                            <div v-for="v in resp.votes" :key="v" class="w-2.5 h-2.5 rounded-full bg-indigo-400 border-2 border-slate-50"></div>
                                        </div>
                                    </div>
                                    <div v-if="!activeEvent.responses?.length" class="text-center py-10 bg-slate-50 rounded-[32px] border-2 border-dashed border-slate-100">
                                        <p class="text-slate-300 font-bold text-sm italic">Aucun vote enregistré.</p>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>

                    <!-- SHARE SHEET MODAL -->
                    <transition name="fade">
                        <div v-if="isShareSheetOpen" class="fixed inset-0 z-[100] flex items-end justify-center">
                            <div class="absolute inset-0 bg-slate-900/60 backdrop-blur-md" @click="isShareSheetOpen = false"></div>
                            <div class="relative bg-white w-full max-w-md rounded-t-[56px] p-10 shadow-2xl animate-in slide-in-from-bottom duration-500">
                                <div class="w-16 h-1.5 bg-slate-100 rounded-full mx-auto mb-10"></div>
                                <div class="flex justify-between items-center mb-10">
                                    <div>
                                        <h3 class="text-2xl font-black text-slate-900">Invitez vos amis</h3>
                                        <p class="text-slate-400 font-bold text-sm">Collectez les réponses en un clic.</p>
                                    </div>
                                    <button @click="isShareSheetOpen = false" class="w-12 h-12 bg-slate-50 rounded-2xl flex items-center justify-center text-slate-400 active:scale-90 transition-all">
                                        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M18 6 6 18"/><path d="m6 6 12 12"/></svg>
                                    </button>
                                </div>
                                <div class="bg-slate-50 p-6 rounded-[32px] border-2 border-slate-100 flex flex-col gap-4 mb-10">
                                    <div class="truncate text-xs text-slate-400 font-mono font-bold">{{ shareUrl }}</div>
                                    <button @click="copyLink" :class="['flex items-center justify-center gap-3 h-14 rounded-2xl font-black text-sm transition-all shadow-xl', copyFeedback ? 'bg-emerald-500 text-white' : 'bg-indigo-600 text-white shadow-indigo-200']">
                                        <svg v-if="copyFeedback" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>
                                        <svg v-else xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><rect width="14" height="14" x="8" y="8" rx="2" ry="2"/><path d="M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2"/></svg>
                                        {{ copyFeedback ? 'COPIÉ !' : 'COPIER LE LIEN' }}
                                    </button>
                                </div>
                                <div class="grid grid-cols-4 gap-6">
                                    <div v-for="app in shareApps" :key="app.label" @click="app.action" class="flex flex-col items-center gap-3 cursor-pointer group">
                                        <div :class="['w-16 h-16 rounded-[22px] flex items-center justify-center shadow-lg active:scale-90 transition-all text-white group-hover:-translate-y-1', app.color]">
                                            <component :is="app.svgIcon"></component>
                                        </div>
                                        <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest">{{ app.label }}</span>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </transition>
                </div>
            </template>
        </div>
    </div>

    <!-- Firebase SDK (Modular via ESM) -->
    <script type="module">
        import { initializeApp } from "https://www.gstatic.com/firebasejs/11.6.1/firebase-app.js";
        import { getFirestore, collection, addDoc, doc, onSnapshot, updateDoc, arrayUnion } from "https://www.gstatic.com/firebasejs/11.6.1/firebase-firestore.js";
        import { getAuth, signInAnonymously, signInWithCustomToken, onAuthStateChanged } from "https://www.gstatic.com/firebasejs/11.6.1/firebase-auth.js";

        const { createApp, ref, computed, onMounted, h } = Vue;

        const firebaseConfig = JSON.parse(__firebase_config);
        const firebaseApp = initializeApp(firebaseConfig);
        const db = getFirestore(firebaseApp);
        const auth = getAuth(firebaseApp);
        const appId = typeof __app_id !== 'undefined' ? __app_id : 'doodle-pro-app';

        // Icon Components
        const WhatsAppIcon = { render: () => h('svg', { xmlns: "http://www.w3.org/2000/svg", width: 28, height: 28, viewBox: "0 0 24 24", fill: "none", stroke: "currentColor", strokeWidth: 2, strokeLinecap: "round", strokeLinejoin: "round" }, [h('path', { d: "M7.9 20A9 9 0 1 0 4 16.1L2 22Z" }), h('path', { d: "M8 12h.01" }), h('path', { d: "M12 12h.01" }), h('path', { d: "M16 12h.01" })]) };
        const SMSIcon = { render: () => h('svg', { xmlns: "http://www.w3.org/2000/svg", width: 24, height: 24, viewBox: "0 0 24 24", fill: "none", stroke: "currentColor", strokeWidth: 2, strokeLinecap: "round", strokeLinejoin: "round" }, [h('path', { d: "m22 2-7 20-4-9-9-4Z" }), h('path', { d: "M22 2 11 13" })]) };
        const MailIcon = { render: () => h('svg', { xmlns: "http://www.w3.org/2000/svg", width: 24, height: 24, viewBox: "0 0 24 24", fill: "none", stroke: "currentColor", strokeWidth: 2, strokeLinecap: "round", strokeLinejoin: "round" }, [h('rect', { width: 20, height: 16, x: 2, y: 4, rx: 2 }), h('path', { d: "m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7" })]) };
        const OpenIcon = { render: () => h('svg', { xmlns: "http://www.w3.org/2000/svg", width: 24, height: 24, viewBox: "0 0 24 24", fill: "none", stroke: "currentColor", strokeWidth: 2, strokeLinecap: "round", strokeLinejoin: "round" }, [h('path', { d: "M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6" }), h('polyline', { points: "15 3 21 3 21 9" }), h('line', { x1: 10, x2: 21, y1: 14, y2: 3 })]) };

        const MOCK_EVENTS = [
            {
                id: "mock-1",
                title: "Dîner de l'Équipe",
                location: "Le Petit Bistro, Paris",
                description: "Fêtons ensemble la fin du projet !",
                organizer: "Sarah",
                createdAt: Date.now() - 1000000,
                options: [
                    { id: 1, date: "Lun 12 Fév", time: "19:30" },
                    { id: 2, date: "Mar 13 Fév", time: "20:00" },
                    { id: 3, date: "Jeu 15 Fév", time: "19:00" },
                ],
                responses: [
                    { name: "Marc", votes: [1, 3], avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=Marc" },
                    { name: "Julie", votes: [3], avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=Julie" },
                ]
            }
        ];

        createApp({
            setup() {
                const user = ref(null);
                const view = ref('home');
                const loading = ref(true);
                const events = ref([...MOCK_EVENTS]);
                const activeEvent = ref(null);
                const selectedSlots = ref([]);
                const userName = ref('');
                const eventViewMode = ref('vote');
                const isShareSheetOpen = ref(false);
                const copyFeedback = ref(false);

                const newTitle = ref('');
                const newLocation = ref('');
                const newOptions = ref([{ date: '', time: '' }]);

                onMounted(async () => {
                    const initAuth = async () => {
                        if (typeof __initial_auth_token !== 'undefined' && __initial_auth_token) {
                            await signInWithCustomToken(auth, __initial_auth_token);
                        } else {
                            await signInAnonymously(auth);
                        }
                    };
                    await initAuth();
                    onAuthStateChanged(auth, (u) => {
                        user.value = u;
                        loading.value = false;
                        if (u) loadEvents();
                    });
                });

                const loadEvents = () => {
                    const eventsRef = collection(db, 'artifacts', appId, 'public', 'data', 'events');
                    onSnapshot(eventsRef, (snapshot) => {
                        const realEvs = snapshot.docs.map(doc => ({ id: doc.id, ...doc.data() }));
                        events.value = [...MOCK_EVENTS, ...realEvs];
                    });
                };

                const openEvent = (event) => {
                    activeEvent.value = event;
                    view.value = 'event';
                    eventViewMode.value = 'vote';
                    selectedSlots.value = [];
                    window.scrollTo(0, 0);
                };

                const createEvent = async () => {
                    if (!newTitle.value) return;
                    const eventData = {
                        title: newTitle.value,
                        location: newLocation.value,
                        organizerId: user.value.uid,
                        createdAt: Date.now(),
                        options: newOptions.value.map((o, idx) => ({ ...o, id: idx + 1 })),
                        responses: []
                    };
                    await addDoc(collection(db, 'artifacts', appId, 'public', 'data', 'events'), eventData);
                    view.value = 'home';
                    newTitle.value = '';
                    newLocation.value = '';
                    newOptions.value = [{ date: '', time: '' }];
                };

                const toggleSlot = (id) => {
                    if (selectedSlots.value.includes(id)) {
                        selectedSlots.value = selectedSlots.value.filter(s => s !== id);
                    } else {
                        selectedSlots.value.push(id);
                    }
                };

                const submitVote = async () => {
                    if (!userName.value || selectedSlots.value.length === 0) return;
                    const response = {
                        name: userName.value,
                        uid: user.value.uid,
                        votes: selectedSlots.value,
                        avatar: `https://api.dicebear.com/7.x/avataaars/svg?seed=${userName.value}`
                    };

                    if (activeEvent.value.id.toString().startsWith('mock-')) {
                        activeEvent.value.responses.push(response);
                        eventViewMode.value = 'results';
                        return;
                    }

                    const eventRef = doc(db, 'artifacts', appId, 'public', 'data', 'events', activeEvent.value.id);
                    await updateDoc(eventRef, { responses: arrayUnion(response) });
                    eventViewMode.value = 'results';
                };

                const shareUrl = computed(() => `https://doodle.pro/event/${activeEvent.value?.id}`);
                
                const copyLink = () => {
                    const text = `Disponibilités pour "${activeEvent.value?.title}" : ${shareUrl.value}`;
                    const el = document.createElement('textarea');
                    el.value = text;
                    document.body.appendChild(el);
                    el.select();
                    document.execCommand('copy');
                    document.body.removeChild(el);
                    copyFeedback.value = true;
                    setTimeout(() => copyFeedback.value = false, 2000);
                };

                const sortedOptions = computed(() => {
                    if (!activeEvent.value) return [];
                    const counts = {};
                    activeEvent.value.options.forEach(o => counts[o.id] = 0);
                    activeEvent.value.responses?.forEach(r => r.votes.forEach(v => counts[v]++));
                    const total = activeEvent.value.responses?.length || 1;
                    const max = Math.max(...Object.values(counts));
                    return activeEvent.value.options.map(o => ({
                        ...o,
                        votes: counts[o.id],
                        percentage: (counts[o.id] / total) * 100,
                        isWinner: counts[o.id] > 0 && counts[o.id] === max
                    }));
                });

                const shareApps = [
                    { label: 'WhatsApp', color: 'bg-emerald-500', svgIcon: WhatsAppIcon, action: () => window.open(`https://wa.me/?text=${encodeURIComponent(shareUrl.value)}`) },
                    { label: 'SMS', color: 'bg-sky-500', svgIcon: SMSIcon, action: () => window.location.href = `sms:?body=${encodeURIComponent(shareUrl.value)}` },
                    { label: 'Email', color: 'bg-slate-800', svgIcon: MailIcon, action: () => window.location.href = `mailto:?subject=Sondage&body=${encodeURIComponent(shareUrl.value)}` },
                    { label: 'Ouvrir', color: 'bg-indigo-600', svgIcon: OpenIcon, action: () => window.open(shareUrl.value) }
                ];

                return {
                    view, loading, events, activeEvent, selectedSlots, userName, eventViewMode,
                    isShareSheetOpen, copyFeedback, newTitle, newLocation, newOptions,
                    openEvent, createEvent, toggleSlot, submitVote, copyLink, shareUrl,
                    sortedOptions, shareApps
                };
            }
        }).mount('#app');
    </script>
</body>
</html>