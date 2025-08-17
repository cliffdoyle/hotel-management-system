package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cliffdoyle/internal/models"
	"github.com/cliffdoyle/internal/repository"
	"github.com/cliffdoyle/internal/tokens"
	validator "github.com/cliffdoyle/internal/validation"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// DTOs for data transfer
type UserRegisterDTO struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UserLoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AuthToken    *tokens.Token `json:"authentication_token"`
	RefreshToken string        `json:"refresh_token"`
}

type UserService interface {
	Register(ctx context.Context, dto UserRegisterDTO, hotelID uuid.UUID) (*models.User, error)
	Login(ctx context.Context, dto UserLoginDTO) (*LoginResponse, error)
}

type userService struct {
	userRepo repository.UserRepository
	redis    *redis.Client
}

func NewUserService(userRepo repository.UserRepository, redis *redis.Client) UserService {
	return &userService{
		userRepo: userRepo,
		redis:    redis,
	}
}

func (s *userService) Register(ctx context.Context, dto UserRegisterDTO, hotelID uuid.UUID) (*models.User, error) {
	// Validate DTO
	v := validator.New()
	validateRegisterDTO(v, dto)
	if !v.Valid() {
		return nil, &validator.ValidationError{Errors: v.Errors}
	}

	// Check for uniqueness
	_, err := s.userRepo.GetByEmail(ctx, dto.Email)
	if err != nil {
		if !errors.Is(err, repository.ErrRecordNotFound) {
			return nil, err // A real DB error occurred
		}
		// ErrRecordNotFound is the good case here, we can proceed.
	} else {
		// If err is nil, a user was found, so it's a duplicate.
		return nil, repository.ErrDuplicateEmail
	}

	user := &models.User{
		FirstName: &dto.FirstName,
		LastName:  &dto.LastName,
		Email:     dto.Email,
	}

	if err := user.SetPassword(dto.Password); err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	if err := s.userRepo.Create(ctx, user, hotelID); err != nil {
		return nil, fmt.Errorf("failed to create user in repo: %w", err)
	}

	return user, nil
}

func (s *userService) Login(ctx context.Context, dto UserLoginDTO) (*LoginResponse, error) {
	v := validator.New()
	validateLoginDTO(v, dto)
	if !v.Valid() {
		return nil, &validator.ValidationError{Errors: v.Errors}
	}

	user, err := s.userRepo.GetByEmail(ctx, dto.Email)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, repository.ErrInvalidCredentials
		}
		return nil, err
	}

	passwordMatches, err := user.MatchesPassword(dto.Password)
	if err != nil {
		return nil, err
	}
	if !passwordMatches {
		return nil, repository.ErrInvalidCredentials
	}

	authToken, err := tokens.GenerateToken(s.redis, user.ID, 30*time.Minute, tokens.ScopeAuthentication, user.Roles)
	if err != nil {
		return nil, fmt.Errorf("failed to generate auth token: %w", err)
	}

	refreshToken, err := tokens.GenerateToken(s.redis, user.ID, 24*time.Hour*30, tokens.ScopeRefreshToken, user.Roles)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &LoginResponse{
		AuthToken:    authToken,
		RefreshToken: refreshToken.Plaintext,
	}, nil
}

// Validation helpers are kept as private functions within the package
func validateRegisterDTO(v *validator.Validator, dto UserRegisterDTO) {
	v.Check(dto.FirstName != "", "first_name", "must be provided")
	v.Check(dto.LastName != "", "last_name", "must be provided")
	v.Check(dto.Email != "", "email", "must be provided")
	v.Check(validator.Matches(dto.Email, validator.EmailRX), "email", "must be a valid email address")
	v.Check(dto.Password != "", "password", "must be provided")
	v.Check(len(dto.Password) >= 8, "password", "must be at least 8 characters long")
}

func validateLoginDTO(v *validator.Validator, dto UserLoginDTO) {
	v.Check(dto.Email != "", "email", "must be provided")
	v.Check(validator.Matches(dto.Email, validator.EmailRX), "email", "must be a valid email address")
	v.Check(dto.Password != "", "password", "must be provided")
}
