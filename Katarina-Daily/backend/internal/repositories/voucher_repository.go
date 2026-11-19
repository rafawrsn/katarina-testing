package repositories

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"fashion-ecommerce-backend/internal/models"
)

type VoucherRepository struct {
	db *gorm.DB
}

func NewVoucherRepository(db *gorm.DB) *VoucherRepository {
	return &VoucherRepository{db: db}
}

func (r *VoucherRepository) FindByCode(ctx context.Context, code string) (*models.Voucher, error) {
	var v models.Voucher
	if err := r.db.WithContext(ctx).First(&v, "code = ?", code).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

// FindByCodeForUpdateTx locks the voucher row so concurrent checkouts
// can't both push used_count past usage_limit.
func (r *VoucherRepository) FindByCodeForUpdateTx(ctx context.Context, tx *gorm.DB, code string) (*models.Voucher, error) {
	var v models.Voucher
	err := tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&v, "code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *VoucherRepository) IncrementUsageTx(ctx context.Context, tx *gorm.DB, id interface{}) error {
	return tx.WithContext(ctx).Model(&models.Voucher{}).
		Where("id = ?", id).
		UpdateColumn("used_count", gorm.Expr("used_count + 1")).Error
}

func (r *VoucherRepository) Create(ctx context.Context, v *models.Voucher) error {
	return r.db.WithContext(ctx).Create(v).Error
}

func (r *VoucherRepository) List(ctx context.Context) ([]models.Voucher, error) {
	var items []models.Voucher
	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&items).Error
	return items, err
}
