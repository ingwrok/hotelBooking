package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ingwrok/hotelBooking/internal/adapters/primary/web/handlers"
	"github.com/ingwrok/hotelBooking/internal/adapters/primary/web/middleware"
	"github.com/ingwrok/hotelBooking/internal/core/services"
)

func AddonRoutes(app *fiber.App, h *handlers.AddonHandler, userSvc *services.UserService) {
	api := app.Group("/api")

	categories := api.Group("/addon-categories")
	categories.Get("/", h.ListAddonCategories)
	categories.Get("/:addon_category_id", h.GetAddonCategory)

	categoriesAdmin := categories.Group("/", middleware.AuthMiddleware(userSvc), middleware.VerifyAdmin())
	categoriesAdmin.Post("/", h.CreateAddonCategory)
	categoriesAdmin.Put("/:addon_category_id", h.UpdateAddonCategory)
	categoriesAdmin.Delete("/:addon_category_id", h.DeleteAddonCategory)

	addons := api.Group("/addons")
	addons.Get("/", h.ListAddons)
	addons.Get("/:addon_id", h.GetAddon)
	addons.Get("/category/:addon_category_id", h.ListAddonsByCategory)

	addonsAdmin := addons.Group("/", middleware.AuthMiddleware(userSvc), middleware.VerifyAdmin())
	addonsAdmin.Post("/upload", h.UploadImage)
	addonsAdmin.Post("/", h.CreateAddon)
	addonsAdmin.Put("/:addon_id", h.UpdateAddon)
	addonsAdmin.Delete("/:addon_id", h.DeleteAddon)
}
