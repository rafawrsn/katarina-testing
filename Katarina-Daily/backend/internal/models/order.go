package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Order statuses.
const (
	OrderStatusPendingPayment = "PENDING_PAYMENT"
	OrderStatusPaid           = "PAID"
	OrderStatusProcessing     = "PROCESSING"
	OrderStatusShipped        = "SHIPPED"
	OrderStatusDelivered      = "DELIVERED"
	OrderStatusCancelled      = "CANCELLED"
	OrderStatusExpired        = "EXPIRED"
)

// Order is the immutable record of a checkout. Shipping address fields
// are a SNAPSHOT taken at checkout time - orders must never depend on
// the customer's current Address rows, so a later address edit/delete
// cannot alter historical orders.
type Order struct {
	BaseModel
	UserID       uuid.UUID       `gorm:"type:uuid;not null;index" json:"user_id"`
	OrderNumber  string          `gorm:"type:varchar(50);uniqueIndex;not null" json:"order_number"`
	Status       string          `gorm:"type:varchar(30);not null;index" json:"status"`
	Subtotal     decimal.Decimal `gorm:"type:numeric(12,2);not null" json:"subtotal"`
	ShippingCost decimal.Decimal `gorm:"type:numeric(12,2);not null;default:0" json:"shipping_cost"`
	Discount     decimal.Decimal `gorm:"type:numeric(12,2);not null;default:0" json:"discount"`
	Total        decimal.Decimal `gorm:"type:numeric(12,2);not null" json:"total"`

	ShippingCourier string `gorm:"type:varchar(50)" json:"shipping_courier"`
	ShippingService string `gorm:"type:varchar(50)" json:"shipping_service"`

	// Snapshot fields - copied from the chosen Address at checkout time.
	ShippingAddress       string `gorm:"type:text;not null" json:"shipping_address"`
	ShippingRecipientName string `gorm:"type:varchar(150);not null" json:"shipping_recipient_name"`
	ShippingPhone         string `gorm:"type:varchar(30);not null" json:"shipping_phone"`
	ShippingCity          string `gorm:"type:varchar(100);not null" json:"shipping_city"`
	ShippingProvince      string `gorm:"type:varchar(100);not null" json:"shipping_province"`
	ShippingPostalCode    string `gorm:"type:varchar(20);not null" json:"shipping_postal_code"`

	VoucherID *uuid.UUID `gorm:"type:uuid;index" json:"voucher_id,omitempty"`

	Items    []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
	Payment  *Payment    `gorm:"foreignKey:OrderID" json:"payment,omitempty"`
	Shipment *Shipment   `gorm:"foreignKey:OrderID" json:"shipment,omitempty"`
	User     User        `gorm:"foreignKey:UserID" json:"-"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) error {
	o.ID = generateID(o.ID)
	if o.Status == "" {
		o.Status = OrderStatusPendingPayment
	}
	return nil
}

func (Order) TableName() string { return "orders" }

// IsCancellable reports whether a customer may still cancel this order.
func (o *Order) IsCancellable() bool {
	return o.Status == OrderStatusPendingPayment
}

// OrderItem is a purchase-time SNAPSHOT of a ProductVariant: name, SKU,
// color, size and price are copied in at order-creation time so that
// later catalog edits never alter historical order records.
type OrderItem struct {
	BaseModelNoDelete
	OrderID          uuid.UUID       `gorm:"type:uuid;not null;index" json:"order_id"`
	ProductVariantID uuid.UUID       `gorm:"type:uuid;not null;index" json:"product_variant_id"`
	ProductName      string          `gorm:"type:varchar(200);not null" json:"product_name"`
	SKU              string          `gorm:"type:varchar(100);not null" json:"sku"`
	Color            string          `gorm:"type:varchar(50)" json:"color"`
	Size             string          `gorm:"type:varchar(20)" json:"size"`
	Price            decimal.Decimal `gorm:"type:numeric(12,2);not null" json:"price"`
	Quantity         int             `gorm:"not null" json:"quantity"`
	Subtotal         decimal.Decimal `gorm:"type:numeric(12,2);not null" json:"subtotal"`
	ProductVariant 	 ProductVariant  `gorm:"foreignKey:ProductVariantID" json:"-"`
}

func (i *OrderItem) BeforeCreate(tx *gorm.DB) error {
	i.ID = generateID(i.ID)
	return nil
}

func (OrderItem) TableName() string { return "order_items" }
