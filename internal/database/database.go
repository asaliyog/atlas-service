package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Initialize initializes the database connection
func Initialize(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// For now, skip AutoMigrate since the tables are created by init.sql
	// The tables should already exist from the init.sql script
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