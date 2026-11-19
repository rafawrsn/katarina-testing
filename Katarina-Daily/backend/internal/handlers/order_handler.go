package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"fashion-ecommerce-backend/internal/dto"
	"fashion-ecommerce-backend/internal/middleware"
	"fashion-ecommerce-backend/internal/services"
	"fashion-ecommerce-backend/internal/utils"
)

type OrderHandler struct {
	orderService *services.OrderService
}

func NewOrderHandler(orderService *services.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

// ---- Customer ----

func (h *OrderHandler) List(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	pagination := utils.ParsePagination(c)

	items, meta, err := h.orderService.ListByUser(c.Request.Context(), userID, c.Query("status"), pagination)
	if err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.SuccessWithMeta(c, http.StatusOK, "orders retrieved", items, meta)
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid order id", "VALIDATION_ERROR")
		return
	}

	resp, err := h.orderService.GetByIDForUser(c.Request.Context(), userID, orderID)
	if err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "order retrieved", resp)
}

func (h *OrderHandler) Cancel(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid order id", "VALIDATION_ERROR")
		return
	}

	resp, err := h.orderService.Cancel(c.Request.Context(), userID, orderID)
	if err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "order cancelled", resp)
}

// ---- Admin ----

func (h *OrderHandler) AdminList(c *gin.Context) {
	pagination := utils.ParsePagination(c)
	items, meta, err := h.orderService.ListAdmin(c.Request.Context(), c.Query("status"), pagination)
	if err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.SuccessWithMeta(c, http.StatusOK, "orders retrieved", items, meta)
}

func (h *OrderHandler) AdminGetByID(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid order id", "VALIDATION_ERROR")
		return
	}

	resp, err := h.orderService.GetByIDAdmin(c.Request.Context(), orderID)
	if err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "order retrieved", resp)
}

func (h *OrderHandler) AdminUpdateStatus(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid order id", "VALIDATION_ERROR")
		return
	}

	var req dto.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error(), "VALIDATION_ERROR")
		return
	}

	resp, err := h.orderService.UpdateStatusAdmin(c.Request.Context(), orderID, req.Status)
	if err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "order status updated", resp)
}
