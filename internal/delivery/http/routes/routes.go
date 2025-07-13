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
	brandHandler *handlers.BrandHandler,
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
	migrationHandler *handlers.MigrationHandler,
	searchHandler *handlers.SearchHandler,
	recommendationHandler *handlers.RecommendationHandler,
	comparisonHandler *handlers.ProductComparisonHandler,
	productFilterHandler *handlers.ProductFilterHandler,
) {
	// Apply global middleware
	router.Use(middleware.SecurityHeadersMiddleware())
	router.Use(middleware.RequestSizeLimitMiddleware(10 << 20)) // 10MB limit
	router.Use(middleware.CORSMiddleware(&cfg.CORS))
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware())
	router.Use(middleware.ValidationMiddleware())
	router.Use(middleware.SessionValidationMiddleware())

	// Create auth middleware instance
	authMiddleware := middleware.NewAuthMiddleware(cfg)

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
			auth.POST("/logout", userHandler.Logout)
			auth.POST("/refresh", userHandler.RefreshToken)
			auth.POST("/forgot-password", userHandler.ForgotPassword)
			auth.POST("/reset-password", userHandler.ResetPassword)
			auth.POST("/verify-email", userHandler.VerifyEmailWithToken)
			auth.POST("/resend-verification", userHandler.ResendVerification)

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
			products.GET("/filters", productHandler.GetProductFilters)
			products.GET("/category/:categoryId", productHandler.GetProductsByCategory)
			products.GET("/featured", productHandler.GetFeaturedProducts)
			products.GET("/trending", productHandler.GetTrendingProducts)
			if reviewHandler != nil {
				products.GET("/:id/reviews", reviewHandler.GetProductReviews)
				products.GET("/:id/rating", reviewHandler.GetProductRating)
			}
			products.GET("/:id/related", productHandler.GetRelatedProducts)

			// Product recommendation routes
			if recommendationHandler != nil {
				products.GET("/:id/recommendations", recommendationHandler.GetRelatedProducts)
				products.GET("/:id/frequently-bought-together", recommendationHandler.GetFrequentlyBoughtTogether)
			}

			// Product comparison routes
			if comparisonHandler != nil {
				products.POST("/compare", comparisonHandler.CreateComparison)
				products.GET("/compare/:id", comparisonHandler.GetComparison)
				products.GET("/compare/user", comparisonHandler.GetUserComparison)
				products.GET("/compare/session", comparisonHandler.GetSessionComparison)
				products.PUT("/compare/:id", comparisonHandler.UpdateComparison)
				products.DELETE("/compare/:id", comparisonHandler.DeleteComparison)
				products.POST("/compare/:id/products/:product_id", comparisonHandler.AddProductToComparison)
				products.DELETE("/compare/:id/products/:product_id", comparisonHandler.RemoveProductFromComparison)
				products.POST("/compare/:id/clear", comparisonHandler.ClearComparison)
				products.GET("/compare/matrix", comparisonHandler.CompareProducts)
				products.GET("/compare/:id/matrix", comparisonHandler.GetComparisonMatrix)
				products.GET("/compare/popular", comparisonHandler.GetPopularComparedProducts)
			}

			// Advanced product filtering routes
			if productFilterHandler != nil {
				products.GET("/filter", productFilterHandler.FilterProducts)
				products.GET("/facets", productFilterHandler.GetFilterFacets)
				products.POST("/filters/dynamic", productFilterHandler.GetDynamicFilters)
				products.GET("/filters/suggestions", productFilterHandler.GetFilterSuggestions)
				products.POST("/filters/related", productFilterHandler.GetRelatedFilters)
				products.GET("/attributes", productFilterHandler.GetAttributeFilters)
				products.GET("/attributes/:attribute_id/terms", productFilterHandler.GetAttributeTerms)

				// Filter sets (public access for session-based)
				products.GET("/filter-sets/session", productFilterHandler.GetSessionFilterSets)
				products.GET("/filter-sets/:id", productFilterHandler.GetFilterSet)
			}

			// Search autocomplete and suggestions
			products.GET("/suggestions", productHandler.GetSearchSuggestions)
			products.GET("/popular-searches", productHandler.GetPopularSearches)

			// Authenticated routes
			authProducts := products.Group("")
			authProducts.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
			{
				authProducts.GET("/search-history", productHandler.GetSearchHistory)

				// Authenticated filter set routes
				if productFilterHandler != nil {
					authProducts.POST("/filter-sets", productFilterHandler.SaveFilterSet)
					authProducts.GET("/filter-sets/user", productFilterHandler.GetUserFilterSets)
					authProducts.PUT("/filter-sets/:id", productFilterHandler.UpdateFilterSet)
					authProducts.DELETE("/filter-sets/:id", productFilterHandler.DeleteFilterSet)
				}
			}
		}

		// Public category routes
		categories := v1.Group("/categories")
		{
			categories.GET("", categoryHandler.GetCategories)
			categories.GET("/:id", categoryHandler.GetCategory)
			categories.GET("/tree", categoryHandler.GetCategoryTree)
			categories.GET("/root", categoryHandler.GetRootCategories)
			categories.GET("/:id/children", categoryHandler.GetCategoryChildren)
			categories.GET("/:id/path", categoryHandler.GetCategoryPath)
			categories.GET("/:id/count", categoryHandler.GetCategoryProductCount)
			categories.GET("/:id/landing", categoryHandler.GetCategoryLandingPage)

			// Enhanced category routes
			categories.GET("/search", categoryHandler.SearchCategories)
			categories.GET("/filter", categoryHandler.GetCategoriesWithFilters)

			// SEO routes (public access for frontend)
			categories.GET("/:id/seo", categoryHandler.GetCategorySEO)
			categories.GET("/slug/validate", categoryHandler.ValidateSlugAvailability)
		}

		// Public search routes
		if searchHandler != nil {
			search := v1.Group("/search")
			{
				search.GET("", searchHandler.FullTextSearch)
				search.GET("/enhanced", searchHandler.EnhancedSearch)
				search.GET("/suggestions", searchHandler.GetSearchSuggestions)
				search.GET("/facets", searchHandler.GetSearchFacets)
				search.GET("/popular", searchHandler.GetPopularSearchTerms)
				search.GET("/autocomplete", searchHandler.GetAutocomplete)
				search.GET("/autocomplete/enhanced", searchHandler.GetEnhancedAutocomplete)
				search.GET("/autocomplete/smart", searchHandler.GetSmartAutocomplete)
				search.POST("/autocomplete/track", searchHandler.TrackAutocompleteInteraction)
				search.GET("/trending", searchHandler.GetTrendingSearches)
				search.GET("/trends", searchHandler.GetSearchTrends)
				search.POST("/record", searchHandler.RecordSearchEvent)
				search.POST("/autocomplete/click", searchHandler.RecordAutocompleteClick)

				// Authenticated search routes
				authSearch := search.Group("")
				authSearch.Use(authMiddleware.RequireAuth())
				{
					authSearch.GET("/autocomplete/personalized", searchHandler.GetPersonalizedAutocomplete)
					authSearch.GET("/preferences", searchHandler.GetUserSearchPreferences)
					authSearch.PUT("/preferences", searchHandler.UpdateUserSearchPreferences)
				}
			}
		}

		// Public recommendation routes
		if recommendationHandler != nil {
			recommendations := v1.Group("/recommendations")
			{
				recommendations.GET("", recommendationHandler.GetRecommendations)
				recommendations.GET("/trending", recommendationHandler.GetTrendingProducts)
				recommendations.POST("/track", recommendationHandler.TrackInteraction)
			}
		}

		// Public brand routes
		brands := v1.Group("/brands")
		{
			brands.GET("", brandHandler.GetBrands)
			brands.GET("/active", brandHandler.GetActiveBrands)
			brands.GET("/popular", brandHandler.GetPopularBrands)
			brands.GET("/search", brandHandler.SearchBrands)
			brands.GET("/:id", brandHandler.GetBrand)
			brands.GET("/slug/:slug", brandHandler.GetBrandBySlug)
		}

		// Public cart routes (guest cart support)
		publicCart := v1.Group("/public/cart")
		publicCart.Use(middleware.SessionValidationMiddleware())
		{
			publicCart.GET("", cartHandler.GetCart)
			publicCart.POST("/items", cartHandler.AddToCart)
			publicCart.PUT("/items", cartHandler.UpdateCartItem)
			publicCart.DELETE("/items/:productId", cartHandler.RemoveFromCart)
			publicCart.DELETE("", cartHandler.ClearCart)
		}

		// Public file upload routes (requires authentication, with strict rate limiting)
		publicUpload := v1.Group("/public/upload")
		publicUpload.Use(middleware.PublicUploadRateLimitMiddleware())
		publicUpload.Use(middleware.PublicUploadAuthMiddleware(cfg.JWT.Secret))
		publicUpload.Use(middleware.FileUploadSecurityMiddleware())
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

		// Public order access for success page
		publicOrders := v1.Group("/orders")
		{
			publicOrders.GET("/:id/public", orderHandler.GetOrderPublic)
		}

		// Public payment confirmation for success page
		publicPayments := v1.Group("/payments")
		{
			publicPayments.POST("/confirm-success", paymentHandler.ConfirmPaymentSuccess)
		}

		// Public review routes (no authentication required)
		if reviewHandler != nil {
			publicReviews := v1.Group("/public/reviews")
			{
				publicReviews.GET("/product/:product_id", reviewHandler.GetProductReviews)
				publicReviews.GET("/product/:product_id/summary", reviewHandler.GetProductRating)
				publicReviews.GET("/:id", reviewHandler.GetReview)
			}
		}

		// Public verification routes (no authentication required)
		publicVerification := v1.Group("/public/verification")
		{
			publicVerification.POST("/email/verify", userHandler.VerifyEmail)
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

				// User preferences routes
				users.GET("/preferences", userHandler.GetUserPreferences)
				users.PUT("/preferences", userHandler.UpdateUserPreferences)
				users.PUT("/preferences/theme", userHandler.UpdateTheme)
				users.PUT("/preferences/language", userHandler.UpdateLanguage)

				// Search history routes
				searchHistory := users.Group("/search-history")
				{
					searchHistory.POST("/track", userHandler.TrackSearch)
					searchHistory.GET("", userHandler.GetSearchHistory)
					searchHistory.DELETE("", userHandler.ClearSearchHistory)
				}

				// Saved searches routes
				savedSearches := users.Group("/saved-searches")
				{
					savedSearches.POST("", userHandler.CreateSavedSearch)
					savedSearches.GET("", userHandler.GetSavedSearches)
				}

				// Browsing history routes
				browsingHistory := users.Group("/browsing-history")
				{
					browsingHistory.POST("/track", userHandler.TrackProductView)
					browsingHistory.GET("", userHandler.GetBrowsingHistory)
				}

				// Personalization routes
				personalization := users.Group("/personalization")
				{
					personalization.GET("", userHandler.GetPersonalization)
				}

				// User verification routes
				verification := users.Group("/verification")
				{
					verification.POST("/email/send", userHandler.SendEmailVerification)
					verification.POST("/phone/send", userHandler.SendPhoneVerification)
					verification.POST("/phone/verify", userHandler.VerifyPhone)
					verification.GET("/status", userHandler.GetVerificationStatus)
				}

				// User session routes
				sessions := users.Group("/sessions")
				{
					sessions.GET("", userHandler.GetUserSessions)
					sessions.DELETE("/:session_id", userHandler.InvalidateSession)
					sessions.DELETE("", userHandler.InvalidateAllSessions)
				}

				if reviewHandler != nil {
					users.GET("/:user_id/reviews", reviewHandler.GetUserReviews)
				}
			}

			// Protected search routes (authentication required)
			if searchHandler != nil {
				searchProtected := protected.Group("/search")
				{
					searchProtected.POST("/history", searchHandler.SaveSearchHistory)
					searchProtected.GET("/history", searchHandler.GetUserSearchHistory)
					searchProtected.DELETE("/history", searchHandler.ClearUserSearchHistory)

					searchProtected.POST("/filters", searchHandler.SaveSearchFilter)
					searchProtected.GET("/filters", searchHandler.GetUserSearchFilters)
					searchProtected.PUT("/filters/:id", searchHandler.UpdateSearchFilter)
					searchProtected.DELETE("/filters/:id", searchHandler.DeleteSearchFilter)
				}
			}

			// Protected recommendation routes (authentication required)
			if recommendationHandler != nil {
				recommendationsProtected := protected.Group("/recommendations")
				{
					recommendationsProtected.GET("/personalized", recommendationHandler.GetPersonalizedRecommendations)
				}
			}

			// Upload routes (authenticated users)
			upload := protected.Group("/upload")
			upload.Use(middleware.UploadRateLimitMiddleware())
			upload.Use(middleware.FileUploadSecurityMiddleware())
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
				cart.POST("/merge", cartHandler.MergeGuestCart)
				cart.POST("/check-conflict", cartHandler.CheckCartConflict)
				// cart.POST("/sync", cartHandler.SyncCart) // TODO: Implement SyncCart method
			}

			// Order routes
			orders := protected.Group("/orders")
			{
				orders.POST("", orderHandler.CreateOrder)
				orders.GET("", orderHandler.GetUserOrders)
				orders.GET("/by-session", orderHandler.GetOrderBySessionID)
				orders.GET("/:id", orderHandler.GetOrder)
				orders.POST("/:id/cancel", orderHandler.CancelOrder)
				orders.GET("/:id/events", orderHandler.GetOrderEvents)
				orders.POST("/:id/notes", orderHandler.AddOrderNote)
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

			// Webhook routes (public - no auth required)
			webhooks := v1.Group("/webhooks")
			{
				webhooks.POST("/payment/:provider", paymentHandler.HandleWebhook)
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
				dashboard.GET("/activity", adminHandler.GetRecentActivity)
			}

			// Admin user management
			adminUsers := admin.Group("/users")
			{
				adminUsers.GET("", adminHandler.GetUsers)
				adminUsers.PUT("/:id/status", adminHandler.UpdateUserStatus)
				adminUsers.PUT("/:id/role", adminHandler.UpdateUserRole)
				adminUsers.GET("/:id/activity", adminHandler.GetUserActivity)

				// Bulk user operations
				adminUsers.POST("/bulk/update", adminHandler.BulkUpdateUsers)
				adminUsers.POST("/bulk/delete", adminHandler.BulkDeleteUsers)
				adminUsers.POST("/bulk/activate", adminHandler.BulkActivateUsers)
				adminUsers.POST("/bulk/deactivate", adminHandler.BulkDeactivateUsers)
				adminUsers.POST("/bulk/roles", adminHandler.BulkUpdateUserRoles)

				// User communication
				adminUsers.POST("/notification", adminHandler.SendUserNotification)
				adminUsers.POST("/bulk/notification", adminHandler.SendBulkNotification)
				adminUsers.POST("/email", adminHandler.SendUserEmail)
				adminUsers.POST("/bulk/email", adminHandler.SendBulkEmail)

				// User audit logs
				adminUsers.GET("/audit-logs", adminHandler.GetUserAuditLogs)

				// Announcements
				adminUsers.POST("/announcements", adminHandler.CreateAnnouncement)
			}

			// Admin customer management and segmentation
			adminCustomers := admin.Group("/customers")
			{
				adminCustomers.GET("/search", adminHandler.SearchCustomers)
				adminCustomers.GET("/segments", adminHandler.GetCustomerSegments)
				adminCustomers.GET("/segments/:segment", adminHandler.GetCustomersBySegment)
				adminCustomers.GET("/analytics", adminHandler.GetCustomerAnalytics)
				adminCustomers.GET("/high-value", adminHandler.GetHighValueCustomers)
				adminCustomers.GET("/:customer_id/lifetime-value", adminHandler.GetCustomerLifetimeValue)
			}

			// Admin product management
			adminProducts := admin.Group("/products")
			{
				adminProducts.GET("", productHandler.GetProducts)
				adminProducts.POST("", productHandler.CreateProduct)
				adminProducts.PUT("/:id", productHandler.UpdateProduct)  // Complete replacement
				adminProducts.PATCH("/:id", productHandler.PatchProduct) // Partial update
				adminProducts.DELETE("/:id", productHandler.DeleteProduct)
				adminProducts.PUT("/:id/stock", productHandler.UpdateStock)
			}

			// Admin category management
			adminCategories := admin.Group("/categories")
			{
				adminCategories.POST("", categoryHandler.CreateCategory)
				adminCategories.PUT("/:id", categoryHandler.UpdateCategory)
				adminCategories.DELETE("/:id", categoryHandler.DeleteCategory)

				// Bulk operations
				adminCategories.POST("/bulk", categoryHandler.BulkCreateCategories)
				adminCategories.PUT("/bulk", categoryHandler.BulkUpdateCategories)
				adminCategories.DELETE("/bulk", categoryHandler.BulkDeleteCategories)

				// Tree operations
				adminCategories.POST("/move", categoryHandler.MoveCategory)
				adminCategories.POST("/reorder", categoryHandler.ReorderCategories)
				adminCategories.GET("/tree/stats", categoryHandler.GetCategoryTreeStats)
				adminCategories.POST("/tree/validate", categoryHandler.ValidateAndRepairTree)

				// Analytics and statistics
				adminCategories.GET("/top", categoryHandler.GetTopCategories)
				adminCategories.GET("/:id/analytics", categoryHandler.GetCategoryAnalytics)
				adminCategories.GET("/:id/performance", categoryHandler.GetCategoryPerformanceMetrics)
				adminCategories.GET("/:id/sales", categoryHandler.GetCategorySalesStats)

				// SEO management
				adminCategories.PUT("/:id/seo", categoryHandler.UpdateCategorySEO)
				adminCategories.POST("/:id/seo/generate", categoryHandler.GenerateCategorySEO)
				adminCategories.GET("/:id/seo/validate", categoryHandler.ValidateCategorySEO)

				// Enhanced SEO & URL optimization
				adminCategories.POST("/:id/slug/optimize", categoryHandler.OptimizeSlug)
				adminCategories.GET("/:id/slug/suggestions", categoryHandler.GenerateSlugSuggestions)
				adminCategories.GET("/:id/slug/history", categoryHandler.GetSlugHistory)
				adminCategories.GET("/:id/seo/insights", categoryHandler.GetSEOInsights)
				adminCategories.GET("/:id/seo/competitor-analysis", categoryHandler.GetSEOCompetitorAnalysis)

				// Bulk SEO operations
				adminCategories.POST("/seo/bulk-update", categoryHandler.BulkUpdateSEO)
				adminCategories.POST("/seo/bulk-generate", categoryHandler.BulkGenerateSEO)
				adminCategories.POST("/seo/bulk-validate", categoryHandler.BulkValidateSEO)
				adminCategories.GET("/seo/analytics", categoryHandler.GetSEOAnalytics)
			}

			// Admin brand management
			adminBrands := admin.Group("/brands")
			{
				adminBrands.POST("", brandHandler.CreateBrand)
				adminBrands.PUT("/:id", brandHandler.UpdateBrand)
				adminBrands.DELETE("/:id", brandHandler.DeleteBrand)
			}

			// Admin file uploads
			adminUpload := admin.Group("/upload")
			adminUpload.Use(middleware.UploadRateLimitMiddleware())
			adminUpload.Use(middleware.FileUploadSecurityMiddleware())
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
				adminOrders.PATCH("/:id/status", adminHandler.UpdateOrderStatus) // Add PATCH route
				adminOrders.PUT("/:id/shipping", orderHandler.UpdateShippingInfo)
				adminOrders.PUT("/:id/delivery", orderHandler.UpdateDeliveryStatus)
				adminOrders.POST("/:id/notes", orderHandler.AddOrderNote)
				adminOrders.GET("/:id/events", orderHandler.GetOrderEvents)
				adminOrders.POST("/:id/refund", adminHandler.ProcessRefund)
			}

			// Review management routes
			adminReviews := admin.Group("/reviews")
			{
				adminReviews.GET("", adminHandler.ManageReviews)
				adminReviews.PUT("/:id/status", adminHandler.UpdateReviewStatus)
				adminReviews.POST("/:id/reply", adminHandler.ReplyToReview)
			}

			// Admin search management routes
			if searchHandler != nil {
				adminSearch := admin.Group("/search")
				{
					adminSearch.GET("/analytics", searchHandler.GetSearchAnalytics)
					adminSearch.GET("/popular", searchHandler.GetPopularSearchTerms)
					adminSearch.POST("/rebuild-index", searchHandler.RebuildAutocompleteIndex)
					adminSearch.POST("/cleanup", searchHandler.CleanupSearchData)
				}
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

				// Filter analytics
				if productFilterHandler != nil {
					analytics.GET("/filters", productFilterHandler.GetFilterAnalytics)
					analytics.GET("/filters/popular", productFilterHandler.GetPopularFilters)
				}
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
				system.GET("/cleanup/stats", adminHandler.GetCleanupStats)
				system.POST("/cleanup/trigger", adminHandler.TriggerCleanup)
			}

			// Migration management routes
			migrations := admin.Group("/migrations")
			{
				migrations.GET("/status", migrationHandler.GetMigrationStatus)
				migrations.POST("/run", migrationHandler.RunMigrations)
				migrations.POST("/rollback", migrationHandler.RollbackMigration)
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
				modProducts.PUT("/:id", productHandler.UpdateProduct)  // Complete replacement
				modProducts.PATCH("/:id", productHandler.PatchProduct) // Partial update
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
