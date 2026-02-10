# ✅ Deployment Checklist

Use this checklist before deploying to production.

## 🔍 Pre-Deployment Validation

### 1. Configuration ✅

- [ ] `.env` file exists and is configured
- [ ] `DB_PASSWORD` is set to a strong password (not "postgres")
- [ ] `REDIS_PASSWORD` is set to a strong password (not "redis123")
- [ ] `JWT_SECRET` is set to a cryptographically random string
- [ ] `APP_ENV` is set to "production"
- [ ] All configuration values are appropriate for production

**Command to check:**
```bash
cat .env | grep -E "PASSWORD|SECRET"
```

---

### 2. Security ✅

- [ ] `.env` is in `.gitignore` (secrets not committed)
- [ ] CORS origins are restricted (not "*")
- [ ] Database uses SSL in production
- [ ] Non-root Docker user is configured
- [ ] Health check endpoint is accessible
- [ ] No sensitive data in logs

**Files to review:**
- `.gitignore`
- `internal/delivery/http/handler.go` (CORS config)

---

### 3. Build & Dependencies ✅

- [ ] `go mod tidy` runs without errors
- [ ] `go build ./cmd/server` succeeds
- [ ] All Go files compile
- [ ] No unused dependencies
- [ ] Dependencies are up to date

**Commands:**
```bash
go mod tidy
go build -o matchmaking ./cmd/server
go test ./...
```

---

### 4. Docker ✅

- [ ] Dockerfile builds successfully
- [ ] Image size is reasonable (< 30MB)
- [ ] Multi-stage build is working
- [ ] Health check is configured
- [ ] Non-root user is set

**Commands:**
```bash
docker build -t matchmaking:test .
docker images matchmaking:test
docker run --rm matchmaking:test /app/matchmaking --help
```

---

### 5. Docker Compose ✅

- [ ] `docker-compose.yaml` is valid
- [ ] All services start successfully
- [ ] Services pass health checks
- [ ] Volumes are configured for persistence
- [ ] Networks are properly configured

**Commands:**
```bash
docker-compose config
docker-compose up -d
docker-compose ps
docker-compose logs app
```

---

### 6. Database ✅

- [ ] Migrations run successfully
- [ ] Tables are created correctly
- [ ] Indexes exist on critical columns
- [ ] Constraints are in place
- [ ] Connection pooling is configured

**Commands:**
```bash
docker exec -it antigravity_postgres psql -U postgres -d antigravity -c "\dt"
docker exec -it antigravity_postgres psql -U postgres -d antigravity -c "\d teams"
docker exec -it antigravity_postgres psql -U postgres -d antigravity -c "\d matches"
```

---

### 7. Redis ✅

- [ ] Redis is accessible
- [ ] Password authentication works
- [ ] Queue operations function
- [ ] Lock mechanism works
- [ ] TTL expiration is correct

**Commands:**
```bash
docker exec -it antigravity_redis redis-cli -a redis123 PING
docker exec -it antigravity_redis redis-cli -a redis123 INFO
```

---

### 8. API Endpoints ✅

- [ ] Health check responds
- [ ] Enqueue endpoint works
- [ ] Ready confirmation works
- [ ] Cancel works
- [ ] WebSocket connects

**Commands:**
```bash
# Health check
curl http://localhost:8080/health

# Enqueue test
curl -X POST http://localhost:8080/api/matchmaking/enqueue \
  -H "Content-Type: application/json" \
  -d '{"captain_name":"Test","team_name":"Test Team","average_rank":75}'
```

---

### 9. Matchmaking Logic ✅

- [ ] Teams match within ±2 rank initially
- [ ] Range expands to ±4 after 30 seconds
- [ ] FIFO ordering is maintained
- [ ] Concurrent matching works
- [ ] No duplicate matches

**Test:**
```bash
# Create two teams with ranks 75 and 76 (within ±2)
# They should match immediately
```

---

### 10. Anti-Ghosting ✅

- [ ] Teams are locked when matched
- [ ] Locks expire after 60 seconds
- [ ] Reputation penalties apply
- [ ] Ghosting monitor is running
- [ ] Expired matches are cleaned up

**Test:**
- Enqueue two teams
- Don't confirm ready
- Wait 60 seconds
- Check logs for penalty message

---

### 11. WebSocket ✅

- [ ] WebSocket connections establish
- [ ] Match found notifications sent
- [ ] Multiple concurrent connections work
- [ ] Connections handle errors gracefully
- [ ] No memory leaks

**Test:**
- Open `test-client.html`
- Create team
- Verify WebSocket connection in browser DevTools
- Check for "Match Found" notification

---

### 12. Concurrency ✅

- [ ] Workers start correctly (4 by default)
- [ ] No race conditions detected
- [ ] Channels don't deadlock
- [ ] Graceful shutdown works
- [ ] No goroutine leaks

**Test:**
```bash
# Run with race detector
go run -race cmd/server/main.go

# Send SIGTERM
kill -TERM <pid>
# Should see "Shutting down gracefully..."
```

---

### 13. Logging ✅

- [ ] Logs are structured
- [ ] Important events are logged
- [ ] No sensitive data in logs
- [ ] Log levels are appropriate
- [ ] Logs are accessible

**Commands:**
```bash
docker-compose logs app | grep "Match created"
docker-compose logs app | grep "ERROR"
```

---

### 14. Performance ✅

- [ ] Response times < 100ms
- [ ] Memory usage stable
- [ ] CPU usage reasonable
- [ ] Database queries optimized
- [ ] Connection pooling working

**Test:**
```bash
# Load test with 100 teams
for i in {1..100}; do
  curl -X POST http://localhost:8080/api/matchmaking/enqueue \
    -H "Content-Type: application/json" \
    -d "{\"captain_name\":\"P$i\",\"team_name\":\"T$i\",\"average_rank\":$((RANDOM%100))}" &
done
```

---

### 15. Documentation ✅

- [ ] README.md is complete
- [ ] QUICKSTART.md is accurate
- [ ] API_TESTING.md has examples
- [ ] ARCHITECTURE.md explains design
- [ ] All code is commented

**Files to review:**
- README.md
- QUICKSTART.md
- API_TESTING.md
- ARCHITECTURE.md

---

## 🚀 Production Deployment Steps

### Step 1: Environment Setup
```bash
# Copy and configure production .env
cp .env .env.production
nano .env.production  # Edit for production

# Update:
# - DB_PASSWORD (strong password)
# - REDIS_PASSWORD (strong password)
# - JWT_SECRET (random 32+ characters)
# - APP_ENV=production
```

### Step 2: Build Production Image
```bash
docker build -t matchmaking:v1.0.0 .
docker tag matchmaking:v1.0.0 matchmaking:latest
```

### Step 3: Push to Registry (if using)
```bash
docker tag matchmaking:v1.0.0 your-registry/matchmaking:v1.0.0
docker push your-registry/matchmaking:v1.0.0
```

### Step 4: Deploy
```bash
# Using Docker Compose
docker-compose -f docker-compose.yaml --env-file .env.production up -d

# Or deploy to cloud platform
# (AWS ECS, GCP Cloud Run, Azure Container Instances, etc.)
```

### Step 5: Verify Deployment
```bash
# Health check
curl https://your-domain.com/health

# Check logs
docker-compose logs -f app

# Monitor for 5 minutes
watch -n 5 'curl -s https://your-domain.com/health'
```

---

## 🛡️ Post-Deployment Monitoring

### Essential Metrics to Track

1. **Application Health**
   - Health endpoint status
   - Response times
   - Error rates

2. **Matchmaking Performance**
   - Average match time
   - Queue length
   - Match success rate
   - Ghost rate (%)

3. **System Resources**
   - CPU usage
   - Memory usage
   - Database connections
   - Redis memory

4. **Business Metrics**
   - Active teams
   - Matches per hour
   - Average reputation score
   - WebSocket connections

---

## 🔥 Emergency Rollback Plan

If issues occur in production:

```bash
# 1. Stop current deployment
docker-compose down

# 2. Switch back to previous version
docker tag matchmaking:v0.9.0 matchmaking:latest

# 3. Restart with old version
docker-compose up -d

# 4. Verify rollback
curl http://localhost:8080/health
```

---

## 📞 Production Support Contacts

- **Database Issues**: Check PostgreSQL logs
- **Queue Issues**: Check Redis logs
- **Application Errors**: Check app logs
- **WebSocket Issues**: Check browser DevTools

**Log Commands:**
```bash
docker-compose logs postgres
docker-compose logs redis
docker-compose logs app
```

---

## ✅ Final Checklist

Before going live:

- [ ] All tests above passed
- [ ] Documentation is up to date
- [ ] Team is trained on operations
- [ ] Monitoring is set up
- [ ] Backup strategy is in place
- [ ] Rollback plan is tested
- [ ] Security review completed
- [ ] Performance testing done
- [ ] Load testing completed
- [ ] Production .env configured

---

**🎉 Once all checkboxes are ticked, you're ready for production deployment!**
