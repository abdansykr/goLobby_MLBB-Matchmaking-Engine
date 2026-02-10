<template>
  <Transition
    enter-active-class="transition-all duration-500 ease-out"
    enter-from-class="opacity-0 scale-75"
    enter-to-class="opacity-100 scale-100"
    leave-active-class="transition-all duration-300 ease-in"
    leave-from-class="opacity-100 scale-100"
    leave-to-class="opacity-0 scale-75"
  >
    <div v-if="isVisible" class="modal-overlay" @click.self="handleDecline">
      <div class="relative max-w-4xl w-full mx-4">
        <!-- Countdown Timer Bar -->
        <div class="mb-4">
          <div class="flex items-center justify-between text-sm mb-2">
            <span class="text-gray-400">Match Confirmation</span>
            <span :class="[
              'font-bold',
              timeLeft > 30 ? 'text-cyan-magic-300' : timeLeft > 15 ? 'text-yellow-400' : 'text-red-400'
            ]">
              {{ timeLeft }}s remaining
            </span>
          </div>
          <div class="h-2 bg-midnight-800 rounded-full overflow-hidden">
            <div 
              :class="[
                'timer-bar h-full transition-all',
                timeLeft <= 15 && 'critical'
              ]"
              :style="{ width: timePercentage + '%' }"
            ></div>
          </div>
        </div>

        <!-- Main Modal Card -->
        <div class="glass-card p-8 md:p-12 relative overflow-hidden">
          <!-- Dramatic Background Effects -->
          <div class="absolute inset-0 opacity-10">
            <div class="absolute top-0 left-0 w-96 h-96 bg-electric-violet-500 rounded-full blur-3xl animate-pulse"></div>
            <div class="absolute bottom-0 right-0 w-96 h-96 bg-cyan-magic-500 rounded-full blur-3xl animate-pulse" style="animation-delay: 0.5s"></div>
          </div>

          <!-- Content -->
          <div class="relative z-10">
            <!-- Match Found Header -->
            <div class="text-center mb-8">
              <div class="inline-block px-6 py-2 bg-gradient-magic rounded-full mb-4 animate-pulse-slow">
                <p class="text-sm font-bold tracking-widest">MATCH FOUND</p>
              </div>
              <h2 class="text-2xl md:text-3xl font-bold text-glow">
                Opponent Located!
              </h2>
            </div>

            <!-- VS Screen -->
            <div class="grid grid-cols-1 md:grid-cols-3 gap-8 items-center mb-8">
              <!-- Your Team -->
              <div class="text-center space-y-4 animate-float">
                <div class="relative inline-block">
                  <!-- Hexagon Avatar with Glow -->
                  <div class="avatar-hexagon mx-auto w-32 h-32 shadow-glow-violet">
                    <img :src="yourTeam.avatar" :alt="yourTeam.name" />
                  </div>
                  <!-- Rank Badge -->
                  <div class="absolute -bottom-2 left-1/2 transform -translate-x-1/2">
                    <div class="rank-badge text-xs px-3 py-1">
                      <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                        <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
                      </svg>
                      {{ yourTeam.rank }}
                    </div>
                  </div>
                </div>
                <div>
                  <p class="text-xl font-bold">{{ yourTeam.name }}</p>
                  <p class="text-sm text-gray-400">Your Team</p>
                </div>
              </div>

              <!-- VS Divider -->
              <div class="flex items-center justify-center">
                <div class="text-center">
                  <div class="vs-text mb-4">VS</div>
                  <div class="flex items-center justify-center gap-2 text-sm">
                    <div class="w-2 h-2 bg-green-400 rounded-full animate-pulse"></div>
                    <span class="text-gray-400">Live Match</span>
                  </div>
                </div>
              </div>

              <!-- Opponent Team -->
              <div class="text-center space-y-4 animate-float" style="animation-delay: 0.2s">
                <div class="relative inline-block">
                  <!-- Hexagon Avatar with Glow -->
                  <div class="avatar-hexagon mx-auto w-32 h-32 shadow-glow-cyan">
                    <img :src="opponentTeam.avatar" :alt="opponentTeam.name" />
                  </div>
                  <!-- Rank Badge -->
                  <div class="absolute -bottom-2 left-1/2 transform -translate-x-1/2">
                    <div class="rank-badge text-xs px-3 py-1">
                      <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                        <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1. 81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
                      </svg>
                      {{ opponentTeam.rank }}
                    </div>
                  </div>
                </div>
                <div>
                  <p class="text-xl font-bold">{{ opponentTeam.name }}</p>
                  <p class="text-sm text-gray-400">Opponent</p>
                </div>
              </div>
            </div>

            <!-- Match Details -->
            <div class="grid grid-cols-3 gap-4 mb-8">
              <div class="glass p-4 rounded-lg text-center">
                <p class="text-xs text-gray-400 mb-1">Category</p>
                <p class="font-bold text-cyan-magic-300">{{ matchDetails.category }}</p>
              </div>
              <div class="glass p-4 rounded-lg text-center">
                <p class="text-xs text-gray-400 mb-1">Rank Diff</p>
                <p class="font-bold text-electric-violet-300">±{{ matchDetails.rankDiff }}</p>
              </div>
              <div class="glass p-4 rounded-lg text-center">
                <p class="text-xs text-gray-400 mb-1">Match ID</p>
                <p class="font-bold text-sm text-antique-gold-400">#{{ matchDetails.id }}</p>
              </div>
            </div>

            <!-- Action Buttons -->
            <div class="space-y-4">
              <!-- Accept Button (Primary) -->
              <button 
                @click="handleAccept"
                class="btn-gold w-full text-lg py-4 flex items-center justify-center gap-3 group">
                <svg class="w-6 h-6 group-hover:scale-110 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
                ACCEPT MATCH
              </button>

              <!-- Decline Button (Secondary) -->
              <button 
                @click="handleDecline"
                class="btn-danger w-full flex items-center justify-center gap-2">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
                Decline
              </button>
            </div>

            <!-- Warning Text -->
            <div class="mt-6 text-center">
              <p class="text-xs text-gray-500">
                <span class="text-red-400">⚠️ Warning:</span> Declining too many matches will affect your reputation score
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  isVisible: {
    type: Boolean,
    default: false
  },
  yourTeam: {
    type: Object,
    default: () => ({
      name: 'Team Alpha',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=TeamAlpha',
      rank: 'Mythic III'
    })
  },
  opponentTeam: {
    type: Object,
    default: () => ({
      name: 'Team Beta',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=TeamBeta',
      rank: 'Mythic II'
    })
  },
  matchDetails: {
    type: Object,
    default: () => ({
      category: 'POKE',
      rankDiff: 1,
      id: 'A2F4B8'
    })
  },
  timeout: {
    type: Number,
    default: 60
  }
})

const emit = defineEmits(['accept', 'decline', 'timeout'])

// Timer state
const timeLeft = ref(props.timeout)
const totalTime = props.timeout
let timerInterval = null

const timePercentage = computed(() => {
  return (timeLeft.value / totalTime) * 100
})

// Timer logic
const startTimer = () => {
  timerInterval = setInterval(() => {
    timeLeft.value--
    
    if (timeLeft.value <= 0) {
      clearInterval(timerInterval)
      emit('timeout')
    }
  }, 1000)
}

const stopTimer = () => {
  if (timerInterval) {
    clearInterval(timerInterval)
  }
}

// Handlers
const handleAccept = () => {
  stopTimer()
  emit('accept')
}

const handleDecline = () => {
  stopTimer()
  emit('decline')
}

// Lifecycle
onMounted(() => {
  if (props.isVisible) {
    startTimer()
  }
})

onUnmounted(() => {
  stopTimer()
})
</script>
