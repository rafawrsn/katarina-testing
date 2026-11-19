package handlers

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"fashion-ecommerce-backend/internal/services"
	"fashion-ecommerce-backend/internal/utils"
)

type PaymentHandler struct {
	paymentService *services.PaymentService
}

func NewPaymentHandler(paymentService *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

// Webhook handles POST /api/v1/payments/webhook - the asynchronous,
// server-to-server payment gateway notification. It always reads the
// body and hands it to the gateway-specific signature verifier before
// trusting anything in it.
func (h *PaymentHandler) Webhook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "could not read request body", "INVALID_BODY")
		return
	}

	headers := make(map[string]string, len(c.Request.Header))
	for k := range c.Request.Header {
		headers[k] = c.Request.Header.Get(k)
	}

	if err := h.paymentService.HandleWebhook(c.Request.Context(), body, headers); err != nil {
		slog.Warn("payment webhook processing failed", "error", err)
		utils.FailFromAppError(c, err)
		return
	}

	// Payment gateways expect a plain 200 OK acknowledging receipt.
	utils.Success(c, http.StatusOK, "webhook processed", nil)
}
