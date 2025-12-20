package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ingwrok/hotelBooking/internal/adapters/primary/web/dto"
	"github.com/ingwrok/hotelBooking/internal/core/domain"
	"github.com/ingwrok/hotelBooking/internal/core/services"
)

type AddonHandler struct {
	svc *services.AddonService
}

func NewAddonHandler(s *services.AddonService) *AddonHandler {
	return &AddonHandler{svc: s}
}

// --- Addon Category ---
func (h *AddonHandler)CreateAddonCategory(c *fiber.Ctx) error{
	ctx,cancel := buildCtx(c)
	defer cancel()

	var req dto.AddonCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request body"})
	}

	addon,err := h.svc.AddAddonCategory(ctx, &domain.AddonCategory{
		Name: req.Name,
	})
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(201).JSON(dto.AddonCategoryResponse{
		CategoryID: addon.CategoryID,
		Name:       addon.Name,
	})
}

func (h *AddonHandler)GetAddonCategory(c *fiber.Ctx) error{
	ctx,cancel := buildCtx(c)
	defer cancel()

	id,err := c.ParamsInt("addon_category_id")
	if err != nil || id <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "invalid addon category ID"})
	}

	addon,err := h.svc.GetAddonCategory(ctx,id)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(200).JSON(dto.AddonCategoryResponse{
		CategoryID: addon.CategoryID,
		Name:       addon.Name,
	})
}

func (h *AddonHandler)ListAddonCategories(c *fiber.Ctx) error{
	ctx,cancel := buildCtx(c)
	defer cancel()

	addons,err := h.svc.ListAddonCategories(ctx)
	if err != nil {
		return handleError(c, err)
	}

	resAddonCategories := make([]dto.AddonCategoryResponse, len(addons))
	for i, addon := range addons {
		resAddonCategories[i] = dto.AddonCategoryResponse{
			CategoryID: addon.CategoryID,
			Name:       addon.Name,
		}
	}

	return c.Status(200).JSON(resAddonCategories)
}

func (h *AddonHandler)UpdateAddonCategory(c *fiber.Ctx) error{
	ctx,cancel := buildCtx(c)
	defer cancel()

	id,err := c.ParamsInt("addon_category_id")
	if err != nil || id <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "invalid addon category ID"})
	}

	var req dto.AddonCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request body"})
	}

	err = h.svc.ChangeAddonCategory(ctx, &domain.AddonCategory{
		CategoryID: id,
		Name: req.Name,
	})

	if err != nil {
		return handleError(c, err)
	}

	return c.Status(200).JSON(fiber.Map{"message": "addon category updated successfully"})
}

func (h *AddonHandler)DeleteAddonCategory(c *fiber.Ctx) error{
	ctx,cancel := buildCtx(c)
	defer cancel()

	id,err := c.ParamsInt("addon_category_id")
	if err != nil || id <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "invalid addon category ID"})
	}

	err = h.svc.RemoveAddonCategory(ctx,id)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(200).JSON(fiber.Map{"message": "addon category deleted successfully"})
}

// // --- Addon ---
func (h *AddonHandler)CreateAddon(c *fiber.Ctx) error{
	ctx,cancel := buildCtx(c)
	defer cancel()

	var req dto.AddonRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request body"})
	}

	addon,err := h.svc.AddAddon(ctx, &domain.Addon{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		CategoryID:  req.CategoryID,
		UnitName:    req.UnitName,
	})

	if err != nil {
		return handleError(c, err)
	}

	return c.Status(201).JSON(dto.AddonResponse{
		AddonID:     addon.AddonID,
		Name:        addon.Name,
		Description: addon.Description,
		Price:       addon.Price,
		CategoryID:  addon.CategoryID,
		UnitName:    addon.UnitName,
	})
}

func (h *AddonHandler)GetAddon(c *fiber.Ctx) error {
	ctx,cancel := buildCtx(c)
	defer cancel()

	id,err := c.ParamsInt("addon_id")
	if err != nil || id <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "invalid addon ID"})
	}

	addon,err := h.svc.GetAddon(ctx,id)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(200).JSON(dto.AddonResponse{
		AddonID:     addon.AddonID,
		Name:        addon.Name,
		Description: addon.Description,
		Price:       addon.Price,
		CategoryID:  addon.CategoryID,
		UnitName:    addon.UnitName,
	})
}

func (h *AddonHandler)ListAddons(c *fiber.Ctx) error{ 
	ctx,cancel := buildCtx(c)
	defer cancel()

	addons,err := h.svc.ListAddons(ctx)
	if err != nil {
		return handleError(c, err)
	}

	resAddons := make([]dto.AddonResponse, len(addons))
	for i, addon := range addons {
		resAddons[i] = dto.AddonResponse{
			AddonID:     addon.AddonID,
			Name:        addon.Name,
			Description: addon.Description,
			Price:       addon.Price,
			CategoryID:  addon.CategoryID,
			UnitName:    addon.UnitName,
		}
	}

	return c.Status(200).JSON(resAddons)
}

func (h *AddonHandler)ListAddonsByCategory(c *fiber.Ctx) error{
	ctx,cancel := buildCtx(c)
	defer cancel()

	id,err := c.ParamsInt("addon_category_id")
	if err != nil || id <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "invalid addon category ID"})
	}

	addons,err := h.svc.ListAddonsByCategory(ctx,id)
	if err != nil {
		return handleError(c, err)
	}

	resAddons := make([]dto.AddonResponse, len(addons))
	for i, addon := range addons {
		resAddons[i] = dto.AddonResponse{
			AddonID:     addon.AddonID,
			Name:        addon.Name,
			Description: addon.Description,
			Price:       addon.Price,
			CategoryID:  addon.CategoryID,
		}
	}

	return c.Status(200).JSON(resAddons)
}

func (h *AddonHandler)UpdateAddon(c *fiber.Ctx) error{
	ctx,cancel := buildCtx(c)
	defer cancel()

	id,err := c.ParamsInt("addon_id")
	if err != nil || id <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "invalid addon ID"})
	}

	var req dto.AddonRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request body"})
	}

	err = h.svc.ChangeAddon(ctx, &domain.Addon{
		AddonID:     id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		CategoryID:  req.CategoryID,
		UnitName:    req.UnitName,
	})
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(200).JSON(fiber.Map{"message": "addon updated successfully"})
}

func (h *AddonHandler)DeleteAddon(c *fiber.Ctx) error{
	ctx,cancel := buildCtx(c)
	defer cancel()

	id,err := c.ParamsInt("addon_id")
	if err != nil || id <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "invalid addon ID"})
	}

	err = h.svc.RemoveAddon(ctx,id)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(200).JSON(fiber.Map{"message": "addon deleted successfully"})
}

