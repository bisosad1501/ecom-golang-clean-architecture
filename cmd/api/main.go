package main

import (
	"log"

	"ecom-golang-clean-architecture/internal/delivery/http/handlers"
	"ecom-golang-clean-architecture/internal/delivery/http/routes"
	"ecom-golang-clean-architecture/internal/domain/services"
	"ecom-golang-clean-architecture/internal/domain/storage"
	"ecom-golang-clean-architecture/internal/infrastructure/config"
	"ecom-golang-clean-architecture/internal/infrastructure/database"
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
		orderService,
	)

	fileUseCase := usecases.NewFileUseCase(fileService)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userUseCase)
	productHandler := handlers.NewProductHandler(productUseCase) // Use consolidated implementation
	categoryHandler := handlers.NewCategoryHandler(categoryUseCase)
	cartHandler := handlers.NewCartHandler(cartUseCase)
	orderHandler := handlers.NewOrderHandler(orderUseCase)
	fileHandler := handlers.NewFileHandler(fileUseCase)

	// Initialize Gin router
	router := gin.New()

	// Setup routes
	routes.SetupRoutes(
		router,
		cfg,
		userHandler,
		productHandler,
		categoryHandler,
		cartHandler,
		orderHandler,
		fileHandler,
	)

	// Start server
	log.Printf("Starting server on %s", cfg.App.GetAddress())
	if err := router.Run(cfg.App.GetAddress()); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
