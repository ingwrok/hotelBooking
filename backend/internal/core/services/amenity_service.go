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

type AmenityService struct {
	repo ports.AmenityRepository
}

func NewAmenityService(repo ports.AmenityRepository) *AmenityService {
	return &AmenityService{repo: repo}
}

func (s *AmenityService)AddAmenity(ctx context.Context, amenity *domain.Amenity) (*domain.Amenity,error){
	logger.Info("AddAmenity called",
		zap.Int("AmenityID", amenity.AmenityID),
		zap.String("AmenityName", amenity.Name),
	)
	if amenity.Name == "" {
		logger.Warn("validation failed: missing amenity name")
		return nil,errs.NewValidationError("amenity name is required")
	}

	err := s.repo.CreateAmenity(ctx,amenity)
	if err != nil {
		logger.ErrorErr(err, "repo.CreateAmenity failed")
		return nil,errs.NewUnexpectedError("failed to create amenity")
	}

	logger.Info("amenity created successfully", zap.Int("amenityID", amenity.AmenityID))
	return amenity,nil
}

func (s *AmenityService)GetAmenity(ctx context.Context, id int) (*domain.Amenity, error){
	logger.Info("GetAmenity called", zap.Int("AmenityID", id))

	amenity,err := s.repo.GetAmenityByID(ctx,id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn("amenity not found", zap.Int("AmenityID", id))
			return nil, errs.NewNotFoundError("amenity not found")
		}
		logger.ErrorErr(err, "repo.GetAmenityByID failed")
		return nil, errs.NewUnexpectedError("failed to get amenity")
	}

	logger.Info("amenity retrieved successfully", zap.Int("amenityID", amenity.AmenityID))
	return amenity,nil
}

func (s *AmenityService)ListAmenities(ctx context.Context) ([]*domain.Amenity, error){
	logger.Info("ListAmenities called")
	amenities,err := s.repo.GetAllAmenities(ctx)
	if err != nil {
		logger.ErrorErr(err, "repo.GetAllAmenities failed")
		return nil, errs.NewUnexpectedError("failed to get amenities")
	}

	logger.Debug("amenity list return", zap.Int("count", len(amenities)))
	return amenities,nil
}

func (s *AmenityService)ChangeAmenity(ctx context.Context, amenity *domain.Amenity) error{
	logger.Info("UpdateAmenity called",
		zap.Int("AmenityID", amenity.AmenityID),
		zap.String("AmenityName", amenity.Name),
	)

	if amenity.Name == "" {
		logger.Warn("validation failed: missing amenity name")
		return errs.NewValidationError("amenity name is required")
	}

	err := s.repo.UpdateAmenity(ctx,amenity)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn("amenity not found", zap.Int("AmenityID", amenity.AmenityID))
			return errs.NewNotFoundError("amenity not found")}
		logger.ErrorErr(err, "repo.UpdateAmenity failed")
		return errs.NewUnexpectedError("failed to update amenity")
	}

	logger.Info("amenity updated successfully", zap.Int("amenityID", amenity.AmenityID))
	return nil
}
func (s *AmenityService) RemoveAmenity(ctx context.Context, id int) error{
	logger.Info("RemoveAmenity called", zap.Int("AmenityID", id))

	err := s.repo.DeleteAmenity(ctx,id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn("amenity not found", zap.Int("AmenityID", id))
			return errs.NewNotFoundError("amenity not found")}
		logger.ErrorErr(err, "repo.DeleteAmenity failed")
		return errs.NewUnexpectedError("failed to delete amenity")
	}

	logger.Info("amenity deleted successfully", zap.Int("amenityID", id))
	return nil
}