# 🚀 Quick Start Guide

Get the Antigravity Matchmaking Engine running in under 5 minutes!

## Prerequisites

- **Docker Desktop** installed and running
- **Git** (optional, for cloning)
- **Web Browser** (for testing WebSocket client)

That's all! Docker will handle everything else.

---

## Step 1: Start the Services

Open your terminal and run:

```bash
cd c:\Users\acer\Development\go-projects\matchMaking_go
docker-compose up -d
```

**What this does:**
- Starts PostgreSQL database
- Starts Redis cache
- Runs database migrations
- Builds and starts the matchmaking app

**Expected output:**
```
Creating network "matchmaking_go_antigravity_network" ... done
Creating antigravity_postgres ... done
Creating antigravity_redis ... done
Creating antigravity_migrate ... done
Creating antigravity_app ... done
```

---

## Step 2: Verify Everything is Running

```bash
docker-compose ps
```

**You should see:**
```
NAME                  STATUS       PORTS
antigravity_app       Up           0.0.0.0:8080->8080/tcp
antigravity_postgres  Up (healthy) 0.0.0.0:5432->5432/tcp
antigravity_redis     Up (healthy) 0.0.0.0:6379->6379/tcp
```

**Health check:**
```bash
curl http://localhost:8080/health
```

**Expected response:**
```json
{"status":"healthy","service":"antigravity-matchmaking"}
```

✅ **Your matchmaking engine is now live!**

---

## Step 3: Test Matchmaking (Interactive)

### Option A: Using the Web Client (Recommended)

1. **Open `test-client.html` in your browser**
   - Double-click the file or drag it into your browser

2. **Fill in the form:**
   - Captain Name: `ProPlayer123`
   - Team Name: `Team Alpha`
   - Average Rank: `75`

3. **Click "Find Match"**

4. **Open another browser tab** (or use incognito mode)

5. **In the second tab, create another team:**
   - Captain Name: `EliteGamer`
   - Team Name: `Team Beta`
   - Average Rank: `76`

6. **Click "Find Match"**

**🎮 BOOM! Both teams should match instantly!**

You'll see:
- "MATCH FOUND!" notification
- Opponent details
- "Confirm Ready" button

7. **Click "Confirm Ready" on both tabs**

---

### Option B: Using cURL (Command Line)

**Terminal 1 - Team Alpha:**
```bash
curl -X POST http://localhost:8080/api/matchmaking/enqueue \
  -H "Content-Type: application/json" \
  -d '{
    "captain_name": "ProPlayer123",
    "team_name": "Team Alpha",
    "average_rank": 75
  }'
```

**Save the `team_id` from response!**

**Terminal 2 - Team Beta:**
```bash
curl -X POST http://localhost:8080/api/matchmaking/enqueue \
  -H "Content-Type: application/json" \
  -d '{
    "captain_name": "EliteGamer",
    "team_name": "Team Beta",
    "average_rank": 76
  }'
```

**Match created! 🎉**

---

## Step 4: View Logs

**See what's happening behind the scenes:**

```bash
docker-compose logs -f app
```

**You'll see:**
```
Worker 0: Processing team Team Alpha (Rank: 75)
Worker 1: Processing team Team Beta (Rank: 76)
Worker 0: Match created! Team Alpha vs Team Beta (found in 0s)
```

Press `Ctrl+C` to stop viewing logs.

---

## Step 5: Test Anti-Ghosting System

1. **Create two teams** (using web client or cURL)
2. **DON'T click "Confirm Ready"**
3. **Wait 60 seconds**
4. **Check logs:**

```bash
docker-compose logs app | grep "ghosting penalty"
```

**You should see:**
```
Applied ghosting penalty to team: Team Alpha (new score: 90)
Applied ghosting penalty to team: Team Beta (new score: 90)
```

Both teams lost 10 reputation points! 💀

---

## Step 6: Explore the Database

**Connect to PostgreSQL:**
```bash
docker exec -it antigravity_postgres psql -U postgres -d antigravity
```

**View all teams:**
```sql
SELECT team_name, average_rank, status, reputation_score FROM teams;
```

**View all matches:**
```sql
SELECT * FROM matches;
```

**Exit:**
```sql
\q
```

---

## Step 7: Clean Up (When Done)

**Stop all services:**
```bash
docker-compose down
```

**Remove ALL data (database, queues, etc.):**
```bash
docker-compose down -v
```

**⚠️ Warning:** This deletes all your data!

---

## Common Commands Cheat Sheet

```bash
# Start services
docker-compose up -d

# Stop services
docker-compose down

# View logs
docker-compose logs -f app

# Restart app only
docker-compose restart app

# Rebuild app (after code changes)
docker-compose up -d --build app

# Check health
curl http://localhost:8080/health

# View database
docker exec -it antigravity_postgres psql -U postgres -d antigravity

# View Redis
docker exec -it antigravity_redis redis-cli -a redis123
```

---

## What's Next?

1. **Read the Documentation:**
   - `README.md` - Full feature overview
   - `API_TESTING.md` - Detailed API examples
   - `ARCHITECTURE.md` - System design deep dive

2. **Customize Configuration:**
   - Edit `.env` file
   - Adjust matchmaking timeouts
   - Change rank ranges

3. **Integrate with Your App:**
   - Use the REST API endpoints
   - Connect via WebSocket for real-time updates
   - Build your own frontend

4. **Deploy to Production:**
   - Use cloud services (AWS, GCP, Azure)
   - Set up monitoring (Prometheus, Grafana)
   - Configure proper security (HTTPS, authentication)

---

## Troubleshooting

### Issue: "Port already in use"
**Solution:**
```bash
# Find what's using port 8080
netstat -ano | findstr :8080  # Windows
lsof -i :8080                 # Mac/Linux

# Stop the process or change APP_PORT in .env
```

### Issue: "Connection refused"
**Solution:**
```bash
# Check if services are running
docker-compose ps

# Restart if needed
docker-compose restart
```

### Issue: Teams not matching
**Solution:**
```bash
# Check logs
docker-compose logs app

# Verify rank difference is ≤4
# Initial: ±2, Extended (after 30s): ±4
```

---

## Support

- **GitHub Issues**: Report bugs or request features
- **Documentation**: Check `README.md` and `ARCHITECTURE.md`
- **Logs**: Always check `docker-compose logs app` first

---

**You're all set! Happy matchmaking! 🎮🚀**
