package dto

import (
	"time"

	"github.com/google/uuid"

	"fashion-ecommerce-backend/internal/models"
)

type CreateAddressRequest struct {
	Label         string `json:"label" binding:"omitempty,max=50"`
	RecipientName string `json:"recipient_name" binding:"required,max=150"`
	Phone         string `json:"phone" binding:"required,max=30"`
	Address       string `json:"address" binding:"required"`
	Province      string `json:"province" binding:"required,max=100"`
	City          string `json:"city" binding:"required,max=100"`
	District      string `json:"district" binding:"omitempty,max=100"`
	PostalCode    string `json:"postal_code" binding:"required,max=20"`
	IsDefault     bool   `json:"is_default"`
}

type UpdateAddressRequest struct {
	Label         *string `json:"label"`
	RecipientName *string `json:"recipient_name" binding:"omitempty,max=150"`
	Phone         *string `json:"phone" binding:"omitempty,max=30"`
	Address       *string `json:"address"`
	Province      *string `json:"province" binding:"omitempty,max=100"`
	City          *string `json:"city" binding:"omitempty,max=100"`
	District      *string `json:"district"`
	PostalCode    *string `json:"postal_code" binding:"omitempty,max=20"`
	IsDefault     *bool   `json:"is_default"`
}

type AddressResponse struct {
	ID            uuid.UUID `json:"id"`
	Label         string    `json:"label"`
	RecipientName string    `json:"recipient_name"`
	Phone         string    `json:"phone"`
	Address       string    `json:"address"`
	Province      string    `json:"province"`
	City          string    `json:"city"`
	District      string    `json:"district"`
	PostalCode    string    `json:"postal_code"`
	IsDefault     bool      `json:"is_default"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func ToAddressResponse(a *models.Address) AddressResponse {
	return AddressResponse{
		ID:            a.ID,
		Label:         a.Label,
		RecipientName: a.RecipientName,
		Phone:         a.Phone,
		Address:       a.AddressLine,
		Province:      a.Province,
		City:          a.City,
		District:      a.District,
		PostalCode:    a.PostalCode,
		IsDefault:     a.IsDefault,
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
	}
}

func ToAddressResponseList(items []models.Address) []AddressResponse {
	out := make([]AddressResponse, 0, len(items))
	for i := range items {
		out = append(out, ToAddressResponse(&items[i]))
	}
	return out
}
