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

type AddonRepository struct {
	db *sqlx.DB
}

func NewAddonRepository(db *sqlx.DB) ports.AddonRepository {
	return &AddonRepository{db: db}
}

// CreateAddonCategory: สร้าง Category ใหม่ใน DB
func (r *AddonRepository) CreateAddonCategory(ctx context.Context, category *domain.AddonCategory) error {
	m := model.FromDomainAddonCategory(category)

	q := `INSERT INTO addon_categories (name)
        VALUES ($1)
        RETURNING category_id
        `
	var newID int

	err := r.db.QueryRowContext(ctx, q, m.Name).Scan(&newID)
	if err != nil {
		return err
	}

	category.CategoryID = newID
	return nil
}

func (r *AddonRepository) GetAddonCategoryByID(ctx context.Context, categoryID int) (*domain.AddonCategory, error) {
	m := &model.AddonCategory{}

	q := `SELECT category_id, name
        FROM addon_categories
        WHERE category_id = $1`

	err := r.db.GetContext(ctx, m, q, categoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("addon category id %d: %w", categoryID, errs.ErrNotFound)
		}
		return nil, err
	}

	return m.ToDomain(), nil
}

func (r *AddonRepository) GetAllAddonCategories(ctx context.Context) ([]*domain.AddonCategory, error) {
	var dbModels []*model.AddonCategory

	q := `SELECT category_id, name
        FROM addon_categories
        ORDER BY name ASC`

	err := r.db.SelectContext(ctx, &dbModels, q)
	if err != nil {
		return nil, err
	}

	domainModels := make([]*domain.AddonCategory, len(dbModels))
	for i, m := range dbModels {
		domainModels[i] = m.ToDomain()
	}

	return domainModels, nil
}

func (r *AddonRepository) UpdateAddonCategory(ctx context.Context, category *domain.AddonCategory) error {
	m := model.FromDomainAddonCategory(category)

	q := `UPDATE addon_categories
        SET name = $2
        WHERE category_id = $1`

	result, err := r.db.ExecContext(ctx, q, m.CategoryID, m.Name)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("no category found with id %d: %w", m.CategoryID, errs.ErrNotFound)
	}

	return nil
}

func (r *AddonRepository) DeleteAddonCategory(ctx context.Context, categoryID int) error {
	q := `DELETE FROM addon_categories
        WHERE category_id = $1`

	result, err := r.db.ExecContext(ctx, q, categoryID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("no category found with id %d: %w", categoryID, errs.ErrNotFound)
	}

	return nil
}

func (r *AddonRepository) CreateAddon(ctx context.Context, addon *domain.Addon) error {
	m := model.FromDomainAddon(addon)
	q := `INSERT INTO addons (category_id, name, description, price, unit_name)
	      VALUES ($1, $2, $3, $4, $5)
				RETURNING addon_id`

	var newID int
	err := r.db.QueryRowContext(ctx, q, m.CategoryID, m.Name, m.Description, m.Price, m.UnitName).Scan(&newID)
	if err != nil {
		return err
	}

	addon.AddonID = newID
	return nil
}

func (r *AddonRepository) DeleteAddon(ctx context.Context, addonID int) error {
	q := `DELETE FROM addons WHERE addon_id = $1`

	result, err := r.db.ExecContext(ctx, q, addonID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("no addon found with id %d", addonID)
	}

	return nil
}

func (r *AddonRepository) UpdateAddon(ctx context.Context, addon *domain.Addon) error {
	m := model.FromDomainAddon(addon)
	q := `UPDATE addons
	      SET category_id = $1,
					name = $2,
					description = $3,
					price = $4,
					unit_name = $5
				WHERE addon_id = $6`

	result, err := r.db.ExecContext(ctx, q,
		m.CategoryID,
		m.Name,
		m.Description,
		m.Price,
		m.UnitName,
		m.AddonID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("no addon found with id %d", m.AddonID)
	}

	return nil
}

func (r *AddonRepository) GetAddonByID(ctx context.Context, addonID int) (*domain.Addon, error) {
	var m model.Addon
	q := `SELECT addon_id, category_id, name, description, price, unit_name
	      FROM addons
				WHERE addon_id = $1`

	err := r.db.GetContext(ctx, &m, q, addonID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("addon id %d: %w", addonID, errs.ErrNotFound)
		}
		return nil, err
	}

	return m.ToDomain(), nil
}

func (r *AddonRepository) GetAllAddons(ctx context.Context) ([]*domain.Addon, error) {
	var models []model.Addon
	q := `SELECT addon_id, category_id, name, description, price, unit_name
	      FROM addons`

	err := r.db.SelectContext(ctx, &models, q)
	if err != nil {
		return nil, err
	}

	addons := make([]*domain.Addon, len(models))
	for i, m := range models {
		addons[i] = m.ToDomain()
	}

	return addons, nil
}

func (r *AddonRepository) GetAddonByCategoryID(ctx context.Context, categoryID int) ([]*domain.Addon, error) {
	var models []model.Addon
	q := `SELECT addon_id, category_id, name, description, price, unit_name
	      FROM addons
				WHERE category_id = $1`

	err := r.db.SelectContext(ctx, &models, q, categoryID)
	if err != nil {
		return nil, err
	}

	addons := make([]*domain.Addon, len(models))
	for i, m := range models {
		addons[i] = m.ToDomain()
	}

	return addons, nil
}