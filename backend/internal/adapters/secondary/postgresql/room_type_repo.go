package postgresql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ingwrok/hotelBooking/internal/adapters/secondary/postgresql/model"
	"github.com/ingwrok/hotelBooking/internal/core/domain"
	"github.com/ingwrok/hotelBooking/internal/core/ports"
	"github.com/jmoiron/sqlx"
)

type RoomTypeRepository struct {
	db *sqlx.DB
}

func NewRoomTypeRepository(db *sqlx.DB) ports.RoomTypeRepository{
	return &RoomTypeRepository{db: db}
}

func (r *RoomTypeRepository) CreateRoomType(ctx context.Context, rt *domain.RoomType) error {
  m := model.FromDomainRoomType(rt)

	// 3. Insert ลงตารางหลัก (roomtypes)
	q := `
        INSERT INTO roomtypes (name, description, size_sqm, bed_type, capacity, picture_url)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING room_type_id
    `
	var newID int
	err := r.db.QueryRowContext(ctx, q, m.Name, m.Description, m.SizeSQM, m.BedType, m.Capacity, m.PictureURL,
	).Scan(&newID)
	if err != nil {
		return err
	}

	rt.RoomTypeID = newID
	return nil
}
func (r *RoomTypeRepository) UpdateRoomType(ctx context.Context, rt *domain.RoomType) error {
	m := model.FromDomainRoomType(rt)
  q := `
      UPDATE roomtypes
			SET name=$1, description=$2, size_sqm=$3, bed_type=$4, capacity=$5, picture_url=$6
      WHERE room_type_id=$7
    `
  result, err := r.db.ExecContext(ctx, q,
    m.Name, m.Description, m.SizeSQM, m.BedType, m.Capacity, m.PictureURL,
		m.RoomTypeID,
  )

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("no room type found with id %d", m.RoomTypeID)
	}

  return nil
}
func (r *RoomTypeRepository)DeleteRoomType(ctx context.Context,id int) error{
  q := `DELETE FROM roomtypes
  WHERE room_type_id = $1`

  result, err := r.db.ExecContext(ctx, q, id)
  if err != nil {
    return err
  }

  rows, err := result.RowsAffected()
  if err != nil {
    return err
  }

  if rows == 0 {
    return fmt.Errorf("no room found with id %d", id)
  }

  return nil
}

// read Room
func (r *RoomTypeRepository)GetRoomTypeByID(ctx context.Context, id int) (*domain.RoomType, error){
  q := `SELECT
        	rt.room_type_id,
          rt.name,
          rt.description,
					rt.size_sqm,
          rt.bed_type,
          rt.capacity,
          rt.picture_url AS picture_url,
        COALESCE(
          ARRAY_AGG(a.name) FILTER (WHERE a.name IS NOT NULL),
            '{}'
          ) AS amenities
        FROM roomtypes rt
        LEFT JOIN roomtype_amenities rta ON rt.room_type_id = rta.room_type_id
        LEFT JOIN amenities a ON rta.amenity_id = a.amenity_id
        WHERE rt.room_type_id = $1
        GROUP BY rt.room_type_id`

  var model model.RoomType

    err := r.db.GetContext(ctx, &model, q, id)
    if err != nil {
      if err == sql.ErrNoRows {
        return nil, fmt.Errorf("room type not found")
      }
      return nil, err
    }

    return model.ToDomain(), nil
}

func (r *RoomTypeRepository)GetAllRoomTypes(ctx context.Context) ([]*domain.RoomType, error){
	 q := `SELECT
        	rt.room_type_id,
          rt.name,
          rt.description,
					rt.size_sqm,
          rt.bed_type,
          rt.capacity,
          rt.picture_url AS picture_url,
        COALESCE(
          ARRAY_AGG(a.name) FILTER (WHERE a.name IS NOT NULL),
            '{}'
          ) AS amenities
        FROM roomtypes rt
        LEFT JOIN roomtype_amenities rta ON rt.room_type_id = rta.room_type_id
        LEFT JOIN amenities a ON rta.amenity_id = a.amenity_id
        GROUP BY rt.room_type_id`

  var models []model.RoomType

    err := r.db.SelectContext(ctx, &models, q)
    if err != nil {
      return nil, err
    }

		result := make([]*domain.RoomType, len(models))
		for i, m := range models {
			result[i] = m.ToDomain()
		}

    return result, nil
}

func (r *RoomTypeRepository)UpdateRoomTypeAmenities(ctx context.Context, roomTypeID int, amenityIDs []int) error{
	tx,err := r.db.BeginTxx(ctx,nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// delete existing amenities
	qDelete := `DELETE FROM roomtype_amenities WHERE room_type_id = $1`
	_, err = tx.ExecContext(ctx, qDelete, roomTypeID)
	if err != nil {
		return err
	}

	// insert new amenities
	if len(amenityIDs) > 0 {
		qInsert := `INSERT INTO roomtype_amenities (room_type_id, amenity_id) VALUES ($1, $2)`
		for _, amenityID := range amenityIDs {
			_, err = tx.ExecContext(ctx, qInsert, roomTypeID, amenityID)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}
func (r *RoomTypeRepository)GetAmenitiesByRoomTypeID(ctx context.Context, roomTypeID int) ([]*domain.Amenity, error){
	q := `SELECT a.amenity_id, a.name
				FROM roomtype_amenities rta
				JOIN amenities a ON rta.amenity_id = a.amenity_id
				WHERE rta.room_type_id = $1`

	var models []model.Amenity

	err := r.db.SelectContext(ctx, &models, q, roomTypeID)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Amenity, len(models))
	for i, m := range models {
		result[i] = m.ToDomain()
	}

	return result, nil
}