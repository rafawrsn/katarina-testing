package models

import "gorm.io/gorm"

const (
	RoleCustomer = "customer"
	RoleAdmin    = "admin"
)

// User represents both customers and admins. Role-based access control
// is driven by the Role field.
type User struct {
	BaseModel
	Name         string `gorm:"type:varchar(150)" json:"name"`
	Email        string `gorm:"type:varchar(150);uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"`
	Phone        string `gorm:"type:varchar(30)" json:"phone"`
	Role         string `gorm:"type:varchar(20);not null;default:customer" json:"role"`

	Addresses []Address `gorm:"foreignKey:UserID" json:"addresses,omitempty"`
	Cart      *Cart     `gorm:"foreignKey:UserID" json:"cart,omitempty"`
	Orders    []Order   `gorm:"foreignKey:UserID" json:"orders,omitempty"`
	Wishlists []Wishlist `gorm:"foreignKey:UserID" json:"wishlists,omitempty"`
	Reviews   []Review  `gorm:"foreignKey:UserID" json:"reviews,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = generateID(u.ID)
	if u.Role == "" {
		u.Role = RoleCustomer
	}
	return nil
}

func (User) TableName() string { return "users" }

func (u *User) IsAdmin() bool { return u.Role == RoleAdmin }
