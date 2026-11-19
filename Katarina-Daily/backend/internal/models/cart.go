package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Cart is a single active cart per user (1:1 with User).
type Cart struct {
	BaseModelNoDelete
	UserID uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`

	Items []CartItem `gorm:"foreignKey:CartID" json:"items,omitempty"`
}

func (c *Cart) BeforeCreate(tx *gorm.DB) error {
	c.ID = generateID(c.ID)
	return nil
}

func (Cart) TableName() string { return "carts" }

// CartItem links a cart to a product variant + quantity. A cart must
// never contain two rows for the same variant - the repository layer
// enforces "add again -> increment quantity" instead of inserting a
// duplicate row.
type CartItem struct {
	BaseModelNoDelete
	CartID           uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_cart_variant" json:"cart_id"`
	ProductVariantID uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_cart_variant" json:"product_variant_id"`
	Quantity         int       `gorm:"not null" json:"quantity"`

	ProductVariant ProductVariant `gorm:"foreignKey:ProductVariantID" json:"product_variant,omitempty"`
}

func (i *CartItem) BeforeCreate(tx *gorm.DB) error {
	i.ID = generateID(i.ID)
	return nil
}

func (CartItem) TableName() string { return "cart_items" }
