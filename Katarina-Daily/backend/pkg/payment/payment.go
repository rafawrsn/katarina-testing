// Package payment defines a provider-agnostic abstraction for creating
// payments, checking their status, and verifying webhook notifications.
// Checkout/order/payment services depend only on the Gateway interface
// below, never on a specific provider's SDK, so a different Indonesian
// PSP can be substituted without touching business logic.
package payment

import (
	"context"
	"errors"
)

var (
	ErrInvalidSignature = errors.New("invalid webhook signature")
	ErrUnknownOrder      = errors.New("order referenced by webhook not found")
)

// CreatePaymentRequest carries everything a gateway needs to open a
// transaction for an order. Amount is always the backend-calculated
// order total - never a client-supplied value.
type CreatePaymentRequest struct {
	OrderID     string
	OrderNumber string
	Amount      int64 // smallest currency unit; IDR has none, so plain Rupiah
	CustomerName  string
	CustomerEmail string
	CustomerPhone string
	Items       []CreatePaymentItem
}

type CreatePaymentItem struct {
	ID       string
	Name     string
	Price    int64
	Quantity int
}

// CreatePaymentResponse is what the caller needs to hand the payment
// off to the customer (redirect URL and/or a client-side token).
type CreatePaymentResponse struct {
	TransactionID string
	RedirectURL   string
	Token         string
	RawResponse   []byte
}

// StatusResult is the gateway's authoritative view of a transaction.
type StatusResult struct {
	TransactionID string
	Status        string // maps to models.Payment status constants
	PaymentMethod string
	RawResponse   []byte
}

// WebhookResult is the parsed, signature-verified outcome of an
// asynchronous payment notification.
type WebhookResult struct {
	OrderID       string
	TransactionID string
	Status        string
	PaymentMethod string
	RawResponse   []byte
}

// Gateway is the abstraction all payment-related services depend on.
type Gateway interface {
	Name() string
	CreatePayment(ctx context.Context, req CreatePaymentRequest) (*CreatePaymentResponse, error)
	GetPaymentStatus(ctx context.Context, transactionID string) (*StatusResult, error)
	// HandleWebhook validates the notification's authenticity (e.g. via
	// signature) and returns the normalized result. It performs no
	// database writes itself - the caller (payment service) owns the
	// idempotent transactional update.
	HandleWebhook(ctx context.Context, payload []byte, headers map[string]string) (*WebhookResult, error)
}
