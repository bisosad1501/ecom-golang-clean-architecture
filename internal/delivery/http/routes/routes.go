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
		}

		// Public product routes
		products := v1.Group("/products")
		{
			products.GET("", productHandler.GetProducts)
			products.GET("/:id", productHandler.GetProduct)
			products.GET("/search", productHandler.SearchProducts)
			products.GET("/category/:categoryId", productHandler.GetProductsByCategory)
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
			}

			// Order routes
			orders := protected.Group("/orders")
			{
				orders.POST("", orderHandler.CreateOrder)
				orders.GET("", orderHandler.GetUserOrders)
				orders.GET("/:id", orderHandler.GetOrder)
				orders.POST("/:id/cancel", orderHandler.CancelOrder)
			}
		}

		// Admin routes (admin authentication required)
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		admin.Use(middleware.AdminMiddleware())
		{
			// Admin user management
			adminUsers := admin.Group("/users")
			{
				adminUsers.GET("", userHandler.GetUsers)
				adminUsers.POST("/:id/activate", userHandler.ActivateUser)
				adminUsers.POST("/:id/deactivate", userHandler.DeactivateUser)
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
				adminOrders.GET("", orderHandler.GetOrders)
				adminOrders.PUT("/:id/status", orderHandler.UpdateOrderStatus)
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
