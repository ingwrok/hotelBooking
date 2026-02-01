package ports

import (
	"context"

	"github.com/ingwrok/hotelBooking/internal/core/domain"
)

type UserRepoPort interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id int) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	Update(ctx context.Context, id int, fields map[string]interface{}) error
	Delete(ctx context.Context, id int) error
	GetAll(ctx context.Context) ([]*domain.User, error)
}
