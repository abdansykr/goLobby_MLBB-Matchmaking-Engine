<template>
  <Teleport to="body">
    <Transition name="modal-fade">
      <div v-if="isVisible" class="fixed inset-0 z-[100] flex items-center justify-center p-4" @click.self="$emit('close')">
        <div class="absolute inset-0 bg-black/70 backdrop-blur-sm"></div>

        <div class="relative w-full max-w-md rounded-2xl border border-blue-500/30 bg-[#0d1224]/95 backdrop-blur-xl shadow-[0_0_60px_rgba(139,92,246,0.25)] overflow-hidden" @click.stop>
          <div class="h-1 w-full bg-gradient-to-r from-electric-violet-500 via-cyan-magic-400 to-electric-violet-500"></div>

          <div class="p-8">
            <h2 class="text-2xl font-['Orbitron'] font-black text-white mb-6">Profil Tim</h2>

            <form @submit.prevent="handleSubmit" class="space-y-4">
              
              <!-- Avatar Upload -->
              <div class="flex flex-col items-center mb-6">
                 <div class="relative w-24 h-24 rounded-full overflow-hidden mb-2 border-2 border-cyan-magic-400/50">
                    <img :src="avatarPreview || user?.avatar_url || `https://api.dicebear.com/7.x/identicon/svg?seed=${user?.username}`" class="w-full h-full object-cover" />
                    <div class="absolute inset-0 bg-black/50 flex items-center justify-center opacity-0 hover:opacity-100 transition-opacity cursor-pointer">
                      <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z"></path></svg>
                    </div>
                    <input type="file" accept="image/*" @change="handleFileChange" class="absolute inset-0 opacity-0 cursor-pointer">
                 </div>
                 <p class="text-xs text-gray-400">Klik gambar untuk mengubah (Maks 2MB)</p>
              </div>

              <div>
                <label class="block text-xs font-bold text-gray-400 uppercase tracking-widest mb-1.5">Username</label>
                <input v-model="form.username" type="text" class="auth-input" :disabled="loading" required>
              </div>
              
              <div>
                <label class="block text-xs font-bold text-gray-400 uppercase tracking-widest mb-1.5">Email</label>
                <input v-model="form.email" type="email" class="auth-input" :disabled="loading" required>
              </div>
              
              <div>
                <input v-model="form.whatsapp_number" type="tel" class="auth-input" :class="{'border-red-500/50 focus:border-red-500/70 focus:shadow-[0_0_0_3px_rgba(239,68,68,0.15)]': whatsappError}" :disabled="loading" required>
                <!-- Validation Error -->
                <p v-if="whatsappError" class="text-xs font-bold text-red-400 mt-1.5 flex items-center gap-1">
                  {{ whatsappError }}
                </p>
              </div>

              <Transition name="msg-fade">
                <div v-if="message.text" class="flex items-start gap-2.5 p-3 rounded-xl text-sm font-medium mt-4" :class="message.type === 'error' ? 'bg-red-900/40 border border-red-700/50 text-red-300' : 'bg-emerald-900/40 border border-emerald-700/50 text-emerald-300'">
                  {{ message.text }}
                </div>
              </Transition>

              <button type="submit" class="w-full py-3.5 rounded-xl font-bold text-white bg-gradient-magic shadow-glow-violet hover:shadow-glow-cyan hover:scale-[1.02] active:scale-95 transition-all duration-200 mt-6 disabled:opacity-50 disabled:cursor-not-allowed" :disabled="loading || !!whatsappError">
                <span v-if="loading">Menyimpan...</span>
                <span v-else>Simpan Perubahan</span>
              </button>
            </form>

            <button class="absolute top-4 right-4 p-2 rounded-lg text-gray-500 hover:text-white transition-all" @click="$emit('close')">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { ref, reactive, watch, computed, onMounted, onUnmounted } from 'vue'
import { useAuth } from '../composables/useAuth'

const props = defineProps({ isVisible: Boolean })
const emit = defineEmits(['close'])

const { user, updateProfile, uploadAvatar } = useAuth()

const loading = ref(false)
const selectedFile = ref(null)
const avatarPreview = ref(null)
const message = reactive({ text: '', type: '' })

const form = reactive({
  username: '',
  email: '',
  whatsapp_number: ''
})

const whatsappError = computed(() => {
  const wa = form.whatsapp_number.trim()
  if (!wa) return ''
  if (!wa.startsWith('628')) return '⚠️ Must start with 628'
  if (wa.length < 10) return '⚠️ Too short (min 10 digits)'
  if (wa.length > 15) return '⚠️ Too long (max 15 digits)'
  if (!/^[0-9]+$/.test(wa)) return '⚠️ Numbers only'
  return ''
})

watch(() => props.isVisible, (val) => {
  if (val && user.value) {
    form.username = user.value.username || ''
    form.email = user.value.email || ''
    form.whatsapp_number = user.value.whatsapp_number || ''
    avatarPreview.value = null
    selectedFile.value = null
    message.text = ''
  }
})

const handleKeydown = (e) => {
  if (props.isVisible && e.key === 'Escape') {
    emit('close')
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
})

const handleFileChange = (e) => {
  const file = e.target.files[0]
  if (file) {
    if (file.size > 2 * 1024 * 1024) {
      message.text = 'File maksimum 2MB'
      message.type = 'error'
      return
    }
    selectedFile.value = file
    const reader = new FileReader()
    reader.onload = (e) => { avatarPreview.value = e.target.result }
    reader.readAsDataURL(file)
  }
}

const handleSubmit = async () => {
  if (loading.value) return
  loading.value = true
  message.text = ''

  try {
    if (selectedFile.value) {
      await uploadAvatar(selectedFile.value)
    }
    await updateProfile(form)
    message.text = 'Profil berhasil diperbarui'
    message.type = 'success'
    setTimeout(() => { emit('close') }, 1500)
  } catch (err) {
    message.text = err.message || 'Gagal menyimpan profil'
    message.type = 'error'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-input {
  @apply w-full px-4 py-3 rounded-xl text-sm text-white bg-white/5 border border-white/10 placeholder-gray-600 focus:outline-none focus:border-blue-500/70 focus:shadow-[0_0_0_3px_rgba(59,130,246,0.15)] transition-all disabled:opacity-50;
}
.modal-fade-enter-active, .modal-fade-leave-active { transition: opacity 0.25s ease; }
.modal-fade-enter-from, .modal-fade-leave-to { opacity: 0; }
.msg-fade-enter-active, .msg-fade-leave-active { transition: all 0.2s ease; }
.msg-fade-enter-from, .msg-fade-leave-to { opacity: 0; transform: translateY(-4px); }
</style>
