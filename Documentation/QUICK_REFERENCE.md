# 📋 GoLobby Scrim - Quick Reference Card

## 🚀 Quick Start (30 Seconds)

```bash
# Start everything
docker-compose up -d

# Open test client
start scrim-test-client.html

# Create POKE match
# Team 1: POKE, Rank 5 → Click "Find Scrim Match"
# Team 2: POKE, Rank 6 → Click "Find Scrim Match"
# Result: MATCHED! 🎉
```

---

## 🎯 API Cheat Sheet

### Create Request (POKE)
```bash
curl -X POST http://localhost:3000/api/scrim/request \
  -H "Content-Type: application/json" \
  -d '{
    "team_name": "Team Alpha",
    "whatsapp_number": "6281234567890",
    "category": "POKE",
    "rank_weight": 5
  }'
```

### Create Request (WARKOP)
```bash
curl -X POST http://localhost:3000/api/scrim/request \
  -H "Content-Type: application/json" \
  -H "X-Forwarded-For: 192.168.1.100" \
  -d '{
    "team_name": "Team Beta",
    "whatsapp_number": "6289876543210",
    "category": "WARKOP",
    "rank_weight": 9
  }'
```

### Check Status
```bash
curl http://localhost:3000/api/scrim/request/REQUEST_ID
```

### Cancel
```bash
curl -X POST http://localhost:3000/api/scrim/request/REQUEST_ID/cancel
```

---

## 🎮 PowerShell Quick Commands

```powershell
# Health check
Invoke-RestMethod http://localhost:3000/health

# Create POKE request
$body = @{team_name="Team1";whatsapp_number="628111";category="POKE";rank_weight=5} | ConvertTo-Json
Invoke-RestMethod -Method Post -Uri "http://localhost:3000/api/scrim/request" -ContentType "application/json" -Body $body

# Create WARKOP (different IP)
$body = @{team_name="Team2";whatsapp_number="628222";category="WARKOP";rank_weight=9} | ConvertTo-Json
$headers = @{"Content-Type"="application/json";"X-Forwarded-For"="192.168.1.100"}
Invoke-RestMethod -Method Post -Uri "http://localhost:3000/api/scrim/request" -Headers $headers -Body $body
```

---

## 🗄️ Database Quick Access

```bash
# Connect to PostgreSQL
docker exec -it golobby_postgres psql -U postgres -d golobby

# View active requests
SELECT team_name, category, rank_weight, status FROM scrim_requests WHERE status='searching';

# View all matches
SELECT * FROM scrim_matches ORDER BY created_at DESC LIMIT 5;

# Clear all data
TRUNCATE scrim_requests, scrim_matches CASCADE;
```

---

## 🔧 Docker Commands

```bash
# Start
docker-compose up -d

# Stop
docker-compose down

# Rebuild
docker-compose up -d --build

# View logs (live)
docker logs golobby_app -f

# Restart just app
docker-compose restart app

# Check status
docker-compose ps
```

---

## 🧪 Testing Shortcuts

### POKE Match (Success)
```
Team 1: POKE, Rank 5
Team 2: POKE, Rank 6
Result: ✅ MATCHED (diff=1, within ±2)
```

### POKE No Match
```
Team 1: POKE, Rank 3
Team 2: POKE, Rank 7
Result: ❌ NO MATCH (diff=4, exceeds ±2)
```

### WARKOP Instant
```
Team 1: WARKOP, Rank 9
Team 2: WARKOP, Rank 10
Result: ✅ INSTANT MATCHED
```

### Rate Limit Test
```bash
# First request (same IP)
curl -X POST http://localhost:3000/api/scrim/request -H "Content-Type: application/json" -d '{"team_name":"T1","whatsapp_number":"628111","category":"POKE","rank_weight":5}'

# Second request (same IP)
curl -X POST http://localhost:3000/api/scrim/request -H "Content-Type: application/json" -d '{"team_name":"T2","whatsapp_number":"628222","category":"POKE","rank_weight":6}'

# Expected: "error": "you already have an active request"
```

---

## 📊 Monitoring

### Check Logs
```bash
# All logs
docker logs golobby_app

# Follow live
docker logs golobby_app -f --tail=50

# Grep for matches
docker logs golobby_app | Select-String "MATCH"

# Grep for errors
docker logs golobby_app | Select-String "error"
```

### Check Workers
```bash
docker logs golobby_app | Select-String "Worker.*started"
# Expected:
# Scrim Worker 0 started
# Scrim Worker 1 started
# Scrim Worker 2 started
# Scrim Worker 3 started
```

### Check Database Health
```bash
docker exec golobby_postgres psql -U postgres -c "SELECT version();"
```

### Check Redis
```bash
docker exec golobby_redis redis-cli -a redis123 PING
# Expected: PONG
```

---

## 🐛 Troubleshooting

### No Match Happening?
```bash
# Check workers
docker logs golobby_app | Select-String "Scrim Worker"

# Check database
docker exec golobby_postgres psql -U postgres -d golobby -c "SELECT * FROM scrim_requests WHERE status='searching';"

# Check if categories match
# POKE only matches POKE
# WARKOP only matches WARKOP
```

### Rate Limit Stuck?
```bash
# Clear Redis
docker exec golobby_redis redis-cli -a redis123 FLUSHDB

# Or cancel specific request
curl -X POST http://localhost:3000/api/scrim/request/REQUEST_ID/cancel
```

### Container Won't Start?
```bash
# Check logs
docker logs golobby_app

# Common issues:
# - Port 3000 already in use
# - Database migration failed
# - Redis connection failed

# Nuclear option
docker-compose down -v
docker-compose up -d --build
```

---

## 📁 File Locations

```
Project Root: c:\Users\acer\Development\go-projects\matchMaking_go\

Key Files:
- scrim-test-client.html     → Test UI
- SCRIM_API_DOCS.md           → API reference
- TESTING_GUIDE.md            → Test scenarios
- SCRIM_PROJECT_COMPLETE.md   → Full summary

Code:
- internal/domain/scrim.go
- internal/repository/scrim_request_repository.go
- internal/usecase/scrim_matchmaking_usecase.go
- internal/delivery/http/scrim_handler.go

Database:
- migrations/000002_create_scrim_requests.up.sql
```

---

## 🔑 Environment Variables

```env
# .env file
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=golobby

REDIS_ADDR=localhost:6379
REDIS_PASSWORD=redis123

APP_PORT=8080
APP_ENV=production
```

---

## 📞 Support

| Issue | Solution |
|-------|----------|
| API not responding | `docker logs golobby_app` |
| Database error | `docker logs golobby_postgres` |
| Rate limit issues | `docker exec golobby_redis redis-cli -a redis123 FLUSHDB` |
| Match not found | Check category match + rank tolerance |
| Port conflict | See `SOLUSI_PORT_POSTGRES.md` |

---

## 🎯 Matching Rules

### POKE (Rank 1-8)
```
Tolerance: ±2

Examples:
Rank 5 CAN match: 3, 4, 5, 6, 7
Rank 5 CANNOT match: 1, 2, 8
Rank 1 CAN match: 1, 2, 3
Rank 8 CAN match: 6, 7, 8
```

### WARKOP (Rank 9-10)
```
Tolerance: NONE (instant match)

Examples:
Rank 9 CAN match: 9, 10 (any WARKOP)
Rank 10 CAN match: 9, 10 (any WARKOP)
```

### Category Mismatch
```
POKE ❌ WARKOP
WARKOP ❌ POKE
```

---

## ⏱️ Timeouts

| Event | Duration |
|-------|----------|
| Request expiry | 30 minutes |
| Match confirmation | 60 seconds |
| Rate limit TTL | 30 minutes |
| Cleanup check | Every 10 seconds |
| Status poll | Every 2 seconds (HTML client) |

---

## 🎨 Status Colors (HTML Client)

- 🟡 **Yellow** → SEARCHING
- 🟢 **Green** → MATCHED
- 🔴 **Red** → EXPIRED
- ⚪ **Gray** → IDLE

---

## 💡 Pro Tips

1. **Use X-Forwarded-For header** for testing different IPs
2. **Clear Redis** if rate limit stuck: `redis-cli FLUSHDB`
3. **Check logs live** with `-f` flag: `docker logs -f golobby_app`
4. **Use HTML test client** - it's 10x faster than curl
5. **Keep database clean** - truncate old data periodically
6. **Monitor worker logs** - they show matching activity

---

**Quick Access:** Open `scrim-test-client.html` for visual testing! 🎮
