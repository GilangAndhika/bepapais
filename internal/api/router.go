package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// SetupRoutes mengkonfigurasi semua rute aplikasi
// Nantinya kita akan tambahkan parameter (handler) di sini
func SetupRoutes(app *fiber.App) {

	// Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Izinkan semua
		AllowMethods: "GET,POST,PUT,DELETE",
	}))
	app.Use(logger.New())

	// Rute publik
	// Rute ini untuk mengecek apakah server hidup
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"message": "Selamat Datang di Papais CCTV API",
		})
	})
	
	// Sajikan folder 'media' (HLS) sebagai '/live'
	// ./media/cctv_01/index.m3u8 -> /live/cctv_01/index.m3u8
	app.Static("/live", "./media")

	// Grup API
	// api := app.Group("/api")

	// Rute Admin (akan kita tambahkan nanti)
	// admin := api.Group("/admin")
	// admin.Post("/login", ...)
}