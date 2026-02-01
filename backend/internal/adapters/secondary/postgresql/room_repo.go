package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ingwrok/hotelBooking/internal/adapters/secondary/postgresql/model"
	"github.com/ingwrok/hotelBooking/internal/common/errs"
	"github.com/ingwrok/hotelBooking/internal/core/domain"
	"github.com/ingwrok/hotelBooking/internal/core/ports"
	"github.com/jmoiron/sqlx"
)

type RoomRepository struct {
	db *sqlx.DB
}

func NewRoomRepository(db *sqlx.DB) ports.RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) CreateRoom(ctx context.Context, room *domain.Room) error {
	m := model.FromDomainRoom(room)
	q := `INSERT INTO rooms (room_type_id, room_number)
        VALUES ($1, $2)
        RETURNING room_id
      `
	var newID int
	err := r.db.QueryRowContext(ctx, q, m.RoomTypeID, m.RoomNumber).Scan(&newID)
	if err != nil {
		return err
	}

	room.RoomID = newID
	return nil
}

func (r *RoomRepository) DeleteRoom(ctx context.Context, id int) error {
	q := `DELETE FROM rooms WHERE room_id=$1`
	result, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("no room found with id %d: %w", id, errs.ErrNotFound)
	}

	return nil
}

func (r *RoomRepository) UpdateRoomStatus(ctx context.Context, roomID int, status string) error {
	q := `UPDATE rooms SET status=$1 WHERE room_id=$2`
	result, err := r.db.ExecContext(ctx, q, status, roomID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("no room found with id %d: %w", roomID, errs.ErrNotFound)
	}

	return nil
}

// read Room
func (r *RoomRepository) GetRoomByID(ctx context.Context, id int) (*domain.RoomDetail, error) {
	var m model.Room

	q := `SELECT r.room_id, r.room_type_id, r.room_number, r.status, rt.name as room_type_name
        FROM rooms r
        JOIN roomtypes rt ON r.room_type_id = rt.room_type_id
        WHERE r.room_id=$1`

	err := r.db.GetContext(ctx, &m, q, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("room id %d: %w", id, errs.ErrNotFound)
		}
		return nil, err
	}

	return m.ToDomain(), nil
}
func (r *RoomRepository) GetAllRooms(ctx context.Context) ([]*domain.RoomDetail, error) {
	var models []model.Room
	q := `SELECT r.room_id, r.room_type_id, r.room_number, r.status, rt.name as room_type_name
        FROM rooms r
        JOIN roomtypes rt ON r.room_type_id = rt.room_type_id
        ORDER BY r.room_number`

	err := r.db.SelectContext(ctx, &models, q)
	if err != nil {
		return nil, err
	}

	rooms := make([]*domain.RoomDetail, len(models))
	for i, m := range models {
		rooms[i] = m.ToDomain()
	}

	return rooms, nil
}

func (r *RoomRepository) CheckIfBlockOverlaps(ctx context.Context, roomID int, startDate, endDate time.Time) (int, error) {
	q := `
		SELECT COUNT(block_id)
		FROM room_blocks
		WHERE room_id = $1
		  AND end_date > $2
		  AND start_date < $3
	  `
	var count int
	err := r.db.QueryRowContext(ctx, q, roomID, startDate, endDate).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *RoomRepository) CreateRoomBlock(ctx context.Context, block *domain.RoomBlock) error {
	model := model.FromDomainRoomBlock(block)

	q := `INSERT INTO room_blocks (room_id, start_date, end_date, reason)
		  VALUES ($1, $2, $3, $4)
		  RETURNING block_id`

	var newID int
	err := r.db.QueryRowContext(ctx, q, model.RoomID, model.StartDate, model.EndDate, model.Reason).Scan(&newID)
	if err != nil {
		return err
	}

	block.RoomBlockID = newID

	return nil
}

func (r *RoomRepository) GetRoomBlocksByRoomID(ctx context.Context, roomID int) ([]*domain.RoomBlock, error) {
	var models []model.RoomBlock

	q := `SELECT * FROM room_blocks WHERE room_id = $1`
	err := r.db.SelectContext(ctx, &models, q, roomID)
	if err != nil {
		return nil, err
	}

	blocks := make([]*domain.RoomBlock, len(models))
	for i, m := range models {
		blocks[i] = m.ToDomain()
	}

	return blocks, nil
}

func (r *RoomRepository) DeleteRoomBlock(ctx context.Context, blockID int) error {
	q := `DELETE FROM room_blocks WHERE block_id = $1`

	result, err := r.db.ExecContext(ctx, q, blockID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("no room block found with id %d: %w", blockID, errs.ErrNotFound)
	}

	return err
}

func (r *RoomRepository) GetAvailableRoomCounts(ctx context.Context, checkIn, checkOut time.Time) (map[int]int, error) {
	q := `
    SELECT
      r.room_type_id,
      COUNT(r.room_id) AS available_count
    FROM rooms r
    WHERE r.status != 'maintenance'
    AND NOT EXISTS (
      SELECT 1
      FROM bookings b
      WHERE b.room_id = r.room_id
        AND b.status != 'cancelled'
        AND b.check_in_date < $2
        AND b.check_out_date > $1
    )
    AND NOT EXISTS (
      SELECT 1
      FROM room_blocks rb
      WHERE rb.room_id = r.room_id
        AND rb.start_date < $2
        AND rb.end_date > $1
    )
    GROUP BY r.room_type_id;
  `
	type result struct {
		TypeID int `db:"room_type_id"`
		Count  int `db:"available_count"`
	}

	var rows []result
	err := r.db.SelectContext(ctx, &rows, q, checkIn, checkOut)
	if err != nil {
		return nil, err
	}
	counts := make(map[int]int)
	for _, row := range rows {
		counts[row.TypeID] = row.Count
	}
	return counts, nil

}

// สุ่มหยิบห้องว่าง 1 ห้องจาก Type ที่ระบุ
func (r *RoomRepository) GetAnyAvailableRoomID(ctx context.Context, roomTypeID int, checkIn, checkOut time.Time) (int, error) {
	q := `
    SELECT r.room_id
    FROM rooms r
    WHERE r.room_type_id = $1
      AND r.status != 'maintenance'
      AND NOT EXISTS (
          SELECT 1 FROM bookings b
          WHERE b.room_id = r.room_id
            AND b.status != 'cancelled'
            AND b.check_in_date < $3
            AND b.check_out_date > $2
      )
      AND NOT EXISTS (
          SELECT 1 FROM room_blocks rb
          WHERE rb.room_id = r.room_id
            AND rb.start_date < $3
            AND rb.end_date > $2
      )
    ORDER BY r.room_number
    LIMIT 1;
  `
	var roomID int
	err := r.db.QueryRowContext(ctx, q, roomTypeID, checkIn, checkOut).Scan(&roomID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("no available room : %w", errs.ErrNotFound)
		}
		return 0, err
	}
	return roomID, nil

}
