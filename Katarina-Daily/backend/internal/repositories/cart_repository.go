package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"fashion-ecommerce-backend/internal/models"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{db: db}
}

// GetOrCreateCart returns the user's single cart, creating it if this
// is their first visit. FirstOrCreate is used instead of a manual
// SELECT+CREATE pair to eliminate the race condition where two concurrent
// requests could both find "not found" and both attempt an INSERT,
// causing a unique-constraint violation on user_id.
func (r *CartRepository) GetOrCreateCart(ctx context.Context, userID uuid.UUID) (*models.Cart, error) {
	cart := models.Cart{UserID: userID}
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		FirstOrCreate(&cart).Error; err != nil {
		return nil, err
	}
	return &cart, nil
}

// FindCartByUserID returns the cart for a user without creating one if
// it does not exist. Returns gorm.ErrRecordNotFound when the user has no
// cart yet, which callers treat as "item cannot belong to this user".
func (r *CartRepository) FindCartByUserID(ctx context.Context, userID uuid.UUID) (*models.Cart, error) {
	var cart models.Cart
	if err := r.db.WithContext(ctx).First(&cart, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &cart, nil
}

// GetCartWithItems loads the cart plus every item's variant and parent
// product, so DTO mapping never needs a second round trip.
func (r *CartRepository) GetCartWithItems(ctx context.Context, cartID uuid.UUID) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.WithContext(ctx).
		Preload("Items", func(db *gorm.DB) *gorm.DB { return db.Order("created_at ASC") }).
		Preload("Items.ProductVariant").
		Preload("Items.ProductVariant.Product").
		First(&cart, "id = ?", cartID).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *CartRepository) FindItemByVariant(ctx context.Context, cartID, variantID uuid.UUID) (*models.CartItem, error) {
	var item models.CartItem
	err := r.db.WithContext(ctx).First(&item, "cart_id = ? AND product_variant_id = ?", cartID, variantID).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *CartRepository) FindItemByID(ctx context.Context, id uuid.UUID) (*models.CartItem, error) {
	var item models.CartItem
	err := r.db.WithContext(ctx).First(&item, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *CartRepository) CreateItem(ctx context.Context, item *models.CartItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *CartRepository) UpdateItem(ctx context.Context, item *models.CartItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *CartRepository) DeleteItem(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.CartItem{}, "id = ?", id).Error
}

func (r *CartRepository) ClearItems(ctx context.Context, cartID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.CartItem{}, "cart_id = ?", cartID).Error
}

// ClearItemsTx clears a cart's items using the supplied transaction -
// used right after an order is successfully created at checkout.
func (r *CartRepository) ClearItemsTx(ctx context.Context, tx *gorm.DB, cartID uuid.UUID) error {
	return tx.WithContext(ctx).Delete(&models.CartItem{}, "cart_id = ?", cartID).Error
}
