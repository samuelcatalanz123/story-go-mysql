// Package config loads runtime configuration from environment variables,
// falling back to sensible development defaults.
package config

import (
	"fmt"
	"os"
)

// Config holds every setting the application needs to start.
type Config struct {
	ServerAddr string
	WebDir     string
	JWTSecret  string
	DB         DBConfig
}

// DBConfig holds the MySQL connection settings.
type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

// DSN builds the MySQL data source name from the configured fields.
func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		c.User, c.Password, c.Host, c.Port, c.Name,
	)
}

// Load reads configuration from the environment, applying defaults that
// match a local MySQL instance so the project runs out of the box.
func Load() Config {
	return Config{
		ServerAddr: serverAddr(),
		WebDir:     env("WEB_DIR", "web/dist"),
		JWTSecret:  env("JWT_SECRET", "dev-insecure-secret-change-me"),
		DB: DBConfig{
			User:     env("DB_USER", "root"),
			Password: env("DB_PASSWORD", ""),
			Host:     env("DB_HOST", "127.0.0.1"),
			Port:     env("DB_PORT", "3306"),
			Name:     env("DB_NAME", "story_go_db"),
		},
	}
}

// serverAddr prefers the PORT variable that platforms like Railway inject,
// falling back to SERVER_ADDR (default :8080) for local development.
func serverAddr() string {
	if port, ok := os.LookupEnv("PORT"); ok {
		return ":" + port
	}
	return env("SERVER_ADDR", ":8080")
}

// env returns the value of the environment variable or a fallback default.
func env(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
