package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"fashion-ecommerce-backend/internal/dto"
	"fashion-ecommerce-backend/internal/models"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(ctx context.Context, p *models.Product) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *ProductRepository) FindByID(ctx context.Context, id uuid.UUID, withRelations bool) (*models.Product, error) {
	q := r.db.WithContext(ctx)
	if withRelations {
		q = q.Preload("Category").
			Preload("Variants", func(db *gorm.DB) *gorm.DB { return db.Order("color, size") }).
			Preload("Images", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order") })
	}
	var p models.Product
	if err := q.First(&p, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) FindBySlug(ctx context.Context, slug string) (*models.Product, error) {
	var p models.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Variants", func(db *gorm.DB) *gorm.DB { return db.Order("color, size") }).
		Preload("Images", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order") }).
		First(&p, "slug = ?", slug).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) ExistsBySlug(ctx context.Context, slug string, excludeID *uuid.UUID) (bool, error) {
	q := r.db.WithContext(ctx).Model(&models.Product{}).Where("slug = ?", slug)
	if excludeID != nil {
		q = q.Where("id <> ?", *excludeID)
	}
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// List returns products matching the given filter, with variants and
// the primary image preloaded, plus a total count for pagination.
func (r *ProductRepository) List(ctx context.Context, filter dto.ProductListFilter, offset, limit int) ([]models.Product, int64, error) {
	base := r.db.WithContext(ctx).Model(&models.Product{})

	if filter.CategoryID != nil {
		base = base.Where("category_id = ?", *filter.CategoryID)
	}
	if filter.Status != "" {
		base = base.Where("status = ?", filter.Status)
	} else {
		base = base.Where("status = ?", models.ProductStatusActive)
	}
	if filter.Search != "" {
		like := "%" + filter.Search + "%"
		base = base.Where("name ILIKE ? OR description ILIKE ?", like, like)
	}
	if filter.Color != "" || filter.Size != "" || filter.MinPrice != nil || filter.MaxPrice != nil {
		sub := r.db.Table("product_variants").Select("product_variants.product_id").Where("product_variants.deleted_at IS NULL")
		if filter.Color != "" {
			sub = sub.Where("product_variants.color ILIKE ?", filter.Color)
		}
		if filter.Size != "" {
			sub = sub.Where("product_variants.size ILIKE ?", filter.Size)
		}
		if filter.MinPrice != nil {
			sub = sub.Where("product_variants.price >= ?", *filter.MinPrice)
		}
		if filter.MaxPrice != nil {
			sub = sub.Where("product_variants.price <= ?", *filter.MaxPrice)
		}
		base = base.Where("products.id IN (?)", sub)
	}

	var total int64
	if err := base.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var products []models.Product
	err := base.
		Preload("Category").
		Preload("Variants").
		Preload("Images", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order").Limit(5) }).
		Order("products.created_at DESC").
		Offset(offset).Limit(limit).
		Find(&products).Error
	if err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

func (r *ProductRepository) Update(ctx context.Context, p *models.Product) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *ProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Product{}, "id = ?", id).Error
}

// AddImage appends a gallery image to a product.
func (r *ProductRepository) AddImage(ctx context.Context, img *models.ProductImage) error {
	return r.db.WithContext(ctx).Create(img).Error
}

func (r *ProductRepository) UnsetPrimaryImages(ctx context.Context, productID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.ProductImage{}).
		Where("product_id = ?", productID).
		Update("is_primary", false).Error
}
