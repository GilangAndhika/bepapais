package location

import (
	"context"
	"errors"
	"time"

	"github.com/GilangAndhika/bepapais/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrLocationExists = errors.New("lokasi dengan slug tersebut sudah ada")
	ErrLocationNotFound = errors.New("lokasi tidak ditemukan")
)

type Service interface {
	CreateLocation(ctx context.Context, name, slug, locType string) error
	GetAllLocations(ctx context.Context) ([]models.Location, error)
	GetLocationByID(ctx context.Context, id string) (*models.Location, error)
	UpdateLocation(ctx context.Context, id, name, slug, locType string) error
	DeleteLocation(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewLocationService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateLocation(ctx context.Context, name, slug, locType string) error {
	// Cek jika slug sudah ada
	existing, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		return err
	}
	if existing != nil {
		return ErrLocationExists
	}

	location := &models.Location{
		// ID akan di-generate otomatis oleh Mongo
		Name:      name,
		Slug:      slug,
		Type:      locType,
		CreatedAt: time.Now(),
	}

	return s.repo.Create(ctx, location)
}

func (s *service) GetAllLocations(ctx context.Context) ([]models.Location, error) {
	return s.repo.GetAll(ctx)
}

func (s *service) GetLocationByID(ctx context.Context, id string) (*models.Location, error) {
	location, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrLocationNotFound
		}
		return nil, err
	}
	return location, nil
}

func (s *service) UpdateLocation(ctx context.Context, id, name, slug, locType string) error {
	// Cek apakah lokasi ada
	_, err := s.GetLocationByID(ctx, id)
	if err != nil {
		return err
	}

	location := &models.Location{
		Name: name,
		Slug: slug,
		Type: locType,
	}
	
	return s.repo.Update(ctx, id, location)
}

func (s *service) DeleteLocation(ctx context.Context, id string) error {
	// Cek apakah lokasi ada
	_, err := s.GetLocationByID(ctx, id)
	if err != nil {
		return err
	}
	
	// (Nanti di sini kita harus cek dulu apakah ada 'camera' yang
	// menggunakan location_id ini sebelum menghapus)

	return s.repo.Delete(ctx, id)
}