# 🎨 GoLobby Frontend - Dark Fantasy MLBB-Style UI

## 🎯 Overview

A **premium dark fantasy gaming interface** inspired by Mobile Legends: Bang Bang, built with **Vue 3 + Vite + Tailwind CSS**.

### ✨ Key Features

- 🌙 **Dark Fantasy Theme** - Deep slate backgrounds with electric violet and cyan accents
- 💎 **Glassmorphism** - Frosted glass cards with backdrop blur
- ✨ **Glow Effects** - Magic-style outer glows on buttons and borders
- 🔷 **Hexagonal Design** - MLBB-inspired geometric shapes
- ⚡ **Micro-interactions** - Smooth hover, click, and transition effects
- 📱 **Fully Responsive** - Mobile-first design with perfect stacking
- 🎭 **Dramatic VS Screen** - High-impact match found modal
- 🎮 **HUD-Style Navbar** - Gaming heads-up display aesthetic

---

## 📦 Project Structure

```
frontend/
├── src/
│   ├── components/
│   │   ├── NavbarHUD.vue           # Top navigation bar with HUD style
│   │   └── MatchFoundModal.vue     # Dramatic VS screen modal
│   ├── views/
│   │   └── DashboardView.vue       # Main dashboard with stats
│   ├── assets/
│   │   └── main.css                # Custom CSS with Tailwind
│   ├── composables/                # Vue composables (for API calls)
│   ├── App.vue                     # Root component
│   └── main.js                     # Entry point
├── public/                         # Static assets
├── index.html                      # HTML entry
├── vite.config.js                  # Vite configuration
├── tailwind.config.js              # Tailwind theme config
├── postcss.config.js               # PostCSS config
└── package.json                    # Dependencies
```

---

## 🎨 Design System

### Color Palette

#### Base Colors (Midnight Slate)
```javascript
'midnight-900': '#0a0e27'  // Darkest background
'midnight-800': '#131829'  // Card backgrounds
'midnight-700': '#1a1f3a'  // Elevated surfaces
```

#### Accent Colors (Magic/Cyber)
```javascript
'electric-violet-500': '#8b5cf6'  // Primary magic glow
'cyan-magic-400': '#22d3ee'        // Secondary cyber glow
'antique-gold-400': '#e5c158'      // Prestige/rank accent
```

#### Rank Colors
```javascript
'mythic-glow': '#ff6b00'     // Mythic Glory
'legend-purple': '#9333ea'   // Legend
'epic-blue': '#3b82f6'       // Epic
```

### Typography

**Headers:** `Orbitron` (Bold, futuristic)  
**Body:** `Rajdhani` (Clean, readable)

### Visual Effects

#### Glassmorphism
```css
.glass-card {
  background: rgba(255,255,255,0.1);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(255,255,255,0.2);
}
```

#### Glow Shadows
```css
.shadow-glow-violet {
  box-shadow: 
    0 0 20px rgba(139, 92, 246, 0.5),
    0 0 40px rgba(139, 92, 246, 0.2);
}
```

#### Animations
- `pulse-slow` - Slow breathing effect (3s)
- `glow` - Alternating glow intensity
- `radar` - Rotating radar scan (3s)
- `float` - Gentle floating motion
- `shimmer` - Loading shimmer effect

---

## 🚀 Installation & Setup

### Prerequisites

node.js >= 18.0  
npm >= 9.0
vite >= 5.0

### Step 1: Install Dependencies

```bash
cd frontend
npm install
```

This will install:
- Vue 3.4
- Vite 5.0
- Tailwind CSS 3.4
- Vue Router 4.2
- Axios (for API calls)

### Step 2: Development Server

```bash
npm run dev
```

Server will start on: `http://localhost:5173`

### Step 3: Build for Production

```bash
npm run build
```

Output: `dist/` folder

### Step 4: Preview Production Build

```bash
npm run preview
```

---

## 🎮 Component Guide

### 1. NavbarHUD.vue

**Purpose:** Gaming-style navigation bar

**Features:**
- Hexagonal logo with crystal effect
- User stats display (Rank, Reputation)
- Avatar with hexagon frame
- Notification bell with badge
- Responsive design

**Props:**
```javascript
{
  userName: String,        // Display name
  userAvatar: String,      // Avatar URL
  userStats: Object,       // {rank, reputation}
  hasNotifications: Boolean
}
```

**Usage:**
```vue
<NavbarHUD 
  userName="Captain Alpha"
  :userStats="{ rank: 'Mythic', reputation: 95 }"
  :hasNotifications="true"
/>
```

---

### 2. DashboardView.vue

**Purpose:** Main hub for matchmaking

**Sections:**
1. **Hero Banner** - Gradient title with tagline
2. **Stat Cards (3):**
   - Average Rank (hexagon icon with Mythic badge)
   - Reputation Score (circular progress bar)
   - Today's Matches (win/loss stats)
3. **Main Action:**
   - Portal animation when idle
   - Radar animation when searching
4. **Recent Matches** - History list with win/loss indicators

**States:**
- `isSearching: false` → Shows portal button
- `isSearching: true` → Shows radar animation

**Interactions:**
```javascript
startSearch()  // Begin matchmaking
cancelSearch() // Cancel matchmaking
```

---

### 3. MatchFoundModal.vue

**Purpose:** Dramatic VS screen when match is found

**Features:**
- 60-second countdown timer
- Team vs Team display with hexagon avatars
- Rank badges on avatars
- Match details grid (Category, Rank Diff, ID)
- Accept/Decline buttons
- Auto-timeout handling

**Props:**
```javascript
{
  isVisible: Boolean,
  yourTeam: {
    name: String,
    avatar: String,
    rank: String
  },
  opponentTeam: {
    name: String,
    avatar: String,
    rank: String
  },
  matchDetails: {
    category: String,  // POKE/WARKOP
    rankDiff: Number,
    id: String
  },
  timeout: Number       // Default: 60 seconds
}
```

**Events:**
```javascript
@accept   // User accepts match
@decline  // User declines
@timeout  // Timer runs out
```

**Usage:**
```vue
<MatchFoundModal 
  :isVisible="showModal"
  :yourTeam="team1"
  :opponentTeam="team2"
  :matchDetails="details"
  @accept="handleAccept"
  @decline="handleDecline"
/>
```

---

## 🎨 Tailwind Custom Classes

### Buttons

```html
<!-- Magic Button (Primary) -->
<button class="btn-magic">FIND MATCH</button>

<!-- Gold Button (Accept) -->
<button class="btn-gold">ACCEPT MATCH</button>

<!-- Danger Button (Decline) -->
<button class="btn-danger">Decline</button>
```

### Cards

```html
<!-- Glass Card -->
<div class="glass-card p-6">Content</div>

<!-- Glass Card with Hover -->
<div class="glass-card-hover p-6">Interactive</div>

<!-- Stat Card -->
<div class="stat-card">Stats</div>
```

### Effects

```html
<!-- Text Glow -->
<h1 class="text-glow">Glowing Text</h1>

<!-- Energy Border -->
<div class="energy-border">Hover me</div>

<!-- Portal Pulse -->
<div class="portal-pulse">Pulsing element</div>

<!-- Gradient Text -->
<span class="gradient-text">Magic Text</span>
```

### Hexagons

```html
<!-- Hexagon Shape -->
<div class="hexagon-border bg-gradient-magic">
  <!-- Content -->
</div>

<!-- Avatar Hexagon -->
<div class="avatar-hexagon">
  <img src="avatar.jpg" alt="Avatar">
</div>
```

---

## 🔌 API Integration (TODO)

### Composable: useScrimAPI.js

```javascript
// frontend/src/composables/useScrimAPI.js
import { ref } from 'vue'
import axios from 'axios'

const API_URL = '/api/scrim'

export function useScrimAPI() {
  const loading = ref(false)
  const error = ref(null)

  const createRequest = async (teamData) => {
    loading.value = true
    try {
      const response = await axios.post(`${API_URL}/request`, teamData)
      return response.data
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const getRequestStatus = async (requestId) => {
    try {
      const response = await axios.get(`${API_URL}/request/${requestId}`)
      return response.data
    } catch (err) {
      error.value = err.message
      throw err
    }
  }

  const cancelRequest = async (requestId) => {
    try {
      const response = await axios.post(`${API_URL}/request/${requestId}/cancel`)
      return response.data
    } catch (err) {
      error.value = err.message
      throw err
    }
  }

  return {
    loading,
    error,
    createRequest,
    getRequestStatus,
    cancelRequest
  }
}
```

### Usage in Component:

```vue
<script setup>
import { useScrimAPI } from '@/composables/useScrimAPI'

const { createRequest, loading } = useScrimAPI()

const startSearch = async () => {
  const result = await createRequest({
    team_name: 'Team Alpha',
    whatsapp_number: '628123456789',
    category: 'POKE',
    rank_weight: 5
  })
  console.log('Request created:', result.request_id)
}
</script>
```

---

## 🎭 Responsive Design

### Breakpoints (Tailwind)

- `sm`: 640px  
- `md`: 768px  
- `lg`: 1024px  
- `xl`: 1280px

### Mobile Optimizations

1. **Navbar:**
   - Hides center stats on mobile
   - Collapses avatar text on small screens

2. **Dashboard:**
   - Stats grid: 3 columns → 1 column stack
   - Portal size: Responsive with max-width

3. **Match Modal:**
   - VS screen: 3 columns → 1 column vertical
   - Avatars scale down on mobile
   - Buttons stack vertically

---

## 🚀 Deployment

### Build & Deploy

```bash
# Build production
npm run build

# Preview locally
npm run preview

# Deploy to static host
# Upload dist/ folder to:
# - Vercel
# - Netlify
# - GitHub Pages
# - Any CDN
```

### Environment Variables

Create `.env.production`:

```env
VITE_API_URL=https://api.golobby.com
VITE_WS_URL=wss://api.golobby.com/ws
```

Access in code:

```javascript
const apiUrl = import.meta.env.VITE_API_URL
```

---

## 📝 Checklist for Integration

### Backend Integration

- [ ] Create `useScrimAPI.js` composable
- [ ] Connect `startSearch()` to API
- [ ] Poll `/request/:id` for status
- [ ] Show MatchFoundModal on match
- [ ] Implement WebSocket for real-time updates
- [ ] Handle Accept/Decline confirmation API calls
- [ ] Navigate to Lobby after accept

### Additional Features

- [ ] Create LobbyView.vue (post-match chat)
- [ ] Add loading skeletons
- [ ] Implement toast notifications
- [ ] Add sound effects (optional)
- [ ] Create settings page
- [ ] Add match history view

---

## 🎨 Design Philosophy

### "The Human Touch"

This UI is designed to feel **alive** and **premium**, not sterile or generic:

1. **Texture:** Glassmorphism adds depth
2. **Light:** Glows simulate magical energy
3. **Motion:** Animations feel organic (float, pulse)
4. **Color:** Gradients create richness
5. **Shape:** Hexagons add geometric interest
6. **Typography:** Bold headers command attention

### MLBB Inspiration

- **HUD aesthetic** in navbar
- **Rank badges** with metallic feel
- **VS screen drama** like hero select
- **Portal/magic circles** for actions
- **Dark map** background gradient

---

## 🐛 Troubleshooting

### Tailwind Not Working

```bash
# Clear cache and rebuild
rm -rf node_modules .vite dist
npm install
npm run dev
```

### Fonts Not Loading

Check `main.css` Google Fonts import:
```css
@import url('https://fonts.googleapis.com/css2?family=Orbitron...')
```

### Animations Janky

Ensure hardware acceleration:
```css
.animated-element {
  will-change: transform;
  transform: translateZ(0);
}
```

---

## 📚 Resources

- **Vue 3 Docs:** https://vuejs.org
- **Tailwind CSS:** https://tailwindcss.com
- **Vite:** https://vitejs.dev
- **Design Inspiration:** Mobile Legends, League of Legends Wild Rift

---

**🎮 Ready to dominate the arena!** 🚀
