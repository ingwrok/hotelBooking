package ports

import (
	"context"

	"github.com/ingwrok/hotelBooking/internal/core/domain"
)

type BookingRepository interface {
	CreateBooking(ctx context.Context, booking *domain.Booking, addons []*domain.BookingAddon) error
	GetBookingWithAddons(ctx context.Context, bookingID int) (*domain.BookingDetail, error)
	UpdateBookingStatus(ctx context.Context, bookingID int, status string) error
	SyncBookingAddons(ctx context.Context, bookingID int, addons []*domain.BookingAddon, newTotalPrice float64) error
	GetBookingAddonsByBookingID(ctx context.Context, bookingID int) ([]*domain.BookingAddon, error)
	CancelExpiredBookings(ctx context.Context) (int64, error)
	GetBookingsByUserID(ctx context.Context, userID int) ([]*domain.BookingDetail, error)
	GetAllBookings(ctx context.Context) ([]*domain.BookingDetail, error)
}
