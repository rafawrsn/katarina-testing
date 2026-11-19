package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

const (
	DiscountTypePercentage = "percentage"
	DiscountTypeFixed      = "fixed"

	VoucherStatusActive   = "active"
	VoucherStatusInactive = "inactive"
)

// Voucher is an optional discount code that checkout can apply.
type Voucher struct {
	BaseModel
	Code             string          `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	DiscountType     string          `gorm:"type:varchar(20);not null" json:"discount_type"`
	DiscountValue    decimal.Decimal `gorm:"type:numeric(12,2);not null" json:"discount_value"`
	MinimumPurchase  decimal.Decimal `gorm:"type:numeric(12,2);default:0" json:"minimum_purchase"`
	MaxDiscount      *decimal.Decimal `gorm:"type:numeric(12,2)" json:"max_discount,omitempty"`
	UsageLimit       *int            `json:"usage_limit,omitempty"`
	UsedCount        int             `gorm:"default:0" json:"used_count"`
	StartAt          *time.Time      `json:"start_at,omitempty"`
	ExpiredAt        *time.Time      `json:"expired_at,omitempty"`
	Status           string          `gorm:"type:varchar(20);default:active" json:"status"`
}

func (v *Voucher) BeforeCreate(tx *gorm.DB) error {
	v.ID = generateID(v.ID)
	if v.Status == "" {
		v.Status = VoucherStatusActive
	}
	return nil
}

func (Voucher) TableName() string { return "vouchers" }

// IsValidForAmount checks status, date window and minimum purchase.
func (v *Voucher) IsValidForAmount(amount decimal.Decimal, now time.Time) bool {
	if v.Status != VoucherStatusActive {
		return false
	}
	if v.StartAt != nil && now.Before(*v.StartAt) {
		return false
	}
	if v.ExpiredAt != nil && now.After(*v.ExpiredAt) {
		return false
	}
	if v.UsageLimit != nil && v.UsedCount >= *v.UsageLimit {
		return false
	}
	if amount.LessThan(v.MinimumPurchase) {
		return false
	}
	return true
}

// CalculateDiscount returns the discount amount for a given subtotal,
// applying max_discount as a cap for percentage-type vouchers.
func (v *Voucher) CalculateDiscount(subtotal decimal.Decimal) decimal.Decimal {
	var discount decimal.Decimal
	switch v.DiscountType {
	case DiscountTypePercentage:
		discount = subtotal.Mul(v.DiscountValue).Div(decimal.NewFromInt(100))
		if v.MaxDiscount != nil && discount.GreaterThan(*v.MaxDiscount) {
			discount = *v.MaxDiscount
		}
	case DiscountTypeFixed:
		discount = v.DiscountValue
	}
	if discount.GreaterThan(subtotal) {
		discount = subtotal
	}
	return discount
}
