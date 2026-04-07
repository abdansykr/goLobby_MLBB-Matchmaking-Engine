package usecase

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/golobby/matchmaking/internal/domain"
	"github.com/google/uuid"
)

type ScrimMatchmakingUsecase struct {
	requestRepo domain.ScrimRequestRepository
	matchRepo   domain.ScrimMatchRepository
	rateLimiter domain.RateLimiter
	cacheRepo   domain.CacheRepository // Redis — used for atomic match cancellation
	workers     int
	matchQueue  chan *domain.ScrimRequest
	stopChan    chan struct{}
	wg          sync.WaitGroup

	// wsNotifier is an optional callback to send WebSocket events.
	// Injected after construction to avoid circular dependency.
	wsNotifier func(targetRequestID string, payload map[string]interface{})
}

// SetWSNotifier injects the WebSocket notification callback.
func (u *ScrimMatchmakingUsecase) SetWSNotifier(fn func(targetRequestID string, payload map[string]interface{})) {
	u.wsNotifier = fn
}

// notifyWS sends a WebSocket message to the given request ID (no-op if notifier not set).
func (u *ScrimMatchmakingUsecase) notifyWS(targetRequestID string, payload map[string]interface{}) {
	if u.wsNotifier != nil {
		u.wsNotifier(targetRequestID, payload)
	}
}

func NewScrimMatchmakingUsecase(
	requestRepo domain.ScrimRequestRepository,
	matchRepo domain.ScrimMatchRepository,
	rateLimiter domain.RateLimiter,
	cacheRepo domain.CacheRepository,
	workers int,
) *ScrimMatchmakingUsecase {
	return &ScrimMatchmakingUsecase{
		requestRepo: requestRepo,
		matchRepo:   matchRepo,
		rateLimiter: rateLimiter,
		cacheRepo:   cacheRepo,
		workers:     workers,
		matchQueue:  make(chan *domain.ScrimRequest, 100),
		stopChan:    make(chan struct{}),
	}
}

// Start starts the matchmaking workers
func (u *ScrimMatchmakingUsecase) Start() {
	log.Printf("Starting %d scrim matchmaking workers...", u.workers)

	for i := 0; i < u.workers; i++ {
		u.wg.Add(1)
		go u.matchmakingWorker(i)
	}

	// Start periodic scanner for matching requests
	u.wg.Add(1)
	go u.scanAndMatchRequests()

	// Start cleanup monitor
	u.wg.Add(1)
	go u.cleanupMonitor()

	log.Println("Scrim matchmaking workers started successfully")
}

// Stop stops all workers gracefully
func (u *ScrimMatchmakingUsecase) Stop() {
	log.Println("Stopping scrim matchmaking workers...")
	close(u.stopChan)
	u.wg.Wait()
	log.Println("Scrim matchmaking workers stopped")
}

// scanAndMatchRequests periodically scans for searching requests and tries to match them
func (u *ScrimMatchmakingUsecase) scanAndMatchRequests() {
	defer u.wg.Done()
	log.Println("🔍 Scrim periodic scanner started")

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-u.stopChan:
			log.Println("Scrim scanner shutting down")
			return
		case <-ticker.C:
			u.performMatchScan()
		}
	}
}

// performMatchScan scans for searching requests and matches them
func (u *ScrimMatchmakingUsecase) performMatchScan() {
	ctx := context.Background()

	// Get all searching requests
	searchingRequests, err := u.requestRepo.GetAllSearching(ctx)
	if err != nil {
		log.Printf("Error fetching searching requests: %v", err)
		return
	}

	if len(searchingRequests) == 0 {
		return
	}

	log.Printf("🔍 Scanning %d searching requests for matches...", len(searchingRequests))

	// Try to match each request
	for i, request := range searchingRequests {
		log.Printf("  [%d/%d] Checking: %s (ID=%s, Cat=%s, Rank=%d, Status=%s)",
			i+1, len(searchingRequests), request.TeamName, request.ID.String()[:8],
			request.Category, request.RankWeight, request.Status)

		potentialMatches, err := u.requestRepo.FindPotentialMatches(ctx, request)
		if err != nil {
			log.Printf("  ❌ Error finding matches for request %s: %v", request.ID, err)
			continue
		}

		log.Printf("  → Found %d potential opponents", len(potentialMatches))

		if len(potentialMatches) > 0 {
			opponent := potentialMatches[0]
			log.Printf("📍 Found match: %s (rank %d) vs %s (rank %d)",
				request.TeamName, request.RankWeight,
				opponent.TeamName, opponent.RankWeight)

			// Create match immediately
			u.createMatch(request, opponent)
		}
	}
}

// matchmakingWorker processes match requests
func (u *ScrimMatchmakingUsecase) matchmakingWorker(id int) {
	defer u.wg.Done()
	log.Printf("Scrim Worker %d started", id)

	for {
		select {
		case <-u.stopChan:
			log.Printf("Scrim Worker %d shutting down", id)
			return
		case request := <-u.matchQueue:
			u.processMatchRequest(id, request)
		}
	}
}

// processMatchRequest finds and creates a match for a request
func (u *ScrimMatchmakingUsecase) processMatchRequest(workerID int, request *domain.ScrimRequest) {
	ctx := context.Background()

	log.Printf("Scrim Worker %d: Processing request %s (Team: %s, Category: %s, Rank: %d)",
		workerID, request.ID, request.TeamName, request.Category, request.RankWeight)

	// Find potential matches
	potentialMatches, err := u.requestRepo.FindPotentialMatches(ctx, request)
	if err != nil {
		log.Printf("Scrim Worker %d: Error finding matches: %v", workerID, err)
		return
	}

	if len(potentialMatches) == 0 {
		log.Printf("Scrim Worker %d: No matches found for %s, staying in queue", workerID, request.TeamName)
		return
	}

	// Take the first match (FIFO)
	opponent := potentialMatches[0]

	// Verify they can still match (status might have changed)
	if !request.CanMatchWith(opponent) {
		log.Printf("Scrim Worker %d: Match validation failed", workerID)
		return
	}

	// Create the match using helper
	u.createMatch(request, opponent)
}

// createMatch creates a match between two requests (Thread-safe with check)
func (u *ScrimMatchmakingUsecase) createMatch(request, opponent *domain.ScrimRequest) {
	ctx := context.Background()

	// Double-check both are still searching (avoid race condition)
	req1, err := u.requestRepo.GetByID(ctx, request.ID)
	if err != nil || req1.Status != domain.StatusSearching {
		log.Printf("⚠️ Request %s is no longer searching, skipping match", request.ID)
		return
	}

	req2, err := u.requestRepo.GetByID(ctx, opponent.ID)
	if err != nil || req2.Status != domain.StatusSearching {
		log.Printf("⚠️ Opponent %s is no longer searching, skipping match", opponent.ID)
		return
	}

	// Create the match
	match := &domain.ScrimMatch{
		Team1ID:  request.ID,
		Team2ID:  opponent.ID,
		Category: request.Category,
	}

	// Calculate rank diff for POKE category
	if request.Category == domain.CategoryPoke {
		diff := abs(request.RankWeight - opponent.RankWeight)
		match.RankDiff = &diff
	}

	// Save match to database
	err = u.matchRepo.Create(ctx, match)
	if err != nil {
		log.Printf("❌ Error creating match: %v", err)
		return
	}

	// Update both requests with match_id
	err = u.requestRepo.UpdateMatchID(ctx, request.ID, match.ID)
	if err != nil {
		log.Printf("❌ Error updating team1 match_id: %v", err)
		return
	}

	err = u.requestRepo.UpdateMatchID(ctx, opponent.ID, match.ID)
	if err != nil {
		log.Printf("❌ Error updating team2 match_id: %v", err)
		return
	}

	// Initialize Double Opt-in consensus
	if u.cacheRepo != nil {
		ttl := 60 // 60 seconds TTL
		_ = u.cacheRepo.SetMatchPending(ctx, match.ID, request.ID, opponent.ID, ttl)
		_ = u.cacheRepo.InitConsensus(ctx, match.ID, request.ID, opponent.ID, ttl)
	}

	log.Printf("✅ Match created! %s vs %s (Category: %s, Match ID: %s)",
		request.TeamName, opponent.TeamName, match.Category, match.ID)
}

// cleanupMonitor periodically cleans up expired requests and matches
func (u *ScrimMatchmakingUsecase) cleanupMonitor() {
	defer u.wg.Done()
	log.Println("Scrim cleanup monitor started")

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-u.stopChan:
			log.Println("Scrim cleanup monitor shutting down")
			return
		case <-ticker.C:
			u.runCleanup()
		}
	}
}

// runCleanup expires old requests and cancels expired matches
func (u *ScrimMatchmakingUsecase) runCleanup() {
	ctx := context.Background()

	// Expire old requests (30+ minutes)
	expiredCount, err := u.requestRepo.ExpireOldRequests(ctx)
	if err != nil {
		log.Printf("Error expiring old requests: %v", err)
	} else if expiredCount > 0 {
		log.Printf("Expired %d old scrim requests", expiredCount)
	}

	// Cancel expired matches (60+ seconds without confirmation)
	cancelledCount, err := u.matchRepo.CancelExpiredMatches(ctx)
	if err != nil {
		log.Printf("Error cancelling expired matches: %v", err)
	} else if cancelledCount > 0 {
		log.Printf("Cancelled %d expired scrim matches", cancelledCount)
	}
}

// CreateRequest creates a new scrim request
func (u *ScrimMatchmakingUsecase) CreateRequest(ctx context.Context, request *domain.ScrimRequest) (*domain.ScrimRequest, error) {
	// Validate category
	if !request.IsValidCategory() {
		return nil, fmt.Errorf("invalid category: %s", request.Category)
	}

	// Validate rank weight for category
	if !request.IsValidRankWeight() {
		return nil, fmt.Errorf("invalid rank_weight %d for category %s", request.RankWeight, request.Category)
	}

	// Check rate limit (skip if disabled in development)
	enableRateLimit := strings.ToLower(os.Getenv("ENABLE_RATE_LIMIT")) != "false"
	if enableRateLimit {
		canRequest, err := u.rateLimiter.CanRequest(ctx, request.IPAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to check rate limit: %w", err)
		}
		if !canRequest {
			return nil, fmt.Errorf("you already have an active request")
		}
	} else {
		log.Printf("⚠️ Rate limiting is DISABLED (development mode)")
	}

	// Create request in database
	var err error
	err = u.requestRepo.Create(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set rate limit (only if enabled)
	if enableRateLimit {
		err = u.rateLimiter.SetActiveRequest(ctx, request.IPAddress, request.ID)
		if err != nil {
			log.Printf("Warning: failed to set rate limit: %v", err)
		}
	}

	log.Printf("Team %s enqueued for %s scrim (Rank: %d)",
		request.TeamName, request.Category, request.RankWeight)

	// Add to matchmaking queue
	select {
	case u.matchQueue <- request:
		log.Printf("Request %s added to matchmaking queue", request.ID)
	default:
		log.Printf("Warning: matchmaking queue is full, request %s not queued immediately", request.ID)
	}

	return request, nil
}

// GetRequest retrieves a scrim request with match details
func (u *ScrimMatchmakingUsecase) GetRequest(ctx context.Context, id uuid.UUID) (*domain.ScrimRequest, *domain.MatchResponse, error) {
	request, err := u.requestRepo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get request: %w", err)
	}

	// If matched, get match details
	var matchResponse *domain.MatchResponse
	if request.Status == domain.StatusMatched && request.MatchID != nil {
		match, err := u.matchRepo.GetWithTeams(ctx, *request.MatchID)
		if err != nil {
			log.Printf("Warning: failed to get match details: %v", err)
		} else {
			// Determine opponent
			var opponent *domain.ScrimRequest
			if match.Team1ID == request.ID {
				opponent = match.Team2
			} else {
				opponent = match.Team1
			}

			// Security: Do not expose WhatsApp URL or Number until match is confirmed
			waURL := ""
			oppNumber := "TERSEMBUNYI"
			if match.Status == "confirmed" {
				oppNumber = opponent.WhatsAppNumber
				message := fmt.Sprintf("Halo %s! Kita sudah di-match untuk scrim nih. Ayo kita bahas detail pertandingannya!", opponent.TeamName)
				waURL = fmt.Sprintf("https://wa.me/%s?text=%s", opponent.WhatsAppNumber, url.QueryEscape(message))
			}

			expiresIn := int(time.Until(match.ExpiresAt).Seconds())
			if expiresIn < 0 {
				expiresIn = 0
			}

			matchResponse = &domain.MatchResponse{
				Match:             match,
				OpponentName:      opponent.TeamName,
				OpponentNumber:    oppNumber,
				OpponentAvatarURL: opponent.AvatarURL,
				WhatsAppURL:       waURL,
				ExpiresIn:         expiresIn,
			}
		}
	}

	return request, matchResponse, nil
}

// CancelRequest cancels a scrim request.
// If the request is currently in a PENDING match, the match is cancelled and
// the opponent is notified via WebSocket so their modal closes automatically.
func (u *ScrimMatchmakingUsecase) CancelRequest(ctx context.Context, id uuid.UUID) error {
	request, err := u.requestRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get request: %w", err)
	}

	// ── If there is an active match, propagate cancellation ───────────────
	if request.MatchID != nil {
		matchID := *request.MatchID

		// Load the match to find the opponent
		match, err := u.matchRepo.GetWithTeams(ctx, matchID)
		if err == nil && match.Status == "pending" {
			// Cancel the match record
			_ = u.matchRepo.UpdateStatus(ctx, matchID, "cancelled")

			// Find the opponent's request ID
			var opponentRequestID uuid.UUID
			if match.Team1ID == id {
				opponentRequestID = match.Team2ID
			} else {
				opponentRequestID = match.Team1ID
			}

			// Reset the opponent back to "searching" so they can re-queue
			_ = u.requestRepo.UpdateStatus(ctx, opponentRequestID, domain.StatusSearching)

			// Send WebSocket notification to opponent
			go u.notifyWS(opponentRequestID.String(), map[string]interface{}{
				"type":            "MATCH_DECLINED",
				"reason":          "Lawan membatalkan pertandingan ini.",
				"cancelled_by_id": id.String(),
			})

			// Re-enqueue opponent so they can find a new match quickly
			if opponentReq, err2 := u.requestRepo.GetByID(ctx, opponentRequestID); err2 == nil {
				select {
				case u.matchQueue <- opponentReq:
				default:
				}
			}

			log.Printf("⚡ Match %s cancelled — opponent %s re-queued", matchID, opponentRequestID)
		}
	}

	// Cancel the requester's own record
	err = u.requestRepo.Cancel(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to cancel request: %w", err)
	}

	// Remove rate limit
	if request.IPAddress != "" {
		_ = u.rateLimiter.RemoveActiveRequest(ctx, request.IPAddress)
	}

	log.Printf("Request %s cancelled by team %s", id, request.TeamName)
	return nil
}

// DeclineMatch is called when a player explicitly clicks "Tolak Lawan" on the
// match-found modal. It cancels the match and notifies the opponent immediately.
func (u *ScrimMatchmakingUsecase) DeclineMatch(ctx context.Context, matchID uuid.UUID, decliningRequestID uuid.UUID) error {
	match, err := u.matchRepo.GetWithTeams(ctx, matchID)
	if err != nil {
		return fmt.Errorf("failed to get match: %w", err)
	}

	if match.Status != "pending" {
		return fmt.Errorf("match is no longer pending")
	}

	// Cancel the match
	if err := u.matchRepo.UpdateStatus(ctx, matchID, "cancelled"); err != nil {
		return fmt.Errorf("failed to cancel match: %w", err)
	}

	// Determine opponent
	var opponentRequestID uuid.UUID
	if match.Team1ID == decliningRequestID {
		opponentRequestID = match.Team2ID
	} else {
		opponentRequestID = match.Team1ID
	}

	// Reset both requests back to cancelled / searching
	_ = u.requestRepo.Cancel(ctx, decliningRequestID)
	_ = u.requestRepo.UpdateStatus(ctx, opponentRequestID, domain.StatusSearching)

	// Notify opponent via WebSocket
	go u.notifyWS(opponentRequestID.String(), map[string]interface{}{
		"type":   "MATCH_DECLINED",
		"reason": "Lawan menolak tawaran pertandingan ini.",
	})

	// Re-enqueue opponent
	if opponentReq, err2 := u.requestRepo.GetByID(ctx, opponentRequestID); err2 == nil {
		select {
		case u.matchQueue <- opponentReq:
		default:
		}
	}

	log.Printf("⚡ Match %s declined by request %s — opponent %s notified", matchID, decliningRequestID, opponentRequestID)
	return nil
}

// RejectMatch is the authoritative, concurrency-safe match rejection endpoint.
//
// Design goals:
//  1. Atomic — uses a Redis Lua script so exactly ONE goroutine performs the
//     cancellation even if both players click "Reject" at the same millisecond.
//  2. Broadcast — sends MATCH_CANCELLED to BOTH participants via WebSocket.
//  3. Re-queue  — the non-rejecting participant is automatically placed back
//     into the matchmaking queue so they find a new opponent without manual action.
//  4. DB consistency — the PostgreSQL match row is updated to 'cancelled' and
//     both scrim_request rows are reset to 'searching'.
//
// Endpoint: POST /api/scrim/match/:id/reject?request_id=<uuid>
func (u *ScrimMatchmakingUsecase) RejectMatch(ctx context.Context, matchID uuid.UUID, rejectingRequestID uuid.UUID) error {
	// ── Step 1: Atomic Redis cancellation ────────────────────────────────────
	// Only the first caller wins; concurrent callers get (false, nil) and return
	// without doing any duplicate DB writes or duplicate WS broadcasts.
	var didCancel bool
	var err error

	if u.cacheRepo != nil {
		didCancel, err = u.cacheRepo.CancelMatchAtomically(ctx, matchID)
		if err != nil {
			log.Printf("⚠️  RejectMatch: Redis atomic cancel error (falling back to DB): %v", err)
			// Do NOT return — fall back to DB-only path below
			didCancel = true // assume we are the one to proceed
		}
		if !didCancel {
			// Another goroutine already handled the cancellation.
			// This call is idempotent: just ensure our own request is cancelled.
			log.Printf("ℹ️  RejectMatch: match %s already cancelled by concurrent request", matchID)
			_ = u.requestRepo.Cancel(ctx, rejectingRequestID)
			_ = u.rateLimiter.RemoveActiveRequest(ctx, "")
			return nil
		}
	} else {
		didCancel = true
	}

	// ── Step 2: Load match with team data from DB ─────────────────────────────
	match, err := u.matchRepo.GetWithTeams(ctx, matchID)
	if err != nil {
		return fmt.Errorf("failed to load match %s: %w", matchID, err)
	}

	// Allow cancellation only if match is still pending
	if match.Status != "pending" {
		log.Printf("ℹ️  RejectMatch: match %s is no longer pending (status=%s)", matchID, match.Status)
		_ = u.requestRepo.Cancel(ctx, rejectingRequestID)
		return nil
	}

	// ── Step 3: Update DB — cancel match, reset both requests ─────────────────
	_ = u.matchRepo.UpdateStatus(ctx, matchID, "cancelled")

	// Determine who is the "other" participant (non-rejecting)
	var otherRequestID uuid.UUID
	if match.Team1ID == rejectingRequestID {
		otherRequestID = match.Team2ID
	} else {
		otherRequestID = match.Team1ID
	}

	// Cancel the rejecting team's request
	_ = u.requestRepo.Cancel(ctx, rejectingRequestID)
	// Reset the other team back to searching
	_ = u.requestRepo.UpdateStatus(ctx, otherRequestID, domain.StatusSearching)

	// ── Step 4: Broadcast MATCH_CANCELLED to BOTH participants ────────────────
	cancelPayload := map[string]interface{}{
		"type":    "MATCH_CANCELLED",
		"reason":  "Salah satu tim menolak pertandingan ini.",
		"matchId": matchID.String(),
	}
	go u.notifyWS(rejectingRequestID.String(), cancelPayload)
	go u.notifyWS(otherRequestID.String(), cancelPayload)

	// ── Step 5: Re-enqueue the other team ────────────────────────────────────
	if otherReq, err2 := u.requestRepo.GetByID(ctx, otherRequestID); err2 == nil {
		select {
		case u.matchQueue <- otherReq:
			log.Printf("♻️  RejectMatch: re-queued team %s", otherRequestID)
		default:
			log.Printf("⚠️  RejectMatch: queue full, team %s not immediately re-queued", otherRequestID)
		}
	}

	// ── Step 6: Clean up Redis keys ───────────────────────────────────────────
	if u.cacheRepo != nil {
		_ = u.cacheRepo.DeleteMatch(ctx, matchID)
	}

	log.Printf("✅ Match %s rejected by request %s — both parties notified", matchID, rejectingRequestID)
	return nil
}

// ConfirmMatch confirms a match (Double Opt-in consensus)
func (u *ScrimMatchmakingUsecase) ConfirmMatch(ctx context.Context, matchID uuid.UUID, requestID uuid.UUID) error {
	match, err := u.matchRepo.GetWithTeams(ctx, matchID)
	if err != nil {
		return fmt.Errorf("failed to get match: %w", err)
	}

	if match.Status != "pending" {
		if match.Status == "confirmed" {
			log.Printf("ℹ️  Match %s is already confirmed, ignoring Accept.", matchID)
			return nil
		}
		return fmt.Errorf("match is not pending")
	}

	var allAccepted bool
	if u.cacheRepo != nil {
		allAccepted, err = u.cacheRepo.RecordAcceptance(ctx, matchID, requestID)
		if err != nil {
			log.Printf("⚠️  ConfirmMatch: Redis error recording acceptance: %v", err)
			return fmt.Errorf("redis error for RecordAcceptance: %w", err)
		}
	} else {
		allAccepted = true
	}

	var otherRequestID uuid.UUID
	if match.Team1ID == requestID {
		otherRequestID = match.Team2ID
	} else {
		otherRequestID = match.Team1ID
	}

	u.notifyWS(otherRequestID.String(), map[string]interface{}{
		"type":    "OPPONENT_ACCEPTED",
		"message": "Lawan sudah menyetujui, menunggu Anda!",
	})

	if allAccepted {
		// Both participants have accepted! Update the database.
		err = u.matchRepo.Confirm(ctx, matchID)
		if err != nil {
			return fmt.Errorf("failed to confirm match in database: %w", err)
		}

		// Notify both parties with MATCH_SUCCESS and WhatsApp URLs
		msg1 := fmt.Sprintf("Halo %s! Kita sudah di-match untuk scrim nih. Ayo kita bahas detail pertandingannya!", match.Team2.TeamName)
		waURL1 := fmt.Sprintf("https://wa.me/%s?text=%s", match.Team2.WhatsAppNumber, url.QueryEscape(msg1))

		msg2 := fmt.Sprintf("Halo %s! Kita sudah di-match untuk scrim nih. Ayo kita bahas detail pertandingannya!", match.Team1.TeamName)
		waURL2 := fmt.Sprintf("https://wa.me/%s?text=%s", match.Team1.WhatsAppNumber, url.QueryEscape(msg2))

		go u.notifyWS(match.Team1ID.String(), map[string]interface{}{
			"type":         "MATCH_SUCCESS",
			"whatsapp_url": waURL1,
		})
		go u.notifyWS(match.Team2ID.String(), map[string]interface{}{
			"type":         "MATCH_SUCCESS",
			"whatsapp_url": waURL2,
		})

		log.Printf("✨ Match %s confirmed via Double Opt-in!", matchID)
	} else {
		log.Printf("➡️  Participant %s accepted match %s, waiting for opponent...", requestID, matchID)
	}

	return nil
}

// Helper function
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
