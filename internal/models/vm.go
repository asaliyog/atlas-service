package models

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

// BaseVM represents common fields across all cloud providers
type BaseVM struct {
	ID            string         `json:"id" gorm:"primarykey"`
	Name          string         `json:"name"`
	Status        string         `json:"status"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
	Location      string         `json:"location" gorm:"index"`
	InstanceType  string         `json:"instanceType" gorm:"column:instance_type"`
}

// AWSEC2Instance represents AWS EC2 instances
type AWSEC2Instance struct {
	BaseVM
	AccountID                string          `json:"accountId" gorm:"column:account_id;index"`
	VpcID                    string          `json:"vpcId" gorm:"column:vpc_id"`
	SubnetID                 string          `json:"subnetId" gorm:"column:subnet_id"`
	SecurityGroupIDs         json.RawMessage `json:"securityGroupIds" gorm:"column:security_group_ids;type:json"`
	PrivateIPAddress         string          `json:"privateIpAddress" gorm:"column:private_ip_address"`
	PublicIPAddress          string          `json:"publicIpAddress" gorm:"column:public_ip_address"`
	KeyName                  string          `json:"keyName" gorm:"column:key_name"`
	ImageID                  string          `json:"imageId" gorm:"column:image_id"`
	LaunchTime               time.Time       `json:"launchTime" gorm:"column:launch_time"`
	AvailabilityZone         string          `json:"availabilityZone" gorm:"column:availability_zone"`
	PublicDnsName            string          `json:"publicDnsName" gorm:"column:public_dns_name"`
	PrivateDnsName           string          `json:"privateDnsName" gorm:"column:private_dns_name"`
	Architecture             string          `json:"architecture"`
	VirtualizationType       string          `json:"virtualizationType" gorm:"column:virtualization_type"`
	Platform                 string          `json:"platform"`
	RootDeviceType           string          `json:"rootDeviceType" gorm:"column:root_device_type"`
	MonitoringState          string          `json:"monitoringState" gorm:"column:monitoring_state"`
	PlacementGroupName       string          `json:"placementGroupName" gorm:"column:placement_group_name"`
	PlacementPartitionNumber int             `json:"placementPartitionNumber" gorm:"column:placement_partition_number"`
	PlacementTenancy         string          `json:"placementTenancy" gorm:"column:placement_tenancy"`
	SpotInstanceRequestID    string          `json:"spotInstanceRequestId" gorm:"column:spot_instance_request_id"`
	SriovNetSupport          string          `json:"sriovNetSupport" gorm:"column:sriov_net_support"`
	EbsOptimized             bool            `json:"ebsOptimized" gorm:"column:ebs_optimized"`
	EnaSupport               bool            `json:"enaSupport" gorm:"column:ena_support"`
	SourceDestCheck          bool            `json:"sourceDestCheck" gorm:"column:source_dest_check"`
	Tags                     json.RawMessage `json:"tags" gorm:"type:json"`
}

// TableName returns the table name for AWSEC2Instance
func (AWSEC2Instance) TableName() string {
	return "aws_ec2_instances"
}

// AzureVMInstance represents Azure Virtual Machines
type AzureVMInstance struct {
	BaseVM
	SubscriptionID    string          `json:"subscriptionId" gorm:"column:subscription_id;index"`
	ResourceGroup     string          `json:"resourceGroup" gorm:"column:resource_group"`
	VMSize            string          `json:"vmSize" gorm:"column:vm_size"`
	PrivateIPAddress  string          `json:"privateIpAddress" gorm:"column:private_ip_address"`
	PublicIPAddress   string          `json:"publicIpAddress" gorm:"column:public_ip_address"`
	NetworkInterfaces json.RawMessage `json:"networkInterfaces" gorm:"column:network_interfaces;type:json"`
	OSDisk            json.RawMessage `json:"osDisk" gorm:"column:os_disk;type:json"`
	DataDisks         json.RawMessage `json:"dataDisks" gorm:"column:data_disks;type:json"`
	OSType            string          `json:"osType" gorm:"column:os_type"`
	OSProfile         json.RawMessage `json:"osProfile" gorm:"column:os_profile;type:json"`
	HardwareProfile   json.RawMessage `json:"hardwareProfile" gorm:"column:hardware_profile;type:json"`
	StorageProfile    json.RawMessage `json:"storageProfile" gorm:"column:storage_profile;type:json"`
	NetworkProfile    json.RawMessage `json:"networkProfile" gorm:"column:network_profile;type:json"`
	SecurityProfile   json.RawMessage `json:"securityProfile" gorm:"column:security_profile;type:json"`
	DiagnosticsProfile json.RawMessage `json:"diagnosticsProfile" gorm:"column:diagnostics_profile;type:json"`
	AvailabilitySet   json.RawMessage `json:"availabilitySet" gorm:"column:availability_set;type:json"`
	VirtualMachineScaleSet json.RawMessage `json:"virtualMachineScaleSet" gorm:"column:virtual_machine_scale_set;type:json"`
	ProximityPlacementGroup json.RawMessage `json:"proximityPlacementGroup" gorm:"column:proximity_placement_group;type:json"`
	Priority          string          `json:"priority"`
	EvictionPolicy    string          `json:"evictionPolicy" gorm:"column:eviction_policy"`
	BillingProfile    json.RawMessage `json:"billingProfile" gorm:"column:billing_profile;type:json"`
	HostId            string          `json:"hostId" gorm:"column:host_id"`
	LicenseType       string          `json:"licenseType" gorm:"column:license_type"`
	VMId              string          `json:"vmId" gorm:"column:vm_id"`
	Tags              json.RawMessage `json:"tags" gorm:"type:json"`
}

// TableName returns the table name for AzureVMInstance
func (AzureVMInstance) TableName() string {
	return "azure_vm_instances"
}

// GCPComputeInstance represents GCP Compute Engine instances
type GCPComputeInstance struct {
	BaseVM
	ProjectID                string          `json:"projectId" gorm:"column:project_id;index"`
	Zone                     string          `json:"zone"`
	MachineType              string          `json:"machineType" gorm:"column:machine_type"`
	PrivateIPAddress         string          `json:"privateIpAddress" gorm:"column:private_ip_address"`
	PublicIPAddress          string          `json:"publicIpAddress" gorm:"column:public_ip_address"`
	NetworkInterfaces        json.RawMessage `json:"networkInterfaces" gorm:"column:network_interfaces;type:json"`
	Disks                    json.RawMessage `json:"disks" gorm:"type:json"`
	Metadata                 json.RawMessage `json:"metadata" gorm:"type:json"`
	Tags                     json.RawMessage `json:"tags" gorm:"type:json"`
	Labels                   json.RawMessage `json:"labels" gorm:"type:json"`
	ServiceAccounts          json.RawMessage `json:"serviceAccounts" gorm:"column:service_accounts;type:json"`
	Scheduling               json.RawMessage `json:"scheduling" gorm:"type:json"`
	CpuPlatform              string          `json:"cpuPlatform" gorm:"column:cpu_platform"`
	MinCpuPlatform           string          `json:"minCpuPlatform" gorm:"column:min_cpu_platform"`
	GuestAccelerators        json.RawMessage `json:"guestAccelerators" gorm:"column:guest_accelerators;type:json"`
	ShieldedInstanceConfig   json.RawMessage `json:"shieldedInstanceConfig" gorm:"column:shielded_instance_config;type:json"`
	ConfidentialInstanceConfig json.RawMessage `json:"confidentialInstanceConfig" gorm:"column:confidential_instance_config;type:json"`
	DisplayDevice            json.RawMessage `json:"displayDevice" gorm:"column:display_device;type:json"`
	KeyRevocationActionType  string          `json:"keyRevocationActionType" gorm:"column:key_revocation_action_type"`
	SourceMachineImage       string          `json:"sourceMachineImage" gorm:"column:source_machine_image"`
	ResourcePolicies         json.RawMessage `json:"resourcePolicies" gorm:"column:resource_policies;type:json"`
	ReservationAffinity      json.RawMessage `json:"reservationAffinity" gorm:"column:reservation_affinity;type:json"`
	AdvancedMachineFeatures  json.RawMessage `json:"advancedMachineFeatures" gorm:"column:advanced_machine_features;type:json"`
	Fingerprint              string          `json:"fingerprint"`
	LastStartTimestamp       time.Time       `json:"lastStartTimestamp" gorm:"column:last_start_timestamp"`
	LastStopTimestamp        time.Time       `json:"lastStopTimestamp" gorm:"column:last_stop_timestamp"`
	LastSuspendedTimestamp   time.Time       `json:"lastSuspendedTimestamp" gorm:"column:last_suspended_timestamp"`
	SatisfiesPzs             bool            `json:"satisfiesPzs" gorm:"column:satisfies_pzs"`
	Hostname                 string          `json:"hostname"`
	InstanceEncryptionKey    json.RawMessage `json:"instanceEncryptionKey" gorm:"column:instance_encryption_key;type:json"`
	PrivateIpv6GoogleAccess  string          `json:"privateIpv6GoogleAccess" gorm:"column:private_ipv6_google_access"`
}

// TableName returns the table name for GCPComputeInstance
func (GCPComputeInstance) TableName() string {
	return "gcp_compute_instances"
}

// VM represents a unified virtual machine view across cloud providers
type VM struct {
	ID                   string                 `json:"id"`
	Name                 string                 `json:"name"`
	CloudType            string                 `json:"cloudType"`
	Status               string                 `json:"status"`
	CreatedAt            time.Time              `json:"createdAt"`
	UpdatedAt            time.Time              `json:"updatedAt"`
	CloudAccountID       string                 `json:"cloudAccountId"`
	Location             string                 `json:"location"`
	InstanceType         string                 `json:"instanceType"`
	CloudSpecificDetails json.RawMessage        `json:"cloudSpecificDetails"`
}

// VMListResponse represents the response for VM list endpoint
type VMListResponse struct {
	Data       []VM       `json:"data"`
	Pagination Pagination `json:"pagination"`
}

// Pagination represents pagination information
type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}

// VMFilter represents a filter for VM queries
type VMFilter struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

// ValidateFilter validates a VMFilter for type safety and operator support
func (f VMFilter) ValidateFilter() error {
	// Validate operator
	validOperators := map[string]bool{
		"eq":      true,
		"ne":      true,
		"lt":      true,
		"lte":     true,
		"gt":      true,
		"gte":     true,
		"in":      true,
		"nin":     true,
		"like":    true,
		"ilike":   true,
		"between": true,
		"null":    true,
	}

	if !validOperators[f.Operator] {
		return fmt.Errorf("unsupported operator: %s", f.Operator)
	}

	// Validate field
	if f.Field == "" {
		return fmt.Errorf("field cannot be empty")
	}

	// Validate value based on operator
	switch f.Operator {
	case "in", "nin":
		// Value should be a slice
		if f.Value == nil {
			return fmt.Errorf("value cannot be nil for operator %s", f.Operator)
		}
		v := reflect.ValueOf(f.Value)
		if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
			return fmt.Errorf("value must be an array for operator %s", f.Operator)
		}
		if v.Len() == 0 {
			return fmt.Errorf("value array cannot be empty for operator %s", f.Operator)
		}
	case "between":
		// Value should be an array with exactly 2 elements
		if f.Value == nil {
			return fmt.Errorf("value cannot be nil for operator %s", f.Operator)
		}
		v := reflect.ValueOf(f.Value)
		if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
			return fmt.Errorf("value must be an array for operator %s", f.Operator)
		}
		if v.Len() != 2 {
			return fmt.Errorf("value array must have exactly 2 elements for operator %s", f.Operator)
		}
		// Check that both values are of the same type
		first := v.Index(0).Interface()
		second := v.Index(1).Interface()
		if !areComparableTypes(first, second) {
			return fmt.Errorf("between values must be of comparable types")
		}
	case "null":
		// Value should be a boolean indicating whether to check for null or not null
		if f.Value == nil {
			return fmt.Errorf("value cannot be nil for operator %s", f.Operator)
		}
		if _, ok := f.Value.(bool); !ok {
			return fmt.Errorf("value must be a boolean for operator %s", f.Operator)
		}
	case "like", "ilike":
		// Value should be a string
		if f.Value == nil {
			return fmt.Errorf("value cannot be nil for operator %s", f.Operator)
		}
		if _, ok := f.Value.(string); !ok {
			return fmt.Errorf("value must be a string for operator %s", f.Operator)
		}
	case "lt", "lte", "gt", "gte":
		// Value should be a comparable type (number, string, time)
		if f.Value == nil {
			return fmt.Errorf("value cannot be nil for operator %s", f.Operator)
		}
		if !isComparableType(f.Value) {
			return fmt.Errorf("value must be a comparable type (number, string, time) for operator %s", f.Operator)
		}
	}

	return nil
}

// areComparableTypes checks if two values are of comparable types
func areComparableTypes(a, b interface{}) bool {
	typeA := reflect.TypeOf(a)
	typeB := reflect.TypeOf(b)
	
	// Same type
	if typeA == typeB {
		return true
	}
	
	// Both are numeric
	if isNumericType(typeA) && isNumericType(typeB) {
		return true
	}
	
	return false
}

// isComparableType checks if a type is comparable for ordering operations
func isComparableType(value interface{}) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64,
		 uint, uint8, uint16, uint32, uint64,
		 float32, float64,
		 string, time.Time:
		return true
	}
	return false
}

// isNumericType checks if a type is numeric
func isNumericType(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		 reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		 reflect.Float32, reflect.Float64:
		return true
	}
	return false
}

// VMQueryParams represents the query parameters for VM endpoints
type VMQueryParams struct {
	Page      int        `json:"page"`
	PageSize  int        `json:"pageSize"`
	SortBy    string     `json:"sortBy"`
	SortOrder string     `json:"sortOrder"`
	Filters   []VMFilter `json:"filters"`
}

// ValidateQueryParams validates the query parameters
func (p VMQueryParams) ValidateQueryParams() error {
	// Validate page
	if p.Page < 1 {
		return fmt.Errorf("page must be greater than 0")
	}

	// Validate page size
	if p.PageSize < 1 {
		return fmt.Errorf("pageSize must be greater than 0")
	}
	if p.PageSize > 1000 {
		return fmt.Errorf("pageSize cannot exceed 1000")
	}

	// Validate sort order
	if p.SortOrder != "" && p.SortOrder != "asc" && p.SortOrder != "desc" {
		return fmt.Errorf("sortOrder must be 'asc' or 'desc'")
	}

	// Validate filters
	for i, filter := range p.Filters {
		if err := filter.ValidateFilter(); err != nil {
			return fmt.Errorf("filter %d: %w", i, err)
		}
	}

	return nil
}

// Error represents an error response
type Error struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// ConvertToValue converts interface{} to appropriate type for database queries
func ConvertToValue(value interface{}) (interface{}, error) {
	if value == nil {
		return nil, nil
	}

	switch v := value.(type) {
	case string:
		// Try to convert string to number if it looks like a number
		if intVal, err := strconv.Atoi(v); err == nil {
			return intVal, nil
		}
		if floatVal, err := strconv.ParseFloat(v, 64); err == nil {
			return floatVal, nil
		}
		return v, nil
	case float64:
		// JSON unmarshaling often converts numbers to float64
		if v == float64(int64(v)) {
			return int64(v), nil
		}
		return v, nil
	default:
		return v, nil
	}
}

// SanitizeFieldName sanitizes field names for database queries
func SanitizeFieldName(field string) string {
	// Remove any potentially dangerous characters
	field = strings.ReplaceAll(field, ";", "")
	field = strings.ReplaceAll(field, "--", "")
	field = strings.ReplaceAll(field, "/*", "")
	field = strings.ReplaceAll(field, "*/", "")
	return field
}