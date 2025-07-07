package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PaginatedResponse represents a generic paginated response
type PaginatedResponse[T any] struct {
	Data       []T `json:"data"`
	Pagination struct {
		Page       int `json:"page"`
		PageSize   int `json:"pageSize"`
		TotalItems int `json:"totalItems"`
		TotalPages int `json:"totalPages"`
	} `json:"pagination"`
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse[T any](data []T, page, pageSize, totalItems int) PaginatedResponse[T] {
	totalPages := totalItems / pageSize
	if totalItems%pageSize != 0 {
		totalPages++
	}

	response := PaginatedResponse[T]{
		Data: data,
	}
	response.Pagination.Page = page
	response.Pagination.PageSize = pageSize
	response.Pagination.TotalItems = totalItems
	response.Pagination.TotalPages = totalPages

	return response
}

// SendPaginatedResponse sends a paginated response with proper HTTP status
func SendPaginatedResponse[T any](c *gin.Context, data []T, page, pageSize, totalItems int) {
	response := NewPaginatedResponse(data, page, pageSize, totalItems)
	c.JSON(http.StatusOK, response)
}

// SendErrorResponse sends an error response with proper HTTP status
func SendErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"error": message,
	})
}

// SendSuccessResponse sends a success response with data
func SendSuccessResponse[T any](c *gin.Context, data T) {
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

// SendListResponse sends a simple list response without pagination
func SendListResponse[T any](c *gin.Context, data []T) {
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
} 