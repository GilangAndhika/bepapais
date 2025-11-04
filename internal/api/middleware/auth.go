package middleware

import (
	"github.com/GilangAndhika/bepapais/internal/config"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Protected adalah middleware untuk memvalidasi JWT
func Protected(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil header Authorization
		authHeader := c.Get("Authorization")

		// Cek format "Bearer [token]"
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or invalid authorization token",
			})
		}

		// Ambil token-nya saja
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse dan validasi token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Pastikan metode signing adalah HS256
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Token valid, simpan claims ke 'locals' agar bisa dipakai handler selanjutnya
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to parse token claims",
			})
		}

		c.Locals("adminClaims", claims)
		
		// Lanjutkan ke handler selanjutnya
		return c.Next()
	}
}