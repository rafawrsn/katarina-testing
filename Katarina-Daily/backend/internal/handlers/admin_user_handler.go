package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"fashion-ecommerce-backend/internal/services"
	"fashion-ecommerce-backend/internal/utils"
)

type AdminUserHandler struct {
	adminUserService *services.AdminUserService
}

func NewAdminUserHandler(adminUserService *services.AdminUserService) *AdminUserHandler {
	return &AdminUserHandler{adminUserService: adminUserService}
}

func (h *AdminUserHandler) List(c *gin.Context) {
	pagination := utils.ParsePagination(c)
	items, meta, err := h.adminUserService.List(c.Request.Context(), c.Query("role"), c.Query("search"), pagination)
	if err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.SuccessWithMeta(c, http.StatusOK, "customers retrieved", items, meta)
}

func (h *AdminUserHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid user id", "VALIDATION_ERROR")
		return
	}

	resp, err := h.adminUserService.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "customer retrieved", resp)
}
