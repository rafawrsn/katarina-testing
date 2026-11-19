package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseModel is embedded by every entity that supports soft delete.
// It provides a UUID primary key and created_at/updated_at/deleted_at
// timestamps, matching the schema rules in the project spec.
type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BaseModelNoDelete is embedded by entities that do not use soft delete
// (carts, cart_items, order_items, payments, shipments, wishlists).
type BaseModelNoDelete struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// generateID is called from every model's BeforeCreate hook so primary
// keys are always populated application-side. This avoids depending on a
// Postgres UUID-generation extension being installed (gen_random_uuid()/
// uuid-ossp), which keeps the schema portable across managed Postgres
// providers such as Railway.
func generateID(existing uuid.UUID) uuid.UUID {
	if existing == uuid.Nil {
		return uuid.New()
	}
	return existing
}
