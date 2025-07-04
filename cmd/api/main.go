package main

import (
	"log"

	"ecom-golang-clean-architecture/internal/delivery/http/handlers"
	"ecom-golang-clean-architecture/internal/delivery/http/routes"
	"ecom-golang-clean-architecture/internal/domain/services"
	"ecom-golang-clean-architecture/internal/domain/storage"
	"ecom-golang-clean-architecture/internal/infrastructure/config"
	"ecom-golang-clean-architecture/internal/infrastructure/database"
	"ecom-golang-clean-architecture/internal/infrastructure/oauth"
	"ecom-golang-clean-architecture/internal/infrastructure/payment"
	infraServices "ecom-golang-clean-architecture/internal/infrastructure/services"
	localStorage "ecom-golang-clean-architecture/internal/infrastructure/storage"
	"ecom-golang-clean-architecture/internal/usecases"
	"github.com/gin-gonic/gin"
)

// @title E-commerce API
// @version 1.0
// @description A modern e-commerce API built with Go and Clean Architecture
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Set Gin mode based on environment
	if cfg.App.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database connection
	db, err := database.NewConnection(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run database migrations
	if err := database.AutoMigrate(db); err != nil {
		log.Fatal("Failed to run database migrations:", err)
	}

	// Create database indexes
	if err := database.CreateIndexes(db); err != nil {
		log.Fatal("Failed to create database indexes:", err)
	}

	// Seed initial data
	if err := database.SeedData(db); err != nil {
		log.Fatal("Failed to seed initial data:", err)
	}

	// Initialize repositories
	userRepo := database.NewUserRepository(db)
	userProfileRepo := database.NewUserProfileRepository(db)
	productRepo := database.NewProductRepository(db)
	categoryRepo := database.NewCategoryRepository(db)
	tagRepo := database.NewTagRepository(db)
	imageRepo := database.NewImageRepository(db)
	cartRepo := database.NewCartRepository(db)
	orderRepo := database.NewOrderRepository(db)
	paymentRepo := database.NewPaymentRepository(db)
	fileRepo := database.NewFileRepository(db)
	reviewRepo := database.NewReviewRepository(db)
	reviewVoteRepo := database.NewReviewVoteRepository(db)
	productRatingRepo := database.NewProductRatingRepository(db)
	couponRepo := database.NewCouponRepository(db)
	wishlistRepo := database.NewWishlistRepository(db)
	inventoryRepo := database.NewInventoryRepository(db)
	notificationRepo := database.NewNotificationRepository(db)
	analyticsRepo := database.NewAnalyticsRepository(db)
	addressRepo := database.NewAddressRepository(db)
	shippingRepo := database.NewShippingRepository(db)
	auditRepo := database.NewAuditRepository(db)
	warehouseRepo := database.NewWarehouseRepository(db)

	// Initialize domain services
	passwordService := services.NewPasswordService()
	orderService := services.NewOrderService()

	// Initialize storage service
	fileStorageConfig := config.LoadFileStorageConfig()
	var storageProvider storage.StorageProvider
	var err2 error
	
	// For now, use local storage. In production, this would be configurable
	storageProvider, err2 = localStorage.NewLocalStorage(&fileStorageConfig.LocalConfig)
	if err2 != nil {
		log.Fatal("Failed to initialize storage provider:", err2)
	}

	fileService := services.NewFileService(storageProvider, fileRepo)

	// Initialize use cases
	userUseCase := usecases.NewUserUseCase(
		userRepo,
		userProfileRepo,
		passwordService,
		cfg.JWT.Secret,
	)

	productUseCase := usecases.NewProductUseCase(
		productRepo,
		categoryRepo,
		tagRepo,
		imageRepo,
		cartRepo,
	)

	categoryUseCase := usecases.NewCategoryUseCase(
		categoryRepo,
		fileService,
	)

	cartUseCase := usecases.NewCartUseCase(
		cartRepo,
		productRepo,
	)

	orderUseCase := usecases.NewOrderUseCase(
		orderRepo,
		cartRepo,
		productRepo,
		paymentRepo,
		inventoryRepo,
		orderService,
	)

	fileUseCase := usecases.NewFileUseCase(fileService)

	// Initialize all use cases
	couponUseCase := usecases.NewCouponUseCase(couponRepo, userRepo)
	reviewUseCase := usecases.NewReviewUseCase(reviewRepo, reviewVoteRepo, productRatingRepo, productRepo, orderRepo)
	wishlistUseCase := usecases.NewWishlistUseCase(wishlistRepo, productRepo)
	inventoryUseCase := usecases.NewInventoryUseCase(inventoryRepo, productRepo, warehouseRepo)
	addressUseCase := usecases.NewAddressUseCase(addressRepo)

	// Initialize notification use case (with nil services for now)
	notificationUseCase := usecases.NewNotificationUseCase(
		notificationRepo, userRepo, orderRepo, inventoryRepo,
		nil, nil, nil, // email, sms, push services - TODO: implement
	)

	analyticsUseCase := usecases.NewAnalyticsUseCase(
		analyticsRepo, orderRepo, productRepo, userRepo, inventoryRepo,
	)

	// Initialize payment gateway services
	stripeService := payment.NewStripeService(cfg.Payment.StripeSecretKey)
	paypalService := payment.NewPayPalService(cfg.Payment.PayPalClientID, cfg.Payment.PayPalClientSecret, cfg.Payment.PayPalSandbox)

	// Initialize payment use case
	paymentUseCase := usecases.NewPaymentUseCase(
		paymentRepo, orderRepo, userRepo,
		stripeService, paypalService,
		notificationUseCase,
	)

	// Initialize shipping use case
	shippingUseCase := usecases.NewShippingUseCase(shippingRepo, orderRepo)

	adminUseCase := usecases.NewAdminUseCase(
		userRepo, orderRepo, productRepo, reviewRepo,
		analyticsRepo, inventoryRepo, paymentRepo, auditRepo,
	)

	// Initialize JWT service
	jwtService := infraServices.NewJWTService(cfg.JWT.Secret)

	// Initialize OAuth configuration and service
	oauthConfig := config.NewOAuthConfig()
	oauthService := oauth.NewService(oauthConfig)

	// Initialize OAuth use case
	oauthUseCase := usecases.NewOAuthUseCase(userRepo, oauthService, jwtService)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userUseCase)
	productHandler := handlers.NewProductHandler(productUseCase)
	categoryHandler := handlers.NewCategoryHandler(categoryUseCase)
	cartHandler := handlers.NewCartHandler(cartUseCase)
	orderHandler := handlers.NewOrderHandler(orderUseCase)
	fileHandler := handlers.NewFileHandler(fileUseCase)
	couponHandler := handlers.NewCouponHandler(couponUseCase)
	reviewHandler := handlers.NewReviewHandler(reviewUseCase)
	wishlistHandler := handlers.NewWishlistHandler(wishlistUseCase)
	inventoryHandler := handlers.NewInventoryHandler(inventoryUseCase)
	notificationHandler := handlers.NewNotificationHandler(notificationUseCase)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsUseCase)
	addressHandler := handlers.NewAddressHandler(addressUseCase)
	paymentHandler := handlers.NewPaymentHandler(paymentUseCase)
	shippingHandler := handlers.NewShippingHandler(shippingUseCase)
	adminHandler := handlers.NewAdminHandler(adminUseCase)
	oauthHandler := handlers.NewOAuthHandler(oauthUseCase)

	// Initialize Gin router
	router := gin.New()

	// Setup routes with all handlers
	routes.SetupRoutes(
		router,
		cfg,
		userHandler,
		productHandler,
		categoryHandler,
		cartHandler,
		orderHandler,
		fileHandler,
		reviewHandler,
		wishlistHandler,
		couponHandler,
		inventoryHandler,
		notificationHandler,
		analyticsHandler,
		addressHandler,
		paymentHandler,
		shippingHandler,
		adminHandler,
		oauthHandler,
	)

	// Start server
	log.Printf("Starting server on %s", cfg.App.GetAddress())
	if err := router.Run(cfg.App.GetAddress()); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
