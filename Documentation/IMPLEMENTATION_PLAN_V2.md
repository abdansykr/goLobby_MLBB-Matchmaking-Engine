# 🚀 GoLobby Scrim Matchmaking v2.0 - Implementation Plan

## 📋 Overview

Pengembangan major update untuk GoLobby Matchmaking dengan:
- **Category-based matching**: POKE (1-8) vs WARKOP (9-10)
- **WhatsApp Integration**: Auto-generate WA link saat match
- **Vue.js 3 Dark Gaming UI**: Premium dark fantasy aesthetic
- **Rate Limiting**: 1 IP = 1 active request
- **Auto-expiry**: 30 menit timeout

---

## ✅ Progress Tracking

### TAHAP 1: Database & Schema ✅
- [x] Migration 000002: `scrim_requests` table
- [x] Migration 000002: `scrim_matches` table
- [x] Enum types: `scrim_category`, `scrim_status`
- [x] Auto-expiry functions
- [x] Indexes untuk performance

### TAHAP 2: Domain Layer ✅
- [x] ScrimRequest entity
- [x] ScrimMatch entity
- [x] Validation logic (CanMatchWith, IsValidRankWeight)
- [x] WhatsApp URL generator
- [x] Repository interfaces

### TAHAP 3: Backend Repository (In Progress)
- [ ] ScrimRequestRepository PostgreSQL implementation
- [ ] ScrimMatchRepository PostgreSQL implementation
- [ ] RateLimiter Redis implementation
- [ ] Unit tests

### TAHAP 4: Backend Usecase/Service
- [ ] ScrimMatchmakingUsecase
- [ ] POKE matching logic (±2 tolerance)
- [ ] WARKOP matching logic (no tolerance)
- [ ] Background worker (Goroutines)
- [ ] Auto-cleanup scheduler (30 min expiry)
- [ ] WhatsApp notification integration
- [ ] Rate limiting middleware

### TAHAP 5: Backend HTTP Handlers
- [ ] POST /api/scrim/request - Create scrim request
- [ ] POST /api/scrim/:id/cancel - Cancel request
- [ ] GET /api/scrim/:id/status - Get request status
- [ ] POST /api/scrim/match/:id/confirm - Confirm match
- [ ] WebSocket /ws/scrim/:id - Real-time updates
- [ ] Rate limiting middleware

### TAHAP 6: Frontend - Vue.js 3 Setup
- [ ] Initialize Vue 3 project
- [ ] Install Tailwind CSS
- [ ] Setup Vite build
- [ ] Configure custom Tailwind colors & effects
- [ ] Create project structure

### TAHAP 7: Frontend - UI Components
- [ ] **DashboardHub.vue** - Main dashboard dengan hexagonal cards
- [ ] **SearchingRadar.vue** - Animated radar saat searching
- [ ] **MatchFoundModal.vue** - Dramatic VS screen
- [ ] **LobbyRoom.vue** - Post-match chat/coordination
- [ ] **NavbarHUD.vue** - Gaming-style navbar
- [ ] **StatCard.vue** - Hexagonal stat displays

### TAHAP 8: Frontend - Styling & Animations
- [ ] Glassmorphism effects
- [ ] Glow/bloom effects pada buttons
- [ ] Micro-interactions (hover, click)
- [ ] Magic circle/radar animations
- [ ] Countdown timer dengan color transition
- [ ] Responsive mobile layout

### TAHAP 9: Docker & Deployment
- [ ] Update Dockerfile untuk Vue build
- [ ] Multi-stage build (Vue → Go)
- [ ] Update docker-compose untuk include frontend
- [ ] Nginx reverse proxy (optional)
- [ ] Environment variables untuk WhatsApp API

### TAHAP 10: Testing & Documentation
- [ ] Integration tests
- [ ] E2E tests dengan Cypress
- [ ] API documentation
- [ ] Deployment guide
- [ ] User guide

---

## 🗂️ File Structure (Target)

```
matchMaking_go/
├── cmd/
│   └── server/
│       └── main.go (updated dengan scrim routes)
│
├── internal/
│   ├── domain/
│   │   ├── scrim.go ✅
│   │   ├── scrim_repository.go ✅
│   │   ├── team.go (existing)
│   │   └── match.go (existing)
│   │
│   ├── repository/
│   │   ├── scrim_request_repository.go (TODO)
│   │   ├── scrim_match_repository.go (TODO)
│   │   ├── rate_limiter_redis.go (TODO)
│   │   ├── team_repository.go (existing)
│   │   └── match_repository.go (existing)
│   │
│   ├── usecase/
│   │   ├── scrim_matchmaking_usecase.go (TODO)
│   │   └── matchmaking_usecase.go (existing)
│   │
│   ├── delivery/
│   │   └── http/
│   │       ├── scrim_handler.go (TODO)
│   │       ├── handler.go (existing)
│   │       └── websocket_hub.go (existing)
│   │
│   ├── middleware/
│   │   └── rate_limiter.go (TODO)
│   │
│   └── config/
│       └── config.go (update dengan WhatsApp config)
│
├── frontend/ (NEW)
│   ├── public/
│   ├── src/
│   │   ├── assets/
│   │   │   ├── images/
│   │   │   └── icons/
│   │   ├── components/
│   │   │   ├── DashboardHub.vue
│   │   │   ├── SearchingRadar.vue
│   │   │   ├── MatchFoundModal.vue
│   │   │   ├── LobbyRoom.vue
│   │   │   ├── NavbarHUD.vue
│   │   │   └── StatCard.vue
│   │   ├── views/
│   │   │   ├── HomeView.vue
│   │   │   └── MatchView.vue
│   │   ├── composables/
│   │   │   ├── useScrimMatch.js
│   │   │   └── useWebSocket.js
│   │   ├── App.vue
│   │   └── Main.js
│   ├── package.json
│   ├── vite.config.js
│   └── tailwind.config.js (dengan custom dark gaming theme)
│
├── migrations/
│   ├── 000001_init_schema.up.sql (existing)
│   ├── 000001_init_schema.down.sql (existing)
│   ├── 000002_create_scrim_requests.up.sql ✅
│   └── 000002_create_scrim_requests.down.sql ✅
│
├── Dockerfile (update untuk multi-stage Vue + Go)
├── docker-compose.yaml (update)
├── .env (update dengan WHATSAPP_API_KEY)
└── README.md (update)
```

---

## 🎨 UI/UX Specifications

### Color Palette (Tailwind Custom)

```javascript
// tailwind.config.js
colors: {
  'midnight': {
    900: '#0a0e27',
    800: '#141937',
    700: '#1e2447',
  },
  'magic': {
    violet: '#8b5cf6',
    cyan: '#06b6d4',
    gradient: 'linear-gradient(135deg, #8b5cf6 0%, #06b6d4 100%)',
  },
  'gold': {
    antique: '#d4af37',
    light: '#f4e4a6',
  },
}
```

### Key UI Elements

1. **Dashboard Hub**
   - Hexagonal stat cards dengan glassmorphism
   - Glow effect pada borders
   - Magic gradient backgrounds

2. **Find Scrim Button**
   - Giant portal-style button
   - Pulse animation saat idle
   - Explode animation saat diklik

3. **Searching State**
   - Rotating radar/magic circle
   - Pulsing rings
   - Estimated time display

4. **Match Found Modal**
   - Dramatic entrance (zoom + fade)
   - VS screen layout
   - 60s countdown dengan color shift (blue → red)
   - Gold "ACCEPT" button
   - Dark red "DECLINE" button

---

## 🔧 Technical Stack

### Backend
- **Language**: Go 1.21+
- **Framework**: Fiber v2
- **Database**: PostgreSQL 15
- **Cache**: Redis 7
- **WebSocket**: gorilla/websocket
- **Migration**: golang-migrate

### Frontend
- **Framework**: Vue 3 (Composition API)
- **Build Tool**: Vite
- **Styling**: Tailwind CSS 3
- **State**: Pinia (optional) atau Composables
- **HTTP**: Axios
- **WebSocket**: native WebSocket API

### DevOps
- **Container**: Docker
- **Orchestration**: Docker Compose
- **Reverse Proxy**: Nginx (optional)

---

## 🔄 Matchmaking Logic Flow

### POKE Category (Rank Weight 1-8)
```
1. Team A (POKE, Weight 5) masuk queue
2. Cari tim POKE dengan Weight 3-7 (±2)
3. Jika ada, create match
4. Jika tidak, tunggu di queue
5. Auto-expire setelah 30 menit
```

### WARKOP Category (Rank Weight 9-10)
```
1. Team B (WARKOP, Weight 10) masuk queue
2. Cari tim WARKOP mana pun (no tolerance
3. Match dengan tim WARKOP pertama yang tersedia
4. Create match immediately
```

### WhatsApp Integration
```
1. Match created
2. Generate WA URLs untuk kedua tim:
   - Team A: wa.me/{TeamB_number}?text=Hi...
   - Team B: wa.me/{TeamA_number}?text=Hi...
3. Send via WebSocket
4. 60 detik untuk confirm
5. Jika ada yang decline/timeout → back to queue
```

---

## 📡 API Endpoints

### Scrim Endpoints

```
POST   /api/scrim/request
Body: {
  "team_name": "Team Alpha",
  "whatsapp_number": "6281234567890",
  "category": "POKE",
  "rank_weight": 5
}
Response: {
  "id": "uuid",
  "status": "searching",
  "expires_at": "2024-..."
}

GET    /api/scrim/:id/status
Response: {
  "id": "uuid",
  "status": "matched",
  "match": {
    "id": "uuid",
    "opponent_name": "Team Beta",
    "whatsapp_url": "https://wa.me/...",
    "expires_in": 45
  }
}

POST   /api/scrim/:id/cancel
Response: { "message": "cancelled" }

POST   /api/scrim/match/:id/confirm
Response: { "message": "confirmed" }

WS     /ws/scrim/:id
Messages:
  - SEARCHING
  - MATCH_FOUND { opponent, whatsapp_url }
  - MATCH_CONFIRMED
  - MATCH_CANCELLED
```

---

## 🚀 Next Steps

1. **Immediate**: Implement repository layer (PostgreSQL + Redis)
2. **Next**: Implement usecase dengan matchmaking logic
3. **Then**: Setup Vue.js project dengan Tailwind
4. **Then**: Build UI components dengan dark gaming theme
5. **Finally**: Integration & testing

---

## 📝 Notes

- Pastikan rate limiting strict: 1 IP = max 1 active request
- Auto-cleanup harus reliable (cronjob atau background worker)
- WebSocket disconnection handling harus robust
- Mobile responsive adalah MUST
- Animasi harus smooth (60fps target)

---

**Status**: 🟡 In Development  
**Last Updated**: 2026-02-07  
**Version**: 2.0.0-alpha
