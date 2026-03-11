package http

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golobby/matchmaking/internal/metrics"
)

// OCRResultRequest is the JSON payload sent by the Python OCR service webhook
type OCRResultRequest struct {
	MatchID    string  `json:"match_id"`
	Result     string  `json:"result"`      // "verified" | "disputed" | "error"
	WinnerTeam *string `json:"winner_team"`
	ScoreTeam1 *string `json:"score_team1"`
	ScoreTeam2 *string `json:"score_team2"`
	Confidence float64 `json:"confidence"`
	RawText    string  `json:"raw_text"`
}

// OCRHandler handles incoming OCR verification results from the Python microservice
type OCRHandler struct {
	wsHub         *WebSocketHub
	webhookSecret string
}

// NewOCRHandler creates a new OCR handler
func NewOCRHandler(wsHub *WebSocketHub) *OCRHandler {
	secret := os.Getenv("OCR_WEBHOOK_SECRET")
	if secret == "" {
		secret = "ocr-webhook-secret-changeme"
		log.Println("⚠️  OCR_WEBHOOK_SECRET not set — using insecure default")
	}
	return &OCRHandler{
		wsHub:         wsHub,
		webhookSecret: secret,
	}
}

// ReceiveOCRResult handles POST /api/ocr/result
// This endpoint is called asynchronously by the Python OCR microservice
// after it finishes analysing a screenshot.
func (h *OCRHandler) ReceiveOCRResult(c *fiber.Ctx) error {
	// ── Authenticate the webhook ──────────────────────────────────────
	secret := c.Get("X-OCR-Secret")
	if secret != h.webhookSecret {
		log.Printf("🔐 OCR webhook rejected — invalid secret from %s", c.IP())
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// ── Parse payload ─────────────────────────────────────────────────
	var req OCRResultRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	log.Printf("🤖 OCR result received: match_id=%s result=%s confidence=%.2f%%",
		req.MatchID, req.Result, req.Confidence)

	// ── Track metrics ─────────────────────────────────────────────────
	metrics.OCRVerificationTotal.WithLabelValues(req.Result).Inc()

	// ── Broadcast result to any connected WebSocket clients ───────────
	// The frontend listens for this event to show the verification badge.
	notification := map[string]interface{}{
		"type":       "OCR_RESULT",
		"match_id":   req.MatchID,
		"result":     req.Result,
		"confidence": req.Confidence,
		"timestamp":  time.Now().Format(time.RFC3339),
	}
	if req.WinnerTeam != nil {
		notification["winner_team"] = *req.WinnerTeam
	}
	if req.ScoreTeam1 != nil && req.ScoreTeam2 != nil {
		notification["score"] = *req.ScoreTeam1 + " - " + *req.ScoreTeam2
	}

	// Broadcast to match-specific channel (match_id used as client key)
	h.wsHub.BroadcastToClient(req.MatchID, notification)

	// ── Apply reputation effects for fraud ────────────────────────────
	if req.Result == "disputed" {
		log.Printf("⚠️  Fraud detected on match %s — flagging for review", req.MatchID)
		metrics.ReputationPenaltiesTotal.WithLabelValues("fraud").Inc()
		// TODO: call reputation usecase to apply suspension
		// This is wired up when the full reputation usecase is implemented
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "OCR result processed",
	})
}

// RegisterOCRRoutes registers OCR-related routes
func RegisterOCRRoutes(app *fiber.App, handler *OCRHandler) {
	app.Post("/api/ocr/result", handler.ReceiveOCRResult)
}
