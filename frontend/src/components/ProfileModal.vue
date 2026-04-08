<template>
  <Teleport to="body">
    <Transition name="modal-fade">
      <div v-if="isVisible" class="fixed inset-0 z-[100] flex items-center justify-center p-4" @click.self="$emit('close')">
        <div class="absolute inset-0 bg-black/70 backdrop-blur-sm"></div>

        <div class="relative w-full max-w-md rounded-2xl border border-blue-500/30 bg-[#0d1224]/95 backdrop-blur-xl shadow-[0_0_60px_rgba(139,92,246,0.25)] overflow-hidden" @click.stop>
          <div class="h-1 w-full bg-gradient-to-r from-electric-violet-500 via-cyan-magic-400 to-electric-violet-500"></div>

          <div class="p-6">
            <h2 class="text-xl font-['Orbitron'] font-black text-white mb-4">Pengaturan Akun</h2>

            <!-- Tab Switcher -->
            <div class="flex gap-1 p-1 bg-white/5 rounded-xl mb-5">
              <button
                class="flex-1 py-2 px-3 rounded-lg text-sm font-bold transition-all duration-200"
                :class="activeTab === 'profile' ? 'bg-electric-violet-500 text-white shadow-glow-violet' : 'text-gray-400 hover:text-white'"
                @click="activeTab = 'profile'"
              >✏️ Edit Profil</button>
              <button
                class="flex-1 py-2 px-3 rounded-lg text-sm font-bold transition-all duration-200"
                :class="activeTab === 'password' ? 'bg-electric-violet-500 text-white shadow-glow-violet' : 'text-gray-400 hover:text-white'"
                @click="activeTab = 'password'; pwMsg.text = ''"
              >🔑 Ganti Password</button>
            </div>

            <!-- ── TAB: Edit Profil ── -->
            <form v-if="activeTab === 'profile'" @submit.prevent="handleSubmit" class="space-y-4">

              <!-- Avatar Upload -->
              <div class="flex flex-col items-center mb-4">
                <div class="relative w-20 h-20 rounded-full overflow-hidden mb-2 border-2 border-cyan-magic-400/50">
                  <img :src="avatarPreview || user?.avatar_url || `https://api.dicebear.com/7.x/identicon/svg?seed=${user?.username}`" class="w-full h-full object-cover" />
                  <div class="absolute inset-0 bg-black/50 flex items-center justify-center opacity-0 hover:opacity-100 transition-opacity cursor-pointer">
                    <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z"></path></svg>
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
                <label class="block text-xs font-bold text-gray-400 uppercase tracking-widest mb-1.5">No. WhatsApp</label>
                <input v-model="form.whatsapp_number" type="tel" class="auth-input" :class="{'border-red-500/50': whatsappError}" :disabled="loading" required>
                <p v-if="whatsappError" class="text-xs font-bold text-red-400 mt-1.5">{{ whatsappError }}</p>
              </div>

              <Transition name="msg-fade">
                <div v-if="message.text" class="flex items-start gap-2.5 p-3 rounded-xl text-sm font-medium" :class="message.type === 'error' ? 'bg-red-900/40 border border-red-700/50 text-red-300' : 'bg-emerald-900/40 border border-emerald-700/50 text-emerald-300'">
                  {{ message.text }}
                </div>
              </Transition>

              <button type="submit" class="w-full py-3 rounded-xl font-bold text-white bg-gradient-magic shadow-glow-violet hover:shadow-glow-cyan hover:scale-[1.02] active:scale-95 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed" :disabled="loading || !!whatsappError">
                <span v-if="loading">Menyimpan...</span>
                <span v-else>Simpan Perubahan</span>
              </button>
            </form>

            <!-- ── TAB: Ganti Password ── -->
            <form v-else-if="activeTab === 'password'" @submit.prevent="handleChangePassword" class="space-y-4">
              <div class="bg-blue-900/20 border border-blue-500/30 rounded-xl p-3 text-xs text-blue-300 flex items-start gap-2">
                <svg class="w-4 h-4 shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
                <span>Masukkan password lama Anda untuk verifikasi, lalu buat password baru minimal 6 karakter.</span>
              </div>

              <div>
                <label class="block text-xs font-bold text-gray-400 uppercase tracking-widest mb-1.5">Password Lama</label>
                <div class="relative">
                  <input v-model="pwForm.current" :type="showCurrentPw ? 'text' : 'password'" class="auth-input pr-10" :disabled="pwLoading" placeholder="Password saat ini" required>
                  <button type="button" class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-200 transition-colors" @click="showCurrentPw = !showCurrentPw" tabindex="-1">
                    <svg v-if="!showCurrentPw" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/></svg>
                    <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"/></svg>
                  </button>
                </div>
              </div>

              <div>
                <label class="block text-xs font-bold text-gray-400 uppercase tracking-widest mb-1.5">Password Baru</label>
                <div class="relative">
                  <input v-model="pwForm.newPw" :type="showNewPw ? 'text' : 'password'" class="auth-input pr-10" :disabled="pwLoading" placeholder="Min. 6 karakter" required>
                  <button type="button" class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-200 transition-colors" @click="showNewPw = !showNewPw" tabindex="-1">
                    <svg v-if="!showNewPw" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/></svg>
                    <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"/></svg>
                  </button>
                </div>
                <p v-if="pwForm.newPw && pwForm.newPw.length < 6" class="text-xs text-red-400 mt-1">⚠️ Minimal 6 karakter</p>
              </div>

              <div>
                <label class="block text-xs font-bold text-gray-400 uppercase tracking-widest mb-1.5">Konfirmasi Password Baru</label>
                <input v-model="pwForm.confirm" :type="showNewPw ? 'text' : 'password'" class="auth-input" :class="{'border-red-500/50': pwForm.confirm && pwForm.confirm !== pwForm.newPw}" :disabled="pwLoading" placeholder="Ulangi password baru" required>
                <p v-if="pwForm.confirm && pwForm.confirm !== pwForm.newPw" class="text-xs text-red-400 mt-1">⚠️ Password tidak cocok</p>
              </div>

              <Transition name="msg-fade">
                <div v-if="pwMsg.text" class="flex items-start gap-2.5 p-3 rounded-xl text-sm font-medium" :class="pwMsg.type === 'error' ? 'bg-red-900/40 border border-red-700/50 text-red-300' : 'bg-emerald-900/40 border border-emerald-700/50 text-emerald-300'">
                  {{ pwMsg.text }}
                </div>
              </Transition>

              <button type="submit" class="w-full py-3 rounded-xl font-bold text-white bg-gradient-magic shadow-glow-violet hover:shadow-glow-cyan hover:scale-[1.02] active:scale-95 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed" :disabled="pwLoading || pwForm.newPw !== pwForm.confirm || pwForm.newPw.length < 6">
                <span v-if="pwLoading">Menyimpan...</span>
                <span v-else>🔒 Perbarui Password</span>
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

const { user, updateProfile, uploadAvatar, changePassword } = useAuth()

// ── Tabs ──────────────────────────────────────────────────────────────────────
const activeTab = ref('profile')

// ── Profile tab state ─────────────────────────────────────────────────────────
const loading = ref(false)
const selectedFile = ref(null)
const avatarPreview = ref(null)
const message = reactive({ text: '', type: '' })
const form = reactive({ username: '', email: '', whatsapp_number: '' })

// ── Password tab state ────────────────────────────────────────────────────────
const pwLoading = ref(false)
const showCurrentPw = ref(false)
const showNewPw = ref(false)
const pwMsg = reactive({ text: '', type: '' })
const pwForm = reactive({ current: '', newPw: '', confirm: '' })

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
    // Reset profile form
    form.username = user.value.username || ''
    form.email = user.value.email || ''
    form.whatsapp_number = user.value.whatsapp_number || ''
    avatarPreview.value = null
    selectedFile.value = null
    message.text = ''
    // Reset password form
    pwForm.current = ''
    pwForm.newPw = ''
    pwForm.confirm = ''
    pwMsg.text = ''
    activeTab.value = 'profile'
    showCurrentPw.value = false
    showNewPw.value = false
  }
})

const handleKeydown = (e) => {
  if (props.isVisible && e.key === 'Escape') emit('close')
}
onMounted(() => document.addEventListener('keydown', handleKeydown))
onUnmounted(() => document.removeEventListener('keydown', handleKeydown))

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
    if (selectedFile.value) await uploadAvatar(selectedFile.value)
    await updateProfile(form)
    message.text = 'Profil berhasil diperbarui ✅'
    message.type = 'success'
    setTimeout(() => emit('close'), 1500)
  } catch (err) {
    message.text = err.message || 'Gagal menyimpan profil'
    message.type = 'error'
  } finally {
    loading.value = false
  }
}

const handleChangePassword = async () => {
  if (pwLoading.value) return
  if (pwForm.newPw !== pwForm.confirm) {
    pwMsg.text = 'Password konfirmasi tidak cocok'
    pwMsg.type = 'error'
    return
  }
  pwLoading.value = true
  pwMsg.text = ''
  try {
    await changePassword({ currentPassword: pwForm.current, newPassword: pwForm.newPw })
    pwMsg.text = '🔒 Password berhasil diperbarui!'
    pwMsg.type = 'success'
    pwForm.current = ''
    pwForm.newPw = ''
    pwForm.confirm = ''
  } catch (err) {
    pwMsg.text = err.message || 'Gagal mengubah password'
    pwMsg.type = 'error'
  } finally {
    pwLoading.value = false
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
