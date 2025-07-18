package main

import (
	"context"
	"log"
	"time"

	"ecom-golang-clean-architecture/internal/delivery/http/handlers"
	"ecom-golang-clean-architecture/internal/delivery/http/routes"
	"ecom-golang-clean-architecture/internal/domain/services"
	"ecom-golang-clean-architecture/internal/domain/storage"
	"ecom-golang-clean-architecture/internal/infrastructure/config"
	"ecom-golang-clean-architecture/internal/infrastructure/database"
	"ecom-golang-clean-architecture/internal/infrastructure/oauth"
	"ecom-golang-clean-architecture/internal/infrastructure/payment"
	"ecom-golang-clean-architecture/internal/infrastructure/repositories"
	infraServices "ecom-golang-clean-architecture/internal/infrastructure/services"
	localStorage "ecom-golang-clean-architecture/internal/infrastructure/storage"
	"ecom-golang-clean-architecture/internal/infrastructure/websocket"
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

	// Initialize default email templates
	emailTemplateRepo := database.NewEmailTemplateRepository(db)
	emailTemplateService := infraServices.NewEmailTemplateService(emailTemplateRepo)
	if err := emailTemplateService.InitializeDefaultTemplates(ctx); err != nil {
		log.Printf("⚠️ Failed to initialize default email templates: %v", err)
	} else {
		log.Printf("✅ Default email templates initialized")
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
	categoryRepo := database.NewCategoryRepository(db)
	productCategoryRepo := repositories.NewProductCategoryRepository(db)
	// Initialize category hierarchy service for optimized category queries
	categoryHierarchyService := services.NewCategoryHierarchyService(categoryRepo)
	productRepo := database.NewProductRepository(db, categoryHierarchyService)
	brandRepo := database.NewBrandRepository(db)
	tagRepo := database.NewTagRepository(db)
	imageRepo := database.NewImageRepository(db)
	cartRepo := database.NewCartRepository(db)
	orderRepo := database.NewOrderRepository(db)
	checkoutRepo := repositories.NewCheckoutSessionRepository(db)
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
	orderEventRepo := database.NewOrderEventRepository(db)

	// Initialize transaction manager
	txManager := database.NewTransactionManager(db)

	// Initialize domain services
	passwordService := services.NewPasswordService()
	orderService := services.NewOrderService(orderRepo)
	simpleStockService := services.NewSimpleStockService(productRepo, inventoryRepo)
	userMetricsService := services.NewUserMetricsService(userRepo, orderRepo)
	_ = services.NewProductCategoryService(productCategoryRepo, productRepo, categoryRepo) // Will be used later
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

	// Initialize Gmail service
	gmailService := infraServices.NewGmailService(&cfg.Email)

	// Validate Gmail configuration
	if err := gmailService.ValidateConfiguration(); err != nil {
		log.Printf("⚠️ Gmail configuration validation failed: %v", err)
		log.Printf("📧 Email functionality will use fallback console logging")
	} else {
		log.Printf("✅ Gmail service configured successfully")
	}

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
		gmailService,
		nil, // notificationService - will be set later
		cfg.JWT.Secret,
	)

	productUseCase := usecases.NewProductUseCase(
		productRepo,
		categoryRepo,
		productCategoryRepo,
		tagRepo,
		imageRepo,
		cartRepo,
		inventoryRepo,
		warehouseRepo,
	)

	categoryUseCase := usecases.NewCategoryUseCase(
		categoryRepo,
		productRepo,
		productCategoryRepo,
		fileService,
	)

	brandUseCase := usecases.NewBrandUseCase(
		brandRepo,
	)

	cartUseCase := usecases.NewCartUseCase(
		cartRepo,
		productRepo,
		simpleStockService, // Use simple stock service instead
	)

	// Initialize WebSocket hub for real-time notifications
	websocketHub := websocket.NewHub(context.Background())

	// Start WebSocket hub in background
	go websocketHub.Run()

	// Initialize notification use case with WebSocket hub
	notificationUseCase := usecases.NewNotificationUseCase(
		notificationRepo, userRepo, orderRepo, paymentRepo, inventoryRepo,
		reviewRepo, productRepo,
		nil, nil, nil, // email, sms, push services - TODO: implement
		websocketHub,  // WebSocket hub for real-time notifications
	)

	// Re-initialize userUseCase with notificationUseCase
	userUseCase = usecases.NewUserUseCase(
		userRepo,
		userProfileRepo,
		userSessionRepo,
		userLoginHistoryRepo,
		userActivityRepo,
		userPreferencesRepo,
		userVerificationRepo,
		passwordResetRepo,
		passwordService,
		gmailService,
		notificationUseCase, // Now we have notificationUseCase
		cfg.JWT.Secret,
	)

	// Initialize notification queue processor
	queueProcessor := infraServices.NewNotificationQueueProcessor(
		notificationRepo,
		notificationUseCase,
		1,              // workers (reduced to 1 to avoid race conditions)
		5,              // batch size (reduced for better control)
		10*time.Second, // poll interval (reduced for faster processing)
		2*time.Minute,  // retry interval
		3,              // max retries
	)

	// Initialize payment gateway services
	stripeService := payment.NewStripeServiceWithWebhook(cfg.Payment.StripeSecretKey, cfg.Payment.StripeWebhookSecret)
	paypalService := payment.NewPayPalService(cfg.Payment.PayPalClientID, cfg.Payment.PayPalClientSecret, cfg.Payment.PayPalSandbox)

	// Initialize payment use case
	paymentUseCase := usecases.NewPaymentUseCase(
		paymentRepo, paymentMethodRepo, orderRepo, userRepo,
		stripeService, paypalService,
		notificationUseCase,
		orderEventService,
		userMetricsService,
		txManager,
		simpleStockService,
	)

	orderUseCase := usecases.NewOrderUseCase(
		orderRepo,
		cartRepo,
		productRepo,
		paymentRepo,
		inventoryRepo,
		orderEventRepo,
		orderService,
		simpleStockService,
		orderEventService,
		userMetricsService,
		notificationUseCase, // Pass notification service
		txManager,
	)

	checkoutUseCase := usecases.NewCheckoutUseCase(
		checkoutRepo,
		cartRepo,
		orderRepo,
		productRepo,
		simpleStockService,
		orderService,
		paymentUseCase,
		txManager,
	)

	fileUseCase := usecases.NewFileUseCase(fileService)

	// Initialize all use cases
	couponUseCase := usecases.NewCouponUseCase(couponRepo, userRepo)
	reviewUseCase := usecases.NewReviewUseCase(reviewRepo, reviewVoteRepo, productRatingRepo, productRepo, orderRepo, userRepo, notificationUseCase)
	wishlistUseCase := usecases.NewWishlistUseCase(wishlistRepo, productRepo, productCategoryRepo)
	inventoryUseCase := usecases.NewInventoryUseCase(inventoryRepo, productRepo, warehouseRepo, notificationUseCase)
	addressUseCase := usecases.NewAddressUseCase(addressRepo)

	analyticsUseCase := usecases.NewAnalyticsUseCase(
		analyticsRepo, orderRepo, productRepo, userRepo, inventoryRepo,
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
		userLoginHistoryRepo, orderUseCase,
	)

	// Initialize email use case (with nil repositories for now)
	emailUseCase := usecases.NewEmailUseCase(
		nil, nil, nil, nil, // email service, repo, template repo, subscription repo - TODO: implement
		userRepo, orderRepo, productRepo,
	)

	// Initialize abandoned cart use case
	abandonedCartUseCase := usecases.NewAbandonedCartUseCase(
		cartRepo, userRepo, emailUseCase, productRepo, orderRepo,
	)

	// Initialize stock cleanup use case - DEPRECATED (using simple stock service now)
	// stockCleanupUseCase := usecases.NewStockCleanupUseCase(
	//	stockReservationService,
	//	orderRepo,
	//	stockReservationRepo,
	//	cartRepo, // Pass the cartRepo
	// )

	// Initialize JWT service
	jwtService := infraServices.NewJWTService(cfg.JWT.Secret)

	// Initialize OAuth configuration and service
	oauthConfig := config.NewOAuthConfig()
	oauthService := oauth.NewService(oauthConfig)

	// Initialize OAuth use case
	oauthUseCase := usecases.NewOAuthUseCase(userRepo, oauthService, jwtService)

	// Initialize search repository and use case
	searchRepo := database.NewSearchRepository(db)
	searchUseCase := usecases.NewSearchUseCase(searchRepo, productRepo, productCategoryRepo)

	// Initialize recommendation repository and use case
	recommendationRepo := database.NewRecommendationRepository(db)
	recommendationUseCase := usecases.NewRecommendationUseCase(recommendationRepo, productRepo, userRepo)

	// Initialize product comparison system
	comparisonRepo := database.NewProductComparisonRepository(db)
	comparisonUseCase := usecases.NewProductComparisonUseCase(comparisonRepo, productRepo, productCategoryRepo)

	// Initialize advanced product filtering system
	productFilterRepo := database.NewProductFilterRepository(db)
	productFilterUseCase := usecases.NewProductFilterUseCase(productFilterRepo, productRepo, productCategoryRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userUseCase)
	productHandler := handlers.NewProductHandler(productUseCase)
	categoryHandler := handlers.NewCategoryHandler(categoryUseCase)
	brandHandler := handlers.NewBrandHandler(brandUseCase)
	cartHandler := handlers.NewCartHandler(cartUseCase)
	orderHandler := handlers.NewOrderHandler(orderUseCase)
	checkoutHandler := handlers.NewCheckoutHandler(checkoutUseCase)
	fileHandler := handlers.NewFileHandler(fileUseCase)
	couponHandler := handlers.NewCouponHandler(couponUseCase)
	reviewHandler := handlers.NewReviewHandler(reviewUseCase, fileUseCase)
	wishlistHandler := handlers.NewWishlistHandler(wishlistUseCase)
	inventoryHandler := handlers.NewInventoryHandler(inventoryUseCase)
	notificationHandler := handlers.NewNotificationHandler(notificationUseCase)
	websocketHandler := handlers.NewWebSocketHandler(websocketHub)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsUseCase)
	addressHandler := handlers.NewAddressHandler(addressUseCase)
	paymentHandler := handlers.NewPaymentHandler(paymentUseCase)
	shippingHandler := handlers.NewShippingHandler(shippingUseCase)
	adminHandler := handlers.NewAdminHandler(adminUseCase)
	oauthHandler := handlers.NewOAuthHandler(oauthUseCase)
	migrationHandler := handlers.NewMigrationHandler(db)
	searchHandler := handlers.NewSearchHandler(searchUseCase)
	recommendationHandler := handlers.NewRecommendationHandler(recommendationUseCase)
	comparisonHandler := handlers.NewProductComparisonHandler(comparisonUseCase)
	productFilterHandler := handlers.NewProductFilterHandler(productFilterUseCase)
	abandonedCartHandler := handlers.NewAbandonedCartHandler(abandonedCartUseCase)

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
		checkoutHandler,
		fileHandler,
		reviewHandler,
		wishlistHandler,
		couponHandler,
		inventoryHandler,
		notificationHandler,
		websocketHandler,
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
		abandonedCartHandler,
	)

	// Background cleanup scheduler removed - using simple stock service

	// Start notification queue processor
	go func() {
		ctx := context.Background()
		if err := queueProcessor.Start(ctx); err != nil {
			log.Printf("Failed to start notification queue processor: %v", err)
		}
	}()

	// Start server
	log.Printf("Starting server on %s", cfg.App.GetAddress())
	if err := router.Run(cfg.App.GetAddress()); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
