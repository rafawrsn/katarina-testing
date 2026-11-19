package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"fashion-ecommerce-backend/internal/models"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) CreateTx(ctx context.Context, tx *gorm.DB, p *models.Payment) error {
	return tx.WithContext(ctx).Create(p).Error
}

func (r *PaymentRepository) FindByOrderID(ctx context.Context, orderID uuid.UUID) (*models.Payment, error) {
	var p models.Payment
	if err := r.db.WithContext(ctx).First(&p, "order_id = ?", orderID).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

// FindByOrderIDForUpdateTx locks the payment row for update - the
// cornerstone of webhook idempotency: whoever gets the lock first
// decides whether this notification has already been processed.
func (r *PaymentRepository) FindByOrderIDForUpdateTx(ctx context.Context, tx *gorm.DB, orderID uuid.UUID) (*models.Payment, error) {
	var p models.Payment
	err := tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&p, "order_id = ?", orderID).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PaymentRepository) UpdateTx(ctx context.Context, tx *gorm.DB, p *models.Payment) error {
	return tx.WithContext(ctx).Save(p).Error
}

func (r *PaymentRepository) Update(ctx context.Context, p *models.Payment) error {
	return r.db.WithContext(ctx).Save(p).Error
}
