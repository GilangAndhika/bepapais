package models

import "time"

// Location merepresentasikan 'locations' collection
type Location struct {
	ID        string    `json:"id" bson:"_id,omitempty"` // omitempty agar bisa auto-generate
	Name      string    `json:"name" bson:"name"`
	Slug      string    `json:"slug" bson:"slug"`
	Type      string    `json:"type" bson:"type"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}