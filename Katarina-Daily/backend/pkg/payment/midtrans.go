package payment

import (
	"bytes"
	"context"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"fashion-ecommerce-backend/internal/models"
)

// MidtransGateway implements Gateway against Midtrans' Snap + Core API
// (https://docs.midtrans.com). It is the reference implementation for
// an Indonesian payment gateway; swap it for another PSP by writing a
// second Gateway implementation - nothing else in the codebase needs
// to change.
type MidtransGateway struct {
	ServerKey   string
	ClientKey   string
	IsProduction bool
	httpClient  *http.Client
}

func NewMidtransGateway(serverKey, clientKey string, isProduction bool) *MidtransGateway {
	return &MidtransGateway{
		ServerKey:    serverKey,
		ClientKey:    clientKey,
		IsProduction: isProduction,
		httpClient:   &http.Client{Timeout: 15 * time.Second},
	}
}

func (g *MidtransGateway) Name() string { return "midtrans" }

func (g *MidtransGateway) snapBaseURL() string {
	if g.IsProduction {
		return "https://app.midtrans.com/snap/v1"
	}
	return "https://app.sandbox.midtrans.com/snap/v1"
}

func (g *MidtransGateway) coreBaseURL() string {
	if g.IsProduction {
		return "https://api.midtrans.com/v2"
	}
	return "https://api.sandbox.midtrans.com/v2"
}

func (g *MidtransGateway) basicAuthHeader() string {
	token := base64.StdEncoding.EncodeToString([]byte(g.ServerKey + ":"))
	return "Basic " + token
}

// --- Snap transaction request/response shapes ---

type snapItem struct {
	ID       string `json:"id"`
	Price    int64  `json:"price"`
	Quantity int    `json:"quantity"`
	Name     string `json:"name"`
}

type snapCreateRequest struct {
	TransactionDetails struct {
		OrderID     string `json:"order_id"`
		GrossAmount int64  `json:"gross_amount"`
	} `json:"transaction_details"`
	CustomerDetails struct {
		FirstName string `json:"first_name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
	} `json:"customer_details"`
	ItemDetails []snapItem `json:"item_details,omitempty"`
}

type snapCreateResponse struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
	ErrorMessages []string `json:"error_messages"`
}

func (g *MidtransGateway) CreatePayment(ctx context.Context, req CreatePaymentRequest) (*CreatePaymentResponse, error) {
	if strings.TrimSpace(g.ServerKey) == "" {
		return nil, errors.New("midtrans server key is not configured")
	}
	body := snapCreateRequest{}
	body.TransactionDetails.OrderID = req.OrderNumber
	body.TransactionDetails.GrossAmount = req.Amount
	body.CustomerDetails.FirstName = req.CustomerName
	body.CustomerDetails.Email = req.CustomerEmail
	body.CustomerDetails.Phone = req.CustomerPhone
	for _, item := range req.Items {
		body.ItemDetails = append(body.ItemDetails, snapItem{
			ID:       item.ID,
			Price:    item.Price,
			Quantity: item.Quantity,
			Name:     truncate(item.Name, 50),
		})
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal snap request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, g.snapBaseURL()+"/transactions", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", g.basicAuthHeader())

	resp, err := g.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("call midtrans snap api: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("midtrans snap api returned status %d: %s", resp.StatusCode, string(raw))
	}

	var parsed snapCreateResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, fmt.Errorf("decode snap response: %w", err)
	}

	return &CreatePaymentResponse{
		TransactionID: req.OrderNumber, // Midtrans uses order_id as the reference until a transaction_id is assigned on payment
		RedirectURL:   parsed.RedirectURL,
		Token:         parsed.Token,
		RawResponse:   raw,
	}, nil
}

// --- Core API status check ---

type coreStatusResponse struct {
	TransactionID     string `json:"transaction_id"`
	OrderID           string `json:"order_id"`
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status"`
	PaymentType       string `json:"payment_type"`
	StatusCode        string `json:"status_code"`
	GrossAmount       string `json:"gross_amount"`
	SignatureKey      string `json:"signature_key"`
}

func (g *MidtransGateway) GetPaymentStatus(ctx context.Context, transactionID string) (*StatusResult, error) {
	if strings.TrimSpace(g.ServerKey) == "" {
		return nil, errors.New("midtrans server key is not configured")
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, g.coreBaseURL()+"/"+transactionID+"/status", nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", g.basicAuthHeader())

	resp, err := g.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("call midtrans status api: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("midtrans status api returned status %d: %s", resp.StatusCode, string(raw))
	}

	var parsed coreStatusResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, fmt.Errorf("decode status response: %w", err)
	}

	return &StatusResult{
		TransactionID: parsed.TransactionID,
		Status:        mapMidtransStatus(parsed.TransactionStatus, parsed.FraudStatus),
		PaymentMethod: parsed.PaymentType,
		RawResponse:   raw,
	}, nil
}

// --- Webhook (HTTP notification) handling ---

type notificationPayload struct {
	OrderID           string `json:"order_id"`
	StatusCode        string `json:"status_code"`
	GrossAmount       string `json:"gross_amount"`
	SignatureKey      string `json:"signature_key"`
	TransactionID     string `json:"transaction_id"`
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status"`
	PaymentType       string `json:"payment_type"`
}

// HandleWebhook verifies the SHA512(order_id+status_code+gross_amount+
// ServerKey) signature Midtrans sends with every HTTP notification,
// per Midtrans' documented scheme, then returns the normalized result.
// It performs NO database writes - the payment service applies the
// result idempotently inside its own transaction.
func (g *MidtransGateway) HandleWebhook(ctx context.Context, payload []byte, headers map[string]string) (*WebhookResult, error) {
	if strings.TrimSpace(g.ServerKey) == "" {
		return nil, errors.New("midtrans server key is not configured")
	}

	var n notificationPayload
	if err := json.Unmarshal(payload, &n); err != nil {
		return nil, fmt.Errorf("decode webhook payload: %w", err)
	}

	if n.OrderID == "" || n.StatusCode == "" || n.GrossAmount == "" || n.SignatureKey == "" {
		return nil, ErrInvalidSignature
	}

	expected := sha512Hex(n.OrderID + n.StatusCode + n.GrossAmount + g.ServerKey)
	if !equalConstantTime(expected, n.SignatureKey) {
		return nil, ErrInvalidSignature
	}

	return &WebhookResult{
		OrderID:       n.OrderID, // this is our order_number, not the internal UUID
		TransactionID: n.TransactionID,
		Status:        mapMidtransStatus(n.TransactionStatus, n.FraudStatus),
		PaymentMethod: n.PaymentType,
		RawResponse:   payload,
	}, nil
}

func mapMidtransStatus(transactionStatus, fraudStatus string) string {
	switch transactionStatus {
	case "capture":
		if fraudStatus == "challenge" {
			return models.PaymentStatusPending
		}
		return models.PaymentStatusPaid
	case "settlement":
		return models.PaymentStatusPaid
	case "pending":
		return models.PaymentStatusPending
	case "deny", "cancel":
		return models.PaymentStatusFailed
	case "expire":
		return models.PaymentStatusExpired
	case "refund", "partial_refund":
		return models.PaymentStatusRefunded
	default:
		return models.PaymentStatusPending
	}
}

func sha512Hex(s string) string {
	sum := sha512.Sum512([]byte(s))
	return hex.EncodeToString(sum[:])
}

// equalConstantTime avoids leaking signature-comparison timing info.
// Uses crypto/subtle.ConstantTimeCompare so even a length mismatch
// does not reveal information through early return timing.
func equalConstantTime(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

// truncate clips s to at most n Unicode code points (runes), not bytes,
// so multi-byte UTF-8 characters (e.g. Indonesian text, emoji) are never
// split mid-character which would produce invalid UTF-8.
func truncate(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n])
}
