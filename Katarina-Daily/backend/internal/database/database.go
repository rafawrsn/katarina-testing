package database

import (
	"fmt"
	"log/slog"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"fashion-ecommerce-backend/internal/config"
	"fashion-ecommerce-backend/internal/models"
)

// Connect opens a PostgreSQL connection via GORM and configures the
// underlying connection pool. PostgreSQL is the only supported database
// for this project.
func Connect(cfg *config.Config) (*gorm.DB, error) {
	logLevel := gormlogger.Warn
	if !cfg.IsProduction() {
		logLevel = gormlogger.Info
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: gormlogger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	slog.Info("database connection established")
	return db, nil
}

// AutoMigrate creates/updates all tables. It intentionally lists models
// in dependency order (referenced tables before referencing tables) so
// foreign keys can be created on first run.
func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.User{},
		&models.Address{},
		&models.Category{},
		&models.Product{},
		&models.ProductVariant{},
		&models.ProductImage{},
		&models.Cart{},
		&models.CartItem{},
		&models.Voucher{},
		&models.Order{},
		&models.OrderItem{},
		&models.Payment{},
		&models.Shipment{},
		&models.Wishlist{},
		&models.Review{},
	)
	if err != nil {
		return fmt.Errorf("automigrate failed: %w", err)
	}
	slog.Info("database migration complete")
	return nil
}
