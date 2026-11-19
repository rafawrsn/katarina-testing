package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"fashion-ecommerce-backend/internal/models"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// DB exposes the underlying *gorm.DB so the checkout service can open
// its own transaction spanning orders + order_items + cart clearing.
func (r *OrderRepository) DB() *gorm.DB { return r.db }

func (r *OrderRepository) CreateTx(ctx context.Context, tx *gorm.DB, order *models.Order) error {
	return tx.WithContext(ctx).Create(order).Error
}

func (r *OrderRepository) CreateItemsTx(ctx context.Context, tx *gorm.DB, items []models.OrderItem) error {
	if len(items) == 0 {
		return nil
	}
	return tx.WithContext(ctx).Create(&items).Error
}

func (r *OrderRepository) ExistsByOrderNumber(ctx context.Context, orderNumber string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Order{}).Where("order_number = ?", orderNumber).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *OrderRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	var o models.Order
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Payment").
		Preload("Shipment").
		First(&o, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *OrderRepository) FindByOrderNumber(ctx context.Context, orderNumber string) (*models.Order, error) {
	var o models.Order
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Payment").
		Preload("Shipment").
		First(&o, "order_number = ?", orderNumber).Error
	if err != nil {
		return nil, err
	}
	return &o, nil
}

// FindByIDForUpdateTx locks the order row for update within tx - used
// during payment/status transitions to serialize concurrent updates
// (e.g. a webhook and an admin status change racing).
func (r *OrderRepository) FindByIDForUpdateTx(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*models.Order, error) {
	var o models.Order
	err := tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("Items").
		First(&o, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *OrderRepository) FindByOrderNumberForUpdateTx(ctx context.Context, tx *gorm.DB, orderNumber string) (*models.Order, error) {
	var o models.Order
	err := tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("Items").
		First(&o, "order_number = ?", orderNumber).Error
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *OrderRepository) UpdateTx(ctx context.Context, tx *gorm.DB, order *models.Order) error {
	return tx.WithContext(ctx).Save(order).Error
}

func (r *OrderRepository) Update(ctx context.Context, order *models.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}

func (r *OrderRepository) ListByUser(ctx context.Context, userID uuid.UUID, status string, offset, limit int) ([]models.Order, int64, error) {
	q := r.db.WithContext(ctx).Model(&models.Order{}).Where("user_id = ?", userID)
	if status != "" {
		q = q.Where("status = ?", status)
	}

	var total int64
	if err := q.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var orders []models.Order
	err := q.Preload("Items").Preload("Payment").Preload("Shipment").
		Order("created_at DESC").Offset(offset).Limit(limit).
		Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}

func (r *OrderRepository) ListAdmin(ctx context.Context, status string, offset, limit int) ([]models.Order, int64, error) {
	q := r.db.WithContext(ctx).Model(&models.Order{})
	if status != "" {
		q = q.Where("status = ?", status)
	}

	var total int64
	if err := q.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var orders []models.Order
	err := q.Preload("Items").Preload("Payment").Preload("Shipment").
		Order("created_at DESC").Offset(offset).Limit(limit).
		Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}

// HasPurchasedProduct checks whether userID has any DELIVERED (or PAID+)
// order containing productID - used to gate review submission.
func (r *OrderRepository) HasPurchasedProduct(ctx context.Context, userID, productID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("order_items").
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Joins("JOIN product_variants ON product_variants.id = order_items.product_variant_id").
		Where("orders.user_id = ? AND product_variants.product_id = ? AND orders.status IN ?",
			userID, productID, []string{models.OrderStatusPaid, models.OrderStatusProcessing, models.OrderStatusShipped, models.OrderStatusDelivered}).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// HasPurchasedProductInOrder is the same check as HasPurchasedProduct
// but scoped to one specific order - a review references a single
// order_id, so this confirms that exact order (a) belongs to the
// customer, (b) actually contains the product, and (c) has reached a
// paid-or-later status, before a review is allowed to reference it.
func (r *OrderRepository) HasPurchasedProductInOrder(ctx context.Context, userID, productID, orderID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("order_items").
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Joins("JOIN product_variants ON product_variants.id = order_items.product_variant_id").
		Where("orders.id = ? AND orders.user_id = ? AND product_variants.product_id = ? AND orders.status IN ?",
			orderID, userID, productID, []string{models.OrderStatusPaid, models.OrderStatusProcessing, models.OrderStatusShipped, models.OrderStatusDelivered}).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
