package routes

import (
	"ecom-golang-clean-architecture/internal/delivery/http/handlers"
	"ecom-golang-clean-architecture/internal/delivery/http/middleware"
	"ecom-golang-clean-architecture/internal/infrastructure/config"
	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up all API routes
func SetupRoutes(
	router *gin.Engine,
	cfg *config.Config,
	userHandler *handlers.UserHandler,
	productHandler *handlers.ProductHandler,
	categoryHandler *handlers.CategoryHandler,
	cartHandler *handlers.CartHandler,
	orderHandler *handlers.OrderHandler,
	fileHandler *handlers.FileHandler,
	reviewHandler *handlers.ReviewHandler,
	wishlistHandler *handlers.WishlistHandler,
	couponHandler *handlers.CouponHandler,
	inventoryHandler *handlers.InventoryHandler,
	notificationHandler *handlers.NotificationHandler,
	analyticsHandler *handlers.AnalyticsHandler,
	addressHandler *handlers.AddressHandler,
	paymentHandler *handlers.PaymentHandler,
	shippingHandler *handlers.ShippingHandler,
	adminHandler *handlers.AdminHandler,
	oauthHandler *handlers.OAuthHandler,
) {
	// Apply global middleware
	router.Use(middleware.CORSMiddleware(&cfg.CORS))
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware())
	router.Use(middleware.ValidationMiddleware())

	// Serve static files for uploads
	router.Static("/uploads", "./uploads")

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "ecom-api",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)

			// OAuth routes
			if oauthHandler != nil {
				// OAuth URL generation (for frontend)
				auth.GET("/google/url", oauthHandler.GetGoogleAuthURL)
				auth.GET("/facebook/url", oauthHandler.GetFacebookAuthURL)

				// OAuth callbacks
				auth.GET("/google/callback", oauthHandler.GoogleCallback)
				auth.GET("/facebook/callback", oauthHandler.FacebookCallback)

				// Direct OAuth login (redirects to provider)
				auth.GET("/google/login", oauthHandler.GoogleLogin)
				auth.GET("/facebook/login", oauthHandler.FacebookLogin)
			}
		}

		// Public product routes
		products := v1.Group("/products")
		{
			products.GET("", productHandler.GetProducts)
			products.GET("/:id", productHandler.GetProduct)
			products.GET("/search", productHandler.SearchProducts)
			products.GET("/category/:categoryId", productHandler.GetProductsByCategory)
			products.GET("/featured", productHandler.GetFeaturedProducts)
			products.GET("/trending", productHandler.GetTrendingProducts)
			if reviewHandler != nil {
				products.GET("/:id/reviews", reviewHandler.GetProductReviews)
				products.GET("/:id/rating", reviewHandler.GetProductRating)
			}
			products.GET("/:id/related", productHandler.GetRelatedProducts)
		}

		// Public category routes
		categories := v1.Group("/categories")
		{
			categories.GET("", categoryHandler.GetCategories)
			categories.GET("/:id", categoryHandler.GetCategory)
			categories.GET("/tree", categoryHandler.GetCategoryTree)
			categories.GET("/root", categoryHandler.GetRootCategories)
			categories.GET("/:id/children", categoryHandler.GetCategoryChildren)
		}

		// Public file upload routes (no authentication required)
		publicUpload := v1.Group("/public/upload")
		{
			publicUpload.POST("/image", fileHandler.UploadImagePublic)
			publicUpload.POST("/document", fileHandler.UploadDocumentPublic)
		}

		// Public file routes
		publicFiles := v1.Group("/public/files")
		{
			publicFiles.GET("/:id", fileHandler.GetFileUpload)
		}

		// Shipping routes (public)
		if shippingHandler != nil {
			shipping := v1.Group("/shipping")
			{
				shipping.GET("/methods", shippingHandler.GetShippingMethods)
				// shipping.POST("/calculate", shippingHandler.CalculateShipping) // TODO: Implement CalculateShipping method
				shipping.GET("/track/:tracking_number", shippingHandler.TrackShipment)
			}
		}

		// Coupon routes (public validation)
		coupons := v1.Group("/coupons")
		{
			// coupons.GET("/public", couponHandler.GetActiveCoupons) // TODO: Implement GetActiveCoupons method
			coupons.POST("/validate", couponHandler.ValidateCoupon)
		}

		// Protected routes (authentication required)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/profile", userHandler.GetProfile)
				users.PUT("/profile", userHandler.UpdateProfile)
				users.POST("/change-password", userHandler.ChangePassword)
				// users.DELETE("/account", userHandler.DeleteAccount) // TODO: Implement DeleteAccount method
				if reviewHandler != nil {
					users.GET("/:user_id/reviews", reviewHandler.GetUserReviews)
				}
			}

			// Upload routes (authenticated users)
			upload := protected.Group("/upload")
			{
				upload.POST("/image", fileHandler.UploadImage)
				upload.POST("/document", fileHandler.UploadDocument)
			}

			// File management routes (authenticated users)
			files := protected.Group("/files")
			{
				files.GET("", fileHandler.GetFileUploads)
				files.GET("/:id", fileHandler.GetFileUpload)
				files.DELETE("/:id", fileHandler.DeleteFile)
			}

			// Cart routes
			cart := protected.Group("/cart")
			{
				cart.GET("", cartHandler.GetCart)
				cart.POST("/items", cartHandler.AddToCart)
				cart.PUT("/items", cartHandler.UpdateCartItem)
				cart.DELETE("/items/:productId", cartHandler.RemoveFromCart)
				cart.DELETE("", cartHandler.ClearCart)
				// cart.POST("/sync", cartHandler.SyncCart) // TODO: Implement SyncCart method
			}

			// Order routes
			orders := protected.Group("/orders")
			{
				orders.POST("", orderHandler.CreateOrder)
				orders.GET("", orderHandler.GetUserOrders)
				orders.GET("/:id", orderHandler.GetOrder)
				orders.POST("/:id/cancel", orderHandler.CancelOrder)
				// orders.GET("/:id/invoice", orderHandler.GetOrderInvoice) // TODO: Implement GetOrderInvoice method
				// orders.POST("/:id/reorder", orderHandler.ReorderItems) // TODO: Implement ReorderItems method
			}

			// Review routes
			reviews := protected.Group("/reviews")
			{
				reviews.POST("", reviewHandler.CreateReview)
				reviews.GET("/:id", reviewHandler.GetReview)
				reviews.PUT("/:id", reviewHandler.UpdateReview)
				reviews.DELETE("/:id", reviewHandler.DeleteReview)
				reviews.POST("/:id/vote", reviewHandler.VoteReview)
			}

			// Wishlist routes
			wishlist := protected.Group("/wishlist")
			{
				wishlist.GET("", wishlistHandler.GetWishlist)
				wishlist.POST("/items", wishlistHandler.AddToWishlist)
				wishlist.DELETE("/items/:id", wishlistHandler.RemoveFromWishlist)
				wishlist.DELETE("/clear", wishlistHandler.ClearWishlist)
				// wishlist.POST("/items/:product_id/move-to-cart", wishlistHandler.MoveToCart) // TODO: Implement MoveToCart method
				wishlist.GET("/count", wishlistHandler.GetWishlistCount)
			}

			// Address routes
			addresses := protected.Group("/addresses")
			{
				// addresses.GET("", addressHandler.GetUserAddresses) // TODO: Implement GetUserAddresses method
				addresses.POST("", addressHandler.CreateAddress)
				addresses.GET("/:id", addressHandler.GetAddress)
				addresses.PUT("/:id", addressHandler.UpdateAddress)
				addresses.DELETE("/:id", addressHandler.DeleteAddress)
				addresses.PUT("/:id/default", addressHandler.SetDefaultAddress)
				// addresses.POST("/validate", addressHandler.ValidateAddress) // TODO: Implement ValidateAddress method
			}

			// Payment routes
			payments := protected.Group("/payments")
			{
				payments.POST("", paymentHandler.ProcessPayment)
				payments.POST("/checkout-session", paymentHandler.CreateCheckoutSession)
				payments.GET("/:id", paymentHandler.GetPayment)
				payments.POST("/:id/refund", paymentHandler.ProcessRefund)
				payments.GET("/methods", paymentHandler.GetUserPaymentMethods)
				payments.POST("/methods", paymentHandler.SavePaymentMethod)
				payments.DELETE("/methods/:id", paymentHandler.DeletePaymentMethod)
				payments.PUT("/methods/:id/default", paymentHandler.SetDefaultPaymentMethod)
			}

			// Notification routes
			notifications := protected.Group("/notifications")
			{
				notifications.GET("", notificationHandler.GetUserNotifications)
				notifications.PUT("/:id/read", notificationHandler.MarkAsRead)
				notifications.PUT("/read-all", notificationHandler.MarkAllAsRead)
				notifications.GET("/count", notificationHandler.GetUnreadCount)
				// notifications.GET("/preferences", notificationHandler.GetUserPreferences) // TODO: Implement GetUserPreferences method
				// notifications.PUT("/preferences", notificationHandler.UpdateUserPreferences) // TODO: Implement UpdateUserPreferences method
			}
		}

		// Admin routes (admin authentication required)
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		admin.Use(middleware.AdminMiddleware())
		{
			// Dashboard routes
			dashboard := admin.Group("/dashboard")
			{
				dashboard.GET("", adminHandler.GetDashboard)
				dashboard.GET("/stats", adminHandler.GetSystemStats)
				dashboard.GET("/real-time", analyticsHandler.GetRealTimeMetrics)
			}

			// Admin user management
			adminUsers := admin.Group("/users")
			{
				adminUsers.GET("", adminHandler.GetUsers)
				adminUsers.PUT("/:id/status", adminHandler.UpdateUserStatus)
				adminUsers.PUT("/:id/role", adminHandler.UpdateUserRole)
				adminUsers.GET("/:id/activity", adminHandler.GetUserActivity)
			}

			// Admin product management
			adminProducts := admin.Group("/products")
			{
				adminProducts.POST("", productHandler.CreateProduct)
				adminProducts.PUT("/:id", productHandler.UpdateProduct)        // Complete replacement
				adminProducts.PATCH("/:id", productHandler.PatchProduct)       // Partial update
				adminProducts.DELETE("/:id", productHandler.DeleteProduct)
				adminProducts.PUT("/:id/stock", productHandler.UpdateStock)
			}

			// Admin category management
			adminCategories := admin.Group("/categories")
			{
				adminCategories.POST("", categoryHandler.CreateCategory)
				adminCategories.PUT("/:id", categoryHandler.UpdateCategory)
				adminCategories.DELETE("/:id", categoryHandler.DeleteCategory)
			}

			// Admin file uploads
			adminUpload := admin.Group("/upload")
			{
				adminUpload.POST("/image", fileHandler.UploadImageAdmin)
				adminUpload.POST("/document", fileHandler.UploadDocumentAdmin)
			}

			// Admin file management
			adminFiles := admin.Group("/files")
			{
				adminFiles.GET("", fileHandler.GetFileUploads)
				adminFiles.GET("/:id", fileHandler.GetFileUpload)
				adminFiles.DELETE("/:id", fileHandler.DeleteFile)
			}

			// Admin order management
			adminOrders := admin.Group("/orders")
			{
				adminOrders.GET("", adminHandler.GetOrders)
				adminOrders.GET("/:id", adminHandler.GetOrderDetails)
				adminOrders.PUT("/:id/status", adminHandler.UpdateOrderStatus)
				adminOrders.POST("/:id/refund", adminHandler.ProcessRefund)
			}

			// Review management routes
			adminReviews := admin.Group("/reviews")
			{
				adminReviews.GET("", adminHandler.ManageReviews)
				adminReviews.PUT("/:id/status", adminHandler.UpdateReviewStatus)
			}

			// Inventory management routes
			inventory := admin.Group("/inventory")
			{
				inventory.GET("", inventoryHandler.GetInventories)
				inventory.GET("/:id", inventoryHandler.GetInventory)
				inventory.PUT("/:id", inventoryHandler.UpdateInventory)
				inventory.POST("/movements", inventoryHandler.RecordMovement)
				inventory.GET("/movements", inventoryHandler.GetMovements)
				inventory.POST("/adjust", inventoryHandler.AdjustStock)
				inventory.POST("/transfer", inventoryHandler.TransferStock)
				inventory.GET("/alerts", inventoryHandler.GetStockAlerts)
				inventory.PUT("/alerts/:id/resolve", inventoryHandler.ResolveAlert)
				inventory.GET("/low-stock", inventoryHandler.GetLowStockItems)
				inventory.GET("/out-of-stock", inventoryHandler.GetOutOfStockItems)
			}

			// Coupon management routes
			adminCoupons := admin.Group("/coupons")
			{
				adminCoupons.GET("", couponHandler.ListCoupons)
				adminCoupons.POST("", couponHandler.CreateCoupon)
				adminCoupons.GET("/:id", couponHandler.GetCoupon)
				adminCoupons.PUT("/:id", couponHandler.UpdateCoupon)
				adminCoupons.DELETE("/:id", couponHandler.DeleteCoupon)
			}

			// Analytics routes
			analytics := admin.Group("/analytics")
			{
				analytics.GET("/sales", analyticsHandler.GetSalesMetrics)
				analytics.GET("/products", analyticsHandler.GetProductMetrics)
				analytics.GET("/users", analyticsHandler.GetUserMetrics)
				analytics.GET("/traffic", analyticsHandler.GetTrafficMetrics)
				analytics.POST("/events", analyticsHandler.TrackEvent)
				analytics.GET("/top-products", analyticsHandler.GetTopProducts)
				analytics.GET("/top-categories", analyticsHandler.GetTopCategories)
			}

			// Reports routes
			reports := admin.Group("/reports")
			{
				reports.POST("/generate", adminHandler.GenerateReport)
				reports.GET("", adminHandler.GetReports)
				reports.GET("/:id/download", adminHandler.DownloadReport)
			}

			// System management routes
			system := admin.Group("/system")
			{
				system.GET("/logs", adminHandler.GetSystemLogs)
				system.GET("/audit", adminHandler.GetAuditLogs)
				system.POST("/backup", adminHandler.BackupDatabase)
			}
		}

		// Moderator routes (moderator/admin authentication required)
		moderator := v1.Group("/moderator")
		moderator.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		moderator.Use(middleware.ModeratorMiddleware())
		{
			// Moderator product management
			modProducts := moderator.Group("/products")
			{
				modProducts.POST("", productHandler.CreateProduct)
				modProducts.PUT("/:id", productHandler.UpdateProduct)          // Complete replacement
				modProducts.PATCH("/:id", productHandler.PatchProduct)         // Partial update  
				modProducts.PUT("/:id/stock", productHandler.UpdateStock)
			}

			// Moderator file uploads
			modUpload := moderator.Group("/upload")
			{
				modUpload.POST("/image", fileHandler.UploadImage)
				modUpload.POST("/document", fileHandler.UploadDocument)
			}
		}
	}
}
