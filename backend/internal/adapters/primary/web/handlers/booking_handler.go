package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ingwrok/hotelBooking/internal/adapters/primary/web/dto"
	"github.com/ingwrok/hotelBooking/internal/adapters/primary/web/middleware"

	"github.com/ingwrok/hotelBooking/internal/core/domain"
	"github.com/ingwrok/hotelBooking/internal/core/services"
	"github.com/ingwrok/hotelBooking/internal/core/utils"
)

type BookingHandler struct {
	svc *services.BookingService
}

func NewBookingHandler(s *services.BookingService) *BookingHandler {
	return &BookingHandler{svc: s}
}

func (h *BookingHandler) CreateBooking(c *fiber.Ctx) error {
	ctx, cancel := buildCtx(c)
	defer cancel()

	var req dto.BookingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request body"})
	}

	// Get authenticated user from context (set by middleware)
	authUser := middleware.GetAuthUser(c)
	if authUser == nil {
		return c.Status(401).JSON(fiber.Map{"message": "unauthorized"})
	}

	domainAddons := dto.ToDomainBookingAddons(req.BookingAddon)

	checkin, err := utils.ParseDate(req.CheckInDate, "check in date")
	if err != nil {
		return handleError(c, err)
	}

	checkout, err := utils.ParseDate(req.CheckOutDate, "check out date")
	if err != nil {
		return handleError(c, err)
	}

	booking, err := h.svc.AddBooking(ctx, &domain.Booking{
		UserID:       authUser.ID,
		RatePlanID:   req.RatePlanID,
		CheckInDate:  checkin,
		CheckOutDate: checkout,
		NumAdults:    req.NumAdults,
		Email:        req.Email,
		BookingAddon: domainAddons,
	}, req.RoomTypeID)
	if err != nil {
		return handleError(c, err)
	}

	// Fetch full details (populated with joins) for the response
	fullBooking, err := h.svc.GetFullDetails(ctx, booking.BookingID)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(201).JSON(dto.ToBookingResponse(fullBooking))
}

func (h *BookingHandler) UpdateStatus(c *fiber.Ctx) error {
	ctx, cancel := buildCtx(c)
	defer cancel()

	id, err := c.ParamsInt("booking_id")
	if err != nil || id <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "invalid booking ID"})
	}

	type update struct {
		Status string `json:"status"`
	}

	var req update
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request body"})
	}

	err = h.svc.ChangeStatus(ctx, id, req.Status)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(200).JSON(fiber.Map{"message": "status updated successfully"})
}

func (h *BookingHandler) GetAddons(c *fiber.Ctx) error {
	ctx, cancel := buildCtx(c)
	defer cancel()

	id, err := c.ParamsInt("booking_id")
	if err != nil || id <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "invalid booking ID"})
	}

	addons, err := h.svc.GetAddonDetails(ctx, id)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(200).JSON(dto.ToBookingAddonResponses(addons))
}

func (h *BookingHandler) GetFullBooking(c *fiber.Ctx) error {
	ctx, cancel := buildCtx(c)
	defer cancel()

	id, err := c.ParamsInt("booking_id")
	if err != nil || id <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "invalid booking ID"})
	}

	booking, err := h.svc.GetFullDetails(ctx, id)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(200).JSON(dto.ToBookingResponse(booking))
}

func (h *BookingHandler) UpdateAddons(c *fiber.Ctx) error {
	ctx, cancel := buildCtx(c)
	defer cancel()

	id, err := c.ParamsInt("booking_id")
	if err != nil || id <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "invalid booking ID"})
	}

	var req []dto.BookingAddonRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request body"})
	}

	domainAddons := dto.ToDomainBookingAddons(req)

	err = h.svc.ModifyBookingAddons(ctx, id, domainAddons)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(200).JSON(fiber.Map{"message": "booking addons updated successfully"})
}

func (h *BookingHandler) GetBookings(c *fiber.Ctx) error {
	ctx, cancel := buildCtx(c)
	defer cancel()

	au := middleware.GetAuthUser(c)
	if au == nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid user ID"})
	}

	bookings, err := h.svc.GetMyHistory(ctx, au.ID)
	if err != nil {
		return handleError(c, err)
	}

	resBookings := make([]*dto.BookingResponse, len(bookings))
	for i, booking := range bookings {
		resBookings[i] = dto.ToBookingResponse(booking)
	}

	return c.Status(200).JSON(resBookings)
}

func (h *BookingHandler) SimulatePayment(c *fiber.Ctx) error {
	ctx, cancel := buildCtx(c)
	defer cancel()

	bookingID, err := c.ParamsInt("booking_id")
	if err != nil || bookingID <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "invalid booking ID"})
	}

	// Just update status to confirmed for mock
	err = h.svc.ChangeStatus(ctx, bookingID, "confirmed")
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(200).JSON(fiber.Map{"message": "payment successful"})
}

func (h *BookingHandler) GetAllBookings(c *fiber.Ctx) error {
	ctx, cancel := buildCtx(c)
	defer cancel()

	bookings, err := h.svc.GetAllBookings(ctx)
	if err != nil {
		return handleError(c, err)
	}

	// Map to DTO
	var res []*dto.BookingResponse
	for _, b := range bookings {
		res = append(res, dto.ToBookingResponse(b))
	}

	return c.Status(200).JSON(res)
}
