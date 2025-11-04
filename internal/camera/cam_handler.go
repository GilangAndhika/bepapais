package camera

import (
	"github.com/GilangAndhika/bepapais/internal/models"
	"github.com/gofiber/fiber/v2"
)

type CameraHandler struct {
	service Service
}

func NewCameraHandler(service Service) *CameraHandler {
	return &CameraHandler{service: service}
}

// CreateCamera (Admin)
func (h *CameraHandler) CreateCamera(c *fiber.Ctx) error {
	var cam models.Camera
	if err := c.BodyParser(&cam); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Validasi (sederhana)
	if cam.ID == "" || cam.Name == "" || cam.LocationID == "" || cam.Source.Path == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID, Name, LocationID, dan Source Path wajib diisi"})
	}

	if err := h.service.CreateCamera(c.Context(), &cam); err != nil {
		if err == ErrCameraExists {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Kamera dengan ID tersebut sudah ada"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(cam)
}

// GetAllCameras (Publik)
func (h *CameraHandler) GetAllCameras(c *fiber.Ctx) error {
	cams, err := h.service.GetAllCameras(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(cams)
}

// GetCameraByID (Publik)
func (h *CameraHandler) GetCameraByID(c *fiber.Ctx) error {
	id := c.Params("id")
	cam, err := h.service.GetCameraByID(c.Context(), id)
	if err != nil {
		if err == ErrCameraNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Kamera tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(cam)
}

// UpdateCamera (Admin)
func (h *CameraHandler) UpdateCamera(c *fiber.Ctx) error {
	id := c.Params("id")
	var cam models.Camera
	if err := c.BodyParser(&cam); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if err := h.service.UpdateCamera(c.Context(), id, &cam); err != nil {
		if err == ErrCameraNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Kamera tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(cam)
}

// DeleteCamera (Admin)
func (h *CameraHandler) DeleteCamera(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.DeleteCamera(c.Context(), id); err != nil {
		if err == ErrCameraNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Kamera tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Kamera berhasil dihapus"})
}