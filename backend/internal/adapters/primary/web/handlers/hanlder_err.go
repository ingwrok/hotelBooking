package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/ingwrok/hotelBooking/internal/common/errs"
)

func handleError(c *fiber.Ctx, err error) error {
    // รองรับทั้ง value และ pointer และรองรับ wrapped errors
    var appErr errs.AppError
    if errors.As(err, &appErr) {
        return c.Status(appErr.Code).JSON(fiber.Map{"message": appErr.Message})
    }
    var appErrPtr *errs.AppError
    if errors.As(err, &appErrPtr) && appErrPtr != nil {
        return c.Status(appErrPtr.Code).JSON(fiber.Map{"message": appErrPtr.Message})
    }
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "unexpected error"})
}