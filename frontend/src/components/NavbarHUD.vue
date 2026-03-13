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

        <!-- User Profile Section -->
        <div class="flex items-center gap-4">
          <!-- Notification Icon -->
          <button class="relative p-2 rounded-lg glass hover:shadow-glow-cyan transition-all">
            <svg class="w-5 h-5 text-cyan-magic-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
            </svg>
            <!-- Notification Dot -->
            <span v-if="hasNotifications" 
              class="absolute top-1 right-1 w-2 h-2 bg-red-500 rounded-full animate-pulse"></span>
          </button>

          <!-- Avatar with Hexagon Frame -->
          <div class="flex items-center gap-3 group cursor-pointer">
            <div class="avatar-hexagon energy-border">
              <img :src="userAvatar" 
                :alt="userName" 
                class="w-full h-full object-cover">
            </div>
            <div class="hidden lg:block">
              <p class="font-bold text-sm">{{ userName }}</p>
              <p class="text-xs text-gray-400">Kapten</p>
            </div>
            <!-- Dropdown Icon -->
            <svg class="w-4 h-4 text-gray-400 group-hover:text-cyan-magic-300 transition-colors" 
              fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
            </svg>
          </div>
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
</script>
