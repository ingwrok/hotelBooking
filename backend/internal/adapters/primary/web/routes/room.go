package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ingwrok/hotelBooking/internal/adapters/primary/web/handlers"
)

func RoomRoutes(app *fiber.App, h *handlers.RoomHandler) {

	rooms := app.Group("/api/rooms")

	rooms.Get("/:room_id/blocks", h.GetRoomBlocks)
  rooms.Get("/:room_id", h.GetRoom)
	rooms.Get("/", h.ListRooms)

	rooms.Post("/block", h.BlockRoom)
	rooms.Post("/availability/count", h.CountAvailableRooms)
	rooms.Post("/availability/find", h.FindAvailableRoom)
	rooms.Post("/:room_type_id", h.CreateRoom)

  rooms.Patch("/:room_id/status", h.ChangeRoomStatus)

  rooms.Delete("/blocks/:block_id", h.UnblockRoom)
	rooms.Delete("/:room_id", h.RemoveRoom)

}