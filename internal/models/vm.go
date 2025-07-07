package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// VM represents a virtual machine across cloud providers
type VM struct {
	ID                   string                 `json:"id" gorm:"primarykey"`
	Name                 string                 `json:"name"`
	CloudType            string                 `json:"cloudType" gorm:"column:cloud_type;index"`
	Status               string                 `json:"status"`
	CreatedAt            time.Time              `json:"createdAt"`
	UpdatedAt            time.Time              `json:"updatedAt"`
	DeletedAt            gorm.DeletedAt         `json:"-" gorm:"index"`
	CloudAccountID       string                 `json:"cloudAccountId" gorm:"column:cloud_account_id;index"`
	Location             string                 `json:"location" gorm:"index"`
	InstanceType         string                 `json:"instanceType" gorm:"column:instance_type"`
	CloudSpecificDetails json.RawMessage        `json:"cloudSpecificDetails" gorm:"column:cloud_specific_details;type:json"`
}

// AWSDetails represents AWS-specific VM details
type AWSDetails struct {
	CloudType        string   `json:"cloudType"`
	VpcID            string   `json:"vpcId"`
	SubnetID         string   `json:"subnetId"`
	SecurityGroupIDs []string `json:"securityGroupIds"`
	PrivateIPAddress string   `json:"privateIpAddress,omitempty"`
	PublicIPAddress  string   `json:"publicIpAddress,omitempty"`
}

// GCPDetails represents GCP-specific VM details
type GCPDetails struct {
	CloudType        string `json:"cloudType"`
	MachineType      string `json:"machineType"`
	Network          string `json:"network"`
	Region           string `json:"region"`
	PrivateIPAddress string `json:"privateIpAddress,omitempty"`
	PublicIPAddress  string `json:"publicIpAddress,omitempty"`
}

// AzureDetails represents Azure-specific VM details
type AzureDetails struct {
	CloudType        string `json:"cloudType"`
	ResourceGroup    string `json:"resourceGroup"`
	VMSize           string `json:"vmSize"`
	PrivateIPAddress string `json:"privateIpAddress,omitempty"`
	PublicIPAddress  string `json:"publicIpAddress,omitempty"`
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

// VMQueryParams represents the query parameters for VM endpoints
type VMQueryParams struct {
	Page      int        `json:"page"`
	PageSize  int        `json:"pageSize"`
	SortBy    string     `json:"sortBy"`
	SortOrder string     `json:"sortOrder"`
	Filters   []VMFilter `json:"filters"`
}

// Error represents an error response
type Error struct {
	Message string `json:"message"`
}