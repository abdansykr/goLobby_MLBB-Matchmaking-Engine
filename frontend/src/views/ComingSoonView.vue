<template>
  <div class="flex flex-col items-center justify-center min-h-[calc(100vh-5rem)] px-4 relative overflow-hidden">
    <!-- Background glow for extra aesthetic -->
    <div class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[500px] h-[500px] bg-electric-violet-600/20 rounded-full blur-[100px] -z-10 pointer-events-none"></div>

    <!-- Image -->
    <div class="mb-8 transform hover:-translate-y-2 transition-transform duration-500 drop-shadow-2xl">
      <img src="/coming_soon_cat.png" alt="Coming Soon Mascot" class="w-64 md:w-80 h-auto object-contain">
    </div>
    
    <!-- Title -->
    <h2 class="text-3xl md:text-5xl font-['Orbitron'] font-black text-white mb-4 text-center">
      Ditunggu yaa....
    </h2>
    
    <!-- Description -->
    <p class="text-gray-300 text-center max-w-md mb-10 leading-relaxed text-sm md:text-base">
      Fitur ini sedang dalam tahap pengembangan. Kami sedang menyiapkan sesuatu yang spesial untuk Anda..
    </p>
    
    <!-- Back Button -->
    <router-link to="/" class="px-10 py-3.5 rounded-xl font-bold text-white bg-gradient-to-r from-pink-500 to-electric-violet-600 hover:from-pink-400 hover:to-electric-violet-500 shadow-lg hover:shadow-glow-violet transition-all duration-300 flex items-center gap-2">
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
      </svg>
      Kembali ke Dashboard
    </router-link>
  </div>
</template>

<script setup>
// Coming soon view component
</script>
