package dto

import "time"

type RatePlanRequest struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	IsSpecialPackage bool   `json:"isSpecialPackage"`
	AllowFreeCancel  bool   `json:"allowFreeCancel"`
	AllowPayLater    bool   `json:"allowPayLater"`
}

type RatePlanResponse struct {
	RatePlanID       int       `json:"ratePlanId"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	IsSpecialPackage bool      `json:"isSpecialPackage"`
	AllowFreeCancel  bool      `json:"allowFreeCancel"`
	AllowPayLater    bool      `json:"allowPayLater"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type RatePlanFullResponse struct {
	RatePlanID       int       `json:"ratePlanId"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	IsSpecialPackage bool      `json:"isSpecialPackage"`
	AllowFreeCancel  bool      `json:"allowFreeCancel"`
	AllowPayLater    bool      `json:"allowPayLater"`
	Price            float64   `json:"price"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type RoomTypeRatePrice struct {
	RoomTypeID int     `json:"roomTypeId"`
	RatePlanID int     `json:"ratePlanId"`
	Price      float64 `json:"price"`
}
