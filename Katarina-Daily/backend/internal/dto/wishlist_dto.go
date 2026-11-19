package dto

import (
	"time"

	"github.com/google/uuid"
)

type AddWishlistRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
}

// WishlistItemResponse embeds the full ProductResponse so the frontend
// can render a wishlist grid without a second round trip per item.
type WishlistItemResponse struct {
	ID        uuid.UUID       `json:"id"`
	ProductID uuid.UUID       `json:"product_id"`
	Product   ProductResponse `json:"product"`
	CreatedAt time.Time       `json:"created_at"`
}
