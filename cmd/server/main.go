package main

import (
	"log"

	"github.com/GilangAndhika/bepapais/internal/api"
	"github.com/GilangAndhika/bepapais/internal/auth"
	"github.com/GilangAndhika/bepapais/internal/config"
	"github.com/GilangAndhika/bepapais/internal/database"
	"github.com/GilangAndhika/bepapais/internal/location"

	"github.com/gofiber/fiber/v2"
	// Kita akan impor 'models' dan 'database' saat dibutuhkan
)

func main() {
	// 1. Muat Konfigurasi
	cfg := config.NewConfig()

	// 2. Hubungkan ke Database
	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Gagal terhubung ke MongoDB: %v", err)
	}

	// 3. Inisialisasi Layer (Dependency Injection)

	// --- Auth ---
	authCollection := db.Collection("admins")
	authRepo := auth.NewAuthRepository(authCollection)
	authService := auth.NewAuthService(authRepo, cfg.JWTSecret)
	authHandler := auth.NewAuthHandler(authService)

	// --- Location ---
	locationCollection := db.Collection("locations") // <-- BARU
	locationRepo := location.NewLocationRepository(locationCollection) // <-- BARU
	locationService := location.NewLocationService(locationRepo)       // <-- BARU
	locationHandler := location.NewLocationHandler(locationService)   // <-- BARU

	// 4. Inisialisasi Server Fiber
	app := fiber.New()

	// 5. Siapkan Rute
	// Teruskan semua handler DAN config ke router
	api.SetupRoutes(app, cfg, authHandler, locationHandler) // <-- PARAMETER DIPERBARUI

	// 6. Jalankan Server
	serverAddr := ":" + cfg.ServerPort
	log.Printf("Server berjalan di http://localhost%s", serverAddr)
	log.Fatal(app.Listen(serverAddr))
}