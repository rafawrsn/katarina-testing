package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"fashion-ecommerce-backend/internal/dto"
	"fashion-ecommerce-backend/internal/services"
	"fashion-ecommerce-backend/internal/utils"
)

type ProductHandler struct {
	productService *services.ProductService
}

func NewProductHandler(productService *services.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) List(c *gin.Context) {
	filter := dto.ProductListFilter{
		Search: c.Query("search"),
		Color:  c.Query("color"),
		Size:   c.Query("size"),
		Status: c.Query("status"),
	}
	if raw := c.Query("category_id"); raw != "" {
		if id, err := uuid.Parse(raw); err == nil {
			filter.CategoryID = &id
		}
	}
	if raw := c.Query("min_price"); raw != "" {
		if v, err := decimal.NewFromString(raw); err == nil {
			filter.MinPrice = &v
		}
	}
	if raw := c.Query("max_price"); raw != "" {
		if v, err := decimal.NewFromString(raw); err == nil {
			filter.MaxPrice = &v
		}
	}

	pagination := utils.ParsePagination(c)
	items, meta, err := h.productService.List(c.Request.Context(), filter, pagination)
	if err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.SuccessWithMeta(c, http.StatusOK, "products retrieved", items, meta)
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		// Allow slug lookups too: GET /products/:id also accepts a slug.
		resp, slugErr := h.productService.GetBySlug(c.Request.Context(), c.Param("id"))
		if slugErr != nil {
			utils.Fail(c, http.StatusBadRequest, "invalid product id", "VALIDATION_ERROR")
			return
		}
		utils.Success(c, http.StatusOK, "product retrieved", resp)
		return
	}

	resp, err := h.productService.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "product retrieved", resp)
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error(), "VALIDATION_ERROR")
		return
	}

	resp, err := h.productService.Create(c.Request.Context(), req)
	if err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, "product created", resp)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid product id", "VALIDATION_ERROR")
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error(), "VALIDATION_ERROR")
		return
	}

	resp, err := h.productService.Update(c.Request.Context(), id, req)
	if err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "product updated", resp)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid product id", "VALIDATION_ERROR")
		return
	}

	if err := h.productService.Delete(c.Request.Context(), id); err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "product deleted", nil)
}

func (h *ProductHandler) ListVariants(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid product id", "VALIDATION_ERROR")
		return
	}

	items, err := h.productService.ListVariants(c.Request.Context(), productID)
	if err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "variants retrieved", items)
}

func (h *ProductHandler) CreateVariant(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid product id", "VALIDATION_ERROR")
		return
	}

	var req dto.CreateVariantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error(), "VALIDATION_ERROR")
		return
	}

	resp, err := h.productService.CreateVariant(c.Request.Context(), productID, req)
	if err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, "variant created", resp)
}

func (h *ProductHandler) UpdateVariant(c *gin.Context) {
	variantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid variant id", "VALIDATION_ERROR")
		return
	}

	var req dto.UpdateVariantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error(), "VALIDATION_ERROR")
		return
	}

	resp, err := h.productService.UpdateVariant(c.Request.Context(), variantID, req)
	if err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "variant updated", resp)
}

func (h *ProductHandler) DeleteVariant(c *gin.Context) {
	variantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid variant id", "VALIDATION_ERROR")
		return
	}

	if err := h.productService.DeleteVariant(c.Request.Context(), variantID); err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "variant deleted", nil)
}

func (h *ProductHandler) AddImage(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid product id", "VALIDATION_ERROR")
		return
	}

	var req dto.CreateProductImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error(), "VALIDATION_ERROR")
		return
	}

	resp, err := h.productService.AddImage(c.Request.Context(), productID, req)
	if err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, "product image added", resp)
}
