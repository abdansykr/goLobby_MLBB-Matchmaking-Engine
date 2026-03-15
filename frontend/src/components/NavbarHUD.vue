<template>
  <nav class="fixed top-0 left-0 right-0 z-50 glass-strong border-b border-white/10">
    <div class="container mx-auto px-4 py-3">
      <div class="flex items-center justify-between">
        <!-- Logo Section with Crystal Effect -->
        <div class="flex items-center gap-3 group">
          <div class="relative">
            <!-- Hexagon Logo Container -->
            <div class="w-12 h-12 hexagon-border bg-gradient-magic flex items-center justify-center shadow-glow-violet">
              <svg class="w-7 h-7 text-white" fill="currentColor" viewBox="0 0 20 20">
                <path d="M10 2L2 7v6l8 5 8-5V7l-8-5zM10 13l-5-3 5-3 5 3-5 3z"/>
              </svg>
            </div>
            <!-- Animated Glow Ring -->
            <div class="absolute inset-0 hexagon-border bg-electric-violet-500/30 animate-ping"></div>
          </div>
          
          <div>
            <h1 class="font-['Orbitron'] text-2xl font-black gradient-text">
              GOLOBBY
            </h1>
            <p class="text-xs text-cyan-magic-300 font-medium tracking-wider">
              SCRIM ARENA
            </p>
          </div>
        </div>

        <!-- Center Stats (Optional - shows when logged in) -->
        <div v-if="userStats" class="hidden md:flex items-center gap-8">
          <div class="text-center">
            <p class="text-xs text-gray-400 uppercase tracking-wide">Rank</p>
            <div class="flex items-center gap-2 mt-1">
              <div class="w-6 h-6 rounded-full bg-legend-purple animate-pulse"></div>
              <span class="font-bold text-lg">{{ userStats.rank }}</span>
            </div>
          </div>
          
          <div class="text-center">
            <p class="text-xs text-gray-400 uppercase tracking-wide">Skor Reputasi</p>
            <div class="flex items-center gap-2 mt-1">
              <span class="font-bold text-lg" :class="reputationColor">
                {{ userStats.reputation }}
              </span>
            </div>
          </div>
        </div>

        <!-- Action & Identity Section -->
        <div class="flex items-center gap-3 sm:gap-4">
          
          <!-- Discord Icon -->
          <a href="#" target="_blank" class="relative p-2 rounded-lg glass hover:shadow-glow-cyan transition-all" title="Join Discord">
            <svg class="w-5 h-5 text-[#5865F2]" fill="currentColor" viewBox="0 0 127.14 96.36">
              <path d="M107.7,8.07A105.15,105.15,0,0,0,81.47,0a72.06,72.06,0,0,0-3.36,6.83A97.68,97.68,0,0,0,49,6.83,72.37,72.37,0,0,0,45.64,0,105.89,105.89,0,0,0,19.39,8.09C2.79,32.65-1.71,56.6.54,80.21h0A105.73,105.73,0,0,0,32.71,96.36,77.7,77.7,0,0,0,39.6,85.25a68.42,68.42,0,0,1-10.85-5.18c.91-.66,1.8-1.34,2.66-2a75.57,75.57,0,0,0,64.32,0c.87.71,1.76,1.39,2.66,2a68.68,68.68,0,0,1-10.87,5.19,77.7,77.7,0,0,0,6.89,11.1A105.25,105.25,0,0,0,126.6,80.22h0C129.24,52.84,122.09,29.11,107.7,8.07ZM42.45,65.69C36.18,65.69,31,60,31,53s5-12.74,11.43-12.74S54,46,53.89,53,48.84,65.69,42.45,65.69Zm42.24,0C78.41,65.69,73.31,60,73.31,53s5-12.74,11.43-12.74S96.33,46,96.22,53,91.08,65.69,84.69,65.69Z"/>
            </svg>
          </a>

          <!-- Saweria Icon -->
          <a href="#" target="_blank" class="relative p-2 rounded-lg glass hover:shadow-glow-yellow transition-all flex items-center gap-2" title="Donasi via Saweria">
            <span class="text-xs font-bold text-yellow-500 hidden sm:block">Saweria</span>
            <svg class="w-5 h-5 text-yellow-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </a>

          <!-- Vertical Divider -->
          <div class="h-8 w-px bg-white/10 mx-1"></div>

          <!-- Login / Profile Button -->
          <button @click="handleLoginClick" class="flex items-center gap-3 group cursor-pointer px-3 py-1.5 rounded-xl glass hover:bg-white/5 transition-all">
            <div class="w-8 h-8 rounded-full bg-gradient-to-tr from-cyan-500 to-blue-500 flex items-center justify-center p-[2px]">
              <img src="https://api.dicebear.com/7.x/avataaars/svg?seed=Guest" alt="Guest" class="w-full h-full rounded-full bg-midnight-900 object-cover">
            </div>
            <span class="font-bold text-sm text-gray-200 group-hover:text-white transition-colors">Login</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Decorative Border Glow -->
    <div class="absolute bottom-0 left-0 right-0 h-px bg-gradient-magic opacity-50"></div>
  </nav>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  userName: {
    type: String,
    default: 'PlayerOne'
  },
  userAvatar: {
    type: String,
    default: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Captain'
  },
  userStats: {
    type: Object,
    default: () => ({
      rank: 'Mythic',
      reputation: 95
    })
  },
  hasNotifications: {
    type: Boolean,
    default: false
  }
})

const reputationColor = computed(() => {
  const rep = props.userStats.reputation
  if (rep >= 80) return 'text-green-400'
  if (rep >= 50) return 'text-yellow-400'
  return 'text-red-400'
})

const handleLoginClick = () => {
  alert('🚀 Fitur Login sedang dalam pengembangan! Nantikan update selanjutnya.')
}
</script>
