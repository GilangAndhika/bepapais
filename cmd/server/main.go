package main

import (
	"context"
	"log"

	"github.com/GilangAndhika/bepapais/internal/api"
	"github.com/GilangAndhika/bepapais/internal/auth"
	"github.com/GilangAndhika/bepapais/internal/config"
	"github.com/GilangAndhika/bepapais/internal/database"
	"github.com/GilangAndhika/bepapais/internal/location"
	"github.com/GilangAndhika/bepapais/internal/camera"
	"github.com/GilangAndhika/bepapais/internal/streaming"

	"github.com/gofiber/fiber/v2"
	// Kita akan impor 'models' dan 'database' saat dibutuhkan
)

// Fungsi helper baru untuk memulai stream saat boot
func startInitialStreams(repo camera.Repository, sm *streaming.Manager) {
	log.Println("[Main] Memulai semua stream yang aktif dari database...")
	
	enabledCameras, err := repo.GetAllEnabled(context.Background())
	if err != nil {
		log.Fatalf("[Main] Gagal mengambil kamera dari DB: %v", err)
	}

	log.Printf("[Main] Ditemukan %d stream untuk dimulai.", len(enabledCameras))
	for _, cam := range enabledCameras {
		// Jalankan dalam goroutine agar tidak memblokir satu sama lain
		go sm.StartStream(cam)
	}
}


func main() {
	// ... (1. Muat Konfigurasi)
	cfg := config.NewConfig()

	// ... (2. Hubungkan ke Database)
	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Gagal terhubung ke MongoDB: %v", err)
	}

	// 3. Inisialisasi Layer (Dependency Injection)

	// --- Streaming ---
	streamManager := streaming.NewManager()

	// --- Auth ---
	authCollection := db.Collection("admins")
	authRepo := auth.NewAuthRepository(authCollection)
	authService := auth.NewAuthService(authRepo, cfg.JWTSecret)
	authHandler := auth.NewAuthHandler(authService)

	// --- Location ---
	locationCollection := db.Collection("locations")
	locationRepo := location.NewLocationRepository(locationCollection)
	locationService := location.NewLocationService(locationRepo)
	locationHandler := location.NewLocationHandler(locationService)
	
	// --- Camera ---
	cameraCollection := db.Collection("cameras") 
	cameraRepo := camera.NewCameraRepository(cameraCollection) 
	cameraService := camera.NewCameraService(cameraRepo, streamManager) 
	cameraHandler := camera.NewCameraHandler(cameraService)

	// 4. Inisialisasi Server Fiber
	app := fiber.New()

	// 5. Siapkan Rute
	api.SetupRoutes(app, cfg, authHandler, locationHandler, cameraHandler)

	// 6. Mulai Stream yang Aktif
	// Ganti blok ini dengan fungsi helper kita
	go startInitialStreams(cameraRepo, streamManager) // <-- INI DIA PERUBAHANNYA

	// 7. Jalankan Server
	serverAddr := ":" + cfg.ServerPort
	log.Printf("Server berjalan di http://localhost%s", serverAddr)
	log.Fatal(app.Listen(serverAddr))
}