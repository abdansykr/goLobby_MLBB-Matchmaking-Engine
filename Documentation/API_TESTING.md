# API Testing Guide

This guide provides ready-to-use commands for testing the Antigravity Matchmaking API.

## Prerequisites

- Server running on `http://localhost:8080`
- `curl` installed (comes with Windows 10+, macOS, Linux)
- `jq` (optional, for pretty JSON output)

## Test Workflow

### 1. Health Check

```bash
curl http://localhost:8080/health
```

**Expected Response:**
```json
{
  "status": "healthy",
  "service": "antigravity-matchmaking"
}
```

---

### 2. Enqueue Team 1 (Rank 75)

```bash
curl -X POST http://localhost:8080/api/matchmaking/enqueue \
  -H "Content-Type: application/json" \
  -d '{
    "captain_name": "ProPlayer123",
    "team_name": "Team Alpha",
    "average_rank": 75
  }'
```

**Expected Response:**
```json
{
  "message": "Team enqueued successfully",
  "team_id": "123e4567-e89b-12d3-a456-426614174000",
  "status": "WAITING"
}
```

**💡 Save the `team_id` - you'll need it!**

---

### 3. Enqueue Team 2 (Rank 76 - within ±2 range)

```bash
curl -X POST http://localhost:8080/api/matchmaking/enqueue \
  -H "Content-Type: application/json" \
  -d '{
    "captain_name": "EliteGamer",
    "team_name": "Team Beta",
    "average_rank": 76
  }'
```

**➡️ These teams should match immediately!**

---

### 4. Test Extended Range Matching

#### Enqueue Team A (Rank 50)

```bash
curl -X POST http://localhost:8080/api/matchmaking/enqueue \
  -H "Content-Type: application/json" \
  -d '{
    "captain_name": "MidRanker",
    "team_name": "Team Gamma",
    "average_rank": 50
  }'
```

**⏳ Wait 30 seconds...**

#### Enqueue Team B (Rank 54)

This is outside the initial ±2 range but within the extended ±4 range.

```bash
curl -X POST http://localhost:8080/api/matchmaking/enqueue \
  -H "Content-Type: application/json" \
  -d '{
    "captain_name": "ChallengSeeker",
    "team_name": "Team Delta",
    "average_rank": 54
  }'
```

**➡️ These teams should match after the 30-second timeout!**

---

### 5. Confirm Ready

Replace `TEAM_ID` and `MATCH_ID` with actual values from your response.

```bash
curl -X POST http://localhost:8080/api/matchmaking/ready \
  -H "Content-Type: application/json" \
  -d '{
    "team_id": "123e4567-e89b-12d3-a456-426614174000",
    "match_id": "987fcdeb-51a2-43d7-9876-543210fedcba"
  }'
```

**Expected Response:**
```json
{
  "message": "Ready confirmation received",
  "status": "READY"
}
```

---

### 6. Cancel Matchmaking

```bash
curl -X POST "http://localhost:8080/api/matchmaking/cancel?team_id=TEAM_ID"
```

**Expected Response:**
```json
{
  "message": "Matchmaking cancelled",
  "status": "CANCELLED"
}
```

---

## WebSocket Testing

### Using websocat (CLI WebSocket Client)

Install: `cargo install websocat` or download from GitHub.

```bash
websocat "ws://localhost:8080/ws?team_id=YOUR_TEAM_ID"
```

You'll receive JSON messages when a match is found:

```json
{
  "type": "MATCH_FOUND",
  "team_id": "123e4567-e89b-12d3-a456-426614174000",
  "match_id": "987fcdeb-51a2-43d7-9876-543210fedcba",
  "data": {
    "opponent_name": "Team Beta",
    "opponent_rank": 76,
    "match_id": "987fcdeb-51a2-43d7-9876-543210fedcba",
    "expires_at": "2024-01-15T10:35:00Z"
  },
  "timestamp": "2024-01-15T10:34:00Z"
}
```

### Using Browser (test-client.html)

1. Open `test-client.html` in your browser
2. Fill in the form
3. Click "Find Match"
4. Wait for opponent

**Open multiple browser tabs to simulate multiple teams!**

---

## Testing Anti-Ghosting System

### Scenario: Team Doesn't Confirm Ready

1. **Enqueue two teams** (they match)
2. **DO NOT click "Confirm Ready"**
3. **Wait 60 seconds**
4. **Check server logs**:

```bash
docker-compose logs -f app
```

You should see:
```
Applied ghosting penalty to team: Team Alpha (new score: 90)
Applied ghosting penalty to team: Team Beta (new score: 90)
```

Both teams lose 10 reputation points!

---

## Testing Rank Matching Logic

### Test Case 1: Perfect Match (±0)

```bash
# Team 1: Rank 80
curl -X POST http://localhost:8080/api/matchmaking/enqueue \
  -H "Content-Type: application/json" \
  -d '{"captain_name": "Player1", "team_name": "Team1", "average_rank": 80}'

# Team 2: Rank 80 (exact match)
curl -X POST http://localhost:8080/api/matchmaking/enqueue \
  -H "Content-Type: application/json" \
  -d '{"captain_name": "Player2", "team_name": "Team2", "average_rank": 80}'
```

**➡️ Instant match**

---

### Test Case 2: Within Initial Range (±2)

```bash
# Team 1: Rank 70
curl -X POST http://localhost:8080/api/matchmaking/enqueue \
  -H "Content-Type: application/json" \
  -d '{"captain_name": "Player3", "team_name": "Team3", "average_rank": 70}'

# Team 2: Rank 72 (difference of 2)
curl -X POST http://localhost:8080/api/matchmaking/enqueue \
  -H "Content-Type: application/json" \
  -d '{"captain_name": "Player4", "team_name": "Team4", "average_rank": 72}'
```

**➡️ Instant match**

---

### Test Case 3: Beyond Initial, Within Extended (±3 to ±4)

```bash
# Team 1: Rank 60
curl -X POST http://localhost:8080/api/matchmaking/enqueue \
  -H "Content-Type: application/json" \
  -d '{"captain_name": "Player5", "team_name": "Team5", "average_rank": 60}'

# Wait 30 seconds (for range expansion)

# Team 2: Rank 63 (difference of 3, outside ±2 but within ±4)
curl -X POST http://localhost:8080/api/matchmaking/enqueue \
  -H "Content-Type: application/json" \
  -d '{"captain_name": "Player6", "team_name": "Team6", "average_rank": 63}'
```

**➡️ Match after 30 seconds**

---

### Test Case 4: Beyond Extended Range (>±4)

```bash
# Team 1: Rank 40
curl -X POST http://localhost:8080/api/matchmaking/enqueue \
  -H "Content-Type: application/json" \
  -d '{"captain_name": "Player7", "team_name": "Team7", "average_rank": 40}'

# Team 2: Rank 50 (difference of 10, beyond ±4)
curl -X POST http://localhost:8080/api/matchmaking/enqueue \
  -H "Content-Type: application/json" \
  -d '{"captain_name": "Player8", "team_name": "Team8", "average_rank": 50}'
```

**➡️ No match - keeps searching**

---

## Monitoring

### Check PostgreSQL Database

```bash
docker exec -it antigravity_postgres psql -U postgres -d antigravity

# List all teams
SELECT id, team_name, average_rank, status, reputation_score FROM teams;

# List all matches
SELECT id, team1_id, team2_id, status, rank_diff FROM matches;

# Exit
\q
```

### Check Redis Queue

```bash
docker exec -it antigravity_redis redis-cli -a redis123

# Check queue length
ZCARD matchmaking:queue

# View queue contents
ZRANGE matchmaking:queue 0 -1

# Check locks
KEYS matchmaking:lock:*

# Exit
exit
```

---

## Performance Testing

### Simulate 100 Teams

```bash
for i in {1..100}; do
  curl -X POST http://localhost:8080/api/matchmaking/enqueue \
    -H "Content-Type: application/json" \
    -d "{\"captain_name\": \"Player$i\", \"team_name\": \"Team$i\", \"average_rank\": $((RANDOM % 100))}" &
done
wait
```

**Check logs to see matchmaking speed!**

---

## Common Issues

### Issue: "Connection refused"
**Solution:** Ensure server is running: `docker-compose up -d`

### Issue: "Team is already in a match"
**Solution:** Team is locked. Wait for match to expire or cancel first.

### Issue: WebSocket disconnects
**Solution:** Check if team_id is valid UUID format.

---

## PowerShell Commands (Windows)

For Windows users, use PowerShell:

```powershell
# Health Check
Invoke-RestMethod -Uri "http://localhost:8080/health"

# Enqueue Team
$body = @{
    captain_name = "ProPlayer123"
    team_name = "Team Alpha"
    average_rank = 75
} | ConvertTo-Json

Invoke-RestMethod -Method Post -Uri "http://localhost:8080/api/matchmaking/enqueue" `
    -ContentType "application/json" -Body $body
```

---

**Happy Testing! 🎮**
