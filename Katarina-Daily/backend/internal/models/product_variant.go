package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

const (
	VariantStatusActive   = "active"
	VariantStatusInactive = "inactive"
)

// ProductVariant is a purchasable SKU: a specific color/size combination
// of a Product, with its own price, stock and weight. This is the unit
// that carts, orders and stock management actually operate on.
type ProductVariant struct {
	BaseModel
	ProductID uuid.UUID       `gorm:"type:uuid;not null;index" json:"product_id"`
	SKU       string          `gorm:"type:varchar(100);uniqueIndex;not null" json:"sku"`
	Color     string          `gorm:"type:varchar(50);not null" json:"color"`
	Size      string          `gorm:"type:varchar(20);not null" json:"size"`
	Price     decimal.Decimal `gorm:"type:numeric(12,2);not null" json:"price"`
	Stock     int             `gorm:"not null;default:0" json:"stock"`
	// WeightGrams is the shipping weight of a single unit, in grams.
	WeightGrams int    `gorm:"not null;column:weight" json:"weight"`
	Status      string `gorm:"type:varchar(20);not null;default:active" json:"status"`

	Product    Product     `gorm:"foreignKey:ProductID" json:"-"`
	CartItems  []CartItem  `gorm:"foreignKey:ProductVariantID" json:"-"`
	OrderItems []OrderItem `gorm:"foreignKey:ProductVariantID" json:"-"`
}

func (v *ProductVariant) BeforeCreate(tx *gorm.DB) error {
	v.ID = generateID(v.ID)
	if v.Status == "" {
		v.Status = VariantStatusActive
	}
	return nil
}

func (ProductVariant) TableName() string { return "product_variants" }

// IsActive reports whether this variant can currently be purchased.
func (v *ProductVariant) IsActive() bool {
	return v.Status == VariantStatusActive
}

// IsPurchasable checks status and stock for a requested quantity.
func (v *ProductVariant) IsPurchasable(quantity int) bool {
	return v.IsActive() && quantity > 0 && v.Stock >= quantity
}
