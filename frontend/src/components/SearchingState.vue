<template>
  <div class="searching-wrapper py-6">

    <!-- ── RADAR ARENA ─────────────────────────────────────────── -->
    <div class="radar-arena mx-auto mb-8">
      <!-- Outer pulsing rings -->
      <div class="ring ring-1"></div>
      <div class="ring ring-2"></div>
      <div class="ring ring-3"></div>
      <div class="ring ring-4"></div>

      <!-- Sweep beam -->
      <div class="sweep-beam" :style="{ background: sweepColor }"></div>

      <!-- Ping dots (simulated players found) -->
      <div
        v-for="dot in pingDots"
        :key="dot.id"
        class="ping-dot"
        :style="{ top: dot.y + '%', left: dot.x + '%', background: dot.color }"
      ></div>

      <!-- Center core -->
      <div class="core-outer" :style="{ boxShadow: `0 0 30px 8px ${coreGlow}` }"></div>
      <div class="core-inner">
        <!-- Category icon -->
        <svg v-if="formData.category === 'POKE'" class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
        </svg>
        <svg v-else class="w-8 h-8 text-midnight-900" fill="currentColor" viewBox="0 0 20 20">
          <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
        </svg>
      </div>
    </div>

    <!-- ── STATUS TEXT ─────────────────────────────────────────── -->
    <div class="text-center space-y-2 mb-8">
      <h3 class="text-2xl md:text-3xl font-bold font-['Orbitron'] tracking-wide"
          :class="isPoke ? 'text-glow' : 'text-glow-gold'">
        {{ statusMessages[statusIdx] }}
      </h3>
      <p class="text-gray-400 text-base">
        Mode: <span class="font-bold" :class="isPoke ? 'text-cyan-magic-300' : 'text-antique-gold-300'">
          {{ isPoke ? 'Ranked Solo' : 'Pro Scrim' }}
        </span>
        <span v-if="isPoke" class="text-gray-500"> · Rank: </span>
        <span v-if="isPoke" class="font-bold text-electric-violet-400">{{ formData.rankName }}</span>
      </p>
    </div>

    <!-- ── STATS ROW ───────────────────────────────────────────── -->
    <div class="stats-row mb-8">
      <!-- Elapsed time -->
      <div class="stat-box" :class="isPoke ? 'border-cyan-magic-500/30' : 'border-antique-gold-400/30'">
        <div class="stat-icon" :class="isPoke ? 'bg-cyan-magic-500/10' : 'bg-antique-gold-400/10'">
          <svg class="w-5 h-5" :class="isPoke ? 'text-cyan-magic-400' : 'text-antique-gold-400'"
               fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
        </div>
        <div>
          <div class="stat-value font-['Orbitron']"
               :class="isPoke ? 'text-cyan-magic-300' : 'text-antique-gold-300'">
            {{ elapsedDisplay }}
          </div>
          <div class="stat-label">Elapsed</div>
        </div>
      </div>

      <!-- Players scanned -->
      <div class="stat-box" :class="isPoke ? 'border-electric-violet-500/30' : 'border-antique-gold-400/20'">
        <div class="stat-icon" :class="isPoke ? 'bg-electric-violet-500/10' : 'bg-antique-gold-400/10'">
          <svg class="w-5 h-5" :class="isPoke ? 'text-electric-violet-400' : 'text-antique-gold-400'"
               fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z"/>
          </svg>
        </div>
        <div>
          <div class="stat-value font-['Orbitron']"
               :class="isPoke ? 'text-electric-violet-400' : 'text-antique-gold-400'">
            {{ teamsScanned }}
          </div>
          <div class="stat-label">Teams Scanned</div>
        </div>
      </div>

      <!-- Connection status -->
      <div class="stat-box border-green-500/30">
        <div class="stat-icon bg-green-500/10">
          <svg class="w-5 h-5 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M8.111 16.404a5.5 5.5 0 017.778 0M12 20h.01m-7.08-7.071c3.904-3.905 10.236-3.905 14.141 0M1.394 9.393c5.857-5.857 15.355-5.857 21.213 0"/>
          </svg>
        </div>
        <div>
          <div class="stat-value text-green-400 font-['Orbitron'] animate-pulse">LIVE</div>
          <div class="stat-label">Real-time</div>
        </div>
      </div>
    </div>

    <!-- ── PROGRESS BAR ────────────────────────────────────────── -->
    <div class="progress-track mb-8">
      <div class="progress-fill"
           :class="isPoke ? 'from-electric-violet-500 via-cyan-magic-400 to-electric-violet-500' : 'from-antique-gold-500 via-yellow-300 to-antique-gold-500'"
           :style="{ width: progressWidth + '%' }">
        <div class="progress-shimmer"></div>
      </div>
    </div>

    <!-- ── ACTIVITY LOG ────────────────────────────────────────── -->
    <div class="activity-log mb-8">
      <div class="activity-header">
        <span class="activity-dot"></span>
        <span class="text-xs text-gray-500 font-mono tracking-wider">MATCHMAKING LOG</span>
      </div>
      <div class="activity-entries">
        <transition-group name="log-slide">
          <div
            v-for="entry in activityLog"
            :key="entry.id"
            class="log-entry"
          >
            <span class="log-time">{{ entry.time }}</span>
            <span class="log-text" :class="entry.color">{{ entry.msg }}</span>
          </div>
        </transition-group>
      </div>
    </div>

    <!-- ── CANCEL BUTTON ───────────────────────────────────────── -->
    <div class="text-center">
      <button
        @click="$emit('cancel')"
        class="cancel-btn group"
      >
        <svg class="w-4 h-4 group-hover:rotate-90 transition-transform duration-300"
             fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
        </svg>
        Cancel Search
      </button>
    </div>

  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  formData: { type: Object, required: true }
})
defineEmits(['cancel'])

const isPoke = computed(() => props.formData.category === 'POKE')
const sweepColor = computed(() =>
  isPoke.value
    ? 'conic-gradient(from 0deg, transparent 70%, rgba(6,182,212,0.6) 100%)'
    : 'conic-gradient(from 0deg, transparent 70%, rgba(212,175,55,0.6) 100%)'
)
const coreGlow = computed(() =>
  isPoke.value ? 'rgba(6,182,212,0.5)' : 'rgba(212,175,55,0.5)'
)

// ── Cycling status messages ────────────────────────────────────
const statusMessages = [
  'Scanning the arena...',
  'Calibrating rank filters...',
  'Connecting to opponents...',
  'Analyzing team rosters...',
  'Almost there...',
]
const statusIdx = ref(0)

// ── Elapsed timer ──────────────────────────────────────────────
const elapsed = ref(0)
const elapsedDisplay = computed(() => {
  const m = Math.floor(elapsed.value / 60)
  const s = elapsed.value % 60
  return m > 0 ? `${m}m ${s}s` : `${s}s`
})

// ── Simulated teams scanned counter ───────────────────────────
const teamsScanned = ref(0)

// ── Progress bar (oscillates) ─────────────────────────────────
const progressWidth = ref(10)
let progressDir = 1

// ── Ping dots on radar ────────────────────────────────────────
const pingDots = ref([])
let dotIdCounter = 0
const addPingDot = () => {
  const angle = Math.random() * 360
  const r = 20 + Math.random() * 35
  const x = 50 + r * Math.cos((angle * Math.PI) / 180)
  const y = 50 + r * Math.sin((angle * Math.PI) / 180)
  const colors = isPoke.value
    ? ['#06b6d4', '#818cf8', '#22d3ee']
    : ['#d4af37', '#fbbf24', '#f59e0b']
  pingDots.value.push({
    id: dotIdCounter++,
    x: Math.max(5, Math.min(90, x)),
    y: Math.max(5, Math.min(90, y)),
    color: colors[Math.floor(Math.random() * colors.length)]
  })
  if (pingDots.value.length > 6) pingDots.value.shift()
}

// ── Activity log ──────────────────────────────────────────────
const activityLog = ref([])
let logId = 0
const logMessages = [
  { msg: '🔍 Mulai mencari lawan...', color: 'text-gray-300' },
  { msg: '📡 Berhasil masuk ke server GoLobby', color: 'text-green-400' },
  { msg: '⚡ Koneksi aman dan stabil', color: 'text-cyan-magic-300' },
  { msg: '🎯 Mencari musuh dengan Rank yang setara...', color: 'text-electric-violet-300' },
  { msg: '🔄 Memperluas radar pencarian', color: 'text-gray-300' },
  { msg: '👥 Mendapatkan beberapa tim yang potensial', color: 'text-yellow-300' },
  { msg: '⚙️ Memeriksa selisih / batas Rank...', color: 'text-gray-300' },
  { msg: '🏆 Mengatur keadilan permainan', color: 'text-green-400' },
  { msg: '📊 Mempersiapkan area pertandingan...', color: 'text-cyan-magic-300' },
  { msg: '🌐 Sedang dalam antrian tunggu...', color: 'text-gray-300' },
]
let logMsgIdx = 0
const pushLog = () => {
  const now = new Date()
  const ts = `${now.getHours().toString().padStart(2, '0')}:${now.getMinutes().toString().padStart(2, '0')}:${now.getSeconds().toString().padStart(2, '0')}`
  const entry = logMessages[logMsgIdx % logMessages.length]
  activityLog.value.unshift({ id: logId++, time: ts, msg: entry.msg, color: entry.color })
  if (activityLog.value.length > 4) activityLog.value.pop()
  logMsgIdx++
}

// ── Timers ────────────────────────────────────────────────────
let elapsedTimer, dotTimer, progressTimer, statusTimer, scanTimer, logTimer

onMounted(() => {
  pushLog()

  elapsedTimer = setInterval(() => { elapsed.value++ }, 1000)

  dotTimer = setInterval(addPingDot, 1200)

  progressTimer = setInterval(() => {
    progressWidth.value += progressDir * (2 + Math.random() * 4)
    if (progressWidth.value >= 90) progressDir = -1
    if (progressWidth.value <= 10) progressDir = 1
  }, 400)

  statusTimer = setInterval(() => {
    statusIdx.value = (statusIdx.value + 1) % statusMessages.length
  }, 3000)

  scanTimer = setInterval(() => {
    teamsScanned.value += Math.floor(1 + Math.random() * 3)
  }, 1500)

  logTimer = setInterval(pushLog, 4000)
})

onUnmounted(() => {
  clearInterval(elapsedTimer)
  clearInterval(dotTimer)
  clearInterval(progressTimer)
  clearInterval(statusTimer)
  clearInterval(scanTimer)
  clearInterval(logTimer)
})
</script>

<style scoped>
/* ── Wrapper ─────────────────────────────────────────────────── */
.searching-wrapper {
  animation: fadeInUp 0.5s ease;
}
@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(20px); }
  to   { opacity: 1; transform: translateY(0); }
}

/* ── Radar ───────────────────────────────────────────────────── */
.radar-arena {
  position: relative;
  width: 240px;
  height: 240px;
}

.ring {
  position: absolute;
  border-radius: 50%;
  border: 1px solid;
  animation: radarPulse 2.4s ease-out infinite;
}
.ring-1 { inset: 0;   border-color: rgba(6,182,212,0.5);  animation-delay: 0s; }
.ring-2 { inset: 20%; border-color: rgba(6,182,212,0.35); animation-delay: 0.6s; }
.ring-3 { inset: 36%; border-color: rgba(6,182,212,0.25); animation-delay: 1.2s; }
.ring-4 { inset: 50%; border-color: rgba(6,182,212,0.15); animation-delay: 1.8s; }

@keyframes radarPulse {
  0%   { opacity: 1; transform: scale(1); }
  100% { opacity: 0; transform: scale(1.06); }
}

.sweep-beam {
  position: absolute;
  inset: 0;
  border-radius: 50%;
  animation: radarSweep 2.4s linear infinite;
  transform-origin: center;
}
@keyframes radarSweep {
  from { transform: rotate(0deg); }
  to   { transform: rotate(360deg); }
}

.ping-dot {
  position: absolute;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  transform: translate(-50%, -50%);
  animation: dotBlink 1.5s ease-out forwards;
  box-shadow: 0 0 8px currentColor;
}
@keyframes dotBlink {
  0%   { opacity: 0; transform: translate(-50%, -50%) scale(0.5); }
  30%  { opacity: 1; transform: translate(-50%, -50%) scale(1.3); }
  100% { opacity: 0; transform: translate(-50%, -50%) scale(1); }
}

.core-outer {
  position: absolute;
  inset: 42%;
  border-radius: 50%;
  background: rgba(6,182,212,0.15);
  animation: corePulse 1.5s ease-in-out infinite;
}
.core-inner {
  position: absolute;
  inset: 44%;
  border-radius: 50%;
  background: linear-gradient(135deg, #0891b2, #6d28d9);
  display: flex;
  align-items: center;
  justify-content: center;
}
@keyframes corePulse {
  0%, 100% { transform: scale(1); }
  50%       { transform: scale(1.15); }
}

/* ── Stats row ───────────────────────────────────────────────── */
.stats-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  max-width: 480px;
  margin-left: auto;
  margin-right: auto;
}
.stat-box {
  display: flex;
  align-items: center;
  gap: 10px;
  background: rgba(0,0,0,0.3);
  border: 1px solid;
  border-radius: 12px;
  padding: 12px;
}
.stat-icon {
  padding: 8px;
  border-radius: 8px;
  flex-shrink: 0;
}
.stat-value {
  font-size: 1rem;
  font-weight: 900;
  line-height: 1;
}
.stat-label {
  font-size: 0.6rem;
  color: #6b7280;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  margin-top: 2px;
}

/* ── Progress bar ────────────────────────────────────────────── */
.progress-track {
  max-width: 440px;
  margin-left: auto;
  margin-right: auto;
  height: 6px;
  background: rgba(255,255,255,0.05);
  border-radius: 999px;
  overflow: hidden;
}
.progress-fill {
  height: 100%;
  border-radius: 999px;
  background-image: linear-gradient(90deg, var(--tw-gradient-stops));
  background-size: 200% 100%;
  position: relative;
  transition: width 0.4s ease;
  animation: progressShift 2s linear infinite;
}
@keyframes progressShift {
  from { background-position: 0% 50%; }
  to   { background-position: 200% 50%; }
}

/* ── Activity log ────────────────────────────────────────────── */
.activity-log {
  max-width: 440px;
  margin-left: auto;
  margin-right: auto;
  background: rgba(0,0,0,0.4);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 12px;
  padding: 12px 16px;
  font-family: 'Courier New', monospace;
}
.activity-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
}
.activity-dot {
  width: 6px;
  height: 6px;
  background: #22c55e;
  border-radius: 50%;
  animation: pulse 1.5s infinite;
}
.activity-entries {
  display: flex;
  flex-direction: column;
  gap: 5px;
  min-height: 80px;
}
.log-entry {
  display: flex;
  gap: 10px;
  font-size: 0.72rem;
  line-height: 1.4;
}
.log-time {
  color: #4b5563;
  flex-shrink: 0;
  font-size: 0.65rem;
  padding-top: 1px;
}

/* ── Log transitions ─────────────────────────────────────────── */
.log-slide-enter-active { transition: all 0.3s ease; }
.log-slide-enter-from   { opacity: 0; transform: translateY(-8px); }
.log-slide-leave-active { transition: all 0.2s ease; }
.log-slide-leave-to     { opacity: 0; }

/* ── Cancel button ───────────────────────────────────────────── */
.cancel-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 10px 28px;
  background: rgba(239,68,68,0.08);
  border: 1px solid rgba(239,68,68,0.3);
  border-radius: 10px;
  color: #f87171;
  font-size: 0.875rem;
  font-weight: 600;
  transition: all 0.2s;
}
.cancel-btn:hover {
  background: rgba(239,68,68,0.18);
  border-color: rgba(239,68,68,0.6);
  box-shadow: 0 0 16px rgba(239,68,68,0.2);
}
</style>
