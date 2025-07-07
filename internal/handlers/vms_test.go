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

	// Auto migrate
	err = db.AutoMigrate(&models.VM{})
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
	suite.db.Exec("DELETE FROM vms")
	
	// Insert test data
	suite.insertTestVMs()
}

func (suite *VMHandlerTestSuite) insertTestVMs() {
	awsDetails, _ := json.Marshal(models.AWSDetails{
		CloudType: "aws",
		VpcID:     "vpc-123",
		SubnetID:  "subnet-123",
		SecurityGroupIDs: []string{"sg-123"},
		PrivateIPAddress: "10.0.1.100",
		PublicIPAddress:  "54.123.45.67",
	})

	gcpDetails, _ := json.Marshal(models.GCPDetails{
		CloudType:   "gcp",
		MachineType: "e2-standard-2",
		Network:     "default",
		Region:      "us-central1",
		PrivateIPAddress: "10.128.0.2",
		PublicIPAddress:  "34.123.45.67",
	})

	azureDetails, _ := json.Marshal(models.AzureDetails{
		CloudType:     "azure",
		ResourceGroup: "my-rg",
		VMSize:        "Standard_D2s_v3",
		PrivateIPAddress: "10.0.0.4",
		PublicIPAddress:  "52.123.45.67",
	})

	vms := []models.VM{
		{
			ID:             "i-1234567890abcdef0",
			Name:           "web-server-01",
			CloudType:      "aws",
			Status:         "running",
			CreatedAt:      time.Now().Add(-24 * time.Hour),
			CloudAccountID: "123456789012",
			Location:       "us-east-1",
			InstanceType:   "t2.micro",
			CloudSpecificDetails: awsDetails,
		},
		{
			ID:             "i-fedcba0987654321",
			Name:           "database-server",
			CloudType:      "aws",
			Status:         "stopped",
			CreatedAt:      time.Now().Add(-48 * time.Hour),
			CloudAccountID: "123456789012",
			Location:       "us-west-2",
			InstanceType:   "t3.medium",
			CloudSpecificDetails: awsDetails,
		},
		{
			ID:             "gcp-web-server",
			Name:           "web-server-gcp",
			CloudType:      "gcp",
			Status:         "running",
			CreatedAt:      time.Now().Add(-12 * time.Hour),
			CloudAccountID: "my-project-123456",
			Location:       "us-central1-a",
			InstanceType:   "e2-standard-2",
			CloudSpecificDetails: gcpDetails,
		},
		{
			ID:             "azure-web-vm",
			Name:           "web-vm",
			CloudType:      "azure",
			Status:         "running",
			CreatedAt:      time.Now().Add(-6 * time.Hour),
			CloudAccountID: "12345678-1234-1234-1234-123456789012",
			Location:       "East US",
			InstanceType:   "Standard_D2s_v3",
			CloudSpecificDetails: azureDetails,
		},
	}

	for _, vm := range vms {
		suite.db.Create(&vm)
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
	// Should be sorted by name in descending order
	assert.Equal(suite.T(), "web-vm", response.Data[0].Name)
	assert.Equal(suite.T(), "web-server-gcp", response.Data[1].Name)
	assert.Equal(suite.T(), "web-server-01", response.Data[2].Name)
	assert.Equal(suite.T(), "database-server", response.Data[3].Name)
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

func (suite *VMHandlerTestSuite) TestGetVMs_FilterByNameContains() {
	filterJSON := `[{"field":"name","operator":"contains","value":"web"}]`
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

func (suite *VMHandlerTestSuite) TestGetVMs_MultipleFilters() {
	filterJSON := `[{"field":"cloudType","operator":"eq","value":"aws"},{"field":"status","operator":"eq","value":"running"}]`
	req, _ := http.NewRequest("GET", "/vms?filter="+filterJSON, nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response models.VMListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 1, len(response.Data))
	if len(response.Data) > 0 {
		assert.Equal(suite.T(), "aws", response.Data[0].CloudType)
		assert.Equal(suite.T(), "running", response.Data[0].Status)
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
	assert.Contains(suite.T(), response.Message, "invalid filter JSON")
}

func (suite *VMHandlerTestSuite) TestGetVMs_InvalidOperator() {
	filterJSON := `[{"field":"cloudType","operator":"invalid","value":"aws"}]`
	req, _ := http.NewRequest("GET", "/vms?filter="+filterJSON, nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response models.Error
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Contains(suite.T(), response.Message, "unsupported operator")
}

func (suite *VMHandlerTestSuite) TestGetVMs_ParameterLimits() {
	// Test page size limit
	req, _ := http.NewRequest("GET", "/vms?pageSize=200", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response models.VMListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)

	// Should be limited to 100
	assert.Equal(suite.T(), 100, response.Pagination.PageSize)
}

func TestVMHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(VMHandlerTestSuite))
}