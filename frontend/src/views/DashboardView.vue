<template>
  <div class="min-h-screen pb-10 px-4 relative overflow-hidden">
    <!-- TRON BACKGROUND LINES -->
    <div class="fixed inset-0 pointer-events-none -z-10 bg-grid-cyan opacity-10"></div>
    
    <!-- TOAST NOTIFICATIONS -->
    <div class="fixed top-4 right-4 z-[100] flex flex-col gap-3">
      <transition-group name="toast">
        <div v-for="toast in toasts" :key="toast.id" 
             :class="[
               'px-6 py-4 rounded-xl border shadow-2xl backdrop-blur-md flex items-center gap-3 w-80',
               toast.type === 'error' ? 'bg-red-900/40 border-red-500/50 text-red-100 shadow-red-500/20' : 
               toast.type === 'success' ? 'bg-green-900/40 border-green-500/50 text-green-100 shadow-green-500/20' :
               'bg-yellow-900/40 border-yellow-500/50 text-yellow-100 shadow-yellow-500/20'
             ]">
          <div v-html="toast.icon" class="flex-shrink-0"></div>
          <div class="flex-1 font-medium text-sm leading-tight">{{ toast.message }}</div>
          <button @click="removeToast(toast.id)" class="opacity-50 hover:opacity-100">×</button>
        </div>
      </transition-group>
    </div>

    <!-- MAIN CONTAINER -->
    <div class="container mx-auto max-w-6xl relative z-10">
      <!-- Logo Header removed - moved to Navbar -->

      <!-- Hero Section -->
      <div class="mb-12 text-center">
        <h2 class="text-5xl md:text-6xl font-['Orbitron'] font-black mb-4 gradient-text">
          SCRIM MATCHMAKING
        </h2>
        <p class="text-xl text-gray-400">
          Find worthy opponents. Prove your dominance.
        </p>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-4 gap-6">
        <!-- Main Area -->
        <div class="lg:col-span-3">
          
          <!-- STEP 1: Category Selection -->
          <div v-if="!selectedCategory" class="space-y-6">
            <h3 class="text-2xl font-bold text-center mb-8">Choose Your Match Type</h3>
            
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <!-- POKE Card -->
              <div class="glass-card p-8 transition-all duration-300">
                <div class="text-center space-y-4">
                  <!-- Icon -->
                  <div class="w-20 h-20 mx-auto rounded-full bg-gradient-to-br from-electric-violet-500 to-electric-violet-700 flex items-center justify-center shadow-glow-violet transition-transform">
                    <svg class="w-10 h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                  </div>
                  
                  <!-- Title -->
                  <h4 class="text-2xl font-bold text-glow">POKE</h4>
                  <p class="text-electric-violet-300 font-semibold">Ranked Solo/Duo Match</p>
                  
                  <!-- Description -->
                  <p class="text-sm text-gray-400 leading-relaxed">
                    Match based on your solo rank with balanced matchmaking (±1 rank tolerance)
                  </p>
                  
                  <!-- Features -->
                  <div class="space-y-2 text-left text-sm">
                    <div class="flex items-center gap-2 text-gray-300">
                      <svg class="w-4 h-4 text-green-400" fill="currentColor" viewBox="0 0 20 20">
                        <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                      </svg>
                      Rank-based matchmaking
                    </div>
                    <div class="flex items-center gap-2 text-gray-300">
                      <svg class="w-4 h-4 text-green-400" fill="currentColor" viewBox="0 0 20 20">
                        <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                      </svg>
                      Fair & balanced teams
                    </div>
                    <div class="flex items-center gap-2 text-gray-300">
                      <svg class="w-4 h-4 text-green-400" fill="currentColor" viewBox="0 0 20 20">
                        <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                      </svg>
                      Warrior to Classic ranks
                    </div>
                  </div>
                  
                  <!-- Button -->
                  <button @click="selectCategory('POKE')" class="w-full btn-magic py-3 text-sm mt-4">
                    SELECT POKE
                  </button>
                </div>
              </div>

              <!-- WARKOP Card -->
              <div class="glass-card p-8 transition-all duration-300">
                <div class="text-center space-y-4">
                  <!-- Icon -->
                  <div class="w-20 h-20 mx-auto rounded-full bg-gradient-gold flex items-center justify-center shadow-glow-gold transition-transform">
                    <svg class="w-10 h-10 text-midnight-900" fill="currentColor" viewBox="0 0 20 20">
                      <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
                    </svg>
                  </div>
                  
                  <!-- Title -->
                  <h4 class="text-2xl font-bold text-glow-gold">WARKOP</h4>
                  <p class="text-antique-gold-300 font-semibold">Pro Scrim Team Match</p>
                  
                  <!-- Description -->
                  <p class="text-sm text-gray-400 leading-relaxed">
                    Professional team scrimmage with instant match for competitive teams
                  </p>
                  
                  <!-- Features -->
                  <div class="space-y-2 text-left text-sm">
                    <div class="flex items-center gap-2 text-gray-300">
                      <svg class="w-4 h-4 text-antique-gold-400" fill="currentColor" viewBox="0 0 20 20">
                        <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                      </svg>
                      Instant matching
                    </div>
                    <div class="flex items-center gap-2 text-gray-300">
                      <svg class="w-4 h-4 text-antique-gold-400" fill="currentColor" viewBox="0 0 20 20">
                        <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                      </svg>
                      Pro team coordination
                    </div>
                    <div class="flex items-center gap-2 text-gray-300">
                      <svg class="w-4 h-4 text-antique-gold-400" fill="currentColor" viewBox="0 0 20 20">
                        <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                      </svg>
                      Competitive environment
                    </div>
                  </div>
                  
                  <!-- Button -->
                  <button @click="selectCategory('WARKOP')" class="w-full btn-gold py-3 text-sm mt-4">
                    SELECT WARKOP
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- STEP 2a: POKE Form -->
          <div v-else-if="selectedCategory === 'POKE'" class="glass-card p-8 md:p-12 relative">
            <div v-if="loadingApi" class="absolute inset-0 flex items-center justify-center bg-black/50 z-20 rounded-2xl">
              <div class="w-16 h-16 border-4 border-electric-violet-500 border-t-cyan-magic-400 rounded-full animate-spin"></div>
            </div>

            <div v-if="!isSearching">
              <!-- Back Button -->
              <button @click="selectedCategory = null" class="mb-6 text-sm text-gray-400 hover:text-cyan-magic-300 transition-colors flex items-center gap-2">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
                </svg>
                Back to selection
              </button>

              <div class="space-y-8">
                <div class="text-center">
                  <div class="inline-flex items-center gap-3 bg-electric-violet-500/10 border border-electric-violet-500/30 rounded-full px-6 py-2 mb-4">
                    <svg class="w-5 h-5 text-electric-violet-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    <span class="text-electric-violet-300 font-bold">POKE - Ranked Solo Match</span>
                  </div>
                  <h3 class="text-2xl font-bold text-glow mb-2">POKE - Ranked Solo Match</h3>
                  <p class="text-sm text-gray-400">Fill in your details for ranked matchmaking</p>
                </div>

                <div class="space-y-6 max-w-xl mx-auto">
                  <!-- Team Name -->
                  <div>
                    <label class="block text-sm text-gray-400 mb-2 font-medium">Your Name</label>
                    <input 
                      type="text" 
                      v-model="formData.teamName" 
                      placeholder="e.g. Solo Warriors" 
                      class="w-full bg-midnight-800 border-2 border-electric-violet-500/30 rounded-lg p-4 text-white text-lg focus:outline-none focus:border-cyan-magic-400 transition-all"
                      @keyup.enter="startSearch"
                    />
                  </div>

                  <!-- WhatsApp -->
                  <div>
                    <label class="block text-sm text-gray-400 mb-2 font-medium">WhatsApp Number</label>
                    <input 
                      type="tel" 
                      v-model="formData.whatsappNumber" 
                      placeholder="e.g. 628123456789" 
                      :class="[
                        'w-full bg-midnight-800 border-2 rounded-lg p-4 text-white text-lg focus:outline-none transition-all',
                        whatsappError ? 'border-red-500' : 'border-electric-violet-500/30 focus:border-cyan-magic-400'
                      ]"
                      @keyup.enter="startSearch"
                    />
                    <p v-if="whatsappError" class="text-xs text-red-400 mt-1">{{ whatsappError }}</p>
                    <p v-else class="text-xs text-gray-500 mt-1">Format: 628xxxxxxxxx (without +)</p>
                  </div>

                  <!-- Rank -->
                  <div>
                    <label class="block text-sm text-gray-400 mb-2 font-medium">Your Solo Rank</label>
                    <select 
                      v-model="formData.rankWeight" 
                      @change="updateRankName" 
                      class="w-full bg-midnight-800 border-2 border-electric-violet-500/30 rounded-lg p-4 text-white text-lg focus:outline-none focus:border-cyan-magic-400 transition-all cursor-pointer"
                    >
                      <option value="1">1 - Warrior</option>
                      <option value="2">2 - Elite</option>
                      <option value="3">3 - Master</option>
                      <option value="4">4 - Grandmaster</option>
                      <option value="5">5 - Epic</option>
                      <option value="6">6 - Legend</option>
                      <option value="7">7 - Mythic</option>
                      <option value="8">8 - Mythical Glory</option>
                      <option value="9">9 - Classic/Fun</option>
                    </select>
                    <p class="text-xs text-gray-500 mt-1">Matchmaking tolerance: ±1 ranks</p>
                  </div>
                </div>

                <!-- Error -->
                <div v-if="errorApi" class="text-center">
                  <div class="inline-block text-red-400 bg-red-900/20 p-3 px-6 rounded-lg border border-red-500/50">
                    ⚠️ {{ errorApi }}
                  </div>
                </div>

                <!-- Portal -->
                <div class="relative w-40 h-40 mx-auto cursor-pointer group" @click="startSearch">
                  <div class="absolute inset-0 rounded-full border-4 border-electric-violet-500/30 animate-spin-slow group-hover:border-cyan-magic-400/50 transition-all"></div>
                  <div class="absolute inset-4 rounded-full bg-gradient-magic opacity-20 animate-pulse-slow group-hover:opacity-40 transition-all"></div>
                  <div class="absolute inset-0 flex items-center justify-center">
                    <svg class="w-16 h-16 text-cyan-magic-300 group-hover:scale-110 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0zM10 7v3m0 0v3m0-3h3m-3 0H7" />
                    </svg>
                  </div>
                </div>

                <!-- CTA -->
                <div class="text-center">
                  <button 
                    @click="startSearch"
                    :disabled="!isValidFormPoke"
                    :class="{'opacity-50 cursor-not-allowed': !isValidFormPoke, 'portal-pulse': isValidFormPoke}"
                    class="btn-magic text-2xl px-16 py-6">
                    <span class="relative z-10 flex items-center gap-3">
                      <svg class="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
                      </svg>
                      FIND RANKED MATCH
                    </span>
                  </button>
                  <p class="text-sm text-gray-500 mt-4">Balanced matchmaking based on rank</p>
                </div>
              </div>
            </div>
            <SearchingState v-else :formData="formData" @cancel="cancelSearch" />
          </div>

          <!-- STEP 2b: WARKOP Form -->
          <div v-else-if="selectedCategory === 'WARKOP'" class="glass-card p-8 md:p-12 relative">
            <div v-if="loadingApi" class="absolute inset-0 flex items-center justify-center bg-black/50 z-20 rounded-2xl">
              <div class="w-16 h-16 border-4 border-antique-gold-400 border-t-cyan-magic-400 rounded-full animate-spin"></div>
            </div>

            <div v-if="!isSearching">
              <!-- Back Button -->
              <button @click="selectedCategory = null" class="mb-6 text-sm text-gray-400 hover:text-antique-gold-300 transition-colors flex items-center gap-2">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
                </svg>
                Back to selection
              </button>

              <div class="space-y-8">
                <div class="text-center">
                  <div class="inline-flex items-center gap-3 bg-antique-gold-400/10 border border-antique-gold-400/30 rounded-full px-6 py-2 mb-4">
                    <svg class="w-5 h-5 text-antique-gold-400" fill="currentColor" viewBox="0 0 20 20">
                      <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
                    </svg>
                    <span class="text-antique-gold-300 font-bold">WARKOP - Pro Scrim</span>
                  </div>
                  <h3 class="text-2xl font-bold text-glow-gold mb-2">Team Registration</h3>
                  <p class="text-sm text-gray-400">Professional team scrim details</p>
                </div>

                <div class="space-y-6 max-w-xl mx-auto">
                  <!-- Team Name -->
                  <div>
                    <label class="block text-sm text-gray-400 mb-2 font-medium">Team Name</label>
                    <input 
                      type="text" 
                      v-model="formData.teamName" 
                      placeholder="e.g. RRQ Hoshi" 
                      class="w-full bg-midnight-800 border-2 border-antique-gold-400/30 rounded-lg p-4 text-white text-lg focus:outline-none focus:border-antique-gold-400 transition-all"
                      @keyup.enter="startSearch"
                    />
                  </div>

                  <!-- Captain Name -->
                  <div>
                    <label class="block text-sm text-gray-400 mb-2 font-medium">Captain Name</label>
                    <input 
                      type="text" 
                      v-model="formData.captainName" 
                      placeholder="e.g. Lemon" 
                      class="w-full bg-midnight-800 border-2 border-antique-gold-400/30 rounded-lg p-4 text-white text-lg focus:outline-none focus:border-antique-gold-400 transition-all"
                      @keyup.enter="startSearch"
                    />
                  </div>

                  <!-- WhatsApp -->
                  <div>
                    <label class="block text-sm text-gray-400 mb-2 font-medium">WhatsApp Number (Captain)</label>
                    <input 
                      type="tel" 
                      v-model="formData.whatsappNumber" 
                      placeholder="e.g. 628123456789" 
                      :class="[
                        'w-full bg-midnight-800 border-2 rounded-lg p-4 text-white text-lg focus:outline-none transition-all',
                        whatsappError ? 'border-red-500' : 'border-antique-gold-400/30 focus:border-antique-gold-400'
                      ]"
                      @keyup.enter="startSearch"
                    />
                    <p v-if="whatsappError" class="text-xs text-red-400 mt-1">{{ whatsappError }}</p>
                    <p v-else class="text-xs text-gray-500 mt-1">Format: 628xxxxxxxxx (without +)</p>
                  </div>

                  <!-- Pro Info -->
                  <div class="bg-antique-gold-400/10 border border-antique-gold-400/30 rounded-lg p-4">
                    <div class="flex items-start gap-3">
                      <svg class="w-6 h-6 text-antique-gold-400 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                        <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
                      </svg>
                      <div>
                        <p class="text-antique-gold-400 font-bold text-sm">Pro Scrim Mode</p>
                        <p class="text-xs text-gray-400 mt-1">Instant match with other professional teams. No rank filtering applied.</p>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- Error -->
                <div v-if="errorApi" class="text-center">
                  <div class="inline-block text-red-400 bg-red-900/20 p-3 px-6 rounded-lg border border-red-500/50">
                    ⚠️ {{ errorApi }}
                  </div>
                </div>

                <!-- Portal -->
                <div class="relative w-40 h-40 mx-auto cursor-pointer group" @click="startSearch">
                  <div class="absolute inset-0 rounded-full border-4 border-antique-gold-400/30 animate-spin-slow group-hover:border-antique-gold-400/70 transition-all"></div>
                  <div class="absolute inset-4 rounded-full bg-gradient-gold opacity-20 animate-pulse-slow group-hover:opacity-40 transition-all"></div>
                  <div class="absolute inset-0 flex items-center justify-center">
                    <svg class="w-16 h-16 text-antique-gold-400 group-hover:scale-110 transition-transform" fill="currentColor" viewBox="0 0 20 20">
                      <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
                    </svg>
                  </div>
                </div>

                <!-- CTA -->
                <div class="text-center">
                  <button 
                    @click="startSearch"
                    :disabled="!isValidFormWarkop"
                    :class="{'opacity-50 cursor-not-allowed': !isValidFormWarkop, 'portal-pulse': isValidFormWarkop}"
                    class="btn-gold text-2xl px-16 py-6">
                    <span class="relative z-10 flex items-center gap-3">
                      <svg class="w-7 h-7" fill="currentColor" viewBox="0 0 20 20">
                        <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
                      </svg>
                      FIND PRO SCRIM
                    </span>
                  </button>
                  <p class="text-sm text-gray-500 mt-4">Instant match with pro teams</p>
                </div>
              </div>
            </div>
            <SearchingState v-else :formData="formData" @cancel="cancelSearch" />
          </div>

        </div>


        <!-- Live Server Status Sidebar -->
        <div class="lg:col-span-1">
          <div class="glass-card p-6 sticky top-6 h-fit border border-white/5 backdrop-blur-2xl bg-midnight-900/30">
            <!-- Header with Live Indicator -->
            <div class="flex items-center justify-between mb-6">
              <h4 class="text-lg font-bold flex items-center gap-3">
                <span class="bg-gradient-to-r from-cyan-magic-300 via-electric-violet-300 to-antique-gold-300 bg-clip-text text-transparent tracking-widest">
                  LIVE STATUS
                </span>
              </h4>
              <div class="flex items-center gap-2">
                <span class="relative flex h-3 w-3">
                  <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
                  <span class="relative inline-flex rounded-full h-3 w-3 bg-green-500"></span>
                </span>
                <span class="text-xs text-green-400 font-mono">ONLINE</span>
              </div>
            </div>
            
            <div class="space-y-5">
              <!-- Average Wait Time (Highlighted with Subtle Animated Border) -->
              <div class="relative group">
                <!-- Glowing effect behind -->
                <div class="absolute -inset-0.5 bg-gradient-to-r from-cyan-magic-500/20 to-blue-500/20 rounded-xl blur-sm opacity-50 group-hover:opacity-100 transition duration-1000 group-hover:duration-200"></div>
                
                <div class="relative glass bg-black/40 p-4 rounded-xl flex items-center gap-4 border border-cyan-magic-500/30">
                  <div class="p-3 rounded-full bg-cyan-magic-500/10 shadow-[0_0_15px_rgba(6,182,212,0.2)]">
                    <!-- Stopwatch Icon -->
                    <svg class="w-6 h-6 text-cyan-magic-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                  </div>
                  <div>
                    <div class="text-2xl font-black text-cyan-magic-300 font-['Orbitron'] tracking-wider transition-all duration-300">
                      < {{ liveStats.waitTime }}s
                    </div>
                    <div class="text-[10px] text-gray-400 font-bold uppercase tracking-widest">Avg Wait Time</div>
                  </div>
                </div>
              </div>

              <!-- Rank Tolerance -->
              <div class="glass bg-black/20 p-4 rounded-xl flex items-center gap-4 border border-electric-violet-500/10 hover:border-electric-violet-500/30 transition-all duration-500 hover:shadow-[0_0_20px_rgba(139,92,246,0.1)]">
                <div class="p-3 rounded-full bg-electric-violet-500/10 shadow-[0_0_10px_rgba(139,92,246,0.15)] pulse-ring-violet">
                  <!-- Medal/Rank Icon -->
                  <svg class="w-6 h-6 text-electric-violet-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </div>
                <div>
                  <div class="text-2xl font-black text-electric-violet-400 animate-pulse font-['Orbitron'] tracking-wider">±1</div>
                  <div class="text-[10px] text-gray-400 font-bold uppercase tracking-widest">Rank Tolerance</div>
                </div>
              </div>

              <!-- Confirmation Time -->
              <div class="glass bg-black/20 p-4 rounded-xl flex items-center gap-4 border border-antique-gold-400/10 hover:border-antique-gold-400/30 transition-all duration-500 hover:shadow-[0_0_20px_rgba(212,175,55,0.1)]">
                <div class="p-3 rounded-full bg-antique-gold-400/10 shadow-[0_0_10px_rgba(212,175,55,0.15)]">
                  <!-- Bell Alert Icon -->
                  <svg class="w-6 h-6 text-antique-gold-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
                  </svg>
                </div>
                <div>
                  <div class="text-2xl font-black text-antique-gold-400 font-['Orbitron'] tracking-wider transition-all duration-300">{{ liveStats.confirm }}s</div>
                  <div class="text-[10px] text-gray-400 font-bold uppercase tracking-widest">Confirm Time</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Recent Matches Section (Simplified - No Results) -->
      <div class="glass-card p-6 mt-8">
        <h4 class="text-xl font-bold mb-6 flex items-center gap-2">
          <svg class="w-5 h-5 text-cyan-magic-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          Recent Pro Scrims
        </h4>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <div v-for="(match, index) in recentMatches" 
            :key="index"
            class="glass p-4 rounded-lg hover:shadow-glow-violet transition-all">
            <div class="flex items-start gap-3">
              <div class="w-3 h-3 rounded-full mt-1 flex-shrink-0 bg-cyan-magic-400"></div>
              <div class="flex-1 min-w-0">
                <p class="font-medium truncate text-white">{{ match.opponent }}</p>
                <p class="text-xs text-gray-400">{{ match.category }} • {{ match.timeAgo }}</p>
              </div>
            </div>
          </div>
          
          <div v-if="recentMatches.length === 0" class="md:col-span-2 lg:col-span-3 text-center py-12 text-gray-500">
            <svg class="w-16 h-16 mx-auto mb-3 opacity-30" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
            </svg>
            <p class="text-sm">No pro scrim matches yet</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Match Found Modal -->
    <MatchFoundModal 
      :isVisible="showMatchModal"
      :yourTeam="matchData.yourTeam"
      :opponentTeam="matchData.opponentTeam"
      :matchDetails="matchData.details"
      :timeout="matchData.timeout"
      @accept="handleMatchAccept"
      @decline="handleMatchDecline"
      @timeout="handleMatchDecline"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useScrimAPI } from '@/composables/useScrimAPI'
import MatchFoundModal from '@/components/MatchFoundModal.vue'
import SearchingState from '@/components/SearchingState.vue'



const { findMatch, checkStatus, cancelMatch, loading: loadingApi, error: errorApi } = useScrimAPI()

// Toasts System
const toasts = ref([])
const addToast = (message, type = 'error') => {
  const id = Date.now() + Math.random()
  const icons = {
    error: `<svg class="w-6 h-6 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/></svg>`,
    success: `<svg class="w-6 h-6 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>`,
    warning: `<svg class="w-6 h-6 text-yellow-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>`
  }
  toasts.value.push({ id, message, type, icon: icons[type] })
  setTimeout(() => removeToast(id), 5000)
}
const removeToast = (id) => {
  toasts.value = toasts.value.filter(t => t.id !== id)
}

// State
const selectedCategory = ref(null)
const isSearching = ref(false)
const showMatchModal = ref(false)
const currentRequestId = ref(null)
let pollInterval = null

// Form Data
const formData = ref({
  teamName: '',
  captainName: '',
  whatsappNumber: '',
  category: 'POKE',
  rankWeight: '5',
  rankName: 'Epic'
})

// WhatsApp Validation
const isValidWhatsApp = computed(() => {
  const wa = formData.value.whatsappNumber.trim()
  // Must start with 628 and have at least 10 digits total
  return wa.startsWith('628') && wa.length >= 10 && wa.length <= 15 && /^[0-9]+$/.test(wa)
})

const whatsappError = computed(() => {
  const wa = formData.value.whatsappNumber.trim()
  if (!wa) return ''
  if (!wa.startsWith('628')) return '⚠️ Must start with 628'
  if (wa.length < 10) return '⚠️ Too short (min 10 digits)'
  if (wa.length > 15) return '⚠️ Too long (max 15 digits)'
  if (!/^[0-9]+$/.test(wa)) return '⚠️ Numbers only'
  return ''
})

// Validation
const isValidFormPoke = computed(() => {
  return formData.value.teamName.trim().length >= 3 && 
         isValidWhatsApp.value
})

const isValidFormWarkop = computed(() => {
  return formData.value.teamName.trim().length >= 3 && 
         formData.value.captainName.trim().length >= 3 &&
         isValidWhatsApp.value
})

const selectCategory = (category) => {
  selectedCategory.value = category
  formData.value.category = category
  
  if (category === 'WARKOP') {
    formData.value.rankWeight = '10'
    formData.value.rankName = 'Pro Scrim'
  } else {
    formData.value.rankWeight = '5'
    updateRankName()
  }
}

const updateRankName = () => {
  const ranks = ['Warrior', 'Elite', 'Master', 'Grandmaster', 'Epic', 'Legend', 'Mythic', 'Mythical Glory', 'Classic/Fun', 'Pro Scrim']
  formData.value.rankName = ranks[parseInt(formData.value.rankWeight) - 1] || 'Unknown'
}

// Match Data
const matchData = ref({
  yourTeam: {},
  opponentTeam: {},
  details: {},
  timeout: 60
})

// Methods
const startSearch = async () => {
  const isPoke = selectedCategory.value === 'POKE'
  const isValid = isPoke ? isValidFormPoke.value : isValidFormWarkop.value
  
  if (!isValid) {
    if (formData.value.teamName.trim().length < 3) addToast("Team Name is too short (min 3 chars).", 'warning')
    else if (!isPoke && formData.value.captainName.trim().length < 3) addToast("Captain Name is too short (min 3 chars).", 'warning')
    else if (whatsappError.value) addToast(`WhatsApp Error: ${whatsappError.value}`, 'warning')
    else if (!isValidWhatsApp.value) addToast("Please enter a valid WhatsApp number.", 'warning')
    else addToast("Please fill in all required fields properly.", 'warning')
    return
  }

  isSearching.value = true // SHOW WAITING UI IMMEDIATELY
  
  try {
    // Pass WebSocket match-found callback to findMatch
    const result = await findMatch(formData.value, (wsData) => {
      // Called instantly when WebSocket delivers match
      isSearching.value = false
      clearInterval(pollInterval)
      const match = wsData
      matchData.value = {
        yourTeam: {
          name: formData.value.teamName,
          avatar: `https://api.dicebear.com/7.x/identicon/svg?seed=${formData.value.teamName}`,
          rank: formData.value.rankName
        },
        opponentTeam: {
          name: match.opponent_name || 'Opponent Team',
          avatar: `https://api.dicebear.com/7.x/identicon/svg?seed=${match.opponent_name}`,
          rank: formData.value.category === 'WARKOP' ? 'Pro Scrim' : 'Similar Rank'
        },
        details: {
          category: formData.value.category,
          rankDiff: 0,
          id: (match.match_id || 'MATCH').substring(0, 6).toUpperCase(),
          whatsapp_url: match.whatsapp_url
        },
        timeout: match.expires_in || 60
      }
      showMatchModal.value = true
    })
    
    if (result && result.request_id) {
       currentRequestId.value = result.request_id
       startPolling(result.request_id)
       addToast("Matchmaking started! Connecting to network...", 'success')
    } else {
       // Request accepted by server but structure weird? fallback
       addToast("Matchmaking starting...", 'success')
    }
  } catch (err) {
    isSearching.value = false // Revert if API blows up
    let errMsg = "Failed to start match search."
    if (errorApi.value) errMsg += ` Reason: ${errorApi.value}`
    addToast(errMsg, 'error')
    console.error("Search failed", err)
  }
}

const startPolling = (requestId) => {
  if (pollInterval) clearInterval(pollInterval)
  
  pollInterval = setInterval(async () => {
    try {
      const status = await checkStatus(requestId)
      
      if (status.status === 'matched') {
         clearInterval(pollInterval)
         isSearching.value = false
         
         const match = status.match
         matchData.value = {
            yourTeam: {
               name: formData.value.teamName,
               avatar: `https://api.dicebear.com/7.x/identicon/svg?seed=${formData.value.teamName}`,
               rank: formData.value.rankName
            },
            opponentTeam: {
               name: match.opponent_name || 'Opponent Team',
               avatar: `https://api.dicebear.com/7.x/identicon/svg?seed=${match.opponent_name}`,
               rank: formData.value.category === 'WARKOP' ? 'Pro Scrim' : 'Similar Rank'
            },
            details: {
               category: status.category,
               rankDiff: 0,
               id: match.match_id.substring(0, 6).toUpperCase(),
               whatsapp_url: match.whatsapp_url
            },
            timeout: match.expires_in || 60
         }
         
         showMatchModal.value = true
      }
      
      if (status.status === 'expired') {
          clearInterval(pollInterval)
          isSearching.value = false
          addToast("Request expired. Time limit reached.", 'warning')
      }

    } catch (err) {
      console.error("Polling error", err)
    }
  }, 2000)
}

const cancelSearch = async () => {
  if (currentRequestId.value) {
    await cancelMatch(currentRequestId.value)
  }
  clearInterval(pollInterval)
  isSearching.value = false
  currentRequestId.value = null
}

const handleMatchAccept = () => {
  if (matchData.value.details.whatsapp_url) {
     window.open(matchData.value.details.whatsapp_url, '_blank')
  }
  showMatchModal.value = false
  
  // Only add WARKOP (Pro Scrim) to history, not POKE
  if (formData.value.category === 'WARKOP') {
    recentMatches.value.unshift({
      opponent: matchData.value.opponentTeam.name,
      category: formData.value.category,
      timeAgo: 'Just now'
    })
    
    if (recentMatches.value.length > 5) {
      recentMatches.value.pop()
    }
  }
  
  // Reset
  selectedCategory.value = null
  formData.value.teamName = ''
  formData.value.captainName = ''
  formData.value.whatsappNumber = ''
  formData.value.category = 'POKE'
  formData.value.rankWeight = '5'
  updateRankName()
}

// Live Status state
const liveStats = ref({
  waitTime: 30,
  tolerance: 1,
  confirm: 60
})
let liveStatusInterval = null

const handleMatchDecline = () => {
  showMatchModal.value = false
  addToast("Match declined/expired. You have been placed back to selection.", 'warning')
}

onMounted(() => {
  // Live Status dynamic animation
  liveStatusInterval = setInterval(() => {
    // fluctuate wait time between 24s and 34s
    liveStats.value.waitTime = Math.floor(24 + Math.random() * 10)
    // occasionally blip confirm time
    if (Math.random() > 0.8) {
      liveStats.value.confirm = 59
      setTimeout(() => liveStats.value.confirm = 60, 1000)
    }
  }, 3500)
})

onUnmounted(() => {
  if (pollInterval) clearInterval(pollInterval)
  if (liveStatusInterval) clearInterval(liveStatusInterval)
})

// Recent Matches
const recentMatches = ref([
  { opponent: 'RRQ Hoshi', category: 'WARKOP', result: 'won', timeAgo: '2 hours ago', duration: '15:34' },
  { opponent: 'EVOS Legends', category: 'POKE', result: 'won', timeAgo: '5 hours ago', duration: '18:21' },
  { opponent: 'Alter Ego', category: 'POKE', result: 'lost', timeAgo: 'Yesterday', duration: '21:12' }
])

updateRankName()
</script>

<style scoped>
.toast-enter-active, .toast-leave-active {
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}
.toast-enter-from {
  opacity: 0;
  transform: translateX(50px) scale(0.9);
}
.toast-leave-to {
  opacity: 0;
  transform: translateY(-20px) scale(0.9);
}
</style>
