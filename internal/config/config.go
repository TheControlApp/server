package config

import (
	"os"
	
	"github.com/spf13/viper"
)

type Config struct {
	Environment string   `mapstructure:"environment"`
	Server      Server   `mapstructure:"server"`
	Database    Database `mapstructure:"database"`
	Redis       Redis    `mapstructure:"redis"`
	Auth        Auth     `mapstructure:"auth"`
	Legacy      Legacy   `mapstructure:"legacy"`
}

type Server struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type Database struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"sslmode"`
}

type Redis struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type Auth struct {
	JWTSecret     string `mapstructure:"jwt_secret"`
	JWTExpiration int    `mapstructure:"jwt_expiration"`
}

type Legacy struct {
	CryptoKey             string `mapstructure:"crypto_key"`
	UpgradeNotifications  bool   `mapstructure:"upgrade_notifications"`
	NotificationFrequency int    `mapstructure:"notification_frequency"`
	SunsetDate            string `mapstructure:"sunset_date"`
}

func Load() (*Config, error) {
	// Check for custom config file from environment
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = "config"
	} else {
		// Remove .yaml extension if present
		if len(configFile) > 5 && configFile[len(configFile)-5:] == ".yaml" {
			configFile = configFile[:len(configFile)-5]
		}
	}
	
	viper.SetConfigName(configFile)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// Set defaults
	viper.SetDefault("environment", "development")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.name", "controlme")
	viper.SetDefault("database.username", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("auth.jwt_expiration", 86400) // 24 hours
	viper.SetDefault("legacy.upgrade_notifications", true)
	viper.SetDefault("legacy.notification_frequency", 24) // hours

	// Read environment variables
	viper.AutomaticEnv()
	
	// Bind specific environment variables to config keys
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.name", "DB_NAME")
	viper.BindEnv("database.username", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("redis.host", "REDIS_HOST")
	viper.BindEnv("redis.port", "REDIS_PORT")
	viper.BindEnv("redis.password", "REDIS_PASSWORD")
	viper.BindEnv("auth.jwt_secret", "JWT_SECRET")
	viper.BindEnv("legacy.crypto_key", "LEGACY_CRYPTO_KEY")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
