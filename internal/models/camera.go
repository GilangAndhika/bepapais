package models

import "time"

// Camera merepresentasikan 'cameras' collection di MongoDB
type Camera struct {
	ID           string       `json:"id" bson:"_id"` // _id di Mongo, 'id' di JSON
	Name         string       `json:"name" bson:"name"`
	LocationText string       `json:"location_text" bson:"location_text"`
	LocationID   string       `json:"location_id" bson:"location_id"`
	Source       CameraSource `json:"source" bson:"source"`
	Features     Features     `json:"features" bson:"features"`
	Enabled      bool         `json:"enabled" bson:"enabled"`
	CreatedAt    time.Time    `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at" bson:"updated_at"`
}

type CameraSource struct {
	Type     string `json:"type" bson:"type"`
	IP       string `json:"ip" bson:"ip"`
	Port     int    `json:"port" bson:"port"`
	Path     string `json:"path" bson:"path"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

type Features struct {
	PTZ   bool   `json:"ptz" bson:"ptz"`
	Audio bool   `json:"audio" bson:"audio"`
	Brand string `json:"brand" bson:"brand"`
}