package http

import (
	"log"
	"net"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golobby/matchmaking/internal/domain"
	"github.com/golobby/matchmaking/internal/usecase"
	"github.com/google/uuid"
)

type ScrimHandler struct {
	scrimUsecase *usecase.ScrimMatchmakingUsecase
	wsHub        *WebSocketHub
}

func NewScrimHandler(scrimUsecase *usecase.ScrimMatchmakingUsecase, wsHub *WebSocketHub) *ScrimHandler {
	h := &ScrimHandler{
		scrimUsecase: scrimUsecase,
		wsHub:        wsHub,
	}

	// Wire the WebSocket notifier into the usecase so it can push
	// MATCH_DECLINED events to the opponent in real-time.
	scrimUsecase.SetWSNotifier(func(targetRequestID string, payload map[string]interface{}) {
		wsHub.BroadcastToClient(targetRequestID, payload)
	})

	return h
}

// Helper functions to extract values from map (support multiple keys)
func getStringValue(data map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		if val, ok := data[key]; ok {
			if str, ok := val.(string); ok {
				return str
			}
		}
	}
	return ""
}

func getIntValue(data map[string]interface{}, keys ...string) int {
	for _, key := range keys {
		if val, ok := data[key]; ok {
			// Handle both number and string
			switch v := val.(type) {
			case float64:
				return int(v)
			case int:
				return v
			case string:
				// Try parse string to int
				if i, err := strconv.Atoi(v); err == nil {
					return i
				}
			}
		}
	}
	return 0
}

// CreateRequest handles POST /api/scrim/request
func (h *ScrimHandler) CreateRequest(c *fiber.Ctx) error {
	// Manual parsing to support both camelCase and snake_case
	var rawData map[string]interface{}

	// Log raw body for debugging
	bodyBytes := c.Body()
	log.Printf("📥 Received request body: %s", string(bodyBytes))

	if err := c.BodyParser(&rawData); err != nil {
		log.Printf("❌ JSON Binding Error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":    "Invalid JSON format",
			"details":  err.Error(),
			"received": string(bodyBytes),
		})
	}

	// Extract values (support both camelCase and snake_case)
	teamName := getStringValue(rawData, "teamName", "team_name")
	whatsappNumber := getStringValue(rawData, "whatsappNumber", "whatsapp_number")
	avatarURL := getStringValue(rawData, "avatarUrl", "avatar_url")
	category := getStringValue(rawData, "category")
	rankWeight := getIntValue(rawData, "rankWeight", "rank_weight")

	log.Printf("✅ Parsed Data: TeamName=%s, WhatsApp=%s, Category=%s, Rank=%d",
		teamName, whatsappNumber, category, rankWeight)

	// Validate input
	if teamName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "teamName is required",
		})
	}
	if whatsappNumber == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "whatsappNumber is required",
		})
	}
	if category != "POKE" && category != "WARKOP" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "category must be POKE or WARKOP",
		})
	}
	if rankWeight < 1 || rankWeight > 10 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "rankWeight must be between 1 and 10",
		})
	}

	// Validate category-specific rank
	if category == "POKE" && (rankWeight < 1 || rankWeight > 8) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "POKE category requires rankWeight between 1 and 8",
		})
	}
	if category == "WARKOP" && (rankWeight < 9 || rankWeight > 10) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "WARKOP category requires rankWeight between 9 and 10",
		})
	}

	// Get client IP - Prioritize X-Forwarded-For for testing
	ipAddress := c.Get("X-Forwarded-For")
	if ipAddress == "" {
		ipAddress = c.Get("X-Real-IP")
	}
	if ipAddress == "" {
		ipAddress = c.IP()
	}
	if ipAddress == "" {
		// Final fallback to connection
		if addr, ok := c.Context().RemoteAddr().(*net.TCPAddr); ok {
			ipAddress = addr.IP.String()
		}
	}

	log.Printf("🌐 Client IP: %s", ipAddress)

	// Create request
	request := &domain.ScrimRequest{
		TeamName:       teamName,
		WhatsAppNumber: whatsappNumber,
		AvatarURL:      avatarURL,
		Category:       domain.ScrimCategory(category),
		RankWeight:     rankWeight,
		IPAddress:      ipAddress,
	}

	createdRequest, err := h.scrimUsecase.CreateRequest(c.Context(), request)
	if err != nil {
		log.Printf("❌ Error creating scrim request: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	log.Printf("✅ Scrim request created: ID=%s, Team=%s, Category=%s, Rank=%d",
		createdRequest.ID, createdRequest.TeamName, createdRequest.Category, createdRequest.RankWeight)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":     "Scrim request created successfully",
		"request_id":  createdRequest.ID,
		"status":      createdRequest.Status,
		"category":    createdRequest.Category,
		"rank_weight": createdRequest.RankWeight,
		"expires_at":  createdRequest.ExpiresAt,
	})
}

// GetRequestStatus handles GET /api/scrim/request/:id
func (h *ScrimHandler) GetRequestStatus(c *fiber.Ctx) error {
	idParam := c.Params("id")
	requestID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request ID",
		})
	}

	request, matchResponse, err := h.scrimUsecase.GetRequest(c.Context(), requestID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Request not found",
		})
	}

	response := fiber.Map{
		"request_id":  request.ID,
		"team_name":   request.TeamName,
		"category":    request.Category,
		"rank_weight": request.RankWeight,
		"status":      request.Status,
		"created_at":  request.CreatedAt,
		"expires_at":  request.ExpiresAt,
	}

	// Add match details if matched
	if matchResponse != nil {
		response["match"] = fiber.Map{
			"match_id":        matchResponse.Match.ID,
			"opponent_name":   matchResponse.OpponentName,
			"opponent_number": matchResponse.OpponentNumber,
			"whatsapp_url":    matchResponse.WhatsAppURL,
			"expires_in":      matchResponse.ExpiresIn,
			"status":          matchResponse.Match.Status,
		}

		// Notify via WebSocket if this is the first time checking after match
		// (This is a simple implementation - in production you'd track this better)
		go h.notifyMatchFound(request.ID.String(), matchResponse)
	}

	return c.JSON(response)
}

// CancelRequest handles POST /api/scrim/request/:id/cancel
func (h *ScrimHandler) CancelRequest(c *fiber.Ctx) error {
	idParam := c.Params("id")
	requestID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request ID",
		})
	}

	err = h.scrimUsecase.CancelRequest(c.Context(), requestID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Request cancelled successfully",
	})
}

// ConfirmMatch handles POST /api/scrim/match/:id/confirm?request_id=<uuid>
func (h *ScrimHandler) ConfirmMatch(c *fiber.Ctx) error {
	idParam := c.Params("id")
	matchID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid match ID",
		})
	}

	requestIDParam := c.Query("request_id")
	requestID, err := uuid.Parse(requestIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "request_id query param is required and must be a valid UUID",
		})
	}

	err = h.scrimUsecase.ConfirmMatch(c.Context(), matchID, requestID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Match acceptance recorded successfully",
	})
}

// DeclineMatch handles POST /api/scrim/match/:id/decline?request_id=<uuid>
// Called when a player clicks "Tolak Lawan" on the match modal.
// It cancels the match and immediately notifies the opponent via WebSocket.
func (h *ScrimHandler) DeclineMatch(c *fiber.Ctx) error {
	idParam := c.Params("id")
	matchID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid match ID",
		})
	}

	requestIDParam := c.Query("request_id")
	requestID, err := uuid.Parse(requestIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "request_id query param is required and must be a valid UUID",
		})
	}

	if err := h.scrimUsecase.DeclineMatch(c.Context(), matchID, requestID); err != nil {
		log.Printf("❌ DeclineMatch error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Match declined — opponent notified",
	})
}

// RejectMatch handles POST /api/scrim/match/:id/reject?request_id=<uuid>
// This is the canonical, concurrency-safe rejection endpoint.
// It uses a Redis Lua script to guarantee that only ONE goroutine performs
// database writes and WebSocket broadcasts even if both players reject
// simultaneously. MATCH_CANCELLED is sent to BOTH participants.
func (h *ScrimHandler) RejectMatch(c *fiber.Ctx) error {
	idParam := c.Params("id")
	matchID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid match ID",
		})
	}

	requestIDParam := c.Query("request_id")
	requestID, err := uuid.Parse(requestIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "request_id query param is required and must be a valid UUID",
		})
	}

	if err := h.scrimUsecase.RejectMatch(c.Context(), matchID, requestID); err != nil {
		log.Printf("❌ RejectMatch error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Match rejected — kedua tim telah diberitahu",
	})
}

// notifyMatchFound sends WebSocket notification for match found
func (h *ScrimHandler) notifyMatchFound(requestID string, matchResponse *domain.MatchResponse) {
	message := map[string]interface{}{
		"type":                "SCRIM_MATCH_FOUND",
		"match_id":            matchResponse.Match.ID,
		"opponent_name":       matchResponse.OpponentName,
		"opponent_number":     matchResponse.OpponentNumber,
		"opponent_avatar_url": matchResponse.OpponentAvatarURL,
		"whatsapp_url":        matchResponse.WhatsAppURL,
		"expires_in":          matchResponse.ExpiresIn,
	}

	h.wsHub.BroadcastToClient(requestID, message)
}

// RegisterScrimRoutes registers all scrim-related routes
func RegisterScrimRoutes(app *fiber.App, handler *ScrimHandler) {
	api := app.Group("/api/scrim")

	api.Post("/request", handler.CreateRequest)
	api.Get("/request/:id", handler.GetRequestStatus)
	api.Post("/request/:id/cancel", handler.CancelRequest)
	api.Post("/match/:id/confirm", handler.ConfirmMatch)
	api.Post("/match/:id/decline", handler.DeclineMatch)  // legacy: single notify
	api.Post("/match/:id/reject", handler.RejectMatch)    // canonical: atomic, broadcasts to both
}
