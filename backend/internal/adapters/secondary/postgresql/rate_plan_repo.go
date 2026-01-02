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

type RatePlanRepository struct {
	db *sqlx.DB
}

func NewRatePlanRepository(db *sqlx.DB) ports.RatePlanRepository {
	return &RatePlanRepository{db: db}
}

func (r *RatePlanRepository) CreateRatePlan(ctx context.Context, rp *domain.RatePlan) error {
	m := model.FromDomainRatePlan(rp)

	q := `INSERT INTO rate_plans (name, description, is_special_package, allow_free_cancel, allow_pay_later)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING rate_plan_id`

	var newID int
	err := r.db.QueryRowContext(ctx, q, m.Name, m.Description, m.IsSpecialPackage, m.AllowFreeCancel, m.AllowPayLater).Scan(&newID)
	if err != nil {
		return err
	}

	rp.RatePlanID = newID
	return nil
}

func (r *RatePlanRepository) UpdateRatePlan(ctx context.Context, rp *domain.RatePlan) error {
	m := model.FromDomainRatePlan(rp)

	q := `UPDATE rate_plans
				SET name = $1,
					description = $2,
					is_special_package = $3,
					allow_free_cancel = $4,
					allow_pay_later = $5
				WHERE rate_plan_id = $6`

	result, err := r.db.ExecContext(ctx, q, m.Name, m.Description, m.IsSpecialPackage, m.AllowFreeCancel, m.AllowPayLater, m.RatePlanID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("no rate plan found with id %d: %w", m.RatePlanID, errs.ErrNotFound)
	}

	return err
}

func (r *RatePlanRepository) DeleteRatePlan(ctx context.Context, ratePlanID int) error {
	q := `DELETE FROM rate_plans WHERE rate_plan_id = $1`

	result, err := r.db.ExecContext(ctx, q, ratePlanID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("no rate plan found with id %d: %w", ratePlanID, errs.ErrNotFound)
	}

	return err
}

func (r *RatePlanRepository) GetRatePlanByID(ctx context.Context, ratePlanID int) (*domain.RatePlan, error) {
	q := `SELECT rate_plan_id, name, description, is_special_package, allow_free_cancel, allow_pay_later, created_at, updated_at
				FROM rate_plans
				WHERE rate_plan_id = $1`

	var m model.RatePlan
	err := r.db.GetContext(ctx, &m, q, ratePlanID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("rate plan id %d: %w", ratePlanID, errs.ErrNotFound)
		}
		return nil, err
	}

	return m.ToDomain(), nil
}

func (r *RatePlanRepository) GetAllRatePlans(ctx context.Context) ([]*domain.RatePlan, error) {
	q := `SELECT rate_plan_id, name, description, is_special_package, allow_free_cancel, allow_pay_later, created_at, updated_at
				FROM rate_plans`

	var ms []model.RatePlan
	err := r.db.SelectContext(ctx, &ms, q)
	if err != nil {
		return nil, err
	}

	rps := make([]*domain.RatePlan, len(ms))
	for i, m := range ms {
		rps[i] = m.ToDomain()
	}

	return rps, nil
}

func (r *RatePlanRepository) SetRoomTypePrice(ctx context.Context, roomTypeID, ratePlanID int, price float64) error {
	q := `INSERT INTO room_type_rate_prices (room_type_id, rate_plan_id, price)
        VALUES ($1, $2, $3)
        ON CONFLICT (room_type_id, rate_plan_id)
        DO UPDATE SET price = EXCLUDED.price;
				`
	_, err := r.db.ExecContext(ctx, q, roomTypeID, ratePlanID, price)
	if err != nil {
		return err
	}

	return err
}

func (r *RatePlanRepository) GetPriceByRoomType(ctx context.Context, roomTypeID, ratePlanID int) (float64, error) {
	q := `SELECT price
				FROM room_type_rate_prices
				WHERE room_type_id = $1 AND rate_plan_id = $2`

	var price float64
	err := r.db.GetContext(ctx, &price, q, roomTypeID, ratePlanID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("room type id %d: %w", roomTypeID, errs.ErrNotFound)
		}
		return 0, err
	}

	return price, nil
}

func (r *RatePlanRepository) DeleteRoomTypePrice(ctx context.Context, roomTypeID, ratePlanID int) error {
	q := `DELETE FROM room_type_rate_prices WHERE room_type_id = $1 AND rate_plan_id = $2`

	result, err := r.db.ExecContext(ctx, q, roomTypeID, ratePlanID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("no price found for room type %d and rate plan %d: %w", roomTypeID, ratePlanID, errs.ErrNotFound)
	}

	return err
}

func (r *RatePlanRepository) GetAllRatePlansByRoomTypeID(ctx context.Context, roomTypeID int) ([]*domain.RatePlanFull, error) {
	q := `SELECT rp.rate_plan_id, rp.name, rp.description, rp.is_special_package, rp.allow_free_cancel, rp.allow_pay_later, rtrp.price, rp.created_at, rp.updated_at
				FROM rate_plans rp
				JOIN room_type_rate_prices rtrp ON rp.rate_plan_id = rtrp.rate_plan_id
				WHERE rtrp.room_type_id = $1
				ORDER BY rtrp.price ASC
				`

	var ms []model.RatePlanFull
	err := r.db.SelectContext(ctx, &ms, q, roomTypeID)
	if err != nil {
		return nil, err
	}

	rps := make([]*domain.RatePlanFull, len(ms))
	for i, m := range ms {
		rps[i] = m.ToDomain()
	}

	return rps, nil
}