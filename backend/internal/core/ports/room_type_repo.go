package ports

import (
    "context"
    "github.com/ingwrok/hotelBooking/internal/core/domain"
)

type RoomTypeRepository interface {
    // --- Room Type (CRUD ประเภทห้อง) ---
  CreateRoomType(ctx context.Context, rt *domain.RoomType) error
	// แก้ Signature ให้รับทั้งก้อน (เพื่อให้แก้ไขชื่อ, ราคา, รูปภาพได้)
	UpdateRoomType(ctx context.Context, rt *domain.RoomType) error 
	DeleteRoomType(ctx context.Context, id int) error
	GetRoomTypeByID(ctx context.Context, id int) (*domain.RoomType, error)
	GetAllRoomTypes(ctx context.Context) ([]*domain.RoomType, error)

	// --- Amenities Management ---
	// (ฟังก์ชันเหล่านี้อาจจะไม่ได้ใช้ถ้าเราจัดการ Amenities ใน Create/Update เลย)
	// แต่เก็บไว้ก็ไม่เสียหายครับ
	UpdateRoomTypeAmenities(ctx context.Context, roomTypeID int, amenityIDs []int) error
	GetAmenitiesByRoomTypeID(ctx context.Context, roomTypeID int) ([]*domain.Amenity, error)
}