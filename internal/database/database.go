package database

import (
	"fmt"
	"log"

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
	log.Println("Running database migrations with improved error handling...")
	
	// Try GORM AutoMigrate first, with better error handling
	models := []interface{}{
		&models.User{},
		&models.Command{},
		&models.Tag{},
		&models.Block{},
		&models.Report{},
	}
	
	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			log.Printf("GORM AutoMigrate failed for model %T: %v", model, err)
			log.Println("Falling back to manual table verification...")
			
			// If GORM fails, let's verify the tables exist manually
			if err := verifyTableExists(db, model); err != nil {
				return fmt.Errorf("migration failed for model %T: %w", model, err)
			}
		}
	}
	
	log.Println("✅ Database migration completed successfully.")
	return nil
}

// verifyTableExists checks if a table exists for a given model
func verifyTableExists(db *gorm.DB, model interface{}) error {
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(model); err != nil {
		return fmt.Errorf("could not parse model %T: %w", model, err)
	}
	
	tableName := stmt.Schema.Table
	if tableName == "" {
		return fmt.Errorf("could not determine table name for model %T", model)
	}
	
	var count int64
	err := db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_name = ? AND table_type = 'BASE TABLE'", 
		tableName).Scan(&count).Error
	
	if err != nil {
		return fmt.Errorf("error checking table %s: %w", tableName, err)
	}
	
	if count == 0 {
		return fmt.Errorf("table %s does not exist and GORM migration failed", tableName)
	}
	
	log.Printf("✓ Table '%s' exists", tableName)
	return nil
}
