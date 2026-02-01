package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ingwrok/hotelBooking/internal/adapters/primary/web/handlers"
	"github.com/ingwrok/hotelBooking/internal/adapters/primary/web/middleware"
	"github.com/ingwrok/hotelBooking/internal/core/services"
)

func UserRoutes(app *fiber.App, h *handlers.UserHandler, userSvc *services.UserService) {
	auth := app.Group("/api/auth")
	auth.Post("/register", h.Register)
	auth.Post("/login", h.Login)
	auth.Post("/logout", h.Logout)

	users := app.Group("/api/users", middleware.AuthMiddleware(userSvc))

	users.Get("/:id", middleware.VerifyUser("id"), h.GetUser)
	users.Put("/:id", middleware.VerifyUser("id"), h.UpdateUser)

	users.Get("/", middleware.VerifyAdmin(), h.GetUsers)
	users.Delete("/:id", middleware.VerifyAdmin(), h.DeleteUser)
}
