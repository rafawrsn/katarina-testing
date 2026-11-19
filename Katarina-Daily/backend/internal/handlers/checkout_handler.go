package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"fashion-ecommerce-backend/internal/dto"
	"fashion-ecommerce-backend/internal/middleware"
	"fashion-ecommerce-backend/internal/services"
	"fashion-ecommerce-backend/internal/utils"
)

type CheckoutHandler struct {
	checkoutService *services.CheckoutService
}

func NewCheckoutHandler(checkoutService *services.CheckoutService) *CheckoutHandler {
	return &CheckoutHandler{checkoutService: checkoutService}
}

// Checkout handles POST /api/v1/orders, which doubles as the checkout
// endpoint: it turns the customer's current cart into an order. See
// dto.CheckoutRequest - no prices are ever accepted from the client.
func (h *CheckoutHandler) Checkout(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req dto.CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error(), "VALIDATION_ERROR")
		return
	}

	resp, err := h.checkoutService.Checkout(c.Request.Context(), userID, req)
	if err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, "order created, proceed to payment", resp)
}
