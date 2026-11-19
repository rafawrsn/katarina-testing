package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"fashion-ecommerce-backend/internal/dto"
	"fashion-ecommerce-backend/internal/services"
	"fashion-ecommerce-backend/internal/utils"
)

type CategoryHandler struct {
	categoryService *services.CategoryService
}

func NewCategoryHandler(categoryService *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

func (h *CategoryHandler) List(c *gin.Context) {
	items, err := h.categoryService.List(c.Request.Context())
	if err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "categories retrieved", items)
}

func (h *CategoryHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid category id", "VALIDATION_ERROR")
		return
	}

	resp, err := h.categoryService.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "category retrieved", resp)
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error(), "VALIDATION_ERROR")
		return
	}

	resp, err := h.categoryService.Create(c.Request.Context(), req)
	if err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, "category created", resp)
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid category id", "VALIDATION_ERROR")
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error(), "VALIDATION_ERROR")
		return
	}

	resp, err := h.categoryService.Update(c.Request.Context(), id, req)
	if err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "category updated", resp)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid category id", "VALIDATION_ERROR")
		return
	}

	if err := h.categoryService.Delete(c.Request.Context(), id); err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "category deleted", nil)
}
