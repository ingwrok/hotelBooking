package email

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ingwrok/hotelBooking/internal/common/logger"
	"github.com/ingwrok/hotelBooking/internal/core/domain"
	"github.com/resend/resend-go/v2"
	"go.uber.org/zap"
)

type ResendAdapter struct {
	client *resend.Client
	from   string
}

func NewResendAdapter() *ResendAdapter {
	apiKey := os.Getenv("RESEND_API_KEY")
	from := os.Getenv("SMTP_FROM")

	if apiKey == "" {
		logger.Warn("RESEND_API_KEY missing. Emails will be LOGGED ONLY.")
		return &ResendAdapter{client: nil}
	}

	client := resend.NewClient(apiKey)
	return &ResendAdapter{
		client: client,
		from:   from,
	}
}

// ใช้ชื่อฟังก์ชันเดิมเป๊ะๆ เพื่อให้ BookingService เรียกใช้งานได้ทันที
func (a *ResendAdapter) SendBookingConfirmation(ctx context.Context, booking *domain.BookingDetail, addons []*domain.BookingAddon) error {
	// Construct Email Body
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Subject: Booking Confirmation #%d - %s\n\n", booking.BookingID, strings.ToUpper(booking.Status)))
	sb.WriteString(fmt.Sprintf("Dear %s,\n\nThank you for choosing our hotel!\n", booking.UserName))
	sb.WriteString(fmt.Sprintf("Here are your booking details:\n\n"))

	sb.WriteString(fmt.Sprintf("Booking ID:  #%d\n", booking.BookingID))
	sb.WriteString(fmt.Sprintf("Status:      %s\n", strings.ToUpper(booking.Status)))
	sb.WriteString(fmt.Sprintf("Room Type:   %s\n", booking.RoomTypeName))
	sb.WriteString(fmt.Sprintf("Rate Plan:   %s\n", booking.RatePlanName))
	sb.WriteString(fmt.Sprintf("Check-in:    %s\n", booking.CheckInDate.Format("02 Jan 2006")))
	sb.WriteString(fmt.Sprintf("Check-out:   %s\n", booking.CheckOutDate.Format("02 Jan 2006")))
	sb.WriteString(fmt.Sprintf("Guests:      %d Adults\n", booking.NumAdults))

	if booking.RoomNumber != "" {
		sb.WriteString(fmt.Sprintf("Room Number: %s\n", booking.RoomNumber))
	}
	sb.WriteString("\n----------------------------------------\n")
	sb.WriteString("PRICE BREAKDOWN\n")
	sb.WriteString("----------------------------------------\n")
	sb.WriteString(fmt.Sprintf("Room Charge:    THB %.2f\n", booking.RoomSubTotal))

	if len(addons) > 0 {
		sb.WriteString("Add-ons:\n")
		for _, ad := range addons {
			sb.WriteString(fmt.Sprintf("- %s (x%d): THB %.2f\n", ad.AddonName, ad.Quantity, ad.PriceAtBooking*float64(ad.Quantity)))
		}
		sb.WriteString(fmt.Sprintf("Addon Subtotal: THB %.2f\n", booking.AddonSubTotal))
	}

	sb.WriteString(fmt.Sprintf("Taxes (7%%):     THB %.2f\n", booking.TaxesAmount))
	sb.WriteString("----------------------------------------\n")
	sb.WriteString(fmt.Sprintf("TOTAL PRICE:    THB %.2f\n", booking.TotalPrice))
	sb.WriteString("----------------------------------------\n\n")

	sb.WriteString("We look forward to welcoming you!\n")

	body := sb.String()

	logger.Info("-------- EMAIL CONTENT START --------")
	fmt.Println(body)
	logger.Info("-------- EMAIL CONTENT END --------")
	if a.client == nil {
		return nil
	}

	recipient := booking.Email
	if recipient == "" {
		recipient = os.Getenv("SMTP_DEBUG_RECIPIENT")
	}

	params := &resend.SendEmailRequest{
		From:    a.from,
		To:      []string{recipient},
		Subject: fmt.Sprintf("Booking Confirmation #%d", booking.BookingID),
		Text:    body,
	}

	_, err := a.client.Emails.SendWithContext(ctx, params)
	if err != nil {
		logger.ErrorErr(err, "Failed to send email via Resend API")
		return nil
	}

	logger.Info("Email sent successfully via Resend", zap.String("to", recipient))
	return nil
}
