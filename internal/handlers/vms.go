package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"golang-service/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetVMs godoc
// @Summary Retrieve a list of virtual machines
// @Description Fetches a paginated list of virtual machines across AWS EC2, GCP Compute, and Azure VMs with filtering and sorting
// @Tags vms
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number for pagination (1-based)" default(1)
// @Param pageSize query int false "Number of VMs per page" default(20)
// @Param sortBy query string false "Field to sort by" default("id")
// @Param sortOrder query string false "Sort order (asc or desc)" default("asc")
// @Param filter query string false "JSON string of filter objects"
// @Success 200 {object} models.VMListResponse
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /vms [get]
func (h *Handler) GetVMs(c *gin.Context) {
	// Parse query parameters
	params, err := h.parseVMQueryParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Message: err.Error()})
		return
	}

	// Build base query
	query := h.db.Model(&models.VM{})

	// Apply filters
	for _, filter := range params.Filters {
		query, err = h.applyVMFilter(query, filter)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Message: fmt.Sprintf("Invalid filter: %s", err.Error())})
			return
		}
	}

	// Get total count
	var totalItems int64
	if err := query.Count(&totalItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{Message: "Failed to count VMs"})
		return
	}

	// Apply sorting
	sortOrder := "ASC"
	if params.SortOrder == "desc" {
		sortOrder = "DESC"
	}
	
	sortField := params.SortBy
	// Handle nested fields for sorting
	if strings.Contains(sortField, ".") {
		// For JSON fields, we need to use JSON extraction
		if strings.HasPrefix(sortField, "cloudSpecificDetails.") {
			jsonPath := strings.TrimPrefix(sortField, "cloudSpecificDetails.")
			sortField = fmt.Sprintf("JSON_EXTRACT(cloud_specific_details, '$.%s')", jsonPath)
		}
	} else {
		// Map field names to database column names
		sortField = h.mapFieldToDBColumn(sortField)
	}
	
	query = query.Order(fmt.Sprintf("%s %s", sortField, sortOrder))

	// Apply pagination
	offset := (params.Page - 1) * params.PageSize
	var vms []models.VM
	if err := query.Offset(offset).Limit(params.PageSize).Find(&vms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{Message: "Failed to fetch VMs"})
		return
	}

	// Calculate total pages
	totalPages := int(totalItems) / params.PageSize
	if int(totalItems)%params.PageSize != 0 {
		totalPages++
	}

	response := models.VMListResponse{
		Data: vms,
		Pagination: models.Pagination{
			Page:       params.Page,
			PageSize:   params.PageSize,
			TotalItems: int(totalItems),
			TotalPages: totalPages,
		},
	}

	c.JSON(http.StatusOK, response)
}

// parseVMQueryParams parses and validates query parameters
func (h *Handler) parseVMQueryParams(c *gin.Context) (models.VMQueryParams, error) {
	params := models.VMQueryParams{
		Page:      1,
		PageSize:  20,
		SortBy:    "id",
		SortOrder: "asc",
		Filters:   []models.VMFilter{},
	}

	// Parse page
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			params.Page = page
		}
	}

	// Parse pageSize
	if pageSizeStr := c.Query("pageSize"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 {
			if pageSize > 100 {
				pageSize = 100
			}
			params.PageSize = pageSize
		}
	}

	// Parse sortBy
	if sortBy := c.Query("sortBy"); sortBy != "" {
		params.SortBy = sortBy
	}

	// Parse sortOrder
	if sortOrder := c.Query("sortOrder"); sortOrder != "" {
		if sortOrder == "asc" || sortOrder == "desc" {
			params.SortOrder = sortOrder
		}
	}

	// Parse filters
	if filterStr := c.Query("filter"); filterStr != "" {
		var filters []models.VMFilter
		if err := json.Unmarshal([]byte(filterStr), &filters); err != nil {
			return params, fmt.Errorf("invalid filter JSON: %w", err)
		}
		params.Filters = filters
	}

	return params, nil
}

// applyVMFilter applies a single filter to the query
func (h *Handler) applyVMFilter(query *gorm.DB, filter models.VMFilter) (*gorm.DB, error) {
	// Handle nested fields (cloudSpecificDetails.*)
	if strings.HasPrefix(filter.Field, "cloudSpecificDetails.") {
		return h.applyJSONFilter(query, filter)
	}

	// Map JSON field names to database column names
	dbField := h.mapFieldToDBColumn(filter.Field)

	// Handle regular fields
	switch filter.Operator {
	case "eq":
		return query.Where(fmt.Sprintf("%s = ?", dbField), filter.Value), nil
	case "neq":
		return query.Where(fmt.Sprintf("%s != ?", dbField), filter.Value), nil
	case "contains":
		return query.Where(fmt.Sprintf("%s LIKE ?", dbField), fmt.Sprintf("%%%s%%", filter.Value)), nil
	case "lt":
		return query.Where(fmt.Sprintf("%s < ?", dbField), filter.Value), nil
	case "gt":
		return query.Where(fmt.Sprintf("%s > ?", dbField), filter.Value), nil
	case "le":
		return query.Where(fmt.Sprintf("%s <= ?", dbField), filter.Value), nil
	case "ge":
		return query.Where(fmt.Sprintf("%s >= ?", dbField), filter.Value), nil
	default:
		return query, fmt.Errorf("unsupported operator: %s", filter.Operator)
	}
}

// mapFieldToDBColumn maps JSON field names to database column names
func (h *Handler) mapFieldToDBColumn(field string) string {
	switch field {
	case "cloudType":
		return "cloud_type"
	case "cloudAccountId":
		return "cloud_account_id"
	case "instanceType":
		return "instance_type"
	case "createdAt":
		return "created_at"
	case "updatedAt":
		return "updated_at"
	default:
		return field
	}
}

// applyJSONFilter applies filters to JSON fields
func (h *Handler) applyJSONFilter(query *gorm.DB, filter models.VMFilter) (*gorm.DB, error) {
	jsonPath := strings.TrimPrefix(filter.Field, "cloudSpecificDetails.")
	
	switch filter.Operator {
	case "eq":
		return query.Where("JSON_EXTRACT(cloud_specific_details, ?) = ?", fmt.Sprintf("$.%s", jsonPath), filter.Value), nil
	case "neq":
		return query.Where("JSON_EXTRACT(cloud_specific_details, ?) != ?", fmt.Sprintf("$.%s", jsonPath), filter.Value), nil
	case "contains":
		return query.Where("JSON_EXTRACT(cloud_specific_details, ?) LIKE ?", fmt.Sprintf("$.%s", jsonPath), fmt.Sprintf("%%%s%%", filter.Value)), nil
	case "lt":
		return query.Where("JSON_EXTRACT(cloud_specific_details, ?) < ?", fmt.Sprintf("$.%s", jsonPath), filter.Value), nil
	case "gt":
		return query.Where("JSON_EXTRACT(cloud_specific_details, ?) > ?", fmt.Sprintf("$.%s", jsonPath), filter.Value), nil
	case "le":
		return query.Where("JSON_EXTRACT(cloud_specific_details, ?) <= ?", fmt.Sprintf("$.%s", jsonPath), filter.Value), nil
	case "ge":
		return query.Where("JSON_EXTRACT(cloud_specific_details, ?) >= ?", fmt.Sprintf("$.%s", jsonPath), filter.Value), nil
	default:
		return query, fmt.Errorf("unsupported operator: %s", filter.Operator)
	}
}