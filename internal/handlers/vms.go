package handlers

import (
	"context"
	"encoding/json"
	"golang-service/internal/cache"
	"golang-service/internal/config"
	"golang-service/internal/models"
	"golang-service/internal/utils"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// VMsHandler handles VM-related HTTP requests
type VMsHandler struct {
	db         *gorm.DB
	cache      *cache.RedisCache
	envService *config.EnvironmentService
	config     *config.Config
}

// NewVMsHandler creates a new VMs handler
func NewVMsHandler(db *gorm.DB, cache *cache.RedisCache, envService *config.EnvironmentService, config *config.Config) *VMsHandler {
	return &VMsHandler{db: db, cache: cache, envService: envService, config: config}
}

// GetVMs handles GET /api/v1/vms
func (h *VMsHandler) GetVMs(c *gin.Context) {
	// Get filter configuration for VMs endpoint
	filterConfig := config.VMsFilterConfig()

	// Parse and validate filters from query parameters
	filters, err := filterConfig.ParseQueryParams(c.Request.URL.Query())
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Filter validation error: %s", err.Error()))
		return
	}

	// Parse pagination and sorting parameters
	page := 1
	pageSize := 10
	sortBy := ""
	sortOrder := "asc"

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := c.Query("pageSize"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 1000 {
			pageSize = ps
		}
	}

	if sortByParam := c.Query("sortBy"); sortByParam != "" {
		sortBy = sortByParam
	}

	if sortOrderParam := c.Query("sortOrder"); sortOrderParam != "" {
		if sortOrderParam == "desc" {
			sortOrder = "desc"
		}
	}

	// Try to get VMs from cache first (if Redis is available)
	var cachedVMs []models.VM
	if h.cache != nil {
		ctx := context.Background()
		cachedVMs, err = h.cache.GetVMs(ctx)
		if err != nil {
			log.Printf("Cache error: %v", err)
		}
	}

	// If cache miss or Redis unavailable, fetch from database and cache the result
	if cachedVMs == nil {
		log.Println("Cache miss or Redis unavailable - fetching VMs from database")
		cachedVMs, err = h.fetchVMsFromDatabase()
		if err != nil {
			utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to fetch VMs")
			return
		}

		// Cache the result (async) if Redis is available
		if h.cache != nil {
			go func() {
				if err := h.cache.SetVMs(context.Background(), cachedVMs); err != nil {
					log.Printf("Failed to cache VMs: %v", err)
				}
			}()
		}
	} else {
		log.Println("Cache hit - using cached VMs")
	}

	// Apply filters using the configurable system (including environment filters)
	filteredVMs := utils.ApplyFilters(cachedVMs, filters)

	// Apply sorting
	sortedVMs := h.applySorting(filteredVMs, sortBy, sortOrder)

	// Calculate pagination
	totalItems := len(sortedVMs)
	paginatedVMs := utils.ApplyPagination(sortedVMs, page, pageSize)

	// Send response using the reusable utility
	utils.SendPaginatedResponse(c, paginatedVMs, page, pageSize, totalItems)
}



// applySorting applies sorting to VMs
func (h *VMsHandler) applySorting(vms []models.VM, sortBy, sortOrder string) []models.VM {
	if sortBy == "" {
		return vms
	}

	sortedVMs := make([]models.VM, len(vms))
	copy(sortedVMs, vms)

	sort.Slice(sortedVMs, func(i, j int) bool {
		valI := utils.GetFieldValue(sortedVMs[i], sortBy)
		valJ := utils.GetFieldValue(sortedVMs[j], sortBy)

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

	return sortedVMs
}

// fetchVMsFromDatabase fetches VMs from all cloud provider tables in parallel
func (h *VMsHandler) fetchVMsFromDatabase() ([]models.VM, error) {
	var allVMs []models.VM
	var mu sync.Mutex
	var wg sync.WaitGroup
	var errors []error

	// Fetch from AWS VMs
	wg.Add(1)
	go func() {
		defer wg.Done()
		var awsVMs []models.AWSEC2Instance
		if err := h.db.Find(&awsVMs).Error; err != nil {
			mu.Lock()
			errors = append(errors, fmt.Errorf("failed to fetch AWS VMs: %w", err))
			mu.Unlock()
			return
		}

		// Convert AWS VMs to unified VM format
		for _, awsVM := range awsVMs {
			// Extract name from tags or use instance_id as fallback
			name := awsVM.InstanceID
			if awsVM.Tags != nil {
				// Try to extract Name tag from JSON
				var tags map[string]interface{}
				if err := json.Unmarshal(awsVM.Tags, &tags); err == nil {
					if nameTag, ok := tags["Name"].(string); ok && nameTag != "" {
						name = nameTag
					}
				}
			}

			// Extract status from state JSON
			status := "unknown"
			if awsVM.State != nil {
				var state map[string]interface{}
				if err := json.Unmarshal(awsVM.State, &state); err == nil {
					if stateName, ok := state["name"].(string); ok {
						status = stateName
					}
				}
			}

			vm := models.VM{
				ID:                   awsVM.ARN,
				Name:                 name,
				CloudType:            "aws",
				Status:               status,
				CloudAccountID:       awsVM.AccountID,
				Location:             awsVM.Region,
				InstanceType:         awsVM.InstanceType,
				CloudSpecificDetails: awsVM.Tags, // Store tags as cloud-specific details
			}
			
			// Resolve environment for this VM if environment service is available
			if h.envService != nil {
				if environment, err := h.envService.ResolveEnvironmentForVM(vm); err == nil {
					vm.Environment = &models.EnvironmentInfo{
						ID:          environment.ID,
						Name:        environment.Name,
						Description: environment.Description,
						Tags:        environment.Tags,
					}
					vm.Env = environment.ID
				}
			}
			
			mu.Lock()
			allVMs = append(allVMs, vm)
			mu.Unlock()
		}
	}()

	// Fetch from Azure VMs
	wg.Add(1)
	go func() {
		defer wg.Done()
		var azureVMs []models.AzureVMInstance
		if err := h.db.Find(&azureVMs).Error; err != nil {
			mu.Lock()
			errors = append(errors, fmt.Errorf("failed to fetch Azure VMs: %w", err))
			mu.Unlock()
			return
		}

		// Convert Azure VMs to unified VM format
		for _, azureVM := range azureVMs {
			// Extract status from properties JSON
			status := "unknown"
			if azureVM.Properties != nil {
				var properties map[string]interface{}
				if err := json.Unmarshal(azureVM.Properties, &properties); err == nil {
					if provisioningState, ok := properties["provisioningState"].(string); ok {
						status = provisioningState
					}
				}
			}

			vm := models.VM{
				ID:                   azureVM.ID,
				Name:                 azureVM.Name,
				CloudType:            "azure",
				Status:               status,
				CloudAccountID:       azureVM.SubscriptionID,
				Location:             azureVM.Location,
				InstanceType:         "", // Will extract from properties if needed
				CloudSpecificDetails: azureVM.Properties, // Store properties as cloud-specific details
			}
			
			// Resolve environment for this VM if environment service is available
			if h.envService != nil {
				if environment, err := h.envService.ResolveEnvironmentForVM(vm); err == nil {
					vm.Environment = &models.EnvironmentInfo{
						ID:          environment.ID,
						Name:        environment.Name,
						Description: environment.Description,
						Tags:        environment.Tags,
					}
					vm.Env = environment.ID
				}
			}
			
			mu.Lock()
			allVMs = append(allVMs, vm)
			mu.Unlock()
		}
	}()

	// Fetch from GCP VMs
	wg.Add(1)
	go func() {
		defer wg.Done()
		var gcpVMs []models.GCPComputeInstance
		if err := h.db.Find(&gcpVMs).Error; err != nil {
			mu.Lock()
			errors = append(errors, fmt.Errorf("failed to fetch GCP VMs: %w", err))
			mu.Unlock()
			return
		}

		// Convert GCP VMs to unified VM format
		for _, gcpVM := range gcpVMs {
			vm := models.VM{
				ID:                   gcpVM.SelfLink,
				Name:                 gcpVM.Name,
				CloudType:            "gcp",
				Status:               gcpVM.Status,
				CloudAccountID:       gcpVM.ProjectID,
				Location:             gcpVM.Zone, // Using zone as location
				InstanceType:         gcpVM.MachineType,
				CloudSpecificDetails: gcpVM.Labels, // Store labels as cloud-specific details
			}
			
			// Resolve environment for this VM if environment service is available
			if h.envService != nil {
				if environment, err := h.envService.ResolveEnvironmentForVM(vm); err == nil {
					vm.Environment = &models.EnvironmentInfo{
						ID:          environment.ID,
						Name:        environment.Name,
						Description: environment.Description,
						Tags:        environment.Tags,
					}
					vm.Env = environment.ID
				}
			}
			
			mu.Lock()
			allVMs = append(allVMs, vm)
			mu.Unlock()
		}
	}()

	// Wait for all goroutines to complete
	wg.Wait()

	// Check for errors
	if len(errors) > 0 {
		return nil, fmt.Errorf("errors fetching VMs: %v", errors)
	}

	log.Printf("Fetched %d VMs from database", len(allVMs))

	return allVMs, nil
}

