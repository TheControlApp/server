package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/thecontrolapp/controlme-go/internal/auth"
	"github.com/thecontrolapp/controlme-go/internal/config"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// User model matching the main application
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	ScreenName   string    `gorm:"size:50;not null" json:"screen_name"`
	LoginName    string    `gorm:"size:50;not null;unique" json:"login_name"`
	Password     string    `gorm:"size:255;not null" json:"-"`
	Role         string    `gorm:"size:50" json:"role"`
	RandOpt      bool      `gorm:"default:false" json:"rand_opt"`
	AnonCmd      bool      `gorm:"default:false" json:"anon_cmd"`
	Verified     bool      `gorm:"default:false" json:"verified"`
	VerifiedCode int       `gorm:"default:0" json:"verified_code"`
	ThumbsUp     int       `gorm:"default:0" json:"thumbs_up"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	LoginDate    time.Time `json:"login_date"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func main() {
	var (
		username     = flag.String("username", "", "Login username (required)")
		password     = flag.String("password", "", "User password (leave empty for no password)")
		passwordType = flag.String("password-type", "encrypt", "Password type: 'encrypt' (legacy AES), 'hash' (bcrypt), 'plain' (store as-is)")
		screenName   = flag.String("screen", "", "Screen/display name (defaults to username)")
		role         = flag.String("role", "user", "User role (user, admin, dom, sub)")
		verified     = flag.Bool("verified", true, "Whether user is verified")
		randOpt      = flag.Bool("rand-opt", false, "Random option setting")
		anonCmd      = flag.Bool("anon-cmd", false, "Anonymous command setting")
		thumbsUp     = flag.Int("thumbs", 0, "Thumbs up count")
		help         = flag.Bool("help", false, "Show help")
		list         = flag.Bool("list", false, "List existing users")
		createTable  = flag.Bool("create-table", false, "Create users table if it doesn't exist")
	)

	flag.Parse()

	if *help {
		showHelp()
		return
	}

	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	// Connect to database
	db, err := connectDB(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// Create table if requested
	if *createTable {
		if err := db.AutoMigrate(&User{}); err != nil {
			log.Fatalf("❌ Failed to create table: %v", err)
		}
		fmt.Println("✅ Users table created/updated successfully")
		return
	}

	// List users if requested
	if *list {
		listUsers(db)
		return
	}

	// Validate required fields
	if *username == "" {
		fmt.Println("❌ Username is required")
		showHelp()
		os.Exit(1)
	}

	// Set defaults
	if *screenName == "" {
		*screenName = *username
	}

	// Hash password if provided
	var processedPassword string
	if *password != "" {
		switch *passwordType {
		case "encrypt":
			// Use legacy AES encryption (for drop-in compatibility)
			encrypted, err := encryptLegacyPassword(*password, cfg)
			if err != nil {
				log.Fatalf("❌ Failed to encrypt password: %v", err)
			}
			processedPassword = encrypted
			fmt.Printf("🔐 Using legacy AES encryption\n")
		case "hash":
			// Use bcrypt hashing (for modern auth)
			hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
			if err != nil {
				log.Fatalf("❌ Failed to hash password: %v", err)
			}
			processedPassword = string(hash)
			fmt.Printf("🔒 Using bcrypt hashing\n")
		case "plain":
			// Store as-is (for importing existing data)
			processedPassword = *password
			fmt.Printf("⚠️  Storing password as plain text\n")
		default:
			log.Fatalf("❌ Invalid password type: %s (must be 'encrypt', 'hash', or 'plain')", *passwordType)
		}
	}

	// Create user
	user := User{
		ID:           uuid.New(),
		ScreenName:   *screenName,
		LoginName:    *username,
		Password:     processedPassword,
		Role:         *role,
		RandOpt:      *randOpt,
		AnonCmd:      *anonCmd,
		Verified:     *verified,
		VerifiedCode: 0,
		ThumbsUp:     *thumbsUp,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		LoginDate:    time.Now(),
	}

	if err := db.Create(&user).Error; err != nil {
		if strings.Contains(err.Error(), "relation \"users\" does not exist") {
			fmt.Println("❌ Users table doesn't exist. Run with --create-table first")
			os.Exit(1)
		}
		log.Fatalf("❌ Failed to create user: %v", err)
	}

	fmt.Printf("✅ Created user successfully!\n")
	fmt.Printf("   Username: %s\n", user.LoginName)
	fmt.Printf("   Screen Name: %s\n", user.ScreenName)
	fmt.Printf("   Role: %s\n", user.Role)
	fmt.Printf("   Verified: %t\n", user.Verified)
	fmt.Printf("   ID: %s\n", user.ID)
	if *password == "" {
		fmt.Printf("   Password: (none)\n")
	} else {
		fmt.Printf("   Password: (encrypted)\n")
	}
}

func connectDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC",
		cfg.Database.Host,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
}

func listUsers(db *gorm.DB) {
	var users []User
	if err := db.Find(&users).Error; err != nil {
		if strings.Contains(err.Error(), "relation \"users\" does not exist") {
			fmt.Println("❌ Users table doesn't exist. Run with --create-table first")
			return
		}
		log.Fatalf("❌ Failed to list users: %v", err)
	}

	if len(users) == 0 {
		fmt.Println("📝 No users found")
		return
	}

	fmt.Printf("📋 Found %d users:\n\n", len(users))
	fmt.Printf("%-20s %-20s %-10s %-10s %-36s\n", "USERNAME", "SCREEN NAME", "ROLE", "VERIFIED", "ID")
	fmt.Printf("%s\n", strings.Repeat("-", 100))

	for _, user := range users {
		verified := "✅"
		if !user.Verified {
			verified = "❌"
		}
		fmt.Printf("%-20s %-20s %-10s %-10s %-36s\n",
			user.LoginName,
			user.ScreenName,
			user.Role,
			verified,
			user.ID.String())
	}
}

func showHelp() {
	fmt.Println("🔧 ControlMe User Creation Tool")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/tools/user-manager/main.go [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --username string         Login username (required)")
	fmt.Println("  --password string         User password (optional, no password if empty)")
	fmt.Println("  --password-type string    Password storage type:")
	fmt.Println("                              'encrypt' - Legacy AES encryption (default, for .NET client compatibility)")
	fmt.Println("                              'hash'    - Modern bcrypt hashing")
	fmt.Println("                              'plain'   - Store as-is (for importing existing data)")
	fmt.Println("  --screen string           Display name (defaults to username)")
	fmt.Println("  --role string             User role: user, admin, dom, sub (default: user)")
	fmt.Println("  --verified                User is verified (default: true)")
	fmt.Println("  --rand-opt                Enable random option (default: false)")
	fmt.Println("  --anon-cmd                Enable anonymous commands (default: false)")
	fmt.Println("  --thumbs int              Thumbs up count (default: 0)")
	fmt.Println("  --list                    List all existing users")
	fmt.Println("  --create-table            Create users table if it doesn't exist")
	fmt.Println("  --help                    Show this help")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Create table first")
	fmt.Println("  go run cmd/tools/user-manager/main.go --create-table")
	fmt.Println()
	fmt.Println("  # Create user with legacy AES encryption (for .NET client compatibility)")
	fmt.Println("  go run cmd/tools/user-manager/main.go --username alice --password secret123")
	fmt.Println("  go run cmd/tools/user-manager/main.go --username alice --password secret123 --password-type encrypt")
	fmt.Println()
	fmt.Println("  # Create user with modern bcrypt hashing")
	fmt.Println("  go run cmd/tools/user-manager/main.go --username bob --password secret123 --password-type hash")
	fmt.Println()
	fmt.Println("  # Import user with existing encrypted password")
	fmt.Println("  go run cmd/tools/user-manager/main.go --username charlie --password 'AAAAACm1LLz77SmR8t3v12QoWhw=' --password-type plain")
	fmt.Println()
	fmt.Println("  # Create admin user")
	fmt.Println("  go run cmd/tools/user-manager/main.go --username admin --password admin123 --role admin")
	fmt.Println()
	fmt.Println("  # Create dom/sub users")
	fmt.Println("  go run cmd/tools/user-manager/main.go --username master --password dom123 --role dom")
	fmt.Println("  go run cmd/tools/user-manager/main.go --username slave --password sub123 --role sub")
	fmt.Println()
	fmt.Println("  # List all users")
	fmt.Println("  go run cmd/tools/user-manager/main.go --list")
}

// encryptLegacyPassword encrypts a password using the legacy AES system
func encryptLegacyPassword(password string, cfg *config.Config) (string, error) {
	legacyCrypto := auth.NewLegacyCrypto(cfg.Legacy.CryptoKey)
	return legacyCrypto.Encrypt(password)
}
