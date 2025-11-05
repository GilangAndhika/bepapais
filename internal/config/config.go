package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config menampung semua konfigurasi aplikasi
type Config struct {
	MongoURI   string
	DBName     string
	ServerPort string
	JWTSecret  string
}

// NewConfig memuat konfigurasi dari file .env
func NewConfig() *Config {
	// Memuat .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Peringatan: Tidak dapat menemukan file .env, menggunakan environment variables sistem.")
	}

	// Ambil "PORT" (disediakan oleh Heroku)
	port := os.Getenv("PORT")
	
	// Jika "PORT" tidak ada, ambil "SERVER_PORT" (dari file .env lokal Anda)
	if port == "" {
		port = getEnv("SERVER_PORT", "8000") // Gunakan helper Anda
	}

	return &Config{
		MongoURI:   getEnv("MONGO_URI", "mongodb://localhost:27017"),
		DBName:     getEnv("DB_NAME", "papais_cctv"),
		ServerPort: port, // <-- Gunakan variabel 'port' yang sudah pintar di sini
		JWTSecret:  getEnv("JWT_SECRET", "super-secret-key"),
	}
}

// getEnv adalah helper untuk membaca env var dengan nilai default
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}