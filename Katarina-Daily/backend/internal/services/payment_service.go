package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"fashion-ecommerce-backend/internal/models"
	"fashion-ecommerce-backend/internal/repositories"
	"fashion-ecommerce-backend/internal/utils"
	"fashion-ecommerce-backend/pkg/payment"
)

type PaymentService struct {
	db       *gorm.DB
	orders   *repositories.OrderRepository
	payments *repositories.PaymentRepository
	variants *repositories.VariantRepository
	gateway  payment.Gateway
}

func NewPaymentService(db *gorm.DB, orders *repositories.OrderRepository, payments *repositories.PaymentRepository, variants *repositories.VariantRepository, gateway payment.Gateway) *PaymentService {
	return &PaymentService{db: db, orders: orders, payments: payments, variants: variants, gateway: gateway}
}

// HandleWebhook verifies the gateway notification's authenticity, then
// applies it idempotently: a duplicate delivery of an already-processed
// successful payment is a safe no-op (no double order, no double stock
// decrement, no duplicate payment record).
func (s *PaymentService) HandleWebhook(ctx context.Context, rawPayload []byte, headers map[string]string) error {
	result, err := s.gateway.HandleWebhook(ctx, rawPayload, headers)
	if err != nil {
		if errors.Is(err, payment.ErrInvalidSignature) {
			return utils.NewAppError(401, "INVALID_SIGNATURE", "webhook signature verification failed")
		}
		return utils.ErrInternal(err)
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		// result.OrderID actually carries our human-readable
		// order_number (that's the field Midtrans' notification calls
		// "order_id" - see MidtransGateway.CreatePayment).
		order, oErr := s.orders.FindByOrderNumberForUpdateTx(ctx, tx, result.OrderID)
		if oErr != nil {
			if errors.Is(oErr, gorm.ErrRecordNotFound) {
				return payment.ErrUnknownOrder
			}
			return oErr
		}

		pay, pErr := s.payments.FindByOrderIDForUpdateTx(ctx, tx, order.ID)
		if pErr != nil {
			return pErr
		}

		// --- Idempotency guard ---
		// A payment that has already reached PAID must never be
		// reprocessed: no re-decrementing stock, no re-transitioning
		// the order. Gateways routinely redeliver the same
		// notification.
		if pay.Status == models.PaymentStatusPaid {
			slog.Info("webhook ignored: payment already PAID", "order_number", order.OrderNumber)
			return nil
		}

		pay.Status = result.Status
		pay.PaymentMethod = result.PaymentMethod
		if result.TransactionID != "" {
			txID := result.TransactionID
			pay.TransactionID = &txID
		}
		pay.RawResponse = datatypes.JSON(result.RawResponse)

		now := time.Now().UTC()
		switch result.Status {
		case models.PaymentStatusPaid:
			pay.PaidAt = &now
		case models.PaymentStatusExpired:
			pay.ExpiredAt = &now
		}

		if err := s.payments.UpdateTx(ctx, tx, pay); err != nil {
			return err
		}

		// Only PENDING_PAYMENT orders should be moved by a webhook;
		// later states are owned by admin/shipment flows.
		if order.Status != models.OrderStatusPendingPayment {
			return nil
		}

		// Guard: an order should always have items. If items somehow
		// failed to preload we must not silently succeed.
		if len(order.Items) == 0 {
			return fmt.Errorf("order %s has no items loaded; cannot update stock", order.OrderNumber)
		}

		switch result.Status {
		case models.PaymentStatusPaid:
			// Payment is now confirmed/successful: deduct stock for every line item atomically.
			for _, item := range order.Items {
				affected, decErr := s.variants.DecrementStock(ctx, tx, item.ProductVariantID, item.Quantity)
				if decErr != nil {
					return decErr
				}
				if affected == 0 {
					slog.Error("insufficient stock when confirming payment; transaction rolled back",
						"order_number", order.OrderNumber,
						"product_variant_id", item.ProductVariantID,
						"quantity", item.Quantity)
					return fmt.Errorf("insufficient stock for item %s (%s/%s) on order %s",
						item.ProductName, item.Color, item.Size, order.OrderNumber)
				}
			}
			order.Status = models.OrderStatusPaid

		case models.PaymentStatusExpired:
			// Payment window lapsed while pending; stock was never decremented, so stock remains unchanged.
			order.Status = models.OrderStatusExpired

		case models.PaymentStatusFailed:
			// Payment failed while pending; stock was never decremented, so stock remains unchanged.
			order.Status = models.OrderStatusCancelled

		default:
			// PENDING stays PENDING_PAYMENT; stock remains unchanged.
			return nil
		}

		return s.orders.UpdateTx(ctx, tx, order)

	})
	if err != nil {
		if errors.Is(err, payment.ErrUnknownOrder) {
			return utils.NewAppError(404, "ORDER_NOT_FOUND", "order referenced by webhook not found")
		}
		if ae, ok := utils.AsAppError(err); ok {
			return ae
		}
		return utils.ErrInternal(err)
	}
	return nil
}
