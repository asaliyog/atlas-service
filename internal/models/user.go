package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Name      string         `json:"name" gorm:"not null"`
	AzureID   string         `json:"azure_id" gorm:"uniqueIndex"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// CreateUserRequest represents the request payload for creating a user
type CreateUserRequest struct {
	Email   string `json:"email" binding:"required,email"`
	Name    string `json:"name" binding:"required"`
	AzureID string `json:"azure_id"`
}

// UpdateUserRequest represents the request payload for updating a user
type UpdateUserRequest struct {
	Email    *string `json:"email" binding:"omitempty,email"`
	Name     *string `json:"name"`
	IsActive *bool   `json:"is_active"`
}