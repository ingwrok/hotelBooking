package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ingwrok/hotelBooking/internal/adapters/primary/web/dto"
	"github.com/ingwrok/hotelBooking/internal/core/domain"
	"github.com/ingwrok/hotelBooking/internal/core/services"
)

type AmenityHandler struct {
	svc *services.AmenityService
}

func NewAmenityHandler(s *services.AmenityService) *AmenityHandler {
	return &AmenityHandler{svc: s}
}

func (h *AmenityHandler) CreateAmenity(c *fiber.Ctx) error {
	ctx, cancel := buildCtx(c)
	defer cancel()

	var req dto.AmenityRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request body"})
	}

	amenity,err := h.svc.AddAmenity(ctx, &domain.Amenity{
		Name: req.Name,
	})
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(201).JSON(dto.AmenityResponse{
		AmenityID: amenity.AmenityID,
		Name:      amenity.Name,
	})
}

func (h *AmenityHandler) GetAmenity(c *fiber.Ctx) error {
	ctx, cancel := buildCtx(c)
	defer cancel()

	id, err := c.ParamsInt("amenity_id")
	if err != nil || id <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "invalid amenity ID"})
	}

	amenity, err := h.svc.GetAmenity(ctx, id)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(200).JSON(dto.AmenityResponse{
		AmenityID: amenity.AmenityID,
		Name:      amenity.Name,
	})
}

func (h *AmenityHandler) ListAmenities(c *fiber.Ctx) error {
	ctx, cancel := buildCtx(c)
	defer cancel()

	amenities, err := h.svc.ListAmenities(ctx)
	if err != nil {
		return handleError(c, err)
	}

	resAmenities := make([]dto.AmenityResponse, len(amenities))
	for i, a := range amenities {
		resAmenities[i] = dto.AmenityResponse{
			AmenityID: a.AmenityID,
			Name:      a.Name,
		}
	}

	return c.Status(200).JSON(resAmenities)
}

func (h *AmenityHandler) UpdateAmenity(c *fiber.Ctx) error {
	ctx, cancel := buildCtx(c)
	defer cancel()

	id, err := c.ParamsInt("amenity_id")
	if err != nil || id <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "invalid amenity ID"})
	}

	var req dto.AmenityRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request body"})
	}

	err = h.svc.ChangeAmenity(ctx, &domain.Amenity{
		AmenityID: id,
		Name:      req.Name,
	})
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(200).JSON(fiber.Map{"message": "amenity updated successfully"})
}

func (h *AmenityHandler) RemoveAmenity(c *fiber.Ctx) error{
	ctx, cancel := buildCtx(c)
	defer cancel()

	id, err := c.ParamsInt("amenity_id")
	if err != nil || id <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "invalid amenity ID"})
	}

	err = h.svc.RemoveAmenity(ctx, id)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(200).JSON(fiber.Map{"message": "amenity removed successfully"})
}