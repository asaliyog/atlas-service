package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestHealthCheck(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Test health check without database
	t.Run("HealthCheck without database", func(t *testing.T) {
		router.GET("/health", HealthCheck)
		
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/health", nil)
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response HealthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		assert.Equal(t, "healthy", response.Status)
		assert.Equal(t, "1.0.0", response.Version)
		assert.Equal(t, "unknown", response.Database.Status)
		assert.Equal(t, "N/A", response.Database.Ping)
		assert.Equal(t, "healthy", response.Services["api"])
	})
	
	// Test health check with database
	t.Run("HealthCheck with database", func(t *testing.T) {
		// Setup in-memory SQLite database
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)
		
		router := gin.New()
		router.GET("/health", DatabaseHealthMiddleware(db), HealthCheck)
		
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/health", nil)
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response HealthResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		assert.Equal(t, "healthy", response.Status)
		assert.Equal(t, "1.0.0", response.Version)
		assert.Equal(t, "healthy", response.Database.Status)
		assert.NotEqual(t, "N/A", response.Database.Ping)
		assert.Equal(t, "healthy", response.Services["api"])
	})
}

func TestDatabaseHealthMiddleware(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Setup in-memory SQLite database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	
	// Test middleware
	router.Use(DatabaseHealthMiddleware(db))
	router.GET("/test", func(c *gin.Context) {
		dbInstance, exists := c.Get("db")
		assert.True(t, exists)
		assert.NotNil(t, dbInstance)
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
}