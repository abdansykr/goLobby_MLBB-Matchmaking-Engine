package http

import (
	"database/sql"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golobby/matchmaking/internal/domain"
	"github.com/golobby/matchmaking/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler handles register and login requests.
type AuthHandler struct {
	userRepo  *repository.UserRepository
	jwtSecret string
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(userRepo *repository.UserRepository, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// ─── Request / Response structs ──────────────────────────────────────────────

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Whatsapp string `json:"whatsapp"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ─── Register ────────────────────────────────────────────────────────────────

// Register handles POST /api/auth/register
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req registerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// ── Validation ──────────────────────────────────────────────────────
	req.Username = strings.TrimSpace(req.Username)
	req.Email    = strings.TrimSpace(strings.ToLower(req.Email))
	req.Password = strings.TrimSpace(req.Password)
	req.Whatsapp = strings.TrimSpace(req.Whatsapp)

	if req.Username == "" || req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "username, email, dan password wajib diisi",
		})
	}
	if len(req.Username) < 3 || len(req.Username) > 30 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username harus antara 3-30 karakter",
		})
	}
	if len(req.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password minimal 6 karakter",
		})
	}

	// ── Check duplicate email ────────────────────────────────────────────
	if _, err := h.userRepo.FindByEmail(c.Context(), req.Email); err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Email sudah terdaftar",
		})
	}

	// ── Check duplicate username ─────────────────────────────────────────
	if _, err := h.userRepo.FindByUsername(c.Context(), req.Username); err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Username sudah digunakan",
		})
	}

	// ── Hash password ────────────────────────────────────────────────────
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal memproses password",
		})
	}

	// ── Persist ──────────────────────────────────────────────────────────
	user := &domain.User{
		Username:       req.Username,
		Email:          req.Email,
		PasswordHash:   string(hash),
		WhatsappNumber: req.Whatsapp,
	}
	if err := h.userRepo.Create(c.Context(), user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal membuat akun",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":  "Akun berhasil dibuat",
		"username": user.Username,
	})
}

// ─── Reset Password ───────────────────────────────────────────────────────────

// ResetPassword handles POST /api/auth/reset-password
// Verifies email + username match, then sets a new password.
func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	var req struct {
		Email       string `json:"email"`
		Username    string `json:"username"`
		NewPassword string `json:"new_password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	req.Email       = strings.TrimSpace(strings.ToLower(req.Email))
	req.Username    = strings.TrimSpace(req.Username)
	req.NewPassword = strings.TrimSpace(req.NewPassword)

	if req.Email == "" || req.Username == "" || req.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email, username, dan password baru wajib diisi"})
	}
	if len(req.NewPassword) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Password baru minimal 6 karakter"})
	}

	// Find user by email
	user, err := h.userRepo.FindByEmail(c.Context(), req.Email)
	if err == sql.ErrNoRows || user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Email atau username tidak cocok"})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Terjadi kesalahan server"})
	}

	// Verify username matches
	if !strings.EqualFold(user.Username, req.Username) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Email atau username tidak cocok"})
	}

	// Hash and save new password
	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memproses password"})
	}
	user.PasswordHash = string(newHash)

	if err := h.userRepo.UpdatePasswordHash(c.Context(), user.ID, string(newHash)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan password baru"})
	}

	return c.JSON(fiber.Map{"message": "Password berhasil direset. Silakan login."})
}

// ─── Login ───────────────────────────────────────────────────────────────────

// Login handles POST /api/auth/login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	req.Email    = strings.TrimSpace(strings.ToLower(req.Email))
	req.Password = strings.TrimSpace(req.Password)

	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email dan password wajib diisi",
		})
	}

	// ── Find user ────────────────────────────────────────────────────────
	user, err := h.userRepo.FindByEmail(c.Context(), req.Email)
	if err == sql.ErrNoRows || user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Email atau password salah",
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Terjadi kesalahan server",
		})
	}

	// ── Verify password ──────────────────────────────────────────────────
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Email atau password salah",
		})
	}

	// ── Sign JWT ─────────────────────────────────────────────────────────
	claims := jwt.MapClaims{
		"user_id":  user.ID.String(),
		"username": user.Username,
		"email":    user.Email,
		"exp":      time.Now().Add(72 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal membuat token",
		})
	}

	return c.JSON(fiber.Map{
		"token":           signed,
		"username":        user.Username,
		"email":           user.Email,
		"user_id":         user.ID,
		"whatsapp_number": user.WhatsappNumber,
		"avatar_url":      user.AvatarURL,
	})
}

// RegisterAuthRoutes registers auth endpoints.
func RegisterAuthRoutes(app *fiber.App, handler *AuthHandler) {
	auth := app.Group("/api/auth")
	auth.Post("/register", handler.Register)
	auth.Post("/login", handler.Login)
	auth.Post("/reset-password", handler.ResetPassword)
}
