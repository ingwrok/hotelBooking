package ports

import (
	"context"

	"github.com/ingwrok/hotelBooking/internal/core/domain"
)

type RoomTypeRepository interface {
	// --- Room Type (CRUD ประเภทห้อง) ---
	CreateRoomType(ctx context.Context, rt *domain.RoomType) error
	UpdateRoomType(ctx context.Context, rt *domain.RoomType) error
	DeleteRoomType(ctx context.Context, id int) error
	GetRoomTypeByID(ctx context.Context, id int) (*domain.RoomType, error)
	GetAllRoomTypes(ctx context.Context) ([]*domain.RoomType, error)

	GetRoomTypeFullDetail(ctx context.Context, id int) (*domain.RoomTypeDetails, error)
}
