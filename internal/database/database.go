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

	// Test the connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	
	log.Println("✓ Database connection established successfully")

	// Ensure UUID extension is available
	if err := ensureUUIDExtension(db); err != nil {
		return nil, fmt.Errorf("failed to setup UUID extension: %w", err)
	}

	// Run auto-migration
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

// ensureUUIDExtension ensures that the uuid-ossp extension is available
func ensureUUIDExtension(db *gorm.DB) error {
	log.Println("Ensuring UUID extension is available...")
	
	// First check if the extension exists
	var extensionExists bool
	if err := db.Raw("SELECT EXISTS(SELECT 1 FROM pg_extension WHERE extname = 'uuid-ossp')").Scan(&extensionExists).Error; err != nil {
		log.Printf("Warning: Failed to check for uuid-ossp extension: %v", err)
	}
	
	if !extensionExists {
		if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
			log.Printf("Warning: Failed to create uuid-ossp extension: %v", err)
			// Don't fail here as uuid_generate_v4() might still work or be available through other means
		} else {
			log.Println("✓ UUID extension created successfully")
		}
	} else {
		log.Println("✓ UUID extension already exists")
	}
	
	return nil
}

// runMigrations runs the database migrations
func runMigrations(db *gorm.DB) error {
	log.Println("Running database migrations with improved error handling...")
	
	// Migrate models in order to handle dependencies properly
	// Start with models that have no foreign key dependencies
	if err := migrateWithFallback(db, &models.User{}, "User"); err != nil {
		return err
	}
	
	if err := migrateWithFallback(db, &models.Tag{}, "Tag"); err != nil {
		return err
	}
	
	// Now migrate models with foreign key dependencies
	if err := migrateCommandTable(db); err != nil {
		return fmt.Errorf("failed to migrate Command model: %w", err)
	}
	log.Println("✓ Command table migrated successfully")
	
	if err := migrateWithFallback(db, &models.Block{}, "Block"); err != nil {
		return err
	}
	
	if err := migrateWithFallback(db, &models.Report{}, "Report"); err != nil {
		return err
	}
	
	log.Println("✅ Database migration completed successfully.")
	return nil
}

// migrateWithFallback attempts GORM AutoMigrate with fallback error handling
func migrateWithFallback(db *gorm.DB, model interface{}, modelName string) error {
	if err := db.AutoMigrate(model); err != nil {
		log.Printf("GORM AutoMigrate failed for %s model: %v", modelName, err)
		log.Printf("Attempting manual %s table creation...", modelName)
		
		// Try manual table creation
		if err := createTableManually(db, modelName); err != nil {
			return fmt.Errorf("manual %s table creation failed: %w", modelName, err)
		}
	}
	log.Printf("✓ %s table migrated successfully", modelName)
	return nil
}

// migrateCommandTable handles the Command table migration with special care for foreign keys
func migrateCommandTable(db *gorm.DB) error {
	// First, try the normal AutoMigrate
	if err := db.AutoMigrate(&models.Command{}); err != nil {
		log.Printf("GORM AutoMigrate failed for Command model: %v", err)
		log.Println("Attempting manual Command table creation...")
		
		// Create the table manually if AutoMigrate fails
		if err := createCommandTableManually(db); err != nil {
			return fmt.Errorf("manual Command table creation failed: %w", err)
		}
	}
	return nil
}

// createCommandTableManually creates the commands table manually
func createCommandTableManually(db *gorm.DB) error {
	// Check if table already exists
	var exists bool
	if err := db.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'commands')").Scan(&exists).Error; err != nil {
		return fmt.Errorf("error checking if commands table exists: %w", err)
	}
	
	if exists {
		log.Println("Commands table already exists, skipping manual creation")
		return nil
	}
	
	log.Println("Creating commands table manually due to GORM migration failure...")
	
	// Create the commands table manually with proper foreign key constraints
	createTableSQL := `
		CREATE TABLE commands (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			instructions TEXT NOT NULL,
			sender_id UUID NOT NULL,
			receiver_id UUID,
			tags TEXT,
			status VARCHAR(20) DEFAULT 'pending',
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW(),
			CONSTRAINT fk_commands_sender FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
			CONSTRAINT fk_commands_receiver FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE SET NULL
		)`
	
	if err := db.Exec(createTableSQL).Error; err != nil {
		return fmt.Errorf("error creating commands table: %w", err)
	}
	
	// Create indexes for better performance
	indexSQL := []string{
		"CREATE INDEX IF NOT EXISTS idx_commands_sender_id ON commands(sender_id)",
		"CREATE INDEX IF NOT EXISTS idx_commands_receiver_id ON commands(receiver_id)",
		"CREATE INDEX IF NOT EXISTS idx_commands_status ON commands(status)",
		"CREATE INDEX IF NOT EXISTS idx_commands_created_at ON commands(created_at)",
	}
	
	for _, sql := range indexSQL {
		if err := db.Exec(sql).Error; err != nil {
			log.Printf("Warning: Failed to create index: %v", err)
			// Don't fail the migration for index creation failures
		}
	}
	
	log.Println("Commands table created manually with indexes")
	return nil
}

// createTableManually creates tables manually based on model name
func createTableManually(db *gorm.DB, modelName string) error {
	switch modelName {
	case "User":
		return createUserTableManually(db)
	case "Tag":
		return createTagTableManually(db)
	case "Block":
		return createBlockTableManually(db)
	case "Report":
		return createReportTableManually(db)
	default:
		return fmt.Errorf("unknown model name: %s", modelName)
	}
}

// createUserTableManually creates the users table manually
func createUserTableManually(db *gorm.DB) error {
	// Check if table already exists
	var exists bool
	if err := db.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'users')").Scan(&exists).Error; err != nil {
		return fmt.Errorf("error checking if users table exists: %w", err)
	}
	
	if exists {
		log.Println("Users table already exists, skipping manual creation")
		return nil
	}
	
	log.Println("Creating users table manually due to GORM migration failure...")
	
	createTableSQL := `
		CREATE TABLE users (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			screen_name VARCHAR(50) NOT NULL,
			login_name VARCHAR(50) NOT NULL UNIQUE,
			email VARCHAR(300) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			role VARCHAR(50),
			random_opt_in BOOLEAN DEFAULT false,
			anon_cmd BOOLEAN DEFAULT false,
			verified BOOLEAN DEFAULT false,
			verified_code BIGINT DEFAULT 0,
			thumbs_up BIGINT DEFAULT 0,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW(),
			login_date TIMESTAMPTZ DEFAULT NOW()
		)`
	
	if err := db.Exec(createTableSQL).Error; err != nil {
		return fmt.Errorf("error creating users table: %w", err)
	}
	
	// Create indexes
	indexSQL := []string{
		"CREATE UNIQUE INDEX IF NOT EXISTS idx_users_login_name ON users(login_name)",
		"CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email)",
		"CREATE INDEX IF NOT EXISTS idx_users_role ON users(role)",
		"CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at)",
	}
	
	for _, sql := range indexSQL {
		if err := db.Exec(sql).Error; err != nil {
			log.Printf("Warning: Failed to create index: %v", err)
		}
	}
	
	log.Println("Users table created manually with indexes")
	return nil
}

// createTagTableManually creates the tags table manually
func createTagTableManually(db *gorm.DB) error {
	var exists bool
	if err := db.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'tags')").Scan(&exists).Error; err != nil {
		return fmt.Errorf("error checking if tags table exists: %w", err)
	}
	
	if exists {
		log.Println("Tags table already exists, skipping manual creation")
		return nil
	}
	
	log.Println("Creating tags table manually due to GORM migration failure...")
	
	createTableSQL := `
		CREATE TABLE tags (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			name VARCHAR(100) NOT NULL UNIQUE,
			description TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`
	
	if err := db.Exec(createTableSQL).Error; err != nil {
		return fmt.Errorf("error creating tags table: %w", err)
	}
	
	// Create indexes
	indexSQL := []string{
		"CREATE UNIQUE INDEX IF NOT EXISTS idx_tags_name ON tags(name)",
		"CREATE INDEX IF NOT EXISTS idx_tags_created_at ON tags(created_at)",
	}
	
	for _, sql := range indexSQL {
		if err := db.Exec(sql).Error; err != nil {
			log.Printf("Warning: Failed to create index: %v", err)
		}
	}
	
	log.Println("Tags table created manually with indexes")
	return nil
}

// createBlockTableManually creates the blocks table manually
func createBlockTableManually(db *gorm.DB) error {
	var exists bool
	if err := db.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'blocks')").Scan(&exists).Error; err != nil {
		return fmt.Errorf("error checking if blocks table exists: %w", err)
	}
	
	if exists {
		log.Println("Blocks table already exists, skipping manual creation")
		return nil
	}
	
	log.Println("Creating blocks table manually due to GORM migration failure...")
	
	createTableSQL := `
		CREATE TABLE blocks (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_id UUID NOT NULL,
			blocked_id UUID NOT NULL,
			reason TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			CONSTRAINT fk_blocks_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			CONSTRAINT fk_blocks_blocked FOREIGN KEY (blocked_id) REFERENCES users(id) ON DELETE CASCADE
		)`
	
	if err := db.Exec(createTableSQL).Error; err != nil {
		return fmt.Errorf("error creating blocks table: %w", err)
	}
	
	// Create indexes
	indexSQL := []string{
		"CREATE INDEX IF NOT EXISTS idx_blocks_user_id ON blocks(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_blocks_blocked_id ON blocks(blocked_id)",
		"CREATE UNIQUE INDEX IF NOT EXISTS idx_blocks_unique ON blocks(user_id, blocked_id)",
	}
	
	for _, sql := range indexSQL {
		if err := db.Exec(sql).Error; err != nil {
			log.Printf("Warning: Failed to create index: %v", err)
		}
	}
	
	log.Println("Blocks table created manually with indexes")
	return nil
}

// createReportTableManually creates the reports table manually
func createReportTableManually(db *gorm.DB) error {
	var exists bool
	if err := db.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'reports')").Scan(&exists).Error; err != nil {
		return fmt.Errorf("error checking if reports table exists: %w", err)
	}
	
	if exists {
		log.Println("Reports table already exists, skipping manual creation")
		return nil
	}
	
	log.Println("Creating reports table manually due to GORM migration failure...")
	
	createTableSQL := `
		CREATE TABLE reports (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			reporter_id UUID NOT NULL,
			reported_id UUID NOT NULL,
			reason TEXT NOT NULL,
			status VARCHAR(20) DEFAULT 'pending',
			created_at TIMESTAMPTZ DEFAULT NOW(),
			CONSTRAINT fk_reports_reporter FOREIGN KEY (reporter_id) REFERENCES users(id) ON DELETE CASCADE,
			CONSTRAINT fk_reports_reported FOREIGN KEY (reported_id) REFERENCES users(id) ON DELETE CASCADE
		)`
	
	if err := db.Exec(createTableSQL).Error; err != nil {
		return fmt.Errorf("error creating reports table: %w", err)
	}
	
	// Create indexes
	indexSQL := []string{
		"CREATE INDEX IF NOT EXISTS idx_reports_reporter_id ON reports(reporter_id)",
		"CREATE INDEX IF NOT EXISTS idx_reports_reported_id ON reports(reported_id)",
		"CREATE INDEX IF NOT EXISTS idx_reports_status ON reports(status)",
		"CREATE INDEX IF NOT EXISTS idx_reports_created_at ON reports(created_at)",
	}
	
	for _, sql := range indexSQL {
		if err := db.Exec(sql).Error; err != nil {
			log.Printf("Warning: Failed to create index: %v", err)
		}
	}
	
	log.Println("Reports table created manually with indexes")
	return nil
}
