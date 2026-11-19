package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Review is a 1-5 star rating + comment left by a customer who
// purchased the product (enforced in the service layer via order_id).
type Review struct {
	BaseModel
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	ProductID uuid.UUID `gorm:"type:uuid;not null;index" json:"product_id"`
	OrderID   uuid.UUID `gorm:"type:uuid;not null;index" json:"order_id"`
	Rating    int       `gorm:"not null" json:"rating"`
	Comment   string    `gorm:"type:text" json:"comment"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (r *Review) BeforeCreate(tx *gorm.DB) error {
	r.ID = generateID(r.ID)
	return nil
}

func (Review) TableName() string { return "reviews" }
