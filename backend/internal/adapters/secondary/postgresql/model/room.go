package model

import (
	"time"

	"github.com/ingwrok/hotelBooking/internal/core/domain"
)

type Room struct {
	RoomID					int `db:"room_id"`
	RoomTypeID		int `db:"room_type_id"`
	RoomNumber		string `db:"room_number"`
	Status 			string `db:"status"`
}

func (m *Room) ToDomain() *domain.Room {
	return &domain.Room{
		RoomID:         m.RoomID,
		RoomTypeID: m.RoomTypeID,
		RoomNumber: m.RoomNumber,
		Status:     m.Status,
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


type RoomType struct {
	RoomTypeID		int `db:"room_type_id"`
	Name				string `db:"name"`
	Description	string `db:"description"`
	SizeSQM				float64 `db:"size_sqm"`
	BedType			string `db:"bed_type"`
	Capacity		int `db:"capacity"`
	PictureURL		[]string `db:"picture_url"`
	Amenities	 []string `db:"amenities"`
}


func (m *RoomType) ToDomain() *domain.RoomType {
	return &domain.RoomType{
		RoomTypeID:  m.RoomTypeID,
		Name:        m.Name,
		Description: m.Description,
		SizeSQM:     m.SizeSQM,
		BedType:     m.BedType,
		Capacity:    m.Capacity,
		PictureURL:  m.PictureURL,
		Amenities:   m.Amenities,
	}
}

// FromDomainRoomType อาจจะไม่ค่อยได้ใช้ ถ้าไม่มีฟังก์ชันสร้าง RoomType
func FromDomainRoomType(d *domain.RoomType) *RoomType {
	return &RoomType{
		RoomTypeID:  d.RoomTypeID,
		Name:        d.Name,
		Description: d.Description,
		SizeSQM:     d.SizeSQM,
		BedType:     d.BedType,
		Capacity:    d.Capacity,
		PictureURL: d.PictureURL,
		Amenities:   d.Amenities,
	}
}

type RoomBlock struct {
	RoomBlockID	int `db:"block_id"`
	RoomID			int `db:"room_id"`
	StartDate	time.Time `db:"start_date"`
	EndDate		time.Time `db:"end_date"`
	Reason			string `db:"reason"`
}

func (m *RoomBlock) ToDomain() *domain.RoomBlock {
	return &domain.RoomBlock{
		RoomBlockID: m.RoomBlockID,
		RoomID:    m.RoomID,
		StartDate: m.StartDate,
		EndDate:   m.EndDate,
		Reason:    m.Reason,
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