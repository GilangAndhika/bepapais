package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


type Admin struct {
	// MODIFIKASI: Ubah 'string' menjadi 'primitive.ObjectID'
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username     string             `json:"username" bson:"username"`
	PasswordHash string             `json:"-" bson:"password_hash"` // Ini sudah benar
	Role         string             `json:"role" bson:"role"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
}
