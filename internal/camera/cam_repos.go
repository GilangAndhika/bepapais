package camera

import (
	"context"

	"github.com/GilangAndhika/bepapais/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Create(ctx context.Context, camera *models.Camera) error
	GetAll(ctx context.Context) ([]models.Camera, error)
	GetAllEnabled(ctx context.Context) ([]models.Camera, error)
	GetByID(ctx context.Context, id string) (*models.Camera, error)
	Update(ctx context.Context, id string, camera *models.Camera) error
	Delete(ctx context.Context, id string) error
	
	SearchEnabledByName(ctx context.Context, query string) ([]models.Camera, error)
	GetEnabledByLocation(ctx context.Context, locationID string) ([]models.Camera, error)
}

type repository struct {
	collection *mongo.Collection
}

func NewCameraRepository(collection *mongo.Collection) Repository {
	return &repository{collection: collection}
}

func (r *repository) Create(ctx context.Context, camera *models.Camera) error {
	_, err := r.collection.InsertOne(ctx, camera)
	return err
}

func (r *repository) GetAll(ctx context.Context) ([]models.Camera, error) {
	var cameras []models.Camera
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &cameras); err != nil {
		return nil, err
	}
	return cameras, nil
}

// GetAllEnabled hanya mengambil kamera yang 'enabled: true'
func (r *repository) GetAllEnabled(ctx context.Context) ([]models.Camera, error) {
	var cameras []models.Camera
	filter := bson.M{"enabled": true}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &cameras); err != nil {
		return nil, err
	}
	return cameras, nil
}


func (r *repository) GetByID(ctx context.Context, id string) (*models.Camera, error) {
	var camera models.Camera
	filter := bson.M{"_id": id} // ID kita adalah string kustom
	err := r.collection.FindOne(ctx, filter).Decode(&camera)
	if err != nil {
		return nil, err
	}
	return &camera, nil
}

func (r *repository) Update(ctx context.Context, id string, camera *models.Camera) error {
	update := bson.M{"$set": camera}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *repository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// SearchEnabledByName mencari kamera berdasarkan nama (case-insensitive)
func (r *repository) SearchEnabledByName(ctx context.Context, query string) ([]models.Camera, error) {
	var cameras []models.Camera
	
	// Filter: enabled=true DAN name=query (regex, case-insensitive)
	filter := bson.M{
		"enabled": true,
		"name":    bson.M{"$regex": query, "$options": "i"},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &cameras); err != nil {
		return nil, err
	}
	return cameras, nil
}

// GetEnabledByLocation mencari kamera berdasarkan location_id
func (r *repository) GetEnabledByLocation(ctx context.Context, locationID string) ([]models.Camera, error) {
	var cameras []models.Camera
	
	// Filter: enabled=true DAN location_id=locationID
	filter := bson.M{
		"enabled":     true,
		"location_id": locationID,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &cameras); err != nil {
		return nil, err
	}
	return cameras, nil
}