package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/ingwrok/hotelBooking/internal/common/errs"
	"github.com/ingwrok/hotelBooking/internal/common/logger"
	"github.com/ingwrok/hotelBooking/internal/core/domain"
	"github.com/ingwrok/hotelBooking/internal/core/ports"
	"github.com/ingwrok/hotelBooking/internal/core/utils"
	"go.uber.org/zap"
)

type RoomService struct {
	repo ports.RoomRepository
}

func NewRoomService(repo ports.RoomRepository) *RoomService {
	return &RoomService{repo: repo}
}

func (s *RoomService)	AddRoom(ctx context.Context, room *domain.Room) (*domain.Room,error){
	logger.Info("AddRoom called",
		zap.Int("RoomTypeID", room.RoomTypeID),
		zap.String("RoomNumber", room.RoomNumber),
	)
	if room.RoomTypeID <= 0 || room.RoomNumber == "" {
		logger.Warn("validation failed: missing roomTypeID or roomNumber")
		return nil,errs.NewValidationError("room type ID and number are required")
	}

	err := s.repo.CreateRoom(ctx,room)
	if err != nil {
		logger.ErrorErr(err, "repo.CreateRoom failed")
		return nil,errs.NewUnexpectedError("failed to create room")
	}

	logger.Info("room created successfully", zap.Int("roomID", room.RoomID))
	return room,nil
}

func (s *RoomService)RemoveRoom(ctx context.Context, id int) error{
	logger.Info("RemoveRoom called", zap.Int("roomID", id))

	if id <= 0 {
		logger.Warn("validation failed: missing roomID")
		return errs.NewValidationError("room ID is required")
	}

	err := s.repo.DeleteRoom(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn("room not found", zap.Int("roomID", id))
			return errs.NewNotFoundError("room not found")
		}
		logger.ErrorErr(err, "DeleteRoom failed")
		return errs.NewUnexpectedError("failed to delete room")
	}
	return nil
}

func (s *RoomService)ChangeRoomStatus(ctx context.Context, roomID int, status string) error{
	logger.Info("ChangeRoomStatus called",
		zap.Int("roomID", roomID),
		zap.String("status", status),
	)

	if roomID <= 0 || status == "" {
		logger.Warn("validation failed: missing roomID or status")
		return errs.NewValidationError("room ID and status are required")
	}

	normalizedStatus := strings.ToLower(status)

	validStatuses := map[string]bool{
		"available":   true,
		"maintenance": true,
		"occupied":    true,
		"dirty":       true,
	}

	if !validStatuses[normalizedStatus] {
		logger.Warn("invalid room status", zap.String("status", status))
		return errs.NewValidationError("invalid room status")
	}

	err := s.repo.UpdateRoomStatus(ctx, roomID, normalizedStatus)
	if err != nil {
    if errors.Is(err, errs.ErrNotFound) {
			logger.Warn("room not found", zap.Int("roomID", roomID))
      return errs.NewNotFoundError("room not found")
    }
		logger.ErrorErr(err, "UpdateRoomStatus failed")
    return errs.NewUnexpectedError("failed to update room")
	}

	logger.Info("room status updated",
		zap.Int("roomID", roomID),
		zap.String("status", normalizedStatus),
	)
	return nil
}

func (s *RoomService)GetRoom(ctx context.Context, id int) (*domain.Room, error){
	logger.Info("GetRoom called", zap.Int("roomID", id))

	if id <= 0 {
		logger.Warn("validation failed: missing roomID")
		return nil, errs.NewValidationError("room ID is required")
	}

	room, err := s.repo.GetRoomByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn("room not found", zap.Int("roomID", id))
			return nil, errs.NewNotFoundError("room not found")
		}
		logger.ErrorErr(err, "GetRoomByID failed")
		return nil, errs.NewUnexpectedError("failed to get room")
	}
	logger.Debug("room fetched", zap.Int("roomID", id))
	return room, nil
}

func (s *RoomService)ListRooms(ctx context.Context) ([]*domain.Room, error){
	logger.Info("ListRooms called")

	rooms, err := s.repo.GetAllRooms(ctx)
	if err != nil {
		logger.ErrorErr(err, "GetAllRooms failed")
		return nil, errs.NewUnexpectedError("failed to retrieve list rooms")
	}

	logger.Debug("room list returned", zap.Int("count", len(rooms)))
	return rooms, nil
}

func (s *RoomService)BlockRoom(ctx context.Context, block *domain.RoomBlock) error{
	logger.Info("BlockRoom called",
		zap.Int("roomID", block.RoomID),
		zap.Time("start", block.StartDate),
		zap.Time("end", block.EndDate),
	)

	if block.RoomID <= 0 {
		logger.Warn("validation failed: missing roomID")
		return errs.NewValidationError("room ID is required")
	}

	if block.StartDate.After(block.EndDate) {
		logger.Warn("block date invalid",
			zap.Time("start", block.StartDate),
			zap.Time("end", block.EndDate),
		)
    return errs.NewValidationError("start date must be before or equal to end date")
	}
	if block.StartDate.Before(time.Now().Truncate(24 * time.Hour)) {
		logger.Warn("block start date in the past")
		return errs.NewValidationError("block start date cannot be in the past")
	}

	count,err := s.repo.CheckIfBlockOverlaps(ctx,block.RoomID,block.StartDate,block.EndDate)
	if err != nil {
		logger.ErrorErr(err, "CheckIfBlockOverlaps failed")
		return errs.NewUnexpectedError("failed to check room block overlaps")
	}
	if count > 0 {
		logger.Warn("room block overlaps with existing block", zap.Int("count", count))
		return errs.NewValidationError("room block overlaps with existing block")
	}

	err = s.repo.CreateRoomBlock(ctx,block)
	if err != nil {
		logger.ErrorErr(err, "CreateRoomBlock failed")
		return errs.NewUnexpectedError("failed to create room block record")
	}

	logger.Info("room blocked successfully", zap.Int("roomID", block.RoomID))
	return nil
}

func (s *RoomService)GetRoomBlocks(ctx context.Context, roomID int) ([]*domain.RoomBlock, error){
	logger.Info("GetRoomBlocks called", zap.Int("roomID", roomID))

	blocks,err := s.repo.GetRoomBlocksByRoomID(ctx, roomID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn("no blocks found", zap.Int("roomID", roomID))
			return nil, errs.NewNotFoundError("no room blocks found for the specified room")
		}
		logger.ErrorErr(err,"GetRoomBlocksByRoomID failed")
		return nil, errs.NewUnexpectedError("failed to retrieve room blocks")
	}
	logger.Debug("room blocks fetched", zap.Int("count", len(blocks)))
	return blocks, nil
}

func (s *RoomService)	UnblockRoom(ctx context.Context, blockID int) error{
	logger.Info("UnblockRoom called", zap.Int("blockID", blockID))

	if blockID <= 0 {
		logger.Warn("validation failed: missing blockID")
		return errs.NewValidationError("block ID is required")
	}

	err := s.repo.DeleteRoomBlock(ctx,blockID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn("room block not found", zap.Int("blockID", blockID))
			return errs.NewNotFoundError("room block not found")
		}
		logger.ErrorErr(err,"DeleteRoomBlock failed")
		return errs.NewUnexpectedError("failed to delete room block")
	}

	logger.Info("room block deleted", zap.Int("blockID", blockID))
	return nil
}

func (s *RoomService)CountAvailableRooms(ctx context.Context, checkInStr, checkOutStr string) (map[int]int, error){

	checkIn, err := utils.ParseDate(checkInStr,"check-in")
	if err != nil {
		logger.Warn(err.Error(), zap.String("checkIn", checkInStr))
		return nil, err
	}

	checkOut, err := utils.ParseDate(checkOutStr,"check-out")
	if err != nil {
		logger.Warn(err.Error(), zap.String("checkIn", checkInStr))
		return nil, err
	}

	logger.Info("CountAvailableRooms called",
			zap.Time("checkIn", checkIn),
			zap.Time("checkOut", checkOut),
	)

	if checkIn.After(checkOut) {
		logger.Warn("Validation failed: check-in date is after check-out date",
			zap.Time("checkIn", checkIn),
			zap.Time("checkOut", checkOut),
		)
		return nil, errs.NewValidationError("check-in date must be before check-out date")
	}

	counts,err := s.repo.GetAvailableRoomCounts(ctx, checkIn, checkOut)
	if err != nil {
		logger.ErrorErr(err, "GetAvailableRoomCounts failed")
		return nil, errs.NewUnexpectedError("failed to get available room counts")
	}

	logger.Debug("Available room counts calculated", zap.Int("total_types", len(counts)))
	return counts, nil
}

func (s *RoomService)FindAvailableRoom(ctx context.Context, roomTypeID int, checkInStr, checkOutStr string) (int, error){

	if roomTypeID <= 0 {
		logger.Warn("validation failed: missing roomTypeID")
		return 0, errs.NewValidationError("room type ID is required")
	}

	checkIn, err := utils.ParseDate(checkInStr,"check-in")
	if err != nil {
		logger.Warn(err.Error(), zap.String("checkIn", checkInStr))
		return 0, err
	}

	checkOut, err := utils.ParseDate(checkOutStr,"check-out")
	if err != nil {
		logger.Warn(err.Error(), zap.String("checkIn", checkInStr))
		return 0, err
	}


	logger.Info("FindAvailableRoom called",
		zap.Int("roomTypeID", roomTypeID),
		zap.Time("checkIn", checkIn),
		zap.Time("checkOut", checkOut),
	)

	if checkIn.After(checkOut) {
		logger.Warn("Validation failed: check-in date is after check-out date")
		return 0, errs.NewValidationError("check-in date must be before check-out date")
	}
	roomID,err := s.repo.GetAnyAvailableRoomID(ctx, roomTypeID, checkIn, checkOut)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn("No rooms found for type", zap.Int("roomTypeID", roomTypeID), zap.Error(err))
			return 0, errs.NewNotFoundError("no available room found for the specified type and dates")
		}
		logger.ErrorErr(err, "GetAnyAvailableRoomID failed")
		return 0, errs.NewUnexpectedError("failed to find available room")
	}

	logger.Debug("Room assigned successfully", zap.Int("roomID", roomID))
	return roomID, nil
}