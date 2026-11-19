package payment

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"testing"
)

func TestHandleWebhook_MissingServerKey(t *testing.T) {
	gw := NewMidtransGateway("", "client-key", false)

	payload := []byte(`{
		"order_id": "KD-123456",
		"status_code": "200",
		"gross_amount": "100000",
		"signature_key": "some_sig"
	}`)

	_, err := gw.HandleWebhook(context.Background(), payload, nil)
	if err == nil {
		t.Fatal("expected error when ServerKey is empty, got nil")
	}
}

func TestHandleWebhook_InvalidSignature(t *testing.T) {
	gw := NewMidtransGateway("real-server-key", "client-key", false)

	payload := []byte(`{
		"order_id": "KD-123456",
		"status_code": "200",
		"gross_amount": "100000",
		"signature_key": "invalidsignaturehash",
		"transaction_status": "settlement"
	}`)

	_, err := gw.HandleWebhook(context.Background(), payload, nil)
	if !errors.Is(err, ErrInvalidSignature) {
		t.Fatalf("expected ErrInvalidSignature, got %v", err)
	}
}

func TestHandleWebhook_ValidSignature(t *testing.T) {
	serverKey := "test-secret-key"
	gw := NewMidtransGateway(serverKey, "client-key", false)

	orderID := "KD-123456"
	statusCode := "200"
	grossAmount := "100000"

	sum := sha512.Sum512([]byte(orderID + statusCode + grossAmount + serverKey))
	validSig := hex.EncodeToString(sum[:])

	payload := []byte(`{
		"order_id": "` + orderID + `",
		"status_code": "` + statusCode + `",
		"gross_amount": "` + grossAmount + `",
		"signature_key": "` + validSig + `",
		"transaction_id": "tx-999",
		"transaction_status": "settlement",
		"fraud_status": "accept",
		"payment_type": "gopay"
	}`)

	res, err := gw.HandleWebhook(context.Background(), payload, nil)
	if err != nil {
		t.Fatalf("expected successful webhook verification, got %v", err)
	}
	if res.OrderID != orderID {
		t.Errorf("expected OrderID %s, got %s", orderID, res.OrderID)
	}
	if res.Status != "PAID" {
		t.Errorf("expected Status PAID, got %s", res.Status)
	}
	if res.PaymentMethod != "gopay" {
		t.Errorf("expected PaymentMethod gopay, got %s", res.PaymentMethod)
	}
}
