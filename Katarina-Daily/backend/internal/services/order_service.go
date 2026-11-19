package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"fashion-ecommerce-backend/internal/dto"
	"fashion-ecommerce-backend/internal/models"
	"fashion-ecommerce-backend/internal/repositories"
	"fashion-ecommerce-backend/internal/utils"
)

type OrderService struct {
	db       *gorm.DB
	orders   *repositories.OrderRepository
	variants *repositories.VariantRepository
}

func NewOrderService(db *gorm.DB, orders *repositories.OrderRepository, variants *repositories.VariantRepository) *OrderService {
	return &OrderService{db: db, orders: orders, variants: variants}
}

// allowedTransitions defines the legal order-status state machine.
// Terminal states map to an empty slice.
var allowedTransitions = map[string][]string{
	models.OrderStatusPendingPayment: {models.OrderStatusPaid, models.OrderStatusCancelled, models.OrderStatusExpired},
	models.OrderStatusPaid:           {models.OrderStatusProcessing, models.OrderStatusCancelled},
	models.OrderStatusProcessing:     {models.OrderStatusShipped, models.OrderStatusCancelled},
	models.OrderStatusShipped:        {models.OrderStatusDelivered},
	models.OrderStatusDelivered:      {},
	models.OrderStatusCancelled:      {},
	models.OrderStatusExpired:        {},
}

// stockDecremented reports whether stock was already deducted for an
// order in this status (true for PAID and subsequent active fulfillment states),
// so cancelling it later must restore stock.
func stockDecremented(status string) bool {
	switch status {
	case models.OrderStatusPaid, models.OrderStatusProcessing, models.OrderStatusShipped, models.OrderStatusDelivered:
		return true
	default:
		return false
	}
}

// ---- Customer ----

func (s *OrderService) ListByUser(ctx context.Context, userID uuid.UUID, status string, pagination utils.PaginationParams) ([]dto.OrderResponse, utils.PaginationMeta, error) {
	items, total, err := s.orders.ListByUser(ctx, userID, status, pagination.Offset(), pagination.Limit)
	if err != nil {
		return nil, utils.PaginationMeta{}, utils.ErrInternal(err)
	}
	return dto.ToOrderResponseList(items), utils.NewPaginationMeta(pagination, total), nil
}

func (s *OrderService) GetByIDForUser(ctx context.Context, userID, orderID uuid.UUID) (*dto.OrderResponse, error) {
	order, err := s.orders.FindByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("order")
		}
		return nil, utils.ErrInternal(err)
	}
	if order.UserID != userID {
		return nil, utils.ErrForbidden
	}
	resp := dto.ToOrderResponse(order)
	return &resp, nil
}

func (s *OrderService) Cancel(ctx context.Context, userID, orderID uuid.UUID) (*dto.OrderResponse, error) {
	var order *models.Order

	err := s.db.Transaction(func(tx *gorm.DB) error {
		var txErr error
		order, txErr = s.orders.FindByIDForUpdateTx(ctx, tx, orderID)
		if txErr != nil {
			return txErr
		}
		if order.UserID != userID {
			return utils.ErrForbidden
		}
		if !order.IsCancellable() {
			return utils.ErrConflict("ORDER_NOT_CANCELLABLE", "only orders awaiting payment can be cancelled by the customer")
		}

		// Only restore stock if the order status had actually decremented stock.
		if stockDecremented(order.Status) {
			for _, item := range order.Items {
				if err := s.variants.RestoreStock(ctx, tx, item.ProductVariantID, item.Quantity); err != nil {
					return err
				}
			}
		}

		order.Status = models.OrderStatusCancelled
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

	resp := dto.ToOrderResponse(order)
	return &resp, nil
}

// ---- Admin ----

func (s *OrderService) ListAdmin(ctx context.Context, status string, pagination utils.PaginationParams) ([]dto.OrderResponse, utils.PaginationMeta, error) {
	items, total, err := s.orders.ListAdmin(ctx, status, pagination.Offset(), pagination.Limit)
	if err != nil {
		return nil, utils.PaginationMeta{}, utils.ErrInternal(err)
	}
	return dto.ToOrderResponseList(items), utils.NewPaginationMeta(pagination, total), nil
}

func (s *OrderService) GetByIDAdmin(ctx context.Context, orderID uuid.UUID) (*dto.OrderResponse, error) {
	order, err := s.orders.FindByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("order")
		}
		return nil, utils.ErrInternal(err)
	}
	resp := dto.ToOrderResponse(order)
	return &resp, nil
}

func (s *OrderService) UpdateStatusAdmin(ctx context.Context, orderID uuid.UUID, newStatus string) (*dto.OrderResponse, error) {
	var order *models.Order

	err := s.db.Transaction(func(tx *gorm.DB) error {
		var txErr error
		order, txErr = s.orders.FindByIDForUpdateTx(ctx, tx, orderID)
		if txErr != nil {
			return txErr
		}

		if order.Status == newStatus {
			return nil // no-op, already in this state
		}

		allowed := allowedTransitions[order.Status]
		permitted := false
		for _, candidate := range allowed {
			if candidate == newStatus {
				permitted = true
				break
			}
		}
		if !permitted {
			return utils.ErrConflict("INVALID_STATUS_TRANSITION", "cannot move an order from "+order.Status+" to "+newStatus)
		}

		// If transitioning from PENDING_PAYMENT to PAID, deduct stock now.
		if order.Status == models.OrderStatusPendingPayment && newStatus == models.OrderStatusPaid {
			for _, item := range order.Items {
				affected, decErr := s.variants.DecrementStock(ctx, tx, item.ProductVariantID, item.Quantity)
				if decErr != nil {
					return decErr
				}
				if affected == 0 {
					return utils.ErrValidation(fmt.Sprintf(
						"insufficient stock for %s (%s/%s)", item.ProductName, item.Color, item.Size))
				}
			}
		}

		// Cancelling an order that had already had stock deducted
		// (PAID or later) must release that stock back to inventory.
		if newStatus == models.OrderStatusCancelled && stockDecremented(order.Status) {
			for _, item := range order.Items {
				if err := s.variants.RestoreStock(ctx, tx, item.ProductVariantID, item.Quantity); err != nil {
					return err
				}
			}
		}

		order.Status = newStatus
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

	resp := dto.ToOrderResponse(order)
	return &resp, nil
}
