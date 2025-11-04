package main

import (
	"log"

	"github.com/GilangAndhika/bepapais/internal/api"
	"github.com/GilangAndhika/bepapais/internal/config"
	"github.com/GilangAndhika/bepapais/internal/database"

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
	
	// (Nanti kita akan teruskan 'db' ini ke Repositories)
	_ = db // Menghindari error "unused variable" untuk saat ini

	// 3. Inisialisasi Server Fiber
	app := fiber.New()

	// 4. Siapkan Rute
	// Kita teruskan 'app' ke router kita
	api.SetupRoutes(app)
	
	// (Nanti di sini kita akan menginisialisasi dan menjalankan
	// 'Streaming Manager' untuk FFmpeg)


	// 5. Jalankan Server
	serverAddr := ":" + cfg.ServerPort
	log.Printf("Server berjalan di http://localhost%s", serverAddr)
	log.Fatal(app.Listen(serverAddr))
}