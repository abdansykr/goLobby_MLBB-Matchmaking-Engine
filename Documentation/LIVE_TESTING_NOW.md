# 🧪 Live Testing - Step by Step

## ✅ Prerequisites Done
- ✅ Docker containers running
- ✅ Database cleared
- ✅ Redis cleared
- ✅ Server healthy

---

## 🎯 Test 1: POKE Matching (LIVE TEST NOW!)

### Step 1: Create Team Alpha (POKE, Rank 5)

```powershell
$body1 = @{
   team_name = "Team Alpha"
    whatsapp_number = "6281111111111"
    category = "POKE"
    rank_weight = 5
} | ConvertTo-Json

Invoke-RestMethod -Method Post `
    -Uri "http://localhost:3000/api/scrim/request" `
    -ContentType "application/json" `
    -Body $body1
```

**Expected Output:**
```
category    : POKE
rank_weight : 5
status      : searching
request_id  : <save this ID>
```

**Save the request_id!**

---

### Step 2: Create Team Beta (POKE, Rank 6, Different IP)

```powershell
$body2 = @{
    team_name = "Team Beta"
    whatsapp_number = "6282222222222"
    category = "POKE"
    rank_weight = 6
} | ConvertTo-Json

$headers = @{
    "Content-Type" = "application/json"
    "X-Forwarded-For" = "192.168.1.100"
}

Invoke-RestMethod -Method Post `
    -Uri "http://localhost:3000/api/scrim/request" `
    -Headers $headers `
    -Body $body2
```

**Expected:** Both teams match instantly!

---

### Step 3: Check Match Status

```powershell
# Use request_id from Step 1
$requestId = "PASTE_REQUEST_ID_HERE"
Invoke-RestMethod -Uri "http://localhost:3000/api/scrim/request/$requestId"
```

**Expected Output:**
```json
{
  "status": "matched",
  "match": {
    "opponent_name": "Team Beta",
    "opponent_number": "6282222222222",
    "whatsapp_url": "https://wa.me/...",
    "expires_in": 60
  }
}
```

---

## 🎨 Using HTML Test Client

### Open Client
```bash
start scrim-test-client.html
```

### Steps:
1. **Clear any old data** - Click "Clear All Requests"

2. **Team 1:**
   - Team Name: Team Alpha
   - WhatsApp: 6281111111111
   - Category: **POKE**
   - Rank: **5**
   - Click "Find Scrim Match"

3. **Team 2:**
   - Team Name: Team Beta  
   - WhatsApp: 6282222222222
   - Category: **POKE**
   - Rank: **6**
   - Click "Find Scrim Match"

4. **Watch:**
   - Status changes to SEARCHING
   - Then MATCHED (usually < 2 seconds!)
   - WhatsApp links appear
   - Countdown timer starts (60s)

---

## 🐛 If You Get Errors

### Error: "you already have an active request"

**Fix:**
```powershell
# Clear Redis
docker exec golobby_redis redis-cli -a redis123 FLUSHALL

# Clear database
docker exec golobby_postgres psql -U postgres -d golobby -c "TRUNCATE scrim_requests, scrim_matches CASCADE;"
```

### Error: "Network error: Failed to fetch"

**Fix:**
```powershell
# Check server is running
Invoke-RestMethod http://localhost:3000/health

# If failed, restart
docker-compose restart app
```

### Error: CORS issue in browser

**Fix:** Already handled in code, but if issues persist:
- Use Chrome/Firefox (not IE)
- Disable browser extensions
- Open dev console (F12) to see exact error

---

## 📊 View Database After Match

```powershell
# View matched requests
docker exec golobby_postgres psql -U postgres -d golobby -c "
SELECT team_name, category, rank_weight, status 
FROM scrim_requests 
ORDER BY created_at DESC;"

# View matches
docker exec golobby_postgres psql -U postgres -d golobby -c "
SELECT id, category, rank_diff, status, created_at 
FROM scrim_matches 
ORDER BY created_at DESC;"
```

---

## 🔥 Quick Reset (Fresh Start)

```powershell
# One command to clear everything
docker exec golobby_redis redis-cli -a redis123 FLUSHALL; docker exec golobby_postgres psql -U postgres -d golobby -c "TRUNCATE scrim_requests, scrim_matches CASCADE;"
```

---

**Ready to test!** 🚀

Try the PowerShell commands above OR use the HTML client!
