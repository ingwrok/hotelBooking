package services

import (
	"context"
	"errors"

	"github.com/ingwrok/hotelBooking/internal/common/errs"
	"github.com/ingwrok/hotelBooking/internal/common/logger"
	"github.com/ingwrok/hotelBooking/internal/core/domain"
	"github.com/ingwrok/hotelBooking/internal/core/ports"
	"go.uber.org/zap"
)

type AddonService struct {
	repo ports.AddonRepository
}

func NewAddonService(repo ports.AddonRepository) *AddonService {
	return &AddonService{repo: repo}
}

func (s *AddonService)AddAddonCategory(ctx context.Context, category *domain.AddonCategory) (*domain.AddonCategory,error){
	logger.Info("AddAddonCategory called",
		zap.String("Name", category.Name),
	)
	if category.Name == ""{
		logger.Warn("validation failed: missing name")
		return nil,errs.NewValidationError("validation failed: missing name")
	}

	err := s.repo.CreateAddonCategory(ctx,category)
	if err != nil {
		logger.ErrorErr(err,"repo.CreateAddonCategory failed")
		return nil,errs.NewUnexpectedError("failed to create addon category")
	}

	logger.Info("addon category created successfully", zap.Int("CategoryID", category.CategoryID))
	return category, nil
}

func (s *AddonService)GetAddonCategory(ctx context.Context, categoryID int) (*domain.AddonCategory, error){
	logger.Info("GetAddonCategory called",
		zap.Int("CategoryID", categoryID),
	)

	category, err := s.repo.GetAddonCategoryByID(ctx, categoryID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound){
			logger.Warn("addon category not found", zap.Int("CategoryID", categoryID))
			return nil, errs.NewNotFoundError("addon category not found")
		}
		logger.ErrorErr(err,"repo.GetAddonCategoryByID failed")
		return nil, errs.NewUnexpectedError("failed to get addon category")
	}

	logger.Debug("addon category fetched", zap.Int("CategoryID", category.CategoryID))
	return category, nil
}

func (s *AddonService)ListAddonCategories(ctx context.Context) ([]*domain.AddonCategory, error){
	logger.Info("ListAddonCategories called")

	addonCategories,err := s.repo.GetAllAddonCategories(ctx)
	if err != nil {
		logger.ErrorErr(err,"repo.GetAllAddonCategories failed")
		return nil,errs.NewUnexpectedError("failed to retrieve list addonCategories")
	}

	logger.Debug("addonCategory list returned", zap.Int("count", len(addonCategories)))
	return addonCategories, nil
}

func (s *AddonService)ChangeAddonCategory(ctx context.Context, category *domain.AddonCategory) error{
	logger.Info("ChangeAddonCategory called",
		zap.Int("CategoryID", category.CategoryID),
		zap.String("Name", category.Name),
	)

	if category.CategoryID <= 0 || category.Name == ""{
		logger.Warn("validation failed: missing or invalid input")
		return errs.NewValidationError("addon category invalid input")
	}

	err := s.repo.UpdateAddonCategory(ctx, category)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound){
			logger.Warn("addon category not found", zap.Int("CategoryID", category.CategoryID))
			return errs.NewNotFoundError("addon category not found")
		}
		logger.ErrorErr(err,"repo.UpdateAddonCategory failed")
		return errs.NewUnexpectedError("failed to update addon category")
	}

	logger.Info("addon category updated",
		zap.Int("CategoryID", category.CategoryID),
	)
	return nil
}

func (s *AddonService)RemoveAddonCategory(ctx context.Context, categoryID int) error {
	logger.Info("RemoveAddonCategory called", zap.Int("CategoryID", categoryID))

  if categoryID <= 0 {
    logger.Warn("validation failed: invalid category ID")
    return errs.NewValidationError("addon category ID is required")
  }

  err := s.repo.DeleteAddonCategory(ctx, categoryID)
  if err != nil {
    if errors.Is(err, errs.ErrNotFound) {
      logger.Warn("addon category not found", zap.Int("CategoryID", categoryID))
      return errs.NewNotFoundError("addon category not found")
    }

    logger.ErrorErr(err, "repo.DeleteAddonCategory failed")
    return errs.NewUnexpectedError("failed to delete addon category")
  }

  logger.Info("addon category deleted successfully", zap.Int("CategoryID", categoryID))
  return nil
}

// addon
func (s *AddonService)AddAddon(ctx context.Context , addon *domain.Addon) (*domain.Addon,error){
	logger.Info("AddAddon called",
		zap.Int("CategoryID", addon.CategoryID),
		zap.String("Name", addon.Name),
		zap.String("Description", addon.Description),
		zap.Float64("Price", addon.Price),
		zap.String("UnitName", addon.UnitName),
	)

	if addon.CategoryID <= 0 || addon.Name == "" || addon.Description == "" || addon.Price <= 0 || addon.UnitName == ""{
		logger.Warn("validation failed: missing or invalid input")
		return nil,errs.NewValidationError("addon invalid input")
	}

	err := s.repo.CreateAddon(ctx, addon)
	if err != nil {
		logger.ErrorErr(err,"repo.CreateAddon failed")
		return nil,errs.NewUnexpectedError("failed to create addon")
	}

	logger.Info("addon created successfully", zap.Int("AddonID", addon.AddonID))
	return addon, nil
}

func (s *AddonService)RemoveAddon(ctx context.Context, addonID int) error{
	logger.Info("RemoveAddon called", zap.Int("AddonID", addonID))

	if addonID <= 0 {
		logger.Warn("validation failed: invalid addon ID")
		return errs.NewValidationError("addon ID is required")
	}

	err := s.repo.DeleteAddon(ctx, addonID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn("addon not found", zap.Int("AddonID", addonID))
			return errs.NewNotFoundError("addon not found")
		}
		logger.ErrorErr(err, "repo.DeleteAddon failed")
		return errs.NewUnexpectedError("failed to delete addon")
	}

	logger.Info("addon deleted successfully", zap.Int("AddonID", addonID))
	return nil
}

func (s *AddonService)ChangeAddon(ctx context.Context, addon *domain.Addon) error{
	logger.Info("ChangeAddon called",
		zap.Int("AddonID", addon.AddonID),
		zap.Int("CategoryID", addon.CategoryID),
		zap.String("Name", addon.Name),
		zap.String("Description", addon.Description),
		zap.Float64("Price", addon.Price),
		zap.String("UnitName", addon.UnitName),
	)

	if addon.AddonID <= 0 || addon.CategoryID <= 0 || addon.Name == "" || addon.Description == "" || addon.Price <= 0 || addon.UnitName == ""{
		logger.Warn("validation failed: missing or invalid input")
		return errs.NewValidationError("addon invalid input")
	}

	err := s.repo.UpdateAddon(ctx, addon)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound){
			logger.Warn("addon not found", zap.Int("AddonID", addon.AddonID))
			return errs.NewNotFoundError("addon not found")
		}
		logger.ErrorErr(err,"repo.UpdateAddon failed")
		return errs.NewUnexpectedError("failed to update addon")
	}

	logger.Info("addon updated successfully", zap.Int("AddonID", addon.AddonID))
	return nil
}

func (s *AddonService)GetAddon(ctx context.Context, addonID int) (*domain.Addon, error){
	logger.Info("GetAddon called", zap.Int("AddonID", addonID))

	if addonID <= 0 {
		logger.Warn("validation failed: invalid addon ID")
		return nil, errs.NewValidationError("addon ID is required")
	}

	addon, err := s.repo.GetAddonByID(ctx, addonID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound){
			logger.Warn("addon not found", zap.Int("AddonID", addonID))
			return nil, errs.NewNotFoundError("addon not found")
		}
		logger.ErrorErr(err,"repo.GetAddonByID failed")
		return nil, errs.NewUnexpectedError("failed to get addon")
	}

	logger.Debug("addon fetched",zap.Int("addonID",addonID))
	return addon, nil
}

func (s *AddonService)ListAddons(ctx context.Context) ([]*domain.Addon, error){
	logger.Info("ListAddons called")

	addons, err := s.repo.GetAllAddons(ctx)
	if err != nil {
		logger.ErrorErr(err,"repo.ListAddons failed")
		return nil, errs.NewUnexpectedError("failed to list addons")
	}

	logger.Debug("addon list returned", zap.Int("count", len(addons)))
	return addons, nil
}

func (s *AddonService)ListAddonsByCategory(ctx context.Context, categoryID int) ([]*domain.Addon, error){
	logger.Info("ListAddonsByCategory called", zap.Int("CategoryID", categoryID))

	if categoryID <= 0 {
		logger.Warn("validation failed: invalid category ID")
		return nil, errs.NewValidationError("category ID is required")
	}

	addons, err := s.repo.GetAddonByCategoryID(ctx, categoryID)
	if err != nil {
		logger.ErrorErr(err,"repo.ListAddonsByCategory failed")
		return nil, errs.NewUnexpectedError("failed to list addons by category")
	}

	logger.Debug("addon list returned", zap.Int("count", len(addons)))
	return addons, nil
}
