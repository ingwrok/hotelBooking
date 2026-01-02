package model

import (
	"time"

	"github.com/ingwrok/hotelBooking/internal/core/domain"
)

type Booking struct {
	BookingID     int       `db:"booking_id"`
	UserID        int       `db:"user_id"`
	RatePlanID    int       `db:"rate_plan_id"`
	RoomID        int       `db:"room_id"`
	CheckInDate   time.Time `db:"check_in_date"`
	CheckOutDate  time.Time `db:"check_out_date"`
	NumAdults     int       `db:"num_adults"`
	Status        string    `db:"status"`
	RoomSubTotal  float64   `db:"room_subtotal"`
	AddonSubTotal float64   `db:"addon_subtotal"`
	TaxesAmount   float64   `db:"taxes_amount"`
	TotalPrice    float64   `db:"total_price"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
	ExpiredAt     time.Time `db:"expired_at"`
}

func (m *Booking) ToDomain(addons []*BookingAddon) *domain.Booking {

	var domainAddons []*domain.BookingAddon
	for _, a := range addons{
		domainAddons = append(domainAddons, &domain.BookingAddon{
			BookingAddonID: a.BookingAddonID,
			BookingID:      a.BookingID,
			AddonID:        a.AddonID,
			Quantity:       a.Quantity,
			PriceAtBooking: a.PriceAtBooking,
		})
	}

	return &domain.Booking{
		BookingID:     m.BookingID,
		UserID:        m.UserID,
		RatePlanID:    m.RatePlanID,
		RoomID:        m.RoomID,
		CheckInDate:   m.CheckInDate,
		CheckOutDate:  m.CheckOutDate,
		NumAdults:     m.NumAdults,
		Status:        m.Status,
		RoomSubTotal:  m.RoomSubTotal,
		AddonSubTotal: m.AddonSubTotal,
		TaxesAmount:   m.TaxesAmount,
		TotalPrice:    m.TotalPrice,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
		ExpiredAt:     m.ExpiredAt,
		BookingAddon: domainAddons,
	}
}

func FromDomainBooking(booking *domain.Booking) *Booking {
	return &Booking{
		BookingID:     booking.BookingID,
		UserID:        booking.UserID,
		RatePlanID:    booking.RatePlanID,
		RoomID:        booking.RoomID,
		CheckInDate:   booking.CheckInDate,
		CheckOutDate:  booking.CheckOutDate,
		NumAdults:     booking.NumAdults,
		Status:        booking.Status,
		RoomSubTotal:  booking.RoomSubTotal,
		AddonSubTotal: booking.AddonSubTotal,
		TaxesAmount:   booking.TaxesAmount,
		TotalPrice:    booking.TotalPrice,
		CreatedAt:     booking.CreatedAt,
		UpdatedAt:     booking.UpdatedAt,
		ExpiredAt:     booking.ExpiredAt,
	}
}

type BookingAddon struct {
	BookingAddonID int     `db:"booking_addon_id"`
	BookingID      int     `db:"booking_id"`
	AddonID        int     `db:"addon_id"`
	AddonName      string  `db:"addon_name"`
	Quantity       int     `db:"quantity"`
	PriceAtBooking float64 `db:"price_at_time_of_booking"`
}

func (m *BookingAddon) ToDomain() *domain.BookingAddon {
	return &domain.BookingAddon{
		BookingAddonID: m.BookingAddonID,
		BookingID:      m.BookingID,
		AddonID:        m.AddonID,
		AddonName:      m.AddonName,
		Quantity:       m.Quantity,
		PriceAtBooking: m.PriceAtBooking,
	}
}

func FromDomainBookingAddon(bookingAddon *domain.BookingAddon) *BookingAddon {
	return &BookingAddon{
		BookingAddonID: bookingAddon.BookingAddonID,
		BookingID:      bookingAddon.BookingID,
		AddonID:        bookingAddon.AddonID,
		Quantity:       bookingAddon.Quantity,
		PriceAtBooking: bookingAddon.PriceAtBooking,
	}
}

type BookingDetail struct {
	Booking
	RatePlanName string `db:"rate_plan_name"`
	RoomNumber   string `db:"room_number"`
  RoomTypeName string `db:"room_type_name"`
}

func (m *BookingDetail) ToDomainDetail(addons []*BookingAddon) *domain.BookingDetail {

	var domainAddons []*domain.BookingAddon
  for _, a := range addons {
    domainAddons = append(domainAddons, a.ToDomain())
  }

  return &domain.BookingDetail{
    BookingID:     m.BookingID,
    UserID:        m.UserID,
    RatePlanID:    m.RatePlanID,
    RoomID:        m.RoomID,
    CheckInDate:   m.CheckInDate,
    CheckOutDate:  m.CheckOutDate,
    NumAdults:     m.NumAdults,
    Status:        m.Status,
    RoomSubTotal:  m.RoomSubTotal,
    AddonSubTotal: m.AddonSubTotal,
    TaxesAmount:   m.TaxesAmount,
    TotalPrice:    m.TotalPrice,
    CreatedAt:     m.CreatedAt,
    UpdatedAt:     m.UpdatedAt,
    ExpiredAt:     m.ExpiredAt,
    BookingAddon:  domainAddons,
    RatePlanName:  m.RatePlanName,
		RoomNumber:    m.RoomNumber,
    RoomTypeName:  m.RoomTypeName,
  }
}