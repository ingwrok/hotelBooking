package ports

import (
	"context"

	"github.com/ingwrok/hotelBooking/internal/core/domain"
)

type AmenityRepository interface {
	CreateAmenity(ctx context.Context, amenity *domain.Amenity) error
	GetAmenityByID(ctx context.Context, id int) (*domain.Amenity, error)
	GetAllAmenities(ctx context.Context) ([]*domain.Amenity, error)
	UpdateAmenity(ctx context.Context, amenity *domain.Amenity) error
	DeleteAmenity(ctx context.Context, id int) error
}
