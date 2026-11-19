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

type ReviewHandler struct {
	reviewService *services.ReviewService
}

func NewReviewHandler(reviewService *services.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewService: reviewService}
}

// ListByProduct handles GET /api/v1/products/:id/reviews - public.
func (h *ReviewHandler) ListByProduct(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid product id", "VALIDATION_ERROR")
		return
	}

	pagination := utils.ParsePagination(c)
	items, meta, err := h.reviewService.ListByProduct(c.Request.Context(), productID, pagination)
	if err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.SuccessWithMeta(c, http.StatusOK, "reviews retrieved", items, meta)
}

// Create handles POST /api/v1/products/:id/reviews - customer only,
// and only for a product they've actually purchased (see dto.CreateReviewRequest).
func (h *ReviewHandler) Create(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid product id", "VALIDATION_ERROR")
		return
	}

	var req dto.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, err.Error(), "VALIDATION_ERROR")
		return
	}

	resp, err := h.reviewService.Create(c.Request.Context(), userID, productID, req)
	if err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, "review submitted", resp)
}

// DeleteOwn handles DELETE /api/v1/reviews/:id - a customer deleting
// their own review.
func (h *ReviewHandler) DeleteOwn(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	reviewID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid review id", "VALIDATION_ERROR")
		return
	}

	if err := h.reviewService.DeleteOwn(c.Request.Context(), userID, reviewID); err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "review deleted", nil)
}

// DeleteAdmin handles DELETE /api/v1/admin/reviews/:id - admin
// moderation, can remove any review.
func (h *ReviewHandler) DeleteAdmin(c *gin.Context) {
	reviewID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid review id", "VALIDATION_ERROR")
		return
	}

	if err := h.reviewService.DeleteAdmin(c.Request.Context(), reviewID); err != nil {
		utils.FailFromAppError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "review deleted", nil)
}
