package main

import (
	"encoding/json"
	"log"
	"time"

	"golang-service/internal/config"
	"golang-service/internal/database"
	"golang-service/internal/models"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Create sample VMs
	sampleVMs := []models.VM{
		// AWS VMs
		{
			ID:             "i-1234567890abcdef0",
			Name:           "web-server-01",
			CloudType:      "aws",
			Status:         "running",
			CreatedAt:      time.Now().Add(-24 * time.Hour),
			CloudAccountID: "123456789012",
			Location:       "us-east-1",
			InstanceType:   "t2.micro",
			CloudSpecificDetails: mustMarshalJSON(models.AWSDetails{
				CloudType:        "aws",
				VpcID:           "vpc-12345678",
				SubnetID:        "subnet-12345678",
				SecurityGroupIDs: []string{"sg-12345678", "sg-98765432"},
				PrivateIPAddress: "10.0.1.100",
				PublicIPAddress:  "54.123.45.67",
			}),
		},
		{
			ID:             "i-fedcba0987654321",
			Name:           "database-server",
			CloudType:      "aws",
			Status:         "running",
			CreatedAt:      time.Now().Add(-48 * time.Hour),
			CloudAccountID: "123456789012",
			Location:       "us-west-2",
			InstanceType:   "t3.medium",
			CloudSpecificDetails: mustMarshalJSON(models.AWSDetails{
				CloudType:        "aws",
				VpcID:           "vpc-87654321",
				SubnetID:        "subnet-87654321",
				SecurityGroupIDs: []string{"sg-87654321"},
				PrivateIPAddress: "10.0.2.200",
			}),
		},
		// GCP VMs
		{
			ID:             "projects/my-project/zones/us-central1-a/instances/web-server-gcp",
			Name:           "web-server-gcp",
			CloudType:      "gcp",
			Status:         "running",
			CreatedAt:      time.Now().Add(-12 * time.Hour),
			CloudAccountID: "my-project-123456",
			Location:       "us-central1-a",
			InstanceType:   "e2-standard-2",
			CloudSpecificDetails: mustMarshalJSON(models.GCPDetails{
				CloudType:        "gcp",
				MachineType:      "e2-standard-2",
				Network:          "default",
				Region:           "us-central1",
				PrivateIPAddress: "10.128.0.2",
				PublicIPAddress:  "34.123.45.67",
			}),
		},
		{
			ID:             "projects/my-project/zones/europe-west1-b/instances/api-server",
			Name:           "api-server",
			CloudType:      "gcp",
			Status:         "stopped",
			CreatedAt:      time.Now().Add(-72 * time.Hour),
			CloudAccountID: "my-project-123456",
			Location:       "europe-west1-b",
			InstanceType:   "n1-standard-1",
			CloudSpecificDetails: mustMarshalJSON(models.GCPDetails{
				CloudType:        "gcp",
				MachineType:      "n1-standard-1",
				Network:          "custom-network",
				Region:           "europe-west1",
				PrivateIPAddress: "10.132.0.5",
			}),
		},
		// Azure VMs
		{
			ID:             "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/my-rg/providers/Microsoft.Compute/virtualMachines/web-vm",
			Name:           "web-vm",
			CloudType:      "azure",
			Status:         "running",
			CreatedAt:      time.Now().Add(-6 * time.Hour),
			CloudAccountID: "12345678-1234-1234-1234-123456789012",
			Location:       "East US",
			InstanceType:   "Standard_D2s_v3",
			CloudSpecificDetails: mustMarshalJSON(models.AzureDetails{
				CloudType:        "azure",
				ResourceGroup:    "my-rg",
				VMSize:           "Standard_D2s_v3",
				PrivateIPAddress: "10.0.0.4",
				PublicIPAddress:  "52.123.45.67",
			}),
		},
		{
			ID:             "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/production/providers/Microsoft.Compute/virtualMachines/prod-app",
			Name:           "prod-app",
			CloudType:      "azure",
			Status:         "running",
			CreatedAt:      time.Now().Add(-96 * time.Hour),
			CloudAccountID: "12345678-1234-1234-1234-123456789012",
			Location:       "West Europe",
			InstanceType:   "Standard_B2s",
			CloudSpecificDetails: mustMarshalJSON(models.AzureDetails{
				CloudType:        "azure",
				ResourceGroup:    "production",
				VMSize:           "Standard_B2s",
				PrivateIPAddress: "10.1.0.10",
			}),
		},
	}

	// Insert sample VMs
	for _, vm := range sampleVMs {
		result := db.Create(&vm)
		if result.Error != nil {
			log.Printf("Failed to create VM %s: %v", vm.ID, result.Error)
		} else {
			log.Printf("Created VM: %s (%s)", vm.Name, vm.ID)
		}
	}

	log.Println("Sample VM data seeded successfully!")
}

func mustMarshalJSON(v interface{}) json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}
	return data
}