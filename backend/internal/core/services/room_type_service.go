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

type RoomTypeService struct {
	repo ports.RoomTypeRepository
}

func NewRoomTypeService(repo ports.RoomTypeRepository) *RoomTypeService {
	return &RoomTypeService{repo: repo}
}

func (s *RoomTypeService) AddRoomType(ctx context.Context, rt *domain.RoomType) (*domain.RoomType, error) {
	logger.Info("AddRoomType called",
		zap.String("Name", rt.Name),
	)
	if rt.Name == "" || rt.Capacity <= 0 {
		logger.Warn("validation failed: missing name or invalid capacity")
		return nil, errs.NewValidationError("room type name and valid capacity are required")
	}

	err := s.repo.CreateRoomType(ctx, rt)
	if err != nil {
		logger.ErrorErr(err, "repo.CreateRoomType failed")
		return nil, errs.NewUnexpectedError("failed to create room type")
	}

	logger.Info("room type created successfully", zap.Int("RoomTypeID", rt.RoomTypeID))
	return rt, nil
}

func (s *RoomTypeService) ChangeRoomType(ctx context.Context, rt *domain.RoomType) error{
	logger.Info("ChangeRoomType called",
		zap.Int("roomTypeID", rt.RoomTypeID),
	)

	if rt.Name == "" || rt.Description == "" || rt.BedType == "" ||
		rt.Capacity <= 0 || rt.SizeSQM <= 0 {
			logger.Warn("validation failed: missing or invalid input")
			return errs.NewValidationError("room type invalid input")
	}

	err := s.repo.UpdateRoomType(ctx, rt)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound){
			logger.Warn("room type not found", zap.Int("roomTypeID",rt.RoomTypeID))
			return errs.NewNotFoundError("room type not found")
		}
		logger.ErrorErr(err, "repo.UpdateRoomType failed")
		return errs.NewUnexpectedError("failed to update room type")
	}

	logger.Info("room status updated",
		zap.Int("roomTypeID",rt.RoomTypeID),
	)
	return nil
}

func (s *RoomTypeService) RemoveRoomType(ctx context.Context, id int) error {
	logger.Info("RemoveRoomType called", zap.Int("RoomTypeID", id))

	err := s.repo.DeleteRoomType(ctx,id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn("room type not found", zap.Int("RoomTypeID", id))
			return errs.NewNotFoundError("room type not found")
		}
		logger.ErrorErr(err, "repo.DeleteRoomType failed")
		return errs.NewUnexpectedError("failed to delete room type")
	}
	return nil
}

func (s *RoomTypeService) GetRoomType(ctx context.Context,id int)(*domain.RoomType,error){
	logger.Info("GetRoomType called",zap.Int("roomID",id))

	rt, err := s.repo.GetRoomTypeByID(ctx,id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound){
			logger.Warn("roomType not found",zap.Int("roomTypeID",id))
			return nil, errs.NewNotFoundError("roomType not found")
		}
		logger.ErrorErr(err,"GetRoomTypeByID failed")
		return nil,errs.NewUnexpectedError("failed to get roomType")
	}
	logger.Debug("roomType fetched", zap.Int("roomTypeID",id))
	return rt,nil
}

func (s *RoomTypeService) ListRoomTypes(ctx context.Context) ([]*domain.RoomType, error) {
	logger.Info("ListRoomTypes called")

	rts, err := s.repo.GetAllRoomTypes(ctx)
	if err != nil {
		logger.ErrorErr(err, "GetAllRomTypes failed")
		return nil, errs.NewUnexpectedError("failed to retrieve list roomTypes")
	}
	logger.Debug("roomType list returned",zap.Int("count",len(rts)))
	return rts,nil
}

func (s *RoomTypeService) GetRoomTypeFullDetail(ctx context.Context,id int)(*domain.RoomTypeDetails,error){
	logger.Info("GetRoomTypeFullDetail called",zap.Int("roomID",id))

	rtf, err := s.repo.GetRoomTypeFullDetail(ctx,id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound){
			logger.Warn("roomType not found",zap.Int("roomTypeID",id))
			return nil, errs.NewNotFoundError("roomType not found")
		}
		logger.ErrorErr(err,"GetRoomTypeByID failed")
		return nil,errs.NewUnexpectedError("failed to get roomType")
	}
	logger.Debug("roomType fetched", zap.Int("roomTypeID",id))
	return rtf,nil
}