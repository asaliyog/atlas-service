package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"golang-service/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetVMs godoc
// @Summary Retrieve a list of virtual machines
// @Description Fetches a paginated list of virtual machines across AWS EC2, GCP Compute, and Azure VMs with advanced filtering and sorting
// @Tags vms
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number for pagination (1-based)" default(1)
// @Param pageSize query int false "Number of VMs per page" default(20)
// @Param sortBy query string false "Field to sort by" default("createdAt")
// @Param sortOrder query string false "Sort order (asc or desc)" default("asc")
// @Param filter query string false "JSON string of filter objects with operators: eq, ne, lt, lte, gt, gte, in, nin, like, ilike, between, null"
// @Success 200 {object} models.VMListResponse
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /vms [get]
func (h *Handler) GetVMs(c *gin.Context) {
	// Parse and validate query parameters
	params, err := h.parseVMQueryParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: "Invalid query parameters",
			Code:    "INVALID_PARAMS",
			Details: err.Error(),
		})
		return
	}

	// Validate query parameters
	if err := params.ValidateQueryParams(); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: "Query parameter validation failed",
			Code:    "VALIDATION_ERROR",
			Details: err.Error(),
		})
		return
	}

	// Fetch VMs from all cloud providers
	vms, totalCount, err := h.fetchVMsFromAllClouds(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: "Failed to fetch VMs",
			Code:    "FETCH_ERROR",
			Details: err.Error(),
		})
		return
	}

	// Calculate total pages
	totalPages := int(totalCount) / params.PageSize
	if int(totalCount)%params.PageSize != 0 {
		totalPages++
	}

	response := models.VMListResponse{
		Data: vms,
		Pagination: models.Pagination{
			Page:       params.Page,
			PageSize:   params.PageSize,
			TotalItems: int(totalCount),
			TotalPages: totalPages,
		},
	}

	c.JSON(http.StatusOK, response)
}

// fetchVMsFromAllClouds fetches VMs from all three cloud providers and combines them
func (h *Handler) fetchVMsFromAllClouds(params models.VMQueryParams) ([]models.VM, int64, error) {
	var allVMs []models.VM
	var totalCount int64

	// Fetch from AWS EC2
	awsVMs, awsCount, err := h.fetchAWSVMs(params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch AWS VMs: %w", err)
	}
	allVMs = append(allVMs, awsVMs...)
	totalCount += awsCount

	// Fetch from Azure VMs
	azureVMs, azureCount, err := h.fetchAzureVMs(params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch Azure VMs: %w", err)
	}
	allVMs = append(allVMs, azureVMs...)
	totalCount += azureCount

	// Fetch from GCP Compute
	gcpVMs, gcpCount, err := h.fetchGCPVMs(params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch GCP VMs: %w", err)
	}
	allVMs = append(allVMs, gcpVMs...)
	totalCount += gcpCount

	// Apply cross-cloud sorting and pagination
	if len(allVMs) > 0 {
		allVMs = h.applyCrossCloudSorting(allVMs, params.SortBy, params.SortOrder)
		allVMs = h.applyCrossCloudPagination(allVMs, params.Page, params.PageSize)
	}

	return allVMs, totalCount, nil
}

// fetchAWSVMs fetches VMs from AWS EC2 table
func (h *Handler) fetchAWSVMs(params models.VMQueryParams) ([]models.VM, int64, error) {
	query := h.db.Model(&models.AWSEC2Instance{})
	
	// Apply filters
	filteredQuery, err := h.applyFiltersToQuery(query, params.Filters, "aws")
	if err != nil {
		return nil, 0, err
	}

	// Get count
	var count int64
	if err := filteredQuery.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	sortedQuery := h.applySortingToQuery(filteredQuery, params.SortBy, params.SortOrder)

	// Fetch instances
	var instances []models.AWSEC2Instance
	if err := sortedQuery.Find(&instances).Error; err != nil {
		return nil, 0, err
	}

	// Convert to unified VM format
	var vms []models.VM
	for _, instance := range instances {
		vm := models.VM{
			ID:             instance.ID,
			Name:           instance.Name,
			CloudType:      "aws",
			Status:         instance.Status,
			CreatedAt:      instance.CreatedAt,
			UpdatedAt:      instance.UpdatedAt,
			CloudAccountID: instance.AccountID,
			Location:       instance.Location,
			InstanceType:   instance.InstanceType,
		}
		
		// Convert cloud-specific details to JSON
		cloudDetails, _ := json.Marshal(instance)
		vm.CloudSpecificDetails = cloudDetails
		
		vms = append(vms, vm)
	}

	return vms, count, nil
}

// fetchAzureVMs fetches VMs from Azure VM table
func (h *Handler) fetchAzureVMs(params models.VMQueryParams) ([]models.VM, int64, error) {
	query := h.db.Model(&models.AzureVMInstance{})
	
	// Apply filters
	filteredQuery, err := h.applyFiltersToQuery(query, params.Filters, "azure")
	if err != nil {
		return nil, 0, err
	}

	// Get count
	var count int64
	if err := filteredQuery.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	sortedQuery := h.applySortingToQuery(filteredQuery, params.SortBy, params.SortOrder)

	// Fetch instances
	var instances []models.AzureVMInstance
	if err := sortedQuery.Find(&instances).Error; err != nil {
		return nil, 0, err
	}

	// Convert to unified VM format
	var vms []models.VM
	for _, instance := range instances {
		vm := models.VM{
			ID:             instance.ID,
			Name:           instance.Name,
			CloudType:      "azure",
			Status:         instance.Status,
			CreatedAt:      instance.CreatedAt,
			UpdatedAt:      instance.UpdatedAt,
			CloudAccountID: instance.SubscriptionID,
			Location:       instance.Location,
			InstanceType:   instance.InstanceType,
		}
		
		// Convert cloud-specific details to JSON
		cloudDetails, _ := json.Marshal(instance)
		vm.CloudSpecificDetails = cloudDetails
		
		vms = append(vms, vm)
	}

	return vms, count, nil
}

// fetchGCPVMs fetches VMs from GCP Compute table
func (h *Handler) fetchGCPVMs(params models.VMQueryParams) ([]models.VM, int64, error) {
	query := h.db.Model(&models.GCPComputeInstance{})
	
	// Apply filters
	filteredQuery, err := h.applyFiltersToQuery(query, params.Filters, "gcp")
	if err != nil {
		return nil, 0, err
	}

	// Get count
	var count int64
	if err := filteredQuery.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	sortedQuery := h.applySortingToQuery(filteredQuery, params.SortBy, params.SortOrder)

	// Fetch instances
	var instances []models.GCPComputeInstance
	if err := sortedQuery.Find(&instances).Error; err != nil {
		return nil, 0, err
	}

	// Convert to unified VM format
	var vms []models.VM
	for _, instance := range instances {
		vm := models.VM{
			ID:             instance.ID,
			Name:           instance.Name,
			CloudType:      "gcp",
			Status:         instance.Status,
			CreatedAt:      instance.CreatedAt,
			UpdatedAt:      instance.UpdatedAt,
			CloudAccountID: instance.ProjectID,
			Location:       instance.Location,
			InstanceType:   instance.InstanceType,
		}
		
		// Convert cloud-specific details to JSON
		cloudDetails, _ := json.Marshal(instance)
		vm.CloudSpecificDetails = cloudDetails
		
		vms = append(vms, vm)
	}

	return vms, count, nil
}

// applyFiltersToQuery applies advanced filters to a GORM query
func (h *Handler) applyFiltersToQuery(query *gorm.DB, filters []models.VMFilter, cloudType string) (*gorm.DB, error) {
	for _, filter := range filters {
		// Validate filter
		if err := filter.ValidateFilter(); err != nil {
			return nil, fmt.Errorf("filter validation failed: %w", err)
		}

		// Special handling for cloudType filter
		if filter.Field == "cloudType" {
			// If filter matches current cloud type, we can skip it (all records in this table are of this type)
			if filter.Operator == "eq" && filter.Value == cloudType {
				continue
			}
			// If filter doesn't match current cloud type with eq, exclude all records
			if filter.Operator == "eq" && filter.Value != cloudType {
				query = query.Where("1 = 0") // This will exclude all records
				continue
			}
			// If filter is "not equal" and matches current cloud type, exclude all records
			if filter.Operator == "ne" && filter.Value == cloudType {
				query = query.Where("1 = 0") // This will exclude all records
				continue
			}
			// If filter is "not equal" and doesn't match current cloud type, include all records (skip filter)
			if filter.Operator == "ne" && filter.Value != cloudType {
				continue
			}
			// For other operators like "in", "nin", check if current cloud type matches
			if filter.Operator == "in" {
				v := reflect.ValueOf(filter.Value)
				if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
					found := false
					for i := 0; i < v.Len(); i++ {
						if v.Index(i).Interface() == cloudType {
							found = true
							break
						}
					}
					if !found {
						query = query.Where("1 = 0") // Exclude all records
					}
					continue
				}
			}
			if filter.Operator == "nin" {
				v := reflect.ValueOf(filter.Value)
				if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
					found := false
					for i := 0; i < v.Len(); i++ {
						if v.Index(i).Interface() == cloudType {
							found = true
							break
						}
					}
					if found {
						query = query.Where("1 = 0") // Exclude all records
					}
					continue
				}
			}
		}

		// Apply filter
		filteredQuery, err := h.applyAdvancedFilter(query, filter, cloudType)
		if err != nil {
			return nil, fmt.Errorf("failed to apply filter: %w", err)
		}
		query = filteredQuery
	}
	return query, nil
}

// applyAdvancedFilter applies a single advanced filter with proper type checking
func (h *Handler) applyAdvancedFilter(query *gorm.DB, filter models.VMFilter, cloudType string) (*gorm.DB, error) {
	// Map field to database column
	dbField := h.mapFieldToDBColumn(filter.Field, cloudType)
	
	// Convert value to appropriate type
	value, err := models.ConvertToValue(filter.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to convert filter value: %w", err)
	}

	// Apply filter based on operator
	switch filter.Operator {
	case "eq":
		return query.Where(fmt.Sprintf("%s = ?", dbField), value), nil
	case "ne":
		return query.Where(fmt.Sprintf("%s != ?", dbField), value), nil
	case "lt":
		return query.Where(fmt.Sprintf("%s < ?", dbField), value), nil
	case "lte":
		return query.Where(fmt.Sprintf("%s <= ?", dbField), value), nil
	case "gt":
		return query.Where(fmt.Sprintf("%s > ?", dbField), value), nil
	case "gte":
		return query.Where(fmt.Sprintf("%s >= ?", dbField), value), nil
	case "in":
		// Handle array values
		v := reflect.ValueOf(filter.Value)
		if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
			var values []interface{}
			for i := 0; i < v.Len(); i++ {
				convertedVal, _ := models.ConvertToValue(v.Index(i).Interface())
				values = append(values, convertedVal)
			}
			return query.Where(fmt.Sprintf("%s IN ?", dbField), values), nil
		}
		return nil, fmt.Errorf("'in' operator requires array value")
	case "nin":
		// Handle array values
		v := reflect.ValueOf(filter.Value)
		if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
			var values []interface{}
			for i := 0; i < v.Len(); i++ {
				convertedVal, _ := models.ConvertToValue(v.Index(i).Interface())
				values = append(values, convertedVal)
			}
			return query.Where(fmt.Sprintf("%s NOT IN ?", dbField), values), nil
		}
		return nil, fmt.Errorf("'nin' operator requires array value")
	case "like":
		return query.Where(fmt.Sprintf("%s LIKE ?", dbField), fmt.Sprintf("%%%s%%", value)), nil
	case "ilike":
		return query.Where(fmt.Sprintf("LOWER(%s) LIKE LOWER(?)", dbField), fmt.Sprintf("%%%s%%", value)), nil
	case "between":
		v := reflect.ValueOf(filter.Value)
		if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
			if v.Len() == 2 {
				val1, _ := models.ConvertToValue(v.Index(0).Interface())
				val2, _ := models.ConvertToValue(v.Index(1).Interface())
				return query.Where(fmt.Sprintf("%s BETWEEN ? AND ?", dbField), val1, val2), nil
			}
		}
		return nil, fmt.Errorf("'between' operator requires array with 2 values")
	case "null":
		isNull, ok := filter.Value.(bool)
		if !ok {
			return nil, fmt.Errorf("'null' operator requires boolean value")
		}
		if isNull {
			return query.Where(fmt.Sprintf("%s IS NULL", dbField)), nil
		}
		return query.Where(fmt.Sprintf("%s IS NOT NULL", dbField)), nil
	default:
		return nil, fmt.Errorf("unsupported operator: %s", filter.Operator)
	}
}

// mapFieldToDBColumn maps JSON field names to database column names for each cloud type
func (h *Handler) mapFieldToDBColumn(field string, cloudType string) string {
	// Sanitize field name
	field = models.SanitizeFieldName(field)
	
	// Common field mappings
	switch field {
	case "cloudType":
		return "cloud_type"
	case "cloudAccountId":
		switch cloudType {
		case "aws":
			return "account_id"
		case "azure":
			return "subscription_id"
		case "gcp":
			return "project_id"
		}
	case "instanceType":
		switch cloudType {
		case "aws":
			return "instance_type"
		case "azure":
			return "vm_size"
		case "gcp":
			return "machine_type"
		}
	case "createdAt":
		return "created_at"
	case "updatedAt":
		return "updated_at"
	case "privateIpAddress":
		return "private_ip_address"
	case "publicIpAddress":
		return "public_ip_address"
	}
	
	// Cloud-specific mappings
	switch cloudType {
	case "aws":
		switch field {
		case "vpcId":
			return "vpc_id"
		case "subnetId":
			return "subnet_id"
		case "securityGroupIds":
			return "security_group_ids"
		case "keyName":
			return "key_name"
		case "imageId":
			return "image_id"
		case "launchTime":
			return "launch_time"
		case "availabilityZone":
			return "availability_zone"
		case "publicDnsName":
			return "public_dns_name"
		case "privateDnsName":
			return "private_dns_name"
		case "virtualizationType":
			return "virtualization_type"
		case "rootDeviceType":
			return "root_device_type"
		case "monitoringState":
			return "monitoring_state"
		case "placementGroupName":
			return "placement_group_name"
		case "placementPartitionNumber":
			return "placement_partition_number"
		case "placementTenancy":
			return "placement_tenancy"
		case "spotInstanceRequestId":
			return "spot_instance_request_id"
		case "sriovNetSupport":
			return "sriov_net_support"
		case "ebsOptimized":
			return "ebs_optimized"
		case "enaSupport":
			return "ena_support"
		case "sourceDestCheck":
			return "source_dest_check"
		}
	case "azure":
		switch field {
		case "resourceGroup":
			return "resource_group"
		case "vmSize":
			return "vm_size"
		case "networkInterfaces":
			return "network_interfaces"
		case "osDisk":
			return "os_disk"
		case "dataDisks":
			return "data_disks"
		case "osType":
			return "os_type"
		case "osProfile":
			return "os_profile"
		case "hardwareProfile":
			return "hardware_profile"
		case "storageProfile":
			return "storage_profile"
		case "networkProfile":
			return "network_profile"
		case "securityProfile":
			return "security_profile"
		case "diagnosticsProfile":
			return "diagnostics_profile"
		case "availabilitySet":
			return "availability_set"
		case "virtualMachineScaleSet":
			return "virtual_machine_scale_set"
		case "proximityPlacementGroup":
			return "proximity_placement_group"
		case "evictionPolicy":
			return "eviction_policy"
		case "billingProfile":
			return "billing_profile"
		case "hostId":
			return "host_id"
		case "licenseType":
			return "license_type"
		case "vmId":
			return "vm_id"
		}
	case "gcp":
		switch field {
		case "projectId":
			return "project_id"
		case "machineType":
			return "machine_type"
		case "networkInterfaces":
			return "network_interfaces"
		case "serviceAccounts":
			return "service_accounts"
		case "cpuPlatform":
			return "cpu_platform"
		case "minCpuPlatform":
			return "min_cpu_platform"
		case "guestAccelerators":
			return "guest_accelerators"
		case "shieldedInstanceConfig":
			return "shielded_instance_config"
		case "confidentialInstanceConfig":
			return "confidential_instance_config"
		case "displayDevice":
			return "display_device"
		case "keyRevocationActionType":
			return "key_revocation_action_type"
		case "sourceMachineImage":
			return "source_machine_image"
		case "resourcePolicies":
			return "resource_policies"
		case "reservationAffinity":
			return "reservation_affinity"
		case "advancedMachineFeatures":
			return "advanced_machine_features"
		case "lastStartTimestamp":
			return "last_start_timestamp"
		case "lastStopTimestamp":
			return "last_stop_timestamp"
		case "lastSuspendedTimestamp":
			return "last_suspended_timestamp"
		case "satisfiesPzs":
			return "satisfies_pzs"
		case "instanceEncryptionKey":
			return "instance_encryption_key"
		case "privateIpv6GoogleAccess":
			return "private_ipv6_google_access"
		}
	}
	
	// Return field as-is if no mapping found
	return field
}

// applySortingToQuery applies sorting to a GORM query
func (h *Handler) applySortingToQuery(query *gorm.DB, sortBy, sortOrder string) *gorm.DB {
	if sortBy == "" {
		sortBy = "created_at"
	}
	
	if sortOrder == "" || (sortOrder != "asc" && sortOrder != "desc") {
		sortOrder = "asc"
	}
	
	// Map field to database column (using aws as default for common fields)
	dbField := h.mapFieldToDBColumn(sortBy, "aws")
	
	return query.Order(fmt.Sprintf("%s %s", dbField, strings.ToUpper(sortOrder)))
}

// applyCrossCloudSorting applies sorting across VMs from different clouds
func (h *Handler) applyCrossCloudSorting(vms []models.VM, sortBy, sortOrder string) []models.VM {
	if len(vms) <= 1 {
		return vms
	}
	
	// Implement sorting logic here
	// For simplicity, we'll sort by the most common fields
	// In production, you might want to implement a more sophisticated sorting mechanism
	
	// Default sorting by createdAt
	if sortBy == "" || sortBy == "createdAt" {
		if sortOrder == "desc" {
			// Sort by CreatedAt descending
			for i := 0; i < len(vms)-1; i++ {
				for j := i + 1; j < len(vms); j++ {
					if vms[i].CreatedAt.Before(vms[j].CreatedAt) {
						vms[i], vms[j] = vms[j], vms[i]
					}
				}
			}
		} else {
			// Sort by CreatedAt ascending
			for i := 0; i < len(vms)-1; i++ {
				for j := i + 1; j < len(vms); j++ {
					if vms[i].CreatedAt.After(vms[j].CreatedAt) {
						vms[i], vms[j] = vms[j], vms[i]
					}
				}
			}
		}
	}
	
	return vms
}

// applyCrossCloudPagination applies pagination to the combined VM list
func (h *Handler) applyCrossCloudPagination(vms []models.VM, page, pageSize int) []models.VM {
	if len(vms) == 0 {
		return vms
	}
	
	start := (page - 1) * pageSize
	end := start + pageSize
	
	if start >= len(vms) {
		return []models.VM{}
	}
	
	if end > len(vms) {
		end = len(vms)
	}
	
	return vms[start:end]
}

// parseVMQueryParams parses and validates query parameters
func (h *Handler) parseVMQueryParams(c *gin.Context) (models.VMQueryParams, error) {
	params := models.VMQueryParams{
		Page:      1,
		PageSize:  20,
		SortBy:    "createdAt",
		SortOrder: "asc",
		Filters:   []models.VMFilter{},
	}

	// Parse page
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			params.Page = page
		} else {
			return params, fmt.Errorf("invalid page parameter: %s", pageStr)
		}
	}

	// Parse pageSize
	if pageSizeStr := c.Query("pageSize"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 {
			if pageSize > 1000 {
				return params, fmt.Errorf("pageSize cannot exceed 1000")
			}
			params.PageSize = pageSize
		} else {
			return params, fmt.Errorf("invalid pageSize parameter: %s", pageSizeStr)
		}
	}

	// Parse sortBy
	if sortBy := c.Query("sortBy"); sortBy != "" {
		params.SortBy = sortBy
	}

	// Parse sortOrder
	if sortOrder := c.Query("sortOrder"); sortOrder != "" {
		if sortOrder != "asc" && sortOrder != "desc" {
			return params, fmt.Errorf("sortOrder must be 'asc' or 'desc'")
		}
		params.SortOrder = sortOrder
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