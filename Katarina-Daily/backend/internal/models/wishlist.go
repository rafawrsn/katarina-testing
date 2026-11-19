package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Wishlist links a user to a product they want to save. The
// (user_id, product_id) pair is unique.
type Wishlist struct {
	BaseModelNoDelete
	UserID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_wishlist_user_product" json:"user_id"`
	ProductID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_wishlist_user_product" json:"product_id"`

	Product Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (w *Wishlist) BeforeCreate(tx *gorm.DB) error {
	w.ID = generateID(w.ID)
	return nil
}

func (Wishlist) TableName() string { return "wishlists" }
