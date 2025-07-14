package main

import (
	"context"
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

	// Run database migrations using Migration Manager
	migrationManager := database.NewMigrationManager(db)
	ctx := context.Background()
	if err := migrationManager.RunMigrations(ctx); err != nil {
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
	userSessionRepo := database.NewUserSessionRepository(db)
	userLoginHistoryRepo := database.NewUserLoginHistoryRepository(db)
	userActivityRepo := database.NewUserActivityRepository(db)
	userPreferencesRepo := database.NewUserPreferencesRepository(db)
	userVerificationRepo := database.NewUserVerificationRepository(db)
	passwordResetRepo := database.NewPasswordResetRepository(db)
	productRepo := database.NewProductRepository(db)
	categoryRepo := database.NewCategoryRepository(db)
	brandRepo := database.NewBrandRepository(db)
	tagRepo := database.NewTagRepository(db)
	imageRepo := database.NewImageRepository(db)
	cartRepo := database.NewCartRepository(db)
	orderRepo := database.NewOrderRepository(db)
	paymentRepo := database.NewPaymentRepository(db)
	paymentMethodRepo := database.NewPaymentMethodRepository(db)
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
	stockReservationRepo := database.NewStockReservationRepository(db)
	orderEventRepo := database.NewOrderEventRepository(db)

	// Initialize transaction manager
	txManager := database.NewTransactionManager(db)

	// Initialize domain services
	passwordService := services.NewPasswordService()
	orderService := services.NewOrderService(orderRepo)
	stockReservationService := services.NewStockReservationService(
		stockReservationRepo,
		productRepo,
		inventoryRepo,
	)
	orderEventService := services.NewOrderEventService(orderEventRepo)

	// Initialize storage service
	fileStorageConfig := config.LoadFileStorageConfig()
	var storageProvider storage.StorageProvider
	var err2 error

	// For now, use local storage. In production, this would be configurable
	storageProvider, err2 = localStorage.NewLocalStorage(&fileStorageConfig.LocalConfig)
	if err2 != nil {
		log.Fatal("Failed to initialize storage provider:", err2)
	}

	// Initialize file security service
	fileSecurityService := services.NewFileSecurityService()

	fileService := services.NewFileService(storageProvider, fileRepo, fileSecurityService)

	// Initialize use cases
	userUseCase := usecases.NewUserUseCase(
		userRepo,
		userProfileRepo,
		userSessionRepo,
		userLoginHistoryRepo,
		userActivityRepo,
		userPreferencesRepo,
		userVerificationRepo,
		passwordResetRepo,
		passwordService,
		cfg.JWT.Secret,
	)

	productUseCase := usecases.NewProductUseCase(
		productRepo,
		categoryRepo,
		tagRepo,
		imageRepo,
		cartRepo,
		inventoryRepo,
		warehouseRepo,
	)

	categoryUseCase := usecases.NewCategoryUseCase(
		categoryRepo,
		productRepo,
		fileService,
	)

	brandUseCase := usecases.NewBrandUseCase(
		brandRepo,
	)

	cartUseCase := usecases.NewCartUseCase(
		cartRepo,
		productRepo,
		stockReservationService, // Pass the stockReservationService
	)

	orderUseCase := usecases.NewOrderUseCase(
		orderRepo,
		cartRepo,
		productRepo,
		paymentRepo,
		inventoryRepo,
		stockReservationRepo,
		orderEventRepo,
		orderService,
		stockReservationService,
		orderEventService,
		txManager,
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
	stripeService := payment.NewStripeServiceWithWebhook(cfg.Payment.StripeSecretKey, cfg.Payment.StripeWebhookSecret)
	paypalService := payment.NewPayPalService(cfg.Payment.PayPalClientID, cfg.Payment.PayPalClientSecret, cfg.Payment.PayPalSandbox)

	// Initialize payment use case
	paymentUseCase := usecases.NewPaymentUseCase(
		paymentRepo, paymentMethodRepo, orderRepo, userRepo,
		stripeService, paypalService,
		notificationUseCase,
		stockReservationService,
		orderEventService,
		txManager,
	)

	// Initialize distance service
	distanceService := services.NewDistanceService()

	// Initialize shipping compatibility service
	compatibilityService := services.NewShippingCompatibilityService()

	// Initialize shipping use case
	shippingUseCase := usecases.NewShippingUseCase(shippingRepo, orderRepo, distanceService, compatibilityService)

	adminUseCase := usecases.NewAdminUseCase(
		userRepo, orderRepo, productRepo, reviewRepo,
		analyticsRepo, inventoryRepo, paymentRepo, auditRepo,
		orderUseCase,
	)

	// Initialize stock cleanup use case
	stockCleanupUseCase := usecases.NewStockCleanupUseCase(
		stockReservationService,
		orderRepo,
		stockReservationRepo,
		cartRepo, // Pass the cartRepo
	)

	// Initialize JWT service
	jwtService := infraServices.NewJWTService(cfg.JWT.Secret)

	// Initialize OAuth configuration and service
	oauthConfig := config.NewOAuthConfig()
	oauthService := oauth.NewService(oauthConfig)

	// Initialize OAuth use case
	oauthUseCase := usecases.NewOAuthUseCase(userRepo, oauthService, jwtService)

	// Initialize search repository and use case
	searchRepo := database.NewSearchRepository(db)
	searchUseCase := usecases.NewSearchUseCase(searchRepo, productRepo)

	// Initialize recommendation repository and use case
	recommendationRepo := database.NewRecommendationRepository(db)
	recommendationUseCase := usecases.NewRecommendationUseCase(recommendationRepo, productRepo, userRepo)

	// Initialize product comparison system
	comparisonRepo := database.NewProductComparisonRepository(db)
	comparisonUseCase := usecases.NewProductComparisonUseCase(comparisonRepo, productRepo)

	// Initialize advanced product filtering system
	productFilterRepo := database.NewProductFilterRepository(db)
	productFilterUseCase := usecases.NewProductFilterUseCase(productFilterRepo, productRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userUseCase)
	productHandler := handlers.NewProductHandler(productUseCase)
	categoryHandler := handlers.NewCategoryHandler(categoryUseCase)
	brandHandler := handlers.NewBrandHandler(brandUseCase)
	cartHandler := handlers.NewCartHandler(cartUseCase)
	orderHandler := handlers.NewOrderHandler(orderUseCase)
	fileHandler := handlers.NewFileHandler(fileUseCase)
	couponHandler := handlers.NewCouponHandler(couponUseCase)
	reviewHandler := handlers.NewReviewHandler(reviewUseCase, fileUseCase)
	wishlistHandler := handlers.NewWishlistHandler(wishlistUseCase)
	inventoryHandler := handlers.NewInventoryHandler(inventoryUseCase)
	notificationHandler := handlers.NewNotificationHandler(notificationUseCase)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsUseCase)
	addressHandler := handlers.NewAddressHandler(addressUseCase)
	paymentHandler := handlers.NewPaymentHandler(paymentUseCase)
	shippingHandler := handlers.NewShippingHandler(shippingUseCase)
	adminHandler := handlers.NewAdminHandler(adminUseCase, stockCleanupUseCase)
	oauthHandler := handlers.NewOAuthHandler(oauthUseCase)
	migrationHandler := handlers.NewMigrationHandler(db)
	searchHandler := handlers.NewSearchHandler(searchUseCase)
	recommendationHandler := handlers.NewRecommendationHandler(recommendationUseCase)
	comparisonHandler := handlers.NewProductComparisonHandler(comparisonUseCase)
	productFilterHandler := handlers.NewProductFilterHandler(productFilterUseCase)

	// Initialize Gin router
	router := gin.New()

	// Setup routes with all handlers
	routes.SetupRoutes(
		router,
		cfg,
		userHandler,
		productHandler,
		categoryHandler,
		brandHandler,
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
		migrationHandler,
		searchHandler,
		recommendationHandler,
		comparisonHandler,
		productFilterHandler,
	)

	// Start background cleanup scheduler
	go func() {
		ctx := context.Background()
		usecases.StartCleanupScheduler(ctx, stockCleanupUseCase)
	}()

	// Start server
	log.Printf("Starting server on %s", cfg.App.GetAddress())
	if err := router.Run(cfg.App.GetAddress()); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
