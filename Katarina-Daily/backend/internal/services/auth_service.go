package services

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"fashion-ecommerce-backend/internal/config"
	"fashion-ecommerce-backend/internal/dto"
	"fashion-ecommerce-backend/internal/models"
	"fashion-ecommerce-backend/internal/utils"
)

// AuthUserRepository describes the persistence operations required by AuthService.
// Keeping the dependency behind an interface allows isolated service-level tests.
type AuthUserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Update(ctx context.Context, user *models.User) error
	List(ctx context.Context, role, search string, offset, limit int) ([]models.User, int64, error)
}

type AuthService struct {
	users AuthUserRepository
	cfg   *config.Config
}

func NewAuthService(users AuthUserRepository, cfg *config.Config) *AuthService {
	return &AuthService{users: users, cfg: cfg}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	exists, err := s.users.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, utils.ErrInternal(err)
	}
	if exists {
		return nil, utils.ErrConflict("EMAIL_ALREADY_REGISTERED", "an account with this email already exists")
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, utils.ErrInternal(err)
	}

	user := &models.User{
		Name:         strings.TrimSpace(req.Name),
		Email:        email,
		PasswordHash: hash,
		Phone:        req.Phone,
		Role:         models.RoleCustomer,
	}
	if err := s.users.Create(ctx, user); err != nil {
		return nil, utils.ErrInternal(err)
	}

	return s.issueToken(user)
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	user, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewAppError(http.StatusUnauthorized, "INVALID_CREDENTIALS", "invalid email or password")
		}
		return nil, utils.ErrInternal(err)
	}

	if !utils.CheckPassword(user.PasswordHash, req.Password) {
		return nil, utils.NewAppError(http.StatusUnauthorized, "INVALID_CREDENTIALS", "invalid email or password")
	}

	return s.issueToken(user)
}

// Me returns the currently authenticated user's public profile.
func (s *AuthService) Me(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("user")
		}
		return nil, utils.ErrInternal(err)
	}
	resp := dto.ToUserResponse(user)
	return &resp, nil
}

func (s *AuthService) issueToken(user *models.User) (*dto.AuthResponse, error) {
	token, err := utils.GenerateToken(s.cfg.JWTSecret, s.cfg.JWTExpiryHrs, user.ID, user.Email, user.Role)
	if err != nil {
		return nil, utils.ErrInternal(err)
	}
	return &dto.AuthResponse{
		Token: token,
		User:  dto.ToUserResponse(user),
	}, nil
}
