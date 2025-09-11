package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBUser string
	DBPass string
	DBHost string
	DBName string
	Port   string
}

func Load() *Config {
	return &Config{
		DBUser: getenv("DB_USER", "root"),
		DBPass: getenv("DB_PASS", "Wsxokn@123"),
		DBHost: getenv("DB_HOST", "localhost:3306"),
		DBName: getenv("DB_NAME", "usersdb"),
		Port:   getenv("PORT", "8080"),
	}
}

func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		c.DBUser, c.DBPass, c.DBHost, c.DBName)
}

func getenv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}
