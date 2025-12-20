package postgresql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ingwrok/hotelBooking/internal/adapters/secondary/postgresql/model"
	"github.com/ingwrok/hotelBooking/internal/common/errs"
	"github.com/ingwrok/hotelBooking/internal/core/domain"
	"github.com/ingwrok/hotelBooking/internal/core/ports"
	"github.com/jmoiron/sqlx"
)

type AmenityRepository struct {
	db *sqlx.DB
}

func NewAmenityRepository(db *sqlx.DB) ports.AmenityRepository{
	return &AmenityRepository{db: db}
}

func (r *AmenityRepository)CreateAmenity(ctx context.Context, amenity *domain.Amenity) error{
	m := model.FromDomainAmenity(amenity)
	q := `INSERT INTO amenities (name)
				VALUES ($1)
				RETURNING amenity_id
			`
	var newID int
	err := r.db.QueryRowContext(ctx, q, m.Name).Scan(&newID)
	if err != nil {
		return err
	}

	amenity.AmenityID = newID
	return nil
}

func (r *AmenityRepository)GetAmenityByID(ctx context.Context, id int) (*domain.Amenity, error){
	var m model.Amenity
	q := `SELECT amenity_id, name FROM amenities WHERE amenity_id = $1`

	err := r.db.GetContext(ctx, &m, q, id)
	if err != nil {
		if err == sql.ErrNoRows{
			return nil, fmt.Errorf("amenity id %d : %w",id,errs.ErrNotFound)
		}
		return nil, err
	}

	return m.ToDomain(), nil
}

func (r *AmenityRepository)GetAllAmenities(ctx context.Context) ([]*domain.Amenity, error){
	var models []model.Amenity
	q := `SELECT amenity_id, name FROM amenities`

	err := r.db.SelectContext(ctx, &models, q)
	if err != nil {
		return nil, err
	}

	amenities := make([]*domain.Amenity, len(models))
	for i, m := range models {
		amenities[i] = m.ToDomain()
	}

	return amenities, nil
}

func (r *AmenityRepository)UpdateAmenity(ctx context.Context, amenity *domain.Amenity) error{
	m := model.FromDomainAmenity(amenity)
	q := `UPDATE amenities SET name = $1 WHERE amenity_id = $2`

	result, err := r.db.ExecContext(ctx, q, m.Name, m.AmenityID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("no amenity found with id %d: %w", m.AmenityID, errs.ErrNotFound)
	}

	return nil
}

func (r *AmenityRepository)DeleteAmenity(ctx context.Context, id int) error{
	q := `DELETE FROM amenities WHERE amenity_id = $1`

	result, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("no amenity found with id %d: %w", id, errs.ErrNotFound)
	}

	return nil
}

func (r *AmenityRepository)GetAmenitiesByRoomTypeID(ctx context.Context, roomTypeID int) ([]*domain.Amenity, error){
	q := `SELECT a.amenity_id, a.name
				FROM roomtype_amenities rta
				JOIN amenities a ON rta.amenity_id = a.amenity_id
				WHERE rta.room_type_id = $1
				`

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