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

// VMsHandler handles VM-related HTTP requests
type VMsHandler struct {
	db *gorm.DB
}

// NewVMsHandler creates a new VMs handler
func NewVMsHandler(db *gorm.DB) *VMsHandler {
	return &VMsHandler{db: db}
}

// GetVMs handles GET /api/v1/vms
func (h *VMsHandler) GetVMs(c *gin.Context) {
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

	// Get VMs from database
	var vms []models.VM
	if err := h.db.Find(&vms).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to fetch VMs")
		return
	}

	// Apply filters
	filteredVMs := h.applyFilters(vms, params.Filters)

	// Apply sorting
	sortedVMs := h.applySorting(filteredVMs, params.SortBy, params.SortOrder)

	// Calculate pagination
	totalItems := len(sortedVMs)
	paginatedVMs := utils.ApplyPagination(sortedVMs, params.Page, params.PageSize)

	// Send response using the reusable utility
	utils.SendPaginatedResponse(c, paginatedVMs, params.Page, params.PageSize, totalItems)
}

// applyFilters applies filters to VMs using the reusable utility
func (h *VMsHandler) applyFilters(vms []models.VM, filters []utils.QueryFilter) []models.VM {
	// Separate DB-level filters from in-memory filters
	var dbFilters []utils.QueryFilter
	var memoryFilters []utils.QueryFilter

	for _, filter := range filters {
		// Fields that can be filtered at DB level
		dbFields := map[string]bool{
			"name":        true,
			"instance_id": true,
			"region":      true,
			"zone":        true,
			"instance_type": true,
			"cpu_cores":   true,
			"memory_gb":   true,
			"storage_gb":  true,
			"private_ip":  true,
			"public_ip":   true,
			"created_at":  true,
			"updated_at":  true,
		}

		if dbFields[filter.Field] {
			dbFilters = append(dbFilters, filter)
		} else {
			memoryFilters = append(memoryFilters, filter)
		}
	}

	// Apply DB-level filters
	filteredVMs := h.applyDBFilters(vms, dbFilters)

	// Apply in-memory filters using the reusable utility
	if len(memoryFilters) > 0 {
		filteredVMs = utils.ApplyInMemoryFilters(filteredVMs, memoryFilters, h.extractVMField)
	}

	return filteredVMs
}

// applyDBFilters applies database-level filters
func (h *VMsHandler) applyDBFilters(vms []models.VM, filters []utils.QueryFilter) []models.VM {
	if len(filters) == 0 {
		return vms
	}

	// For now, we'll apply all filters in memory since we're using a simple approach
	// In a production environment, you'd want to build dynamic SQL queries
	return utils.ApplyInMemoryFilters(vms, filters, h.extractVMField)
}

// extractVMField extracts field value from VM struct for filtering
func (h *VMsHandler) extractVMField(vm models.VM, field string) interface{} {
	v := reflect.ValueOf(vm)
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

// applySorting applies sorting to VMs
func (h *VMsHandler) applySorting(vms []models.VM, sortBy, sortOrder string) []models.VM {
	if sortBy == "" {
		return vms
	}

	sortedVMs := make([]models.VM, len(vms))
	copy(sortedVMs, vms)

	sort.Slice(sortedVMs, func(i, j int) bool {
		valI := h.extractVMField(sortedVMs[i], sortBy)
		valJ := h.extractVMField(sortedVMs[j], sortBy)

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