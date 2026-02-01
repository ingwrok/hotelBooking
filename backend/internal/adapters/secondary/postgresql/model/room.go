package model

import (
	"time"

	"github.com/ingwrok/hotelBooking/internal/core/domain"
	"github.com/lib/pq"
)

type Room struct {
	RoomID       int    `db:"room_id"`
	RoomTypeID   int    `db:"room_type_id"`
	RoomTypeName string `db:"room_type_name"`
	RoomNumber   string `db:"room_number"`
	Status       string `db:"status"`
}

func (m *Room) ToDomain() *domain.RoomDetail {
	return &domain.RoomDetail{
		RoomID:       m.RoomID,
		RoomTypeID:   m.RoomTypeID,
		RoomTypeName: m.RoomTypeName,
		RoomNumber:   m.RoomNumber,
		Status:       m.Status,
	}
}

func FromDomainRoom(d *domain.Room) *Room {
	return &Room{
		RoomID:     d.RoomID,
		RoomTypeID: d.RoomTypeID,
		RoomNumber: d.RoomNumber,
		Status:     d.Status,
	}
}

type RoomTypeDetails struct {
	RoomTypeID  int            `db:"room_type_id"`
	Name        string         `db:"name"`
	Description string         `db:"description"`
	SizeSQM     float64        `db:"size_sqm"`
	BedType     string         `db:"bed_type"`
	Capacity    int            `db:"capacity"`
	PictureURL  pq.StringArray `db:"picture_url"`
	Amenities   pq.StringArray `db:"amenities"`
}

func (m *RoomTypeDetails) ToDomain() *domain.RoomTypeDetails {
	return &domain.RoomTypeDetails{
		RoomTypeID:  m.RoomTypeID,
		Name:        m.Name,
		Description: m.Description,
		SizeSQM:     m.SizeSQM,
		BedType:     m.BedType,
		Capacity:    m.Capacity,
		PictureURL:  []string(m.PictureURL),
		Amenities:   []string(m.Amenities),
	}
}

func FromDomainRoomTypeDetails(d *domain.RoomTypeDetails) *RoomTypeDetails {
	return &RoomTypeDetails{
		RoomTypeID:  d.RoomTypeID,
		Name:        d.Name,
		Description: d.Description,
		SizeSQM:     d.SizeSQM,
		BedType:     d.BedType,
		Capacity:    d.Capacity,
		PictureURL:  pq.StringArray(d.PictureURL),
		Amenities:   pq.StringArray(d.Amenities),
	}
}

type RoomType struct {
	RoomTypeID  int            `db:"room_type_id"`
	Name        string         `db:"name"`
	Description string         `db:"description"`
	SizeSQM     float64        `db:"size_sqm"`
	BedType     string         `db:"bed_type"`
	Capacity    int            `db:"capacity"`
	PictureURL  pq.StringArray `db:"picture_url"`
	AmenityIDs  []int          `db:"amenity_ids"`
	TotalRooms  int            `db:"total_rooms"`
}

func (m *RoomType) ToDomain() *domain.RoomType {
	return &domain.RoomType{
		RoomTypeID:  m.RoomTypeID,
		Name:        m.Name,
		Description: m.Description,
		SizeSQM:     m.SizeSQM,
		BedType:     m.BedType,
		Capacity:    m.Capacity,
		PictureURL:  []string(m.PictureURL),
		AmenityIDs:  m.AmenityIDs,
		TotalRooms:  m.TotalRooms,
	}
}

func FromDomainRoomType(d *domain.RoomType) *RoomType {
	return &RoomType{
		RoomTypeID:  d.RoomTypeID,
		Name:        d.Name,
		Description: d.Description,
		SizeSQM:     d.SizeSQM,
		BedType:     d.BedType,
		Capacity:    d.Capacity,
		PictureURL:  pq.StringArray(d.PictureURL),
		AmenityIDs:  d.AmenityIDs,
		TotalRooms:  d.TotalRooms,
	}
}

type RoomBlock struct {
	RoomBlockID int       `db:"block_id"`
	RoomID      int       `db:"room_id"`
	StartDate   time.Time `db:"start_date"`
	EndDate     time.Time `db:"end_date"`
	Reason      string    `db:"reason"`
}

func (m *RoomBlock) ToDomain() *domain.RoomBlock {
	return &domain.RoomBlock{
		RoomBlockID: m.RoomBlockID,
		RoomID:      m.RoomID,
		StartDate:   m.StartDate,
		EndDate:     m.EndDate,
		Reason:      m.Reason,
	}
}

func FromDomainRoomBlock(d *domain.RoomBlock) *RoomBlock {
	return &RoomBlock{
		RoomBlockID: d.RoomBlockID,
		RoomID:      d.RoomID,
		StartDate:   d.StartDate,
		EndDate:     d.EndDate,
		Reason:      d.Reason,
	}
}
