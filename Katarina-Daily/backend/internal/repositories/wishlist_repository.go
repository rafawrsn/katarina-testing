package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"fashion-ecommerce-backend/internal/models"
)

type WishlistRepository struct {
	db *gorm.DB
}

func NewWishlistRepository(db *gorm.DB) *WishlistRepository {
	return &WishlistRepository{db: db}
}

func (r *WishlistRepository) Create(ctx context.Context, w *models.Wishlist) error {
	return r.db.WithContext(ctx).Create(w).Error
}

func (r *WishlistRepository) Exists(ctx context.Context, userID, productID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Wishlist{}).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Count(&count).Error
	return count > 0, err
}

func (r *WishlistRepository) Delete(ctx context.Context, userID, productID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Delete(&models.Wishlist{}).Error
}

func (r *WishlistRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]models.Wishlist, error) {
	var items []models.Wishlist
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Product.Category").
		Preload("Product.Variants", func(db *gorm.DB) *gorm.DB { return db.Order("color, size") }).
		Preload("Product.Images", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order") }).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&items).Error
	return items, err
}
