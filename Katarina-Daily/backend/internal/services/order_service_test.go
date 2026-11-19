package services

import (
	"testing"

	"fashion-ecommerce-backend/internal/models"
)

func TestStockDecremented(t *testing.T) {
	tests := []struct {
		status   string
		expected bool
	}{
		{models.OrderStatusPendingPayment, false},
		{models.OrderStatusPaid, true},
		{models.OrderStatusProcessing, true},
		{models.OrderStatusShipped, true},
		{models.OrderStatusDelivered, true},
		{models.OrderStatusCancelled, false},
		{models.OrderStatusExpired, false},
		{"UNKNOWN_STATUS", false},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			got := stockDecremented(tt.status)
			if got != tt.expected {
				t.Errorf("stockDecremented(%q) = %v; want %v", tt.status, got, tt.expected)
			}
		})
	}
}
