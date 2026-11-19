package database

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"fashion-ecommerce-backend/internal/models"
)

// Seed populates the database with initial catalog and user data if empty.
func Seed(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.Product{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		slog.Info("database already has products; skipping seed")
		return nil
	}

	slog.Info("seeding database with Katarina Daily catalog data...")

	// 1. Categories
	categories := []models.Category{
		{
			Name:        "Tops",
			Slug:        "tops",
			Description: "Minimalist shirts and blouses with subtle embroidery details for daily wear.",
			ImageURL:    "https://images.pexels.com/photos/22441297/pexels-photo-22441297.jpeg?auto=compress&cs=tinysrgb&h=900",
		},
		{
			Name:        "Outerwear",
			Slug:        "outerwear",
			Description: "Comfortable cardigans, lightweight cotton outers, and tailored coats.",
			ImageURL:    "https://images.pexels.com/photos/9603628/pexels-photo-9603628.jpeg?auto=compress&cs=tinysrgb&h=900",
		},
		{
			Name:        "Dresses",
			Slug:        "dresses",
			Description: "Korean minimalist silhouettes with delicate embroidered touches.",
			ImageURL:    "https://images.pexels.com/photos/8070398/pexels-photo-8070398.jpeg?auto=compress&cs=tinysrgb&h=900",
		},
		{
			Name:        "Bottoms",
			Slug:        "bottoms",
			Description: "Neutral linen pants and flowy minimal skirts.",
			ImageURL:    "https://images.pexels.com/photos/6373455/pexels-photo-6373455.jpeg?auto=compress&cs=tinysrgb&h=900",
		},
	}

	catMap := make(map[string]uuid.UUID)
	for i := range categories {
		if err := db.Where(models.Category{Slug: categories[i].Slug}).
			Attrs(categories[i]).
			FirstOrCreate(&categories[i]).Error; err != nil {
			slog.Error("failed to seed category", "name", categories[i].Name, "error", err)
			continue
		}
		catMap[categories[i].Slug] = categories[i].ID
	}

	// 3. Products with Variants and Images
	type variantSeed struct {
		SKU    string
		Color  string
		Size   string
		Price  string
		Stock  int
		Weight int
	}

	type productSeed struct {
		Name        string
		Slug        string
		CatSlug     string
		Description string
		ImageURL    string
		Variants    []variantSeed
	}

	products := []productSeed{
		{
			Name:        "Katarina Embroidery Shirt",
			Slug:        "katarina-embroidery-shirt",
			CatSlug:     "tops",
			Description: "Minimalist shirt with special embroidery details. Self-manufactured with breathable soft cotton.",
			ImageURL:    "https://images.pexels.com/photos/22441297/pexels-photo-22441297.jpeg?auto=compress&cs=tinysrgb&h=900",
			Variants: []variantSeed{
				{SKU: "KD-SHIRT-CRM-S", Color: "Cream", Size: "S", Price: "129000", Stock: 25, Weight: 250},
				{SKU: "KD-SHIRT-CRM-M", Color: "Cream", Size: "M", Price: "129000", Stock: 30, Weight: 250},
				{SKU: "KD-SHIRT-CRM-L", Color: "Cream", Size: "L", Price: "129000", Stock: 20, Weight: 250},
				{SKU: "KD-SHIRT-BGE-S", Color: "Beige", Size: "S", Price: "129000", Stock: 15, Weight: 250},
				{SKU: "KD-SHIRT-BGE-M", Color: "Beige", Size: "M", Price: "129000", Stock: 20, Weight: 250},
			},
		},
		{
			Name:        "Daily Minimalist Blouse",
			Slug:        "daily-minimalist-blouse",
			CatSlug:     "tops",
			Description: "Elegant Korean-inspired blouse for daily style with delicate cuff stitching.",
			ImageURL:    "https://images.pexels.com/photos/7945547/pexels-photo-7945547.jpeg?auto=compress&cs=tinysrgb&h=900",
			Variants: []variantSeed{
				{SKU: "KD-BLOUSE-WHT-S", Color: "Off White", Size: "S", Price: "149000", Stock: 30, Weight: 220},
				{SKU: "KD-BLOUSE-WHT-M", Color: "Off White", Size: "M", Price: "149000", Stock: 25, Weight: 220},
				{SKU: "KD-BLOUSE-BLU-S", Color: "Soft Blue", Size: "S", Price: "149000", Stock: 15, Weight: 220},
				{SKU: "KD-BLOUSE-BLU-M", Color: "Soft Blue", Size: "M", Price: "149000", Stock: 18, Weight: 220},
			},
		},
		{
			Name:        "Soft Cotton Outer",
			Slug:        "soft-cotton-outer",
			CatSlug:     "outerwear",
			Description: "Comfortable minimalist outerwear with an oversized drape and subtle lapel detail.",
			ImageURL:    "https://images.pexels.com/photos/9603628/pexels-photo-9603628.jpeg?auto=compress&cs=tinysrgb&h=900",
			Variants: []variantSeed{
				{SKU: "KD-OUTER-OAT-ALL", Color: "Oatmeal", Size: "All Size", Price: "179000", Stock: 25, Weight: 350},
				{SKU: "KD-OUTER-CHR-ALL", Color: "Charcoal", Size: "All Size", Price: "179000", Stock: 15, Weight: 350},
			},
		},
		{
			Name:        "Neutral Knit Cardigan",
			Slug:        "neutral-knit-cardigan",
			CatSlug:     "outerwear",
			Description: "Soft knit cardigan in warm neutral tones with custom horn-finish buttons.",
			ImageURL:    "https://images.pexels.com/photos/4641824/pexels-photo-4641824.jpeg?auto=compress&cs=tinysrgb&h=900",
			Variants: []variantSeed{
				{SKU: "KD-CARD-ALM-SM", Color: "Almond", Size: "S-M", Price: "159000", Stock: 20, Weight: 400},
				{SKU: "KD-CARD-ALM-LXL", Color: "Almond", Size: "L-XL", Price: "159000", Stock: 18, Weight: 400},
				{SKU: "KD-CARD-MCH-SM", Color: "Mocha", Size: "S-M", Price: "159000", Stock: 12, Weight: 400},
			},
		},
		{
			Name:        "Linen Everyday Shirt",
			Slug:        "linen-everyday-shirt",
			CatSlug:     "tops",
			Description: "Breathable linen shirt for effortless daily wear, featuring fine collar embroidery.",
			ImageURL:    "https://images.pexels.com/photos/6373455/pexels-photo-6373455.jpeg?auto=compress&cs=tinysrgb&h=900",
			Variants: []variantSeed{
				{SKU: "KD-LINEN-NAT-S", Color: "Natural", Size: "S", Price: "139000", Stock: 20, Weight: 230},
				{SKU: "KD-LINEN-NAT-M", Color: "Natural", Size: "M", Price: "139000", Stock: 25, Weight: 230},
				{SKU: "KD-LINEN-OLV-S", Color: "Olive", Size: "S", Price: "139000", Stock: 14, Weight: 230},
			},
		},
		{
			Name:        "Beige Tailored Coat",
			Slug:        "beige-tailored-coat",
			CatSlug:     "outerwear",
			Description: "Refined tailored coat with a calm minimalist silhouette and inner pocket embroidery.",
			ImageURL:    "https://images.pexels.com/photos/4456728/pexels-photo-4456728.jpeg?auto=compress&cs=tinysrgb&h=900",
			Variants: []variantSeed{
				{SKU: "KD-COAT-BGE-M", Color: "Warm Beige", Size: "M", Price: "199000", Stock: 12, Weight: 600},
				{SKU: "KD-COAT-BGE-L", Color: "Warm Beige", Size: "L", Price: "199000", Stock: 10, Weight: 600},
				{SKU: "KD-COAT-BLK-M", Color: "Black", Size: "M", Price: "199000", Stock: 15, Weight: 600},
			},
		},
	}

	for _, ps := range products {
		catID, ok := catMap[ps.CatSlug]
		if !ok {
			catID = catMap["tops"]
		}

		prod := models.Product{
			CategoryID:  catID,
			Name:        ps.Name,
			Slug:        ps.Slug,
			Description: ps.Description,
			Status:      models.ProductStatusActive,
		}

		if err := db.Where(models.Product{Slug: prod.Slug}).
			Attrs(prod).
			FirstOrCreate(&prod).Error; err != nil {
			slog.Error("failed to seed product", "name", prod.Name, "error", err)
			continue
		}

		// Image
		img := models.ProductImage{
			ProductID: prod.ID,
			ImageURL:  ps.ImageURL,
			IsPrimary: true,
			SortOrder: 0,
		}
		_ = db.FirstOrCreate(&img, models.ProductImage{ProductID: prod.ID, ImageURL: ps.ImageURL})

		// Variants
		for _, vs := range ps.Variants {
			price, _ := decimal.NewFromString(vs.Price)
			variant := models.ProductVariant{
				ProductID:   prod.ID,
				SKU:         vs.SKU,
				Color:       vs.Color,
				Size:        vs.Size,
				Price:       price,
				Stock:       vs.Stock,
				WeightGrams: vs.Weight,
				Status:      models.VariantStatusActive,
			}
			_ = db.FirstOrCreate(&variant, models.ProductVariant{SKU: vs.SKU})
		}
	}

	slog.Info("database seeding completed successfully")
	return nil
}
