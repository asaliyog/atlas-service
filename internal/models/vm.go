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
	CqSyncTime                              time.Time       `json:"-" gorm:"column:_cq_sync_time"`
	CqSourceName                            string          `json:"-" gorm:"column:_cq_source_name"`
	CqID                                    string          `json:"-" gorm:"column:_cq_id;primarykey"`
	CqParentID                              string          `json:"-" gorm:"column:_cq_parent_id"`
	AccountID                               string          `json:"accountId" gorm:"column:account_id;index"`
	Region                                  string          `json:"region"`
	ARN                                     string          `json:"arn" gorm:"column:arn;primarykey"`
	StateTransitionReasonTime               *time.Time      `json:"-" gorm:"column:state_transition_reason_time"`
	Tags                                    json.RawMessage `json:"tags" gorm:"type:json"`
	AmiLaunchIndex                          *int64          `json:"-" gorm:"column:ami_launch_index"`
	Architecture                            string          `json:"architecture"`
	BlockDeviceMappings                     json.RawMessage `json:"-" gorm:"column:block_device_mappings;type:json"`
	BootMode                                string          `json:"-" gorm:"column:boot_mode"`
	CapacityReservationID                   string          `json:"-" gorm:"column:capacity_reservation_id"`
	CapacityReservationSpecification        json.RawMessage `json:"-" gorm:"column:capacity_reservation_specification;type:json"`
	ClientToken                             string          `json:"-" gorm:"column:client_token"`
	CpuOptions                              json.RawMessage `json:"-" gorm:"column:cpu_options;type:json"`
	CurrentInstanceBootMode                 string          `json:"-" gorm:"column:current_instance_boot_mode"`
	EbsOptimized                            *bool           `json:"ebsOptimized" gorm:"column:ebs_optimized"`
	ElasticGpuAssociations                  json.RawMessage `json:"-" gorm:"column:elastic_gpu_associations;type:json"`
	ElasticInferenceAcceleratorAssociations json.RawMessage `json:"-" gorm:"column:elastic_inference_accelerator_associations;type:json"`
	EnaSupport                              *bool           `json:"enaSupport" gorm:"column:ena_support"`
	EnclaveOptions                          json.RawMessage `json:"-" gorm:"column:enclave_options;type:json"`
	HibernationOptions                      json.RawMessage `json:"-" gorm:"column:hibernation_options;type:json"`
	Hypervisor                              string          `json:"hypervisor"`
	IamInstanceProfile                      json.RawMessage `json:"-" gorm:"column:iam_instance_profile;type:json"`
	ImageID                                 string          `json:"imageId" gorm:"column:image_id"`
	InstanceID                              string          `json:"instanceId" gorm:"column:instance_id"`
	InstanceLifecycle                       string          `json:"-" gorm:"column:instance_lifecycle"`
	InstanceType                            string          `json:"instanceType" gorm:"column:instance_type"`
	Ipv6Address                             string          `json:"-" gorm:"column:ipv6_address"`
	KernelID                                string          `json:"-" gorm:"column:kernel_id"`
	KeyName                                 string          `json:"keyName" gorm:"column:key_name"`
	LaunchTime                              *time.Time      `json:"launchTime" gorm:"column:launch_time"`
	Licenses                                json.RawMessage `json:"-" gorm:"column:licenses;type:json"`
	MaintenanceOptions                      json.RawMessage `json:"-" gorm:"column:maintenance_options;type:json"`
	MetadataOptions                         json.RawMessage `json:"-" gorm:"column:metadata_options;type:json"`
	Monitoring                              json.RawMessage `json:"-" gorm:"column:monitoring;type:json"`
	NetworkInterfaces                       json.RawMessage `json:"-" gorm:"column:network_interfaces;type:json"`
	OutpostARN                              string          `json:"-" gorm:"column:outpost_arn"`
	Placement                               json.RawMessage `json:"-" gorm:"column:placement;type:json"`
	Platform                                string          `json:"platform"`
	PlatformDetails                         string          `json:"-" gorm:"column:platform_details"`
	PrivateDnsName                          string          `json:"privateDnsName" gorm:"column:private_dns_name"`
	PrivateDnsNameOptions                   json.RawMessage `json:"-" gorm:"column:private_dns_name_options;type:json"`
	PrivateIPAddress                        string          `json:"privateIpAddress" gorm:"column:private_ip_address"`
	ProductCodes                            json.RawMessage `json:"-" gorm:"column:product_codes;type:json"`
	PublicDnsName                           string          `json:"publicDnsName" gorm:"column:public_dns_name"`
	PublicIPAddress                         string          `json:"publicIpAddress" gorm:"column:public_ip_address"`
	RamdiskID                               string          `json:"-" gorm:"column:ramdisk_id"`
	RootDeviceName                          string          `json:"-" gorm:"column:root_device_name"`
	RootDeviceType                          string          `json:"rootDeviceType" gorm:"column:root_device_type"`
	SecurityGroups                          json.RawMessage `json:"securityGroups" gorm:"column:security_groups;type:json"`
	SourceDestCheck                         *bool           `json:"sourceDestCheck" gorm:"column:source_dest_check"`
	SpotInstanceRequestID                   string          `json:"-" gorm:"column:spot_instance_request_id"`
	SriovNetSupport                         string          `json:"-" gorm:"column:sriov_net_support"`
	State                                   json.RawMessage `json:"state" gorm:"type:json"`
	StateReason                             json.RawMessage `json:"-" gorm:"column:state_reason;type:json"`
	StateTransitionReason                   string          `json:"-" gorm:"column:state_transition_reason"`
	SubnetID                                string          `json:"subnetId" gorm:"column:subnet_id"`
	TpmSupport                              string          `json:"-" gorm:"column:tpm_support"`
	UsageOperation                          string          `json:"-" gorm:"column:usage_operation"`
	UsageOperationUpdateTime                *time.Time      `json:"-" gorm:"column:usage_operation_update_time"`
	VirtualizationType                      string          `json:"virtualizationType" gorm:"column:virtualization_type"`
	VpcID                                   string          `json:"vpcId" gorm:"column:vpc_id"`
	CreatedAt                               time.Time       `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt                               time.Time       `json:"updatedAt" gorm:"column:updated_at"`
	DeletedAt                               gorm.DeletedAt  `json:"-" gorm:"column:deleted_at;index"`
	Name                                    string          `json:"name"`
	Status                                  string          `json:"status"`
	Location                                string          `json:"location"`
	InstanceTypeAlt                         string          `json:"instanceTypeAlt" gorm:"column:instance_type_alt"`
}

// TableName returns the table name for AWSEC2Instance
func (AWSEC2Instance) TableName() string {
	return "aws_ec2_instances"
}

// AzureVMInstance represents Azure Virtual Machines
type AzureVMInstance struct {
	CqSyncTime        time.Time       `json:"-" gorm:"column:_cq_sync_time"`
	CqSourceName      string          `json:"-" gorm:"column:_cq_source_name"`
	CqID              string          `json:"-" gorm:"column:_cq_id;primarykey"`
	CqParentID        string          `json:"-" gorm:"column:_cq_parent_id"`
	SubscriptionID    string          `json:"subscriptionId" gorm:"column:subscription_id;index"`
	InstanceView      json.RawMessage `json:"-" gorm:"column:instance_view;type:json"`
	Location          string          `json:"location"`
	ExtendedLocation  json.RawMessage `json:"-" gorm:"column:extended_location;type:json"`
	Identity          json.RawMessage `json:"-" gorm:"column:identity;type:json"`
	Plan              json.RawMessage `json:"-" gorm:"column:plan;type:json"`
	Properties        json.RawMessage `json:"-" gorm:"column:properties;type:json"`
	Tags              json.RawMessage `json:"tags" gorm:"type:json"`
	Zones             []string        `json:"-" gorm:"column:zones;type:text[]"`
	ID                string          `json:"id" gorm:"primarykey"`
	Name              string          `json:"name"`
	Resources         json.RawMessage `json:"-" gorm:"column:resources;type:json"`
	Type              string          `json:"-" gorm:"column:type"`
	CreatedAt         time.Time       `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt         time.Time       `json:"updatedAt" gorm:"column:updated_at"`
	DeletedAt         gorm.DeletedAt  `json:"-" gorm:"column:deleted_at;index"`
	Status            string          `json:"status"`
	InstanceTypeAlt   string          `json:"instanceTypeAlt" gorm:"column:instance_type_alt"`
}

// TableName returns the table name for AzureVMInstance
func (AzureVMInstance) TableName() string {
	return "azure_compute_virtual_machines"
}

// GCPComputeInstance represents GCP Compute Engine instances
type GCPComputeInstance struct {
	CqSyncTime                              time.Time       `json:"-" gorm:"column:_cq_sync_time"`
	CqSourceName                            string          `json:"-" gorm:"column:_cq_source_name"`
	CqID                                    string          `json:"-" gorm:"column:_cq_id;primarykey"`
	CqParentID                              string          `json:"-" gorm:"column:_cq_parent_id"`
	ProjectID                               string          `json:"projectId" gorm:"column:project_id;index"`
	AdvancedMachineFeatures                 json.RawMessage `json:"-" gorm:"column:advanced_machine_features;type:json"`
	CanIpForward                            *bool           `json:"-" gorm:"column:can_ip_forward"`
	ConfidentialInstanceConfig              json.RawMessage `json:"-" gorm:"column:confidential_instance_config;type:json"`
	CpuPlatform                             string          `json:"-" gorm:"column:cpu_platform"`
	CreationTimestamp                       string          `json:"-" gorm:"column:creation_timestamp"`
	DeletionProtection                      *bool           `json:"-" gorm:"column:deletion_protection"`
	Description                             string          `json:"-" gorm:"column:description"`
	Disks                                   json.RawMessage `json:"-" gorm:"column:disks;type:json"`
	DisplayDevice                           json.RawMessage `json:"-" gorm:"column:display_device;type:json"`
	Fingerprint                             string          `json:"-" gorm:"column:fingerprint"`
	GuestAccelerators                       json.RawMessage `json:"-" gorm:"column:guest_accelerators;type:json"`
	Hostname                                string          `json:"-" gorm:"column:hostname"`
	ID                                     *int64          `json:"-" gorm:"column:id"`
	InstanceEncryptionKey                   json.RawMessage `json:"-" gorm:"column:instance_encryption_key;type:json"`
	KeyRevocationActionType                 string          `json:"-" gorm:"column:key_revocation_action_type"`
	Kind                                    string          `json:"-" gorm:"column:kind"`
	LabelFingerprint                        string          `json:"-" gorm:"column:label_fingerprint"`
	Labels                                  json.RawMessage `json:"labels" gorm:"type:json"`
	LastStartTimestamp                      string          `json:"-" gorm:"column:last_start_timestamp"`
	LastStopTimestamp                       string          `json:"-" gorm:"column:last_stop_timestamp"`
	LastSuspendedTimestamp                  string          `json:"-" gorm:"column:last_suspended_timestamp"`
	MachineType                             string          `json:"machineType" gorm:"column:machine_type"`
	Metadata                                json.RawMessage `json:"metadata" gorm:"type:json"`
	MinCpuPlatform                          string          `json:"-" gorm:"column:min_cpu_platform"`
	Name                                    string          `json:"name"`
	NetworkInterfaces                       json.RawMessage `json:"-" gorm:"column:network_interfaces;type:json"`
	NetworkPerformanceConfig                json.RawMessage `json:"-" gorm:"column:network_performance_config;type:json"`
	Params                                  json.RawMessage `json:"-" gorm:"column:params;type:json"`
	PrivateIpv6GoogleAccess                 string          `json:"-" gorm:"column:private_ipv6_google_access"`
	ReservationAffinity                     json.RawMessage `json:"-" gorm:"column:reservation_affinity;type:json"`
	ResourcePolicies                        []string        `json:"-" gorm:"column:resource_policies;type:text[]"`
	ResourceStatus                          json.RawMessage `json:"-" gorm:"column:resource_status;type:json"`
	SatisfiesPzs                            *bool           `json:"-" gorm:"column:satisfies_pzs"`
	Scheduling                              json.RawMessage `json:"-" gorm:"column:scheduling;type:json"`
	SelfLink                                string          `json:"selfLink" gorm:"column:self_link;primarykey"`
	ServiceAccounts                         json.RawMessage `json:"-" gorm:"column:service_accounts;type:json"`
	ShieldedInstanceConfig                  json.RawMessage `json:"-" gorm:"column:shielded_instance_config;type:json"`
	ShieldedInstanceIntegrityPolicy         json.RawMessage `json:"-" gorm:"column:shielded_instance_integrity_policy;type:json"`
	SourceMachineImage                      string          `json:"-" gorm:"column:source_machine_image"`
	SourceMachineImageEncryptionKey         json.RawMessage `json:"-" gorm:"column:source_machine_image_encryption_key;type:json"`
	StartRestricted                         *bool           `json:"-" gorm:"column:start_restricted"`
	Status                                  string          `json:"status"`
	StatusMessage                           string          `json:"-" gorm:"column:status_message"`
	Tags                                    json.RawMessage `json:"tags" gorm:"type:json"`
	Zone                                    string          `json:"zone"`
	CreatedAt                               time.Time       `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt                               time.Time       `json:"updatedAt" gorm:"column:updated_at"`
	DeletedAt                               gorm.DeletedAt  `json:"-" gorm:"column:deleted_at;index"`
	InstanceTypeAlt                         string          `json:"instanceTypeAlt" gorm:"column:instance_type_alt"`
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
		"eq":           true,
		"ne":           true,
		"lt":           true,
		"lte":          true,
		"gt":           true,
		"gte":          true,
		"in":           true,
		"not_in":       true,
		"contains":     true,
		"starts_with":  true,
		"ends_with":    true,
		"like":         true,
		"ilike":        true,
		"between":      true,
		"is_null":      true,
		"is_not_null":  true,
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
	case "in", "not_in":
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
	case "is_null", "is_not_null":
		// These operators don't need a value
		// Value can be nil or any value, it will be ignored
	case "contains", "starts_with", "ends_with", "like", "ilike":
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