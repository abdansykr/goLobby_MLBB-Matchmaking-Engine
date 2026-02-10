package http

import (
	"log"

	"github.com/golobby/matchmaking/internal/domain"
	"github.com/golobby/matchmaking/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

// MatchmakingHandler handles HTTP requests for matchmaking
type MatchmakingHandler struct {
	usecase *usecase.MatchmakingUsecase
	hub     *WebSocketHub
}

// NewMatchmakingHandler creates a new matchmaking handler
func NewMatchmakingHandler(uc *usecase.MatchmakingUsecase, hub *WebSocketHub) *MatchmakingHandler {
	return &MatchmakingHandler{
		usecase: uc,
		hub:     hub,
	}
}

// EnqueueTeamRequest represents the request to enqueue a team
type EnqueueTeamRequest struct {
	CaptainID   string `json:"captain_id"`
	CaptainName string `json:"captain_name"`
	TeamName    string `json:"team_name"`
	AverageRank int    `json:"average_rank"`
}

// EnqueueTeam handles POST /api/matchmaking/enqueue
func (h *MatchmakingHandler) EnqueueTeam(c *fiber.Ctx) error {
	var req EnqueueTeamRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate input
	if req.TeamName == "" || req.CaptainName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Team name and captain name are required",
		})
	}

	if req.AverageRank < 0 || req.AverageRank > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Average rank must be between 0 and 100",
		})
	}

	// Parse captain ID
	captainID, err := uuid.Parse(req.CaptainID)
	if err != nil {
		captainID = uuid.New() // Generate new if not provided
	}

	// Create team
	team := &domain.Team{
		ID:              uuid.New(),
		CaptainID:       captainID,
		CaptainName:     req.CaptainName,
		TeamName:        req.TeamName,
		AverageRank:     req.AverageRank,
		Status:          domain.TeamStatusWaiting,
		ReputationScore: 100,
	}

	// Enqueue team
	if err := h.usecase.EnqueueTeam(c.Context(), team); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Team enqueued successfully",
		"team_id": team.ID,
		"status":  "WAITING",
	})
}

// ConfirmReadyRequest represents the ready confirmation request
type ConfirmReadyRequest struct {
	TeamID  string `json:"team_id"`
	MatchID string `json:"match_id"`
}

// ConfirmReady handles POST /api/matchmaking/ready
func (h *MatchmakingHandler) ConfirmReady(c *fiber.Ctx) error {
	var req ConfirmReadyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	teamID, err := uuid.Parse(req.TeamID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid team ID",
		})
	}

	matchID, err := uuid.Parse(req.MatchID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid match ID",
		})
	}

	if err := h.usecase.ConfirmReady(c.Context(), teamID, matchID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Ready confirmation received",
		"status":  "READY",
	})
}

// CancelMatchmaking handles POST /api/matchmaking/cancel
func (h *MatchmakingHandler) CancelMatchmaking(c *fiber.Ctx) error {
	teamIDStr := c.Query("team_id")
	if teamIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "team_id is required",
		})
	}

	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid team ID",
		})
	}

	if err := h.usecase.CancelMatchmaking(c.Context(), teamID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Matchmaking cancelled",
		"status":  "CANCELLED",
	})
}

// WebSocketHandler handles WebSocket connections
func (h *MatchmakingHandler) WebSocketHandler(c *websocket.Conn) {
	teamID := c.Query("team_id")
	if teamID == "" {
		log.Println("WebSocket connection without team_id")
		c.Close()
		return
	}

	// Register connection
	h.hub.Register(teamID, c)
	defer h.hub.Unregister(teamID)

	// Keep connection alive
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error for team %s: %v", teamID, err)
			break
		}
	}
}

// HealthCheck handles GET /health
func (h *MatchmakingHandler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "healthy",
		"service": "golobby-matchmaking",
	})
}
