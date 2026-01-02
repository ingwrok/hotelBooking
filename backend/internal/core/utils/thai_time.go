package utils

import "time"

func ToThaiTime(t time.Time) time.Time {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return t
	}
	return t.In(loc)
}
