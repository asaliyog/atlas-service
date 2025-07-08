package main

import (
	"log"
	"os"

	"golang-service/internal/config"
	"golang-service/internal/database"
	"golang-service/internal/handlers"
	"golang-service/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
)

// @title Cloud Inventory API
// @version 1.0.0
// @description API for retrieving unified virtual machine (VM) data across AWS EC2, GCP Compute, and Azure VMs, with detailed cloud-specific fields and flexible filtering
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Setup Gin router
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Add CORS middleware
	router.Use(middleware.CORS())

	// Health check endpoint (no auth required, but with DB context)
	router.GET("/health", handlers.DatabaseHealthMiddleware(db), handlers.HealthCheck)

	// API routes with authentication
	api := router.Group("/api/v1")
	api.Use(middleware.AzureEntraAuth(cfg))
	{
		// Initialize handlers
		usersHandler := handlers.NewUsersHandler(db)
		vmsHandler := handlers.NewVMsHandler(db)

		// User management endpoints
		api.GET("/users", usersHandler.GetUsers)

		// VM management endpoints
		api.GET("/vms", vmsHandler.GetVMs)
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Health check available at: http://localhost:%s/health", port)
	log.Printf("Swagger documentation available at: http://localhost:%s/swagger/index.html", port)
	
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}