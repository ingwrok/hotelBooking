package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ingwrok/hotelBooking/internal/adapters/primary/web/handlers"
	"github.com/ingwrok/hotelBooking/internal/adapters/primary/web/middleware"
	"github.com/ingwrok/hotelBooking/internal/core/services"
)

func AmenityRoutes(app *fiber.App, h *handlers.AmenityHandler, userSvc *services.UserService) {
	amenities := app.Group("/api/amenities")

	amenities.Get("/:amenity_id", h.GetAmenity)
	amenities.Get("/", h.ListAmenities)

	admin := amenities.Group("/", middleware.AuthMiddleware(userSvc), middleware.VerifyAdmin())
	admin.Post("/", h.CreateAmenity)
	admin.Patch("/:amenity_id", h.UpdateAmenity)
	admin.Delete("/:amenity_id", h.RemoveAmenity)
}
