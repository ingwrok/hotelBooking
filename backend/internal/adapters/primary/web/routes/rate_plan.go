package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ingwrok/hotelBooking/internal/adapters/primary/web/handlers"
	"github.com/ingwrok/hotelBooking/internal/adapters/primary/web/middleware"
	"github.com/ingwrok/hotelBooking/internal/core/services"
)

func RatePlanRoutes(app *fiber.App, h *handlers.RatePlanHandler, userSvc *services.UserService) {
	ratePlans := app.Group("/api/rate_plans")

  ratePlans.Get("/", h.ListRatePlans)
  ratePlans.Get("/:rate_plan_id", h.GetRatePlan)
  ratePlans.Get("/:rate_plan_id/room-types/:room_type_id", h.GetPrice)
  ratePlans.Get("/room-types/:room_type_id", h.ListRatePlansByRoomType)

  admin := ratePlans.Group("/", middleware.AuthMiddleware(userSvc), middleware.VerifyAdmin())
  admin.Post("/", h.CreateRatePlan)
  admin.Put("/:rate_plan_id", h.UpdateRatePlan)
  admin.Delete("/:rate_plan_id", h.RemoveRatePlan)
  admin.Put("/:rate_plan_id/room-types/:room_type_id", h.UpdateRoomTypePrice)
  admin.Delete("/:rate_plan_id/room-types/:room_type_id", h.RemoveRoomTypePrice)
}
