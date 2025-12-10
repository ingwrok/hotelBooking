package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ingwrok/hotelBooking/internal/adapters/primary/web/dto"
	"github.com/ingwrok/hotelBooking/internal/core/domain"
	"github.com/ingwrok/hotelBooking/internal/core/services"
	"github.com/ingwrok/hotelBooking/internal/core/utils"
)

type RoomHandler struct {
	svc *services.RoomService
}

func NewRoomHandler(s *services.RoomService) *RoomHandler {
	return &RoomHandler{svc: s}
}

func (h * RoomHandler) CreateRoom(c *fiber.Ctx) error {
	ctx, cancel := buildCtx(c)
	defer cancel()

	rtid,err := c.ParamsInt("room_type_id")
	if err != nil || rtid <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "invalid room type ID"})
	}

	var req dto.RoomRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request body"})
	}

	const InitialRoomStatus = "available"

	room,err := h.svc.AddRoom(ctx, &domain.Room{
		RoomTypeID: rtid,
		RoomNumber: req.RoomNumber,
		Status:     InitialRoomStatus,
	})
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(201).JSON(dto.RoomResponse{
		RoomID:     room.RoomID,
		RoomTypeID: room.RoomTypeID,
		RoomNumber: room.RoomNumber,
		Status:     room.Status,
	})
}

func (h * RoomHandler) RemoveRoom(c *fiber.Ctx) error {
	ctx, cancel := buildCtx(c)
	defer cancel()

	id,err := c.ParamsInt("room_id")
	if err != nil || id <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "invalid room ID"})
	}

	err = h.svc.RemoveRoom(ctx, id)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(200).JSON(fiber.Map{"message": "room deleted successfully"})
}

func (h *RoomHandler) ChangeRoomStatus(c *fiber.Ctx) error {
	ctx, cancel := buildCtx(c)
	defer cancel()

	id, err := c.ParamsInt("room_id")
	if err != nil || id <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "invalid room ID"})
	}

	var req dto.RoomStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request body"})
	}

	err = h.svc.ChangeRoomStatus(ctx, id, req.Status)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(200).JSON(fiber.Map{"message": "room status updated"})
}

func (h *RoomHandler) GetRoom(c *fiber.Ctx) error {
	ctx, cancel := buildCtx(c)
	defer cancel()

	id, err := c.ParamsInt("room_id")
	if err != nil || id <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "invalid room ID"})
	}

	room, err := h.svc.GetRoom(ctx, id)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(200).JSON(dto.RoomResponse{
		RoomID:     room.RoomID,
		RoomTypeID: room.RoomTypeID,
		RoomNumber: room.RoomNumber,
		Status:     room.Status,
	})
}

func (h *RoomHandler) ListRooms(c *fiber.Ctx) error {
	ctx, cancel := buildCtx(c)
	defer cancel()

	rooms, err := h.svc.ListRooms(ctx)
	if err != nil {
		return handleError(c, err)
	}

	resRoom := make([]dto.RoomResponse, len(rooms))
	for i, r := range rooms {
		resRoom[i] = dto.RoomResponse{
			RoomID:     r.RoomID,
			RoomTypeID: r.RoomTypeID,
			RoomNumber: r.RoomNumber,
			Status:     r.Status,
		}
	}

	return c.Status(200).JSON(resRoom)
}


func (h *RoomHandler) BlockRoom(c *fiber.Ctx) error {
	ctx, cancel := buildCtx(c)
	defer cancel()

	var req dto.RoomBlockRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request body"})
	}

	startDate, err := utils.ParseDate(req.StartDate,"block start date")
	if err != nil {
		return handleError(c,err)
	}

	endDate, err := utils.ParseDate(req.EndDate,"block end date")
	if err != nil {
		return handleError(c,err)
	}

	block := &domain.RoomBlock{
		RoomID:    req.RoomID,
		StartDate: startDate,
		EndDate:   endDate,
		Reason:    req.Reason,
	}

	err = h.svc.BlockRoom(ctx, block)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(201).JSON(fiber.Map{"message": "room blocked"})
}

func (h *RoomHandler) GetRoomBlocks(c *fiber.Ctx) error {
	ctx, cancel := buildCtx(c)
	defer cancel()

	roomID, err := c.ParamsInt("room_id")
	if err != nil || roomID <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "invalid room ID"})
	}

	blocks, err := h.svc.GetRoomBlocks(ctx, roomID)
	if err != nil {
		return handleError(c, err)
	}

	resBlocks := make([]dto.RoomBlockResponse, len(blocks))
	for i, r := range blocks {
		resBlocks[i] = dto.RoomBlockResponse{
			RoomBlockID: r.RoomBlockID,
			RoomID:      r.RoomID,
			StartDate:   r.StartDate,
			EndDate:     r.EndDate,
			Reason:      r.Reason,
		}
	}

	return c.Status(200).JSON(resBlocks)
}


func (h *RoomHandler) UnblockRoom(c *fiber.Ctx) error {
	ctx, cancel := buildCtx(c)
	defer cancel()

	blockID, err := c.ParamsInt("block_id")
	if err != nil || blockID <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "invalid block ID"})
	}

	err = h.svc.UnblockRoom(ctx, blockID)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(200).JSON(fiber.Map{"message": "room unblocked"})
}


func (h *RoomHandler) CountAvailableRooms(c *fiber.Ctx) error {
	ctx, cancel := buildCtx(c)
	defer cancel()

	var req dto.AvailabilityRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request body"})
	}

	result, err := h.svc.CountAvailableRooms(ctx, req.CheckIn, req.CheckOut)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(200).JSON(result)
}

func (h *RoomHandler) FindAvailableRoom(c *fiber.Ctx) error {
	ctx, cancel := buildCtx(c)
	defer cancel()

	var req dto.FindRoomRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request body"})
	}

	roomID, err := h.svc.FindAvailableRoom(ctx, req.RoomTypeID, req.CheckIn, req.CheckOut)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(200).JSON(fiber.Map{"room_id": roomID})
}