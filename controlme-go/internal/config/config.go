package config

import (
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
	viper.SetConfigName("config")
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
