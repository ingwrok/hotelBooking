package middleware

import (
	"context"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ingwrok/hotelBooking/internal/common/logger"
	"github.com/ingwrok/hotelBooking/internal/core/domain"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type userGetter interface {
	GetUser(ctx context.Context, id int) (*domain.User, error)
}
type bookingGetter interface {
	GetFullDetails(ctx context.Context, bookingID int) (*domain.BookingDetail, error)
}

type AuthUser struct {
	ID      int
	IsAdmin bool
}

type MyCustomClaims struct {
	UserID  int  `json:"uid"`
	IsAdmin bool `json:"is_admin"`
	jwt.RegisteredClaims
}

func AuthMiddleware(ug userGetter) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var tokenString string
		authHeader := c.Get("Authorization")

		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = authHeader[7:]
		} else {
			tokenString = c.Cookies("Authorization")
		}

		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
		}

		secret := []byte(viper.GetString("secret"))
		claims := &MyCustomClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return secret, nil
		})

		if err != nil || !token.Valid {
			logger.ErrorErr(err, "ParseWithClaims failed in AuthMiddleware")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired token"})
		}

		u, err := ug.GetUser(c.UserContext(), claims.UserID)
		if err != nil {
			logger.ErrorErr(err, "GetUser failed in AuthMiddleware", zap.Int("UserID", claims.UserID))
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user not found or inactive"})
		}

		c.Locals("authUser", &AuthUser{
			ID:      int(u.UserID),
			IsAdmin: u.IsAdmin,
		})

		// DEBUG
		fmt.Printf(">>> GLOBAL AUTH DEBUG: UserID=%d IsAdmin=%v <<<\n", u.UserID, u.IsAdmin)

		return c.Next()
	}
}

func GetAuthUser(c *fiber.Ctx) *AuthUser {
	if v := c.Locals("authUser"); v != nil {
		return v.(*AuthUser)
	}
	return nil
}

func VerifyUser(paramName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		au := GetAuthUser(c)
		if au == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		if au.IsAdmin {
			return c.Next()
		}

		targetID, err := c.ParamsInt(paramName)
		if err != nil || targetID <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id parameter"})
		}

		if au.ID != targetID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "permission denied: you can only access your own data",
			})
		}

		return c.Next()
	}
}

func VerifyAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		au := GetAuthUser(c)
		if au == nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		if !au.IsAdmin {
			fmt.Printf(">>> VERIFY ADMIN FAILED: UserID=%d IsAdmin=%v <<<\n", au.ID, au.IsAdmin)
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "admin privilege required"})
		}
		return c.Next()
	}
}

func VerifyBookingOwner(bookingSvc bookingGetter) fiber.Handler {
	return func(c *fiber.Ctx) error {
		au := GetAuthUser(c)
		if au == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		if au.IsAdmin {
			return c.Next()
		}

		bookingID, err := c.ParamsInt("booking_id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid booking id"})
		}

		booking, err := bookingSvc.GetFullDetails(c.UserContext(), bookingID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "booking not found"})
		}

		if booking.UserID != au.ID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "you don't have permission to this booking"})
		}

		return c.Next()
	}
}
