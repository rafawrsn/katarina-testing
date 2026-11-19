package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"fashion-ecommerce-backend/internal/dto"
	"fashion-ecommerce-backend/internal/models"
	"fashion-ecommerce-backend/internal/repositories"
	"fashion-ecommerce-backend/internal/utils"
)

type ShipmentService struct {
	db        *gorm.DB
	orders    *repositories.OrderRepository
	shipments *repositories.ShipmentRepository
}

func NewShipmentService(db *gorm.DB, orders *repositories.OrderRepository, shipments *repositories.ShipmentRepository) *ShipmentService {
	return &ShipmentService{db: db, orders: orders, shipments: shipments}
}

func (s *ShipmentService) List(ctx context.Context, status string, pagination utils.PaginationParams) ([]dto.ShipmentResponse, utils.PaginationMeta, error) {
	items, total, err := s.shipments.List(ctx, status, pagination.Offset(), pagination.Limit)
	if err != nil {
		return nil, utils.PaginationMeta{}, utils.ErrInternal(err)
	}
	out := make([]dto.ShipmentResponse, 0, len(items))
	for i := range items {
		out = append(out, dto.ToShipmentResponse(&items[i]))
	}
	return out, utils.NewPaginationMeta(pagination, total), nil
}

// Create opens a shipment for an order that is ready to ship
// (PROCESSING). Providing a tracking number immediately also advances
// the order to SHIPPED, matching the required end-to-end flow.
func (s *ShipmentService) Create(ctx context.Context, orderID uuid.UUID, req dto.CreateShipmentRequest) (*dto.ShipmentResponse, error) {
	var shipment *models.Shipment

	err := s.db.Transaction(func(tx *gorm.DB) error {
		order, err := s.orders.FindByIDForUpdateTx(ctx, tx, orderID)
		if err != nil {
			return err
		}
		if order.Status != models.OrderStatusProcessing && order.Status != models.OrderStatusPaid {
			return utils.ErrConflict("ORDER_NOT_READY_TO_SHIP", "order must be PAID or PROCESSING before a shipment can be created")
		}

		var existing int64
		if err := tx.WithContext(ctx).Model(&models.Shipment{}).Where("order_id = ?", orderID).Count(&existing).Error; err != nil {
			return err
		}
		if existing > 0 {
			return utils.ErrConflict("SHIPMENT_ALREADY_EXISTS", "this order already has a shipment")
		}

		status := models.ShipmentStatusPending
		var shippedAt *time.Time
		if req.TrackingNumber != "" {
			status = models.ShipmentStatusShipped
			now := time.Now().UTC()
			shippedAt = &now
		}

		shipment = &models.Shipment{
			OrderID:        orderID,
			Courier:        req.Courier,
			Service:        req.Service,
			TrackingNumber: req.TrackingNumber,
			Status:         status,
			ShippedAt:      shippedAt,
		}
		if err := tx.WithContext(ctx).Create(shipment).Error; err != nil {
			return err
		}

		if status == models.ShipmentStatusShipped {
			order.Status = models.OrderStatusShipped
		} else {
			order.Status = models.OrderStatusProcessing
		}
		return s.orders.UpdateTx(ctx, tx, order)
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("order")
		}
		if ae, ok := utils.AsAppError(err); ok {
			return nil, ae
		}
		return nil, utils.ErrInternal(err)
	}

	resp := dto.ToShipmentResponse(shipment)
	return &resp, nil
}

// UpdateStatus updates a shipment's tracking number/status and keeps
// the parent order's status in sync (SHIPPED / DELIVERED).
func (s *ShipmentService) UpdateStatus(ctx context.Context, orderID uuid.UUID, req dto.UpdateShipmentRequest) (*dto.ShipmentResponse, error) {
	var shipment *models.Shipment

	err := s.db.Transaction(func(tx *gorm.DB) error {
		order, err := s.orders.FindByIDForUpdateTx(ctx, tx, orderID)
		if err != nil {
			return err
		}

		var sh models.Shipment
		if err := tx.WithContext(ctx).First(&sh, "order_id = ?", orderID).Error; err != nil {
			return err
		}

		if req.TrackingNumber != nil {
			sh.TrackingNumber = *req.TrackingNumber
		}
		if req.Status != nil {
			switch *req.Status {
			case models.ShipmentStatusPending, models.ShipmentStatusShipped, models.ShipmentStatusInTransit, models.ShipmentStatusDelivered, models.ShipmentStatusFailed:
				// valid status
			default:
				return utils.ErrValidation("invalid shipment status: " + *req.Status)
			}

			sh.Status = *req.Status
			now := time.Now().UTC()
			switch *req.Status {
			case models.ShipmentStatusShipped:
				sh.ShippedAt = &now
			case models.ShipmentStatusDelivered:
				sh.DeliveredAt = &now
			}
		}
		if err := tx.WithContext(ctx).Save(&sh).Error; err != nil {
			return err
		}
		shipment = &sh

		switch sh.Status {
		case models.ShipmentStatusShipped, models.ShipmentStatusInTransit:
			order.Status = models.OrderStatusShipped
		case models.ShipmentStatusDelivered:
			order.Status = models.OrderStatusDelivered
		}
		return s.orders.UpdateTx(ctx, tx, order)
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("shipment")
		}
		if ae, ok := utils.AsAppError(err); ok {
			return nil, ae
		}
		return nil, utils.ErrInternal(err)
	}

	resp := dto.ToShipmentResponse(shipment)
	return &resp, nil
}
