package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang-service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database")
	}
	
	// Auto-migrate the schema
	db.AutoMigrate(&models.User{})
	
	return db
}

func TestGetUsers(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	
	// Create test users
	testUsers := []models.User{
		{Email: "user1@example.com", Name: "User 1", AzureID: "azure-user1", IsActive: true},
		{Email: "user2@example.com", Name: "User 2", AzureID: "azure-user2", IsActive: false},
		{Email: "user3@example.com", Name: "User 3", AzureID: "azure-user3", IsActive: true},
	}
	

	
	t.Run("Get all users with default pagination", func(t *testing.T) {
		db := setupTestDB()
		handler := New(db)
		
		// Create test users
		for _, user := range testUsers {
			userCopy := user
			db.Create(&userCopy)
			// If the original was false, explicitly set it to false after creation
			if !user.IsActive {
				userCopy.IsActive = false
				db.Save(&userCopy)
			}
		}
		
		router := gin.New()
		router.GET("/users", handler.GetUsers)
		
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users", nil)
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response PaginatedResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		assert.Equal(t, 1, response.Page)
		assert.Equal(t, 10, response.PerPage)
		assert.Equal(t, int64(3), response.Total)
		assert.Equal(t, 1, response.TotalPages)
		
		// Check that we have users in the response
		usersData := response.Data.([]interface{})
		assert.Equal(t, 3, len(usersData))
	})
	
	t.Run("Get users with custom pagination", func(t *testing.T) {
		db := setupTestDB()
		handler := New(db)
		
		// Create test users
		for _, user := range testUsers {
			userCopy := user
			db.Create(&userCopy)
			// If the original was false, explicitly set it to false after creation
			if !user.IsActive {
				userCopy.IsActive = false
				db.Save(&userCopy)
			}
		}
		
		router := gin.New()
		router.GET("/users", handler.GetUsers)
		
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users?page=1&per_page=2", nil)
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response PaginatedResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		assert.Equal(t, 1, response.Page)
		assert.Equal(t, 2, response.PerPage)
		assert.Equal(t, int64(3), response.Total)
		assert.Equal(t, 2, response.TotalPages)
		
		usersData := response.Data.([]interface{})
		assert.Equal(t, 2, len(usersData))
	})
	
	t.Run("Get users with search filter", func(t *testing.T) {
		db := setupTestDB()
		handler := New(db)
		
		// Create test users
		for _, user := range testUsers {
			userCopy := user
			db.Create(&userCopy)
			// If the original was false, explicitly set it to false after creation
			if !user.IsActive {
				userCopy.IsActive = false
				db.Save(&userCopy)
			}
		}
		
		router := gin.New()
		router.GET("/users", handler.GetUsers)
		
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users?search=user1", nil)
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response PaginatedResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		assert.Equal(t, int64(1), response.Total)
		
		usersData := response.Data.([]interface{})
		assert.Equal(t, 1, len(usersData))
	})
	
	t.Run("Get users with active filter", func(t *testing.T) {
		db := setupTestDB()
		handler := New(db)
		
		// Create test users
		for _, user := range testUsers {
			// Create a copy to avoid the range variable issue
			userCopy := user
			result := db.Create(&userCopy)
			assert.NoError(t, result.Error)
			
			// If the original was false, explicitly set it to false after creation
			// This is needed because GORM's default:true overrides false values
			if !user.IsActive {
				userCopy.IsActive = false
				db.Save(&userCopy)
			}
		}
		
		router := gin.New()
		router.GET("/users", handler.GetUsers)
		
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users?active=true", nil)
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response PaginatedResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		assert.Equal(t, int64(2), response.Total)
		
		usersData := response.Data.([]interface{})
		assert.Equal(t, 2, len(usersData))
	})
}

func TestGetUser(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	db := setupTestDB()
	handler := New(db)
	
	// Create test user
	user := models.User{
		Email:   "test@example.com",
		Name:    "Test User",
		AzureID: "azure-123",
	}
	db.Create(&user)
	
	t.Run("Get existing user", func(t *testing.T) {
		router := gin.New()
		router.GET("/users/:id", handler.GetUser)
		
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/users/%d", user.ID), nil)
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response models.User
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		assert.Equal(t, user.ID, response.ID)
		assert.Equal(t, user.Email, response.Email)
		assert.Equal(t, user.Name, response.Name)
	})
	
	t.Run("Get non-existing user", func(t *testing.T) {
		router := gin.New()
		router.GET("/users/:id", handler.GetUser)
		
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users/999", nil)
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
	
	t.Run("Get user with invalid ID", func(t *testing.T) {
		router := gin.New()
		router.GET("/users/:id", handler.GetUser)
		
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users/invalid", nil)
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCreateUser(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	db := setupTestDB()
	handler := New(db)
	
	t.Run("Create valid user", func(t *testing.T) {
		router := gin.New()
		router.POST("/users", handler.CreateUser)
		
		userReq := models.CreateUserRequest{
			Email:   "newuser@example.com",
			Name:    "New User",
			AzureID: "azure-456",
		}
		
		jsonData, _ := json.Marshal(userReq)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusCreated, w.Code)
		
		var response models.User
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		assert.Equal(t, userReq.Email, response.Email)
		assert.Equal(t, userReq.Name, response.Name)
		assert.Equal(t, userReq.AzureID, response.AzureID)
		assert.True(t, response.IsActive)
	})
	
	t.Run("Create user with invalid email", func(t *testing.T) {
		router := gin.New()
		router.POST("/users", handler.CreateUser)
		
		userReq := models.CreateUserRequest{
			Email: "invalid-email",
			Name:  "Test User",
		}
		
		jsonData, _ := json.Marshal(userReq)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	
	t.Run("Create user with missing name", func(t *testing.T) {
		router := gin.New()
		router.POST("/users", handler.CreateUser)
		
		userReq := models.CreateUserRequest{
			Email: "test@example.com",
		}
		
		jsonData, _ := json.Marshal(userReq)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	
	t.Run("Create user with duplicate email", func(t *testing.T) {
		router := gin.New()
		router.POST("/users", handler.CreateUser)
		
		// Create first user
		user1 := models.User{Email: "duplicate@example.com", Name: "User 1", AzureID: "azure-duplicate1"}
		db.Create(&user1)
		
		// Try to create second user with same email
		userReq := models.CreateUserRequest{
			Email: "duplicate@example.com",
			Name:  "User 2",
		}
		
		jsonData, _ := json.Marshal(userReq)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusConflict, w.Code)
	})
}

func TestUpdateUser(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	db := setupTestDB()
	handler := New(db)
	
	// Create test user
	user := models.User{
		Email:   "update@example.com",
		Name:    "Update User",
		AzureID: "azure-789",
	}
	db.Create(&user)
	
	t.Run("Update user name", func(t *testing.T) {
		router := gin.New()
		router.PUT("/users/:id", handler.UpdateUser)
		
		newName := "Updated Name"
		updateReq := models.UpdateUserRequest{
			Name: &newName,
		}
		
		jsonData, _ := json.Marshal(updateReq)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/users/%d", user.ID), bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response models.User
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		assert.Equal(t, newName, response.Name)
		assert.Equal(t, user.Email, response.Email) // Should remain unchanged
	})
	
	t.Run("Update user with invalid ID", func(t *testing.T) {
		router := gin.New()
		router.PUT("/users/:id", handler.UpdateUser)
		
		updateReq := models.UpdateUserRequest{}
		jsonData, _ := json.Marshal(updateReq)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/users/invalid", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	
	t.Run("Update non-existing user", func(t *testing.T) {
		router := gin.New()
		router.PUT("/users/:id", handler.UpdateUser)
		
		updateReq := models.UpdateUserRequest{}
		jsonData, _ := json.Marshal(updateReq)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/users/999", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestDeleteUser(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	db := setupTestDB()
	handler := New(db)
	
	// Create test user
	user := models.User{
		Email:   "delete@example.com",
		Name:    "Delete User",
		AzureID: "azure-delete",
	}
	db.Create(&user)
	
	t.Run("Delete existing user", func(t *testing.T) {
		router := gin.New()
		router.DELETE("/users/:id", handler.DeleteUser)
		
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/users/%d", user.ID), nil)
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusNoContent, w.Code)
		
		// Verify user is soft deleted
		var deletedUser models.User
		err := db.First(&deletedUser, user.ID).Error
		assert.Error(t, err) // Should not find the user (soft deleted)
	})
	
	t.Run("Delete non-existing user", func(t *testing.T) {
		router := gin.New()
		router.DELETE("/users/:id", handler.DeleteUser)
		
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/users/999", nil)
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
	
	t.Run("Delete user with invalid ID", func(t *testing.T) {
		router := gin.New()
		router.DELETE("/users/:id", handler.DeleteUser)
		
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/users/invalid", nil)
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}