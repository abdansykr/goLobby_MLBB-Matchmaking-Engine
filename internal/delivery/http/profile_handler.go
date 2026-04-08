package http

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/golobby/matchmaking/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type ProfileHandler struct {
	userRepo *repository.UserRepository
}

func NewProfileHandler(userRepo *repository.UserRepository) *ProfileHandler {
	return &ProfileHandler{
		userRepo: userRepo,
	}
}

// GetProfile handles GET /api/user/profile
func (h *ProfileHandler) GetProfile(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	user, err := h.userRepo.FindByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(fiber.Map{
		"username":        user.Username,
		"email":           user.Email,
		"whatsapp_number": user.WhatsappNumber,
		"avatar_url":      user.AvatarURL,
	})
}

// UpdateProfile handles PUT /api/user/profile
func (h *ProfileHandler) UpdateProfile(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var req struct {
		Username       string `json:"username"`
		Email          string `json:"email"`
		WhatsappNumber string `json:"whatsapp_number"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	user, err := h.userRepo.FindByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Update fields if provided
	if req.Username != "" {
		user.Username = strings.TrimSpace(req.Username)
	}
	if req.Email != "" {
		user.Email = strings.TrimSpace(strings.ToLower(req.Email))
	}
	if req.WhatsappNumber != "" {
		user.WhatsappNumber = strings.TrimSpace(req.WhatsappNumber)
	}

	if err := h.userRepo.UpdateProfile(c.Context(), user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update profile"})
	}

	return c.JSON(fiber.Map{
		"message": "Profile updated successfully",
		"user": fiber.Map{
			"username":        user.Username,
			"email":           user.Email,
			"whatsapp_number": user.WhatsappNumber,
			"avatar_url":      user.AvatarURL,
		},
	})
}

// UploadAvatar handles POST /api/user/avatar
func (h *ProfileHandler) UploadAvatar(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to get avatar file"})
	}

	// Simple validation
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" && ext != ".webp" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Only images (png, jpg, webp) are allowed"})
	}
	if file.Size > 2*1024*1024 { // 2MB
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File size exceeds limit of 2MB"})
	}

	fileName := fmt.Sprintf("%s-%d%s", userIDStr, time.Now().Unix(), ext)
	savePath := fmt.Sprintf("./uploads/avatars/%s", fileName)

	if err := c.SaveFile(file, savePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save file"})
	}

	avatarURL := fmt.Sprintf("/uploads/avatars/%s", fileName)

	user, err := h.userRepo.FindByID(c.Context(), userID)
	if err == nil {
		user.AvatarURL = avatarURL
		_ = h.userRepo.UpdateProfile(c.Context(), user)
	}

	return c.JSON(fiber.Map{
		"message":    "Avatar uploaded successfully",
		"avatar_url": avatarURL,
	})
}

// ChangePassword handles PUT /api/user/password
func (h *ProfileHandler) ChangePassword(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	req.CurrentPassword = strings.TrimSpace(req.CurrentPassword)
	req.NewPassword = strings.TrimSpace(req.NewPassword)

	if req.CurrentPassword == "" || req.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Password lama dan baru wajib diisi"})
	}
	if len(req.NewPassword) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Password baru minimal 6 karakter"})
	}

	user, err := h.userRepo.FindByID(c.Context(), userID)
	if err == sql.ErrNoRows || user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Terjadi kesalahan server"})
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Password lama tidak sesuai"})
	}

	// Hash new password
	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memproses password"})
	}

	user.PasswordHash = string(newHash)
	if err := h.userRepo.UpdateProfile(c.Context(), user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan password"})
	}

	return c.JSON(fiber.Map{"message": "Password berhasil diperbarui"})
}

func RegisterProfileRoutes(app *fiber.App, handler *ProfileHandler, authMiddleware fiber.Handler) {
	api := app.Group("/api/user", authMiddleware)
	api.Get("/profile", handler.GetProfile)
	api.Put("/profile", handler.UpdateProfile)
	api.Put("/password", handler.ChangePassword)
	api.Post("/avatar", handler.UploadAvatar)
}
