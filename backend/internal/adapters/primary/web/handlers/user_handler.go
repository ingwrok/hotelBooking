package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ingwrok/hotelBooking/internal/adapters/primary/web/dto"
	"github.com/ingwrok/hotelBooking/internal/common/errs"
	"github.com/ingwrok/hotelBooking/internal/core/services"
)

type UserHandler struct {
	svc *services.UserService
}

func NewUserHandler(s *services.UserService) *UserHandler {
	return &UserHandler{svc: s}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	ctx, cancel := buildCtx(c)
	defer cancel()

	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return handleError(c, errs.NewValidationError("invalid request body"))
	}

	created, err := h.svc.Register(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(dto.UserResponse{
		UserID:   created.UserID,
		Username: created.Username,
		Email:    created.Email,
		IsAdmin:  created.IsAdmin,
	})
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	ctx, cancel := buildCtx(c)
	defer cancel()

	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return handleError(c, errs.NewValidationError("request body incorrect format"))
	}

	token, user, err := h.svc.Login(ctx, req.Username, req.Password)
	if err != nil {
		return handleError(c, err)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
	})

	return c.Status(fiber.StatusOK).JSON(dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			UserID:   user.UserID,
			Username: user.Username,
			Email:    user.Email,
			IsAdmin:  user.IsAdmin,
		},
	})
}

func (h *UserHandler) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour), // Expire immediately
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
	})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "logged out successfully"})
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	ctx, cancel := buildCtx(c)
	defer cancel()

	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return handleError(c, errs.NewValidationError("invalid user id"))
	}

	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return handleError(c, errs.NewValidationError("request body incorrect format"))
	}

	if err := h.svc.UpdateUser(ctx, id, body); err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "user updated"})
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	ctx, cancel := buildCtx(c)
	defer cancel()

	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return handleError(c, errs.NewValidationError("invalid user id"))
	}

	if err := h.svc.DeleteUser(ctx, id); err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "user deleted"})
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	ctx, cancel := buildCtx(c)
	defer cancel()

	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return handleError(c, errs.NewValidationError("invalid user id"))
	}

	u, err := h.svc.GetUser(ctx, id)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(dto.UserResponse{
		UserID:   u.UserID,
		Username: u.Username,
		Email:    u.Email,
		IsAdmin:  u.IsAdmin,
	})
}

func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	ctx, cancel := buildCtx(c)
	defer cancel()

	users, err := h.svc.GetUsers(ctx)
	if err != nil {
		return handleError(c, err)
	}

	out := make([]dto.UserResponse, 0, len(users))
	for _, u := range users {
		out = append(out, dto.UserResponse{
			UserID:   u.UserID,
			Username: u.Username,
			Email:    u.Email,
			IsAdmin:  u.IsAdmin,
		})
	}
	return c.JSON(out)
}
