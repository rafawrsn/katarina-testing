package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"fashion-ecommerce-backend/internal/dto"
	"fashion-ecommerce-backend/internal/models"
	"fashion-ecommerce-backend/internal/repositories"
	"fashion-ecommerce-backend/internal/utils"
)

type CategoryService struct {
	categories *repositories.CategoryRepository
}

func NewCategoryService(categories *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{categories: categories}
}

func (s *CategoryService) List(ctx context.Context) ([]dto.CategoryResponse, error) {
	items, err := s.categories.List(ctx)
	if err != nil {
		return nil, utils.ErrInternal(err)
	}
	return dto.ToCategoryResponseList(items), nil
}

func (s *CategoryService) GetByID(ctx context.Context, id uuid.UUID) (*dto.CategoryResponse, error) {
	c, err := s.categories.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("category")
		}
		return nil, utils.ErrInternal(err)
	}
	resp := dto.ToCategoryResponse(c)
	return &resp, nil
}

func (s *CategoryService) Create(ctx context.Context, req dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	slug := req.Slug
	if slug == "" {
		slug = utils.Slugify(req.Name)
	} else {
		slug = utils.Slugify(slug)
	}

	slug, err := s.uniqueSlug(ctx, slug, nil)
	if err != nil {
		return nil, err
	}

	c := &models.Category{
		Name:        req.Name,
		Slug:        slug,
		Description: req.Description,
		ImageURL:    req.ImageURL,
	}
	if err := s.categories.Create(ctx, c); err != nil {
		return nil, utils.ErrInternal(err)
	}
	resp := dto.ToCategoryResponse(c)
	return &resp, nil
}

func (s *CategoryService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	c, err := s.categories.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("category")
		}
		return nil, utils.ErrInternal(err)
	}

	if req.Name != nil {
		c.Name = *req.Name
	}
	if req.Description != nil {
		c.Description = *req.Description
	}
	if req.ImageURL != nil {
		c.ImageURL = *req.ImageURL
	}
	if req.Slug != nil && *req.Slug != "" {
		newSlug, uErr := s.uniqueSlug(ctx, utils.Slugify(*req.Slug), &id)
		if uErr != nil {
			return nil, uErr
		}
		c.Slug = newSlug
	}

	if err := s.categories.Update(ctx, c); err != nil {
		return nil, utils.ErrInternal(err)
	}
	resp := dto.ToCategoryResponse(c)
	return &resp, nil
}

func (s *CategoryService) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.categories.FindByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrNotFound("category")
		}
		return utils.ErrInternal(err)
	}

	hasProducts, err := s.categories.HasProducts(ctx, id)
	if err != nil {
		return utils.ErrInternal(err)
	}
	if hasProducts {
		return utils.ErrConflict("CATEGORY_HAS_PRODUCTS", "cannot delete a category that still has products")
	}

	if err := s.categories.Delete(ctx, id); err != nil {
		return utils.ErrInternal(err)
	}
	return nil
}

func (s *CategoryService) uniqueSlug(ctx context.Context, base string, excludeID *uuid.UUID) (string, error) {
	slug := base
	for i := 1; i < 50; i++ {
		exists, err := s.categories.ExistsBySlug(ctx, slug, excludeID)
		if err != nil {
			return "", utils.ErrInternal(err)
		}
		if !exists {
			return slug, nil
		}
		slug = base + "-" + utils.RandomCode(4)
	}
	return "", utils.NewAppError(500, "SLUG_GENERATION_FAILED", "could not generate a unique slug")
}
