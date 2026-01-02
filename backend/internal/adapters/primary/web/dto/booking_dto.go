package dto

import (
	"time"

	"github.com/ingwrok/hotelBooking/internal/core/domain"
)

type BookingRequest struct {
	UserID       int                   `json:"user_id"`
	RatePlanID   int                   `json:"rate_plan_id"`
	RoomTypeID   int                   `json:"room_type_id"`
	CheckInDate  string                `json:"check_in_date"`
	CheckOutDate string                `json:"check_out_date"`
	NumAdults    int                   `json:"num_adults"`
	BookingAddon []BookingAddonRequest `json:"booking_addon"`
}

type BookingResponse struct {
	BookingID     int                    `json:"booking_id"`
	UserID        int                    `json:"user_id"`
	RatePlanID    int                    `json:"rate_plan_id"`
	RoomID        int                    `json:"room_id"`
	CheckInDate   time.Time              `json:"check_in_date"`
	CheckOutDate  time.Time              `json:"check_out_date"`
	NumAdults     int                    `json:"num_adults"`
	Status        string                 `json:"status"`
	RoomSubTotal  float64                `json:"room_sub_total"`
	AddonSubTotal float64                `json:"addon_sub_total"`
	TaxesAmount   float64                `json:"taxes_amount"`
	TotalPrice    float64                `json:"total_price"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
	ExpiredAt     time.Time              `json:"expired_at"`
	RatePlanName  string                 `json:"rate_plan_name"`
	RoomNumber    string                 `json:"room_number"`
	RoomTypeName  string                 `json:"room_type_name"`
	BookingAddon  []BookingAddonResponse `json:"booking_addon"`
}

func ToBookingResponse(b *domain.BookingDetail) *BookingResponse {
	return &BookingResponse{
		BookingID:     b.BookingID,
		UserID:        b.UserID,
		RatePlanID:    b.RatePlanID,
		RoomID:        b.RoomID,
		CheckInDate:   b.CheckInDate,
		CheckOutDate:  b.CheckOutDate,
		NumAdults:     b.NumAdults,
		Status:        b.Status,
		RoomSubTotal:  b.RoomSubTotal,
		AddonSubTotal: b.AddonSubTotal,
		TaxesAmount:   b.TaxesAmount,
		TotalPrice:    b.TotalPrice,
		CreatedAt:     b.CreatedAt,
		UpdatedAt:     b.UpdatedAt,
		ExpiredAt:     b.ExpiredAt,
		RatePlanName:  b.RatePlanName,
		RoomNumber:    b.RoomNumber,
		RoomTypeName:  b.RoomTypeName,
		BookingAddon:  (ToBookingAddonResponses(b.BookingAddon)),
	}
}

func ToDomainBookingAddons(addons []BookingAddonRequest) []*domain.BookingAddon {
	if addons == nil {
		return []*domain.BookingAddon{}
	}

	res := make([]*domain.BookingAddon, len(addons))
	for i, a := range addons {
		res[i] = &domain.BookingAddon{
			AddonID:  a.AddonID,
			Quantity: a.Quantity,
		}
	}
	return res
}

type BookingAddonRequest struct {
	AddonID  int `json:"addon_id"`
	Quantity int `json:"quantity"`
}

type BookingAddonResponse struct {
	BookingAddonID int     `json:"booking_addon_id"`
	BookingID      int     `json:"booking_id"`
	AddonID        int     `json:"addon_id"`
	AddonName      string  `json:"addon_name"`
	Quantity       int     `json:"quantity"`
	PriceAtBooking float64 `json:"price_at_booking"`
}

func ToBookingAddonResponses(addons []*domain.BookingAddon) []BookingAddonResponse {
	if addons == nil {
		return []BookingAddonResponse{}
	}

	res := make([]BookingAddonResponse, len(addons))
	for i, a := range addons {
		res[i] = BookingAddonResponse{
			BookingAddonID: a.BookingAddonID,
			BookingID:      a.BookingID,
			AddonID:        a.AddonID,
			AddonName:      a.AddonName,
			Quantity:       a.Quantity,
			PriceAtBooking: a.PriceAtBooking,
		}
	}
	return res
}
