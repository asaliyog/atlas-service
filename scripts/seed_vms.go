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

	// Create sample AWS EC2 instances
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
			AccountID:                "123456789012",
			VpcID:                    "vpc-12345678",
			SubnetID:                 "subnet-12345678",
			SecurityGroupIDs:         mustMarshalJSON([]string{"sg-12345678", "sg-98765432"}),
			PrivateIPAddress:         "10.0.1.100",
			PublicIPAddress:          "54.123.45.67",
			KeyName:                  "my-key-pair",
			ImageID:                  "ami-0abcdef1234567890",
			LaunchTime:               time.Now().Add(-24 * time.Hour),
			AvailabilityZone:         "us-east-1a",
			PublicDnsName:            "ec2-54-123-45-67.compute-1.amazonaws.com",
			PrivateDnsName:           "ip-10-0-1-100.ec2.internal",
			Architecture:             "x86_64",
			VirtualizationType:       "hvm",
			Platform:                 "linux",
			RootDeviceType:           "ebs",
			MonitoringState:          "enabled",
			EbsOptimized:             true,
			Tags:                     mustMarshalJSON(map[string]string{"Environment": "production", "Owner": "team-web"}),
		},
		{
			BaseVM: models.BaseVM{
				ID:           "i-fedcba0987654321",
				Name:         "database-server",
				Status:       "running",
				CreatedAt:    time.Now().Add(-48 * time.Hour),
				UpdatedAt:    time.Now().Add(-2 * time.Hour),
				Location:     "us-west-2",
				InstanceType: "t3.medium",
			},
			AccountID:                "123456789012",
			VpcID:                    "vpc-87654321",
			SubnetID:                 "subnet-87654321",
			SecurityGroupIDs:         mustMarshalJSON([]string{"sg-87654321"}),
			PrivateIPAddress:         "10.0.2.200",
			KeyName:                  "db-key-pair",
			ImageID:                  "ami-0fedcba0987654321",
			LaunchTime:               time.Now().Add(-48 * time.Hour),
			AvailabilityZone:         "us-west-2b",
			PrivateDnsName:           "ip-10-0-2-200.us-west-2.compute.internal",
			Architecture:             "x86_64",
			VirtualizationType:       "hvm",
			Platform:                 "linux",
			RootDeviceType:           "ebs",
			MonitoringState:          "disabled",
			EbsOptimized:             false,
			Tags:                     mustMarshalJSON(map[string]string{"Environment": "production", "Owner": "team-db"}),
		},
	}

	// Create sample GCP Compute instances
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
			ProjectID:                "my-project-123456",
			Zone:                     "us-central1-a",
			MachineType:              "e2-standard-2",
			PrivateIPAddress:         "10.128.0.2",
			PublicIPAddress:          "34.123.45.67",
			NetworkInterfaces:        mustMarshalJSON([]string{"default"}),
			CpuPlatform:              "Intel Broadwell",
			LastStartTimestamp:       time.Now().Add(-12 * time.Hour),
			Hostname:                 "web-server-gcp.c.my-project-123456.internal",
			Tags:                     mustMarshalJSON(map[string]string{"http-server": "true", "https-server": "true"}),
			Labels:                   mustMarshalJSON(map[string]string{"env": "production", "team": "web"}),
		},
		{
			BaseVM: models.BaseVM{
				ID:           "projects/my-project/zones/europe-west1-b/instances/api-server",
				Name:         "api-server",
				Status:       "stopped",
				CreatedAt:    time.Now().Add(-72 * time.Hour),
				UpdatedAt:    time.Now().Add(-6 * time.Hour),
				Location:     "europe-west1-b",
				InstanceType: "n1-standard-1",
			},
			ProjectID:                "my-project-123456",
			Zone:                     "europe-west1-b",
			MachineType:              "n1-standard-1",
			PrivateIPAddress:         "10.132.0.5",
			NetworkInterfaces:        mustMarshalJSON([]string{"custom-network"}),
			CpuPlatform:              "Intel Skylake",
			LastStopTimestamp:        time.Now().Add(-6 * time.Hour),
			Hostname:                 "api-server.c.my-project-123456.internal",
			Tags:                     mustMarshalJSON(map[string]string{"api-server": "true"}),
			Labels:                   mustMarshalJSON(map[string]string{"env": "development", "team": "api"}),
		},
	}

	// Create sample Azure VM instances
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
			SubscriptionID:           "12345678-1234-1234-1234-123456789012",
			ResourceGroup:            "my-rg",
			VMSize:                   "Standard_D2s_v3",
			PrivateIPAddress:         "10.0.0.4",
			PublicIPAddress:          "52.123.45.67",
			OSType:                   "Linux",
			VMId:                     "12345678-9012-3456-7890-123456789012",
			Tags:                     mustMarshalJSON(map[string]string{"Environment": "production", "Owner": "team-web"}),
		},
		{
			BaseVM: models.BaseVM{
				ID:           "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/production/providers/Microsoft.Compute/virtualMachines/prod-app",
				Name:         "prod-app",
				Status:       "running",
				CreatedAt:    time.Now().Add(-96 * time.Hour),
				UpdatedAt:    time.Now().Add(-1 * time.Hour),
				Location:     "West Europe",
				InstanceType: "Standard_B2s",
			},
			SubscriptionID:           "12345678-1234-1234-1234-123456789012",
			ResourceGroup:            "production",
			VMSize:                   "Standard_B2s",
			PrivateIPAddress:         "10.1.0.10",
			OSType:                   "Linux",
			VMId:                     "87654321-2109-6543-0987-210987654321",
			Tags:                     mustMarshalJSON(map[string]string{"Environment": "production", "Owner": "team-app"}),
		},
	}

	// Insert AWS EC2 instances
	for _, instance := range awsInstances {
		result := db.Create(&instance)
		if result.Error != nil {
			log.Printf("Failed to create AWS EC2 instance %s: %v", instance.ID, result.Error)
		} else {
			log.Printf("Created AWS EC2 instance: %s (%s)", instance.Name, instance.ID)
		}
	}

	// Insert GCP Compute instances
	for _, instance := range gcpInstances {
		result := db.Create(&instance)
		if result.Error != nil {
			log.Printf("Failed to create GCP Compute instance %s: %v", instance.ID, result.Error)
		} else {
			log.Printf("Created GCP Compute instance: %s (%s)", instance.Name, instance.ID)
		}
	}

	// Insert Azure VM instances
	for _, instance := range azureInstances {
		result := db.Create(&instance)
		if result.Error != nil {
			log.Printf("Failed to create Azure VM instance %s: %v", instance.ID, result.Error)
		} else {
			log.Printf("Created Azure VM instance: %s (%s)", instance.Name, instance.ID)
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