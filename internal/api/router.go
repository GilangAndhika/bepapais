package api

import (
	"github.com/GilangAndhika/bepapais/internal/api/middleware"
	"github.com/GilangAndhika/bepapais/internal/auth"
	"github.com/GilangAndhika/bepapais/internal/camera"
	"github.com/GilangAndhika/bepapais/internal/config"
	"github.com/GilangAndhika/bepapais/internal/location"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupRoutes mengkonfigurasi semua rute aplikasi
func SetupRoutes(
	app *fiber.App,
	cfg *config.Config,
	authHandler *auth.AuthHandler,
	locationHandler *location.LocationHandler,
	cameraHandler *camera.CameraHandler,
) {

	// Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
	}))
	app.Use(logger.New())

	// Middleware autentikasi
	authWare := middleware.Protected(cfg) // <-- INISIALISASI MIDDLEWARE

	// Rute publik
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"message": "Selamat Datang di Papais CCTV API",
		})
	})
	app.Static("/live", "./media")

	// --- Grup API Publik (untuk frontend React) ---
	api := app.Group("/api")
	api.Get("/locations", locationHandler.GetAllLocations)     // GET /api/locations
	api.Get("/locations/:id", locationHandler.GetLocationByID) // GET /api/locations/123
	api.Get("/cameras", cameraHandler.GetAllCameras)           // GET /api/cameras
	api.Get("/cameras/:id", cameraHandler.GetCameraByID)       // GET /api/cameras/123

	// --- Grup API Admin (perlu token) ---
	admin := api.Group("/admin")
	admin.Post("/login", authHandler.Login) // POST /api/admin/login

	// --- Rute Admin yang Dilindungi ---
	// Rute di bawah ini memerlukan "Bearer [token]" di header

	// Rute Admin Lokasi
	admin.Post("/locations", authWare, locationHandler.CreateLocation)       // POST /api/admin/locations
	admin.Put("/locations/:id", authWare, locationHandler.UpdateLocation)    // PUT /api/admin/locations/123
	admin.Delete("/locations/:id", authWare, locationHandler.DeleteLocation) // DELETE /api/admin/locations/123

	// Rute Admin Kamera
	admin.Post("/cameras", authWare, cameraHandler.CreateCamera)       // POST /api/admin/cameras
	admin.Put("/cameras/:id", authWare, cameraHandler.UpdateCamera)    // PUT /api/admin/cameras/123
	admin.Delete("/cameras/:id", authWare, cameraHandler.DeleteCamera) // DELETE /api/admin/cameras/123
}
