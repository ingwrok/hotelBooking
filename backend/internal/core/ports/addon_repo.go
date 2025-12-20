package ports

import (
	"context"

	"github.com/ingwrok/hotelBooking/internal/core/domain"
)

type AddonRepository interface {
	// category
	CreateAddonCategory(ctx context.Context, category *domain.AddonCategory) error
	GetAddonCategoryByID(ctx context.Context, categoryID int) (*domain.AddonCategory, error)
	GetAllAddonCategories(ctx context.Context) ([]*domain.AddonCategory, error)
	UpdateAddonCategory(ctx context.Context, category *domain.AddonCategory) error
	DeleteAddonCategory(ctx context.Context, categoryID int) error

	// addon
	CreateAddon(ctx context.Context , addon *domain.Addon) error
	DeleteAddon(ctx context.Context, addonID int) error
	UpdateAddon(ctx context.Context, addon *domain.Addon) error
	GetAddonByID(ctx context.Context, addonID int) (*domain.Addon, error)
	GetAllAddons(ctx context.Context) ([]*domain.Addon, error)
	GetAddonByCategoryID(ctx context.Context, categoryID int) ([]*domain.Addon, error)

}