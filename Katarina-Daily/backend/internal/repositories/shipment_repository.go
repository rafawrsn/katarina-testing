package repositories

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"fashion-ecommerce-backend/internal/models"
)

type ShipmentRepository struct {
	db *gorm.DB
}

func NewShipmentRepository(db *gorm.DB) *ShipmentRepository {
	return &ShipmentRepository{db: db}
}

func (r *ShipmentRepository) Create(ctx context.Context, s *models.Shipment) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *ShipmentRepository) FindByOrderID(ctx context.Context, orderID uuid.UUID) (*models.Shipment, error) {
	var s models.Shipment
	err := r.db.WithContext(ctx).First(&s, "order_id = ?", orderID).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *ShipmentRepository) Exists(ctx context.Context, orderID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Shipment{}).Where("order_id = ?", orderID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ShipmentRepository) Update(ctx context.Context, s *models.Shipment) error {
	return r.db.WithContext(ctx).Save(s).Error
}

func (r *ShipmentRepository) List(ctx context.Context, status string, offset, limit int) ([]models.Shipment, int64, error) {
	q := r.db.WithContext(ctx).Model(&models.Shipment{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	var total int64
	if err := q.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []models.Shipment
	err := q.Order("created_at DESC").Offset(offset).Limit(limit).Find(&items).Error
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

var ErrShipmentNotFound = errors.New("shipment not found")
