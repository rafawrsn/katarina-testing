package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"fashion-ecommerce-backend/internal/models"
)

type AddCartItemRequest struct {
	ProductVariantID uuid.UUID `json:"product_variant_id" binding:"required"`
	Quantity         int       `json:"quantity" binding:"required,gt=0"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,gt=0"`
}

type CartItemResponse struct {
	ID        uuid.UUID       `json:"id"`
	Variant   VariantResponse `json:"variant"`
	ProductID uuid.UUID       `json:"product_id"`
	ProductName string        `json:"product_name"`
	Quantity  int             `json:"quantity"`
	UnitPrice decimal.Decimal `json:"unit_price"`
	Subtotal  decimal.Decimal `json:"subtotal"`
	Available bool            `json:"available"`
	CreatedAt time.Time       `json:"created_at"`
}

type CartResponse struct {
	ID          uuid.UUID           `json:"id"`
	Items       []CartItemResponse  `json:"items"`
	TotalItems  int                 `json:"total_items"`
	TotalWeight int                 `json:"total_weight_grams"`
	Subtotal    decimal.Decimal     `json:"subtotal"`
}

// ToCartResponse builds the API-facing cart, computing totals from the
// current variant prices (never from any client-supplied values).
func ToCartResponse(cart *models.Cart) CartResponse {
	resp := CartResponse{
		ID:       cart.ID,
		Items:    make([]CartItemResponse, 0, len(cart.Items)),
		Subtotal: decimal.Zero,
	}
	for _, item := range cart.Items {
		lineSubtotal := item.ProductVariant.Price.Mul(decimal.NewFromInt(int64(item.Quantity)))
		resp.Items = append(resp.Items, CartItemResponse{
			ID:          item.ID,
			Variant:     ToVariantResponse(&item.ProductVariant),
			ProductID:   item.ProductVariant.ProductID,
			ProductName: item.ProductVariant.Product.Name,
			Quantity:    item.Quantity,
			UnitPrice:   item.ProductVariant.Price,
			Subtotal:    lineSubtotal,
			Available:   item.ProductVariant.IsPurchasable(item.Quantity),
		})
		resp.TotalItems += item.Quantity
		resp.TotalWeight += item.ProductVariant.WeightGrams * item.Quantity
		resp.Subtotal = resp.Subtotal.Add(lineSubtotal)
	}
	return resp
}
