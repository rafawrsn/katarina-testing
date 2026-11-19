package dto

import (
	"time"

	"github.com/google/uuid"

	"fashion-ecommerce-backend/internal/models"
)

type CreateReviewRequest struct {
	// OrderID identifies which of the customer's paid orders this
	// review is attached to - the service verifies that order actually
	// belongs to the customer and contains this product before
	// accepting the review.
	OrderID uuid.UUID `json:"order_id" binding:"required"`
	Rating  int       `json:"rating" binding:"required,min=1,max=5"`
	Comment string    `json:"comment" binding:"omitempty,max=2000"`
}

type ReviewResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	UserName  string    `json:"user_name,omitempty"`
	ProductID uuid.UUID `json:"product_id"`
	OrderID   uuid.UUID `json:"order_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

func ToReviewResponse(r *models.Review) ReviewResponse {
	resp := ReviewResponse{
		ID:        r.ID,
		UserID:    r.UserID,
		ProductID: r.ProductID,
		OrderID:   r.OrderID,
		Rating:    r.Rating,
		Comment:   r.Comment,
		CreatedAt: r.CreatedAt,
	}
	if r.User.ID != uuid.Nil {
		resp.UserName = r.User.Name
	}
	return resp
}

func ToReviewResponseList(items []models.Review) []ReviewResponse {
	out := make([]ReviewResponse, 0, len(items))
	for i := range items {
		out = append(out, ToReviewResponse(&items[i]))
	}
	return out
}
