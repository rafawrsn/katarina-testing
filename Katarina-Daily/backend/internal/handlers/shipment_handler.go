package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"fashion-ecommerce-backend/internal/dto"
	"fashion-ecommerce-backend/internal/services"
	"fashion-ecommerce-backend/internal/utils"
)

// ShipmentHandler exposes admin-only shipment management:
// creating a shipment for an order and updating its tracking
// number/status (which keeps the order status in sync).
type ShipmentHandler struct {
	shipmentService *services.ShipmentService
}

func NewShipmentHandler(shipmentService *services.ShipmentService) *ShipmentHandler {
	return &ShipmentHandler{shipmentService: shipmentService}
}

func (h *ShipmentHandler) List(c *gin.Context) {
	pagination := utils.ParsePagination(c)
	items, meta, err := h.shipmentService.List(c.Request.Context(), c.Query("status"), pagination)
	if err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.SuccessWithMeta(c, http.StatusOK, "shipments retrieved", items, meta)
}

func (h *ShipmentHandler) Create(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid order id", "VALIDATION_ERROR")
		return
	}

	var req dto.CreateShipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error(), "VALIDATION_ERROR")
		return
	}

	resp, err := h.shipmentService.Create(c.Request.Context(), orderID, req)
	if err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, "shipment created", resp)
}

func (h *ShipmentHandler) UpdateStatus(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid order id", "VALIDATION_ERROR")
		return
	}

	var req dto.UpdateShipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error(), "VALIDATION_ERROR")
		return
	}

	resp, err := h.shipmentService.UpdateStatus(c.Request.Context(), orderID, req)
	if err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "shipment updated", resp)
}
