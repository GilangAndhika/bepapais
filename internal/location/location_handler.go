package location

import (
	"github.com/gofiber/fiber/v2"
)

type LocationHandler struct {
	service Service
}

func NewLocationHandler(service Service) *LocationHandler {
	return &LocationHandler{service: service}
}

type CreateLocationRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
	Type string `json:"type"` // e.g., "Kecamatan"
}

// CreateLocation handler
func (h *LocationHandler) CreateLocation(c *fiber.Ctx) error {
	var req CreateLocationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Validasi
	if req.Name == "" || req.Slug == "" || req.Type == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Name, Slug, dan Type tidak boleh kosong"})
	}

	err := h.service.CreateLocation(c.Context(), req.Name, req.Slug, req.Type)
	if err != nil {
		if err == ErrLocationExists {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Slug sudah digunakan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Lokasi berhasil dibuat"})
}

// GetAllLocations handler (Rute Publik)
func (h *LocationHandler) GetAllLocations(c *fiber.Ctx) error {
	locations, err := h.service.GetAllLocations(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(locations)
}

// GetLocationByID handler (Rute Publik)
func (h *LocationHandler) GetLocationByID(c *fiber.Ctx) error {
	id := c.Params("id")
	location, err := h.service.GetLocationByID(c.Context(), id)
	if err != nil {
		if err == ErrLocationNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Lokasi tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(location)
}

// UpdateLocation handler
func (h *LocationHandler) UpdateLocation(c *fiber.Ctx) error {
	id := c.Params("id")
	var req CreateLocationRequest // Pakai struct yang sama dengan create
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	err := h.service.UpdateLocation(c.Context(), id, req.Name, req.Slug, req.Type)
	if err != nil {
		if err == ErrLocationNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Lokasi tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Lokasi berhasil diupdate"})
}

// DeleteLocation handler
func (h *LocationHandler) DeleteLocation(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.service.DeleteLocation(c.Context(), id)
	if err != nil {
		if err == ErrLocationNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Lokasi tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Lokasi berhasil dihapus"})
}