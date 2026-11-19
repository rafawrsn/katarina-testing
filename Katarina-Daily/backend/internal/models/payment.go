package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Payment statuses. Only a verified payment-gateway webhook may move a
// payment into PAID/FAILED/EXPIRED/REFUNDED - the frontend's opinion of
// payment status is never trusted.
const (
	PaymentStatusPending  = "PENDING"
	PaymentStatusPaid     = "PAID"
	PaymentStatusFailed   = "FAILED"
	PaymentStatusExpired  = "EXPIRED"
	PaymentStatusRefunded = "REFUNDED"
)

// Payment stores payment-gateway transaction data for an Order (1:1).
type Payment struct {
	BaseModelNoDelete
	OrderID         uuid.UUID          `gorm:"type:uuid;uniqueIndex;not null" json:"order_id"`
	PaymentGateway  string             `gorm:"type:varchar(50);not null" json:"payment_gateway"`
	TransactionID   *string            `gorm:"type:varchar(150);uniqueIndex" json:"transaction_id,omitempty"`
	PaymentMethod   string             `gorm:"type:varchar(50)" json:"payment_method"`
	Amount          decimal.Decimal    `gorm:"type:numeric(12,2);not null" json:"amount"`
	Status          string             `gorm:"type:varchar(20);not null;index" json:"status"`
	PaidAt          *time.Time         `json:"paid_at,omitempty"`
	ExpiredAt       *time.Time         `json:"expired_at,omitempty"`
	RawResponse     datatypes.JSON     `gorm:"type:jsonb" json:"raw_response,omitempty"`
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	p.ID = generateID(p.ID)
	if p.Status == "" {
		p.Status = PaymentStatusPending
	}
	return nil
}

func (Payment) TableName() string { return "payments" }
