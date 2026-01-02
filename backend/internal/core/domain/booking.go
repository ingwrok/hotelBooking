package domain

import "time"

type Booking struct {
	BookingID     int
	UserID        int
	RatePlanID    int
	RoomID        int
	CheckInDate   time.Time
	CheckOutDate  time.Time
	NumAdults     int
	Status        string
	RoomSubTotal  float64
	AddonSubTotal float64
	TaxesAmount   float64
	TotalPrice    float64
	CreatedAt     time.Time
  UpdatedAt     time.Time
  ExpiredAt     time.Time
	BookingAddon []*BookingAddon
}

type BookingDetail struct {
	BookingID     int
	UserID        int
	RatePlanID    int
	RoomID        int
	CheckInDate   time.Time
	CheckOutDate  time.Time
	NumAdults     int
	Status        string
	RoomSubTotal  float64
	AddonSubTotal float64
	TaxesAmount   float64
	TotalPrice    float64
	CreatedAt     time.Time
  UpdatedAt     time.Time
  ExpiredAt     time.Time
	BookingAddon []*BookingAddon
	RatePlanName string
	RoomNumber   string
	RoomTypeName string
}


type BookingAddon struct {
	BookingAddonID int
	BookingID      int
	AddonID        int
	AddonName      string
	Quantity       int
	PriceAtBooking float64
}