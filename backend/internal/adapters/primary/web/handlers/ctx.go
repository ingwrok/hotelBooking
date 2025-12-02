package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func buildCtx(c *fiber.Ctx) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	reqID := c.Get("X-Request-ID")
	if reqID == "" {
		reqID = uuid.New().String()
	}

	ctx = context.WithValue(ctx, "requestID", reqID)
	return ctx, cancel
}
