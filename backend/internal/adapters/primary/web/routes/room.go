package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ingwrok/hotelBooking/internal/adapters/primary/web/handlers"
	"github.com/ingwrok/hotelBooking/internal/adapters/primary/web/middleware"
	"github.com/ingwrok/hotelBooking/internal/core/services"
)

func RoomRoutes(app *fiber.App, h *handlers.RoomHandler, userSvc *services.UserService) {
	rooms := app.Group("/api/rooms")

	rooms.Get("/:room_id", h.GetRoom)
	rooms.Get("/", h.ListRooms)
	rooms.Post("/availability/count", h.CountAvailableRooms)
	rooms.Post("/availability/find", h.FindAvailableRoom)

	admin := rooms.Group("/", middleware.AuthMiddleware(userSvc), middleware.VerifyAdmin())
	admin.Get("/:room_id/blocks", h.GetRoomBlocks)
	admin.Post("/block", h.BlockRoom)
	admin.Post("/:room_type_id", h.CreateRoom)
	admin.Patch("/:room_id/status", h.ChangeRoomStatus)
	admin.Delete("/blocks/:block_id", h.UnblockRoom)
	admin.Delete("/:room_id", h.RemoveRoom)
}
