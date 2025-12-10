package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ingwrok/hotelBooking/internal/adapters/primary/web/handlers"
)

func RoomTypeRoutes(app *fiber.App, h *handlers.RoomTypeHandler) {
	roomTypes := app.Group("/api/room_types")

	roomTypes.Get("/:room_type_id/full", h.GetRoomTypeFullDetail)
	roomTypes.Get("/:room_type_id",h.GetRoomType)
	roomTypes.Get("/", h.ListRoomTypes)

	roomTypes.Post("/", h.CreateRoomType)

	roomTypes.Patch("/:room_type_id", h.UpdateRoomType)

	roomTypes.Delete("/:room_type_id", h.RemoveRoomType)
}