package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ingwrok/hotelBooking/internal/common/errs"
	"github.com/ingwrok/hotelBooking/internal/common/logger"
	"github.com/ingwrok/hotelBooking/internal/core/domain"
	"github.com/ingwrok/hotelBooking/internal/core/ports"
	"go.uber.org/zap"
)

type BookingService struct {
	bookingRepo  ports.BookingRepository
	roomRepo     ports.RoomRepository
	rateplanRepo ports.RatePlanRepository
	addonRepo    ports.AddonRepository
	emailRepo    ports.EmailRepository
}

func NewBookingService(b ports.BookingRepository, r ports.RoomRepository, rp ports.RatePlanRepository, a ports.AddonRepository, e ports.EmailRepository) *BookingService {
	return &BookingService{
		bookingRepo:  b,
		roomRepo:     r,
		rateplanRepo: rp,
		addonRepo:    a,
		emailRepo:    e,
	}
}

func (s *BookingService) AddBooking(ctx context.Context, booking *domain.Booking, roomTypeID int) (*domain.Booking, error) {
	logger.Info("AddBooking called",
		zap.Int("UserID", booking.UserID),
		zap.Int("RatePlanID", booking.RatePlanID),
		zap.Int("RoomTypeID", roomTypeID),
	)

	price, err := s.rateplanRepo.GetPriceByRoomType(ctx, roomTypeID, booking.RatePlanID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return nil, fmt.Errorf("rate plan id %d: %w", booking.RatePlanID, errs.ErrNotFound)
		}
		logger.ErrorErr(err, "GetPriceByRoomType failed")
		return nil, err
	}

	roomID, err := s.roomRepo.GetAnyAvailableRoomID(ctx, roomTypeID, booking.CheckInDate, booking.CheckOutDate)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn("No rooms found for type", zap.Int("roomTypeID", roomTypeID), zap.Error(err))
			return nil, errs.NewNotFoundError("no available room found for the specified type and dates")
		}
		logger.ErrorErr(err, "GetAnyAvailableRoomID failed")
		return nil, errs.NewUnexpectedError("failed to find available room")
	}
	booking.RoomID = roomID

	numNights := int(booking.CheckOutDate.Sub(booking.CheckInDate).Hours() / 24)
	if numNights <= 0 {
		logger.Warn("numNight must more 1")
		return nil, fmt.Errorf("invalid stay duration")
	}
	booking.RoomSubTotal = price * float64(numNights)

	var addonTotal float64
	for i := range booking.BookingAddon {
		addon, err := s.addonRepo.GetAddonByID(ctx, booking.BookingAddon[i].AddonID)
		if err != nil {
			if errors.Is(err, errs.ErrNotFound) {
				logger.Warn("addon not found", zap.Int("AddonID", addon.AddonID))
				return nil, errs.NewNotFoundError("addon not found")
			}
			logger.ErrorErr(err, "repo.GetAddonByID failed")
			return nil, errs.NewUnexpectedError("failed to get addon")
		}
		booking.BookingAddon[i].PriceAtBooking = addon.Price
		addonTotal += (addon.Price * float64(booking.BookingAddon[i].Quantity))
	}
	booking.AddonSubTotal = addonTotal

	booking.TaxesAmount = (booking.RoomSubTotal + booking.AddonSubTotal) * 0.07
	booking.TotalPrice = booking.RoomSubTotal + booking.AddonSubTotal + booking.TaxesAmount

	booking.Status = "pending"
	booking.ExpiredAt = time.Now().Add(30 * time.Minute)

	err = s.bookingRepo.CreateBooking(ctx, booking, booking.BookingAddon)
	if err != nil {
		logger.ErrorErr(err, "repo.CreateBooking failed")
		return nil, err
	}

	// Send Confirmation Email
	go func() {
		// Fetch full details (populated with joins) for the email
		// Note: At this point, the transaction is committed, so we can read from DB.
		details, dbErr := s.bookingRepo.GetBookingWithAddons(context.Background(), booking.BookingID)
		if dbErr != nil {
			logger.ErrorErr(dbErr, "failed to fetch booking for email")
			return
		}

		emailCtx := context.Background()
		if emailErr := s.emailRepo.SendBookingConfirmation(emailCtx, details, details.BookingAddon); emailErr != nil {
			logger.ErrorErr(emailErr, "failed to send confirmation email")
		}
	}()

	logger.Info("booking created successfully", zap.Int("BookingID", booking.BookingID))
	return booking, err
}

func (s *BookingService) GetFullDetails(ctx context.Context, bookingID int) (*domain.BookingDetail, error) {
	logger.Info("GetFullDetails called", zap.Int("BookingID", bookingID))

	if bookingID <= 0 {
		logger.Warn("validation failed: missing booking id")
		return nil, errs.NewValidationError("invalid booking id")
	}

	booking, err := s.bookingRepo.GetBookingWithAddons(ctx, bookingID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return nil, fmt.Errorf("booking id %d: %w", bookingID, errs.ErrNotFound)
		}
		logger.ErrorErr(err, "GetFullDetails failed")
		return nil, err
	}

	return booking, nil
}

func (s *BookingService) ChangeStatus(ctx context.Context, bookingID int, status string) error {
	logger.Info("ChangeStatus called", zap.Int("BookingID", bookingID), zap.String("Status", status))

	if bookingID <= 0 || status == "" {
		logger.Warn("validation failed: missing booking id or status")
		return errs.NewValidationError("invalid booking id or status")
	}

	normalizedStatus := strings.ToLower(status)

	validStatues := map[string]bool{
		"pending":     true,
		"confirmed":   true,
		"cancelled":   true,
		"checked-in":  true,
		"checked-out": true,
	}

	if !validStatues[normalizedStatus] {
		logger.Warn("invalid status", zap.String("status", status))
		return errs.NewValidationError("invalid status")
	}

	// Just update status to confirmed for mock
	if err := s.bookingRepo.UpdateBookingStatus(ctx, bookingID, normalizedStatus); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return fmt.Errorf("booking id %d: %w", bookingID, errs.ErrNotFound)
		}
		logger.ErrorErr(err, "ChangeStatus failed")
		return err
	}

	// If status changed to confirmed, resend email (receipt)
	if normalizedStatus == "confirmed" {
		go func() {
			// Fetch full details to get email and addons
			details, dbErr := s.bookingRepo.GetBookingWithAddons(context.Background(), bookingID)
			if dbErr != nil {
				logger.ErrorErr(dbErr, "failed to fetch booking for email")
				return
			}

			details.Status = "confirmed"

			emailCtx := context.Background()
			if emailErr := s.emailRepo.SendBookingConfirmation(emailCtx, details, details.BookingAddon); emailErr != nil {
				logger.ErrorErr(emailErr, "failed to send confirmation email")
			} else {
				logger.Info("confirmation email sent (payment success)", zap.String("email", details.Email))
			}
		}()
	}

	logger.Info("booking status changed successfully", zap.Int("BookingID", bookingID), zap.String("Status", status))
	return nil
}

func (s *BookingService) ModifyBookingAddons(ctx context.Context, bookingID int, newAddons []*domain.BookingAddon) error {
	booking, err := s.bookingRepo.GetBookingWithAddons(ctx, bookingID)
	if err != nil {
		return err
	}

	var newAddonTotal float64
	for i := range newAddons {
		addon, err := s.addonRepo.GetAddonByID(ctx, newAddons[i].AddonID)
		if err != nil {
			return err
		}

		newAddons[i].PriceAtBooking = addon.Price
		newAddonTotal += (addon.Price * float64(newAddons[i].Quantity))
	}

	newTaxes := (booking.RoomSubTotal + newAddonTotal) * 0.07
	newTotalPrice := booking.RoomSubTotal + newAddonTotal + newTaxes

	err = s.bookingRepo.SyncBookingAddons(ctx, bookingID, newAddons, newTotalPrice)
	if err != nil {
		logger.ErrorErr(err, "repo.SyncBookingAddons failed")
		return err
	}

	booking.AddonSubTotal = newAddonTotal
	booking.TaxesAmount = newTaxes
	booking.TotalPrice = newTotalPrice

	return nil
}

func (s *BookingService) GetAddonDetails(ctx context.Context, bookingID int) ([]*domain.BookingAddon, error) {
	logger.Info("GetAddonDetails called", zap.Int("BookingID", bookingID))

	if bookingID <= 0 {
		logger.Warn("validation failed: missing booking id")
		return nil, errs.NewValidationError("invalid booking id")
	}

	addons, err := s.bookingRepo.GetBookingAddonsByBookingID(ctx, bookingID)
	if err != nil {
		logger.ErrorErr(err, "GetAddonDetails failed")
		return nil, err
	}

	logger.Debug("addon details returned", zap.Int("count", len(addons)))
	return addons, nil
}

func (s *BookingService) CleanupExpiredBookings(ctx context.Context) (int64, error) {
	logger.Info("CleanupExpiredBookings called")

	rows, err := s.bookingRepo.CancelExpiredBookings(ctx)
	if err != nil {
		logger.ErrorErr(err, "CleanupExpiredBookings failed")
		return 0, err
	}
	logger.Info("expired bookings cleaned up", zap.Int64("rowsAffected", rows))
	return rows, nil
}

func (s *BookingService) GetMyHistory(ctx context.Context, userID int) ([]*domain.BookingDetail, error) {
	logger.Info("GetMyHistory called", zap.Int("UserID", userID))

	if userID <= 0 {
		logger.Warn("validation failed: missing booking id")
		return nil, errs.NewValidationError("invalid booking id")
	}

	bookings, err := s.bookingRepo.GetBookingsByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return nil, fmt.Errorf("booking id %d: %w", userID, errs.ErrNotFound)
		}
		logger.ErrorErr(err, "GetFullDetails failed")
		return nil, err
	}
	logger.Debug("bookings returned", zap.Int("count", len(bookings)))
	return bookings, nil
}

func (s *BookingService) GetAllBookings(ctx context.Context) ([]*domain.BookingDetail, error) {
	logger.Info("GetAllBookings called (Admin)")
	bookings, err := s.bookingRepo.GetAllBookings(ctx)
	if err != nil {
		logger.ErrorErr(err, "GetAllBookings failed")
		return nil, errs.NewUnexpectedError("failed to retrieve all bookings")
	}
	logger.Debug("admin bookings returned", zap.Int("count", len(bookings)))
	return bookings, nil
}
