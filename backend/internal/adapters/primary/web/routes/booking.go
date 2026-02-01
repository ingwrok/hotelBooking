package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ingwrok/hotelBooking/internal/adapters/primary/web/handlers"
	"github.com/ingwrok/hotelBooking/internal/adapters/primary/web/middleware"
	"github.com/ingwrok/hotelBooking/internal/core/services"
)

func BookingRoutes(app *fiber.App, h *handlers.BookingHandler, userSvc *services.UserService, bookingSvc *services.BookingService) {
	bookings := app.Group("/api/bookings", middleware.AuthMiddleware(userSvc))

	bookings.Get("/my", h.GetBookings)

	bookings.Post("/", h.CreateBooking)

	// Admin Routes
	bookings.Get("/all", middleware.VerifyAdmin(), h.GetAllBookings)

	bookings.Get("/:booking_id", middleware.VerifyBookingOwner(bookingSvc), h.GetFullBooking)
	bookings.Get("/:booking_id/addons", middleware.VerifyBookingOwner(bookingSvc), h.GetAddons)
	bookings.Put("/:booking_id/addons", middleware.VerifyBookingOwner(bookingSvc), h.UpdateAddons)
	bookings.Post("/:booking_id/pay", middleware.VerifyBookingOwner(bookingSvc), h.SimulatePayment)

	bookings.Patch("/:booking_id/status", middleware.VerifyAdmin(), h.UpdateStatus)
}
