package config

import (
	"fmt"
	"os"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	TimeZone string
}

func DBFromEnv() DBConfig {
	return DBConfig{
		Host:     get("POSTGRES_HOST", "localhost"),
		Port:     get("POSTGRES_PORT", "5432"),
		User:     get("POSTGRES_USER", "postgres"),
		Password: get("POSTGRES_PASSWORD", ""),
		DBName:   get("POSTGRES_DB", "postgres"),
		SSLMode:  get("POSTGRES_SSLMODE", "disable"),
		TimeZone: get("DB_TIMEZONE", "UTC"),
	}
}

func (dc DBConfig) DSN() string {
	// gorm postgres DSN
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		dc.Host, dc.User, dc.Password, dc.DBName, dc.Port, dc.SSLMode, dc.TimeZone,
	)
}

func get(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
