package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"fashion-ecommerce-backend/internal/models"
)

// ---- Products ----

type CreateProductRequest struct {
	CategoryID  uuid.UUID `json:"category_id" binding:"required"`
	Name        string    `json:"name" binding:"required,max=200"`
	Slug        string    `json:"slug" binding:"omitempty,max=220"`
	Description string    `json:"description"`
	Status      string    `json:"status" binding:"omitempty,oneof=active inactive draft"`
}

type UpdateProductRequest struct {
	CategoryID  *uuid.UUID `json:"category_id"`
	Name        *string    `json:"name" binding:"omitempty,max=200"`
	Slug        *string    `json:"slug" binding:"omitempty,max=220"`
	Description *string    `json:"description"`
	Status      *string    `json:"status" binding:"omitempty,oneof=active inactive draft"`
}

type ProductResponse struct {
	ID          uuid.UUID              `json:"id"`
	CategoryID  uuid.UUID              `json:"category_id"`
	Category    *CategoryResponse      `json:"category,omitempty"`
	Name        string                 `json:"name"`
	Slug        string                 `json:"slug"`
	Description string                 `json:"description"`
	Status      string                 `json:"status"`
	MinPrice    decimal.Decimal        `json:"min_price"`
	MaxPrice    decimal.Decimal        `json:"max_price"`
	TotalStock  int                    `json:"total_stock"`
	Variants    []VariantResponse      `json:"variants,omitempty"`
	Images      []ProductImageResponse `json:"images,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

func ToProductResponse(p *models.Product) ProductResponse {
	resp := ProductResponse{
		ID:          p.ID,
		CategoryID:  p.CategoryID,
		Name:        p.Name,
		Slug:        p.Slug,
		Description: p.Description,
		Status:      p.Status,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
	if p.Category.ID != uuid.Nil {
		cat := ToCategoryResponse(&p.Category)
		resp.Category = &cat
	}
	if len(p.Variants) > 0 {
		resp.Variants = ToVariantResponseList(p.Variants)
		resp.MinPrice = p.Variants[0].Price
		resp.MaxPrice = p.Variants[0].Price
		for _, v := range p.Variants {
			if v.Price.LessThan(resp.MinPrice) {
				resp.MinPrice = v.Price
			}
			if v.Price.GreaterThan(resp.MaxPrice) {
				resp.MaxPrice = v.Price
			}
			resp.TotalStock += v.Stock
		}
	}
	if len(p.Images) > 0 {
		resp.Images = ToProductImageResponseList(p.Images)
	}
	return resp
}

func ToProductResponseList(items []models.Product) []ProductResponse {
	out := make([]ProductResponse, 0, len(items))
	for i := range items {
		out = append(out, ToProductResponse(&items[i]))
	}
	return out
}

// ---- Variants ----

type CreateVariantRequest struct {
	SKU    string          `json:"sku" binding:"required,max=100"`
	Color  string          `json:"color" binding:"required,max=50"`
	Size   string          `json:"size" binding:"required,max=20"`
	Price  decimal.Decimal `json:"price" binding:"required"`
	Stock  int             `json:"stock" binding:"gte=0"`
	Weight int             `json:"weight" binding:"required,gt=0"`
	Status string          `json:"status" binding:"omitempty,oneof=active inactive"`
}

type UpdateVariantRequest struct {
	SKU    *string          `json:"sku" binding:"omitempty,max=100"`
	Color  *string          `json:"color" binding:"omitempty,max=50"`
	Size   *string          `json:"size" binding:"omitempty,max=20"`
	Price  *decimal.Decimal `json:"price"`
	Stock  *int             `json:"stock" binding:"omitempty,gte=0"`
	Weight *int             `json:"weight" binding:"omitempty,gt=0"`
	Status *string          `json:"status" binding:"omitempty,oneof=active inactive"`
}

type VariantResponse struct {
	ID        uuid.UUID       `json:"id"`
	ProductID uuid.UUID       `json:"product_id"`
	SKU       string          `json:"sku"`
	Color     string          `json:"color"`
	Size      string          `json:"size"`
	Price     decimal.Decimal `json:"price"`
	Stock     int             `json:"stock"`
	Weight    int             `json:"weight"`
	Status    string          `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

func ToVariantResponse(v *models.ProductVariant) VariantResponse {
	return VariantResponse{
		ID:        v.ID,
		ProductID: v.ProductID,
		SKU:       v.SKU,
		Color:     v.Color,
		Size:      v.Size,
		Price:     v.Price,
		Stock:     v.Stock,
		Weight:    v.WeightGrams,
		Status:    v.Status,
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}
}

func ToVariantResponseList(items []models.ProductVariant) []VariantResponse {
	out := make([]VariantResponse, 0, len(items))
	for i := range items {
		out = append(out, ToVariantResponse(&items[i]))
	}
	return out
}

// ---- Images ----

type CreateProductImageRequest struct {
	ImageURL  string `json:"image_url" binding:"required,url"`
	IsPrimary bool   `json:"is_primary"`
	SortOrder int    `json:"sort_order"`
}

type ProductImageResponse struct {
	ID        uuid.UUID `json:"id"`
	ImageURL  string    `json:"image_url"`
	IsPrimary bool      `json:"is_primary"`
	SortOrder int       `json:"sort_order"`
}

func ToProductImageResponseList(items []models.ProductImage) []ProductImageResponse {
	out := make([]ProductImageResponse, 0, len(items))
	for _, img := range items {
		out = append(out, ProductImageResponse{
			ID:        img.ID,
			ImageURL:  img.ImageURL,
			IsPrimary: img.IsPrimary,
			SortOrder: img.SortOrder,
		})
	}
	return out
}

// ProductListFilter captures the supported query filters for
// GET /api/v1/products.
type ProductListFilter struct {
	CategoryID *uuid.UUID
	Search     string
	Color      string
	Size       string
	MinPrice   *decimal.Decimal
	MaxPrice   *decimal.Decimal
	Status     string
}
