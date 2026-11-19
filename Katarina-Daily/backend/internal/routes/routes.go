package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"fashion-ecommerce-backend/internal/config"
	"fashion-ecommerce-backend/internal/handlers"
	"fashion-ecommerce-backend/internal/middleware"
	"fashion-ecommerce-backend/internal/utils"
)

// Handlers bundles every HTTP handler the router needs. Built once in
// main.go and passed in here to keep route registration purely
// declarative.
type Handlers struct {
	Auth       *handlers.AuthHandler
	Address    *handlers.AddressHandler
	Category   *handlers.CategoryHandler
	Product    *handlers.ProductHandler
	Cart       *handlers.CartHandler
	Shipping   *handlers.ShippingHandler
	Checkout   *handlers.CheckoutHandler
	Order      *handlers.OrderHandler
	Payment    *handlers.PaymentHandler
	Shipment   *handlers.ShipmentHandler
	AdminUser  *handlers.AdminUserHandler
	Wishlist   *handlers.WishlistHandler
	Review     *handlers.ReviewHandler
}

// Register mounts every route from the spec onto the given engine.
func Register(r *gin.Engine, h *Handlers, cfg *config.Config) {
	r.GET("/health", func(c *gin.Context) {
		utils.Success(c, http.StatusOK, "ok", gin.H{"status": "healthy"})
	})

	auth := middleware.Auth(cfg.JWTSecret)
	admin := middleware.RequireAdmin()
	authLimiter := middleware.RateLimit(30) // stricter limit on auth endpoints

	api := r.Group("/api/v1")
	{
		// ---- Auth ----
		authGroup := api.Group("/auth")
		authGroup.Use(authLimiter)
		{
			authGroup.POST("/register", h.Auth.Register)
			authGroup.POST("/login", h.Auth.Login)
			authGroup.GET("/me", auth, h.Auth.Me)
		}

		// ---- Categories ----
		categories := api.Group("/categories")
		{
			categories.GET("", h.Category.List)
			categories.GET("/:id", h.Category.GetByID)
			categories.POST("", auth, admin, h.Category.Create)
			categories.PUT("/:id", auth, admin, h.Category.Update)
			categories.DELETE("/:id", auth, admin, h.Category.Delete)
		}

		// ---- Products ----
		products := api.Group("/products")
		{
			products.GET("", h.Product.List)
			products.GET("/:id", h.Product.GetByID)
			products.POST("", auth, admin, h.Product.Create)
			products.PUT("/:id", auth, admin, h.Product.Update)
			products.DELETE("/:id", auth, admin, h.Product.Delete)

			products.GET("/:id/variants", h.Product.ListVariants)
			products.POST("/:id/variants", auth, admin, h.Product.CreateVariant)
			products.POST("/:id/images", auth, admin, h.Product.AddImage)

			// Reviews are nested under their product: public to read,
			// authenticated (and purchase-verified) to write.
			products.GET("/:id/reviews", h.Review.ListByProduct)
			products.POST("/:id/reviews", auth, h.Review.Create)
		}

		// ---- Variants (top-level, per spec) ----
		variants := api.Group("/variants")
		variants.Use(auth, admin)
		{
			variants.PUT("/:id", h.Product.UpdateVariant)
			variants.DELETE("/:id", h.Product.DeleteVariant)
		}

		// ---- Addresses (customer, own data only) ----
		addresses := api.Group("/addresses")
		addresses.Use(auth)
		{
			addresses.GET("", h.Address.List)
			addresses.POST("", h.Address.Create)
			addresses.PUT("/:id", h.Address.Update)
			addresses.DELETE("/:id", h.Address.Delete)
			addresses.POST("/:id/default", h.Address.SetDefault)
		}

		// ---- Cart (customer) ----
		cart := api.Group("/cart")
		cart.Use(auth)
		{
			cart.GET("", h.Cart.Get)
			cart.POST("/items", h.Cart.AddItem)
			cart.PUT("/items/:id", h.Cart.UpdateItem)
			cart.DELETE("/items/:id", h.Cart.DeleteItem)
			cart.DELETE("", h.Cart.Clear)
		}

		// ---- Shipping rates (customer, needs an address on file) ----
		shipping := api.Group("/shipping")
		shipping.Use(auth)
		{
			shipping.GET("/rates", h.Shipping.GetRates)
		}

		// ---- Wishlist (customer, own data only) ----
		wishlist := api.Group("/wishlist")
		wishlist.Use(auth)
		{
			wishlist.GET("", h.Wishlist.List)
			wishlist.POST("", h.Wishlist.Add)
			wishlist.DELETE("/:product_id", h.Wishlist.Remove)
		}

		// ---- Reviews (top-level: deleting your own review) ----
		reviews := api.Group("/reviews")
		reviews.Use(auth)
		{
			reviews.DELETE("/:id", h.Review.DeleteOwn)
		}

		// ---- Orders (customer). POST here IS checkout: it turns the
		// customer's current cart into an order server-side. ----
		orders := api.Group("/orders")
		orders.Use(auth)
		{
			orders.GET("", h.Order.List)
			orders.GET("/:id", h.Order.GetByID)
			orders.POST("", h.Checkout.Checkout)
			orders.POST("/:id/cancel", h.Order.Cancel)
		}

		// ---- Payments ----
		payments := api.Group("/payments")
		{
			// Server-to-server notification from the payment gateway.
			// Deliberately NOT behind JWT auth (the gateway can't
			// supply one) - authenticity comes from signature
			// verification inside the handler/service instead. It gets
			// its own generous rate limit to absorb legitimate retries.
			payments.POST("/webhook", middleware.RateLimit(120), h.Payment.Webhook)
		}

		// ---- Admin ----
		adminGroup := api.Group("/admin")
		adminGroup.Use(auth, admin)
		{
			// Mirrors of the public catalog admin routes, per spec
			// section 32 - identical handlers to /products, /categories,
			// /variants above.
			adminGroup.GET("/products", h.Product.List)
			adminGroup.POST("/products", h.Product.Create)
			adminGroup.PUT("/products/:id", h.Product.Update)
			adminGroup.DELETE("/products/:id", h.Product.Delete)

			adminGroup.GET("/categories", func(c *gin.Context) { h.Category.List(c) })
			adminGroup.POST("/categories", h.Category.Create)
			adminGroup.PUT("/categories/:id", h.Category.Update)
			adminGroup.DELETE("/categories/:id", h.Category.Delete)

			adminGroup.PUT("/variants/:id", h.Product.UpdateVariant)
			adminGroup.DELETE("/variants/:id", h.Product.DeleteVariant)

			// Order management.
			adminGroup.GET("/orders", h.Order.AdminList)
			adminGroup.GET("/orders/:id", h.Order.AdminGetByID)
			adminGroup.PUT("/orders/:id/status", h.Order.AdminUpdateStatus)

			// Shipment management - one shipment per order.
			adminGroup.GET("/shipments", h.Shipment.List)
			adminGroup.POST("/orders/:id/shipment", h.Shipment.Create)
			adminGroup.PUT("/orders/:id/shipment", h.Shipment.UpdateStatus)

			// Customer management (read-only from the admin panel).
			adminGroup.GET("/users", h.AdminUser.List)
			adminGroup.GET("/users/:id", h.AdminUser.GetByID)

			// Review moderation - admins can remove any review.
			adminGroup.DELETE("/reviews/:id", h.Review.DeleteAdmin)
		}
	}
}
