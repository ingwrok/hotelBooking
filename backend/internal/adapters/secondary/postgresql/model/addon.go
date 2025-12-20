package model

import "github.com/ingwrok/hotelBooking/internal/core/domain"

type AddonCategory struct {
	CategoryID int    `db:"category_id"`
	Name       string `db:"name"`
}

func (m *AddonCategory) ToDomain() *domain.AddonCategory {
	return &domain.AddonCategory{
		CategoryID: m.CategoryID,
		Name:       m.Name,
	}
}

func FromDomainAddonCategory(d *domain.AddonCategory) *AddonCategory{
	return &AddonCategory{
		CategoryID: d.CategoryID,
		Name:       d.Name,
	}
}

type Addon struct {
	AddonID   int    `db:"addon_id"`
	CategoryID int    `db:"category_id"`
	Name      string `db:"name"`
	Description string `db:"description"`
	Price     float64 `db:"price"`
	UnitName	string `db:"unit_name"`
}

func (m *Addon) ToDomain() *domain.Addon {
	return &domain.Addon{
		AddonID:     m.AddonID,
		CategoryID:  m.CategoryID,
		Name:        m.Name,
		Description: m.Description,
		Price:       m.Price,
		UnitName:    m.UnitName,
	}
}

func FromDomainAddon(d *domain.Addon) *Addon{
	return &Addon{
		AddonID:     d.AddonID,
		CategoryID:  d.CategoryID,
		Name:        d.Name,
		Description: d.Description,
		Price:       d.Price,
		UnitName:    d.UnitName,
	}
}