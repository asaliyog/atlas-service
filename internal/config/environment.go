package config

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"

	"golang-service/internal/models"
)

// EnvironmentService manages environment configurations
type EnvironmentService struct {
	configPath string
	config     *models.EnvironmentConfig
	mu         sync.RWMutex
	lastLoad   time.Time
}

// NewEnvironmentService creates a new environment service
func NewEnvironmentService(configPath string) *EnvironmentService {
	return &EnvironmentService{
		configPath: configPath,
	}
}

// LoadConfig loads the environment configuration from the YAML file
func (s *EnvironmentService) LoadConfig() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Read the configuration file
	data, err := os.ReadFile(s.configPath)
	if err != nil {
		return fmt.Errorf("failed to read environment config file: %w", err)
	}

	// Parse the YAML configuration
	var config models.EnvironmentConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse environment config YAML: %w", err)
	}

	// Add timestamps to environments
	now := time.Now()
	for i := range config.Environments {
		config.Environments[i].CreatedAt = now
		config.Environments[i].UpdatedAt = now
	}

	s.config = &config
	s.lastLoad = now

	return nil
}

// GetEnvironments returns all environments
func (s *EnvironmentService) GetEnvironments() ([]models.Environment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.config == nil {
		return nil, fmt.Errorf("environment configuration not loaded")
	}

	// Return a copy to avoid external modifications
	environments := make([]models.Environment, len(s.config.Environments))
	copy(environments, s.config.Environments)

	return environments, nil
}

// GetEnvironmentByID returns a specific environment by ID
func (s *EnvironmentService) GetEnvironmentByID(id string) (*models.Environment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.config == nil {
		return nil, fmt.Errorf("environment configuration not loaded")
	}

	for _, env := range s.config.Environments {
		if env.ID == id {
			return &env, nil
		}
	}

	return nil, fmt.Errorf("environment with ID '%s' not found", id)
}

// GetEnvironmentsByCriteria returns environments matching specific criteria
func (s *EnvironmentService) GetEnvironmentsByCriteria(account, region, vpc string) ([]models.Environment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.config == nil {
		return nil, fmt.Errorf("environment configuration not loaded")
	}

	var matching []models.Environment

	for _, env := range s.config.Environments {
		matches := true

		if account != "" && env.Criteria.Account != account {
			matches = false
		}
		if region != "" && env.Criteria.Region != region {
			matches = false
		}
		if vpc != "" && env.Criteria.VPC != vpc {
			matches = false
		}

		if matches {
			matching = append(matching, env)
		}
	}

	return matching, nil
}

// GetEnvironmentsByTag returns environments with specific tags
func (s *EnvironmentService) GetEnvironmentsByTag(tags []string) ([]models.Environment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.config == nil {
		return nil, fmt.Errorf("environment configuration not loaded")
	}

	var matching []models.Environment

	for _, env := range s.config.Environments {
		// Check if environment has any of the specified tags
		for _, tag := range tags {
			for _, envTag := range env.Tags {
				if envTag == tag {
					matching = append(matching, env)
					break
				}
			}
		}
	}

	return matching, nil
}

// ResolveEnvironmentForResource resolves which environment a resource belongs to
func (s *EnvironmentService) ResolveEnvironmentForResource(accountID, region, vpc string) (*models.Environment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.config == nil {
		return nil, fmt.Errorf("environment configuration not loaded")
	}

	// Find the best matching environment
	var bestMatch *models.Environment
	var bestScore int

	for _, env := range s.config.Environments {
		score := 0
		
		// Check account match (highest priority)
		if env.Criteria.Account == accountID {
			score += 10
		} else if env.Criteria.Account != "" {
			continue // Skip if account doesn't match and is required
		}
		
		// Check region match
		if env.Criteria.Region == region {
			score += 5
		} else if env.Criteria.Region != "" {
			// For GCP, check if zone starts with the region (e.g., us-central1-a matches us-central1)
			if strings.HasPrefix(region, env.Criteria.Region) {
				score += 5
			} else {
				continue // Skip if region doesn't match and is required
			}
		}
		
		// Check VPC match (if specified)
		if env.Criteria.VPC != "" {
			if env.Criteria.VPC == vpc {
				score += 3
			} else if vpc != "" {
				// Only skip if VPC is specified in environment AND VM has a VPC that doesn't match
				// This allows Azure/GCP environments (with empty VPC) to match VMs without VPC info
				continue
			}
			// If environment has VPC but VM doesn't, still allow match (for Azure/GCP)
		}
		
		// Update best match if this environment has a higher score
		if score > bestScore {
			bestScore = score
			bestMatch = &env
		}
	}

	if bestMatch == nil {
		return nil, fmt.Errorf("no environment found matching criteria: account=%s, region=%s, vpc=%s", accountID, region, vpc)
	}

	return bestMatch, nil
}

// ResolveEnvironmentForVM resolves environment for a VM based on its cloud-specific properties
func (s *EnvironmentService) ResolveEnvironmentForVM(vm models.VM) (*models.Environment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.config == nil {
		return nil, fmt.Errorf("environment configuration not loaded")
	}

	// Find the best matching environment based on cloud type
	var bestMatch *models.Environment
	var bestScore int

	for _, env := range s.config.Environments {
		// First check if cloud type matches
		if env.Criteria.CloudType != "" && env.Criteria.CloudType != vm.CloudType {
			continue // Skip if cloud type doesn't match
		}

		score := 0
		
		switch vm.CloudType {
		case "aws":
			// AWS criteria: account, region, vpc
			if env.Criteria.Account == vm.CloudAccountID {
				score += 10
			} else if env.Criteria.Account != "" {
				continue // Skip if account doesn't match
			}
			
			if env.Criteria.Region == vm.Location {
				score += 5
			} else if env.Criteria.Region != "" {
				continue // Skip if region doesn't match
			}
			
			// VPC matching would need to be extracted from VM data
			// For now, we'll skip VPC matching as it requires additional data extraction
			
		case "azure":
			// Azure criteria: subscription, location
			if env.Criteria.Subscription == vm.CloudAccountID {
				score += 10
			} else if env.Criteria.Subscription != "" {
				continue // Skip if subscription doesn't match
			}
			
			if env.Criteria.Location == vm.Location {
				score += 5
			} else if env.Criteria.Location != "" {
				continue // Skip if location doesn't match
			}
			
		case "gcp":
			// GCP criteria: project, zone
			if env.Criteria.Project == vm.CloudAccountID {
				score += 10
			} else if env.Criteria.Project != "" {
				continue // Skip if project doesn't match
			}
			
			if env.Criteria.Zone != "" {
				// Check if zone starts with the environment zone (e.g., us-central1-a matches us-central1)
				if strings.HasPrefix(vm.Location, env.Criteria.Zone) {
					score += 5
				} else {
					continue // Skip if zone doesn't match
				}
			}
		}
		
		// Update best match if this environment has a higher score
		if score > bestScore {
			bestScore = score
			bestMatch = &env
		}
	}

	if bestMatch == nil {
		return nil, fmt.Errorf("no environment found matching criteria for VM: cloudType=%s, account=%s, location=%s", vm.CloudType, vm.CloudAccountID, vm.Location)
	}

	return bestMatch, nil
}



// ReloadConfig reloads the configuration from the file
func (s *EnvironmentService) ReloadConfig() error {
	return s.LoadConfig()
}

// GetLastLoadTime returns when the configuration was last loaded
func (s *EnvironmentService) GetLastLoadTime() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastLoad
}

// GetConfigPath returns the configuration file path
func (s *EnvironmentService) GetConfigPath() string {
	return s.configPath
}

// ValidateConfig validates the loaded configuration
func (s *EnvironmentService) ValidateConfig() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.config == nil {
		return fmt.Errorf("environment configuration not loaded")
	}

	// Check for duplicate IDs
	ids := make(map[string]bool)
	for _, env := range s.config.Environments {
		if ids[env.ID] {
			return fmt.Errorf("duplicate environment ID: %s", env.ID)
		}
		ids[env.ID] = true

		// Validate required fields
		if env.ID == "" {
			return fmt.Errorf("environment missing required field: id")
		}
		if env.Name == "" {
			return fmt.Errorf("environment '%s' missing required field: name", env.ID)
		}
		if env.Criteria.Account == "" {
			return fmt.Errorf("environment '%s' missing required field: criteria.account", env.ID)
		}
		if env.Criteria.Region == "" {
			return fmt.Errorf("environment '%s' missing required field: criteria.region", env.ID)
		}
	}

	return nil
} 