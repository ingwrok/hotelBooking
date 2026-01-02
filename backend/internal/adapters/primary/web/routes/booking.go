package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ingwrok/hotelBooking/internal/adapters/primary/web/handlers"
)

func BookingRoutes(app *fiber.App, h *handlers.BookingHandler) {
	bookings := app.Group("/api/bookings")

	bookings.Post("/", h.CreateBooking)
	bookings.Get("/:booking_id", h.GetFullBooking)
	bookings.Patch("/:booking_id/status", h.UpdateStatus)
	bookings.Get("/:booking_id/addons", h.GetAddons)
	bookings.Put("/:booking_id/addons", h.UpdateAddons)
}