package services

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"fashion-ecommerce-backend/internal/config"
	"fashion-ecommerce-backend/internal/dto"
	"fashion-ecommerce-backend/internal/models"
	"fashion-ecommerce-backend/internal/utils"
)

// userRepositoryStub is a Stub Test Double: it returns predetermined data.
type userRepositoryStub struct {
	user *models.User
	err  error
}

func (s *userRepositoryStub) Create(context.Context, *models.User) error { return nil }
func (s *userRepositoryStub) FindByID(context.Context, uuid.UUID) (*models.User, error) {
	return s.user, s.err
}
func (s *userRepositoryStub) FindByEmail(context.Context, string) (*models.User, error) {
	return s.user, s.err
}
func (s *userRepositoryStub) ExistsByEmail(context.Context, string) (bool, error) { return false, nil }
func (s *userRepositoryStub) Update(context.Context, *models.User) error          { return nil }
func (s *userRepositoryStub) List(context.Context, string, string, int, int) ([]models.User, int64, error) {
	return nil, 0, nil
}

// userRepositorySpy is a Spy/Mock-style Test Double: besides returning
// configured data, it records calls and arguments for interaction assertions.
type userRepositorySpy struct {
	user      *models.User
	err       error
	calls     int
	lastEmail string
}

func (s *userRepositorySpy) Create(context.Context, *models.User) error { return nil }
func (s *userRepositorySpy) FindByID(context.Context, uuid.UUID) (*models.User, error) {
	return s.user, s.err
}
func (s *userRepositorySpy) FindByEmail(_ context.Context, email string) (*models.User, error) {
	s.calls++
	s.lastEmail = email
	return s.user, s.err
}
func (s *userRepositorySpy) ExistsByEmail(context.Context, string) (bool, error) { return false, nil }
func (s *userRepositorySpy) Update(context.Context, *models.User) error          { return nil }
func (s *userRepositorySpy) List(context.Context, string, string, int, int) ([]models.User, int64, error) {
	return nil, 0, nil
}

func testAuthUser(t *testing.T) *models.User {
	t.Helper()
	hash, err := utils.HashPassword("Password123")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	return &models.User{BaseModel: models.BaseModel{ID: uuid.New()}, Name: "Test User", Email: "user@example.com", PasswordHash: hash, Role: models.RoleCustomer}
}

func TestAuthServiceLogin(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{JWTSecret: "unit-test-secret", JWTExpiryHrs: 1}

	t.Run("success with normalized email", func(t *testing.T) {
		repo := &userRepositorySpy{user: testAuthUser(t)}
		svc := NewAuthService(repo, cfg)
		got, err := svc.Login(ctx, dto.LoginRequest{Email: "  USER@EXAMPLE.COM ", Password: "Password123"})
		if err != nil {
			t.Fatalf("Login() error = %v", err)
		}
		if got == nil || got.Token == "" {
			t.Fatal("Login() did not return a JWT")
		}
		if got.User.Email != "user@example.com" {
			t.Errorf("user email = %q", got.User.Email)
		}
		if repo.calls != 1 {
			t.Errorf("FindByEmail calls = %d, want 1", repo.calls)
		}
		if repo.lastEmail != "user@example.com" {
			t.Errorf("repository email = %q, want normalized email", repo.lastEmail)
		}
	})

	t.Run("unknown email", func(t *testing.T) {
		repo := &userRepositoryStub{err: gorm.ErrRecordNotFound}
		svc := NewAuthService(repo, cfg)
		got, err := svc.Login(ctx, dto.LoginRequest{Email: "missing@example.com", Password: "Password123"})
		if got != nil {
			t.Errorf("response = %#v, want nil", got)
		}
		ae, ok := utils.AsAppError(err)
		if !ok || ae.Status != 401 || ae.Code != "INVALID_CREDENTIALS" {
			t.Errorf("error = %#v, want 401 INVALID_CREDENTIALS", err)
		}
	})

	t.Run("wrong password", func(t *testing.T) {
		repo := &userRepositoryStub{user: testAuthUser(t)}
		svc := NewAuthService(repo, cfg)
		got, err := svc.Login(ctx, dto.LoginRequest{Email: "user@example.com", Password: "wrong-password"})
		if got != nil {
			t.Errorf("response = %#v, want nil", got)
		}
		ae, ok := utils.AsAppError(err)
		if !ok || ae.Status != 401 || ae.Code != "INVALID_CREDENTIALS" {
			t.Errorf("error = %#v, want 401 INVALID_CREDENTIALS", err)
		}
	})

	t.Run("repository failure", func(t *testing.T) {
		repo := &userRepositoryStub{err: errors.New("database unavailable")}
		svc := NewAuthService(repo, cfg)
		got, err := svc.Login(ctx, dto.LoginRequest{Email: "user@example.com", Password: "Password123"})
		if got != nil {
			t.Errorf("response = %#v, want nil", got)
		}
		ae, ok := utils.AsAppError(err)
		if !ok || ae.Status != 500 || ae.Code != "INTERNAL_ERROR" {
			t.Errorf("error = %#v, want 500 INTERNAL_ERROR", err)
		}
	})
}
