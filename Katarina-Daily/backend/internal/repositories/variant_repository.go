package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"fashion-ecommerce-backend/internal/models"
)

type VariantRepository struct {
	db *gorm.DB
}

func NewVariantRepository(db *gorm.DB) *VariantRepository {
	return &VariantRepository{db: db}
}

func (r *VariantRepository) Create(ctx context.Context, v *models.ProductVariant) error {
	return r.db.WithContext(ctx).Create(v).Error
}

func (r *VariantRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.ProductVariant, error) {
	var v models.ProductVariant
	if err := r.db.WithContext(ctx).First(&v, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

// FindByIDForUpdate locks the variant row (SELECT ... FOR UPDATE) within
// tx, so concurrent checkouts/webhooks can't both decrement the same
// stock past zero. Must be called inside a transaction.
func (r *VariantRepository) FindByIDForUpdate(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*models.ProductVariant, error) {
	var v models.ProductVariant
	err := tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&v, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *VariantRepository) ExistsBySKU(ctx context.Context, sku string, excludeID *uuid.UUID) (bool, error) {
	q := r.db.WithContext(ctx).Model(&models.ProductVariant{}).Where("sku = ?", sku)
	if excludeID != nil {
		q = q.Where("id <> ?", *excludeID)
	}
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *VariantRepository) ListByProduct(ctx context.Context, productID uuid.UUID) ([]models.ProductVariant, error) {
	var items []models.ProductVariant
	err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("color, size").
		Find(&items).Error
	return items, err
}

func (r *VariantRepository) Update(ctx context.Context, v *models.ProductVariant) error {
	return r.db.WithContext(ctx).Save(v).Error
}

// UpdateWithTx saves a variant using the given transaction handle.
func (r *VariantRepository) UpdateWithTx(ctx context.Context, tx *gorm.DB, v *models.ProductVariant) error {
	return tx.WithContext(ctx).Save(v).Error
}

func (r *VariantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.ProductVariant{}, "id = ?", id).Error
}

// DecrementStock atomically reduces stock by quantity, guarded so it can
// never go negative (WHERE stock >= quantity). Returns the number of
// rows affected: 0 means insufficient stock (caller must treat as a
// failure and roll back the enclosing transaction).
func (r *VariantRepository) DecrementStock(ctx context.Context, tx *gorm.DB, variantID uuid.UUID, quantity int) (int64, error) {
	result := tx.WithContext(ctx).Model(&models.ProductVariant{}).
		Where("id = ? AND stock >= ?", variantID, quantity).
		UpdateColumn("stock", gorm.Expr("stock - ?", quantity))
	return result.RowsAffected, result.Error
}

// RestoreStock increments stock back (e.g. on order cancellation after
// stock had already been deducted).
func (r *VariantRepository) RestoreStock(ctx context.Context, tx *gorm.DB, variantID uuid.UUID, quantity int) error {
	return tx.WithContext(ctx).Model(&models.ProductVariant{}).
		Where("id = ?", variantID).
		UpdateColumn("stock", gorm.Expr("stock + ?", quantity)).Error
}
