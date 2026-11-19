package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"fashion-ecommerce-backend/internal/dto"
	"fashion-ecommerce-backend/internal/models"
	"fashion-ecommerce-backend/internal/repositories"
	"fashion-ecommerce-backend/internal/utils"
	"fashion-ecommerce-backend/pkg/payment"
)

type CheckoutService struct {
	db          *gorm.DB
	users       *repositories.UserRepository
	addresses   *repositories.AddressRepository
	carts       *repositories.CartRepository
	orders      *repositories.OrderRepository
	payments    *repositories.PaymentRepository
	vouchers    *repositories.VoucherRepository
	variants    *repositories.VariantRepository
	shippingSvc *ShippingService
	gateway     payment.Gateway
}

func NewCheckoutService(
	db *gorm.DB,
	users *repositories.UserRepository,
	addresses *repositories.AddressRepository,
	carts *repositories.CartRepository,
	orders *repositories.OrderRepository,
	payments *repositories.PaymentRepository,
	vouchers *repositories.VoucherRepository,
	variants *repositories.VariantRepository,
	shippingSvc *ShippingService,
	gateway payment.Gateway,
) *CheckoutService {
	return &CheckoutService{
		db: db, users: users, addresses: addresses, carts: carts,
		orders: orders, payments: payments, vouchers: vouchers,
		variants: variants,
		shippingSvc: shippingSvc, gateway: gateway,
	}
}

// Checkout implements the full server-side checkout flow from the spec:
// re-price the live cart, recompute weight and shipping, apply a
// voucher if valid, create the order + snapshot order items, clear the
// cart, open a payment transaction with the gateway, and return
// everything the client needs to complete payment. Nothing monetary is
// ever taken from the request body.
func (s *CheckoutService) Checkout(ctx context.Context, userID uuid.UUID, req dto.CheckoutRequest) (*dto.OrderResponse, error) {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, utils.ErrInternal(err)
	}

	addr, err := s.addresses.FindByID(ctx, req.AddressID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("address")
		}
		return nil, utils.ErrInternal(err)
	}
	if addr.UserID != userID {
		return nil, utils.ErrForbidden
	}

	cart, err := s.carts.GetOrCreateCart(ctx, userID)
	if err != nil {
		return nil, utils.ErrInternal(err)
	}
	fullCart, err := s.carts.GetCartWithItems(ctx, cart.ID)
	if err != nil {
		return nil, utils.ErrInternal(err)
	}
	if len(fullCart.Items) == 0 {
		return nil, utils.ErrValidation("cart is empty")
	}

	subtotal := decimal.Zero
	totalWeight := 0
	orderItems := make([]models.OrderItem, 0, len(fullCart.Items))
	for _, item := range fullCart.Items {
		v := item.ProductVariant
		if !v.IsActive() {
			return nil, utils.ErrValidation(fmt.Sprintf("%s (%s/%s) is no longer available", v.Product.Name, v.Color, v.Size))
		}
		if item.Quantity > v.Stock {
			return nil, utils.ErrValidation(fmt.Sprintf("insufficient stock for %s (%s/%s)", v.Product.Name, v.Color, v.Size))
		}
		lineSubtotal := v.Price.Mul(decimal.NewFromInt(int64(item.Quantity)))
		subtotal = subtotal.Add(lineSubtotal)
		totalWeight += v.WeightGrams * item.Quantity

		orderItems = append(orderItems, models.OrderItem{
			ProductVariantID: v.ID,
			ProductName:      v.Product.Name,
			SKU:              v.SKU,
			Color:            v.Color,
			Size:             v.Size,
			Price:            v.Price,
			Quantity:         item.Quantity,
			Subtotal:         lineSubtotal,
		})
	}

	shippingCost, err := s.shippingSvc.QuoteForCheckout(ctx, addr.Province, addr.City, addr.District, addr.PostalCode, totalWeight, req.ShippingCourier, req.ShippingService)
	if err != nil {
		if ae, ok := utils.AsAppError(err); ok {
			return nil, ae
		}
		return nil, utils.ErrInternal(err)
	}

	var order models.Order
	var paymentRecord models.Payment

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		discount := decimal.Zero
		var voucher *models.Voucher

		if req.VoucherCode != "" {
			v, vErr := s.vouchers.FindByCodeForUpdateTx(ctx, tx, req.VoucherCode)
			if vErr != nil {
				if errors.Is(vErr, gorm.ErrRecordNotFound) {
					return utils.ErrValidation("invalid voucher code")
				}
				return vErr
			}
			if !v.IsValidForAmount(subtotal, time.Now().UTC()) {
				return utils.ErrValidation("voucher is not valid for this order")
			}
			discount = v.CalculateDiscount(subtotal)
			voucher = v
		}

		total := subtotal.Add(shippingCost).Sub(discount)
		if total.IsNegative() {
			total = decimal.Zero
		}

		orderNumber, genErr := s.generateUniqueOrderNumber(ctx, tx)
		if genErr != nil {
			return genErr
		}

		order = models.Order{
			UserID:                userID,
			OrderNumber:           orderNumber,
			Status:                models.OrderStatusPendingPayment,
			Subtotal:              subtotal,
			ShippingCost:          shippingCost,
			Discount:              discount,
			Total:                 total,
			ShippingCourier:       req.ShippingCourier,
			ShippingService:       req.ShippingService,
			ShippingAddress:       addr.AddressLine,
			ShippingRecipientName: addr.RecipientName,
			ShippingPhone:         addr.Phone,
			ShippingCity:          addr.City,
			ShippingProvince:      addr.Province,
			ShippingPostalCode:    addr.PostalCode,
		}
		if voucher != nil {
			order.VoucherID = &voucher.ID
		}

		if err := s.orders.CreateTx(ctx, tx, &order); err != nil {
			return err
		}

		for i := range orderItems {
			orderItems[i].OrderID = order.ID
		}
		if err := s.orders.CreateItemsTx(ctx, tx, orderItems); err != nil {
			return err
		}

		// Note: Stock remains unchanged while payment status is PENDING.
		// Stock is only deducted once payment is confirmed as PAID / Successful.

		if voucher != nil {
			if err := s.vouchers.IncrementUsageTx(ctx, tx, voucher.ID); err != nil {
				return err
			}
		}

		// Cart is cleared once the order (with its own snapshot of
		// everything) has been durably created.
		if err := s.carts.ClearItemsTx(ctx, tx, cart.ID); err != nil {
			return err
		}

		paymentRecord = models.Payment{
			OrderID:        order.ID,
			PaymentGateway: s.gateway.Name(),
			Amount:         total,
			Status:         models.PaymentStatusPending,
		}
		return s.payments.CreateTx(ctx, tx, &paymentRecord)
	})
	if txErr != nil {
		if ae, ok := utils.AsAppError(txErr); ok {
			return nil, ae
		}
		return nil, utils.ErrInternal(txErr)
	}

	// The order, order items, voucher usage, cart clear and PENDING
	// payment row are already durably committed above. Opening the
	// transaction with the payment gateway is an external network call
	// and intentionally happens outside the DB transaction.
	gwItems := make([]payment.CreatePaymentItem, 0, len(orderItems))
	for _, oi := range orderItems {
		gwItems = append(gwItems, payment.CreatePaymentItem{
			ID:       oi.ProductVariantID.String(),
			Name:     oi.ProductName,
			Price:    oi.Price.Round(0).IntPart(),
			Quantity: oi.Quantity,
		})
	}

	gatewayResp, gwErr := s.gateway.CreatePayment(ctx, payment.CreatePaymentRequest{
		OrderID:       order.ID.String(),
		OrderNumber:   order.OrderNumber,
		Amount:        order.Total.Round(0).IntPart(),
		CustomerName:  user.Name,
		CustomerEmail: user.Email,
		CustomerPhone: user.Phone,
		Items:         gwItems,
	})
	if gwErr != nil {
		// The order and a PENDING payment row already exist; the
		// customer can retry initiating payment (e.g. via a "retry
		// payment" endpoint) without re-checking out. We surface a
		// 502 so the frontend knows this step specifically failed.
		slog.Error("payment gateway CreatePayment failed", "order_number", order.OrderNumber, "error", gwErr)
		return nil, utils.WrapAppError(502, "PAYMENT_GATEWAY_ERROR", "order was created but starting payment failed - please retry payment for this order", gwErr)
	}

	paymentRecord.TransactionID = &gatewayResp.TransactionID
	paymentRecord.RawResponse = datatypes.JSON(gatewayResp.RawResponse)
	if err := s.payments.Update(ctx, &paymentRecord); err != nil {
		slog.Error("failed to persist gateway response on payment record", "order_number", order.OrderNumber, "error", err)
	}

	loaded, err := s.orders.FindByID(ctx, order.ID)
	if err != nil {
		return nil, utils.ErrInternal(err)
	}

	resp := dto.ToOrderResponse(loaded)
	if resp.Payment != nil {
		resp.Payment.RedirectURL = gatewayResp.RedirectURL
		resp.Payment.SnapToken = gatewayResp.Token
	}
	return &resp, nil
}

func (s *CheckoutService) generateUniqueOrderNumber(ctx context.Context, tx *gorm.DB) (string, error) {
	for i := 0; i < 10; i++ {
		candidate := utils.GenerateOrderNumber()
		var count int64
		if err := tx.WithContext(ctx).Model(&models.Order{}).Where("order_number = ?", candidate).Count(&count).Error; err != nil {
			return "", err
		}
		if count == 0 {
			return candidate, nil
		}
	}
	return "", errors.New("failed to generate a unique order number after 10 attempts")
}
