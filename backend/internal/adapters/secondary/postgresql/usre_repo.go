package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ingwrok/hotelBooking/internal/adapters/secondary/postgresql/model"
	"github.com/ingwrok/hotelBooking/internal/common/errs"
	"github.com/ingwrok/hotelBooking/internal/core/domain"
	"github.com/ingwrok/hotelBooking/internal/core/ports"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct{ db *sqlx.DB }

func NewUserRepository(db *sqlx.DB) ports.UserRepoPort {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *domain.User) error {
	m := model.FromDomain(u)

	q := `INSERT INTO users(username, email, password_hash, is_admin) VALUES($1,$2,$3,$4) RETURNING user_id`

	var newID int
	err := r.db.QueryRowContext(ctx, q, m.Username, m.Email, m.PasswordHash, m.IsAdmin).Scan(&newID)
	if err != nil {
		return err
	}
	u.UserID = newID
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	var m model.User

	q := `SELECT user_id, username, email, password_hash, is_admin FROM users WHERE user_id=$1`

	err := r.db.GetContext(ctx, &m, q, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", errs.ErrNotFound)
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var m model.User
	q := `SELECT user_id, username, email, password_hash, is_admin FROM users WHERE username=$1`

	err := r.db.GetContext(ctx, &m, q, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", errs.ErrNotFound)
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *UserRepository) Update(ctx context.Context, id int, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return nil
	}

	allowed := map[string]bool{
		"username":      true,
		"email":         true,
		"password_hash": true,
	}

	var set []string
	var args []interface{}
	i := 1
	for k, v := range fields {
		if !allowed[k] {
			return fmt.Errorf("column %q not allowed", k)
		}
		set = append(set, fmt.Sprintf("%s=$%d", k, i))
		args = append(args, v)
		i++
	}

	args = append(args, id)
	q := fmt.Sprintf("UPDATE users SET %s WHERE user_id=$%d", strings.Join(set, ", "), i)

	res, err := r.db.ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}
	if aff, _ := res.RowsAffected(); aff == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `DELETE FROM bookings WHERE user_id=$1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete user bookings: %w", err)
	}

	q := `DELETE FROM users WHERE user_id=$1`
	res, err := tx.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("no user found with id: %d: %w", id, errs.ErrNotFound)
	}

	return tx.Commit()
}

func (r *UserRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	var models []model.User
	q := `SELECT user_id, username, email, password_hash, is_admin FROM users`

	err := r.db.SelectContext(ctx, &models, q)
	if err != nil {
		return nil, err
	}

	users := make([]*domain.User, len(models))
	for i, m := range models {
		users[i] = m.ToDomain()
	}
	return users, nil
}
