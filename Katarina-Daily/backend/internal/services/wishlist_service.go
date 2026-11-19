package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"fashion-ecommerce-backend/internal/dto"
	"fashion-ecommerce-backend/internal/models"
	"fashion-ecommerce-backend/internal/repositories"
	"fashion-ecommerce-backend/internal/utils"
)

type WishlistService struct {
	wishlists *repositories.WishlistRepository
	products  *repositories.ProductRepository
}

func NewWishlistService(wishlists *repositories.WishlistRepository, products *repositories.ProductRepository) *WishlistService {
	return &WishlistService{wishlists: wishlists, products: products}
}

func (s *WishlistService) List(ctx context.Context, userID uuid.UUID) ([]dto.WishlistItemResponse, error) {
	items, err := s.wishlists.ListByUser(ctx, userID)
	if err != nil {
		return nil, utils.ErrInternal(err)
	}

	out := make([]dto.WishlistItemResponse, 0, len(items))
	for _, w := range items {
		out = append(out, dto.WishlistItemResponse{
			ID:        w.ID,
			ProductID: w.ProductID,
			Product:   dto.ToProductResponse(&w.Product),
			CreatedAt: w.CreatedAt,
		})
	}
	return out, nil
}

func (s *WishlistService) Add(ctx context.Context, userID uuid.UUID, req dto.AddWishlistRequest) (*dto.WishlistItemResponse, error) {
	product, err := s.products.FindByID(ctx, req.ProductID, true)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("product")
		}
		return nil, utils.ErrInternal(err)
	}

	exists, err := s.wishlists.Exists(ctx, userID, req.ProductID)
	if err != nil {
		return nil, utils.ErrInternal(err)
	}
	if exists {
		return nil, utils.ErrConflict("ALREADY_IN_WISHLIST", "this product is already in your wishlist")
	}

	w := &models.Wishlist{UserID: userID, ProductID: req.ProductID}
	if err := s.wishlists.Create(ctx, w); err != nil {
		return nil, utils.ErrInternal(err)
	}

	return &dto.WishlistItemResponse{
		ID:        w.ID,
		ProductID: w.ProductID,
		Product:   dto.ToProductResponse(product),
		CreatedAt: w.CreatedAt,
	}, nil
}

func (s *WishlistService) Remove(ctx context.Context, userID, productID uuid.UUID) error {
	exists, err := s.wishlists.Exists(ctx, userID, productID)
	if err != nil {
		return utils.ErrInternal(err)
	}
	if !exists {
		return utils.ErrNotFound("wishlist_item")
	}
	if err := s.wishlists.Delete(ctx, userID, productID); err != nil {
		return utils.ErrInternal(err)
	}
	return nil
}
