package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Address belongs to exactly one User. A user may have many addresses,
// one of which can be flagged as the default.
type Address struct {
	BaseModel
	UserID         uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Label          string    `gorm:"type:varchar(50)" json:"label"`
	RecipientName  string    `gorm:"type:varchar(150);not null" json:"recipient_name"`
	Phone          string    `gorm:"type:varchar(30);not null" json:"phone"`
	AddressLine    string    `gorm:"type:text;not null;column:address" json:"address"`
	Province       string    `gorm:"type:varchar(100);not null" json:"province"`
	City           string    `gorm:"type:varchar(100);not null" json:"city"`
	District       string    `gorm:"type:varchar(100)" json:"district"`
	PostalCode     string    `gorm:"type:varchar(20);not null" json:"postal_code"`
	IsDefault      bool      `gorm:"default:false" json:"is_default"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

func (a *Address) BeforeCreate(tx *gorm.DB) error {
	a.ID = generateID(a.ID)
	return nil
}

func (Address) TableName() string { return "addresses" }
