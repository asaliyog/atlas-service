package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version"`
	Uptime    string            `json:"uptime"`
	Database  DatabaseStatus    `json:"database"`
	Services  map[string]string `json:"services"`
}

// DatabaseStatus represents the database connectivity status
type DatabaseStatus struct {
	Status      string `json:"status"`
	Ping        string `json:"ping"`
	Connections int    `json:"connections"`
}

var startTime = time.Now()

// HealthCheck godoc
// @Summary Health check endpoint
// @Description Get the health status of the service
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Failure 503 {object} HealthResponse
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	uptime := time.Since(startTime)
	
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Uptime:    uptime.String(),
		Database: DatabaseStatus{
			Status: "unknown",
			Ping:   "N/A",
		},
		Services: map[string]string{
			"api": "healthy",
		},
	}

	// Check database connectivity if available
	if db, exists := c.Get("db"); exists {
		if gormDB, ok := db.(*gorm.DB); ok {
			start := time.Now()
			sqlDB, err := gormDB.DB()
			if err != nil {
				response.Database.Status = "error"
				response.Database.Ping = err.Error()
				response.Status = "unhealthy"
				c.JSON(http.StatusServiceUnavailable, response)
				return
			}
			
			if err := sqlDB.Ping(); err != nil {
				response.Database.Status = "error"
				response.Database.Ping = err.Error()
				response.Status = "unhealthy"
				c.JSON(http.StatusServiceUnavailable, response)
				return
			}
			
			response.Database.Status = "healthy"
			response.Database.Ping = time.Since(start).String()
			response.Database.Connections = sqlDB.Stats().OpenConnections
		}
	}

	c.JSON(http.StatusOK, response)
}

// DatabaseHealthMiddleware adds database instance to context for health checks
func DatabaseHealthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}