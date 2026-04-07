<template>
  <Teleport to="body">
    <Transition name="modal-fade">
      <div v-if="isVisible" class="fixed inset-0 z-[110] flex items-center justify-center p-4" @click.self="$emit('cancel')">
        <!-- Backdrop -->
        <div class="absolute inset-0 bg-black/70 backdrop-blur-sm"></div>

        <!-- Modal Card -->
        <div class="relative w-full max-w-sm rounded-2xl border border-red-500/30 bg-[#0d1224]/95 backdrop-blur-xl shadow-[0_0_60px_rgba(220,38,38,0.25)] overflow-hidden" @click.stop>
          <!-- Top gradient bar -->
          <div class="h-1 w-full bg-gradient-to-r from-red-500 via-rose-400 to-red-500"></div>

          <div class="p-8 text-center">
            <!-- Warning Icon -->
            <div class="w-16 h-16 rounded-full bg-red-500/20 text-red-400 flex items-center justify-center mx-auto mb-4 border border-red-500/30 shadow-glow-red">
              <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
              </svg>
            </div>

            <h2 class="text-xl font-['Orbitron'] font-black text-white mb-2">Keluar dari Akun?</h2>
            <p class="text-sm text-gray-400 mb-8">Anda harus login kembali untuk masuk ke dashboard Matchmaking.</p>

            <div class="flex gap-3">
              <button
                class="flex-1 py-3 rounded-xl border border-white/10 bg-white/5 hover:bg-white/10 text-white font-bold transition-all"
                @click="$emit('cancel')"
              >
                Batal
              </button>
              <button
                class="flex-1 py-3 rounded-xl bg-gradient-to-r from-red-600 to-rose-600 hover:from-red-500 hover:to-rose-500 text-white font-bold shadow-[0_0_20px_rgba(225,29,72,0.4)] transition-all"
                @click="$emit('confirm')"
              >
                Ya, Keluar
              </button>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { onMounted, onUnmounted } from 'vue'

const props = defineProps({
  isVisible: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['confirm', 'cancel'])

const handleKeydown = (e) => {
  if (props.isVisible && e.key === 'Escape') {
    emit('cancel')
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
})
</script>

<style scoped>
.shadow-glow-red {
  box-shadow: 0 0 30px rgba(239, 68, 68, 0.4);
}
.modal-fade-enter-active,
.modal-fade-leave-active { transition: opacity 0.25s ease; }
.modal-fade-enter-from,
.modal-fade-leave-to { opacity: 0; }
</style>
