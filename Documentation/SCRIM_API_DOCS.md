# 📚 GoLobby Scrim API Documentation

## 🔗 Base URL
```
http://localhost:3000/api/scrim
```

---

## 📡 Endpoints

### 1. Create Scrim Request
Creates a new scrim matchmaking request.

**Endpoint:** `POST /api/scrim/request`

**Request Body:**
```json
{
  "team_name": "Team Alpha",
  "whatsapp_number": "6281234567890",
  "category": "POKE",
  "rank_weight": 5
}
```

**Field Specifications:**
| Field | Type | Required | Validation |
|-------|------|----------|------------|
| `team_name` | string | Yes | Non-empty |
| `whatsapp_number` | string | Yes | Valid phone number (e.g., 628xxx) |
| `category` | string | Yes | Must be "POKE" or "WARKOP" |
| `rank_weight` | integer | Yes | 1-10 (POKE: 1-8, WARKOP: 9-10) |

**Category Rules:**
- **POKE** (Rank 1-8):
  - Matches with teams within ±2 rank tolerance
  - Example: Rank 5 can match with Rank 3-7
  
- **WARKOP** (Rank 9-10):
  - Instant matching with any WARKOP team
  - No rank tolerance check

**Success Response:** `201 Created`
```json
{
  "message": "Scrim request created successfully",
  "request_id": "ab745781-f0eb-41c5-942b-75a6280fe2c5",
  "status": "searching",
  "category": "POKE",
  "rank_weight": 5,
  "expires_at": "2026-02-06T20:15:51.174707Z"
}
```

**Error Responses:**

`400 Bad Request` - Invalid input
```json
{
  "error": "POKE category requires rank_weight between 1 and 8"
}
```

`400 Bad Request` - Rate limit exceeded
```json
{
  "error": "you already have an active request"
}
```

**Rate Limiting:**
- **1 active request per IP address**
- Limit resets when request is cancelled, expired, or matched
- TTL: 30 minutes

---

### 2. Get Request Status
Retrieves the current status of a scrim request.

**Endpoint:** `GET /api/scrim/request/:id`

**URL Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | UUID | Request ID from create response |

**Success Response:** `200 OK`

**When Searching:**
```json
{
  "request_id": "ab745781-f0eb-41c5-942b-75a6280fe2c5",
  "team_name": "Team Alpha",
  "category": "POKE",
  "rank_weight": 5,
  "status": "searching",
  "created_at": "2026-02-06T19:45:51.174707Z",
  "expires_at": "2026-02-06T20:15:51.174707Z"
}
```

**When Matched:**
```json
{
  "request_id": "ab745781-f0eb-41c5-942b-75a6280fe2c5",
  "team_name": "Team Alpha",
  "category": "POKE",
  "rank_weight": 5,
  "status": "matched",
  "created_at": "2026-02-06T19:45:51Z",
  "expires_at": "2026-02-06T20:15:51Z",
  "match": {
    "match_id": "c7d8e9f0-1234-5678-90ab-cdef12345678",
    "opponent_name": "Team Beta",
    "opponent_number": "6289876543210",
    "whatsapp_url": "https://wa.me/6289876543210?text=Hi%20Team%20Beta!%20We've%20been%20matched...",
    "expires_in": 45,
    "status": "pending"
  }
}
```

**Status Values:**
| Status | Description |
|--------|-------------|
| `searching` | Actively looking for match |
| `matched` | Match found, awaiting confirmation |
| `expired` | Request timed out (30 min) |
| `cancelled` | Cancelled by user |

**Error Response:** `404 Not Found`
```json
{
  "error": "Request not found"
}
```

---

### 3. Cancel Request
Cancels an active scrim request.

**Endpoint:** `POST /api/scrim/request/:id/cancel`

**URL Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | UUID | Request ID to cancel |

**Success Response:** `200 OK`
```json
{
  "message": "Request cancelled successfully"
}
```

**Error Responses:**

`400 Bad Request` - Invalid ID format
```json
{
  "error": "Invalid request ID"
}
```

`400 Bad Request` - Request not found
```json
{
  "error": "failed to get request: scrim request not found"
}
```

**Side Effects:**
- Removes rate limit for IP address
- Updates status to `cancelled`
- Match becomes invalid if already matched

---

### 4. Confirm Match
Confirms participation in a matched scrim.

**Endpoint:** `POST /api/scrim/match/:id/confirm`

**URL Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | UUID | Match ID from status response |

**Success Response:** `200 OK`
```json
{
  "message": "Match confirmed successfully"
}
```

**Error Responses:**

`400 Bad Request` - Match not pending
```json
{
  "error": "match is not pending"
}
```

`400 Bad Request` - Match expired
```json
{
  "error": "failed to get match: scrim match not found"
}
```

**Notes:**
- Matches expire after **60 seconds** if not confirmed
- Both teams should confirm independently
- After confirmation, use WhatsApp link to coordinate

---

## 🔄 Request Lifecycle

```
[CREATE REQUEST]
       ↓
[STATUS: searching] ← (polling every 2s recommended)
       ↓ (match found)
[STATUS: matched]
       ↓
[CONFIRM MATCH] (within 60s)
       ↓
[coordinated via WhatsApp]

Parallel paths:
- [CANCEL] → status: cancelled
- [30 min timeout] → status: expired
- [Match timeout] → back to searching
```

---

## 🧪 Testing Examples

### Example 1: POKE Match (Success)

**Step 1:** Team A Creates Request
```bash
curl -X POST http://localhost:3000/api/scrim/request \
  -H "Content-Type: application/json" \
  -d '{
    "team_name": "Team Alpha",
    "whatsapp_number": "6281111111111",
    "category": "POKE",
    "rank_weight": 5
  }'
```

**Step 2:** Team B Creates Request (Different IP!)
```bash
curl -X POST http://localhost:3000/api/scrim/request \
  -H "Content-Type: application/json" \
  -H "X-Forwarded-For: 192.168.1.100" \
  -d '{
    "team_name": "Team Beta",
    "whatsapp_number": "6282222222222",
    "category": "POKE",
    "rank_weight": 6
  }'
```

**Step 3:** Check Status (Team A)
```bash
curl http://localhost:3000/api/scrim/request/REQUEST_ID_HERE
```

**Expected:** Both teams status = `matched` with opponent details

---

### Example 2: POKE No Match (Rank Too Far)

**Team A:** Rank 3
**Team B:** Rank 7

**Result:** No match (±4 exceeds ±2 tolerance)

---

### Example 3: WARKOP Instant Match

**Team A:**
```json
{
  "team_name": "Pro Team",
  "whatsapp_number": "6283333333333",
  "category": "WARKOP",
  "rank_weight": 9
}
```

**Team B:**
```json
{
  "team_name": "Elite Squad",
  "whatsapp_number": "6284444444444",
  "category": "WARKOP",
  "rank_weight": 10
}
```

**Result:** Instant match (WARKOP ignores rank diff)

---

### Example 4: Rate Limiting

```bash
# First request
curl -X POST http://localhost:3000/api/scrim/request \
  -H "Content-Type: application/json" \
  -d '{"team_name":"Team1","whatsapp_number":"628111","category":"POKE","rank_weight":5}'

# Second request from SAME IP
curl -X POST http://localhost:3000/api/scrim/request \
  -H "Content-Type: application/json" \
  -d '{"team_name":"Team2","whatsapp_number":"628222","category":"POKE","rank_weight":6}'
```

**Result:** Second request returns `"error": "you already have an active request"`

---

## 🎯 PowerShell Examples

### Create POKE Request
```powershell
$body = @{
    team_name = "Team Alpha"
    whatsapp_number = "6281234567890"
    category = "POKE"
    rank_weight = 5
} | ConvertTo-Json

Invoke-RestMethod -Method Post `
    -Uri "http://localhost:3000/api/scrim/request" `
    -ContentType "application/json" `
    -Body $body
```

### Create WARKOP Request (Simulated Different IP)
```powershell
$body = @{
    team_name = "Team Beta"
    whatsapp_number = "6289876543210"
    category = "WARKOP"
    rank_weight = 9
} | ConvertTo-Json

$headers = @{
    "Content-Type" = "application/json"
    "X-Forwarded-For" = "192.168.1.100"
}

Invoke-RestMethod -Method Post `
    -Uri "http://localhost:3000/api/scrim/request" `
    -Headers $headers `
    -Body $body
```

### Check Status
```powershell
$requestId = "ab745781-f0eb-41c5-942b-75a6280fe2c5"
Invoke-RestMethod -Uri "http://localhost:3000/api/scrim/request/$requestId"
```

### Cancel Request
```powershell
$requestId = "ab745781-f0eb-41c5-942b-75a6280fe2c5"
Invoke-RestMethod -Method Post `
    -Uri "http://localhost:3000/api/scrim/request/$requestId/cancel"
```

---

## ⚙️ System Behavior

### Auto-Expiry (30 Minutes)
```
Request created at: 19:45:00
Expires at:         20:15:00

Status changes automatically:
19:45:00 → searching
20:15:01 → expired (cleanup monitor runs every 10s)
```

### Match Timeout (60 Seconds)
```
Match created at:  19:50:00
Expires at:        19:51:00

If not confirmed by 19:51:00:
- Match status → cancelled
- Both teams → back to searching (automatically)
```

### Matching Algorithm

**POKE Category:**
```python
def can_match_poke(team1_rank, team2_rank):
    return abs(team1_rank - team2_rank) <= 2
```

**WARKOP Category:**
```python
def can_match_warkop(team1, team2):
    return team1.category == "WARKOP" and team2.category == "WARKOP"
```

---

## 🔧 Configuration

### Environment Variables
```bash
# Database
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=golobby

# Redis (Rate Limiting)
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=redis123

# Server
APP_PORT=8080
APP_ENV=production
```

### Workers
- **Matchmaking Workers:** 4
- **Cleanup Monitor:** 1 (runs every 10s)

---

## 🚨 Error Codes Reference

| HTTP Status | Error Message | Cause | Solution |
|-------------|---------------|-------|----------|
| 400 | "team_name is required" | Missing field | Provide team_name |
| 400 | "POKE category requires rank_weight between 1 and 8" | Invalid rank | Use correct rank for category |
| 400 | "you already have an active request" | Rate limited | Wait or cancel existing request |
| 404 | "Request not found" | Invalid ID | Check request ID |
| 400 | "match is not pending" | Match already confirmed/expired | Create new request |

---

## 📊 Database Schema

### scrim_requests Table
```sql
id              UUID PRIMARY KEY
team_name       VARCHAR(100)
whatsapp_number VARCHAR(20)
category        scrim_category (POKE/WARKOP)
rank_weight     INTEGER (1-10)
status          scrim_status (searching/matched/expired/cancelled)
match_id        UUID (nullable)
ip_address      VARCHAR(45)
created_at      TIMESTAMP
updated_at      TIMESTAMP
expires_at      TIMESTAMP (created_at + 30 min)
matched_at      TIMESTAMP (nullable)
```

### scrim_matches Table
```sql
id           UUID PRIMARY KEY
team1_id     UUID (FK → scrim_requests)
team2_id     UUID (FK → scrim_requests)
category     scrim_category
rank_diff    INTEGER (nullable, POKE only)
status       VARCHAR(50) (pending/confirmed/cancelled)
created_at   TIMESTAMP
confirmed_at TIMESTAMP (nullable)
expires_at   TIMESTAMP (created_at + 60s)
```

---

## 🎮 Interactive Testing

Use the included test client:
```bash
# Open in browser
start scrim-test-client.html

# Or serve via HTTP server
python -m http.server 8000
# Then visit: http://localhost:8000/scrim-test-client.html
```

**Features:**
- ✅ Two-team simulation
- ✅ POKE/WARKOP category selection
- ✅ Real-time status polling
- ✅ WhatsApp link display
- ✅ Automatic IP spoofing for testing
- ✅ Live logs

---

## 🔐 Security Notes

### Rate Limiting
- **Mechanism:** Redis-based IP tracking
- **TTL:** 30 minutes (same as request expiry)
- **Bypass:** Use `X-Forwarded-For` header (testing only!)

### Production Recommendations
1. Remove `X-Forwarded-For` bypass or add whitelist
2. Use HTTPS for API
3. Implement authentication (JWT/OAuth)
4. Add request signing for WhatsApp URLs
5. Implement CAPTCHA for spam prevention
6. Add database indices for performance
7. Monitor Redis memory usage

---

**API Version:** 1.0.0  
**Last Updated:** 2026-02-07  
**Backend Version:** GoLobby Matchmaking v2.0
