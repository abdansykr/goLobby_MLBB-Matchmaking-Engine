package usecase

import (
	"context"
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"github.com/golobby/matchmaking/internal/domain"
	"github.com/google/uuid"
)

// MatchmakingConfig holds configuration for the matchmaking system
type MatchmakingConfig struct {
	InitialRankRange  int // ±2 points initially
	ExtendedRankRange int // ±4 points after timeout
	MatchTimeout      int // 30 seconds before extending range
	ReadyTimeout      int // 60 seconds for ready confirmation
	GhostingPenalty   int // -10 reputation points
}

// MatchmakingUsecase handles the business logic for matchmaking
type MatchmakingUsecase struct {
	teamRepo  domain.TeamRepository
	matchRepo domain.MatchRepository
	cache     domain.CacheRepository
	config    MatchmakingConfig
	
	// Channels for concurrent operations
	matchmakerChan chan *domain.Team
	resultChan     chan *MatchResult
	stopChan       chan struct{}
	
	// WebSocket broadcast channel
	broadcastChan chan *BroadcastMessage
	
	// Thread-safe map of active teams
	activeTeams sync.Map
	
	// Worker pool
	workerPool *sync.WaitGroup
}

// MatchResult represents the result of a matchmaking operation
type MatchResult struct {
	Match *domain.Match
	Team1 *domain.Team
	Team2 *domain.Team
	Error error
}

// BroadcastMessage represents a message to broadcast via WebSocket
type BroadcastMessage struct {
	Type      string      `json:"type"`
	TeamID    uuid.UUID   `json:"team_id,omitempty"`
	MatchID   uuid.UUID   `json:"match_id,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// NewMatchmakingUsecase creates a new matchmaking usecase
func NewMatchmakingUsecase(
	teamRepo domain.TeamRepository,
	matchRepo domain.MatchRepository,
	cache domain.CacheRepository,
	config MatchmakingConfig,
) *MatchmakingUsecase {
	uc := &MatchmakingUsecase{
		teamRepo:       teamRepo,
		matchRepo:      matchRepo,
		cache:          cache,
		config:         config,
		matchmakerChan: make(chan *domain.Team, 100),
		resultChan:     make(chan *MatchResult, 100),
		stopChan:       make(chan struct{}),
		broadcastChan:  make(chan *BroadcastMessage, 100),
		workerPool:     &sync.WaitGroup{},
	}
	
	return uc
}

// StartMatchmakingWorkers starts the matchmaking worker goroutines
func (uc *MatchmakingUsecase) StartMatchmakingWorkers(numWorkers int) {
	log.Printf("Starting %d matchmaking workers...", numWorkers)
	
	for i := 0; i < numWorkers; i++ {
		uc.workerPool.Add(1)
		go uc.matchmakingWorker(i)
	}
	
	// Start the ghosting monitor
	uc.workerPool.Add(1)
	go uc.ghostingMonitor()
	
	log.Println("Matchmaking workers started successfully")
}

// StopMatchmakingWorkers gracefully stops all workers
func (uc *MatchmakingUsecase) StopMatchmakingWorkers() {
	log.Println("Stopping matchmaking workers...")
	close(uc.stopChan)
	uc.workerPool.Wait()
	log.Println("All workers stopped")
}

// GetBroadcastChannel returns the broadcast channel for WebSocket updates
func (uc *MatchmakingUsecase) GetBroadcastChannel() <-chan *BroadcastMessage {
	return uc.broadcastChan
}

// EnqueueTeam adds a team to the matchmaking queue
func (uc *MatchmakingUsecase) EnqueueTeam(ctx context.Context, team *domain.Team) error {
	// Check if team is already locked
	locked, err := uc.cache.IsTeamLocked(ctx, team.ID)
	if err != nil {
		return fmt.Errorf("failed to check team lock: %w", err)
	}
	if locked {
		return fmt.Errorf("team is already in a match")
	}
	
	// Create team in database
	if err := uc.teamRepo.Create(ctx, team); err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}
	
	// Add to Redis queue for fast matching
	if err := uc.cache.EnqueueTeam(ctx, team); err != nil {
		return fmt.Errorf("failed to enqueue team: %w", err)
	}
	
	// Track active team
	uc.activeTeams.Store(team.ID.String(), team)
	
	// Send to matchmaking channel
	select {
	case uc.matchmakerChan <- team:
		log.Printf("Team %s enqueued for matchmaking (Rank: %d)", team.TeamName, team.AverageRank)
	default:
		return fmt.Errorf("matchmaking queue is full")
	}
	
	return nil
}

// matchmakingWorker is the core concurrent matchmaking logic
func (uc *MatchmakingUsecase) matchmakingWorker(workerID int) {
	defer uc.workerPool.Done()
	
	log.Printf("Worker %d started", workerID)
	
	for {
		select {
		case <-uc.stopChan:
			log.Printf("Worker %d stopping...", workerID)
			return
			
		case team := <-uc.matchmakerChan:
			// Process matchmaking with timeout and dynamic range scaling
			uc.processMatchmaking(workerID, team)
		}
	}
}

// processMatchmaking handles the smart matchmaking algorithm
func (uc *MatchmakingUsecase) processMatchmaking(workerID int, team *domain.Team) {
	ctx := context.Background()
	startTime := time.Now()
	
	log.Printf("Worker %d: Processing team %s (Rank: %d)", workerID, team.TeamName, team.AverageRank)
	
	// Initial search with ±2 rank range
	opponent := uc.findOpponent(ctx, team, uc.config.InitialRankRange)
	
	if opponent == nil {
		// Wait for timeout, then extend search range
		time.Sleep(time.Duration(uc.config.MatchTimeout) * time.Second)
		
		elapsed := time.Since(startTime).Seconds()
		log.Printf("Worker %d: No match found for %s in %.0fs, extending range to ±%d", 
			workerID, team.TeamName, elapsed, uc.config.ExtendedRankRange)
		
		// Extended search with ±4 rank range
		opponent = uc.findOpponent(ctx, team, uc.config.ExtendedRankRange)
	}
	
	if opponent != nil {
		// Create match
		match := uc.createMatch(ctx, team, opponent)
		if match != nil {
			elapsed := time.Since(startTime).Seconds()
			log.Printf("Worker %d: Match created! %s vs %s (found in %.0fs)", 
				workerID, team.TeamName, opponent.TeamName, elapsed)
			
			// Broadcast match found notification
			uc.broadcastMatchFound(match, team, opponent)
		}
	} else {
		// Re-enqueue team for next attempt
		log.Printf("Worker %d: No opponent found for %s, re-enqueueing", workerID, team.TeamName)
		_ = uc.cache.EnqueueTeam(ctx, team)
		
		// Send back to channel for retry
		select {
		case uc.matchmakerChan <- team:
		default:
			log.Printf("Worker %d: Failed to re-enqueue team %s", workerID, team.TeamName)
		}
	}
}

// findOpponent searches for a suitable opponent within rank range
func (uc *MatchmakingUsecase) findOpponent(ctx context.Context, team *domain.Team, rankRange int) *domain.Team {
	queueLength, err := uc.cache.GetQueueLength(ctx)
	if err != nil || queueLength == 0 {
		return nil
	}
	
	// Try to find a match from queue
	for i := int64(0); i < queueLength; i++ {
		candidate, err := uc.cache.DequeueTeam(ctx)
		if err != nil || candidate == nil {
			continue
		}
		
		// Skip if it's the same team
		if candidate.ID == team.ID {
			continue
		}
		
		// Check if candidate is locked
		locked, _ := uc.cache.IsTeamLocked(ctx, candidate.ID)
		if locked {
			continue
		}
		
		// Check rank difference
		rankDiff := int(math.Abs(float64(team.AverageRank - candidate.AverageRank)))
		if rankDiff <= rankRange {
			return candidate
		}
		
		// Re-enqueue if not matched
		_ = uc.cache.EnqueueTeam(ctx, candidate)
	}
	
	return nil
}

// createMatch creates a new match and locks both teams
func (uc *MatchmakingUsecase) createMatch(ctx context.Context, team1, team2 *domain.Team) *domain.Match {
	rankDiff := int(math.Abs(float64(team1.AverageRank - team2.AverageRank)))
	match := domain.NewMatch(team1.ID, team2.ID, rankDiff, uc.config.ReadyTimeout)
	
	// Create match in database
	if err := uc.matchRepo.Create(ctx, match); err != nil {
		log.Printf("Failed to create match: %v", err)
		return nil
	}
	
	// Lock both teams in Redis (anti-ghosting)
	if err := uc.cache.LockTeam(ctx, team1.ID, match.ID, uc.config.ReadyTimeout); err != nil {
		log.Printf("Failed to lock team1: %v", err)
		return nil
	}
	
	if err := uc.cache.LockTeam(ctx, team2.ID, match.ID, uc.config.ReadyTimeout); err != nil {
		log.Printf("Failed to lock team2: %v", err)
		_ = uc.cache.UnlockTeam(ctx, team1.ID)
		return nil
	}
	
	// Update team statuses
	team1.Lock(match.ID)
	team2.Lock(match.ID)
	
	_ = uc.teamRepo.Update(ctx, team1)
	_ = uc.teamRepo.Update(ctx, team2)
	
	// Set match as pending in cache (also stores both team IDs for conflict-free cancellation)
	_ = uc.cache.SetMatchPending(ctx, match.ID, team1.ID, team2.ID, uc.config.ReadyTimeout)
	
	return match
}

// broadcastMatchFound sends WebSocket notification to both captains
func (uc *MatchmakingUsecase) broadcastMatchFound(match *domain.Match, team1, team2 *domain.Team) {
	// Notify team 1
	uc.broadcastChan <- &BroadcastMessage{
		Type:      "MATCH_FOUND",
		TeamID:    team1.ID,
		MatchID:   match.ID,
		Data: map[string]interface{}{
			"opponent_name": team2.TeamName,
			"opponent_rank": team2.AverageRank,
			"match_id":      match.ID,
			"expires_at":    match.ExpiresAt,
		},
		Timestamp: time.Now(),
	}
	
	// Notify team 2
	uc.broadcastChan <- &BroadcastMessage{
		Type:      "MATCH_FOUND",
		TeamID:    team2.ID,
		MatchID:   match.ID,
		Data: map[string]interface{}{
			"opponent_name": team1.TeamName,
			"opponent_rank": team1.AverageRank,
			"match_id":      match.ID,
			"expires_at":    match.ExpiresAt,
		},
		Timestamp: time.Now(),
	}
}

// ConfirmReady marks a team as ready for the match
func (uc *MatchmakingUsecase) ConfirmReady(ctx context.Context, teamID, matchID uuid.UUID) error {
	// Get team
	team, err := uc.teamRepo.GetByID(ctx, teamID)
	if err != nil {
		return fmt.Errorf("team not found: %w", err)
	}
	
	// Verify team is locked with this match
	lockedMatchID, err := uc.cache.GetTeamLock(ctx, teamID)
	if err != nil || lockedMatchID == nil || *lockedMatchID != matchID {
		return fmt.Errorf("team is not locked to this match")
	}
	
	// Mark as ready
	team.MarkAsReady()
	if err := uc.teamRepo.Update(ctx, team); err != nil {
		return fmt.Errorf("failed to update team: %w", err)
	}
	
	log.Printf("Team %s confirmed ready for match %s", team.TeamName, matchID)
	
	// Check if both teams are ready
	match, err := uc.matchRepo.GetByID(ctx, matchID)
	if err != nil {
		return err
	}
	
	team1, _ := uc.teamRepo.GetByID(ctx, match.Team1ID)
	team2, _ := uc.teamRepo.GetByID(ctx, match.Team2ID)
	
	if team1.Status == domain.TeamStatusReady && team2.Status == domain.TeamStatusReady {
		match.Confirm()
		_ = uc.matchRepo.Update(ctx, match)
		
		// Unlock teams
		_ = uc.cache.UnlockTeam(ctx, team1.ID)
		_ = uc.cache.UnlockTeam(ctx, team2.ID)
		
		log.Printf("Match %s confirmed! Both teams ready", matchID)
	}
	
	return nil
}

// ghostingMonitor monitors for expired matches and applies penalties
func (uc *MatchmakingUsecase) ghostingMonitor() {
	defer uc.workerPool.Done()
	
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	
	log.Println("Ghosting monitor started")
	
	for {
		select {
		case <-uc.stopChan:
			log.Println("Ghosting monitor stopping...")
			return
			
		case <-ticker.C:
			uc.checkExpiredMatches()
		}
	}
}

// checkExpiredMatches finds and handles expired matches
func (uc *MatchmakingUsecase) checkExpiredMatches() {
	ctx := context.Background()
	
	expiredMatches, err := uc.matchRepo.GetExpiredMatches(ctx)
	if err != nil {
		log.Printf("Failed to get expired matches: %v", err)
		return
	}
	
	for _, match := range expiredMatches {
		log.Printf("Processing expired match: %s", match.ID)
		
		team1, _ := uc.teamRepo.GetByID(ctx, match.Team1ID)
		team2, _ := uc.teamRepo.GetByID(ctx, match.Team2ID)
		
		// Apply penalty to teams that didn't confirm
		if team1 != nil && team1.Status != domain.TeamStatusReady {
			team1.ApplyGhostingPenalty(uc.config.GhostingPenalty)
			_ = uc.teamRepo.Update(ctx, team1)
			log.Printf("Applied ghosting penalty to team: %s (new score: %d)", team1.TeamName, team1.ReputationScore)
		}
		
		if team2 != nil && team2.Status != domain.TeamStatusReady {
			team2.ApplyGhostingPenalty(uc.config.GhostingPenalty)
			_ = uc.teamRepo.Update(ctx, team2)
			log.Printf("Applied ghosting penalty to team: %s (new score: %d)", team2.TeamName, team2.ReputationScore)
		}
		
		// Cancel match
		match.Cancel()
		_ = uc.matchRepo.Update(ctx, match)
		
		// Unlock teams
		_ = uc.cache.UnlockTeam(ctx, match.Team1ID)
		_ = uc.cache.UnlockTeam(ctx, match.Team2ID)
		
		// Remove from cache
		_ = uc.cache.DeleteMatch(ctx, match.ID)
	}
}

// CancelMatchmaking cancels matchmaking for a team
func (uc *MatchmakingUsecase) CancelMatchmaking(ctx context.Context, teamID uuid.UUID) error {
	// Remove from queue
	if err := uc.cache.RemoveFromQueue(ctx, teamID); err != nil {
		log.Printf("Failed to remove from queue: %v", err)
	}
	
	// Get and update team
	team, err := uc.teamRepo.GetByID(ctx, teamID)
	if err != nil {
		return err
	}
	
	team.Cancel()
	return uc.teamRepo.Update(ctx, team)
}
