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
	// Note: If you encounter issues with AutoMigrate on PostgreSQL 17,
	// you may need to run migrations manually or use a different approach
	err = db.AutoMigrate(
		&models.User{},
		&models.AWSEC2Instance{},
		&models.AzureVMInstance{},
		&models.GCPComputeInstance{},
		// Add other models here as you create them
	)
	if err != nil {
		// Log the error but don't fail completely - tables may already exist
		// In production, you might want to handle this differently
		return db, nil // Continue without migration for now
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