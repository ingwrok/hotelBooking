package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ingwrok/hotelBooking/internal/common/logger"
	"github.com/ingwrok/hotelBooking/internal/core/domain"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

type GomailAdapter struct {
	dialer *gomail.Dialer
	from   string
}

func NewGomailAdapter() *GomailAdapter {
	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM")

	logger.Info("Initializing Email Adapter",
		zap.String("host", host),
		zap.String("port", portStr),
		zap.String("user", username),
	)

	if host == "" || portStr == "" || username == "" || password == "" {
		logger.Warn("SMTP configuration missing. Emails will be LOGGED ONLY.")
		return &GomailAdapter{dialer: nil}
	}

	port, _ := strconv.Atoi(portStr)
	d := gomail.NewDialer(host, port, username, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true} // Simplify for dev

	return &GomailAdapter{
		dialer: d,
		from:   from, 
	}
}

func (a *GomailAdapter) SendBookingConfirmation(ctx context.Context, booking *domain.BookingDetail, addons []*domain.BookingAddon) error {
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

	// Always log for debugging
	logger.Info("-------- EMAIL CONTENT START --------")
	fmt.Println(body)
	logger.Info("-------- EMAIL CONTENT END --------")

	// If no dialer, we are done
	if a.dialer == nil {
		return nil
	}

	// Send Real Email
	m := gomail.NewMessage()
	recipient := booking.Email
	if recipient == "" {
		recipient = os.Getenv("SMTP_DEBUG_RECIPIENT")
		if recipient == "" {
			logger.Warn("No recipient email found (User email lookup not implemented). Skipping actual send.")
			return nil
		}
	}

	m.SetHeader("From", a.from)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", fmt.Sprintf("Booking Confirmation #%d", booking.BookingID))
	m.SetBody("text/plain", body)

	if err := a.dialer.DialAndSend(m); err != nil {
		logger.ErrorErr(err, "Failed to send email via SMTP")
		return err
	}

	logger.Info("Email sent successfully", zap.String("to", recipient))
	return nil
}
