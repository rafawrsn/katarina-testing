package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"fashion-ecommerce-backend/internal/config"
	"fashion-ecommerce-backend/internal/database"
	"fashion-ecommerce-backend/internal/handlers"
	"fashion-ecommerce-backend/internal/middleware"
	"fashion-ecommerce-backend/internal/repositories"
	"fashion-ecommerce-backend/internal/routes"
	"fashion-ecommerce-backend/internal/services"
	"fashion-ecommerce-backend/pkg/payment"
	"fashion-ecommerce-backend/pkg/shipping"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cfg := config.Load()
	if cfg.IsProduction() {
		if cfg.JWTSecret == "" || cfg.JWTSecret == "change_me" {
			slog.Error("JWT_SECRET must be securely set in production environment")
			os.Exit(1)
		}
		if cfg.PaymentServerKey == "" {
			slog.Error("PAYMENT_SERVER_KEY must be set in production environment")
			os.Exit(1)
		}
	}

	db, err := database.Connect(cfg)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	if err := database.AutoMigrate(db); err != nil {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	if err := database.Seed(db); err != nil {
		slog.Warn("database seeding encountered an error", "error", err)
	}

	// ---- Repositories ----
	userRepo := repositories.NewUserRepository(db)
	addressRepo := repositories.NewAddressRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	productRepo := repositories.NewProductRepository(db)
	variantRepo := repositories.NewVariantRepository(db)
	cartRepo := repositories.NewCartRepository(db)
	orderRepo := repositories.NewOrderRepository(db)
	paymentRepo := repositories.NewPaymentRepository(db)
	shipmentRepo := repositories.NewShipmentRepository(db)
	voucherRepo := repositories.NewVoucherRepository(db)
	wishlistRepo := repositories.NewWishlistRepository(db)
	reviewRepo := repositories.NewReviewRepository(db)

	// ---- External provider abstractions ----
	// Swap these two constructors to plug in a real courier aggregator
	// or a different payment gateway - nothing else in the app changes.
	shippingProvider := shipping.NewMockProvider()
	paymentGateway := payment.NewMidtransGateway(cfg.PaymentServerKey, cfg.PaymentClientKey, cfg.PaymentEnvironment == "production")

	// ---- Services ----
	authService := services.NewAuthService(userRepo, cfg)
	addressService := services.NewAddressService(addressRepo, db)
	categoryService := services.NewCategoryService(categoryRepo)
	productService := services.NewProductService(productRepo, variantRepo, categoryRepo)
	cartService := services.NewCartService(cartRepo, variantRepo)
	shippingService := services.NewShippingService(addressRepo, cartRepo, shippingProvider, cfg)
	checkoutService := services.NewCheckoutService(db, userRepo, addressRepo, cartRepo, orderRepo, paymentRepo, voucherRepo, variantRepo, shippingService, paymentGateway)
	orderService := services.NewOrderService(db, orderRepo, variantRepo)
	paymentService := services.NewPaymentService(db, orderRepo, paymentRepo, variantRepo, paymentGateway)
	shipmentService := services.NewShipmentService(db, orderRepo, shipmentRepo)
	adminUserService := services.NewAdminUserService(userRepo)
	wishlistService := services.NewWishlistService(wishlistRepo, productRepo)
	reviewService := services.NewReviewService(reviewRepo, productRepo, orderRepo)

	// ---- Handlers ----
	h := &routes.Handlers{
		Auth:      handlers.NewAuthHandler(authService),
		Address:   handlers.NewAddressHandler(addressService),
		Category:  handlers.NewCategoryHandler(categoryService),
		Product:   handlers.NewProductHandler(productService),
		Cart:      handlers.NewCartHandler(cartService),
		Shipping:  handlers.NewShippingHandler(shippingService),
		Checkout:  handlers.NewCheckoutHandler(checkoutService),
		Order:     handlers.NewOrderHandler(orderService),
		Payment:   handlers.NewPaymentHandler(paymentService),
		Shipment:  handlers.NewShipmentHandler(shipmentService),
		AdminUser: handlers.NewAdminUserHandler(adminUserService),
		Wishlist:  handlers.NewWishlistHandler(wishlistService),
		Review:    handlers.NewReviewHandler(reviewService),
	}

	// ---- Gin engine ----
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(middleware.Recovery())
	r.Use(middleware.StructuredLogger())
	r.Use(middleware.SecureHeaders())
	r.Use(middleware.CORS(cfg.CORSAllowedOrigins))
	r.Use(middleware.RateLimit(cfg.RateLimitPerMinute))

	routes.Register(r, h, cfg)

	srv := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		slog.Info("server starting", "port", cfg.AppPort, "env", cfg.AppEnv)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown on SIGINT/SIGTERM (Railway sends SIGTERM on deploy).
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("forced server shutdown", "error", err)
	}

	if sqlDB, err := db.DB(); err == nil {
		_ = sqlDB.Close()
	}

	slog.Info("server exited cleanly")
}
