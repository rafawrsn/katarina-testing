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

type WishlistHandler struct {
	wishlistService *services.WishlistService
}

func NewWishlistHandler(wishlistService *services.WishlistService) *WishlistHandler {
	return &WishlistHandler{wishlistService: wishlistService}
}

func (h *WishlistHandler) List(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	items, err := h.wishlistService.List(c.Request.Context(), userID)
	if err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "wishlist retrieved", items)
}

func (h *WishlistHandler) Add(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req dto.AddWishlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error(), "VALIDATION_ERROR")
		return
	}

	resp, err := h.wishlistService.Add(c.Request.Context(), userID, req)
	if err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, "product added to wishlist", resp)
}

func (h *WishlistHandler) Remove(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	productID, err := uuid.Parse(c.Param("product_id"))
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid product id", "VALIDATION_ERROR")
		return
	}

	if err := h.wishlistService.Remove(c.Request.Context(), userID, productID); err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "product removed from wishlist", nil)
}
