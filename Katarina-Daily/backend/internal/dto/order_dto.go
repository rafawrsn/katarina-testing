package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"fashion-ecommerce-backend/internal/models"
)

type OrderItemResponse struct {
	ID          uuid.UUID       `json:"id"`
	ProductVariantID uuid.UUID  `json:"product_variant_id"`
	ProductName string          `json:"product_name"`
	SKU         string          `json:"sku"`
	Color       string          `json:"color"`
	Size        string          `json:"size"`
	Price       decimal.Decimal `json:"price"`
	Quantity    int             `json:"quantity"`
	Subtotal    decimal.Decimal `json:"subtotal"`
}

func toOrderItemResponse(i *models.OrderItem) OrderItemResponse {
	return OrderItemResponse{
		ID:               i.ID,
		ProductVariantID: i.ProductVariantID,
		ProductName:      i.ProductName,
		SKU:              i.SKU,
		Color:            i.Color,
		Size:             i.Size,
		Price:            i.Price,
		Quantity:         i.Quantity,
		Subtotal:         i.Subtotal,
	}
}

type OrderResponse struct {
	ID                    uuid.UUID           `json:"id"`
	OrderNumber           string              `json:"order_number"`
	Status                string              `json:"status"`
	Subtotal              decimal.Decimal     `json:"subtotal"`
	ShippingCost          decimal.Decimal     `json:"shipping_cost"`
	Discount              decimal.Decimal     `json:"discount"`
	Total                 decimal.Decimal     `json:"total"`
	ShippingCourier       string              `json:"shipping_courier"`
	ShippingService       string              `json:"shipping_service"`
	ShippingAddress       string              `json:"shipping_address"`
	ShippingRecipientName string              `json:"shipping_recipient_name"`
	ShippingPhone         string              `json:"shipping_phone"`
	ShippingCity          string              `json:"shipping_city"`
	ShippingProvince      string              `json:"shipping_province"`
	ShippingPostalCode    string              `json:"shipping_postal_code"`
	Items                 []OrderItemResponse `json:"items"`
	Payment               *PaymentResponse    `json:"payment,omitempty"`
	Shipment              *ShipmentResponse   `json:"shipment,omitempty"`
	CreatedAt             time.Time           `json:"created_at"`
	UpdatedAt             time.Time           `json:"updated_at"`
}

func ToOrderResponse(o *models.Order) OrderResponse {
	resp := OrderResponse{
		ID:                    o.ID,
		OrderNumber:           o.OrderNumber,
		Status:                o.Status,
		Subtotal:              o.Subtotal,
		ShippingCost:          o.ShippingCost,
		Discount:              o.Discount,
		Total:                 o.Total,
		ShippingCourier:       o.ShippingCourier,
		ShippingService:       o.ShippingService,
		ShippingAddress:       o.ShippingAddress,
		ShippingRecipientName: o.ShippingRecipientName,
		ShippingPhone:         o.ShippingPhone,
		ShippingCity:          o.ShippingCity,
		ShippingProvince:      o.ShippingProvince,
		ShippingPostalCode:    o.ShippingPostalCode,
		Items:                 make([]OrderItemResponse, 0, len(o.Items)),
		CreatedAt:             o.CreatedAt,
		UpdatedAt:             o.UpdatedAt,
	}
	for i := range o.Items {
		resp.Items = append(resp.Items, toOrderItemResponse(&o.Items[i]))
	}
	if o.Payment != nil {
		pr := ToPaymentResponse(o.Payment)
		resp.Payment = &pr
	}
	if o.Shipment != nil {
		sr := ToShipmentResponse(o.Shipment)
		resp.Shipment = &sr
	}
	return resp
}

func ToOrderResponseList(items []models.Order) []OrderResponse {
	out := make([]OrderResponse, 0, len(items))
	for i := range items {
		out = append(out, ToOrderResponse(&items[i]))
	}
	return out
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=PENDING_PAYMENT PAID PROCESSING SHIPPED DELIVERED CANCELLED EXPIRED"`
}

type CreateShipmentRequest struct {
	Courier        string `json:"courier" binding:"required"`
	Service        string `json:"service"`
	TrackingNumber string `json:"tracking_number"`
}

type UpdateShipmentRequest struct {
	Status         *string `json:"status" binding:"omitempty,oneof=PENDING SHIPPED IN_TRANSIT DELIVERED FAILED"`
	TrackingNumber *string `json:"tracking_number"`
}

type ShipmentResponse struct {
	ID             uuid.UUID  `json:"id"`
	Courier        string     `json:"courier"`
	Service        string     `json:"service"`
	TrackingNumber string     `json:"tracking_number"`
	Status         string     `json:"status"`
	ShippedAt      *time.Time `json:"shipped_at,omitempty"`
	DeliveredAt    *time.Time `json:"delivered_at,omitempty"`
}

func ToShipmentResponse(s *models.Shipment) ShipmentResponse {
	return ShipmentResponse{
		ID:             s.ID,
		Courier:        s.Courier,
		Service:        s.Service,
		TrackingNumber: s.TrackingNumber,
		Status:         s.Status,
		ShippedAt:      s.ShippedAt,
		DeliveredAt:    s.DeliveredAt,
	}
}
