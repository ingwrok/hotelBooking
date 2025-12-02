package model

import "github.com/ingwrok/hotelBooking/internal/core/domain"

type Amenity struct {
	AmenityID	int `db:"amenity_id"`
	Name			string `db:"name"`
}

func (m *Amenity) ToDomain() *domain.Amenity {
	return &domain.Amenity{
		AmenityID: m.AmenityID,
		Name:      m.Name,
	}
}

func FromDomainAmenity(d *domain.Amenity) *Amenity {
	return &Amenity{
		AmenityID: d.AmenityID,
		Name:      d.Name,
	}
}