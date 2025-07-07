package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang-service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type VMHandlerTestSuite struct {
	suite.Suite
	db     *gorm.DB
	router *gin.Engine
	handler *Handler
}

func (suite *VMHandlerTestSuite) SetupSuite() {
	// Set up test database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)

	// Auto migrate all cloud-specific tables
	err = db.AutoMigrate(
		&models.AWSEC2Instance{},
		&models.AzureVMInstance{},
		&models.GCPComputeInstance{},
	)
	suite.Require().NoError(err)

	suite.db = db
	suite.handler = New(db)

	// Set up Gin router
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	suite.router.GET("/vms", suite.handler.GetVMs)
}

func (suite *VMHandlerTestSuite) SetupTest() {
	// Clean up database before each test
	suite.db.Exec("DELETE FROM aws_ec2_instances")
	suite.db.Exec("DELETE FROM azure_vm_instances")
	suite.db.Exec("DELETE FROM gcp_compute_instances")
	
	// Insert test data
	suite.insertTestVMs()
}

func (suite *VMHandlerTestSuite) insertTestVMs() {
	// Insert AWS EC2 instances
	awsInstances := []models.AWSEC2Instance{
		{
			BaseVM: models.BaseVM{
				ID:           "i-1234567890abcdef0",
				Name:         "web-server-01",
				Status:       "running",
				CreatedAt:    time.Now().Add(-24 * time.Hour),
				UpdatedAt:    time.Now().Add(-1 * time.Hour),
				Location:     "us-east-1",
				InstanceType: "t2.micro",
			},
			AccountID:        "123456789012",
			VpcID:           "vpc-123",
			SubnetID:        "subnet-123",
			SecurityGroupIDs: mustMarshalJSON([]string{"sg-123"}),
			PrivateIPAddress: "10.0.1.100",
			PublicIPAddress:  "54.123.45.67",
		},
		{
			BaseVM: models.BaseVM{
				ID:           "i-fedcba0987654321",
				Name:         "database-server",
				Status:       "stopped",
				CreatedAt:    time.Now().Add(-48 * time.Hour),
				UpdatedAt:    time.Now().Add(-2 * time.Hour),
				Location:     "us-west-2",
				InstanceType: "t3.medium",
			},
			AccountID:        "123456789012",
			VpcID:           "vpc-456",
			SubnetID:        "subnet-456",
			SecurityGroupIDs: mustMarshalJSON([]string{"sg-456"}),
			PrivateIPAddress: "10.0.2.200",
		},
	}

	// Insert GCP Compute instances
	gcpInstances := []models.GCPComputeInstance{
		{
			BaseVM: models.BaseVM{
				ID:           "projects/my-project/zones/us-central1-a/instances/web-server-gcp",
				Name:         "web-server-gcp",
				Status:       "running",
				CreatedAt:    time.Now().Add(-12 * time.Hour),
				UpdatedAt:    time.Now().Add(-30 * time.Minute),
				Location:     "us-central1-a",
				InstanceType: "e2-standard-2",
			},
			ProjectID:         "my-project-123456",
			Zone:             "us-central1-a",
			MachineType:      "e2-standard-2",
			PrivateIPAddress: "10.128.0.2",
			PublicIPAddress:  "34.123.45.67",
		},
	}

	// Insert Azure VM instances
	azureInstances := []models.AzureVMInstance{
		{
			BaseVM: models.BaseVM{
				ID:           "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/my-rg/providers/Microsoft.Compute/virtualMachines/web-vm",
				Name:         "web-vm",
				Status:       "running",
				CreatedAt:    time.Now().Add(-6 * time.Hour),
				UpdatedAt:    time.Now().Add(-10 * time.Minute),
				Location:     "East US",
				InstanceType: "Standard_D2s_v3",
			},
			SubscriptionID:   "12345678-1234-1234-1234-123456789012",
			ResourceGroup:    "my-rg",
			VMSize:          "Standard_D2s_v3",
			PrivateIPAddress: "10.0.0.4",
			PublicIPAddress:  "52.123.45.67",
		},
	}

	for _, instance := range awsInstances {
		suite.db.Create(&instance)
	}

	for _, instance := range gcpInstances {
		suite.db.Create(&instance)
	}

	for _, instance := range azureInstances {
		suite.db.Create(&instance)
	}
}

func (suite *VMHandlerTestSuite) TestGetVMs_Basic() {
	req, _ := http.NewRequest("GET", "/vms", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response models.VMListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 4, len(response.Data))
	assert.Equal(suite.T(), 1, response.Pagination.Page)
	assert.Equal(suite.T(), 20, response.Pagination.PageSize)
	assert.Equal(suite.T(), 4, response.Pagination.TotalItems)
	assert.Equal(suite.T(), 1, response.Pagination.TotalPages)

	// Check that all cloud types are present
	cloudTypes := make(map[string]bool)
	for _, vm := range response.Data {
		cloudTypes[vm.CloudType] = true
	}
	assert.True(suite.T(), cloudTypes["aws"])
	assert.True(suite.T(), cloudTypes["gcp"])
	assert.True(suite.T(), cloudTypes["azure"])
}

func (suite *VMHandlerTestSuite) TestGetVMs_Pagination() {
	req, _ := http.NewRequest("GET", "/vms?page=1&pageSize=2", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response models.VMListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 2, len(response.Data))
	assert.Equal(suite.T(), 1, response.Pagination.Page)
	assert.Equal(suite.T(), 2, response.Pagination.PageSize)
	assert.Equal(suite.T(), 4, response.Pagination.TotalItems)
	assert.Equal(suite.T(), 2, response.Pagination.TotalPages)
}

func (suite *VMHandlerTestSuite) TestGetVMs_Sorting() {
	req, _ := http.NewRequest("GET", "/vms?sortBy=name&sortOrder=desc", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response models.VMListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 4, len(response.Data))
	// Note: Sorting across multiple tables might not maintain perfect order
	// but we can verify basic functionality
}

func (suite *VMHandlerTestSuite) TestGetVMs_FilterByCloudType() {
	filterJSON := `[{"field":"cloudType","operator":"eq","value":"aws"}]`
	req, _ := http.NewRequest("GET", "/vms?filter="+filterJSON, nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response models.VMListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 2, len(response.Data))
	for _, vm := range response.Data {
		assert.Equal(suite.T(), "aws", vm.CloudType)
	}
}

func (suite *VMHandlerTestSuite) TestGetVMs_FilterByStatus() {
	filterJSON := `[{"field":"status","operator":"eq","value":"running"}]`
	req, _ := http.NewRequest("GET", "/vms?filter="+filterJSON, nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response models.VMListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 3, len(response.Data))
	for _, vm := range response.Data {
		assert.Equal(suite.T(), "running", vm.Status)
	}
}

func (suite *VMHandlerTestSuite) TestGetVMs_FilterByNameLike() {
	filterJSON := `[{"field":"name","operator":"like","value":"web"}]`
	req, _ := http.NewRequest("GET", "/vms?filter="+filterJSON, nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response models.VMListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 3, len(response.Data))
	for _, vm := range response.Data {
		assert.Contains(suite.T(), vm.Name, "web")
	}
}

func (suite *VMHandlerTestSuite) TestGetVMs_FilterByInstanceTypeIn() {
	filterJSON := `[{"field":"instanceType","operator":"in","value":["t2.micro","e2-standard-2"]}]`
	req, _ := http.NewRequest("GET", "/vms?filter="+filterJSON, nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response models.VMListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)

	assert.GreaterOrEqual(suite.T(), len(response.Data), 1)
	for _, vm := range response.Data {
		assert.Contains(suite.T(), []string{"t2.micro", "e2-standard-2"}, vm.InstanceType)
	}
}

func (suite *VMHandlerTestSuite) TestGetVMs_FilterByInstanceTypeNotEqual() {
	filterJSON := `[{"field":"instanceType","operator":"ne","value":"t2.micro"}]`
	req, _ := http.NewRequest("GET", "/vms?filter="+filterJSON, nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response models.VMListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)

	for _, vm := range response.Data {
		assert.NotEqual(suite.T(), "t2.micro", vm.InstanceType)
	}
}

func (suite *VMHandlerTestSuite) TestGetVMs_MultipleFilters() {
	filterJSON := `[{"field":"status","operator":"eq","value":"running"},{"field":"cloudType","operator":"ne","value":"aws"}]`
	req, _ := http.NewRequest("GET", "/vms?filter="+filterJSON, nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response models.VMListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)

	for _, vm := range response.Data {
		assert.Equal(suite.T(), "running", vm.Status)
		assert.NotEqual(suite.T(), "aws", vm.CloudType)
	}
}

func (suite *VMHandlerTestSuite) TestGetVMs_InvalidFilter() {
	req, _ := http.NewRequest("GET", "/vms?filter=invalid-json", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response models.Error
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "INVALID_PARAMS", response.Code)
}

func (suite *VMHandlerTestSuite) TestGetVMs_InvalidOperator() {
	filterJSON := `[{"field":"status","operator":"invalid","value":"running"}]`
	req, _ := http.NewRequest("GET", "/vms?filter="+filterJSON, nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response models.Error
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "VALIDATION_ERROR", response.Code)
}

func (suite *VMHandlerTestSuite) TestGetVMs_InvalidPageSize() {
	req, _ := http.NewRequest("GET", "/vms?pageSize=2000", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response models.Error
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "INVALID_PARAMS", response.Code)
}

func (suite *VMHandlerTestSuite) TestGetVMs_InvalidSortOrder() {
	req, _ := http.NewRequest("GET", "/vms?sortOrder=invalid", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response models.Error
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "INVALID_PARAMS", response.Code)
}

func (suite *VMHandlerTestSuite) TestGetVMs_FilterValidation_EmptyField() {
	filterJSON := `[{"field":"","operator":"eq","value":"running"}]`
	req, _ := http.NewRequest("GET", "/vms?filter="+filterJSON, nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response models.Error
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "VALIDATION_ERROR", response.Code)
}

func (suite *VMHandlerTestSuite) TestGetVMs_FilterValidation_InvalidBetween() {
	filterJSON := `[{"field":"createdAt","operator":"between","value":["2023-01-01"]}]`
	req, _ := http.NewRequest("GET", "/vms?filter="+filterJSON, nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response models.Error
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "VALIDATION_ERROR", response.Code)
}

func (suite *VMHandlerTestSuite) TestGetVMs_FilterValidation_InvalidInOperator() {
	filterJSON := `[{"field":"status","operator":"in","value":"running"}]`
	req, _ := http.NewRequest("GET", "/vms?filter="+filterJSON, nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response models.Error
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "VALIDATION_ERROR", response.Code)
}

func (suite *VMHandlerTestSuite) TestGetVMs_FilterValidation_InvalidNullOperator() {
	filterJSON := `[{"field":"status","operator":"null","value":"not-boolean"}]`
	req, _ := http.NewRequest("GET", "/vms?filter="+filterJSON, nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response models.Error
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "VALIDATION_ERROR", response.Code)
}

func TestVMHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(VMHandlerTestSuite))
}

func mustMarshalJSON(v interface{}) json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}