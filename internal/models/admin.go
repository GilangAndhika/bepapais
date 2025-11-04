package models

import "time"

// Admin merepresentasikan 'admins' collection
type Admin struct {
	ID           string    `json:"id" bson:"_id,omitempty"`
	Username     string    `json:"username" bson:"username"`
	PasswordHash string    `json:"-" bson:"password_hash"` // '-' di JSON agar tidak terkirim
	Role         string    `json:"role" bson:"role"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
}