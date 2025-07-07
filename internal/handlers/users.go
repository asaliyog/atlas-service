package handlers

import (
	"golang-service/internal/models"
	"golang-service/internal/utils"
	"fmt"
	"net/http"
	"reflect"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UsersHandler handles user-related HTTP requests
type UsersHandler struct {
	db *gorm.DB
}

// NewUsersHandler creates a new users handler
func NewUsersHandler(db *gorm.DB) *UsersHandler {
	return &UsersHandler{db: db}
}

// GetUsers handles GET /api/v1/users
func (h *UsersHandler) GetUsers(c *gin.Context) {
	// Parse query parameters using the reusable utility
	params, err := utils.ParseQueryParams(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Validate query parameters
	if err := utils.ValidateQueryParams(params); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Get users from database
	var users []models.User
	if err := h.db.Find(&users).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	// Apply filters using the reusable utility
	filteredUsers := utils.ApplyInMemoryFilters(users, params.Filters, h.extractUserField)

	// Apply sorting
	sortedUsers := h.applySorting(filteredUsers, params.SortBy, params.SortOrder)

	// Calculate pagination
	totalItems := len(sortedUsers)
	paginatedUsers := utils.ApplyPagination(sortedUsers, params.Page, params.PageSize)

	// Send response using the reusable utility
	utils.SendPaginatedResponse(c, paginatedUsers, params.Page, params.PageSize, totalItems)
}

// extractUserField extracts field value from User struct for filtering
func (h *UsersHandler) extractUserField(user models.User, field string) interface{} {
	v := reflect.ValueOf(user)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldType := t.Field(i)
		fieldValue := v.Field(i)

		// Check JSON tag first, then field name
		jsonTag := fieldType.Tag.Get("json")
		if jsonTag != "" {
			jsonField := strings.Split(jsonTag, ",")[0]
			if jsonField == field {
				return fieldValue.Interface()
			}
		}

		// Check field name
		if strings.EqualFold(fieldType.Name, field) {
			return fieldValue.Interface()
		}
	}

	return nil
}

// applySorting applies sorting to users
func (h *UsersHandler) applySorting(users []models.User, sortBy, sortOrder string) []models.User {
	if sortBy == "" {
		return users
	}

	sortedUsers := make([]models.User, len(users))
	copy(sortedUsers, users)

	sort.Slice(sortedUsers, func(i, j int) bool {
		valI := h.extractUserField(sortedUsers[i], sortBy)
		valJ := h.extractUserField(sortedUsers[j], sortBy)

		// Handle nil values
		if valI == nil && valJ == nil {
			return false
		}
		if valI == nil {
			return sortOrder == "asc"
		}
		if valJ == nil {
			return sortOrder == "desc"
		}

		// Convert to string for comparison
		strI := strings.ToLower(fmt.Sprintf("%v", valI))
		strJ := strings.ToLower(fmt.Sprintf("%v", valJ))

		if sortOrder == "desc" {
			return strI > strJ
		}
		return strI < strJ
	})

	return sortedUsers
}