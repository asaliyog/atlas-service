package models

import (
	"time"
)

// Environment represents an environment configuration
type Environment struct {
	ID          string                 `json:"id" yaml:"id"`
	Name        string                 `json:"name" yaml:"name"`
	Description string                 `json:"description" yaml:"description"`
	Criteria    EnvironmentCriteria    `json:"criteria" yaml:"criteria"`
	Tags        []string               `json:"tags" yaml:"tags"`
	Metadata    map[string]interface{} `json:"metadata" yaml:"metadata"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
}

// EnvironmentCriteria defines the criteria for environment matching
type EnvironmentCriteria struct {
	CloudType   string `json:"cloud_type" yaml:"cloud_type"`
	Account     string `json:"account" yaml:"account"`
	Region      string `json:"region" yaml:"region"`
	VPC         string `json:"vpc" yaml:"vpc"`
	Subscription string `json:"subscription" yaml:"subscription"`
	Location    string `json:"location" yaml:"location"`
	Project     string `json:"project" yaml:"project"`
	Zone        string `json:"zone" yaml:"zone"`
}

// EnvironmentListResponse represents the response for listing environments
type EnvironmentListResponse struct {
	Data       []Environment `json:"data"`
	Pagination Pagination    `json:"pagination"`
	Links      HATEOASLinks  `json:"_links"`
}

// HATEOASLinks represents HATEOAS links for navigation
type HATEOASLinks struct {
	Self     string `json:"self"`
	VMs      string `json:"vms,omitempty"`
	Next     string `json:"next,omitempty"`
	Previous string `json:"previous,omitempty"`
}

// EnvironmentConfig represents the root configuration structure
type EnvironmentConfig struct {
	Environments []Environment `yaml:"environments"`
} 