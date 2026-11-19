package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"fashion-ecommerce-backend/internal/models"
)

type ReviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) Create(ctx context.Context, rev *models.Review) error {
	return r.db.WithContext(ctx).Create(rev).Error
}

func (r *ReviewRepository) ExistsForOrder(ctx context.Context, userID, productID, orderID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Review{}).
		Where("user_id = ? AND product_id = ? AND order_id = ?", userID, productID, orderID).
		Count(&count).Error
	return count > 0, err
}

func (r *ReviewRepository) ListByProduct(ctx context.Context, productID uuid.UUID, offset, limit int) ([]models.Review, int64, error) {
	q := r.db.WithContext(ctx).Model(&models.Review{}).Where("product_id = ?", productID)

	var total int64
	if err := q.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []models.Review
	err := q.Preload("User").Order("created_at DESC").Offset(offset).Limit(limit).Find(&items).Error
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *ReviewRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Review, error) {
	var rev models.Review
	if err := r.db.WithContext(ctx).First(&rev, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &rev, nil
}

func (r *ReviewRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Review{}, "id = ?", id).Error
}
