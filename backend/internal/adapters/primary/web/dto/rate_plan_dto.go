package dto

import "time"

type RatePlanRequest struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	IsSpecialPackage bool   `json:"is_special_package"`
	AllowFreeCancel  bool   `json:"allow_free_cancel"`
	AllowPayLater    bool   `json:"allow_pay_later"`
}

type RatePlanResponse struct {
	RatePlanID       int       `json:"rate_plan_id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	IsSpecialPackage bool      `json:"is_special_package"`
	AllowFreeCancel  bool      `json:"allow_free_cancel"`
	AllowPayLater    bool      `json:"allow_pay_later"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type RatePlanFullResponse struct {
	RatePlanID       int       `json:"rate_plan_id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	IsSpecialPackage bool      `json:"is_special_package"`
	AllowFreeCancel  bool      `json:"allow_free_cancel"`
	AllowPayLater    bool      `json:"allow_pay_later"`
	Price            float64   `json:"price"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type RoomTypeRatePrice struct {
	RoomTypeID int     `json:"room_type_id"`
	RatePlanID int     `json:"rate_plan_id"`
	Price      float64 `json:"price"`
}
