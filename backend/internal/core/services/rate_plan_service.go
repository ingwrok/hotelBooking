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

type RatePlanService struct {
	repo ports.RatePlanRepository
}

func NewRatePlanService(repo ports.RatePlanRepository) *RatePlanService {
	return &RatePlanService{repo: repo}
}

func (s *RatePlanService) AddRatePlan(ctx context.Context, rp *domain.RatePlan) (*domain.RatePlan, error) {
	logger.Info("AddRatePlan called",
		zap.Int("RatePlanID", rp.RatePlanID),
		zap.String("Name", rp.Name),
	)
	if rp.Name == "" || rp.Description == "" {
		logger.Warn("validation failed: missing rate plan name or description")
		return nil, errs.NewValidationError("rate plan name and description is required")
	}

	err := s.repo.CreateRatePlan(ctx, rp)
	if err != nil {
		logger.ErrorErr(err, "repo.CreateRatePlan failed")
		return nil, errs.NewUnexpectedError("failed to create rate plan")
	}

	logger.Info("rate plan created successfully", zap.Int("RatePlanID", rp.RatePlanID))
	return rp, nil
}

func (s *RatePlanService) ChangeRatePlan(ctx context.Context, rp *domain.RatePlan) error {
	logger.Info("ChangeRatePlan called",
		zap.Int("RatePlanID", rp.RatePlanID),
		zap.String("Name", rp.Name),
	)

	if rp.RatePlanID <= 0 {
		return errs.NewValidationError("invalid rate plan ID")
	}
	if rp.Name == "" {
		return errs.NewValidationError("rate plan name is required")
	}

	err := s.repo.UpdateRatePlan(ctx, rp)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn("rate plan not found", zap.Int("RatePlanID", rp.RatePlanID))
			return errs.NewNotFoundError("rate plan not found")
		}
		logger.ErrorErr(err, "repo.UpdateRatePlan failed")
		return errs.NewUnexpectedError("failed to update rate plan")
	}

	logger.Info("rate plan updated successfully", zap.Int("RatePlanID", rp.RatePlanID))
	return nil
}

func (s *RatePlanService) RemoveRatePlan(ctx context.Context, ratePlanID int) error {
	logger.Info("RemoveRatePlan called", zap.Int("RatePlanID", ratePlanID))
	err := s.repo.DeleteRatePlan(ctx, ratePlanID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn("rate plan not found", zap.Int("RatePlanID", ratePlanID))
			return errs.NewNotFoundError("rate plan not found")
		}
		logger.ErrorErr(err, "repo.DeleteRatePlan failed")
		return errs.NewUnexpectedError("failed to delete rate plan")
	}

	logger.Info("rate plan deleted successfully", zap.Int("RatePlanID", ratePlanID))
	return nil
}

func (s *RatePlanService) GetRatePlan(ctx context.Context, ratePlanID int) (*domain.RatePlan, error) {
	logger.Info("GetRatePlan called", zap.Int("RatePlanID", ratePlanID))
	rp, err := s.repo.GetRatePlanByID(ctx, ratePlanID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn("rate plan not found", zap.Int("RatePlanID", ratePlanID))
			return nil, errs.NewNotFoundError("rate plan not found")
		}
		logger.ErrorErr(err, "repo.GetRatePlanByID failed")
		return nil, errs.NewUnexpectedError("failed to get rate plan")
	}

	logger.Debug("rate plan fetched", zap.Int("RatePlanID", rp.RatePlanID))
	return rp, nil
}

func (s *RatePlanService) ListRatePlans(ctx context.Context) ([]*domain.RatePlan, error) {
	logger.Info("ListRatePlans called")
	rps, err := s.repo.GetAllRatePlans(ctx)
	if err != nil {
		logger.ErrorErr(err, "repo.GetAllRatePlans failed")
		return nil, errs.NewUnexpectedError("failed to get all rate plans")
	}

	logger.Debug("rate plan list return", zap.Int("count", len(rps)))
	return rps, nil
}

func (s *RatePlanService) ChangeRoomTypePrice(ctx context.Context, roomTypeID, ratePlanID int, price float64) error {
	logger.Info("ChangeRoomTypePrice called",
		zap.Int("RoomTypeID", roomTypeID),
		zap.Int("RatePlanID", ratePlanID),
	)

	if roomTypeID <= 0 || ratePlanID <= 0 {
		logger.Warn("validation failed: missing room type ID or rate plan ID")
		return errs.NewValidationError("room type ID or rate plan ID is required")
	}

	if price <= 0 {
		logger.Warn("validation failed: missing price")
		return errs.NewValidationError("price is required")
	}

	err := s.repo.SetRoomTypePrice(ctx, roomTypeID, ratePlanID, price)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn("room type not found", zap.Int("RoomTypeID", roomTypeID))
			return errs.NewNotFoundError("room type not found")
		}
		logger.ErrorErr(err, "repo.UpdateRoomTypePrice failed")
		return errs.NewUnexpectedError("failed to update room type price")
	}

	logger.Info("room type price updated",
		zap.Int("RoomTypeID", roomTypeID),
		zap.Int("RatePlanID", ratePlanID),
	)
	return nil
}

func (s *RatePlanService) GetPrice(ctx context.Context, roomTypeID, ratePlanID int) (float64, error) {
	logger.Info("GetPrice called",
		zap.Int("RoomTypeID", roomTypeID),
		zap.Int("RatePlanID", ratePlanID),
	)

	if roomTypeID <= 0 || ratePlanID <= 0 {
		logger.Warn("validation failed: missing room type ID or rate plan ID")
		return 0, errs.NewValidationError("room type ID or rate plan ID is required")
	}

	price, err := s.repo.GetPriceByRoomType(ctx, roomTypeID, ratePlanID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn("room type not found", zap.Int("RoomTypeID", roomTypeID))
			return 0, errs.NewNotFoundError("room type not found")
		}
		logger.ErrorErr(err, "repo.GetPriceByRoomType failed")
		return 0, errs.NewUnexpectedError("failed to get room type price")
	}

	logger.Info("room type price fetched",
		zap.Int("RoomTypeID", roomTypeID),
		zap.Int("RatePlanID", ratePlanID),
	)
	return price, nil
}

func (s *RatePlanService) RemoveRoomTypePrice(ctx context.Context, roomTypeID, ratePlanID int) error {
	logger.Info("RemoveRoomTypePrice called",
		zap.Int("RoomTypeID", roomTypeID),
		zap.Int("RatePlanID", ratePlanID),
	)

	if roomTypeID <= 0 || ratePlanID <= 0 {
		logger.Warn("validation failed: missing room type ID or rate plan ID")
		return errs.NewValidationError("room type ID or rate plan ID is required")
	}

	err := s.repo.DeleteRoomTypePrice(ctx, roomTypeID, ratePlanID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn("room type not found", zap.Int("RoomTypeID", roomTypeID))
			return errs.NewNotFoundError("room type not found")
		}
		logger.ErrorErr(err, "repo.DeleteRoomTypePrice failed")
		return errs.NewUnexpectedError("failed to delete room type price")
	}

	logger.Info("room type price deleted",
		zap.Int("RoomTypeID", roomTypeID),
		zap.Int("RatePlanID", ratePlanID),
	)
	return nil
}

func (s *RatePlanService) ListRatePlansByRoomType(ctx context.Context, roomTypeID int) ([]*domain.RatePlanFull, error) {
	logger.Info("ListRatePlansByRoomType called",
		zap.Int("RoomTypeID", roomTypeID),
	)

	if roomTypeID <= 0 {
		logger.Warn("validation failed: missing room type ID")
		return nil, errs.NewValidationError("room type ID is required")
	}

	rps, err := s.repo.GetAllRatePlansByRoomTypeID(ctx, roomTypeID)
	if err != nil {
		logger.ErrorErr(err, "repo.ListRatePlansByRoomType failed")
		return nil, errs.NewUnexpectedError("failed to list rate plans by room type")
	}

	return rps, nil
}
