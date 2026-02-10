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
	workers     int
	matchQueue  chan *domain.ScrimRequest
	stopChan    chan struct{}
	wg          sync.WaitGroup
}

func NewScrimMatchmakingUsecase(
	requestRepo domain.ScrimRequestRepository,
	matchRepo domain.ScrimMatchRepository,
	rateLimiter domain.RateLimiter,
	workers int,
) *ScrimMatchmakingUsecase {
	return &ScrimMatchmakingUsecase{
		requestRepo: requestRepo,
		matchRepo:   matchRepo,
		rateLimiter: rateLimiter,
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

			// Generate WhatsApp URL
			message := fmt.Sprintf("Hi %s! We've been matched for a scrim. Let's coordinate the match details!", opponent.TeamName)
			waURL := fmt.Sprintf("https://wa.me/%s?text=%s", opponent.WhatsAppNumber, url.QueryEscape(message))

			expiresIn := int(time.Until(match.ExpiresAt).Seconds())
			if expiresIn < 0 {
				expiresIn = 0
			}

			matchResponse = &domain.MatchResponse{
				Match:          match,
				OpponentName:   opponent.TeamName,
				OpponentNumber: opponent.WhatsAppNumber,
				WhatsAppURL:    waURL,
				ExpiresIn:      expiresIn,
			}
		}
	}

	return request, matchResponse, nil
}

// CancelRequest cancels a scrim request
func (u *ScrimMatchmakingUsecase) CancelRequest(ctx context.Context, id uuid.UUID) error {
	request, err := u.requestRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get request: %w", err)
	}

	// Cancel the request
	err = u.requestRepo.Cancel(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to cancel request: %w", err)
	}

	// Remove rate limit
	if request.IPAddress != "" {
		err = u.rateLimiter.RemoveActiveRequest(ctx, request.IPAddress)
		if err != nil {
			log.Printf("Warning: failed to remove rate limit: %v", err)
		}
	}

	log.Printf("Request %s cancelled by team %s", id, request.TeamName)

	return nil
}

// ConfirmMatch confirms a match
func (u *ScrimMatchmakingUsecase) ConfirmMatch(ctx context.Context, matchID uuid.UUID) error {
	match, err := u.matchRepo.GetByID(ctx, matchID)
	if err != nil {
		return fmt.Errorf("failed to get match: %w", err)
	}

	if match.Status != "pending" {
		return fmt.Errorf("match is not pending")
	}

	err = u.matchRepo.Confirm(ctx, matchID)
	if err != nil {
		return fmt.Errorf("failed to confirm match: %w", err)
	}

	log.Printf("Match %s confirmed", matchID)

	return nil
}

// Helper function
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
