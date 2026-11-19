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

type AddressHandler struct {
	addressService *services.AddressService
}

func NewAddressHandler(addressService *services.AddressService) *AddressHandler {
	return &AddressHandler{addressService: addressService}
}

func (h *AddressHandler) List(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	items, err := h.addressService.List(c.Request.Context(), userID)
	if err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "addresses retrieved", items)
}

func (h *AddressHandler) Create(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req dto.CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error(), "VALIDATION_ERROR")
		return
	}

	resp, err := h.addressService.Create(c.Request.Context(), userID, req)
	if err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, "address created", resp)
}

func (h *AddressHandler) Update(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	addressID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid address id", "VALIDATION_ERROR")
		return
	}

	var req dto.UpdateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error(), "VALIDATION_ERROR")
		return
	}

	resp, err := h.addressService.Update(c.Request.Context(), userID, addressID, req)
	if err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "address updated", resp)
}

func (h *AddressHandler) Delete(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	addressID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid address id", "VALIDATION_ERROR")
		return
	}

	if err := h.addressService.Delete(c.Request.Context(), userID, addressID); err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "address deleted", nil)
}

func (h *AddressHandler) SetDefault(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	addressID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid address id", "VALIDATION_ERROR")
		return
	}

	resp, err := h.addressService.SetDefault(c.Request.Context(), userID, addressID)
	if err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "default address updated", resp)
}
