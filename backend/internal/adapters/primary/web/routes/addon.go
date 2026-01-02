package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ingwrok/hotelBooking/internal/adapters/primary/web/handlers"
)

func AddonRoutes(app *fiber.App, h *handlers.AddonHandler){
	api := app.Group("/api")

	categories := api.Group("/addon-categories")
	categories.Post("/", h.CreateAddonCategory)
	categories.Get("/", h.ListAddonCategories)
	categories.Get("/:addon_category_id", h.GetAddonCategory)
	categories.Put("/:addon_category_id", h.UpdateAddonCategory)
	categories.Delete("/:addon_category_id", h.DeleteAddonCategory)

	addons := api.Group("/addons")
	addons.Post("/", h.CreateAddon)
	addons.Get("/", h.ListAddons)
	addons.Get("/:addon_id", h.GetAddon)
	addons.Put("/:addon_id", h.UpdateAddon)
	addons.Delete("/:addon_id", h.DeleteAddon)
	addons.Get("/category/:addon_category_id", h.ListAddonsByCategory)
}