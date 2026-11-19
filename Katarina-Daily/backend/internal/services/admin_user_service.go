package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"fashion-ecommerce-backend/internal/dto"
	"fashion-ecommerce-backend/internal/repositories"
	"fashion-ecommerce-backend/internal/utils"
)

type AdminUserService struct {
	users *repositories.UserRepository
}

func NewAdminUserService(users *repositories.UserRepository) *AdminUserService {
	return &AdminUserService{users: users}
}

func (s *AdminUserService) List(ctx context.Context, role, search string, pagination utils.PaginationParams) ([]dto.UserResponse, utils.PaginationMeta, error) {
	items, total, err := s.users.List(ctx, role, search, pagination.Offset(), pagination.Limit)
	if err != nil {
		return nil, utils.PaginationMeta{}, utils.ErrInternal(err)
	}
	return dto.ToUserResponseList(items), utils.NewPaginationMeta(pagination, total), nil
}

func (s *AdminUserService) GetByID(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
	u, err := s.users.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("user")
		}
		return nil, utils.ErrInternal(err)
	}
	resp := dto.ToUserResponse(u)
	return &resp, nil
}
