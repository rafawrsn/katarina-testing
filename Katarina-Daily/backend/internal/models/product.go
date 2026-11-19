package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	ProductStatusActive   = "active"
	ProductStatusInactive = "inactive"
	ProductStatusDraft    = "draft"
)

// Product is the parent catalog entity. Price, stock, color, size and
// weight intentionally do NOT live here - they belong on ProductVariant,
// since a fashion product is sold as color/size combinations that can
// each have a different price and stock level.
type Product struct {
	BaseModel
	CategoryID  uuid.UUID `gorm:"type:uuid;not null;index" json:"category_id"`
	Name        string    `gorm:"type:varchar(200);not null" json:"name"`
	Slug        string    `gorm:"type:varchar(220);uniqueIndex;not null" json:"slug"`
	Description string    `gorm:"type:text" json:"description"`
	Status      string    `gorm:"type:varchar(20);not null;default:active" json:"status"`

	Category Category         `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Variants []ProductVariant `gorm:"foreignKey:ProductID" json:"variants,omitempty"`
	Images   []ProductImage   `gorm:"foreignKey:ProductID" json:"images,omitempty"`
	Wishlists []Wishlist      `gorm:"foreignKey:ProductID" json:"-"`
	Reviews   []Review        `gorm:"foreignKey:ProductID" json:"-"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	p.ID = generateID(p.ID)
	if p.Status == "" {
		p.Status = ProductStatusActive
	}
	return nil
}

func (Product) TableName() string { return "products" }
