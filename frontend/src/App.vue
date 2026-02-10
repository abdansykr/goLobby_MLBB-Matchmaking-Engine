<template>
  <div id="app" class="min-h-screen">
    <!-- Navbar -->
    <nav class="fixed top-0 left-0 right-0 z-50 glass border-b border-electric-violet-500/20 backdrop-blur-md bg-midnight-900/80">
      <div class="container mx-auto px-4 h-20 flex items-center justify-between">
        <!-- Logo and Brand -->
        <router-link to="/" class="flex items-center gap-4 group cursor-pointer">
          <!-- Logo Icon -->
          <div class="w-10 h-10 rounded-full bg-gradient-magic flex items-center justify-center shadow-glow-violet group-hover:scale-110 transition-transform duration-300">
            <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
          </div>
          
          <!-- Brand Text -->
          <div class="flex flex-col">
            <h1 class="text-2xl font-['Orbitron'] font-black tracking-wider text-white group-hover:text-cyan-magic-300 transition-colors">GOLOBBY</h1>
            <span class="text-[10px] text-cyan-magic-400 font-bold tracking-[0.2em] uppercase -mt-1">Scrim Arena</span>
          </div>
        </router-link>

        <!-- Right Side: Removed -->
      </div>
    </nav>
    
    <!-- Main Content with Padding for Navbar -->
    <div class="pt-24 min-h-screen">
      <router-view v-slot="{ Component }">
        <Transition
          mode="out-in"
          enter-active-class="transition-opacity duration-300"
          enter-from-class="opacity-0"
          enter-to-class="opacity-100"
          leave-active-class="transition-opacity duration-200"
          leave-from-class="opacity-100"
          leave-to-class="opacity-0"
        >
          <component :is="Component" />
        </Transition>
      </router-view>
    </div>

    <!-- Global Match Found Modal -->
    <MatchFoundModal 
      :isVisible="showMatchModal"
      :yourTeam="matchData.yourTeam"
      :opponentTeam="matchData.opponentTeam"
      :matchDetails="matchData.details"
      @accept="handleMatchAccept"
      @decline="handleMatchDecline"
      @timeout="handleMatchTimeout"
    />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import MatchFoundModal from './components/MatchFoundModal.vue'

// Match modal state
const showMatchModal = ref(false)
const matchData = ref({
  yourTeam: {},
  opponentTeam: {},
  details: {}
})

// Match handlers
const handleMatchAccept = () => {
  console.log('Match accepted!')
  showMatchModal.value = false
  // TODO: Navigate to lobby
}

const handleMatchDecline = () => {
  console.log('Match declined')
  showMatchModal.value = false
  // TODO: Return to search or show cooldown
}

const handleMatchTimeout = () => {
  console.log('Match timed out')
  showMatchModal.value = false
  // TODO: Show penalty notification
}

// Demo: Show modal after 5 seconds (remove in production)
// setTimeout(() => {
//   matchData.value = {
//     yourTeam: {
//       name: 'Team Alpha',
//       avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Alpha',
//       rank: 'Mythic III'
//     },
//     opponentTeam: {
//       name: 'Team Beta',
//       avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Beta',
//       rank: 'Mythic II'
//     },
//     details: {
//       category: 'POKE',
//       rankDiff: 1,
//       id: 'A2F4B8'
//     }
//   }
//   showMatchModal.value = true
// }, 5000)
</script>
