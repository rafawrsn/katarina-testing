package dto

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ShippingRatesRequest is used for GET /api/v1/shipping/rates.
type ShippingRatesRequest struct {
	AddressID uuid.UUID `form:"address_id" binding:"required"`
	Courier   string    `form:"courier"`
}

type ShippingRateOption struct {
	Courier      string          `json:"courier"`
	Service      string          `json:"service"`
	Description  string          `json:"description"`
	Cost         decimal.Decimal `json:"cost"`
	ETD          string          `json:"etd"`
}

// CheckoutRequest is the only input the client sends for checkout - it
// intentionally contains no prices, totals or stock. Everything
// monetary is recalculated server-side from the customer's live cart.
type CheckoutRequest struct {
	AddressID       uuid.UUID `json:"address_id" binding:"required"`
	ShippingCourier string    `json:"shipping_courier" binding:"required"`
	ShippingService string    `json:"shipping_service" binding:"required"`
	VoucherCode     string    `json:"voucher_code" binding:"omitempty"`
}
