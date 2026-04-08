<template>
  <!-- ══════════════════════════════════════════════════════════════
       AUTH MODAL — Login & Register
       Design: Glassmorphism dark, blue/violet accents, no layout change
       ══════════════════════════════════════════════════════════════ -->
  <Teleport to="body">
    <Transition name="modal-fade">
      <div
        v-if="isVisible"
        class="fixed inset-0 z-[100] flex items-center justify-center p-4"
        @click.self="$emit('close')"
      >
        <!-- Backdrop -->
        <div class="absolute inset-0 bg-black/70 backdrop-blur-sm"></div>

        <!-- Modal Card -->
        <div
          class="relative w-full max-w-md rounded-2xl border border-blue-500/30 bg-[#0d1224]/95 backdrop-blur-xl shadow-[0_0_60px_rgba(139,92,246,0.25)] overflow-hidden"
          @click.stop
        >
          <!-- Top gradient bar -->
          <div class="h-1 w-full bg-gradient-to-r from-electric-violet-500 via-cyan-magic-400 to-electric-violet-500"></div>

          <!-- Header -->
          <div class="px-8 pt-8 pb-4">
            <!-- Logo row -->
            <div class="flex items-center gap-3 mb-6">
              <div class="w-9 h-9 rounded-full bg-gradient-magic flex items-center justify-center shadow-glow-violet">
                <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
                </svg>
              </div>
              <span class="font-['Orbitron'] font-black text-xl gradient-text">GOLOBBY</span>
            </div>

            <!-- Tab switcher (hidden in forgot mode) -->
            <div v-if="mode !== 'forgot'" class="flex gap-1 p-1 bg-white/5 rounded-xl mb-6">
              <button
                id="auth-tab-login"
                class="flex-1 py-2.5 px-4 rounded-lg text-sm font-bold transition-all duration-200"
                :class="mode === 'login'
                  ? 'bg-electric-violet-500 text-white shadow-glow-violet'
                  : 'text-gray-400 hover:text-white'"
                @click="switchMode('login')"
              >Login</button>
              <button
                id="auth-tab-register"
                class="flex-1 py-2.5 px-4 rounded-lg text-sm font-bold transition-all duration-200"
                :class="mode === 'register'
                  ? 'bg-electric-violet-500 text-white shadow-glow-violet'
                  : 'text-gray-400 hover:text-white'"
                @click="switchMode('register')"
              >Daftar</button>
            </div>

            <!-- Forgot mode: back button -->
            <div v-else class="mb-4">
              <button type="button" class="flex items-center gap-1.5 text-sm text-gray-400 hover:text-cyan-magic-300 transition-colors" @click="switchMode('login')">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/></svg>
                Kembali ke Login
              </button>
            </div>

            <h2 class="text-2xl font-['Orbitron'] font-black text-white mb-1">
              {{ mode === 'login' ? 'Selamat Datang' : mode === 'register' ? 'Buat Akun' : 'Reset Password' }}
            </h2>
            <p class="text-sm text-gray-400 mb-6">
              {{ mode === 'login'
                ? 'Login untuk melacak histori dan profil scrim kamu.'
                : mode === 'register'
                ? 'Bergabunglah dengan komunitas GoLobby Scrim Arena.'
                : 'Verifikasi identitasmu dengan email & username, lalu set password baru.' }}
            </p>
          </div>

          <!-- Form: Login & Register -->
          <form v-if="mode !== 'forgot'" class="px-8 pb-8 space-y-4" @submit.prevent="handleSubmit">

            <!-- Username and WhatsApp (register only) -->
            <Transition name="field-slide">
              <div v-if="mode === 'register'" class="space-y-4">
                <div>
                  <label class="block text-xs font-bold text-gray-400 uppercase tracking-widest mb-1.5">Username</label>
                  <input
                    id="auth-input-username"
                    v-model="form.username"
                    type="text"
                    autocomplete="username"
                    placeholder="contoh: ProPlayer99"
                    class="auth-input"
                    :disabled="loading"
                    required
                  >
                </div>
                <div>
                  <label class="block text-xs font-bold text-gray-400 uppercase tracking-widest mb-1.5">No. WhatsApp</label>
                  <input
                    id="auth-input-whatsapp"
                    v-model="form.whatsapp"
                    type="tel"
                    placeholder="contoh: 62812345678"
                    class="auth-input"
                    :class="{'border-red-500/50 focus:border-red-500/70 focus:shadow-[0_0_0_3px_rgba(239,68,68,0.15)]': whatsappError}"
                    :disabled="loading"
                    required
                  >
                  <!-- Validation Error -->
                  <p v-if="whatsappError" class="text-xs font-bold text-red-400 mt-1.5 flex items-center gap-1">
                    {{ whatsappError }}
                  </p>
                </div>
              </div>
            </Transition>

            <!-- Email -->
            <div>
              <label class="block text-xs font-bold text-gray-400 uppercase tracking-widest mb-1.5">Email</label>
              <input
                id="auth-input-email"
                v-model="form.email"
                type="email"
                autocomplete="email"
                placeholder="kamu@email.com"
                class="auth-input"
                :disabled="loading"
                required
              >
            </div>

            <!-- Password -->
            <div>
              <label class="block text-xs font-bold text-gray-400 uppercase tracking-widest mb-1.5">Password</label>
              <div class="relative">
                <input
                  id="auth-input-password"
                  v-model="form.password"
                  :type="showPassword ? 'text' : 'password'"
                  autocomplete="current-password"
                  placeholder="••••••••"
                  class="auth-input pr-12"
                  :disabled="loading"
                  required
                >
                <!-- Eye toggle -->
                <button
                  type="button"
                  class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-200 transition-colors"
                  @click="showPassword = !showPassword"
                  tabindex="-1"
                >
                  <svg v-if="!showPassword" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                  </svg>
                  <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
                  </svg>
                </button>
              </div>
            </div>

            <!-- Lupa Password (login only) -->
            <div v-if="mode === 'login'" class="text-right -mt-2">
              <button
                type="button"
                class="text-xs text-gray-500 hover:text-cyan-magic-400 transition-colors"
                @click="switchMode('forgot')"
              >Lupa Password?</button>
            </div>

            <!-- Error / Success message -->
            <Transition name="msg-fade">
              <div v-if="message.text"
                class="flex items-start gap-2.5 p-3 rounded-xl text-sm font-medium"
                :class="message.type === 'error'
                  ? 'bg-red-900/40 border border-red-700/50 text-red-300'
                  : 'bg-emerald-900/40 border border-emerald-700/50 text-emerald-300'"
              >
                <svg v-if="message.type === 'error'" class="w-4 h-4 mt-0.5 shrink-0" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
                </svg>
                <svg v-else class="w-4 h-4 mt-0.5 shrink-0" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                </svg>
                {{ message.text }}
              </div>
            </Transition>

            <!-- Submit button -->
            <button
              id="auth-btn-submit"
              type="submit"
              class="w-full py-3.5 rounded-xl font-bold text-white bg-gradient-magic shadow-glow-violet hover:shadow-glow-cyan hover:scale-[1.02] active:scale-95 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed disabled:scale-100"
              :disabled="loading || !!whatsappError"
            >
              <span v-if="loading" class="flex items-center justify-center gap-2">
                <svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
                </svg>
                Memproses...
              </span>
              <span v-else>{{ mode === 'login' ? 'Masuk ke GoLobby' : 'Buat Akun Sekarang' }}</span>
            </button>

            <!-- Switch mode link -->
            <p class="text-center text-sm text-gray-500 pt-1">
              {{ mode === 'login' ? 'Belum punya akun?' : 'Sudah punya akun?' }}
              <button
                type="button"
                class="text-cyan-magic-400 hover:text-cyan-magic-300 font-bold ml-1 transition-colors"
                @click="switchMode(mode === 'login' ? 'register' : 'login')"
              >{{ mode === 'login' ? 'Daftar gratis' : 'Login di sini' }}</button>
            </p>
          </form>

          <!-- ── Form: Forgot Password ── -->
          <form v-else class="px-8 pb-8 space-y-4" @submit.prevent="handleForgotSubmit">
            <div>
              <label class="block text-xs font-bold text-gray-400 uppercase tracking-widest mb-1.5">Email Terdaftar</label>
              <input
                v-model="forgotForm.email"
                type="email"
                placeholder="email@kamu.com"
                class="auth-input"
                :disabled="forgotLoading"
                required
              >
            </div>

            <div>
              <label class="block text-xs font-bold text-gray-400 uppercase tracking-widest mb-1.5">Username</label>
              <input
                v-model="forgotForm.username"
                type="text"
                placeholder="Username akun kamu"
                class="auth-input"
                :disabled="forgotLoading"
                required
              >
            </div>

            <div>
              <label class="block text-xs font-bold text-gray-400 uppercase tracking-widest mb-1.5">Password Baru</label>
              <div class="relative">
                <input
                  v-model="forgotForm.newPassword"
                  :type="showForgotPw ? 'text' : 'password'"
                  placeholder="Min. 6 karakter"
                  class="auth-input pr-10"
                  :disabled="forgotLoading"
                  required
                >
                <button type="button" class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-200 transition-colors" @click="showForgotPw = !showForgotPw" tabindex="-1">
                  <svg v-if="!showForgotPw" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/></svg>
                  <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"/></svg>
                </button>
              </div>
              <p v-if="forgotForm.newPassword && forgotForm.newPassword.length < 6" class="text-xs text-red-400 mt-1">⚠️ Minimal 6 karakter</p>
            </div>

            <div>
              <label class="block text-xs font-bold text-gray-400 uppercase tracking-widest mb-1.5">Konfirmasi Password Baru</label>
              <input
                v-model="forgotForm.confirmPassword"
                :type="showForgotPw ? 'text' : 'password'"
                placeholder="Ulangi password baru"
                class="auth-input"
                :class="{'border-red-500/50': forgotForm.confirmPassword && forgotForm.confirmPassword !== forgotForm.newPassword}"
                :disabled="forgotLoading"
                required
              >
              <p v-if="forgotForm.confirmPassword && forgotForm.confirmPassword !== forgotForm.newPassword" class="text-xs text-red-400 mt-1">⚠️ Password tidak cocok</p>
            </div>

            <!-- Message -->
            <Transition name="msg-fade">
              <div v-if="forgotMessage.text"
                class="flex items-start gap-2.5 p-3 rounded-xl text-sm font-medium"
                :class="forgotMessage.type === 'error' ? 'bg-red-900/40 border border-red-700/50 text-red-300' : 'bg-emerald-900/40 border border-emerald-700/50 text-emerald-300'"
              >
                <svg v-if="forgotMessage.type === 'error'" class="w-4 h-4 mt-0.5 shrink-0" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd"/></svg>
                <svg v-else class="w-4 h-4 mt-0.5 shrink-0" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/></svg>
                {{ forgotMessage.text }}
              </div>
            </Transition>

            <button
              type="submit"
              class="w-full py-3.5 rounded-xl font-bold text-white bg-gradient-magic shadow-glow-violet hover:shadow-glow-cyan hover:scale-[1.02] active:scale-95 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
              :disabled="forgotLoading || forgotForm.newPassword !== forgotForm.confirmPassword || forgotForm.newPassword.length < 6"
            >
              <span v-if="forgotLoading" class="flex items-center justify-center gap-2">
                <svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/></svg>
                Memproses...
              </span>
              <span v-else>🔐 Reset Password</span>
            </button>
          </form>

          <!-- Close button -->
          <button
            id="auth-btn-close"
            class="absolute top-4 right-4 p-2 rounded-lg text-gray-500 hover:text-white hover:bg-white/10 transition-all"
            @click="$emit('close')"
            aria-label="Tutup"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { ref, reactive, watch, computed, onMounted, onUnmounted } from 'vue'
import { useAuth } from '../composables/useAuth'

// ── Props & Emits ─────────────────────────────────────────────────────────────
const props = defineProps({
  isVisible: { type: Boolean, default: false },
  initialMode: { type: String, default: 'login' }, // 'login' | 'register'
})
const emit = defineEmits(['close', 'success'])

// ── State ─────────────────────────────────────────────────────────────────────
const { login, register, resetPassword } = useAuth()
const mode         = ref(props.initialMode)
const loading      = ref(false)
const showPassword = ref(false)
const showForgotInfo = ref(false)
const message      = reactive({ text: '', type: '' }) // type: 'error' | 'success'
const form         = reactive({ username: '', email: '', password: '', whatsapp: '' })

// ── Forgot Password State ─────────────────────────────────────────────────────
const forgotLoading  = ref(false)
const showForgotPw   = ref(false)
const forgotMessage  = reactive({ text: '', type: '' })
const forgotForm     = reactive({ email: '', username: '', newPassword: '', confirmPassword: '' })

// ── Watchers ──────────────────────────────────────────────────────────────────
watch(() => props.isVisible, (val) => {
  if (val) {
    // Reset state when modal opens
    resetForm()
  }
})

// ── Lifecycle ─────────────────────────────────────────────────────────────────
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

// ── Computed ──────────────────────────────────────────────────────────────────
const whatsappError = computed(() => {
  if (mode.value !== 'register') return ''
  const wa = form.whatsapp.trim()
  if (!wa) return ''
  if (!wa.startsWith('628')) return '⚠️ Must start with 628'
  if (wa.length < 10) return '⚠️ Too short (min 10 digits)'
  if (wa.length > 15) return '⚠️ Too long (max 15 digits)'
  if (!/^[0-9]+$/.test(wa)) return '⚠️ Numbers only'
  return ''
})

// ── Methods ───────────────────────────────────────────────────────────────────
const switchMode = (newMode) => {
  mode.value = newMode
  resetForm()
}

const resetForm = () => {
  form.username  = ''
  form.email     = ''
  form.password  = ''
  form.whatsapp  = ''
  message.text   = ''
  message.type   = ''
  showPassword.value = false
  // Reset forgot form too
  forgotForm.email           = ''
  forgotForm.username        = ''
  forgotForm.newPassword     = ''
  forgotForm.confirmPassword = ''
  forgotMessage.text         = ''
  forgotMessage.type         = ''
  showForgotPw.value         = false
}

const setMessage = (text, type = 'error') => {
  message.text = text
  message.type = type
}

const handleSubmit = async () => {
  if (loading.value) return
  loading.value = true
  message.text  = ''

  try {
    if (mode.value === 'register') {
      await register({ username: form.username, email: form.email, password: form.password, whatsapp: form.whatsapp })
      setMessage('Akun berhasil dibuat! Silakan login.', 'success')
      // Auto-switch to login after 1.5s
      setTimeout(() => switchMode('login'), 1500)
    } else {
      const data = await login({ email: form.email, password: form.password })
      setMessage(`Selamat datang, ${data.username}! 🎮`, 'success')
      // Close modal and emit success
      setTimeout(() => {
        emit('success', data)
        emit('close')
      }, 800)
    }
  } catch (err) {
    setMessage(err.message)
  } finally {
    loading.value = false
  }
}

const handleForgotSubmit = async () => {
  if (forgotLoading.value) return
  forgotLoading.value = true
  forgotMessage.text = ''

  try {
    await resetPassword({
      email: forgotForm.email,
      username: forgotForm.username,
      newPassword: forgotForm.newPassword
    })
    forgotMessage.text = '🔓 Password berhasil direset! Silakan login dengan password baru.'
    forgotMessage.type = 'success'
    // Auto-switch to login after 2s
    setTimeout(() => switchMode('login'), 2000)
  } catch (err) {
    forgotMessage.text = err.message || 'Gagal mereset password'
    forgotMessage.type = 'error'
  } finally {
    forgotLoading.value = false
  }
}
</script>

<style>
.auth-input {
  @apply w-full px-4 py-3 rounded-xl text-sm text-white;
  @apply bg-white/5 border border-white/10;
  @apply placeholder-gray-600;
  @apply focus:outline-none focus:border-blue-500/70 focus:shadow-[0_0_0_3px_rgba(59,130,246,0.15)];
  @apply transition-all duration-200;
  @apply disabled:opacity-50;
}

/* Transitions */
.modal-fade-enter-active,
.modal-fade-leave-active { transition: opacity 0.25s ease; }
.modal-fade-enter-from,
.modal-fade-leave-to { opacity: 0; }

.field-slide-enter-active,
.field-slide-leave-active { transition: all 0.2s ease; overflow: hidden; }
.field-slide-enter-from { opacity: 0; max-height: 0; }
.field-slide-enter-to   { opacity: 1; max-height: 100px; }
.field-slide-leave-from { opacity: 1; max-height: 100px; }
.field-slide-leave-to   { opacity: 0; max-height: 0; }

.msg-fade-enter-active,
.msg-fade-leave-active { transition: all 0.2s ease; }
.msg-fade-enter-from,
.msg-fade-leave-to { opacity: 0; transform: translateY(-4px); }
</style>
