package handlers

import (
	"net/http"
	"sort"
	"strconv"
	"strings"

	"golang-service/internal/config"
	"golang-service/internal/models"
	"golang-service/internal/utils"

	"github.com/gin-gonic/gin"
)

// EnvironmentHandler handles environment-related requests
type EnvironmentHandler struct {
	envService *config.EnvironmentService
}

// NewEnvironmentHandler creates a new environment handler
func NewEnvironmentHandler(envService *config.EnvironmentService) *EnvironmentHandler {
	return &EnvironmentHandler{
		envService: envService,
	}
}

// ListEnvironments handles GET /api/v1/environments
func (h *EnvironmentHandler) ListEnvironments(c *gin.Context) {
	// Check if environment service is available
	if h.envService == nil {
		utils.SendErrorResponse(c, http.StatusServiceUnavailable, "Environment configuration service is not available")
		return
	}

	// Get query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	sortBy := c.DefaultQuery("sortBy", "id")
	sortOrder := c.DefaultQuery("sortOrder", "asc")

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 1000 {
		pageSize = 20
	}

	// Get all environments
	environments, err := h.envService.GetEnvironments()
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to load environments")
		return
	}

	// Apply filters if provided
	filteredEnvs := h.applyFilters(environments, c)

	// Apply sorting
	sort.Slice(filteredEnvs, func(i, j int) bool {
		var less bool
		switch sortBy {
		case "name":
			less = strings.ToLower(filteredEnvs[i].Name) < strings.ToLower(filteredEnvs[j].Name)
		case "description":
			less = strings.ToLower(filteredEnvs[i].Description) < strings.ToLower(filteredEnvs[j].Description)
		default:
			less = strings.ToLower(filteredEnvs[i].ID) < strings.ToLower(filteredEnvs[j].ID)
		}
		if sortOrder == "desc" {
			return !less
		}
		return less
	})

	// Apply pagination
	totalItems := len(filteredEnvs)
	totalPages := (totalItems + pageSize - 1) / pageSize

	start := (page - 1) * pageSize
	end := start + pageSize
	if start >= totalItems {
		start = totalItems
	}
	if end > totalItems {
		end = totalItems
	}

	paginatedEnvs := filteredEnvs[start:end]

	// Build HATEOAS links
	baseURL := getBaseURL(c)
	links := h.buildHATEOASLinks(baseURL, page, pageSize, totalPages, totalItems)

	// Create response
	response := models.EnvironmentListResponse{
		Data: paginatedEnvs,
		Pagination: models.Pagination{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: totalItems,
			TotalPages: totalPages,
		},
		Links: links,
	}

	c.JSON(http.StatusOK, response)
}

// GetEnvironment handles GET /api/v1/environments/:id
func (h *EnvironmentHandler) GetEnvironment(c *gin.Context) {
	// Check if environment service is available
	if h.envService == nil {
		utils.SendErrorResponse(c, http.StatusServiceUnavailable, "Environment configuration service is not available")
		return
	}

	envID := c.Param("id")

	environment, err := h.envService.GetEnvironmentByID(envID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Environment not found")
		return
	}

	// Build HATEOAS links for single environment
	baseURL := getBaseURL(c)
	links := models.HATEOASLinks{
		Self: baseURL + "/api/v1/environments/" + envID,
		VMs:  baseURL + "/api/v1/vms?env=" + envID,
	}

	// Add links to environment response
	response := gin.H{
		"data":  environment,
		"_links": links,
	}

	c.JSON(http.StatusOK, response)
}

// applyFilters applies query parameter filters to environments
func (h *EnvironmentHandler) applyFilters(environments []models.Environment, c *gin.Context) []models.Environment {
	filtered := environments

	// Filter by tag
	if tags := c.QueryArray("tag"); len(tags) > 0 {
		var tagFiltered []models.Environment
		for _, env := range filtered {
			for _, tag := range tags {
				for _, envTag := range env.Tags {
					if envTag == tag {
						tagFiltered = append(tagFiltered, env)
						break
					}
				}
			}
		}
		filtered = tagFiltered
	}

	// Filter by account
	if account := c.Query("account"); account != "" {
		var accountFiltered []models.Environment
		for _, env := range filtered {
			if env.Criteria.Account == account {
				accountFiltered = append(accountFiltered, env)
			}
		}
		filtered = accountFiltered
	}

	// Filter by region
	if region := c.Query("region"); region != "" {
		var regionFiltered []models.Environment
		for _, env := range filtered {
			if env.Criteria.Region == region {
				regionFiltered = append(regionFiltered, env)
			}
		}
		filtered = regionFiltered
	}

	// Filter by VPC
	if vpc := c.Query("vpc"); vpc != "" {
		var vpcFiltered []models.Environment
		for _, env := range filtered {
			if env.Criteria.VPC == vpc {
				vpcFiltered = append(vpcFiltered, env)
			}
		}
		filtered = vpcFiltered
	}

	// Filter by name (contains)
	if name := c.Query("name"); name != "" {
		var nameFiltered []models.Environment
		for _, env := range filtered {
			if strings.Contains(strings.ToLower(env.Name), strings.ToLower(name)) {
				nameFiltered = append(nameFiltered, env)
			}
		}
		filtered = nameFiltered
	}

	// Filter by description (contains)
	if description := c.Query("description"); description != "" {
		var descFiltered []models.Environment
		for _, env := range filtered {
			if strings.Contains(strings.ToLower(env.Description), strings.ToLower(description)) {
				descFiltered = append(descFiltered, env)
			}
		}
		filtered = descFiltered
	}

	return filtered
}

// buildHATEOASLinks builds HATEOAS links for navigation
func (h *EnvironmentHandler) buildHATEOASLinks(baseURL string, page, pageSize, totalPages, totalItems int) models.HATEOASLinks {
	links := models.HATEOASLinks{
		Self: baseURL + "/api/v1/environments",
		VMs:  baseURL + "/api/v1/vms",
	}

	// Add pagination links
	if page > 1 {
		links.Previous = baseURL + "/api/v1/environments?page=" + strconv.Itoa(page-1) + "&pageSize=" + strconv.Itoa(pageSize)
	}

	if page < totalPages {
		links.Next = baseURL + "/api/v1/environments?page=" + strconv.Itoa(page+1) + "&pageSize=" + strconv.Itoa(pageSize)
	}

	return links
}

// ReloadConfig handles POST /api/v1/environments/reload
func (h *EnvironmentHandler) ReloadConfig(c *gin.Context) {
	// Check if environment service is available
	if h.envService == nil {
		utils.SendErrorResponse(c, http.StatusServiceUnavailable, "Environment configuration service is not available")
		return
	}

	if err := h.envService.ReloadConfig(); err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to reload configuration")
		return
	}

	// Validate the reloaded configuration
	if err := h.envService.ValidateConfig(); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid configuration")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Configuration reloaded successfully",
		"timestamp": h.envService.GetLastLoadTime(),
	})
}

// GetConfigInfo handles GET /api/v1/environments/config/info
func (h *EnvironmentHandler) GetConfigInfo(c *gin.Context) {
	// Check if environment service is available
	if h.envService == nil {
		utils.SendErrorResponse(c, http.StatusServiceUnavailable, "Environment configuration service is not available")
		return
	}

	environments, err := h.envService.GetEnvironments()
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get configuration info")
		return
	}

	// Count environments by various criteria
	accountCount := make(map[string]int)
	regionCount := make(map[string]int)
	tagCount := make(map[string]int)

	for _, env := range environments {
		accountCount[env.Criteria.Account]++
		regionCount[env.Criteria.Region]++
		for _, tag := range env.Tags {
			tagCount[tag]++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"totalEnvironments": len(environments),
		"configPath":        h.envService.GetConfigPath(),
		"lastLoaded":        h.envService.GetLastLoadTime(),
		"accounts":          accountCount,
		"regions":           regionCount,
		"tags":              tagCount,
	})
}

// getBaseURL builds the base URL from the request
func getBaseURL(c *gin.Context) string {
	scheme := c.Request.URL.Scheme
	if scheme == "" {
		scheme = "http"
	}
	host := c.Request.Host
	return scheme + "://" + host
} 