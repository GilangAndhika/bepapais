package location

import (
	"context"

	"github.com/GilangAndhika/bepapais/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Create(ctx context.Context, location *models.Location) error
	GetAll(ctx context.Context) ([]models.Location, error)
	GetByID(ctx context.Context, id string) (*models.Location, error)
	Update(ctx context.Context, id string, location *models.Location) error
	Delete(ctx context.Context, id string) error
	FindBySlug(ctx context.Context, slug string) (*models.Location, error)
}

type repository struct {
	collection *mongo.Collection
}

func NewLocationRepository(collection *mongo.Collection) Repository {
	return &repository{collection: collection}
}

func (r *repository) Create(ctx context.Context, location *models.Location) error {
	_, err := r.collection.InsertOne(ctx, location)
	return err
}

func (r *repository) GetAll(ctx context.Context) ([]models.Location, error) {
	var locations []models.Location
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &locations); err != nil {
		return nil, err
	}
	return locations, nil
}

func (r *repository) GetByID(ctx context.Context, id string) (*models.Location, error) {
	var location models.Location
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err // ID tidak valid
	}

	filter := bson.M{"_id": objID}
	err = r.collection.FindOne(ctx, filter).Decode(&location)
	if err != nil {
		return nil, err
	}
	return &location, nil
}

func (r *repository) Update(ctx context.Context, id string, location *models.Location) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"name": location.Name,
			"slug": location.Slug,
			"type": location.Type,
		},
	}
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

func (r *repository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

func (r *repository) FindBySlug(ctx context.Context, slug string) (*models.Location, error) {
	var location models.Location
	filter := bson.M{"slug": slug}
	err := r.collection.FindOne(ctx, filter).Decode(&location)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &location, nil
}