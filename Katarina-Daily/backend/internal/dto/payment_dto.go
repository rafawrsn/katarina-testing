package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"fashion-ecommerce-backend/internal/models"
)

type PaymentResponse struct {
	ID             uuid.UUID  `json:"id"`
	OrderID        uuid.UUID  `json:"order_id"`
	PaymentGateway string     `json:"payment_gateway"`
	TransactionID  *string    `json:"transaction_id,omitempty"`
	PaymentMethod  string     `json:"payment_method"`
	Amount         decimal.Decimal `json:"amount"`
	Status         string     `json:"status"`
	PaidAt         *time.Time `json:"paid_at,omitempty"`
	ExpiredAt      *time.Time `json:"expired_at,omitempty"`
	// RedirectURL / Instructions are populated only right after
	// creation, sourced from the gateway response - never stored back
	// onto the Payment model itself.
	RedirectURL  string `json:"redirect_url,omitempty"`
	SnapToken    string `json:"snap_token,omitempty"`
}

func ToPaymentResponse(p *models.Payment) PaymentResponse {
	return PaymentResponse{
		ID:             p.ID,
		OrderID:        p.OrderID,
		PaymentGateway: p.PaymentGateway,
		TransactionID:  p.TransactionID,
		PaymentMethod:  p.PaymentMethod,
		Amount:         p.Amount,
		Status:         p.Status,
		PaidAt:         p.PaidAt,
		ExpiredAt:      p.ExpiredAt,
	}
}
