package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ingwrok/hotelBooking/internal/adapters/primary/web/middleware"
	"github.com/ingwrok/hotelBooking/internal/common/errs"
	"github.com/ingwrok/hotelBooking/internal/common/logger"
	"github.com/ingwrok/hotelBooking/internal/core/domain"
	"github.com/ingwrok/hotelBooking/internal/core/ports"
	"golang.org/x/crypto/bcrypt"

	"github.com/spf13/viper"
)

type UserService struct {
	repo ports.UserRepoPort
}

func NewUserService(repo ports.UserRepoPort) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, username, email, password string) (*domain.User, error) {
	if username == "" || email == "" || password == "" {
		return nil, errs.NewValidationError("all fields are required")
	}

	existing, _ := s.repo.GetByUsername(ctx, username)
	if existing != nil {
		return nil, errs.NewValidationError("username already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.ErrorErr(err, "failed to hash password")
		return nil, errs.NewUnexpectedError("internal server error")
	}

	newUser := &domain.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
		IsAdmin:      false,
	}

	if err := s.repo.Create(ctx, newUser); err != nil {
		logger.ErrorErr(err, "repo.Create failed in Register")
		return nil, errs.NewUnexpectedError("failed to create user")
	}
	newUser.PasswordHash = ""
	return newUser, nil
}

func (s *UserService) Login(ctx context.Context, username string, password string) (string, *domain.User, error) {
	u, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil, errs.NewUnauthorizedError("invalid username or password")
		}
		logger.ErrorErr(err, "repo.GetByUsername failed in Login")
		return "", nil, errs.NewUnexpectedError("internal server error")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", nil, errs.NewUnauthorizedError("invalid username or password")
	}

	sec := viper.GetString("secret")
	if sec == "" {
		logger.Error("jwt secret missing")
		return "", nil, errs.NewUnexpectedError("internal server error")
	}

	claims := &middleware.MyCustomClaims{
		UserID:  u.UserID,
		IsAdmin: u.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", u.UserID),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    "hotel-booking",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(sec))
	if err != nil {
		return "", nil, errs.NewUnexpectedError("failed to generate token")
	}

	u.PasswordHash = ""
	return tokenString, u, nil
}

func (s *UserService) UpdateUser(ctx context.Context, userID int, fields map[string]interface{}) error {
	if userID == 0 {
		return errs.NewValidationError("user_id is required")
	}
	if len(fields) == 0 {
		return nil
	}

	if v, ok := fields["password"]; ok {
		if raw, ok2 := v.(string); ok2 && raw != "" {
			hash, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
			if err != nil {
				logger.ErrorErr(err, "failed to hash password during update")
				return errs.NewUnexpectedError("internal server error")
			}
			fields["password_hash"] = string(hash)
			delete(fields, "password")
		} else {
			delete(fields, "password")
		}
	}

	delete(fields, "is_admin")

	if err := s.repo.Update(ctx, userID, fields); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errs.NewNotFoundError("user not found")
		}
		logger.ErrorErr(err, "repo.Update failed")
		return errs.NewUnexpectedError("internal server error")
	}
	return nil
}

func (s *UserService) DeleteUser(ctx context.Context, userID int) error {
	if err := s.repo.Delete(ctx, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errs.NewNotFoundError("user not found")
		}
		logger.ErrorErr(err, "repo.Delete failed")
		return errs.NewUnexpectedError("internal server error")
	}
	return nil
}

func (s *UserService) GetUser(ctx context.Context, userID int) (*domain.User, error) {
	u, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("user not found")
		}
		logger.ErrorErr(err, "repo.GetByID failed")
		return nil, errs.NewUnexpectedError("internal server error")
	}

	u.PasswordHash = ""
	return u, nil
}

func (s *UserService) GetUsers(ctx context.Context) ([]*domain.User, error) {
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		logger.ErrorErr(err, "repo.GetAll failed")
		return nil, errs.NewUnexpectedError("internal server error")
	}

	for i := range users {
		users[i].PasswordHash = ""
	}
	return users, nil
}
