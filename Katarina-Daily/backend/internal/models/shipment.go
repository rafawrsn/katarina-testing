package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Shipment statuses.
const (
	ShipmentStatusPending   = "PENDING"
	ShipmentStatusShipped   = "SHIPPED"
	ShipmentStatusInTransit = "IN_TRANSIT"
	ShipmentStatusDelivered = "DELIVERED"
	ShipmentStatusFailed    = "FAILED"
)

// Shipment tracks courier/tracking-number info for an Order (1:1).
type Shipment struct {
	BaseModelNoDelete
	OrderID        uuid.UUID  `gorm:"type:uuid;uniqueIndex;not null" json:"order_id"`
	Courier        string     `gorm:"type:varchar(50);not null" json:"courier"`
	Service        string     `gorm:"type:varchar(50)" json:"service"`
	TrackingNumber string     `gorm:"type:varchar(100)" json:"tracking_number"`
	Status         string     `gorm:"type:varchar(20);not null" json:"status"`
	ShippedAt      *time.Time `json:"shipped_at,omitempty"`
	DeliveredAt    *time.Time `json:"delivered_at,omitempty"`
}

func (s *Shipment) BeforeCreate(tx *gorm.DB) error {
	s.ID = generateID(s.ID)
	if s.Status == "" {
		s.Status = ShipmentStatusPending
	}
	return nil
}

func (Shipment) TableName() string { return "shipments" }
