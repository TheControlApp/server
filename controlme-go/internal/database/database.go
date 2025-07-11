package database

import (
	"fmt"

	"github.com/thecontrolapp/controlme-go/internal/config"
	"github.com/thecontrolapp/controlme-go/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Initialize sets up the database connection and runs migrations
func Initialize(cfg *config.Config) (*gorm.DB, error) {
	// Build connection string
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC",
		cfg.Database.Host,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)

	// Set log level based on environment
	logLevel := logger.Info
	if cfg.Environment == "production" {
		logLevel = logger.Error
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run auto-migration
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

// runMigrations runs the database migrations
func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Command{},
		&models.ControlAppCmd{},
		&models.ChatLog{},
		&models.Group{},
		&models.GroupMember{},
		&models.Relationship{},
		&models.Block{},
		&models.Report{},
		&models.Invite{},
	)
}
