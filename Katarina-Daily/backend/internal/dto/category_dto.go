package dto

import (
	"time"

	"github.com/google/uuid"

	"fashion-ecommerce-backend/internal/models"
)

type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,max=150"`
	Slug        string `json:"slug" binding:"omitempty,max=180"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url" binding:"omitempty,url"`
}

type UpdateCategoryRequest struct {
	Name        *string `json:"name" binding:"omitempty,max=150"`
	Slug        *string `json:"slug" binding:"omitempty,max=180"`
	Description *string `json:"description"`
	ImageURL    *string `json:"image_url" binding:"omitempty,url"`
}

type CategoryResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToCategoryResponse(c *models.Category) CategoryResponse {
	return CategoryResponse{
		ID:          c.ID,
		Name:        c.Name,
		Slug:        c.Slug,
		Description: c.Description,
		ImageURL:    c.ImageURL,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

func ToCategoryResponseList(items []models.Category) []CategoryResponse {
	out := make([]CategoryResponse, 0, len(items))
	for i := range items {
		out = append(out, ToCategoryResponse(&items[i]))
	}
	return out
}
