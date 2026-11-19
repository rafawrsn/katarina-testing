package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"fashion-ecommerce-backend/internal/models"
)

type AddressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) *AddressRepository {
	return &AddressRepository{db: db}
}

func (r *AddressRepository) Create(ctx context.Context, a *models.Address) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *AddressRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Address, error) {
	var a models.Address
	if err := r.db.WithContext(ctx).First(&a, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AddressRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]models.Address, error) {
	var items []models.Address
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("is_default DESC, created_at DESC").
		Find(&items).Error
	return items, err
}

func (r *AddressRepository) Update(ctx context.Context, a *models.Address) error {
	return r.db.WithContext(ctx).Save(a).Error
}

func (r *AddressRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Address{}, "id = ?", id).Error
}

// ClearDefaultForUser unsets is_default on every address for a user
// except the one being promoted, so exactly one default ever exists.
func (r *AddressRepository) ClearDefaultForUser(ctx context.Context, tx *gorm.DB, userID uuid.UUID, exceptID uuid.UUID) error {
	d := r.db
	if tx != nil {
		d = tx
	}
	return d.WithContext(ctx).Model(&models.Address{}).
		Where("user_id = ? AND id <> ?", userID, exceptID).
		Update("is_default", false).Error
}

func (r *AddressRepository) WithTx(tx *gorm.DB) *AddressRepository {
	return &AddressRepository{db: tx}
}
