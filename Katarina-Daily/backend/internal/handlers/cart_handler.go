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

type CartHandler struct {
	cartService *services.CartService
}

func NewCartHandler(cartService *services.CartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

func (h *CartHandler) Get(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	resp, err := h.cartService.GetCart(c.Request.Context(), userID)
	if err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "cart retrieved", resp)
}

func (h *CartHandler) AddItem(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req dto.AddCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error(), "VALIDATION_ERROR")
		return
	}

	resp, err := h.cartService.AddItem(c.Request.Context(), userID, req)
	if err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, "item added to cart", resp)
}

func (h *CartHandler) UpdateItem(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	itemID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid cart item id", "VALIDATION_ERROR")
		return
	}

	var req dto.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error(), "VALIDATION_ERROR")
		return
	}

	resp, err := h.cartService.UpdateItem(c.Request.Context(), userID, itemID, req)
	if err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "cart item updated", resp)
}

func (h *CartHandler) DeleteItem(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	itemID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid cart item id", "VALIDATION_ERROR")
		return
	}

	if err := h.cartService.DeleteItem(c.Request.Context(), userID, itemID); err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "cart item removed", nil)
}

func (h *CartHandler) Clear(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	if err := h.cartService.ClearCart(c.Request.Context(), userID); err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "cart cleared", nil)
}
