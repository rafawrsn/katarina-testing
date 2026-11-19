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

type ProductService struct {
	products   *repositories.ProductRepository
	variants   *repositories.VariantRepository
	categories *repositories.CategoryRepository
}

func NewProductService(products *repositories.ProductRepository, variants *repositories.VariantRepository, categories *repositories.CategoryRepository) *ProductService {
	return &ProductService{products: products, variants: variants, categories: categories}
}

// ---- Public catalog ----

func (s *ProductService) List(ctx context.Context, filter dto.ProductListFilter, pagination utils.PaginationParams) ([]dto.ProductResponse, utils.PaginationMeta, error) {
	items, total, err := s.products.List(ctx, filter, pagination.Offset(), pagination.Limit)
	if err != nil {
		return nil, utils.PaginationMeta{}, utils.ErrInternal(err)
	}
	return dto.ToProductResponseList(items), utils.NewPaginationMeta(pagination, total), nil
}

func (s *ProductService) GetByID(ctx context.Context, id uuid.UUID) (*dto.ProductResponse, error) {
	p, err := s.products.FindByID(ctx, id, true)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("product")
		}
		return nil, utils.ErrInternal(err)
	}
	resp := dto.ToProductResponse(p)
	return &resp, nil
}

func (s *ProductService) GetBySlug(ctx context.Context, slug string) (*dto.ProductResponse, error) {
	p, err := s.products.FindBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("product")
		}
		return nil, utils.ErrInternal(err)
	}
	resp := dto.ToProductResponse(p)
	return &resp, nil
}

func (s *ProductService) ListVariants(ctx context.Context, productID uuid.UUID) ([]dto.VariantResponse, error) {
	if _, err := s.products.FindByID(ctx, productID, false); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("product")
		}
		return nil, utils.ErrInternal(err)
	}
	items, err := s.variants.ListByProduct(ctx, productID)
	if err != nil {
		return nil, utils.ErrInternal(err)
	}
	return dto.ToVariantResponseList(items), nil
}

// ---- Admin: products ----

func (s *ProductService) Create(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error) {
	if _, err := s.categories.FindByID(ctx, req.CategoryID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrValidation("category_id does not reference an existing category")
		}
		return nil, utils.ErrInternal(err)
	}

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

	p := &models.Product{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Slug:        slug,
		Description: req.Description,
		Status:      req.Status,
	}
	if err := s.products.Create(ctx, p); err != nil {
		return nil, utils.ErrInternal(err)
	}

	created, err := s.products.FindByID(ctx, p.ID, true)
	if err != nil {
		return nil, utils.ErrInternal(err)
	}
	resp := dto.ToProductResponse(created)
	return &resp, nil
}

func (s *ProductService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	p, err := s.products.FindByID(ctx, id, false)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("product")
		}
		return nil, utils.ErrInternal(err)
	}

	if req.CategoryID != nil {
		if _, catErr := s.categories.FindByID(ctx, *req.CategoryID); catErr != nil {
			if errors.Is(catErr, gorm.ErrRecordNotFound) {
				return nil, utils.ErrValidation("category_id does not reference an existing category")
			}
			return nil, utils.ErrInternal(catErr)
		}
		p.CategoryID = *req.CategoryID
	}
	if req.Name != nil {
		p.Name = *req.Name
	}
	if req.Description != nil {
		p.Description = *req.Description
	}
	if req.Status != nil {
		p.Status = *req.Status
	}
	if req.Slug != nil && *req.Slug != "" {
		newSlug, sErr := s.uniqueSlug(ctx, utils.Slugify(*req.Slug), &id)
		if sErr != nil {
			return nil, sErr
		}
		p.Slug = newSlug
	}

	if err := s.products.Update(ctx, p); err != nil {
		return nil, utils.ErrInternal(err)
	}

	updated, err := s.products.FindByID(ctx, id, true)
	if err != nil {
		return nil, utils.ErrInternal(err)
	}
	resp := dto.ToProductResponse(updated)
	return &resp, nil
}

func (s *ProductService) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.products.FindByID(ctx, id, false); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrNotFound("product")
		}
		return utils.ErrInternal(err)
	}
	if err := s.products.Delete(ctx, id); err != nil {
		return utils.ErrInternal(err)
	}
	return nil
}

func (s *ProductService) AddImage(ctx context.Context, productID uuid.UUID, req dto.CreateProductImageRequest) (*dto.ProductImageResponse, error) {
	if _, err := s.products.FindByID(ctx, productID, false); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("product")
		}
		return nil, utils.ErrInternal(err)
	}

	if req.IsPrimary {
		if err := s.products.UnsetPrimaryImages(ctx, productID); err != nil {
			return nil, utils.ErrInternal(err)
		}
	}

	img := &models.ProductImage{
		ProductID: productID,
		ImageURL:  req.ImageURL,
		IsPrimary: req.IsPrimary,
		SortOrder: req.SortOrder,
	}
	if err := s.products.AddImage(ctx, img); err != nil {
		return nil, utils.ErrInternal(err)
	}

	return &dto.ProductImageResponse{
		ID:        img.ID,
		ImageURL:  img.ImageURL,
		IsPrimary: img.IsPrimary,
		SortOrder: img.SortOrder,
	}, nil
}

// ---- Admin: variants ----

func (s *ProductService) CreateVariant(ctx context.Context, productID uuid.UUID, req dto.CreateVariantRequest) (*dto.VariantResponse, error) {
	if _, err := s.products.FindByID(ctx, productID, false); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("product")
		}
		return nil, utils.ErrInternal(err)
	}

	skuTaken, err := s.variants.ExistsBySKU(ctx, req.SKU, nil)
	if err != nil {
		return nil, utils.ErrInternal(err)
	}
	if skuTaken {
		return nil, utils.ErrConflict("SKU_ALREADY_EXISTS", "a variant with this SKU already exists")
	}

	v := &models.ProductVariant{
		ProductID:   productID,
		SKU:         req.SKU,
		Color:       req.Color,
		Size:        req.Size,
		Price:       req.Price,
		Stock:       req.Stock,
		WeightGrams: req.Weight,
		Status:      req.Status,
	}
	if err := s.variants.Create(ctx, v); err != nil {
		return nil, utils.ErrInternal(err)
	}
	resp := dto.ToVariantResponse(v)
	return &resp, nil
}

func (s *ProductService) UpdateVariant(ctx context.Context, variantID uuid.UUID, req dto.UpdateVariantRequest) (*dto.VariantResponse, error) {
	v, err := s.variants.FindByID(ctx, variantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("variant")
		}
		return nil, utils.ErrInternal(err)
	}

	if req.SKU != nil && *req.SKU != v.SKU {
		taken, sErr := s.variants.ExistsBySKU(ctx, *req.SKU, &variantID)
		if sErr != nil {
			return nil, utils.ErrInternal(sErr)
		}
		if taken {
			return nil, utils.ErrConflict("SKU_ALREADY_EXISTS", "a variant with this SKU already exists")
		}
		v.SKU = *req.SKU
	}
	if req.Color != nil {
		v.Color = *req.Color
	}
	if req.Size != nil {
		v.Size = *req.Size
	}
	if req.Price != nil {
		v.Price = *req.Price
	}
	if req.Stock != nil {
		v.Stock = *req.Stock
	}
	if req.Weight != nil {
		v.WeightGrams = *req.Weight
	}
	if req.Status != nil {
		v.Status = *req.Status
	}

	if err := s.variants.Update(ctx, v); err != nil {
		return nil, utils.ErrInternal(err)
	}
	resp := dto.ToVariantResponse(v)
	return &resp, nil
}

func (s *ProductService) DeleteVariant(ctx context.Context, variantID uuid.UUID) error {
	if _, err := s.variants.FindByID(ctx, variantID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrNotFound("variant")
		}
		return utils.ErrInternal(err)
	}
	if err := s.variants.Delete(ctx, variantID); err != nil {
		return utils.ErrInternal(err)
	}
	return nil
}

func (s *ProductService) uniqueSlug(ctx context.Context, base string, excludeID *uuid.UUID) (string, error) {
	slug := base
	for i := 1; i < 50; i++ {
		exists, err := s.products.ExistsBySlug(ctx, slug, excludeID)
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
