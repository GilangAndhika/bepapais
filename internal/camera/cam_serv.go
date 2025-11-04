package camera

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/GilangAndhika/bepapais/internal/models"
	"github.com/GilangAndhika/bepapais/internal/streaming"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCameraNotFound = errors.New("kamera tidak ditemukan")
	ErrCameraExists   = errors.New("kamera dengan ID tersebut sudah ada")
)

type Service interface {
	CreateCamera(ctx context.Context, cam *models.Camera) error
	GetAllCameras(ctx context.Context) ([]models.Camera, error)
	GetCameraByID(ctx context.Context, id string) (*models.Camera, error)
	UpdateCamera(ctx context.Context, id string, cam *models.Camera) error
	DeleteCamera(ctx context.Context, id string) error
	
	SearchCameras(ctx context.Context, query string) ([]models.Camera, error)
	GetCamerasByLocation(ctx context.Context, locationID string) ([]models.Camera, error)
}

type service struct {
	repo         Repository
	streamManager *streaming.Manager 
}

func NewCameraService(repo Repository, sm *streaming.Manager) Service {
	return &service{
		repo:         repo,
		streamManager: sm,
	}
}

func (s *service) CreateCamera(ctx context.Context, cam *models.Camera) error {
	// Cek jika ID sudah ada
	existing, err := s.repo.GetByID(ctx, cam.ID)
	if existing != nil {
		return ErrCameraExists
	}
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	cam.CreatedAt = time.Now()
	cam.UpdatedAt = time.Now()

	// (SECURITY: Di sinilah Anda seharusnya mengenkripsi cam.Source.Password)
	
	if err := s.repo.Create(ctx, cam); err != nil {
		return err
	}

	// Jika 'enabled', otomatis mulai stream
	if cam.Enabled {
		s.streamManager.StartStream(*cam)
	}
	return nil
}

func (s *service) GetAllCameras(ctx context.Context) ([]models.Camera, error) {
	return s.repo.GetAll(ctx)
}

func (s *service) GetCameraByID(ctx context.Context, id string) (*models.Camera, error) {
	cam, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrCameraNotFound
		}
		return nil, err
	}
	// (SECURITY: Di sinilah Anda seharusnya mendekripsi password jika diperlukan)
	return cam, nil
}

func (s *service) UpdateCamera(ctx context.Context, id string, cam *models.Camera) error {
	// Dapatkan data lama untuk perbandingan
	oldCam, err := s.GetCameraByID(ctx, id)
	if err != nil {
		return err
	}

	// Update data
	cam.UpdatedAt = time.Now()
	cam.CreatedAt = oldCam.CreatedAt // Pastikan created_at tidak berubah
	cam.ID = oldCam.ID // Pastikan ID tidak berubah

	if err := s.repo.Update(ctx, id, cam); err != nil {
		return err
	}

	// --- Logika Orkestrasi Streaming ---
	
	// Kasus 1: Dimatikan (was enabled, now disabled)
	if oldCam.Enabled && !cam.Enabled {
		s.streamManager.StopStream(cam.ID)
	}

	// Kasus 2: Dinyalakan (was disabled, now enabled)
	if !oldCam.Enabled && cam.Enabled {
		s.streamManager.StartStream(*cam)
	}

	// Kasus 3: Diedit saat sedang nyala (source changed)
	if oldCam.Enabled && cam.Enabled {
		// (Kita bisa cek source, atau gampangnya, restart saja)
		s.streamManager.StartStream(*cam) // StartStream sudah handle stop stream lama
	}
	
	return nil
}

func (s *service) DeleteCamera(ctx context.Context, id string) error {
	// Cek jika ada
	_, err := s.GetCameraByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Hentikan stream jika sedang berjalan
	s.streamManager.StopStream(id)
	
	// (Bonus: Hapus folder ./media/[id]
	os.RemoveAll(filepath.Join("media", id))
	
	return nil
}

func (s *service) SearchCameras(ctx context.Context, query string) ([]models.Camera, error) {
	return s.repo.SearchEnabledByName(ctx, query)
}

func (s *service) GetCamerasByLocation(ctx context.Context, locationID string) ([]models.Camera, error) {
	return s.repo.GetEnabledByLocation(ctx, locationID)
}