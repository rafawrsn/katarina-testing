package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	DefaultPage  = 1
	DefaultLimit = 20
	MaxLimit     = 100
)

// PaginationParams is parsed from ?page= & ?limit= query params.
type PaginationParams struct {
	Page  int
	Limit int
}

// Offset returns the SQL OFFSET for these params.
func (p PaginationParams) Offset() int {
	return (p.Page - 1) * p.Limit
}

// ParsePagination reads page/limit from the query string, applying
// sane defaults and an upper bound on limit.
func ParsePagination(c *gin.Context) PaginationParams {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = DefaultPage
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}
	return PaginationParams{Page: page, Limit: limit}
}

// PaginationMeta is returned alongside list endpoints.
type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

func NewPaginationMeta(p PaginationParams, totalItems int64) PaginationMeta {
	totalPages := int(totalItems) / p.Limit
	if int(totalItems)%p.Limit != 0 {
		totalPages++
	}
	if totalPages < 1 {
		totalPages = 1
	}
	return PaginationMeta{
		Page:       p.Page,
		Limit:      p.Limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}
