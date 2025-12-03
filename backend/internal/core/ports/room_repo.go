package ports

import (
    "context"
    "time"
    "github.com/ingwrok/hotelBooking/internal/core/domain"
)

type RoomRepository interface {
    // write Room
    CreateRoom(ctx context.Context, room *domain.Room) error
    DeleteRoom(ctx context.Context, id int) error
    UpdateRoomStatus(ctx context.Context, roomID int, status string) error

		// read Room
		GetRoomByID(ctx context.Context, id int) (*domain.Room, error) // คืนค่า Room นะ ไม่ใช่ RoomType
    GetAllRooms(ctx context.Context) ([]*domain.Room, error)

    // Room Block
    CheckIfBlockOverlaps(ctx context.Context,roomID int,startDate,endDate time.Time,) (int, error)
    CreateRoomBlock(ctx context.Context, block *domain.RoomBlock) error
    GetRoomBlocksByRoomID(ctx context.Context, roomID int) ([]*domain.RoomBlock, error)
    DeleteRoomBlock(ctx context.Context, blockID int) error

    // Room Availability
    // หาจำนวนห้องว่างของแต่ละ Type ในช่วงเวลาที่กำหนด
    GetAvailableRoomCounts(ctx context.Context, checkIn, checkOut time.Time) (map[int]int, error)
    // สุ่มหยิบห้องว่าง 1 ห้องจาก Type ที่ระบุ
    GetAnyAvailableRoomID(ctx context.Context, roomTypeID int, checkIn, checkOut time.Time) (int, error)
}