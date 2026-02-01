package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ingwrok/hotelBooking/internal/adapters/primary/web/handlers"
	"github.com/ingwrok/hotelBooking/internal/adapters/primary/web/middleware"
	"github.com/ingwrok/hotelBooking/internal/core/services"
)

func RoomTypeRoutes(app *fiber.App, h *handlers.RoomTypeHandler, userSvc *services.UserService) {
	roomTypes := app.Group("/api/room_types")

	roomTypes.Get("/:room_type_id/full", h.GetRoomTypeFullDetail)
	roomTypes.Get("/:room_type_id", h.GetRoomType)
	roomTypes.Get("/", h.ListRoomTypes)

	admin := roomTypes.Group("", middleware.AuthMiddleware(userSvc), middleware.VerifyAdmin())
	admin.Post("/upload", h.UploadImage)
	admin.Post("/", h.CreateRoomType)
	admin.Patch("/:room_type_id", h.UpdateRoomType)
	admin.Delete("/:room_type_id", h.RemoveRoomType)
}
