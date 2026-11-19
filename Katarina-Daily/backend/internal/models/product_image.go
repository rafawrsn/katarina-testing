package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProductImage stores one or more gallery images for a Product.
type ProductImage struct {
	BaseModel
	ProductID uuid.UUID `gorm:"type:uuid;not null;index" json:"product_id"`
	ImageURL  string    `gorm:"type:text;not null" json:"image_url"`
	IsPrimary bool      `gorm:"default:false" json:"is_primary"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
}

func (i *ProductImage) BeforeCreate(tx *gorm.DB) error {
	i.ID = generateID(i.ID)
	return nil
}

func (ProductImage) TableName() string { return "product_images" }
