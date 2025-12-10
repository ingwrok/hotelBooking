package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ingwrok/hotelBooking/internal/adapters/primary/web/dto"
	"github.com/ingwrok/hotelBooking/internal/core/domain"
	"github.com/ingwrok/hotelBooking/internal/core/services"
)

type RoomTypeHandler struct {
	svc *services.RoomTypeService
}

func NewRoomTypeHandler(s *services.RoomTypeService) *RoomTypeHandler {
	return &RoomTypeHandler{svc: s}
}

func (h *RoomTypeHandler) CreateRoomType(c *fiber.Ctx) error {
	ctx,cancel := buildCtx(c)
	defer cancel()

	var req dto.RoomTypeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request body"})
	}

	roomType,err := h.svc.AddRoomType(ctx, &domain.RoomType{
		Name:        req.Name,
		Description: req.Description,
		SizeSQM:     req.SizeSQM,
		BedType:     req.BedType,
		Capacity:    req.Capacity,
		PictureURL:  req.PictureURL,
		AmenityIDs:  req.AmenityIDs,
	})
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(201).JSON(dto.RoomTypeResponse{
		RoomTypeID:  roomType.RoomTypeID,
		Name:        roomType.Name,
		Description: roomType.Description,
		SizeSQM:     roomType.SizeSQM,
		BedType:     roomType.BedType,
		Capacity:    roomType.Capacity,
		PictureURL:  roomType.PictureURL,
	})
}
func (h *RoomTypeHandler) UpdateRoomType(c *fiber.Ctx) error {
	ctx,cancel := buildCtx(c)
	defer cancel()

	id,err := c.ParamsInt("room_type_id")
	if err != nil || id <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "invalid room type ID"})
	}

	var req dto.RoomTypeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request body"})
	}

	err = h.svc.ChangeRoomType(ctx, &domain.RoomType{
		RoomTypeID:  id,
		Name:        req.Name,
		Description: req.Description,
		SizeSQM:     req.SizeSQM,
		BedType:     req.BedType,
		Capacity:    req.Capacity,
		PictureURL:  req.PictureURL,
		AmenityIDs:  req.AmenityIDs,
	})
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(200).JSON(fiber.Map{"message": "room type updated successfully"})
}

func (h *RoomTypeHandler) RemoveRoomType(c *fiber.Ctx) error {
	ctx,cancel := buildCtx(c)
	defer cancel()

	id,err := c.ParamsInt("room_type_id")
	if err != nil || id <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "invalid room type ID"})
	}

	err = h.svc.RemoveRoomType(ctx,id)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(200).JSON(fiber.Map{"message": "room type deleted successfully"})
}

func (h *RoomTypeHandler) GetRoomType(c *fiber.Ctx) error {
	ctx,cancel := buildCtx(c)
	defer cancel()

	id,err := c.ParamsInt("room_type_id")
	if err != nil || id <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "invalid room type ID"})
	}

	roomType,err := h.svc.GetRoomType(ctx,id)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(200).JSON(dto.RoomTypeResponse{
		RoomTypeID:  roomType.RoomTypeID,
		Name:        roomType.Name,
		Description: roomType.Description,
		SizeSQM:     roomType.SizeSQM,
		BedType:     roomType.BedType,
		Capacity:    roomType.Capacity,
		PictureURL:  roomType.PictureURL,
	})
}

func (h *RoomTypeHandler) ListRoomTypes(c *fiber.Ctx) error {
	ctx,cancel := buildCtx(c)
	defer cancel()

	roomTypes,err := h.svc.ListRoomTypes(ctx)
	if err != nil {
		return handleError(c, err)
	}

	resRoomTypes := make([]dto.RoomTypeResponse,len(roomTypes))
	for i,rt := range roomTypes {
		resRoomTypes[i] = dto.RoomTypeResponse{
			RoomTypeID:  rt.RoomTypeID,
			Name:        rt.Name,
			Description: rt.Description,
			SizeSQM:     rt.SizeSQM,
			BedType:     rt.BedType,
			Capacity:    rt.Capacity,
			PictureURL:  rt.PictureURL,
		}
	}

	return c.Status(200).JSON(resRoomTypes)
}

func (h *RoomTypeHandler) GetRoomTypeFullDetail(c *fiber.Ctx) error {
	ctx,cancel := buildCtx(c)
	defer cancel()

	id,err := c.ParamsInt("room_type_id")
	if err != nil || id <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "invalid room type ID"})
	}

	roomTypeDetails,err := h.svc.GetRoomTypeFullDetail(ctx,id)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(200).JSON(dto.RoomTypeDetailResponse{
		RoomTypeID:  roomTypeDetails.RoomTypeID,
		Name:        roomTypeDetails.Name,
		Description: roomTypeDetails.Description,
		SizeSQM:     roomTypeDetails.SizeSQM,
		BedType:     roomTypeDetails.BedType,
		Capacity:    roomTypeDetails.Capacity,
		PictureURL:  roomTypeDetails.PictureURL,
		Amenities:   roomTypeDetails.Amenities,
	})
}