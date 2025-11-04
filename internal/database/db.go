package database

import (
	"context"
	"log"
	"time"

	"github.com/GilangAndhika/bepapais/internal/config" 
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectDB menginisialisasi koneksi ke MongoDB dan mengembalikan database
func ConnectDB(cfg *config.Config) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		return nil, err
	}

	// Ping ke database untuk memastikan koneksi
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	log.Println("Berhasil terhubung ke MongoDB!")
	
	db := client.Database(cfg.DBName)
	return db, nil
}