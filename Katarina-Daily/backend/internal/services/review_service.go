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

type ReviewService struct {
	reviews  *repositories.ReviewRepository
	products *repositories.ProductRepository
	orders   *repositories.OrderRepository
}

func NewReviewService(reviews *repositories.ReviewRepository, products *repositories.ProductRepository, orders *repositories.OrderRepository) *ReviewService {
	return &ReviewService{reviews: reviews, products: products, orders: orders}
}

func (s *ReviewService) ListByProduct(ctx context.Context, productID uuid.UUID, pagination utils.PaginationParams) ([]dto.ReviewResponse, utils.PaginationMeta, error) {
	items, total, err := s.reviews.ListByProduct(ctx, productID, pagination.Offset(), pagination.Limit)
	if err != nil {
		return nil, utils.PaginationMeta{}, utils.ErrInternal(err)
	}
	return dto.ToReviewResponseList(items), utils.NewPaginationMeta(pagination, total), nil
}

// Create enforces the spec's rule that only a customer who purchased
// the product may review it: the referenced order must belong to the
// customer, actually contain this product, and have reached a
// paid-or-later status. One review per (user, product, order).
func (s *ReviewService) Create(ctx context.Context, userID, productID uuid.UUID, req dto.CreateReviewRequest) (*dto.ReviewResponse, error) {
	if _, err := s.products.FindByID(ctx, productID, false); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("product")
		}
		return nil, utils.ErrInternal(err)
	}

	purchased, err := s.orders.HasPurchasedProductInOrder(ctx, userID, productID, req.OrderID)
	if err != nil {
		return nil, utils.ErrInternal(err)
	}
	if !purchased {
		return nil, utils.NewAppError(403, "PURCHASE_NOT_VERIFIED", "you can only review a product from an order that contains it and has been paid")
	}

	alreadyReviewed, err := s.reviews.ExistsForOrder(ctx, userID, productID, req.OrderID)
	if err != nil {
		return nil, utils.ErrInternal(err)
	}
	if alreadyReviewed {
		return nil, utils.ErrConflict("REVIEW_ALREADY_EXISTS", "you have already reviewed this product for this order")
	}

	rev := &models.Review{
		UserID:    userID,
		ProductID: productID,
		OrderID:   req.OrderID,
		Rating:    req.Rating,
		Comment:   req.Comment,
	}
	if err := s.reviews.Create(ctx, rev); err != nil {
		return nil, utils.ErrInternal(err)
	}

	resp := dto.ToReviewResponse(rev)
	return &resp, nil
}

// DeleteOwn lets a customer delete their own review.
func (s *ReviewService) DeleteOwn(ctx context.Context, userID, reviewID uuid.UUID) error {
	rev, err := s.reviews.FindByID(ctx, reviewID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrNotFound("review")
		}
		return utils.ErrInternal(err)
	}
	if rev.UserID != userID {
		return utils.ErrForbidden
	}
	if err := s.reviews.Delete(ctx, reviewID); err != nil {
		return utils.ErrInternal(err)
	}
	return nil
}

// DeleteAdmin lets an admin moderate/remove any review.
func (s *ReviewService) DeleteAdmin(ctx context.Context, reviewID uuid.UUID) error {
	if _, err := s.reviews.FindByID(ctx, reviewID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrNotFound("review")
		}
		return utils.ErrInternal(err)
	}
	if err := s.reviews.Delete(ctx, reviewID); err != nil {
		return utils.ErrInternal(err)
	}
	return nil
}
