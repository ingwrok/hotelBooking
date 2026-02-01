package model

import "github.com/ingwrok/hotelBooking/internal/core/domain"

type User struct {
	UserID       int    `db:"user_id"`
	Username     string `db:"username"`
	Email        string `db:"email"`
	PasswordHash string `db:"password_hash"`
	IsAdmin      bool   `db:"is_admin"`
}

func (m *User) ToDomain() *domain.User {
	return &domain.User{
		UserID:       m.UserID,
		Username:     m.Username,
		Email:        m.Email,
		PasswordHash: m.PasswordHash,
		IsAdmin:      m.IsAdmin,
	}
}

func FromDomain(d *domain.User) *User {
	return &User{
		UserID:       d.UserID,
		Username:     d.Username,
		Email:        d.Email,
		PasswordHash: d.PasswordHash,
		IsAdmin:      d.IsAdmin,
	}
}