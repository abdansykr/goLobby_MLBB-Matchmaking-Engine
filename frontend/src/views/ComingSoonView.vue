<template>
  <div class="min-h-[calc(100vh-4rem)] flex flex-col items-center justify-center px-4 py-16 relative overflow-hidden">
    <!-- Ambient background glows -->
    <div class="absolute top-1/4 left-1/4 w-96 h-96 bg-electric-violet-600/15 rounded-full blur-[120px] pointer-events-none"></div>
    <div class="absolute bottom-1/4 right-1/4 w-80 h-80 bg-cyan-magic-500/10 rounded-full blur-[100px] pointer-events-none"></div>

    <!-- Floating particles -->
    <div class="absolute inset-0 overflow-hidden pointer-events-none">
      <div v-for="p in particles" :key="p.id"
        class="absolute w-1 h-1 rounded-full opacity-40 animate-float"
        :style="{ left: p.x + '%', top: p.y + '%', backgroundColor: p.color, animationDelay: p.delay + 's', animationDuration: p.duration + 's' }"
      ></div>
    </div>

    <!-- Main card -->
    <div class="glass-card p-8 sm:p-12 max-w-lg w-full text-center relative border border-electric-violet-500/20 shadow-[0_0_60px_rgba(139,92,246,0.15)]">

      <!-- Top badge -->
      <div class="inline-flex items-center gap-2 bg-electric-violet-500/10 border border-electric-violet-500/30 rounded-full px-4 py-1.5 mb-8">
        <span class="relative flex h-2 w-2">
          <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-electric-violet-400 opacity-75"></span>
          <span class="relative inline-flex rounded-full h-2 w-2 bg-electric-violet-500"></span>
        </span>
        <span class="text-xs font-bold text-electric-violet-300 uppercase tracking-widest">Dalam Pengembangan</span>
      </div>

      <!-- Mascot image -->
      <div class="relative w-44 h-44 sm:w-56 sm:h-56 mx-auto mb-8">
        <!-- Glow ring -->
        <div class="absolute inset-0 rounded-full bg-gradient-to-br from-electric-violet-500/30 to-cyan-magic-500/30 blur-2xl animate-pulse-slow"></div>
        <!-- Inner ring -->
        <div class="absolute inset-2 rounded-full border-2 border-dashed border-electric-violet-500/30 animate-spin-slow"></div>
        <!-- Image wrapper -->
        <div class="relative h-full w-full flex items-center justify-center">
          <img
            src="/coming_soon_cat.png"
            alt="Coming Soon Mascot"
            class="w-36 h-36 sm:w-44 sm:h-44 object-contain transform hover:-translate-y-2 transition-transform duration-500 drop-shadow-2xl animate-float"
          >
        </div>
      </div>

      <!-- Title -->
      <h2 class="text-3xl sm:text-4xl font-['Orbitron'] font-black text-white mb-3 leading-tight">
        Ditunggu ya<span class="gradient-text">...</span>
      </h2>

      <!-- Subtitle -->
      <p class="text-gray-400 text-sm sm:text-base leading-relaxed mb-8">
        Fitur ini sedang dalam tahap pengembangan aktif.<br>
        <span class="text-electric-violet-300 font-semibold">Kami sedang menyiapkan sesuatu yang spesial</span> untukmu! 🚀
      </p>

      <!-- Feature chips -->
      <div class="flex flex-wrap justify-center gap-2 mb-10">
        <span v-for="f in features" :key="f"
          class="px-3 py-1 rounded-full bg-white/5 border border-white/10 text-xs text-gray-400 font-medium"
        >{{ f }}</span>
      </div>

      <!-- CTA Buttons -->
      <div class="flex flex-col sm:flex-row gap-3 justify-center">
        <router-link
          to="/"
          class="flex items-center justify-center gap-2 px-6 py-3 rounded-xl font-bold text-white bg-gradient-magic shadow-glow-violet hover:shadow-glow-cyan hover:scale-105 active:scale-95 transition-all duration-200"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
          Kembali ke Dashboard
        </router-link>

        <a
          href="https://saweria.co/golobby"
          target="_blank"
          rel="noopener noreferrer"
          class="flex items-center justify-center gap-2 px-6 py-3 rounded-xl font-bold text-yellow-400 glass border border-yellow-400/20 hover:bg-yellow-400/10 hover:border-yellow-400/50 hover:scale-105 active:scale-95 transition-all duration-200"
        >
          <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
            <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z" />
          </svg>
          Dukung GoLobby
        </a>
      </div>
    </div>

    <!-- Bottom text -->
    <p class="mt-8 text-center text-xs text-gray-600">
      GoLobby v2.0 — MLBB Scrim Matchmaking Engine
    </p>
  </div>
</template>

<script setup>
// Generate floating particles
const particles = Array.from({ length: 20 }, (_, i) => ({
  id: i,
  x: Math.random() * 100,
  y: Math.random() * 100,
  color: ['#8b5cf6', '#06b6d4', '#d4af37'][Math.floor(Math.random() * 3)],
  delay: Math.random() * 4,
  duration: 3 + Math.random() * 4,
}))

const features = [
  '🔐 Login Akun', '📊 Histori Match', '🏆 Leaderboard',
  '📱 Push Notif', '🎮 Lobby Voice', '⭐ Rating Tim',
]
</script>
