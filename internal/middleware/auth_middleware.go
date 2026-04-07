package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// NewAuthMiddleware returns a Fiber middleware that validates a JWT Bearer token.
// If valid, it stores the user_id string in c.Locals("user_id").
// If invalid or missing, it returns 401 Unauthorized.
//
// Usage in main.go:
//
//	authMW := middleware.NewAuthMiddleware(cfg.JWT.Secret)
//	api.Post("/scrim/request", authMW, scrimHandler.CreateRequest)
func NewAuthMiddleware(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// ── Extract Bearer token ─────────────────────────────────────────
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header is required",
			})
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization format — expected: Bearer <token>",
			})
		}
		tokenString := parts[1]

		// ── Parse & validate ────────────────────────────────────────────
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "unexpected signing method")
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token tidak valid atau sudah kedaluwarsa",
			})
		}

		// ── Store claims in context ──────────────────────────────────────
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token claims tidak valid",
			})
		}

		c.Locals("user_id", claims["user_id"])
		c.Locals("username", claims["username"])

		return c.Next()
	}
}
