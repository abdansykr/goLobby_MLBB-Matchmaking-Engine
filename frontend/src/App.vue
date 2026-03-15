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

        <!-- Action & Identity Section -->
        <div class="flex items-center gap-3 sm:gap-4">
          <!-- Discord Icon (Coming Soon) -->
          <router-link to="/coming-soon" class="relative p-2.5 rounded-xl glass hover:bg-white/5 hover:shadow-glow-cyan transition-all group" title="Join Discord">
            <svg class="w-5 h-5 text-[#5865F2] group-hover:scale-110 transition-transform" fill="currentColor" viewBox="0 0 127.14 96.36">
              <path d="M107.7,8.07A105.15,105.15,0,0,0,81.47,0a72.06,72.06,0,0,0-3.36,6.83A97.68,97.68,0,0,0,49,6.83,72.37,72.37,0,0,0,45.64,0,105.89,105.89,0,0,0,19.39,8.09C2.79,32.65-1.71,56.6.54,80.21h0A105.73,105.73,0,0,0,32.71,96.36,77.7,77.7,0,0,0,39.6,85.25a68.42,68.42,0,0,1-10.85-5.18c.91-.66,1.8-1.34,2.66-2a75.57,75.57,0,0,0,64.32,0c.87.71,1.76,1.39,2.66,2a68.68,68.68,0,0,1-10.87,5.19,77.7,77.7,0,0,0,6.89,11.1A105.25,105.25,0,0,0,126.6,80.22h0C129.24,52.84,122.09,29.11,107.7,8.07ZM42.45,65.69C36.18,65.69,31,60,31,53s5-12.74,11.43-12.74S54,46,53.89,53,48.84,65.69,42.45,65.69Zm42.24,0C78.41,65.69,73.31,60,73.31,53s5-12.74,11.43-12.74S96.33,46,96.22,53,91.08,65.69,84.69,65.69Z"/>
            </svg>
          </router-link>

          <!-- Saweria Icon -->
          <a href="https://saweria.co/golobby" target="_blank" class="relative px-4 py-2 rounded-xl glass hover:bg-white/10 hover:shadow-glow-yellow transition-all flex items-center gap-2 group animate-glow-pulse" title="Support Us via Saweria">
            <svg class="w-5 h-5 text-yellow-400 group-hover:scale-110 transition-transform" fill="currentColor" viewBox="0 0 24 24">
              <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z" />
            </svg>
            <span class="text-xs font-bold text-yellow-400 hidden sm:block tracking-wide uppercase group-hover:text-yellow-300 transition-colors">Support Us</span>
          </a>

          <!-- Vertical Divider -->
          <div class="h-8 w-px bg-white/10 mx-1 hidden sm:block"></div>

          <!-- Login / Profile Button (Coming Soon) -->
          <router-link to="/coming-soon" class="flex items-center gap-3 group cursor-pointer px-3 py-1.5 rounded-xl glass hover:bg-white/5 transition-all">
            <div class="w-8 h-8 rounded-full bg-gradient-to-tr from-cyan-500 to-electric-violet-500 flex items-center justify-center p-[2px] group-hover:shadow-glow-violet transition-shadow">
              <img src="https://api.dicebear.com/7.x/avataaars/svg?seed=Guest" alt="Guest Avatar" class="w-full h-full rounded-full bg-midnight-900 object-cover">
            </div>
            <span class="font-bold text-sm text-gray-300 group-hover:text-white transition-colors tracking-wide">Login</span>
          </router-link>
        </div>
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
