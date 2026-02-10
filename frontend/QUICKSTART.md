# 🚀 QUICK START - GoLobby Frontend

## ✅ Error Fixed!

The `border-border` class error has been fixed in `main.css`.

---

## 📋 Setup Steps

### 1. Install Node.js (if not installed)

**Download:** https://nodejs.org/en/download

**Required version:** Node.js >= 18.0

**Verify installation:**
```bash
node --version   # Should show v18.x or higher
npm --version    # Should show v9.x or higher
```

---

### 2. Navigate to Frontend Directory

```bash
cd c:\Users\acer\Development\go-projects\matchMaking_go\frontend
```

---

### 3. Install Dependencies

```bash
npm install
```

**This will install:**
- Vue 3.4.15
- Vite 5.0.11
- Tailwind CSS 3.4.1
- Vue Router 4.2.5
- Axios 1.6.5

**Time:** ~2-3 minutes

---

### 4. Start Development Server

```bash
npm run dev
```

**Expected output:**
```
  VITE v5.0.11  ready in 1234 ms

  ➜  Local:   http://localhost:5173/
  ➜  Network: use --host to expose
  ➜  press h to show help
```

---

### 5. Open in Browser

Navigate to: **http://localhost:5173**

---

## 🎨 What You'll See

### Dark Fantasy Interface
✅ **Midnight blue** background with subtle purple/cyan gradients  
✅ **HUD-style navbar** with hexagon logo (crystal glow effect)  
✅ **Glassmorphism cards** (frosted glass effect)  
✅ **3 Stat Cards:**
   - Average Rank (Mythic Glory with glow)
   - Reputation Score (circular progress 95%)
   - Today's Matches (8W-2L, 80% WR)

✅ **Giant Portal Button** - "FIND SCRIM MATCH" with purple pulse animation  
✅ **Recent Matches List** - Win/loss history with glass cards

---

## 🧪 Testing the UI

### Test Search Animation

1. Click **"FIND SCRIM MATCH"** button
2. Portal transforms into **animated radar** (3 rotating circles)
3. Shows "Searching for Opponents..." with countdown
4. Click "Cancel Search" to return

### Test Match Modal (Manual Trigger)

Open `src/App.vue` and uncomment lines 43-59 to auto-show modal after 5 seconds.

---

## 🏗️ Build for Production

```bash
npm run build
```

**Output:** `dist/` folder  
**Deploy:** Upload to Vercel, Netlify, or any static host

---

## 🔧 Troubleshooting

### Error: "npm is not recognized"

**Fix:** Install Node.js from https://nodejs.org

### Error: Port 5173 already in use

**Fix:**
```bash
# Kill process on port 5173
npx kill-port 5173

# Or use different port
npm run dev -- --port 3001
```

### Styles not loading

**Fix:**
```bash
# Clear cache
rm -rf node_modules .vite
npm install
npm run dev
```

### Fonts not showing

**Fix:** Check internet connection (Google Fonts loads from CDN)

---

## 📁 Project Structure

```
frontend/
├── src/
│   ├── components/        # Reusable components
│   │   ├── NavbarHUD.vue
│   │   └── MatchFoundModal.vue
│   ├── views/             # Page components
│   │   └── DashboardView.vue
│   ├── assets/            # CSS & images
│   │   └── main.css       # 200+ lines custom styles!
│   ├── App.vue            # Root component
│   └── main.js            # Entry point
├── public/                # Static files
├── index.html             # HTML template
├── package.json           # Dependencies
├── vite.config.js         # Vite config + API proxy
├── tailwind.config.js     # Dark Fantasy theme
└── postcss.config.js      # PostCSS config
```

---

## 🎯 Next Steps

### Option A: Just View the UI
```bash
npm run dev
# Explore the beautiful interface!
```

### Option B: Connect to Backend
1. Ensure backend is running on `localhost:3000`
2. Create `src/composables/useScrimAPI.js`
3. Connect Dashboard search button to API
4. Poll for match status
5. Show MatchFoundModal when matched

### Option C: Customize Design
1. Edit `tailwind.config.js` for colors
2. Modify `src/assets/main.css` for effects
3. Hot reload will update instantly!

---

## 🎨 Design System Quick Reference

### Colors (Tailwind)
```html
<!-- Backgrounds -->
<div class="bg-midnight-900">
<div class="bg-midnight-800">

<!-- Accents -->
<div class="bg-electric-violet-500">
<div class="bg-cyan-magic-400">
<div class="bg-antique-gold-400">
```

### Components
```html
<!-- Glass Card -->
<div class="glass-card p-6">Content</div>

<!-- Magic Button -->
<button class="btn-magic">Click me</button>

<!-- Hexagon Avatar -->
<div class="avatar-hexagon">
  <img src="avatar.jpg" />
</div>
```

### Effects
```html
<!-- Glow Text -->
<h1 class="text-glow">Glowing!</h1>

<!-- Gradient Text -->
<span class="gradient-text">Magic</span>

<!-- Portal Pulse -->
<div class="portal-pulse">Pulsing</div>
```

---

## 🎮 Enjoy Your Dark Fantasy Arena! 🔮

**Everything is ready to go!** Just run:

```bash
cd frontend
npm install
npm run dev
```

**Open:** http://localhost:5173

---

**Status:** 🟢 **READY FOR DEVELOPMENT**

**Issues?** Check the troubleshooting section above! 🚀
