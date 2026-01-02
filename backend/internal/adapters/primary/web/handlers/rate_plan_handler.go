package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ingwrok/hotelBooking/internal/adapters/primary/web/dto"
	"github.com/ingwrok/hotelBooking/internal/core/domain"
	"github.com/ingwrok/hotelBooking/internal/core/services"
	"github.com/ingwrok/hotelBooking/internal/core/utils"
)

type RatePlanHandler struct {
	svc *services.RatePlanService
}

func NewRatePlanHandler(s *services.RatePlanService) *RatePlanHandler {
	return &RatePlanHandler{svc: s}
}

func (h *RatePlanHandler) CreateRatePlan(c *fiber.Ctx) error {
	ctx, cancel := buildCtx(c)
	defer cancel()

	var req dto.RatePlanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request body"})
	}

	ratePlan, err := h.svc.AddRatePlan(ctx, &domain.RatePlan{
		Name:             req.Name,
		Description:      req.Description,
		IsSpecialPackage: req.IsSpecialPackage,
		AllowFreeCancel:  req.AllowFreeCancel,
		AllowPayLater:    req.AllowPayLater,
	})
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(201).JSON(dto.RatePlanResponse{
		RatePlanID:       ratePlan.RatePlanID,
		Name:             ratePlan.Name,
		Description:      ratePlan.Description,
		IsSpecialPackage: ratePlan.IsSpecialPackage,
		AllowFreeCancel:  ratePlan.AllowFreeCancel,
		AllowPayLater:    ratePlan.AllowPayLater,
		CreatedAt:        utils.ToThaiTime(ratePlan.CreatedAt),
		UpdatedAt:        utils.ToThaiTime(ratePlan.UpdatedAt),
	})
}

func (h *RatePlanHandler) UpdateRatePlan(c *fiber.Ctx) error {
	ctx, cancel := buildCtx(c)
	defer cancel()

	id, err := c.ParamsInt("rate_plan_id")
	if err != nil || id <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "invalid rate plan ID"})
	}

	var req dto.RatePlanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request body"})
	}

	err = h.svc.ChangeRatePlan(ctx, &domain.RatePlan{
		RatePlanID:       id,
		Name:             req.Name,
		Description:      req.Description,
		IsSpecialPackage: req.IsSpecialPackage,
		AllowFreeCancel:  req.AllowFreeCancel,
		AllowPayLater:    req.AllowPayLater,
	})
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(200).JSON(fiber.Map{"message": "rate plan updated successfully"})
}

func (h *RatePlanHandler) RemoveRatePlan(c *fiber.Ctx) error {
	ctx, cancel := buildCtx(c)
	defer cancel()

	id, err := c.ParamsInt("rate_plan_id")
	if err != nil || id <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "invalid rate plan ID"})
	}

	err = h.svc.RemoveRatePlan(ctx, id)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(200).JSON(fiber.Map{"message": "rate plan deleted successfully"})
}

func (h *RatePlanHandler) GetRatePlan(c *fiber.Ctx) error {
	ctx, cancel := buildCtx(c)
	defer cancel()

	id, err := c.ParamsInt("rate_plan_id")
	if err != nil || id <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "invalid rate plan ID"})
	}

	ratePlan, err := h.svc.GetRatePlan(ctx, id)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(200).JSON(dto.RatePlanResponse{
		RatePlanID:       ratePlan.RatePlanID,
		Name:             ratePlan.Name,
		Description:      ratePlan.Description,
		IsSpecialPackage: ratePlan.IsSpecialPackage,
		AllowFreeCancel:  ratePlan.AllowFreeCancel,
		AllowPayLater:    ratePlan.AllowPayLater,
		CreatedAt:        utils.ToThaiTime(ratePlan.CreatedAt),
		UpdatedAt:        utils.ToThaiTime(ratePlan.UpdatedAt),
	})
}

func (h *RatePlanHandler) ListRatePlans(c *fiber.Ctx) error {
	ctx, cancel := buildCtx(c)
	defer cancel()

	ratePlans, err := h.svc.ListRatePlans(ctx)
	if err != nil {
		return handleError(c, err)
	}

	resRatePlans := make([]dto.RatePlanResponse, len(ratePlans))
	for i, rp := range ratePlans {
		resRatePlans[i] = dto.RatePlanResponse{
			RatePlanID:       rp.RatePlanID,
			Name:             rp.Name,
			Description:      rp.Description,
			IsSpecialPackage: rp.IsSpecialPackage,
			AllowFreeCancel:  rp.AllowFreeCancel,
			AllowPayLater:    rp.AllowPayLater,
			CreatedAt:        utils.ToThaiTime(rp.CreatedAt),
			UpdatedAt:        utils.ToThaiTime(rp.UpdatedAt),
		}
	}

	return c.Status(200).JSON(resRatePlans)
}

func (h *RatePlanHandler) GetPrice(c *fiber.Ctx) error {
	ctx, cancel := buildCtx(c)
	defer cancel()

	rateplanID, err := c.ParamsInt("rate_plan_id")
	if err != nil || rateplanID <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "invalid rate plan ID"})
	}
	roomTypeID, err := c.ParamsInt("room_type_id")
	if err != nil || roomTypeID <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "invalid room type ID"})
	}

	price, err := h.svc.GetPrice(ctx, roomTypeID, rateplanID)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(200).JSON(fiber.Map{"price": price})
}

func (h *RatePlanHandler) UpdateRoomTypePrice(c *fiber.Ctx) error {
	ctx, cancel := buildCtx(c)
	defer cancel()

	rateplanID, err := c.ParamsInt("rate_plan_id")
	if err != nil || rateplanID <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "invalid rate plan ID"})
	}
	roomTypeID, err := c.ParamsInt("room_type_id")
	if err != nil || roomTypeID <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "invalid room type ID"})
	}

	type priceRequest struct {
		Price float64 `json:"price"`
	}

	var req priceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request body"})
	}

	err = h.svc.ChangeRoomTypePrice(ctx, roomTypeID, rateplanID, req.Price)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(200).JSON(fiber.Map{"message": "room type price updated successfully"})
}

func (h *RatePlanHandler) RemoveRoomTypePrice(c *fiber.Ctx) error {
	ctx, cancel := buildCtx(c)
	defer cancel()

	rateplanID, err := c.ParamsInt("rate_plan_id")
	if err != nil || rateplanID <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "invalid rate plan ID"})
	}
	roomTypeID, err := c.ParamsInt("room_type_id")
	if err != nil || roomTypeID <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "invalid room type ID"})
	}

	err = h.svc.RemoveRoomTypePrice(ctx, roomTypeID, rateplanID)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(200).JSON(fiber.Map{"message": "room type price deleted successfully"})
}

func (h *RatePlanHandler) ListRatePlansByRoomType(c *fiber.Ctx) error {
	ctx, cancel := buildCtx(c)
	defer cancel()

	roomTypeID, err := c.ParamsInt("room_type_id")
	if err != nil || roomTypeID <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "invalid room type ID"})
	}

	ratePlans, err := h.svc.ListRatePlansByRoomType(ctx, roomTypeID)
	if err != nil {
		return handleError(c, err)
	}

	resRatePlans := make([]dto.RatePlanFullResponse, len(ratePlans))
	for i, rp := range ratePlans {
		resRatePlans[i] = dto.RatePlanFullResponse{
			RatePlanID:       rp.RatePlanID,
			Name:             rp.Name,
			Description:      rp.Description,
			IsSpecialPackage: rp.IsSpecialPackage,
			AllowFreeCancel:  rp.AllowFreeCancel,
			AllowPayLater:    rp.AllowPayLater,
			Price:           rp.Price,
			CreatedAt:        utils.ToThaiTime(rp.CreatedAt),
			UpdatedAt:        utils.ToThaiTime(rp.UpdatedAt),
		}
	}

	return c.Status(200).JSON(resRatePlans)
}
