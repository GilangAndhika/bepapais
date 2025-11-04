package auth

import (
	"context"

	"github.com/GilangAndhika/bepapais/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	FindByUsername(ctx context.Context, username string) (*models.Admin, error)
}

type repository struct {
	collection *mongo.Collection
}

// NewAuthRepository membuat instance repository baru
func NewAuthRepository(collection *mongo.Collection) Repository {
	return &repository{collection: collection}
}

// FindByUsername mencari admin berdasarkan username
func (r *repository) FindByUsername(ctx context.Context, username string) (*models.Admin, error) {
	var admin models.Admin
	filter := bson.M{"username": username}

	err := r.collection.FindOne(ctx, filter).Decode(&admin)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // User tidak ditemukan, tapi bukan error
		}
		return nil, err
	}
	return &admin, nil
}