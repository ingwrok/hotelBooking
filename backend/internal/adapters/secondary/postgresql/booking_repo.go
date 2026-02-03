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

type BookingRepository struct {
	db *sqlx.DB
}

func NewBookingRepository(db *sqlx.DB) ports.BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) CreateBooking(ctx context.Context, booking *domain.Booking, baddons []*domain.BookingAddon) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	mb := model.FromDomainBooking(booking)
	queryBooking := `
		INSERT INTO bookings (
			user_id, rate_plan_id, room_id, check_in_date, check_out_date,
			num_adults, room_subtotal, addon_subtotal,
			taxes_amount, total_price, expired_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING booking_id`

	var bookingID int
	err = tx.QueryRowContext(ctx, queryBooking,
		mb.UserID, mb.RatePlanID, mb.RoomID, mb.CheckInDate, mb.CheckOutDate,
		mb.NumAdults, mb.RoomSubTotal, mb.AddonSubTotal,
		mb.TaxesAmount, mb.TotalPrice, mb.ExpiredAt,
	).Scan(&bookingID)

	if err != nil {
		return err
	}

	if len(baddons) > 0 {
		queryAddon := `
			INSERT INTO booking_addons (booking_id, addon_id, quantity, price_at_time_of_booking)
			VALUES ($1, $2, $3, $4)`

		for _, daddon := range baddons {
			mAddon := model.FromDomainBookingAddon(daddon)
			_, err := tx.ExecContext(ctx, queryAddon, bookingID, mAddon.AddonID, mAddon.Quantity, mAddon.PriceAtBooking)
			if err != nil {
				return err
			}
		}
	}

	booking.BookingID = bookingID

	return tx.Commit()
}

func (r *BookingRepository) GetBookingWithAddons(ctx context.Context, bookingID int) (*domain.BookingDetail, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var mBookingDetail model.BookingDetail
	queryBooking := `
			SELECT 
				b.*, 
				rp.name as rate_plan_name, 
				r.room_number, 
				rt.name as room_type_name,
				u.email as user_email,
				u.username as user_name
			FROM bookings b
			JOIN rate_plans rp ON b.rate_plan_id = rp.rate_plan_id
			JOIN rooms r ON b.room_id = r.room_id
			JOIN roomtypes rt ON r.room_type_id = rt.room_type_id
			JOIN users u ON b.user_id = u.user_id
			WHERE b.booking_id = $1`

	err = tx.GetContext(ctx, &mBookingDetail, queryBooking, bookingID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("booking id %d: %w", bookingID, errs.ErrNotFound)
		}
		return nil, err
	}

	var mAddons []*model.BookingAddon
	queryAddons := `SELECT ba.*,a.name as addon_name
									FROM booking_addons ba
									JOIN addons a ON ba.addon_id = a.addon_id
									WHERE ba.booking_id = $1`
	err = tx.SelectContext(ctx, &mAddons, queryAddons, bookingID)
	if err != nil {
		return nil, err
	}

	booking := mBookingDetail.ToDomainDetail(mAddons)
	return booking, tx.Commit()
}

func (r *BookingRepository) UpdateBookingStatus(ctx context.Context, bookingID int, status string) error {
	q := `UPDATE bookings SET status=$1 WHERE booking_id=$2`
	result, err := r.db.ExecContext(ctx, q, status, bookingID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("no booking found with id %d: %w", bookingID, errs.ErrNotFound)
	}

	return nil
}

func (r *BookingRepository) SyncBookingAddons(ctx context.Context, bookingID int, addons []*domain.BookingAddon, newTotalPrice float64) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	queryUpdateTotal := `UPDATE bookings SET total_price = $1, updated_at = NOW() WHERE booking_id = $2`
	result, err := tx.ExecContext(ctx, queryUpdateTotal, newTotalPrice, bookingID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("booking id %d not found", bookingID)
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM booking_addons WHERE booking_id = $1", bookingID)
	if err != nil {
		return err
	}

	queryInsert := `INSERT INTO booking_addons (booking_id, addon_id, quantity, price_at_time_of_booking) VALUES ($1, $2, $3, $4)`
	for _, a := range addons {
		mAddon := model.FromDomainBookingAddon(a)

		_, err = tx.ExecContext(ctx, queryInsert, bookingID, mAddon.AddonID, mAddon.Quantity, mAddon.PriceAtBooking)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *BookingRepository) GetBookingAddonsByBookingID(ctx context.Context, bookingID int) ([]*domain.BookingAddon, error) {
	var mAddons []*model.BookingAddon
	query := `SELECT
            	ba.booking_addon_id,ba.booking_id, ba.addon_id, ba.quantity, ba.price_at_time_of_booking,
            	a.name as addon_name
       			FROM booking_addons ba
        		JOIN addons a ON ba.addon_id = a.addon_id
        		WHERE ba.booking_id = $1`

	err := r.db.SelectContext(ctx, &mAddons, query, bookingID)
	if err != nil {
		return nil, err
	}

	var dAddons []*domain.BookingAddon
	for _, ma := range mAddons {
		dAddons = append(dAddons, ma.ToDomain())
	}
	return dAddons, nil
}

func (r *BookingRepository) CancelExpiredBookings(ctx context.Context) (int64, error) {
	q := `
        UPDATE bookings
        SET status = 'cancelled'
        WHERE status = 'pending'
          AND expired_at < NOW()`

	result, err := r.db.ExecContext(ctx, q)
	if err != nil {
		return 0, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rows, nil
}

func (r *BookingRepository) GetBookingsByUserID(ctx context.Context, userID int) ([]*domain.BookingDetail, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var mBookingDetail []model.BookingDetail
	queryBooking := `
			SELECT
				b.*,
				rp.name as rate_plan_name,
				r.room_number,
				rt.name as room_type_name
			FROM bookings b
			JOIN rate_plans rp ON b.rate_plan_id = rp.rate_plan_id
			JOIN rooms r ON b.room_id = r.room_id
			JOIN roomtypes rt ON r.room_type_id = rt.room_type_id
			WHERE b.user_id = $1
			ORDER BY b.created_at DESC`

	err = tx.SelectContext(ctx, &mBookingDetail, queryBooking, userID)
	if err != nil {
		return nil, err
	}

	var result []*domain.BookingDetail
	for _, m := range mBookingDetail {
		var mAddons []*model.BookingAddon
		queryAddons := `SELECT ba.*,a.name as addon_name
										FROM booking_addons ba
										JOIN addons a ON ba.addon_id = a.addon_id
										WHERE ba.booking_id = $1`
		err = tx.SelectContext(ctx, &mAddons, queryAddons, m.BookingID)
		if err != nil {
			return nil, err
		}

		result = append(result, m.ToDomainDetail(mAddons))
	}
	return result, tx.Commit()
}
func (r *BookingRepository) GetAllBookings(ctx context.Context) ([]*domain.BookingDetail, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var mBookingDetail []model.BookingDetail
	// Join with Users to get Booker info
	queryBooking := `
			SELECT
				b.*,
				rp.name as rate_plan_name,
				r.room_number,
				rt.name as room_type_name,
				u.username as user_name,
				u.email as user_email
			FROM bookings b
			JOIN rate_plans rp ON b.rate_plan_id = rp.rate_plan_id
			JOIN rooms r ON b.room_id = r.room_id
			JOIN roomtypes rt ON r.room_type_id = rt.room_type_id
			JOIN users u ON b.user_id = u.user_id
			ORDER BY b.created_at DESC`

	err = tx.SelectContext(ctx, &mBookingDetail, queryBooking)
	if err != nil {
		return nil, err
	}

	var result []*domain.BookingDetail
	for _, m := range mBookingDetail {
		result = append(result, m.ToDomainDetail(nil))
	}
	return result, tx.Commit()
}
