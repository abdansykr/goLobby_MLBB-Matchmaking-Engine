<template>
  <Transition
    enter-active-class="transition-all duration-500 ease-out"
    enter-from-class="opacity-0 scale-75"
    enter-to-class="opacity-100 scale-100"
    leave-active-class="transition-all duration-300 ease-in"
    leave-from-class="opacity-100 scale-100"
    leave-to-class="opacity-0 scale-75"
  >
    <!-- ═══════════════════════════════════════════════
         STATE 1: LAWAN DITEMUKAN (pending / confirmed)
         ═══════════════════════════════════════════════ -->
    <div v-if="isVisible && internalStatus !== 'cancelled'" class="modal-overlay">
      <div class="relative max-w-4xl w-full mx-4">
        <!-- Countdown Timer Bar -->
        <div class="mb-4">
          <div class="flex items-center justify-between text-sm mb-2">
            <span class="text-gray-400">Konfirmasi Pertandingan</span>
            <span :class="[
              'font-bold',
              timeLeft > 30 ? 'text-cyan-magic-300' : timeLeft > 15 ? 'text-yellow-400' : 'text-red-400'
            ]">
              Tersisa {{ timeLeft }} detik
            </span>
          </div>
          <div class="h-2 bg-midnight-800 rounded-full overflow-hidden">
            <div
              :class="['timer-bar h-full transition-all', timeLeft <= 15 && 'critical']"
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

          <div class="relative z-10">
            <!-- Header -->
            <div class="text-center mb-8">
              <div class="inline-block px-6 py-2 bg-gradient-magic rounded-full mb-4 animate-pulse-slow">
                <p class="text-sm font-bold tracking-widest">PERSIAPAN MABAR</p>
              </div>
              <h2 class="text-2xl md:text-3xl font-bold text-glow">
                {{ internalStatus === 'confirmed' ? '✅ Pertandingan Dikonfirmasi!' : internalStatus === 'waiting' ? 'Menunggu Lawan...' : 'Lawan Ditemukan!' }}
              </h2>
            </div>

            <!-- VS Screen -->
            <div class="grid grid-cols-1 md:grid-cols-3 gap-8 items-center mb-8">
              <!-- Your Team -->
              <div class="text-center space-y-4 animate-float">
                <div class="relative inline-block">
                  <div class="avatar-hexagon mx-auto w-32 h-32 shadow-glow-violet">
                    <img :src="yourTeam.avatar" :alt="yourTeam.name" />
                  </div>
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
                  <p class="text-sm text-gray-400">Tim Anda</p>
                </div>
              </div>

              <!-- VS Divider -->
              <div class="flex items-center justify-center">
                <div class="text-center">
                  <div class="vs-text mb-4">VS</div>
                  <div class="flex items-center justify-center gap-2 text-sm">
                    <div :class="['w-2 h-2 rounded-full animate-pulse', internalStatus === 'confirmed' ? 'bg-green-400' : 'bg-yellow-400']"></div>
                    <span class="text-gray-400">{{ internalStatus === 'confirmed' ? 'Terkonfirmasi' : internalStatus === 'waiting' ? 'Menunggu Persetujuan' : 'Sedang Terhubung' }}</span>
                  </div>
                </div>
              </div>

              <!-- Opponent Team -->
              <div class="text-center space-y-4 animate-float" style="animation-delay: 0.2s">
                <div class="relative inline-block">
                  <div class="avatar-hexagon mx-auto w-32 h-32 shadow-glow-cyan">
                    <img :src="opponentTeam.avatar" :alt="opponentTeam.name" />
                  </div>
                  <div class="absolute -bottom-2 left-1/2 transform -translate-x-1/2">
                    <div class="rank-badge text-xs px-3 py-1">
                      <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                        <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
                      </svg>
                      {{ opponentTeam.rank }}
                    </div>
                  </div>
                </div>
                <div>
                  <p class="text-xl font-bold">{{ opponentTeam.name }}</p>
                  <p class="text-sm text-gray-400">Tim Lawan</p>
                </div>
              </div>
            </div>

            <!-- Match Details -->
            <div class="grid grid-cols-3 gap-4 mb-8">
              <div class="glass p-4 rounded-lg text-center">
                <p class="text-xs text-gray-400 mb-1">Kategori</p>
                <p class="font-bold text-cyan-magic-300">{{ matchDetails.category }}</p>
              </div>
              <div class="glass p-4 rounded-lg text-center">
                <p class="text-xs text-gray-400 mb-1">Selisih Rank</p>
                <p class="font-bold text-electric-violet-300">±{{ matchDetails.rankDiff }}</p>
              </div>
              <div class="glass p-4 rounded-lg text-center">
                <p class="text-xs text-gray-400 mb-1">Match ID</p>
                <p class="font-bold text-sm text-antique-gold-400">#{{ matchDetails.id }}</p>
              </div>
            </div>

            <!-- Action Buttons -->
            <div class="space-y-4">
              <!-- PENDING: Terima Match -->
              <button
                v-if="internalStatus === 'pending'"
                @click="handleAccept"
                :disabled="accepting"
                class="btn-gold w-full text-lg py-4 flex items-center justify-center gap-3 group disabled:opacity-60 disabled:cursor-wait">
                <svg v-if="!accepting" class="w-6 h-6 group-hover:scale-110 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
                <svg v-else class="w-6 h-6 animate-spin" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
                </svg>
                {{ accepting ? 'Mengkonfirmasi...' : 'TERIMA MATCH' }}
              </button>

              <!-- WAITING: Menunggu konfirmasi lawan -->
              <div v-if="internalStatus === 'waiting'" class="w-full text-lg py-4 flex items-center justify-center gap-3 bg-midnight-700/80 text-yellow-500 font-bold rounded-xl animate-pulse cursor-wait border border-yellow-500/30">
                <svg class="w-6 h-6 animate-spin" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
                </svg>
                Menunggu Konfirmasi Lawan...
              </div>

              <!-- CONFIRMED: WhatsApp Button -->
              <template v-if="internalStatus === 'confirmed'">
                <div class="text-center py-2 text-green-400 font-semibold text-lg">
                  ✅ Kamu sudah menerima pertandingan ini!
                </div>
                <a
                  v-if="whatsappUrl"
                  :href="whatsappUrl"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="flex items-center justify-center gap-3 w-full py-4 px-6 rounded-xl
                         bg-green-600 hover:bg-green-500 text-white font-bold text-lg
                         transition-all duration-300 shadow-lg hover:shadow-green-500/30 animate-pulse-slow">
                  <svg class="w-6 h-6" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 005.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 00-3.48-8.413Z"/>
                  </svg>
                  Hubungi Lawan via WhatsApp
                </a>
                <button @click="handleClose" class="w-full py-3 rounded-xl border border-gray-600 text-gray-400 hover:text-white hover:border-gray-400 transition-all">
                  Tutup
                </button>
              </template>

              <!-- Decline Button — only shown on pending -->
              <button
                v-if="internalStatus === 'pending'"
                @click="handleDecline"
                class="btn-danger w-full flex items-center justify-center gap-2">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
                Tolak Lawan
              </button>
            </div>

            <!-- Warning Text -->
            <div class="mt-6 text-center">
              <p class="text-xs text-gray-500">
                <span class="text-red-400">⚠️ Peringatan:</span> Sering membatalkan Match akan memberlakukan pengurangan Skor Reputasi/Kredit.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- ═══════════════════════════════════════════════
         STATE 2: MATCH DIBATALKAN OLEH LAWAN
         ═══════════════════════════════════════════════ -->
    <div v-else-if="isVisible && internalStatus === 'cancelled'" class="modal-overlay">
      <div class="relative max-w-lg w-full mx-4">
        <div class="glass-card p-8 md:p-12 relative overflow-hidden text-center">
          <div class="absolute inset-0 opacity-10">
            <div class="absolute inset-0 bg-red-900 rounded-3xl blur-3xl"></div>
          </div>
          <div class="relative z-10 space-y-6">
            <!-- Icon -->
            <div class="flex justify-center">
              <div class="w-20 h-20 rounded-full bg-red-900/50 border-2 border-red-500 flex items-center justify-center animate-pulse">
                <svg class="w-10 h-10 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
                </svg>
              </div>
            </div>

            <div>
              <h2 class="text-2xl font-bold text-red-400 mb-2">Pertandingan Dibatalkan</h2>
              <p class="text-gray-400 text-sm">{{ cancelReason || 'Lawan menolak pertandingan ini.' }}</p>
            </div>

            <!-- Tombol Cari Lagi -->
            <button
              @click="handleSearchAgain"
              class="btn-gold w-full text-lg py-4 flex items-center justify-center gap-3 group">
              <svg class="w-6 h-6 group-hover:rotate-180 transition-transform duration-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
              </svg>
              Cari Lawan Baru
            </button>

            <button @click="handleClose" class="w-full py-3 rounded-xl border border-gray-600 text-gray-400 hover:text-white hover:border-gray-400 transition-all text-sm">
              Kembali ke Menu
            </button>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup>
import { ref, computed, watch, onUnmounted } from 'vue'

const props = defineProps({
  isVisible: {
    type: Boolean,
    default: false
  },
  yourTeam: {
    type: Object,
    default: () => ({ name: 'Tim Anda', avatar: '', rank: 'Epic' })
  },
  opponentTeam: {
    type: Object,
    default: () => ({ name: 'Tim Lawan', avatar: '', rank: 'Epic' })
  },
  matchDetails: {
    type: Object,
    default: () => ({ category: 'POKE', rankDiff: 0, id: 'MATCH' })
  },
  timeout: {
    type: Number,
    default: 60
  },
  // 'pending' | 'waiting' | 'confirmed' | 'cancelled'
  matchStatus: {
    type: String,
    default: 'pending'
  },
  whatsappUrl: {
    type: String,
    default: ''
  },
  cancelReason: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['accept', 'decline', 'timeout', 'close', 'search-again'])

// ── Internal state ───────────────────────────────────────────────
// Mirror matchStatus so we can show cancelled state without closing the modal
const internalStatus = ref(props.matchStatus)
const accepting = ref(false)

// When parent changes matchStatus (e.g. opponent rejected), update internal state
watch(() => props.matchStatus, (val) => {
  internalStatus.value = val
  // If match is cancelled while modal is open, stop the timer
  if (val === 'cancelled') stopTimer()
  // If match is confirmed, stop the timer
  if (val === 'confirmed') stopTimer()
})

// ── Timer ────────────────────────────────────────────────────────
const timeLeft = ref(props.timeout)
let timerInterval = null

const timePercentage = computed(() => (timeLeft.value / props.timeout) * 100)

const startTimer = () => {
  stopTimer()
  timeLeft.value = props.timeout
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
    timerInterval = null
  }
}

// Restart timer every time modal becomes visible with a fresh match
watch(() => props.isVisible, (visible) => {
  if (visible && internalStatus.value === 'pending') {
    internalStatus.value = props.matchStatus
    startTimer()
  } else {
    stopTimer()
  }
}, { immediate: true })

// ── Handlers ─────────────────────────────────────────────────────
const handleAccept = async () => {
  if (accepting.value || internalStatus.value !== 'pending') return
  accepting.value = true
  stopTimer()
  emit('accept')
  // Parent will set matchStatus → 'confirmed', which our watcher picks up
  accepting.value = false
}

const handleDecline = () => {
  stopTimer()
  emit('decline')
}

const handleClose = () => {
  stopTimer()
  emit('close')
}

const handleSearchAgain = () => {
  stopTimer()
  emit('search-again')
}

onUnmounted(stopTimer)
</script>
