package utils

import (
	"fmt"
	"time"

	"github.com/ingwrok/hotelBooking/internal/common/errs"
)

const DateFormat = "2006-01-02"

func ParseDate(dateStr string,fieldName string) (time.Time, error) {

	if dateStr == "" {
		return time.Time{}, errs.NewValidationError(fmt.Sprintf("%s is required", fieldName))
	}

	t,err := time.Parse(DateFormat, dateStr)
	if err != nil {
		return time.Time{}, errs.NewValidationError(fmt.Sprintf("invalid %s format: expected %s", fieldName, DateFormat))
	}

	return t.Truncate(24*time.Hour), nil

}