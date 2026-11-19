package models

import "gorm.io/gorm"

// Category groups products. A category has many products.
type Category struct {
	BaseModel
	Name        string `gorm:"type:varchar(150);not null" json:"name"`
	Slug        string `gorm:"type:varchar(180);uniqueIndex;not null" json:"slug"`
	Description string `gorm:"type:text" json:"description"`
	ImageURL    string `gorm:"type:varchar(500)" json:"image_url"`

	Products []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	c.ID = generateID(c.ID)
	return nil
}

func (Category) TableName() string { return "categories" }
