package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"fashion-ecommerce-backend/internal/dto"
	"fashion-ecommerce-backend/internal/middleware"
	"fashion-ecommerce-backend/internal/services"
	"fashion-ecommerce-backend/internal/utils"
)

type ShippingHandler struct {
	shippingService *services.ShippingService
}

func NewShippingHandler(shippingService *services.ShippingService) *ShippingHandler {
	return &ShippingHandler{shippingService: shippingService}
}

// GetRates handles GET /api/v1/shipping/rates?address_id=...&courier=...
func (h *ShippingHandler) GetRates(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req dto.ShippingRatesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error(), "VALIDATION_ERROR")
		return
	}

	rates, err := h.shippingService.GetRatesForCart(c.Request.Context(), userID, req.AddressID, req.Courier)
	if err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "shipping rates retrieved", rates)
}
