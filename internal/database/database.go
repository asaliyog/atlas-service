package database

import (
	"golang-service/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Initialize initializes the database connection and runs migrations
func Initialize(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate the schema
	err = db.AutoMigrate(
		&models.User{},
		&models.AWSEC2Instance{},
		&models.AzureVMInstance{},
		&models.GCPComputeInstance{},
		// Add other models here as you create them
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Health checks database connectivity
func Health(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}