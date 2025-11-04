package auth

import (
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service Service
}

// NewAuthHandler membuat instance handler baru
func NewAuthHandler(service Service) *AuthHandler {
	return &AuthHandler{service: service}
}

// LoginRequest adalah struct untuk parsing body JSON
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login adalah handler untuk endpoint /login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest

	// 1. Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// 2. Validasi input
	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username dan password tidak boleh kosong",
		})
	}

	// 3. Panggil service
	token, err := h.service.Login(c.Context(), req.Username, req.Password)
	if err != nil {
		// Cek jenis error
		switch err {
		case ErrUserNotFound, ErrInvalidCredentials:
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Username atau password salah",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}
	}

	// 4. Kirim token sebagai response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login berhasil",
		"token":   token,
	})
}