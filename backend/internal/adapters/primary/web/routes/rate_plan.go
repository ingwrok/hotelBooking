package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ingwrok/hotelBooking/internal/adapters/primary/web/handlers"
)

func RatePlanRoutes(app *fiber.App, h *handlers.RatePlanHandler) {

	ratePlans := app.Group("/api/rate_plans")

	ratePlans.Post("/", h.CreateRatePlan)
	ratePlans.Get("/", h.ListRatePlans)
	ratePlans.Get("/:rate_plan_id", h.GetRatePlan)
	ratePlans.Put("/:rate_plan_id", h.UpdateRatePlan)
	ratePlans.Delete("/:rate_plan_id", h.RemoveRatePlan)

	ratePlans.Put("/:rate_plan_id/room-types/:room_type_id", h.UpdateRoomTypePrice)
	ratePlans.Get("/:rate_plan_id/room-types/:room_type_id", h.GetPrice)
	ratePlans.Get("/room-types/:room_type_id",h.ListRatePlansByRoomType)
	ratePlans.Delete("/:rate_plan_id/room-types/:room_type_id", h.RemoveRoomTypePrice)
}
