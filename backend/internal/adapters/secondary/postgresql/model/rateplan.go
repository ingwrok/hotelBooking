package model

import (
	"time"

	"github.com/ingwrok/hotelBooking/internal/core/domain"
)

type RatePlan struct {
	RatePlanID       int       `db:"rate_plan_id"`
	Name             string    `db:"name"`
	Description      string    `db:"description"`
	IsSpecialPackage bool      `db:"is_special_package"`
	AllowFreeCancel  bool      `db:"allow_free_cancel"`
	AllowPayLater    bool      `db:"allow_pay_later"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

func (m *RatePlan) ToDomain() *domain.RatePlan {
	return &domain.RatePlan{
		RatePlanID:       m.RatePlanID,
		Name:             m.Name,
		Description:      m.Description,
		IsSpecialPackage: m.IsSpecialPackage,
		AllowFreeCancel:  m.AllowFreeCancel,
		AllowPayLater:    m.AllowPayLater,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}

func FromDomainRatePlan(ratePlan *domain.RatePlan) *RatePlan {
	return &RatePlan{
		RatePlanID:       ratePlan.RatePlanID,
		Name:             ratePlan.Name,
		Description:      ratePlan.Description,
		IsSpecialPackage: ratePlan.IsSpecialPackage,
		AllowFreeCancel:  ratePlan.AllowFreeCancel,
		AllowPayLater:    ratePlan.AllowPayLater,
	}
}

type RoomTypeRatePrice struct {
	RoomTypeID int     `db:"room_type_id"`
	RatePlanID int     `db:"rate_plan_id"`
	Price      float64 `db:"price"`
}

func (m *RoomTypeRatePrice) ToDomain() *domain.RoomTypeRatePrice {
	return &domain.RoomTypeRatePrice{
		RoomTypeID: m.RoomTypeID,
		RatePlanID: m.RatePlanID,
		Price:      m.Price,
	}
}

func FromDomainRoomTypeRatePrice(roomTypeRatePrice *domain.RoomTypeRatePrice) *RoomTypeRatePrice {
	return &RoomTypeRatePrice{
		RoomTypeID: roomTypeRatePrice.RoomTypeID,
		RatePlanID: roomTypeRatePrice.RatePlanID,
		Price:      roomTypeRatePrice.Price,
	}
}

type RatePlanFull struct {
	RatePlanID       int       `db:"rate_plan_id"`
	Name             string    `db:"name"`
	Description      string    `db:"description"`
	IsSpecialPackage bool      `db:"is_special_package"`
	AllowFreeCancel  bool      `db:"allow_free_cancel"`
	AllowPayLater    bool      `db:"allow_pay_later"`
	Price            float64   `db:"price"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

func (m *RatePlanFull) ToDomain() *domain.RatePlanFull {
	return &domain.RatePlanFull{
		RatePlanID:       m.RatePlanID,
		Name:             m.Name,
		Description:      m.Description,
		IsSpecialPackage: m.IsSpecialPackage,
		AllowFreeCancel:  m.AllowFreeCancel,
		AllowPayLater:    m.AllowPayLater,
		Price:            m.Price,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}

func FromDomainRatePlanFull(ratePlanFull *domain.RatePlanFull) *RatePlanFull {
	return &RatePlanFull{
		RatePlanID:       ratePlanFull.RatePlanID,
		Name:             ratePlanFull.Name,
		Description:      ratePlanFull.Description,
		IsSpecialPackage: ratePlanFull.IsSpecialPackage,
		AllowFreeCancel:  ratePlanFull.AllowFreeCancel,
		AllowPayLater:    ratePlanFull.AllowPayLater,
		Price:            ratePlanFull.Price,
	}
}
