package domain

import "time"

type RatePlan struct {
	RatePlanID       int
	Name             string
	Description      string
	IsSpecialPackage bool
	AllowFreeCancel  bool
	AllowPayLater    bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type RatePlanFull struct {
	RatePlanID       int
	Name             string
	Description      string
	IsSpecialPackage bool
	AllowFreeCancel  bool
	AllowPayLater    bool
	Price            float64
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type RoomTypeRatePrice struct {
	RoomTypeID int
	RatePlanID int
	Price      float64
}
