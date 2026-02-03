package dto

import (
	"time"

	"github.com/ingwrok/hotelBooking/internal/core/domain"
)

type BookingRequest struct {
	UserID       int                   `json:"userId"`
	RatePlanID   int                   `json:"ratePlanId"`
	RoomTypeID   int                   `json:"roomTypeId"`
	CheckInDate  string                `json:"checkInDate"`
	CheckOutDate string                `json:"checkOutDate"`
	NumAdults    int                   `json:"numAdults"`
	Email        string                `json:"email"`
	BookingAddon []BookingAddonRequest `json:"bookingAddon"`
}

type BookingResponse struct {
	BookingID     int                    `json:"bookingId"`
	UserID        int                    `json:"userId"`
	RatePlanID    int                    `json:"ratePlanId"`
	RoomID        int                    `json:"roomId"`
	CheckInDate   time.Time              `json:"checkInDate"`
	CheckOutDate  time.Time              `json:"checkOutDate"`
	NumAdults     int                    `json:"numAdults"`
	Status        string                 `json:"status"`
	RoomSubTotal  float64                `json:"roomSubTotal"`
	AddonSubTotal float64                `json:"addonSubTotal"`
	TaxesAmount   float64                `json:"taxesAmount"`
	TotalPrice    float64                `json:"totalPrice"`
	CreatedAt     time.Time              `json:"createdAt"`
	UpdatedAt     time.Time              `json:"updatedAt"`
	ExpiredAt     time.Time              `json:"expiredAt"`
	RatePlanName  string                 `json:"ratePlanName"`
	RoomNumber    string                 `json:"roomNumber"`
	RoomTypeName  string                 `json:"roomTypeName"`
	GuestDetails  *GuestInfoResponse     `json:"guestDetails"`
	BookingAddon  []BookingAddonResponse `json:"bookingAddon"`
}

type GuestInfoResponse struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

func ToBookingResponse(b *domain.BookingDetail) *BookingResponse {
	// Simple split of username for demo purposes, as we only store username string
	// In a real app, User struct would have FirstName/LastName separately.
	// We'll reuse UserName for FirstName and "-" for LastName if single word.
	firstName := b.UserName
	lastName := ""

	// If needed we could split string here, but keeping it simple is safer.
	// Frontend expects { firstName, lastName }

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
		GuestDetails: &GuestInfoResponse{
			FirstName: firstName,
			LastName:  lastName,
			Email:     b.Email,
		},
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
	AddonID  int `json:"addonId"`
	Quantity int `json:"quantity"`
}

type BookingAddonResponse struct {
	BookingAddonID int     `json:"bookingAddonId"`
	BookingID      int     `json:"bookingId"`
	AddonID        int     `json:"addonId"`
	AddonName      string  `json:"addonName"`
	Quantity       int     `json:"quantity"`
	PriceAtBooking float64 `json:"priceAtBooking"`
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
