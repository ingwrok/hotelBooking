package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ingwrok/hotelBooking/internal/adapters/primary/web/handlers"
)

func AmenityRoutes(app *fiber.App, h *handlers.AmenityHandler) {

	amenities := app.Group("/api/amenities")

	amenities.Post("/", h.CreateAmenity)

	amenities.Get("/:amenity_id",  h.GetAmenity)
	amenities.Get("/", h.ListAmenities)

	amenities.Patch("/:amenity_id", h.UpdateAmenity)

	amenities.Delete("/:amenity_id", h.RemoveAmenity)
}