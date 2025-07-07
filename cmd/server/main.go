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

// @title Golang Service API
// @version 1.0
// @description RESTful API service with Azure Entra ID authentication
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
		h := handlers.New(db)

		// User management endpoints
		api.GET("/users", h.GetUsers)
		api.POST("/users", h.CreateUser)
		api.GET("/users/:id", h.GetUser)
		api.PUT("/users/:id", h.UpdateUser)
		api.DELETE("/users/:id", h.DeleteUser)
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