package ports

import (
	"context"

	"github.com/ingwrok/hotelBooking/internal/core/domain"
)

type EmailRepository interface {
	SendBookingConfirmation(ctx context.Context, booking *domain.BookingDetail, addons []*domain.BookingAddon) error
}
