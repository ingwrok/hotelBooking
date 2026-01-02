package ports

import (
	"context"

	"github.com/ingwrok/hotelBooking/internal/core/domain"
)

type RatePlanRepository interface {
	CreateRatePlan(ctx context.Context, rp *domain.RatePlan) error
	UpdateRatePlan(ctx context.Context, rp *domain.RatePlan) error
	DeleteRatePlan(ctx context.Context, ratePlanID int) error
	GetRatePlanByID(ctx context.Context, ratePlanID int) (*domain.RatePlan, error)
	GetAllRatePlans(ctx context.Context) ([]*domain.RatePlan, error)

	SetRoomTypePrice(ctx context.Context, roomTypeID, ratePlanID int, price float64) error
	GetPriceByRoomType(ctx context.Context, roomTypeID, ratePlanID int) (float64, error)
	DeleteRoomTypePrice(ctx context.Context, roomTypeID, ratePlanID int) error
	GetAllRatePlansByRoomTypeID(ctx context.Context, roomTypeID int) ([]*domain.RatePlanFull, error)
}
