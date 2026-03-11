# Frontend / UI Architecture (Vue 3 + Tailwind)

## 1. Stack & Philosophy
Frontend **GoLobby** (*Matchmaking Client*) ditransformasikan dari kerangka prototipe ke arah Antarmuka e-Sports berdaya kompetitif tinggi menggunakan *Premium Dark Mode* + *Glassmorphism*.

- **Framework**: Vue.js 3 (Composition API / `<script setup>`)
- **Build Tool / Pipeline**: Vite
- **Styling**: TailwindCSS dengan konfigurasi *Utility-First* + Ekstensi Vanilla CSS murni
- **Host Port Docker**: `:5173`
- **Reverse Proxy**: NGINX (Static file servers & API routing layer)

## 2. Struktur Desain UI & Animasi

Penyesuaian visualisasi modern dengan standar industri Gaming:

### Tema Warna
Sebuah kustom palet yang diatur di Tailwind Configuration, menghindari "Warna Klasik generik" (`red`, `blue`, `green` plain):

- **Midnight**: `slate/zinc/gray` kustom untuk kanvas (Background) yang suram pekat.
- **Electric Violet** (`violet-500` - `#8b5cf6`): *Accent color* yang menghidupkan UI Mode Ranked Match (POKE).
- **Cyan Magic** (`cyan-400` - `#22d3ee`): Gradasi interaktif penyanding *Violet* yang memberi efek neon "Cyberpunk" terang menyala.
- **Antique Gold** (`#d4af37`): Gradasi Elegan eksklusif pendaftaran *Pro Scrim* (WARKOP).

### Asset CSS Kompleks (`main.css`)
Efek "Premium" dicapai dengan kombinasi `box-shadow` terang, *blur filter* (Glassmorphism card), dan Micro-Animasi *Keyframe*:

### Key Components

- **`DashboardView.vue`**: Master layouter / *Orchestrator*. Pengatur alur state pilihan (Pilih Kategori `POKE`/`WARKOP` -> Form Pengisian -> Matchmaking Spinner -> Modal/Notifikasi Toast Error -> Modal "Found").
- **`SearchingState.vue`**: UI "Pencarian Lawan".
  - Dibekali dengan Animasi **Radar Sweep Ping** dan elemen-elemen canggih seolah "Sistem mencari pemain real-time": *Loading Progress Bar Oscillate*, *Live Console Terminal Logs*, & *Stats Counters*.
  - Menghubungi Backend via plugin REST `axios`. Dan juga mendengarkan Sinyal WebSocket.
- **`MatchFoundModal.vue`**: Dialog "Alert" atau pop-up kritis saat sistem mengembalikan pertandingan berhasil. Memberi UI *Head-to-Head* / *Versus* screen (Kamu Vs Lawan). Dilengkapi Hitungan mundur animasi kritis progres-bar detik penentuan.

## 3. Composable `useScrimAPI.js`
Inilah abstraksi asinkron komunikasi klien dan server (*Hook*).

- Melindungi UI Vue dari baris Fetch yang membanjiri Component.
- Menangani *Fail-Safe Mechanism*. Saat Websocket Terputus, `useScrimAPI.js` secara mandiri menjadwalkan "Polling Cadangan" untuk mengirim `axios.get()` meminta info setiap tiga detiknya, menjaga aplikasi bebas crash.

## 4. Toast Notifications
Pendekatan Alert elegan non-intrusif pengganti `alert()` native DOM: Toast Notification.
Diletakkan absolut fixed pada `top-right`, membawa feedback warna kritis terhadap aksi Pengguna: Peringkat, Input tidak sah, atau "Match Canceled".

## 5. Deployment / Nginx Multi-stage Docker
Terdapat 2 layer build (Node.js & Nginx). Direktori `frontend/nginx.conf` me-*reverse-proxy* `location /` untuk menyalurkan asset SPA Vue (History Mode Fallback) sementara `location /api/` di *proxy-pass* secara mulus ke Backend Fiber port `3000`. `location /ws` me-*proxy-pass* koneksi websockets.
